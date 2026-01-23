package application

import (
	"fmt"
	"math/rand"

	"github.com/vigmiranda/dragonki/internal/ports"
)

type PlayerService struct {
	Repo ports.PlayerRepository
}

func NewPlayerService(repo ports.PlayerRepository) *PlayerService {
	return &PlayerService{Repo: repo}
}

func (s *PlayerService) Status(userID int64) (string, error) {
	player, _ := s.Repo.GetOrCreate(userID)
	player.RegenEnergy()
	s.Repo.Save(player)
	logText := fmt.Sprintf(
		"🧙 Status\n❤️ HP: %d\n⚔️ ATK: %d\n✨ XP: %d\n🏆 Nível: %d\n⚡ Energia: %d",
		player.HP, player.ATK, player.XP, player.Level, player.Energy,
	)
	return logText, nil
}

func (s *PlayerService) Battle(userID int64) (string, error) {
	player, _ := s.Repo.GetOrCreate(userID)

	player.RegenEnergy()

	if player.Energy <= 0 {
		return "⚡ Energia esgotada", nil
	}

	player.Energy--
	monsterHP := rand.Intn(40) + 30
	logText := "🐉 Um monstro apareceu!\n"

	for monsterHP > 0 && player.HP > 0 {
		monsterHP -= player.ATK
		player.HP -= rand.Intn(10) + 5
	}

	if player.HP > 0 {
		xp := rand.Intn(20) + 10
		player.XP += xp
		logText += fmt.Sprintf("🏆 Vitória! +%d XP\n", xp)

		if player.XP >= player.Level*50 {
			player.Level++
			player.ATK += 5
			player.HP = 100
			logText += "⬆️ Subiu de nível!"
		}
	} else {
		player.HP = 100
		logText += "💀 Derrota! HP restaurado."
	}

	s.Repo.Save(player)
	return logText, nil
}
