
# DragonKi — Estratégia de Equipamentos e Bônus de Conjunto (Set Bonus)

Documento de design para o sistema de equipamentos com sinergia no DragonKi.

---

# 1. Objetivo

Introduzir uma mecânica de **sinergia entre equipamentos**, permitindo que combinações de itens concedam bônus adicionais ao jogador.

Esse sistema aumenta:

- profundidade estratégica
- diversidade de builds
- valor dos drops
- motivação para colecionar equipamentos
- progressão de longo prazo

Essa mecânica é conhecida em RPGs como **Set Bonus**.

---

# 2. Conceito

Quando o jogador equipa múltiplas peças pertencentes ao **mesmo conjunto de equipamentos**, ele recebe bônus adicionais de atributos.

Exemplo:

Set Tempestade

| Peças Equipadas | Bônus |
|-----------------|------|
2 peças | +3 Speed |
3 peças | +5 Speed |
4 peças | +10 Speed +5 ATK |

---

# 3. Estrutura de Dados

## Tabela: equipment_sets

Representa um conjunto de equipamentos.

```
id
name
description
```

---

## Tabela: equipment_set_bonuses

Define os bônus concedidos por número de peças.

```
id
set_id
pieces_required
atk_bonus
def_bonus
speed_bonus
hp_bonus
```

---

## Tabela: items

Adicionar campo:

```
set_id
```

Isso conecta o item a um conjunto.

---

# 4. Exemplo de Sets do DragonKi

## 🌪 Set Tempestade
Classe recomendada: **Wind Striker**

Peças:

- Elmo da Tempestade
- Armadura da Tempestade
- Luvas da Tempestade
- Botas da Tempestade

Bônus:

| Peças | Bônus |
|------|------|
2 | +4 Speed |
3 | +6 Speed |
4 | +8 Speed +4 ATK |

---

## 🪨 Set Guardião da Terra
Classe recomendada: **Earth Sentinel**

Peças:

- Elmo do Guardião
- Couraça do Guardião
- Manoplas do Guardião
- Grevas do Guardião

Bônus:

| Peças | Bônus |
|------|------|
2 | +6 DEF |
3 | +10 DEF |
4 | +12 DEF +5 HP |

---

## 🔥 Set Vulcânico
Classe recomendada: **Ki Vanguard**

Peças:

- Máscara Vulcânica
- Armadura Vulcânica
- Cinto Vulcânico
- Botas Vulcânicas

Bônus:

| Peças | Bônus |
|------|------|
2 | +5 ATK |
3 | +7 ATK |
4 | +10 ATK +3 Speed |

---

# 5. Lógica de Aplicação no Jogo

Durante o cálculo final de atributos do jogador:

1. Identificar os itens equipados
2. Contar quantas peças pertencem a cada conjunto
3. Aplicar bônus correspondentes

---

# 6. Exemplo de Implementação (Go)

```go
func ApplySetBonus(player *Player, equippedItems []Item) {

    setCount := map[int]int{}

    for _, item := range equippedItems {
        if item.SetID != 0 {
            setCount[item.SetID]++
        }
    }

    for setID, pieces := range setCount {

        bonuses := getSetBonuses(setID)

        for _, bonus := range bonuses {
            if pieces >= bonus.PiecesRequired {
                player.ATK += bonus.AtkBonus
                player.DEF += bonus.DefBonus
                player.Speed += bonus.SpeedBonus
                player.MaxHP += bonus.HpBonus
            }
        }
    }
}
```

---

# 7. Impacto no Gameplay

O sistema cria **três camadas de progressão**:

### Equipamentos Individuais
Itens fortes isoladamente.

### Equipamentos Raros
Itens com atributos superiores.

### Conjuntos Completos
Combinações que liberam bônus poderosos.

Isso incentiva:

- coleta de equipamentos
- troca entre jogadores
- builds especializadas

---

# 8. Possíveis Builds no DragonKi

## Build Evasiva
Alta velocidade e esquiva.

Set recomendado:
Tempestade

---

## Build Tanque
Alta defesa e vida.

Set recomendado:
Guardião da Terra

---

## Build Agressiva
Alto dano.

Set recomendado:
Vulcânico

---

# 9. Sets Regionais (Expansão Futura)

Sets podem ser vinculados a regiões específicas.

### 🌲 Forest Spirit Set
Bônus:

- chance de esquiva aumentada na floresta

### ❄ Frozen Peaks Set
Bônus:

- resistência ao gelo

### 🌋 Volcano Lands Set
Bônus:

- bônus de dano contra criaturas vulcânicas

---

# 10. Benefícios para o DragonKi

Esse sistema:

- aumenta o valor de equipamentos
- melhora o balanceamento de builds
- incentiva exploração de mapas
- cria metas de progressão
- fortalece a economia do jogo

---

DragonKi Equipment System Design
