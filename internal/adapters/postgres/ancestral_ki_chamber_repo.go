package postgres

import (
	"errors"
	"strings"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *AncestralKiChamberPostgresRepo) StartRunIfEligible(playerID int64, chamberLevel int, energyCost int, goldCost int, startedAt time.Time) (*domain.AncestralKiChamberRun, error) {
	var run domain.AncestralKiChamberRun

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var player domain.Player
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", playerID).
			Take(&player).Error; err != nil {
			return err
		}

		var activeRun domain.AncestralKiChamberRun
		err := tx.
			Model(&domain.AncestralKiChamberRun{}).
			Where("player_id = ?", playerID).
			Where("status = ?", domain.AncestralKiChamberRunStatusActive).
			Order("id DESC").
			First(&activeRun).Error
		if err == nil {
			return domain.ErrAncestralKiChamberRunActive
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if player.Energy < energyCost || player.Gold < goldCost {
			return domain.ErrAncestralKiChamberInsufficientResource
		}

		now := startedAt.UTC()
		if err := tx.
			Model(&domain.Player{}).
			Where("id = ?", playerID).
			Updates(map[string]any{
				"energy":      gorm.Expr("energy - ?", energyCost),
				"gold":        gorm.Expr("gold - ?", goldCost),
				"last_update": now,
			}).Error; err != nil {
			return err
		}

		run = domain.AncestralKiChamberRun{
			PlayerID:      playerID,
			ChamberLevel:  chamberLevel,
			CurrentBattle: 1,
			Victories:     0,
			Status:        domain.AncestralKiChamberRunStatusActive,
			StartedAt:     now,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if err := tx.Create(&run).Error; err != nil {
			if isAncestralKiChamberActiveRunConstraint(err) {
				return domain.ErrAncestralKiChamberRunActive
			}
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &run, nil
}

func (r *AncestralKiChamberPostgresRepo) ResolveBattleResult(playerID int64, monsterID int64, monsterLevel int, won bool, xpReward int, goldReward int, totalBattles int, successPct float64, resolvedAt time.Time) (*domain.AncestralKiChamberResolution, error) {
	var resolution domain.AncestralKiChamberResolution

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var run domain.AncestralKiChamberRun
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Model(&domain.AncestralKiChamberRun{}).
			Where("player_id = ?", playerID).
			Where("status = ?", domain.AncestralKiChamberRunStatusActive).
			Order("id DESC").
			First(&run).Error; err != nil {
			return err
		}

		result := "LOSE"
		if won {
			result = "WIN"
		}
		if err := tx.Create(&domain.AncestralKiChamberBattleLog{
			ChamberRunID: run.ID,
			BattleNumber: run.CurrentBattle,
			MonsterID:    monsterID,
			MonsterLevel: monsterLevel,
			Result:       result,
			CreatedAt:    resolvedAt.UTC(),
		}).Error; err != nil {
			return err
		}

		now := resolvedAt.UTC()
		if !won {
			run.Status = domain.AncestralKiChamberRunStatusFailed
			run.FinishedAt = &now
			run.UpdatedAt = now
			if err := tx.
				Model(&domain.AncestralKiChamberRun{}).
				Where("id = ?", run.ID).
				Updates(map[string]any{
					"status":      run.Status,
					"finished_at": run.FinishedAt,
					"updated_at":  run.UpdatedAt,
				}).Error; err != nil {
				return err
			}
			resolution.Run = &run
			return nil
		}

		run.Victories++
		if xpReward > 0 {
			run.AccumulatedXP += xpReward
		}
		if goldReward > 0 {
			run.AccumulatedGold += goldReward
		}

		updateFields := map[string]any{
			"victories":        run.Victories,
			"accumulated_xp":   run.AccumulatedXP,
			"accumulated_gold": run.AccumulatedGold,
			"updated_at":       now,
		}

		if run.Victories >= totalBattles {
			run.Status = domain.AncestralKiChamberRunStatusCompleted
			run.FinishedAt = &now
			updateFields["status"] = run.Status
			updateFields["finished_at"] = run.FinishedAt
		} else {
			run.CurrentBattle++
			updateFields["current_battle"] = run.CurrentBattle
		}

		run.UpdatedAt = now
		if err := tx.
			Model(&domain.AncestralKiChamberRun{}).
			Where("id = ?", run.ID).
			Updates(updateFields).Error; err != nil {
			return err
		}

		if run.Status == domain.AncestralKiChamberRunStatusCompleted {
			progress, err := getOrCreateAncestralKiChamberProgressWithTx(tx, playerID, now)
			if err != nil {
				return err
			}
			if run.ChamberLevel > progress.HighestCompletedLevel {
				progress.HighestCompletedLevel = run.ChamberLevel
				progress.UpdatedAt = now
				if err := tx.
					Model(&domain.AncestralKiChamberProgress{}).
					Where("id = ?", progress.ID).
					Updates(map[string]any{
						"highest_completed_level": progress.HighestCompletedLevel,
						"updated_at":              progress.UpdatedAt,
					}).Error; err != nil {
					return err
				}
			}

			xpPayout := int(float64(run.AccumulatedXP) * successPct)
			goldPayout := int(float64(run.AccumulatedGold) * successPct)
			if xpPayout > 0 || goldPayout > 0 {
				if err := tx.
					Model(&domain.Player{}).
					Where("id = ?", playerID).
					Updates(map[string]any{
						"xp":          gorm.Expr("xp + ?", xpPayout),
						"gold":        gorm.Expr("gold + ?", goldPayout),
						"last_update": now,
					}).Error; err != nil {
					return err
				}
			}

			resolution.Progress = progress
			resolution.XPPayout = xpPayout
			resolution.GoldPayout = goldPayout
		}

		resolution.Run = &run
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &resolution, nil
}

func (r *AncestralKiChamberPostgresRepo) AbandonActiveRun(playerID int64, abandonedAt time.Time) (*domain.AncestralKiChamberRun, error) {
	var run domain.AncestralKiChamberRun

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Model(&domain.AncestralKiChamberRun{}).
			Where("player_id = ?", playerID).
			Where("status = ?", domain.AncestralKiChamberRunStatusActive).
			Order("id DESC").
			First(&run).Error; err != nil {
			return err
		}

		now := abandonedAt.UTC()
		run.Status = domain.AncestralKiChamberRunStatusAbandoned
		run.FinishedAt = &now
		run.UpdatedAt = now

		return tx.
			Model(&domain.AncestralKiChamberRun{}).
			Where("id = ?", run.ID).
			Updates(map[string]any{
				"status":      run.Status,
				"finished_at": run.FinishedAt,
				"updated_at":  run.UpdatedAt,
			}).Error
	})
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

func isAncestralKiChamberActiveRunConstraint(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), "uq_ak_chamber_runs_player_active")
}

func getOrCreateAncestralKiChamberProgressWithTx(tx *gorm.DB, playerID int64, now time.Time) (*domain.AncestralKiChamberProgress, error) {
	var progress domain.AncestralKiChamberProgress
	err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
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
		CreatedAt:             now,
		UpdatedAt:             now,
	}
	if err := tx.Create(&progress).Error; err != nil {
		return nil, err
	}
	return &progress, nil
}
