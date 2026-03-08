# DragonKi -- Big Boss Generation Standard

Este documento define o **padrão oficial para criação de Big Boss no
DragonKi**. Ele serve como template para gerar artes, conteúdo e design
consistente para os bosses semanais do jogo.

------------------------------------------------------------------------

# 🎨 Estilo Visual

Os Big Boss devem seguir sempre o mesmo padrão visual:

-   Estilo **fantasy RPG card**
-   Moldura **ornamentada temática da região**
-   Logo **DragonKi no topo**
-   Boss central dominante na arte
-   Estatísticas exibidas no rodapé
-   Atmosfera da região (lava, gelo, floresta, pântano etc.)
-   Arte altamente detalhada estilo **fantasy anime / game card**

Formato recomendado:

-   **Vertical**
-   **Aspect Ratio: 4:5**
-   **Alta definição**
-   **Estilo dark fantasy / epic fantasy**

------------------------------------------------------------------------

# 🧠 Estrutura do Boss

Cada Big Boss deve possuir:

    Nome do Boss
    Título do Boss

    HP: X,XXX,XXX
    ATK: XXX

Exemplo:

    Leviathan
    Terror das Profundezas

    HP: 10,000,000
    ATK: 650

------------------------------------------------------------------------

# 🧾 Prompt Template Oficial

Use este prompt para gerar novos bosses:

    Fantasy RPG boss card illustration for the game DragonKi.

    A massive boss creature themed around {REGION_THEME}. 
    The boss name is "{BOSS_NAME}" and its title is "{BOSS_TITLE}".

    Environment: {REGION_DESCRIPTION}

    The creature should look extremely powerful and intimidating, designed as a weekly raid boss.

    Style:
    - epic fantasy illustration
    - dark fantasy atmosphere
    - ultra detailed
    - glowing magical effects
    - dramatic lighting
    - highly detailed textures

    Card layout:
    - ornate fantasy frame themed around the region
    - DragonKi title at the top
    - boss name in large fantasy font
    - subtitle below the name
    - stats at the bottom

    Stats displayed:

    HP: {HP_VALUE}
    ATK: {ATK_VALUE}

    Format:
    vertical fantasy game card
    4:5 aspect ratio
    very high detail

------------------------------------------------------------------------

# 🌍 Regiões Oficiais do DragonKi

  Região               Estilo Visual
  -------------------- -------------------------------------------
  Volcano Lands        lava, magma, fogo
  Forest of Spirits    energia verde, espíritos, árvores antigas
  Frozen Peaks         gelo, tempestades de neve
  Swamp of Shadows     pântano sombrio, necromancia
  Celestial Ruins      ruínas antigas flutuantes
  Desert of Eternity   deserto místico, titãs antigos

------------------------------------------------------------------------

# 🧩 Exemplo de Boss

    BOSS_NAME: Thaldrion
    BOSS_TITLE: Colosso da Floresta Ancestral

    REGION_THEME: Forest of Spirits
    REGION_DESCRIPTION: ancient mystical forest with glowing spirits and giant twisted trees

    HP: 8,200,000
    ATK: 710

------------------------------------------------------------------------

# 📁 Estrutura recomendada no projeto

Sugestão de organização no repositório:

    /docs
       big_boss_prompt_template.md
       big_boss_regions.md
       big_boss_stats.md

------------------------------------------------------------------------

# 🔁 Ciclo semanal de Big Boss

Uma sugestão de rotação semanal:

    Semana 1 — Volcano Lands
    Semana 2 — Frozen Peaks
    Semana 3 — Forest of Spirits
    Semana 4 — Swamp of Shadows

Esse sistema cria **conteúdo recorrente para os jogadores** e incentiva
participação semanal.

------------------------------------------------------------------------

# 🐉 DragonKi

Documento oficial de padronização de criação de Big Boss.
