package ports

import "github.com/vigmiranda/dragonki/internal/domain"

type InventoryRepository interface {
	UpdateInventoryItem(item *domain.InventoryItem) error
	CreateOrUpdateInventory(
		userID int64,
		itemID uint,
		amount int,
	) error
	PurchaseItem(playerID int64, itemID uint, amount int, price int) error
	SellInventoryItem(playerID int64, inventoryItemID uint, goldValue int) (*domain.InventoryItem, error)
	ListItemsByUserAndType(
		playerID int64,
		itemType domain.ItemType,
		page int,
		pageSize int,
	) ([]domain.InventoryItem, int64, error)
	ListEquipmentItemsByUserAndSlot(
		playerID int64,
		slot domain.EquipmentSlot,
		page int,
		pageSize int,
	) ([]domain.InventoryItem, int64, error)
	CountEquipmentItemsBySlot(playerID int64) (map[domain.EquipmentSlot]int64, error)
	GetInventoryItemByUser(
		playerID int64,
		inventoryItemID uint,
	) (*domain.InventoryItem, error)
	GetOrCreatePlayerEquipment(playerID int64) (*domain.PlayerEquipment, error)
	UpdatePlayerEquipment(equipment *domain.PlayerEquipment) error
	SavePlayerAndInventoryItem(player *domain.Player, item *domain.InventoryItem) error
	SavePlayerAndEquipment(player *domain.Player, equipment *domain.PlayerEquipment) error
}
