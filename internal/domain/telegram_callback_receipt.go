package domain

import "time"

type TelegramCallbackReceipt struct {
	CallbackID string    `gorm:"column:callback_id;primaryKey"`
	UserID     int64     `gorm:"column:user_id;not null"`
	ChatID     int64     `gorm:"column:chat_id;not null"`
	MessageID  int       `gorm:"column:message_id;not null"`
	Action     string    `gorm:"column:action;not null"`
	RawData    string    `gorm:"column:raw_data;not null"`
	ReceivedAt time.Time `gorm:"column:received_at;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (TelegramCallbackReceipt) TableName() string {
	return "telegram_callback_receipts"
}
