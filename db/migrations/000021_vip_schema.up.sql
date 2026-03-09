CREATE TABLE IF NOT EXISTS vip_config (
    id BIGSERIAL PRIMARY KEY,
    enabled BOOLEAN NOT NULL DEFAULT FALSE,
    vip_price_ki_crystals INT NOT NULL DEFAULT 10,
    duration_days INT NOT NULL DEFAULT 7,
    xp_bonus_percent DOUBLE PRECISION NOT NULL DEFAULT 0,
    gold_bonus_percent DOUBLE PRECISION NOT NULL DEFAULT 0,
    sell_bonus_percent DOUBLE PRECISION NOT NULL DEFAULT 0,
    shop_discount_percent DOUBLE PRECISION NOT NULL DEFAULT 0,
    big_boss_daily_attack_limit INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE vip_config
    ADD COLUMN IF NOT EXISTS enabled BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS vip_price_ki_crystals INT NOT NULL DEFAULT 10,
    ADD COLUMN IF NOT EXISTS duration_days INT NOT NULL DEFAULT 7,
    ADD COLUMN IF NOT EXISTS xp_bonus_percent DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS gold_bonus_percent DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS sell_bonus_percent DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS shop_discount_percent DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS big_boss_daily_attack_limit INT NOT NULL DEFAULT 1,
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE vip_config
    ALTER COLUMN created_at SET DEFAULT NOW(),
    ALTER COLUMN updated_at SET DEFAULT NOW();

UPDATE vip_config
SET created_at = COALESCE(created_at, NOW()),
    updated_at = COALESCE(updated_at, NOW())
WHERE created_at IS NULL OR updated_at IS NULL;

CREATE TABLE IF NOT EXISTS player_vip (
    id BIGSERIAL PRIMARY KEY,
    player_id BIGINT NOT NULL UNIQUE,
    started_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE player_vip
    ADD COLUMN IF NOT EXISTS player_id BIGINT NOT NULL,
    ADD COLUMN IF NOT EXISTS started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS expires_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE player_vip
    ALTER COLUMN created_at SET DEFAULT NOW(),
    ALTER COLUMN updated_at SET DEFAULT NOW();

UPDATE player_vip
SET created_at = COALESCE(created_at, NOW()),
    updated_at = COALESCE(updated_at, NOW())
WHERE created_at IS NULL OR updated_at IS NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'player_vip_player_id_key'
    ) THEN
        ALTER TABLE player_vip
            ADD CONSTRAINT player_vip_player_id_key UNIQUE (player_id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_player_vip_expires_at
    ON player_vip (expires_at);

INSERT INTO vip_config (
    enabled,
    vip_price_ki_crystals,
    duration_days,
    xp_bonus_percent,
    gold_bonus_percent,
    sell_bonus_percent,
    shop_discount_percent,
    big_boss_daily_attack_limit
)
SELECT
    FALSE, 10, 7, 0, 0, 0, 0, 1
WHERE NOT EXISTS (SELECT 1 FROM vip_config);
