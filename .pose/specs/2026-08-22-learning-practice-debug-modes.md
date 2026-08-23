---
slug: learning-practice-debug-modes
status: draft
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: session-orchestration-disclosure, assistance-hints-detours, feedback-evaluation-progression, mastery-review-scheduling, tutor-skill-host-integration
priority: 150
components: session-modes, tutor-skill
delivers:
---

# Spec: learning-practice-debug-modes

## 1. Intent

### Goal
Implementar os modos ensino, prática, revisão, depuração e exploração com políticas independentes e adaptação transparente.

### Business value
Permitir que a mesma árvore curricular atenda perfis e intenções diferentes sem duplicar desafios.

### Constraints
- Modo não determina sozinho profundidade, ajuda ou avaliação.
- Ajuste automático nunca supera escolha manual.
- Depuração preserva a oportunidade de formular hipóteses.

### Non-goals
- Modo entrevista, tratado em spec separada.
- Personalização baseada em modelo probabilístico.

## 2. Requirements

### Functional
- R1: Definir defaults independentes para os cinco modos sem acoplar suas dimensões.
- R2: Permitir override explícito de profundidade, ajuda, avaliação, avanço e tempo.
- R3: No ensino, contextualizar conceitos e permitir pistas progressivas.
- R4: Na prática, reduzir contexto e registrar dependência de pistas.
- R5: Na revisão, selecionar competências vencidas e variações curtas.
- R6: Na depuração, ordenar reprodução, divergência, hipótese, observação e correção.
- R7: Na exploração, permitir feedback livre sem obrigação de tentativa ou avanço.
- R8: Adaptar granularidade por evidência repetida e explicar toda mudança.
- R9: Permitir que o aluno proponha o próximo passo e registrar a autonomia.
- R10: Autorar ao menos um desafio `kind: debug` com `fixture` (campo
  entregue por administrative-cli-fixtures) e validar `workspace
  prepare` de ponta a ponta com o conteúdo real desse desafio.

### Non-functional
- Políticas devem ser versionadas e testáveis por matriz.
- Mudanças de modo não podem perder evidências.

### Security
- Nenhum modo reduz as restrições de workspace e execução.

### Compatibility
- Novos modos devem compor as mesmas dimensões.

## 3. Technical Plan

### Affected areas
- internal/session/, internal/mastery/, .agents/skills/codinho/

### Artifacts
- created: internal/session/modes.go
- created: internal/session/adaptation.go
- created: internal/session/modes_test.go
- created: internal/session/adaptation_test.go
- modified: internal/mastery/scheduler.go
- modified: .agents/skills/codinho/references/session-modes.md
- created: testdata/host/mode-transcripts/
- created: packs/go-debugging.yaml (desafio(s) `kind: debug` com `fixture`, requirement R10)

### Delivery targets
Nenhum novo; amplia a capability de tutoria planejada.

### API/contract changes
- Completar SessionPolicy e session_configure com mode presets explicáveis.

### Data/storage changes
- Persistir mode_changed, granularity_changed e learner_next_step_proposed.

### Technical risks
- Adaptação agressiva pode frustrar ou criar dependência.
- Defaults podem virar rótulos rígidos de nível.

## 4. Tasks

### Planning
- [ ] Definir matriz de defaults e critérios de adaptação.
- [ ] Definir protocolo de depuração e evidência de autonomia.

### Implementation
- [ ] Implementar presets e overrides.
- [ ] Implementar adaptação conservadora e opt-out.
- [ ] Integrar revisão vencida e protocolo de depuração.
- [ ] Atualizar skill e transcripts.
- [ ] Cobrir troca de modo no meio da sessão.

### Validation
- [ ] Executar matriz modo versus dimensão.
- [ ] Testar que segurança nunca relaxa.
- [ ] Validar transcripts com iniciantes e experientes.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Modo e nível do aluno não são a mesma coisa.
- Options considered: presets rígidos; configuração livre sem defaults; presets sobre dimensões.
- Decision: modos fornecem defaults substituíveis.
- Rationale: preserva flexibilidade e clareza.
- Consequences: respostas devem mostrar configuração efetiva, não apenas nome do modo.

## 6. Validation

### Strategy
Testar matriz de políticas, eventos e transcripts comportamentais.

### Deterministic checks
- Test: go test ./internal/session/... ./internal/mastery/...
- Lint: gofmt -l internal/session internal/mastery
- Typecheck: go vet ./internal/session/... ./internal/mastery/...
- Build: go build ./cmd/codinho
- Security / Contract: golden policies e adversarial transcripts.

### Execution log
- Pendente.

### Results summary
- Nenhum modo especializado entregue.

### Requirement trace
- Mapear R1–R9 a matriz, tests e transcripts.

### Known gaps
- Calibração humana de adaptação ficará pendente até piloto.

## 7. Final Report

### Delivered scope
Nenhum; spec draft.

### Files and modules changed
- Planejados em session, mastery e skill.

### Validation executed
- Command: pose lint-spec learning-practice-debug-modes --ready-check
- Result: registrar após gate.

### Residual risks
- Ajustes automáticos devem começar desabilitados ou conservadores.

### Follow-ups
- [covered: v1-integrated-acceptance] Medir progressão de granularidade no piloto.
