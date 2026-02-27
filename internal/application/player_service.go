package application

import (
	"fmt"
	"log"
	"os"
	"strconv"

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

	// 🔹 fallback seguro
	recoveryStr := os.Getenv("ENERGY_REGEN_TIME")
	if recoveryStr == "" {
		recoveryStr = "10m"
	}

	skin := fmt.Sprintf("skin_%s", player.Class.Skin)
	genderIcon := "🧝‍♂️"
	if player.Gender == "WOMAN" {
		genderIcon = "🧝‍♀️️"
	}
	if player.Form == "Ki Awakened" {
		player = domain.KiAwakened(player)
		skin = "ki_awakened"
	}

	msg := fmt.Sprintf(
		genderIcon+"<b>%s</b> - <b>Classe</b>: %s \n"+
			"🗺️ <b>Local</b>: %s - (Lv %d+)\n"+
			"💪 <b>Level</b>: %d   🧬 <b>XP</b>: %d/%d\n"+
			"❤️ <b>HP</b>: %d/%d   ⚡️ <b>Energia</b>: %d/%d\n"+
			"⚔️ <b>ATK</b>: %d   🛡️ <b>DEF</b>: %d   🏃🏻 <b>VEL</b>: %d \n"+
			"💰 <b>Gold</b>: %d\n"+
			"💎 <b>Ki Crystals</b>: %d\n"+
			"<b>+1 de ⚡️ em </b>: %s\n",
		player.Name,
		player.Class.Name,
		region.Name,
		region.MinLevel,
		player.Level,
		player.XP,
		player.Level*80,
		player.HP,
		player.MaxHP,
		player.Energy,
		player.MaxEnergy,
		player.ATK,
		player.DEF,
		player.Speed,
		player.Gold,
		player.KiCrystals,
		recoveryStr,
	)
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/skins/%s_%s.png",
			domain.ConvertLowCase(player.Gender),
			skin),
	}
	menu := "MAIN_MENU"

	viewResponse := view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: menu,
	}

	return viewResponse, nil
}

func (s *PlayerService) UpdateClass(playerId int64, params map[string]string) (view.GameResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(playerId)
	if err != nil {
		log.Fatal(err)
	}
	region, err := s.RegionRepository.GetByID(player.RegionID)
	if err != nil {
		log.Fatal(err)
	}

	classIDStr, _ := params["classID"]
	classID, err := strconv.ParseInt(classIDStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	class, err := s.PlayerRepository.GetClassesByID(classID)

	player.ClassID = uint(classID)
	s.PlayerRepository.Update(player)

	msg := fmt.Sprintf("Agora voce é %s", class.Name)
	//fmt.Sprintf(
	//genderIcon+"<b>%s</b> - <b>Classe</b>: %s \n"+
	//	"🗺️ <b>Local</b>: %s - (Lv %d+)\n"+
	//	"💪 <b>Level</b>: %d   🧬 <b>XP</b>: %d/%d\n"+
	//	"❤️ <b>HP</b>: %d/%d   ⚡️ <b>Energia</b>: %d/%d\n"+
	//	"⚔️ <b>ATK</b>: %d     🛡️ <b>DEF</b>: %d \n"+
	//	"💰 <b>Gold</b>: %d\n"+
	//	"💎 <b>Ki Crystals</b>: %d\n"+
	//	"<b>+1 de ⚡️ em </b>: %s\n",
	//player.Name,
	//player.ClassName,
	//region.Name,
	//region.MinLevel,
	//player.Level,
	//player.XP,
	//player.Level*80,
	//player.HP,
	//player.MaxHP,
	//player.Energy,
	//player.MaxEnergy,
	//player.ATK,
	//player.DEF,
	//player.Gold,
	//player.KiCrystals,
	//recoveryStr,

	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/maps/%s.png", domain.ConvertLowCase(region.Tag)),
	}
	menu := "MAIN_MENU"

	viewResponse := view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: menu,
	}

	return viewResponse, nil
}
