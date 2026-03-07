package domain

type Equipment struct {
	ID uint `gorm:"primaryKey"`

	ItemID uint `gorm:"uniqueIndex;not null"`
	Item   Item `gorm:"foreignKey:ItemID"`

	Slot           string // weapon, gloves, boots, ring, amulet
	RequiredLevel  int
	AtkBonus       int
	DefBonus       int
	SpeedBonus     int
	HPBonus        int
	EnergyBonus    int
	AllowedClasses string
}

func (Equipment) TableName() string {
	return "equipments"
}
