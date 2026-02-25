package domain

type Monster struct {
	MonsterID   int
	Name        string
	Tag         string
	Level       int64
	HP          int
	ATK         int
	DEF         int
	MaxXP       int
	MinXP       int `gorm:"default:1"`
	MaxMoney    int
	MinMoney    int `gorm:"default:1"`
	SpawnWeight int
	RegionID    int64
}
