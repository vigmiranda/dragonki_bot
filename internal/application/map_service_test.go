package application

import (
	"errors"
	"strings"
	"testing"

	"github.com/vigmiranda/dragonki/internal/domain"
)

func TestNewMapService(t *testing.T) {
	playerRepo := &fakePlayerRepo{}
	regionRepo := &fakeRegionRepo{}

	svc := NewMapService(playerRepo, regionRepo)
	if svc == nil {
		t.Fatalf("expected non-nil service")
	}
	if svc.PlayerRepository != playerRepo {
		t.Fatalf("expected player repo to be set")
	}
	if svc.RegionRepository != regionRepo {
		t.Fatalf("expected region repo to be set")
	}
}

func TestMapService_ListMap_PlayerError(t *testing.T) {
	svc := MapService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) {
				return nil, errors.New("db fail")
			},
		},
		RegionRepository: &fakeRegionRepo{},
	}

	_, err := svc.ListMap(1)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMapService_ListMap_ListByLevelError(t *testing.T) {
	player := basePlayer()
	svc := MapService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
		RegionRepository: &fakeRegionRepo{
			listByLevelFn: func(level int64) ([]domain.Region, error) {
				return nil, errors.New("region fail")
			},
		},
	}

	_, err := svc.ListMap(1)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMapService_ListMap_Success(t *testing.T) {
	player := basePlayer()
	player.Level = 20
	regions := []domain.Region{
		{RegionID: 1, Name: "Plains", MinLevel: 1},
		{RegionID: 2, Name: "Forest", MinLevel: 15},
	}

	svc := MapService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
		RegionRepository: &fakeRegionRepo{
			listByLevelFn: func(level int64) ([]domain.Region, error) {
				if level != int64(player.Level) {
					t.Fatalf("expected level %d, got %d", player.Level, level)
				}
				return regions, nil
			},
		},
	}

	resp, err := svc.ListMap(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "REGIONS_MENU" {
		t.Fatalf("expected REGIONS_MENU, got %q", resp.MenuKey)
	}
	if resp.Image == nil || resp.Image.Path != "assets/maps/word.png" {
		t.Fatalf("unexpected image: %+v", resp.Image)
	}
	gotRegions, ok := resp.MenuParams["Regions"].([]domain.Region)
	if !ok {
		t.Fatalf("expected Regions in menu params")
	}
	if len(gotRegions) != 2 {
		t.Fatalf("expected 2 regions, got %d", len(gotRegions))
	}
}

func TestMapService_ToRegion_PlayerError(t *testing.T) {
	svc := MapService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) {
				return nil, errors.New("db fail")
			},
		},
		RegionRepository: &fakeRegionRepo{},
	}
	_, err := svc.ToRegion(1, map[string]string{"regionID": "1"})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMapService_ToRegion_MissingRegionID(t *testing.T) {
	svc := MapService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return basePlayer(), nil },
		},
		RegionRepository: &fakeRegionRepo{},
	}
	_, err := svc.ToRegion(1, map[string]string{})
	if err == nil || !strings.Contains(err.Error(), "missing param regionID") {
		t.Fatalf("expected missing param error, got %v", err)
	}
}

func TestMapService_ToRegion_InvalidRegionID(t *testing.T) {
	svc := MapService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return basePlayer(), nil },
		},
		RegionRepository: &fakeRegionRepo{},
	}
	_, err := svc.ToRegion(1, map[string]string{"regionID": "abc"})
	if err == nil {
		t.Fatalf("expected parse error")
	}
}

func TestMapService_ToRegion_UpdateError(t *testing.T) {
	svc := MapService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return basePlayer(), nil },
			updateFn: func(player *domain.Player) error {
				return errors.New("update fail")
			},
		},
		RegionRepository: &fakeRegionRepo{},
	}
	_, err := svc.ToRegion(1, map[string]string{"regionID": "2"})
	if err == nil {
		t.Fatalf("expected update error")
	}
}

func TestMapService_ToRegion_GetByIDError(t *testing.T) {
	player := basePlayer()
	svc := MapService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
		RegionRepository: &fakeRegionRepo{
			getByIDFn: func(regionID int64) (*domain.Region, error) {
				return nil, errors.New("region fail")
			},
		},
	}
	_, err := svc.ToRegion(1, map[string]string{"regionID": "2"})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMapService_ToRegion_SuccessDefaultRegionMsg(t *testing.T) {
	player := basePlayer()
	svc := MapService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
		RegionRepository: &fakeRegionRepo{
			getByIDFn: func(regionID int64) (*domain.Region, error) {
				return &domain.Region{RegionID: 99, Tag: "MYSTERY", Nerf: 10}, nil
			},
		},
	}

	resp, err := svc.ToRegion(1, map[string]string{"regionID": "99"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "MAIN_MENU" {
		t.Fatalf("expected MAIN_MENU, got %q", resp.MenuKey)
	}
	if player.RegionID != 99 {
		t.Fatalf("expected player region updated to 99, got %d", player.RegionID)
	}
	if resp.Image == nil || !strings.Contains(resp.Image.Path, "assets/maps/mystery.png") {
		t.Fatalf("unexpected image path: %+v", resp.Image)
	}
	if !strings.Contains(resp.Text, "Viagem concluida") {
		t.Fatalf("unexpected text: %q", resp.Text)
	}
}

func TestGetMsgAndImgByRegion_Region1And2(t *testing.T) {
	msg1, img1 := getMsgAndImgByRegion(&domain.Region{RegionID: 1, Tag: "PLAINS", Nerf: 20})
	if !strings.Contains(msg1, "Plan") {
		t.Fatalf("expected plains narrative, got %q", msg1)
	}
	if img1 == nil || img1.Path != "assets/npcs/lyra.png" {
		t.Fatalf("unexpected region1 image: %+v", img1)
	}

	msg2, img2 := getMsgAndImgByRegion(&domain.Region{RegionID: 2, Tag: "FOREST", Nerf: 30})
	if !strings.Contains(msg2, "Forest Spirit") {
		t.Fatalf("expected forest narrative, got %q", msg2)
	}
	if img2 == nil || img2.Path != "assets/npcs/nyssara_guardia_das_raizes.png" {
		t.Fatalf("unexpected region2 image: %+v", img2)
	}
}
