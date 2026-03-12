DROP INDEX IF EXISTS public.uq_weekly_boss_attacks_daily_seq;

ALTER TABLE IF EXISTS public.weekly_boss_attacks
    DROP COLUMN IF EXISTS daily_seq;
