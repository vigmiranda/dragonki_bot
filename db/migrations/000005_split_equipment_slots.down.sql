ALTER TABLE IF EXISTS public.player_equipments
    DROP COLUMN IF EXISTS gloves_inventory_item_id,
    DROP COLUMN IF EXISTS ring_inventory_item_id,
    DROP COLUMN IF EXISTS amulet_inventory_item_id;

UPDATE public.items
SET slot = 'ARMOR'
WHERE slot = 'GLOVES';

UPDATE public.items
SET slot = 'ACCESSORY'
WHERE slot = 'RING';

UPDATE public.items
SET slot = 'HELMET'
WHERE slot = 'AMULET';
