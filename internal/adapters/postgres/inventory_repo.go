package postgres

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InventoryPostgres struct {
	db *gorm.DB
}

func NewInventoryPostgresRepo(db *gorm.DB) *InventoryPostgres {
	return &InventoryPostgres{db: db}
}

func (r *InventoryPostgres) UpdateInventoryItem(item *domain.InventoryItem) error {
	return r.db.Save(item).Error
}

func (r *InventoryPostgres) ListItemsByUserAndType(
	playerID int64,
	itemType domain.ItemType,
	page int,
	pageSize int,
) ([]domain.InventoryItem, int64, error) {
	var inventory domain.Inventory
	if err := r.db.Where("player_id = ?", playerID).First(&inventory).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []domain.InventoryItem{}, 0, nil
		}
		return nil, 0, err
	}

	excludedIDs, err := r.getEquippedInventoryItemIDs(playerID)
	if err != nil {
		return nil, 0, err
	}

	var (
		items []domain.InventoryItem
		total int64
	)
	offset := (page - 1) * pageSize

	countQuery := r.db.
		Model(&domain.InventoryItem{}).
		Joins("JOIN items ON items.id = inventory_items.item_id").
		Where("inventory_items.inventory_id = ?", inventory.ID).
		Where("items.type = ?", itemType).
		Where("inventory_items.quantity > 0")
	if itemType == domain.ItemEquipment && len(excludedIDs) > 0 {
		countQuery = countQuery.Where("inventory_items.id NOT IN ?", excludedIDs)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	findQuery := r.db.
		Model(&domain.InventoryItem{}).
		Joins("JOIN items ON items.id = inventory_items.item_id").
		Where("inventory_items.inventory_id = ?", inventory.ID).
		Where("items.type = ?", itemType).
		Where("inventory_items.quantity > 0")
	if itemType == domain.ItemEquipment && len(excludedIDs) > 0 {
		findQuery = findQuery.Where("inventory_items.id NOT IN ?", excludedIDs)
	}
	if err := findQuery.
		Preload("Item").
		Preload("Item.Equipment").
		Order("inventory_items.id ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	if err := r.attachEquipmentByItemID(items); err != nil {
		return nil, 0, err
	}
	hydrateInventoryItemsEquipmentMetadata(items)

	return items, total, nil
}

func (r *InventoryPostgres) ListEquipmentItemsByUserAndSlot(
	playerID int64,
	slot domain.EquipmentSlot,
	page int,
	pageSize int,
) ([]domain.InventoryItem, int64, error) {
	var inventory domain.Inventory
	if err := r.db.Where("player_id = ?", playerID).First(&inventory).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []domain.InventoryItem{}, 0, nil
		}
		return nil, 0, err
	}

	excludedIDs, err := r.getEquippedInventoryItemIDs(playerID)
	if err != nil {
		return nil, 0, err
	}

	var (
		items []domain.InventoryItem
		total int64
	)
	offset := (page - 1) * pageSize

	countQuery := r.db.
		Model(&domain.InventoryItem{}).
		Joins("JOIN items ON items.id = inventory_items.item_id").
		Joins("JOIN equipments ON equipments.item_id = items.id").
		Where("inventory_items.inventory_id = ?", inventory.ID).
		Where("items.type = ?", domain.ItemEquipment).
		Where("equipments.slot = ?", slot).
		Where("inventory_items.quantity > 0")
	if len(excludedIDs) > 0 {
		countQuery = countQuery.Where("inventory_items.id NOT IN ?", excludedIDs)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	findQuery := r.db.
		Model(&domain.InventoryItem{}).
		Joins("JOIN items ON items.id = inventory_items.item_id").
		Joins("JOIN equipments ON equipments.item_id = items.id").
		Where("inventory_items.inventory_id = ?", inventory.ID).
		Where("items.type = ?", domain.ItemEquipment).
		Where("equipments.slot = ?", slot).
		Where("inventory_items.quantity > 0")
	if len(excludedIDs) > 0 {
		findQuery = findQuery.Where("inventory_items.id NOT IN ?", excludedIDs)
	}
	if err := findQuery.
		Preload("Item").
		Preload("Item.Equipment").
		Order("inventory_items.id ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	if err := r.attachEquipmentByItemID(items); err != nil {
		return nil, 0, err
	}
	hydrateInventoryItemsEquipmentMetadata(items)

	return items, total, nil
}

