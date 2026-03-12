# DragonKi - Plano Implementavel de Stress Test

Este plano adapta os testes de carga para a arquitetura real do DragonKi em 2026-03-11.
O projeto nao expoe endpoints HTTP de gameplay como `/api/battle`. Hoje a superficie real e:

- `cmd/bot` -> `POST /webhook/telegram`
- `cmd/webhook` -> `POST /webhook/mercadopago`
- `cmd/worker` -> processamento assincrono em Redis/Asynq

Por isso, os stress tests precisam simular updates do Telegram, webhooks do Mercado Pago e drenagem de fila, alem de observar Postgres, Redis e o proprio processo Go.

---

# 1. Objetivo

Validar que o DragonKi suporta crescimento de usuarios sem perder consistencia em:

- batalhas e callbacks do Telegram
- Big Boss semanal
- operacoes economicas e de inventario
- webhooks de pagamento
- processamento assincrono no worker

Os testes devem encontrar:

- gargalos de CPU e memoria
- saturacao de conexoes no Postgres
- pressao no Redis e fila
- regressao de latencia
- erros HTTP e retries
- falhas de idempotencia e concorrencia

---

# 2. Ferramenta escolhida

Ferramenta principal: `k6`

Motivos:

- cobre bem webhooks HTTP
- suporta cenarios com payload dinamico
- tem thresholds objetivos
- funciona bem em CI/CD

Ferramentas complementares:

- `go test` com cenarios concorrentes para service/repo
- `pprof` para CPU e heap
- metricas de processo e de aplicacao
- `pg_stat_statements` e `EXPLAIN ANALYZE` para banco
- `redis-cli INFO` e metricas do Asynq para fila

---

# 3. Escopo real de carga

## 3.1 Bot Telegram

Alvo:

- `POST /webhook/telegram`

Carga realista:

- updates de callback: batalha, Big Boss, inventario, loja, camara ancestral
- updates de mensagem: onboarding e `/start ref_<id>`

Riscos principais:

- callbacks duplicados
- dedupe excessivo ou insuficiente
- latencia no render de painel
- contencao no banco

## 3.2 Webhook Mercado Pago

Alvo:

- `POST /webhook/mercadopago`

Riscos principais:

- retries do parceiro
- settle parcial
- erro sob burst de notificacoes

## 3.3 Worker Big Boss

Alvo:

- fila Redis/Asynq do `cmd/worker`

Riscos principais:

- backlog crescendo sem drenagem
- tempo alto de processamento por job
- reprocessamento duplicado

## 3.4 Testes concorrentes fora do HTTP

Alguns fluxos importantes devem continuar sendo validados por testes concorrentes em Go, nao apenas por `k6`:

- compra/venda de inventario
- equipar/desequipar
- payout do Big Boss
- referral reward
- settle de pagamento

---

# 4. Cenarios prioritarios

## 4.1 Telegram callback de batalha

Simular callbacks como `action_fight?monsterID=<id>` no `/webhook/telegram`.

Objetivo:

- medir throughput da rota principal de gameplay
- observar latencia do bot sob carga sustentada
- validar consistencia de XP, gold, HP e logs de batalha

Perfis:

- load: `100`, `200`, `500` VUs
- stress: `800`, `1000+` VUs
- spike: `0 -> 500` callbacks em poucos segundos

Validacoes:

- erro HTTP < 1%
- p95 < 300ms no alvo inicial
- ausencia de escrita parcial no banco

## 4.2 Telegram callback de Big Boss

Simular `big_boss_attack` no `/webhook/telegram`.

Objetivo:

- validar concorrencia no boss semanal
- detectar HP inconsistente
- validar limite diario e ranking

Validacoes:

- HP nunca negativo
- contribuicoes coerentes com o total de dano
- callbacks duplicados absorvidos por dedupe

## 4.3 Telegram callback de inventario e loja

Simular:

- `buy_item`
- `sell_item`
- `use_item`
- `equip_item`
- `unequip_item`

Objetivo:

- pressionar operacoes economicas e inventario
- validar ausencia de duplicacao de item e gold

Validacoes:

- saldo economico consistente
- nenhum item vendido e equipado ao mesmo tempo
- rollback correto em erro

## 4.4 Telegram callback da Camara Ancestral

Simular:

- `start_ancestral_ki_run`
- `ancestral_ki_chamber_continue`
- `abandon_ancestral_ki_run`

Objetivo:

- validar concorrencia de run
- medir custo das resolucoes de batalha e progresso

## 4.5 Webhook Mercado Pago

Simular burst de `payment.updated` no `/webhook/mercadopago`.

Objetivo:

- validar throughput do webhook
- medir idempotencia e settle
- observar fila de retries do parceiro

Validacoes:

