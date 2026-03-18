package postgres

import (
	"context"
	"fmt"

	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BattlePostgresRepo struct {
	db *gorm.DB
}

func NewBattlePostgresRepo(db *gorm.DB) *BattlePostgresRepo {
	return &BattlePostgresRepo{
		db: db,
	}
}

func (r *BattlePostgresRepo) SaveBattleLog(ctx context.Context, battleLog *domain.BattleLog) error {
	return r.db.WithContext(ctx).Save(battleLog).Error
}

func (r *BattlePostgresRepo) SaveBattleLogAndPlayer(ctx context.Context, battleLog *domain.BattleLog, player *domain.Player) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existingPlayer domain.Player
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", player.ID).
			Take(&existingPlayer).Error; err != nil {
			return err
		}

		if err := tx.Save(battleLog).Error; err != nil {
			return err
		}

		updates := map[string]interface{}{
			"name":              player.Name,
			"gender":            player.Gender,
			"class_id":          player.ClassID,
			"level":             player.Level,
			"xp":                player.XP,
			"hp":                player.HP,
			"max_hp":            player.MaxHP,
			"atk":               player.ATK,
			"def":               player.DEF,
			"speed":             player.Speed,
			"energy":            player.Energy,
			"max_energy":        player.MaxEnergy,
			"gold":              player.Gold,
			"ki_crystals":       player.KiCrystals,
			"region_id":         player.RegionID,
			"last_energy_at":    player.LastEnergyAt,
			"form":              player.Form,
			"form_battles_left": player.FormBattlesLeft,
			"onboarding_step":   player.OnboardingStep,
			"last_update":       battleLog.CreatedAt.UTC(),
			"user_id_telegram":  player.UserIDTelegram,
			"name_telegram":     player.NameTelegram,
			"user_name":         player.UserName,
			"dodge_accumulator": player.DodgeAccumulator,
		}

		result := tx.
			Model(&domain.Player{}).
			Where("id = ?", player.ID).
			Updates(updates)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func (r *BattlePostgresRepo) GetWeeklyLeaderboard(metric string, limit int) ([]domain.WeeklyBattleRank, error) {
	orderField := "total_battles"
	switch metric {
	case "total_battles":
		orderField = "total_battles"
	case "wins":
		orderField = "wins"
	case "losses":
		orderField = "losses"
	case "escape":
		orderField = "escape"
	default:
		return nil, fmt.Errorf("invalid ranking metric: %s", metric)
	}

	if limit <= 0 {
		limit = 5
	}

	query := fmt.Sprintf(`
SELECT
    p.id as player_id,
    p.name,
    COUNT(*) as total_battles,
    SUM(CASE WHEN result = 'WIN' THEN 1 ELSE 0 END) as wins,
    SUM(CASE WHEN result in ('LOSE','ESCAPE/LOSS') THEN 1 ELSE 0 END) as losses,
    SUM(CASE WHEN result = 'ESCAPE' THEN 1 ELSE 0 END) as escape
FROM battle_logs bl
JOIN players p on p.id = bl.player_id
WHERE bl.created_at >= date_trunc('week', NOW())
GROUP BY p.id, p.name
ORDER BY %s DESC, total_battles DESC, p.name ASC
LIMIT ?`, orderField)

	var result []domain.WeeklyBattleRank
	if err := r.db.Raw(query, limit).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
