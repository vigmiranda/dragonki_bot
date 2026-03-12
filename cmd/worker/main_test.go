package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/gorm"
)

func TestBuildWorkerApp_WithObservabilityEnabledExposesDebugVars(t *testing.T) {
	app, err := buildWorkerApp(workerAppConfig{
		DB:                   &gorm.DB{},
		RedisURL:             "redis://localhost:6379/0",
		Concurrency:          1,
		ObservabilityEnabled: true,
		PprofEnabled:         false,
		ObservabilityAddr:    ":0",
	})
	if err != nil {
		t.Fatalf("build worker app: %v", err)
	}
	defer app.Cleanup()

	if app.ObservabilityHandler == nil {
		t.Fatal("expected observability handler")
	}
	if app.ObservabilityServer == nil {
		t.Fatal("expected observability server")
	}

	req := httptest.NewRequest(http.MethodGet, "/debug/vars", nil)
	rec := httptest.NewRecorder()
	app.ObservabilityHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode vars payload: %v", err)
	}
	if payload["service"] != "worker" {
		t.Fatalf("expected service worker, got %v", payload["service"])
	}
}
