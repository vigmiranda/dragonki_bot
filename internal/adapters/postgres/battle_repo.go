package postgres

import (
	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
)

type BattlePostgresRepo struct {
	db *gorm.DB
}

func NewBattlePostgresRepo(db *gorm.DB) *BattlePostgresRepo {
	return &BattlePostgresRepo{
		db: db,
	}
}

func (r *BattlePostgresRepo) SaveBattleLog(battleLog *domain.BattleLog) error {
	return r.db.Save(battleLog).Error
}

func (r *BattlePostgresRepo) SaveBattleLogAndPlayer(battleLog *domain.BattleLog, player *domain.Player) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
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

		return tx.
			Model(&domain.Player{}).
			Where("id = ?", player.ID).
			Updates(updates).
			Error
	})
}
