package application

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
)

func TestNewMenuService(t *testing.T) {
	playerRepo := &fakePlayerRepo{}
	regionRepo := &fakeRegionRepo{}

	svc := NewMenuService(playerRepo, regionRepo)
	if svc == nil {
		t.Fatalf("expected non-nil service")
	}
	if svc.PlayerRepository != playerRepo {
		t.Fatalf("expected player repo set")
	}
	if svc.RegionRepository != regionRepo {
		t.Fatalf("expected region repo set")
	}
}

func TestMenuService_Menu_GetOrCreateError(t *testing.T) {
	svc := &MenuService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) {
				return nil, errors.New("db fail")
			},
		},
		RegionRepository: &fakeRegionRepo{},
	}

	_, err := svc.Menu(1)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMenuService_Menu_RegionError(t *testing.T) {
	player := basePlayer()
	svc := &MenuService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
		RegionRepository: &fakeRegionRepo{
			getByIDFn: func(regionID int64) (*domain.Region, error) {
				return nil, errors.New("region fail")
			},
		},
	}

	_, err := svc.Menu(1)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMenuService_Menu_PassiveUpdateError(t *testing.T) {
	player := basePlayer()
	svc := &MenuService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
			updateFn: func(player *domain.Player) error {
				return errors.New("update fail")
			},
		},
		RegionRepository: &fakeRegionRepo{
			getByIDFn: func(regionID int64) (*domain.Region, error) {
				return &domain.Region{RegionID: 1, Name: "Plains", Tag: "PLAINS", MinLevel: 1}, nil
			},
		},
	}

	_, err := svc.Menu(1)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMenuService_Menu_Success(t *testing.T) {
	player := basePlayer()
	player.Name = "Trunks"
	player.Gold = 333
	player.LastEnergyAt = time.Now().UTC()
	player.Energy = player.MaxEnergy
	player.HP = player.MaxHP

	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	regionRepo := &fakeRegionRepo{
		getByIDFn: func(regionID int64) (*domain.Region, error) {
			return &domain.Region{RegionID: 1, Name: "Planicie", Tag: "PLAINS", MinLevel: 1}, nil
		},
	}
	svc := &MenuService{
		PlayerRepository: playerRepo,
		RegionRepository: regionRepo,
	}

	resp, err := svc.Menu(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "MAIN_MENU" {
		t.Fatalf("expected MAIN_MENU, got %q", resp.MenuKey)
	}
	if !strings.Contains(resp.Text, "Trunks") {
		t.Fatalf("expected player name in text, got %q", resp.Text)
	}
	if !strings.Contains(resp.Text, "333") {
		t.Fatalf("expected player gold in text, got %q", resp.Text)
	}
	if resp.Image == nil || resp.Image.Path != "assets/maps/plains.png" {
		t.Fatalf("unexpected image path: %+v", resp.Image)
	}
	if len(playerRepo.updatedPlayers) != 1 {
		t.Fatalf("expected one passive update call, got %d", len(playerRepo.updatedPlayers))
	}
}
