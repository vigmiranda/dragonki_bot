package ports

import "github.com/vigmiranda/dragonki/internal/domain"

type MonsterRepository interface {
	ListByRegion(userID int64) ([]domain.Monster, error)
	ListAll() ([]domain.Monster, error)
	GetByMonsterID(monsterID int64) (*domain.Monster, error)
	GetMonsterDrops(monsterID int) ([]domain.MonsterDrop, error)
}
