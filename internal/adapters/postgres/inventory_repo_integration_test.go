package postgres

import (
	"testing"

	"github.com/vigmiranda/dragonki/internal/domain"
)

func TestInventoryPostgresIntegration_ConsumablesAggregateAndEquipmentRespectsEquippedSlots(t *testing.T) {
	testDB := newIntegrationTestDB(t)

	playerRepo := NewPlayerPostgresRepo(testDB.DB)
	inventoryRepo := NewInventoryPostgresRepo(testDB.DB)

	consumableItem := domain.Item{
		Code:        "IT_POTION_SMALL",
		Name:        "Integration Potion",
		Description: "Potion used by integration tests.",
		Type:        domain.ItemConsumable,
		Rarity:      "COMMON",
		Amount:      1,
		Price:       100,
	}
	if err := testDB.DB.Create(&consumableItem).Error; err != nil {
		t.Fatalf("create consumable item fixture: %v", err)
	}

	weaponItem := domain.Item{
		Code:        "IT_WEAPON_SMALL",
		Name:        "Integration Sword",
		Description: "Weapon used by integration tests.",
		Type:        domain.ItemEquipment,
		Rarity:      "COMMON",
		Amount:      1,
		Price:       500,
	}
	if err := testDB.DB.Create(&weaponItem).Error; err != nil {
		t.Fatalf("create weapon item fixture: %v", err)
	}

	weaponEquipment := domain.Equipment{
		ItemID:         weaponItem.ID,
		Slot:           string(domain.SlotWeapon),
		RequiredLevel:  1,
		AtkBonus:       5,
		AllowedClasses: "Wind Striker,Earth Sentinel,Ki Vanguard",
	}
	if err := testDB.DB.Create(&weaponEquipment).Error; err != nil {
		t.Fatalf("create weapon equipment fixture: %v", err)
	}

	player, err := playerRepo.GetOrCreate(710001)
	if err != nil {
		t.Fatalf("create player: %v", err)
	}

	playerID := int64(player.ID)

	if err := inventoryRepo.CreateOrUpdateInventory(playerID, consumableItem.ID, 2); err != nil {
		t.Fatalf("add consumable inventory first batch: %v", err)
	}
	if err := inventoryRepo.CreateOrUpdateInventory(playerID, consumableItem.ID, 3); err != nil {
		t.Fatalf("add consumable inventory second batch: %v", err)
	}

	consumables, totalConsumables, err := inventoryRepo.ListItemsByUserAndType(playerID, domain.ItemConsumable, 1, 10)
	if err != nil {
		t.Fatalf("list consumables: %v", err)
	}

	if totalConsumables != 1 {
		t.Fatalf("unexpected consumable total: got %d want 1", totalConsumables)
	}
	if len(consumables) != 1 {
		t.Fatalf("unexpected consumable list size: got %d want 1", len(consumables))
	}
	if consumables[0].Quantity != 5 {
		t.Fatalf("unexpected consumable quantity: got %d want 5", consumables[0].Quantity)
	}
	if consumables[0].ItemID != consumableItem.ID {
		t.Fatalf("unexpected consumable item id: got %d want %d", consumables[0].ItemID, consumableItem.ID)
	}

	if err := inventoryRepo.CreateOrUpdateInventory(playerID, weaponItem.ID, 2); err != nil {
		t.Fatalf("add equipment inventory: %v", err)
	}

	weapons, totalWeapons, err := inventoryRepo.ListEquipmentItemsByUserAndSlot(playerID, domain.SlotWeapon, 1, 10)
	if err != nil {
		t.Fatalf("list weapon equipments before equip: %v", err)
	}

	if totalWeapons != 2 || len(weapons) != 2 {
		t.Fatalf("unexpected weapon totals before equip: total=%d len=%d", totalWeapons, len(weapons))
	}
	for _, item := range weapons {
		if item.Quantity != 1 {
			t.Fatalf("expected equipment quantity 1 per row, got %d", item.Quantity)
		}
		if item.Item.Slot != string(domain.SlotWeapon) {
			t.Fatalf("expected hydrated weapon slot, got %q", item.Item.Slot)
		}
	}

	selectedWeapon, err := inventoryRepo.GetInventoryItemByUser(playerID, weapons[0].ID)
	if err != nil {
		t.Fatalf("get inventory item by user: %v", err)
	}

	if selectedWeapon.Item.Slot != string(domain.SlotWeapon) {
		t.Fatalf("expected selected weapon hydrated slot, got %q", selectedWeapon.Item.Slot)
	}

	playerEquipment, err := inventoryRepo.GetOrCreatePlayerEquipment(playerID)
	if err != nil {
		t.Fatalf("get or create player equipment: %v", err)
	}

	playerEquipment.WeaponInventoryItemID = &selectedWeapon.ID
	if err := inventoryRepo.UpdatePlayerEquipment(playerEquipment); err != nil {
		t.Fatalf("update player equipment: %v", err)
	}

	weapons, totalWeapons, err = inventoryRepo.ListEquipmentItemsByUserAndSlot(playerID, domain.SlotWeapon, 1, 10)
	if err != nil {
		t.Fatalf("list weapon equipments after equip: %v", err)
	}

	if totalWeapons != 1 || len(weapons) != 1 {
		t.Fatalf("unexpected weapon totals after equip: total=%d len=%d", totalWeapons, len(weapons))
	}
	if weapons[0].ID == selectedWeapon.ID {
		t.Fatalf("equipped inventory item should be excluded from list")
	}

	slotCounts, err := inventoryRepo.CountEquipmentItemsBySlot(playerID)
	if err != nil {
		t.Fatalf("count equipment items by slot: %v", err)
	}

	if slotCounts[domain.SlotWeapon] != 1 {
		t.Fatalf("unexpected available weapon count after equip: got %d want 1", slotCounts[domain.SlotWeapon])
	}
	if slotCounts[domain.SlotArmor] != 0 {
		t.Fatalf("unexpected armor count: got %d want 0", slotCounts[domain.SlotArmor])
	}
}
