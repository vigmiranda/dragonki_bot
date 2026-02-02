package domain

type Equipment struct {
	ID uint `gorm:"primaryKey"`

	ItemID uint
	Item   Item `gorm:"foreignKey:ItemID"`

	Slot string // weapon, gloves, boots, ring, amulet

	BonusHP     int
	BonusEnergy int
	BonusATK    int
	BonusDEF    int
}
