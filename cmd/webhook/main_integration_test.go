package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vigmiranda/dragonki/internal/adapters/mercadopago"
	"github.com/vigmiranda/dragonki/internal/adapters/postgres"
	"github.com/vigmiranda/dragonki/internal/domain"
	integrationtest "github.com/vigmiranda/dragonki/internal/testutil/integration"
)

func TestBuildWebhookAppIntegration_ConfirmsPaymentAndRemainsIdempotent(t *testing.T) {
	testDB := integrationtest.NewPostgresDB(t, integrationtest.PostgresOptions{
		DatabaseNamePrefix: "dragonki_webhook_it",
	})

	playerRepo := postgres.NewPlayerPostgresRepo(testDB.DB)
	paymentRepo := postgres.NewPaymentRepo(testDB.DB)

	player, err := playerRepo.GetOrCreate(930001)
	if err != nil {
		t.Fatalf("create player: %v", err)
	}
	initialCrystals := player.KiCrystals

	paymentRecord := &domain.Payment{
		ID:           "payment_it_001",
		PlayerID:     player.UserID,
		Amount:       3,
		Status:       domain.PaymentPending,
		QRCode:       "qr-code",
		QRCodeBase64: "qr-base64",
		PartnerId:    987654,
	}
	if err := paymentRepo.Save(context.Background(), paymentRecord); err != nil {
		t.Fatalf("save payment fixture: %v", err)
	}

	t.Setenv("MP_WEBHOOK_TOKEN", "secret-token")

	app, err := buildWebhookApp(webhookAppConfig{
		DB:       testDB.DB,
		MPClient: &mercadopago.Client{Ctx: context.Background()},
		FeatureFlagResolver: func(name string) (bool, bool) {
			if name == "PAYMENTS" {
				return true, true
			}
			return false, false
		},
	})
	if err != nil {
		t.Fatalf("build webhook app: %v", err)
	}
	defer app.Cleanup()

	payload := `{"type":"payment","action":"payment.updated","data":{"id":"987654"}}`
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/webhook/mercadopago", strings.NewReader(payload))
		req.Header.Set("X-Webhook-Token", "secret-token")
		rec := httptest.NewRecorder()

		app.Handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 on webhook call %d, got %d", i+1, rec.Code)
		}
	}

	reloadedPlayer, err := playerRepo.GetOrCreate(player.UserID)
	if err != nil {
		t.Fatalf("reload player: %v", err)
	}
	if reloadedPlayer.KiCrystals != initialCrystals+3 {
		t.Fatalf("unexpected ki crystals after webhook: got %d want %d", reloadedPlayer.KiCrystals, initialCrystals+3)
	}

	var persistedPayment domain.Payment
	if err := testDB.DB.Where("id = ?", paymentRecord.ID).First(&persistedPayment).Error; err != nil {
		t.Fatalf("load persisted payment: %v", err)
	}
	if persistedPayment.Status != domain.PaymentApproved {
		t.Fatalf("unexpected payment status: got %q want %q", persistedPayment.Status, domain.PaymentApproved)
	}

	var audits []domain.PaymentWebhookAudit
	if err := testDB.DB.Order("id asc").Find(&audits).Error; err != nil {
		t.Fatalf("load webhook audits: %v", err)
	}
	if len(audits) != 2 {
		t.Fatalf("expected 2 webhook audits, got %d", len(audits))
	}
	for i, audit := range audits {
		if audit.ProcessingStatus != "CONFIRMED" {
			t.Fatalf("unexpected processing status on audit %d: %q", i+1, audit.ProcessingStatus)
		}
		if audit.HTTPStatus != http.StatusOK {
			t.Fatalf("unexpected http status on audit %d: %d", i+1, audit.HTTPStatus)
		}
		if audit.PartnerPaymentID != 987654 {
			t.Fatalf("unexpected partner payment id on audit %d: %d", i+1, audit.PartnerPaymentID)
		}
		if !audit.Authorized || !audit.FeatureEnabled {
			t.Fatalf("expected authorized and feature-enabled audit on row %d", i+1)
		}
	}
}
