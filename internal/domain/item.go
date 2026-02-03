package domain

type ItemType string

const (
	ItemConsumable ItemType = "CONSUMABLE"
	ItemEquipment  ItemType = "EQUIPMENT"
	ItemKey        ItemType = "KEY"
)

type Item struct {
	ID uint `gorm:"primaryKey"`

	Code        string
	Name        string
	Description string
	Type        ItemType
	Rarity      string
	Amount      int
	Price       int
}
