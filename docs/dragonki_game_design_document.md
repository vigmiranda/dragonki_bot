# DragonKi --- Documento Oficial do Projeto

## 1. Visão Geral

**DragonKi** é um RPG baseado em texto operado através de um **bot no
Telegram**, inspirado em elementos de anime (principalmente Dragon Ball)
e mecânicas clássicas de RPG.

O jogo combina: - Progressão de personagens - Sistema de combate
dinâmico - Classes especializadas - Exploração de regiões - Coleta de
itens - Evolução estratégica

O objetivo é criar uma experiência de RPG acessível diretamente dentro
do Telegram.

------------------------------------------------------------------------

# 2. Tecnologias do Projeto

## Backend

-   **Go (Golang)**

## Banco de Dados

-   **PostgreSQL**
-   ORM: **GORM**

## Plataforma

-   **Telegram Bot API**

## Arquitetura

Arquitetura modular baseada em camadas:

-   Domain
-   Application / Services
-   Adapters (DB / Telegram)
-   Infrastructure

------------------------------------------------------------------------

# 3. Conceito do Mundo

O mundo de **DragonKi** é alimentado pela energia ancestral chamada
**Ki**.

Antigamente existiam **Dragões Ancestrais**, seres capazes de manipular
o Ki primordial. Após sua extinção, apenas fragmentos de sua energia
permaneceram.

Esses fragmentos são conhecidos como:

### Ki Crystals

Pedras raríssimas formadas das lágrimas dos dragões ancestrais.\
Elas não aumentam diretamente o poder dos guerreiros, mas possuem enorme
valor e importância no mundo.

------------------------------------------------------------------------

# 4. Sistema de Atributos

Cada personagem possui atributos principais.

| Atributo \| Descrição \|

\|--------\|--------\|HP\| Vida total do personagem \| \|ATK\| Poder de
ataque \| \|DEF\| Defesa contra ataques \| \|Speed\| Velocidade, usada
para esquiva \| \|Energy\| Energia usada para habilidades \|

------------------------------------------------------------------------

# 5. Sistema de Combate

O combate **não é por turno tradicional**.

A batalha acontece em loops de ataque até que um dos lados seja
derrotado.

Exemplo simplificado:

-   Player ataca
-   Monstro contra-ataca
-   Chance de esquiva baseada em Speed

### Fórmula de Esquiva

Chance dinâmica baseada em Speed:

    dodgeChance = min(25, player.Speed * 0.5)

Limite máximo recomendado:

    25%

------------------------------------------------------------------------

# 6. Sistema de Progressão

## Níveis

Nível máximo atual planejado:

    Level 50

XP necessário cresce progressivamente.

Exemplo de progressão:

  Level   XP
  ------- -------
  1       100
  5       350
  10      900
  20      2500
  30      5000
  40      9000
  50      15000

------------------------------------------------------------------------

# 7. Classes (Desbloqueadas no Level 15)

Ao atingir **Level 15**, o jogador escolhe uma especialização.

## Wind Striker

Foco: **Velocidade e esquiva**

Características: - Alto Speed - Ataques rápidos - Chance maior de
esquiva

------------------------------------------------------------------------

## Earth Sentinel

Foco: **Defesa e resistência**

Características: - Alta DEF - Alto HP - Tank natural

------------------------------------------------------------------------

## Ki Vanguard

Foco: **Ataque bruto**

Características: - Alto ATK - Dano explosivo - Menos defesa

------------------------------------------------------------------------

# 8. Regiões do Mundo

## Volcano Lands

Faixa de nível:

    Level 35 — 45

Ambiente: - Vulcões ativos - Rios de lava - Criaturas do fogo

Estética: - Pedra vulcânica - Atmosfera infernal - Energia instável

------------------------------------------------------------------------

# 9. Tipos de Inimigos

## Inimigos Comuns

Criaturas frequentes da região.

Exemplos:

-   **Ignivar**
-   **Magmarok**
-   **Salamorth**
-   **Volkar de Muspel**
-   **Cinder Salamander**

------------------------------------------------------------------------

## Mini Boss

Criaturas mais raras e poderosas.

Drops melhores e mais XP.

------------------------------------------------------------------------

## Boss

Chefes da região.

Possuem:

-   Mecânicas especiais
-   Alto HP
-   Drops raros

------------------------------------------------------------------------

# 10. Sistema de Drops

Monstros podem dropar:

-   Ouro
-   Itens
-   Materiais

Drops são definidos por tabela.

Exemplo:

| Item \| Chance \|

\|----\|----\|Gold Coin\|60%\| \|Fire Fragment\|25%\| \|Rare Core\|10%\|
\|Legendary Drop\|5%\|

------------------------------------------------------------------------

# 11. Sistema de XP

XP recebido depende de:

-   Level do player
-   Level do monstro

Exemplo:

    xpBase = 25 * monsterLevel ^ 1.12

Multiplicadores ajustam a recompensa baseado na diferença de nível.

Limite recomendado:

    XP máximo por batalha: 120

------------------------------------------------------------------------

# 12. Sistema de Eventos

Possibilidades futuras:

-   Caçada global
-   Boss mundial
-   Eventos temporários
-   Invasões de regiões

------------------------------------------------------------------------

# 13. Direção Artística

Estilo visual:

**DragonKi Style**

Mistura de:

-   Anime
-   RPG Fantasy
-   Influência de Dragon Ball
-   Criaturas mitológicas

Regiões possuem identidade própria.

------------------------------------------------------------------------

# 14. Roadmap do Projeto

### Curto prazo

-   Balanceamento de XP
-   Criação de inimigos
-   Sistema de drops
-   Melhorias no combate

### Médio prazo

-   Sistema de itens
-   Equipamentos
-   Mais regiões

### Longo prazo

-   Multiplayer
-   Raids
-   Guildas

------------------------------------------------------------------------

# 15. Objetivo do Projeto

DragonKi busca criar um RPG acessível e divertido dentro do Telegram,
combinando:

-   Progressão satisfatória
-   Combate dinâmico
-   Exploração contínua
