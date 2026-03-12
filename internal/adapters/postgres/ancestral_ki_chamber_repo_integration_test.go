package postgres

import (
	"errors"
	"testing"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
)

func TestAncestralKiChamberPostgresRepoIntegration_StartRunIfEligible(t *testing.T) {
	testDB := newIntegrationTestDB(t)
	chamberRepo := NewAncestralKiChamberPostgresRepo(testDB.DB)
	playerRepo := NewPlayerPostgresRepo(testDB.DB)

	userID := time.Now().UnixNano()
	player, err := playerRepo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("create player fixture: %v", err)
	}

	player.Energy = 7
	player.Gold = 900
	if err := playerRepo.Update(player); err != nil {
		t.Fatalf("persist player resources: %v", err)
	}

	startedAt := time.Now().UTC()
	run, err := chamberRepo.StartRunIfEligible(int64(player.ID), 2, 1, 150, startedAt)
	if err != nil {
		t.Fatalf("start run atomically: %v", err)
	}
	if run == nil || run.ID == 0 {
		t.Fatalf("expected persisted run, got %+v", run)
	}
	if run.PlayerID != int64(player.ID) || run.ChamberLevel != 2 {
		t.Fatalf("unexpected run persisted: %+v", run)
	}
	if run.Status != domain.AncestralKiChamberRunStatusActive {
		t.Fatalf("expected active run status, got %s", run.Status)
	}

	reloadedPlayer, err := playerRepo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("reload player after start run: %v", err)
	}
	if reloadedPlayer.Energy != 6 || reloadedPlayer.Gold != 750 {
		t.Fatalf("unexpected resource debit after start run: energy=%d gold=%d", reloadedPlayer.Energy, reloadedPlayer.Gold)
	}

	activeRun, err := chamberRepo.GetActiveRunByPlayer(int64(player.ID))
	if err != nil {
		t.Fatalf("load active run: %v", err)
	}
	if activeRun.ID != run.ID {
		t.Fatalf("expected same active run id, got %d want %d", activeRun.ID, run.ID)
	}

	secondRun, err := chamberRepo.StartRunIfEligible(int64(player.ID), 2, 1, 150, time.Now().UTC())
	if !errors.Is(err, domain.ErrAncestralKiChamberRunActive) {
		t.Fatalf("expected active run error on retry, got run=%+v err=%v", secondRun, err)
	}

	reloadedPlayer, err = playerRepo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("reload player after retry: %v", err)
	}
	if reloadedPlayer.Energy != 6 || reloadedPlayer.Gold != 750 {
		t.Fatalf("expected retry to keep previous resources, got energy=%d gold=%d", reloadedPlayer.Energy, reloadedPlayer.Gold)
	}

	var runCount int64
	if err := testDB.DB.
		Table("ancestral_ki_chamber_runs").
		Where("player_id = ?", player.ID).
		Count(&runCount).Error; err != nil {
		t.Fatalf("count chamber runs: %v", err)
	}
	if runCount != 1 {
		t.Fatalf("expected one persisted run, got %d", runCount)
	}
}

func TestAncestralKiChamberPostgresRepoIntegration_StartRunIfEligible_RollsBackOnInsufficientResources(t *testing.T) {
	testDB := newIntegrationTestDB(t)
	chamberRepo := NewAncestralKiChamberPostgresRepo(testDB.DB)
	playerRepo := NewPlayerPostgresRepo(testDB.DB)

	userID := time.Now().UnixNano()
	player, err := playerRepo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("create player fixture: %v", err)
	}

	player.Energy = 0
	player.Gold = 149
	if err := playerRepo.Update(player); err != nil {
		t.Fatalf("persist low resources: %v", err)
	}

	run, err := chamberRepo.StartRunIfEligible(int64(player.ID), 1, 1, 150, time.Now().UTC())
	if !errors.Is(err, domain.ErrAncestralKiChamberInsufficientResource) {
		t.Fatalf("expected insufficient resource error, got run=%+v err=%v", run, err)
	}

	reloadedPlayer, err := playerRepo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("reload player after failed start run: %v", err)
	}
	if reloadedPlayer.Energy != 0 || reloadedPlayer.Gold != 149 {
		t.Fatalf("expected failed start run to keep resources, got energy=%d gold=%d", reloadedPlayer.Energy, reloadedPlayer.Gold)
	}

	var runCount int64
	if err := testDB.DB.
		Table("ancestral_ki_chamber_runs").
		Where("player_id = ?", player.ID).
		Count(&runCount).Error; err != nil {
		t.Fatalf("count chamber runs after failure: %v", err)
	}
	if runCount != 0 {
		t.Fatalf("expected no chamber run after insufficient resources, got %d", runCount)
	}
}

