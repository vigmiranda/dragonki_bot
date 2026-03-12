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
	svc := &ShopService{}

	resp, err := svc.ListEquipments()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "SHOP_EQUIPMENT_ACTION_MENU" {
		t.Fatalf("expected SHOP_EQUIPMENT_ACTION_MENU, got %q", resp.MenuKey)
	}
	if resp.Image == nil || resp.Image.Path != "assets/npcs/kiara_mercadora.png" {
		t.Fatalf("unexpected image: %+v", resp.Image)
	}
	if !strings.Contains(strings.ToLower(resp.Text), "comprar ou vender") {
		t.Fatalf("unexpected text: %q", resp.Text)
	}
}

func TestShopService_OpenBuyEquipmentTypes_Success(t *testing.T) {
	svc := &ShopService{}

	resp, err := svc.OpenBuyEquipmentTypes()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "SHOP_EQUIPMENT_TYPES_MENU" {
		t.Fatalf("expected SHOP_EQUIPMENT_TYPES_MENU, got %q", resp.MenuKey)
	}
}

func TestShopService_ListEquipmentsByType_Success(t *testing.T) {
	items := []domain.Item{
		{
			ID:             3,
			Name:           "Espada A",
			Amount:         1,
			Price:          700,
			Type:           domain.ItemEquipment,
			RequiredLevel:  15,
			AllowedClasses: "Wind Striker,Earth Sentinel",
			AtkBonus:       9,
		},
		{
			ID:             4,
			Name:           "Espada B",
			Amount:         1,
			Price:          1900,
			Type:           domain.ItemEquipment,
			RequiredLevel:  35,
			AllowedClasses: "Ki Vanguard",
			AtkBonus:       20,
			SpeedBonus:     2,
		},
	}

	svc := &ShopService{
		ItemRepository: &fakeItemsRepo{
			getRandomBySlotAndLimitFn: func(slot domain.EquipmentSlot, limit int) ([]domain.Item, error) {
				if slot != domain.SlotWeapon {
					t.Fatalf("expected WEAPON, got %s", slot)
				}
				if limit != 3 {
					t.Fatalf("expected limit 3, got %d", limit)
				}
				return items, nil
			},
		},
	}

	resp, err := svc.ListEquipmentsByType("WEAPON")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "SHOP_ITEMS_MENU" {
		t.Fatalf("expected SHOP_ITEMS_MENU, got %q", resp.MenuKey)
	}
	if len(resp.MenuParams) != 3 { // 2 itens + _backAction
		t.Fatalf("expected 3 params, got %d", len(resp.MenuParams))
	}
	foundCatalogPrice := false
	for label := range resp.MenuParams {
		if strings.Contains(label, "Espada A - 700 Gold") {
			foundCatalogPrice = true
		}
	}
	if !foundCatalogPrice {
		t.Fatalf("expected catalog equipment price in label, params=%+v", resp.MenuParams)
	}
	if !strings.Contains(resp.Text, "Espada A") || !strings.Contains(resp.Text, "Espada B") {
		t.Fatalf("expected listed item names in text, got %q", resp.Text)
	}
	if !strings.Contains(resp.Text, "Requisitos: Nivel 15+ | Classes: Wind Striker,Earth Sentinel") {
		t.Fatalf("expected requirements for first item in text, got %q", resp.Text)
	}
	if !strings.Contains(resp.Text, "Bonus: ATK +9") {
		t.Fatalf("expected bonus line for first item in text, got %q", resp.Text)
	}
}

