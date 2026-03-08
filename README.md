# DragonKi Bot

Bot de RPG para Telegram com backend em Go, Postgres e Redis.

## Stack
- Go 1.22
- Postgres
- Redis (fila assíncrona)
- Telegram Bot API
- GORM
- Asynq (jobs no Redis)

## Serviços do projeto
- `cmd/bot`: serviço principal (webhook Telegram + regras de jogo)
- `cmd/webhook`: webhook de pagamentos (Mercado Pago)
- `cmd/worker`: consumidor de jobs assíncronos (premiação Big Boss)

## Funcionalidades importantes
- Feature flags em banco (`feature_flags`)
- Inventário/equipamentos normalizados
- Loja com compra e venda de equipamentos
- Big Boss semanal
  - geração semanal automática
  - ataque diário (1x por dia)
  - ranking por dano
  - premiação Top 3
  - premiação assíncrona via fila Redis

## Estrutura de dados (Big Boss)
- `weekly_bosses`: instâncias semanais
- `weekly_boss_attacks`: histórico de ataques diários
- `weekly_boss_contributions`: dano acumulado por jogador
- `weekly_boss_templates`: catálogo dinâmico de bosses (sorteio)
- `weekly_boss_rewards`: configuração de premiação por posição
- `weekly_boss_reward_items`: itens da premiação
- `weekly_boss_payouts`: auditoria/idempotência dos pagamentos

## Variáveis de ambiente
Base:
- `APP_ENV`
- `DATABASE_URL`
- `BOT_TOKEN`
- `RUN_MIGRATIONS`
- `RUN_AUTO_MIGRATE`
- `MIGRATIONS_DIR`

Fila / Worker:
- `REDIS_URL` (ex: `redis://localhost:6379`)
- `BIG_BOSS_QUEUE_ENABLED` (`true|false`)
- `BIG_BOSS_QUEUE_NAME` (default: `big_boss_rewards`)
- `QUEUE_CONCURRENCY` (default: `10`)

Pagamentos:
- `MP_ACCESS_TOKEN`
- `MP_NOTIFICATION_URL`
- `MP_WEBHOOK_TOKEN`
- `KI_CRYSTALS_VALUE`
- `KI_CRYSTALS_MULT`
- `KI_CRYSTALS_CONVERT`

## Rodando local (Docker)
Subir Postgres + Redis:

```bash
docker compose up -d
```

## Rodando local (Go)
Aplicar migrations (uma vez):

```bash
$env:RUN_MIGRATIONS="true"
go run ./cmd/bot
```

Bot:

```bash
$env:RUN_MIGRATIONS="false"
$env:BIG_BOSS_QUEUE_ENABLED="true"
$env:REDIS_URL="redis://localhost:6379"
go run ./cmd/bot
```

Worker:

```bash
$env:REDIS_URL="redis://localhost:6379"
$env:QUEUE_CONCURRENCY="10"
go run ./cmd/worker
```

## Deploy Railway
### Bot
- Dockerfile: `Dockerfile.bot`
- Habilitar `BIG_BOSS_QUEUE_ENABLED=true` e `REDIS_URL`

### Worker
- Criar serviço separado no mesmo repo
- Dockerfile path: `Dockerfile.worker`
- Start command: vazio (usar `CMD` do Dockerfile)
- Configurar `DATABASE_URL`, `REDIS_URL`, `QUEUE_CONCURRENCY`, `APP_ENV`

### Webhook pagamentos
- Dockerfile: `Dockerfile.webhook`

## Notas operacionais
- O scheduler do Big Boss roda no `cmd/bot`.
- A premiação Top 3 é enfileirada ao derrotar o boss.
- `weekly_boss_payouts` evita duplicidade de pagamento.

## Testes

```bash
go test ./internal/application ./internal/adapters/postgres ./internal/adapters/telegram ./cmd/bot ./cmd/worker
```
