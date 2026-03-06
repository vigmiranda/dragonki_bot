package application

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
)

func TestNewPlayerService(t *testing.T) {
	playerRepo := &fakePlayerRepo{}
	regionRepo := &fakeRegionRepo{}
	invRepo := &fakeInventoryRepo{}

	svc := NewPlayerService(playerRepo, regionRepo, invRepo)
	if svc == nil {
		t.Fatalf("expected non-nil service")
	}
	if svc.PlayerRepository != playerRepo {
		t.Fatalf("expected player repo to be set")
	}
	if svc.RegionRepository != regionRepo {
		t.Fatalf("expected region repo to be set")
	}
	if svc.InventoryRepository != invRepo {
		t.Fatalf("expected inventory repo to be set")
	}
}

func TestPlayerService_Status_PlayerRepoError(t *testing.T) {
	svc := &PlayerService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) {
				return nil, errors.New("db fail")
			},
		},
		RegionRepository: &fakeRegionRepo{},
	}

	_, err := svc.Status(1)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestPlayerService_Status_RegionRepoError(t *testing.T) {
	player := basePlayer()
	svc := &PlayerService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
		RegionRepository: &fakeRegionRepo{
			getByIDFn: func(regionID int64) (*domain.Region, error) {
				return nil, errors.New("region fail")
			},
		},
	}

	_, err := svc.Status(1)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestPlayerService_Status_ApplyPassiveUpdateError(t *testing.T) {
	player := basePlayer()
	svc := &PlayerService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
			updateFn: func(player *domain.Player) error {
				return errors.New("update fail")
			},
		},
		RegionRepository: &fakeRegionRepo{
			getByIDFn: func(regionID int64) (*domain.Region, error) {
				return &domain.Region{RegionID: 1, Name: "Plains", MinLevel: 1}, nil
			},
		},
	}

	_, err := svc.Status(1)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestPlayerService_Status_SuccessDefaultRegenAndMainMenu(t *testing.T) {
	resetAppConfigForTest()
	t.Setenv("ENERGY_REGEN_TIME", "")

	player := basePlayer()
	player.Name = "Goku"
	player.Gender = "MAN"
	player.Class.Name = "Blade Monk"
	player.Class.Skin = "basic"
	player.HP = player.MaxHP
	player.Energy = player.MaxEnergy
	player.FormBattlesLeft = 1
	player.LastEnergyAt = time.Now().UTC()

	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	regionRepo := &fakeRegionRepo{
		getByIDFn: func(regionID int64) (*domain.Region, error) {
			return &domain.Region{RegionID: 1, Name: "Planicie", MinLevel: 1}, nil
		},
	}
	svc := &PlayerService{
		PlayerRepository: playerRepo,
		RegionRepository: regionRepo,
	}

	resp, err := svc.Status(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "MAIN_MENU" {
		t.Fatalf("expected MAIN_MENU, got %q", resp.MenuKey)
	}
	if resp.Image == nil || !strings.Contains(resp.Image.Path, "assets/skins/man/basicplayer.png") {
		t.Fatalf("unexpected image path: %+v", resp.Image)
	}
	if !strings.Contains(resp.Text, "10m") {
		t.Fatalf("expected default regen text 10m, got %q", resp.Text)
	}
	if len(playerRepo.updatedPlayers) != 1 {
		t.Fatalf("expected one passive update save, got %d", len(playerRepo.updatedPlayers))
	}
}

func TestPlayerService_Status_WomanIconAndKiAwakenedImage(t *testing.T) {
	resetAppConfigForTest()
	t.Setenv("ENERGY_REGEN_TIME", "7m")

	player := basePlayer()
	player.Gender = "WOMAN"
	player.Form = "Ki Awakened"
	player.FormBattlesLeft = 2
	player.Class.Skin = "basic"

	svc := &PlayerService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
		RegionRepository: &fakeRegionRepo{
			getByIDFn: func(regionID int64) (*domain.Region, error) {
				return &domain.Region{RegionID: 1, Name: "Planicie", MinLevel: 1}, nil
			},
		},
	}

	resp, err := svc.Status(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Text, "<b>Bonus Equip</b>") {
		t.Fatalf("expected equipment bonus line in text, got %q", resp.Text)
	}
	if resp.Image == nil || !strings.Contains(resp.Image.Path, "ki_awakened.png") {
		t.Fatalf("expected awakened image, got %+v", resp.Image)
	}
}

func TestPlayerService_UpdateClass_GetOrCreateError(t *testing.T) {
	svc := &PlayerService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) {
				return nil, errors.New("db fail")
			},
		},
	}
	_, err := svc.UpdateClass(1, map[string]string{"classID": "1"})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestPlayerService_UpdateClass_InvalidClassID(t *testing.T) {
	svc := &PlayerService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return basePlayer(), nil },
		},
	}
	_, err := svc.UpdateClass(1, map[string]string{"classID": "abc"})
	if err == nil {
		t.Fatalf("expected parse error")
	}
}

func TestPlayerService_UpdateClass_GetClassError(t *testing.T) {
	svc := &PlayerService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return basePlayer(), nil },
			getClassesByIDFn: func(id int64) (*domain.Class, error) {
				return nil, errors.New("class fail")
			},
		},
	}
	_, err := svc.UpdateClass(1, map[string]string{"classID": "1"})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestPlayerService_UpdateClass_UpdateError(t *testing.T) {
	svc := &PlayerService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return basePlayer(), nil },
			getClassesByIDFn: func(id int64) (*domain.Class, error) {
				return &domain.Class{ID: id, Name: "Wind Striker"}, nil
			},
			updateFn: func(player *domain.Player) error {
				return errors.New("update fail")
			},
		},
	}
	_, err := svc.UpdateClass(1, map[string]string{"classID": "2"})
	if err == nil {
		t.Fatalf("expected update error")
	}
}

func TestPlayerService_UpdateClass_Level15AppliesStatsAndReturnsMainMenu(t *testing.T) {
	player := basePlayer()
	player.Level = 15
	player.Gender = "WOMAN"
	player.Class.Skin = "earth_sentinel"

	repo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		getClassesByIDFn: func(id int64) (*domain.Class, error) {
			return &domain.Class{ID: id, Name: "Wind Striker"}, nil
		},
	}
	svc := &PlayerService{PlayerRepository: repo}

	resp, err := svc.UpdateClass(1, map[string]string{"classID": "3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "MAIN_MENU" {
		t.Fatalf("expected MAIN_MENU, got %q", resp.MenuKey)
	}
	if !strings.Contains(resp.Text, "Wind Striker") {
		t.Fatalf("expected class name in message, got %q", resp.Text)
	}
	if resp.Image == nil || !strings.Contains(resp.Image.Path, "assets/skins/woman/earth_sentinelplayer.png") {
		t.Fatalf("unexpected image path: %+v", resp.Image)
	}
	if player.ClassID != 3 {
		t.Fatalf("expected class id 3, got %d", player.ClassID)
	}
	if player.ATK != 110 || player.DEF != 50 || player.Speed != 50 || player.MaxHP != 160 || player.MaxEnergy != 23 {
		t.Fatalf("expected lvl15 stat template applied, got ATK=%d DEF=%d SPD=%d HP=%d EN=%d",
			player.ATK, player.DEF, player.Speed, player.MaxHP, player.MaxEnergy)
	}
}
