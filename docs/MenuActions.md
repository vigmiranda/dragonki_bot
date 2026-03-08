# Wiki de Acoes de Menu e Submenu

Este documento lista as acoes de callback do bot Telegram, onde aparecem no menu e qual fluxo/servico executam.

## Mapa geral de dispatch
Arquivo base: `internal/adapters/telegram/callbacks.go`

- Onboarding: `to_start`, `select_man`, `select_woman`
- Navegacao de menu: `action_menu`, `menu_adventure`, `menu_character`, `menu_inventory_scope`, `menu_economy`, `menu_events`
- Mapa: `list_map`, `to_region`
- Player: `action_status`, `to_class`
- Batalha: `action_hunt`, `action_fight`, `action_escape`
- Rank: `open_rank`, `rank_top_battles`, `rank_top_wins`, `rank_top_losses`, `rank_top_escape`
- Transformacao/treino: `to_training`, `list_transform`, `action_transform`
- Energia/pousada: `to_sleep`, `action_sleep`, `action_rest`
- Loja: `to_shop`, `list_potions`, `list_equipments`, `open_shop_equipment_types`, `list_buy_equipments`, `list_equipments_type`, `list_sell_equipments`, `buy_item`, `sell_item`
- Cambio: `to_exchange`, `action_convert`
- Big Boss: `to_big_boss`, `big_boss_attack`, `big_boss_ranking`
- Inventario: `open_inventory`, `open_inventory_potions`, `open_inventory_equipaments`, `open_inventory_equipaments_type`, `open_inventory_equipped`, `use_item`, `equip_item`, `unequip_item`, `unequip_all`
- Pagamento: `action_buy_ki`

## 1) Onboarding

### `to_start`
- Origem: menu inicial (`onboardingMenu`)
- Handler: `handleOnboarding`
- Efeito: abre selecao de genero.

### `select_man`
- Origem: menu de genero (`choosingGenderMenu`)
- Handler: `handleOnboarding`
- Efeito: define genero `MAN`.

### `select_woman`
- Origem: menu de genero (`choosingGenderMenu`)
- Handler: `handleOnboarding`
- Efeito: define genero `WOMAN`.

## 2) Navegacao por escopo

### `action_menu`
- Origem: botoes "Voltar" em varios menus
- Handler: `handleMenu`
- Efeito: recarrega painel principal com status/regiao do player.

### `menu_adventure`
- Origem: menu principal
- Handler: `handleMenu`
- Efeito: abre submenu Aventura + informacoes atuais do player/regiao.

### `menu_character`
- Origem: menu principal
- Handler: `handleMenu`
- Efeito: abre submenu Personagem.

### `menu_inventory_scope`
- Origem: menu principal
- Handler: `handleMenu`
- Efeito: abre submenu Inventario.

### `menu_economy`
- Origem: menu principal
- Handler: `handleMenu`
- Efeito: abre submenu Economia.

### `menu_events`
- Origem: menu principal
- Handler: `handleMenu`
- Efeito: abre submenu Eventos (Big Boss quando flag ligada).

## 3) Aventura / Mapa / Batalha

### `action_hunt`
- Origem: submenu Aventura
- Handler: `handleBattle`
- Efeito: inicia encontro com monstro.

### `list_map`
- Origem: submenu Aventura
- Handler: `handleMap`
- Efeito: lista regioes disponiveis.

### `to_region?regionID=<id>`
- Origem: menu de regioes (`RegionsMenu`)
- Handler: `handleMap`
- Params:
  - `regionID` (obrigatorio)
- Efeito: move o player para regiao selecionada.

### `action_fight?monsterID=<id>`
- Origem: menu de batalha (`battleMenu`)
- Handler: `handleBattle`
- Params:
  - `monsterID` (obrigatorio)
- Efeito: executa batalha contra o monstro.

### `action_escape?monsterID=<id>`
- Origem: menu de batalha (`battleMenu`)
- Handler: `handleBattle`
- Params:
  - `monsterID` (obrigatorio)
- Efeito: tenta escapar da batalha.

## 4) Personagem / Treino / Energia

### `action_status`
- Origem: submenu Personagem
- Handler: `handlePlayer`
- Efeito: exibe status detalhado.

### `to_class?classID=<id>`
- Origem: escolha de classe (`ChoiceClassMenu`)
- Handler: `handlePlayer`
- Params:
  - `classID` (obrigatorio)
- Efeito: atualiza classe do player.

### `to_training`
- Origem: submenu Aventura
- Handler: `handleTraining`
- Efeito: abre tela de treinamento.

### `list_transform`
- Origem: submenu Personagem
- Handler: `handleTraining`
- Efeito: lista informacoes de transformacao.

### `action_transform`
- Origem: menu de transformacao
- Handler: `handleTraining`
- Efeito: tenta transformar (consome energia).

### `to_sleep`
- Origem: submenu Personagem
- Handler: `handleEnergy`
- Efeito: abre menu de descanso.

### `action_sleep`
- Origem: menu de energia
- Handler: `handleEnergy`
- Efeito: dorme (custo em gold / recuperacao alta).

### `action_rest`
- Origem: menu de energia
- Handler: `handleEnergy`
- Efeito: descanso leve (custo baixo / recuperacao parcial).

## 5) Inventario

### `open_inventory`
- Origem: submenu Inventario
- Handler: `handleInventory`
- Efeito: abre menu principal do inventario.

