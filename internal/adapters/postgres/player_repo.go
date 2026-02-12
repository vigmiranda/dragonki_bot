package postgres

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
)

type PlayerPostgresRepo struct {
	db *gorm.DB
}

func NewPlayerPostgresRepo(db *gorm.DB) *PlayerPostgresRepo {
	return &PlayerPostgresRepo{
		db: db,
	}
}
func (r *PlayerPostgresRepo) GetOrCreate(userID int64) (*domain.Player, error) {

	var existing domain.Player

	err := r.db.
		Where("user_id = ?", userID).
		First(&existing).
		Error

	// ✅ Encontrou
	if err == nil {
		return &existing, nil
	}

	// ❌ Erro diferente de "não encontrado"
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 🆕 NOVA struct (não reutilizar!)
	newPlayer := domain.Player{
		UserID:          userID,
		Name:            "User_" + uuid.New().String()[:6],
		Gender:          "MAN",
		ClassID:         1,
		ClassName:       "Blade Monk",
		Level:           1,
		XP:              0,
		HP:              100,
		MaxHP:           100,
		ATK:             10,
		DEF:             5,
		Energy:          10,
		MaxEnergy:       10,
		Gold:            0,
		KiCrystals:      0,
		RegionID:        1,
		LastEnergyAt:    time.Now().UTC(),
		Form:            "Base",
		FormBattlesLeft: 0,
		OnboardingStep:  "ONBOARDING",
		LastUpdate:      time.Now().UTC(),
	}
	newPlayer.ID = 0

	if err := r.db.Create(&newPlayer).Error; err != nil {
		return nil, err
	}

	return &newPlayer, nil
}

func (r *PlayerPostgresRepo) Update(p *domain.Player) error {
	updates := map[string]interface{}{
		"class_id":          p.ClassID,
		"class_name":        p.ClassName,
		"level":             p.Level,
		"xp":                p.XP,
		"hp":                p.HP,
		"max_hp":            p.MaxHP,
		"atk":               p.ATK,
		"def":               p.DEF,
		"energy":            p.Energy,
		"max_energy":        p.MaxEnergy,
		"gold":              p.Gold,
		"ki_crystals":       p.KiCrystals,
		"region_id":         p.RegionID,
		"last_energy_at":    time.Now(),
		"form":              p.Form,
		"form_battles_left": p.FormBattlesLeft,
		"last_update":       time.Now(),
	}

	return r.db.
		Model(&domain.Player{}).
		Where("id = ?", p.ID).
		Updates(updates).
		Error
}
func (r *PlayerPostgresRepo) CountOnlineLast15Min() (int, error) {
	var count int64

	err := r.db.
		Model(&domain.Player{}).
		Where("last_update >= NOW() - INTERVAL '15 minutes'").
		Count(&count).
		Error

	return int(count), err
}
