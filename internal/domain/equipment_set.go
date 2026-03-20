package domain

type EquipmentSet struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Description string
}

func (EquipmentSet) TableName() string {
	return "equipment_sets"
}

type EquipmentSetBonus struct {
	ID             uint `gorm:"primaryKey"`
	SetID          uint
	PiecesRequired int
	AtkBonus       int
	DefBonus       int
	SpeedBonus     int
	HPBonus        int
}

func (EquipmentSetBonus) TableName() string {
	return "equipment_set_bonuses"
}
