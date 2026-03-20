package postgres

import (
	"errors"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WeeklyBossPostgresRepo struct {
	db *gorm.DB
}

var weeklyBossAttackDateLocation = loadWeeklyBossAttackDateLocation()

func NewWeeklyBossPostgresRepo(db *gorm.DB) *WeeklyBossPostgresRepo {
	return &WeeklyBossPostgresRepo{db: db}
}

func loadWeeklyBossAttackDateLocation() *time.Location {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		return time.FixedZone("BRT", -3*60*60)
	}
	return loc
}

func normalizeWeeklyBossAttackDate(at time.Time) time.Time {
	localDate := at.In(weeklyBossAttackDateLocation)
	return time.Date(localDate.Year(), localDate.Month(), localDate.Day(), 0, 0, 0, 0, time.UTC)
}

func weeklyBossAttackDateValue(at time.Time) string {
	return normalizeWeeklyBossAttackDate(at).Format("2006-01-02")
}

func (r *WeeklyBossPostgresRepo) GetActive(at time.Time) (*domain.WeeklyBoss, error) {
	var boss domain.WeeklyBoss
	err := r.db.
		Model(&domain.WeeklyBoss{}).
		Where("status = ?", domain.WeeklyBossStatusActive).
		Where("starts_at <= ?", at).
		Where("ends_at >= ?", at).
		Order("id DESC").
		First(&boss).Error
	if err != nil {
		return nil, err
	}
	return &boss, nil
}

func (r *WeeklyBossPostgresRepo) GetByID(id uint) (*domain.WeeklyBoss, error) {
	var boss domain.WeeklyBoss
	if err := r.db.First(&boss, id).Error; err != nil {
		return nil, err
	}
	return &boss, nil
}

func (r *WeeklyBossPostgresRepo) GetByWeekKey(weekKey string) (*domain.WeeklyBoss, error) {
	var boss domain.WeeklyBoss
	if err := r.db.
		Model(&domain.WeeklyBoss{}).
		Where("week_key = ?", weekKey).
		First(&boss).Error; err != nil {
		return nil, err
	}
	return &boss, nil
}

func (r *WeeklyBossPostgresRepo) Create(boss *domain.WeeklyBoss) error {
	if boss == nil {
		return gorm.ErrInvalidData
	}
	return r.db.Create(boss).Error
}

func (r *WeeklyBossPostgresRepo) ExpireActiveBefore(at time.Time) error {
	return r.db.Exec(`
		UPDATE weekly_bosses
		   SET status = ?, updated_at = ?
		 WHERE status = ? AND ends_at < ?`,
		domain.WeeklyBossStatusExpired, at, domain.WeeklyBossStatusActive, at,
	).Error
}

