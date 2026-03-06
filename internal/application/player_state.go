package application

import (
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

func applyPassivePlayerUpdates(repo ports.PlayerRepository, player *domain.Player, regenEvery time.Duration) error {
	if player.Energy < player.MaxEnergy || player.HP < player.MaxHP {
		player.RegenEnergy(regenEvery)
	}
	player.SetTransformBase()
	return repo.Update(player)
}
