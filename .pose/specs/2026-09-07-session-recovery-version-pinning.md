---
slug: session-recovery-version-pinning
status: in-progress
created_at: 2026-09-07
completed_at:
supersedes:
depends_on: eventstore-idempotency-scope, session-orchestration-disclosure, workspace-observation-baselines, reliability-observability-compatibility
priority: 10
components: sessions, eventstore, workspace, mcp-server
delivers: capability:session-recovery
---

# Spec: session-recovery-version-pinning

## 1. Intent

### Goal
Restaurar sessões, conteúdo fixado e escopo de evidências após reinício real do servidor.

### Business value
Fechar uma lacuna verificável da V1 antes de ampliar a superfície que depende dela.

### Constraints
- Preservar autoria do aluno, operação local, JSONL e separação pedagógica.
- Usar o [relatório de auditoria](../reports/2026-09-07-doc-audit-auditoria-do-planejamento-v1.md) como baseline, não como evidência de entrega.
- Consumir knowledge:planning-audit-2026-09 antes de iniciar a implementação.

### Non-goals
- Reabrir ou reescrever atestações históricas de specs done.
- Publicar release, alterar código do aluno ou ampliar a V1 para nuvem/multiusuário.

## 2. Requirements

### Functional
- R1: Após fechar e reabrir codinho serve sobre o mesmo estado, restaurar sessão, política completa, profundidade, nó, avaliações, pistas, detours e revisão confirmados.
- R2: Preservar IDs e idempotência após restart: retry retorna o resultado original; request novo cria sessão distinta; rejeitar reutilização incompatível sem anexar evento.
- R3: Persistir identidade e conteúdo recuperável do catálogo fixado; atualização ou remoção de pack não troca silenciosamente instruções, checks ou critérios de sessões existentes.
- R4: Restaurar baselines, roots autorizadas e vínculo sessão–evidência; evidence_get bloqueia leitura cross-session e informa obsolescência usando o escopo persistido após restart. A autorização no consumo por avaliação exige evaluation-evidence-lineage (Decision 2).
- R5: Validar transições antes do append; falha de domínio não pode deixar evento confirmado que o replay não consegue aplicar. Testar crash entre append e aplicação.
- R6: Ler logs legados sem inventar políticas/conteúdo ausentes; retornar diagnóstico explícito de sessão irrecuperável e permitir novas sessões sem apagar o histórico.
- R7: Provar os cenários por dois processos MCP distintos, incluindo atualização de pack entre processos, corrupção/truncamento e retry após perda da resposta.

### Non-functional
- Manter determinismo e falhas explícitas; provar o caminho composto, não apenas helpers isolados.
- Identificar resultados por versão, commit e cenário; skips obrigatórios não são sucesso.

### Security
- Não copiar estado real do aluno para fixtures, logs ou relatórios.
- Aplicar confinamento de paths, redaction e consentimento nos novos caminhos.

### Compatibility
Documentar em ADR antes da implementação a extensão dos payloads, snapshot do catálogo e política para eventos legados incompletos. Preservar JSONL e interfaces MCP existentes; nenhum banco novo.

## 3. Technical Plan

### Affected areas
- sessions
- eventstore
- workspace
- mcp-server

### Artifacts
- modified: internal/session/service.go
- created: internal/session/recovery.go
- created: internal/session/recovery_test.go
- modified: internal/application/workspace.go
- modified: internal/application/session.go
- modified: internal/application/checks.go
- modified: internal/cli/state.go
- modified: internal/cli/cli_test.go
- created: cmd/codinho/recovery_integration_test.go
- modified: internal/learning/session.go
- modified: internal/session/interview.go
- modified: internal/eventstore/store.go
- modified: internal/eventstore/recovery.go
- modified: internal/eventstore/recovery_test.go
- modified: internal/mcpserver/errors.go
- created: internal/application/workspace_recovery_test.go
- created: .pose/adr/2026-09-06-durable-session-replay-with-pinned-content.md
- created: .pose/knowledge/2026-09-07-decision-log-adr-durable-session-replay-review.md
- modified: docs/compatibility.md
- modified: .pose/indexes/validation-matrix.json
- modified: .pose/roadmaps/codinho-v1.md
- modified: .pose/specs/2026-08-22-v1-integrated-acceptance.md
- created: .pose/specs/2026-09-07-evaluation-evidence-lineage.md
- created: .pose/changelogs/unreleased/session-recovery-version-pinning.md

