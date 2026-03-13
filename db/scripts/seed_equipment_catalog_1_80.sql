-- Equipment catalog seed for levels 1..80.
-- Current domain model support:
--   public.items
--   public.equipments
--
-- Notes:
-- - This script does not create equipment sets tables because the current
--   domain model does not persist set bonuses yet.
-- - The naming follows the class fantasy from
--   docs/discovery/dragonki_equipment_set_strategy.md.
-- - Melee weapons are intentionally limited to two models: sword and staff.
-- - Wind Striker and Earth Sentinel use swords.
-- - Ki Vanguard uses staffs.
-- - The script is idempotent by item code.

DO $$
DECLARE
    v_class_code TEXT;
    v_class_name TEXT;
    v_family_label TEXT;
    v_weapon_label TEXT;
    v_level INT;
    v_slot TEXT;
    v_rarity TEXT;
    v_rarity_bonus INT;
    v_base_price INT;
    v_price_step INT;
    v_price INT;
    v_item_code TEXT;
    v_item_name TEXT;
    v_item_desc TEXT;
    v_allowed_classes TEXT;
    v_atk_bonus INT;
    v_def_bonus INT;
    v_speed_bonus INT;
    v_hp_bonus INT;
    v_energy_bonus INT;
    v_item_id BIGINT;
    v_equipment_id BIGINT;
