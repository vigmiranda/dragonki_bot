package postgres

import (
	"errors"
	"testing"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
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

func TestInventoryPostgresIntegration_PurchaseAndSellAreAtomic(t *testing.T) {
	testDB := newIntegrationTestDB(t)

	playerRepo := NewPlayerPostgresRepo(testDB.DB)
	inventoryRepo := NewInventoryPostgresRepo(testDB.DB)

	item := domain.Item{
		Code:        "IT_SHOP_SWORD",
		Name:        "Atomic Sword",
		Description: "Equipment sold in integration tests.",
		Type:        domain.ItemEquipment,
		Rarity:      "COMMON",
		Amount:      1,
		Price:       500,
	}
	if err := testDB.DB.Create(&item).Error; err != nil {
		t.Fatalf("create shop item fixture: %v", err)
	}

	equipment := domain.Equipment{
		ItemID:         item.ID,
		Slot:           string(domain.SlotWeapon),
		RequiredLevel:  1,
		AtkBonus:       7,
		AllowedClasses: "Wind Striker,Earth Sentinel,Ki Vanguard",
	}
	if err := testDB.DB.Create(&equipment).Error; err != nil {
		t.Fatalf("create equipment fixture: %v", err)
	}

	player, err := playerRepo.GetOrCreate(710002)
	if err != nil {
		t.Fatalf("create player: %v", err)
	}

	player.Gold = 1000
	if err := playerRepo.Update(player); err != nil {
		t.Fatalf("persist player gold: %v", err)
	}

	if err := inventoryRepo.PurchaseItem(int64(player.ID), item.ID, 1, 500); err != nil {
		t.Fatalf("purchase item atomically: %v", err)
	}

	reloadedPlayer, err := playerRepo.GetOrCreate(710002)
	if err != nil {
		t.Fatalf("reload player after purchase: %v", err)
	}
	if reloadedPlayer.Gold != 500 {
		t.Fatalf("expected gold debited after purchase, got %d", reloadedPlayer.Gold)
	}

	items, total, err := inventoryRepo.ListItemsByUserAndType(int64(player.ID), domain.ItemEquipment, 1, 10)
	if err != nil {
		t.Fatalf("list purchased equipment: %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("expected one purchased equipment row, got total=%d len=%d", total, len(items))
	}

	soldItem, err := inventoryRepo.SellInventoryItem(int64(player.ID), items[0].ID, 200)
	if err != nil {
		t.Fatalf("sell item atomically: %v", err)
	}
	if soldItem == nil || soldItem.Quantity != 0 {
		t.Fatalf("expected sold item quantity 0, got %+v", soldItem)
	}

	reloadedPlayer, err = playerRepo.GetOrCreate(710002)
	if err != nil {
		t.Fatalf("reload player after sell: %v", err)
	}
	if reloadedPlayer.Gold != 700 {
		t.Fatalf("expected gold credited after sell, got %d", reloadedPlayer.Gold)
	}
}

func TestInventoryPostgresIntegration_SellEquippedItemIsRejectedWithoutCreditingGold(t *testing.T) {
	testDB := newIntegrationTestDB(t)

	playerRepo := NewPlayerPostgresRepo(testDB.DB)
	inventoryRepo := NewInventoryPostgresRepo(testDB.DB)

	item := domain.Item{
		Code:        "IT_EQUIPPED_SWORD",
		Name:        "Equipped Sword",
		Description: "Equipment protected during integration tests.",
		Type:        domain.ItemEquipment,
		Rarity:      "COMMON",
		Amount:      1,
		Price:       300,
	}
	if err := testDB.DB.Create(&item).Error; err != nil {
		t.Fatalf("create equipped item fixture: %v", err)
	}

	equipment := domain.Equipment{
		ItemID:         item.ID,
		Slot:           string(domain.SlotWeapon),
		RequiredLevel:  1,
		AtkBonus:       5,
		AllowedClasses: "Wind Striker,Earth Sentinel,Ki Vanguard",
	}
	if err := testDB.DB.Create(&equipment).Error; err != nil {
		t.Fatalf("create equipped equipment fixture: %v", err)
	}

	player, err := playerRepo.GetOrCreate(710003)
	if err != nil {
		t.Fatalf("create player: %v", err)
	}

	player.Gold = 400
	if err := playerRepo.Update(player); err != nil {
		t.Fatalf("persist player gold: %v", err)
	}

	if err := inventoryRepo.CreateOrUpdateInventory(int64(player.ID), item.ID, 1); err != nil {
		t.Fatalf("seed inventory item: %v", err)
	}

	items, _, err := inventoryRepo.ListItemsByUserAndType(int64(player.ID), domain.ItemEquipment, 1, 10)
	if err != nil {
		t.Fatalf("list equipment items: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected one equipment item, got %d", len(items))
	}

	playerEquipment, err := inventoryRepo.GetOrCreatePlayerEquipment(int64(player.ID))
	if err != nil {
		t.Fatalf("get or create player equipment: %v", err)
	}
	playerEquipment.WeaponInventoryItemID = &items[0].ID
	if err := inventoryRepo.UpdatePlayerEquipment(playerEquipment); err != nil {
		t.Fatalf("equip inventory item: %v", err)
	}

	soldItem, err := inventoryRepo.SellInventoryItem(int64(player.ID), items[0].ID, 120)
	if !errors.Is(err, domain.ErrInventoryItemEquipped) {
		t.Fatalf("expected equipped item error, got soldItem=%+v err=%v", soldItem, err)
	}

	reloadedPlayer, err := playerRepo.GetOrCreate(710003)
	if err != nil {
		t.Fatalf("reload player after rejected sell: %v", err)
	}
	if reloadedPlayer.Gold != 400 {
		t.Fatalf("expected gold unchanged after rejected sell, got %d", reloadedPlayer.Gold)
	}

	reloadedItem, err := inventoryRepo.GetInventoryItemByUser(int64(player.ID), items[0].ID)
	if err != nil {
		t.Fatalf("reload equipped inventory item: %v", err)
	}
	if reloadedItem.Quantity != 1 {
		t.Fatalf("expected equipped item quantity unchanged, got %d", reloadedItem.Quantity)
	}
}

func TestInventoryPostgresIntegration_PurchaseItem_InsufficientGoldRollsBack(t *testing.T) {
	testDB := newIntegrationTestDB(t)

	playerRepo := NewPlayerPostgresRepo(testDB.DB)
	inventoryRepo := NewInventoryPostgresRepo(testDB.DB)

	item := domain.Item{
		Code:        "IT_POOR_SWORD",
		Name:        "Expensive Sword",
		Description: "Equipment used to validate rollback.",
		Type:        domain.ItemEquipment,
		Rarity:      "COMMON",
		Amount:      1,
		Price:       900,
	}
	if err := testDB.DB.Create(&item).Error; err != nil {
		t.Fatalf("create expensive item fixture: %v", err)
	}

	equipment := domain.Equipment{
		ItemID:         item.ID,
		Slot:           string(domain.SlotWeapon),
		RequiredLevel:  1,
		AtkBonus:       9,
		AllowedClasses: "Wind Striker,Earth Sentinel,Ki Vanguard",
	}
	if err := testDB.DB.Create(&equipment).Error; err != nil {
		t.Fatalf("create expensive equipment fixture: %v", err)
	}

	player, err := playerRepo.GetOrCreate(710004)
	if err != nil {
		t.Fatalf("create player: %v", err)
	}

	player.Gold = 100
	if err := playerRepo.Update(player); err != nil {
		t.Fatalf("persist low gold player: %v", err)
	}

	if err := inventoryRepo.PurchaseItem(int64(player.ID), item.ID, 1, 900); !errors.Is(err, domain.ErrInventoryInsufficientGold) {
		t.Fatalf("expected insufficient gold error, got %v", err)
	}

	reloadedPlayer, err := playerRepo.GetOrCreate(710004)
	if err != nil {
		t.Fatalf("reload player after rejected purchase: %v", err)
	}
	if reloadedPlayer.Gold != 100 {
		t.Fatalf("expected gold unchanged after rejected purchase, got %d", reloadedPlayer.Gold)
	}

	items, total, err := inventoryRepo.ListItemsByUserAndType(int64(player.ID), domain.ItemEquipment, 1, 10)
	if err != nil {
		t.Fatalf("list inventory after rejected purchase: %v", err)
	}
	if total != 0 || len(items) != 0 {
		t.Fatalf("expected no inventory rows after rejected purchase, got total=%d len=%d", total, len(items))
	}
}

func TestInventoryPostgresIntegration_SavePlayerAndEquipment_PersistsTogether(t *testing.T) {
	testDB := newIntegrationTestDB(t)

	playerRepo := NewPlayerPostgresRepo(testDB.DB)
	inventoryRepo := NewInventoryPostgresRepo(testDB.DB)

	item := domain.Item{
		Code:        "IT_EQUIP_SAVE",
		Name:        "Persist Sword",
		Description: "Equipment used to validate atomic player/equipment save.",
		Type:        domain.ItemEquipment,
		Rarity:      "COMMON",
		Amount:      1,
		Price:       250,
	}
	if err := testDB.DB.Create(&item).Error; err != nil {
		t.Fatalf("create save item fixture: %v", err)
	}

	equipmentMeta := domain.Equipment{
		ItemID:         item.ID,
		Slot:           string(domain.SlotWeapon),
		RequiredLevel:  1,
		AtkBonus:       6,
		AllowedClasses: "Wind Striker,Earth Sentinel,Ki Vanguard",
	}
	if err := testDB.DB.Create(&equipmentMeta).Error; err != nil {
		t.Fatalf("create save equipment fixture: %v", err)
	}

	userID := time.Now().UnixNano()
	player, err := playerRepo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("create player: %v", err)
	}

	if err := inventoryRepo.CreateOrUpdateInventory(int64(player.ID), item.ID, 1); err != nil {
		t.Fatalf("seed inventory item: %v", err)
	}
	items, _, err := inventoryRepo.ListItemsByUserAndType(int64(player.ID), domain.ItemEquipment, 1, 10)
	if err != nil {
		t.Fatalf("list equipment items: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected one inventory item, got %d", len(items))
	}

	playerEquipment, err := inventoryRepo.GetOrCreatePlayerEquipment(int64(player.ID))
	if err != nil {
		t.Fatalf("get player equipment: %v", err)
	}

	player.ATK = 123
	player.DEF = 77
	player.Speed = 18
	player.MaxHP = 155
	player.MaxEnergy = 28
	player.HP = 144
	player.Energy = 19
	player.Gold = 321
	player.UserName = "persist_sword"
	playerEquipment.WeaponInventoryItemID = &items[0].ID

	if err := inventoryRepo.SavePlayerAndEquipment(player, playerEquipment); err != nil {
		t.Fatalf("save player and equipment atomically: %v", err)
	}

	reloadedPlayer, err := playerRepo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("reload player: %v", err)
	}
	if reloadedPlayer.ATK != 123 || reloadedPlayer.DEF != 77 || reloadedPlayer.Speed != 18 {
		t.Fatalf("unexpected persisted stats: atk=%d def=%d speed=%d", reloadedPlayer.ATK, reloadedPlayer.DEF, reloadedPlayer.Speed)
	}
	if reloadedPlayer.MaxHP != 155 || reloadedPlayer.MaxEnergy != 28 || reloadedPlayer.HP != 144 || reloadedPlayer.Energy != 19 {
		t.Fatalf("unexpected persisted vitality: hp=%d max_hp=%d energy=%d max_energy=%d", reloadedPlayer.HP, reloadedPlayer.MaxHP, reloadedPlayer.Energy, reloadedPlayer.MaxEnergy)
	}
	if reloadedPlayer.Gold != 321 || reloadedPlayer.UserName != "persist_sword" {
		t.Fatalf("unexpected persisted economy/profile data: gold=%d username=%q", reloadedPlayer.Gold, reloadedPlayer.UserName)
	}

	reloadedEquipment, err := inventoryRepo.GetOrCreatePlayerEquipment(int64(player.ID))
	if err != nil {
		t.Fatalf("reload player equipment: %v", err)
	}
	if reloadedEquipment.WeaponInventoryItemID == nil || *reloadedEquipment.WeaponInventoryItemID != items[0].ID {
		t.Fatalf("expected persisted weapon inventory item id %d, got %+v", items[0].ID, reloadedEquipment.WeaponInventoryItemID)
	}
}

