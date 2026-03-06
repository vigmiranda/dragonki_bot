ALTER TABLE IF EXISTS public.player_equipments
    ADD COLUMN IF NOT EXISTS weapon_inventory_item_id BIGINT,
    ADD COLUMN IF NOT EXISTS armor_inventory_item_id BIGINT,
    ADD COLUMN IF NOT EXISTS helmet_inventory_item_id BIGINT,
    ADD COLUMN IF NOT EXISTS boots_inventory_item_id BIGINT,
    ADD COLUMN IF NOT EXISTS accessory_inventory_item_id BIGINT;

UPDATE public.player_equipments pe
SET weapon_inventory_item_id = (
    SELECT ii.id
    FROM public.inventory_items ii
    JOIN public.inventories inv ON inv.id = ii.inventory_id
    WHERE inv.player_id = pe.player_id
      AND ii.item_id = pe.weapon_id
    ORDER BY ii.id ASC
    LIMIT 1
)
WHERE pe.weapon_inventory_item_id IS NULL
  AND pe.weapon_id IS NOT NULL
  AND pe.weapon_id > 0;

UPDATE public.player_equipments pe
SET armor_inventory_item_id = (
    SELECT ii.id
    FROM public.inventory_items ii
    JOIN public.inventories inv ON inv.id = ii.inventory_id
    WHERE inv.player_id = pe.player_id
      AND ii.item_id = pe.armor_id
    ORDER BY ii.id ASC
    LIMIT 1
)
WHERE pe.armor_inventory_item_id IS NULL
  AND pe.armor_id IS NOT NULL
  AND pe.armor_id > 0;

UPDATE public.player_equipments pe
SET helmet_inventory_item_id = (
    SELECT ii.id
    FROM public.inventory_items ii
    JOIN public.inventories inv ON inv.id = ii.inventory_id
    WHERE inv.player_id = pe.player_id
      AND ii.item_id = pe.helmet_id
    ORDER BY ii.id ASC
    LIMIT 1
)
WHERE pe.helmet_inventory_item_id IS NULL
  AND pe.helmet_id IS NOT NULL
  AND pe.helmet_id > 0;

UPDATE public.player_equipments pe
SET boots_inventory_item_id = (
    SELECT ii.id
    FROM public.inventory_items ii
    JOIN public.inventories inv ON inv.id = ii.inventory_id
    WHERE inv.player_id = pe.player_id
      AND ii.item_id = pe.boots_id
    ORDER BY ii.id ASC
    LIMIT 1
)
WHERE pe.boots_inventory_item_id IS NULL
  AND pe.boots_id IS NOT NULL
  AND pe.boots_id > 0;

UPDATE public.player_equipments pe
SET accessory_inventory_item_id = (
    SELECT ii.id
    FROM public.inventory_items ii
    JOIN public.inventories inv ON inv.id = ii.inventory_id
    WHERE inv.player_id = pe.player_id
      AND ii.item_id = pe.accessory_id
    ORDER BY ii.id ASC
    LIMIT 1
)
WHERE pe.accessory_inventory_item_id IS NULL
  AND pe.accessory_id IS NOT NULL
  AND pe.accessory_id > 0;

CREATE INDEX IF NOT EXISTS idx_player_equipments_weapon_inventory_item_id
    ON public.player_equipments (weapon_inventory_item_id);
CREATE INDEX IF NOT EXISTS idx_player_equipments_armor_inventory_item_id
    ON public.player_equipments (armor_inventory_item_id);
CREATE INDEX IF NOT EXISTS idx_player_equipments_helmet_inventory_item_id
    ON public.player_equipments (helmet_inventory_item_id);
CREATE INDEX IF NOT EXISTS idx_player_equipments_boots_inventory_item_id
    ON public.player_equipments (boots_inventory_item_id);
CREATE INDEX IF NOT EXISTS idx_player_equipments_accessory_inventory_item_id
    ON public.player_equipments (accessory_inventory_item_id);
