package postgres

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
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

	err := r.db.QueryRow(`
	SELECT user_id,name, gender, class, hp, max_hp, atk, def, xp, level, energy,
	    max_energy,gold, ki_crystals, region_id, last_energy_at, form, 
	    form_battles, onboarding_step, last_update
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
		&p.LastUpdate,
	)

	if errors.Is(err, sql.ErrNoRows) {
		p = domain.Player{
			UserID:          userID,
			Name:            "User_" + uuid.New().String()[:6],
			Gender:          "MAN",
			Class:           "Blade Monk",
			HP:              100,
			MaxHP:           100,
			ATK:             10,
			DEF:             5,
			XP:              0,
			Level:           1,
			Energy:          10,
			MaxEnergy:       10,
			Gold:            0,
			KiCrystals:      0,
			RegionID:        1,
			LastEnergyAt:    time.Now().UTC(),
			Form:            "Base",
			FormBattlesLeft: 0,
			OnboardingStep:  "DONE",
			LastUpdate:      time.Now().UTC(),
		}

		_, err = r.db.Exec(`
		INSERT INTO players (
			user_id, name, gender, class, hp, max_hp, atk, def, xp, level, energy,
	    	max_energy,gold, ki_crystals, region_id, last_energy_at, form, 
	    	form_battles, onboarding_step, last_update
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`,
			p.UserID,
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
			p.LastEnergyAt,
			p.Form,
			p.FormBattlesLeft,
			p.OnboardingStep,
			p.LastUpdate,
		)

		return &p, err
	}

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *PlayerPostgresRepo) Save(p *domain.Player) error {
	_, err := r.db.Exec(`
	UPDATE players
	SET name=$1, gender=$2, class=$3, hp=$4, max_hp=$5, atk=$6, def=$7, xp=$8, level=$9, energy=$10,
	    max_energy=$11,gold=$12, ki_crystals=$13, region_id=$14, last_energy_at=$15, form=$16, 
	    form_battles=$17, onboarding_step=$18, last_update= now()
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
		p.LastEnergyAt,
		p.Form,
		p.FormBattlesLeft,
		p.OnboardingStep,
		p.UserID,
	)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
