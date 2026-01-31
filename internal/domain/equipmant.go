package domain

type Equipment struct {
	Item
	Slot EquipmentSlot

	BonusHP     int
	BonusEnergy int
	BonusATK    int
	BonusDEF    int
}
