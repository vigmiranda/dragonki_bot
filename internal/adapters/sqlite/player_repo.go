package sqlite

import (
	"database/sql"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
)

type PlayerSQLiteRepo struct {
	db *sql.DB
}

func NewPlayerSQLiteRepo(sqlite *SQLite) *PlayerSQLiteRepo {
	return &PlayerSQLiteRepo{
		db: sqlite.DB,
	}
}
func (r *PlayerSQLiteRepo) GetOrCreate(userID int64) (*domain.Player, error) {
	var p domain.Player
	var last string

	err := r.db.QueryRow(`
	SELECT user_id, name, gender, class, hp, atk, xp, level, energy, last_energy_at, form, form_battles
	FROM players WHERE user_id = ?`,
		userID,
	).Scan(
		&p.UserID,
		&p.Name,
		&p.Gender,
		&p.Class,
		&p.HP,
		&p.ATK,
		&p.XP,
		&p.Level,
		&p.Energy,
		&last,
		&p.Form,
		&p.FormBattlesLeft,
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
			energy, last_energy_at, form, form_battles
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
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
		)

		return &p, err
	}

	if err != nil {
		return nil, err
	}

	p.LastEnergyAt, _ = time.Parse(time.RFC3339, last)
	return &p, nil
}

func (r *PlayerSQLiteRepo) Save(p *domain.Player) error {
	_, err := r.db.Exec(`
	UPDATE players
	SET name=?, gender=?, hp=?, atk=?, xp=?, level=?, energy=?, last_energy_at=?, form=?, form_battles=?
	WHERE user_id=?`,
		p.Name,
		p.Gender,
		p.HP,
		p.ATK,
		p.XP,
		p.Level,
		p.Energy,
		p.LastEnergyAt.Format(time.RFC3339),
		p.Form,
		p.FormBattlesLeft,
		p.UserID,
	)
	return err
}
