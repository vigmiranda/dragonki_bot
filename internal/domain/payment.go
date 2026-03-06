package domain

import (
	"errors"
	"fmt"
	"time"
)

type PaymentStatus string

const (
	PaymentPending  PaymentStatus = "pending"
	PaymentApproved PaymentStatus = "approved"
	PaymentFailed   PaymentStatus = "failed"
)

var ErrPaymentNotFound = errors.New("payment not found")

type Payment struct {
	ID           string
	PlayerID     int64
	Amount       float64
	Status       PaymentStatus
	QRCode       string `json:"qr_code"`
	QRCodeBase64 string `json:"qr_code_base64"`
	PartnerId    int64
	CreatedAt    time.Time
	LastUpdate   time.Time
}

func BuildNumericPaymentID(userID int64) string {
	timestamp := time.Now().UnixMilli()
	return fmt.Sprintf("%d%d", userID, timestamp)
}
