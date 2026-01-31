package domain

type PlayerEquipment struct {
	PlayerID int64

	Weapon Equipment
	Gloves Equipment
	Boots  Equipment
	Ring   Equipment
	Amulet Equipment
}
