ALTER TABLE IF EXISTS public.player_equipments
    ADD COLUMN IF NOT EXISTS gloves_inventory_item_id BIGINT,
    ADD COLUMN IF NOT EXISTS ring_inventory_item_id BIGINT,
    ADD COLUMN IF NOT EXISTS amulet_inventory_item_id BIGINT;

UPDATE public.player_equipments pe
SET gloves_inventory_item_id = pe.armor_inventory_item_id
WHERE pe.gloves_inventory_item_id IS NULL
  AND pe.armor_inventory_item_id IS NOT NULL
  AND EXISTS (
    SELECT 1
    FROM public.inventory_items ii
    JOIN public.items it ON it.id = ii.item_id
    WHERE ii.id = pe.armor_inventory_item_id
      AND (it.code LIKE 'GLOVES_%' OR it.slot = 'ARMOR')
  );

UPDATE public.player_equipments pe
SET ring_inventory_item_id = pe.accessory_inventory_item_id
WHERE pe.ring_inventory_item_id IS NULL
  AND pe.accessory_inventory_item_id IS NOT NULL
  AND EXISTS (
    SELECT 1
    FROM public.inventory_items ii
    JOIN public.items it ON it.id = ii.item_id
    WHERE ii.id = pe.accessory_inventory_item_id
      AND (it.code LIKE 'RING_%' OR it.slot = 'ACCESSORY')
  );

UPDATE public.player_equipments pe
SET amulet_inventory_item_id = pe.helmet_inventory_item_id
WHERE pe.amulet_inventory_item_id IS NULL
  AND pe.helmet_inventory_item_id IS NOT NULL
  AND EXISTS (
    SELECT 1
    FROM public.inventory_items ii
    JOIN public.items it ON it.id = ii.item_id
    WHERE ii.id = pe.helmet_inventory_item_id
      AND it.code LIKE 'AMULET_%'
  );

UPDATE public.player_equipments
SET armor_inventory_item_id = NULL
WHERE gloves_inventory_item_id IS NOT NULL
  AND armor_inventory_item_id = gloves_inventory_item_id;

UPDATE public.player_equipments
SET accessory_inventory_item_id = NULL
WHERE ring_inventory_item_id IS NOT NULL
  AND accessory_inventory_item_id = ring_inventory_item_id;

UPDATE public.player_equipments
SET helmet_inventory_item_id = NULL
WHERE amulet_inventory_item_id IS NOT NULL
  AND helmet_inventory_item_id = amulet_inventory_item_id;

UPDATE public.items
SET slot = 'GLOVES'
WHERE code LIKE 'GLOVES_%'
   OR (slot = 'ARMOR' AND code LIKE 'GLOVES%');

UPDATE public.items
SET slot = 'RING'
WHERE code LIKE 'RING_%'
   OR (slot = 'ACCESSORY' AND code LIKE 'RING%');

UPDATE public.items
SET slot = 'AMULET'
WHERE code LIKE 'AMULET_%'
   OR (slot = 'HELMET' AND code LIKE 'AMULET%');
