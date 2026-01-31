package domain

type InventoryItem struct {
	ItemID   int64
	Quantity int
}

type Inventory struct {
	PlayerID int64
	Items    []InventoryItem
}
