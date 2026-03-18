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
	vipService       *VIPService
	cfg              appConfig
}

func NewMenuService(repo ports.PlayerRepository, regionRepo ports.RegionRepository, vipService ...*VIPService) *MenuService {
	var v *VIPService
	if len(vipService) > 0 {
		v = vipService[0]
	}
	return &MenuService{PlayerRepository: repo, RegionRepository: regionRepo, vipService: v, cfg: getAppConfig()}
}

func (s *MenuService) Menu(userID int64) (view.GameResponse, error) {
	cfg := resolveAppConfig(s.cfg)

	var player, err = s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		return view.GameResponse{}, err
	}
	//var playersOnline, _ = s.PlayerRepository.CountOnlineLast15Min()
	region, err := s.RegionRepository.GetByID(player.RegionID)
	if err != nil {
		return view.GameResponse{}, err
	}

	if err := applyPassivePlayerUpdates(s.PlayerRepository, player, cfg.EnergyRegenDuration); err != nil {
		return view.GameResponse{}, err
	}
	//multiPlayersOnlineDB := os.Getenv("MULT_PLAYERS")
	//multiPlayersOnline, _ := strconv.Atoi(multiPlayersOnlineDB)
	playerName := displayPlayerName(player, s.vipService)

	msg := fmt.Sprintf(
		"<b> %s </b> \n\n"+
			//"👥 <b>Online (últimos 15 min)</b>: %d\n"+
			"🏰 <b>Local</b>:  %s - (Lv %d+)\n"+
			"❤️ <b>HP</b>: %d/%d\n"+
			"⚡️ <b>Energia</b>: %d/%d\n"+
			"💪 <b>Level</b>: %d\n"+
			"💰 <b>Gold</b>: %d\n",
		playerName,
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
