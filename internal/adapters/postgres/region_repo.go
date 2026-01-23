package postgres

import (
	"database/sql"

	"github.com/vigmiranda/dragonki/internal/domain"
)

type RegionPostgres struct {
	db *sql.DB
}

func NewRegionPostgresRepo(postgres *Postgres) *RegionPostgres {
	return &RegionPostgres{
		db: postgres.DB,
	}
}

func (r *RegionPostgres) GetByID(regionID int64) (*domain.Region, error) {
	var region domain.Region

	err := r.db.QueryRow(`
		SELECT region_id, name, tag, min_level
		FROM regions
		WHERE region_id = $1
	`, regionID).Scan(
		&region.RegionID,
		&region.Name,
		&region.Tag,
		&region.MinLevel,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &region, nil
}

func (r *RegionPostgres) ListAll() ([]domain.Region, error) {
	rows, err := r.db.Query(`
		SELECT region_id, name, tag, min_level
		FROM regions
		ORDER BY region_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var regions []domain.Region

	for rows.Next() {
		var region domain.Region
		err := rows.Scan(
			&region.RegionID,
			&region.Name,
			&region.Tag,
			&region.MinLevel,
		)
		if err != nil {
			return nil, err
		}
		regions = append(regions, region)
	}

	return regions, nil
}
