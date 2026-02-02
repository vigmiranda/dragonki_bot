package domain

type PlayerEquipment struct {
	ID       uint  `gorm:"primaryKey"`
	PlayerID int64 `gorm:"uniqueIndex"`

	WeaponID uint
	Weapon   Equipment `gorm:"foreignKey:WeaponID"`

	GlovesID uint
	Gloves   Equipment `gorm:"foreignKey:GlovesID"`

	BootsID uint
	Boots   Equipment `gorm:"foreignKey:BootsID"`

	RingID uint
	Ring   Equipment `gorm:"foreignKey:RingID"`

	AmuletID uint
	Amulet   Equipment `gorm:"foreignKey:AmuletID"`
}
