package postgres

type migration struct {
	version int
	sql     string
}

var migrations = []migration{
	{
		version: 1,
		sql: `
		CREATE TABLE IF NOT EXISTS players (
			user_id BIGINT PRIMARY KEY,
			name TEXT,
			gender TEXT,
			class TEXT,

			hp INTEGER,
			max_hp INTEGER DEFAULT 100,

			atk INTEGER,
			def INTEGER DEFAULT 20,

			xp INTEGER,
			level INTEGER,

			energy INTEGER,
			max_energy INTEGER DEFAULT 10,

			gold INTEGER DEFAULT 0,
			ki_crystals INTEGER DEFAULT 0,

			region_id BIGINT DEFAULT 1,

			last_energy_at TIMESTAMP,
			form TEXT,
			form_battles INTEGER
		);
		`,
	},
	{
		version: 2,
		sql: `
		CREATE TABLE IF NOT EXISTS monsters (
			monster_id BIGSERIAL PRIMARY KEY,
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
		VALUES
		('Slime Verde', 'SLIME_VERDE', 1, 100, 45, 35, 50, 100, 5),
		('Guerreiro Lagarto', 'GUERREIRO_LAGARTO', 1, 80, 30, 20, 40, 80, 10),
		('Goblin Errante', 'GOBLIN_ERRANTE', 1, 50, 20, 15, 25, 60, 20),
		('Lobo das Planícies', 'LOBO_DAS_PLANICIES', 1, 40, 12, 5, 20, 40, 25),
		('Morcego Sombrio', 'MORCEGO_SOMBRIO', 1, 30, 7, 6, 10, 20, 40);
		`,
	},
	{
		version: 3,
		sql: `
		CREATE TABLE IF NOT EXISTS regions (
			region_id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			tag TEXT NOT NULL,
			min_level INTEGER DEFAULT 1
		);

		INSERT INTO regions (region_id, name, tag, min_level)
		SELECT 1, 'Planícies', 'PLAIN_LV1', 1
		WHERE NOT EXISTS (
			SELECT 1 FROM regions WHERE region_id = 1
		);
		`,
	},
	{
		version: 4,
		sql: `
		ALTER TABLE players ALTER COLUMN region_id SET DEFAULT 1;
		ALTER TABLE players ADD COLUMN onboarding_step TEXT DEFAULT 'START';
		ALTER TABLE players ADD COLUMN last_update TIMESTAMP DEFAULT now();

		CREATE OR REPLACE FUNCTION set_last_update()
		RETURNS TRIGGER AS $$
		BEGIN
		  NEW.last_update = now();
		  RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
		   
		CREATE TRIGGER trigger_set_last_update
		BEFORE UPDATE ON players
		FOR EACH ROW
		EXECUTE FUNCTION set_last_update();   

		`,
	},
	{
		version: 5,
		sql: `
		ALTER TABLE players ALTER COLUMN onboarding_step SET DEFAULT 'CHOOSING_NAME';

		`,
	},
}
