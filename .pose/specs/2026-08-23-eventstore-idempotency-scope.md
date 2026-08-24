---
slug: eventstore-idempotency-scope
status: done
created_at: 2026-08-23
completed_at: 2026-08-23
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
- [x] Confirmar que nenhum consumidor atual depende do comportamento cross-stream (grep por reuse de request_id entre sessões) — todo chamador em `internal/mcpserver` passa `request_id` por chamada já escopada a um `session_id`/stream; nenhum reusa deliberadamente entre streams.

### Implementation
- [x] Trocar a chave de `s.seen` para `streamID + "\x00" + requestID` (ou struct composta) — helper `seenKey`.
- [x] Atualizar `Open`/`ReadEvents`'s reconstrução do cache para usar a mesma chave composta.

### Validation
- [x] Teste replicando o cenário de colisão (duas streams, mesmo request_id) confirmando isolamento.
- [x] Suíte completa de internal/eventstore, internal/session, internal/assistance, internal/assessment sem regressão.

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
- `go test ./internal/eventstore/... -run TestAppendIdempotencyIsScopedPerStream -v` antes da correção → FAIL, reproduzindo a colisão exatamente como descrita (2026-08-23).
- Correção aplicada em `internal/eventstore/store.go` (helper `seenKey`, três pontos de uso: replay em `Open`, checagem em `Append`, gravação em `Append`).
- `go build ./...` → ok (2026-08-23).
- `gofmt -l internal/eventstore` → sem saída (2026-08-23).
- `go vet ./internal/eventstore/...` → sem diagnósticos (2026-08-23).
- `go test ./internal/eventstore/... ./internal/session/... ./internal/assistance/... ./internal/assessment/... -race -v` → `ok` em todos os pacotes, sem data races (2026-08-23).
- `govulncheck ./...` → "No vulnerabilities found." (2026-08-23).

### Results summary
Cache de idempotência de `Append` corrigido para ser indexado por
`(stream_id, request_id)` via `seenKey(streamID, requestID)` em vez de
`request_id` global. Teste de regressão
`TestAppendIdempotencyIsScopedPerStream` reproduz a colisão (falhava antes,
passa depois) e confirma que uma retentativa real dentro da mesma stream
continua idempotente. Nenhuma mudança de formato persistido; `Open`/
`ReadEvents` reconstroem o cache corrigido a partir de logs já gravados sem
migração.

### Requirement trace
- R1 [satisfied] internal/eventstore/store.go (seenKey usado em Append) + test:TestAppendIdempotencyIsScopedPerStream.
- R2 [satisfied] internal/eventstore/store.go (Append) + test:TestAppendIsIdempotentByRequestID, test:TestAppendIdempotencyIsScopedPerStream (retry dentro de ses_2).
- R3 [satisfied] test:TestAppendIdempotencyIsScopedPerStream (duas streams, mesmo request_id, IDs de evento distintos).
- R4 [satisfied] internal/eventstore/store.go (Open) + test:TestOpenRecoversRevisionsAndIdempotencyFromExistingLog (não alterado, continua passando com a chave composta).

### Known gaps
- Nenhum.

## 7. Final Report

### Delivered scope
Correção do cache de idempotência de `eventstore.Store.Append`, agora
escopado por `(stream_id, request_id)`. Nenhuma mudança de contrato público
nem de formato persistido.

### Files and modules changed
- internal/eventstore/store.go (helper `seenKey`, três pontos de uso).
- internal/eventstore/store_test.go (teste de regressão `TestAppendIdempotencyIsScopedPerStream`).

### Validation executed
- Command: go test ./internal/eventstore/... ./internal/session/... ./internal/assistance/... ./internal/assessment/... -race
- Result: ok em todos os pacotes.
- Command: go build ./... && gofmt -l internal/eventstore && go vet ./internal/eventstore/...
- Result: sem erros/diagnósticos.
- Command: govulncheck ./...
- Result: No vulnerabilities found.

### Residual risks
- Nenhum adicional além do já descrito em Technical risks.

### Follow-ups
- [wont-do: correção completa e coberta por regressão; nenhum trabalho residual identificado]