func (r *WeeklyBossPostgresRepo) GetRandomTemplate() (*domain.WeeklyBossTemplate, error) {
	var template domain.WeeklyBossTemplate
	err := r.db.
		Model(&domain.WeeklyBossTemplate{}).
		Where("enabled = ?", true).
		Order("RANDOM()").
		Limit(1).
		First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *WeeklyBossPostgresRepo) ApplyDamage(id uint, damage int64, at time.Time) (*domain.WeeklyBoss, error) {
	if damage < 0 {
		damage = 0
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`
			UPDATE weekly_bosses
			   SET current_hp = GREATEST(current_hp - ?, 0),
			       updated_at = ?
			 WHERE id = ?`, damage, at, id).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			UPDATE weekly_bosses
			   SET status = ?, defeated_at = COALESCE(defeated_at, ?), updated_at = ?
			 WHERE id = ? AND current_hp = 0 AND status <> ?`,
			domain.WeeklyBossStatusDefeated, at, at, id, domain.WeeklyBossStatusDefeated,
		).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *WeeklyBossPostgresRepo) CountPlayerAttacksOnDate(bossID uint, playerID int64, attackDate time.Time) (int64, error) {
	var count int64
	err := r.db.
		Model(&domain.WeeklyBossAttack{}).
		Where("boss_id = ?", bossID).
		Where("player_id = ?", playerID).
		Where("attack_date = ?", weeklyBossAttackDateValue(attackDate)).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *WeeklyBossPostgresRepo) CreateAttack(attack *domain.WeeklyBossAttack) error {
	if attack == nil {
		return gorm.ErrInvalidData
	}
	attack.AttackAt = normalizeWeeklyBossAttackDate(attack.AttackAt)
	return r.db.Transaction(func(tx *gorm.DB) error {
		if attack.DailySeq <= 0 {
			var count int64
			if err := tx.
				Model(&domain.WeeklyBossAttack{}).
				Where("boss_id = ?", attack.BossID).
				Where("player_id = ?", attack.PlayerID).
				Where("attack_date = ?", weeklyBossAttackDateValue(attack.AttackAt)).
				Count(&count).Error; err != nil {
				return err
			}
			attack.DailySeq = int(count) + 1
		}
		return tx.Create(attack).Error
	})
}

func (r *WeeklyBossPostgresRepo) ExecuteAttack(bossID uint, playerID int64, damage int64, attackLimit int, energyCost int, goldCost int, attackedAt time.Time) (*domain.WeeklyBoss, int64, error) {
	if bossID == 0 || playerID == 0 {
		return nil, 0, gorm.ErrInvalidData
	}
	if attackLimit < 1 {
		attackLimit = 1
	}
	if damage < 0 {
		damage = 0
	}

	attackDate := normalizeWeeklyBossAttackDate(attackedAt)
	var updatedBoss domain.WeeklyBoss
	var totalDamage int64

	err := r.db.Transaction(func(tx *gorm.DB) error {
		type playerState struct {
			ID uint
		}
		var player playerState
		if err := tx.
			Raw(`SELECT id FROM players WHERE id = ? AND energy >= ? AND gold >= ? FOR UPDATE`, playerID, energyCost, goldCost).
			Scan(&player).Error; err != nil {
			return err
		}
		if player.ID == 0 {
			return domain.ErrWeeklyBossInsufficientResource
		}

		var attacksToday int64
		if err := tx.
			Model(&domain.WeeklyBossAttack{}).
			Where("boss_id = ?", bossID).
			Where("player_id = ?", playerID).
			Where("attack_date = ?", attackDate.Format("2006-01-02")).
			Count(&attacksToday).Error; err != nil {
			return err
		}
		if attacksToday >= int64(attackLimit) {
			return domain.ErrWeeklyBossAttackLimitReached
		}

		attack := &domain.WeeklyBossAttack{
			BossID:   bossID,
			PlayerID: playerID,
			AttackAt: attackDate,
			DailySeq: int(attacksToday) + 1,
			Damage:   damage,
		}
		if err := tx.Create(attack).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			UPDATE weekly_bosses
			   SET current_hp = GREATEST(current_hp - ?, 0),
			       updated_at = ?
			 WHERE id = ?`, damage, attackedAt, bossID).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			UPDATE weekly_bosses
			   SET status = ?, defeated_at = COALESCE(defeated_at, ?), updated_at = ?
			 WHERE id = ? AND current_hp = 0 AND status <> ?`,
			domain.WeeklyBossStatusDefeated, attackedAt, attackedAt, bossID, domain.WeeklyBossStatusDefeated,
		).Error; err != nil {
			return err
		}

		entry := domain.WeeklyBossContribution{
			BossID:       bossID,
			PlayerID:     playerID,
			TotalDamage:  damage,
			LastAttackAt: attackedAt,
		}
		if err := tx.
			Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "boss_id"},
					{Name: "player_id"},
				},
				DoUpdates: clause.Assignments(map[string]any{
					"total_damage": clause.Expr{
						SQL:  "weekly_boss_contributions.total_damage + ?",
						Vars: []any{damage},
					},
					"last_attack_at": attackedAt,
					"updated_at":     attackedAt,
				}),
			}).
			Create(&entry).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			UPDATE players
			   SET energy = energy - ?,
			       gold = gold - ?,
			       last_update = ?
			 WHERE id = ?`, energyCost, goldCost, attackedAt, playerID).Error; err != nil {
			return err
		}

		if err := tx.
			Model(&domain.WeeklyBoss{}).
			Where("id = ?", bossID).
			First(&updatedBoss).Error; err != nil {
			return err
		}

		var contribution domain.WeeklyBossContribution
		if err := tx.
			Model(&domain.WeeklyBossContribution{}).
			Where("boss_id = ?", bossID).
			Where("player_id = ?", playerID).
			First(&contribution).Error; err != nil {
			return err
		}
		totalDamage = contribution.TotalDamage
		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	return &updatedBoss, totalDamage, nil
}

func (r *WeeklyBossPostgresRepo) GetContribution(bossID uint, playerID int64) (*domain.WeeklyBossContribution, error) {
	var contribution domain.WeeklyBossContribution
	err := r.db.
		Model(&domain.WeeklyBossContribution{}).
		Where("boss_id = ?", bossID).
		Where("player_id = ?", playerID).
		First(&contribution).Error
	if err != nil {
		return nil, err
	}
	return &contribution, nil
}

func (r *WeeklyBossPostgresRepo) AddContributionDamage(bossID uint, playerID int64, damage int64, attackedAt time.Time) error {
	if damage < 0 {
		damage = 0
	}

	entry := domain.WeeklyBossContribution{
		BossID:       bossID,
		PlayerID:     playerID,
		TotalDamage:  damage,
		LastAttackAt: attackedAt,
	}

	return r.db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "boss_id"},
				{Name: "player_id"},
			},
			DoUpdates: clause.Assignments(map[string]any{
				"total_damage": clause.Expr{
					SQL:  "weekly_boss_contributions.total_damage + ?",
					Vars: []any{damage},
				},
				"last_attack_at": attackedAt,
				"updated_at":     attackedAt,
			}),
		}).
		Create(&entry).Error
}

