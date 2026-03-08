package postgres

import (
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
)

type ReferralPostgresRepo struct {
	db *gorm.DB
}

func NewReferralPostgresRepo(db *gorm.DB) *ReferralPostgresRepo {
	return &ReferralPostgresRepo{db: db}
}

func (r *ReferralPostgresRepo) Create(referral *domain.Referral) error {
	return r.db.Create(referral).Error
}

func (r *ReferralPostgresRepo) GetByReferredPlayerID(referredPlayerID int64) (*domain.Referral, error) {
	var referral domain.Referral
	err := r.db.
		Model(&domain.Referral{}).
		Where("referred_player_id = ?", referredPlayerID).
		First(&referral).Error
	if err != nil {
		return nil, err
	}
	return &referral, nil
}

func (r *ReferralPostgresRepo) MarkRewarded(referralID uint) error {
	now := time.Now().UTC()
	return r.db.
		Model(&domain.Referral{}).
		Where("id = ?", referralID).
		Updates(map[string]any{
			"status":       domain.ReferralStatusRewarded,
			"completed_at": now,
			"rewarded_at":  now,
		}).Error
}

func (r *ReferralPostgresRepo) GrantReward(referrerPlayerID int64, gold int, kiCrystals int) error {
	updates := map[string]any{
		"last_update": time.Now().UTC(),
	}
	if gold > 0 {
		updates["gold"] = gorm.Expr("gold + ?", gold)
	}
	if kiCrystals > 0 {
		updates["ki_crystals"] = gorm.Expr("ki_crystals + ?", kiCrystals)
	}

	return r.db.
		Model(&domain.Player{}).
		Where("id = ?", referrerPlayerID).
		Updates(updates).Error
}

func (r *ReferralPostgresRepo) CountByReferrer(referrerPlayerID int64) (int64, int64, error) {
	var pending int64
	if err := r.db.
		Model(&domain.Referral{}).
		Where("referrer_player_id = ?", referrerPlayerID).
		Where("status = ?", domain.ReferralStatusPending).
		Count(&pending).Error; err != nil {
		return 0, 0, err
	}

	var rewarded int64
	if err := r.db.
		Model(&domain.Referral{}).
		Where("referrer_player_id = ?", referrerPlayerID).
		Where("status = ?", domain.ReferralStatusRewarded).
		Count(&rewarded).Error; err != nil {
		return 0, 0, err
	}

	return pending, rewarded, nil
}

func (r *ReferralPostgresRepo) GetPlayerUserIDByPlayerID(playerID int64) (int64, error) {
	type playerUserIDRow struct {
		UserID int64
	}
	var row playerUserIDRow
	err := r.db.
		Model(&domain.Player{}).
		Select("user_id").
		Where("id = ?", playerID).
		Take(&row).Error
	if err != nil {
		return 0, err
	}
	return row.UserID, nil
}
