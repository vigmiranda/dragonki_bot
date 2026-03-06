package application

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
)

func TestNewEnergyService(t *testing.T) {
	repo := &fakePlayerRepo{}
	svc := NewEnergyService(repo)
	if svc == nil {
		t.Fatalf("expected non-nil service")
	}
	if svc.PlayerRepository != repo {
		t.Fatalf("expected repo to be set")
	}
}

func TestEnergyService_Energy_GetOrCreateError(t *testing.T) {
	svc := &EnergyService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) {
				return nil, errors.New("db fail")
			},
		},
	}

	_, err := svc.Energy(1)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestEnergyService_Energy_PassiveUpdateError(t *testing.T) {
	player := basePlayer()
	svc := &EnergyService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
			updateFn: func(player *domain.Player) error {
				return errors.New("update fail")
			},
		},
	}

	_, err := svc.Energy(1)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestEnergyService_Energy_Success(t *testing.T) {
	player := basePlayer()
	player.Name = "Vegeta"
	player.LastEnergyAt = time.Now().UTC()
	player.Energy = player.MaxEnergy
	player.HP = player.MaxHP

	repo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	svc := &EnergyService{PlayerRepository: repo}

	resp, err := svc.Energy(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "ENERGY_MENU" {
		t.Fatalf("expected ENERGY_MENU, got %q", resp.MenuKey)
	}
	if !strings.Contains(resp.Text, "Vegeta") {
		t.Fatalf("expected player name in message, got %q", resp.Text)
	}
	if resp.Image == nil || resp.Image.Path != "assets/npcs/sir_soneca.png" {
		t.Fatalf("unexpected image: %+v", resp.Image)
	}
	if len(repo.updatedPlayers) != 1 {
		t.Fatalf("expected one passive update call, got %d", len(repo.updatedPlayers))
	}
}

func TestEnergyService_Rest_GetOrCreateError(t *testing.T) {
	svc := &EnergyService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) {
				return nil, errors.New("db fail")
			},
		},
		cfg: appConfig{EnergyRegenTime: "10m"},
	}
	_, err := svc.Rest(1)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestEnergyService_Rest_PassiveUpdateError(t *testing.T) {
	player := basePlayer()
	svc := &EnergyService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
			updateFn: func(player *domain.Player) error {
				return errors.New("update fail")
			},
		},
		cfg: appConfig{EnergyRegenTime: "10m"},
	}
	_, err := svc.Rest(1)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestEnergyService_Rest_NoEnergy(t *testing.T) {
	player := basePlayer()
	player.Energy = 0
	player.LastEnergyAt = time.Now().UTC()
	player.HP = player.MaxHP

	repo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	svc := &EnergyService{
		PlayerRepository: repo,
		cfg:              appConfig{EnergyRegenTime: "12m"},
	}

	resp, err := svc.Rest(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Text, "12m") {
		t.Fatalf("expected regen time in response, got %q", resp.Text)
	}
	if resp.Image == nil || resp.Image.Path != "assets/npcs/sir_soneca.png" {
		t.Fatalf("unexpected image: %+v", resp.Image)
	}
	if len(repo.updatedPlayers) != 1 {
		t.Fatalf("expected only passive update call, got %d", len(repo.updatedPlayers))
	}
}

func TestEnergyService_Rest_WithEnergyUpdatesHPAndEnergy(t *testing.T) {
	player := basePlayer()
	player.Energy = 2
	player.HP = 20
	player.MaxHP = 100
	player.Gender = "MAN"
	player.Class.Skin = "basic"
	player.LastEnergyAt = time.Now().UTC()

	repo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	svc := &EnergyService{
		PlayerRepository: repo,
		cfg:              appConfig{EnergyRegenTime: "10m"},
	}

	resp, err := svc.Rest(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if player.Energy != 1 {
		t.Fatalf("expected energy 1, got %d", player.Energy)
	}
	if player.HP != 50 {
		t.Fatalf("expected hp +30%% => 50, got %d", player.HP)
	}
	if resp.Image == nil || !strings.Contains(resp.Image.Path, "assets/skins/man/basicplayer.png") {
		t.Fatalf("unexpected image: %+v", resp.Image)
	}
	if len(repo.updatedPlayers) != 2 {
		t.Fatalf("expected passive+rest update calls, got %d", len(repo.updatedPlayers))
	}
}

func TestEnergyService_Rest_UpdateErrorOnSecondUpdate(t *testing.T) {
	player := basePlayer()
	player.Energy = 1
	player.HP = 10
	player.MaxHP = 100
	player.LastEnergyAt = time.Now().UTC()

	call := 0
	repo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		updateFn: func(player *domain.Player) error {
			call++
			if call == 2 {
				return errors.New("save fail")
			}
			return nil
		},
	}
	svc := &EnergyService{PlayerRepository: repo, cfg: appConfig{EnergyRegenTime: "10m"}}

	_, err := svc.Rest(1)
	if err == nil {
		t.Fatalf("expected error on second update")
	}
}

func TestEnergyService_Sleep_GetOrCreateError(t *testing.T) {
	svc := &EnergyService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) {
				return nil, errors.New("db fail")
			},
		},
	}
	_, err := svc.Sleep(1)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestEnergyService_Sleep_PassiveUpdateError(t *testing.T) {
	player := basePlayer()
	svc := &EnergyService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
			updateFn: func(player *domain.Player) error {
				return errors.New("update fail")
			},
		},
	}
	_, err := svc.Sleep(1)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestEnergyService_Sleep_NotEnoughGold(t *testing.T) {
	player := basePlayer()
	player.Gold = 499
	player.HP = 10
	player.LastEnergyAt = time.Now().UTC()

	repo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	svc := &EnergyService{PlayerRepository: repo}

	resp, err := svc.Sleep(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Text, "Junte mais gold") {
		t.Fatalf("expected not-enough-gold message, got %q", resp.Text)
	}
	if player.Gold != 499 || player.HP != 10 {
		t.Fatalf("expected unchanged player stats, got gold=%d hp=%d", player.Gold, player.HP)
	}
	if len(repo.updatedPlayers) != 1 {
		t.Fatalf("expected only passive update call, got %d", len(repo.updatedPlayers))
	}
}

func TestEnergyService_Sleep_WithGold(t *testing.T) {
	player := basePlayer()
	player.Gold = 800
	player.HP = 10
	player.MaxHP = 120
	player.LastEnergyAt = time.Now().UTC()

	repo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	svc := &EnergyService{PlayerRepository: repo}

	resp, err := svc.Sleep(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if player.Gold != 300 {
		t.Fatalf("expected gold 300, got %d", player.Gold)
	}
	if player.HP != 120 {
		t.Fatalf("expected hp full 120, got %d", player.HP)
	}
	if !strings.Contains(resp.Text, "HP recarregado") {
		t.Fatalf("unexpected response text: %q", resp.Text)
	}
	if len(repo.updatedPlayers) != 2 {
		t.Fatalf("expected passive+sleep update calls, got %d", len(repo.updatedPlayers))
	}
}

func TestEnergyService_Sleep_UpdateErrorOnSecondUpdate(t *testing.T) {
	player := basePlayer()
	player.Gold = 500
	player.HP = 5
	player.MaxHP = 50
	player.LastEnergyAt = time.Now().UTC()

	call := 0
	repo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		updateFn: func(player *domain.Player) error {
			call++
			if call == 2 {
				return errors.New("save fail")
			}
			return nil
		},
	}
	svc := &EnergyService{PlayerRepository: repo}

	_, err := svc.Sleep(1)
	if err == nil {
		t.Fatalf("expected error on second update")
	}
}
