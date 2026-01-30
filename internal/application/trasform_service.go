package application

import (
	"fmt"
	"log"

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

func (s *TransformService) Transform(userID int64) (view.GameResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}

	//todo retirar daqui
	if player.Energy < player.MaxEnergy || player.HP < player.MaxHP {
		player.RegenEnergy()
	}
	err = s.PlayerRepository.Save(player)
	if err != nil {
		log.Fatal(err)
	}
	// ---------------------

	msg := fmt.Sprintf("❌ Transformação desbloqueia no nível 3")
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/skins/%s_skin_basic.png", domain.ConvertLowCase(player.Gender)),
	}
	menu := "MAIN_MENU"

	if player.Energy < 3 && player.FormBattlesLeft == 0 {
		msg = fmt.Sprintf("⚡ Energia insuficiente")
		img.Path = fmt.Sprintf("assets/skins/%s_tired.png", domain.ConvertLowCase(player.Gender))
	}
	if player.FormBattlesLeft > 0 {
		msg = fmt.Sprintf(
			"🔥 Você já ascendeu o <b>Ki Awakened</b>!\n"+
				"Restam ainda %d lutas",
			player.FormBattlesLeft,
		)
		img.Path = fmt.Sprintf("assets/skins/%s_ki_awakened.png", domain.ConvertLowCase(player.Gender))
	}

	if player.Level >= 3 && player.Level <= 15 && player.Energy >= 3 && player.FormBattlesLeft == 0 {
		player.Form = "Ki Awakened"
		player.FormBattlesLeft = 3
		player.Energy -= 3
		msg = fmt.Sprintf(
			"🔥 Você ascendeu o <b>Ki Awakened</b>! \n\n" +
				"Essa transformação te dá 10%% em ATK, DEF \n" +
				"Gasta 3 de energia e dura 3 batalhas \n" +
				"Ahhh! Caso você morra ou suba de nível.\n" +
				"Você irá retornar à forma base",
		)
		img.Path = fmt.Sprintf("assets/skins/%s_ki_awakened.png", domain.ConvertLowCase(player.Gender))
	}

	err = s.PlayerRepository.Save(player)
	if err != nil {
		log.Fatal(err)
	}

	viewResponse := view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: menu,
	}

	return viewResponse, nil
}
