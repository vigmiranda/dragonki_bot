package ports

type InventoryRepository interface {
	CreateOrUpdateInventory(
		userID int64,
		itemID uint,
		amount int,
	) error
}