- `payment_webhook_audits` gravando corretamente
- pagamentos nao aplicados em duplicidade
- rollback quando player/pagamento for invalido

## 4.6 Worker e fila do Big Boss

Simular enfileiramento massivo de rewards.

Objetivo:

- medir tempo de drenagem
- medir concorrencia efetiva do worker
- validar reprocessamento idempotente

Validacoes:

- backlog reduzindo dentro da janela esperada
- jobs duplicados nao pagam duas vezes

---

# 5. Casos de observabilidade

Os testes de carga devem sempre vir acompanhados de observabilidade minima. Sem isso, o time mede latencia, mas nao entende por que o sistema degradou.

## 5.1 Casos obrigatorios de observabilidade

Durante cada execucao, coletar:

- latencia por rota
- taxa de erro por rota
- throughput por rota
- uso de CPU e memoria do processo Go
- contagem de goroutines
- tempo e taxa de erro de queries do Postgres
- backlog e throughput do Redis/Asynq
- contadores de dedupe de callback por `receipt` e por janela semantica

## 5.2 Casos especificos por fluxo

### Telegram batalha

Observar:

- quantidade de callbacks recebidos
- quantidade de callbacks ignorados por dedupe
- p95/p99 do `/webhook/telegram`
- top queries por tempo total no Postgres

### Big Boss

Observar:

- throughput de ataques por segundo
- lock contention no banco
- crescimento de `weekly_boss_attacks` e `weekly_boss_contributions`
- tempo total para refletir contribuicao no ranking

### Mercado Pago

Observar:

- volume de `payment_webhook_audits`
- taxa de `CONFIRMED`, `CONFIRM_ERROR`, `UNAUTHORIZED`
- tempo de settle

### Worker

Observar:

- profundidade da fila
- tempo medio de execucao do job
- taxa de retry
- quantidade de jobs idempotentes ignorados

## 5.3 Instrumentacao minima recomendada

Antes dos testes pesados, implementar ou habilitar:

- `/debug/pprof` em bot e webhook
- endpoint `/metrics` ou alternativa equivalente
- logs estruturados para dedupe, queue e settle
- consulta periodica a `pg_stat_statements`
- coleta de `redis INFO`

Sem essa instrumentacao, rodar stress test pesado tende a produzir conclusoes fracas.

---

# 6. Metas iniciais

Metas de partida para a primeira rodada:

- `100` batalhas por segundo equivalentes em callback
- p95 < `300ms` no `/webhook/telegram`
- p99 < `700ms` no `/webhook/telegram`
- erro HTTP < `1%`
- webhook Mercado Pago com p95 < `250ms`
- fila do Big Boss drenando sem backlog indefinido

Metas de crescimento:

- `200+` batalhas por segundo equivalentes
- spike de `500` callbacks sem degradacao severa

---

# 7. Estrutura recomendada no repositorio

Criar e manter:

- `tests/load/k6/lib/`
- `tests/load/k6/scenarios/`
- `tests/load/k6/README.md`

Scripts iniciais:

- `telegram_battle.js`
- `telegram_big_boss.js`
- `mercadopago_webhook.js`

Complementos futuros:

- `telegram_inventory_shop.js`
- `telegram_chamber.js`
- `worker_big_boss_queue.js`

---

# 8. Fases de implementacao

## Fase 1 - Baseline

- criar estrutura `k6`
- criar cenarios iniciais de Telegram batalha, Big Boss e Mercado Pago
- definir thresholds
- documentar execucao local

## Fase 2 - Observabilidade

- expor `pprof`
- expor metricas por rota
- contar dedupe por callback
- documentar coleta de Postgres e Redis

## Fase 3 - Carga moderada

- rodar `load test`
- registrar baseline de p95/p99
- identificar queries criticas

## Fase 4 - Stress e spike

- aumentar VUs
- rodar bursts de Big Boss
- medir comportamento do worker

## Fase 5 - Correcao e repeticao

- otimizar gargalos
- rerodar mesmos cenarios
- comparar baseline antes/depois

---

# 9. Frequencia recomendada

Executar:

- antes de releases grandes
- antes de eventos globais
- antes de mudancas na economia
- antes de novos sistemas concorrentes no Telegram

Rodar baseline menor em CI e cenarios pesados sob demanda.

---

# 10. Proximo passo implementavel

O proximo passo tecnico recomendado no DragonKi e:

1. criar `tests/load/k6` com os tres cenarios base
2. adicionar observabilidade minima em bot e webhook
3. rodar baseline local com Postgres e Redis do `docker-compose`
4. registrar metricas e gargalos encontrados

Esse plano substitui o modelo antigo baseado em endpoints REST de gameplay e fica alinhado com a arquitetura webhook-driven atual do projeto.
