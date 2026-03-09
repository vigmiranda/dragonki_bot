DROP INDEX IF EXISTS idx_weekly_boss_attacks_daily_lookup;

ALTER TABLE IF EXISTS public.weekly_boss_attacks
    ADD CONSTRAINT uq_weekly_boss_attacks_daily UNIQUE (boss_id, player_id, attack_date);
