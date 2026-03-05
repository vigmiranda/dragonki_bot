package application

import (
	"errors"
	"strings"
	"testing"

	"github.com/vigmiranda/dragonki/internal/domain"
)

func TestNewExchangeService(t *testing.T) {
	repo := &fakePlayerRepo{}
	svc := NewExchangeService(repo)
	if svc == nil {
		t.Fatalf("expected non-nil service")
	}
	if svc.PlayerRepository != repo {
		t.Fatalf("expected repo to be set")
	}
}

func TestExchangeService_Open(t *testing.T) {
	svc := ExchangeService{}
	resp, err := svc.Open()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "EXCHANGE_MENU" {
		t.Fatalf("expected EXCHANGE_MENU, got %q", resp.MenuKey)
	}
	if resp.Image == nil || resp.Image.Path != "assets/npcs/mestre_trikan.png" {
		t.Fatalf("unexpected image: %+v", resp.Image)
	}
	if !strings.Contains(resp.Text, "Casa de C") {
		t.Fatalf("unexpected text: %q", resp.Text)
	}
}

func TestExchangeService_Convert_GetOrCreateError(t *testing.T) {
	svc := ExchangeService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) {
				return nil, errors.New("db fail")
			},
		},
	}
	_, err := svc.Convert(1, map[string]string{"amountConvert": "1"})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestExchangeService_Convert_MissingAmountParam(t *testing.T) {
	svc := ExchangeService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return basePlayer(), nil },
		},
	}
	_, err := svc.Convert(1, map[string]string{})
	if err == nil || !strings.Contains(err.Error(), "missing param amountConvert") {
		t.Fatalf("expected missing param error, got %v", err)
	}
}

func TestExchangeService_Convert_InvalidAmountParam(t *testing.T) {
	svc := ExchangeService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return basePlayer(), nil },
		},
	}
	_, err := svc.Convert(1, map[string]string{"amountConvert": "abc"})
	if err == nil {
		t.Fatalf("expected parse error")
	}
}

func TestExchangeService_Convert_InvalidEnv(t *testing.T) {
	t.Setenv("KI_CRYSTALS_CONVERT", "abc")
	svc := ExchangeService{
		PlayerRepository: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return basePlayer(), nil },
		},
	}
	_, err := svc.Convert(1, map[string]string{"amountConvert": "2"})
	if err == nil {
		t.Fatalf("expected env parse error")
	}
}

func TestExchangeService_Convert_InsufficientCrystals(t *testing.T) {
	t.Setenv("KI_CRYSTALS_CONVERT", "100")
	player := basePlayer()
	player.KiCrystals = 1
	player.Gold = 50

	repo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	svc := ExchangeService{PlayerRepository: repo}

	resp, err := svc.Convert(1, map[string]string{"amountConvert": "5"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Text, "n") {
		t.Fatalf("expected insufficient message, got %q", resp.Text)
	}
	if player.KiCrystals != 1 || player.Gold != 50 {
		t.Fatalf("expected no player change, got crystals=%d gold=%d", player.KiCrystals, player.Gold)
	}
	if len(repo.updatedPlayers) != 0 {
		t.Fatalf("expected no update call")
	}
}

func TestExchangeService_Convert_Success(t *testing.T) {
	t.Setenv("KI_CRYSTALS_CONVERT", "200")
	player := basePlayer()
	player.Name = "Piccolo"
	player.KiCrystals = 10
	player.Gold = 100

	repo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	svc := ExchangeService{PlayerRepository: repo}

	resp, err := svc.Convert(1, map[string]string{"amountConvert": "3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if player.KiCrystals != 7 {
		t.Fatalf("expected crystals 7, got %d", player.KiCrystals)
	}
	if player.Gold != 700 {
		t.Fatalf("expected gold 700, got %d", player.Gold)
	}
	if len(repo.updatedPlayers) != 1 {
		t.Fatalf("expected one update call, got %d", len(repo.updatedPlayers))
	}
	if !strings.Contains(resp.Text, "Piccolo") {
		t.Fatalf("expected player name in success text, got %q", resp.Text)
	}
}

func TestExchangeService_Convert_UpdateError(t *testing.T) {
	t.Setenv("KI_CRYSTALS_CONVERT", "100")
	player := basePlayer()
	player.KiCrystals = 5

	repo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		updateFn: func(player *domain.Player) error {
			return errors.New("save fail")
		},
	}
	svc := ExchangeService{PlayerRepository: repo}

	_, err := svc.Convert(1, map[string]string{"amountConvert": "2"})
	if err == nil {
		t.Fatalf("expected update error")
	}
}
