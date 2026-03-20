-- PASS/FAIL audit for db/scripts/seed_shop_low_level_equipments_2_14.sql

SELECT
    'low_level_seed_total_items' AS check_name,
    CASE WHEN COUNT(*) = 35 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*)::TEXT AS actual_value,
    '35' AS expected_value
FROM public.items
WHERE code LIKE 'LL\_%' ESCAPE '\';

SELECT
    'low_level_seed_total_equipments' AS check_name,
    CASE WHEN COUNT(*) = 35 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*)::TEXT AS actual_value,
    '35' AS expected_value
FROM public.equipments e
JOIN public.items i
    ON i.id = e.item_id
WHERE i.code LIKE 'LL\_%' ESCAPE '\';

SELECT
    'low_level_seed_valid_levels' AS check_name,
    CASE WHEN COUNT(*) = 0 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*)::TEXT AS actual_value,
    '0' AS expected_value
FROM public.items i
JOIN public.equipments e
    ON e.item_id = i.id
WHERE i.code LIKE 'LL\_%' ESCAPE '\'
  AND (e.required_level < 2 OR e.required_level > 14);

SELECT
    'low_level_seed_weapon_types' AS check_name,
    CASE WHEN COUNT(*) = 0 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*)::TEXT AS actual_value,
    '0' AS expected_value
FROM public.items i
JOIN public.equipments e
    ON e.item_id = i.id
WHERE i.code LIKE 'LL\_%' ESCAPE '\'
  AND e.slot = 'WEAPON'
  AND i.name NOT LIKE '%Espada%'
  AND i.name NOT LIKE '%Bastao%';

SELECT
    'low_level_seed_allowed_classes' AS check_name,
    CASE WHEN COUNT(*) = 35 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*)::TEXT AS actual_value,
    '35' AS expected_value
FROM public.items i
JOIN public.equipments e
    ON e.item_id = i.id
WHERE i.code LIKE 'LL\_%' ESCAPE '\'
  AND e.allowed_classes = '';

SELECT
    'low_level_seed_slots_have_5_items' AS check_name,
    CASE WHEN COUNT(*) = 7 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*)::TEXT AS actual_value,
    '7' AS expected_value
FROM (
    SELECT e.slot
    FROM public.items i
    JOIN public.equipments e
        ON e.item_id = i.id
    WHERE i.code LIKE 'LL\_%' ESCAPE '\'
    GROUP BY e.slot
    HAVING COUNT(*) >= 5
) valid_slots;