func TestAncestralKiChamberPostgresRepoIntegration_ResolveBattleResult_CompleteRun(t *testing.T) {
	testDB := newIntegrationTestDB(t)
	chamberRepo := NewAncestralKiChamberPostgresRepo(testDB.DB)
	playerRepo := NewPlayerPostgresRepo(testDB.DB)

	userID := time.Now().UnixNano()
	player, err := playerRepo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("create player fixture: %v", err)
	}

	player.XP = 25
	player.Gold = 500
	if err := playerRepo.Update(player); err != nil {
		t.Fatalf("persist player state: %v", err)
	}

	startedAt := time.Now().UTC().Add(-5 * time.Minute)
	run := &domain.AncestralKiChamberRun{
		PlayerID:        int64(player.ID),
		ChamberLevel:    1,
		CurrentBattle:   5,
		Victories:       4,
		AccumulatedXP:   100,
		AccumulatedGold: 200,
		Status:          domain.AncestralKiChamberRunStatusActive,
		StartedAt:       startedAt,
		CreatedAt:       startedAt,
		UpdatedAt:       startedAt,
	}
	if err := chamberRepo.CreateRun(run); err != nil {
		t.Fatalf("create active chamber run: %v", err)
	}

	resolution, err := chamberRepo.ResolveBattleResult(int64(player.ID), 777, 20, true, 50, 25, 5, 0.40, time.Now().UTC())
	if err != nil {
		t.Fatalf("resolve chamber run completion: %v", err)
	}
	if resolution == nil || resolution.Run == nil {
		t.Fatalf("expected resolution with run, got %+v", resolution)
	}
	if resolution.Run.Status != domain.AncestralKiChamberRunStatusCompleted {
		t.Fatalf("expected completed run, got %s", resolution.Run.Status)
	}
	if resolution.Run.Victories != 5 || resolution.Run.AccumulatedXP != 150 || resolution.Run.AccumulatedGold != 225 {
		t.Fatalf("unexpected completed run totals: %+v", resolution.Run)
	}
	if resolution.XPPayout != 60 || resolution.GoldPayout != 90 {
		t.Fatalf("unexpected payout totals: xp=%d gold=%d", resolution.XPPayout, resolution.GoldPayout)
	}
	if resolution.Progress == nil || resolution.Progress.HighestCompletedLevel != 1 {
		t.Fatalf("expected unlocked progress, got %+v", resolution.Progress)
	}

	reloadedPlayer, err := playerRepo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("reload player after completion: %v", err)
	}
	if reloadedPlayer.XP != 85 || reloadedPlayer.Gold != 590 {
		t.Fatalf("unexpected player payout after completion: xp=%d gold=%d", reloadedPlayer.XP, reloadedPlayer.Gold)
	}

	reloadedRun, err := chamberRepo.GetActiveRunByPlayer(int64(player.ID))
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected no active run after completion, got run=%+v err=%v", reloadedRun, err)
	}

	var persistedRun domain.AncestralKiChamberRun
	if err := testDB.DB.
		Where("id = ?", run.ID).
		Take(&persistedRun).Error; err != nil {
		t.Fatalf("load persisted chamber run: %v", err)
	}
	if persistedRun.Status != domain.AncestralKiChamberRunStatusCompleted {
		t.Fatalf("expected persisted completed run, got %s", persistedRun.Status)
	}

	var battleLog domain.AncestralKiChamberBattleLog
	if err := testDB.DB.
		Where("chamber_run_id = ?", run.ID).
		Take(&battleLog).Error; err != nil {
		t.Fatalf("load chamber battle log: %v", err)
	}
	if battleLog.Result != "WIN" || battleLog.BattleNumber != 5 {
		t.Fatalf("unexpected battle log persisted: %+v", battleLog)
	}
}

