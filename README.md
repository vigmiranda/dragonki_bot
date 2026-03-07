# 🐉 DragonKi — Telegram RPG Bot

**DragonKi** é um RPG por texto para Telegram inspirado em **Dungeons & Dragons** e **Dragon Ball**, combinando fantasia clássica, progressão de poder, energia (stamina) e transformações baseadas em Ki.

Cada jogador possui seu próprio personagem, evolução persistente em banco de dados e ações limitadas por energia que se regenera com o tempo.

---

## 🎮 Funcionalidades

- 🤖 Bot de Telegram multiplayer
- 🧙 Personagem individual por usuário
- ⚔️ Sistema de batalhas PvE
- ⚡ Energia limitada com recarga automática por tempo
- 🔥 Sistema de transformação (Ki Awakened)
- 🏆 Progressão por nível e XP
- 💾 Persistência com SQLite
- 🐳 Pronto para Docker
- ☁️ Deploy gratuito (Railway)

---

## 🐉 Lore Resumido

No mundo de **Aethoria**, os Dragões Primordiais libertaram uma energia ancestral conhecida como **Ki**.  
Aqueles capazes de dominá-la se tornam Portadores do Ki, guerreiros capazes de evoluir, despertar e transcender seus limites.

Seu destino é ascender.

---
## 🧩 NPC`s
- 🙎‍♂️ Sir Soneca - Dono da pousada
---
## 🧩 Comandos do Jogo

| Comando        | Descrição |
|---------------|----------|
| `/start`      | Cria o personagem |
| `/status`     | Mostra status completo |
| `/batalhar`   | Luta contra um monstro (consome energia) |
| `/transform`  | Ativa transformação de Ki (nível 3+) |

---

## ⚔️ Classes (v1)

> Atualmente o personagem inicia como **Blade Monk**  
> Sistema de escolha de classes será expandido na v2.

**Blade Monk**
- Especialista em Ki e combate corpo-a-corpo
- Focado em progressão e transformações

---

## 🔥 Sistema de Transformação

### ⚡ Ki Awakened
- Desbloqueia no nível 3
- Custo: 3 de energia
- Duração: 3 batalhas
- Bônus:
    - +50% ATK
    - +20% XP

Transformações são temporárias e não acumulam.

---

## ⚡ Sistema de Energia

- Energia máxima: **10**
- Recupera **1 ponto a cada 10 minutos**
- Batalhas consomem energia
- Transformações consomem energia

A recarga é calculada automaticamente com base no tempo, sem uso de cron jobs.

---

## 🧠 Tecnologias Utilizadas

- **Go 1.21+**
- **Telegram Bot API**
- **PostgreSQL**
- **Docker**
- **Railway (Deploy) -> https://railway.app/** 

---
## ⚙️ Configuração de Ambiente

O projeto carrega variáveis nesta ordem:

1. Variáveis já definidas no sistema (não são sobrescritas)
2. `.env`
3. `.env.<APP_ENV>`

Se `APP_ENV` não for definido, o padrão é `development`.

Arquivos de referência:

- `.env.example`
- `.env.development`
- `.env.staging`
- `.env.production`

Variáveis críticas:

- `DATABASE_URL`
- `BOT_TOKEN`
- `MP_ACCESS_TOKEN`
- `MP_NOTIFICATION_URL`
- `MP_WEBHOOK_TOKEN`

Deploy seguro:

- `RUN_MIGRATIONS=true` somente quando quiser executar migração automática.
- Em `staging`/`production`, manter `RUN_MIGRATIONS=false` por padrão.
- Defina `MP_WEBHOOK_TOKEN` para proteger o endpoint de webhook.
- Feature flags agora são lidas da tabela `public.feature_flags` no Postgres.
- Para desativar compra de Ki sem novo deploy:
  - `UPDATE public.feature_flags SET enabled = FALSE, updated_at = NOW() WHERE name = 'PAYMENTS';`
- Para aplicar seeds SQL manualmente (ex.: atualizar `monster_drops`):
  - `go run ./cmd/seed`
  - Opcional: `SQL_SEED_FILE=internal/adapters/postgres/SQL/Inserts/monster_drops.sql go run ./cmd/seed`

---

## 📁 Estrutura do Projeto

```text
dragonki/
├── main.go
├── go.mod
├── Dockerfile
└── database.db (gerado automaticamente)

## Database Migrations (golang-migrate)

Variaveis:
- RUN_MIGRATIONS: aplica migrations versionadas no boot do bot quando true.
- RUN_AUTO_MIGRATE: controla AutoMigrate do GORM (recomendado true apenas em development).
- MIGRATIONS_DIR: diretorio das migrations (padrao: db/migrations).

Comandos:
- go run ./cmd/migrate -action up
- go run ./cmd/migrate -action version
- go run ./cmd/migrate -action down -steps 1
- go run ./cmd/migrate -action force -version 1
