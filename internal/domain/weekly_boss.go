package domain

import (
	"errors"
	"time"
)

const (
	WeeklyBossStatusActive   = "active"
	WeeklyBossStatusDefeated = "defeated"
	WeeklyBossStatusExpired  = "expired"
)

type WeeklyBoss struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	Name    string `gorm:"size:120;not null"`
	Tag     string `gorm:"size:64;not null"`
	WeekKey string `gorm:"size:16;uniqueIndex;not null"`

	MaxHP     int64 `gorm:"not null"`
	CurrentHP int64 `gorm:"not null"`
	ATK       int   `gorm:"not null;column:atk"`

	Status string `gorm:"size:20;not null;default:active"`

	StartsAt   time.Time
	EndsAt     time.Time
	DefeatedAt *time.Time

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (WeeklyBoss) TableName() string {
	return "weekly_bosses"
}

type WeeklyBossAttack struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	BossID   uint      `gorm:"not null;index:uq_weekly_boss_attacks_daily_seq,unique"`
	PlayerID int64     `gorm:"not null;index:uq_weekly_boss_attacks_daily_seq,unique"`
	AttackAt time.Time `gorm:"column:attack_date;type:date;not null;index:uq_weekly_boss_attacks_daily_seq,unique"`
	DailySeq int       `gorm:"not null;default:1;index:uq_weekly_boss_attacks_daily_seq,unique"`
	Damage   int64     `gorm:"not null"`

	CreatedAt time.Time `gorm:"not null"`
}

func (WeeklyBossAttack) TableName() string {
	return "weekly_boss_attacks"
}

var (
	ErrWeeklyBossAttackLimitReached   = errors.New("weekly boss attack limit reached")
	ErrWeeklyBossInsufficientResource = errors.New("weekly boss insufficient resources")
)

type WeeklyBossContribution struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	BossID   uint  `gorm:"not null;index:uq_weekly_boss_contributions,unique"`
	PlayerID int64 `gorm:"not null;index:uq_weekly_boss_contributions,unique"`

	TotalDamage  int64     `gorm:"not null;default:0"`
	LastAttackAt time.Time `gorm:"not null"`

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (WeeklyBossContribution) TableName() string {
	return "weekly_boss_contributions"
}

type WeeklyBossRankingEntry struct {
	PlayerID    int64
	PlayerName  string
	TotalDamage int64
}

type WeeklyBossTemplate struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	Name    string `gorm:"size:120;not null"`
	Tag     string `gorm:"size:64;uniqueIndex;not null"`
	MaxHP   int64  `gorm:"not null"`
	ATK     int    `gorm:"not null;column:atk"`
	Enabled bool   `gorm:"not null;default:true"`

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (WeeklyBossTemplate) TableName() string {
	return "weekly_boss_templates"
}

const (
	WeeklyBossPayoutPending = "pending"
	WeeklyBossPayoutPaid    = "paid"
	WeeklyBossPayoutFailed  = "failed"
)

type WeeklyBossReward struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	Position int   `gorm:"not null;uniqueIndex"`
	Gold     int64 `gorm:"not null;default:0"`
	Enabled  bool  `gorm:"not null;default:true"`

	Items []WeeklyBossRewardItem `gorm:"foreignKey:RewardID"`

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (WeeklyBossReward) TableName() string {
	return "weekly_boss_rewards"
}

type WeeklyBossRewardItem struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	RewardID uint `gorm:"not null;index"`
	ItemID   uint `gorm:"not null"`
	Quantity int  `gorm:"not null"`

	CreatedAt time.Time `gorm:"not null"`
}

func (WeeklyBossRewardItem) TableName() string {
	return "weekly_boss_reward_items"
}

type WeeklyBossPayout struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	BossID   uint  `gorm:"not null;index"`
	PlayerID int64 `gorm:"not null"`
	Position int   `gorm:"not null"`
	GoldPaid int64 `gorm:"not null;default:0"`

	Status       string `gorm:"size:20;not null;default:pending"`
	ErrorMessage string `gorm:"type:text;not null;default:''"`
	PaidAt       *time.Time

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (WeeklyBossPayout) TableName() string {
	return "weekly_boss_payouts"
}
