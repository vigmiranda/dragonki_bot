package postgres

import (
	"testing"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
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

func TestPlayerPostgresRepoIntegration_ClassEvolutionChain_Level15To35To50(t *testing.T) {
	testDB := newIntegrationTestDB(t)
	repo := NewPlayerPostgresRepo(testDB.DB)

	classes := []domain.Class{
		{ID: 101, Name: "Blade Monk", BaseClass: "Blade Monk", Tier: domain.TierBase, MinLevel: 1, Skin: "basic"},
		{ID: 102, Name: "Ki Vanguard", BaseClass: "Blade Monk", Tier: domain.TierAdvanced, MinLevel: 15, Skin: "ki_vanguard"},
		{ID: 103, Name: "Wind Striker", BaseClass: "Blade Monk", Tier: domain.TierAdvanced, MinLevel: 15, Skin: "wind_striker"},
		{ID: 104, Name: "Earth Sentinel", BaseClass: "Blade Monk", Tier: domain.TierAdvanced, MinLevel: 15, Skin: "earth_sentinel"},
		{ID: 105, Name: "Storm Vanguard", BaseClass: "Ki Vanguard", Tier: domain.TierElite, MinLevel: 35, Skin: "storm_vanguard"},
		{ID: 106, Name: "Tempest Avatar", BaseClass: "Storm Vanguard", Tier: domain.TierMaster, MinLevel: 50, Skin: "tempest_avatar"},
	}

	for _, class := range classes {
		if err := testDB.DB.Create(&class).Error; err != nil {
			t.Fatalf("create class fixture %s: %v", class.Name, err)
		}
	}

	level15Classes, err := repo.GetClassesByLvl(15)
	if err != nil {
		t.Fatalf("get classes by level 15: %v", err)
	}
	if len(level15Classes) != 3 {
		t.Fatalf("expected 3 class choices at level 15, got %d", len(level15Classes))
	}

	advancedClass, err := repo.FindClassEvolutionByBaseClassAndLevel("Blade Monk", 15)
	if err != nil {
		t.Fatalf("find advanced class evolution: %v", err)
	}
	if advancedClass.Name != "Ki Vanguard" {
		t.Fatalf("unexpected first advanced evolution at level 15: %q", advancedClass.Name)
	}

	level35Class, err := repo.FindClassEvolutionByBaseClassAndLevel("Ki Vanguard", 35)
	if err != nil {
		t.Fatalf("find elite class evolution: %v", err)
	}
	if level35Class.Name != "Storm Vanguard" {
		t.Fatalf("unexpected evolution at level 35: %q", level35Class.Name)
	}

	level50Class, err := repo.FindClassEvolutionByBaseClassAndLevel("Storm Vanguard", 50)
	if err != nil {
		t.Fatalf("find master class evolution: %v", err)
	}
	if level50Class.Name != "Tempest Avatar" {
		t.Fatalf("unexpected evolution at level 50: %q", level50Class.Name)
	}
}
