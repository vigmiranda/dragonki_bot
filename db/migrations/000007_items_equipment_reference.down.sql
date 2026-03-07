DROP INDEX IF EXISTS idx_items_equipment_id;

ALTER TABLE public.items
    DROP CONSTRAINT IF EXISTS fk_items_equipment_id;

ALTER TABLE public.items
    DROP COLUMN IF EXISTS equipment_id;
