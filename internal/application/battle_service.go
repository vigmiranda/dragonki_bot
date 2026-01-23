package application

import (
	"fmt"
	"math/rand"

	"github.com/vigmiranda/dragonki/internal/adapters/postgres"
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

type BattleService struct {
	playerRepo  ports.PlayerRepository
	monsterRepo ports.MonsterRepository
}

func NewBattleService(playerRepo ports.PlayerRepository, monsterRepo *postgres.MonsterPostgresRepo) *BattleService {
	return &BattleService{
		playerRepo:  playerRepo,
		monsterRepo: monsterRepo,
	}
}

func (s *BattleService) Battle(userID int64) (MsgResponse, error) {
	player, err := s.playerRepo.GetOrCreate(userID)
	if err != nil {
		return MsgResponse{}, err
	}

	player.RegenEnergy()

	if player.Energy <= 0 {
		return MsgResponse{
			Message: "⚡ Energia esgotada",
		}, nil
	}

	player.Energy--

	monsters, err := s.monsterRepo.ListByLevel(int64(player.Level))
	if err != nil || len(monsters) == 0 {
		return MsgResponse{}, fmt.Errorf("nenhum monstro encontrado")
	}

	monster := pickMonsterByWeight(monsters)
	if monster == nil {
		return MsgResponse{}, fmt.Errorf("falha ao sortear monstro")
	}

	msg := MsgResponse{
		Message: fmt.Sprintf("🐉 Um %s apareceu!\n", monster.Name),
		Img:     fmt.Sprintf("assets/monsters/map_%d/%s.png", monster.Level, monster.Tag),
	}

	monsterHP := monster.HP

	//logica de dano
	for monsterHP > 0 && player.HP > 0 {
		monsterHP -= player.ATK
		player.HP -= rand.Intn(monster.ATK/2+1) + monster.ATK/2
	}

	if player.HP > 0 {
		xpReceived := rand.Intn(monster.MaxXP)
		moneyReceived := rand.Intn(monster.MaxMoney)

		player.XP += xpReceived
		player.Gold += moneyReceived

		msg.Message += fmt.Sprintf(
			"🏆 Vitória! +%d XP\n💰 +%d Gold\n",
			xpReceived,
			moneyReceived,
		)

		if player.XP >= player.Level*50 {
			player.Level++
			player.ATK += 5
			player.HP = 100
			msg.Message += "⬆️ Subiu de nível!\n"
		}

		if player.FormBattlesLeft > 0 {
			player.FormBattlesLeft--
		}

	} else {
		player.HP = 100
		msg.Message += "💀 Derrota! HP restaurado.\n"
	}

	if player.FormBattlesLeft == 0 {
		player.Form = "Base"
	}

	if err := s.playerRepo.Save(player); err != nil {
		return MsgResponse{}, err
	}

	return msg, nil
}

func pickMonsterByWeight(monsters []domain.Monster) *domain.Monster {
	totalWeight := 0
	for _, m := range monsters {
		if m.SpawnWeight > 0 {
			totalWeight += m.SpawnWeight
		}
	}

	if totalWeight == 0 {
		return nil
	}

	roll := rand.Intn(totalWeight)
	acc := 0

	for _, m := range monsters {
		acc += m.SpawnWeight
		if roll < acc {
			return &m
		}
	}

	return nil
}
