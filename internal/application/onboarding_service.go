package application

import (
	"fmt"
	"log"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

type OnboardingService struct {
	PlayerRepository ports.PlayerRepository
}

func NewOnboardingService(repo ports.PlayerRepository) *OnboardingService {
	return &OnboardingService{PlayerRepository: repo}
}

func (s *OnboardingService) Onboarding(userID int64) (*domain.Player, error) {
	return s.PlayerRepository.GetOrCreate(userID)
}
func (s *OnboardingService) HasCompletedOnboarding(
	userID int64,
	userIDTelegram int64,
	nameTelegram string,
	userName string,
) (string, error) {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		return "", err
	}
	player.UserIDTelegram = userIDTelegram
	player.NameTelegram = nameTelegram
	player.UserName = userName
	s.PlayerRepository.Update(player)
	return player.OnboardingStep, nil
}
func (s *OnboardingService) HandleSelectGender() (view.GameResponse, error) {
	msg := fmt.Sprintf("Selecione o genero do seu personagem")
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/skins/man_woman_select.png"),
	}
	menu := "GENDER_MENU"
	viewResponse := view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: menu,
	}
	return viewResponse, nil
}

func (s *OnboardingService) SetGender(userID int64, gender string) (view.GameResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}

	player.Gender = gender
	player.OnboardingStep = "CHOOSING_NAME"
	err = s.PlayerRepository.Update(player)
	if err != nil {
		log.Fatal(err)
	}

	msg := fmt.Sprintf("✍️ Ótima escolha!\n\nAgora digite o nome do seu personagem:")
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/skins/%s_skin_basic.png", domain.ConvertLowCase(player.Gender)),
	}

	viewResponse := view.GameResponse{
		Text:  msg,
		Image: img,
	}

	return viewResponse, nil
}

func (s *OnboardingService) SetName(userID int64, name string) error {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}
	player.Name = name
	player.OnboardingStep = "DONE"

	return s.PlayerRepository.Update(player)
}
