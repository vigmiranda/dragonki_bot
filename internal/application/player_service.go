package application

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

type MsgResponse struct {
	Message    string
	Img        string
	Menu       tgbotapi.InlineKeyboardMarkup
	Conditions bool
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
	//player.RegenEnergy()

	var msgResponse MsgResponse
	player.Gender = gender
	player.OnboardingStep = "CHOOSING_NAME"

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
	//player.RegenEnergy()

	var msgResponse MsgResponse

	if err != nil {
		msgResponse.Message = fmt.Sprintf("Voce está perdido!")
		msgResponse.Img = fmt.Sprintf("assets/error.png")
		return msgResponse, err
	}

	msgResponse.Message = fmt.Sprintf(
		"<b>%s</b> \n\n"+
			"🏰 <b>Local</b>:  %s (L%d)\n"+
			"❤️ <b>HP</b>: %d/%d\n"+
			"⚡️ <b>Energia</b>: %d/%d\n"+
			"💪 <b>Level</b>: %d\n"+
			"💰 <b>Gold</b>: %d\n",
		player.Name,
		region.Name,
		region.MinLevel,
		player.HP,
		player.MaxHP,
		player.Energy,
		player.MaxEnergy,
		player.Level,
		player.Gold,
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
	region, _ := s.RegionRepository.GetByID(player.RegionID)

	//player.RegenEnergy()

	genderIcon := "🧝‍♂️"
	if player.Gender == "WOMAN" {
		genderIcon = "🧝‍♀️️"
	}
	if player.Form == "Ki Awakened" {
		player = domain.KiAwakened(player)
	}
	var msgResponse MsgResponse
	msgResponse.Message = fmt.Sprintf(
		genderIcon+"<b>%s</b> - <b>Classe</b>: %s \n"+
			"💪 <b>Level</b>: %d\t\t\t\t\t\t\t\t\t\t\t🧬 <b>XP</b>: %d\n"+
			"❤️ <b>HP</b>: %d/%d\t\t\t\t\t\t⚡️ <b>Energia</b>: %d/%d\n"+
			"⚔️ <b>ATK</b>: %d \t\t\t\t\t\t\t\t\t\t🛡️ <b>DEF</b>: %d \n"+
			"🗺️ <b>Local</b>: %s(Lv%d)\n"+
			"💰 <b>Gold</b>: %d\n"+
			"💎 <b>KiCrystals</b>: %d\n",
		player.Name,
		player.Class,
		player.Level,
		player.XP,
		player.HP,
		player.MaxHP,
		player.Energy,
		player.MaxEnergy,
		player.ATK,
		player.DEF,
		region.Name,
		region.MinLevel,
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
	//player.RegenEnergy()

	var msgResponse MsgResponse

	logText := ""

	if player.Level < 3 {
		logText += fmt.Sprintf("❌ Transformação desbloqueia no nível 3")
		msgResponse.Img = "assets/skins/" + strings.ToLower(player.Gender) + "_skin_basic.png"
	}
	if player.Energy < 3 {
		logText += fmt.Sprintf("⚡ Energia insuficiente")
		msgResponse.Img = "assets/skins/" + strings.ToLower(player.Gender) + "_tired.png"
	}

	if player.Level >= 3 && player.Level <= 15 && player.Energy >= 3 {
		player.Form = "Ki Awakened"
		player.FormBattlesLeft = 3
		player.Energy -= 3
		logText += fmt.Sprintf(
			"🔥 Você ascendeu o <b>Ki Awakened</b>! \n\n" +
				"Essa transformação te dá 10%% em ATK, DEF \n" +
				"Gasta 3 de energia e dura 3 batalhas \n" +
				"Ahhh! Caso você morra ou suba de nível.\n" +
				"Você irá retornar à forma base",
		)
		msgResponse.Img = "assets/skins/" + strings.ToLower(player.Gender) + "_ki_awakened.png"
	}

	s.PlayerRepository.Save(player)
	msgResponse.Message = logText
	return msgResponse, nil
}

func (s *PlayerService) Energy(userID int64) (MsgResponse, error) {
	player, _ := s.PlayerRepository.GetOrCreate(userID)

	var msgResponse MsgResponse

	//	player.RegenEnergy()
	msgResponse.Message =
		fmt.Sprintf("📜%d-<b>%s</b> \n"+
			"Você ainda tem HP e energia para mais umas batalhas.\n"+
			"❤️ <b>HP</b>: %d/%d\t\t\t\t\t\t⚡️ <b>Energia</b>: %d/%d\n",
			player.Level,
			player.Name,
			player.HP,
			player.MaxHP,
			player.Energy,
			player.MaxEnergy,
		)
	msgResponse.Img = fmt.Sprintf("assets/skins/%s_%s.png", strings.ToLower(player.Gender), "skin_basic")

	if player.HP < 30 {
		msgResponse.Message = fmt.Sprintf(
			"<b>%s</b> seu HP esta baixo 😫. \n "+
				"❤️ <b>HP</b>: %d/%d \n",
			player.Name,
			player.HP,
			player.MaxHP,
		)
		msgResponse.Img = fmt.Sprintf("assets/skins/%s_%s.png", strings.ToLower(player.Gender), "tired")
	}
	return msgResponse, nil
}
func (s *PlayerService) Recovery(userID int64) (MsgResponse, error) {
	player, _ := s.PlayerRepository.GetOrCreate(userID)
	var msgResponse MsgResponse
	msgResponse = MsgResponse{
		Message: fmt.Sprintf(
			"🙎‍♂️ Ei eu sou Sir Soneca (Dono da pousada)  \n" +
				"Infelizmente você nao podera ficar \n" +
				"Junte mais gold e retorne!!\n",
		),
		Img:        "assets/npcs/sir_soneca.png",
		Conditions: false,
	}

	if player.Gold >= 500 {
		player.Gold = player.Gold - 500
		player.HP = player.MaxHP
		msgResponse = MsgResponse{
			Message: fmt.Sprintf(
				"HP recarregado!!! \n" +
					"Bora voltar pra Luta!! \n",
			),
			Img:        fmt.Sprintf("assets/skins/%s_skin_basic.png", strings.ToLower(player.Gender)),
			Conditions: true,
		}
		err := s.PlayerRepository.Save(player)
		if err != nil {
			log.Fatal(err)
		}
	}

	return msgResponse, nil
}
func (s *PlayerService) Onboarding(userID int64) (*domain.Player, error) {
	return s.PlayerRepository.GetOrCreate(userID)
}
