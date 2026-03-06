INSERT INTO public.items (
    id, code, name, description, type, rarity, amount, price,
    slot, required_level, atk_bonus, def_bonus, speed_bonus, hp_bonus, energy_bonus, allowed_classes
) VALUES
    (60, 'WEAPON_INITIATE', 'Lamina do Iniciante', 'Arma basica para novos guerreiros.', 'EQUIPMENT', 'COMMON', 1, 1500000, 'WEAPON', 1, 7, 1, 0, 0, 0, 'Wind Striker,Earth Sentinel,Ki Vanguard'),
    (61, 'ARMOR_INITIATE', 'Armadura do Iniciante', 'Protecao leve para treinamento.', 'EQUIPMENT', 'COMMON', 1, 1500000, 'ARMOR', 1, 0, 7, 0, 10, 0, 'Wind Striker,Earth Sentinel,Ki Vanguard'),
    (62, 'HELMET_INITIATE', 'Elmo do Iniciante', 'Capacete simples para defesa.', 'EQUIPMENT', 'COMMON', 1, 1500000, 'HELMET', 1, 0, 4, 0, 8, 1, 'Wind Striker,Earth Sentinel,Ki Vanguard'),
    (63, 'BOOTS_INITIATE', 'Botas do Iniciante', 'Botas leves para mobilidade.', 'EQUIPMENT', 'COMMON', 1, 1500000, 'BOOTS', 1, 0, 0, 3, 0, 2, 'Wind Striker,Earth Sentinel,Ki Vanguard'),
    (64, 'ACCESSORY_INITIATE', 'Amuleto do Iniciante', 'Amuleto simples com ki estavel.', 'EQUIPMENT', 'COMMON', 1, 1500000, 'ACCESSORY', 1, 2, 1, 1, 6, 1, 'Wind Striker,Earth Sentinel,Ki Vanguard')
ON CONFLICT (id) DO UPDATE SET
    code = EXCLUDED.code,
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    type = EXCLUDED.type,
    rarity = EXCLUDED.rarity,
    amount = EXCLUDED.amount,
    price = EXCLUDED.price,
    slot = EXCLUDED.slot,
    required_level = EXCLUDED.required_level,
    atk_bonus = EXCLUDED.atk_bonus,
    def_bonus = EXCLUDED.def_bonus,
    speed_bonus = EXCLUDED.speed_bonus,
    hp_bonus = EXCLUDED.hp_bonus,
    energy_bonus = EXCLUDED.energy_bonus,
    allowed_classes = EXCLUDED.allowed_classes;
