package ports

import "context"
import "github.com/vigmiranda/dragonki/internal/domain"

type PaymentGateway interface {
	CreatePix(ctx context.Context, playerID int64, amount float64) (*domain.Payment, string, error)
}
