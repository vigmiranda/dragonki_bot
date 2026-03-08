CREATE TABLE IF NOT EXISTS public.weekly_bosses (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(120) NOT NULL,
    week_key VARCHAR(16) NOT NULL UNIQUE,
    max_hp BIGINT NOT NULL CHECK (max_hp > 0),
    current_hp BIGINT NOT NULL CHECK (current_hp >= 0),
    atk INTEGER NOT NULL CHECK (atk > 0),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    defeated_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_weekly_bosses_status
        CHECK (status IN ('active', 'defeated', 'expired')),
    CONSTRAINT chk_weekly_bosses_time_window
        CHECK (ends_at > starts_at)
);

CREATE INDEX IF NOT EXISTS idx_weekly_bosses_status ON public.weekly_bosses (status);
CREATE INDEX IF NOT EXISTS idx_weekly_bosses_starts_at ON public.weekly_bosses (starts_at);
CREATE INDEX IF NOT EXISTS idx_weekly_bosses_ends_at ON public.weekly_bosses (ends_at);

CREATE TABLE IF NOT EXISTS public.weekly_boss_attacks (
    id BIGSERIAL PRIMARY KEY,
    boss_id BIGINT NOT NULL REFERENCES public.weekly_bosses (id) ON DELETE CASCADE,
    player_id BIGINT NOT NULL,
    attack_date DATE NOT NULL,
    damage BIGINT NOT NULL CHECK (damage >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_weekly_boss_attacks_daily UNIQUE (boss_id, player_id, attack_date)
);

CREATE INDEX IF NOT EXISTS idx_weekly_boss_attacks_boss_id ON public.weekly_boss_attacks (boss_id);
CREATE INDEX IF NOT EXISTS idx_weekly_boss_attacks_player_id ON public.weekly_boss_attacks (player_id);
CREATE INDEX IF NOT EXISTS idx_weekly_boss_attacks_attack_date ON public.weekly_boss_attacks (attack_date);

CREATE TABLE IF NOT EXISTS public.weekly_boss_contributions (
    id BIGSERIAL PRIMARY KEY,
    boss_id BIGINT NOT NULL REFERENCES public.weekly_bosses (id) ON DELETE CASCADE,
    player_id BIGINT NOT NULL,
    total_damage BIGINT NOT NULL DEFAULT 0 CHECK (total_damage >= 0),
    last_attack_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_weekly_boss_contributions UNIQUE (boss_id, player_id)
);

CREATE INDEX IF NOT EXISTS idx_weekly_boss_contributions_boss_id ON public.weekly_boss_contributions (boss_id);
CREATE INDEX IF NOT EXISTS idx_weekly_boss_contributions_player_id ON public.weekly_boss_contributions (player_id);
CREATE INDEX IF NOT EXISTS idx_weekly_boss_contributions_total_damage ON public.weekly_boss_contributions (boss_id, total_damage DESC);
