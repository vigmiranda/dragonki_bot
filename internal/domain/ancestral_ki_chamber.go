package domain

import (
	"errors"
	"time"
)

var (
	ErrAncestralKiChamberRunActive            = errors.New("ancestral ki chamber active run exists")
	ErrAncestralKiChamberInsufficientResource = errors.New("ancestral ki chamber insufficient resources")
)

const (
	AncestralKiChamberRunStatusActive    = "ACTIVE"
	AncestralKiChamberRunStatusCompleted = "COMPLETED"
	AncestralKiChamberRunStatusFailed    = "FAILED"
	AncestralKiChamberRunStatusAbandoned = "ABANDONED"
)

type AncestralKiChamberProgress struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	PlayerID              int64 `gorm:"not null;uniqueIndex"`
	HighestCompletedLevel int   `gorm:"not null;default:0"`

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (AncestralKiChamberProgress) TableName() string {
	return "ancestral_ki_chamber_progress"
}

type AncestralKiChamberRun struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	PlayerID        int64  `gorm:"not null;index"`
	ChamberLevel    int    `gorm:"not null"`
	CurrentBattle   int    `gorm:"not null;default:1"`
	Victories       int    `gorm:"not null;default:0"`
	AccumulatedXP   int    `gorm:"not null;default:0"`
	AccumulatedGold int    `gorm:"not null;default:0"`
	Status          string `gorm:"size:20;not null;default:ACTIVE"`

	StartedAt  time.Time
	FinishedAt *time.Time

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (AncestralKiChamberRun) TableName() string {
	return "ancestral_ki_chamber_runs"
}

type AncestralKiChamberBattleLog struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	ChamberRunID uint  `gorm:"not null;index"`
	BattleNumber int   `gorm:"not null"`
	MonsterID    int64 `gorm:"not null"`
	MonsterLevel int   `gorm:"not null"`
	Result       string

	CreatedAt time.Time `gorm:"not null"`
}

func (AncestralKiChamberBattleLog) TableName() string {
	return "ancestral_ki_chamber_battle_logs"
}

type AncestralKiChamberResolution struct {
	Run        *AncestralKiChamberRun
	Progress   *AncestralKiChamberProgress
	XPPayout   int
	GoldPayout int
}
