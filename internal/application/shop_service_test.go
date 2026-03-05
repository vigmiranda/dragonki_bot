package application

import (
	"errors"
	"strings"
	"testing"

	"github.com/vigmiranda/dragonki/internal/domain"
)

func TestNewShopService(t *testing.T) {
	playerRepo := &fakePlayerRepo{}
	itemRepo := &fakeItemsRepo{}
	inventoryRepo := &fakeInventoryRepo{}

	svc := NewShopService(playerRepo, itemRepo, inventoryRepo)
	if svc == nil {
		t.Fatalf("expected non-nil service")
	}
	if svc.PlayerRepository != playerRepo {
		t.Fatalf("expected player repo set")
	}
	if svc.ItemRepository != itemRepo {
		t.Fatalf("expected item repo set")
	}
	if svc.InventoryRepository != inventoryRepo {
		t.Fatalf("expected inventory repo set")
	}
}

func TestShopService_OpenShop_Success(t *testing.T) {
	player := basePlayer()
	player.Name = "Gohan"

	svc := &ShopService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
	}

	resp, err := svc.OpenShop(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "SHOP_MENU" {
		t.Fatalf("expected SHOP_MENU, got %q", resp.MenuKey)
	}
	if !strings.Contains(resp.Text, "Gohan") {
		t.Fatalf("expected player name in text, got %q", resp.Text)
	}
	if resp.Image == nil || resp.Image.Path != "assets/npcs/kiara_mercadora.png" {
		t.Fatalf("unexpected image: %+v", resp.Image)
	}
}

func TestShopService_ListPotions_Success(t *testing.T) {
	items := []domain.Item{
		{ID: 11, Name: "Pocao HP", Amount: 1, Price: 100},
		{ID: 22, Name: "Pocao Energia", Amount: 2, Price: 250},
	}

	svc := &ShopService{
		ItemRepository: &fakeItemsRepo{
			getRandomListFn: func(itemType domain.ItemType, rarity domain.ItemRarity) ([]domain.Item, error) {
				if itemType != domain.ItemConsumable {
					t.Fatalf("expected CONSUMABLE, got %s", itemType)
				}
				if rarity != domain.Common {
					t.Fatalf("expected COMMON, got %s", rarity)
				}
				return items, nil
			},
		},
	}

	resp, err := svc.ListPotions()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "SHOP_ITEMS_MENU" {
		t.Fatalf("expected SHOP_ITEMS_MENU, got %q", resp.MenuKey)
	}
	if resp.Image == nil || resp.Image.Path != "assets/npcs/kiara_mercadora.png" {
		t.Fatalf("unexpected image: %+v", resp.Image)
	}
	if len(resp.MenuParams) != 2 {
		t.Fatalf("expected 2 menu params, got %d", len(resp.MenuParams))
	}
}

func TestShopService_ListEquipments_Success(t *testing.T) {
	items := []domain.Item{
		{ID: 3, Name: "Espada", Amount: 1, Price: 700},
	}

	svc := &ShopService{
		ItemRepository: &fakeItemsRepo{
			getRandomListFn: func(itemType domain.ItemType, rarity domain.ItemRarity) ([]domain.Item, error) {
				if itemType != domain.ItemEquipment {
					t.Fatalf("expected EQUIPMENT, got %s", itemType)
				}
				if rarity != domain.Common {
					t.Fatalf("expected COMMON, got %s", rarity)
				}
				return items, nil
			},
		},
	}

	resp, err := svc.ListEquipments()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "SHOP_ITEMS_MENU" {
		t.Fatalf("expected SHOP_ITEMS_MENU, got %q", resp.MenuKey)
	}
	if resp.Image == nil || resp.Image.Path != "assets/npcs/kiara_mercadora.png" {
		t.Fatalf("unexpected image: %+v", resp.Image)
	}
	if !strings.Contains(resp.Text, "equipamento") && !strings.Contains(resp.Text, "equip") {
		t.Fatalf("unexpected text: %q", resp.Text)
	}
}

