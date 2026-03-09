ALTER TABLE ancestral_ki_chamber_runs
    ADD COLUMN IF NOT EXISTS accumulated_xp INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS accumulated_gold INT NOT NULL DEFAULT 0;