### Delivery targets
- capability:session-recovery module:cmd/codinho profile:composed-capability entrypoint:cmd/codinho/main.go

### API/contract changes
Documentar em ADR antes da implementação a extensão dos payloads, snapshot do catálogo e política para eventos legados incompletos. Preservar JSONL e interfaces MCP existentes; nenhum banco novo.

### Data/storage changes
Aplicar o [ADR de replay durável](../adr/2026-09-06-durable-session-replay-with-pinned-content.md): início versionado com política/conteúdo fixados, replay de deltas, baselines persistidas e diagnóstico conservador para legado. Envelope JSONL permanece v1.

### Technical risks
Replay parcial pode perder consentimento, atribuir evidência incorretamente ou reutilizar IDs. A ausência de dados em eventos antigos não admite recuperação perfeita por inferência.

## 4. Tasks

### Planning
- [x] Reproduzir o achado e revisar contratos/ADRs aplicáveis.
- [x] Completar decisões de formato e plano de testes negativos antes de modificar código.
- [x] Reconciliar esta lista de artefatos com os arquivos efetivos; declarar arquivos adicionais antes de alterá-los.

### Implementation
- [x] Implementar primeiro o menor fluxo que fecha a lacuna.
- [x] Integrar entradas reais, persistência/compatibilidade e diagnósticos.
- [x] Atualizar documentação e checks declarativos junto com o contrato.

### Validation
- [x] Executar os cenários de cada R-ID, incluindo negativos.
- [ ] Executar pose assess integrate e validação estruturada no candidato.
- [ ] Reconciliar artifacts, surface e revisão independente antes do closeout.

## 5. Decisions

### Decision 1
- Date: 2026-09-07
- Context: O probe MCP de 2026-09-07 UTC iniciou ses_1; após reinício limpo, session_get retornou SESSION_NOT_ACTIVE e um novo session_start retornou STATE_CONFLICT. internal/session.New cria mapas e contador vazios; Start persiste apenas challenge_id/mode e usa CatalogRef vazio. A recuperação de JSONL existente não recompõe os serviços.
- Options considered: deixar implementação implícita no aceite; criar remediação focada.
- Decision: planejar remediação própria e exigir sua entrega antes do aceite.
- Rationale: v1-integrated-acceptance tem non-goal de adicionar features ou correções.
- Consequences: iniciar por reprodução; o ADR de replay durável governa a implementação.

### Decision 2
- Date: 2026-09-07
- Context: Revisão independente verificou que o R4 original pressupunha uma proteção inexistente: EvidenceGet isola leitura, mas StepEvaluate consulta evidências sem validar origem/freshness.
- Options considered: ampliar recuperação para redefinir avaliação; explicitar retrieval e abrir remediação de consumo.
- Decision: delimitar R4 ao contrato real de leitura/freshness e registrar evaluation-evidence-lineage, prioridade 15 e dependência bloqueante do aceite V1.
- Rationale: autorização de evidência qualitativa/externa precisa de política própria; não alegar que persistência resolve consumo.
- Consequences: produto não tem isolamento completo de avaliação; esta pendência continua visível no roadmap e no aceite. Consultar knowledge:adr-durable-session-replay-review.

## 6. Validation

### Strategy
Validar cada contrato com cenário positivo e negativo pela entrada de produção; usar somente dados sintéticos.

| Cenário | Comando obrigatório | Evidência esperada |
|---|---|---|
| Estado, IDs, retries, legado e append inválido | go test ./internal/session -run 'TestRecovery\|TestRejectedTransition\|TestLegacySession' | Estado restaurado; zero evento em chamada inválida |
| Baseline e isolamento após restart | go test ./internal/application -run TestWorkspaceRecovery | Diff preservado, root trocado negado, evidência cross-session negada |
| Cauda truncada seguida de append | go test ./internal/eventstore -run TestRecovery | Prefixo e novos eventos legíveis, backup da cauda |
| Dois processos e pack atualizado | go test ./cmd/codinho -run TestSessionRecoveryOverRealStdio | Sessão antiga preservada e sessão nova usa catálogo novo |
| Regressões concorrentes | go test -race ./... | Todos os pacotes passam |

