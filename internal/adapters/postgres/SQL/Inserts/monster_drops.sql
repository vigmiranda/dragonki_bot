DELETE FROM public.monster_drops WHERE id BETWEEN 1 AND 30;

INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (1, 1, 32, 0.05, 1, 1, 0);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (2, 4, 32, 0.04, 1, 1, 0);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (3, 6, 32, 0.04, 1, 1, 0);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (4, 9, 32, 0.08, 1, 1, 0);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (5, 10, 32, 0.2, 1, 1, 0);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (6, 1, 35, 0.05, 1, 1, 0);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (7, 4, 35, 0.04, 1, 1, 0);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (8, 6, 35, 0.04, 1, 1, 0);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (9, 9, 35, 0.08, 1, 1, 0);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (10, 10, 35, 0.2, 1, 1, 0);

-- Equipment drops strategy:
-- - 1 equipment entry per monster
-- - low chance (rare reward), qty fixed 1
-- - rarity progression by monster level/region
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (11, 1, 40, 0.015, 1, 1, 10);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (12, 2, 43, 0.012, 1, 1, 10);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (13, 3, 46, 0.012, 1, 1, 10);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (14, 4, 49, 0.010, 1, 1, 10);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (15, 5, 52, 0.010, 1, 1, 10);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (16, 6, 40, 0.010, 1, 1, 10);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (17, 7, 43, 0.010, 1, 1, 10);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (18, 8, 46, 0.010, 1, 1, 10);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (19, 9, 41, 0.008, 1, 1, 6);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (20, 10, 47, 0.018, 1, 1, 6);

INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (21, 11, 50, 0.012, 1, 1, 6);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (22, 12, 53, 0.011, 1, 1, 6);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (23, 13, 41, 0.011, 1, 1, 6);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (24, 14, 44, 0.011, 1, 1, 6);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (25, 15, 47, 0.011, 1, 1, 6);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (26, 16, 50, 0.010, 1, 1, 6);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (27, 17, 42, 0.005, 1, 1, 2);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (28, 18, 45, 0.005, 1, 1, 2);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (29, 19, 48, 0.005, 1, 1, 2);
INSERT INTO public.monster_drops (id, monster_id, item_id, drop_rate, min_qty, max_qty, rarity_weight) VALUES (30, 20, 54, 0.007, 1, 1, 2);
