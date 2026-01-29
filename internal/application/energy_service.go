package application

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

type EnergyService struct {
	PlayerRepository ports.PlayerRepository
}

func NewEnergyService(repo ports.PlayerRepository) *EnergyService {
	return &EnergyService{PlayerRepository: repo}
}

func (s *EnergyService) Energy(userID int64) (view.MsgResponse, error) {
	var player, err = s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}
	if player.Energy < player.MaxEnergy || player.HP < player.MaxHP {
		player.RegenEnergy()
	}
	player.SetTransformBase()
	err = s.PlayerRepository.Save(player)
	if err != nil {
		log.Fatal(err)
	}

	var msgResponse view.MsgResponse
	msgResponse.Message =
		fmt.Sprintf(
			"🛏️ Sir Soneca, Guardião do Descanso 🛏️ \n\n"+
				"“Ooooh… hã? Um aventureiro?” 😴 \n"+
				"Bem-vindo %s à minha humilde pousada… \n\n"+
				"Sou Sir Soneca, cavaleiro aposentado e mestre supremo da arte de dormir bem. \n\n"+
				"Aqui você pode descansar, recuperar suas forças e deixar que o cansaço da batalha desapareça. \n"+
				"Uma boa noite de sono cura feridas, restaura o HP e prepara seu espírito para a próxima luta. \n\n"+
				"🕯️ Camas macias, silêncio garantido… e roncos inclusos no pacote. \n\n"+
				"Então, jovem guerreiro(a) … vai apenas descansar um pouco… \n"+
				"ou dormir como um verdadeiro herói? 😴💤 \n\n",
			player.Name)
	msgResponse.Img =
		tgbotapi.FilePath(fmt.Sprintf("assets/npcs/sir_soneca.png"))
	return msgResponse, nil
}
func (s *EnergyService) Rest(userID int64) (view.MsgResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}
	if player.Energy < player.MaxEnergy || player.HP < player.MaxHP {
		player.RegenEnergy()
	}
	player.SetTransformBase()
	err = s.PlayerRepository.Save(player)
	if err != nil {
		log.Fatal(err)
	}
	var msgResponse view.MsgResponse
	msgResponse = view.MsgResponse{
		Message: fmt.Sprintf(
			"Você nao tem energia\n"+
				"Aguarde %s!! e volter aqui para descansar!!\n",
			os.Getenv("ENERGY_REGEN_TIME"),
		),
		Img:        tgbotapi.FilePath(fmt.Sprintf("assets/npcs/sir_soneca.png")),
		Conditions: false,
	}

	if player.Energy >= 1 {
		player.Energy -= 1

		player.HP += int(float64(player.MaxHP) * 0.3)
		if player.HP > player.MaxHP {
			player.HP = player.MaxHP
		}

		msgResponse = view.MsgResponse{
			Message: fmt.Sprintf(
				"❤️ HP recarregado!!! %d/%d\n"+
					"Bora voltar pra Luta ⚔️!! \n",
				player.HP,
				player.MaxHP,
			),
			Img:        tgbotapi.FilePath(fmt.Sprintf("assets/skins/%s_skin_basic.png", domain.ConvertLowCase(player.Gender))),
			Conditions: true,
		}
		err := s.PlayerRepository.Save(player)
		if err != nil {
			log.Fatal(err)
		}
	}
	return msgResponse, nil

}
func (s *EnergyService) Sleep(userID int64) (view.MsgResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}
	if player.Energy < player.MaxEnergy || player.HP < player.MaxHP {
		player.RegenEnergy()
	}
	player.SetTransformBase()
	err = s.PlayerRepository.Save(player)
	if err != nil {
		log.Fatal(err)
	}

	var msgResponse view.MsgResponse
	msgResponse = view.MsgResponse{
		Message: fmt.Sprintf(
			"Infelizmente você não podera ficar \n" +
				"Junte mais gold e retorne!!\n",
		),
		Img:        tgbotapi.FilePath(fmt.Sprintf("assets/npcs/sir_soneca.png")),
		Conditions: false,
	}

	if player.Gold >= 500 {
		player.Gold = player.Gold - 500
		player.HP = player.MaxHP
		msgResponse = view.MsgResponse{
			Message: fmt.Sprintf(
				"❤️ HP recarregado!!! \n" +
					"Bora voltar pra Luta ⚔️!! \n",
			),
			Img:        tgbotapi.FilePath(fmt.Sprintf("assets/skins/%s_skin_basic.png", domain.ConvertLowCase(player.Gender))),
			Conditions: true,
		}
		err := s.PlayerRepository.Save(player)
		if err != nil {
			log.Fatal(err)
		}
	}

	return msgResponse, nil
}
