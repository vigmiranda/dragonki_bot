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
	Img        tgbotapi.RequestFileData
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
func (s *PlayerService) Onboarding(userID int64) (*domain.Player, error) {
	return s.PlayerRepository.GetOrCreate(userID)
}
func (s *PlayerService) HasCompletedOnboarding(userID int64) (string, error) {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		return "", err
	}
	return player.OnboardingStep, nil
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
	msgResponse.Img = tgbotapi.FilePath(fmt.Sprintf("assets/skins/%s_skin_basic.png", domain.ConvertLowCase(gender)))
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
		"<b> %s </b> \n\n"+
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
	msgResponse.Img = tgbotapi.FilePath(fmt.Sprintf("assets/maps/%s.png", domain.ConvertLowCase(region.Tag)))
	return msgResponse, nil
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
			"🗺️ <b>Local</b>: %s(Lv%d)\n"+
			"💪 <b>Level</b>: %d\t\t\t\t\t\t\t\t\t\t\t🧬 <b>XP</b>: %d\n"+
			"❤️ <b>HP</b>: %d/%d\t\t\t\t\t\t⚡️ <b>Energia</b>: %d/%d\n"+
			"⚔️ <b>ATK</b>: %d \t\t\t\t\t\t\t\t\t\t🛡️ <b>DEF</b>: %d \n"+
			"💰 <b>Gold</b>: %d\n"+
			"💎 <b>KiCrystals</b>: %d\n",
		player.Name,
		player.Class,
		region.Name,
		player.Level,
		player.XP,
		player.HP,
		player.MaxHP,
		player.Energy,
		player.MaxEnergy,
		player.ATK,
		player.DEF,
		region.MinLevel,
		player.Gold,
		player.KiCrystals,
	)
	msgResponse.Img = tgbotapi.FilePath(fmt.Sprintf("assets/skins/%s_skin_basic.png", domain.ConvertLowCase(player.Gender)))
	if player.FormBattlesLeft > 0 {
		msgResponse.Img = tgbotapi.FilePath(fmt.Sprintf("assets/skins/%s_ki_awakened.png", domain.ConvertLowCase(player.Gender)))
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
		msgResponse.Img = tgbotapi.FilePath(fmt.Sprintf("assets/skins/%s_skin_basic.png", domain.ConvertLowCase(player.Gender)))
	}
	if player.Energy < 3 && player.FormBattlesLeft == 0 {
		logText += fmt.Sprintf("⚡ Energia insuficiente")
		msgResponse.Img = tgbotapi.FilePath(fmt.Sprintf("assets/skins/%s_tired.png", domain.ConvertLowCase(player.Gender)))
	}
	if player.FormBattlesLeft > 0 {
		logText += fmt.Sprintf(
			"🔥 Você já ascendeu o <b>Ki Awakened</b>!\n"+
				"Restam ainda %d lutas",
			player.FormBattlesLeft,
		)
		msgResponse.Img = tgbotapi.FilePath(fmt.Sprintf("assets/skins/%s_ki_awakened.png", domain.ConvertLowCase(player.Gender)))
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
		msgResponse.Img = tgbotapi.FilePath(fmt.Sprintf("assets/skins/%s_ki_awakened.png", domain.ConvertLowCase(player.Gender)))
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
		fmt.Sprintf(
			"🛏️ Sir Soneca, Guardião do Descanso 🛏️ \n\n"+
				"“Ooooh… hã? Um aventureiro?” 😴 \n"+
				"Bem-vindo %s à minha humilde pousada… \n\n"+
				"Sou Sir Soneca, cavaleiro aposentado e mestre supremo da arte de dormir bem. \n\n"+
				"Aqui você pode descansar, recuperar suas forças e deixar que o cansaço da batalha desapareça. \n"+
				"Uma boa noite de sono cura feridas, restaura o HP e prepara seu espírito para a próxima luta. \n\n"+
				"🕯️ Camas macias, silêncio garantido… e roncos inclusos no pacote. \n\n"+
				"Então, jovem guerreiro(a) … vai apenas descansar um pouco… \n"+
				"ou dormir como um verdadeiro herói? 😴💤 \n\n",
			player.Name)
	msgResponse.Img =
		tgbotapi.FilePath(fmt.Sprintf("assets/npcs/sir_soneca.png"))
	return msgResponse, nil
}
func (s *PlayerService) Rest(userID int64) (MsgResponse, error) {
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
			"Você nao tem energia\n"+
				"Aguarde %s!! e volter aqui para descansar!!\n",
			os.Getenv("ENERGY_REGEN_TIME"),
		),
		Img:        tgbotapi.FilePath(fmt.Sprintf("assets/npcs/sir_soneca.png")),
		Conditions: false,
	}

	if player.Energy >= 1 {
		player.Energy -= 1

		player.HP += int(float64(player.MaxHP) * 0.3)
		if player.HP > player.MaxHP {
			player.HP = player.MaxHP
		}

		msgResponse = MsgResponse{
			Message: fmt.Sprintf(
				"❤️ HP recarregado!!! %d/%d\n"+
					"Bora voltar pra Luta ⚔️!! \n",
				player.HP,
				player.MaxHP,
			),
			Img:        tgbotapi.FilePath(fmt.Sprintf("assets/skins/%s_skin_basic.png", domain.ConvertLowCase(player.Gender))),
			Conditions: true,
		}
		err := s.PlayerRepository.Save(player)
		if err != nil {
			log.Fatal(err)
		}
	}
	return msgResponse, nil

}
func (s *PlayerService) Sleep(userID int64) (MsgResponse, error) {
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
			"Infelizmente você não podera ficar \n" +
				"Junte mais gold e retorne!!\n",
		),
		Img:        tgbotapi.FilePath(fmt.Sprintf("assets/npcs/sir_soneca.png")),
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
			Img:        tgbotapi.FilePath(fmt.Sprintf("assets/skins/%s_skin_basic.png", domain.ConvertLowCase(player.Gender))),
			Conditions: true,
		}
		err := s.PlayerRepository.Save(player)
		if err != nil {
			log.Fatal(err)
		}
	}

	return msgResponse, nil
}

func (s *PlayerService) OpenShop(userID int64) (MsgResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}

	var msgResponse MsgResponse
	msgResponse = MsgResponse{
		Message: fmt.Sprintf(
			"🪙 Kiara, a Mercadora 🪙 \n\n"+
				"Seja bem-vindo %s à minha lojinha! \n\n"+
				"Eu sou Kiara, mercadora viajante e guardiã de artefatos raros. \n\n"+
				"Aqui você encontrará poções, equipamentos e itens especiais que vão fortalecer seu corpo, "+
				"renovar seu Ki e te ajudar a sobreviver aos desafios que te esperam lá fora. \n\n"+
				"💎 Cada item aqui foi escolhido para quem busca evoluir, não apenas lutar. \n"+
				"Prepare-se bem… a jornada de um verdadeiro guerreiro (a) começa com as escolhas certas. \n\n"+
				"O que posso fazer por você hoje? 😄 \n\n",
			player.Name,
		),
		Img:        tgbotapi.FilePath(fmt.Sprintf("assets/npcs/kiara_mercadora.png")),
		Conditions: false,
	}

	return msgResponse, nil

}
