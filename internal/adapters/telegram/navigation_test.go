package telegram

import (
	"strings"
	"testing"

	"github.com/vigmiranda/dragonki/internal/application"
)

func TestMainMenu_InventoryButtonOpensInventoryDirectly(t *testing.T) {
	menu := mainMenu()

	if !menuHasCallback(menu, "open_inventory") {
		t.Fatalf("expected main menu inventory shortcut to open inventory directly")
	}
	if menuHasCallback(menu, "menu_inventory_scope") {
		t.Fatalf("expected main menu inventory shortcut to stop using intermediate inventory scope")
	}
}

func TestInventoryMenu_ShowsEquippedShortcut(t *testing.T) {
	menu := inventoryMenu()

	if !menuHasCallback(menu, "open_inventory_equipped") {
		t.Fatalf("expected inventory menu to expose equipped items shortcut")
	}
}

func TestHandleInventory_MenuInventoryScopeUsesOpenInventoryFlow(t *testing.T) {
	h := &Handler{
		inventoryService: application.NewInventoryService(nil, nil, nil),
	}

	resp := h.handleInventory("menu_inventory_scope", 1, nil)

	if resp.CameFrom != "menu_inventory_scope" {
		t.Fatalf("expected legacy inventory scope to keep action origin, got %q", resp.CameFrom)
	}
	if !strings.Contains(resp.Message, "Arkanis") {
		t.Fatalf("expected legacy inventory scope to reuse inventory open response, got %q", resp.Message)
	}
	if resp.Menu == nil || !menuHasCallback(resp.Menu, "open_inventory_equipped") {
		t.Fatalf("expected legacy inventory scope to render the flattened inventory menu")
	}
}
