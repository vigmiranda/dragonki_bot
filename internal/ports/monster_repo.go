package ports

import "github.com/vigmiranda/dragonki/internal/domain"

type MonsterRepository interface {
	ListByRegionAndPlayerLvl(userID int64, lvl int) ([]domain.Monster, error)
	GetByMonsterID(monsterID int64) (*domain.Monster, error)
}
