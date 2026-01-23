package ports

import "github.com/vigmiranda/dragonki/internal/domain"

type RegionRepository interface {
	GetByID(regionID int64) (*domain.Region, error)
	ListAll() ([]domain.Region, error)
}
