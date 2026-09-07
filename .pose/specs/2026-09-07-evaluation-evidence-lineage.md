---
slug: evaluation-evidence-lineage
status: draft
created_at: 2026-09-07
completed_at:
supersedes:
depends_on: session-recovery-version-pinning, safe-check-executor, feedback-evaluation-progression
priority: 15
components: sessions, workspace, mcp-server
delivers: capability:evaluation-evidence-lineage
---

# Spec: evaluation-evidence-lineage

## 1. Intent

### Goal
Validar origem, escopo e atualidade de evidências no consumo por avaliação.

### Business value
Impedir aprovação baseada em check de outra sessão ou de arquivos já alterados.

### Constraints
- Preservar operação local, JSONL e autoria do aluno.
- Consumir knowledge:adr-durable-session-replay-review.
- Definir em ADR a política de evidência qualitativa/externa antes da implementação.

### Non-goals
- Alterar navegação curricular, gerar solução ou executar comandos livres.
- Reabrir atestações históricas de specs done.

## 2. Requirements

### Functional
- R1: Rejeitar evidência desconhecida ou pertencente a outra sessão antes de anexar avaliação/tentativa.
- R2: Verificar vínculo ao desafio fixado, check e nó aplicável; não transferir resultado estrutural entre escopos sem regra explícita.
- R3: Revalidar fingerprint, root e globs no consumo; evidência estrutural obsoleta, ausente ou com root indisponível não pode produzir verdict met.
- R4: Definir registro autorizado para evidência qualitativa/externa, sem aceitar IDs arbitrários como prova; preservar rubrica e origem.
- R5: Preservar replay histórico e retry de avaliação já confirmada, sem reavaliar retrospectivamente fatos contra o workspace atual.
- R6: Aplicar os mesmos gates antes e após reinício; expor erro MCP estável e seguro sem paths ou conteúdo privado.
- R7: Testar o fluxo real observe/check/evaluate com duas sessões, alteração de arquivo e restart; rejeições não avançam revisão nem criam tentativa.

### Non-functional
- Evitar janela de aprovação indevida entre verificar fingerprint e consumir resultado.
- Manter checks determinísticos e explicitar limitações de concorrência no ADR.

### Security
Não confiar em evidence_id fornecido pelo tutor como autorização. Distinguir recuperação de histórico de aprovação de evidência nova.

### Compatibility
Planejar migração dos testes que usam IDs sintéticos sem registro. Não converter ausência de prova em aprovação; documentar o tratamento de avaliações legadas sem reescrevê-las.

## 3. Technical Plan

### Affected areas
- internal/session, internal/application, internal/mcpserver e composição em cmd/codinho.

### Artifacts
- modified: internal/session/service.go
- modified: internal/application/session.go
- modified: internal/application/assessment.go
- modified: internal/application/workspace.go
- modified: internal/mcpserver/errors.go
- modified: cmd/codinho/main.go
- created: internal/application/evaluation_evidence_test.go
- created: cmd/codinho/evaluation_evidence_integration_test.go
- modified: .pose/indexes/validation-matrix.json
- modified: docs/compatibility.md

### Delivery targets
- capability:evaluation-evidence-lineage module:cmd/codinho profile:composed-capability entrypoint:cmd/codinho/main.go

### API/contract changes
Definir interface de validação de evidência na composição; evitar dependência circular entre serviços. Declarar qualquer novo artefato de ADR/teste antes de alterá-lo.

### Data/storage changes
Reutilizar eventos de observação/check com baseline/root/globs persistidos; não criar banco paralelo.

### Technical risks
EvidenceGet protege leitura, mas não é chamado por StepEvaluate. CheckOutcomeOverride consulta o store global; presença de um ID e verdict pass não provam escopo nem atualidade.

## 4. Tasks

### Planning
- [ ] Reproduzir consumo indevido por dois session IDs e por fingerprint obsoleto.
- [ ] Definir política estrutural/qualitativa e ADR com compatibilidade.

### Implementation
- [ ] Integrar gate no caminho de consumo, antes do append.
- [ ] Preservar retries históricos e reavaliar evidências novas.
- [ ] Expor diagnósticos seguros e atualizar fixtures de testes.

### Validation
- [ ] Provar R1–R7 por unidades e processos MCP reais.
- [ ] Registrar check permanente na matriz e obter revisão independente.

## 5. Decisions

### Decision 1
- Date: 2026-09-07
- Context: Revisão independente da recuperação encontrou bypass pré-existente em StepEvaluate/checkOutcomeOverride.
- Options considered: ampliar recuperação; remediação própria antes do aceite.
- Decision: isolar o consumo de evidência nesta spec e torná-la dependência do aceite V1.
- Rationale: persistir escopo de leitura não define autorização para avaliação qualitativa ou estrutural.
- Consequences: não reivindicar isolamento completo enquanto esta spec não estiver done.

## 6. Validation

### Strategy
Executar cenários positivo, cross-session, stale, root substituído, evidência ausente e retry após restart.

### Deterministic checks
- Test: go test -race ./...
- Integration: go test ./cmd/codinho -run TestEvaluationEvidenceOverRealStdio
- Lint: pose lint-spec evaluation-evidence-lineage --ready-check; pose check --strict
- Typecheck: go vet ./...
- Build: go build ./cmd/codinho
- Security / Contract: pose assess integrate; pose validate --strict

### Execution log
- 2026-09-07 UTC: criada a partir de inspeção do código e revisão independente; implementação não iniciada.

### Results summary
Planejamento pronto para reprodução; nenhuma correção de consumo atribuída a esta spec.

### Requirement trace
Preencher R1–R7 com evidência do candidato durante a implementação.

### Known gaps
O produto atual não bloqueia todo consumo cross-session/obsoleto. A spec de recuperação resolve somente leitura/freshness.

## 7. Final Report

### Delivered scope
Somente planejamento e dependência no roadmap.

### Files and modules changed
- Esta spec; roadmap codinho-v1; dependência de v1-integrated-acceptance.

### Validation executed
Validar estrutura e readiness junto à atualização do planejamento.

### Residual risks
Não usar check pass como prova suficiente de aprovação até entregar este gate.

### Follow-ups
- [covered: v1-integrated-acceptance] Reexecutar isolamento e atualidade no candidato composto V1.