func (r *WeeklyBossPostgresRepo) ListTopContributions(bossID uint, limit int) ([]domain.WeeklyBossRankingEntry, error) {
	if limit <= 0 {
		limit = 10
	}

	query := `
SELECT
    c.player_id,
    COALESCE(NULLIF(TRIM(p.name), ''), 'Player ' || c.player_id::text) AS player_name,
    c.total_damage
FROM weekly_boss_contributions c
LEFT JOIN players p ON p.id = c.player_id
WHERE c.boss_id = ?
ORDER BY c.total_damage DESC, c.last_attack_at ASC
LIMIT ?`

	var rows []domain.WeeklyBossRankingEntry
	err := r.db.Raw(query, bossID, limit).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *WeeklyBossPostgresRepo) GetRewardByPosition(position int) (*domain.WeeklyBossReward, error) {
	var reward domain.WeeklyBossReward
	err := r.db.
		Model(&domain.WeeklyBossReward{}).
		Where("position = ?", position).
		Where("enabled = ?", true).
		Preload("Items").
		First(&reward).Error
	if err != nil {
		return nil, err
	}
	return &reward, nil
}

func (r *WeeklyBossPostgresRepo) CreatePayoutIfNotExists(payout *domain.WeeklyBossPayout) (bool, error) {
	if payout == nil {
		return false, gorm.ErrInvalidData
	}
	res := r.db.
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(payout)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (r *WeeklyBossPostgresRepo) CreateAndPayPayoutIfNotExists(payout *domain.WeeklyBossPayout, items []domain.WeeklyBossRewardItem, paidAt time.Time) (bool, error) {
	if payout == nil {
		return false, gorm.ErrInvalidData
	}

	created := false
	err := r.db.Transaction(func(tx *gorm.DB) error {
		res := tx.
			Clauses(clause.OnConflict{DoNothing: true}).
			Create(payout)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return nil
		}

		if payout.GoldPaid > 0 {
			if err := tx.Exec(`
				UPDATE players
				   SET gold = gold + ?,
				       last_update = ?
				 WHERE id = ?`,
				payout.GoldPaid, paidAt, payout.PlayerID,
			).Error; err != nil {
				return err
			}
		}

		for _, itemReward := range items {
			if itemReward.ItemID == 0 || itemReward.Quantity <= 0 {
				continue
			}
			if err := createOrUpdateInventoryWithDB(tx, payout.PlayerID, itemReward.ItemID, itemReward.Quantity); err != nil {
				return err
			}
		}

		if err := tx.
			Model(&domain.WeeklyBossPayout{}).
			Where("id = ?", payout.ID).
			Updates(map[string]any{
				"status":        domain.WeeklyBossPayoutPaid,
				"error_message": "",
				"paid_at":       paidAt,
				"updated_at":    paidAt,
			}).Error; err != nil {
			return err
		}

		created = true
		return nil
	})
	if err != nil {
		return false, err
	}
	return created, nil
}

func (r *WeeklyBossPostgresRepo) MarkPayoutPaid(payoutID uint, paidAt time.Time) error {
	return r.db.
		Model(&domain.WeeklyBossPayout{}).
		Where("id = ?", payoutID).
		Updates(map[string]any{
			"status":        domain.WeeklyBossPayoutPaid,
			"error_message": "",
			"paid_at":       paidAt,
			"updated_at":    paidAt,
		}).Error
}

func (r *WeeklyBossPostgresRepo) MarkPayoutFailed(payoutID uint, reason string) error {
	return r.db.
		Model(&domain.WeeklyBossPayout{}).
		Where("id = ?", payoutID).
		Updates(map[string]any{
			"status":        domain.WeeklyBossPayoutFailed,
			"error_message": reason,
			"updated_at":    time.Now().UTC(),
		}).Error
}

func (r *WeeklyBossPostgresRepo) AddGoldToPlayer(playerID int64, gold int64) error {
	if gold <= 0 {
		return nil
	}
	return r.db.Exec(`
		UPDATE players
		   SET gold = gold + ?,
		       last_update = NOW()
		 WHERE id = ?`, gold, playerID).Error
}

func (r *WeeklyBossPostgresRepo) HasActive() (bool, error) {
	var count int64
	err := r.db.Model(&domain.WeeklyBoss{}).
		Where("status = ?", domain.WeeklyBossStatusActive).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *WeeklyBossPostgresRepo) FindLatestByStatus(status string) (*domain.WeeklyBoss, error) {
	var boss domain.WeeklyBoss
	err := r.db.Model(&domain.WeeklyBoss{}).
		Where("status = ?", status).
		Order("id DESC").
		First(&boss).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &boss, nil
}
