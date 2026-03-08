CREATE TABLE IF NOT EXISTS public.weekly_boss_rewards (
    id BIGSERIAL PRIMARY KEY,
    position INTEGER NOT NULL UNIQUE CHECK (position BETWEEN 1 AND 3),
    gold BIGINT NOT NULL DEFAULT 0 CHECK (gold >= 0),
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS public.weekly_boss_reward_items (
    id BIGSERIAL PRIMARY KEY,
    reward_id BIGINT NOT NULL REFERENCES public.weekly_boss_rewards (id) ON DELETE CASCADE,
    item_id BIGINT NOT NULL REFERENCES public.items (id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_weekly_boss_reward_items UNIQUE (reward_id, item_id)
);

CREATE TABLE IF NOT EXISTS public.weekly_boss_payouts (
    id BIGSERIAL PRIMARY KEY,
    boss_id BIGINT NOT NULL REFERENCES public.weekly_bosses (id) ON DELETE CASCADE,
    player_id BIGINT NOT NULL,
    position INTEGER NOT NULL CHECK (position BETWEEN 1 AND 3),
    gold_paid BIGINT NOT NULL DEFAULT 0 CHECK (gold_paid >= 0),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    error_message TEXT NOT NULL DEFAULT '',
    paid_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_weekly_boss_payouts_boss_player UNIQUE (boss_id, player_id),
    CONSTRAINT uq_weekly_boss_payouts_boss_position UNIQUE (boss_id, position),
    CONSTRAINT chk_weekly_boss_payouts_status
        CHECK (status IN ('pending', 'paid', 'failed'))
);

CREATE INDEX IF NOT EXISTS idx_weekly_boss_payouts_boss_id
    ON public.weekly_boss_payouts (boss_id);

INSERT INTO public.weekly_boss_rewards (position, gold, enabled, created_at, updated_at)
VALUES
    (1, 10000, TRUE, NOW(), NOW()),
    (2, 5000, TRUE, NOW(), NOW()),
    (3, 2500, TRUE, NOW(), NOW())
ON CONFLICT (position) DO UPDATE
SET gold = EXCLUDED.gold,
    enabled = EXCLUDED.enabled,
    updated_at = NOW();

WITH r1 AS (
    SELECT id FROM public.weekly_boss_rewards WHERE position = 1
), i1 AS (
    SELECT id FROM public.items WHERE code = 'WEAPON_MASTER' ORDER BY id LIMIT 1
)
INSERT INTO public.weekly_boss_reward_items (reward_id, item_id, quantity)
SELECT r1.id, i1.id, 1 FROM r1, i1
ON CONFLICT (reward_id, item_id) DO NOTHING;

WITH r1 AS (
    SELECT id FROM public.weekly_boss_rewards WHERE position = 1
), i1 AS (
    SELECT id FROM public.items WHERE code = 'RING_MASTER' ORDER BY id LIMIT 1
)
INSERT INTO public.weekly_boss_reward_items (reward_id, item_id, quantity)
SELECT r1.id, i1.id, 1 FROM r1, i1
ON CONFLICT (reward_id, item_id) DO NOTHING;

WITH r2 AS (
    SELECT id FROM public.weekly_boss_rewards WHERE position = 2
), i2 AS (
    SELECT id FROM public.items WHERE code = 'WEAPON_WARRIOR' ORDER BY id LIMIT 1
)
INSERT INTO public.weekly_boss_reward_items (reward_id, item_id, quantity)
SELECT r2.id, i2.id, 1 FROM r2, i2
ON CONFLICT (reward_id, item_id) DO NOTHING;
