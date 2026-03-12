package application

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
)

type chamberRepoMock struct {
	progress        *domain.AncestralKiChamberProgress
	activeRun       *domain.AncestralKiChamberRun
	battleLogs      []*domain.AncestralKiChamberBattleLog
	startRunErr     error
	startRunCalls   int
	lastStartRunAt  time.Time
	lastStartPlayer int64
	lastStartLevel  int
	resolveErr      error
	resolveCalls    int
	abandonErr      error
	abandonCalls    int
}

func (m *chamberRepoMock) GetOrCreateProgress(playerID int64) (*domain.AncestralKiChamberProgress, error) {
	if m.progress == nil {
		m.progress = &domain.AncestralKiChamberProgress{
			ID:                    1,
			PlayerID:              playerID,
			HighestCompletedLevel: 0,
		}
	}
	return m.progress, nil
}

func (m *chamberRepoMock) UpdateProgress(progress *domain.AncestralKiChamberProgress) error {
	m.progress = progress
	return nil
}

func (m *chamberRepoMock) GetActiveRunByPlayer(playerID int64) (*domain.AncestralKiChamberRun, error) {
	if m.activeRun == nil || m.activeRun.PlayerID != playerID || m.activeRun.Status != domain.AncestralKiChamberRunStatusActive {
		return nil, gorm.ErrRecordNotFound
	}
	return m.activeRun, nil
}

func (m *chamberRepoMock) StartRunIfEligible(playerID int64, chamberLevel int, energyCost int, goldCost int, startedAt time.Time) (*domain.AncestralKiChamberRun, error) {
	m.startRunCalls++
	m.lastStartRunAt = startedAt
	m.lastStartPlayer = playerID
	m.lastStartLevel = chamberLevel
	_ = energyCost
	_ = goldCost

	if m.startRunErr != nil {
		return nil, m.startRunErr
	}
	if m.activeRun != nil && m.activeRun.PlayerID == playerID && m.activeRun.Status == domain.AncestralKiChamberRunStatusActive {
		return nil, domain.ErrAncestralKiChamberRunActive
	}

	run := &domain.AncestralKiChamberRun{
		ID:            1,
		PlayerID:      playerID,
		ChamberLevel:  chamberLevel,
		CurrentBattle: 1,
		Victories:     0,
		Status:        domain.AncestralKiChamberRunStatusActive,
		StartedAt:     startedAt,
		CreatedAt:     startedAt,
		UpdatedAt:     startedAt,
	}
	m.activeRun = run
	return run, nil
}

func (m *chamberRepoMock) CreateRun(run *domain.AncestralKiChamberRun) error {
	run.ID = 1
	m.activeRun = run
	return nil
}

func (m *chamberRepoMock) AbandonActiveRun(playerID int64, abandonedAt time.Time) (*domain.AncestralKiChamberRun, error) {
	m.abandonCalls++
	if m.abandonErr != nil {
		return nil, m.abandonErr
	}

	run, err := m.GetActiveRunByPlayer(playerID)
	if err != nil {
		return nil, err
	}

	run.Status = domain.AncestralKiChamberRunStatusAbandoned
	run.FinishedAt = &abandonedAt
	run.UpdatedAt = abandonedAt
	return run, nil
}

func (m *chamberRepoMock) ResolveBattleResult(playerID int64, monsterID int64, monsterLevel int, won bool, xpReward int, goldReward int, totalBattles int, successPct float64, resolvedAt time.Time) (*domain.AncestralKiChamberResolution, error) {
	m.resolveCalls++
	if m.resolveErr != nil {
		return nil, m.resolveErr
	}

	run, err := m.GetActiveRunByPlayer(playerID)
	if err != nil {
		return nil, err
	}

	result := "LOSE"
	if won {
		result = "WIN"
	}
	m.battleLogs = append(m.battleLogs, &domain.AncestralKiChamberBattleLog{
		ChamberRunID: run.ID,
		BattleNumber: run.CurrentBattle,
		MonsterID:    monsterID,
		MonsterLevel: monsterLevel,
		Result:       result,
		CreatedAt:    resolvedAt,
	})

	if !won {
		run.Status = domain.AncestralKiChamberRunStatusFailed
		run.FinishedAt = &resolvedAt
		return &domain.AncestralKiChamberResolution{Run: run}, nil
	}

	run.Victories++
	run.AccumulatedXP += xpReward
	run.AccumulatedGold += goldReward

	resolution := &domain.AncestralKiChamberResolution{Run: run}
	if run.Victories >= totalBattles {
		run.Status = domain.AncestralKiChamberRunStatusCompleted
		run.FinishedAt = &resolvedAt
		if m.progress == nil {
			m.progress = &domain.AncestralKiChamberProgress{ID: 1, PlayerID: playerID}
		}
		if run.ChamberLevel > m.progress.HighestCompletedLevel {
			m.progress.HighestCompletedLevel = run.ChamberLevel
		}
		resolution.Progress = m.progress
		resolution.XPPayout = int(float64(run.AccumulatedXP) * successPct)
		resolution.GoldPayout = int(float64(run.AccumulatedGold) * successPct)
		return resolution, nil
	}

	run.CurrentBattle++
	return resolution, nil
}

