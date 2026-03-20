-- Rollback for db/scripts/seed_shop_low_level_equipments_2_14.sql
--
-- Removes only the generic low-level equipment catalog created with code prefix LL_.
-- It also clears inventory references before deleting the generated items.

DO $$
DECLARE
    v_items_seq TEXT;
    v_equipments_seq TEXT;
BEGIN
    UPDATE public.player_equipments pe
       SET weapon_inventory_item_id = NULL
     WHERE weapon_inventory_item_id IN (
         SELECT ii.id
           FROM public.inventory_items ii
           JOIN public.items i ON i.id = ii.item_id
          WHERE i.code LIKE 'LL\_%' ESCAPE '\'
     );

    UPDATE public.player_equipments pe
       SET armor_inventory_item_id = NULL
     WHERE armor_inventory_item_id IN (
         SELECT ii.id
           FROM public.inventory_items ii
           JOIN public.items i ON i.id = ii.item_id
          WHERE i.code LIKE 'LL\_%' ESCAPE '\'
     );

    UPDATE public.player_equipments pe
       SET gloves_inventory_item_id = NULL
     WHERE gloves_inventory_item_id IN (
         SELECT ii.id
           FROM public.inventory_items ii
           JOIN public.items i ON i.id = ii.item_id
          WHERE i.code LIKE 'LL\_%' ESCAPE '\'
     );

    UPDATE public.player_equipments pe
       SET helmet_inventory_item_id = NULL
     WHERE helmet_inventory_item_id IN (
         SELECT ii.id
           FROM public.inventory_items ii
           JOIN public.items i ON i.id = ii.item_id
          WHERE i.code LIKE 'LL\_%' ESCAPE '\'
     );

    UPDATE public.player_equipments pe
       SET boots_inventory_item_id = NULL
     WHERE boots_inventory_item_id IN (
         SELECT ii.id
           FROM public.inventory_items ii
           JOIN public.items i ON i.id = ii.item_id
          WHERE i.code LIKE 'LL\_%' ESCAPE '\'
     );

    UPDATE public.player_equipments pe
       SET ring_inventory_item_id = NULL
     WHERE ring_inventory_item_id IN (
         SELECT ii.id
           FROM public.inventory_items ii
           JOIN public.items i ON i.id = ii.item_id
          WHERE i.code LIKE 'LL\_%' ESCAPE '\'
     );

    UPDATE public.player_equipments pe
       SET amulet_inventory_item_id = NULL
     WHERE amulet_inventory_item_id IN (
         SELECT ii.id
           FROM public.inventory_items ii
           JOIN public.items i ON i.id = ii.item_id
          WHERE i.code LIKE 'LL\_%' ESCAPE '\'
     );

    UPDATE public.player_equipments pe
       SET accessory_inventory_item_id = NULL
     WHERE accessory_inventory_item_id IN (
         SELECT ii.id
           FROM public.inventory_items ii
           JOIN public.items i ON i.id = ii.item_id
          WHERE i.code LIKE 'LL\_%' ESCAPE '\'
     );

    DELETE FROM public.inventory_items
     WHERE item_id IN (
         SELECT id
           FROM public.items
          WHERE code LIKE 'LL\_%' ESCAPE '\'
     );

    UPDATE public.items
       SET equipment_id = NULL
     WHERE code LIKE 'LL\_%' ESCAPE '\';

    DELETE FROM public.equipments
     WHERE item_id IN (
         SELECT id
           FROM public.items
          WHERE code LIKE 'LL\_%' ESCAPE '\'
     );

    DELETE FROM public.items
     WHERE code LIKE 'LL\_%' ESCAPE '\';

    v_items_seq := pg_get_serial_sequence('public.items', 'id');
    IF v_items_seq IS NOT NULL THEN
        EXECUTE format(
            'SELECT setval(%L, COALESCE((SELECT MAX(id) FROM public.items), 0) + 1, false)',
            v_items_seq
        );
    END IF;

    v_equipments_seq := pg_get_serial_sequence('public.equipments', 'id');
    IF v_equipments_seq IS NOT NULL THEN
        EXECUTE format(
            'SELECT setval(%L, COALESCE((SELECT MAX(id) FROM public.equipments), 0) + 1, false)',
            v_equipments_seq
        );
    END IF;
END $$;
