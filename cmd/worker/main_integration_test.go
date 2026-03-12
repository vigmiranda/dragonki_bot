package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/vigmiranda/dragonki/internal/adapters/postgres"
	"github.com/vigmiranda/dragonki/internal/adapters/queue"
	"github.com/vigmiranda/dragonki/internal/application"
	"github.com/vigmiranda/dragonki/internal/domain"
	integrationtest "github.com/vigmiranda/dragonki/internal/testutil/integration"
)

func TestBuildWorkerAppIntegration_ProcessesBigBossRewardQueue(t *testing.T) {
	testDB := integrationtest.NewPostgresDB(t, integrationtest.PostgresOptions{
		DatabaseNamePrefix: "dragonki_worker_it",
	})
	redisInstance := integrationtest.NewRedisInstance(t)

	playerRepo := postgres.NewPlayerPostgresRepo(testDB.DB)
	inventoryRepo := postgres.NewInventoryPostgresRepo(testDB.DB)
	weeklyBossRepo := postgres.NewWeeklyBossPostgresRepo(testDB.DB)

	player, err := playerRepo.GetOrCreate(940001)
	if err != nil {
		t.Fatalf("create player: %v", err)
	}
	player.Gold = 1000
	player.Energy = 20
	player.MaxEnergy = 20
	player.HP = 100
	player.MaxHP = 100
	player.ATK = 50
	player.DEF = 0
	player.Speed = 0
	player.OnboardingStep = "DONE"
	if err := playerRepo.Update(player); err != nil {
		t.Fatalf("update player fixture: %v", err)
	}
	initialGold := player.Gold

	now := time.Now().UTC()
	year, week := now.ISOWeek()
	boss := &domain.WeeklyBoss{
		Name:      "Integration Big Boss",
		Tag:       "IT_BIG_BOSS",
		WeekKey:   fmt.Sprintf("%04d-W%02d", year, week),
		MaxHP:     10,
		CurrentHP: 10,
		ATK:       1,
		Status:    domain.WeeklyBossStatusActive,
		StartsAt:  now.Add(-1 * time.Hour),
		EndsAt:    now.Add(24 * time.Hour),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := weeklyBossRepo.Create(boss); err != nil {
		t.Fatalf("create boss fixture: %v", err)
	}

	reward := &domain.WeeklyBossReward{
		Position:  1,
		Gold:      777,
		Enabled:   true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := testDB.DB.Create(reward).Error; err != nil {
		t.Fatalf("create reward fixture: %v", err)
	}

	workerApp, err := buildWorkerApp(workerAppConfig{
		DB:          testDB.DB,
		RedisURL:    redisInstance.URL,
		Concurrency: 1,
	})
	if err != nil {
		t.Fatalf("build worker app: %v", err)
	}
	defer workerApp.Cleanup()
	if err := workerApp.Server.Start(workerApp.Handler); err != nil {
		t.Fatalf("start worker app: %v", err)
	}
	defer workerApp.Server.Shutdown()

	rewardQueue, err := queue.NewAsynqBigBossQueue(redisInstance.URL, "")
	if err != nil {
		t.Fatalf("create reward queue: %v", err)
	}
	defer rewardQueue.Close()

	service := application.NewWeeklyBossService(weeklyBossRepo, playerRepo, inventoryRepo, rewardQueue)

	response, err := service.Attack(player.UserID)
	if err != nil {
		t.Fatalf("attack big boss: %v", err)
	}
	if !strings.Contains(response.Text, "Premiacao do Top 3 em processamento") {
		t.Fatalf("expected async reward message, got %q", response.Text)
	}

	waitForWorkerCondition(t, 15*time.Second, func() bool {
		var payout domain.WeeklyBossPayout
		if err := testDB.DB.Where("boss_id = ?", boss.ID).First(&payout).Error; err != nil {
			return false
		}
		if payout.Status != domain.WeeklyBossPayoutPaid {
			return false
		}

		reloadedPlayer, reloadErr := playerRepo.GetOrCreate(player.UserID)
		if reloadErr != nil {
			return false
		}
		expectedGold := initialGold - 200 + int(reward.Gold)
		return reloadedPlayer.Gold == expectedGold
	})

	var payout domain.WeeklyBossPayout
	if err := testDB.DB.Where("boss_id = ?", boss.ID).First(&payout).Error; err != nil {
		t.Fatalf("load payout: %v", err)
	}
	if payout.Status != domain.WeeklyBossPayoutPaid {
		t.Fatalf("unexpected payout status: %q", payout.Status)
	}
	if payout.Position != 1 || payout.PlayerID != int64(player.ID) || payout.GoldPaid != reward.Gold {
		t.Fatalf("unexpected payout row: %+v", payout)
	}

	reloadedBoss, err := weeklyBossRepo.GetByID(boss.ID)
	if err != nil {
		t.Fatalf("reload boss: %v", err)
	}
	if reloadedBoss.Status != domain.WeeklyBossStatusDefeated || reloadedBoss.CurrentHP != 0 {
		t.Fatalf("unexpected boss state after defeat: %+v", reloadedBoss)
	}

	reloadedPlayer, err := playerRepo.GetOrCreate(player.UserID)
	if err != nil {
		t.Fatalf("reload player after payout: %v", err)
	}
	expectedGold := initialGold - 200 + int(reward.Gold)
	if reloadedPlayer.Gold != expectedGold {
		t.Fatalf("unexpected player gold after payout: got %d want %d", reloadedPlayer.Gold, expectedGold)
	}
}

func waitForWorkerCondition(t *testing.T, timeout time.Duration, fn func() bool) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if fn() {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	t.Fatalf("condition not met within %s", timeout)
}
