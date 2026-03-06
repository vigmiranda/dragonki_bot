DROP INDEX IF EXISTS public.idx_player_equipments_weapon_inventory_item_id;
DROP INDEX IF EXISTS public.idx_player_equipments_armor_inventory_item_id;
DROP INDEX IF EXISTS public.idx_player_equipments_helmet_inventory_item_id;
DROP INDEX IF EXISTS public.idx_player_equipments_boots_inventory_item_id;
DROP INDEX IF EXISTS public.idx_player_equipments_accessory_inventory_item_id;

ALTER TABLE IF EXISTS public.player_equipments
    DROP COLUMN IF EXISTS weapon_inventory_item_id,
    DROP COLUMN IF EXISTS armor_inventory_item_id,
    DROP COLUMN IF EXISTS helmet_inventory_item_id,
    DROP COLUMN IF EXISTS boots_inventory_item_id,
    DROP COLUMN IF EXISTS accessory_inventory_item_id;
