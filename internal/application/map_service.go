package application

import (
	"fmt"
	"log"
	"strconv"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

type MapService struct {
	PlayerRepository ports.PlayerRepository
	RegionRepository ports.RegionRepository
}

func NewMapService(repo ports.PlayerRepository, regionRepo ports.RegionRepository) *MapService {
	return &MapService{PlayerRepository: repo, RegionRepository: regionRepo}
}

func (s MapService) ListMap(playerID int64) (view.GameResponse, error) {
	var player, err = s.PlayerRepository.GetOrCreate(playerID)
	if err != nil {
		log.Fatal(err)
	}
	regions, err := s.RegionRepository.ListByLevel(int64(player.Level))
	if err != nil {
		log.Fatal(err)
	}

	msg := fmt.Sprintf(
		"<b>Word:</b>\n\n" +
			"Plains lvl 1-15 \n" +
			"Forest of Spirits lvl 15-34 \n" +
			"Volcano Lands lvl 35-50 \n" +
			"Frozen Peaks lvl 50-70 \n" +
			"Swamp of Shadows lvl 70+" +
			"\n\n<b>Lugares que você pode explorar!!!</b> \n\n" +
			"",
	)
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/maps/word.png"),
	}
	menu := "REGIONS_MENU"

	viewResponse := view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: menu,
		MenuParams: map[string]any{
			"Regions": regions,
		},
	}

	return viewResponse, nil
}

func (s MapService) ToRegion(playerID int64, params map[string]string) (view.GameResponse, error) {
	var player, err = s.PlayerRepository.GetOrCreate(playerID)
	if err != nil {
		log.Fatal(err)
	}
	regionIDStr, _ := params["regionID"]
	regionId, err := strconv.ParseInt(regionIDStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	player.RegionID = regionId
	err = s.PlayerRepository.Update(player)

	var region, _ = s.RegionRepository.GetByID(regionId)

	if err != nil {
		return view.GameResponse{}, err
	}

	msg := fmt.Sprintf(
		"\n\n <b>Viagem concluida!!!</b> \n\n",
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
