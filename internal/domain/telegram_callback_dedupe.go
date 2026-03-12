package domain

import "time"

type TelegramCallbackDedupe struct {
	DedupeKey   string    `gorm:"column:dedupe_key;primaryKey"`
	Action      string    `gorm:"column:action;not null"`
	WindowStart time.Time `gorm:"column:window_start;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (TelegramCallbackDedupe) TableName() string {
	return "telegram_callback_dedupes"
}
