package postgres

import (
	"database/sql"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
)

type PlayerPostgresRepo struct {
	db *sql.DB
}

func NewPlayerPostgresRepo(postgres *Postgres) *PlayerPostgresRepo {
	return &PlayerPostgresRepo{
		db: postgres.DB,
	}
}
func (r *PlayerPostgresRepo) GetOrCreate(userID int64) (*domain.Player, error) {
	var p domain.Player
	var last string

	err := r.db.QueryRow(`
	SELECT user_id, name, gender, class, hp, max_hp, atk, def, xp, level, energy, max_energy,
	       gold, ki_crystals,  region_id, last_energy_at, form, form_battles, onboarding_step
	FROM players WHERE user_id = $1`,
		userID,
	).Scan(
		&p.UserID,
		&p.Name,
		&p.Gender,
		&p.Class,
		&p.HP,
		&p.MaxHP,
		&p.ATK,
		&p.DEF,
		&p.XP,
		&p.Level,
		&p.Energy,
		&p.MaxEnergy,
		&p.Gold,
		&p.KiCrystals,
		&p.RegionID,
		&p.LastEnergyAt,
		&p.Form,
		&p.FormBattlesLeft,
		&p.OnboardingStep,
	)

	if err == sql.ErrNoRows {
		p = domain.Player{
			UserID:       userID,
			Name:         "User_",
			Gender:       "MAN",
			Class:        "Blade Monk",
			HP:           100,
			ATK:          10,
			Level:        1,
			Energy:       10,
			LastEnergyAt: time.Now(),
			Form:         "Base",
		}

		_, err = r.db.Exec(`
		INSERT INTO players (
			user_id, name, gender, class, hp, atk, xp, level,
			energy, last_energy_at, form, form_battles, region_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
			p.UserID,
			p.Name,
			p.Gender,
			p.Class,
			p.HP,
			p.ATK,
			p.XP,
			p.Level,
			p.Energy,
			p.LastEnergyAt.Format(time.RFC3339),
			p.Form,
			p.FormBattlesLeft,
			p.RegionID,
		)

		return &p, err
	}

	if err != nil {
		return nil, err
	}

	p.LastEnergyAt, _ = time.Parse(time.RFC3339, last)
	return &p, nil
}

func (r *PlayerPostgresRepo) Save(p *domain.Player) error {
	_, err := r.db.Exec(`
	UPDATE players
	SET name=$1, gender=$2, class=$3, hp=$4, max_hp=$5, atk=$6, def=$7, xp=$8, level=$9, energy=$10, max_energy=$11,
	    gold=$12, ki_crystals=$13, region_id=$14, last_energy_at=$15, form=$16, form_battles=$17, onboarding_step=$18
	WHERE user_id=$19`,
		p.Name,
		p.Gender,
		p.Class,
		p.HP,
		p.MaxHP,
		p.ATK,
		p.DEF,
		p.XP,
		p.Level,
		p.Energy,
		p.MaxEnergy,
		p.Gold,
		p.KiCrystals,
		p.RegionID,
		p.LastEnergyAt.Format(time.RFC3339),
		p.Form,
		p.FormBattlesLeft,
		p.OnboardingStep,
		p.UserID,
	)
	return err
}
