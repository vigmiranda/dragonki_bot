package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterMercadoPagoWebhook_RegisteredRoute(t *testing.T) {
	mux := http.NewServeMux()
	called := false
	registerMercadoPagoWebhook(mux, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusAccepted)
	})

	req := httptest.NewRequest(http.MethodPost, "/webhook/mercadopago", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", rec.Code)
	}
	if !called {
		t.Fatalf("expected handler call")
	}
}

func TestRegisterMercadoPagoWebhook_UnregisteredRoute(t *testing.T) {
	mux := http.NewServeMux()
	registerMercadoPagoWebhook(mux, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/webhook/other", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}
