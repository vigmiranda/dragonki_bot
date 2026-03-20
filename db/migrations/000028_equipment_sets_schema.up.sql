CREATE TABLE IF NOT EXISTS public.equipment_sets (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_equipment_sets_name
    ON public.equipment_sets(name);

CREATE TABLE IF NOT EXISTS public.equipment_set_bonuses (
    id BIGSERIAL PRIMARY KEY,
    set_id BIGINT NOT NULL REFERENCES public.equipment_sets(id) ON DELETE CASCADE,
    pieces_required INTEGER NOT NULL,
    atk_bonus INTEGER NOT NULL DEFAULT 0,
    def_bonus INTEGER NOT NULL DEFAULT 0,
    speed_bonus INTEGER NOT NULL DEFAULT 0,
    hp_bonus INTEGER NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_equipment_set_bonuses_set_piece
    ON public.equipment_set_bonuses(set_id, pieces_required);

ALTER TABLE public.items
    ADD COLUMN IF NOT EXISTS set_id BIGINT;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_items_set_id'
    ) THEN
        ALTER TABLE public.items
            ADD CONSTRAINT fk_items_set_id
            FOREIGN KEY (set_id)
            REFERENCES public.equipment_sets(id)
            ON DELETE SET NULL;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_items_set_id
    ON public.items(set_id);
