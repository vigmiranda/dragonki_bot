package application

import (
	"fmt"
	"strings"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

const (
	rankMetricBattles = "total_battles"
	rankMetricWins    = "wins"
	rankMetricLosses  = "losses"
	rankMetricEscape  = "escape"
)

type RankService struct {
	BattleRepository ports.BattleRepository
}

func NewRankService(battleRepo ports.BattleRepository) *RankService {
	return &RankService{BattleRepository: battleRepo}
}

func (s *RankService) OpenWeeklyRank() (view.GameResponse, error) {
	return view.GameResponse{
		Text:    "Ranking semanal de batalhas.\nEscolha uma categoria:",
		Image:   &view.Image{Type: view.ImageAsset, Path: "assets/npcs/mestre_kalunga.png"},
		MenuKey: "RANK_MENU",
	}, nil
}

func (s *RankService) WeeklyTop(metric string) (view.GameResponse, error) {
	title, valueLabel, err := rankMetricMeta(metric)
	if err != nil {
		return view.GameResponse{}, err
	}

	rows, err := s.BattleRepository.GetWeeklyLeaderboard(metric, 5)
	if err != nil {
		return view.GameResponse{}, err
	}

	text := buildRankText(title, valueLabel, rows, metric)
	return view.GameResponse{
		Text:    text,
		Image:   &view.Image{Type: view.ImageAsset, Path: "assets/npcs/mestre_kalunga.png"},
		MenuKey: "RANK_MENU",
	}, nil
}

func rankMetricMeta(metric string) (title string, valueLabel string, err error) {
	switch metric {
	case rankMetricBattles:
		return "Top 5 - Mais batalhas na semana", "Batalhas", nil
	case rankMetricWins:
		return "Top 5 - Mais vitorias na semana", "Vitorias", nil
	case rankMetricLosses:
		return "Top 5 - Mais derrotas na semana", "Derrotas", nil
	case rankMetricEscape:
		return "Top 5 - Mais fugas na semana", "Escapes", nil
	default:
		return "", "", fmt.Errorf("rank metric invalida: %s", metric)
	}
}

func buildRankText(title string, valueLabel string, rows []domain.WeeklyBattleRank, metric string) string {
	var b strings.Builder
	b.WriteString(title)
	b.WriteString("\n\n")

	if len(rows) == 0 {
		b.WriteString("Nenhuma batalha registrada nesta semana.")
		return b.String()
	}

	for i, row := range rows {
		fmt.Fprintf(&b, "%d. %s - %s: %d", i+1, row.Name, valueLabel, rankMetricValue(row, metric))
		if i < len(rows)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

func rankMetricValue(row domain.WeeklyBattleRank, metric string) int64 {
	switch metric {
	case rankMetricWins:
		return row.Wins
	case rankMetricLosses:
		return row.Losses
	case rankMetricEscape:
		return row.Escape
	default:
		return row.TotalBattles
	}
}
