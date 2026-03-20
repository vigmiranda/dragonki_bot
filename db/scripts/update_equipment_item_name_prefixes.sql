-- Prefix equipment item names with slot-specific icons.
--
-- Rules:
-- - ARMOR  -> 🛡
-- - BOOTS  -> 🥾
-- - GLOVES -> 🧤
-- - RING   -> 💍
-- - HELMET -> 🪖
-- - AMULET -> 📿
-- - ACCESSORY -> 🎗
-- - WEAPON starting with Bastao/Bastão -> 🪄
-- - WEAPON starting with Lamina/Lâmina/Espada -> 🗡
--
-- Idempotency:
-- - strips any known equipment icon prefix before applying the correct one
-- - updates only rows that actually need a change

WITH normalized AS (
    SELECT
        i.id,
        regexp_replace(i.name, '^(🛡|🥾|🧤|💍|🪖|📿|🎗|🪄|🗡)\s*', '') AS clean_name,
        CASE
            WHEN e.slot = 'ARMOR' THEN '🛡 '
            WHEN e.slot = 'BOOTS' THEN '🥾 '
            WHEN e.slot = 'GLOVES' THEN '🧤 '
            WHEN e.slot = 'RING' THEN '💍 '
            WHEN e.slot = 'HELMET' THEN '🪖 '
            WHEN e.slot = 'AMULET' THEN '📿 '
            WHEN e.slot = 'ACCESSORY' THEN '🎗 '
            WHEN e.slot = 'WEAPON'
                 AND regexp_replace(i.name, '^(🛡|🥾|🧤|💍|🪖|📿|🎗|🪄|🗡)\s*', '') ~* '^(Bastão|Bastao|🦯 Bastão|🦯 Bastao)'
                THEN '🪄 '
            WHEN e.slot = 'WEAPON'
                 AND regexp_replace(i.name, '^(🛡|🥾|🧤|💍|🪖|📿|🎗|🪄|🗡)\s*', '') ~* '^(Lâmina|Lamina|Espada)'
                THEN '🗡 '
            ELSE NULL
        END AS prefix
    FROM public.items i
    JOIN public.equipments e
      ON e.item_id = i.id
    WHERE i.type = 'EQUIPMENT'
)
UPDATE public.items i
   SET name = normalized.prefix || normalized.clean_name
  FROM normalized
 WHERE i.id = normalized.id
   AND normalized.prefix IS NOT NULL
   AND i.name IS DISTINCT FROM normalized.prefix || normalized.clean_name;