func (r *InventoryPostgres) CountEquipmentItemsBySlot(playerID int64) (map[domain.EquipmentSlot]int64, error) {
	result := map[domain.EquipmentSlot]int64{
		domain.SlotWeapon:    0,
		domain.SlotArmor:     0,
		domain.SlotGloves:    0,
		domain.SlotHelmet:    0,
		domain.SlotBoots:     0,
		domain.SlotRing:      0,
		domain.SlotAmulet:    0,
		domain.SlotAccessory: 0,
	}

	var inventory domain.Inventory
	if err := r.db.Where("player_id = ?", playerID).First(&inventory).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, nil
		}
		return nil, err
	}

	excludedIDs, err := r.getEquippedInventoryItemIDs(playerID)
	if err != nil {
		return nil, err
	}

	type slotCount struct {
		Slot  string
		Total int64
	}
	var rows []slotCount

	query := r.db.
		Table("inventory_items").
		Select("equipments.slot AS slot, COUNT(inventory_items.id) AS total").
		Joins("JOIN items ON items.id = inventory_items.item_id").
		Joins("JOIN equipments ON equipments.item_id = items.id").
		Where("inventory_items.inventory_id = ?", inventory.ID).
		Where("items.type = ?", domain.ItemEquipment).
		Where("inventory_items.quantity > 0")
	if len(excludedIDs) > 0 {
		query = query.Where("inventory_items.id NOT IN ?", excludedIDs)
	}
	if err := query.Group("equipments.slot").Scan(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		s, ok := normalizeEquipmentSlot(row.Slot)
		if ok {
			result[s] += row.Total
		}
	}
	return result, nil
}

func (r *InventoryPostgres) CreateOrUpdateInventory(userID int64, itemID uint, amount int) error {
	return createOrUpdateInventoryWithDB(r.db, userID, itemID, amount)
}

func (r *InventoryPostgres) PurchaseItem(playerID int64, itemID uint, amount int, price int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var player domain.Player
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", playerID).
			Take(&player).Error; err != nil {
			return err
		}

		if player.Gold < price {
			return domain.ErrInventoryInsufficientGold
		}

		if err := createOrUpdateInventoryWithDB(tx, playerID, itemID, amount); err != nil {
			return err
		}

		return tx.
			Model(&domain.Player{}).
			Where("id = ?", playerID).
			Updates(map[string]any{
				"gold":        gorm.Expr("gold - ?", price),
				"last_update": time.Now().UTC(),
			}).Error
	})
}

func (r *InventoryPostgres) SellInventoryItem(playerID int64, inventoryItemID uint, goldValue int) (*domain.InventoryItem, error) {
	var soldItem domain.InventoryItem

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var player domain.Player
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", playerID).
			Take(&player).Error; err != nil {
			return err
		}

		var equipment domain.PlayerEquipment
		err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("player_id = ?", playerID).
			Take(&equipment).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if equipmentUsesInventoryItem(&equipment, inventoryItemID) {
			return domain.ErrInventoryItemEquipped
		}

		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Joins("JOIN inventories ON inventories.id = inventory_items.inventory_id").
			Preload("Item").
			Where("inventory_items.id = ?", inventoryItemID).
			Where("inventories.player_id = ?", playerID).
			First(&soldItem).Error; err != nil {
			return err
		}

		if soldItem.Quantity <= 0 || soldItem.Item.Type != domain.ItemEquipment {
			return domain.ErrInventoryItemNotSellable
		}

		soldItem.Quantity--
		if err := tx.Save(&soldItem).Error; err != nil {
			return err
		}

		return tx.
			Model(&domain.Player{}).
			Where("id = ?", playerID).
			Updates(map[string]any{
				"gold":        gorm.Expr("gold + ?", goldValue),
				"last_update": time.Now().UTC(),
			}).Error
	})
	if err != nil {
		return nil, err
	}

	return &soldItem, nil
}

