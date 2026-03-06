package main

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type fakeTelegramUpdateHandler struct {
	callCount int
	last      tgbotapi.Update
}

func (f *fakeTelegramUpdateHandler) HandleUpdate(update tgbotapi.Update) {
	f.callCount++
	f.last = update
}

func TestRegisterTelegramWebhook_DecodeError(t *testing.T) {
	mux := http.NewServeMux()
	h := &fakeTelegramUpdateHandler{}
	registerTelegramWebhook(mux, h, func(r io.Reader, update *tgbotapi.Update) error {
		return errors.New("decode fail")
	})

	req := httptest.NewRequest(http.MethodPost, "/webhook/telegram", strings.NewReader("{"))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	if h.callCount != 0 {
		t.Fatalf("expected no handler call, got %d", h.callCount)
	}
}

func TestRegisterTelegramWebhook_Success(t *testing.T) {
	mux := http.NewServeMux()
	h := &fakeTelegramUpdateHandler{}
	registerTelegramWebhook(mux, h, func(r io.Reader, update *tgbotapi.Update) error {
		update.UpdateID = 123
		return nil
	})

	req := httptest.NewRequest(http.MethodPost, "/webhook/telegram", strings.NewReader("{}"))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if h.callCount != 1 {
		t.Fatalf("expected one handler call, got %d", h.callCount)
	}
	if h.last.UpdateID != 123 {
		t.Fatalf("expected update id 123, got %d", h.last.UpdateID)
	}
}

func TestRegisterTelegramWebhook_UnregisteredRoute(t *testing.T) {
	mux := http.NewServeMux()
	h := &fakeTelegramUpdateHandler{}
	registerTelegramWebhook(mux, h, decodeTelegramUpdate)

	req := httptest.NewRequest(http.MethodPost, "/webhook/other", strings.NewReader("{}"))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
	if h.callCount != 0 {
		t.Fatalf("expected no handler call, got %d", h.callCount)
	}
}

func TestRegisterTelegramWebhook_MethodNotAllowed(t *testing.T) {
	mux := http.NewServeMux()
	h := &fakeTelegramUpdateHandler{}
	registerTelegramWebhook(mux, h, decodeTelegramUpdate)

	req := httptest.NewRequest(http.MethodGet, "/webhook/telegram", strings.NewReader("{}"))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
	if h.callCount != 0 {
		t.Fatalf("expected no handler call, got %d", h.callCount)
	}
}

func TestShouldRunMigrations(t *testing.T) {
	t.Setenv("RUN_MIGRATIONS", "true")
	if !shouldRunMigrations() {
		t.Fatalf("expected true when env is true")
	}

	t.Setenv("RUN_MIGRATIONS", "false")
	if shouldRunMigrations() {
		t.Fatalf("expected false when env is false")
	}
}
