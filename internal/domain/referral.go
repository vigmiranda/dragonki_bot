package domain

import "time"

const (
	ReferralStatusPending  = "PENDING"
	ReferralStatusRewarded = "REWARDED"
)

type Referral struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	ReferrerPlayerID int64 `gorm:"not null;index"`
	ReferredPlayerID int64 `gorm:"not null;uniqueIndex"`

	Status        string `gorm:"size:20;not null;default:PENDING"`
	RequiredLevel int    `gorm:"not null;default:7"`

	CreatedAt   time.Time
	CompletedAt *time.Time
	RewardedAt  *time.Time
}

func (Referral) TableName() string {
	return "referrals"
}
