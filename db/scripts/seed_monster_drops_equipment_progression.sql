-- Equipment progression drops seed aligned with:
-- - db/scripts/seed_shop_low_level_equipments_2_14.sql
-- - db/scripts/seed_equipment_catalog_1_80.sql
--
-- Strategy:
-- - Uses only equipment families seeded by the scripts above: LL_* and SET_*.
-- - Excludes money-only monsters such as Aurum Karun.
-- - For each monster, keeps the nearest 3 equipment levels at or below the monster level.
--   Boss monsters keep the nearest 4 levels.
-- - Within the same monster, higher equipment level always receives a lower drop chance.
-- - Bosses receive a modest multiplier, but the progression ordering is preserved.
--
-- Idempotency:
-- - Removes only monster_drops rows owned by this seed (LL_* and SET_* items).
-- - Rebuilds the progression rows from the current monsters/items/equipments state.

DROP TABLE IF EXISTS tmp_seed_monster_drops_equipment_progression;

WITH seq AS (
    SELECT
        pg_get_serial_sequence('public.monster_drops', 'id') AS seq_name,
        COALESCE((SELECT MAX(id) FROM public.monster_drops), 0) + 1 AS next_id
)
SELECT setval(seq_name::regclass, next_id, false)
FROM seq
WHERE seq_name IS NOT NULL;

CREATE TEMP TABLE tmp_seed_monster_drops_equipment_progression (
    monster_id INTEGER NOT NULL,
    item_id BIGINT NOT NULL,
    drop_rate NUMERIC(10, 6) NOT NULL,
    min_qty INTEGER NOT NULL,
    max_qty INTEGER NOT NULL,
    rarity_weight INTEGER NOT NULL
);

WITH monster_pool AS (
    SELECT
        m.monster_id,
        m.level AS monster_level,
        m.name,
        COALESCE(m.tag, '') AS tag,
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
    FROM public.monsters m
),
equipment_pool AS (
    SELECT
        i.id AS item_id,
        i.code,
        i.rarity,
        COALESCE(i.set_id, 0) AS set_id,
        e.required_level,
        e.slot
    FROM public.items i
    JOIN public.equipments e
      ON e.item_id = i.id
    WHERE i.type = 'EQUIPMENT'
      AND (
          i.code LIKE 'LL\_%' ESCAPE '\'
          OR i.code LIKE 'SET\_%' ESCAPE '\'
      )
),
ranked_candidates AS (
    SELECT
        m.monster_id,
        m.monster_level,
        m.is_boss,
        e.item_id,
        e.code,
        e.rarity,
        e.set_id,
        e.required_level,
        e.slot,
        DENSE_RANK() OVER (
            PARTITION BY m.monster_id
            ORDER BY e.required_level DESC
        ) AS level_rank
    FROM monster_pool m
    JOIN equipment_pool e
      ON e.required_level <= m.monster_level
    WHERE NOT m.is_money_monster
),
selected_candidates AS (
    SELECT
        rc.monster_id,
        rc.monster_level,
        rc.is_boss,
        rc.item_id,
        rc.code,
        rc.rarity,
        rc.set_id,
        rc.required_level,
        rc.slot,
        rc.level_rank,
        CASE rc.rarity
            WHEN 'COMMON' THEN 12
            WHEN 'UNCOMMON' THEN 8
            WHEN 'RARE' THEN 5
            WHEN 'EPIC' THEN 3
            ELSE 2
        END AS rarity_weight,
        ROUND(
            LEAST(
                0.030000::NUMERIC,
                GREATEST(
                    0.000060::NUMERIC,
                    (
                        CASE
                            WHEN rc.monster_level <= 14 THEN 0.001800
                            WHEN rc.monster_level <= 32 THEN 0.001100
                            WHEN rc.monster_level <= 45 THEN 0.000850
                            ELSE 0.000700
                        END
                        * CASE rc.level_rank
                            WHEN 1 THEN 0.55
                            WHEN 2 THEN 0.80
                            WHEN 3 THEN 1.05
                            ELSE 1.30
                          END
                        * CASE rc.rarity
                            WHEN 'COMMON' THEN 1.00
                            WHEN 'UNCOMMON' THEN 0.82
                            WHEN 'RARE' THEN 0.64
                            WHEN 'EPIC' THEN 0.50
                            ELSE 0.40
                          END
                        * CASE
                            WHEN rc.is_boss THEN 1.45
                            ELSE 1.00
                          END
                        * CASE
                            WHEN rc.set_id > 0 THEN 0.90
                            ELSE 1.00
                          END
                        * CASE rc.slot
                            WHEN 'WEAPON' THEN 0.90
                            WHEN 'ARMOR' THEN 0.92
                            WHEN 'ACCESSORY' THEN 0.95
                            WHEN 'RING' THEN 1.05
                            WHEN 'AMULET' THEN 1.05
                            ELSE 1.00
                          END
                    )::NUMERIC
                )
            ),
            6
        ) AS drop_rate
    FROM ranked_candidates rc
    WHERE rc.level_rank <= CASE
        WHEN rc.is_boss THEN 4
        ELSE 3
    END
)
INSERT INTO tmp_seed_monster_drops_equipment_progression (
    monster_id,
    item_id,
    drop_rate,
    min_qty,
    max_qty,
    rarity_weight
)
SELECT
    sc.monster_id,
    sc.item_id,
    sc.drop_rate,
    1,
    1,
    sc.rarity_weight
FROM selected_candidates sc
ORDER BY sc.monster_id, sc.required_level DESC, sc.item_id;

DELETE FROM public.monster_drops md
USING public.items i
WHERE md.item_id = i.id
  AND (
      i.code LIKE 'LL\_%' ESCAPE '\'
      OR i.code LIKE 'SET\_%' ESCAPE '\'
  );

INSERT INTO public.monster_drops (
    monster_id,
    item_id,
    drop_rate,
    min_qty,
    max_qty,
    rarity_weight
)
SELECT
    t.monster_id,
    t.item_id,
    t.drop_rate,
    t.min_qty,
    t.max_qty,
    t.rarity_weight
FROM tmp_seed_monster_drops_equipment_progression t
ORDER BY t.monster_id, t.item_id;

DROP TABLE IF EXISTS tmp_seed_monster_drops_equipment_progression;
