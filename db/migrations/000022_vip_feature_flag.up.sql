INSERT INTO public.feature_flags (name, enabled, description, created_at, updated_at)
VALUES
    ('VIP', FALSE, 'Habilita menu e compra de assinatura VIP', NOW(), NOW())
ON CONFLICT (name) DO NOTHING;
