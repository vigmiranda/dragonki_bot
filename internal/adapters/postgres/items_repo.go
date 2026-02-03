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
		Where("type = ?", itemType).
		Where("rarity = ?", rarity).
		Order("RANDOM()").
		Limit(3).
		Find(&items).
		Error

	if err != nil {
		return nil, err
	}

	return items, nil

}
func (r *ItemsPostgresRepo) GetItemById(itemID int) (*domain.Item, error) {
	var item *domain.Item
	err := r.db.Model(&domain.Item{}).
		Where("id = ?", itemID).
		First(&item).
		Error
	if err != nil {
		return nil, err
	}

	return item, nil
}
