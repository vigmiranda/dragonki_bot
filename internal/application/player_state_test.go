package application

import (
	"errors"
	"testing"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
)

func TestApplyPassivePlayerUpdates_ErrorOnUpdate(t *testing.T) {
	player := basePlayer()
	repo := &fakePlayerRepo{
		updateFn: func(player *domain.Player) error {
			return errors.New("update fail")
		},
	}

	err := applyPassivePlayerUpdates(repo, player)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestApplyPassivePlayerUpdates_NoRegenWhenFullButResetForm(t *testing.T) {
	t.Setenv("ENERGY_REGEN_TIME", "10m")
	last := time.Now().UTC().Add(-30 * time.Minute)
	player := basePlayer()
	player.Energy = player.MaxEnergy
	player.HP = player.MaxHP
	player.Form = "Ki Awakened"
	player.FormBattlesLeft = 0
	player.LastEnergyAt = last

	repo := &fakePlayerRepo{}
	err := applyPassivePlayerUpdates(repo, player)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if player.Energy != player.MaxEnergy || player.HP != player.MaxHP {
		t.Fatalf("expected no regen changes")
	}
	if !player.LastEnergyAt.Equal(last) {
		t.Fatalf("expected LastEnergyAt unchanged")
	}
	if player.Form != "Base" {
		t.Fatalf("expected form reset to Base, got %q", player.Form)
	}
	if len(repo.updatedPlayers) != 1 {
		t.Fatalf("expected one update call")
	}
}

func TestApplyPassivePlayerUpdates_RegenWhenNeeded(t *testing.T) {
	t.Setenv("ENERGY_REGEN_TIME", "10m")
	player := basePlayer()
	player.Energy = 1
	player.MaxEnergy = 10
	player.HP = player.MaxHP
	player.LastEnergyAt = time.Now().UTC().Add(-25 * time.Minute)
	player.Form = "Ki Awakened"
	player.FormBattlesLeft = 2

	repo := &fakePlayerRepo{}
	err := applyPassivePlayerUpdates(repo, player)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if player.Energy != 3 {
		t.Fatalf("expected energy recovered to 3, got %d", player.Energy)
	}
	// Form should remain because FormBattlesLeft > 0.
	if player.Form != "Ki Awakened" {
		t.Fatalf("expected form unchanged, got %q", player.Form)
	}
	if len(repo.updatedPlayers) != 1 {
		t.Fatalf("expected one update call")
	}
}
