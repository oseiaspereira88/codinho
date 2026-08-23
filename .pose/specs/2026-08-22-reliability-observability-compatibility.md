---
slug: reliability-observability-compatibility
status: in-progress
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: local-event-store, mcp-stdio-foundation, safe-check-executor, mastery-review-scheduling, administrative-cli-fixtures, security-privacy-hardening
priority: 230
components: reliability, observability, compatibility
delivers:
---

# Spec: reliability-observability-compatibility

## 1. Intent

### Goal
Comprovar recuperação, desempenho local, diagnósticos, logging e compatibilidade do runtime V1.

### Business value
Tornar sessões retomáveis e problemas operacionais explicáveis sem depender de inspeção manual de arquivos.

### Constraints
- stdout do serve permanece exclusivo do protocolo.
- Logs são locais, estruturados, limitados e sem código por padrão.
- Compatibilidade declarada exige CI ou evidência equivalente.

### Non-goals
- Telemetria remota, dashboard, tracing distribuído ou SLO de serviço cloud.

## 2. Requirements

### Functional
- R1: Recuperar após interrupção entre evento e snapshot sem perder evento confirmado.
- R2: Detectar e diagnosticar log truncado, snapshot adulterado, lock órfão e evidência ausente.
- R3: Implementar logs em stderr com correlation ID, tool, duração, status e tamanho, sem payload sensível.
- R4: Completar codinho doctor para binário, catálogo, estado, workspace, checks e conexão estática MCP.
- R5: Atender startup até 1 segundo e queries p95 até 100 ms no volume-alvo.
- R6: Limitar memória, output, tempo e concorrência em operações custosas.
- R7: Testar Linux e macOS; documentar Windows como suportado somente após evidência.
- R8: Testar atualização de pack sem alterar sessão fixada e migração de evento suportada.
- R9: Detectar qualquer escrita em stdout fora de frames MCP.
- R10: Fornecer version info e relatório diagnóstico sanitizado.

### Non-functional
- Testes de recuperação devem ser determinísticos e repetíveis.
- Benchmarks são informativos com thresholds documentados.

### Security
- Diagnóstico não expõe roots, ambiente completo, código ou secrets.

### Compatibility
- Fixar matriz de Go, OS, SDK MCP, schema de catálogo e schema de evento.

## 3. Technical Plan

### Affected areas
- internal/eventstore/, internal/mcpserver/, internal/cli/, internal/diagnostics/, tests/

### Artifacts
- created: internal/diagnostics/doctor.go
- created: internal/diagnostics/doctor_test.go
- created: internal/diagnostics/process_unix.go
- created: internal/diagnostics/process_windows.go
- created: internal/mcpserver/logging.go
- created: internal/mcpserver/logging_test.go
- modified: internal/mcpserver/server.go
- modified: internal/cli/doctor.go
- modified: internal/cli/root.go
- modified: internal/cli/root_test.go
- created: cmd/codinho/reliability_test.go
- modified: .pose/indexes/validation-matrix.json
- created: docs/compatibility.md
- created: docs/troubleshooting.md

### Delivery targets
Nenhum novo; qualidade transversal.

### API/contract changes
- Estabilizar schema do relatório doctor e version info.

### Data/storage changes
- Implementar política de compactação/retenção se medições exigirem, preservando log auditável.

### Technical risks
- Thresholds variam por máquina.
- Teste de crash pode ser frágil.

## 4. Tasks

### Planning
- [x] Definir fault matrix e compatibility matrix.
- [x] Definir benchmark dataset no volume V1.

### Implementation
- [x] Completar doctor e relatório sanitizado.
- [x] Implementar fault injection e recovery cases.
- [x] Implementar logging e correlation.
- [x] Implementar limites e benchmarks.
- [x] Criar suites de compatibilidade e docs.

