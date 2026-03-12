# k6 Load Tests

Estrutura inicial de stress/load test para a arquitetura atual do DragonKi.

Alvos atuais:

- `cmd/bot` -> `POST /webhook/telegram`
- `cmd/webhook` -> `POST /webhook/mercadopago`

## Requisitos

- `k6` instalado
- bot e webhook rodando localmente
- Postgres e Redis locais via `docker compose up -d`

## Variaveis uteis

- `BASE_URL` default `http://localhost:8080`
- `BotObservabilityUrl` default `http://localhost:8080`
- `WebhookObservabilityUrl` default `http://localhost:8080`
- `TELEGRAM_CHAT_ID` default `200`
- `TELEGRAM_USER_ID` default `200`
- `TELEGRAM_MESSAGE_ID` default `1000`
- `MONSTER_ID` default `1`
- `MP_WEBHOOK_TOKEN` para cenarios do Mercado Pago
- `TELEGRAM_DELIVERY_MODE=noop` no `cmd/bot` para cenarios locais de carga via webhook

## Exemplos

Instalar `k6` no Windows:

```powershell
powershell -ExecutionPolicy Bypass -File tests/load/k6/install_k6_windows.ps1
```

Forcar `winget`:

```powershell
powershell -ExecutionPolicy Bypass -File tests/load/k6/install_k6_windows.ps1 -ForceWinget
```

Forcar Chocolatey:

```powershell
powershell -ExecutionPolicy Bypass -File tests/load/k6/install_k6_windows.ps1 -ForceChoco
```

Bot Telegram com callback de batalha:

```bash
k6 run tests/load/k6/scenarios/telegram_battle.js
```

Bot Telegram com Big Boss:

```bash
k6 run tests/load/k6/scenarios/telegram_big_boss.js
```

Webhook Mercado Pago:

```bash
k6 run tests/load/k6/scenarios/mercadopago_webhook.js
```

Runner baseline com captura de `/debug/vars` antes e depois:

```powershell
powershell -ExecutionPolicy Bypass -File tests/load/k6/run_baseline.ps1 `
  -Scenario telegram_battle `
  -BaseUrl http://localhost:8080 `
  -BotObservabilityUrl http://localhost:8080 `
  -WebhookObservabilityUrl http://localhost:8081 `
  -WorkerObservabilityUrl http://localhost:9091
```

Saida gerada por execucao:

- `k6_stdout.txt`
- `k6_summary.json`
- `bot_debug_vars_before.json`
- `bot_debug_vars_after.json`
- `webhook_debug_vars_before.json`
- `webhook_debug_vars_after.json`
- `worker_debug_vars_before.json`
- `worker_debug_vars_after.json`

Comparar snapshots da execucao:

```powershell
powershell -ExecutionPolicy Bypass -File tests/load/k6/compare_debug_vars.ps1 `
  -RunDir tests/load/k6/results/telegram_battle-20260312-120000
```

Arquivos gerados pelo comparador:

- `debug_vars_diff.csv`
- `debug_vars_diff.json`

Checklist de leitura:

- [analysis_checklist.md](C:/Users/vigmi/OneDrive/Documentos/develop/dragonki/dragonki_bot/tests/load/k6/analysis_checklist.md)
- [LOCAL_BASELINE_PROTOCOL.md](C:/Users/vigmi/OneDrive/Documentos/develop/dragonki/dragonki_bot/tests/load/k6/LOCAL_BASELINE_PROTOCOL.md)

Observacao:

- `BASE_URL` e a URL usada pelo cenario do `k6`
- `BotObservabilityUrl` e `WebhookObservabilityUrl` podem apontar para processos diferentes
- o worker e lido por `WORKER_OBSERVABILITY_ADDR`, default `http://localhost:9091`
- se um endpoint nao estiver disponivel, o snapshot correspondente salva o erro em JSON

Runner completo com testes Go + baseline de carga:

```powershell
powershell -ExecutionPolicy Bypass -File tests/load/k6/run_validation_suite.ps1 `
  -Scenario telegram_battle `
  -BaseUrl http://localhost:8080 `
  -BotObservabilityUrl http://localhost:8080 `
  -WebhookObservabilityUrl http://localhost:8081 `
  -WorkerObservabilityUrl http://localhost:9091
```

Para incluir as integracoes Go no mesmo fluxo:

```powershell
powershell -ExecutionPolicy Bypass -File tests/load/k6/run_validation_suite.ps1 `
  -Scenario telegram_battle `
  -BaseUrl http://localhost:8080 `
  -BotObservabilityUrl http://localhost:8080 `
  -WebhookObservabilityUrl http://localhost:8081 `
  -WorkerObservabilityUrl http://localhost:9091 `
  -RunIntegration
```

## Observabilidade recomendada durante a execucao

- logs do `cmd/bot`
- logs do `cmd/webhook`
- `pprof`
- metricas da aplicacao
- `pg_stat_statements`
- `redis-cli INFO`

Os scripts daqui sao baselines. Eles nao substituem testes concorrentes em Go para fluxos economicos e de consistencia.