func TestInventoryPostgresIntegration_SavePlayerAndEquipment_RollsBackEquipmentWhenPlayerMissing(t *testing.T) {
	testDB := newIntegrationTestDB(t)

	playerRepo := NewPlayerPostgresRepo(testDB.DB)
	inventoryRepo := NewInventoryPostgresRepo(testDB.DB)

	item := domain.Item{
		Code:        "IT_EQUIP_ROLLBACK",
		Name:        "Rollback Sword",
		Description: "Equipment used to validate rollback.",
		Type:        domain.ItemEquipment,
		Rarity:      "COMMON",
		Amount:      1,
		Price:       250,
	}
	if err := testDB.DB.Create(&item).Error; err != nil {
		t.Fatalf("create rollback item fixture: %v", err)
	}

	equipmentMeta := domain.Equipment{
		ItemID:         item.ID,
		Slot:           string(domain.SlotWeapon),
		RequiredLevel:  1,
		AtkBonus:       6,
		AllowedClasses: "Wind Striker,Earth Sentinel,Ki Vanguard",
	}
	if err := testDB.DB.Create(&equipmentMeta).Error; err != nil {
		t.Fatalf("create rollback equipment fixture: %v", err)
	}

	userID := time.Now().UnixNano()
	player, err := playerRepo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("create player: %v", err)
	}

	if err := inventoryRepo.CreateOrUpdateInventory(int64(player.ID), item.ID, 1); err != nil {
		t.Fatalf("seed inventory item: %v", err)
	}
	items, _, err := inventoryRepo.ListItemsByUserAndType(int64(player.ID), domain.ItemEquipment, 1, 10)
	if err != nil {
		t.Fatalf("list equipment items: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected one inventory item, got %d", len(items))
	}

	playerEquipment, err := inventoryRepo.GetOrCreatePlayerEquipment(int64(player.ID))
	if err != nil {
		t.Fatalf("get player equipment: %v", err)
	}
	if playerEquipment.WeaponInventoryItemID != nil {
		t.Fatalf("expected empty equipment before rollback scenario")
	}

	missingPlayer := *player
	missingPlayer.ID = player.ID + 999999
	missingPlayer.ATK = 999

	rollbackEquipment := *playerEquipment
	rollbackEquipment.WeaponInventoryItemID = &items[0].ID

	err = inventoryRepo.SavePlayerAndEquipment(&missingPlayer, &rollbackEquipment)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record not found on missing player update, got %v", err)
	}

	reloadedPlayer, err := playerRepo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("reload original player: %v", err)
	}
	if reloadedPlayer.ATK == 999 {
		t.Fatalf("expected player rollback to keep original stats")
	}

	reloadedEquipment, err := inventoryRepo.GetOrCreatePlayerEquipment(int64(player.ID))
	if err != nil {
		t.Fatalf("reload original equipment: %v", err)
	}
	if reloadedEquipment.WeaponInventoryItemID != nil {
		t.Fatalf("expected equipment rollback to keep weapon slot empty, got %+v", reloadedEquipment.WeaponInventoryItemID)
	}
}

