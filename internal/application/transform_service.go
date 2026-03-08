package application

import (
	"fmt"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

type TransformService struct {
	PlayerRepository ports.PlayerRepository
	cfg              appConfig
}

func NewTransformService(repo ports.PlayerRepository) *TransformService {
	return &TransformService{PlayerRepository: repo, cfg: getAppConfig()}
}
func (s *TransformService) Training() (view.GameResponse, error) {
	msg := fmt.Sprintf(
		"Area de Treinamento\n\n" +
			"Aqui voce aprimora o controle do Ki e prepara sua transformacao.\n" +
			"Use esta area para revisar as regras e ativar sua forma quando estiver pronto.",
	)
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/maps/ki_ancestral_hall.png"),
	}
	menu := "MAIN_MENU"

	viewResponse := view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: menu,
	}
	return viewResponse, nil
}

func (s *TransformService) ListTransform(playerID int64) (view.GameResponse, error) {
	msg := fmt.Sprintf(
		"“Ouça com atenção, discípulo…” \n\n" +
			"A transformação não é apenas força… é disciplina.\n" +
			"Ao concentrar sua energia vital, você desperta um poder oculto dentro de si.\n" +
			"Seu Ki se intensifica, seus golpes se tornam mais devastadores… mas há um preço.\n" +
			"Cada transformação consome sua energia.\n" +
			"Quanto maior o poder, maior o desgaste.\n" +
			"A partir do nível 5, seu espírito já é forte o bastante para sustentar essa elevação.\n" +
			"Use com sabedoria.\n" +
			"Não é a fúria que vence batalhas…\n" +
			"É o controle do próprio Ki.\n\n" +
			"Ah, e ela dura apenas 3 batalhas!!!\n\n",
	)

	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/npcs/mestre_kalunga.png"),
	}
	menu := "TRANSFORM_MENU"

	viewResponse := view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: menu,
	}

	return viewResponse, nil
}

func (s *TransformService) Transform(userID int64) (view.GameResponse, error) {
	cfg := resolveAppConfig(s.cfg)

	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		return view.GameResponse{}, err
	}

	if err := applyPassivePlayerUpdates(s.PlayerRepository, player, cfg.EnergyRegenDuration); err != nil {
		return view.GameResponse{}, err
	}

	msg := fmt.Sprintf("❌ Transformação desbloqueia no nível 5")
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/skins/%s/%s/tired.png", domain.ConvertLowCase(player.Gender), player.Class.Skin),
	}
	menu := "MAIN_MENU"

	if player.Energy < 3 && player.FormBattlesLeft == 0 {
		msg = fmt.Sprintf("⚡ Energia insuficiente")
		img.Path = fmt.Sprintf("assets/skins/%s_tired.png", domain.ConvertLowCase(player.Gender))
	}
	if player.FormBattlesLeft > 0 {
		msg = fmt.Sprintf(
			"🔥 Você já ascendeu o <b>Ki Awakened</b>!\n"+
				"Restam ainda %d lutas",
			player.FormBattlesLeft,
		)
		img.Path = fmt.Sprintf("assets/skins/%s/%s/ki_awakened.png",
			domain.ConvertLowCase(player.Gender),
			player.Class.Skin)
	}

	if player.Level >= 5 && player.Energy >= 3 && player.FormBattlesLeft == 0 {
		player.Form = "Ki Awakened"
		player.FormBattlesLeft = 3
		player.Energy -= 3
		msg = fmt.Sprintf(
			"🔥 Você ascendeu o <b>Ki Awakened</b>! \n\n" +
				"Essa transformação te dá 20%% em ATK, DEF \n" +
				"Gasta 3 de energia e dura 3 batalhas \n" +
				"Ahhh! Caso você morra ou suba de nível.\n" +
				"Você irá retornar à forma base",
		)
		img.Path = fmt.Sprintf("assets/skins/%s/%s/ki_awakened.png",
			domain.ConvertLowCase(player.Gender),
			player.Class.Skin)
	}

	err = s.PlayerRepository.Update(player)
	if err != nil {
		return view.GameResponse{}, err
	}

	viewResponse := view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: menu,
	}

	return viewResponse, nil
}
