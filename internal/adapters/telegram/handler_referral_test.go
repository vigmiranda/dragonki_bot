package telegram

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/application"
	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
)

type tgSentMessage struct {
	chatID string
	text   string
}

type fakePlayerRepoTG struct {
	player *domain.Player
}

func (f *fakePlayerRepoTG) GetOrCreate(userID int64) (*domain.Player, error) {
	return f.player, nil
}

func (f *fakePlayerRepoTG) Update(player *domain.Player) error { return nil }

func (f *fakePlayerRepoTG) CountOnlineLast15Min() (int, error) { return 0, nil }

func (f *fakePlayerRepoTG) GetClassesByLvl(lvl int) ([]domain.Class, error) { return nil, nil }

func (f *fakePlayerRepoTG) GetClassesByID(id int64) (*domain.Class, error) { return nil, nil }

type fakeReferralRepoTG struct {
	referral           *domain.Referral
	referrerUserIDByID map[int64]int64
}

func (f *fakeReferralRepoTG) Create(referral *domain.Referral) error {
	ref := *referral
	ref.ID = 1
	f.referral = &ref
	return nil
}

func (f *fakeReferralRepoTG) GetByReferredPlayerID(referredPlayerID int64) (*domain.Referral, error) {
	if f.referral == nil || f.referral.ReferredPlayerID != referredPlayerID {
		return nil, gorm.ErrRecordNotFound
	}
	ref := *f.referral
	return &ref, nil
}

func (f *fakeReferralRepoTG) MarkRewarded(referralID uint) error { return nil }

func (f *fakeReferralRepoTG) GrantReward(referrerPlayerID int64, gold int, kiCrystals int) error {
	return nil
}

func (f *fakeReferralRepoTG) CountByReferrer(referrerPlayerID int64) (pending int64, rewarded int64, err error) {
	return 0, 0, nil
}

func (f *fakeReferralRepoTG) GetPlayerUserIDByPlayerID(playerID int64) (int64, error) {
	return f.referrerUserIDByID[playerID], nil
}

func TestHandleUpdate_StartWithReferral_SendsMessagesToReferredAndReferrer(t *testing.T) {
	var mu sync.Mutex
	sent := make([]tgSentMessage, 0, 2)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseMultipartForm(10 << 20)
		_ = r.ParseForm()

		switch {
		case strings.HasSuffix(r.URL.Path, "/getMe"):
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":true,"result":{"id":1,"is_bot":true,"first_name":"DragonKi","username":"dragonki_bot"}}`))
			return
		case strings.HasSuffix(r.URL.Path, "/sendMessage"):
			mu.Lock()
			sent = append(sent, tgSentMessage{
				chatID: r.FormValue("chat_id"),
				text:   r.FormValue("text"),
			})
			mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":1,"date":0,"chat":{"id":1,"type":"private"},"text":"ok"}}`))
			return
		case strings.HasSuffix(r.URL.Path, "/sendPhoto"):
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":2,"date":0,"chat":{"id":1,"type":"private"},"caption":"ok"}}`))
			return
		default:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":true,"result":{}}`))
			return
		}
	}))
	defer server.Close()

	bot, err := tgbotapi.NewBotAPIWithAPIEndpoint("TEST_TOKEN", server.URL+"/bot%s/%s")
	if err != nil {
		t.Fatalf("failed to init bot api: %v", err)
	}

	playerRepo := &fakePlayerRepoTG{
		player: &domain.Player{
			ID:             200,
			UserID:         200,
			OnboardingStep: "ONBOARDING",
		},
	}
	referralRepo := &fakeReferralRepoTG{
		referrerUserIDByID: map[int64]int64{
			100: 999,
		},
	}

	h := NewHandler(
		bot,
		nil,
		nil,
		nil,
		application.NewOnboardingService(playerRepo),
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		application.NewReferralService(referralRepo, playerRepo, bot.Self.UserName),
	)

	h.HandleUpdate(tgbotapi.Update{
		UpdateID: 1,
		Message: &tgbotapi.Message{
			MessageID: 1,
			Chat:      &tgbotapi.Chat{ID: 200, Type: "private"},
			From:      &tgbotapi.User{ID: 200, FirstName: "Goku", UserName: "goku"},
			Text:      "/start ref_100",
		},
	})

	mu.Lock()
	defer mu.Unlock()

	if len(sent) < 2 {
		t.Fatalf("expected at least 2 sendMessage calls, got %d", len(sent))
	}

	var hasReferredMsg bool
	var hasReferrerMsg bool
	for _, m := range sent {
		if m.chatID == "200" && strings.Contains(m.text, "Convite registrado com sucesso") {
			hasReferredMsg = true
		}
		if m.chatID == "999" && strings.Contains(m.text, "novo jogador entrou pelo seu link") {
			hasReferrerMsg = true
		}
	}
	if !hasReferredMsg {
		t.Fatalf("expected confirmation message to referred player")
	}
	if !hasReferrerMsg {
		t.Fatalf("expected notification message to referrer")
	}
}

func TestHandleUpdate_StartWithReferral_DuplicateDoesNotSendReferralMessages(t *testing.T) {
	var mu sync.Mutex
	sent := make([]tgSentMessage, 0, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseMultipartForm(10 << 20)
		_ = r.ParseForm()

		switch {
		case strings.HasSuffix(r.URL.Path, "/getMe"):
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":true,"result":{"id":1,"is_bot":true,"first_name":"DragonKi","username":"dragonki_bot"}}`))
			return
		case strings.HasSuffix(r.URL.Path, "/sendMessage"):
			mu.Lock()
			sent = append(sent, tgSentMessage{
				chatID: r.FormValue("chat_id"),
				text:   r.FormValue("text"),
			})
			mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":1,"date":0,"chat":{"id":1,"type":"private"},"text":"ok"}}`))
			return
		case strings.HasSuffix(r.URL.Path, "/sendPhoto"):
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":2,"date":0,"chat":{"id":1,"type":"private"},"caption":"ok"}}`))
			return
		default:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":true,"result":{}}`))
			return
		}
	}))
	defer server.Close()

	bot, err := tgbotapi.NewBotAPIWithAPIEndpoint("TEST_TOKEN", server.URL+"/bot%s/%s")
	if err != nil {
		t.Fatalf("failed to init bot api: %v", err)
	}

	playerRepo := &fakePlayerRepoTG{
		player: &domain.Player{
			ID:             200,
			UserID:         200,
			OnboardingStep: "ONBOARDING",
		},
	}
	referralRepo := &fakeReferralRepoTG{
		referral: &domain.Referral{
			ID:               99,
			ReferrerPlayerID: 100,
			ReferredPlayerID: 200,
			Status:           domain.ReferralStatusPending,
			RequiredLevel:    7,
		},
		referrerUserIDByID: map[int64]int64{
			100: 999,
		},
	}

	h := NewHandler(
		bot,
		nil,
		nil,
		nil,
		application.NewOnboardingService(playerRepo),
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		application.NewReferralService(referralRepo, playerRepo, bot.Self.UserName),
	)

	h.HandleUpdate(tgbotapi.Update{
		UpdateID: 2,
		Message: &tgbotapi.Message{
			MessageID: 2,
			Chat:      &tgbotapi.Chat{ID: 200, Type: "private"},
			From:      &tgbotapi.User{ID: 200, FirstName: "Goku", UserName: "goku"},
			Text:      "/start ref_100",
		},
	})

	mu.Lock()
	defer mu.Unlock()
	if len(sent) != 0 {
		t.Fatalf("expected no referral sendMessage on duplicate referral, got %d", len(sent))
	}
}
