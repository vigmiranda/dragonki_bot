package ports

import (
	"github.com/vigmiranda/dragonki/internal/domain"
)

type VIPConfigRepository interface {
	GetEnabledConfig() (*domain.VIPConfig, error)
}

type PlayerVIPRepository interface {
	GetByPlayerID(playerID int64) (*domain.PlayerVIP, error)
	Upsert(vip *domain.PlayerVIP) error
}
