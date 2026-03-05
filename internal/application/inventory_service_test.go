package application

import (
	"errors"
	"strings"
	"testing"

	"github.com/vigmiranda/dragonki/internal/domain"
)

type fakeItemsRepo struct {
	getRandomListFn func(itemType domain.ItemType, rarity domain.ItemRarity) ([]domain.Item, error)
	getItemByIDFn   func(id int) (*domain.Item, error)
}

func (f *fakeItemsRepo) GetRandomList(itemType domain.ItemType, rarity domain.ItemRarity) ([]domain.Item, error) {
	if f.getRandomListFn == nil {
		return nil, nil
	}
	return f.getRandomListFn(itemType, rarity)
}

func (f *fakeItemsRepo) GetItemById(id int) (*domain.Item, error) {
	if f.getItemByIDFn == nil {
		return nil, nil
	}
	return f.getItemByIDFn(id)
}

func TestNewInventoryService(t *testing.T) {
	invRepo := &fakeInventoryRepo{}
	playerRepo := &fakePlayerRepo{}
	itemsRepo := &fakeItemsRepo{}

	svc := NewInventoryService(invRepo, playerRepo, itemsRepo)
	if svc == nil {
		t.Fatalf("expected non-nil service")
	}
	if svc.InventoryRepository != invRepo {
		t.Fatalf("expected inventory repo to be set")
	}
	if svc.PlayerRepository != playerRepo {
		t.Fatalf("expected player repo to be set")
	}
	if svc.ItemsRepository != itemsRepo {
		t.Fatalf("expected items repo to be set")
	}
}

func TestInventoryService_Open(t *testing.T) {
	svc := InventoryService{}
	resp, err := svc.Open(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "INVENTORY_MENU" {
		t.Fatalf("expected INVENTORY_MENU, got %q", resp.MenuKey)
	}
	if resp.Image == nil || resp.Image.Path != "assets/npcs/arkanis_inventario.png" {
		t.Fatalf("unexpected image: %+v", resp.Image)
	}
}

func TestInventoryService_ListItemsByUserAndType_Consumables(t *testing.T) {
	player := basePlayer()
	player.ID = 42
	inventoryItems := []domain.InventoryItem{
		{ID: 1, Item: domain.Item{Name: "Potion"}, Quantity: 2},
	}

	invRepo := &fakeInventoryRepo{
		listByTypeFn: func(playerID int64, itemType domain.ItemType, page int, pageSize int) ([]domain.InventoryItem, int64, error) {
			if playerID != 42 {
				t.Fatalf("expected playerID 42, got %d", playerID)
			}
			if itemType != domain.ItemConsumable {
				t.Fatalf("expected CONSUMABLE, got %s", itemType)
			}
			if page != 1 || pageSize != 5 {
				t.Fatalf("expected pagination 1/5, got %d/%d", page, pageSize)
			}
			return inventoryItems, 6, nil
		},
	}

	svc := InventoryService{
		InventoryRepository: invRepo,
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
	}

	resp, err := svc.ListItemsByUserAndType(1, domain.ItemConsumable)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "POTION_MENU" {
		t.Fatalf("expected POTION_MENU, got %q", resp.MenuKey)
	}
	if !strings.Contains(resp.Text, "po") {
		t.Fatalf("expected consumable text, got %q", resp.Text)
	}
	if resp.MenuParams["page"] != 1 {
		t.Fatalf("expected page=1, got %+v", resp.MenuParams["page"])
	}
	if resp.MenuParams["totalPages"] != 2 {
		t.Fatalf("expected totalPages=2, got %+v", resp.MenuParams["totalPages"])
	}
}

func TestInventoryService_ListItemsByUserAndType_RepoErrorIgnoredCurrentBehavior(t *testing.T) {
	player := basePlayer()
	player.ID = 7
	invRepo := &fakeInventoryRepo{
		listByTypeFn: func(playerID int64, itemType domain.ItemType, page int, pageSize int) ([]domain.InventoryItem, int64, error) {
			return nil, 0, errors.New("repo fail")
		},
	}

	svc := InventoryService{
		InventoryRepository: invRepo,
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
	}

	resp, err := svc.ListItemsByUserAndType(1, domain.ItemEquipment)
	if err != nil {
		t.Fatalf("expected nil error due current behavior, got %v", err)
	}
	if resp.MenuKey != "POTION_MENU" {
		t.Fatalf("expected POTION_MENU, got %q", resp.MenuKey)
	}
}

func TestInventoryService_UseItem_GetPlayerError(t *testing.T) {
	svc := InventoryService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) {
				return nil, errors.New("db fail")
			},
		},
		InventoryRepository: &fakeInventoryRepo{},
	}
	_, err := svc.UseItem(1, 10)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestInventoryService_UseItem_GetInventoryItemError(t *testing.T) {
	player := basePlayer()
	svc := InventoryService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
		InventoryRepository: &fakeInventoryRepo{
			getByUserFn: func(playerID int64, inventoryItemID uint) (*domain.InventoryItem, error) {
				return nil, errors.New("not found")
			},
		},
	}
	_, err := svc.UseItem(1, 10)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestInventoryService_UseItem_NoQuantity(t *testing.T) {
	player := basePlayer()
	inv := &domain.InventoryItem{ID: 1, Quantity: 0, Item: domain.Item{Name: "Potion"}}

	invRepo := &fakeInventoryRepo{
		getByUserFn: func(playerID int64, inventoryItemID uint) (*domain.InventoryItem, error) {
			return inv, nil
		},
	}
	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}

	svc := InventoryService{
		PlayerRepository:    playerRepo,
		InventoryRepository: invRepo,
	}

	resp, err := svc.UseItem(1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "MAIN_MENU" {
		t.Fatalf("expected MAIN_MENU, got %q", resp.MenuKey)
	}
	if !strings.Contains(resp.Text, "n") {
		t.Fatalf("expected out-of-stock message, got %q", resp.Text)
	}
	if len(invRepo.updateInventoryCalls) != 0 || len(playerRepo.updatedPlayers) != 0 {
		t.Fatalf("expected no updates")
	}
}

