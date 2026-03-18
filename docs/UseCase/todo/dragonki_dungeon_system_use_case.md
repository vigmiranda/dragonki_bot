
# DragonKi — Sistema de Dungeons
Documento de Caso de Uso

## Visão Geral
O sistema de **Dungeons** adiciona conteúdo cooperativo ao DragonKi, permitindo que até **3 jogadores** enfrentem uma sequência de inimigos e um **Boss final** para obter recompensas especiais.

As Dungeons funcionam de forma **assíncrona**, permitindo que os jogadores avancem mesmo sem estarem online ao mesmo tempo.

---

# 1. Acesso à Dungeon
Para iniciar uma dungeon é necessário possuir um item especial:

**Dungeon Key**

### Regras
- Apenas **1 chave é consumida**
- O jogador que usar a chave torna-se o **Líder da Dungeon**
- O líder pode convidar até **2 jogadores adicionais**

### Limite de participantes
1 a 3 jogadores

---

# 2. Criação da Dungeon
Fluxo:

Menu → Dungeon → Criar Dungeon

Passos:

1. Jogador usa **Dungeon Key**
2. Sistema cria uma nova dungeon
3. Jogador convida até **2 aliados**
4. Dungeon é iniciada

---

# 3. Escalonamento de Dificuldade
A dungeon é ajustada pelo **nível do jogador mais forte**.

Exemplo:

Player A: Lv 28  
Player B: Lv 34  
Player C: Lv 31  

Dungeon será configurada para **Lv 34**.

---

# 4. Estrutura da Dungeon

Wave 1 – 2 inimigos  
Wave 2 – 3 inimigos  
Wave 3 – inimigo elite  
Wave 4 – inimigo elite  
Boss final  

Total aproximado: **6 a 8 combates**

---

# 5. Sistema de Combate
A dungeon utiliza o **mesmo sistema de batalha do jogo**.

Inclui:

- ATK
- DEF
- SPEED
- Esquiva
- Uso de poções
- Transformações
- Equipamentos

---

# 6. Restrição Durante Dungeon
Enquanto estiver em dungeon ativa o jogador **não poderá**:

- Caçada
- Câmara do Ki Ancestral
- Outra Dungeon

Mas **poderá**:

- Continuar dungeon
- Usar itens
- Participar do Big Boss semanal

---

# 7. Big Boss Durante Dungeon
O **Big Boss semanal permanece acessível** mesmo com dungeon ativa.

Fluxo:

Dungeon ativa  
↓  
Player acessa Big Boss  
↓  
Após ataque pode voltar para dungeon  

---

# 8. Condição de Derrota
A dungeon falha quando:

HP dos 3 jogadores = 0

Resultado:

Dungeon Falhou  
Sem recompensa do boss  

---

# 9. Condição de Vitória
A dungeon é concluída quando:

Boss final é derrotado

---

# 10. Sistema de Recompensas
A recompensa do boss é baseada no **dano causado por cada jogador**.

Exemplo:

Boss HP: 100000

Dano causado:

Player A: 50000  
Player B: 30000  
Player C: 20000  

Recompensas distribuídas proporcionalmente.

---

# 11. Tipos de Recompensa

Possíveis recompensas:

- Gold
- XP
- Equipamentos
- Materiais
- Ki Crystals
- Dungeon Keys

---

# 12. Persistência (Banco de Dados)

### dungeon_runs

id  
leader_id  
status  
max_level  
current_wave  
boss_alive  
created_at  

---

### dungeon_players

id  
dungeon_id  
player_id  
damage_done  
current_hp  
is_alive  

---

### dungeon_enemies

id  
dungeon_id  
enemy_id  
hp  
wave  
is_boss  

---

# 13. Fluxo Completo

Menu  
↓  
Dungeon  
↓  
Criar Dungeon  
↓  
Consumir chave  
↓  
Convidar jogadores  
↓  
Iniciar  
↓  
Combater waves  
↓  
Boss final  
↓  
Distribuição de recompensas  

---

# 14. Ideias Futuras

## Dungeons por região

Caverna Vulcânica — Volcano Lands  
Templo Congelado — Frozen Peaks  
Ruínas Sombrias — Swamp of Shadows  

---

## Chaves raras

Common Key — normal  
Ancient Key — elite  
Mythic Key — boss raro  

---

# 15. Balanceamento Recomendado

Players: 1–3  
Waves: 4  
Boss: 1  
Tempo médio: 5–10 minutos  

Ideal para gameplay em **Telegram RPG**.

---

# Conclusão
O sistema de Dungeons adiciona:

- Cooperação entre jogadores
- Conteúdo PvE estruturado
- Recompensas especiais
- Progressão cooperativa

Mantendo compatibilidade com a arquitetura atual do **DragonKi**.
