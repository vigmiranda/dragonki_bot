package webhook

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vigmiranda/dragonki/internal/bootstrap"
	"github.com/vigmiranda/dragonki/internal/domain"
)

type fakePaymentConfirmer struct {
	confirmFn   func(ctx context.Context, paymentID int64) error
	callCount   int
	lastPayment int64
}

type fakeWebhookAuditStore struct {
	saveFn    func(ctx context.Context, audit *domain.PaymentWebhookAudit) error
	callCount int
	last      *domain.PaymentWebhookAudit
}

func (f *fakeWebhookAuditStore) Save(ctx context.Context, audit *domain.PaymentWebhookAudit) error {
	f.callCount++
	f.last = audit
	if f.saveFn == nil {
		return nil
	}
	return f.saveFn(ctx, audit)
}

func (f *fakePaymentConfirmer) ConfirmPayment(ctx context.Context, paymentID int64) error {
	f.callCount++
	f.lastPayment = paymentID
	if f.confirmFn == nil {
		return nil
	}
	return f.confirmFn(ctx, paymentID)
}

func TestMercadoPagoWebhook_Handle_InvalidJSON(t *testing.T) {
	bootstrap.SetFeatureFlagResolver(func(name string) (bool, bool) { return true, true })
	t.Cleanup(func() { bootstrap.SetFeatureFlagResolver(nil) })

	w := NewMercadoPagoWebhook(&fakePaymentConfirmer{}, "")
	req := httptest.NewRequest(http.MethodPost, "/webhook/mercadopago", strings.NewReader("{"))
	rec := httptest.NewRecorder()

	w.Handle(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestMercadoPagoWebhook_Handle_IgnoresNonPaymentEvents(t *testing.T) {
	bootstrap.SetFeatureFlagResolver(func(name string) (bool, bool) { return true, true })
	t.Cleanup(func() { bootstrap.SetFeatureFlagResolver(nil) })

	fake := &fakePaymentConfirmer{}
	w := NewMercadoPagoWebhook(fake, "")
	req := httptest.NewRequest(http.MethodPost, "/webhook/mercadopago", strings.NewReader(`{"type":"merchant_order","data":{"id":"1"}}`))
	rec := httptest.NewRecorder()

	w.Handle(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if fake.callCount != 0 {
		t.Fatalf("expected no confirm call, got %d", fake.callCount)
	}
}

func TestMercadoPagoWebhook_Handle_PaymentCreatedDoesNotConfirm(t *testing.T) {
	bootstrap.SetFeatureFlagResolver(func(name string) (bool, bool) { return true, true })
	t.Cleanup(func() { bootstrap.SetFeatureFlagResolver(nil) })

	fake := &fakePaymentConfirmer{}
	w := NewMercadoPagoWebhook(fake, "")
	req := httptest.NewRequest(http.MethodPost, "/webhook/mercadopago", strings.NewReader(`{"type":"payment","action":"payment.created","data":{"id":"123"}}`))
	rec := httptest.NewRecorder()

	w.Handle(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if fake.callCount != 0 {
		t.Fatalf("expected no confirm call, got %d", fake.callCount)
	}
}

func TestMercadoPagoWebhook_Handle_InvalidPaymentID(t *testing.T) {
	bootstrap.SetFeatureFlagResolver(func(name string) (bool, bool) { return true, true })
	t.Cleanup(func() { bootstrap.SetFeatureFlagResolver(nil) })

	fake := &fakePaymentConfirmer{}
	w := NewMercadoPagoWebhook(fake, "")
	req := httptest.NewRequest(http.MethodPost, "/webhook/mercadopago", strings.NewReader(`{"type":"payment","action":"payment.updated","data":{"id":"abc"}}`))
	rec := httptest.NewRecorder()

	w.Handle(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if fake.callCount != 0 {
		t.Fatalf("expected no confirm call, got %d", fake.callCount)
	}
}

func TestMercadoPagoWebhook_Handle_ConfirmError(t *testing.T) {
	bootstrap.SetFeatureFlagResolver(func(name string) (bool, bool) { return true, true })
	t.Cleanup(func() { bootstrap.SetFeatureFlagResolver(nil) })
	bootstrap.ResetAppMetricsForTesting()

	fake := &fakePaymentConfirmer{
		confirmFn: func(ctx context.Context, paymentID int64) error { return errors.New("db fail") },
	}
	w := NewMercadoPagoWebhook(fake, "")
	req := httptest.NewRequest(http.MethodPost, "/webhook/mercadopago", strings.NewReader(`{"type":"payment","action":"payment.updated","data":{"id":"123"}}`))
	rec := httptest.NewRecorder()

	w.Handle(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
	if fake.callCount != 1 || fake.lastPayment != 123 {
		t.Fatalf("unexpected confirm payload: calls=%d id=%d", fake.callCount, fake.lastPayment)
	}
	if got := bootstrap.SnapshotAppMetrics()["payments.webhook.confirm_error_total"]; got != 1 {
		t.Fatalf("expected webhook confirm error counter 1, got %d", got)
	}
}

func TestMercadoPagoWebhook_Handle_Success(t *testing.T) {
	bootstrap.SetFeatureFlagResolver(func(name string) (bool, bool) { return true, true })
	t.Cleanup(func() { bootstrap.SetFeatureFlagResolver(nil) })
	bootstrap.ResetAppMetricsForTesting()

	fake := &fakePaymentConfirmer{}
	w := NewMercadoPagoWebhook(fake, "")
	req := httptest.NewRequest(http.MethodPost, "/webhook/mercadopago", strings.NewReader(`{"type":"payment","action":"payment.updated","data":{"id":"456"}}`))
	rec := httptest.NewRecorder()

	w.Handle(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if fake.callCount != 1 || fake.lastPayment != 456 {
		t.Fatalf("unexpected confirm payload: calls=%d id=%d", fake.callCount, fake.lastPayment)
	}
	if got := bootstrap.SnapshotAppMetrics()["payments.webhook.confirmed_total"]; got != 1 {
		t.Fatalf("expected webhook confirmed counter 1, got %d", got)
	}
}

func TestMercadoPagoWebhook_Handle_AuditSaved(t *testing.T) {
	bootstrap.SetFeatureFlagResolver(func(name string) (bool, bool) { return true, true })
	t.Cleanup(func() { bootstrap.SetFeatureFlagResolver(nil) })

	fake := &fakePaymentConfirmer{}
	auditStore := &fakeWebhookAuditStore{}
	w := NewMercadoPagoWebhook(fake, "", auditStore)
	req := httptest.NewRequest(http.MethodPost, "/webhook/mercadopago?source=test", strings.NewReader(`{"type":"payment","action":"payment.updated","data":{"id":"456"}}`))
	req.Header.Set("X-Test", "ok")
	rec := httptest.NewRecorder()

	w.Handle(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if auditStore.callCount != 1 {
		t.Fatalf("expected one audit record, got %d", auditStore.callCount)
	}
	if auditStore.last == nil {
		t.Fatalf("expected audit payload")
	}
	if auditStore.last.ProcessingStatus != "CONFIRMED" {
		t.Fatalf("expected CONFIRMED status, got %q", auditStore.last.ProcessingStatus)
	}
	if auditStore.last.PartnerPaymentID != 456 {
		t.Fatalf("expected partner_payment_id=456, got %d", auditStore.last.PartnerPaymentID)
	}
	if !strings.Contains(auditStore.last.RawBody, "\"payment.updated\"") {
		t.Fatalf("expected raw body to contain webhook payload")
	}
}

func TestMercadoPagoWebhook_Handle_UnauthorizedWhenTokenMismatch(t *testing.T) {
	bootstrap.SetFeatureFlagResolver(func(name string) (bool, bool) { return true, true })
	t.Cleanup(func() { bootstrap.SetFeatureFlagResolver(nil) })

	fake := &fakePaymentConfirmer{}
	w := NewMercadoPagoWebhook(fake, "secret")
	req := httptest.NewRequest(http.MethodPost, "/webhook/mercadopago", strings.NewReader(`{"type":"payment","action":"payment.updated","data":{"id":"456"}}`))
	rec := httptest.NewRecorder()

	w.Handle(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
	if fake.callCount != 0 {
		t.Fatalf("expected no confirm call, got %d", fake.callCount)
	}
}

func TestMercadoPagoWebhook_Handle_DisabledByFeatureFlag(t *testing.T) {
	bootstrap.SetFeatureFlagResolver(func(name string) (bool, bool) { return false, true })
	t.Cleanup(func() { bootstrap.SetFeatureFlagResolver(nil) })

	fake := &fakePaymentConfirmer{}
	w := NewMercadoPagoWebhook(fake, "")
	req := httptest.NewRequest(http.MethodPost, "/webhook/mercadopago", strings.NewReader(`{"type":"payment","action":"payment.updated","data":{"id":"456"}}`))
	rec := httptest.NewRecorder()

	w.Handle(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 when payments feature is disabled, got %d", rec.Code)
	}
	if fake.callCount != 0 {
		t.Fatalf("expected no confirm call, got %d", fake.callCount)
	}
}