func TestShopService_BuyItem_InsufficientGold(t *testing.T) {
	player := basePlayer()
	player.ID = 9
	player.Gold = 50

	invRepo := &fakeInventoryRepo{
		purchaseFn: func(playerID int64, itemID uint, amount int, price int) error {
			return domain.ErrInventoryInsufficientGold
		},
	}
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
	if len(invRepo.purchaseCalls) != 1 {
		t.Fatalf("expected one purchase attempt, got %d", len(invRepo.purchaseCalls))
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
			return &domain.Item{
				ID:             7,
				Name:           "Pocao",
				Price:          300,
				Amount:         2,
				RequiredLevel:  15,
				AllowedClasses: "Wind Striker,Earth Sentinel",
				AtkBonus:       3,
				DefBonus:       0,
				SpeedBonus:     1,
				HPBonus:        0,
				EnergyBonus:    2,
			}, nil
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
	if player.Gold != 1000 {
		t.Fatalf("expected no direct gold mutation, got %d", player.Gold)
	}
	if len(invRepo.purchaseCalls) != 1 {
		t.Fatalf("expected one atomic purchase call, got %d", len(invRepo.purchaseCalls))
	}
	call := invRepo.purchaseCalls[0]
	if call.playerID != int64(player.ID) || call.itemID != 7 || call.amount != 2 || call.price != 300 {
		t.Fatalf("unexpected atomic purchase call: %+v", call)
	}
	if len(playerRepo.updatedPlayers) != 0 {
		t.Fatalf("expected no direct player update")
	}
	if !strings.Contains(resp.Text, "Obrigada") {
		t.Fatalf("expected success message, got %q", resp.Text)
	}
	if !strings.Contains(resp.Text, "Item: Pocao") {
		t.Fatalf("expected item name in success message, got %q", resp.Text)
	}
	if !strings.Contains(resp.Text, "Requisitos: Nivel 15+ | Classes: Wind Striker,Earth Sentinel") {
		t.Fatalf("expected requirements text in success message, got %q", resp.Text)
	}
	if !strings.Contains(resp.Text, "ATK +3") || !strings.Contains(resp.Text, "VEL +1") || !strings.Contains(resp.Text, "ENERGIA +2") {
		t.Fatalf("expected positive bonus text in success message, got %q", resp.Text)
	}
	if strings.Contains(resp.Text, "DEF +0") || strings.Contains(resp.Text, "HP +0") {
		t.Fatalf("expected only bonuses > 0 in success message, got %q", resp.Text)
	}
}

func TestShopService_BuyItem_EquipmentUsesCatalogPrice(t *testing.T) {
	player := basePlayer()
	player.ID = 30
	player.Gold = 1600

	invRepo := &fakeInventoryRepo{}
	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	itemRepo := &fakeItemsRepo{
		getItemByIDFn: func(id int) (*domain.Item, error) {
			return &domain.Item{ID: 7, Name: "Arma", Price: 300, Amount: 1, Type: domain.ItemEquipment}, nil
		},
	}
	svc := &ShopService{
		PlayerRepository:    playerRepo,
		ItemRepository:      itemRepo,
		InventoryRepository: invRepo,
	}

	_, err := svc.BuyItem(1, map[string]string{"itemID": "7"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(invRepo.purchaseCalls) != 1 {
		t.Fatalf("expected one atomic purchase call, got %d", len(invRepo.purchaseCalls))
	}
	if invRepo.purchaseCalls[0].price != 300 {
		t.Fatalf("expected catalog price 300, got %d", invRepo.purchaseCalls[0].price)
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
			purchaseFn: func(playerID int64, itemID uint, amount int, price int) error {
				return errors.New("inventory fail")
			},
		},
	}

	_, err := svc.BuyItem(1, map[string]string{"itemID": "1"})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestShopService_BuyItem_PurchaseError(t *testing.T) {
	player := basePlayer()
	player.ID = 8
	player.Gold = 1000

	svc := &ShopService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
		ItemRepository: &fakeItemsRepo{
			getItemByIDFn: func(id int) (*domain.Item, error) {
				return &domain.Item{ID: 2, Price: 200, Amount: 1}, nil
			},
		},
		InventoryRepository: &fakeInventoryRepo{
			purchaseFn: func(playerID int64, itemID uint, amount int, price int) error {
				return errors.New("purchase fail")
			},
		},
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

func TestShopService_ListSellEquipments_Success(t *testing.T) {
	player := basePlayer()
	player.ID = 50
	svc := &ShopService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
		InventoryRepository: &fakeInventoryRepo{
			listByTypeFn: func(playerID int64, itemType domain.ItemType, page int, pageSize int) ([]domain.InventoryItem, int64, error) {
				if playerID != int64(player.ID) {
					t.Fatalf("unexpected playerID %d", playerID)
				}
				if itemType != domain.ItemEquipment {
					t.Fatalf("expected equipment type, got %s", itemType)
				}
				return []domain.InventoryItem{
					{ID: 1, Quantity: 1, Item: domain.Item{Name: "Espada", Price: 500, Type: domain.ItemEquipment}},
				}, 1, nil
			},
		},
	}

	resp, err := svc.ListSellEquipments(99, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "SHOP_SELL_ITEMS_MENU" {
		t.Fatalf("expected SHOP_SELL_ITEMS_MENU, got %q", resp.MenuKey)
	}
}

func TestShopService_SellItem_Success(t *testing.T) {
	player := basePlayer()
	player.ID = 51
	player.Gold = 100

	invItem := &domain.InventoryItem{
		ID:       10,
		ItemID:   200,
		Quantity: 1,
		Item: domain.Item{
			ID:    200,
			Name:  "Lanca",
			Type:  domain.ItemEquipment,
			Price: 500,
		},
	}

	invRepo := &fakeInventoryRepo{
		getByUserFn: func(playerID int64, inventoryItemID uint) (*domain.InventoryItem, error) {
			return invItem, nil
		},
		listByTypeFn: func(playerID int64, itemType domain.ItemType, page int, pageSize int) ([]domain.InventoryItem, int64, error) {
			return []domain.InventoryItem{}, 0, nil
		},
	}
	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}

	svc := &ShopService{
		PlayerRepository:    playerRepo,
		InventoryRepository: invRepo,
	}

	resp, err := svc.SellItem(99, map[string]string{"inventoryItemID": "10", "page": "1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if invItem.Quantity != 1 {
		t.Fatalf("expected no direct quantity mutation, got %d", invItem.Quantity)
	}
	if player.Gold != 100 {
		t.Fatalf("expected no direct gold mutation, got %d", player.Gold)
	}
	if len(invRepo.sellCalls) != 1 {
		t.Fatalf("expected one atomic sell call, got %d", len(invRepo.sellCalls))
	}
	if !strings.Contains(resp.Text, "vendido por 200 Gold") {
		t.Fatalf("expected sell confirmation text, got %q", resp.Text)
	}
}
