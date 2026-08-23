---
slug: eventstore-idempotency-scope
status: draft
created_at: 2026-08-23
completed_at:
supersedes:
depends_on: local-event-store
priority: 15
components: eventstore
delivers:
---

# Spec: eventstore-idempotency-scope

## 1. Intent

### Goal
Corrigir o cache de idempotência de `eventstore.Store.Append` para ser
indexado por `(stream_id, request_id)` em vez de `request_id` global,
eliminando a possibilidade de duas streams diferentes colidirem ao
reaproveitar o mesmo `request_id`.

### Business value
Idempotência é uma garantia central do event store (local-event-store R4);
uma colisão entre sessões diferentes que reutilizam o mesmo `request_id`
devolveria a um cliente o evento de OUTRA sessão, quebrando essa garantia
silenciosamente — sem erro, sem log, apenas o dado errado.

### Constraints
- `local-event-store` já está `done` e selada; esta correção deve tratar
  seu formato de log persistido como compatível (não é uma migração de
  schema, `request_id` continua uma string opaca por evento).
- A correção não pode quebrar a leitura de logs já gravados por versões
  anteriores do código (replay em `Open`/`ReadEvents`).

### Non-goals
- Mudar o formato serializado do `Event` em si (schema_version permanece).
- Introduzir namespacing de `request_id` no lado do cliente/MCP.

## 2. Requirements

### Functional
- R1: `Append` deve tratar `request_id` como único por stream, não globalmente.
- R2: Uma retentativa com o mesmo `(stream_id, request_id)` deve continuar retornando o evento original cacheado, sem gravar de novo.
- R3: Duas streams diferentes reutilizando o mesmo `request_id` nunca devem observar o evento uma da outra.
- R4: `Open`/`ReadEvents` devem reconstruir o cache corrigido a partir de logs existentes sem exigir migração de arquivo.

### Non-functional
- A correção não pode regredir nenhum teste existente de `internal/eventstore`, `internal/session`, `internal/assistance` ou `internal/assessment` que dependa do comportamento de idempotência atual dentro de uma única stream.

### Compatibility
- Logs JSONL já gravados continuam sendo lidos sem alteração de formato.

## 3. Technical Plan

### Affected areas
- internal/eventstore/

### Artifacts
- modified: internal/eventstore/store.go
- modified: internal/eventstore/store_test.go

### Delivery targets
Nenhum; correção interna sem contrato público novo.

### API/contract changes
Nenhuma mudança de assinatura pública esperada (`Append`/`Open` mantêm a mesma interface).

### Data/storage changes
Nenhuma no formato persistido; apenas o índice em memória (`s.seen`) muda de chave.

### Technical risks
- Retrocompatibilidade com qualquer consumidor que dependa (mesmo que acidentalmente) do comportamento atual de colisão cross-stream.

## 4. Tasks

### Planning
- [ ] Confirmar que nenhum consumidor atual depende do comportamento cross-stream (grep por reuse de request_id entre sessões).

### Implementation
- [ ] Trocar a chave de `s.seen` para `streamID + "\x00" + requestID` (ou struct composta).
- [ ] Atualizar `Open`/`ReadEvents`'s reconstrução do cache para usar a mesma chave composta.

### Validation
- [ ] Teste replicando o cenário de colisão (duas streams, mesmo request_id) confirmando isolamento.
- [ ] Suíte completa de internal/eventstore, internal/session, internal/assistance, internal/assessment sem regressão.

## 5. Decisions

### Decision 1
- Date: 2026-08-23
- Context: achado durante feedback-evaluation-progression ao investigar como `step_evaluate` poderia gravar dois eventos por chamada de forma idempotente.
- Options considered: (a) corrigir agora, tocando evidência selada de local-event-store; (b) reportar e sequenciar como spec própria.
- Decision: (b), a pedido explícito do usuário.
- Rationale: local-event-store já está fechada/selada; uma correção isolada e testável merece seu próprio ciclo de review em vez de reabrir aquele fechamento.
- Consequences: até esta spec ser implementada, clientes MCP devem evitar reutilizar `request_id` entre sessões diferentes (mitigação prática: IDs gerados incluindo o session_id, já é o padrão natural de qualquer cliente razoável).

## 6. Validation

### Strategy
Teste unitário reproduzindo a colisão antes da correção (deve falhar) e confirmando isolamento depois.

### Deterministic checks
- Test: go test ./internal/eventstore/... ./internal/session/... ./internal/assistance/... ./internal/assessment/...
- Lint: gofmt -l internal/eventstore
- Typecheck: go vet ./internal/eventstore/...
- Build: go build ./...

### Execution log
- Pendente.

### Results summary
- Nenhuma correção aplicada ainda; spec criada para sequenciar o trabalho.

### Requirement trace
- Mapear R1–R4 a testes de colisão cross-stream e replay.

### Known gaps
- Nenhum até a implementação começar.

## 7. Final Report

### Delivered scope
Nenhum; spec draft aguardando implementação.

### Files and modules changed
- Planejados em internal/eventstore.

### Validation executed
- Command: pose lint-spec eventstore-idempotency-scope --ready-check
- Result: registrar após validação.

### Residual risks
- Nenhum adicional além do já descrito em Technical risks.

### Follow-ups
- [open]
