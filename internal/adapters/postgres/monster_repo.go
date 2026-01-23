package postgres

import (
	"database/sql"

	"github.com/vigmiranda/dragonki/internal/domain"
)

type MonsterPostgresRepo struct {
	db *sql.DB
}

func NewMonsterPostgresRepo(postgres *Postgres) *MonsterPostgresRepo {
	return &MonsterPostgresRepo{
		db: postgres.DB,
	}
}

func (r *MonsterPostgresRepo) ListByLevel(level int64) ([]domain.Monster, error) {
	rows, err := r.db.Query(`
		SELECT monster_id, name, tag, level, hp, atk, def,
		       max_xp, max_money, spawn_weight
		FROM monsters
		WHERE level <= $1
		ORDER BY monster_id`,
		level,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monsters []domain.Monster

	for rows.Next() {
		var m domain.Monster
		err := rows.Scan(
			&m.MonsterID,
			&m.Name,
			&m.Tag,
			&m.Level,
			&m.HP,
			&m.ATK,
			&m.DEF,
			&m.MaxXP,
			&m.MaxMoney,
			&m.SpawnWeight,
		)
		if err != nil {
			return nil, err
		}
		monsters = append(monsters, m)
	}

	return monsters, nil
}
