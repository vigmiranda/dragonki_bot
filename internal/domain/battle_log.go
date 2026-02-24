package domain

import "time"

type BattleLog struct {
	ID uint `gorm:"primaryKey"`

	PlayerID  int64 `gorm:"index"`
	MonsterID int

	Result string // WIN | LOSE | ESCAPE

	OldXP       int
	XPGained    int
	XPCurrent   int
	OldGold     int
	GoldGained  int
	GoldCurrent int

	PlayerDamageTotal  int
	MonsterDamageTotal int

	CreatedAt time.Time `gorm:"index"`
}