BEGIN
    FOR v_class_code, v_class_name, v_family_label, v_weapon_label IN
        SELECT *
        FROM (
            VALUES
                ('WIND',  'Wind Striker',   'Tempestade',          'Espada'),
                ('EARTH', 'Earth Sentinel', 'Guardiao da Terra',   'Espada'),
                ('KI',    'Ki Vanguard',    'Vulcanico',           'Bastao')
        ) AS class_cfg(class_code, class_name, family_label, weapon_label)
    LOOP
        v_allowed_classes := v_class_name;

        FOR v_level IN 1..80 LOOP
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

            FOR v_slot IN
                SELECT *
                FROM (
                    VALUES
                        ('WEAPON'),
                        ('HELMET'),
                        ('ARMOR'),
                        ('GLOVES'),
                        ('BOOTS'),
                        ('RING'),
                        ('AMULET')
                ) AS slot_cfg(slot_name)
            LOOP
                v_item_code := format(
                    'EQ_%s_%s_L%s',
                    v_class_code,
                    v_slot,
                    lpad(v_level::TEXT, 2, '0')
                );

                v_item_name := CASE
                    WHEN v_class_code = 'WIND' AND v_slot = 'WEAPON' THEN format('%s da Tempestade Nv. %s', v_weapon_label, v_level)
                    WHEN v_class_code = 'WIND' AND v_slot = 'HELMET' THEN format('Elmo da Tempestade Nv. %s', v_level)
                    WHEN v_class_code = 'WIND' AND v_slot = 'ARMOR' THEN format('Armadura da Tempestade Nv. %s', v_level)
                    WHEN v_class_code = 'WIND' AND v_slot = 'GLOVES' THEN format('Luvas da Tempestade Nv. %s', v_level)
                    WHEN v_class_code = 'WIND' AND v_slot = 'BOOTS' THEN format('Botas da Tempestade Nv. %s', v_level)
                    WHEN v_class_code = 'WIND' AND v_slot = 'RING' THEN format('Anel da Tempestade Nv. %s', v_level)
                    WHEN v_class_code = 'WIND' AND v_slot = 'AMULET' THEN format('Amuleto da Tempestade Nv. %s', v_level)

                    WHEN v_class_code = 'EARTH' AND v_slot = 'WEAPON' THEN format('%s do Guardiao Nv. %s', v_weapon_label, v_level)
                    WHEN v_class_code = 'EARTH' AND v_slot = 'HELMET' THEN format('Elmo do Guardiao Nv. %s', v_level)
                    WHEN v_class_code = 'EARTH' AND v_slot = 'ARMOR' THEN format('Couraca do Guardiao Nv. %s', v_level)
                    WHEN v_class_code = 'EARTH' AND v_slot = 'GLOVES' THEN format('Manoplas do Guardiao Nv. %s', v_level)
                    WHEN v_class_code = 'EARTH' AND v_slot = 'BOOTS' THEN format('Grevas do Guardiao Nv. %s', v_level)
                    WHEN v_class_code = 'EARTH' AND v_slot = 'RING' THEN format('Anel do Guardiao Nv. %s', v_level)
                    ELSE format('Amuleto do Guardiao Nv. %s', v_level)
                END;

                IF v_class_code = 'KI' THEN
                    v_item_name := CASE v_slot
                        WHEN 'WEAPON' THEN format('%s Vulcanico Nv. %s', v_weapon_label, v_level)
                        WHEN 'HELMET' THEN format('Elmo Vulcanico Nv. %s', v_level)
                        WHEN 'ARMOR' THEN format('Armadura Vulcanica Nv. %s', v_level)
                        WHEN 'GLOVES' THEN format('Manoplas Vulcanicas Nv. %s', v_level)
                        WHEN 'BOOTS' THEN format('Botas Vulcanicas Nv. %s', v_level)
                        WHEN 'RING' THEN format('Anel Vulcanico Nv. %s', v_level)
                        ELSE format('Amuleto Vulcanico Nv. %s', v_level)
                    END;
                END IF;

                v_item_desc := format(
                    'Equipamento %s para %s no nivel %s.',
                    lower(v_family_label),
                    v_class_name,
                    v_level
                );

                v_base_price := CASE v_slot
                    WHEN 'WEAPON' THEN 180
                    WHEN 'ARMOR' THEN 160
                    WHEN 'HELMET' THEN 130
                    WHEN 'GLOVES' THEN 110
                    WHEN 'BOOTS' THEN 110
                    WHEN 'RING' THEN 125
                    ELSE 125
                END;

                v_price_step := CASE v_slot
                    WHEN 'WEAPON' THEN 45
                    WHEN 'ARMOR' THEN 38
                    WHEN 'HELMET' THEN 28
                    WHEN 'GLOVES' THEN 24
                    WHEN 'BOOTS' THEN 24
                    WHEN 'RING' THEN 26
                    ELSE 26
                END;

                v_price := v_base_price + (v_level * v_price_step) + v_rarity_bonus;

                v_atk_bonus := 0;
                v_def_bonus := 0;
                v_speed_bonus := 0;
                v_hp_bonus := 0;
                v_energy_bonus := 0;

                IF v_class_code = 'WIND' THEN
                    CASE v_slot
                        WHEN 'WEAPON' THEN
                            v_atk_bonus := ceil(v_level::NUMERIC * 1.40)::INT + 5;
                            v_def_bonus := ceil(v_level::NUMERIC * 0.15)::INT;
                            v_speed_bonus := ceil(v_level::NUMERIC * 0.35)::INT;
                            v_hp_bonus := v_level * 2;
                        WHEN 'HELMET' THEN
                            v_def_bonus := ceil(v_level::NUMERIC * 0.55)::INT + 2;
                            v_speed_bonus := ceil(v_level::NUMERIC * 0.10)::INT;
                            v_hp_bonus := v_level * 2;
                            v_energy_bonus := ceil(v_level::NUMERIC * 0.08)::INT;
                        WHEN 'ARMOR' THEN
                            v_def_bonus := ceil(v_level::NUMERIC * 0.90)::INT + 4;
                            v_speed_bonus := ceil(v_level::NUMERIC * 0.15)::INT;
                            v_hp_bonus := v_level * 4;
                        WHEN 'GLOVES' THEN
                            v_atk_bonus := ceil(v_level::NUMERIC * 0.50)::INT + 1;
                            v_speed_bonus := ceil(v_level::NUMERIC * 0.25)::INT;
                            v_hp_bonus := v_level;
                        WHEN 'BOOTS' THEN
                            v_def_bonus := ceil(v_level::NUMERIC * 0.20)::INT;
                            v_speed_bonus := ceil(v_level::NUMERIC * 0.45)::INT + 1;
                            v_energy_bonus := ceil(v_level::NUMERIC * 0.10)::INT;
                        WHEN 'RING' THEN
                            v_atk_bonus := ceil(v_level::NUMERIC * 0.45)::INT;
                            v_speed_bonus := ceil(v_level::NUMERIC * 0.20)::INT;
                            v_hp_bonus := v_level;
                            v_energy_bonus := ceil(v_level::NUMERIC * 0.15)::INT;
                        WHEN 'AMULET' THEN
                            v_def_bonus := ceil(v_level::NUMERIC * 0.40)::INT;
                            v_speed_bonus := ceil(v_level::NUMERIC * 0.15)::INT;
                            v_hp_bonus := v_level * 2;
                            v_energy_bonus := ceil(v_level::NUMERIC * 0.20)::INT;
                    END CASE;
                ELSIF v_class_code = 'EARTH' THEN
                    CASE v_slot
                        WHEN 'WEAPON' THEN
                            v_atk_bonus := ceil(v_level::NUMERIC * 1.25)::INT + 5;
                            v_def_bonus := ceil(v_level::NUMERIC * 0.30)::INT;
                            v_hp_bonus := v_level * 3;
                        WHEN 'HELMET' THEN
                            v_def_bonus := ceil(v_level::NUMERIC * 0.80)::INT + 3;
                            v_hp_bonus := v_level * 3;
                            v_energy_bonus := ceil(v_level::NUMERIC * 0.08)::INT;
                        WHEN 'ARMOR' THEN
                            v_def_bonus := ceil(v_level::NUMERIC * 1.20)::INT + 5;
                            v_hp_bonus := v_level * 5;
                        WHEN 'GLOVES' THEN
                            v_atk_bonus := ceil(v_level::NUMERIC * 0.40)::INT + 1;
                            v_def_bonus := ceil(v_level::NUMERIC * 0.45)::INT;
                            v_hp_bonus := v_level * 2;
                        WHEN 'BOOTS' THEN
                            v_def_bonus := ceil(v_level::NUMERIC * 0.50)::INT;
                            v_speed_bonus := ceil(v_level::NUMERIC * 0.10)::INT;
                            v_hp_bonus := v_level * 2;
                            v_energy_bonus := ceil(v_level::NUMERIC * 0.10)::INT;
                        WHEN 'RING' THEN
                            v_atk_bonus := ceil(v_level::NUMERIC * 0.30)::INT;
                            v_def_bonus := ceil(v_level::NUMERIC * 0.35)::INT;
                            v_hp_bonus := v_level * 2;
                            v_energy_bonus := ceil(v_level::NUMERIC * 0.15)::INT;
                        WHEN 'AMULET' THEN
                            v_def_bonus := ceil(v_level::NUMERIC * 0.45)::INT;
                            v_hp_bonus := v_level * 3;
                            v_energy_bonus := ceil(v_level::NUMERIC * 0.25)::INT;
                    END CASE;
                ELSE
                    CASE v_slot
                        WHEN 'WEAPON' THEN
                            v_atk_bonus := ceil(v_level::NUMERIC * 1.50)::INT + 5;
                            v_speed_bonus := ceil(v_level::NUMERIC * 0.20)::INT;
                            v_energy_bonus := ceil(v_level::NUMERIC * 0.30)::INT;
                        WHEN 'HELMET' THEN
                            v_atk_bonus := ceil(v_level::NUMERIC * 0.15)::INT;
                            v_def_bonus := ceil(v_level::NUMERIC * 0.45)::INT + 1;
                            v_hp_bonus := v_level * 2;
                            v_energy_bonus := ceil(v_level::NUMERIC * 0.25)::INT;
                        WHEN 'ARMOR' THEN
                            v_atk_bonus := ceil(v_level::NUMERIC * 0.20)::INT;
                            v_def_bonus := ceil(v_level::NUMERIC * 0.75)::INT + 3;
                            v_hp_bonus := v_level * 3;
                            v_energy_bonus := ceil(v_level::NUMERIC * 0.25)::INT;
                        WHEN 'GLOVES' THEN
                            v_atk_bonus := ceil(v_level::NUMERIC * 0.60)::INT + 2;
                            v_speed_bonus := ceil(v_level::NUMERIC * 0.15)::INT;
                            v_energy_bonus := ceil(v_level::NUMERIC * 0.15)::INT;
                        WHEN 'BOOTS' THEN
                            v_def_bonus := ceil(v_level::NUMERIC * 0.15)::INT;
                            v_speed_bonus := ceil(v_level::NUMERIC * 0.25)::INT;
                            v_energy_bonus := ceil(v_level::NUMERIC * 0.20)::INT;
                        WHEN 'RING' THEN
                            v_atk_bonus := ceil(v_level::NUMERIC * 0.50)::INT;
                            v_speed_bonus := ceil(v_level::NUMERIC * 0.15)::INT;
                            v_hp_bonus := v_level;
                            v_energy_bonus := ceil(v_level::NUMERIC * 0.25)::INT;
                        WHEN 'AMULET' THEN
                            v_atk_bonus := ceil(v_level::NUMERIC * 0.25)::INT;
                            v_def_bonus := ceil(v_level::NUMERIC * 0.20)::INT;
                            v_hp_bonus := v_level;
                            v_energy_bonus := ceil(v_level::NUMERIC * 0.40)::INT;
                    END CASE;
                END IF;

                v_item_id := NULL;
                SELECT id
                  INTO v_item_id
                  FROM public.items
                 WHERE code = v_item_code
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
                        price
                    ) VALUES (
                        v_item_code,
                        v_item_name,
                        v_item_desc,
                        'EQUIPMENT',
                        v_rarity,
                        1,
                        v_price
                    )
                    RETURNING id INTO v_item_id;
                ELSE
                    UPDATE public.items
                       SET name = v_item_name,
                           description = v_item_desc,
                           type = 'EQUIPMENT',
                           rarity = v_rarity,
                           amount = 1,
                           price = v_price
                     WHERE id = v_item_id;
                END IF;

                v_equipment_id := NULL;
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
                        v_slot,
                        v_level,
                        v_atk_bonus,
                        v_def_bonus,
                        v_speed_bonus,
                        v_hp_bonus,
                        v_energy_bonus,
                        v_allowed_classes
                    )
                    RETURNING id INTO v_equipment_id;
                ELSE
                    UPDATE public.equipments
                       SET slot = v_slot,
                           required_level = v_level,
                           atk_bonus = v_atk_bonus,
                           def_bonus = v_def_bonus,
                           speed_bonus = v_speed_bonus,
                           hp_bonus = v_hp_bonus,
                           energy_bonus = v_energy_bonus,
                           allowed_classes = v_allowed_classes
                     WHERE id = v_equipment_id;
                END IF;

                UPDATE public.items
                   SET equipment_id = v_equipment_id
                 WHERE id = v_item_id;
            END LOOP;
        END LOOP;
    END LOOP;
END $$;
