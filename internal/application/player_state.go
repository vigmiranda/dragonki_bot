package application

import (
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

func applyPassivePlayerUpdates(repo ports.PlayerRepository, player *domain.Player) error {
	if player.Energy < player.MaxEnergy || player.HP < player.MaxHP {
		player.RegenEnergy()
	}
	player.SetTransformBase()
	return repo.Update(player)
}
