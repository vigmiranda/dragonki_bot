# DragonKi Bot - Mapa de Manutencao

## Objetivo
Este documento descreve onde alterar o codigo com seguranca, quais pontos tem maior risco de regressao, e como validar mudancas antes de deploy.

## 1) Mapa rapido do sistema

### Entradas executaveis
- `cmd/bot/main.go`: servico principal (webhook Telegram, regras de jogo, scheduler do Big Boss).
- `cmd/webhook/main.go`: webhook de pagamentos Mercado Pago.
- `cmd/worker/main.go`: worker Asynq para jobs assincronos (premiacao Big Boss).
- `cmd/migrate/main.go`: controle de migrations (`up`, `down`, `version`, `force`).
- `cmd/seed/main.go`: aplicacao de seed SQL.

### Camadas e responsabilidades
- `internal/domain`: entidades e modelos de negocio.
- `internal/application`: regras de negocio e casos de uso.
- `internal/ports`: contratos (interfaces) usados pela aplicacao.
- `internal/adapters`:
  - `telegram`: interface com Telegram.
  - `postgres`: persistencia.
  - `mercadopago`: gateway de pagamento.
  - `queue`: Redis/Asynq.
- `internal/bootstrap`: carga de env, feature flags, migrations.
- `db/migrations`: evolucao de schema.

## 2) Fluxos principais

### Telegram (core do jogo)
1. Recebe update em `/webhook/telegram` (`cmd/bot/main.go`).
2. `Handler.HandleUpdate` identifica mensagem/callback (`internal/adapters/telegram/handler.go`).
3. `handleCallback` roteia por `action` (`internal/adapters/telegram/callbacks.go`).
4. `handle_*.go` chama service em `internal/application`.
5. Service persiste via interfaces de `internal/ports`, implementadas em `internal/adapters/postgres`.

### Pagamentos
1. Usuario aciona compra no Telegram (`action_buy_ki`).
2. `PaymentService` cria PIX via gateway e salva pagamento pendente.
3. Mercado Pago envia webhook para `/webhook/mercadopago`.
4. `MercadoPagoWebhook` valida token e confirma pagamento.
5. Repositorio liquida pagamento com idempotencia.

### Big Boss semanal
1. `WeeklyBossService` garante boss da semana (`EnsureCurrentWeeklyBoss`).
2. Jogador ataca 1x por dia (energia/gold).
3. Ao derrotar boss, sistema enfileira distribuicao de premio Top 3.
4. Worker processa job e aplica payout com auditoria/idempotencia.

## 3) Onde mexer por tipo de mudanca

### Nova acao no bot (botao/menu)
- Editar menu em `internal/adapters/telegram/menus.go`.
- Adicionar roteamento em `internal/adapters/telegram/callbacks.go`.
- Implementar comportamento no `handle_*.go` apropriado.
- Se houver regra nova, criar/ajustar service em `internal/application`.

### Nova regra de negocio
- Criar/editar service em `internal/application`.
- Se precisar de acesso a dados novo: atualizar interface em `internal/ports`.
- Implementar contrato no adapter (`internal/adapters/postgres` ou outro).
- Cobrir com testes de service.

### Mudanca de banco
- Criar migration `up/down` em `db/migrations`.
- Ajustar modelos/queries/repositorios.
- Validar `cmd/migrate` localmente.

### Nova feature com rollout gradual
- Usar `bootstrap.IsFeatureEnabled(...)` nos pontos de entrada.
- Garantir consistencia: menu, handler e webhook/servico devem respeitar a mesma flag.

### Novo job assincrono
- Definir payload e task type na aplicacao.
- Enfileirar no service.
- Registrar handler no worker (Asynq).
- Tratar idempotencia e retries.

## 4) Hotspots de regressao

### Acoes Telegram por string
- Risco: typo em `action` quebra fluxo e cai em fallback.
- Mitigacao: revisar menu + callback + handler no mesmo PR.

### Idempotencia e concorrencia
- Big Boss: ataque diario unico, payout unico por posicao/jogador.
- Pagamentos: webhook pode repetir notificacao.
- Mitigacao: manter guardas no banco + checks de reprocessamento.

### Flags e configuracao
- Flags inconsistentes podem expor opcao no menu sem backend habilitado.
- Env ausente pode derrubar servico no startup (`DATABASE_URL`, `BOT_TOKEN`, `REDIS_URL`).

### Mudancas cross-layer
- Alterar `ports` sem atualizar adapter quebra compilacao ou runtime.
- Mitigacao: tratar alteracao como "pacote fechado" (service + port + adapter + teste).

## 5) Checklist de PR (obrigatorio)
- [ ] Mudanca limitada ao escopo (sem side effects ocultos).
- [ ] Menu/callback/handler sincronizados (quando for Telegram UI).
- [ ] Services com teste de unidade cobrindo caminho feliz + erro.
- [ ] Repositorios cobertos em caso de query/regra critica.
- [ ] Migration `up/down` criada e validada (quando houver schema).
- [ ] Feature flag aplicada em todos os pontos de entrada relevantes.
- [ ] Idempotencia validada em pagamentos/eventos/jobs.
- [ ] Variaveis de ambiente novas documentadas no `README.md`/`.env.example`.
- [ ] `go test` dos pacotes impactados passou.

## 6) Script de validacao minima antes de merge
Executar no repo:

```bash
go test ./internal/application ./internal/adapters/postgres ./internal/adapters/telegram ./cmd/bot ./cmd/worker
```

Se houve alteracao em pagamentos, incluir `./cmd/webhook`.

## 7) Template de rollout (producao)

### Pre-deploy
- [ ] Revisar flags relacionadas (ex.: `payments`, `big_boss_weekly`, `ancestral_ki_chamber`).
- [ ] Confirmar env vars obrigatorias no ambiente alvo.
- [ ] Rodar migrations necessarias.
- [ ] Definir estrategia de rollback (flag off, rollback de migration, revert de deploy).

### Deploy
- [ ] Subir bot.
- [ ] Subir webhook (se pagamento foi alterado).
- [ ] Subir worker (se fila/job foi alterado).
- [ ] Verificar health basico via logs de inicializacao.

### Pos-deploy (primeiros 30 min)
- [ ] Conferir logs de erro no handler Telegram.
- [ ] Validar 1 fluxo funcional ponta a ponta da feature alterada.
- [ ] Em features de pagamento/evento, monitorar duplicidade e retries.
- [ ] Se houver anomalia, desabilitar flag e executar rollback planejado.

## 8) Matriz rapida de ownership por area
- Gameplay (batalha, energia, shop, inventario): `internal/application/*service.go` + handlers Telegram.
- Eventos (Big Boss, Camara): services dedicados + menus/handlers + repositorios.
- Pagamentos: `payment_service`, webhook Mercado Pago, repositorio de pagamentos.
- Persistencia/SQL: adapters Postgres + migrations.
- Operacao: `cmd/*`, envs, flags, Redis queue, worker.

## 9) Melhorias recomendadas (backlog tecnico)
- Centralizar constantes de `action` e `menu key` para reduzir erro de typo.
- Padronizar logging (`log.Printf`) em vez de mistura com `fmt.Printf`.
- Adicionar testes de integracao para fluxos criticos (Big Boss e webhook de pagamento).
- Criar runbook de incidentes com consultas SQL prontas para auditoria de payout/payment.
