package application

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

type TransformService struct {
	PlayerRepository ports.PlayerRepository
}

func NewTransformService(repo ports.PlayerRepository) *TransformService {
	return &TransformService{PlayerRepository: repo}
}

func (s *TransformService) Transform(userID int64) (view.MsgResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}
	if player.Energy < player.MaxEnergy || player.HP < player.MaxHP {
		player.RegenEnergy()
	}
	err = s.PlayerRepository.Save(player)
	if err != nil {
		log.Fatal(err)
	}

	var msgResponse view.MsgResponse

	logText := ""

	if player.Level < 3 {
		logText += fmt.Sprintf("❌ Transformação desbloqueia no nível 3")
		msgResponse.Img = tgbotapi.FilePath(fmt.Sprintf("assets/skins/%s_skin_basic.png", domain.ConvertLowCase(player.Gender)))
	}
	if player.Energy < 3 && player.FormBattlesLeft == 0 {
		logText += fmt.Sprintf("⚡ Energia insuficiente")
		msgResponse.Img = tgbotapi.FilePath(fmt.Sprintf("assets/skins/%s_tired.png", domain.ConvertLowCase(player.Gender)))
	}
	if player.FormBattlesLeft > 0 {
		logText += fmt.Sprintf(
			"🔥 Você já ascendeu o <b>Ki Awakened</b>!\n"+
				"Restam ainda %d lutas",
			player.FormBattlesLeft,
		)
		msgResponse.Img = tgbotapi.FilePath(fmt.Sprintf("assets/skins/%s_ki_awakened.png", domain.ConvertLowCase(player.Gender)))
	}

	if player.Level >= 3 && player.Level <= 15 && player.Energy >= 3 && player.FormBattlesLeft == 0 {
		player.Form = "Ki Awakened"
		player.FormBattlesLeft = 3
		player.Energy -= 3
		logText += fmt.Sprintf(
			"🔥 Você ascendeu o <b>Ki Awakened</b>! \n\n" +
				"Essa transformação te dá 10%% em ATK, DEF \n" +
				"Gasta 3 de energia e dura 3 batalhas \n" +
				"Ahhh! Caso você morra ou suba de nível.\n" +
				"Você irá retornar à forma base",
		)
		msgResponse.Img = tgbotapi.FilePath(fmt.Sprintf("assets/skins/%s_ki_awakened.png", domain.ConvertLowCase(player.Gender)))
	}

	err = s.PlayerRepository.Save(player)
	if err != nil {
		log.Fatal(err)
	}
	msgResponse.Message = logText
	return msgResponse, nil
}
