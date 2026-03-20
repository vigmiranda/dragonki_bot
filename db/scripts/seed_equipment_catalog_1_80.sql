-- Equipment catalog seed faithful to:
-- docs/discovery/dragonki_equipment_set_strategy.md
--
-- Scope:
-- - Creates the structures required by the strategy document:
--   public.equipment_sets
--   public.equipment_set_bonuses
--   public.items.set_id
-- - Seeds the three canonical DragonKi sets:
--   Set Tempestade
--   Set Guardiao da Terra
--   Set Vulcanico
-- - Seeds level-based equipment entries from level 15 to 80.
--
-- Important compatibility note:
-- - The current game schema has no BELT slot.
-- - To preserve the exact piece name from the document, "Cinto Vulcanico"
--   is persisted using slot ACCESSORY.
--
-- Idempotency:
-- - equipment sets are upserted by name
-- - set bonuses are upserted by (set_id, pieces_required)
-- - items are upserted by code
-- - allowed_classes follows the strategy document's "Classe recomendada"

CREATE TABLE IF NOT EXISTS public.equipment_sets (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_equipment_sets_name
    ON public.equipment_sets(name);

CREATE TABLE IF NOT EXISTS public.equipment_set_bonuses (
    id BIGSERIAL PRIMARY KEY,
    set_id BIGINT NOT NULL REFERENCES public.equipment_sets(id) ON DELETE CASCADE,
    pieces_required INTEGER NOT NULL,
    atk_bonus INTEGER NOT NULL DEFAULT 0,
    def_bonus INTEGER NOT NULL DEFAULT 0,
    speed_bonus INTEGER NOT NULL DEFAULT 0,
    hp_bonus INTEGER NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_equipment_set_bonuses_set_piece
    ON public.equipment_set_bonuses(set_id, pieces_required);

ALTER TABLE public.items
    ADD COLUMN IF NOT EXISTS set_id BIGINT;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_items_set_id'
    ) THEN
        ALTER TABLE public.items
            ADD CONSTRAINT fk_items_set_id
            FOREIGN KEY (set_id)
            REFERENCES public.equipment_sets(id)
            ON DELETE SET NULL;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_items_set_id
    ON public.items(set_id);

DO $$
DECLARE
    v_items_seq TEXT;
    v_equipments_seq TEXT;
    v_sets_seq TEXT;
    v_set_bonuses_seq TEXT;
BEGIN
    v_items_seq := pg_get_serial_sequence('public.items', 'id');
    IF v_items_seq IS NOT NULL THEN
        EXECUTE format(
            'SELECT setval(%L, COALESCE((SELECT MAX(id) FROM public.items), 0) + 1, false)',
            v_items_seq
        );
    END IF;

    v_equipments_seq := pg_get_serial_sequence('public.equipments', 'id');
    IF v_equipments_seq IS NOT NULL THEN
        EXECUTE format(
            'SELECT setval(%L, COALESCE((SELECT MAX(id) FROM public.equipments), 0) + 1, false)',
            v_equipments_seq
        );
    END IF;

    v_sets_seq := pg_get_serial_sequence('public.equipment_sets', 'id');
    IF v_sets_seq IS NOT NULL THEN
        EXECUTE format(
            'SELECT setval(%L, COALESCE((SELECT MAX(id) FROM public.equipment_sets), 0) + 1, false)',
            v_sets_seq
        );
    END IF;

    v_set_bonuses_seq := pg_get_serial_sequence('public.equipment_set_bonuses', 'id');
    IF v_set_bonuses_seq IS NOT NULL THEN
        EXECUTE format(
            'SELECT setval(%L, COALESCE((SELECT MAX(id) FROM public.equipment_set_bonuses), 0) + 1, false)',
            v_set_bonuses_seq
        );
    END IF;
END $$;

DO $$
DECLARE
    v_level INT;
    v_rarity TEXT;
    v_rarity_bonus INT;
    v_price INT;
    v_allowed_classes TEXT;
    v_item_id BIGINT;
    v_equipment_id BIGINT;
    v_set_id BIGINT;
    cfg RECORD;
    bonus_cfg RECORD;
BEGIN
    FOR cfg IN
        SELECT *
        FROM (
            VALUES
                (
                    'Set Tempestade',
                    'Conjunto de mobilidade e evasao para builds focadas em velocidade.',
                    'Wind Striker',
                    'TEMPESTADE',
                    'HELMET',
                    'Elmo da Tempestade',
                    130,
                    28,
                    0.00::NUMERIC,
                    0.45::NUMERIC,
                    0.20::NUMERIC,
                    1.50::NUMERIC,
                    0.08::NUMERIC,
                    0,
                    2,
                    0,
                    0,
                    0
                ),
                (
                    'Set Tempestade',
                    'Conjunto de mobilidade e evasao para builds focadas em velocidade.',
                    'Wind Striker',
                    'TEMPESTADE',
                    'ARMOR',
                    'Armadura da Tempestade',
                    160,
                    36,
                    0.00::NUMERIC,
                    0.70::NUMERIC,
                    0.25::NUMERIC,
                    2.50::NUMERIC,
                    0.00::NUMERIC,
                    0,
                    3,
                    0,
                    0,
                    0
                ),
                (
                    'Set Tempestade',
                    'Conjunto de mobilidade e evasao para builds focadas em velocidade.',
                    'Wind Striker',
                    'TEMPESTADE',
                    'GLOVES',
                    'Luvas da Tempestade',
                    115,
                    24,
                    0.35::NUMERIC,
                    0.10::NUMERIC,
                    0.30::NUMERIC,
                    0.80::NUMERIC,
                    0.05::NUMERIC,
                    1,
                    0,
                    0,
                    0,
                    0
                ),
                (
                    'Set Tempestade',
                    'Conjunto de mobilidade e evasao para builds focadas em velocidade.',
                    'Wind Striker',
                    'TEMPESTADE',
                    'BOOTS',
                    'Botas da Tempestade',
                    115,
                    24,
                    0.00::NUMERIC,
                    0.15::NUMERIC,
                    0.45::NUMERIC,
                    0.70::NUMERIC,
                    0.12::NUMERIC,
                    0,
                    0,
                    1,
                    0,
                    0
                ),

                (
                    'Set Guardiao da Terra',
                    'Conjunto defensivo para builds tanque com foco em defesa e vida.',
                    'Earth Sentinel',
                    'GUARDIAO',
                    'HELMET',
                    'Elmo do Guardiao',
                    135,
                    28,
                    0.00::NUMERIC,
                    0.65::NUMERIC,
                    0.00::NUMERIC,
                    2.20::NUMERIC,
                    0.08::NUMERIC,
                    0,
                    3,
                    0,
                    0,
                    0
                ),
                (
                    'Set Guardiao da Terra',
                    'Conjunto defensivo para builds tanque com foco em defesa e vida.',
                    'Earth Sentinel',
                    'GUARDIAO',
                    'ARMOR',
                    'Couraca do Guardiao',
                    170,
                    38,
                    0.00::NUMERIC,
                    1.10::NUMERIC,
                    0.00::NUMERIC,
                    4.50::NUMERIC,
                    0.00::NUMERIC,
                    0,
                    5,
                    0,
                    0,
                    0
                ),
                (
                    'Set Guardiao da Terra',
                    'Conjunto defensivo para builds tanque com foco em defesa e vida.',
                    'Earth Sentinel',
                    'GUARDIAO',
                    'GLOVES',
                    'Manoplas do Guardiao',
                    120,
                    24,
                    0.25::NUMERIC,
                    0.55::NUMERIC,
                    0.00::NUMERIC,
                    2.00::NUMERIC,
                    0.05::NUMERIC,
                    1,
                    0,
                    0,
                    0,
                    0
                ),
                (
                    'Set Guardiao da Terra',
                    'Conjunto defensivo para builds tanque com foco em defesa e vida.',
                    'Earth Sentinel',
                    'GUARDIAO',
                    'BOOTS',
                    'Grevas do Guardiao',
                    120,
                    24,
                    0.00::NUMERIC,
                    0.50::NUMERIC,
                    0.08::NUMERIC,
                    1.80::NUMERIC,
                    0.10::NUMERIC,
                    0,
                    0,
                    0,
                    0,
                    0
                ),

                (
                    'Set Vulcanico',
                    'Conjunto agressivo para builds focadas em dano e pressao ofensiva.',
                    'Ki Vanguard',
                    'VULCANICO',
                    'HELMET',
                    'Mascara Vulcanica',
                    135,
                    28,
                    0.30::NUMERIC,
                    0.20::NUMERIC,
                    0.00::NUMERIC,
                    1.50::NUMERIC,
                    0.20::NUMERIC,
                    1,
                    0,
                    0,
                    0,
                    0
                ),
                (
                    'Set Vulcanico',
                    'Conjunto agressivo para builds focadas em dano e pressao ofensiva.',
                    'Ki Vanguard',
                    'VULCANICO',
                    'ARMOR',
                    'Armadura Vulcanica',
                    165,
                    36,
                    0.35::NUMERIC,
                    0.65::NUMERIC,
                    0.00::NUMERIC,
                    3.00::NUMERIC,
                    0.15::NUMERIC,
                    2,
                    2,
                    0,
                    0,
                    0
                ),
                (
                    'Set Vulcanico',
                    'Conjunto agressivo para builds focadas em dano e pressao ofensiva.',
                    'Ki Vanguard',
                    'VULCANICO',
                    'ACCESSORY',
                    'Cinto Vulcanico',
                    125,
                    26,
                    0.55::NUMERIC,
                    0.10::NUMERIC,
                    0.20::NUMERIC,
                    1.00::NUMERIC,
                    0.25::NUMERIC,
                    2,
                    0,
                    0,
                    0,
                    0
                ),
                (
                    'Set Vulcanico',
                    'Conjunto agressivo para builds focadas em dano e pressao ofensiva.',
                    'Ki Vanguard',
                    'VULCANICO',
                    'BOOTS',
                    'Botas Vulcanicas',
                    115,
                    24,
                    0.15::NUMERIC,
                    0.10::NUMERIC,
                    0.30::NUMERIC,
                    0.70::NUMERIC,
                    0.20::NUMERIC,
                    0,
                    0,
                    0,
                    0,
                    0
                )
        ) AS set_cfg(
            set_name,
            set_description,
            recommended_class,
            set_code,
            slot_name,
            piece_name,
            base_price,
            price_step,
            atk_scale,
            def_scale,
            speed_scale,
            hp_scale,
            energy_scale,
            flat_atk,
            flat_def,
            flat_speed,
            flat_hp,
            flat_energy
        )
    LOOP
        INSERT INTO public.equipment_sets (
            name,
            description
        ) VALUES (
            cfg.set_name,
            cfg.set_description
        )
        ON CONFLICT (name) DO UPDATE
        SET description = EXCLUDED.description;
    END LOOP;

    FOR bonus_cfg IN
        SELECT *
        FROM (
            VALUES
                ('Set Tempestade', 2, 0, 0, 4, 0),
                ('Set Tempestade', 3, 0, 0, 6, 0),
                ('Set Tempestade', 4, 4, 0, 8, 0),

                ('Set Guardiao da Terra', 2, 0, 6, 0, 0),
                ('Set Guardiao da Terra', 3, 0, 10, 0, 0),
                ('Set Guardiao da Terra', 4, 0, 12, 0, 5),

                ('Set Vulcanico', 2, 5, 0, 0, 0),
                ('Set Vulcanico', 3, 7, 0, 0, 0),
                ('Set Vulcanico', 4, 10, 0, 3, 0)
        ) AS bonus_values(
            set_name,
            pieces_required,
            atk_bonus,
            def_bonus,
            speed_bonus,
            hp_bonus
        )
    LOOP
        SELECT id
          INTO v_set_id
          FROM public.equipment_sets
         WHERE name = bonus_cfg.set_name
         LIMIT 1;

        INSERT INTO public.equipment_set_bonuses (
            set_id,
            pieces_required,
            atk_bonus,
            def_bonus,
            speed_bonus,
            hp_bonus
        ) VALUES (
            v_set_id,
            bonus_cfg.pieces_required,
            bonus_cfg.atk_bonus,
            bonus_cfg.def_bonus,
            bonus_cfg.speed_bonus,
            bonus_cfg.hp_bonus
        )
        ON CONFLICT (set_id, pieces_required) DO UPDATE
        SET atk_bonus = EXCLUDED.atk_bonus,
            def_bonus = EXCLUDED.def_bonus,
            speed_bonus = EXCLUDED.speed_bonus,
            hp_bonus = EXCLUDED.hp_bonus;
    END LOOP;

    FOR v_level IN 15..80 LOOP
        v_rarity := CASE
            WHEN v_level <= 20 THEN 'COMMON'
            WHEN v_level <= 40 THEN 'UNCOMMON'
            WHEN v_level <= 60 THEN 'RARE'
            ELSE 'EPIC'
        END;

        v_rarity_bonus := CASE v_rarity
            WHEN 'COMMON' THEN 0
            WHEN 'UNCOMMON' THEN 400
            WHEN 'RARE' THEN 1200
            ELSE 3000
        END;

        FOR cfg IN
            SELECT *
            FROM (
                VALUES
                    (
                        'Set Tempestade',
                        'Wind Striker',
                        'TEMPESTADE',
                        'HELMET',
                        'Elmo da Tempestade',
                        130,
                        28,
                        0.00::NUMERIC,
                        0.45::NUMERIC,
                        0.20::NUMERIC,
                        1.50::NUMERIC,
                        0.08::NUMERIC,
                        0,
                        2,
                        0,
                        0,
                        0
                    ),
                    (
                        'Set Tempestade',
                        'Wind Striker',
                        'TEMPESTADE',
                        'ARMOR',
                        'Armadura da Tempestade',
                        160,
                        36,
                        0.00::NUMERIC,
                        0.70::NUMERIC,
                        0.25::NUMERIC,
                        2.50::NUMERIC,
                        0.00::NUMERIC,
                        0,
                        3,
                        0,
                        0,
                        0
                    ),
                    (
                        'Set Tempestade',
                        'Wind Striker',
                        'TEMPESTADE',
                        'GLOVES',
                        'Luvas da Tempestade',
                        115,
                        24,
                        0.35::NUMERIC,
                        0.10::NUMERIC,
                        0.30::NUMERIC,
                        0.80::NUMERIC,
                        0.05::NUMERIC,
                        1,
                        0,
                        0,
                        0,
                        0
                    ),
                    (
                        'Set Tempestade',
                        'Wind Striker',
                        'TEMPESTADE',
                        'BOOTS',
                        'Botas da Tempestade',
                        115,
                        24,
                        0.00::NUMERIC,
                        0.15::NUMERIC,
                        0.45::NUMERIC,
                        0.70::NUMERIC,
                        0.12::NUMERIC,
                        0,
                        0,
                        1,
                        0,
                        0
                    ),

                    (
                        'Set Guardiao da Terra',
                        'Earth Sentinel',
                        'GUARDIAO',
                        'HELMET',
                        'Elmo do Guardiao',
                        135,
                        28,
                        0.00::NUMERIC,
                        0.65::NUMERIC,
                        0.00::NUMERIC,
                        2.20::NUMERIC,
                        0.08::NUMERIC,
                        0,
                        3,
                        0,
                        0,
                        0
                    ),
                    (
                        'Set Guardiao da Terra',
                        'Earth Sentinel',
                        'GUARDIAO',
                        'ARMOR',
                        'Couraca do Guardiao',
                        170,
                        38,
                        0.00::NUMERIC,
                        1.10::NUMERIC,
                        0.00::NUMERIC,
                        4.50::NUMERIC,
                        0.00::NUMERIC,
                        0,
                        5,
                        0,
                        0,
                        0
                    ),
                    (
                        'Set Guardiao da Terra',
                        'Earth Sentinel',
                        'GUARDIAO',
                        'GLOVES',
                        'Manoplas do Guardiao',
                        120,
                        24,
                        0.25::NUMERIC,
                        0.55::NUMERIC,
                        0.00::NUMERIC,
                        2.00::NUMERIC,
                        0.05::NUMERIC,
                        1,
                        0,
                        0,
                        0,
                        0
                    ),
                    (
                        'Set Guardiao da Terra',
                        'Earth Sentinel',
                        'GUARDIAO',
                        'BOOTS',
                        'Grevas do Guardiao',
                        120,
                        24,
                        0.00::NUMERIC,
                        0.50::NUMERIC,
                        0.08::NUMERIC,
                        1.80::NUMERIC,
                        0.10::NUMERIC,
                        0,
                        0,
                        0,
                        0,
                        0
                    ),

                    (
                        'Set Vulcanico',
                        'Ki Vanguard',
                        'VULCANICO',
                        'HELMET',
                        'Mascara Vulcanica',
                        135,
                        28,
                        0.30::NUMERIC,
                        0.20::NUMERIC,
                        0.00::NUMERIC,
                        1.50::NUMERIC,
                        0.20::NUMERIC,
                        1,
                        0,
                        0,
                        0,
                        0
                    ),
                    (
                        'Set Vulcanico',
                        'Ki Vanguard',
                        'VULCANICO',
                        'ARMOR',
                        'Armadura Vulcanica',
                        165,
                        36,
                        0.35::NUMERIC,
                        0.65::NUMERIC,
                        0.00::NUMERIC,
                        3.00::NUMERIC,
                        0.15::NUMERIC,
                        2,
                        2,
                        0,
                        0,
                        0
                    ),
                    (
                        'Set Vulcanico',
                        'Ki Vanguard',
                        'VULCANICO',
                        'ACCESSORY',
                        'Cinto Vulcanico',
                        125,
                        26,
                        0.55::NUMERIC,
                        0.10::NUMERIC,
                        0.20::NUMERIC,
                        1.00::NUMERIC,
                        0.25::NUMERIC,
                        2,
                        0,
                        0,
                        0,
                        0
                    ),
                    (
                        'Set Vulcanico',
                        'Ki Vanguard',
                        'VULCANICO',
                        'BOOTS',
                        'Botas Vulcanicas',
                        115,
                        24,
                        0.15::NUMERIC,
                        0.10::NUMERIC,
                        0.30::NUMERIC,
                        0.70::NUMERIC,
                        0.20::NUMERIC,
                        0,
                        0,
                        0,
                        0,
                        0
                    )
            ) AS item_cfg(
                set_name,
                recommended_class,
                set_code,
                slot_name,
                piece_name,
                base_price,
                price_step,
                atk_scale,
                def_scale,
                speed_scale,
                hp_scale,
                energy_scale,
                flat_atk,
                flat_def,
                flat_speed,
                flat_hp,
                flat_energy
            )
        LOOP
            SELECT id
              INTO v_set_id
              FROM public.equipment_sets
             WHERE name = cfg.set_name
             LIMIT 1;

            v_allowed_classes := CASE cfg.recommended_class
                WHEN 'Ki Vanguard' THEN 'Ki Vanguard,Storm Vanguard,Tempest Avatar'
                WHEN 'Wind Striker' THEN 'Wind Striker,Sky Striker,Hurricane Ascendant'
                WHEN 'Earth Sentinel' THEN 'Earth Sentinel,Terra Guardian,Titan Sentinel'
                ELSE cfg.recommended_class
            END;

            v_price := cfg.base_price + (v_level * cfg.price_step) + v_rarity_bonus;

            SELECT id
              INTO v_item_id
              FROM public.items
             WHERE code = format(
                 'SET_%s_%s_L%s',
                 cfg.set_code,
                 cfg.slot_name,
                 lpad(v_level::TEXT, 2, '0')
             )
             ORDER BY id
             LIMIT 1;

            IF v_item_id IS NULL THEN
                INSERT INTO public.items (
                    code,
                    name,
                    description,
                    type,
                    rarity,
                    amount,
                    price,
                    set_id
                ) VALUES (
                    format(
                        'SET_%s_%s_L%s',
                        cfg.set_code,
                        cfg.slot_name,
                        lpad(v_level::TEXT, 2, '0')
                    ),
                    format('%s Nv. %s', cfg.piece_name, v_level),
                    format(
                        '%s. Parte do %s. Classe recomendada: %s. Nivel %s.',
                        cfg.piece_name,
                        cfg.set_name,
                        cfg.recommended_class,
                        v_level
                    ),
                    'EQUIPMENT',
                    v_rarity,
                    1,
                    v_price,
                    v_set_id
                )
                RETURNING id INTO v_item_id;
            ELSE
                UPDATE public.items
                   SET name = format('%s Nv. %s', cfg.piece_name, v_level),
                       description = format(
                           '%s. Parte do %s. Classe recomendada: %s. Nivel %s.',
                           cfg.piece_name,
                           cfg.set_name,
                           cfg.recommended_class,
                           v_level
                       ),
                       type = 'EQUIPMENT',
                       rarity = v_rarity,
                       amount = 1,
                       price = v_price,
                       set_id = v_set_id
                 WHERE id = v_item_id;
            END IF;

            SELECT id
              INTO v_equipment_id
              FROM public.equipments
             WHERE item_id = v_item_id
             LIMIT 1;

            IF v_equipment_id IS NULL THEN
                INSERT INTO public.equipments (
                    item_id,
                    slot,
                    required_level,
                    atk_bonus,
                    def_bonus,
                    speed_bonus,
                    hp_bonus,
                    energy_bonus,
                    allowed_classes
                ) VALUES (
                    v_item_id,
                    cfg.slot_name,
                    v_level,
                    ceil(v_level::NUMERIC * cfg.atk_scale)::INT + cfg.flat_atk,
                        ceil(v_level::NUMERIC * cfg.def_scale)::INT + cfg.flat_def,
                        ceil(v_level::NUMERIC * cfg.speed_scale)::INT + cfg.flat_speed,
                        ceil(v_level::NUMERIC * cfg.hp_scale)::INT + cfg.flat_hp,
                        ceil(v_level::NUMERIC * cfg.energy_scale)::INT + cfg.flat_energy,
                        v_allowed_classes
                )
                RETURNING id INTO v_equipment_id;
            ELSE
                UPDATE public.equipments
                   SET slot = cfg.slot_name,
                       required_level = v_level,
                       atk_bonus = ceil(v_level::NUMERIC * cfg.atk_scale)::INT + cfg.flat_atk,
                       def_bonus = ceil(v_level::NUMERIC * cfg.def_scale)::INT + cfg.flat_def,
                       speed_bonus = ceil(v_level::NUMERIC * cfg.speed_scale)::INT + cfg.flat_speed,
                       hp_bonus = ceil(v_level::NUMERIC * cfg.hp_scale)::INT + cfg.flat_hp,
                       energy_bonus = ceil(v_level::NUMERIC * cfg.energy_scale)::INT + cfg.flat_energy,
                       allowed_classes = v_allowed_classes
                 WHERE id = v_equipment_id;
            END IF;

            UPDATE public.items
               SET equipment_id = v_equipment_id
             WHERE id = v_item_id;
        END LOOP;
    END LOOP;
END $$;