func TestInventoryService_UseItem_ApplyHPSmallAndPersist(t *testing.T) {
	player := basePlayer()
	player.HP = 70
	player.MaxHP = 100
	player.Gender = "MAN"
	player.Class.Skin = "basic"

	inv := &domain.InventoryItem{
		ID:       1,
		Code:     "POTION_HP_SMALL",
		Quantity: 2,
		Item:     domain.Item{Name: "Small HP Potion"},
	}

	invRepo := &fakeInventoryRepo{
		getByUserFn: func(playerID int64, inventoryItemID uint) (*domain.InventoryItem, error) {
			return inv, nil
		},
	}
	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	svc := InventoryService{
		PlayerRepository:    playerRepo,
		InventoryRepository: invRepo,
	}

	resp, err := svc.UseItem(1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if player.HP != 90 {
		t.Fatalf("expected HP 90, got %d", player.HP)
	}
	if inv.Quantity != 1 {
		t.Fatalf("expected inventory qty decremented to 1, got %d", inv.Quantity)
	}
	if len(invRepo.updateInventoryCalls) != 1 || len(playerRepo.updatedPlayers) != 1 {
		t.Fatalf("expected one inventory update and one player update")
	}
	if resp.Image == nil || !strings.Contains(resp.Image.Path, "assets/skins/man/basicplayer.png") {
		t.Fatalf("unexpected image path: %+v", resp.Image)
	}
}

func TestInventoryService_UseItem_ApplyEnergyLargeCap(t *testing.T) {
	player := basePlayer()
	player.Energy = 1
	player.MaxEnergy = 10

	inv := &domain.InventoryItem{
		ID:       1,
		Code:     "POTION_ENERGY_LARGE",
		Quantity: 1,
		Item:     domain.Item{Name: "Large Energy Potion"},
	}

	invRepo := &fakeInventoryRepo{
		getByUserFn: func(playerID int64, inventoryItemID uint) (*domain.InventoryItem, error) {
			return inv, nil
		},
	}
	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	svc := InventoryService{
		PlayerRepository:    playerRepo,
		InventoryRepository: invRepo,
	}

	_, err := svc.UseItem(1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if player.Energy != player.MaxEnergy {
		t.Fatalf("expected energy capped at %d, got %d", player.MaxEnergy, player.Energy)
	}
	if inv.Quantity != 0 {
		t.Fatalf("expected quantity 0, got %d", inv.Quantity)
	}
}

func TestInventoryService_UseItem_UpdateInventoryError(t *testing.T) {
	player := basePlayer()
	inv := &domain.InventoryItem{
		ID:       1,
		Code:     "POTION_HP_SMALL",
		Quantity: 1,
		Item:     domain.Item{Name: "Potion"},
	}
	invRepo := &fakeInventoryRepo{
		getByUserFn: func(playerID int64, inventoryItemID uint) (*domain.InventoryItem, error) {
			return inv, nil
		},
		updateFn: func(item *domain.InventoryItem) error {
			return errors.New("save inv fail")
		},
	}
	svc := InventoryService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
		InventoryRepository: invRepo,
	}
	_, err := svc.UseItem(1, 1)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestInventoryService_UseItem_UpdatePlayerError(t *testing.T) {
	player := basePlayer()
	inv := &domain.InventoryItem{
		ID:       1,
		Code:     "POTION_HP_SMALL",
		Quantity: 1,
		Item:     domain.Item{Name: "Potion"},
	}
	invRepo := &fakeInventoryRepo{
		getByUserFn: func(playerID int64, inventoryItemID uint) (*domain.InventoryItem, error) {
			return inv, nil
		},
	}
	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		updateFn: func(player *domain.Player) error {
			return errors.New("save player fail")
		},
	}
	svc := InventoryService{
		PlayerRepository:    playerRepo,
		InventoryRepository: invRepo,
	}
	_, err := svc.UseItem(1, 1)
	if err == nil {
		t.Fatalf("expected error")
	}
}
