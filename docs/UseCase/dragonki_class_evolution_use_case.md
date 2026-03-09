
# DragonKi — Caso de Uso
## Módulo: Evolução de Classe

---

# Descrição

O módulo de **Evolução de Classe** permite que o jogador avance sua especialização em marcos específicos da progressão do jogo.

A evolução ocorre em três estágios principais:

- **Level 15 → Escolha da especialização**
- **Level 35 → Evolução automática**
- **Level 50 → Evolução automática final**

A escolha do jogador acontece **somente no nível 15**.
A partir desse momento, toda evolução futura ocorre automaticamente dentro da mesma árvore de classe escolhida.

---

# Objetivo do módulo

Criar um sistema de progressão que:

- fortaleça a identidade do personagem
- ofereça escolha estratégica no início da progressão
- simplifique a evolução posterior
- mantenha consistência na árvore de classes
- permita expansão futura de classes

---

# Estrutura de progressão

Level 15 → escolha da especialização  
Level 35 → evolução automática  
Level 50 → evolução automática final

---

# Classe Base

Blade Monk

---

# Evolução Level 15 (Escolha do jogador)

Blade Monk
 ├─ Ki Vanguard
 ├─ Earth Sentinel
 └─ Wind Striker

---

# Evolução Level 35 (Automática)

Ki Vanguard → Storm Vanguard  
Earth Sentinel → Terra Guardian  
Wind Striker → Sky Striker

---

# Evolução Level 50 (Automática)

Storm Vanguard → Tempest Avatar  
Terra Guardian → Titan Sentinel  
Sky Striker → Hurricane Ascendant

---

# Regras de Negócio

1. Evolução ocorre apenas nos níveis 15, 35 e 50
2. A escolha do jogador ocorre somente no nível 15
3. Evoluções posteriores são automáticas
4. A evolução respeita a árvore de classes
5. A classe do jogador é atualizada após cada evolução

---

# Estrutura de dados sugerida

```go
type Class struct {
    ID          int64
    Name        string
    BaseClass   string
    Tier        string
    UnlockLevel int
    ParentID    *int64
}
```

---

# Consulta de evolução automática

```sql
SELECT *
FROM classes
WHERE parent_id = player.class_id
AND unlock_level = player.level
LIMIT 1;
```

---

# Pseudocódigo

```go
func AutoEvolveClass(player *Player) {

    if player.Level != 35 && player.Level != 50 {
        return
    }

    nextClass := classRepo.FindByParentAndLevel(player.ClassID, player.Level)

    if nextClass == nil {
        return
    }

    player.ClassID = nextClass.ID
    playerRepo.Update(player)
}
```

---

DragonKi Game Design Document
Class Evolution System
