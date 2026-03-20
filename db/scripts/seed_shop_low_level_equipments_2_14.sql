-- Generic low-level equipment seed for shop rotation (levels 2..14).
--
-- Goals:
-- - Add at least 5 items per shop category:
--   WEAPON, ARMOR, GLOVES, HELMET, BOOTS, RING, AMULET
-- - Weapons are limited to sword and staff themes.
-- - Items are generic and can be used by the whole progression tree.
-- - Bonuses are varied item-to-item to create light randomness while staying deterministic.
--
-- Idempotency:
-- - upsert by item code

WITH seq AS (
    SELECT
        pg_get_serial_sequence('public.items', 'id') AS seq_name,
        COALESCE((SELECT MAX(id) FROM public.items), 0) + 1 AS next_id
)
SELECT setval(seq_name::regclass, next_id, false)
FROM seq
WHERE seq_name IS NOT NULL;

WITH seq AS (
    SELECT
        pg_get_serial_sequence('public.equipments', 'id') AS seq_name,
        COALESCE((SELECT MAX(id) FROM public.equipments), 0) + 1 AS next_id
)
SELECT setval(seq_name::regclass, next_id, false)
FROM seq
WHERE seq_name IS NOT NULL;

DROP TABLE IF EXISTS tmp_seed_shop_low_level_equipments_2_14;

CREATE TEMP TABLE tmp_seed_shop_low_level_equipments_2_14 (
    code TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    slot TEXT NOT NULL,
    required_level INTEGER NOT NULL,
    rarity TEXT NOT NULL,
    price INTEGER NOT NULL,
    atk_bonus INTEGER NOT NULL,
    def_bonus INTEGER NOT NULL,
    speed_bonus INTEGER NOT NULL,
    hp_bonus INTEGER NOT NULL,
    energy_bonus INTEGER NOT NULL,
    allowed_classes TEXT NOT NULL DEFAULT ''
);