func TestInventoryPostgresIntegration_SavePlayerAndInventoryItem_PersistsTogether(t *testing.T) {
	testDB := newIntegrationTestDB(t)

	playerRepo := NewPlayerPostgresRepo(testDB.DB)
	inventoryRepo := NewInventoryPostgresRepo(testDB.DB)

	item := domain.Item{
		Code:        "IT_POTION_ATOMIC_SAVE",
		Name:        "Atomic Potion",
		Description: "Consumable used to validate atomic player/item save.",
		Type:        domain.ItemConsumable,
		Rarity:      "COMMON",
		Amount:      1,
		Price:       100,
	}
	if err := testDB.DB.Create(&item).Error; err != nil {
		t.Fatalf("create atomic potion fixture: %v", err)
	}

	userID := time.Now().UnixNano()
	player, err := playerRepo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("create player: %v", err)
	}

	player.HP = 50
	player.Energy = 7
	if err := playerRepo.Update(player); err != nil {
		t.Fatalf("persist player state: %v", err)
	}

	if err := inventoryRepo.CreateOrUpdateInventory(int64(player.ID), item.ID, 2); err != nil {
		t.Fatalf("seed consumable inventory item: %v", err)
	}
	items, _, err := inventoryRepo.ListItemsByUserAndType(int64(player.ID), domain.ItemConsumable, 1, 10)
	if err != nil {
		t.Fatalf("list consumable items: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected one consumable row, got %d", len(items))
	}

	player.HP = 90
	player.Energy = 12
	items[0].Quantity = 1

	if err := inventoryRepo.SavePlayerAndInventoryItem(player, &items[0]); err != nil {
		t.Fatalf("save player and inventory item atomically: %v", err)
	}

	reloadedPlayer, err := playerRepo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("reload player: %v", err)
	}
	if reloadedPlayer.HP != 90 || reloadedPlayer.Energy != 12 {
		t.Fatalf("unexpected persisted player after atomic save: hp=%d energy=%d", reloadedPlayer.HP, reloadedPlayer.Energy)
	}

	reloadedItem, err := inventoryRepo.GetInventoryItemByUser(int64(player.ID), items[0].ID)
	if err != nil {
		t.Fatalf("reload inventory item: %v", err)
	}
	if reloadedItem.Quantity != 1 {
		t.Fatalf("expected persisted quantity 1, got %d", reloadedItem.Quantity)
	}
}

