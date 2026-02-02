package ports

import "github.com/vigmiranda/dragonki/internal/domain"

type MonsterRepository interface {
	ListByLevel(userID int64) ([]domain.Monster, error)
	GetByMonsterID(monsterID int64) (*domain.Monster, error)
}