func TestShopService_BuyItem_InsufficientGold(t *testing.T) {
	player := basePlayer()
	player.ID = 9
	player.Gold = 50

	invRepo := &fakeInventoryRepo{}
	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	itemRepo := &fakeItemsRepo{
		getItemByIDFn: func(id int) (*domain.Item, error) {
			return &domain.Item{ID: 99, Name: "Elixir", Price: 100, Amount: 1}, nil
		},
	}
	svc := &ShopService{
		PlayerRepository:    playerRepo,
		ItemRepository:      itemRepo,
		InventoryRepository: invRepo,
	}

	resp, err := svc.BuyItem(1, map[string]string{"itemID": "99"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Text, "n") {
		t.Fatalf("expected insufficient gold message, got %q", resp.Text)
	}
	if player.Gold != 50 {
		t.Fatalf("expected unchanged gold, got %d", player.Gold)
	}
	if len(invRepo.createOrUpdateCalls) != 0 {
		t.Fatalf("expected no inventory writes")
	}
	if len(playerRepo.updatedPlayers) != 0 {
		t.Fatalf("expected no player update")
	}
}

func TestShopService_BuyItem_Success(t *testing.T) {
	player := basePlayer()
	player.ID = 12
	player.Gold = 1000

	invRepo := &fakeInventoryRepo{}
	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	itemRepo := &fakeItemsRepo{
		getItemByIDFn: func(id int) (*domain.Item, error) {
			return &domain.Item{ID: 7, Name: "Pocao", Price: 300, Amount: 2}, nil
		},
	}
	svc := &ShopService{
		PlayerRepository:    playerRepo,
		ItemRepository:      itemRepo,
		InventoryRepository: invRepo,
	}

	resp, err := svc.BuyItem(1, map[string]string{"itemID": "7"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if player.Gold != 700 {
		t.Fatalf("expected gold 700, got %d", player.Gold)
	}
	if len(invRepo.createOrUpdateCalls) != 1 {
		t.Fatalf("expected one inventory write, got %d", len(invRepo.createOrUpdateCalls))
	}
	call := invRepo.createOrUpdateCalls[0]
	if call.playerID != int64(player.ID) || call.itemID != 7 || call.amount != 2 {
		t.Fatalf("unexpected inventory write: %+v", call)
	}
	if len(playerRepo.updatedPlayers) != 1 {
		t.Fatalf("expected one player update")
	}
	if !strings.Contains(resp.Text, "Obrigada") {
		t.Fatalf("expected success message, got %q", resp.Text)
	}
}

func TestShopService_BuyItem_CreateInventoryError(t *testing.T) {
	player := basePlayer()
	player.ID = 5
	player.Gold = 1000

	svc := &ShopService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
		ItemRepository: &fakeItemsRepo{
			getItemByIDFn: func(id int) (*domain.Item, error) {
				return &domain.Item{ID: 1, Price: 100, Amount: 1}, nil
			},
		},
		InventoryRepository: &fakeInventoryRepo{
			err: errors.New("inventory fail"),
		},
	}

	_, err := svc.BuyItem(1, map[string]string{"itemID": "1"})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestShopService_BuyItem_UpdatePlayerError(t *testing.T) {
	player := basePlayer()
	player.ID = 8
	player.Gold = 1000

	svc := &ShopService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
			updateFn: func(player *domain.Player) error {
				return errors.New("save player fail")
			},
		},
		ItemRepository: &fakeItemsRepo{
			getItemByIDFn: func(id int) (*domain.Item, error) {
				return &domain.Item{ID: 2, Price: 200, Amount: 1}, nil
			},
		},
		InventoryRepository: &fakeInventoryRepo{},
	}

	_, err := svc.BuyItem(1, map[string]string{"itemID": "2"})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestShopService_BuyItem_InvalidItemID_CurrentBehaviorUsesZero(t *testing.T) {
	player := basePlayer()
	player.ID = 21
	player.Gold = 1000

	capturedID := -1
	svc := &ShopService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
		ItemRepository: &fakeItemsRepo{
			getItemByIDFn: func(id int) (*domain.Item, error) {
				capturedID = id
				return &domain.Item{ID: 0, Price: 10, Amount: 1}, nil
			},
		},
		InventoryRepository: &fakeInventoryRepo{},
	}

	_, err := svc.BuyItem(1, map[string]string{"itemID": "abc"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedID != 0 {
		t.Fatalf("expected parsed id fallback to 0, got %d", capturedID)
	}
}
