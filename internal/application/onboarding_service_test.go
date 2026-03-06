package application

import (
	"errors"
	"strings"
	"testing"

	"github.com/vigmiranda/dragonki/internal/domain"
)

func TestNewOnboardingService(t *testing.T) {
	repo := &fakePlayerRepo{}
	svc := NewOnboardingService(repo)
	if svc == nil {
		t.Fatalf("expected non-nil service")
	}
	if svc.PlayerRepository != repo {
		t.Fatalf("expected repo set")
	}
}

func TestOnboardingService_Onboarding(t *testing.T) {
	player := basePlayer()
	repo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	svc := &OnboardingService{PlayerRepository: repo}

	got, err := svc.Onboarding(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != player {
		t.Fatalf("expected same player pointer")
	}
}

func TestOnboardingService_HasCompletedOnboarding_GetPlayerError(t *testing.T) {
	svc := &OnboardingService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) {
				return nil, errors.New("db fail")
			},
		},
	}
	_, err := svc.HasCompletedOnboarding(1, 10, "name", "username")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestOnboardingService_HasCompletedOnboarding_UpdateError(t *testing.T) {
	player := basePlayer()
	svc := &OnboardingService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
			updateFn: func(player *domain.Player) error {
				return errors.New("update fail")
			},
		},
	}
	_, err := svc.HasCompletedOnboarding(1, 10, "name", "username")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestOnboardingService_HasCompletedOnboarding_Success(t *testing.T) {
	player := basePlayer()
	player.OnboardingStep = "CHOOSING_NAME"
	repo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	svc := &OnboardingService{PlayerRepository: repo}

	step, err := svc.HasCompletedOnboarding(1, 123, "Goku", "kakarot")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if step != "CHOOSING_NAME" {
		t.Fatalf("expected onboarding step CHOOSING_NAME, got %q", step)
	}
	if player.UserIDTelegram != 123 || player.NameTelegram != "Goku" || player.UserName != "kakarot" {
		t.Fatalf("expected telegram fields updated, got %+v", player)
	}
	if len(repo.updatedPlayers) != 1 {
		t.Fatalf("expected one update call")
	}
}

func TestOnboardingService_HandleSelectGender(t *testing.T) {
	svc := &OnboardingService{}
	resp, err := svc.HandleSelectGender()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "GENDER_MENU" {
		t.Fatalf("expected GENDER_MENU, got %q", resp.MenuKey)
	}
	if resp.Image == nil || resp.Image.Path != "assets/skins/man_woman_select.png" {
		t.Fatalf("unexpected image: %+v", resp.Image)
	}
}

func TestOnboardingService_SetGender_GetPlayerError(t *testing.T) {
	svc := &OnboardingService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) {
				return nil, errors.New("db fail")
			},
		},
	}
	_, err := svc.SetGender(1, "MAN")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestOnboardingService_SetGender_UpdateError(t *testing.T) {
	player := basePlayer()
	svc := &OnboardingService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
			updateFn: func(player *domain.Player) error {
				return errors.New("update fail")
			},
		},
	}
	_, err := svc.SetGender(1, "WOMAN")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestOnboardingService_SetGender_Success(t *testing.T) {
	player := basePlayer()
	player.Class.Skin = "basic"
	repo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	svc := &OnboardingService{PlayerRepository: repo}

	resp, err := svc.SetGender(1, "WOMAN")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if player.Gender != "WOMAN" || player.OnboardingStep != "CHOOSING_NAME" {
		t.Fatalf("expected gender and step updated, got gender=%s step=%s", player.Gender, player.OnboardingStep)
	}
	if resp.Image == nil || !strings.Contains(resp.Image.Path, "assets/skins/woman/basicplayer.png") {
		t.Fatalf("unexpected image path: %+v", resp.Image)
	}
	if len(repo.updatedPlayers) != 1 {
		t.Fatalf("expected one update call")
	}
}

func TestOnboardingService_SetName_GetPlayerError(t *testing.T) {
	svc := &OnboardingService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) {
				return nil, errors.New("db fail")
			},
		},
	}
	err := svc.SetName(1, "Vegeta")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestOnboardingService_SetName_UpdateError(t *testing.T) {
	player := basePlayer()
	svc := &OnboardingService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
			updateFn: func(player *domain.Player) error {
				return errors.New("update fail")
			},
		},
	}
	err := svc.SetName(1, "Vegeta")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestOnboardingService_SetName_Success(t *testing.T) {
	player := basePlayer()
	repo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	svc := &OnboardingService{PlayerRepository: repo}

	err := svc.SetName(1, "Vegeta")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if player.Name != "Vegeta" || player.OnboardingStep != "DONE" {
		t.Fatalf("expected name and step updated, got name=%s step=%s", player.Name, player.OnboardingStep)
	}
	if len(repo.updatedPlayers) != 1 {
		t.Fatalf("expected one update call")
	}
}
