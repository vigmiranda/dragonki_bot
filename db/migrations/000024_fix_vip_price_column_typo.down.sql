DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'vip_config'
          AND column_name = 'vip_price_ki_crystals'
    ) AND NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'vip_config'
          AND column_name = 'vi_price_ki_crystals'
    ) THEN
        EXECUTE 'ALTER TABLE public.vip_config RENAME COLUMN vip_price_ki_crystals TO vi_price_ki_crystals';
    END IF;
END $$;