func (m *chamberRepoMock) UpdateRun(run *domain.AncestralKiChamberRun) error {
	m.activeRun = run
	return nil
}

func (m *chamberRepoMock) CreateBattleLog(log *domain.AncestralKiChamberBattleLog) error {
	m.battleLogs = append(m.battleLogs, log)
	return nil
}

func TestGetChamberMonsterLevelRange_LevelOneClamp(t *testing.T) {
	min, max := GetChamberMonsterLevelRange(1, 1)
	if min != 1 || max != 2 {
		t.Fatalf("expected (1,2), got (%d,%d)", min, max)
	}
}

func TestGetChamberMonsterLevelRange_GenericLevel(t *testing.T) {
	min, max := GetChamberMonsterLevelRange(20, 5)
	if min != 20 || max != 25 {
		t.Fatalf("expected (20,25), got (%d,%d)", min, max)
	}
}

func TestAncestralKiChamberService_StartRun_InsufficientResources(t *testing.T) {
	repo := &chamberRepoMock{startRunErr: domain.ErrAncestralKiChamberInsufficientResource}
	player := &domain.Player{ID: 10, Level: 20, Energy: 10, Gold: 1000}
	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	svc := NewAncestralKiChamberService(repo, playerRepo, nil, nil)

	resp, err := svc.StartRun(99, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Text, "Recursos insuficientes") {
		t.Fatalf("unexpected response: %q", resp.Text)
	}
	if len(playerRepo.updatedPlayers) != 0 {
		t.Fatalf("expected no direct player update, got %d", len(playerRepo.updatedPlayers))
	}
}

func TestAncestralKiChamberService_StartRun_BlockedLevel(t *testing.T) {
	repo := &chamberRepoMock{
		progress: &domain.AncestralKiChamberProgress{
			ID:                    1,
			PlayerID:              10,
			HighestCompletedLevel: 1,
		},
	}
	player := &domain.Player{ID: 10, Level: 20, Energy: 10, Gold: 1000}
	svc := NewAncestralKiChamberService(repo, &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}, nil, nil)

	resp, err := svc.StartRun(99, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Text, "bloqueado") {
		t.Fatalf("unexpected response: %q", resp.Text)
	}
}

func TestAncestralKiChamberService_StartRun_UsesAtomicRepo(t *testing.T) {
	repo := &chamberRepoMock{}
	player := &domain.Player{ID: 10, Level: 20, Energy: 10, Gold: 1000}
	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	svc := NewAncestralKiChamberService(repo, playerRepo, nil, nil)

	resp, err := svc.StartRun(99, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Text, "Voce entrou na Camara do Ki Ancestral") {
		t.Fatalf("unexpected response: %q", resp.Text)
	}
	if repo.startRunCalls != 1 {
		t.Fatalf("expected one atomic start call, got %d", repo.startRunCalls)
	}
	if repo.activeRun == nil || repo.activeRun.Status != domain.AncestralKiChamberRunStatusActive {
		t.Fatalf("expected active run to be created, got %+v", repo.activeRun)
	}
	if len(playerRepo.updatedPlayers) != 0 {
		t.Fatalf("expected no direct player update, got %d", len(playerRepo.updatedPlayers))
	}
}

func TestAncestralKiChamberService_StartRun_ReturnsActiveRunMessageFromAtomicExecution(t *testing.T) {
	repo := &chamberRepoMock{startRunErr: domain.ErrAncestralKiChamberRunActive}
	player := &domain.Player{ID: 10, Level: 20, Energy: 10, Gold: 1000}
	svc := NewAncestralKiChamberService(repo, &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}, nil, nil)

	resp, err := svc.StartRun(99, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Text, "ja possui uma run ativa") {
		t.Fatalf("unexpected response: %q", resp.Text)
	}
}

