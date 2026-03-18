package application

import (
	"fmt"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

type EnergyService struct {
	PlayerRepository ports.PlayerRepository
	vipService       *VIPService
	cfg              appConfig
}

func NewEnergyService(repo ports.PlayerRepository, vipService ...*VIPService) *EnergyService {
	var v *VIPService
	if len(vipService) > 0 {
		v = vipService[0]
	}
	return &EnergyService{
		PlayerRepository: repo,
		vipService:       v,
		cfg:              getAppConfig(),
	}
}

func (s *EnergyService) Energy(userID int64) (view.GameResponse, error) {
	cfg := resolveAppConfig(s.cfg)

	var player, err = s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		return view.GameResponse{}, err
	}

	//todo remover daqui
	if err := applyPassivePlayerUpdates(s.PlayerRepository, player, cfg.EnergyRegenDuration); err != nil {
		return view.GameResponse{}, err
	}
	// --------
	playerName := displayPlayerName(player, s.vipService)

	msg := fmt.Sprintf(
		"🛏️ Sir Soneca, Guardião do Descanso 🛏️ \n\n"+
			"“Ooooh… hã? Um aventureiro(a)?” 😴 \n"+
			"Bem-vindo(a) %s à minha humilde pousada… \n\n"+
			"Sou Sir Soneca, cavaleiro aposentado e mestre supremo da arte de dormir bem. \n\n"+
			"Aqui você pode descansar, recuperar suas forças e deixar que o cansaço da batalha desapareça. \n\n"+
			"🕯️ Camas macias, silêncio garantido… e roncos inclusos no pacote. \n\n"+
			"Então, jovem guerreiro(a) … vai apenas descansar um pouco… \n"+
			"ou dormir como quem teve um dia de grandes batalhas? 😴💤 \n\n",
		playerName)
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/npcs/sir_soneca.png"),
	}
	menu := "ENERGY_MENU"

	viewResponse := view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: menu,
	}

	return viewResponse, nil
}
func (s *EnergyService) Rest(userID int64) (view.GameResponse, error) {
	cfg := resolveAppConfig(s.cfg)

	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		return view.GameResponse{}, err
	}

	//todo remover daqui
	if err := applyPassivePlayerUpdates(s.PlayerRepository, player, cfg.EnergyRegenDuration); err != nil {
		return view.GameResponse{}, err
	}
	// --------

	msg := fmt.Sprintf(
		"Você nao tem energia\n"+
			"Aguarde %s!! e volte aqui para descansar!!\n",
		cfg.EnergyRegenTime)
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/npcs/sir_soneca.png"),
	}
	menu := "ENERGY_MENU"

	if player.Energy >= 1 {
		player.Energy -= 1

		player.HP += int(float64(player.MaxHP) * 0.3)
		if player.HP > player.MaxHP {
			player.HP = player.MaxHP
		}

		msg = fmt.Sprintf(
			"❤️ HP recarregado!!! %d/%d\n"+
				"Bora voltar pra Luta ⚔️!! \n",
			player.HP,
			player.MaxHP,
		)

		img = &view.Image{
			Type: view.ImageAsset,
			Path: fmt.Sprintf("assets/skins/%s/%splayer.png", domain.ConvertLowCase(player.Gender), player.Class.Skin),
		}

		if err := s.PlayerRepository.Update(player); err != nil {
			return view.GameResponse{}, err
		}
	}

	viewResponse := view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: menu,
	}
	return viewResponse, nil

}
func (s *EnergyService) Sleep(userID int64) (view.GameResponse, error) {
	cfg := resolveAppConfig(s.cfg)

	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		return view.GameResponse{}, err
	}

	//todo remover daqui
	if err := applyPassivePlayerUpdates(s.PlayerRepository, player, cfg.EnergyRegenDuration); err != nil {
		return view.GameResponse{}, err
	}
	// --------

	msg := fmt.Sprintf(
		"Infelizmente você não podera ficar \n" +
			"Junte mais gold e retorne!!\n",
	)
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/npcs/sir_soneca.png"),
	}
	menu := "ENERGY_MENU"

	if player.Gold >= 500 {
		player.Gold = player.Gold - 500
		player.HP = player.MaxHP
		msg = fmt.Sprintf(
			"❤️ HP recarregado!!! \n" +
				"Bora voltar pra Luta ⚔️!! \n",
		)
		if err := s.PlayerRepository.Update(player); err != nil {
			return view.GameResponse{}, err
		}
	}
	viewResponse := view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: menu,
	}

	return viewResponse, nil
}