func createOrUpdateInventoryWithDB(db *gorm.DB, userID int64, itemID uint, amount int) error {
	var item domain.Item
	if err := db.First(&item, itemID).Error; err != nil {
		return err
	}

	var inventory domain.Inventory
	err := db.Where("player_id = ?", userID).First(&inventory).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		inventory = domain.Inventory{PlayerID: userID}
		if err := db.Create(&inventory).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	switch item.Type {
	case domain.ItemConsumable:
		var invItem domain.InventoryItem
		err := db.Where("inventory_id = ? AND code = ?", inventory.ID, item.Code).First(&invItem).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			invItem = domain.InventoryItem{InventoryID: inventory.ID, ItemID: item.ID, Code: item.Code, Quantity: amount}
			return db.Create(&invItem).Error
		}
		if err != nil {
			return err
		}
		invItem.Quantity += amount
		return db.Save(&invItem).Error

	case domain.ItemEquipment:
		if amount < 1 {
			amount = 1
		}
		for i := 0; i < amount; i++ {
			invItem := domain.InventoryItem{InventoryID: inventory.ID, ItemID: item.ID, Code: item.Code, Quantity: 1}
			if err := db.Create(&invItem).Error; err != nil {
				return err
			}
		}
		return nil

	case domain.ItemKey:
		invItem := domain.InventoryItem{InventoryID: inventory.ID, ItemID: item.ID, Code: item.Code, Quantity: 1}
		return db.Create(&invItem).Error

	default:
		return fmt.Errorf("tipo de item invalido: %s", item.Type)
	}
}

func (r *InventoryPostgres) GetInventoryItemByUser(playerID int64, inventoryItemID uint) (*domain.InventoryItem, error) {
	var invItem domain.InventoryItem
	err := r.db.
		Joins("JOIN inventories ON inventories.id = inventory_items.inventory_id").
		Preload("Item").
		Preload("Item.Equipment").
		Where("inventory_items.id = ?", inventoryItemID).
		Where("inventories.player_id = ?", playerID).
		First(&invItem).Error
	if err != nil {
		return nil, err
	}
	if err := r.attachEquipmentByInventoryItem(&invItem); err != nil {
		return nil, err
	}
	invItem.Item.HydrateEquipmentMetadata()
	return &invItem, nil
}

func (r *InventoryPostgres) GetOrCreatePlayerEquipment(playerID int64) (*domain.PlayerEquipment, error) {
	var pe domain.PlayerEquipment
	err := r.db.Where("player_id = ?", playerID).First(&pe).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		pe = domain.PlayerEquipment{PlayerID: playerID}
		if err := r.db.Create(&pe).Error; err != nil {
			return nil, err
		}
		return &pe, nil
	}
	if err != nil {
		return nil, err
	}
	return &pe, nil
}

func (r *InventoryPostgres) UpdatePlayerEquipment(equipment *domain.PlayerEquipment) error {
	return r.db.Save(equipment).Error
}

func (r *InventoryPostgres) SavePlayerAndInventoryItem(player *domain.Player, item *domain.InventoryItem) error {
	if player == nil || item == nil {
		return gorm.ErrInvalidData
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(item).Error; err != nil {
			return err
		}

		return updatePlayerWithTx(tx, player)
	})
}

func (r *InventoryPostgres) SavePlayerAndEquipment(player *domain.Player, equipment *domain.PlayerEquipment) error {
	if player == nil || equipment == nil {
		return gorm.ErrInvalidData
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(equipment).Error; err != nil {
			return err
		}

		return updatePlayerWithTx(tx, player)
	})
}

func (r *InventoryPostgres) getEquippedInventoryItemIDs(playerID int64) ([]uint, error) {
	var pe domain.PlayerEquipment
	err := r.db.Where("player_id = ?", playerID).First(&pe).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	seen := map[uint]struct{}{}
	for _, ptr := range []*uint{
		pe.WeaponInventoryItemID,
		pe.ArmorInventoryItemID,
		pe.GlovesInventoryItemID,
		pe.HelmetInventoryItemID,
		pe.BootsInventoryItemID,
		pe.RingInventoryItemID,
		pe.AmuletInventoryItemID,
		pe.AccessoryInventoryItemID,
	} {
		if ptr != nil && *ptr > 0 {
			seen[*ptr] = struct{}{}
		}
	}

	out := make([]uint, 0, len(seen))
	for id := range seen {
		out = append(out, id)
	}
	return out, nil
}

