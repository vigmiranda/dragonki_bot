package postgres

import (
	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
)

type MonsterPostgresRepo struct {
	db *gorm.DB
}

func NewMonsterPostgresRepo(db *gorm.DB) *MonsterPostgresRepo {
	return &MonsterPostgresRepo{
		db: db,
	}
}
func (r *MonsterPostgresRepo) GetByMonsterID(monsterId int64) (*domain.Monster, error) {
	var monster domain.Monster

	err := r.db.Model(&domain.Monster{}).
		Where("monster_id = ?", monsterId).
		First(&monster).Error
	if err != nil {
		return nil, err
	}
	return &monster, nil
}
func (r *MonsterPostgresRepo) ListByLevel(level int64) ([]domain.Monster, error) {
	var monsters []domain.Monster

	err := r.db.Model(&domain.Monster{}).
		Where("level <= ?", level).
		Scan(&monsters).Error
	if err != nil {
		return nil, err
	}

	return monsters, nil
}
