ALTER TABLE IF EXISTS public.weekly_boss_attacks
    DROP CONSTRAINT IF EXISTS uq_weekly_boss_attacks_daily;

CREATE INDEX IF NOT EXISTS idx_weekly_boss_attacks_daily_lookup
    ON public.weekly_boss_attacks (boss_id, player_id, attack_date);
