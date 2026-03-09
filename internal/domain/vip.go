package domain

import "time"

type VIPConfig struct {
	ID                      uint      `gorm:"primaryKey;autoIncrement"`
	Enabled                 bool      `gorm:"not null;default:false"`
	VIPPriceKiCrystals      int       `gorm:"not null;default:10"`
	DurationDays            int       `gorm:"not null;default:7"`
	XPBonusPercent          float64   `gorm:"not null;default:0"`
	GoldBonusPercent        float64   `gorm:"not null;default:0"`
	SellBonusPercent        float64   `gorm:"not null;default:0"`
	ShopDiscountPercent     float64   `gorm:"not null;default:0"`
	BigBossDailyAttackLimit int       `gorm:"not null;default:1"`
	CreatedAt               time.Time `gorm:"not null"`
	UpdatedAt               time.Time `gorm:"not null"`
}

func (VIPConfig) TableName() string {
	return "vip_config"
}

type PlayerVIP struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	PlayerID  int64     `gorm:"not null;uniqueIndex"`
	StartedAt time.Time `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (PlayerVIP) TableName() string {
	return "player_vip"
}
