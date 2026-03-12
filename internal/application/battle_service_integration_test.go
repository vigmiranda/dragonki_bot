package application

import (
	"context"
	"math/rand"
	"strings"
	"testing"

	postgresadapter "github.com/vigmiranda/dragonki/internal/adapters/postgres"
	"github.com/vigmiranda/dragonki/internal/domain"
	integrationtest "github.com/vigmiranda/dragonki/internal/testutil/integration"
	"gorm.io/gorm"
)

type battleIntegrationDB struct {
	DB *gorm.DB
}

func newBattleIntegrationDB(t *testing.T) *battleIntegrationDB {
	t.Helper()

	testDB := integrationtest.NewPostgresDB(t, integrationtest.PostgresOptions{
		DatabaseNamePrefix: "dragonki_app_it",
	})

	return &battleIntegrationDB{DB: testDB.DB}
}

type battleRepoRecorder struct {
	logs    []*domain.BattleLog
	players []*domain.Player
}

func (r *battleRepoRecorder) SaveBattleLog(ctx context.Context, battleLog *domain.BattleLog) error {
	r.logs = append(r.logs, battleLog)
	return nil
}

func (r *battleRepoRecorder) SaveBattleLogAndPlayer(ctx context.Context, battleLog *domain.BattleLog, player *domain.Player) error {
	r.logs = append(r.logs, battleLog)
	r.players = append(r.players, player)
	return nil
}

func (r *battleRepoRecorder) GetWeeklyLeaderboard(metric string, limit int) ([]domain.WeeklyBattleRank, error) {
	return nil, nil
}

type battleMonsterRepoStub struct {
	monster *domain.Monster
}

func (s *battleMonsterRepoStub) ListByRegion(userID int64) ([]domain.Monster, error) {
	return []domain.Monster{*s.monster}, nil
}

func (s *battleMonsterRepoStub) GetByMonsterID(monsterID int64) (*domain.Monster, error) {
	return s.monster, nil
}

func (s *battleMonsterRepoStub) GetMonsterDrops(monsterID int) ([]domain.MonsterDrop, error) {
	return nil, nil
}

type battleRegionRepoStub struct {
	region *domain.Region
}

func (s *battleRegionRepoStub) GetByID(regionID int64) (*domain.Region, error) {
	return s.region, nil
}

func (s *battleRegionRepoStub) ListByLevel(level int64) ([]domain.Region, error) {
	return []domain.Region{*s.region}, nil
}

func TestBattleServiceIntegration_AutoEvolvesClassAt35And50(t *testing.T) {
	testDB := newBattleIntegrationDB(t)
	playerRepo := postgresadapter.NewPlayerPostgresRepo(testDB.DB)
	battleRepo := &battleRepoRecorder{}
	regionRepo := &battleRegionRepoStub{region: &domain.Region{RegionID: 1, Name: "Starter Plains", MinLevel: 1, Nerf: 10}}
	monsterRepo := &battleMonsterRepoStub{monster: &domain.Monster{
		MonsterID: 999,
		Name:      "Training Dummy",
		Tag:       "DUMMY",
		RegionID:  1,
		Level:     50,
		HP:        1,
		ATK:       1,
		DEF:       0,
		MinXP:     200000,
		MaxXP:     200000,
		MinMoney:  1,
		MaxMoney:  1,
	}}

	classes := []domain.Class{
		{ID: 201, Name: "Ki Vanguard", BaseClass: "Blade Monk", Tier: domain.TierAdvanced, MinLevel: 15, Skin: "ki_vanguard"},
		{ID: 202, Name: "Storm Vanguard", BaseClass: "Ki Vanguard", Tier: domain.TierElite, MinLevel: 35, Skin: "storm_vanguard"},
		{ID: 203, Name: "Tempest Avatar", BaseClass: "Storm Vanguard", Tier: domain.TierMaster, MinLevel: 50, Skin: "tempest_avatar"},
	}
	for _, class := range classes {
		if err := testDB.DB.Create(&class).Error; err != nil {
			t.Fatalf("create class fixture %s: %v", class.Name, err)
		}
	}

	svc := NewBattleService(playerRepo, monsterRepo, regionRepo, battleRepo, &fakeInventoryRepo{})
	svc.rng = randForIntegration()
	svc.cfg.MaxPlayerLevel = 50

	player35, err := playerRepo.GetOrCreate(740001)
	if err != nil {
		t.Fatalf("create player for level 35 evolution: %v", err)
	}
	player35.ClassID = uint(classes[0].ID)
	player35.Level = 34
	player35.XP = 0
	player35.HP = 100
	player35.MaxHP = 100
	player35.ATK = 100
	player35.DEF = 100
	player35.Speed = 20
	player35.Energy = 20
	player35.MaxEnergy = 20
	player35.OnboardingStep = "DONE"
	if err := playerRepo.Update(player35); err != nil {
		t.Fatalf("update player35 fixture: %v", err)
	}

	resp35, err := svc.Battle(player35.UserID, map[string]string{"monsterID": "999"})
	if err != nil {
		t.Fatalf("battle for level 35 evolution: %v", err)
	}
	if !strings.Contains(resp35.Text, "Evolucao automatica") {
		t.Fatalf("expected auto evolution message at level 35, got %q", resp35.Text)
	}

	reloaded35, err := playerRepo.GetOrCreate(player35.UserID)
	if err != nil {
		t.Fatalf("reload player35 after battle: %v", err)
	}
	if reloaded35.Level != 35 {
		t.Fatalf("expected player35 level to become 35, got %d", reloaded35.Level)
	}
	if reloaded35.ClassID != uint(classes[1].ID) || reloaded35.Class.Name != classes[1].Name {
		t.Fatalf("expected class to evolve to %s, got id=%d name=%q", classes[1].Name, reloaded35.ClassID, reloaded35.Class.Name)
	}

	player50, err := playerRepo.GetOrCreate(740002)
	if err != nil {
		t.Fatalf("create player for level 50 evolution: %v", err)
	}
	player50.ClassID = uint(classes[1].ID)
	player50.Level = 49
	player50.XP = 0
	player50.HP = 100
	player50.MaxHP = 100
	player50.ATK = 120
	player50.DEF = 110
	player50.Speed = 25
	player50.Energy = 20
	player50.MaxEnergy = 20
	player50.OnboardingStep = "DONE"
	if err := playerRepo.Update(player50); err != nil {
		t.Fatalf("update player50 fixture: %v", err)
	}

	resp50, err := svc.Battle(player50.UserID, map[string]string{"monsterID": "999"})
	if err != nil {
		t.Fatalf("battle for level 50 evolution: %v", err)
	}
	if !strings.Contains(resp50.Text, "Evolucao automatica") {
		t.Fatalf("expected auto evolution message at level 50, got %q", resp50.Text)
	}

	reloaded50, err := playerRepo.GetOrCreate(player50.UserID)
	if err != nil {
		t.Fatalf("reload player50 after battle: %v", err)
	}
	if reloaded50.Level != 50 {
		t.Fatalf("expected player50 level to become 50, got %d", reloaded50.Level)
	}
	if reloaded50.XP != 0 {
		t.Fatalf("expected player50 xp to reset to 0, got %d", reloaded50.XP)
	}
	if reloaded50.ClassID != uint(classes[2].ID) || reloaded50.Class.Name != classes[2].Name {
		t.Fatalf("expected class to evolve to %s, got id=%d name=%q", classes[2].Name, reloaded50.ClassID, reloaded50.Class.Name)
	}
}

func randForIntegration() *rand.Rand {
	return rand.New(rand.NewSource(1))
}
