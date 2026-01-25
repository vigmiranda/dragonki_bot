package application

import (
	"fmt"
	"log"
	"os"
	"strconv"

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
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}

	var msgResponse MsgResponse
	player.Gender = gender
	player.OnboardingStep = "CHOOSING_NAME"

	err = s.PlayerRepository.Save(player)
	if err != nil {
		log.Fatal(err)
	}
	msgResponse.Message = fmt.Sprintf(
		"🧙 Status\n❤️ HP: %d\n⚔️ ATK: %d\n✨ XP: %d\n🏆 Nível: %d\n⚡ Energia: %d",
		player.HP, player.ATK, player.XP, player.Level, player.Energy,
	)
	msgResponse.Img = "assets/skins/" + domain.ConvertLowCase(gender) + "_skin_basic.png"

	return msgResponse, nil
}
func (s *PlayerService) SetName(userID int64, name string) error {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}
	player.Name = name
	player.OnboardingStep = "DONE"

	return s.PlayerRepository.Save(player)
}

func (s *PlayerService) Menu(userID int64) (MsgResponse, error) {
	var player, err = s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}
	var playersOnline, _ = s.PlayerRepository.CountOnlineLast15Min()
	region, err := s.RegionRepository.GetByID(player.RegionID)
	if err != nil {
		log.Fatal(err)
	}
	if player.Energy < player.MaxEnergy || player.HP < player.MaxHP {
		player.RegenEnergy()
	}
	player.SetTransformBase()
	err = s.PlayerRepository.Save(player)
	if err != nil {
		log.Fatal(err)
	}

	var msgResponse MsgResponse
	multiPlayersOnlineDB := os.Getenv("MULT_PLAYERS")
	multiPlayersOnline, _ := strconv.Atoi(multiPlayersOnlineDB)
	msgResponse.Message = fmt.Sprintf(
		"<b>%s</b> \n\n"+
			"👥 <b>Online (últimos 15 min)</b>: %d\n"+
			"🏰 <b>Local</b>:  %s (L%d)\n"+
			"❤️ <b>HP</b>: %d/%d\n"+
			"⚡️ <b>Energia</b>: %d/%d\n"+
			"💪 <b>Level</b>: %d\n"+
			"💰 <b>Gold</b>: %d\n",
		player.Name,
		playersOnline*multiPlayersOnline,
		region.Name,
		region.MinLevel,
		player.HP,
		player.MaxHP,
		player.Energy,
		player.MaxEnergy,
		player.Level,
		player.Gold,
	)
	msgResponse.Img = fmt.Sprintf("assets/maps/%s.png", domain.ConvertLowCase(region.Tag))
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
	err = s.PlayerRepository.Save(player)
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
	msgResponse.Img = "assets/skins/" + domain.ConvertLowCase(player.Gender) + "_skin_basic.png"
	if player.FormBattlesLeft > 0 {
		msgResponse.Img = "assets/skins/" + domain.ConvertLowCase(player.Gender) + "_ki_awakened.png"
	}

	return msgResponse, nil
}

func (s *PlayerService) Transform(userID int64) (MsgResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}
	if player.Energy < player.MaxEnergy || player.HP < player.MaxHP {
		player.RegenEnergy()
	}
	err = s.PlayerRepository.Save(player)
	if err != nil {
		log.Fatal(err)
	}

	var msgResponse MsgResponse

	logText := ""

	if player.Level < 3 {
		logText += fmt.Sprintf("❌ Transformação desbloqueia no nível 3")
		msgResponse.Img = "assets/skins/" + domain.ConvertLowCase(player.Gender) + "_skin_basic.png"
	}
	if player.Energy < 3 && player.FormBattlesLeft == 0 {
		logText += fmt.Sprintf("⚡ Energia insuficiente")
		msgResponse.Img = "assets/skins/" + domain.ConvertLowCase(player.Gender) + "_tired.png"
	}
	if player.FormBattlesLeft > 0 {
		logText += fmt.Sprintf(
			"🔥 Você já ascendeu o <b>Ki Awakened</b>!\n"+
				"Restam ainda %d lutas",
			player.FormBattlesLeft,
		)
		msgResponse.Img = "assets/skins/" + domain.ConvertLowCase(player.Gender) + "_ki_awakened.png"
	}

	if player.Level >= 3 && player.Level <= 15 && player.Energy >= 3 && player.FormBattlesLeft == 0 {
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
		msgResponse.Img = "assets/skins/" + domain.ConvertLowCase(player.Gender) + "_ki_awakened.png"
	}

	err = s.PlayerRepository.Save(player)
	if err != nil {
		log.Fatal(err)
	}
	msgResponse.Message = logText
	return msgResponse, nil
}

func (s *PlayerService) Energy(userID int64) (MsgResponse, error) {
	var player, err = s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}
	if player.Energy < player.MaxEnergy || player.HP < player.MaxHP {
		player.RegenEnergy()
	}
	player.SetTransformBase()
	err = s.PlayerRepository.Save(player)
	if err != nil {
		log.Fatal(err)
	}

	var msgResponse MsgResponse
	msgResponse.Message =
		fmt.Sprintf("📜%d-<b>%s</b> \n"+
			"Você ainda tem HP e energia para mais umas batalhas.\n"+
			"❤️ <b>HP</b>: %d/%d\t\t\t\t\t\t⚡️ <b>Energia</b>: %d/%d\n"+
			"💰 <b>Gold</b>: %d\n",
			player.Level,
			player.Name,
			player.HP,
			player.MaxHP,
			player.Energy,
			player.MaxEnergy,
			player.Gold,
		)
	msgResponse.Img = fmt.Sprintf("assets/skins/%s_%s.png", domain.ConvertLowCase(player.Gender), "skin_basic")

	if player.HP < 30 {
		msgResponse.Message = fmt.Sprintf(
			"<b>%s</b> seu HP esta baixo 😫. \n "+
				"❤️ <b>HP</b>: %d/%d \n",
			player.Name,
			player.HP,
			player.MaxHP,
		)
		msgResponse.Img = fmt.Sprintf("assets/skins/%s_%s.png", domain.ConvertLowCase(player.Gender), "tired")
	}
	return msgResponse, nil
}
func (s *PlayerService) Recovery(userID int64) (MsgResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}
	if player.Energy < player.MaxEnergy || player.HP < player.MaxHP {
		player.RegenEnergy()
	}
	player.SetTransformBase()
	err = s.PlayerRepository.Save(player)
	if err != nil {
		log.Fatal(err)
	}

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
				"❤️ HP recarregado!!! \n" +
					"Bora voltar pra Luta ⚔️!! \n",
			),
			Img:        fmt.Sprintf("assets/skins/%s_skin_basic.png", domain.ConvertLowCase(player.Gender)),
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
