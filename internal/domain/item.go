package domain

type ItemType string

const (
	ItemConsumable ItemType = "CONSUMABLE"
	ItemEquipment  ItemType = "EQUIPMENT"
	ItemKey        ItemType = "KEY"
)

type Item struct {
	ID uint `gorm:"primaryKey"`

	Code        string `gorm:"uniqueIndex"`
	Name        string
	Description string
	Type        ItemType
	Rarity      string
}
