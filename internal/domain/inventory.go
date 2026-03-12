package domain

import "errors"

var (
	ErrInventoryInsufficientGold = errors.New("inventory insufficient gold")
	ErrInventoryItemNotSellable  = errors.New("inventory item not sellable")
	ErrInventoryItemEquipped     = errors.New("inventory item equipped")
)

type InventoryItem struct {
	ID          uint `gorm:"primaryKey"`
	InventoryID uint

	ItemID uint
	Item   Item `gorm:"foreignKey:ItemID"`
	Code   string

	Quantity int
}

type Inventory struct {
	ID       uint  `gorm:"primaryKey"`
	PlayerID int64 `gorm:"uniqueIndex"`

	Items []InventoryItem `gorm:"foreignKey:InventoryID"`
}