func normalizeEquipmentSlot(raw string) (domain.EquipmentSlot, bool) {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "WEAPON":
		return domain.SlotWeapon, true
	case "ARMOR":
		return domain.SlotArmor, true
	case "GLOVES":
		return domain.SlotGloves, true
	case "HELMET":
		return domain.SlotHelmet, true
	case "BOOTS":
		return domain.SlotBoots, true
	case "RING":
		return domain.SlotRing, true
	case "AMULET":
		return domain.SlotAmulet, true
	case "ACCESSORY":
		return domain.SlotAccessory, true
	default:
		return "", false
	}
}

func hydrateInventoryItemsEquipmentMetadata(items []domain.InventoryItem) {
	for i := range items {
		items[i].Item.HydrateEquipmentMetadata()
	}
}

func (r *InventoryPostgres) attachEquipmentByItemID(items []domain.InventoryItem) error {
	itemIDs := make([]uint, 0, len(items))
	itemIndexByID := make(map[uint][]int, len(items))
	for i := range items {
		itemID := items[i].Item.ID
		if itemID == 0 {
			itemID = items[i].ItemID
		}
		if itemID == 0 {
			continue
		}
		itemIDs = append(itemIDs, itemID)
		itemIndexByID[itemID] = append(itemIndexByID[itemID], i)
	}
	if len(itemIDs) == 0 {
		return nil
	}

	var equipments []domain.Equipment
	if err := r.db.
		Model(&domain.Equipment{}).
		Where("item_id IN ?", itemIDs).
		Find(&equipments).Error; err != nil {
		return err
	}

	for i := range equipments {
		indexes, ok := itemIndexByID[equipments[i].ItemID]
		if !ok {
			continue
		}
		for _, idx := range indexes {
			eq := equipments[i]
			items[idx].Item.Equipment = &eq
		}
	}
	return nil
}

func (r *InventoryPostgres) attachEquipmentByInventoryItem(invItem *domain.InventoryItem) error {
	if invItem == nil {
		return nil
	}
	itemID := invItem.Item.ID
	if itemID == 0 {
		itemID = invItem.ItemID
	}
	if itemID == 0 {
		return nil
	}

	var equipment domain.Equipment
	if err := r.db.
		Model(&domain.Equipment{}).
		Where("item_id = ?", itemID).
		First(&equipment).Error; err != nil {
		return nil
	}
	invItem.Item.Equipment = &equipment
	return nil
}

func equipmentUsesInventoryItem(equipment *domain.PlayerEquipment, inventoryItemID uint) bool {
	if equipment == nil || inventoryItemID == 0 {
		return false
	}

	for _, ptr := range []*uint{
		equipment.WeaponInventoryItemID,
		equipment.ArmorInventoryItemID,
		equipment.GlovesInventoryItemID,
		equipment.HelmetInventoryItemID,
		equipment.BootsInventoryItemID,
		equipment.RingInventoryItemID,
		equipment.AmuletInventoryItemID,
		equipment.AccessoryInventoryItemID,
	} {
		if ptr != nil && *ptr == inventoryItemID {
			return true
		}
	}
	return false
}

func updatePlayerWithTx(tx *gorm.DB, player *domain.Player) error {
	if tx == nil || player == nil {
		return gorm.ErrInvalidData
	}

	updates := map[string]any{
		"name":              player.Name,
		"gender":            player.Gender,
		"class_id":          player.ClassID,
		"level":             player.Level,
		"xp":                player.XP,
		"hp":                player.HP,
		"max_hp":            player.MaxHP,
		"atk":               player.ATK,
		"def":               player.DEF,
		"speed":             player.Speed,
		"energy":            player.Energy,
		"max_energy":        player.MaxEnergy,
		"gold":              player.Gold,
		"ki_crystals":       player.KiCrystals,
		"region_id":         player.RegionID,
		"last_energy_at":    player.LastEnergyAt,
		"form":              player.Form,
		"form_battles_left": player.FormBattlesLeft,
		"onboarding_step":   player.OnboardingStep,
		"last_update":       time.Now().UTC(),
		"user_id_telegram":  player.UserIDTelegram,
		"name_telegram":     player.NameTelegram,
		"user_name":         player.UserName,
		"dodge_accumulator": player.DodgeAccumulator,
	}

	result := tx.
		Model(&domain.Player{}).
		Where("id = ?", player.ID).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