func TestAncestralKiChamberService_ResolveBattleResult_CompleteRun(t *testing.T) {
	repo := &chamberRepoMock{
		progress: &domain.AncestralKiChamberProgress{
			ID:                    1,
			PlayerID:              10,
			HighestCompletedLevel: 0,
		},
		activeRun: &domain.AncestralKiChamberRun{
			ID:            50,
			PlayerID:      10,
			ChamberLevel:  1,
			CurrentBattle: 5,
			Victories:     4,
			Status:        domain.AncestralKiChamberRunStatusActive,
			StartedAt:     time.Now().UTC(),
		},
	}

	player := &domain.Player{ID: 10, Level: 20, Energy: 10, Gold: 1000}
	initialXP := player.XP
	initialGold := player.Gold
	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	svc := NewAncestralKiChamberService(repo, playerRepo, nil, nil)

	resp, err := svc.ResolveBattleResult(99, 1, 20, true, 100, 200)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Text, "Camara concluida") {
		t.Fatalf("unexpected response: %q", resp.Text)
	}
	if repo.progress.HighestCompletedLevel != 1 {
		t.Fatalf("expected highest completed level = 1, got %d", repo.progress.HighestCompletedLevel)
	}
	if repo.activeRun.Status != domain.AncestralKiChamberRunStatusCompleted {
		t.Fatalf("expected run status completed, got %s", repo.activeRun.Status)
	}
	if player.XP != initialXP {
		t.Fatalf("expected no direct player xp mutation, got %d", player.XP-initialXP)
	}
	if player.Gold != initialGold {
		t.Fatalf("expected no direct player gold mutation, got %d", player.Gold-initialGold)
	}
	if repo.resolveCalls != 1 {
		t.Fatalf("expected one atomic resolve call, got %d", repo.resolveCalls)
	}
	if len(playerRepo.updatedPlayers) != 0 {
		t.Fatalf("expected no direct player update, got %d", len(playerRepo.updatedPlayers))
	}
}

func TestAncestralKiChamberService_ResolveBattleResult_FailedRunUsesAtomicRepo(t *testing.T) {
	repo := &chamberRepoMock{
		activeRun: &domain.AncestralKiChamberRun{
			ID:            60,
			PlayerID:      10,
			ChamberLevel:  2,
			CurrentBattle: 3,
			Victories:     2,
			Status:        domain.AncestralKiChamberRunStatusActive,
			StartedAt:     time.Now().UTC(),
		},
	}
	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) {
			return &domain.Player{ID: 10, Level: 20}, nil
		},
	}
	svc := NewAncestralKiChamberService(repo, playerRepo, nil, nil)

	resp, err := svc.ResolveBattleResult(99, 1, 20, false, 100, 200)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Text, "foi derrotado") {
		t.Fatalf("unexpected response: %q", resp.Text)
	}
	if repo.resolveCalls != 1 {
		t.Fatalf("expected one atomic resolve call, got %d", repo.resolveCalls)
	}
	if repo.activeRun.Status != domain.AncestralKiChamberRunStatusFailed {
		t.Fatalf("expected failed run status, got %s", repo.activeRun.Status)
	}
	if len(playerRepo.updatedPlayers) != 0 {
		t.Fatalf("expected no direct player update, got %d", len(playerRepo.updatedPlayers))
	}
}

func TestAncestralKiChamberService_ContinueRun_ProcessesBattle(t *testing.T) {
	repo := &chamberRepoMock{
		activeRun: &domain.AncestralKiChamberRun{
			ID:            51,
			PlayerID:      10,
			ChamberLevel:  1,
			CurrentBattle: 1,
			Victories:     0,
			Status:        domain.AncestralKiChamberRunStatusActive,
			StartedAt:     time.Now().UTC(),
		},
		progress: &domain.AncestralKiChamberProgress{
			ID:                    1,
			PlayerID:              10,
			HighestCompletedLevel: 0,
		},
	}
	player := &domain.Player{ID: 10, Level: 20, RegionID: 1, HP: 120, MaxHP: 120, ATK: 300, DEF: 80, Speed: 20}
	svc := NewAncestralKiChamberService(
		repo,
		&fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
		&fakeMonsterRepo{
			listByRegionFn: func(regionID int64) ([]domain.Monster, error) {
				return []domain.Monster{
					{MonsterID: 1001, Name: "Dummy", Level: 20, HP: 50, ATK: 10, DEF: 1, RegionID: 1},
				}, nil
			},
		},
		nil,
	)

	resp, err := svc.ContinueRun(99)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Text, "Camara - Batalha") {
		t.Fatalf("expected chamber battle message, got %q", resp.Text)
	}
	if len(repo.battleLogs) != 1 {
		t.Fatalf("expected one battle log, got %d", len(repo.battleLogs))
	}
}

