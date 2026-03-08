CREATE TABLE IF NOT EXISTS public.weekly_boss_templates (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(120) NOT NULL,
    tag VARCHAR(64) NOT NULL UNIQUE,
    max_hp BIGINT NOT NULL CHECK (max_hp > 0),
    atk INTEGER NOT NULL CHECK (atk > 0),
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_weekly_boss_templates_enabled
    ON public.weekly_boss_templates (enabled);

INSERT INTO public.weekly_boss_templates (name, tag, max_hp, atk, enabled, created_at, updated_at)
VALUES
    ('Astarion, o Senhor da Ruina', 'ASTARION', 2150000, 190, TRUE, NOW(), NOW()),
    ('Surtr, o Devorador do Mundo', 'SURTR', 2000000, 180, TRUE, NOW(), NOW()),
    ('Vorgath, o Colosso Eterno', 'VORGATH', 2300000, 205, TRUE, NOW(), NOW()),
    ('Nyxara, a Rainha do Vazio', 'NYXARA', 2450000, 220, TRUE, NOW(), NOW())
ON CONFLICT (tag) DO UPDATE
SET name = EXCLUDED.name,
    max_hp = EXCLUDED.max_hp,
    atk = EXCLUDED.atk,
    enabled = EXCLUDED.enabled,
    updated_at = NOW();