### Deterministic checks
- Test: go test -race ./internal/session/... ./internal/application/... ./internal/eventstore/... ./cmd/codinho/...
- Lint: pose check --strict; pose lint-spec session-recovery-version-pinning --ready-check
- Typecheck: go vet ./...
- Build: go build ./cmd/codinho
- Security / Contract: pose assess integrate; pose validate --strict --json .pose/results/delivery-validation.json; pose surface-check --spec session-recovery-version-pinning --strict

### Execution log
- 2026-09-07 UTC: criada em revisão de planejamento; implementação e gates de entrega não executados.
- 2026-09-07 UTC: regressões iniciais falharam como esperado: sessão perdida, append de transição inválida e diagnóstico de legado ausente.
- 2026-09-07 UTC: go test -race ./... passou em todos os pacotes; MCP real cobriu update, remoção do desafio, perda da resposta, cauda truncada e diagnóstico legado/corrupção.
- 2026-09-07 UTC: revisão independente corrigiu retry de checks, reserva legacy, freshness runtime, CLI somente leitura, integridade de revisões, conflito atômico de request_id e reparo conservador de tentativa interrompida.
- 2026-09-07 UTC, retomada: go test -race ./... passou; pose check --strict, readiness, skills-check e recurrence-check passaram. assess discover encontrou um módulo; assess tech-debt encontrou zero marcadores. assess integrate retornou zero contratos reconhecidos (limitação conhecida em contributions/20260907-004842-detectar-contratos-mcp-go-e-distinguir-i.md); a integração é comprovada pelos testes MCP reais.
- 2026-09-07 UTC, retomada: primeira validação estruturada executou 11 checks com sucesso e um erro de ambiente: govulncheck instalado em /home/go/go/bin estava ausente do PATH. Corrigir o PATH do processo de validação, preservando a matriz portátil.

### Results summary
Implementação e regressões passaram; fechar somente após validação POSE estruturada, atribuição Git, superfície e revisão selada no candidato final.

### Requirement trace
- R1 [satisfied] test:TestRecoveryRestoresPoliciesHintsDetoursAndHistoricalRetries test:TestRecoveryRestoresGranularityEvaluationAndAdvance
- R2 [satisfied] test:TestRecoveryRestoresSessionAndStartRetry test:TestRecoveryConcurrentRequestReuseIsRejectedAtomically test:TestSessionRecoveryOverRealStdio
- R3 [satisfied] capability:session-recovery evidence:integration check:session-recovery test:TestRecoveryRejectsDamagedPinnedContent
- R4 [satisfied] test:TestWorkspaceRecoveryBaselineScopeAndRetry test:TestWorkspaceRecoveryRejectsReplacedRoot test:TestWorkspaceRecoveryCheckEvidenceAndRetry
- R5 [satisfied] test:TestRejectedTransitionDoesNotAppend test:TestRecoveryCompletesInterruptedSubmission
- R6 [satisfied] test:TestLegacySessionDoesNotBlockNewStarts test:TestRecoveryRejectsDamagedPinnedContent test:TestSessionRecoveryOverRealStdio
- R7 [satisfied] capability:session-recovery evidence:integration check:session-recovery test:TestSessionRecoveryOverRealStdio

### Known gaps
Logs legados incompletos não admitem recuperação segura. Startup exige catálogo/manifest válidos. Após SIGKILL, lock órfão continua sujeito a confirmação operacional antes de liberação; o E2E remove somente seu lock sintético após Wait. Consumo de evidências em avaliação depende de evaluation-evidence-lineage.

## 7. Final Report

### Delivered scope
Recuperação de sessões com conteúdo fixado, idempotência, baselines/evidências, reparo de cauda sob lock e CLI de inspeção somente leitura. Implementado e testado; lifecycle aguarda gates finais.

### Files and modules changed
- internal/session, learning, eventstore, application, cli e mcpserver; E2E real, compatibilidade, ADR e matriz de validação.

### Validation executed
- go test -race ./...: passou em 2026-09-07; validação estruturada POSE e revisão final registradas no closeout.

### Residual risks
Manter baseline de startup sob crescimento do log. Não usar binário anterior como writer do novo estado; não inferir autorização de consumo a partir da existência de uma evidência.

### Follow-ups
- [covered: v1-integrated-acceptance] Reexecutar os requisitos desta spec no candidato composto da V1.
- [open] Confirmar disposition do achado de consumo em evaluation-evidence-lineage com o responsável; owner: @oseiaspereira; due: 2026-09-14. A spec draft já bloqueia o aceite V1.
