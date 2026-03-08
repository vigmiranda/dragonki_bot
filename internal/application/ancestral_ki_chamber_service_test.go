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
	progress   *domain.AncestralKiChamberProgress
	activeRun  *domain.AncestralKiChamberRun
	battleLogs []*domain.AncestralKiChamberBattleLog
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

func (m *chamberRepoMock) CreateRun(run *domain.AncestralKiChamberRun) error {
	run.ID = 1
	m.activeRun = run
	return nil
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
	repo := &chamberRepoMock{}
	player := &domain.Player{ID: 10, Level: 20, Energy: 0, Gold: 100}
	svc := NewAncestralKiChamberService(repo, &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}, nil, nil)

	resp, err := svc.StartRun(99, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Text, "Recursos insuficientes") {
		t.Fatalf("unexpected response: %q", resp.Text)
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
	svc := NewAncestralKiChamberService(repo, &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}, nil, nil)

	resp, err := svc.ResolveBattleResult(99, 1, 20, true)
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
