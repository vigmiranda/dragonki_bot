package application

import (
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
func (s *OnboardingService) HasCompletedOnboarding(userID int64) (string, error) {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		return "", err
	}
	return player.OnboardingStep, nil
}

func (s *OnboardingService) SetGender(userID int64, gender string) (view.MsgResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}

	var msgResponse view.MsgResponse
	player.Gender = gender
	player.OnboardingStep = "CHOOSING_NAME"

	err = s.PlayerRepository.Save(player)
	if err != nil {
		log.Fatal(err)
	}
	return msgResponse, nil
}

func (s *OnboardingService) SetName(userID int64, name string) error {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}
	player.Name = name
	player.OnboardingStep = "DONE"

	return s.PlayerRepository.Save(player)
}