func TestInventoryPostgresIntegration_SavePlayerAndInventoryItem_RollsBackItemWhenPlayerMissing(t *testing.T) {
	testDB := newIntegrationTestDB(t)

	playerRepo := NewPlayerPostgresRepo(testDB.DB)
	inventoryRepo := NewInventoryPostgresRepo(testDB.DB)

	item := domain.Item{
		Code:        "IT_POTION_ATOMIC_ROLLBACK",
		Name:        "Rollback Potion",
		Description: "Consumable used to validate rollback.",
		Type:        domain.ItemConsumable,
		Rarity:      "COMMON",
		Amount:      1,
		Price:       100,
	}
	if err := testDB.DB.Create(&item).Error; err != nil {
		t.Fatalf("create rollback potion fixture: %v", err)
	}

	userID := time.Now().UnixNano()
	player, err := playerRepo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("create player: %v", err)
	}

	if err := inventoryRepo.CreateOrUpdateInventory(int64(player.ID), item.ID, 2); err != nil {
		t.Fatalf("seed rollback consumable item: %v", err)
	}
	items, _, err := inventoryRepo.ListItemsByUserAndType(int64(player.ID), domain.ItemConsumable, 1, 10)
	if err != nil {
		t.Fatalf("list rollback consumables: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected one consumable row, got %d", len(items))
	}

	missingPlayer := *player
	missingPlayer.ID = player.ID + 999999
	missingPlayer.HP = 1

	mutatedItem := items[0]
	mutatedItem.Quantity = 0

	err = inventoryRepo.SavePlayerAndInventoryItem(&missingPlayer, &mutatedItem)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record not found on missing player update, got %v", err)
	}

	reloadedPlayer, err := playerRepo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("reload original player: %v", err)
	}
	if reloadedPlayer.HP == 1 {
		t.Fatalf("expected player rollback to keep original hp")
	}

	reloadedItem, err := inventoryRepo.GetInventoryItemByUser(int64(player.ID), items[0].ID)
	if err != nil {
		t.Fatalf("reload original inventory item: %v", err)
	}
	if reloadedItem.Quantity != 2 {
		t.Fatalf("expected inventory rollback to keep quantity 2, got %d", reloadedItem.Quantity)
	}
}
