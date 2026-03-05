package utils

import (
	"encoding/base64"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/domain"
)

func TestDecodeBase64(t *testing.T) {
	raw := []byte("dragonki")
	enc := base64.StdEncoding.EncodeToString(raw)

	got := DecodeBase64(enc)
	if string(got) != "dragonki" {
		t.Fatalf("expected dragonki, got %q", string(got))
	}

	gotPng := DecodeBase64("data:image/png;base64," + enc)
	if string(gotPng) != "dragonki" {
		t.Fatalf("expected dragonki with png prefix")
	}

	gotJpg := DecodeBase64("data:image/jpeg;base64," + enc)
	if string(gotJpg) != "dragonki" {
		t.Fatalf("expected dragonki with jpeg prefix")
	}

	if bad := DecodeBase64("%%%"); bad != nil {
		t.Fatalf("expected nil for invalid base64")
	}
}

func TestParseCallbackData(t *testing.T) {
	action, params := ParseCallbackData("action_menu")
	if action != "action_menu" {
		t.Fatalf("expected action_menu, got %q", action)
	}
	if len(params) != 0 {
		t.Fatalf("expected no params")
	}

	action, params = ParseCallbackData("buy_item?itemID=10&foo=bar")
	if action != "buy_item" {
		t.Fatalf("expected buy_item, got %q", action)
	}
	if params["itemID"] != "10" || params["foo"] != "bar" {
		t.Fatalf("unexpected params: %+v", params)
	}

	action, params = ParseCallbackData("a?%%%")
	if action != "a" {
		t.Fatalf("expected action a, got %q", action)
	}
	if len(params) != 0 {
		t.Fatalf("expected empty params on parse error")
	}
}

func TestNewPhotoMedia(t *testing.T) {
	media := NewPhotoMedia(tgbotapi.FilePath("assets/a.png"), "hello")
	if media == nil {
		t.Fatalf("expected non-nil media")
	}
	if media.Caption != "hello" || media.ParseMode != "HTML" {
		t.Fatalf("expected caption and HTML parse mode, got caption=%q parseMode=%q", media.Caption, media.ParseMode)
	}

	noCaption := NewPhotoMedia(tgbotapi.FilePath("assets/a.png"), "")
	if noCaption.ParseMode != "" || noCaption.Caption != "" {
		t.Fatalf("expected empty caption/parse mode")
	}
}

func TestBuildMenuItems(t *testing.T) {
	menu := BuildMenuItems(map[string]any{
		"Potion A": "10",
		"Potion B": "20",
	})
	if menu == nil {
		t.Fatalf("expected non-nil menu")
	}
	if len(menu.InlineKeyboard) != 3 {
		t.Fatalf("expected 3 rows (2 items + back), got %d", len(menu.InlineKeyboard))
	}
	last := menu.InlineKeyboard[len(menu.InlineKeyboard)-1][0]
	if last.Text == "" || last.CallbackData == nil || *last.CallbackData != "to_shop" {
		t.Fatalf("expected back button to_shop, got %+v", last)
	}
}

func TestChoiceClassMenu(t *testing.T) {
	menu := ChoiceClassMenu(map[string]any{
		"Classes": []domain.Class{
			{ID: 1, Name: "Wind"},
			{ID: 2, Name: "Earth"},
		},
	})
	if menu == nil {
		t.Fatalf("expected non-nil menu")
	}
	if len(menu.InlineKeyboard) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(menu.InlineKeyboard))
	}
	cb := menu.InlineKeyboard[0][0].CallbackData
	if cb == nil || *cb != "to_class?classID=1" {
		t.Fatalf("unexpected callback: %+v", cb)
	}
}

func TestRegionsMenu(t *testing.T) {
	menu := RegionsMenu(map[string]any{
		"Regions": []domain.Region{
			{RegionID: 1, Name: "Plains", MinLevel: 1},
			{RegionID: 2, Name: "Forest", MinLevel: 15},
		},
	})
	if menu == nil {
		t.Fatalf("expected non-nil menu")
	}
	if len(menu.InlineKeyboard) != 3 {
		t.Fatalf("expected 3 rows (2 regions + back), got %d", len(menu.InlineKeyboard))
	}
	last := menu.InlineKeyboard[len(menu.InlineKeyboard)-1][0]
	if last.CallbackData == nil || *last.CallbackData != "action_menu" {
		t.Fatalf("expected back callback action_menu")
	}
}

func TestPotionMenu(t *testing.T) {
	menu := PotionMenu(map[string]any{
		"Items": []domain.InventoryItem{
			{ID: 10, Quantity: 2, Item: domain.Item{Name: "HP"}},
		},
		"page":       2,
		"totalPages": 3,
	})
	if menu == nil {
		t.Fatalf("expected non-nil menu")
	}
	// item row + nav row + back row
	if len(menu.InlineKeyboard) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(menu.InlineKeyboard))
	}
	itemCb := menu.InlineKeyboard[0][0].CallbackData
	if itemCb == nil || *itemCb != "use_item?inventoryItemID=10" {
		t.Fatalf("unexpected item callback: %+v", itemCb)
	}
	backCb := menu.InlineKeyboard[2][0].CallbackData
	if backCb == nil || *backCb != "open_inventory" {
		t.Fatalf("expected open_inventory back callback")
	}
}

func TestToInt(t *testing.T) {
	if got := toInt(2); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
	if got := toInt(int64(3)); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
	if got := toInt(float64(4.9)); got != 4 {
		t.Fatalf("expected 4, got %d", got)
	}
	if got := toInt("x"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}
