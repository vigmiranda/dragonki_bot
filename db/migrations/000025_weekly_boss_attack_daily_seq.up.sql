ALTER TABLE IF EXISTS public.weekly_boss_attacks
    ADD COLUMN IF NOT EXISTS daily_seq INTEGER NOT NULL DEFAULT 1;

ALTER TABLE IF EXISTS public.weekly_boss_attacks
    DROP CONSTRAINT IF EXISTS uq_weekly_boss_attacks_daily;

DROP INDEX IF EXISTS public.uq_weekly_boss_attacks_daily;

WITH ranked_attacks AS (
    SELECT
        id,
        ROW_NUMBER() OVER (
            PARTITION BY boss_id, player_id, attack_date
            ORDER BY created_at, id
        ) AS seq
    FROM public.weekly_boss_attacks
)
UPDATE public.weekly_boss_attacks AS attacks
   SET daily_seq = ranked_attacks.seq
  FROM ranked_attacks
 WHERE attacks.id = ranked_attacks.id;

CREATE UNIQUE INDEX IF NOT EXISTS uq_weekly_boss_attacks_daily_seq
    ON public.weekly_boss_attacks (boss_id, player_id, attack_date, daily_seq);
