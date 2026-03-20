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
		"Potion A":    "10",
		"Potion B":    "20",
		"_backAction": "list_equipments",
	})
	if menu == nil {
		t.Fatalf("expected non-nil menu")
	}
	if len(menu.InlineKeyboard) != 3 {
		t.Fatalf("expected 3 rows (2 items + back), got %d", len(menu.InlineKeyboard))
	}
	last := menu.InlineKeyboard[len(menu.InlineKeyboard)-1][0]
	if last.Text == "" || last.CallbackData == nil || *last.CallbackData != "list_equipments" {
		t.Fatalf("expected back button list_equipments, got %+v", last)
	}
}

func TestShopEquipmentTypesMenu(t *testing.T) {
	menu := ShopEquipmentTypesMenu()
	if menu == nil {
		t.Fatalf("expected non-nil menu")
	}
	if len(menu.InlineKeyboard) != 9 {
		t.Fatalf("expected 9 rows, got %d", len(menu.InlineKeyboard))
	}
	firstCb := menu.InlineKeyboard[0][0].CallbackData
	if firstCb == nil || *firstCb != "list_equipments_type?slot=WEAPON" {
		t.Fatalf("unexpected first callback: %+v", firstCb)
	}
	accessoryCb := menu.InlineKeyboard[6][0].CallbackData
	if accessoryCb == nil || *accessoryCb != "list_equipments_type?slot=ACCESSORY" {
		t.Fatalf("unexpected accessory callback: %+v", accessoryCb)
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
		"page":             2,
		"totalPages":       3,
		"paginationAction": "open_inventory_potions",
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
	prevCb := menu.InlineKeyboard[1][0].CallbackData
	if prevCb == nil || *prevCb != "open_inventory_potions?page=1" {
		t.Fatalf("unexpected prev callback: %+v", prevCb)
	}
	nextCb := menu.InlineKeyboard[1][1].CallbackData
	if nextCb == nil || *nextCb != "open_inventory_potions?page=3" {
		t.Fatalf("unexpected next callback: %+v", nextCb)
	}
	backCb := menu.InlineKeyboard[2][0].CallbackData
	if backCb == nil || *backCb != "open_inventory" {
		t.Fatalf("expected open_inventory back callback")
	}
}

func TestPotionMenu_EquipmentUsesEquipCallback(t *testing.T) {
	menu := PotionMenu(map[string]any{
		"Items": []domain.InventoryItem{
			{ID: 77, Quantity: 1, Item: domain.Item{Name: "Sword"}},
		},
		"page":             1,
		"totalPages":       1,
		"itemType":         string(domain.ItemEquipment),
		"paginationAction": "open_inventory_equipaments_type?slot=WEAPON",
		"backAction":       "open_inventory_equipaments",
	})
	if menu == nil {
		t.Fatalf("expected non-nil menu")
	}
	if len(menu.InlineKeyboard) != 2 {
		t.Fatalf("expected 2 rows (item + back), got %d", len(menu.InlineKeyboard))
	}
	itemCb := menu.InlineKeyboard[0][0].CallbackData
	if itemCb == nil || *itemCb != "equip_item?inventoryItemID=77" {
		t.Fatalf("unexpected item callback: %+v", itemCb)
	}
	backCb := menu.InlineKeyboard[1][0].CallbackData
	if backCb == nil || *backCb != "open_inventory_equipaments" {
		t.Fatalf("unexpected back callback: %+v", backCb)
	}
}

func TestPotionMenu_WithExistingQueryInPaginationAction(t *testing.T) {
	menu := PotionMenu(map[string]any{
		"Items": []domain.InventoryItem{
			{ID: 77, Quantity: 1, Item: domain.Item{Name: "Sword"}},
		},
		"page":             2,
		"totalPages":       3,
		"itemType":         string(domain.ItemEquipment),
		"paginationAction": "open_inventory_equipaments_type?slot=WEAPON",
	})
	nav := menu.InlineKeyboard[1]
	prevCb := nav[0].CallbackData
	nextCb := nav[1].CallbackData
	if prevCb == nil || *prevCb != "open_inventory_equipaments_type?slot=WEAPON&page=1" {
		t.Fatalf("unexpected prev callback: %+v", prevCb)
	}
	if nextCb == nil || *nextCb != "open_inventory_equipaments_type?slot=WEAPON&page=3" {
		t.Fatalf("unexpected next callback: %+v", nextCb)
	}
}

func TestEquipmentTypesMenu(t *testing.T) {
	menu := EquipmentTypesMenu(map[string]any{
		"weaponCount":    int64(2),
		"armorCount":     int64(1),
		"glovesCount":    int64(5),
		"helmetCount":    int64(0),
		"bootsCount":     int64(3),
		"ringCount":      int64(4),
		"accessoryCount": int64(2),
		"amuletCount":    int64(6),
	})
	if menu == nil {
		t.Fatalf("expected non-nil menu")
	}
	if len(menu.InlineKeyboard) != 10 {
		t.Fatalf("expected 10 rows, got %d", len(menu.InlineKeyboard))
	}
	if menu.InlineKeyboard[0][0].Text != "Armas (2)" {
		t.Fatalf("unexpected label: %q", menu.InlineKeyboard[0][0].Text)
	}
	firstCb := menu.InlineKeyboard[0][0].CallbackData
	if firstCb == nil || *firstCb != "open_inventory_equipaments_type?slot=WEAPON&page=1" {
		t.Fatalf("unexpected first callback: %+v", firstCb)
	}
	accessoryCb := menu.InlineKeyboard[6][0].CallbackData
	if accessoryCb == nil || *accessoryCb != "open_inventory_equipaments_type?slot=ACCESSORY&page=1" {
		t.Fatalf("unexpected accessory callback: %+v", accessoryCb)
	}
	equippedCb := menu.InlineKeyboard[8][0].CallbackData
	if equippedCb == nil || *equippedCb != "open_inventory_equipped" {
		t.Fatalf("unexpected equipped callback: %+v", equippedCb)
	}
}

func TestEquippedItemsMenu(t *testing.T) {
	menu := EquippedItemsMenu(map[string]any{
		"weaponAction":     "unequip_item?slot=WEAPON",
		"ringAction":       "unequip_item?slot=RING",
		"accessoryAction":  "unequip_item?slot=ACCESSORY",
		"unequipAllAction": "unequip_all",
	})
	if menu == nil {
		t.Fatalf("expected non-nil menu")
	}
	if len(menu.InlineKeyboard) != 5 {
		t.Fatalf("expected 5 rows (4 actions + back), got %d", len(menu.InlineKeyboard))
	}
	firstCb := menu.InlineKeyboard[0][0].CallbackData
	if firstCb == nil || *firstCb != "unequip_item?slot=WEAPON" {
		t.Fatalf("unexpected first callback: %+v", firstCb)
	}
	accessoryCb := menu.InlineKeyboard[2][0].CallbackData
	if accessoryCb == nil || *accessoryCb != "unequip_item?slot=ACCESSORY" {
		t.Fatalf("expected accessory callback, got %+v", accessoryCb)
	}
	allCb := menu.InlineKeyboard[3][0].CallbackData
	if allCb == nil || *allCb != "unequip_all" {
		t.Fatalf("expected unequip_all callback, got %+v", allCb)
	}
	lastCb := menu.InlineKeyboard[4][0].CallbackData
	if lastCb == nil || *lastCb != "open_inventory" {
		t.Fatalf("expected open_inventory back callback, got %+v", lastCb)
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
