package domain

type MonsterDrop struct {
	ID uint `gorm:"primary_key"`

	MonsterID int
	ItemID    uint

	Item Item `grom:"foreignkey:ItemID"`

	DropRate float64
	MinQty   int
	MaxQty   int

	RarityWeight int
}
