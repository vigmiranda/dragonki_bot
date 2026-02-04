package application

import (
	"fmt"
	"log"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

type PlayerService struct {
	PlayerRepository ports.PlayerRepository
	RegionRepository ports.RegionRepository
}

func NewPlayerService(repo ports.PlayerRepository, regionRepo ports.RegionRepository) *PlayerService {
	return &PlayerService{PlayerRepository: repo, RegionRepository: regionRepo}
}
func (s *PlayerService) Status(userID int64) (view.GameResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}
	region, err := s.RegionRepository.GetByID(player.RegionID)
	if err != nil {
		log.Fatal(err)
	}
	if player.Energy < player.MaxEnergy || player.HP < player.MaxHP {
		player.RegenEnergy()
	}
	player.SetTransformBase()
	err = s.PlayerRepository.Update(player)
	if err != nil {
		log.Fatal(err)
	}

	genderIcon := "🧝‍♂️"
	if player.Gender == "WOMAN" {
		genderIcon = "🧝‍♀️️"
	}
	if player.Form == "Ki Awakened" {
		player = domain.KiAwakened(player)
	}

	msg := fmt.Sprintf(
		genderIcon+"<b>%s</b> - <b>Classe</b>: %s \n"+
			"🗺️ <b>Local</b>: %s - (Lv %d+)\n"+
			"💪 <b>Level</b>: %d   🧬 <b>XP</b>: %d\n"+
			"❤️ <b>HP</b>: %d/%d   ⚡️ <b>Energia</b>: %d/%d\n"+
			"⚔️ <b>ATK</b>: %d     🛡️ <b>DEF</b>: %d \n"+
			"💰 <b>Gold</b>: %d\n"+
			"💎 <b>Ki Crystals</b>: %d\n",
		player.Name,
		player.ClassName,
		region.Name,
		region.MinLevel,
		player.Level,
		player.XP,
		player.HP,
		player.MaxHP,
		player.Energy,
		player.MaxEnergy,
		player.ATK,
		player.DEF,
		player.Gold,
		player.KiCrystals,
	)
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/skins/%s_skin_basic.png", domain.ConvertLowCase(player.Gender)),
	}
	menu := "MAIN_MENU"

	viewResponse := view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: menu,
	}

	return viewResponse, nil
}
