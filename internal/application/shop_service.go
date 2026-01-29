package application

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/ports"
)

type ShopService struct {
	PlayerRepository ports.PlayerRepository
}

func NewShopService(repo ports.PlayerRepository) *ShopService {
	return &ShopService{PlayerRepository: repo}
}
func (s *ShopService) OpenShop(userID int64) (view.MsgResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}

	var msgResponse view.MsgResponse
	msgResponse = view.MsgResponse{
		Message: fmt.Sprintf(
			"🪙 Kiara, a Mercadora 🪙 \n\n"+
				"Seja bem-vindo %s à minha lojinha! \n\n"+
				"Eu sou Kiara, mercadora viajante e guardiã de artefatos raros. \n\n"+
				"Aqui você encontrará poções, equipamentos e itens especiais que vão fortalecer seu corpo, "+
				"renovar seu Ki e te ajudar a sobreviver aos desafios que te esperam lá fora. \n\n"+
				"💎 Cada item aqui foi escolhido para quem busca evoluir, não apenas lutar. \n"+
				"Prepare-se bem… a jornada de um verdadeiro guerreiro (a) começa com as escolhas certas. \n\n"+
				"O que posso fazer por você hoje? 😄 \n\n",
			player.Name,
		),
		Img:        tgbotapi.FilePath(fmt.Sprintf("assets/npcs/kiara_mercadora.png")),
		Conditions: false,
	}
	return msgResponse, nil
}
