package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/adapters/mercadopago"
	"github.com/vigmiranda/dragonki/internal/adapters/postgres"
	"github.com/vigmiranda/dragonki/internal/domain"
	integrationtest "github.com/vigmiranda/dragonki/internal/testutil/integration"
)

type telegramRequestRecord struct {
	Path    string
	ChatID  string
	Text    string
	Caption string
}

func TestBuildBotAppIntegration_TelegramStartReferralPersistsOnboardingAndReferral(t *testing.T) {
	testDB := integrationtest.NewPostgresDB(t, integrationtest.PostgresOptions{
		DatabaseNamePrefix: "dragonki_bot_it",
	})

	playerRepo := postgres.NewPlayerPostgresRepo(testDB.DB)
	referralRepo := postgres.NewReferralPostgresRepo(testDB.DB)

	referrerUserID := int64(910001)
	referrer, err := playerRepo.GetOrCreate(referrerUserID)
	if err != nil {
		t.Fatalf("create referrer player: %v", err)
	}

	var (
		mu       sync.Mutex
		requests []telegramRequestRecord
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseMultipartForm(10 << 20)
		_ = r.ParseForm()

		mu.Lock()
		requests = append(requests, telegramRequestRecord{
			Path:    r.URL.Path,
			ChatID:  r.FormValue("chat_id"),
			Text:    r.FormValue("text"),
			Caption: r.FormValue("caption"),
		})
		mu.Unlock()

		switch {
		case strings.HasSuffix(r.URL.Path, "/getMe"):
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":true,"result":{"id":1,"is_bot":true,"first_name":"DragonKi","username":"dragonki_bot"}}`))
		case strings.HasSuffix(r.URL.Path, "/sendMessage"):
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":1,"date":0,"chat":{"id":1,"type":"private"},"text":"ok"}}`))
		case strings.HasSuffix(r.URL.Path, "/sendPhoto"):
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":2,"date":0,"chat":{"id":1,"type":"private"},"caption":"ok"}}`))
		default:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":true,"result":{}}`))
		}
	}))
	defer server.Close()

	bot, err := tgbotapi.NewBotAPIWithAPIEndpoint("TEST_TOKEN", server.URL+"/bot%s/%s")
	if err != nil {
		t.Fatalf("init fake telegram bot: %v", err)
	}

	t.Setenv("REFERRAL_REQUIRED_LEVEL", "9")
	t.Setenv("REFERRAL_REWARD_GOLD", "1500")
	t.Setenv("REFERRAL_REWARD_KI", "3")

	app, err := buildBotApp(botAppConfig{
		DB:             testDB.DB,
		Bot:            bot,
		MPClient:       &mercadopago.Client{Ctx: context.Background()},
		AutoMigrate:    false,
		RunMigrations:  false,
		StartScheduler: false,
		FeatureFlagResolver: func(name string) (bool, bool) {
			return false, false
		},
	})
	if err != nil {
		t.Fatalf("build bot app: %v", err)
	}
	defer app.Cleanup()

	update := tgbotapi.Update{
		UpdateID: 1001,
		Message: &tgbotapi.Message{
			MessageID: 1,
			Chat:      &tgbotapi.Chat{ID: 920001, Type: "private"},
			From:      &tgbotapi.User{ID: 920001, FirstName: "Goku", UserName: "goku"},
			Text:      "/start ref_" + strconv.FormatUint(uint64(referrer.ID), 10),
		},
	}

	body, err := json.Marshal(update)
	if err != nil {
		t.Fatalf("marshal telegram update: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/webhook/telegram", strings.NewReader(string(body)))
	rec := httptest.NewRecorder()
	app.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	referred, err := playerRepo.GetOrCreate(920001)
	if err != nil {
		t.Fatalf("load referred player: %v", err)
	}
	if referred.OnboardingStep != "ONBOARDING" {
		t.Fatalf("expected onboarding step ONBOARDING, got %q", referred.OnboardingStep)
	}
	if referred.UserIDTelegram != 920001 || referred.NameTelegram != "Goku" || referred.UserName != "goku" {
		t.Fatalf("unexpected telegram fields persisted: %+v", referred)
	}

	referral, err := referralRepo.GetByReferredPlayerID(int64(referred.ID))
	if err != nil {
		t.Fatalf("load referral by referred player id: %v", err)
	}
	if referral.ReferrerPlayerID != int64(referrer.ID) {
		t.Fatalf("unexpected referrer player id: got %d want %d", referral.ReferrerPlayerID, referrer.ID)
	}
	if referral.ReferredPlayerID != int64(referred.ID) {
		t.Fatalf("unexpected referred player id: got %d want %d", referral.ReferredPlayerID, referred.ID)
	}
	if referral.Status != domain.ReferralStatusPending {
		t.Fatalf("unexpected referral status: %q", referral.Status)
	}
	if referral.RequiredLevel != 9 {
		t.Fatalf("unexpected required level: got %d want 9", referral.RequiredLevel)
	}

	mu.Lock()
	defer mu.Unlock()

	var referredPanel bool
	var referrerMessage bool
	for _, request := range requests {
		if request.ChatID == "920001" && strings.HasSuffix(request.Path, "/sendPhoto") && strings.Contains(request.Caption, "Convite registrado com sucesso") {
			referredPanel = true
		}
		if request.ChatID == "910001" && strings.HasSuffix(request.Path, "/sendMessage") && strings.Contains(request.Text, "novo jogador entrou pelo seu link") {
			referrerMessage = true
		}
	}
	if !referredPanel {
		t.Fatal("expected referral notice embedded in referred player panel")
	}
	if !referrerMessage {
		t.Fatal("expected notification message to referrer")
	}
}
