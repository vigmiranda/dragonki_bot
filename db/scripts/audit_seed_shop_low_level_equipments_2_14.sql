-- Detailed audit for db/scripts/seed_shop_low_level_equipments_2_14.sql

SELECT
    e.slot,
    COUNT(*) AS total_items,
    MIN(e.required_level) AS min_level,
    MAX(e.required_level) AS max_level
FROM public.items i
JOIN public.equipments e
    ON e.item_id = i.id
WHERE i.code LIKE 'LL\_%' ESCAPE '\'
GROUP BY e.slot
ORDER BY e.slot;

SELECT
    i.id,
    i.code,
    i.name,
    i.rarity,
    i.price,
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
WHERE i.code LIKE 'LL\_%' ESCAPE '\'
ORDER BY e.slot, e.required_level, i.code;

SELECT DISTINCT
    e.slot,
    e.allowed_classes
FROM public.items i
JOIN public.equipments e
    ON e.item_id = i.id
WHERE i.code LIKE 'LL\_%' ESCAPE '\'
ORDER BY e.slot, e.allowed_classes;

SELECT
    COUNT(*) AS total_seeded_items
FROM public.items
WHERE code LIKE 'LL\_%' ESCAPE '\';

SELECT
    COUNT(*) AS total_seeded_equipments
FROM public.equipments e
JOIN public.items i
    ON i.id = e.item_id
WHERE i.code LIKE 'LL\_%' ESCAPE '\';

SELECT
    COUNT(*) AS invalid_weapon_count
FROM public.items i
JOIN public.equipments e
    ON e.item_id = i.id
WHERE i.code LIKE 'LL\_%' ESCAPE '\'
  AND e.slot = 'WEAPON'
  AND i.name NOT LIKE '%Espada%'
  AND i.name NOT LIKE '%Bastao%';

SELECT
    COUNT(*) AS invalid_level_count
FROM public.items i
JOIN public.equipments e
    ON e.item_id = i.id
WHERE i.code LIKE 'LL\_%' ESCAPE '\'
  AND (e.required_level < 2 OR e.required_level > 14);
