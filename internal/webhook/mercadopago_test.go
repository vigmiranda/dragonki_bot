package webhook

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vigmiranda/dragonki/internal/bootstrap"
)

type fakePaymentConfirmer struct {
	confirmFn   func(ctx context.Context, paymentID int64) error
	callCount   int
	lastPayment int64
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
}

func TestMercadoPagoWebhook_Handle_Success(t *testing.T) {
	bootstrap.SetFeatureFlagResolver(func(name string) (bool, bool) { return true, true })
	t.Cleanup(func() { bootstrap.SetFeatureFlagResolver(nil) })

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
