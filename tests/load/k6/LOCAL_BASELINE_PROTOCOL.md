# Protocolo Local de Baseline

Este protocolo serve para a primeira rodada real de carga no ambiente local.

## 1. Preparar infraestrutura

Subir Postgres e Redis:

```powershell
docker compose up -d
```

## 2. Configurar ambiente

Bot:

```powershell
$env:DATABASE_URL="postgres://dragonki:dragonki@localhost:5432/dragonki?sslmode=disable"
$env:REDIS_URL="redis://localhost:6379"
$env:RUN_MIGRATIONS="true"
$env:BIG_BOSS_QUEUE_ENABLED="true"
$env:TELEGRAM_DELIVERY_MODE="noop"
$env:OBSERVABILITY_ENABLED="true"
$env:PPROF_ENABLED="true"
go run ./cmd/bot
```

Webhook de pagamento em outro terminal:

```powershell
$env:DATABASE_URL="postgres://dragonki:dragonki@localhost:5432/dragonki?sslmode=disable"
$env:OBSERVABILITY_ENABLED="true"
$env:PPROF_ENABLED="true"
go run ./cmd/webhook
```

Worker em outro terminal:

```powershell
$env:DATABASE_URL="postgres://dragonki:dragonki@localhost:5432/dragonki?sslmode=disable"
$env:REDIS_URL="redis://localhost:6379"
$env:OBSERVABILITY_ENABLED="true"
$env:PPROF_ENABLED="true"
$env:WORKER_OBSERVABILITY_ADDR=":9091"
go run ./cmd/worker
```

## 3. Verificacao rapida antes da carga

Checar observabilidade:

```powershell
Invoke-WebRequest http://localhost:8080/debug/vars -UseBasicParsing
Invoke-WebRequest http://localhost:8081/debug/vars -UseBasicParsing
Invoke-WebRequest http://localhost:9091/debug/vars -UseBasicParsing
```

Checar `pprof`:

```powershell
Invoke-WebRequest http://localhost:8080/debug/pprof/ -UseBasicParsing
Invoke-WebRequest http://localhost:8081/debug/pprof/ -UseBasicParsing
Invoke-WebRequest http://localhost:9091/debug/pprof/ -UseBasicParsing
```

## 4. Ordem recomendada de cenarios

Rodar nesta ordem:

1. `telegram_battle`
2. `telegram_big_boss`
3. `mercadopago_webhook`

Motivo:

- batalha cria o baseline principal de gameplay
- Big Boss pressiona concorrencia e dedupe
- Mercado Pago valida webhook e settle sem misturar com gameplay

## 5. Primeira rodada recomendada

Primeiro cenario:

```powershell
powershell -ExecutionPolicy Bypass -File tests/load/k6/run_baseline.ps1 `
  -Scenario telegram_battle `
  -BaseUrl http://localhost:8080 `
  -BotObservabilityUrl http://localhost:8080 `
  -WebhookObservabilityUrl http://localhost:8081 `
  -WorkerObservabilityUrl http://localhost:9091
```

Depois comparar:

```powershell
powershell -ExecutionPolicy Bypass -File tests/load/k6/compare_debug_vars.ps1 `
  -RunDir tests/load/k6/results/telegram_battle-YYYYMMDD-HHMMSS
```

## 6. O que observar na primeira rodada

Olhar primeiro:

- `k6_summary.json`
- `debug_vars_diff.csv`
- `bot_debug_vars_after.json`

Sinais de alerta iniciais:

- `http_req_failed` acima de `1%`
- `http_req_duration p95` acima de `300ms`
- crescimento forte de `db.query.slow_total`
- crescimento inesperado de `telegram.callback.duplicate_*`

Observacao importante:

- para baseline local do Telegram, mantenha `TELEGRAM_DELIVERY_MODE=noop`
- isso evita `chat not found` e remove o custo da API real do Telegram da medicao

## 7. Rodadas seguintes

Se a primeira rodada estiver estavel:

1. aumentar stages/VUs no script alvo
2. repetir `telegram_battle`
3. rodar `telegram_big_boss`
4. por fim rodar `mercadopago_webhook`

Nao mude dois cenarios ao mesmo tempo na fase inicial. O objetivo e gerar baseline comparavel.
