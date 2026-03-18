package application

import (
	"strings"
	"testing"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
)

type fakeVIPConfigRepo struct {
	cfg *domain.VIPConfig
	err error
}

func (f *fakeVIPConfigRepo) GetEnabledConfig() (*domain.VIPConfig, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.cfg, nil
}

type fakePlayerVIPRepo struct {
	byPlayerID map[int64]*domain.PlayerVIP
	err        error
	upserted   *domain.PlayerVIP
}

func (f *fakePlayerVIPRepo) GetByPlayerID(playerID int64) (*domain.PlayerVIP, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.byPlayerID == nil {
		return nil, nil
	}
	return f.byPlayerID[playerID], nil
}

func (f *fakePlayerVIPRepo) Upsert(vip *domain.PlayerVIP) error {
	f.upserted = vip
	return f.err
}

func TestVIPService_DecoratePlayerName_ActiveVIP(t *testing.T) {
	svc := NewVIPService(nil, &fakeVIPConfigRepo{}, &fakePlayerVIPRepo{
		byPlayerID: map[int64]*domain.PlayerVIP{
			1: {PlayerID: 1, ExpiresAt: time.Now().UTC().Add(time.Hour)},
		},
	})

	if got := svc.DecoratePlayerName(1, "Goku"); got != "👑 Goku" {
		t.Fatalf("expected crowned name, got %q", got)
	}
}

func TestVIPService_DecoratePlayerName_DoesNotDuplicatePrefix(t *testing.T) {
	svc := NewVIPService(nil, &fakeVIPConfigRepo{}, &fakePlayerVIPRepo{
		byPlayerID: map[int64]*domain.PlayerVIP{
			1: {PlayerID: 1, ExpiresAt: time.Now().UTC().Add(time.Hour)},
		},
	})

	if got := svc.DecoratePlayerName(1, "👑 Goku"); got != "👑 Goku" {
		t.Fatalf("expected existing prefix preserved, got %q", got)
	}
}

func TestRankService_WeeklyTop_AddsVIPPrefixToRankingName(t *testing.T) {
	svc := NewRankService(&fakeRankBattleRepo{
		rows: []domain.WeeklyBattleRank{
			{PlayerID: 1, Name: "Goku", TotalBattles: 20, Wins: 14},
			{PlayerID: 2, Name: "Vegeta", TotalBattles: 18, Wins: 12},
		},
	}, NewVIPService(nil, &fakeVIPConfigRepo{}, &fakePlayerVIPRepo{
		byPlayerID: map[int64]*domain.PlayerVIP{
			1: {PlayerID: 1, ExpiresAt: time.Now().UTC().Add(time.Hour)},
		},
	}))

	resp, err := svc.WeeklyTop(rankMetricWins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Text, "1. 👑 Goku - Vitorias: 14") {
		t.Fatalf("expected VIP prefix in ranking, got %q", resp.Text)
	}
	if !strings.Contains(resp.Text, "2. Vegeta - Vitorias: 12") {
		t.Fatalf("expected non-VIP player without prefix, got %q", resp.Text)
	}
}

func TestWeeklyBossService_Ranking_AddsVIPPrefixToPlayerName(t *testing.T) {
	svc := NewWeeklyBossService(&weeklyBossRepoMock{
		activeBoss: &domain.WeeklyBoss{ID: 1, Name: "Surtr", Tag: "SURTR", ATK: 100, CurrentHP: 1000, MaxHP: 1000},
		rankRows: []domain.WeeklyBossRankingEntry{
			{PlayerID: 1, PlayerName: "Vitor", TotalDamage: 1000},
			{PlayerID: 2, PlayerName: "Kael", TotalDamage: 900},
		},
	}, &weeklyBossPlayerRepoMock{player: &domain.Player{ID: 1}}, nil, nil)
	svc.SetVIPService(NewVIPService(nil, &fakeVIPConfigRepo{}, &fakePlayerVIPRepo{
		byPlayerID: map[int64]*domain.PlayerVIP{
			1: {PlayerID: 1, ExpiresAt: time.Now().UTC().Add(time.Hour)},
		},
	}))

	resp, err := svc.Ranking()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !strings.Contains(resp.Text, "1. 👑 Vitor - 1000") {
		t.Fatalf("expected VIP prefix in big boss ranking, got %q", resp.Text)
	}
	if !strings.Contains(resp.Text, "2. Kael - 900") {
		t.Fatalf("expected non-VIP player without prefix, got %q", resp.Text)
	}
}
