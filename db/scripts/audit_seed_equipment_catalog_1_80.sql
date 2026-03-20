-- Audit queries for db/scripts/seed_equipment_catalog_1_80.sql
--
-- Usage:
-- - Run each query block independently, or execute the whole file and inspect
--   the result sets in sequence.

-- 1. Sets and set bonuses
SELECT
    es.id,
    es.name,
    es.description,
    esb.pieces_required,
    esb.atk_bonus,
    esb.def_bonus,
    esb.speed_bonus,
    esb.hp_bonus
FROM public.equipment_sets es
LEFT JOIN public.equipment_set_bonuses esb
    ON esb.set_id = es.id
ORDER BY es.name, esb.pieces_required;

-- 2. Seeded items with set, slot, stats and allowed classes
SELECT
    i.id,
    i.code,
    i.name,
    i.rarity,
    i.price,
    es.name AS set_name,
    e.slot,
    e.required_level,
    e.atk_bonus,
    e.def_bonus,
    e.speed_bonus,
    e.hp_bonus,
    e.energy_bonus,
    e.allowed_classes
FROM public.items i
JOIN public.equipments e
    ON e.item_id = i.id
LEFT JOIN public.equipment_sets es
    ON es.id = i.set_id
WHERE i.code LIKE 'SET\_%' ESCAPE '\'
ORDER BY es.name, e.required_level, e.slot;

-- 3. Distinct allowed_classes per set
SELECT DISTINCT
    es.name AS set_name,
    e.allowed_classes
FROM public.items i
JOIN public.equipments e
    ON e.item_id = i.id
JOIN public.equipment_sets es
    ON es.id = i.set_id
WHERE i.code LIKE 'SET\_%' ESCAPE '\'
ORDER BY es.name, e.allowed_classes;

-- 4. Total counts expected from the seed
-- Expected:
--   3 sets
--   9 set bonus rows
--   792 items
--   792 equipments
SELECT COUNT(*) AS total_sets
FROM public.equipment_sets
WHERE name IN (
    'Set Tempestade',
    'Set Guardiao da Terra',
    'Set Vulcanico'
);

SELECT COUNT(*) AS total_set_bonus_rows
FROM public.equipment_set_bonuses esb
JOIN public.equipment_sets es
    ON es.id = esb.set_id
WHERE es.name IN (
    'Set Tempestade',
    'Set Guardiao da Terra',
    'Set Vulcanico'
);

SELECT COUNT(*) AS total_items_seed
FROM public.items
WHERE code LIKE 'SET\_%' ESCAPE '\';

SELECT COUNT(*) AS total_equipments_seed
FROM public.equipments e
JOIN public.items i
    ON i.id = e.item_id
WHERE i.code LIKE 'SET\_%' ESCAPE '\';

-- 5. Safety check: no seeded item below level 15
SELECT
    i.code,
    i.name,
    e.required_level
FROM public.items i
JOIN public.equipments e
    ON e.item_id = i.id
WHERE i.code LIKE 'SET\_%' ESCAPE '\'
  AND e.required_level < 15
ORDER BY e.required_level, i.code;

-- 6. Distribution by set and level band
SELECT
    es.name AS set_name,
    MIN(e.required_level) AS min_level,
    MAX(e.required_level) AS max_level,
    COUNT(*) AS total_pieces_seeded
FROM public.items i
JOIN public.equipments e
    ON e.item_id = i.id
JOIN public.equipment_sets es
    ON es.id = i.set_id
WHERE i.code LIKE 'SET\_%' ESCAPE '\'
GROUP BY es.name
ORDER BY es.name;

-- 7. Canonical set-piece names from the strategy doc
SELECT DISTINCT
    es.name AS set_name,
    regexp_replace(i.name, ' Nv\. [0-9]+$', '') AS base_item_name,
    e.slot
FROM public.items i
JOIN public.equipments e
    ON e.item_id = i.id
JOIN public.equipment_sets es
    ON es.id = i.set_id
WHERE i.code LIKE 'SET\_%' ESCAPE '\'
ORDER BY es.name, base_item_name;
