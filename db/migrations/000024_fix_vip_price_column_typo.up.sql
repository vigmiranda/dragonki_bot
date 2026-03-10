DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'vip_config'
          AND column_name = 'vi_price_ki_crystals'
    ) THEN
        IF EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'vip_config'
              AND column_name = 'vip_price_ki_crystals'
        ) THEN
            EXECUTE '
                UPDATE public.vip_config
                   SET vip_price_ki_crystals = vi_price_ki_crystals
                 WHERE (vip_price_ki_crystals IS NULL OR vip_price_ki_crystals = 0)
                   AND vi_price_ki_crystals IS NOT NULL
            ';
            EXECUTE 'ALTER TABLE public.vip_config DROP COLUMN vi_price_ki_crystals';
        ELSE
            EXECUTE 'ALTER TABLE public.vip_config RENAME COLUMN vi_price_ki_crystals TO vip_price_ki_crystals';
        END IF;
    END IF;
END $$;

ALTER TABLE IF EXISTS public.vip_config
    ADD COLUMN IF NOT EXISTS vip_price_ki_crystals INT NOT NULL DEFAULT 10;