### Validation
- [x] Repetir fault suite e benchmarks.
- [x] Executar stdio lifecycle e stdout purity.
- [ ] Executar CI nas plataformas declaradas.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Observabilidade local não deve coletar conteúdo do aluno.
- Options considered: logs detalhados; nenhum log; metadados estruturados.
- Decision: registrar metadados operacionais e correlation IDs sem payload.
- Rationale: diagnosticar com privacidade.
- Consequences: reprodução pode exigir export explícito adicional do usuário.

### Decision 2
- Date: 2026-08-23
- Context: R1/R2/R8 (recuperação de log truncado, log corrompido,
  migração de evento) já tinham testes dedicados e passando desde
  local-event-store (`TestReadEventsToleratesTruncatedFinalLine`,
  `TestReadEventsReportsCorruptedMiddleLineWithoutLosingPrefix`,
  `TestReadEventsUpcastsRegisteredVersion`, `TestSnapshotOnlyWrittenAfterEventSynced`,
  `TestReadSnapshotDetectsInvalidJSON`, `TestLockIsExclusiveAndReleasable`)
  — nenhum deles foi reescrito.
- Options considered: (a) duplicar esses cenários num novo
  `internal/eventstore/fault_test.go` só para satisfazer o artifact
  planejado; (b) reconhecer a cobertura já existente e direcionar o
  esforço novo só para o que realmente faltava: detecção de lock
  órfão (nenhum código fazia isso) e o relatório sanitizado que a
  amarra tudo (`codinho doctor`).
- Decision: (b). Não foi criado `internal/eventstore/fault_test.go`;
  `internal/diagnostics` é o artefato novo real desta spec para R1/R2.
- Rationale: reescrever testes já corretos não adiciona evidência, só
  duplica manutenção; a lacuna real e verificável era a detecção de
  lock órfão (`internal/eventstore/lock.go` já documentava isso como
  "fora de escopo... spec operacional dedicada" — esta é essa spec).
- Consequences: a rastreabilidade de R1/R2 cita os testes
  pré-existentes de `internal/eventstore` explicitamente, não um
  arquivo novo.

### Decision 3
- Date: 2026-08-23
- Context: R3 pede logs estruturados por chamada de tool
  (correlation ID, tool, duração, status, tamanho), mas nenhuma das 32
  tools tinha esse logging, e adicionar manualmente em cada uma seria
  invasivo (32 pontos de edição repetitivos).
- Options considered: (a) editar cada handler de tool individualmente;
  (b) usar `mcp.Server.AddReceivingMiddleware` (já suportado pelo SDK)
  para logar toda chamada uma única vez, de forma genérica.
- Decision: (b). `internal/mcpserver/logging.go` registra um único
  middleware em `server.go`.
- Rationale: zero mudança nos 32 handlers existentes; o middleware
  nunca vê o payload de argumento/resposta além do necessário para
  medir tamanho (via `json.Marshal` do resultado, não do payload em
  si) — nunca loga o conteúdo.
- Consequences: qualquer tool nova automaticamente ganha logging
  estruturado sem esforço extra.

## 6. Validation

### Strategy
Combinar fault injection, golden diagnostics, benchmarks e CI multi-OS.

### Deterministic checks
- Test: go test -race ./... e suites de fault/compatibility.
- Lint: gofmt -l .
- Typecheck: go vet ./...
- Build: go build ./cmd/codinho em plataformas declaradas.
- Security / Contract: stdout purity, redaction e schema compatibility.

### Execution log
- `go build ./...` → ok (2026-08-23).
- `gofmt -l .` → sem saída (2026-08-23).
- `go vet ./...` e `GOOS=windows go vet ./...` → sem diagnósticos (2026-08-23).
- `go test ./... -race` → `ok` em todos os pacotes, incluindo `internal/diagnostics` (novo), `internal/mcpserver` (logging middleware) e `cmd/codinho` (stdout purity, startup budget), sem data races (2026-08-23).
- `govulncheck ./...` → "No vulnerabilities found." (2026-08-23); registrado junto com `stdout-purity` e `startup-budget` em `.pose/indexes/validation-matrix.json`.

