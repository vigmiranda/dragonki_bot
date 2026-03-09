package application

import (
	"fmt"
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
	player, err := s.PlayerRepository.GetOrCreate(playerID)
	if err != nil {
		return view.GameResponse{}, err
	}

	regions, err := s.RegionRepository.ListByLevel(int64(player.Level))
	if err != nil {
		return view.GameResponse{}, err
	}

	msg := fmt.Sprintf(
		"<b>Word:</b>\n\n" +
			"Plains lvl 1-15 \n" +
			"Forest of Spirits lvl 15-34 \n" +
			"Volcano Lands lvl 35-50 \n" +
			"Frozen Peaks lvl 50-70 \n" +
			"Swamp of Shadows lvl 70+" +
			"\n\n<b>Lugares que você pode explorar!!!</b> \n\n",
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
	player, err := s.PlayerRepository.GetOrCreate(playerID)
	if err != nil {
		return view.GameResponse{}, err
	}

	regionIDStr, ok := params["regionID"]
	if !ok || regionIDStr == "" {
		return view.GameResponse{}, fmt.Errorf("missing param regionID")
	}

	regionId, err := strconv.ParseInt(regionIDStr, 10, 64)
	if err != nil {
		return view.GameResponse{}, err
	}

	player.RegionID = regionId
	if err := s.PlayerRepository.Update(player); err != nil {
		return view.GameResponse{}, err
	}

	region, err := s.RegionRepository.GetByID(regionId)
	if err != nil {
		return view.GameResponse{}, err
	}

	msg, img := getMsgAndImgByRegion(region)
	menu := "MAIN_MENU"

	viewResponse := view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: menu,
	}

	return viewResponse, nil
}

func getMsgAndImgByRegion(region *domain.Region) (string, *view.Image) {
	msg := fmt.Sprintf(
		"\n\n <b>Viagem concluida!!!</b> \n\n",
	)
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/maps/%s.png", domain.ConvertLowCase(region.Tag)),
	}

	switch region.RegionID {
	case 1:
		msg = fmt.Sprintf(
			"\n\n <b>🌿 Bem-vindo às Planícies de Aethoria, guerreiro.</b> \n\n"+
				"Meu nome é Lyra, e sou a guardiã desta rota.\n"+
				"É aqui que muitos iniciam sua jornada… e onde apenas os determinados avançam.\n"+
				"Os inimigos que habitam estas terras variam do nível 1 ao 15.\n"+
				"Cada batalha será um teste da sua força e do seu controle sobre o Ki.\n\n"+
				"Mas lembre-se…\n"+
				"⚔️ A derrota tem preço.\n"+
				"Se cair em combate, perderá <b> %d%% </b> do seu XP e do seu Gold.\n"+
				"Respire fundo.\n"+
				"Mantenha o foco.\n"+
				"E avance com coragem.\n"+
				"As Planícies estão observando.",
			region.Nerf,
		)
		img = &view.Image{
			Type: view.ImageAsset,
			Path: fmt.Sprintf("assets/npcs/lyra.png"),
		}
	case 2:
		msg = fmt.Sprintf(
			"\n\n <b>🌲 Bem-vindo à Forest Spirit, viajante.</b> \n\n"+
				"Eu sou Nyssara, guardiã destas terras ancestrais.\n"+
				"Aqui, a floresta não é apenas cenário… ela observa, sente e reage.\n"+
				"Os espíritos que habitam estas matas são mais astutos e perigosos.\n"+
				"Cada passo descuidado pode atrair algo que você não está preparado para enfrentar.\n"+
				"A floresta recompensa os fortes…"+
				"mas pune os imprudentes.\n\n"+
				"Saiba que aqui…\n"+
				"⚔️ Cada derrota te custará <b> %d%% </b> do seu XP e do seu Gold.\n\n"+
				"Mantenha seu Ki sob controle.\n"+
				"Escute o vento entre as folhas.\n"+
				"E jamais subestime o silêncio.\n\n"+
				"<b>A Forest Spirit já sabe que você chegou. </b>🌿✨",
			region.Nerf,
		)
		img = &view.Image{
			Type: view.ImageAsset,
			Path: fmt.Sprintf("assets/npcs/nyssara_guardia_das_raizes.png"),
		}
	case 3:
		msg = fmt.Sprintf(
			"\n\n <b>🌋 Bem-vindo às Volcano Lands, viajante.</b> \n\n"+
				"Eu sou Pyressa, guardiã das chamas eternas.\n"+
				"Aqui, a terra não dorme… ela ferve.\n"+
				"Rios de lava cortam estas montanhas e criaturas moldadas pelo fogo caminham entre as rochas.\n"+
				"Cada passo em falso pode despertar algo que dorme sob o magma.\n"+
				"As Volcano Lands forjam guerreiros…\n"+
				"mas reduzem os fracos a cinzas.\n\n"+
				"Saiba que aqui…\n"+
				"⚔️ Cada derrota te custará <b> %d%% </b> do seu XP e do seu Gold.\n\n"+
				"Mantenha seu Ki firme.\n"+
				"Sinta o calor das profundezas.\n"+
				"E jamais vire as costas para as chamas.\n\n"+
				"<b>As Volcano Lands já sentiram sua presença. </b>🔥🌋",
			region.Nerf,
		)
		img = &view.Image{
			Type: view.ImageAsset,
			Path: fmt.Sprintf("assets/npcs/pyressa_guardian_das_chamas.png"),
		}
	}
	return msg, img
}
