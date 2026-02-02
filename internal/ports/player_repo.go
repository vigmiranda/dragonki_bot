package ports

import "github.com/vigmiranda/dragonki/internal/domain"

type PlayerRepository interface {
	GetOrCreate(userID int64) (*domain.Player, error)
	Update(player *domain.Player) error
	CountOnlineLast15Min() (int, error)
}
