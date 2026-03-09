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

func NewWeeklyBossPostgresRepo(db *gorm.DB) *WeeklyBossPostgresRepo {
	return &WeeklyBossPostgresRepo{db: db}
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
		Where("attack_date = ?", attackDate.UTC().Format("2006-01-02")).
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
	return r.db.Create(attack).Error
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
