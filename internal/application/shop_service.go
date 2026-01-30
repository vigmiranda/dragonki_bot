package application

import (
	"fmt"
	"log"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/ports"
)

type ShopService struct {
	PlayerRepository ports.PlayerRepository
}

func NewShopService(repo ports.PlayerRepository) *ShopService {
	return &ShopService{PlayerRepository: repo}
}

func (s *ShopService) OpenShop(userID int64) (view.GameResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}
	msg := fmt.Sprintf(
		"🪙 Kiara, a Mercadora 🪙 \n\n"+
			"Seja bem-vindo %s à minha lojinha! \n\n"+
			"Eu sou Kiara, mercadora viajante e guardiã de artefatos raros. \n\n"+
			"Aqui você encontrará poções, equipamentos e itens especiais que vão fortalecer seu corpo, "+
			"renovar seu Ki e te ajudar a sobreviver aos desafios que te esperam lá fora. \n\n"+
			"💎 Cada item aqui foi escolhido para quem busca evoluir, não apenas lutar. \n"+
			"Prepare-se bem… a jornada de um verdadeiro guerreiro (a) começa com as escolhas certas. \n\n"+
			"O que posso fazer por você hoje? 😄 \n\n",
		player.Name)

	viewResponse := view.GameResponse{
		Text: msg,
		Image: &view.Image{
			Type: view.ImageAsset,
			Path: fmt.Sprintf("assets/npcs/kiara_mercadora.png"),
		},
		MenuKey: "SHOP_MENU",
	}
	return viewResponse, nil
}
