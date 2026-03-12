package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
)

func TestPaymentRepoIntegration_SaveGetLatestPendingAndSettleOnce(t *testing.T) {
	testDB := newIntegrationTestDB(t)

	playerRepo := NewPlayerPostgresRepo(testDB.DB)
	paymentRepo := NewPaymentRepo(testDB.DB)

	player, err := playerRepo.GetOrCreate(720001)
	if err != nil {
		t.Fatalf("create player: %v", err)
	}

	ctx := context.Background()

	olderPending := &domain.Payment{
		ID:           "pay-old-pending",
		PlayerID:     player.UserID,
		Amount:       10,
		Status:       domain.PaymentPending,
		QRCode:       "qr-old",
		QRCodeBase64: "base64-old",
		PartnerId:    990001,
		CreatedAt:    time.Now().UTC().Add(-2 * time.Minute),
		LastUpdate:   time.Now().UTC().Add(-2 * time.Minute),
	}
	if err := paymentRepo.Save(ctx, olderPending); err != nil {
		t.Fatalf("save older pending payment: %v", err)
	}

	approved := &domain.Payment{
		ID:           "pay-approved",
		PlayerID:     player.UserID,
		Amount:       4,
		Status:       domain.PaymentApproved,
		QRCode:       "qr-approved",
		QRCodeBase64: "base64-approved",
		PartnerId:    990002,
		CreatedAt:    time.Now().UTC().Add(-1 * time.Minute),
		LastUpdate:   time.Now().UTC().Add(-1 * time.Minute),
	}
	if err := paymentRepo.Save(ctx, approved); err != nil {
		t.Fatalf("save approved payment: %v", err)
	}

	latestPending := &domain.Payment{
		ID:           "pay-latest-pending",
		PlayerID:     player.UserID,
		Amount:       12,
		Status:       domain.PaymentPending,
		QRCode:       "qr-latest",
		QRCodeBase64: "base64-latest",
		PartnerId:    990003,
		CreatedAt:    time.Now().UTC(),
		LastUpdate:   time.Now().UTC(),
	}
	if err := paymentRepo.Save(ctx, latestPending); err != nil {
		t.Fatalf("save latest pending payment: %v", err)
	}

	gotPending, err := paymentRepo.GetLatestPendingByPlayerID(ctx, player.UserID)
	if err != nil {
		t.Fatalf("get latest pending payment: %v", err)
	}

	if gotPending.ID != latestPending.ID {
		t.Fatalf("unexpected latest pending id: got %s want %s", gotPending.ID, latestPending.ID)
	}
	if gotPending.QRCode != "qr-latest" {
		t.Fatalf("unexpected qr code for latest pending payment: %q", gotPending.QRCode)
	}

	settledNow, err := paymentRepo.SettlePendingByPartnerID(ctx, latestPending.PartnerId, 3)
	if err != nil {
		t.Fatalf("settle pending payment: %v", err)
	}
	if !settledNow {
		t.Fatal("expected first settlement to settle now")
	}

	var persistedPayment domain.Payment
	if err := testDB.DB.Where("id = ?", latestPending.ID).First(&persistedPayment).Error; err != nil {
		t.Fatalf("load persisted payment: %v", err)
	}
	if persistedPayment.Status != domain.PaymentApproved {
		t.Fatalf("unexpected persisted payment status: got %s want %s", persistedPayment.Status, domain.PaymentApproved)
	}

	reloadedPlayer, err := playerRepo.GetOrCreate(player.UserID)
	if err != nil {
		t.Fatalf("reload player after settlement: %v", err)
	}
	if reloadedPlayer.KiCrystals != 36 {
		t.Fatalf("unexpected ki crystals after settlement: got %d want 36", reloadedPlayer.KiCrystals)
	}

	settledNow, err = paymentRepo.SettlePendingByPartnerID(ctx, latestPending.PartnerId, 3)
	if err != nil {
		t.Fatalf("settle same payment second time: %v", err)
	}
	if settledNow {
		t.Fatal("expected second settlement to be idempotent")
	}

	reloadedPlayer, err = playerRepo.GetOrCreate(player.UserID)
	if err != nil {
		t.Fatalf("reload player after second settlement: %v", err)
	}
	if reloadedPlayer.KiCrystals != 36 {
		t.Fatalf("unexpected ki crystals after second settlement: got %d want 36", reloadedPlayer.KiCrystals)
	}
}

