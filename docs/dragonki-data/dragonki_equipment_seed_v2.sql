
BEGIN;

INSERT INTO items (id, code, name, description, type, rarity, required_level, price, sort_order)
VALUES
(65,'WEAPON_TEMPESTADE','Espada da Tempestade','Arma do conjunto Tempestade','EQUIPMENT','RARE',15,1800,11),
(66,'HELMET_TEMPESTADE','Elmo da Tempestade','Elmo do conjunto Tempestade','EQUIPMENT','RARE',15,1500,12),
(67,'ARMOR_TEMPESTADE','Armadura da Tempestade','Armadura do conjunto Tempestade','EQUIPMENT','RARE',15,2200,13),

(68,'WEAPON_GUARDIAO','Bastão do Guardião','Arma do conjunto Guardião','EQUIPMENT','RARE',15,1750,14),
(69,'HELMET_GUARDIAO','Elmo do Guardião','Elmo do conjunto Guardião','EQUIPMENT','RARE',15,1500,15),
(70,'ARMOR_GUARDIAO','Couraça do Guardião','Armadura do conjunto Guardião','EQUIPMENT','RARE',15,2300,16),

(71,'WEAPON_VULCANICO','Espada Vulcânica','Arma do conjunto Vulcânico','EQUIPMENT','EPIC',35,4200,17),
(72,'HELMET_VULCANICO','Máscara Vulcânica','Elmo do conjunto Vulcânico','EQUIPMENT','EPIC',35,3200,18),
(73,'ARMOR_VULCANICO','Armadura Vulcânica','Armadura do conjunto Vulcânico','EQUIPMENT','EPIC',35,5000,19),

(74,'WEAPON_ASCETA','Bastão Ancestral','Arma do conjunto Asceta','EQUIPMENT','EPIC',35,4000,20),
(75,'HELMET_ASCETA','Capuz Ancestral','Elmo do conjunto Asceta','EQUIPMENT','EPIC',35,3000,21),
(76,'ARMOR_ASCETA','Veste Ancestral','Armadura do conjunto Asceta','EQUIPMENT','EPIC',35,4700,22),

(77,'WEAPON_GELIDO','Espada Gélida','Arma do conjunto Gélido','EQUIPMENT','LEGENDARY',50,7000,23),
(78,'HELMET_GELIDO','Capuz Gélido','Elmo do conjunto Gélido','EQUIPMENT','LEGENDARY',50,5200,24),
(79,'ARMOR_GELIDO','Armadura Gélida','Armadura do conjunto Gélido','EQUIPMENT','LEGENDARY',50,8600,25),

(80,'WEAPON_SOMBRAS','Bastão Sombrio','Arma do conjunto Sombras','EQUIPMENT','LEGENDARY',70,9800,26),
(81,'HELMET_SOMBRAS','Máscara das Sombras','Elmo do conjunto Sombras','EQUIPMENT','LEGENDARY',70,7200,27),
(82,'ARMOR_SOMBRAS','Armadura das Sombras','Armadura do conjunto Sombras','EQUIPMENT','LEGENDARY',70,12000,28);

COMMIT;