func TestAncestralKiChamberPostgresRepoIntegration_ResolveBattleResult_FailedRun(t *testing.T) {
	testDB := newIntegrationTestDB(t)
	chamberRepo := NewAncestralKiChamberPostgresRepo(testDB.DB)
	playerRepo := NewPlayerPostgresRepo(testDB.DB)

	userID := time.Now().UnixNano()
	player, err := playerRepo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("create player fixture: %v", err)
	}

	player.XP = 80
	player.Gold = 700
	if err := playerRepo.Update(player); err != nil {
		t.Fatalf("persist player state: %v", err)
	}

	startedAt := time.Now().UTC().Add(-5 * time.Minute)
	run := &domain.AncestralKiChamberRun{
		PlayerID:      int64(player.ID),
		ChamberLevel:  2,
		CurrentBattle: 3,
		Victories:     2,
		Status:        domain.AncestralKiChamberRunStatusActive,
		StartedAt:     startedAt,
		CreatedAt:     startedAt,
		UpdatedAt:     startedAt,
	}
	if err := chamberRepo.CreateRun(run); err != nil {
		t.Fatalf("create active chamber run: %v", err)
	}

	resolution, err := chamberRepo.ResolveBattleResult(int64(player.ID), 888, 22, false, 120, 240, 5, 0.40, time.Now().UTC())
	if err != nil {
		t.Fatalf("resolve failed chamber run: %v", err)
	}
	if resolution == nil || resolution.Run == nil {
		t.Fatalf("expected resolution with run, got %+v", resolution)
	}
	if resolution.Run.Status != domain.AncestralKiChamberRunStatusFailed {
		t.Fatalf("expected failed run, got %s", resolution.Run.Status)
	}
	if resolution.XPPayout != 0 || resolution.GoldPayout != 0 {
		t.Fatalf("expected no payout on failed run, got xp=%d gold=%d", resolution.XPPayout, resolution.GoldPayout)
	}

	reloadedPlayer, err := playerRepo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("reload player after failed run: %v", err)
	}
	if reloadedPlayer.XP != 80 || reloadedPlayer.Gold != 700 {
		t.Fatalf("expected failed run to keep player rewards unchanged, got xp=%d gold=%d", reloadedPlayer.XP, reloadedPlayer.Gold)
	}

	var battleLog domain.AncestralKiChamberBattleLog
	if err := testDB.DB.
		Where("chamber_run_id = ?", run.ID).
		Take(&battleLog).Error; err != nil {
		t.Fatalf("load failed battle log: %v", err)
	}
	if battleLog.Result != "LOSE" || battleLog.BattleNumber != 3 {
		t.Fatalf("unexpected failed battle log persisted: %+v", battleLog)
	}
}

func TestAncestralKiChamberPostgresRepoIntegration_AbandonActiveRun(t *testing.T) {
	testDB := newIntegrationTestDB(t)
	chamberRepo := NewAncestralKiChamberPostgresRepo(testDB.DB)
	playerRepo := NewPlayerPostgresRepo(testDB.DB)

	userID := time.Now().UnixNano()
	player, err := playerRepo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("create player fixture: %v", err)
	}

	startedAt := time.Now().UTC().Add(-3 * time.Minute)
	run := &domain.AncestralKiChamberRun{
		PlayerID:      int64(player.ID),
		ChamberLevel:  3,
		CurrentBattle: 2,
		Victories:     1,
		Status:        domain.AncestralKiChamberRunStatusActive,
		StartedAt:     startedAt,
		CreatedAt:     startedAt,
		UpdatedAt:     startedAt,
	}
	if err := chamberRepo.CreateRun(run); err != nil {
		t.Fatalf("create active run fixture: %v", err)
	}

	abandonedAt := time.Now().UTC()
	abandonedRun, err := chamberRepo.AbandonActiveRun(int64(player.ID), abandonedAt)
	if err != nil {
		t.Fatalf("abandon active run atomically: %v", err)
	}
	if abandonedRun == nil {
		t.Fatal("expected abandoned run")
	}
	if abandonedRun.Status != domain.AncestralKiChamberRunStatusAbandoned {
		t.Fatalf("expected abandoned status, got %s", abandonedRun.Status)
	}
	if abandonedRun.FinishedAt == nil {
		t.Fatal("expected finished_at to be set")
	}

	var persistedRun domain.AncestralKiChamberRun
	if err := testDB.DB.Where("id = ?", run.ID).Take(&persistedRun).Error; err != nil {
		t.Fatalf("reload abandoned run: %v", err)
	}
	if persistedRun.Status != domain.AncestralKiChamberRunStatusAbandoned {
		t.Fatalf("expected persisted abandoned status, got %s", persistedRun.Status)
	}
	if persistedRun.FinishedAt == nil {
		t.Fatal("expected persisted finished_at to be set")
	}

	activeRun, err := chamberRepo.GetActiveRunByPlayer(int64(player.ID))
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected no active run after abandon, got run=%+v err=%v", activeRun, err)
	}
}

func TestAncestralKiChamberPostgresRepoIntegration_AbandonActiveRun_NoActiveRun(t *testing.T) {
	testDB := newIntegrationTestDB(t)
	chamberRepo := NewAncestralKiChamberPostgresRepo(testDB.DB)
	playerRepo := NewPlayerPostgresRepo(testDB.DB)

	userID := time.Now().UnixNano()
	player, err := playerRepo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("create player fixture: %v", err)
	}

	run, err := chamberRepo.AbandonActiveRun(int64(player.ID), time.Now().UTC())
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record not found without active run, got run=%+v err=%v", run, err)
	}
}
