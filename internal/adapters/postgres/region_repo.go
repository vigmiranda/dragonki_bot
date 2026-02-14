package postgres

import (
	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
)

type RegionPostgres struct {
	db *gorm.DB
}

func NewRegionPostgresRepo(db *gorm.DB) *RegionPostgres {
	return &RegionPostgres{
		db: db,
	}
}

func (r *RegionPostgres) GetByID(regionID int64) (*domain.Region, error) {
	var region domain.Region
	err := r.db.Where("region_id = ?", regionID).First(&region).Error
	if err != nil {
		return nil, err
	}
	return &region, nil
}

func (r *RegionPostgres) ListByLevel(level int64) ([]domain.Region, error) {
	var regions []domain.Region
	err := r.db.
		Where("min_level <= ?", level).
		Order("min_level ASC").
		Find(&regions).
		Error
	if err != nil {
		return nil, err
	}
	return regions, nil
}
