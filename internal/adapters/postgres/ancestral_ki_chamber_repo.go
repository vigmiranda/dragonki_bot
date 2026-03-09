package postgres

import (
	"errors"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
)

type AncestralKiChamberPostgresRepo struct {
	db *gorm.DB
}

func NewAncestralKiChamberPostgresRepo(db *gorm.DB) *AncestralKiChamberPostgresRepo {
	return &AncestralKiChamberPostgresRepo{db: db}
}

func (r *AncestralKiChamberPostgresRepo) GetOrCreateProgress(playerID int64) (*domain.AncestralKiChamberProgress, error) {
	var progress domain.AncestralKiChamberProgress
	err := r.db.
		Model(&domain.AncestralKiChamberProgress{}).
		Where("player_id = ?", playerID).
		First(&progress).Error
	if err == nil {
		return &progress, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	progress = domain.AncestralKiChamberProgress{
		PlayerID:              playerID,
		HighestCompletedLevel: 0,
		CreatedAt:             time.Now().UTC(),
		UpdatedAt:             time.Now().UTC(),
	}
	if err := r.db.Create(&progress).Error; err != nil {
		return nil, err
	}
	return &progress, nil
}

func (r *AncestralKiChamberPostgresRepo) UpdateProgress(progress *domain.AncestralKiChamberProgress) error {
	if progress == nil {
		return gorm.ErrInvalidData
	}
	progress.UpdatedAt = time.Now().UTC()
	return r.db.
		Model(&domain.AncestralKiChamberProgress{}).
		Where("id = ?", progress.ID).
		Updates(map[string]any{
			"highest_completed_level": progress.HighestCompletedLevel,
			"updated_at":              progress.UpdatedAt,
		}).Error
}

func (r *AncestralKiChamberPostgresRepo) GetActiveRunByPlayer(playerID int64) (*domain.AncestralKiChamberRun, error) {
	var run domain.AncestralKiChamberRun
	err := r.db.
		Model(&domain.AncestralKiChamberRun{}).
		Where("player_id = ?", playerID).
		Where("status = ?", domain.AncestralKiChamberRunStatusActive).
		Order("id DESC").
		First(&run).Error
	if err != nil {
		return nil, err
	}
	return &run, nil
}

func (r *AncestralKiChamberPostgresRepo) CreateRun(run *domain.AncestralKiChamberRun) error {
	if run == nil {
		return gorm.ErrInvalidData
	}
	return r.db.Create(run).Error
}

func (r *AncestralKiChamberPostgresRepo) UpdateRun(run *domain.AncestralKiChamberRun) error {
	if run == nil {
		return gorm.ErrInvalidData
	}
	run.UpdatedAt = time.Now().UTC()
	return r.db.
		Model(&domain.AncestralKiChamberRun{}).
		Where("id = ?", run.ID).
		Updates(map[string]any{
			"current_battle":   run.CurrentBattle,
			"victories":        run.Victories,
			"accumulated_xp":   run.AccumulatedXP,
			"accumulated_gold": run.AccumulatedGold,
			"status":           run.Status,
			"finished_at":      run.FinishedAt,
			"updated_at":       run.UpdatedAt,
		}).Error
}

func (r *AncestralKiChamberPostgresRepo) CreateBattleLog(logRow *domain.AncestralKiChamberBattleLog) error {
	if logRow == nil {
		return gorm.ErrInvalidData
	}
	return r.db.Create(logRow).Error
}
