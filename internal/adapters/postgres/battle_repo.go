package postgres

import (
	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
)

type BattlePostgresRepo struct {
	db *gorm.DB
}

func NewBattlePostgresRepo(db *gorm.DB) *BattlePostgresRepo {
	return &BattlePostgresRepo{
		db: db,
	}
}

func (r *BattlePostgresRepo) SaveBattleLog(battleLog *domain.BattleLog) error {
	return r.db.Save(battleLog).Error
}
