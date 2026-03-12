
# DragonKi — Plano de Ação para Escalabilidade no Telegram
Guia em etapas para preparar o jogo para crescimento de jogadores

---

# 1. Objetivo

Este documento define um plano prático para preparar o DragonKi para crescer com segurança dentro do ecossistema do Telegram, reduzindo riscos de:

- rate limit do bot
- gargalos de banco de dados
- lentidão em batalhas
- excesso de mensagens
- falhas de processamento em picos de uso

A ideia é evoluir a arquitetura por etapas, sem complexidade desnecessária no início.

---

# 2. Visão Geral dos Riscos

Os principais pontos de atenção para um jogo no Telegram são:

- limite de envio de mensagens do bot
- excesso de mensagens por usuário
- uso de long polling em escala
- muitas leituras e escritas no banco
- processamento concorrente de batalhas e eventos
- crescimento de rankings, bosses e sistemas globais

---

# 3. Estratégia Geral

A estratégia recomendada é evoluir o DragonKi em 4 fases:

## Fase 1 — Base sólida
Preparar o bot para sair do estágio inicial sem gargalos simples.

## Fase 2 — Escala moderada
Suportar crescimento para centenas e poucos milhares de jogadores ativos.

## Fase 3 — Escala avançada
Preparar o jogo para eventos globais, big boss e picos de uso.

## Fase 4 — Escala alta
Arquitetura pronta para dezenas de milhares de jogadores.

---

# 4. Plano de Ação por Etapas

## ETAPA 1 — Trocar Long Polling por Webhook
### Problema
Long polling pode se tornar ineficiente conforme o número de jogadores cresce.

### Ação
Migrar o bot para **Webhook**.

### Benefícios
- menor latência
- menor consumo de recursos
- melhor escalabilidade
- mais compatível com ambiente de produção

### Prioridade
**Alta**

### Resultado esperado
O bot passa a receber eventos do Telegram de forma mais eficiente.

---

## ETAPA 2 — Criar Fila de Envio de Mensagens
### Problema
Telegram possui limite aproximado de envio de mensagens por segundo por bot.

### Ação
Criar uma **message queue interna** para o envio das respostas do bot.

### Estrutura sugerida

Game Logic  
↓  
Message Queue  
↓  
Telegram Sender Worker

### Benefícios
- evita erro 429
- controla bursts de mensagens
- melhora estabilidade do bot

### Prioridade
**Muito Alta**

### Resultado esperado
O bot passa a responder com controle de taxa e não quebra em picos.

---

## ETAPA 3 — Consolidar Respostas do Jogo
### Problema
Muitas mensagens pequenas por ação aumentam o risco de spam e rate limit.

### Ação
Agrupar resultados em mensagens únicas.

### Exemplo
Em vez de enviar:

- XP ganho
- Gold ganho
- item dropado
- dano causado

Enviar:

```
⚔️ Vitória!

+120 XP
+50 Gold
🎁 Item obtido: Espada de Ki
``` 

### Benefícios
- menos mensagens por usuário
- melhor UX
- menor carga no Telegram

### Prioridade
**Alta**

---

## ETAPA 4 — Separar Camada de Processamento e Envio
### Problema
Lógica do jogo misturada com envio de mensagens dificulta escala.

### Ação
Separar em duas camadas:

- **Game Engine**
- **Delivery Layer (Telegram)**

### Benefícios
- código mais limpo
- mais fácil testar
- mais fácil evoluir para app/web no futuro

### Prioridade
**Alta**

---

## ETAPA 5 — Otimizar Leituras do Banco
### Problema
Cada batalha e evento faz muitas operações no banco.

### Ação
Mapear queries críticas e reduzir leituras repetidas.

### Foco inicial
- player
- inventário
- equipamentos
- monstros
- ranking
- big boss

### Benefícios
- menos latência
- menos carga no Postgres
- melhor tempo de resposta

### Prioridade
**Muito Alta**

---

## ETAPA 6 — Introduzir Cache para Dados Quentes
### Problema
Algumas informações são consultadas repetidamente.

### Ação
Adicionar cache para dados de leitura frequente.

### Dados candidatos a cache
- configuração global
- classes
- monstros
- itens
- regiões
- regras de drop
- vip_config

### Tecnologia sugerida
- Redis (futuro)
- cache em memória (início)

### Benefícios
- menos queries
- respostas mais rápidas
- menor carga no banco

### Prioridade
**Média**

---

## ETAPA 7 — Garantir Segurança de Concorrência
### Problema
Eventos globais e ações repetidas podem causar race conditions.

