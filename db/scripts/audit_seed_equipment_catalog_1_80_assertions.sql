-- Assertion-style audit for db/scripts/seed_equipment_catalog_1_80.sql
--
-- Each query returns PASS/FAIL plus diagnostic values.

-- 1. Expected canonical sets exist
SELECT
    'expected_sets_count' AS check_name,
    CASE WHEN COUNT(*) = 3 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*) AS actual_count,
    3 AS expected_count
FROM public.equipment_sets
WHERE name IN (
    'Set Tempestade',
    'Set Guardiao da Terra',
    'Set Vulcanico'
);

-- 2. Expected set bonus rows exist
SELECT
    'expected_set_bonus_rows' AS check_name,
    CASE WHEN COUNT(*) = 9 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*) AS actual_count,
    9 AS expected_count
FROM public.equipment_set_bonuses esb
JOIN public.equipment_sets es
    ON es.id = esb.set_id
WHERE es.name IN (
    'Set Tempestade',
    'Set Guardiao da Terra',
    'Set Vulcanico'
);

-- 3. Expected seeded item count: 3 sets x 4 pieces x 66 levels (15..80) = 792
SELECT
    'expected_seeded_items' AS check_name,
    CASE WHEN COUNT(*) = 792 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*) AS actual_count,
    792 AS expected_count
FROM public.items
WHERE code LIKE 'SET\_%' ESCAPE '\';

-- 4. Expected seeded equipment count = 792
SELECT
    'expected_seeded_equipments' AS check_name,
    CASE WHEN COUNT(*) = 792 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*) AS actual_count,
    792 AS expected_count
FROM public.equipments e
JOIN public.items i
    ON i.id = e.item_id
WHERE i.code LIKE 'SET\_%' ESCAPE '\';

-- 5. No seeded item below level 15
SELECT
    'no_items_below_level_15' AS check_name,
    CASE WHEN COUNT(*) = 0 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*) AS offending_rows,
    0 AS expected_offending_rows
FROM public.items i
JOIN public.equipments e
    ON e.item_id = i.id
WHERE i.code LIKE 'SET\_%' ESCAPE '\'
  AND e.required_level < 15;

-- 6. Every seeded item is linked to a set
SELECT
    'all_seeded_items_have_set_id' AS check_name,
    CASE WHEN COUNT(*) = 0 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*) AS missing_set_rows,
    0 AS expected_missing_rows
FROM public.items i
WHERE i.code LIKE 'SET\_%' ESCAPE '\'
  AND i.set_id IS NULL;

-- 7. Tempestade allowed_classes chain
SELECT
    'tempestade_allowed_classes' AS check_name,
    CASE WHEN COUNT(*) = 1 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*) AS matching_distinct_rows,
    1 AS expected_distinct_rows
FROM (
    SELECT DISTINCT e.allowed_classes
    FROM public.items i
    JOIN public.equipments e
        ON e.item_id = i.id
    JOIN public.equipment_sets es
        ON es.id = i.set_id
    WHERE es.name = 'Set Tempestade'
      AND e.allowed_classes = 'Wind Striker,Sky Striker,Hurricane Ascendant'
) t;

-- 8. Guardiao da Terra allowed_classes chain
SELECT
    'guardiao_allowed_classes' AS check_name,
    CASE WHEN COUNT(*) = 1 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*) AS matching_distinct_rows,
    1 AS expected_distinct_rows
FROM (
    SELECT DISTINCT e.allowed_classes
    FROM public.items i
    JOIN public.equipments e
        ON e.item_id = i.id
    JOIN public.equipment_sets es
        ON es.id = i.set_id
    WHERE es.name = 'Set Guardiao da Terra'
      AND e.allowed_classes = 'Earth Sentinel,Terra Guardian,Titan Sentinel'
) t;

-- 9. Vulcanico allowed_classes chain
SELECT
    'vulcanico_allowed_classes' AS check_name,
    CASE WHEN COUNT(*) = 1 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*) AS matching_distinct_rows,
    1 AS expected_distinct_rows
FROM (
    SELECT DISTINCT e.allowed_classes
    FROM public.items i
    JOIN public.equipments e
        ON e.item_id = i.id
    JOIN public.equipment_sets es
        ON es.id = i.set_id
    WHERE es.name = 'Set Vulcanico'
      AND e.allowed_classes = 'Ki Vanguard,Storm Vanguard,Tempest Avatar'
) t;

-- 10. Distinct canonical pieces per set
SELECT
    'tempestade_piece_names' AS check_name,
    CASE WHEN COUNT(*) = 4 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*) AS actual_count,
    4 AS expected_count
FROM (
    SELECT DISTINCT regexp_replace(i.name, ' Nv\. [0-9]+$', '') AS piece_name
    FROM public.items i
    JOIN public.equipment_sets es
        ON es.id = i.set_id
    WHERE es.name = 'Set Tempestade'
      AND regexp_replace(i.name, ' Nv\. [0-9]+$', '') IN (
          'Elmo da Tempestade',
          'Armadura da Tempestade',
          'Luvas da Tempestade',
          'Botas da Tempestade'
      )
) t;