func TestPaymentRepoIntegration_Save_IsIdempotentByPaymentID(t *testing.T) {
	testDB := newIntegrationTestDB(t)

	paymentRepo := NewPaymentRepo(testDB.DB)
	ctx := context.Background()

	payment := &domain.Payment{
		ID:           "pay-duplicate-safe",
		PlayerID:     880001,
		Amount:       7,
		Status:       domain.PaymentPending,
		QRCode:       "qr-1",
		QRCodeBase64: "base64-1",
		PartnerId:    123456,
		CreatedAt:    time.Now().UTC(),
		LastUpdate:   time.Now().UTC(),
	}
	if err := paymentRepo.Save(ctx, payment); err != nil {
		t.Fatalf("save first payment: %v", err)
	}

	duplicate := &domain.Payment{
		ID:           payment.ID,
		PlayerID:     payment.PlayerID,
		Amount:       99,
		Status:       domain.PaymentApproved,
		QRCode:       "qr-2",
		QRCodeBase64: "base64-2",
		PartnerId:    999999,
		CreatedAt:    time.Now().UTC(),
		LastUpdate:   time.Now().UTC(),
	}
	if err := paymentRepo.Save(ctx, duplicate); err != nil {
		t.Fatalf("save duplicate payment: %v", err)
	}

	var count int64
	if err := testDB.DB.Model(&domain.Payment{}).Where("id = ?", payment.ID).Count(&count).Error; err != nil {
		t.Fatalf("count payments by id: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected one payment row after duplicate save, got %d", count)
	}

	var persisted domain.Payment
	if err := testDB.DB.Where("id = ?", payment.ID).First(&persisted).Error; err != nil {
		t.Fatalf("load persisted payment: %v", err)
	}
	if persisted.Amount != payment.Amount || persisted.Status != payment.Status || persisted.PartnerId != payment.PartnerId {
		t.Fatalf("expected original payment to remain unchanged after duplicate save, got %+v", persisted)
	}
}

func TestPaymentRepoIntegration_SettleRollbackWhenPlayerMissing(t *testing.T) {
	testDB := newIntegrationTestDB(t)

	paymentRepo := NewPaymentRepo(testDB.DB)
	ctx := context.Background()

	payment := &domain.Payment{
		ID:           "pay-missing-player",
		PlayerID:     999999,
		Amount:       5,
		Status:       domain.PaymentPending,
		QRCode:       "qr-missing-player",
		QRCodeBase64: "base64-missing-player",
		PartnerId:    555000,
		CreatedAt:    time.Now().UTC(),
		LastUpdate:   time.Now().UTC(),
	}
	if err := paymentRepo.Save(ctx, payment); err != nil {
		t.Fatalf("save payment fixture: %v", err)
	}

	settledNow, err := paymentRepo.SettlePendingByPartnerID(ctx, payment.PartnerId, 3)
	if err == nil {
		t.Fatal("expected settlement to fail when player is missing")
	}
	if settledNow {
		t.Fatal("expected settlement flag to remain false on rollback")
	}

	var persisted domain.Payment
	if err := testDB.DB.Where("id = ?", payment.ID).First(&persisted).Error; err != nil {
		t.Fatalf("load persisted payment after rollback: %v", err)
	}
	if persisted.Status != domain.PaymentPending {
		t.Fatalf("expected payment status to rollback to pending, got %q", persisted.Status)
	}
}
