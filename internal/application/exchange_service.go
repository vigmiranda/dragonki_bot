package application

import (
	"fmt"
	"os"
	"strconv"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/ports"
)

type ExchangeService struct {
	PlayerRepository ports.PlayerRepository
}

func NewExchangeService(repo ports.PlayerRepository) *ExchangeService {
	return &ExchangeService{PlayerRepository: repo}
}

func (s ExchangeService) Open() (view.GameResponse, error) {
	msg := fmt.Sprintf(
		"💠 Casa de Câmbio de DragonKi 💠\n\n" +
			"“Saudações, viajante.”\n\n" +
			"Eu sou Mestre Tríkan, guardião do fluxo do Ki e das riquezas deste mundo.\n" +
			"Aqui, seus Ki Crystals 💎 podem ser convertidos em ouro… " +
			"e seu progresso, acelerado.\n\n" +
			"Use com sabedoria.\n" +
			"O equilíbrio entre poder e ganância define o verdadeiro guerreiro. ⚖️🔥",
	)
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/npcs/mestre_trikan.png"),
	}
	menu := "EXCHANGE_MENU"

	viewResponse := view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: menu,
	}
	return viewResponse, nil
}

func (s ExchangeService) Convert(userID int64, params map[string]string) (view.GameResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		return view.GameResponse{}, err
	}

	amountConvertStr, ok := params["amountConvert"]
	if !ok || amountConvertStr == "" {
		return view.GameResponse{}, fmt.Errorf("missing param amountConvert")
	}
	amountConvert, err := strconv.Atoi(amountConvertStr)
	if err != nil {
		return view.GameResponse{}, err
	}

	kiCrystalsConvert, err := strconv.Atoi(os.Getenv("KI_CRYSTALS_CONVERT"))
	if err != nil {
		return view.GameResponse{}, err
	}

	msg := fmt.Sprintf(
		"💠 Casa de Câmbio de DragonKi 💠\n\n" +
			"Infelizmente, você não tem Ki Crystals 💎 suficiente.\n" +
			"Volte quando conseguir alguns!!!”\n\n")

	if player.KiCrystals >= amountConvert {
		player.KiCrystals -= amountConvert
		player.Gold += amountConvert * kiCrystalsConvert

		if err := s.PlayerRepository.Update(player); err != nil {
			return view.GameResponse{}, err
		}

		msg = fmt.Sprintf(
			"💠 Casa de Câmbio de DragonKi 💠\n\n"+
				"Feito %s, agora você tem %d gold 💰 e %d de Ki Crystals 💎\n"+
				"Volte sempre que precisar!!!”\n\n",
			player.Name,
			player.Gold,
			player.KiCrystals,
		)
	}

	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/npcs/mestre_trikan.png"),
	}
	menu := "EXCHANGE_MENU"

	viewResponse := view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: menu,
	}
	return viewResponse, nil
}
