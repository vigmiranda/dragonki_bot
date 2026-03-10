# DragonKi Bot

Bot de RPG para Telegram com backend em Go, Postgres e Redis.

## Stack
- Go 1.22
- Postgres
- Redis (fila assincrona)
- Telegram Bot API
- GORM
- Asynq (jobs no Redis)

## Servicos do projeto
- `cmd/bot`: servico principal (webhook Telegram + regras de jogo)
- `cmd/webhook`: webhook de pagamentos (Mercado Pago)
- `cmd/worker`: consumidor de jobs assincronos (premiacao Big Boss)

## Funcionalidades importantes
- Feature flags em banco (`feature_flags`)
- Inventario/equipamentos normalizados
- Loja com compra e venda de equipamentos
- Sistema de referral (indicacao)
  - Link de convite por jogador (`/start ref_<player_id>`)
  - Vinculo apenas no onboarding inicial
  - Recompensa quando indicado atinge nivel minimo
  - Protecoes: sem autoindicacao e sem duplicidade
- Big Boss semanal
  - Geracao semanal automatica
  - Ataque diario (1x por dia)
  - Ranking por dano
  - Premiacao Top 3
  - Premiacao assincrona via fila Redis

## Estrutura de dados (Big Boss)
- `weekly_bosses`: instancias semanais
- `weekly_boss_attacks`: historico de ataques diarios
- `weekly_boss_contributions`: dano acumulado por jogador
- `weekly_boss_templates`: catalogo dinamico de bosses (sorteio)
- `weekly_boss_rewards`: configuracao de premiacao por posicao
- `weekly_boss_reward_items`: itens da premiacao
- `weekly_boss_payouts`: auditoria/idempotencia dos pagamentos

## Estrutura de dados (Referral)
- `referrals`
  - `referrer_player_id`: player que convidou
  - `referred_player_id`: player convidado
  - `status`: `PENDING` | `REWARDED`
  - `required_level`: nivel minimo para liberar recompensa
  - `completed_at`, `rewarded_at`: rastreio de conclusao/pagamento

## Variaveis de ambiente
Base:
- `APP_ENV`
- `DATABASE_URL`
- `BOT_TOKEN`
- `RUN_MIGRATIONS`
- `RUN_AUTO_MIGRATE`
- `MIGRATIONS_DIR`
- `MAX_PLAYER_LEVEL` (default: `50`)

Fila / Worker:
- `REDIS_URL` (ex: `redis://localhost:6379`)
- `BIG_BOSS_QUEUE_ENABLED` (`true|false`)
- `BIG_BOSS_QUEUE_NAME` (default: `big_boss_rewards`)
- `QUEUE_CONCURRENCY` (default: `10`)

Referral:
- `REFERRAL_REQUIRED_LEVEL` (default: `7`)
- `REFERRAL_REWARD_GOLD` (default: `1000`)
- `REFERRAL_REWARD_KI` (default: `1`)

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
- Criar servico separado no mesmo repo
- Dockerfile path: `Dockerfile.worker`
- Start command: vazio (usar `CMD` do Dockerfile)
- Configurar `DATABASE_URL`, `REDIS_URL`, `QUEUE_CONCURRENCY`, `APP_ENV`

### Webhook pagamentos
- Dockerfile: `Dockerfile.webhook`

## Notas operacionais
- O scheduler do Big Boss roda no `cmd/bot`.
- A premiacao Top 3 e enfileirada ao derrotar o boss.
- `weekly_boss_payouts` evita duplicidade de pagamento.
- O menu de referral aparece em `?? Convites` no menu principal.
- O vinculo do convite ocorre no `/start` com `ref_<player_id>`.
- O pagamento do referral ocorre no evento de level up da batalha.

## Testes

```bash
go test ./internal/application ./internal/adapters/postgres ./internal/adapters/telegram ./cmd/bot ./cmd/worker
```

Teste de integracao Postgres com validacao real de persistencia:

```bash
$env:TEST_DATABASE_URL="postgres://dragonki:dragonki@localhost:5432/postgres?sslmode=disable"
go test ./internal/adapters/postgres -run Integration
```

O teste cria um banco temporario, aplica as migrations, valida escrita/leitura real e remove o banco ao final.

## Git Hooks

Para rodar testes automaticamente antes de cada `git push`, configure o repositorio para usar os hooks versionados:

```bash
git config core.hooksPath .githooks
```

O hook `pre-push` executa:

```bash
go test ./...
```

E, se `TEST_DATABASE_URL` estiver definida, tambem executa:

```bash
go test ./internal/adapters/postgres -run Integration -count=1
```

Exemplo de ambiente local para incluir as integracoes no `pre-push`:

```bash
$env:TEST_DATABASE_URL="postgres://dragonki:dragonki@localhost:5432/postgres?sslmode=disable"
git push
```

Por seguranca, o hook e os testes de integracao recusam executar se `TEST_DATABASE_URL` apontar para um host que nao pareca local/efemero, como `localhost`, `127.0.0.1`, `postgres`, `dragonki_postgres`, `*.local` ou `*.internal`.
