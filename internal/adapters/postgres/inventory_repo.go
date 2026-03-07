package postgres

import (
	"errors"
	"fmt"
	"strings"

	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
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
	var item domain.Item
	if err := r.db.First(&item, itemID).Error; err != nil {
		return err
	}

	var inventory domain.Inventory
	err := r.db.Where("player_id = ?", userID).First(&inventory).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		inventory = domain.Inventory{PlayerID: userID}
		if err := r.db.Create(&inventory).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	switch item.Type {
	case domain.ItemConsumable:
		var invItem domain.InventoryItem
		err := r.db.Where("inventory_id = ? AND code = ?", inventory.ID, item.Code).First(&invItem).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			invItem = domain.InventoryItem{InventoryID: inventory.ID, ItemID: item.ID, Code: item.Code, Quantity: amount}
			return r.db.Create(&invItem).Error
		}
		if err != nil {
			return err
		}
		invItem.Quantity += amount
		return r.db.Save(&invItem).Error

	case domain.ItemEquipment:
		if amount < 1 {
			amount = 1
		}
		for i := 0; i < amount; i++ {
			invItem := domain.InventoryItem{InventoryID: inventory.ID, ItemID: item.ID, Code: item.Code, Quantity: 1}
			if err := r.db.Create(&invItem).Error; err != nil {
				return err
			}
		}
		return nil

	case domain.ItemKey:
		invItem := domain.InventoryItem{InventoryID: inventory.ID, ItemID: item.ID, Code: item.Code, Quantity: 1}
		return r.db.Create(&invItem).Error

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
