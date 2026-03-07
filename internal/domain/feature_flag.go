package domain

import "time"

type FeatureFlag struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"uniqueIndex;size:120;not null"`
	Enabled     bool      `gorm:"not null;default:false"`
	Description string    `gorm:"size:500"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}
