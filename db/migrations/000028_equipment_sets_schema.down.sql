DROP INDEX IF EXISTS idx_items_set_id;

ALTER TABLE public.items
    DROP CONSTRAINT IF EXISTS fk_items_set_id;

ALTER TABLE public.items
    DROP COLUMN IF EXISTS set_id;

DROP INDEX IF EXISTS idx_equipment_set_bonuses_set_piece;
DROP TABLE IF EXISTS public.equipment_set_bonuses;

DROP INDEX IF EXISTS idx_equipment_sets_name;
DROP TABLE IF EXISTS public.equipment_sets;
