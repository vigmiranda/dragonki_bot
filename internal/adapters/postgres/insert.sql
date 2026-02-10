INSERT INTO public.classes (id, name, base_class, tier, min_level, atk_bonus, hp_bonus, energy_bonus) VALUES (1, 'Blade Monk', 'Blade Monk', 'BASE', 1, 0, 0, 0);
INSERT INTO public.classes (id, name, base_class, tier, min_level, atk_bonus, hp_bonus, energy_bonus) VALUES (2, 'Sword Disciple', 'Blade Monk', 'ADVANCED', 15, 5, 0, 0);
INSERT INTO public.classes (id, name, base_class, tier, min_level, atk_bonus, hp_bonus, energy_bonus) VALUES (3, 'Spirit Guard', 'Blade Monk', 'ADVANCED', 15, 0, 10, 0);
INSERT INTO public.classes (id, name, base_class, tier, min_level, atk_bonus, hp_bonus, energy_bonus) VALUES (4, 'Ki Adept', 'Blade Monk', 'ADVANCED', 15, 0, 0, 5);
INSERT INTO public.classes (id, name, base_class, tier, min_level, atk_bonus, hp_bonus, energy_bonus) VALUES (5, 'Blade Master', 'Sword Disciple', 'ELITE', 35, 10, 0, 0);
INSERT INTO public.classes (id, name, base_class, tier, min_level, atk_bonus, hp_bonus, energy_bonus) VALUES (6, 'Soul Protector', 'Spirit Guard', 'ELITE', 35, 0, 20, 0);
INSERT INTO public.classes (id, name, base_class, tier, min_level, atk_bonus, hp_bonus, energy_bonus) VALUES (7, 'Ki Channeler', 'Ki Adept', 'ELITE', 35, 0, 0, 15);
INSERT INTO public.classes (id, name, base_class, tier, min_level, atk_bonus, hp_bonus, energy_bonus) VALUES (8, 'Dragon Slayer', 'Blade Master', 'MASTER', 50, 20, 0, 0);
INSERT INTO public.classes (id, name, base_class, tier, min_level, atk_bonus, hp_bonus, energy_bonus) VALUES (9, 'Eternal Sentinel', 'Soul Protector', 'MASTER', 50, 0, 35, 0);
INSERT INTO public.classes (id, name, base_class, tier, min_level, atk_bonus, hp_bonus, energy_bonus) VALUES (10, 'Ascended Sage', 'Ki Channeler', 'MASTER', 50, 0, 0, 25);
INSERT INTO public.monsters (monster_id, name, tag, level, hp, atk, def, max_xp, max_money, spawn_weight, region_id) VALUES (1, 'Slime Verde', 'SLIME_VERDE', 1, 100, 45, 35, 50, 100, 5, 1);
INSERT INTO public.monsters (monster_id, name, tag, level, hp, atk, def, max_xp, max_money, spawn_weight, region_id) VALUES (2, 'Guerreiro Lagarto', 'GUERREIRO_LAGARTO', 1, 80, 30, 20, 40, 80, 10,1);
INSERT INTO public.monsters (monster_id, name, tag, level, hp, atk, def, max_xp, max_money, spawn_weight, region_id) VALUES (3, 'Goblin Errante', 'GOBLIN_ERRANTE', 1, 50, 20, 15, 25, 60, 20,1);
INSERT INTO public.monsters (monster_id, name, tag, level, hp, atk, def, max_xp, max_money, spawn_weight, region_id) VALUES (4, 'Lobo das Planícies', 'LOBO_DAS_PLANICIES', 1, 40, 12, 5, 20, 40, 25,1);
INSERT INTO public.monsters (monster_id, name, tag, level, hp, atk, def, max_xp, max_money, spawn_weight, region_id) VALUES (5, 'Morcego Sombrio', 'MORCEGO_SOMBRIO', 1, 30, 7, 6, 10, 20, 40,1);
INSERT INTO public.regions (region_id, name, tag, min_level) VALUES (1, 'Planícies', 'PLAIN_LV1', 1);
INSERT INTO public.regions (region_id, name, tag, min_level) VALUES (2, 'Forest of Spirits', 'FLOREST_LV15', 15);
INSERT INTO public.regions (region_id, name, tag, min_level) VALUES (3, 'Volcano Lands', 'VULCAO_LV35', 35);
INSERT INTO public.regions (region_id, name, tag, min_level) VALUES (4, 'Frozen Peaks', 'FROZEN_LV50', 50);
INSERT INTO public.regions (region_id, name, tag, min_level) VALUES (5, 'Swamp of Shadows', 'SWAMP_LV70', 70);


