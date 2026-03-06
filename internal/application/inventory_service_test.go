package application

import (
	"errors"
	"strings"
	"testing"

	"github.com/vigmiranda/dragonki/internal/domain"
)

type fakeItemsRepo struct {
	getRandomListFn           func(itemType domain.ItemType, rarity domain.ItemRarity) ([]domain.Item, error)
	getRandomBySlotAndLimitFn func(slot domain.EquipmentSlot, limit int) ([]domain.Item, error)
	getItemByIDFn             func(id int) (*domain.Item, error)
}

func (f *fakeItemsRepo) GetRandomList(itemType domain.ItemType, rarity domain.ItemRarity) ([]domain.Item, error) {
	if f.getRandomListFn == nil {
		return nil, nil
	}
	return f.getRandomListFn(itemType, rarity)
}

func (f *fakeItemsRepo) GetRandomEquipmentsBySlot(slot domain.EquipmentSlot, limit int) ([]domain.Item, error) {
	if f.getRandomBySlotAndLimitFn == nil {
		return nil, nil
	}
	return f.getRandomBySlotAndLimitFn(slot, limit)
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
	if resp.MenuKey != "INVENTORY_ITEMS_MENU" {
		t.Fatalf("expected INVENTORY_ITEMS_MENU, got %q", resp.MenuKey)
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

func TestInventoryService_ListItemsByUserAndType_RepoError(t *testing.T) {
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

	_, err := svc.ListItemsByUserAndType(1, domain.ItemEquipment)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestInventoryService_OpenEquipmentTypes(t *testing.T) {
	player := basePlayer()
	invRepo := &fakeInventoryRepo{
		countEquipmentBySlotFn: func(playerID int64) (map[domain.EquipmentSlot]int64, error) {
			return map[domain.EquipmentSlot]int64{
				domain.SlotWeapon: 2,
				domain.SlotArmor:  1,
				domain.SlotGloves: 5,
				domain.SlotHelmet: 0,
				domain.SlotBoots:  3,
				domain.SlotRing:   4,
				domain.SlotAmulet: 6,
			}, nil
		},
	}
	svc := InventoryService{
		InventoryRepository: invRepo,
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
	}

	resp, err := svc.OpenEquipmentTypes(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "INVENTORY_EQUIPMENT_TYPES_MENU" {
		t.Fatalf("expected INVENTORY_EQUIPMENT_TYPES_MENU, got %q", resp.MenuKey)
	}
	if resp.MenuParams["weaponCount"] != int64(2) || resp.MenuParams["ringCount"] != int64(4) || resp.MenuParams["amuletCount"] != int64(6) {
		t.Fatalf("unexpected counts: %+v", resp.MenuParams)
	}
}

func TestInventoryService_ListEquipmentItemsBySlot(t *testing.T) {
	player := basePlayer()
	player.ID = 77
	inventoryItems := []domain.InventoryItem{
		{ID: 1, Quantity: 1, Item: domain.Item{Name: "Espada"}},
	}
	invRepo := &fakeInventoryRepo{
		listEquipmentBySlotFn: func(playerID int64, slot domain.EquipmentSlot, page int, pageSize int) ([]domain.InventoryItem, int64, error) {
			if playerID != 77 {
				t.Fatalf("expected playerID 77, got %d", playerID)
			}
			if slot != domain.SlotWeapon {
				t.Fatalf("expected slot WEAPON, got %s", slot)
			}
			if page != 1 || pageSize != 5 {
				t.Fatalf("expected pagination 1/5, got %d/%d", page, pageSize)
			}
			return inventoryItems, 1, nil
		},
	}
	svc := InventoryService{
		InventoryRepository: invRepo,
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
	}

	resp, err := svc.ListEquipmentItemsBySlot(1, "WEAPON", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "INVENTORY_ITEMS_MENU" {
		t.Fatalf("expected INVENTORY_ITEMS_MENU, got %q", resp.MenuKey)
	}
	if resp.MenuParams["paginationAction"] != "open_inventory_equipaments_type?slot=WEAPON" {
		t.Fatalf("unexpected pagination action: %v", resp.MenuParams["paginationAction"])
	}
	if resp.MenuParams["backAction"] != "open_inventory_equipaments" {
		t.Fatalf("unexpected back action: %v", resp.MenuParams["backAction"])
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

func TestInventoryService_UseItem_NilPlayer(t *testing.T) {
	svc := InventoryService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) {
				return nil, nil
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

func TestInventoryService_UseItem_NilInventoryItem(t *testing.T) {
	player := basePlayer()
	svc := InventoryService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
		InventoryRepository: &fakeInventoryRepo{
			getByUserFn: func(playerID int64, inventoryItemID uint) (*domain.InventoryItem, error) {
				return nil, nil
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

func TestInventoryService_EquipItem_LevelInsufficient(t *testing.T) {
	player := basePlayer()
	player.Level = 4
	player.Class.Name = "Wind Striker"

	inv := &domain.InventoryItem{
		ID:       1,
		Quantity: 1,
		Item: domain.Item{
			ID:            200,
			Name:          "Espada Flamejante",
			Type:          domain.ItemEquipment,
			Slot:          "WEAPON",
			RequiredLevel: 10,
		},
	}

	svc := InventoryService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
		InventoryRepository: &fakeInventoryRepo{
			getByUserFn: func(playerID int64, inventoryItemID uint) (*domain.InventoryItem, error) { return inv, nil },
		},
	}

	resp, err := svc.UseItem(1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(strings.ToLower(resp.Text), "nivel 10") {
		t.Fatalf("expected level restriction message, got %q", resp.Text)
	}
}

func TestInventoryService_EquipItem_ClassRestriction(t *testing.T) {
	player := basePlayer()
	player.Level = 20
	player.Class.Name = "Earth Sentinel"

	inv := &domain.InventoryItem{
		ID:       1,
		Quantity: 1,
		Item: domain.Item{
			ID:             201,
			Name:           "Lâmina dos Ventos",
			Type:           domain.ItemEquipment,
			Slot:           "WEAPON",
			RequiredLevel:  5,
			AllowedClasses: "Wind Striker,Ki Vanguard",
		},
	}

	svc := InventoryService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
		InventoryRepository: &fakeInventoryRepo{
			getByUserFn: func(playerID int64, inventoryItemID uint) (*domain.InventoryItem, error) { return inv, nil },
		},
	}

	resp, err := svc.UseItem(1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(strings.ToLower(resp.Text), "nao pode equipar") {
		t.Fatalf("expected class restriction message, got %q", resp.Text)
	}
}

func TestInventoryService_EquipItem_ReplacesSlotAndRecalculatesStats(t *testing.T) {
	player := basePlayer()
	player.ID = 99
	player.Level = 30
	player.Class.Name = "Wind Striker"
	player.ATK = 100
	player.DEF = 40
	player.Speed = 20
	player.MaxHP = 200
	player.MaxEnergy = 25

	newItem := domain.Item{
		ID:             300,
		Name:           "Espada Flamejante",
		Type:           domain.ItemEquipment,
		Slot:           "WEAPON",
		RequiredLevel:  10,
		AllowedClasses: "Wind Striker",
		AtkBonus:       15,
		DefBonus:       8,
		SpeedBonus:     4,
		HPBonus:        10,
	}
	inv := &domain.InventoryItem{
		ID:       1,
		ItemID:   300,
		Quantity: 1,
		Item:     newItem,
	}
	oldInv := &domain.InventoryItem{
		ID:       250,
		ItemID:   251,
		Quantity: 1,
		Item: domain.Item{
			ID:         251,
			AtkBonus:   5,
			DefBonus:   2,
			SpeedBonus: 1,
			HPBonus:    5,
		},
	}

	oldID := uint(250)
	updatedPlayer := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	invRepo := &fakeInventoryRepo{
		getByUserFn: func(playerID int64, inventoryItemID uint) (*domain.InventoryItem, error) {
			if inventoryItemID == oldID {
				return oldInv, nil
			}
			return inv, nil
		},
		getOrCreateEquipFn: func(playerID int64) (*domain.PlayerEquipment, error) {
			return &domain.PlayerEquipment{PlayerID: playerID, WeaponInventoryItemID: &oldID}, nil
		},
	}

	svc := InventoryService{
		PlayerRepository:    updatedPlayer,
		InventoryRepository: invRepo,
		ItemsRepository:     &fakeItemsRepo{},
	}

	resp, err := svc.UseItem(1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 100-5+15, 40-2+8, 20-1+4, 200-5+10
	if player.ATK != 110 || player.DEF != 46 || player.Speed != 23 || player.MaxHP != 205 {
		t.Fatalf("unexpected recalculated stats: atk=%d def=%d spd=%d hp=%d", player.ATK, player.DEF, player.Speed, player.MaxHP)
	}
	if len(invRepo.updateEquipCalls) != 1 {
		t.Fatalf("expected equipment update call")
	}
	if invRepo.updateEquipCalls[0].WeaponInventoryItemID == nil || *invRepo.updateEquipCalls[0].WeaponInventoryItemID != inv.ID {
		t.Fatalf("expected weapon slot updated to new inventory item id")
	}
	if len(updatedPlayer.updatedPlayers) != 1 {
		t.Fatalf("expected player update call")
	}
	if !strings.Contains(resp.Text, "Espada Flamejante equipada") {
		t.Fatalf("expected equip success text, got %q", resp.Text)
	}
}

func TestInventoryService_EquipItem_IgnoresZeroEquippedID(t *testing.T) {
	player := basePlayer()
	player.ID = 55
	player.Level = 10
	player.Class.Name = "Wind Striker"

	inv := &domain.InventoryItem{
		ID:       10,
		ItemID:   60,
		Quantity: 1,
		Item: domain.Item{
			ID:             60,
			Name:           "Lamina do Iniciante",
			Type:           domain.ItemEquipment,
			Slot:           "WEAPON",
			RequiredLevel:  1,
			AllowedClasses: "Wind Striker",
			AtkBonus:       7,
		},
	}

	zero := uint(0)
	invRepo := &fakeInventoryRepo{
		getByUserFn: func(playerID int64, inventoryItemID uint) (*domain.InventoryItem, error) { return inv, nil },
		getOrCreateEquipFn: func(playerID int64) (*domain.PlayerEquipment, error) {
			return &domain.PlayerEquipment{PlayerID: playerID, WeaponInventoryItemID: &zero}, nil
		},
	}
	svc := InventoryService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
		InventoryRepository: invRepo,
		ItemsRepository: &fakeItemsRepo{
			getItemByIDFn: func(id int) (*domain.Item, error) {
				t.Fatalf("should not query old item for id=0")
				return nil, nil
			},
		},
	}

	_, err := svc.UseItem(1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInventoryService_EquipItem_LoadsItemWhenPreloadMissing(t *testing.T) {
	player := basePlayer()
	player.ID = 56
	player.Level = 10
	player.Class.Name = "Wind Striker"

	inv := &domain.InventoryItem{
		ID:       11,
		ItemID:   60,
		Quantity: 1,
		Item:     domain.Item{}, // missing preload
	}

	invRepo := &fakeInventoryRepo{
		getByUserFn: func(playerID int64, inventoryItemID uint) (*domain.InventoryItem, error) { return inv, nil },
		getOrCreateEquipFn: func(playerID int64) (*domain.PlayerEquipment, error) {
			return &domain.PlayerEquipment{PlayerID: playerID}, nil
		},
	}
	svc := InventoryService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
		InventoryRepository: invRepo,
		ItemsRepository: &fakeItemsRepo{
			getItemByIDFn: func(id int) (*domain.Item, error) {
				if id != 60 {
					t.Fatalf("expected id 60, got %d", id)
				}
				return &domain.Item{
					ID:             60,
					Name:           "Lamina",
					Type:           domain.ItemEquipment,
					Slot:           "WEAPON",
					RequiredLevel:  1,
					AllowedClasses: "Wind Striker",
					AtkBonus:       7,
				}, nil
			},
		},
	}

	_, err := svc.UseItem(1, 11)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInventoryService_ViewEquippedItems(t *testing.T) {
	player := basePlayer()
	player.ID = 100
	weaponID := uint(10)
	armorID := uint(11)
	gloveID := uint(12)
	ringID := uint(13)
	amuletID := uint(14)

	svc := InventoryService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
		InventoryRepository: &fakeInventoryRepo{
			getByUserFn: func(playerID int64, inventoryItemID uint) (*domain.InventoryItem, error) {
				switch inventoryItemID {
				case 10:
					return &domain.InventoryItem{
						ID:       10,
						ItemID:   300,
						Quantity: 1,
						Item:     domain.Item{ID: 300, Name: "Espada Flamejante", AtkBonus: 15, SpeedBonus: 4},
					}, nil
				case 11:
					return &domain.InventoryItem{
						ID:       11,
						ItemID:   301,
						Quantity: 1,
						Item:     domain.Item{ID: 301, Name: "Armadura de Couro", DefBonus: 8, HPBonus: 20},
					}, nil
				case 12:
					return &domain.InventoryItem{
						ID:       12,
						ItemID:   302,
						Quantity: 1,
						Item:     domain.Item{ID: 302, Name: "Luva de Aco", AtkBonus: 2, DefBonus: 1},
					}, nil
				case 13:
					return &domain.InventoryItem{
						ID:       13,
						ItemID:   303,
						Quantity: 1,
						Item:     domain.Item{ID: 303, Name: "Anel de Ki", AtkBonus: 3, EnergyBonus: 2},
					}, nil
				case 14:
					return &domain.InventoryItem{
						ID:       14,
						ItemID:   304,
						Quantity: 1,
						Item:     domain.Item{ID: 304, Name: "Amuleto Lunar", DefBonus: 5, HPBonus: 10},
					}, nil
				default:
					return nil, nil
				}
			},
			getOrCreateEquipFn: func(playerID int64) (*domain.PlayerEquipment, error) {
				return &domain.PlayerEquipment{
					PlayerID:              playerID,
					WeaponInventoryItemID: &weaponID,
					ArmorInventoryItemID:  &armorID,
					GlovesInventoryItemID: &gloveID,
					RingInventoryItemID:   &ringID,
					AmuletInventoryItemID: &amuletID,
				}, nil
			},
		},
		ItemsRepository: &fakeItemsRepo{},
	}

	resp, err := svc.ViewEquippedItems(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "INVENTORY_EQUIPPED_MENU" {
		t.Fatalf("expected INVENTORY_EQUIPPED_MENU, got %q", resp.MenuKey)
	}
	if !strings.Contains(resp.Text, "Arma: Espada Flamejante | ATK +15 | Speed +4") {
		t.Fatalf("expected weapon line in text, got %q", resp.Text)
	}
	if !strings.Contains(resp.Text, "Armadura:\n  Armadura: Armadura de Couro | DEF +8 | MaxHP +20\n  Luva: Luva de Aco | ATK +2 | DEF +1") {
		t.Fatalf("expected armor group in text, got %q", resp.Text)
	}
	if !strings.Contains(resp.Text, "Acessorio:\n  Anel: Anel de Ki | ATK +3 | Energy +2\n  Amuleto: Amuleto Lunar | DEF +5 | MaxHP +10") {
		t.Fatalf("expected accessory group in text, got %q", resp.Text)
	}
}

func TestInventoryService_ViewEquippedItems_AllSlotsFree(t *testing.T) {
	player := basePlayer()
	player.ID = 101

	svc := InventoryService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
		InventoryRepository: &fakeInventoryRepo{
			getOrCreateEquipFn: func(playerID int64) (*domain.PlayerEquipment, error) {
				return &domain.PlayerEquipment{PlayerID: playerID}, nil
			},
		},
		ItemsRepository: &fakeItemsRepo{},
	}

	resp, err := svc.ViewEquippedItems(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Text, "Arma: vazio") {
		t.Fatalf("expected empty weapon text, got %q", resp.Text)
	}
	if !strings.Contains(resp.Text, "Armadura:\n  Armadura: vazio\n  Luva: vazio") {
		t.Fatalf("expected empty armor group text, got %q", resp.Text)
	}
	if !strings.Contains(resp.Text, "Acessorio:\n  Anel: vazio\n  Amuleto: vazio") {
		t.Fatalf("expected empty accessory group text, got %q", resp.Text)
	}
}

func TestInventoryService_UnequipItem_Success(t *testing.T) {
	player := basePlayer()
	player.ID = 200
	player.Level = 20
	player.ATK = 115
	player.DEF = 55
	player.Speed = 24
	player.MaxHP = 210
	player.MaxEnergy = 27

	weaponInvID := uint(88)
	pe := &domain.PlayerEquipment{
		PlayerID:              int64(player.ID),
		WeaponInventoryItemID: &weaponInvID,
	}
	inv := &domain.InventoryItem{
		ID:       weaponInvID,
		ItemID:   501,
		Quantity: 1,
		Item: domain.Item{
			ID:          501,
			Name:        "Espada Flamejante",
			Type:        domain.ItemEquipment,
			Slot:        "WEAPON",
			AtkBonus:    15,
			DefBonus:    5,
			SpeedBonus:  4,
			HPBonus:     10,
			EnergyBonus: 2,
		},
	}

	invRepo := &fakeInventoryRepo{
		getOrCreateEquipFn: func(playerID int64) (*domain.PlayerEquipment, error) { return pe, nil },
		getByUserFn: func(playerID int64, inventoryItemID uint) (*domain.InventoryItem, error) {
			if inventoryItemID == weaponInvID {
				return inv, nil
			}
			return nil, nil
		},
	}
	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	svc := InventoryService{
		PlayerRepository:    playerRepo,
		InventoryRepository: invRepo,
		ItemsRepository:     &fakeItemsRepo{},
	}

	resp, err := svc.UnequipItem(1, "WEAPON")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pe.WeaponInventoryItemID != nil {
		t.Fatalf("expected weapon slot unequipped")
	}
	if player.ATK != 100 || player.DEF != 50 || player.Speed != 20 || player.MaxHP != 200 || player.MaxEnergy != 25 {
		t.Fatalf("unexpected stats after unequip: atk=%d def=%d spd=%d hp=%d en=%d", player.ATK, player.DEF, player.Speed, player.MaxHP, player.MaxEnergy)
	}
	if len(invRepo.updateEquipCalls) != 1 {
		t.Fatalf("expected one equipment update")
	}
	if len(playerRepo.updatedPlayers) != 1 {
		t.Fatalf("expected one player update")
	}
	if !strings.Contains(resp.Text, "desequipado com sucesso") {
		t.Fatalf("expected success message, got %q", resp.Text)
	}
	if resp.MenuKey != "INVENTORY_EQUIPPED_MENU" {
		t.Fatalf("expected INVENTORY_EQUIPPED_MENU, got %q", resp.MenuKey)
	}
}

func TestInventoryService_UnequipItem_EmptySlot(t *testing.T) {
	player := basePlayer()
	player.ID = 201

	invRepo := &fakeInventoryRepo{
		getOrCreateEquipFn: func(playerID int64) (*domain.PlayerEquipment, error) {
			return &domain.PlayerEquipment{PlayerID: playerID}, nil
		},
	}
	svc := InventoryService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
		InventoryRepository: invRepo,
		ItemsRepository:     &fakeItemsRepo{},
	}

	resp, err := svc.UnequipItem(1, "WEAPON")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Text, "Nenhum item equipado neste slot") {
		t.Fatalf("expected empty-slot message, got %q", resp.Text)
	}
	if resp.MenuKey != "INVENTORY_EQUIPPED_MENU" {
		t.Fatalf("expected INVENTORY_EQUIPPED_MENU, got %q", resp.MenuKey)
	}
}

func TestInventoryService_UnequipAll_Success(t *testing.T) {
	player := basePlayer()
	player.ID = 300
	player.ATK = 120
	player.DEF = 60
	player.Speed = 30
	player.MaxHP = 220
	player.MaxEnergy = 28

	weaponID := uint(91)
	ringID := uint(92)
	pe := &domain.PlayerEquipment{
		PlayerID:              int64(player.ID),
		WeaponInventoryItemID: &weaponID,
		RingInventoryItemID:   &ringID,
	}

	invRepo := &fakeInventoryRepo{
		getOrCreateEquipFn: func(playerID int64) (*domain.PlayerEquipment, error) { return pe, nil },
		getByUserFn: func(playerID int64, inventoryItemID uint) (*domain.InventoryItem, error) {
			switch inventoryItemID {
			case weaponID:
				return &domain.InventoryItem{
					ID:       weaponID,
					ItemID:   601,
					Quantity: 1,
					Item: domain.Item{
						ID:          601,
						Name:        "Espada",
						Type:        domain.ItemEquipment,
						AtkBonus:    15,
						DefBonus:    5,
						SpeedBonus:  4,
						HPBonus:     10,
						EnergyBonus: 2,
					},
				}, nil
			case ringID:
				return &domain.InventoryItem{
					ID:       ringID,
					ItemID:   602,
					Quantity: 1,
					Item: domain.Item{
						ID:         602,
						Name:       "Anel",
						Type:       domain.ItemEquipment,
						AtkBonus:   5,
						DefBonus:   5,
						SpeedBonus: 1,
						HPBonus:    10,
					},
				}, nil
			default:
				return nil, nil
			}
		},
	}
	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	svc := InventoryService{
		PlayerRepository:    playerRepo,
		InventoryRepository: invRepo,
		ItemsRepository:     &fakeItemsRepo{},
	}

	resp, err := svc.UnequipAll(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pe.WeaponInventoryItemID != nil || pe.RingInventoryItemID != nil {
		t.Fatalf("expected all slots unequipped")
	}
	if player.ATK != 100 || player.DEF != 50 || player.Speed != 25 || player.MaxHP != 200 || player.MaxEnergy != 26 {
		t.Fatalf("unexpected stats after unequip all: atk=%d def=%d spd=%d hp=%d en=%d", player.ATK, player.DEF, player.Speed, player.MaxHP, player.MaxEnergy)
	}
	if len(invRepo.updateEquipCalls) != 1 {
		t.Fatalf("expected one equipment update")
	}
	if len(playerRepo.updatedPlayers) != 1 {
		t.Fatalf("expected one player update")
	}
	if !strings.Contains(resp.Text, "Todos os equipamentos foram desequipados") {
		t.Fatalf("expected unequip-all message, got %q", resp.Text)
	}
}

func TestInventoryService_UnequipAll_NoItems(t *testing.T) {
	player := basePlayer()
	player.ID = 301

	invRepo := &fakeInventoryRepo{
		getOrCreateEquipFn: func(playerID int64) (*domain.PlayerEquipment, error) {
			return &domain.PlayerEquipment{PlayerID: playerID}, nil
		},
	}
	svc := InventoryService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
		InventoryRepository: invRepo,
		ItemsRepository:     &fakeItemsRepo{},
	}

	resp, err := svc.UnequipAll(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Text, "nao tem equipamentos para desequipar") {
		t.Fatalf("expected no-items message, got %q", resp.Text)
	}
}
