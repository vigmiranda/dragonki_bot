package ports

import "github.com/vigmiranda/dragonki/internal/domain"

type AncestralKiChamberRepository interface {
	GetOrCreateProgress(playerID int64) (*domain.AncestralKiChamberProgress, error)
	UpdateProgress(progress *domain.AncestralKiChamberProgress) error

	GetActiveRunByPlayer(playerID int64) (*domain.AncestralKiChamberRun, error)
	CreateRun(run *domain.AncestralKiChamberRun) error
	UpdateRun(run *domain.AncestralKiChamberRun) error

	CreateBattleLog(log *domain.AncestralKiChamberBattleLog) error
}
