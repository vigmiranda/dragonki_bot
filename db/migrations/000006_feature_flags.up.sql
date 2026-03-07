CREATE TABLE IF NOT EXISTS public.feature_flags (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(120) NOT NULL UNIQUE,
    enabled BOOLEAN NOT NULL DEFAULT FALSE,
    description VARCHAR(500) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE public.feature_flags
    ALTER COLUMN created_at SET DEFAULT NOW(),
    ALTER COLUMN updated_at SET DEFAULT NOW();

UPDATE public.feature_flags
SET created_at = COALESCE(created_at, NOW()),
    updated_at = COALESCE(updated_at, NOW())
WHERE created_at IS NULL OR updated_at IS NULL;

INSERT INTO public.feature_flags (name, enabled, description, created_at, updated_at)
VALUES
    ('PAYMENTS', TRUE, 'Habilita compra de Ki Crystals e webhook de pagamentos', NOW(), NOW())
ON CONFLICT (name) DO NOTHING;
