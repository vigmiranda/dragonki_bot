package telegram

import (
	"strings"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/bootstrap"
)

func TestMainMenu_HidesBigBossWhenFeatureDisabled(t *testing.T) {
	bootstrap.SetFeatureFlagResolver(func(name string) (bool, bool) {
		if name == "BIG_BOSS_WEEKLY" {
			return false, true
		}
		return false, false
	})
	t.Cleanup(func() { bootstrap.SetFeatureFlagResolver(nil) })

	menu := mainMenu()
	if !menuHasCallback(menu, "menu_events") {
		t.Fatalf("expected events scope button visible on main menu")
	}

	events := eventsMenu()
	if menuHasCallback(events, "to_big_boss") {
		t.Fatalf("expected Big Boss hidden inside events menu when feature is disabled")
	}
}

func TestMainMenu_ShowsBigBossWhenFeatureEnabled(t *testing.T) {
	bootstrap.SetFeatureFlagResolver(func(name string) (bool, bool) {
		if name == "BIG_BOSS_WEEKLY" {
			return true, true
		}
		return false, false
	})
	t.Cleanup(func() { bootstrap.SetFeatureFlagResolver(nil) })

	menu := mainMenu()
	if !menuHasCallback(menu, "menu_events") {
		t.Fatalf("expected events scope button visible on main menu")
	}

	events := eventsMenu()
	if !menuHasCallback(events, "to_big_boss") {
		t.Fatalf("expected Big Boss button visible inside events menu when feature is enabled")
	}
}

func TestHandleBigBoss_BlockedWhenFeatureDisabled(t *testing.T) {
	bootstrap.SetFeatureFlagResolver(func(name string) (bool, bool) {
		if name == "BIG_BOSS_WEEKLY" {
			return false, true
		}
		return false, false
	})
	t.Cleanup(func() { bootstrap.SetFeatureFlagResolver(nil) })

	h := &Handler{}
	got := h.handleBigBoss("to_big_boss", 1)
	if !strings.Contains(got.Message, "indisponivel") {
		t.Fatalf("expected disabled message, got %q", got.Message)
	}
}

func TestHandleBigBoss_FeatureEnabledWithoutServiceShowsUnavailable(t *testing.T) {
	bootstrap.SetFeatureFlagResolver(func(name string) (bool, bool) {
		if name == "BIG_BOSS_WEEKLY" {
			return true, true
		}
		return false, false
	})
	t.Cleanup(func() { bootstrap.SetFeatureFlagResolver(nil) })

	h := &Handler{}
	got := h.handleBigBoss("to_big_boss", 1)
	if !strings.Contains(got.Message, "indisponivel") {
		t.Fatalf("expected unavailable message, got %q", got.Message)
	}
}

func menuHasCallback(menu *tgbotapi.InlineKeyboardMarkup, callback string) bool {
	if menu == nil {
		return false
	}
	for _, row := range menu.InlineKeyboard {
		for _, btn := range row {
			if btn.CallbackData != nil && *btn.CallbackData == callback {
				return true
			}
		}
	}
	return false
}
