package sqlite

import (
	"database/sql"

	"github.com/vigmiranda/dragonki/internal/domain"
)

type RegionSQLiteRepo struct {
	db *sql.DB
}

func NewRegionSQLiteRepo(sqlite *SQLite) *RegionSQLiteRepo {
	return &RegionSQLiteRepo{
		db: sqlite.DB,
	}
}

func (r *RegionSQLiteRepo) GetByID(regionID int64) (*domain.Region, error) {
	var region domain.Region

	err := r.db.QueryRow(`
		SELECT region_id, name, tag, min_level
		FROM regions
		WHERE region_id = ?
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

func (r *RegionSQLiteRepo) ListAll() ([]domain.Region, error) {
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
