-- PASS/FAIL audit for db/scripts/seed_monster_drops_equipment_progression.sql
--
-- Assumes the following seeds were already applied:
-- - db/scripts/seed_shop_low_level_equipments_2_14.sql
-- - db/scripts/seed_equipment_catalog_1_80.sql
-- - db/scripts/seed_monster_drops_equipment_progression.sql

DROP VIEW IF EXISTS tmp_audit_seeded_equipment_drops;
DROP VIEW IF EXISTS tmp_audit_eligible_monsters;
DROP VIEW IF EXISTS tmp_audit_eligible_levels;
DROP VIEW IF EXISTS tmp_audit_actual_levels;

CREATE TEMP VIEW tmp_audit_seeded_equipment_drops AS
SELECT
    md.id,
    md.monster_id,
    md.item_id,
    md.drop_rate,
    md.min_qty,
    md.max_qty,
    md.rarity_weight,
    m.level AS monster_level,
    m.name AS monster_name,
    COALESCE(m.tag, '') AS monster_tag,
    i.code AS item_code,
    i.rarity AS item_rarity,
    COALESCE(i.set_id, 0) AS set_id,
    e.required_level,
    e.slot,
    CASE
        WHEN COALESCE(m.tag, '') = 'AURUM_KARUN' THEN TRUE
        WHEN m.name LIKE '%💰%' THEN TRUE
        ELSE FALSE
    END AS is_money_monster,
    CASE
        WHEN COALESCE(m.tag, '') LIKE '%BOSS%' THEN TRUE
        WHEN m.name LIKE '%👑%' THEN TRUE
        ELSE FALSE
    END AS is_boss
FROM public.monster_drops md
JOIN public.monsters m
    ON m.monster_id = md.monster_id
JOIN public.items i
    ON i.id = md.item_id
JOIN public.equipments e
    ON e.item_id = i.id
WHERE i.code LIKE 'LL\_%' ESCAPE '\'
   OR i.code LIKE 'SET\_%' ESCAPE '\';

CREATE TEMP VIEW tmp_audit_eligible_monsters AS
SELECT
    m.monster_id,
    m.level AS monster_level,
    m.name AS monster_name,
    COALESCE(m.tag, '') AS monster_tag,
    CASE
        WHEN COALESCE(m.tag, '') = 'AURUM_KARUN' THEN TRUE
        WHEN m.name LIKE '%💰%' THEN TRUE
        ELSE FALSE
    END AS is_money_monster,
    CASE
        WHEN COALESCE(m.tag, '') LIKE '%BOSS%' THEN TRUE
        WHEN m.name LIKE '%👑%' THEN TRUE
        ELSE FALSE
    END AS is_boss
FROM public.monsters m;

CREATE TEMP VIEW tmp_audit_eligible_levels AS
SELECT
    em.monster_id,
    COUNT(DISTINCT e.required_level) AS available_distinct_levels
FROM tmp_audit_eligible_monsters em
JOIN public.items i
    ON i.type = 'EQUIPMENT'
   AND (
        i.code LIKE 'LL\_%' ESCAPE '\'
        OR i.code LIKE 'SET\_%' ESCAPE '\'
   )
JOIN public.equipments e
    ON e.item_id = i.id
   AND e.required_level <= em.monster_level
WHERE NOT em.is_money_monster
GROUP BY em.monster_id;

CREATE TEMP VIEW tmp_audit_actual_levels AS
SELECT
    sed.monster_id,
    COUNT(DISTINCT sed.required_level) AS actual_distinct_levels
FROM tmp_audit_seeded_equipment_drops sed
GROUP BY sed.monster_id;

SELECT
    'progression_drop_rows_exist' AS check_name,
    CASE WHEN COUNT(*) > 0 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*)::TEXT AS actual_value,
    '> 0' AS expected_value
FROM tmp_audit_seeded_equipment_drops;

SELECT
    'no_money_monster_equipment_drops' AS check_name,
    CASE WHEN COUNT(*) = 0 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*)::TEXT AS actual_value,
    '0' AS expected_value
FROM tmp_audit_seeded_equipment_drops
WHERE is_money_monster;

