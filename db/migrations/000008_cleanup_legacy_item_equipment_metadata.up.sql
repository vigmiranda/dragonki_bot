UPDATE public.items
SET
    slot = NULL,
    required_level = 0,
    atk_bonus = 0,
    def_bonus = 0,
    speed_bonus = 0,
    hp_bonus = 0,
    energy_bonus = 0,
    allowed_classes = ''
WHERE type = 'EQUIPMENT';
