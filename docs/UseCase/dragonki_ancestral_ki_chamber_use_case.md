# DragonKi --- Caso de Uso

## Módulo: Câmara do Ki Ancestral

------------------------------------------------------------------------

# Descrição

A **Câmara do Ki Ancestral** é um modo especial de desafio progressivo
no DragonKi onde o jogador enfrenta uma sequência de batalhas para
provar sua evolução.

Cada nível da câmara possui:

-   5 batalhas consecutivas
-   custo de entrada em energia e gold (1 de energia e 150 gold até sair ou concluir/morrer)
-   dificuldade escalável baseada no nível do player
-   progressão obrigatória entre níveis

O jogador pode:

-   usar consumíveis entre batalhas
-   se transformar entre batalhas
-   escolher continuar ou sair após cada vitória

O nível da câmara **só é concluído ao vencer as 5 batalhas**.

Se o jogador sair antes disso, o progresso é perdido.

------------------------------------------------------------------------

# Objetivo do módulo

Criar um modo de desafio que:

-   consuma recursos do jogador
-   incentive estratégia
-   valorize consumíveis
-   ofereça progressão escalável
-   funcione como conteúdo de longo prazo

------------------------------------------------------------------------

# Atores

-   Player
-   Sistema DragonKi
-   Câmara do Ki Ancestral

------------------------------------------------------------------------

# Regras principais

## Custo de entrada

Cada tentativa custa:

-   ⚡ 1 Energy
-   🪙 150 Gold

Se o jogador não possuir recursos suficientes:

❌ A entrada na câmara é bloqueada.

------------------------------------------------------------------------

# Estrutura da Câmara

A câmara possui níveis:

Nível 1 → Nível 10

Cada nível possui:

5 batalhas obrigatórias.

------------------------------------------------------------------------

# Progressão

-   Nível 1 começa liberado
-   Nível 2 só libera após concluir Nível 1
-   Nível 3 só libera após concluir Nível 2
-   ...
-   Nível 10 só libera após concluir Nível 9

------------------------------------------------------------------------

# Regra de conclusão

Um nível da câmara é concluído apenas quando:

5 vitórias consecutivas são obtidas.

Se o player sair antes disso:

-   o progresso é perdido
-   a próxima tentativa reinicia na batalha 1

------------------------------------------------------------------------

# Escolha após vitória

Após cada batalha vencida o sistema pergunta:

Continuar ou sair?

Opções:

⚔️ Continuar\
🚪 Sair

------------------------------------------------------------------------

# Uso de itens

Entre batalhas o jogador pode:

-   usar poções
-   restaurar HP
-   restaurar energia
-   ativar buffs

------------------------------------------------------------------------

# Transformações

Entre batalhas o jogador pode ativar transformações normalmente.

------------------------------------------------------------------------

# Escalonamento de dificuldade

A dificuldade depende do:

-   nível da câmara
-   nível atual do player

------------------------------------------------------------------------

# Nível 1

Monstros possíveis:

playerLevel -1\
playerLevel\
playerLevel +1

(mínimo nível 1)

------------------------------------------------------------------------

# Nível 2

Monstros entre:

playerLevel\
até\
playerLevel +2

------------------------------------------------------------------------

# Nível 3

Monstros entre:

playerLevel\
até\
playerLevel +3

------------------------------------------------------------------------

# Regra geral

Para câmaras nível **N**:

monstros entre

playerLevel\
até\
playerLevel + N

------------------------------------------------------------------------

# Exemplo

Player nível 20

Câmara nível 1:

19 → 21

Câmara nível 5:

20 → 25

Câmara nível 10:

20 → 30

------------------------------------------------------------------------

# Fluxo principal

## Entrar na Câmara

1.  Player acessa o menu
2.  Sistema mostra custo
3.  Player confirma entrada
4.  Sistema valida recursos
5.  Recursos são consumidos
6.  Run é iniciada
7.  Primeiro inimigo é sorteado

------------------------------------------------------------------------

# Estados da Run

ACTIVE\
COMPLETED\
FAILED\
ABANDONED

------------------------------------------------------------------------

# Estrutura de Domínio

## ChamberProgress

``` go
type ChamberProgress struct {
    ID int64
    PlayerID int64
    HighestCompletedLv int
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

------------------------------------------------------------------------

## ChamberRun

``` go
type ChamberRun struct {
    ID int64
    PlayerID int64
    ChamberLevel int
    CurrentBattle int
    Victories int
    Status string
    StartedAt time.Time
    FinishedAt *time.Time
}
```

------------------------------------------------------------------------

## ChamberBattleLog

``` go
type ChamberBattleLog struct {
    ID int64
    ChamberRunID int64
    BattleNumber int
    MonsterID int64
    MonsterLevel int
    Result string
    CreatedAt time.Time
}
```

------------------------------------------------------------------------

# Função de cálculo de nível do monstro

``` go
func GetChamberMonsterLevelRange(playerLevel int, chamberLevel int) (int, int) {

    if chamberLevel == 1 {
        min := playerLevel - 1
        if min < 1 {
            min = 1
        }

        return min, playerLevel + 1
    }

    return playerLevel, playerLevel + chamberLevel
}
```

------------------------------------------------------------------------

# Mensagem inicial

🌀 Você entrou na Câmara do Ki Ancestral!

Prepare-se para enfrentar 5 adversários.

------------------------------------------------------------------------

# Mensagem de vitória

✅ Vitória!

Batalhas vencidas: 3 / 5

Deseja continuar ou sair da câmara?

------------------------------------------------------------------------

# Conclusão

🏆 Câmara concluída!

Você venceu todas as 5 batalhas.

O próximo nível foi desbloqueado.

------------------------------------------------------------------------

# Nome técnico sugerido

ancestral_ki_chamber