INSERT INTO tmp_seed_shop_low_level_equipments_2_14 (
    code,
    name,
    description,
    slot,
    required_level,
    rarity,
    price,
    atk_bonus,
    def_bonus,
    speed_bonus,
    hp_bonus,
    energy_bonus
) VALUES
    ('LL_WEAPON_02_OAK_SWORD',     'Espada de Carvalho',         'Lamina curta de treino com fio confiavel.',            'WEAPON',  2, 'COMMON',   180,  5, 0, 1, 0, 0),
    ('LL_WEAPON_05_HERMIT_STAFF',  'Bastao do Eremita',         'Bastao leve usado por viajantes e monges.',           'WEAPON',  5, 'COMMON',   260,  7, 1, 0, 0, 1),
    ('LL_WEAPON_08_SCOUT_SWORD',   'Espada do Batedor',         'Espada agil para quem valoriza velocidade.',          'WEAPON',  8, 'UNCOMMON', 390, 10, 0, 2, 0, 0),
    ('LL_WEAPON_11_GUARD_STAFF',   'Bastao da Guarda',          'Bastao reforcado para combate disciplinado.',         'WEAPON', 11, 'UNCOMMON', 520, 11, 2, 0, 0, 2),
    ('LL_WEAPON_14_DAWN_SWORD',    'Espada da Aurora',          'Espada equilibrada para guerreiros em ascensao.',     'WEAPON', 14, 'UNCOMMON', 680, 13, 1, 1, 0, 1),

    ('LL_ARMOR_02_PADDED_VEST',    'Colete Acolchoado',         'Protecao simples para aventuras iniciais.',           'ARMOR',   2, 'COMMON',   170,  0, 4, 0, 6, 0),
    ('LL_ARMOR_05_HIDE_CUIRASS',   'Couraca de Couro',          'Couraca leve e flexivel.',                            'ARMOR',   5, 'COMMON',   255,  1, 6, 0, 9, 0),
    ('LL_ARMOR_08_IRON_MESH',      'Malha de Ferro',            'Malha resistente para confrontos mais duros.',        'ARMOR',   8, 'UNCOMMON', 380,  1, 8, 0, 12, 0),
    ('LL_ARMOR_11_SENTRY_PLATE',   'Armadura de Sentinela',     'Armadura firme para linha de frente.',                'ARMOR',  11, 'UNCOMMON', 500,  2, 10, 0, 15, 0),
    ('LL_ARMOR_14_STEEL_WARD',     'Armadura de Aco Leve',      'Placas leves com boa sustentacao em combate.',        'ARMOR',  14, 'UNCOMMON', 660,  2, 12, 0, 18, 1),

    ('LL_GLOVES_02_TRAINING_WRAP', 'Luvas de Treino',           'Luvas simples para golpes controlados.',              'GLOVES',  2, 'COMMON',   140,  2, 0, 1, 0, 0),
    ('LL_GLOVES_05_HUNTER_GRIP',   'Luvas do Cacador',          'Pegada firme para ataques rapidos.',                  'GLOVES',  5, 'COMMON',   210,  3, 1, 1, 2, 0),
    ('LL_GLOVES_08_RUSH_GAUNTLET', 'Manoplas de Avanco',        'Luvas reforcadas para pressao ofensiva.',             'GLOVES',  8, 'UNCOMMON', 310,  5, 1, 1, 3, 0),
    ('LL_GLOVES_11_SWIFT_GRIP',    'Luvas do Impulso',          'Modelo leve para manter o ritmo da luta.',            'GLOVES', 11, 'UNCOMMON', 430,  4, 1, 2, 4, 1),
    ('LL_GLOVES_14_IRON_FIST',     'Manoplas Punho Firme',      'Canalizam impacto e estabilidade.',                   'GLOVES', 14, 'UNCOMMON', 560,  6, 2, 2, 5, 1),

    ('LL_HELMET_02_SCOUT_HOOD',    'Capuz de Vigia',            'Cobertura leve para patrulhas iniciais.',             'HELMET',  2, 'COMMON',   145,  0, 3, 0, 5, 1),
    ('LL_HELMET_05_BRONZE_HELM',   'Elmo de Bronze',            'Elmo simples, mas confiavel.',                        'HELMET',  5, 'COMMON',   220,  0, 4, 0, 7, 1),
    ('LL_HELMET_08_IRON_VISOR',    'Elmo de Ferro',             'Protecao firme para o rosto e a mente.',              'HELMET',  8, 'UNCOMMON', 320,  0, 6, 0, 9, 1),
    ('LL_HELMET_11_GUARD_VISOR',   'Visor da Guarda',           'Modelo tatico usado por soldados experientes.',       'HELMET', 11, 'UNCOMMON', 435,  1, 7, 0, 11, 2),
    ('LL_HELMET_14_STEEL_CROWN',   'Coroa de Aco',              'Elmo refinado para combatentes promissores.',         'HELMET', 14, 'UNCOMMON', 575,  1, 9, 0, 13, 2),

    ('LL_BOOTS_02_ROAMER_BOOTS',   'Botas de Caminhante',       'Botas leves para longas trilhas.',                    'BOOTS',   2, 'COMMON',   135,  0, 0, 2, 0, 1),
    ('LL_BOOTS_05_TRAIL_BOOTS',    'Botas de Trilha',           'Mantem equilibrio em terrenos irregulares.',          'BOOTS',   5, 'COMMON',   205,  0, 1, 3, 0, 1),
    ('LL_BOOTS_08_RUNNER_BOOTS',   'Botas do Corredor',         'Favorecem deslocamento rapido.',                      'BOOTS',   8, 'UNCOMMON', 300,  0, 1, 4, 0, 2),
    ('LL_BOOTS_11_WINDSTEP_BOOTS', 'Botas Passo Veloz',         'Aprimoram arranque e reposicionamento.',              'BOOTS',  11, 'UNCOMMON', 420,  0, 1, 5, 0, 2),
    ('LL_BOOTS_14_DAWNSTRIDE',     'Botas da Alvorada',         'Botas estaveis para o combate avancado inicial.',     'BOOTS',  14, 'UNCOMMON', 540,  1, 1, 6, 1, 2),

    ('LL_RING_02_COPPER_RING',     'Anel de Cobre',             'Anel simples com leve ressonancia de Ki.',            'RING',    2, 'COMMON',   145,  1, 1, 0, 2, 1),
    ('LL_RING_05_FOCUS_RING',      'Anel do Foco',              'Ajuda o guerreiro a manter a presenca em batalha.',   'RING',    5, 'COMMON',   225,  2, 1, 1, 3, 1),
    ('LL_RING_08_SPIRIT_RING',     'Anel do Espirito',          'Fortalece o fluxo de energia do usuario.',            'RING',    8, 'UNCOMMON', 335,  3, 1, 1, 4, 2),
    ('LL_RING_11_VIGOR_RING',      'Anel do Vigor',             'Amplifica presenca ofensiva e resistencia.',          'RING',   11, 'UNCOMMON', 455,  3, 2, 1, 5, 2),
    ('LL_RING_14_AURIC_RING',      'Anel Aureo',                'Anel refinado para guerreiros em evolucao.',          'RING',   14, 'UNCOMMON', 590,  4, 2, 1, 6, 2),

    ('LL_AMULET_02_SOFT_CHARM',    'Amuleto Sereno',            'Amuleto discreto com protecao basica.',               'AMULET',  2, 'COMMON',   150,  0, 2, 0, 4, 1),
    ('LL_AMULET_05_ROOT_CHARM',    'Amuleto da Raiz',           'Oferece estabilidade e um pouco mais de vigor.',      'AMULET',  5, 'COMMON',   230,  0, 3, 0, 6, 1),
    ('LL_AMULET_08_KI_CHARM',      'Amuleto de Ki',             'Ajuda a manter a energia em confrontos longos.',      'AMULET',  8, 'UNCOMMON', 340,  1, 4, 0, 8, 2),
    ('LL_AMULET_11_GUARD_CHARM',   'Amuleto Guardiao',          'Canaliza defesa e resistencia em equilibrio.',        'AMULET', 11, 'UNCOMMON', 470,  1, 5, 0, 10, 2),
    ('LL_AMULET_14_LUMEN_CHARM',   'Amuleto Lumen',             'Amuleto raro para a transicao ao jogo avancado.',     'AMULET', 14, 'UNCOMMON', 610,  1, 6, 0, 12, 2);

