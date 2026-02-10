package domain

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
