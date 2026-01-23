package application

import (
	"fmt"
	"strings"

	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

type MsgResponse struct {
	Message string
	Img     string
}

type PlayerService struct {
	PlayerRepository ports.PlayerRepository
	RegionRepository ports.RegionRepository
}

func NewPlayerService(repo ports.PlayerRepository, regionRepo ports.RegionRepository) *PlayerService {
	return &PlayerService{PlayerRepository: repo, RegionRepository: regionRepo}
}

func (s *PlayerService) SetGender(userID int64, gender string) (MsgResponse, error) {
	player, _ := s.PlayerRepository.GetOrCreate(userID)

	var msgResponse MsgResponse

	player.Gender = gender
	player.OnboardingStep = "CHOOSING_NAME"
	player.RegenEnergy()
	s.PlayerRepository.Save(player)
	msgResponse.Message = fmt.Sprintf(
		"🧙 Status\n❤️ HP: %d\n⚔️ ATK: %d\n✨ XP: %d\n🏆 Nível: %d\n⚡ Energia: %d",
		player.HP, player.ATK, player.XP, player.Level, player.Energy,
	)
	msgResponse.Img = "assets/skins/" + strings.ToLower(gender) + "_skin_basic.png"

	return msgResponse, nil
}
func (s *PlayerService) SetName(userID int64, name string) error {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		return err
	}

	player.Name = name
	player.OnboardingStep = "DONE"

	return s.PlayerRepository.Save(player)
}

func (s *PlayerService) Menu(userID int64) (MsgResponse, error) {
	player, _ := s.PlayerRepository.GetOrCreate(userID)
	region, err := s.RegionRepository.GetByID(player.RegionID)
	var msgResponse MsgResponse

	if err != nil {
		msgResponse.Message = fmt.Sprintf("Voce está perdido!")
		msgResponse.Img = fmt.Sprintf("assets/error.png")
		return msgResponse, err
	}

	player.RegenEnergy()
	msgResponse.Message = fmt.Sprintf(
		"%s você está em: %s (L%d)\n"+
			"❤️ HP: %d/%d\n"+
			"⚡️ Energia: %d/%d\n"+
			"💪 Level: %d\n",
		player.Name,
		region.Name,
		region.MinLevel,
		player.HP,
		player.MaxHP,
		player.Energy,
		player.MaxEnergy,
		player.Level,
	)
	msgResponse.Img = fmt.Sprintf("assets/maps/%s.png", strings.ToLower(region.Tag))
	return msgResponse, nil
}

func (s *PlayerService) HasCompletedOnboarding(userID int64) (bool, error) {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		return false, err
	}

	return player.OnboardingStep == "DONE", nil
}

func (s *PlayerService) Status(userID int64) (MsgResponse, error) {
	player, _ := s.PlayerRepository.GetOrCreate(userID)

	var msgResponse MsgResponse

	player.RegenEnergy()
	s.PlayerRepository.Save(player)
	msgResponse.Message = fmt.Sprintf(
		"%s\n"+
			"💪 Level: %d\n"+
			"❤️ HP: %d/%d\n"+
			"⚡️ Energia: %d/%d\n"+
			"💰 Gold: %d\n"+
			"💎 KiCrystals: %d\n",
		player.Name,
		player.Level,
		player.HP,
		player.MaxHP,
		player.Energy,
		player.MaxEnergy,
		player.Gold,
		player.KiCrystals,
	)
	msgResponse.Img = "assets/skins/" + strings.ToLower(player.Gender) + "_skin_basic.png"
	if player.FormBattlesLeft > 0 {
		msgResponse.Img = "assets/skins/" + strings.ToLower(player.Gender) + "_ki_awakened.png"
	}

	return msgResponse, nil
}

func (s *PlayerService) Transform(userID int64) (MsgResponse, error) {
	player, _ := s.PlayerRepository.GetOrCreate(userID)

	var msgResponse MsgResponse

	logText := ""

	if player.Level < 3 {
		logText += fmt.Sprintf("❌ Transformação desbloqueia no nível 3")
		msgResponse.Img = "assets/skins/" + strings.ToLower(player.Gender) + "_skin_basic.png"
	}
	if player.Energy < 3 {
		logText += fmt.Sprintf("⚡ Energia insuficiente")
		msgResponse.Img = "assets/skins/" + strings.ToLower(player.Gender) + "_skin_basic.png"
	}

	if player.Level >= 3 && player.Level <= 15 && player.Energy >= 3 {
		player.Form = "Ki Awakened"
		player.FormBattlesLeft = 3
		player.Energy -= 3
		logText += fmt.Sprintf("🔥 Você despertou o *Ki Awakened*!")
		msgResponse.Img = "assets/skins/" + strings.ToLower(player.Gender) + "_ki_awakened.png"
	}
	player.RegenEnergy()
	s.PlayerRepository.Save(player)
	msgResponse.Message = logText
	return msgResponse, nil
}

func (s *PlayerService) Energy(userID int64) (MsgResponse, error) {
	player, _ := s.PlayerRepository.GetOrCreate(userID)

	var msgResponse MsgResponse

	player.RegenEnergy()
	msgResponse.Message = fmt.Sprintf("📜 %d - %s \n"+
		"Você ainda tem energia para mais umas batalhas.\n"+
		"⚡️ Energia: %d\n", player.Level, player.Name, player.Energy)
	msgResponse.Img = fmt.Sprintf("assets/skins/%s_%s.png", strings.ToLower(player.Gender), "skin_basic")

	s.PlayerRepository.Save(player)

	if player.Energy <= 2 {
		msgResponse.Message = fmt.Sprintf(
			"%s sua energia esta baixa 😫. \n "+
				"⚡️ Energia: %d\n", player.Name, player.Energy)
		msgResponse.Img = fmt.Sprintf("assets/skins/%s_%s.png", strings.ToLower(player.Gender), "tired")
	}
	return msgResponse, nil
}
func (s *PlayerService) Onboarding(userID int64) (*domain.Player, error) {
	return s.PlayerRepository.GetOrCreate(userID)
}
