CREATE TABLE IF NOT EXISTS public.referrals (
    id BIGSERIAL PRIMARY KEY,
    referrer_player_id BIGINT NOT NULL REFERENCES public.players (id) ON DELETE CASCADE,
    referred_player_id BIGINT NOT NULL REFERENCES public.players (id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    required_level INTEGER NOT NULL DEFAULT 7 CHECK (required_level > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ NULL,
    rewarded_at TIMESTAMPTZ NULL,
    CONSTRAINT uq_referrals_referred_player UNIQUE (referred_player_id),
    CONSTRAINT chk_referrals_status CHECK (status IN ('PENDING', 'REWARDED')),
    CONSTRAINT chk_referrals_no_self CHECK (referrer_player_id <> referred_player_id)
);

CREATE INDEX IF NOT EXISTS idx_referrals_referrer_player_id
    ON public.referrals (referrer_player_id);

CREATE INDEX IF NOT EXISTS idx_referrals_status
    ON public.referrals (status);
