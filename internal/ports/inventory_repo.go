package ports

import "github.com/vigmiranda/dragonki/internal/domain"

type InventoryRepository interface {
	UpdateInventoryItem(item *domain.InventoryItem) error
	CreateOrUpdateInventory(
		userID int64,
		itemID uint,
		amount int,
	) error
	ListItemsByUserAndType(
		playerID int64,
		itemType domain.ItemType,
		page int,
		pageSize int,
	) ([]domain.InventoryItem, int64, error)
	GetInventoryItemByUser(
		playerID int64,
		inventoryItemID uint,
	) (*domain.InventoryItem, error)
}
