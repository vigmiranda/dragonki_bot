CREATE TABLE IF NOT EXISTS public.ancestral_ki_chamber_progress (
    id BIGSERIAL PRIMARY KEY,
    player_id BIGINT NOT NULL REFERENCES public.players (id) ON DELETE CASCADE,
    highest_completed_level INTEGER NOT NULL DEFAULT 0 CHECK (highest_completed_level >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_ak_chamber_progress_player UNIQUE (player_id)
);

CREATE INDEX IF NOT EXISTS idx_ak_chamber_progress_player_id
    ON public.ancestral_ki_chamber_progress (player_id);

CREATE TABLE IF NOT EXISTS public.ancestral_ki_chamber_runs (
    id BIGSERIAL PRIMARY KEY,
    player_id BIGINT NOT NULL REFERENCES public.players (id) ON DELETE CASCADE,
    chamber_level INTEGER NOT NULL CHECK (chamber_level BETWEEN 1 AND 10),
    current_battle INTEGER NOT NULL DEFAULT 1 CHECK (current_battle BETWEEN 1 AND 5),
    victories INTEGER NOT NULL DEFAULT 0 CHECK (victories BETWEEN 0 AND 5),
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_ak_chamber_runs_status
        CHECK (status IN ('ACTIVE', 'COMPLETED', 'FAILED', 'ABANDONED'))
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_ak_chamber_runs_player_active
    ON public.ancestral_ki_chamber_runs (player_id)
    WHERE status = 'ACTIVE';

CREATE INDEX IF NOT EXISTS idx_ak_chamber_runs_player_id
    ON public.ancestral_ki_chamber_runs (player_id);

CREATE INDEX IF NOT EXISTS idx_ak_chamber_runs_status
    ON public.ancestral_ki_chamber_runs (status);

CREATE INDEX IF NOT EXISTS idx_ak_chamber_runs_started_at
    ON public.ancestral_ki_chamber_runs (started_at DESC);

CREATE TABLE IF NOT EXISTS public.ancestral_ki_chamber_battle_logs (
    id BIGSERIAL PRIMARY KEY,
    chamber_run_id BIGINT NOT NULL REFERENCES public.ancestral_ki_chamber_runs (id) ON DELETE CASCADE,
    battle_number INTEGER NOT NULL CHECK (battle_number BETWEEN 1 AND 5),
    monster_id BIGINT NULL,
    monster_level INTEGER NOT NULL CHECK (monster_level >= 1),
    result VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_ak_chamber_battle_logs_result
        CHECK (result IN ('WIN', 'LOSE'))
);

CREATE INDEX IF NOT EXISTS idx_ak_chamber_battle_logs_run_id
    ON public.ancestral_ki_chamber_battle_logs (chamber_run_id);

CREATE INDEX IF NOT EXISTS idx_ak_chamber_battle_logs_created_at
    ON public.ancestral_ki_chamber_battle_logs (created_at DESC);
