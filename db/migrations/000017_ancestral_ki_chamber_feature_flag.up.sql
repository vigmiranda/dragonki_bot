INSERT INTO public.feature_flags (name, enabled, description, created_at, updated_at)
VALUES
    ('ANCESTRAL_KI_CHAMBER', FALSE, 'Habilita o modulo Camara do Ki Ancestral', NOW(), NOW())
ON CONFLICT (name) DO NOTHING;