func TestAncestralKiChamberService_ContinueRun_SkipsAurumKarun(t *testing.T) {
	repo := &chamberRepoMock{
		activeRun: &domain.AncestralKiChamberRun{
			ID:            52,
			PlayerID:      10,
			ChamberLevel:  1,
			CurrentBattle: 1,
			Victories:     0,
			Status:        domain.AncestralKiChamberRunStatusActive,
			StartedAt:     time.Now().UTC(),
		},
		progress: &domain.AncestralKiChamberProgress{
			ID:                    1,
			PlayerID:              10,
			HighestCompletedLevel: 0,
		},
	}
	player := &domain.Player{ID: 10, Level: 20, RegionID: 1, HP: 120, MaxHP: 120, ATK: 300, DEF: 80, Speed: 20}
	svc := NewAncestralKiChamberService(
		repo,
		&fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		},
		&fakeMonsterRepo{
			listByRegionFn: func(regionID int64) ([]domain.Monster, error) {
				return []domain.Monster{
					{MonsterID: 999, Name: "Aurum Karun", Tag: "AURUM_KARUN", Level: 20, HP: 80, ATK: 20, DEF: 5, RegionID: 1},
					{MonsterID: 1002, Name: "Bandit", Tag: "BANDIT", Level: 20, HP: 60, ATK: 10, DEF: 1, RegionID: 1},
				}, nil
			},
		},
		nil,
	)

	resp, err := svc.ContinueRun(99)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(strings.ToUpper(resp.Text), "AURUM") {
		t.Fatalf("expected Aurum Karun to be skipped, got %q", resp.Text)
	}
}

func TestAncestralKiChamberService_PickMonsterForRun_RespectsLevelRange(t *testing.T) {
	player := &domain.Player{ID: 10, Level: 10, RegionID: 1}
	run := &domain.AncestralKiChamberRun{
		ID:            90,
		PlayerID:      10,
		ChamberLevel:  1,
		CurrentBattle: 1,
		Victories:     0,
		Status:        domain.AncestralKiChamberRunStatusActive,
	}

	svc := NewAncestralKiChamberService(
		&chamberRepoMock{},
		&fakePlayerRepo{getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil }},
		&fakeMonsterRepo{
			listByRegionFn: func(regionID int64) ([]domain.Monster, error) {
				return []domain.Monster{
					{MonsterID: 1, Name: "SALAMORTH", Level: 35, HP: 100, ATK: 20, DEF: 10, RegionID: 1},
					{MonsterID: 2, Name: "LOBO", Level: 10, HP: 50, ATK: 10, DEF: 5, RegionID: 1},
				}, nil
			},
		},
		nil,
	)

	monster, _, err := svc.pickMonsterForRun(player, run)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if monster.Name == "SALAMORTH" {
		t.Fatalf("expected out-of-range monster to be ignored")
	}
}

func TestAncestralKiChamberService_AbandonRun_NoActiveRun(t *testing.T) {
	svc := NewAncestralKiChamberService(&chamberRepoMock{}, &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) {
			return &domain.Player{ID: 10}, nil
		},
	}, nil, nil)

	resp, err := svc.AbandonRun(99)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Text, "não possui run ativa") && !strings.Contains(resp.Text, "nao possui run ativa") {
		t.Fatalf("unexpected response: %q", resp.Text)
	}
}

func TestAncestralKiChamberService_AbandonRun_UsesAtomicRepo(t *testing.T) {
	repo := &chamberRepoMock{
		activeRun: &domain.AncestralKiChamberRun{
			ID:            77,
			PlayerID:      10,
			ChamberLevel:  2,
			CurrentBattle: 3,
			Victories:     2,
			Status:        domain.AncestralKiChamberRunStatusActive,
			StartedAt:     time.Now().UTC(),
		},
	}
	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) {
			return &domain.Player{ID: 10}, nil
		},
	}
	svc := NewAncestralKiChamberService(repo, playerRepo, nil, nil)

	resp, err := svc.AbandonRun(99)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.abandonCalls != 1 {
		t.Fatalf("expected one atomic abandon call, got %d", repo.abandonCalls)
	}
	if repo.activeRun.Status != domain.AncestralKiChamberRunStatusAbandoned {
		t.Fatalf("expected abandoned status, got %s", repo.activeRun.Status)
	}
	if repo.activeRun.FinishedAt == nil {
		t.Fatalf("expected finished_at to be set")
	}
	if len(playerRepo.updatedPlayers) != 0 {
		t.Fatalf("expected no direct player update, got %d", len(playerRepo.updatedPlayers))
	}
	if !strings.Contains(resp.Text, "saiu da Camara do Ki Ancestral") {
		t.Fatalf("unexpected response: %q", resp.Text)
	}
}

func TestAncestralKiChamberService_Open_PlayerError(t *testing.T) {
	svc := NewAncestralKiChamberService(&chamberRepoMock{}, &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) {
			return nil, errors.New("db fail")
		},
	}, nil, nil)

	_, err := svc.Open(99)
	if err == nil {
		t.Fatalf("expected error")
	}
}
