package application

import (
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

type PlayerService struct {
	PlayerRepository    ports.PlayerRepository
	RegionRepository    ports.RegionRepository
	InventoryRepository ports.InventoryRepository
	cfg                 appConfig
}

func NewPlayerService(repo ports.PlayerRepository, regionRepo ports.RegionRepository, inventoryRepo ports.InventoryRepository) *PlayerService {
	return &PlayerService{
		PlayerRepository:    repo,
		RegionRepository:    regionRepo,
		InventoryRepository: inventoryRepo,
		cfg:                 getAppConfig(),
	}
}

func (s *PlayerService) Status(userID int64) (view.GameResponse, error) {
	cfg := resolveAppConfig(s.cfg)

	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		return view.GameResponse{}, err
	}
	region, err := s.RegionRepository.GetByID(player.RegionID)
	if err != nil {
		return view.GameResponse{}, err
	}
	if err := applyPassivePlayerUpdates(s.PlayerRepository, player, cfg.EnergyRegenDuration); err != nil {
		return view.GameResponse{}, err
	}

	stats, statsErr := calculateEffectivePlayerStats(player, s.InventoryRepository)
	if statsErr != nil {
		log.Printf("status equipment bonus unavailable | playerID=%d | err=%v", player.ID, statsErr)
		stats = effectivePlayerStats{
			ATK:       player.ATK,
			DEF:       player.DEF,
			Speed:     player.Speed,
			MaxHP:     player.MaxHP,
			MaxEnergy: player.MaxEnergy,
		}
	}

	currentHP := player.HP
	if currentHP > stats.MaxHP {
		currentHP = stats.MaxHP
	}
	currentEnergy := player.Energy
	if currentEnergy > stats.MaxEnergy {
		currentEnergy = stats.MaxEnergy
	}

	transform := "player"
	genderIcon := "🧝‍♂️"
	if player.Gender == "WOMAN" {
		genderIcon = "🧝‍♀️️"
	}
	if player.Form == "Ki Awakened" {
		atkBoost := int(float64(stats.ATK) * 0.2)
		defBoost := int(float64(stats.DEF) * 0.2)
		if atkBoost < 1 {
			atkBoost = 1
		}
		if defBoost < 1 {
			defBoost = 1
		}
		stats.ATK += atkBoost
		stats.DEF += defBoost
		transform = "ki_awakened"
	}

	equipmentBonusLine := fmt.Sprintf(
		"<b>Bonus Equip</b>: ATK %+d | DEF %+d | VEL %+d\n",
		stats.Bonus.ATK,
		stats.Bonus.DEF,
		stats.Bonus.Speed,
	)

	XPLvlUp := 45 * math.Pow(float64(player.Level), 1.65)
	msg := fmt.Sprintf(
		genderIcon+"<b>%s</b> - <b>Classe</b>: %s \n"+
			"🗺️ <b>Local</b>: %s - (Lv %d+)\n"+
			"💪 <b>Level</b>: %d   🧬 <b>XP</b>: %d/%d\n"+
			"❤️ <b>HP</b>: %d/%d   ⚡️ <b>Energia</b>: %d/%d\n"+
			"⚔️ <b>ATK</b>: %d   🛡️ <b>DEF</b>: %d   🏃🏻 <b>VEL</b>: %d \n"+
			"%s"+
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
		currentHP,
		stats.MaxHP,
		currentEnergy,
		stats.MaxEnergy,
		stats.ATK,
		stats.DEF,
		stats.Speed,
		equipmentBonusLine,
		player.Gold,
		player.KiCrystals,
		cfg.EnergyRegenTime,
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
			"Otima escolha!\n"+
			"Seu Ki agora segue um novo caminho.\n"+
			"Seus poderes comecaram a se alinhar com sua essencia.\n"+
			"Treine. Evolua. Supere seus limites.\n",
		class.Name,
	)

	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/skins/%s/%splayer.png",
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
