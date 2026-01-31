package domain

type ItemType string

const (
	ItemConsumable ItemType = "CONSUMABLE"
	ItemEquipment  ItemType = "EQUIPMENT"
	ItemKey        ItemType = "KEY"
)

type Item struct {
	ID          int64
	Code        string
	Name        string
	Description string
	Type        ItemType
	Rarity      ItemRarity
}
