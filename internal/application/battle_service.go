package application

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/vigmiranda/dragonki/internal/adapters/postgres"
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

type BattleService struct {
	playerRepo  ports.PlayerRepository
	monsterRepo ports.MonsterRepository
	regionRepo  ports.RegionRepository
}

func NewBattleService(playerRepo ports.PlayerRepository, monsterRepo *postgres.MonsterPostgresRepo,
	regionRepo *postgres.RegionPostgres) *BattleService {
	return &BattleService{
		playerRepo:  playerRepo,
		monsterRepo: monsterRepo,
		regionRepo:  regionRepo,
	}
}

func (s *BattleService) GetBattle(userID int64) (MsgResponse, error) {
	player, err := s.playerRepo.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
		return MsgResponse{
			Message: "Sinto muito, erro ao iniciar Batalha!",
			Img:     "assets/error.png",
		}, err
	}

	//if player.Energy < player.MaxEnergy {player.RegenEnergy()}

	msg := MsgResponse{
		Message:    "⚡ Energia esgotada \n Vá ate a pousada para descansar!",
		Img:        fmt.Sprintf("assets/skins/%s_tired.png", domain.ConvertLowCase(player.Gender)),
		Conditions: true,
	}

	if player.Energy > 0 {
		monsters, err := s.monsterRepo.ListByLevel(int64(player.Level))
		if err != nil || len(monsters) == 0 {
			return MsgResponse{}, fmt.Errorf("nenhum monstro encontrado")
		}

		monster := pickMonsterByWeight(monsters)
		if monster == nil {
			return MsgResponse{}, fmt.Errorf("falha ao sortear monstro")
		}

		msg = MsgResponse{
			Message: fmt.Sprintf(
				"🐉<b> Um %s apareceu! </b>\n\n"+
					"💪 <b>Lv:</b> %d\n"+
					"❤️ <b>HP:</b> %d\n"+
					"⚔️ <b>ATK:</b> %d  🛡️<b>DEF:</> %d\n",
				monster.Name,
				monster.Level,
				monster.HP,
				monster.ATK,
				monster.DEF,
			),
			Img:        fmt.Sprintf("assets/monsters/map_%d/%s.png", monster.Level, monster.Tag),
			Conditions: false,
		}
	}

	return msg, nil
}
func (s *BattleService) Battle(userID int64) (MsgResponse, error) {
	//player, _ := s.playerRepo.GetOrCreate(userID)
	//////////////////////separar daqui pra baixo
	//monsterHP := monster.HP
	//playerHp := player.HP
	//player.Energy--
	//
	////logica de dano
	//for monsterHP > 0 && player.HP > 0 {
	//	monsterHP -= player.ATK
	//	player.HP -= rand.Intn(monster.ATK/2+1) + monster.ATK/2
	//}
	//
	//if player.HP > 0 {
	//	xpReceived := rand.Intn(monster.MaxXP) + 1
	//	moneyReceived := rand.Intn(monster.MaxMoney) + 1
	//	googNews := ""
	//
	//	player.XP += xpReceived
	//	player.Gold += moneyReceived
	//
	//	if player.XP >= player.Level*50 {
	//		player.Level++
	//		player.ATK += int(float64(player.ATK) * 0.2)
	//		player.DEF += int(float64(player.DEF) * 0.2)
	//		player.MaxHP += int(float64(player.ATK) * 0.1)
	//		player.Energy = player.MaxEnergy
	//		player.XP = 0
	//		player.FormBattlesLeft = 0
	//		player.Form = "Base"
	//		player.HP = player.MaxHP
	//		googNews = "⬆️ Subiu de nível!\n"
	//
	//	}
	//	if player.FormBattlesLeft > 0 {
	//		player.FormBattlesLeft--
	//	}
	//
	//	msg.Message = fmt.Sprintf(
	//		"🏆 Vitória! %s\n"+
	//			"%s"+
	//			"Voce recebeu de dano: %d\n"+
	//			"⚔️ +%d XP\n"+
	//			"💰 +%d Gold\n"+
	//			"❤️ HP: %d/%d\n"+
	//			"⚡️ Energia: %d/%d\n"+
	//			"💪 Level: %d\n",
	//		player.Name,
	//		googNews,
	//		playerHp-player.HP,
	//		xpReceived,
	//		moneyReceived,
	//		player.HP,
	//		player.MaxHP,
	//		player.Energy,
	//		player.MaxEnergy,
	//		player.Level,
	//	)
	//	msg.Img = fmt.Sprintf("assets/skins/%s_battle_win.png", domain.ConvertLowCase(player.Gender))
	//
	//} else {
	//	if player.FormBattlesLeft == 0 {
	//		player.Form = "Base"
	//	}
	//
	//	player.XP -= int(float64(player.XP) * 0.9)
	//	player.Gold -= int(float64(player.Gold) * 0.9)
	//	player.HP = 1
	//	player.Form = "Base"
	//	msg.Message = fmt.Sprintf(
	//		"💀 Derrota! %s\n"+
	//			"⚔️ -10%% de XP = %d\n"+
	//			"💰 -10%% de Gold = %d\n"+
	//			"⚡️ Energia: %d/%d\n"+
	//			"💪 Level: %d\n",
	//		player.Name,
	//		player.XP,
	//		player.Gold,
	//		player.Energy,
	//		player.MaxEnergy,
	//		player.Level,
	//	)
	//	msg.Img = fmt.Sprintf("assets/skins/%s_battle_loss.png", domain.ConvertLowCase(player.Gender))
	//}
	//if err := s.playerRepo.Save(player); err != nil {
	//	return MsgResponse{
	//		Message: "💀 Ocorreu um erro ao fim da batalha!!!",
	//		Img:     "assets/error.png",
	//	}, err
	//}
	return MsgResponse{
		Message: "💀 Ocorreu um erro ao fim da batalha!!!",
		Img:     "assets/error.png",
	}, nil
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