SELECT
    'guardiao_piece_names' AS check_name,
    CASE WHEN COUNT(*) = 4 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*) AS actual_count,
    4 AS expected_count
FROM (
    SELECT DISTINCT regexp_replace(i.name, ' Nv\. [0-9]+$', '') AS piece_name
    FROM public.items i
    JOIN public.equipment_sets es
        ON es.id = i.set_id
    WHERE es.name = 'Set Guardiao da Terra'
      AND regexp_replace(i.name, ' Nv\. [0-9]+$', '') IN (
          'Elmo do Guardiao',
          'Couraca do Guardiao',
          'Manoplas do Guardiao',
          'Grevas do Guardiao'
      )
) t;

SELECT
    'vulcanico_piece_names' AS check_name,
    CASE WHEN COUNT(*) = 4 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*) AS actual_count,
    4 AS expected_count
FROM (
    SELECT DISTINCT regexp_replace(i.name, ' Nv\. [0-9]+$', '') AS piece_name
    FROM public.items i
    JOIN public.equipment_sets es
        ON es.id = i.set_id
    WHERE es.name = 'Set Vulcanico'
      AND regexp_replace(i.name, ' Nv\. [0-9]+$', '') IN (
          'Mascara Vulcanica',
          'Armadura Vulcanica',
          'Cinto Vulcanico',
          'Botas Vulcanicas'
      )
) t;

-- 11. Exact set bonus definitions
SELECT
    'tempestade_bonus_values' AS check_name,
    CASE WHEN COUNT(*) = 3 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*) AS matching_rows,
    3 AS expected_matching_rows
FROM (
    SELECT esb.pieces_required, esb.atk_bonus, esb.def_bonus, esb.speed_bonus, esb.hp_bonus
    FROM public.equipment_set_bonuses esb
    JOIN public.equipment_sets es
        ON es.id = esb.set_id
    WHERE es.name = 'Set Tempestade'
      AND (
          (esb.pieces_required = 2 AND esb.atk_bonus = 0 AND esb.def_bonus = 0 AND esb.speed_bonus = 4 AND esb.hp_bonus = 0) OR
          (esb.pieces_required = 3 AND esb.atk_bonus = 0 AND esb.def_bonus = 0 AND esb.speed_bonus = 6 AND esb.hp_bonus = 0) OR
          (esb.pieces_required = 4 AND esb.atk_bonus = 4 AND esb.def_bonus = 0 AND esb.speed_bonus = 8 AND esb.hp_bonus = 0)
      )
) t;

SELECT
    'guardiao_bonus_values' AS check_name,
    CASE WHEN COUNT(*) = 3 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*) AS matching_rows,
    3 AS expected_matching_rows
FROM (
    SELECT esb.pieces_required, esb.atk_bonus, esb.def_bonus, esb.speed_bonus, esb.hp_bonus
    FROM public.equipment_set_bonuses esb
    JOIN public.equipment_sets es
        ON es.id = esb.set_id
    WHERE es.name = 'Set Guardiao da Terra'
      AND (
          (esb.pieces_required = 2 AND esb.atk_bonus = 0 AND esb.def_bonus = 6 AND esb.speed_bonus = 0 AND esb.hp_bonus = 0) OR
          (esb.pieces_required = 3 AND esb.atk_bonus = 0 AND esb.def_bonus = 10 AND esb.speed_bonus = 0 AND esb.hp_bonus = 0) OR
          (esb.pieces_required = 4 AND esb.atk_bonus = 0 AND esb.def_bonus = 12 AND esb.speed_bonus = 0 AND esb.hp_bonus = 5)
      )
) t;

SELECT
    'vulcanico_bonus_values' AS check_name,
    CASE WHEN COUNT(*) = 3 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*) AS matching_rows,
    3 AS expected_matching_rows
FROM (
    SELECT esb.pieces_required, esb.atk_bonus, esb.def_bonus, esb.speed_bonus, esb.hp_bonus
    FROM public.equipment_set_bonuses esb
    JOIN public.equipment_sets es
        ON es.id = esb.set_id
    WHERE es.name = 'Set Vulcanico'
      AND (
          (esb.pieces_required = 2 AND esb.atk_bonus = 5 AND esb.def_bonus = 0 AND esb.speed_bonus = 0 AND esb.hp_bonus = 0) OR
          (esb.pieces_required = 3 AND esb.atk_bonus = 7 AND esb.def_bonus = 0 AND esb.speed_bonus = 0 AND esb.hp_bonus = 0) OR
          (esb.pieces_required = 4 AND esb.atk_bonus = 10 AND esb.def_bonus = 0 AND esb.speed_bonus = 3 AND esb.hp_bonus = 0)
      )
) t;