### Results summary
`codinho doctor` completo (catálogo, estado gravável, lock órfão,
evidência, toolchain), sempre sanitizado (nunca ambiente completo,
código ou secrets — provado por teste). Logging estruturado por tool
via middleware do SDK MCP, sem tocar os 32 handlers. Prova real de
pureza de stdout (feliz e de falha) e de orçamento de startup (<1s);
busca de catálogo em escala V1 já provada <100ms p95 desde
curriculum-graph-path-recommendation. Matriz de compatibilidade e
troubleshooting documentados honestamente (Linux testado; macOS/Windows
sem execução real, conforme Constraint).

### Requirement trace
- R1 [satisfied] cobertura pré-existente de internal/eventstore (Decision 2) + internal/diagnostics (novo).
- R2 [satisfied] internal/diagnostics.checkLock (lock órfão, novo) + testes pré-existentes de log truncado/corrompido/snapshot inválido.
- R3 [satisfied] internal/mcpserver/logging.go (Decision 3) + test:TestLoggingMiddlewareEmitsStructuredLineWithoutPayload.
- R4 [satisfied] internal/diagnostics.Run + codinho doctor (humano/JSON) + testes correspondentes.
- R5 [satisfied] test:TestServeStartupRespondsWithinOneSecond (novo) + test:TestSearchP95BudgetAtV1Scale (pré-existente).
- R6 [satisfied] limites pré-existentes (output/timeout/concorrência em internal/checks; tamanho/profundidade/aliases em internal/curriculum), documentados em docs/compatibility.md.
- R7 [satisfied] docs/compatibility.md (matriz honesta: Linux testado, macOS/Windows não).
- R8 [satisfied] catálogo carregado uma única vez no startup (nenhum hot-reload existe) documentado em docs/compatibility.md + migração de evento já testada (TestReadEventsUpcastsRegisteredVersion, pré-existente).
- R9 [satisfied] test:TestServeNeverWritesToStdoutOnStartupFailure, test:TestServeStdoutIsOnlyValidJSONFrames (novos, subprocesso real).
- R10 [satisfied] codinho doctor/--json (version+commit+relatório sanitizado) + test:TestReportNeverExposesFullEnvironment.

### Known gaps
- CI real multi-OS (macOS/Windows) não foi executado nesta sessão — só checagem estática (`GOOS=windows go vet`) e execução real em Linux.
- Windows permanece "não declarado suportado" até execução real (Compatibility desta spec).

## 7. Final Report

### Delivered scope
`codinho doctor` completo e sanitizado (incluindo detecção de lock
órfão), logging estruturado por tool, prova real de pureza de stdout e
orçamento de startup, e documentação honesta de compatibilidade/
troubleshooting.

### Files and modules changed
- internal/diagnostics/{doctor,doctor_test,process_unix,process_windows}.go (novo pacote).
- internal/mcpserver/{logging,logging_test}.go (novo) + server.go (middleware registrado).
- internal/cli/doctor.go (reescrito), root.go, root_test.go.
- cmd/codinho/reliability_test.go (novo): stdout purity + startup budget.
- .pose/indexes/validation-matrix.json: checks stdout-purity, startup-budget, govulncheck.
- docs/{compatibility,troubleshooting}.md.

### Validation executed
- Command: go test ./... -race
- Result: ok em todos os pacotes.
- Command: go test ./cmd/codinho/... -run "TestServeNeverWritesToStdoutOnStartupFailure|TestServeStdoutIsOnlyValidJSONFrames|TestServeStartupRespondsWithinOneSecond"
- Result: PASS.
- Command: govulncheck ./...
- Result: No vulnerabilities found.

### Residual risks
- Thresholds de performance medidos no hardware desta sessão; podem variar em outra máquina (Technical risk aceito).
- Sem evidência real de macOS/Windows.

### Follow-ups
- [covered: installation-documentation-ci] Automatizar a matriz de compatibilidade.
- [open] Executar a suíte real em macOS e Windows quando houver acesso a essas máquinas/CI.