### Casos críticos
- big boss
- compra de item
- venda de item
- recompensa de indicação
- compra de VIP
- conclusão de câmara
- drop de recompensa

### Ação
Aplicar:
- transações
- lock lógico
- constraints únicas
- validações idempotentes

### Benefícios
- evita duplicação
- evita inconsistência de economia
- protege progressão

### Prioridade
**Muito Alta**

---

## ETAPA 8 — Criar Workers para Processos Pesados
### Problema
Algumas tarefas não precisam bloquear a resposta do usuário.

### Ação
Mover tarefas assíncronas para workers.

### Exemplos
- envio de notificações
- limpeza de registros
- backup
- consolidação de ranking
- encerramento de eventos
- spawn de big boss
- expiração de sistemas temporários

### Benefícios
- bot mais rápido
- menos timeout
- melhor separação de responsabilidades

### Prioridade
**Alta**

---

## ETAPA 9 — Monitorar Métricas Técnicas
### Problema
Sem observabilidade, o crescimento vira tentativa e erro.

### Ação
Criar monitoramento mínimo.

### Métricas importantes
- tempo médio de resposta
- mensagens por segundo
- erros 429
- queries lentas
- número de batalhas por minuto
- uso de CPU / memória
- falhas no webhook

### Benefícios
- decisões baseadas em dados
- identificação rápida de gargalos

### Prioridade
**Alta**

---

## ETAPA 10 — Preparar Banco para Crescimento
### Problema
Com mais jogadores, tabelas como battles, logs e ranking crescem rápido.

### Ação
Aplicar boas práticas de banco:

- índices corretos
- tabelas de histórico separadas
- paginação
- limpeza/arquivamento de logs
- retenção de auditoria definida

### Tabelas críticas
- players
- inventory
- battles
- transactions
- big_boss_attacks
- chamber_runs
- referrals

### Benefícios
- melhor desempenho de longo prazo
- menos risco de degradação progressiva

### Prioridade
**Muito Alta**

---

## ETAPA 11 — Planejar Infra para Eventos Globais
### Problema
Big Boss, rankings e eventos criam picos reais de uso.

### Ação
Criar arquitetura específica para módulos globais.

### Recomendado
- service dedicado para big boss
- tabelas agregadas de ranking
- jobs de consolidação
- limite por usuário por dia
- queries otimizadas por evento

### Benefícios
- picos mais previsíveis
- menor risco de travar o jogo todo

### Prioridade
**Alta**

---

## ETAPA 12 — Estruturar Backup e Recuperação
### Problema
Quanto mais o jogo cresce, maior o impacto de perder dados.

### Ação
Implementar:
- backup diário automático
- storage externo
- retenção de backups
- rotina documentada de restore

### Benefícios
- proteção do progresso dos jogadores
- segurança da economia do jogo

### Prioridade
**Muito Alta**

---

# 5. Ordem Recomendada de Implementação

## Curto Prazo
1. Migrar para Webhook
2. Criar fila de mensagens
3. Consolidar respostas do bot
4. Revisar queries críticas
5. Garantir transações e constraints

## Médio Prazo
6. Workers assíncronos
7. Cache para dados quentes
8. Monitoramento técnico
9. Otimizações de banco

## Longo Prazo
10. Infra específica para eventos globais
11. Redis
12. shard lógico / separação de serviços se necessário

---

# 6. Meta por Faixa de Crescimento

## Até 500 jogadores ativos
- webhook
- fila de mensagens
- respostas consolidadas
- queries revisadas

## Até 2.000 jogadores ativos
- workers
- cache
- monitoramento
- índices fortes no banco

## Até 10.000 jogadores ativos
- serviços separados para eventos
- Redis
- agregação de rankings
- planejamento de scaling horizontal

---

# 7. Tecnologias Recomendadas

## Hoje
- Go
- Railway
- PostgreSQL
- Webhook Telegram

## Próximas
- Redis
- Cloudflare R2
- worker de background
- observabilidade (logs + métricas)

---

# 8. Resultado Esperado

Seguindo esse plano, o DragonKi ficará preparado para:

- crescer com segurança
- suportar picos de uso
- manter boa experiência para os jogadores
- proteger economia e progressão
- evoluir para produto com escala real

---

# 9. Resumo Executivo

Os passos mais importantes para agora são:

- **Webhook**
- **fila de mensagens**
- **respostas consolidadas**
- **otimização de banco**
- **controle de concorrência**
- **backup automático**

Esses 6 pontos já colocam o DragonKi em um nível muito mais profissional para crescimento dentro do Telegram.

---

DragonKi — Telegram Scalability Action Plan
