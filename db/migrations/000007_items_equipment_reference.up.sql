CREATE TABLE IF NOT EXISTS public.equipments (
    id BIGSERIAL PRIMARY KEY,
    item_id BIGINT NOT NULL,
    slot VARCHAR(32) NOT NULL DEFAULT '',
    required_level INTEGER NOT NULL DEFAULT 1,
    atk_bonus INTEGER NOT NULL DEFAULT 0,
    def_bonus INTEGER NOT NULL DEFAULT 0,
    speed_bonus INTEGER NOT NULL DEFAULT 0,
    hp_bonus INTEGER NOT NULL DEFAULT 0,
    energy_bonus INTEGER NOT NULL DEFAULT 0,
    allowed_classes TEXT NOT NULL DEFAULT ''
);

ALTER TABLE public.equipments
    ADD COLUMN IF NOT EXISTS required_level INTEGER NOT NULL DEFAULT 1,
    ADD COLUMN IF NOT EXISTS atk_bonus INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS def_bonus INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS speed_bonus INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS hp_bonus INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS energy_bonus INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS allowed_classes TEXT NOT NULL DEFAULT '';

ALTER TABLE public.equipments
    ADD COLUMN IF NOT EXISTS slot VARCHAR(32) NOT NULL DEFAULT '';

CREATE UNIQUE INDEX IF NOT EXISTS idx_equipments_item_id ON public.equipments(item_id);

INSERT INTO public.equipments (
    item_id,
    slot,
    required_level,
    atk_bonus,
    def_bonus,
    speed_bonus,
    hp_bonus,
    energy_bonus,
    allowed_classes
)
SELECT
    i.id,
    COALESCE(i.slot, ''),
    COALESCE(i.required_level, 1),
    COALESCE(i.atk_bonus, 0),
    COALESCE(i.def_bonus, 0),
    COALESCE(i.speed_bonus, 0),
    COALESCE(i.hp_bonus, 0),
    COALESCE(i.energy_bonus, 0),
    COALESCE(i.allowed_classes, '')
FROM public.items i
WHERE i.type = 'EQUIPMENT'
  AND COALESCE(i.slot, '') <> ''
ON CONFLICT (item_id) DO UPDATE
SET
    slot = EXCLUDED.slot,
    required_level = EXCLUDED.required_level,
    atk_bonus = EXCLUDED.atk_bonus,
    def_bonus = EXCLUDED.def_bonus,
    speed_bonus = EXCLUDED.speed_bonus,
    hp_bonus = EXCLUDED.hp_bonus,
    energy_bonus = EXCLUDED.energy_bonus,
    allowed_classes = EXCLUDED.allowed_classes;

ALTER TABLE public.items
    ADD COLUMN IF NOT EXISTS equipment_id BIGINT;

UPDATE public.items i
SET equipment_id = e.id
FROM public.equipments e
WHERE e.item_id = i.id
  AND i.type = 'EQUIPMENT';

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_items_equipment_id'
    ) THEN
        ALTER TABLE public.items
            ADD CONSTRAINT fk_items_equipment_id
            FOREIGN KEY (equipment_id)
            REFERENCES public.equipments(id)
            ON DELETE SET NULL;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_items_equipment_id ON public.items(equipment_id);
