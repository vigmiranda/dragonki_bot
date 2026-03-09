package postgres

import (
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VIPConfigRepo struct {
	db *gorm.DB
}

func NewVIPConfigRepo(db *gorm.DB) *VIPConfigRepo {
	return &VIPConfigRepo{db: db}
}

func (r *VIPConfigRepo) GetEnabledConfig() (*domain.VIPConfig, error) {
	var cfg domain.VIPConfig
	err := r.db.Model(&domain.VIPConfig{}).
		Where("enabled = ?", true).
		Order("id DESC").
		First(&cfg).Error
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

type PlayerVIPRepo struct {
	db *gorm.DB
}

func NewPlayerVIPRepo(db *gorm.DB) *PlayerVIPRepo {
	return &PlayerVIPRepo{db: db}
}

func (r *PlayerVIPRepo) GetByPlayerID(playerID int64) (*domain.PlayerVIP, error) {
	var vip domain.PlayerVIP
	err := r.db.Model(&domain.PlayerVIP{}).
		Where("player_id = ?", playerID).
		First(&vip).Error
	if err != nil {
		return nil, err
	}
	return &vip, nil
}

func (r *PlayerVIPRepo) Upsert(vip *domain.PlayerVIP) error {
	if vip == nil {
		return gorm.ErrInvalidData
	}
	now := time.Now().UTC()
	if vip.CreatedAt.IsZero() {
		vip.CreatedAt = now
	}
	vip.UpdatedAt = now

	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "player_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"started_at": vip.StartedAt,
			"expires_at": vip.ExpiresAt,
			"updated_at": vip.UpdatedAt,
		}),
	}).Create(vip).Error
}
