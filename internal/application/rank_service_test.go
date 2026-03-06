package application

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/vigmiranda/dragonki/internal/domain"
)

type fakeRankBattleRepo struct {
	rows []domain.WeeklyBattleRank
	err  error
}

func (f *fakeRankBattleRepo) SaveBattleLog(ctx context.Context, battleLog *domain.BattleLog) error {
	return nil
}

func (f *fakeRankBattleRepo) SaveBattleLogAndPlayer(ctx context.Context, battleLog *domain.BattleLog, player *domain.Player) error {
	return nil
}

func (f *fakeRankBattleRepo) GetWeeklyLeaderboard(metric string, limit int) ([]domain.WeeklyBattleRank, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.rows, nil
}

func TestRankService_OpenWeeklyRank(t *testing.T) {
	svc := NewRankService(&fakeRankBattleRepo{})
	resp, err := svc.OpenWeeklyRank()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "RANK_MENU" {
		t.Fatalf("expected RANK_MENU, got %q", resp.MenuKey)
	}
}

func TestRankService_WeeklyTop_Success(t *testing.T) {
	svc := NewRankService(&fakeRankBattleRepo{
		rows: []domain.WeeklyBattleRank{
			{Name: "Goku", TotalBattles: 20, Wins: 14, Losses: 5, Escape: 1},
			{Name: "Vegeta", TotalBattles: 18, Wins: 12, Losses: 4, Escape: 2},
		},
	})

	resp, err := svc.WeeklyTop(rankMetricWins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != "RANK_MENU" {
		t.Fatalf("expected RANK_MENU, got %q", resp.MenuKey)
	}
	if !strings.Contains(resp.Text, "Top 5 - Mais vitorias na semana") {
		t.Fatalf("unexpected title: %q", resp.Text)
	}
	if !strings.Contains(resp.Text, "1. Goku - Vitorias: 14") {
		t.Fatalf("unexpected ranking text: %q", resp.Text)
	}
}

func TestRankService_WeeklyTop_InvalidMetric(t *testing.T) {
	svc := NewRankService(&fakeRankBattleRepo{})
	_, err := svc.WeeklyTop("unknown")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestRankService_WeeklyTop_RepoError(t *testing.T) {
	svc := NewRankService(&fakeRankBattleRepo{err: errors.New("db fail")})
	_, err := svc.WeeklyTop(rankMetricBattles)
	if err == nil {
		t.Fatalf("expected error")
	}
}
