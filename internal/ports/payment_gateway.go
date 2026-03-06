package ports

import "github.com/vigmiranda/dragonki/internal/domain"

type PaymentGateway interface {
	CreatePix(playerID int64, amount float64) (*domain.Payment, string, error)
}
