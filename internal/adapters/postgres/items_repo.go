package postgres

import (
	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
)

type ItemsPostgresRepo struct {
	db *gorm.DB
}

func NewItemsPostgresRepo(db *gorm.DB) *ItemsPostgresRepo {
	return &ItemsPostgresRepo{
		db: db,
	}
}

func (r *ItemsPostgresRepo) GetRandomList(itemType domain.ItemType, rarity domain.ItemRarity) ([]domain.Item, error) {
	var items []domain.Item
	err := r.db.Model(&domain.Item{}).
		Preload("Equipment").
		Where("type = ?", itemType).
		Where("rarity = ?", rarity).
		Order("RANDOM()").
		Limit(3).
		Find(&items).
		Error

	if err != nil {
		return nil, err
	}

	if err := r.attachEquipmentByItemID(items); err != nil {
		return nil, err
	}
	hydrateItemSliceEquipmentMetadata(items)

	return items, nil

}

func (r *ItemsPostgresRepo) GetRandomEquipmentsBySlot(slot domain.EquipmentSlot, limit int) ([]domain.Item, error) {
	if limit <= 0 {
		limit = 3
	}

	var items []domain.Item
	err := r.db.Model(&domain.Item{}).
		Preload("Equipment").
		Joins("JOIN equipments ON equipments.item_id = items.id").
		Where("type = ?", domain.ItemEquipment).
		Where("equipments.slot = ?", slot).
		Order("RANDOM()").
		Limit(limit).
		Find(&items).
		Error

	if err != nil {
		return nil, err
	}

	if err := r.attachEquipmentByItemID(items); err != nil {
		return nil, err
	}
	hydrateItemSliceEquipmentMetadata(items)

	return items, nil
}

func (r *ItemsPostgresRepo) GetItemById(itemID int) (*domain.Item, error) {
	var item domain.Item
	err := r.db.Model(&domain.Item{}).
		Preload("Equipment").
		Where("id = ?", itemID).
		First(&item).
		Error
	if err != nil {
		return nil, err
	}

	if err := r.attachEquipmentByItemIDSingle(&item); err != nil {
		return nil, err
	}
	item.HydrateEquipmentMetadata()
	return &item, nil
}

func hydrateItemSliceEquipmentMetadata(items []domain.Item) {
	for i := range items {
		items[i].HydrateEquipmentMetadata()
	}
}

func (r *ItemsPostgresRepo) attachEquipmentByItemID(items []domain.Item) error {
	itemIDs := make([]uint, 0, len(items))
	indexByID := make(map[uint]int, len(items))
	for i := range items {
		itemIDs = append(itemIDs, items[i].ID)
		indexByID[items[i].ID] = i
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
		idx, ok := indexByID[equipments[i].ItemID]
		if !ok {
			continue
		}
		eq := equipments[i]
		items[idx].Equipment = &eq
	}
	return nil
}

func (r *ItemsPostgresRepo) attachEquipmentByItemIDSingle(item *domain.Item) error {
	if item == nil || item.ID == 0 {
		return nil
	}
	var equipment domain.Equipment
	if err := r.db.
		Model(&domain.Equipment{}).
		Where("item_id = ?", item.ID).
		First(&equipment).Error; err != nil {
		return nil
	}
	item.Equipment = &equipment
	return nil
}
