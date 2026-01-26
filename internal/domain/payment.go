package domain

import "time"

type PaymentStatus string

const (
	PaymentPending  PaymentStatus = "pending"
	PaymentApproved PaymentStatus = "approved"
	PaymentFailed   PaymentStatus = "failed"
)

type Payment struct {
	ID           string
	PlayerID     int64
	Amount       float64
	Status       PaymentStatus
	QRCode       string `json:"qr_code"`
	QRCodeBase64 string `json:"qr_code_base64"`
	PartnerId    string
	CreatedAt    time.Time
}