SELECT
    'all_seeded_rows_qty_fixed_one' AS check_name,
    CASE WHEN COUNT(*) = 0 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*)::TEXT AS actual_value,
    '0' AS expected_value
FROM tmp_audit_seeded_equipment_drops
WHERE min_qty <> 1
   OR max_qty <> 1;

SELECT
    'all_seeded_rows_match_equipment_level_gate' AS check_name,
    CASE WHEN COUNT(*) = 0 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*)::TEXT AS actual_value,
    '0' AS expected_value
FROM tmp_audit_seeded_equipment_drops
WHERE required_level > monster_level;

SELECT
    'no_duplicate_monster_item_rows' AS check_name,
    CASE WHEN COUNT(*) = 0 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*)::TEXT AS actual_value,
    '0' AS expected_value
FROM (
    SELECT monster_id, item_id
    FROM tmp_audit_seeded_equipment_drops
    GROUP BY monster_id, item_id
    HAVING COUNT(*) > 1
) dup;

SELECT
    'all_non_money_monsters_with_eligible_equipment_received_drops' AS check_name,
    CASE WHEN COUNT(*) = 0 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*)::TEXT AS actual_value,
    '0' AS expected_value
FROM tmp_audit_eligible_levels el
LEFT JOIN tmp_audit_actual_levels al
    ON al.monster_id = el.monster_id
WHERE COALESCE(al.actual_distinct_levels, 0) = 0;

SELECT
    'distinct_level_buckets_per_monster_match_target' AS check_name,
    CASE WHEN COUNT(*) = 0 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*)::TEXT AS actual_value,
    '0' AS expected_value
FROM (
    SELECT
        em.monster_id,
        COALESCE(al.actual_distinct_levels, 0) AS actual_distinct_levels,
        LEAST(
            COALESCE(el.available_distinct_levels, 0),
            CASE
                WHEN em.is_boss THEN 4
                ELSE 3
            END
        ) AS expected_distinct_levels
    FROM tmp_audit_eligible_monsters em
    LEFT JOIN tmp_audit_eligible_levels el
        ON el.monster_id = em.monster_id
    LEFT JOIN tmp_audit_actual_levels al
        ON al.monster_id = em.monster_id
    WHERE NOT em.is_money_monster
      AND COALESCE(el.available_distinct_levels, 0) > 0
) t
WHERE actual_distinct_levels <> expected_distinct_levels;

SELECT
    'higher_level_equipment_always_has_lower_drop_rate' AS check_name,
    CASE WHEN COUNT(*) = 0 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*)::TEXT AS actual_value,
    '0' AS expected_value
FROM (
    SELECT DISTINCT
        higher.monster_id,
        higher.required_level AS higher_required_level,
        lower.required_level AS lower_required_level
    FROM tmp_audit_seeded_equipment_drops higher
    JOIN tmp_audit_seeded_equipment_drops lower
        ON lower.monster_id = higher.monster_id
       AND lower.required_level < higher.required_level
    WHERE higher.drop_rate >= lower.drop_rate
) violations;

SELECT
    'rarity_weight_matches_rarity' AS check_name,
    CASE WHEN COUNT(*) = 0 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*)::TEXT AS actual_value,
    '0' AS expected_value
FROM tmp_audit_seeded_equipment_drops
WHERE rarity_weight <> CASE item_rarity
    WHEN 'COMMON' THEN 12
    WHEN 'UNCOMMON' THEN 8
    WHEN 'RARE' THEN 5
    WHEN 'EPIC' THEN 3
    ELSE 2
END;

SELECT
    'drop_rate_within_expected_bounds' AS check_name,
    CASE WHEN COUNT(*) = 0 THEN 'PASS' ELSE 'FAIL' END AS status,
    COUNT(*)::TEXT AS actual_value,
    '0' AS expected_value
FROM tmp_audit_seeded_equipment_drops
WHERE drop_rate < 0.000060
   OR drop_rate > 0.030000;

DROP VIEW IF EXISTS tmp_audit_actual_levels;
DROP VIEW IF EXISTS tmp_audit_eligible_levels;
DROP VIEW IF EXISTS tmp_audit_eligible_monsters;
DROP VIEW IF EXISTS tmp_audit_seeded_equipment_drops;