### `open_inventory_potions?page=<n>`
- Origem: inventario
- Handler: `handleInventory`
- Params:
  - `page` (opcional, default 1)
- Efeito: lista consumiveis paginados.

### `open_inventory_equipaments`
- Origem: inventario
- Handler: `handleInventory`
- Efeito: abre tipos de equipamento no inventario.

### `open_inventory_equipaments_type?slot=<slot>&page=<n>`
- Origem: menu de tipos de equipamento
- Handler: `handleInventory`
- Params:
  - `slot` (obrigatorio: `WEAPON|ARMOR|GLOVES|HELMET|BOOTS|RING|AMULET`)
  - `page` (opcional, default 1)
- Efeito: lista equipamentos daquele slot.

### `open_inventory_equipped`
- Origem: submenu Inventario e menu de tipos
- Handler: `handleInventory`
- Efeito: exibe itens equipados e opcoes de desequipar.

### `use_item?inventoryItemID=<id>`
- Origem: lista de pocoes
- Handler: `handleInventory`
- Params:
  - `inventoryItemID` (obrigatorio)
- Efeito: usa item consumivel.

### `equip_item?inventoryItemID=<id>`
- Origem: lista de equipamentos no inventario
- Handler: `handleInventory`
- Params:
  - `inventoryItemID` (obrigatorio)
- Efeito: equipa item.

### `unequip_item?slot=<slot>`
- Origem: menu de equipados
- Handler: `handleInventory`
- Params:
  - `slot` (obrigatorio)
- Efeito: desequipa slot especifico.

### `unequip_all`
- Origem: menu de equipados
- Handler: `handleInventory`
- Efeito: desequipa todos os slots.

## 6) Economia (Loja + Cambio + Pagamento)

### `to_shop`
- Origem: submenu Economia
- Handler: `handleShop`
- Efeito: abre menu da loja.

### `list_potions`
- Origem: loja
- Handler: `handleShop`
- Efeito: lista pocoes disponiveis para compra.

### `list_equipments`
- Origem: loja
- Handler: `handleShop`
- Efeito: abre menu de acoes de equipamento (compra/venda).

### `open_shop_equipment_types`
- Origem: loja equipamentos (compra)
- Handler: `handleShop`
- Efeito: abre tipos de equipamentos para compra.

### `list_equipments_type?slot=<slot>`
- Origem: tipos de equipamento da loja
- Handler: `handleShop`
- Params:
  - `slot` (obrigatorio)
- Efeito: lista itens daquele tipo para compra.

### `buy_item?itemID=<id>`
- Origem: menus de itens da loja (`BuildMenuItems`)
- Handler: `handleShop`
- Params:
  - `itemID` (obrigatorio)
- Efeito: compra item.

### `list_sell_equipments?page=<n>`
- Origem: loja equipamentos (venda)
- Handler: `handleShop`
- Params:
  - `page` (opcional, default 1)
- Efeito: lista equipamentos nao equipados para venda.

### `sell_item?inventoryItemID=<id>&page=<n>`
- Origem: menu de venda (`ShopSellItemsMenu`)
- Handler: `handleShop`
- Params:
  - `inventoryItemID` (obrigatorio)
  - `page` (opcional)
- Efeito: vende item (40% do preco configurado).

### `to_exchange`
- Origem: submenu Economia
- Handler: `handleExchange`
- Efeito: abre casa de cambio.

### `action_convert?amountConvert=<n>`
- Origem: casa de cambio
- Handler: `handleExchange`
- Params:
  - `amountConvert` (obrigatorio)
- Efeito: converte Ki Crystals em Gold.

### `action_buy_ki`
- Origem: casa de cambio (quando `payments` habilitada)
- Handler: `handlePayment`
- Efeito: inicia compra de Ki Crystals (Mercado Pago).

## 7) Eventos (Big Boss)

### `to_big_boss`
- Origem: submenu Eventos (feature flag `big_boss_weekly`)
- Handler: `handleBigBoss`
- Efeito: abre painel do Big Boss da semana.

### `big_boss_attack`
- Origem: menu do Big Boss
- Handler: `handleBigBoss`
- Efeito: realiza ataque diario no boss.

### `big_boss_ranking`
- Origem: menu do Big Boss
- Handler: `handleBigBoss`
- Efeito: exibe ranking de dano da semana.

## 8) Rankings gerais

### `open_rank`
- Origem: menu principal
- Handler: `handleRank`
- Efeito: abre menu de ranking semanal de batalhas.

### `rank_top_battles`
- Origem: menu de rank
- Handler: `handleRank`
- Efeito: top batalhas da semana.

### `rank_top_wins`
- Origem: menu de rank
- Handler: `handleRank`
- Efeito: top vitorias da semana.

### `rank_top_losses`
- Origem: menu de rank
- Handler: `handleRank`
- Efeito: top derrotas da semana.

### `rank_top_escape`
- Origem: menu de rank
- Handler: `handleRank`
- Efeito: top escapes da semana.

## 9) Acao de URL (sem callback)

### `Entrar na Comunidade` (URL)
- Origem: menu principal
- Tipo: `NewInlineKeyboardButtonURL`
- Destino: `https://t.me/dragonki_community`
- Observacao: nao passa por `callbacks.go`.

## Observacoes
- Alguns callbacks antigos continuam ativos para compatibilidade (ex.: nomenclatura `equipaments` nos callbacks de inventario).
- Menus dinamicos (paginacao, lista de itens e lista de regioes/classes) geram callbacks em runtime com os parametros indicados acima.
