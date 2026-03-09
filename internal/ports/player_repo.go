package ports

import "github.com/vigmiranda/dragonki/internal/domain"

type PlayerRepository interface {
	GetOrCreate(userID int64) (*domain.Player, error)
	Update(player *domain.Player) error
	CountOnlineLast15Min() (int, error)
	GetClassesByLvl(lvl int) ([]domain.Class, error)
	GetClassesByID(id int64) (*domain.Class, error)
	FindClassEvolutionByBaseClassAndLevel(baseClass string, lvl int) (*domain.Class, error)
}
