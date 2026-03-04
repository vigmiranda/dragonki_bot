package application

import (
	"fmt"
	"math"
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
		return view.GameResponse{}, err
	}
	region, err := s.RegionRepository.GetByID(player.RegionID)
	if err != nil {
		return view.GameResponse{}, err
	}
	if err := applyPassivePlayerUpdates(s.PlayerRepository, player); err != nil {
		return view.GameResponse{}, err
	}

	// 🔹 fallback seguro
	recoveryStr := os.Getenv("ENERGY_REGEN_TIME")
	if recoveryStr == "" {
		recoveryStr = "10m"
	}

	transform := "player"
	genderIcon := "🧝‍♂️"
	if player.Gender == "WOMAN" {
		genderIcon = "🧝‍♀️️"
	}
	if player.Form == "Ki Awakened" {
		player = domain.KiAwakened(player)
		transform = "ki_awakened"
	}
	XPLvlUp := 45 * math.Pow(float64(player.Level), 1.65)
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
		int(XPLvlUp),
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
		Path: fmt.Sprintf("assets/skins/%s/%s%s.png",
			domain.ConvertLowCase(player.Gender),
			player.Class.Skin,
			transform),
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
		return view.GameResponse{}, err
	}

	switch player.Level {
	case 15:
		player.ATK = 110
		player.DEF = 50
		player.Speed = 50
		player.MaxHP = 160
		player.MaxEnergy = 23
	}

	classIDStr, _ := params["classID"]
	classID, err := strconv.ParseInt(classIDStr, 10, 64)
	if err != nil {
		return view.GameResponse{}, err
	}
	class, err := s.PlayerRepository.GetClassesByID(classID)
	if err != nil {
		return view.GameResponse{}, err
	}

	player.ClassID = uint(classID)
	if err := s.PlayerRepository.Update(player); err != nil {
		return view.GameResponse{}, err
	}

	msg := fmt.Sprintf(
		"<b>%s</b>\n\n"+
			"🔥 Ótima escolha!\n"+
			"Seu Ki agora segue um novo caminho.\n"+
			"Seus poderes começaram a se alinhar com sua essência.\n"+
			"Treine. Evolua. Supere seus limites.\n",
		class.Name,
	)

	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/skins/%s/%s/player.png",
			domain.ConvertLowCase(player.Gender),
			player.Class.Skin),
	}
	menu := "MAIN_MENU"

	viewResponse := view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: menu,
	}

	return viewResponse, nil
}
