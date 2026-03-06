package webhook

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakePaymentConfirmer struct {
	confirmFn  func(paymentID int64) error
	callCount  int
	lastPaymen int64
}

func (f *fakePaymentConfirmer) ConfirmPayment(paymentID int64) error {
	f.callCount++
	f.lastPaymen = paymentID
	if f.confirmFn == nil {
		return nil
	}
	return f.confirmFn(paymentID)
}

func TestMercadoPagoWebhook_Handle_InvalidJSON(t *testing.T) {
	w := NewMercadoPagoWebhook(&fakePaymentConfirmer{})
	req := httptest.NewRequest(http.MethodPost, "/webhook/mercadopago", strings.NewReader("{"))
	rec := httptest.NewRecorder()

	w.Handle(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestMercadoPagoWebhook_Handle_IgnoresNonPaymentEvents(t *testing.T) {
	fake := &fakePaymentConfirmer{}
	w := NewMercadoPagoWebhook(fake)
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
	fake := &fakePaymentConfirmer{}
	w := NewMercadoPagoWebhook(fake)
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
	fake := &fakePaymentConfirmer{}
	w := NewMercadoPagoWebhook(fake)
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
	fake := &fakePaymentConfirmer{
		confirmFn: func(paymentID int64) error { return errors.New("db fail") },
	}
	w := NewMercadoPagoWebhook(fake)
	req := httptest.NewRequest(http.MethodPost, "/webhook/mercadopago", strings.NewReader(`{"type":"payment","action":"payment.updated","data":{"id":"123"}}`))
	rec := httptest.NewRecorder()

	w.Handle(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
	if fake.callCount != 1 || fake.lastPaymen != 123 {
		t.Fatalf("unexpected confirm payload: calls=%d id=%d", fake.callCount, fake.lastPaymen)
	}
}

func TestMercadoPagoWebhook_Handle_Success(t *testing.T) {
	fake := &fakePaymentConfirmer{}
	w := NewMercadoPagoWebhook(fake)
	req := httptest.NewRequest(http.MethodPost, "/webhook/mercadopago", strings.NewReader(`{"type":"payment","action":"payment.updated","data":{"id":"456"}}`))
	rec := httptest.NewRecorder()

	w.Handle(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if fake.callCount != 1 || fake.lastPaymen != 456 {
		t.Fatalf("unexpected confirm payload: calls=%d id=%d", fake.callCount, fake.lastPaymen)
	}
}