-- =========================
-- CONSUMABLE ITEMS
-- =========================
INSERT INTO items (code, name, description, type, rarity, price) VALUES
-- HP
('POTION_HP_SMALL', 'Poção de Vida Pequena', 'Recupera uma pequena quantidade de HP.', 'CONSUMABLE', 'COMMON', 150),
('POTION_HP_MEDIUM', 'Poção de Vida Média', 'Recupera uma quantidade moderada de HP.', 'CONSUMABLE', 'UNCOMMON', 325),
('POTION_HP_LARGE', 'Poção de Vida Grande', 'Recupera uma grande quantidade de HP.', 'CONSUMABLE', 'RARE', 500),

-- ENERGY
('POTION_ENERGY_SMALL', 'Poção de Energia Pequena', 'Recupera uma pequena quantidade de Energia.', 'CONSUMABLE', 'COMMON', 400),
('POTION_ENERGY_MEDIUM', 'Poção de Energia Média', 'Recupera uma quantidade moderada de Energia.', 'CONSUMABLE', 'UNCOMMON', 700),
('POTION_ENERGY_LARGE', 'Poção de Energia Grande', 'Recupera uma grande quantidade de Energia.', 'CONSUMABLE', 'RARE',1000),

-- BUFFS
('POTION_ATK', 'Essência de Força', 'Aumenta o ATK temporariamente.', 'CONSUMABLE', 'RARE', 350),
('POTION_DEF', 'Essência de Defesa', 'Aumenta o DEF temporariamente.', 'CONSUMABLE', 'RARE', 350);

-- =========================
-- EQUIPMENT ITEMS (BASE)
-- =========================
INSERT INTO items (code, name, description, type, rarity, price) VALUES
-- WEAPONS
('WEAPON_TRAINING', 'Lâmina de Treino', 'Uma arma simples para iniciantes.', 'EQUIPMENT', 'COMMON', 150),
('WEAPON_WARRIOR', 'Lâmina do Guerreiro', 'Uma arma equilibrada e confiável.', 'EQUIPMENT', 'UNCOMMON', 500),
('WEAPON_MASTER', 'Lâmina do Mestre', 'Uma arma imbuída de Ki poderoso.', 'EQUIPMENT', 'EPIC', 1500),

-- GLOVES
('GLOVES_TRAINING', 'Luvas de Treino', 'Luvas simples para combate básico.', 'EQUIPMENT', 'COMMON', 150),
('GLOVES_WARRIOR', 'Luvas do Guerreiro', 'Luvas reforçadas com Ki.', 'EQUIPMENT', 'UNCOMMON', 500),
('GLOVES_MASTER', 'Luvas do Mestre', 'Luvas que canalizam energia pura.', 'EQUIPMENT', 'EPIC', 1500),

-- BOOTS
('BOOTS_TRAINING', 'Botas de Treino', 'Botas simples para longas caminhadas.', 'EQUIPMENT', 'COMMON', 150),
('BOOTS_WARRIOR', 'Botas do Guerreiro', 'Botas leves e resistentes.', 'EQUIPMENT', 'UNCOMMON', 500),
('BOOTS_MASTER', 'Botas do Mestre', 'Botas que aumentam a mobilidade.', 'EQUIPMENT', 'EPIC', 1500),

-- RINGS
('RING_TRAINING', 'Anel Simples', 'Um anel comum.', 'EQUIPMENT', 'COMMON', 150),
('RING_WARRIOR', 'Anel do Guerreiro', 'Um anel que fortalece o espírito.', 'EQUIPMENT', 'UNCOMMON', 500),
('RING_MASTER', 'Anel do Mestre', 'Um anel carregado de Ki.', 'EQUIPMENT', 'EPIC', 1500),

-- AMULETS
('AMULET_TRAINING', 'Amuleto Simples', 'Amuleto básico de proteção.', 'EQUIPMENT', 'COMMON', 150),
('AMULET_WARRIOR', 'Amuleto do Guerreiro', 'Amuleto energizado.', 'EQUIPMENT', 'UNCOMMON', 500),
('AMULET_MASTER', 'Amuleto do Mestre', 'Amuleto ancestral de grande poder.', 'EQUIPMENT', 'EPIC', 1500);
