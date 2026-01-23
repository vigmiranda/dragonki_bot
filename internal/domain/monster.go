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
	MaxMoney    int
	SpawnWeight int
}
