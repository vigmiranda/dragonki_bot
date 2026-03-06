package ports

import "github.com/vigmiranda/dragonki/internal/domain"

type ItemsRepository interface {
	GetRandomList(itemType domain.ItemType, rarity domain.ItemRarity) ([]domain.Item, error)
	GetRandomEquipmentsBySlot(slot domain.EquipmentSlot, limit int) ([]domain.Item, error)
	GetItemById(id int) (*domain.Item, error)
}
