package postgres

import (
	"testing"
	"time"
)

func TestPlayerPostgresRepoIntegration_PersistsAndReadsFromDatabase(t *testing.T) {
	testDB := newIntegrationTestDB(t)
	repo := NewPlayerPostgresRepo(testDB.DB)
	userID := time.Now().UnixNano()

	player, err := repo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("GetOrCreate should create player: %v", err)
	}

	if player.ID == 0 {
		t.Fatal("expected persisted player with generated id")
	}

	player.Name = "Integration Hero"
	player.Level = 7
	player.Gold = 321
	player.KiCrystals = 9
	player.Energy = 15
	player.UserName = "integration_hero"
	player.NameTelegram = "Integration Hero TG"
	player.UserIDTelegram = 99887766

	if err := repo.Update(player); err != nil {
		t.Fatalf("Update should persist player changes: %v", err)
	}

	samePlayer, err := repo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("GetOrCreate should read existing player: %v", err)
	}

	if samePlayer.ID != player.ID {
		t.Fatalf("expected same player id, got %d want %d", samePlayer.ID, player.ID)
	}
	if samePlayer.Level != 7 || samePlayer.Gold != 321 || samePlayer.KiCrystals != 9 {
		t.Fatalf("unexpected repo state after update: level=%d gold=%d ki=%d", samePlayer.Level, samePlayer.Gold, samePlayer.KiCrystals)
	}

	var persisted struct {
		Name         string
		Level        int
		Gold         int
		KiCrystals   int
		Energy       int
		UserName     string
		NameTelegram string
	}

	err = testDB.DB.
		Table("players").
		Select("name, level, gold, ki_crystals, energy, user_name, name_telegram").
		Where("id = ?", player.ID).
		Take(&persisted).
		Error
	if err != nil {
		t.Fatalf("query persisted player row: %v", err)
	}

	if persisted.Name != "Integration Hero" {
		t.Fatalf("unexpected persisted name: %q", persisted.Name)
	}
	if persisted.Level != 7 || persisted.Gold != 321 || persisted.KiCrystals != 9 || persisted.Energy != 15 {
		t.Fatalf("unexpected persisted numeric state: %+v", persisted)
	}
	if persisted.UserName != "integration_hero" || persisted.NameTelegram != "Integration Hero TG" {
		t.Fatalf("unexpected persisted telegram data: %+v", persisted)
	}
}
