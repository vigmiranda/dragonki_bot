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
		`,
	},
	{
		version: 5,
		sql: `
		ALTER TABLE players ALTER COLUMN onboarding_step SET DEFAULT 'CHOOSING_NAME';

		`,
	},
	{
		version: 6,
		sql: `
		ALTER TABLE players ALTER COLUMN last_energy_at SET DEFAULT now();
		ALTER TABLE players ALTER COLUMN last_energy_at SET NOT NULL;
		ALTER TABLE players ALTER COLUMN last_update SET NOT NULL;

		`,
	},
	{
		version: 7,
		sql: `
		CREATE TABLE payments (
			id TEXT PRIMARY KEY,
			player_id BIGINT NOT NULL,
			amount NUMERIC NOT NULL,
			status TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL
		);
		`,
	},
	{
		version: 8,
		sql: `
		ALTER TABLE payments ADD COLUMN qr_code TEXT;
		ALTER TABLE payments ADD COLUMN qr_code_base64 TEXT;
		`,
	},
	{
		version: 9,
		sql: `
		ALTER TABLE payments ADD COLUMN partner_id TEXT;
		`,
	},
	{
		version: 10,
		sql: `
		ALTER TABLE payments
		ALTER COLUMN partner_id TYPE BIGINT USING (partner_id::BIGINT);
		`,
	},
	{
		version: 11,
		sql: `		
		INSERT INTO regions (region_id, name, tag, min_level)
		SELECT 2, 'Forest of Spirits', 'FLOREST_LV15', 15
		WHERE NOT EXISTS (
			SELECT 2 FROM regions WHERE region_id = 2
		);
		INSERT INTO regions (region_id, name, tag, min_level)
		SELECT 3, 'Volcano Lands', 'VULCAO_LV35', 35
		WHERE NOT EXISTS (
			SELECT 3 FROM regions WHERE region_id = 3
		);
		INSERT INTO regions (region_id, name, tag, min_level)
		SELECT 4, 'Frozen Peaks', 'FROZEN_LV50', 50
		WHERE NOT EXISTS (
			SELECT 4 FROM regions WHERE region_id = 4
		);
		INSERT INTO regions (region_id, name, tag, min_level)
		SELECT 5, 'Swamp of Shadows', 'SWAMP_LV70', 70
		WHERE NOT EXISTS (
			SELECT 5 FROM regions WHERE region_id = 5
		);
		`,
	},
	{
		version: 12,
		sql: `	
		CREATE TABLE IF NOT EXISTS classes (
			id SERIAL PRIMARY KEY,
			name VARCHAR(50) NOT NULL,
			base_class VARCHAR(50) NOT NULL,
			tier VARCHAR(20) NOT NULL,
			min_level INTEGER NOT NULL,
		
			atk_bonus INTEGER DEFAULT 0,
			hp_bonus INTEGER DEFAULT 0,
			energy_bonus INTEGER DEFAULT 0
		);

		ALTER TABLE players RENAME COLUMN class TO ClassName ;

		ALTER TABLE players
		ADD COLUMN class_id INTEGER REFERENCES classes(id) DEFAULT 1,
		ADD COLUMN can_evolve BOOLEAN DEFAULT false;
		    
		
		-- Primeira escolha
		INSERT INTO classes (name, base_class, tier, min_level)
		VALUES ('Blade Monk', 'Blade Monk', 'BASE', 1);
		
		-- Primeira escolha
		INSERT INTO classes (name, base_class, tier, min_level, atk_bonus)
		VALUES ('Sword Disciple', 'Blade Monk', 'ADVANCED', 15, 5);
		
		INSERT INTO classes (name, base_class, tier, min_level, hp_bonus)
		VALUES ('Spirit Guard', 'Blade Monk', 'ADVANCED', 15, 10);
		
		INSERT INTO classes (name, base_class, tier, min_level, energy_bonus)
		VALUES ('Ki Adept', 'Blade Monk', 'ADVANCED', 15, 5);
		
		-- Evoluções
		INSERT INTO classes (name, base_class, tier, min_level, atk_bonus)
		VALUES ('Blade Master', 'Sword Disciple', 'ELITE', 35, 10);
		
		INSERT INTO classes (name, base_class, tier, min_level, hp_bonus)
		VALUES ('Soul Protector', 'Spirit Guard', 'ELITE', 35, 20);
		
		INSERT INTO classes (name, base_class, tier, min_level, energy_bonus)
		VALUES ('Ki Channeler', 'Ki Adept', 'ELITE', 35, 15);
		`,
	},
	{
		version: 13,
		sql: `	
		INSERT INTO classes (name, base_class, tier, min_level, atk_bonus)
		VALUES ('Dragon Slayer', 'Blade Master', 'MASTER', 50, 20);
		
		INSERT INTO classes (name, base_class, tier, min_level, hp_bonus)
		VALUES ('Eternal Sentinel', 'Soul Protector', 'MASTER', 50, 35);
		
		INSERT INTO classes (name, base_class, tier, min_level, energy_bonus)
		VALUES ('Ascended Sage', 'Ki Channeler', 'MASTER', 50, 25);
		`,
	},
}
