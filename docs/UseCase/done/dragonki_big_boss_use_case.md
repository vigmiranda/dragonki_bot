📜 DragonKi — Caso de Uso Oficial

👑 Módulo: Big Boss Semanal

Descrição

O módulo Big Boss representa um evento semanal onde um inimigo extremamente poderoso aparece no mundo de DragonKi.

Este inimigo:
 • não aparece nos mapas tradicionais
 • possui HP global compartilhado
 • pode ser atacado por todos os jogadores
 • cada jogador pode atacar apenas 1 vez por dia e descontará 5 de energia e 200 de gold do player
 • o dano causado é resultado de uma batalha completa
 • o boss permanece ativo até domingo ou até seu HP chegar a zero

Durante a semana, o sistema registra o dano total causado por cada jogador, permitindo gerar um ranking semanal de contribuição.

⸻

🎮 Atores
 • Player
 • Sistema DragonKi
 • Big Boss semanal

⸻

🧭 Acesso no jogo

Menu principal:

👑 Big Boss

Tela exibida:

👑 Big Boss da Semana

Nome: Surtr, o Devorador do Mundo
HP: 1.450.000 / 2.000.000

Tempo restante: 3d 12h

Seu dano total: 8.340
Ataque diário: disponível

[⚔️ Atacar]
[🏆 Ranking]


⸻

📜 Regras de Negócio

1️⃣ Existe apenas um Big Boss ativo por semana

Cada semana possui uma instância única de boss global.

Exemplo:

Week: 2026-W12
Boss: Surtr


⸻

2️⃣ O boss possui HP global compartilhado

Todos os jogadores reduzem o mesmo HP.

HP restante = HP máximo - soma de todos os danos


⸻

3️⃣ Cada jogador pode atacar apenas 1 vez por dia

Regra de bloqueio:

weekly_boss_id + player_id + attack_date

Se já atacou:

❌ Você já enfrentou o Big Boss hoje.
Volte amanhã para atacar novamente.


⸻

4️⃣ O combate segue o sistema normal de batalha

O combate usa o mesmo motor de combate do jogo.

A batalha ocorre em rounds:

player ataca
boss contra-ataca
player ataca
boss contra-ataca

até:

player.hp <= 0


⸻

5️⃣ Regra de esquiva baseada em Speed

Antes do boss causar dano, ocorre uma rolagem de esquiva.

Fórmula usada no jogo:

dodgeChance := float64(player.Speed) * 0.5

if dodgeChance > 40 {
    dodgeChance = 40
}

roll := rand.Float64() * 100

Se:

roll < dodgeChance

Resultado:

💨 Player esquiva do ataque

O boss não causa dano naquele turno.

⸻

6️⃣ A batalha termina quando o HP do player chega a zero

O boss não morre na batalha individual.

Ele apenas acumula o dano causado.

⸻

7️⃣ O dano total da batalha é registrado

Durante o combate:

totalDamage += damage

Após o combate:

boss.current_hp -= totalDamage


⸻

8️⃣ O dano do jogador é acumulado durante a semana

Cada jogador possui um registro:

player_total_damage

Que representa:

soma de todos os ataques da semana


⸻

9️⃣ O boss permanece ativo até domingo

O evento possui:

start_at
end_at

Mesmo que não morra, o evento encerra automaticamente.

⸻

🔟 O boss pode ser derrotado antes do domingo

Se:

current_hp <= 0

então:

status = defeated
defeated_at = now

E novos ataques são bloqueados.

⸻

⚔️ Fluxo principal do caso de uso

Caso de uso: Atacar Big Boss

1️⃣ Jogador acessa menu Big Boss

Sistema carrega:
 • boss ativo
 • HP restante
 • tempo restante
 • dano acumulado do jogador
 • status do ataque diário

⸻

2️⃣ Jogador seleciona “Atacar”

Sistema valida:

boss ativo
evento não expirado
boss não derrotado
player não atacou hoje


⸻

3️⃣ Sistema inicia batalha

Inicializa:

playerHP = player.HP
totalDamage = 0


⸻

4️⃣ Loop de combate

for playerHP > 0 {

    // Player ataca
    damage := player.ATK
    totalDamage += damage

    // Boss tenta atacar
    dodgeChance := float64(player.Speed) * 0.5

    if dodgeChance > 40 {
        dodgeChance = 40
    }

    roll := rand.Float64() * 100

    if roll < dodgeChance {
        continue
    }

    bossAtk := boss.ATK - player.DEF

    if bossAtk < 1 {
        bossAtk = 1
    }

    playerHP -= bossAtk
}


⸻

5️⃣ Final da batalha

Player é derrotado.

O dano total é calculado.

⸻

6️⃣ Atualizar HP global do boss

boss.current_hp -= totalDamage

Se ficar negativo:

boss.current_hp = 0


⸻

7️⃣ Registrar ataque

Salvar registro:

weekly_boss_attack

Campos:

player_id
boss_id
attack_date
damage


⸻

8️⃣ Atualizar contribuição semanal

player_total_damage += damage


⸻

9️⃣ Retornar resultado

Mensagem exibida:

⚔️ Você enfrentou Surtr!

💥 Dano causado: 3.420
[06/03/2026 17:07] Vitor: 🩸 HP restante do Boss: 1.245.700

💀 Você foi derrotado.

📊 Seu dano total nesta semana: 12.840


⸻

🏆 Ranking semanal

Ranking ordenado por:

total_damage DESC

Exemplo:

🏆 Ranking Big Boss

1️⃣ Vitor — 42.900
2️⃣ Kael — 38.120
3️⃣ Arion — 34.500
4️⃣ Lina — 29.880


⸻

👑 MVP semanal

Ao final da semana:

player com maior dano total

Exibição:

👑 MVP da Semana

Vitor
Dano total: 42.900


⸻

🗄️ Estrutura de Dados

WeeklyBoss

type WeeklyBoss struct {
 ID int64
 Name string
 WeekKey string

 MaxHP int64
 CurrentHP int64

 Status string

 StartsAt time.Time
 EndsAt time.Time
 DefeatedAt *time.Time
}


⸻

WeeklyBossAttack

type WeeklyBossAttack struct {
 ID int64
 BossID int64
 PlayerID int64

 AttackDate string
 Damage int64

 CreatedAt time.Time
}


⸻

WeeklyBossContribution

type WeeklyBossContribution struct {
 ID int64
 BossID int64
 PlayerID int64

 TotalDamage int64

 LastAttackAt time.Time
}


⸻

🔐 Constraints importantes

1 ataque por dia

unique (boss_id, player_id, attack_date)


⸻

1 contribuição por jogador

unique (boss_id, player_id)


⸻

📦 Serviços envolvidos

CombatEngine
WeeklyBossService
WeeklyBossRankingService
WeeklyBossRepository
WeeklyBossAttackRepository


⸻