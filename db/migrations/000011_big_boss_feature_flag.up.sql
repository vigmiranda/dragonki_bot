INSERT INTO public.feature_flags (name, enabled, description, created_at, updated_at)
VALUES
    ('BIG_BOSS_WEEKLY', FALSE, 'Habilita o modulo semanal de Big Boss', NOW(), NOW())
ON CONFLICT (name) DO NOTHING;