UPDATE public.items i
   SET name = cfg.name,
       description = cfg.description,
       type = 'EQUIPMENT',
       rarity = cfg.rarity,
       amount = 1,
       price = cfg.price
  FROM tmp_seed_shop_low_level_equipments_2_14 cfg
 WHERE i.code = cfg.code;

INSERT INTO public.items (
    code,
    name,
    description,
    type,
    rarity,
    amount,
    price
)
SELECT
    cfg.code,
    cfg.name,
    cfg.description,
    'EQUIPMENT',
    cfg.rarity,
    1,
    cfg.price
FROM tmp_seed_shop_low_level_equipments_2_14 cfg
WHERE NOT EXISTS (
    SELECT 1
    FROM public.items i
    WHERE i.code = cfg.code
);

UPDATE public.equipments e
   SET slot = cfg.slot,
       required_level = cfg.required_level,
       atk_bonus = cfg.atk_bonus,
       def_bonus = cfg.def_bonus,
       speed_bonus = cfg.speed_bonus,
       hp_bonus = cfg.hp_bonus,
       energy_bonus = cfg.energy_bonus,
       allowed_classes = cfg.allowed_classes
  FROM public.items i
  JOIN tmp_seed_shop_low_level_equipments_2_14 cfg
    ON cfg.code = i.code
 WHERE e.item_id = i.id;

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
    cfg.slot,
    cfg.required_level,
    cfg.atk_bonus,
    cfg.def_bonus,
    cfg.speed_bonus,
    cfg.hp_bonus,
    cfg.energy_bonus,
    cfg.allowed_classes
FROM public.items i
JOIN tmp_seed_shop_low_level_equipments_2_14 cfg
  ON cfg.code = i.code
WHERE NOT EXISTS (
    SELECT 1
    FROM public.equipments e
    WHERE e.item_id = i.id
);

UPDATE public.items i
   SET equipment_id = e.id
  FROM public.equipments e
 WHERE e.item_id = i.id
   AND EXISTS (
       SELECT 1
       FROM tmp_seed_shop_low_level_equipments_2_14 cfg
       WHERE cfg.code = i.code
   );

UPDATE public.equipments
   SET allowed_classes = ''
 WHERE required_level <= 14
   AND EXISTS (
       SELECT 1
       FROM public.items i
       WHERE i.id = public.equipments.item_id
         AND i.code LIKE 'LL\_%' ESCAPE '\'
   );

DROP TABLE IF EXISTS tmp_seed_shop_low_level_equipments_2_14;
