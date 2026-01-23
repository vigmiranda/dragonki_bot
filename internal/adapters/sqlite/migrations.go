package sqlite

type migration struct {
	version int
	sql     string
}

var migrations = []migration{
	{
		version: 1,
		sql: `
		CREATE TABLE IF NOT EXISTS players (
		user_id INTEGER PRIMARY KEY,
		name TEXT,
		gender TEXT,
		class TEXT,
		hp INTEGER,
		atk INTEGER,
		xp INTEGER,
		level INTEGER,
		energy INTEGER,
		last_energy_at TEXT,
		form TEXT,
		form_battles INTEGER
	);`,
	},
	{
		version: 2,
		sql: `
		CREATE TABLE IF NOT EXISTS monsters (
		monster_id INTEGER PRIMARY KEY,
		name TEXT,
		tag TEXT,
		level INTEGER,
		hp INTEGER,
		atk INTEGER,
		def INTEGER,
		max_xp INTEGER,
		max_money INTEGER, 
		spawn_weight INTEGER		
		);

		INSERT INTO monsters (name, tag, level, hp, atk, def, max_xp, max_money, spawn_weight)
		VALUES ('Slime Verde', 'SLIME_VERDE', 1, 100, 45, 35, 50, 100, 5);

		INSERT INTO monsters (name, tag, level, hp, atk, def, max_xp, max_money, spawn_weight)
		VALUES ('Guerreiro Lagarto', 'GUERREIRO_LAGARTO', 1, 80, 30, 20, 40, 80, 10);

		INSERT INTO monsters (name, tag, level, hp, atk, def, max_xp, max_money, spawn_weight)
		VALUES ('Goblin Errante', 'GOBLIN_ERRANTE', 1, 50, 20, 15, 25, 60, 20);

		INSERT INTO monsters (name, tag, level, hp, atk, def, max_xp, max_money, spawn_weight)
		VALUES ('Lobo das Planícies','LOBO_DAS_PLANICIES', 1, 40, 12, 5, 20, 40, 25);

		INSERT INTO monsters (name, tag, level, hp, atk, def, max_xp, max_money, spawn_weight)
		VALUES ('Morcego Sombrio', 'MORCEGO_SOMBRIO', 1, 30, 7, 6, 10, 20, 40);
		

	`,
	},
	{
		version: 3,
		sql: `
		CREATE TABLE IF NOT EXISTS regions (
			region_id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			tag TEXT NOT NULL,
			min_level INTEGER DEFAULT 1
		);
		ALTER TABLE players
		ADD COLUMN region_id INTEGER DEFAULT 1;
		INSERT INTO regions (region_id, name, tag, min_level)
		SELECT 1, 'Planícies', 'PLAIN_LV1', 1
		WHERE NOT EXISTS (
			SELECT 1 FROM regions WHERE region_id = 1
		);
		`,
	},
}
