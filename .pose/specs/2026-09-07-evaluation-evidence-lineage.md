---
slug: evaluation-evidence-lineage
status: in-progress
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
- modified: internal/session/recovery.go
- modified: internal/session/checks_override_test.go
- modified: internal/application/session.go
- modified: internal/application/assessment.go
- modified: internal/application/workspace.go
- modified: internal/application/checks.go
- modified: internal/application/checks_test.go
- modified: internal/mcpserver/assessment_tools.go
- modified: internal/mcpserver/contract_test.go
- modified: internal/mcpserver/errors.go
- modified: internal/mcpserver/instructions.go
- modified: internal/assessment/criteria.go
- modified: internal/eventstore/event.go
- modified: cmd/codinho/integration_test.go
- modified: .pose/indexes/validation-matrix.json
- modified: docs/compatibility.md
- created: internal/session/evaluation_evidence.go
- created: internal/application/evaluation_evidence.go
- created: internal/application/evaluation_evidence_test.go
- created: cmd/codinho/evaluation_evidence_integration_test.go
- created: .pose/adr/2026-09-07-scoped-evaluation-evidence-and-qualitative-registration.md
- created: .pose/knowledge/2026-09-07-decision-log-adr-scoped-evaluation-evidence-review.md
- created: .pose/changelogs/unreleased/evaluation-evidence-lineage.md

### Delivery targets
- capability:evaluation-evidence-lineage module:cmd/codinho profile:composed-capability entrypoint:cmd/codinho/main.go

### API/contract changes
Injetar EvaluationEvidenceValidator no serviço de sessão pela composição em application.NewSessionService. Implementar origem/freshness no adaptador de aplicação; o núcleo não consulta filesystem. Adicionar evidence_record para registro qualitativo e check_id nos critérios estruturais baseados em checks. Campos novos de entrada usam omitempty para preservar identidade de retries legados.

### Data/storage changes
Reutilizar eventos de observação/check com baseline/root/globs persistidos. Adicionar evidence_recorded e anexar evidence_lineage às avaliações novas; preservar replay antigo. Registro qualitativo guarda blob redigido, origem, rubrica e vínculo ao desafio/passo; não cria banco paralelo.

### Technical risks
EvidenceGet protege leitura, mas não é chamado por StepEvaluate. CheckOutcomeOverride consulta o store global; presença de um ID e verdict pass não provam escopo nem atualidade.

## 4. Tasks

### Planning
- [x] Reproduzir consumo indevido por dois session IDs e por fingerprint obsoleto.
- [x] Definir política estrutural/qualitativa e ADR com compatibilidade.

### Implementation
- [x] Integrar gate no caminho de consumo, antes do append.
- [x] Preservar retries históricos e reavaliar evidências novas.
- [x] Expor diagnósticos seguros e atualizar fixtures de testes.

### Validation
- [x] Provar R1–R7 por unidades e processos MCP reais.
- [ ] Registrar check permanente na matriz e obter revisão independente.

## 5. Decisions

### Decision 2
- Date: 2026-09-07
- Decision: aplicar o [ADR de evidência com escopo](../adr/2026-09-07-scoped-evaluation-evidence-and-qualitative-registration.md). knowledge:adr-scoped-evaluation-evidence-review.
- Rationale: autorização pertence ao log da sessão e a presença de um blob não prova execução nem atualidade.
- Consequences: IDs externos exigem registro; checks exigem check_id. Sem evidência estrutural verificável, resultado é unverifiable. Política de cobertura de todos os critérios autorados permanece fora desta remediação.

### Decision 1
- Date: 2026-09-07
- Context: Revisão independente da recuperação encontrou bypass pré-existente em StepEvaluate/checkOutcomeOverride.
- Options considered: ampliar recuperação; remediação própria antes do aceite.
- Decision: isolar o consumo de evidência nesta spec e torná-la dependência do aceite V1.
- Rationale: persistir escopo de leitura não define autorização para avaliação qualitativa ou estrutural.
- Consequences: não reivindicar isolamento completo enquanto esta spec não estiver done.

## 6. Validation

### Strategy
Risco alto: autorização de consumo e persistência cruzam sessão, aplicação e MCP.

| Cenário | Comando obrigatório | Resultado esperado |
|---|---|---|
| Origem, nó/check, blob ausente/corrompido, root/globs e drift | go test ./internal/application -run TestEvaluationEvidence | Rejeitar antes do append; revisão e tentativa intactas |
| Registro qualitativo, rubrica/origem, isolamento e retry | go test ./internal/application -run TestEvaluationEvidence | Registro explícito; estrutural nunca usa nota como check |
| Ausência do adaptador e retries históricos | go test ./internal/session | Falhar fechado para citações; preservar avaliações confirmadas |
| MCP observe/check/evaluate e restart | go test ./cmd/codinho -run TestEvaluationEvidenceOverRealStdio | Duas sessões, drift e revalidação após reinício; erro seguro |
| Concorrência, domínio, contratos e instalação | go test -race ./...; go vet ./...; go build ./cmd/codinho | Sem regressão ou data race |
| Gates de entrega | pose validate --strict --json .pose/results/delivery-validation.json; pose surface-check --spec evaluation-evidence-lineage --strict | Evidência atual e atribuída |

