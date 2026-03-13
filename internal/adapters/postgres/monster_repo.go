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
func (r *MonsterPostgresRepo) ListByRegion(regionId int64) ([]domain.Monster, error) {
	var monsters []domain.Monster

	err := r.db.Model(&domain.Monster{}).
		Where("region_id = ?", regionId).
		Scan(&monsters).Error
	if err != nil {
		return nil, err
	}

	return monsters, nil
}

func (r *MonsterPostgresRepo) ListAll() ([]domain.Monster, error) {
	var monsters []domain.Monster

	err := r.db.Model(&domain.Monster{}).
		Scan(&monsters).Error
	if err != nil {
		return nil, err
	}

	return monsters, nil
}

func (r *MonsterPostgresRepo) GetMonsterDrops(monsterID int) ([]domain.MonsterDrop, error) {
	var drops []domain.MonsterDrop
	err := r.db.Model(&domain.MonsterDrop{}).
		Where("monster_id = ?", monsterID).
		Preload("Item").
		Find(&drops).Error
	if err != nil {
		return nil, err
	}

	return drops, nil
}
