# Checklist de Analise de Baseline

Use este checklist depois de cada rodada executada por `run_baseline.ps1`.

## 1. Sanidade da execucao

- confirmar que `k6_summary.json` existe
- confirmar que `bot_debug_vars_before.json` e `bot_debug_vars_after.json` existem
- confirmar que `webhook_debug_vars_before.json` e `webhook_debug_vars_after.json` existem
- confirmar que `worker_debug_vars_before.json` e `worker_debug_vars_after.json` existem
- confirmar se algum snapshot salvou erro em vez de metricas

## 2. Saida do k6

Verificar no `k6_summary.json`:

- `http_req_failed`
- `http_req_duration`
- `iterations`
- `vus_max`
- taxa de checks falhos

Perguntas:

- o erro HTTP ficou abaixo da meta?
- p95 e p99 ficaram dentro do esperado?
- o throughput bateu o alvo da rodada?

## 3. Bot

Comparar `bot_debug_vars_before.json` e `bot_debug_vars_after.json`.

Olhar principalmente:

- `total_requests`
- `routes./webhook/telegram.request_count`
- `routes./webhook/telegram.status_count`
- `routes./webhook/telegram.duration_ms_avg`
- `counters.telegram.callback.received_total`
- `counters.telegram.callback.duplicate_local_total`
- `counters.telegram.callback.duplicate_receipt_total`
- `counters.telegram.callback.duplicate_semantic_total`
- `counters.db.query.total`
- `counters.db.query.error_total`
- `counters.db.query.slow_total`

Perguntas:

- a rota do Telegram concentrou os requests esperados?
- houve pico anormal de dedupe?
- erros de banco apareceram?
- queries lentas cresceram demais para a rodada?

## 4. Webhook Mercado Pago

Comparar `webhook_debug_vars_before.json` e `webhook_debug_vars_after.json`.

Olhar:

- `routes./webhook/mercadopago.request_count`
- `routes./webhook/mercadopago.status_count`
- `counters.payments.webhook.received_total`
- `counters.payments.webhook.confirmed_total`
- `counters.payments.webhook.confirm_error_total`
- `counters.payments.webhook.unauthorized_total`
- `counters.payments.confirm.attempt_total`
- `counters.payments.confirm.settled_total`
- `counters.payments.confirm.idempotent_total`

Perguntas:

- houve retry excessivo do parceiro?
- confirmacao virou gargalo?
- chamadas idempotentes cresceram como esperado ou em excesso?

## 5. Worker

Comparar `worker_debug_vars_before.json` e `worker_debug_vars_after.json`.

Olhar:

- `counters.worker.big_boss.enqueue_total`
- `counters.worker.big_boss.enqueue_error_total`
- `counters.worker.big_boss.enqueue_idempotent_total`
- `counters.worker.big_boss.enqueue_duration_ms_total`
- `counters.worker.big_boss.job_received_total`
- `counters.worker.big_boss.job_success_total`
- `counters.worker.big_boss.job_error_total`
- `counters.worker.big_boss.job_retried_total`
- `counters.worker.big_boss.job_retry_attempt_total`
- `counters.worker.big_boss.job_duration_ms_total`

Perguntas:

- o worker drenou o volume esperado?
- houve retry real?
- o tempo medio por job parece aceitavel?

## 6. Banco

Cruzar snapshots com Postgres:

- `pg_stat_statements`
- top queries por tempo total
- top queries por tempo medio
- quantidade de chamadas

Perguntas:

- a query mais cara bate com o fluxo exercitado?
- o gargalo e CPU de app ou SQL?

## 7. Redis e fila

Cruzar snapshots com Redis/Asynq:

- `redis-cli INFO`
- tamanho das filas
- jobs retry/scheduled/active

Perguntas:

- o backlog subiu e desceu?
- houve fila acumulada apos o fim da carga?

## 8. Conclusao da rodada

Registrar:

- cenario executado
- configuracao de VUs/stages
- p95/p99
- counters anormais
- queries suspeitas
- proximo ajuste recomendado

Sem esse registro, as proximas rodadas viram observacao ad hoc e perdem comparabilidade.
