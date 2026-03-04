package application

import (
	"fmt"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

type MenuService struct {
	PlayerRepository ports.PlayerRepository
	RegionRepository ports.RegionRepository
}

func NewMenuService(repo ports.PlayerRepository, regionRepo ports.RegionRepository) *MenuService {
	return &MenuService{PlayerRepository: repo, RegionRepository: regionRepo}
}

func (s *MenuService) Menu(userID int64) (view.GameResponse, error) {
	var player, err = s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		return view.GameResponse{}, err
	}
	//var playersOnline, _ = s.PlayerRepository.CountOnlineLast15Min()
	region, err := s.RegionRepository.GetByID(player.RegionID)
	if err != nil {
		return view.GameResponse{}, err
	}

	if err := applyPassivePlayerUpdates(s.PlayerRepository, player); err != nil {
		return view.GameResponse{}, err
	}
	//multiPlayersOnlineDB := os.Getenv("MULT_PLAYERS")
	//multiPlayersOnline, _ := strconv.Atoi(multiPlayersOnlineDB)

	msg := fmt.Sprintf(
		"<b> %s </b> \n\n"+
			//"👥 <b>Online (últimos 15 min)</b>: %d\n"+
			"🏰 <b>Local</b>:  %s - (Lv %d+)\n"+
			"❤️ <b>HP</b>: %d/%d\n"+
			"⚡️ <b>Energia</b>: %d/%d\n"+
			"💪 <b>Level</b>: %d\n"+
			"💰 <b>Gold</b>: %d\n",
		player.Name,
		//playersOnline*2,
		region.Name,
		region.MinLevel,
		player.HP,
		player.MaxHP,
		player.Energy,
		player.MaxEnergy,
		player.Level,
		player.Gold,
	)
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
