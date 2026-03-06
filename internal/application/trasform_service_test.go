package application

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
)

func TestNewTransformService(t *testing.T) {
	repo := &fakePlayerRepo{}
	svc := NewTransformService(repo)
	if svc == nil {
		t.Fatalf("expected non-nil service")
	}
	if svc.PlayerRepository != repo {
		t.Fatalf("expected repo set")
	}
}

func TestTransformService_Training(t *testing.T) {
	svc := &TransformService{}
	resp, err := svc.Training()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "MAIN_MENU" {
		t.Fatalf("expected MAIN_MENU, got %q", resp.MenuKey)
	}
	if resp.Image == nil || resp.Image.Path != "assets/maps/ki_ancestral_hall.png" {
		t.Fatalf("unexpected image: %+v", resp.Image)
	}
}

func TestTransformService_ListTransform(t *testing.T) {
	svc := &TransformService{}
	resp, err := svc.ListTransform(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "TRANSFORM_MENU" {
		t.Fatalf("expected TRANSFORM_MENU, got %q", resp.MenuKey)
	}
	if resp.Image == nil || resp.Image.Path != "assets/npcs/mestre_kalunga.png" {
		t.Fatalf("unexpected image: %+v", resp.Image)
	}
	if !strings.Contains(resp.Text, "Ki") {
		t.Fatalf("unexpected text: %q", resp.Text)
	}
}

func TestTransformService_Transform_GetPlayerError(t *testing.T) {
	svc := &TransformService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) {
				return nil, errors.New("db fail")
			},
		},
	}
	_, err := svc.Transform(1)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestTransformService_Transform_FirstUpdateError(t *testing.T) {
	player := basePlayer()
	player.Energy = 1
	player.HP = player.MaxHP

	svc := &TransformService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
			updateFn: func(player *domain.Player) error {
				return errors.New("update fail")
			},
		},
	}
	_, err := svc.Transform(1)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestTransformService_Transform_EnergyInsufficient(t *testing.T) {
	player := basePlayer()
	player.Gender = "MAN"
	player.Class.Skin = "basic"
	player.Level = 10
	player.Energy = 2
	player.MaxEnergy = 2
	player.FormBattlesLeft = 0
	player.HP = player.MaxHP

	repo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	svc := &TransformService{PlayerRepository: repo}

	resp, err := svc.Transform(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Text, "Energia insuficiente") {
		t.Fatalf("expected insufficient energy text, got %q", resp.Text)
	}
	if resp.Image == nil || !strings.Contains(resp.Image.Path, "assets/skins/man_tired.png") {
		t.Fatalf("unexpected image path: %+v", resp.Image)
	}
	if len(repo.updatedPlayers) != 2 {
		t.Fatalf("expected 2 updates (pre and final), got %d", len(repo.updatedPlayers))
	}
}

func TestTransformService_Transform_AlreadyTransformed(t *testing.T) {
	player := basePlayer()
	player.Gender = "WOMAN"
	player.Class.Skin = "basic"
	player.Level = 10
	player.Energy = 10
	player.Form = "Ki Awakened"
	player.FormBattlesLeft = 2
	player.HP = player.MaxHP
	player.LastEnergyAt = time.Now().UTC()

	repo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	svc := &TransformService{PlayerRepository: repo}

	resp, err := svc.Transform(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Text, "j") {
		t.Fatalf("expected already transformed text, got %q", resp.Text)
	}
	if resp.Image == nil || !strings.Contains(resp.Image.Path, "ki_awakened.png") {
		t.Fatalf("unexpected image path: %+v", resp.Image)
	}
}

func TestTransformService_Transform_Success(t *testing.T) {
	player := basePlayer()
	player.Gender = "MAN"
	player.Class.Skin = "basic"
	player.Level = 6
	player.Energy = 6
	player.FormBattlesLeft = 0
	player.HP = player.MaxHP
	player.LastEnergyAt = time.Now().UTC()

	repo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	svc := &TransformService{PlayerRepository: repo}

	resp, err := svc.Transform(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if player.Form != "Ki Awakened" || player.FormBattlesLeft != 3 || player.Energy != 3 {
		t.Fatalf("unexpected transformed state: form=%s battles=%d energy=%d", player.Form, player.FormBattlesLeft, player.Energy)
	}
	if resp.Image == nil || !strings.Contains(resp.Image.Path, "assets/skins/man/basic/ki_awakened.png") {
		t.Fatalf("unexpected image path: %+v", resp.Image)
	}
	if len(repo.updatedPlayers) != 2 {
		t.Fatalf("expected 2 updates (pre and final), got %d", len(repo.updatedPlayers))
	}
}
