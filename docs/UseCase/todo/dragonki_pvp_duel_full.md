
# DragonKi — Sistema PvP (Duelo entre Jogadores)
Documento Oficial Completo

---

# 🧭 1. Visão Geral

O sistema de PvP permite que jogadores se enfrentem em duelos através de matchmaking automático.

O combate utiliza os atributos reais do jogador:

- Classe
- Equipamentos
- Transformações
- Bônus ativos

---

# 🎯 2. Objetivos

- Criar competição entre jogadores
- Incentivar builds e estratégias
- Gerar economia (gold sink)
- Aumentar retenção

---

# ⚔️ 3. Acesso ao PvP

Menu → PvP

---

# 🔎 4. Busca por Duelo

Ao clicar em "Buscar duelo":

- Player entra na fila
- Sistema inicia matchmaking

---

# 📊 5. Matchmaking

## Range inicial

±3 níveis

## Expansão

±5 → ±7 → ±10

## Condição de parada

- Encontrar match
- Ou retornar erro

---

# 🤝 6. Pareamento

- Dois jogadores compatíveis são vinculados
- Criado registro de match
- Status: PENDING_ACCEPT

---

# 👁️ 7. Exibição do Adversário

Cada jogador verá:

- Nome
- Classe
- Nível
- HP
- ATK
- DEF
- SPEED
- Transformação ativa

Todos valores já incluem bônus de equipamentos e buffs.

---

# ✅ 8. Aceite do Duelo

## Regras

- Ambos devem aceitar
- Timeout: 30 segundos

## Fluxos

✔ Ambos aceitam → combate  
❌ Um recusa → volta para fila  

---

# ❌ 9. Cancelamento

- Player pode cancelar busca
- Sai da fila
- Sem penalidade

---

# ⚔️ 10. Sistema de Combate

- Mesmo sistema PvE
- Loop até HP = 0
- Considera speed (esquiva)

---

# 🔋 11. Energia

- -1 energia após combate
- Apenas se duelo ocorrer

---

# ❤️ 12. Regeneração

- HP restaurado após duelo

---

# 💰 13. Economia

## Vencedor

+5% gold do adversário

## Perdedor

-7% gold

## Sistema

2% removido (gold sink)

---

# 🛡️ 14. Proteções

- Gold >= 0
- Piso opcional
- Teto opcional

---

# 🚫 15. Restrições

- Não pode estar em dungeon
- Não pode estar em caçada
- Não pode estar em outro PvP
- Deve ter energia

---

# ⚠️ 16. Edge Cases

## Timeout de aceite

- Conta como recusa

## Player desconecta

- Match cancelado

## Duplo match

- Prevenir via lock/transação

## Race condition

- Matchmaking deve ser atômico

---

# 🗄️ 17. Persistência

## pvp_queue

id  
player_id  
level  
status  
created_at  

---

## pvp_matches

id  
player_one_id  
player_two_id  
status  
winner_id  
loser_id  
created_at  

---

## pvp_acceptance

id  
match_id  
player_id  
accepted  

---

## pvp_results

id  
match_id  
winner_id  
loser_id  
winner_gold  
loser_gold  
created_at  

---

# 🔄 18. Fluxo Completo

Entrar PvP  
↓  
Buscar duelo  
↓  
Match encontrado  
↓  
Mostrar stats  
↓  
Aceite  
↓  
Combate  
↓  
Resultado  
↓  
Aplicar economia  
↓  
Reset HP  
↓  
Consumir energia  

---

# 🧠 19. Estratégia Técnica

## Matchmaking

- Fila em memória ou DB
- Ordenação por level
- Busca incremental

## Concorrência

- Lock por player
- Evitar múltiplos matches

## Transações

- Atualização de gold deve ser ACID

---

# 🚀 20. Evoluções Futuras

- Ranking PvP
- Sistema de Elo
- Ligas
- Histórico
- Recompensas semanais

---

# 🏁 21. Conclusão

O sistema PvP:

- Cria competição
- Movimenta economia
- Incentiva builds
- Aumenta retenção

Fundamental para evolução do DragonKi.
