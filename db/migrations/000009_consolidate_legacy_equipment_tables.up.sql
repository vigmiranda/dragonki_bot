CREATE TABLE IF NOT EXISTS public.equipments (
    id BIGSERIAL PRIMARY KEY,
    item_id BIGINT NOT NULL,
    slot VARCHAR(32) NOT NULL DEFAULT '',
    required_level INTEGER NOT NULL DEFAULT 1,
    atk_bonus INTEGER NOT NULL DEFAULT 0,
    def_bonus INTEGER NOT NULL DEFAULT 0,
    speed_bonus INTEGER NOT NULL DEFAULT 0,
    hp_bonus INTEGER NOT NULL DEFAULT 0,
    energy_bonus INTEGER NOT NULL DEFAULT 0,
    allowed_classes TEXT NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_equipments_item_id ON public.equipments(item_id);

DO $$
DECLARE
    legacy_table TEXT;
    has_required_level BOOLEAN;
    has_atk_bonus BOOLEAN;
    has_def_bonus BOOLEAN;
    has_speed_bonus BOOLEAN;
    has_hp_bonus BOOLEAN;
    has_energy_bonus BOOLEAN;
    has_allowed_classes BOOLEAN;
    has_bonus_atk BOOLEAN;
    has_bonus_def BOOLEAN;
    has_bonus_hp BOOLEAN;
    has_bonus_energy BOOLEAN;
    insert_sql TEXT;
BEGIN
    FOR legacy_table IN SELECT unnest(ARRAY['equipament', 'equipaments']) LOOP
        IF to_regclass('public.' || legacy_table) IS NULL THEN
            CONTINUE;
        END IF;

        SELECT EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = 'public' AND table_name = legacy_table AND column_name = 'required_level'
        ) INTO has_required_level;
        SELECT EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = 'public' AND table_name = legacy_table AND column_name = 'atk_bonus'
        ) INTO has_atk_bonus;
        SELECT EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = 'public' AND table_name = legacy_table AND column_name = 'def_bonus'
        ) INTO has_def_bonus;
        SELECT EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = 'public' AND table_name = legacy_table AND column_name = 'speed_bonus'
        ) INTO has_speed_bonus;
        SELECT EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = 'public' AND table_name = legacy_table AND column_name = 'hp_bonus'
        ) INTO has_hp_bonus;
        SELECT EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = 'public' AND table_name = legacy_table AND column_name = 'energy_bonus'
        ) INTO has_energy_bonus;
        SELECT EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = 'public' AND table_name = legacy_table AND column_name = 'allowed_classes'
        ) INTO has_allowed_classes;

        SELECT EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = 'public' AND table_name = legacy_table AND column_name = 'bonus_atk'
        ) INTO has_bonus_atk;
        SELECT EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = 'public' AND table_name = legacy_table AND column_name = 'bonus_def'
        ) INTO has_bonus_def;
        SELECT EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = 'public' AND table_name = legacy_table AND column_name = 'bonus_hp'
        ) INTO has_bonus_hp;
        SELECT EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = 'public' AND table_name = legacy_table AND column_name = 'bonus_energy'
        ) INTO has_bonus_energy;

        insert_sql := format(
            'INSERT INTO public.equipments (item_id, slot, required_level, atk_bonus, def_bonus, speed_bonus, hp_bonus, energy_bonus, allowed_classes) ' ||
            'SELECT item_id, COALESCE(slot, ''''), %s, %s, %s, %s, %s, %s, %s ' ||
            'FROM public.%I WHERE item_id IS NOT NULL ' ||
            'ON CONFLICT (item_id) DO UPDATE SET ' ||
            'slot = EXCLUDED.slot, required_level = EXCLUDED.required_level, atk_bonus = EXCLUDED.atk_bonus, ' ||
            'def_bonus = EXCLUDED.def_bonus, speed_bonus = EXCLUDED.speed_bonus, hp_bonus = EXCLUDED.hp_bonus, ' ||
            'energy_bonus = EXCLUDED.energy_bonus, allowed_classes = EXCLUDED.allowed_classes',
            CASE WHEN has_required_level THEN 'COALESCE(required_level, 1)' ELSE '1' END,
            CASE
                WHEN has_atk_bonus THEN 'COALESCE(atk_bonus, 0)'
                WHEN has_bonus_atk THEN 'COALESCE(bonus_atk, 0)'
                ELSE '0'
            END,
            CASE
                WHEN has_def_bonus THEN 'COALESCE(def_bonus, 0)'
                WHEN has_bonus_def THEN 'COALESCE(bonus_def, 0)'
                ELSE '0'
            END,
            CASE WHEN has_speed_bonus THEN 'COALESCE(speed_bonus, 0)' ELSE '0' END,
            CASE
                WHEN has_hp_bonus THEN 'COALESCE(hp_bonus, 0)'
                WHEN has_bonus_hp THEN 'COALESCE(bonus_hp, 0)'
                ELSE '0'
            END,
            CASE
                WHEN has_energy_bonus THEN 'COALESCE(energy_bonus, 0)'
                WHEN has_bonus_energy THEN 'COALESCE(bonus_energy, 0)'
                ELSE '0'
            END,
            CASE WHEN has_allowed_classes THEN 'COALESCE(allowed_classes, '''')' ELSE '''''' END,
            legacy_table
        );

        EXECUTE insert_sql;
        EXECUTE format('DROP TABLE public.%I', legacy_table);
    END LOOP;
END $$;
