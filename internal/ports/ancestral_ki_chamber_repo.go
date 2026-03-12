package ports

import (
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
)

type AncestralKiChamberRepository interface {
	GetOrCreateProgress(playerID int64) (*domain.AncestralKiChamberProgress, error)
	UpdateProgress(progress *domain.AncestralKiChamberProgress) error

	GetActiveRunByPlayer(playerID int64) (*domain.AncestralKiChamberRun, error)
	StartRunIfEligible(playerID int64, chamberLevel int, energyCost int, goldCost int, startedAt time.Time) (*domain.AncestralKiChamberRun, error)
	ResolveBattleResult(playerID int64, monsterID int64, monsterLevel int, won bool, xpReward int, goldReward int, totalBattles int, successPct float64, resolvedAt time.Time) (*domain.AncestralKiChamberResolution, error)
	AbandonActiveRun(playerID int64, abandonedAt time.Time) (*domain.AncestralKiChamberRun, error)
	CreateRun(run *domain.AncestralKiChamberRun) error
	UpdateRun(run *domain.AncestralKiChamberRun) error

	CreateBattleLog(log *domain.AncestralKiChamberBattleLog) error
}