Timeout/erro de executor continua unverifiable; rede permanece controlada pelo catálogo. Erro de leitura/fingerprint recusa consumo. Validação faz amostragem antes da resolução e imediatamente antes do append; escritores externos não participam do mutex e alterações ABA não são detectáveis sem snapshot isolado. Não prometer atomicidade do filesystem. Check novo exige fingerprints iguais antes/depois da execução.

### Deterministic checks
- Test: go test -race ./...
- Integration: go test ./cmd/codinho -run TestEvaluationEvidenceOverRealStdio
- Lint: pose lint-spec evaluation-evidence-lineage --ready-check; pose check --strict
- Typecheck: go vet ./...
- Build: go build ./cmd/codinho
- Security / Contract: pose assess integrate; pose validate --strict

### Execution log
- 2026-09-07 UTC: criada a partir de inspeção do código e revisão independente; implementação não iniciada.
- 2026-09-07 UTC: assess discover encontrou um módulo Go (risco alto). Regressão TestEvaluationEvidenceRejectsCrossSessionAndDrift falhou nos dois casos por aceitação indevida antes da implementação.
- 2026-09-07 UTC: regressões de aplicação/sessão e MCP real passaram após integrar autorização, registro qualitativo e check_id. Retry legado testado com o formato JSON original sem o campo novo.
- 2026-09-07 UTC: revisor independente agent:evidence_review, modelo gpt-5.6-luna, aprovou inicialmente sem ressalvas após testes completos, E2E e vet. Aplicado o padrão de .claude/skills/agent-batch-review/SKILL.md por delegação nativa da sessão. Gates humanos não são substituídos.
- 2026-09-07 UTC: assess tech-debt encontrou zero marcadores; recurrence-check zero recorrências. assess integrate reconheceu zero contratos (limitação conhecida de MCP Go); prova de integração vem do E2E real.

### Results summary
Implementação e testes focados passaram; gates estruturados, atribuição Git e revisão final condicionam o fechamento.

### Requirement trace
- R1 [satisfied] test:TestEvaluationEvidenceRejectsCrossSessionAndDrift test:TestEvaluationEvidenceRejectsInvalidScopeAndBlobs
- R2 [satisfied] test:TestEvaluationEvidenceRejectsInvalidScopeAndBlobs test:TestEvaluationEvidenceCheckOutcomesAndObservationLimits
- R3 [satisfied] test:TestEvaluationEvidenceRejectsCrossSessionAndDrift test:TestEvaluationEvidenceRechecksImmediatelyBeforeAppend test:TestEvaluationEvidenceRejectsDriftDuringCheck
- R4 [satisfied] test:TestEvaluationEvidenceQualitativeRegistrationAndRecovery test:TestEvaluationEvidenceRegistrationValidationAndRedaction
- R5 [satisfied] test:TestEvaluationEvidenceLegacyRetryPreservesInputIdentity test:TestEvaluationEvidenceQualitativeRegistrationAndRecovery
- R6 [satisfied] capability:evaluation-evidence-lineage evidence:integration check:evaluation-evidence test:TestEvaluationEvidenceOverRealStdio
- R7 [satisfied] capability:evaluation-evidence-lineage evidence:integration check:evaluation-evidence test:TestEvaluationEvidenceOverRealStdio

### Known gaps
Fingerprint amostra arquivos no escopo dos globs; não congela escritores externos nem detecta ABA. Evidências externas registradas são declarações, não provas autenticadas. A política de cobertura completa de critérios autorados continua distinta.

## 7. Final Report

### Delivered scope
Autorização de evidências no consumo, registro qualitativo, replay/retry compatíveis, diagnóstico MCP seguro e regressões compostas. Implementado e testado; lifecycle aguarda gates finais.

### Files and modules changed
- Núcleo de sessão, adaptador de aplicação, workspace/checks, contrato MCP, testes e documentação declarados em Artifacts.

### Validation executed
Testes focados de aplicação/sessão, contrato MCP, E2E real e go vet passaram na revisão independente inicial; suíte race e validação estruturada serão anexadas ao candidato.

### Residual risks
Limitações de amostragem e seleção de critérios estão explícitas no ADR e em docs/compatibility.md; nenhuma aprovação de catálogo ou aceite humano V1 é inferida.

### Follow-ups
- [covered: v1-integrated-acceptance] Reexecutar isolamento e atualidade no candidato composto V1.
