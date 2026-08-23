---
slug: learning-practice-debug-modes
status: done
created_at: 2026-08-22
completed_at: 2026-08-23
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
- created: internal/session/modes_test.go
- created: internal/session/adaptation.go
- created: internal/session/adaptation_test.go
- modified: internal/session/service.go
- modified: internal/session/service_test.go
- modified: internal/application/session.go
- modified: internal/eventstore/event.go
- modified: internal/mcpserver/session_tools.go
- modified: internal/mcpserver/errors.go
- modified: internal/mcpserver/contract_test.go
- modified: internal/mcpserver/session_contract_test.go
- modified: .agents/skills/codinho/references/session-modes.md
- modified: .agents/skills/codinho/references/tutor-contract.md
- modified: .agents/skills/codinho/references/mcp-tool-routing.md
- modified: .agents/skills/codinho/SKILL.md
- created: testdata/host/mode-transcripts/debug-mode-session.md
- created: packs/go-debugging.yaml
- modified: packs/manifest.yaml
- modified: cmd/codinho/cli_integration_test.go

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
- [x] Definir matriz de defaults e critérios de adaptação.
- [x] Definir protocolo de depuração e evidência de autonomia.

### Implementation
- [x] Implementar presets e overrides.
- [x] Implementar adaptação conservadora e opt-out.
- [x] Integrar revisão vencida e protocolo de depuração.
- [x] Atualizar skill e transcripts.
- [x] Cobrir troca de modo no meio da sessão.

### Validation
- [x] Executar matriz modo versus dimensão.
- [x] Testar que segurança nunca relaxa.
- [x] Validar transcripts com iniciantes e experientes.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Modo e nível do aluno não são a mesma coisa.
- Options considered: presets rígidos; configuração livre sem defaults; presets sobre dimensões.
- Decision: modos fornecem defaults substituíveis.
- Rationale: preserva flexibilidade e clareza.
- Consequences: respostas devem mostrar configuração efetiva, não apenas nome do modo.

### Decision 2
- Date: 2026-08-23
- Context: a task "Cobrir troca de modo no meio da sessão" sugeria
  permitir mudar `mode` via `session_configure`. PROJECT.md §15.6 lista
  explicitamente o que `session_configure` pode alterar: "granularidade,
  política de auxílio, avaliação, tempo e preferências" — `mode` não
  está nessa lista, e granularidade já tem sua própria tool
  (`granularity_adjust`).
- Options considered: (a) adicionar `Mode` a `ConfigureInput`,
  permitindo trocar o modo pedagógico no meio da sessão; (b) manter
  `mode` imutável após `session_start`, como o contrato já declara.
- Decision: (b). `mode` permanece fixado em `session_start`; para trocar
  de modo o tutor inicia uma nova sessão.
- Rationale: introduzir mutação de `mode` violaria um contrato já
  publicado em PROJECT.md §15.6, sem que nenhum requirement desta spec
  exija essa mutação — R1/R2 pedem defaults independentes por modo e
  override explícito no início, não troca dinâmica.
- Consequences: a "troca de modo no meio da sessão" é resolvida como
  não-suportada por design, documentada aqui em vez de implementada;
  nenhuma mudança de código foi necessária além desta constatação.

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
- `go build ./...` → ok (2026-08-23).
- `gofmt -l internal/session internal/mastery internal/mcpserver internal/application internal/eventstore` → sem saída (2026-08-23).
- `go vet ./...` e `GOOS=windows go vet ./...` → sem diagnósticos (2026-08-23).
- `go test ./... -race` → `ok` em todos os pacotes, incluindo `internal/session` (modes.go/adaptation.go novos) e `internal/mcpserver` (learner_next_step_propose, granularity_adjust com reason), sem data races (2026-08-23).
- `go run ./cmd/codinho catalog validate` → `catalog: ok, no diagnostics` com `packs/go-debugging.yaml` carregado (2026-08-23).
- `go test ./cmd/codinho/... -run TestWorkspacePrepareMaterializesRealAuthoredDebugChallenge -v` → PASS; `workspace prepare` real materializou o `main.go` com bug autêntico do desafio `go-debug.slice-off-by-one` (2026-08-23).
- `govulncheck ./...` → "No vulnerabilities found." (2026-08-23).

### Results summary
Defaults independentes por modo (R1/R2) aplicados em `session.Start`;
protocolo de depuração de 6 estágios (R6) exposto por
`DebugProtocolStages` e documentado na skill; granularidade adaptativa
com explicação obrigatória (R8) via `SuggestGranularity` (função pura,
nunca muta nada) e `GranularityAdjust` agora persistindo `reason`;
autonomia do aluno (R9) via nova tool `learner_next_step_propose`;
primeiro desafio `kind: debug` real com `fixture` autoral
(`packs/go-debugging.yaml`, requirement R10) validado ponta a ponta via
`workspace prepare` real. Troca de modo no meio da sessão foi investigada
e resolvida como não-suportada por design (Decision 2), consistente com
PROJECT.md §15.6.

### Requirement trace
- R1 [satisfied] internal/session/modes.go (DefaultsForMode) + test:TestDefaultsForModeAreIndependentPerMode, test:TestModeDefaultsAreAllDistinctAcrossTheFiveNonInterviewModes.
- R2 [satisfied] internal/session/service.go (Start usa defaults só quando campo vazio) + test:TestStartExplicitFieldsOverrideModeDefaults.
- R3 [satisfied] defaults de teaching (macro/progressive/concept_or_api) + .agents/skills/codinho/references/session-modes.md.
- R4 [satisfied] defaults de practice (micro/progressive) já existentes, preservados; help/pistas já registrados via hint_requested (assistance-hints-detours).
- R5 [satisfied] defaults de review (meso/limited/on_step_complete) + reaproveita review_due (mastery-review-scheduling, já done).
- R6 [satisfied] internal/session/modes.go (DebugProtocolStages) + packs/go-debugging.yaml (desafio real com os 6 estágios como macro/meso/micro steps) + test:TestDebugProtocolStagesIsOrderedAndComplete.
- R7 [satisfied] defaults de exploration (free/on_demand, nada força tentativa ou avanço) + test:TestDefaultsForModeAreIndependentPerMode.
- R8 [satisfied] internal/session/adaptation.go (SuggestGranularity, EvidenceThreshold configurável) + GranularityAdjust persiste reason + test:TestSuggestGranularityWidensOnRepeatedEase, test:TestGranularityAdjustPersistsReason, test:TestContractGranularityAdjustAcceptsReason.
- R9 [satisfied] eventstore.EventLearnerNextStepProposed + session.ProposeNextStep + tool learner_next_step_propose + test:TestProposeNextStepRecordsSignalWithoutAdvancing, test:TestContractLearnerNextStepProposeRecordsWithoutAdvancing.
- R10 [satisfied] packs/go-debugging.yaml (fixture real) + test:TestWorkspacePrepareMaterializesRealAuthoredDebugChallenge.

### Known gaps
- Calibração humana de `EvidenceThreshold` (hoje fixo em 3) ficará pendente até piloto real — variável exportada, fácil de ajustar sem mudar a API.
- `SuggestGranularity` é uma função pura consultada pelo raciocínio da skill a partir de `progress_get`; nenhuma tool MCP nova a expõe diretamente (evitou-se crescer a superfície de tools além do necessário para esta spec).
- Mecânica completa de `interview` (briefing, cronômetro, rubrica) permanece para `interview-mode`, como já era non-goal desta spec.

## 7. Final Report

### Delivered scope
Defaults independentes por modo, protocolo de depuração estruturado,
granularidade adaptativa explicável, sinal de autonomia do aluno e
primeiro desafio `kind: debug` real com fixture.

### Files and modules changed
- internal/session/{modes,adaptation}.go (novos) + service.go (Start usa defaults por modo; GranularityAdjust ganha reason; novo ProposeNextStep).
- internal/eventstore/event.go (EventLearnerNextStepProposed).
- internal/application/session.go, internal/mcpserver/{session_tools,errors}.go (nova tool learner_next_step_propose; reason em granularity_adjust).
- packs/go-debugging.yaml, packs/manifest.yaml (primeiro desafio kind:debug com fixture real).
- .agents/skills/codinho/references/{session-modes,tutor-contract,mcp-tool-routing}.md, SKILL.md (documentação atualizada, 30 tools).
- testdata/host/mode-transcripts/debug-mode-session.md.

### Validation executed
- Command: go test ./... -race
- Result: ok em todos os pacotes.
- Command: go test ./cmd/codinho/... -run TestWorkspacePrepareMaterializesRealAuthoredDebugChallenge
- Result: PASS (fixture real materializada ponta a ponta).
- Command: govulncheck ./...
- Result: No vulnerabilities found.

### Residual risks
- Ajustes automáticos de granularidade começam com threshold conservador (3 evidências) e nunca sobrepõem escolha manual (constraint preservada).
- `packs/go-debugging.yaml` é o primeiro e único desafio `kind: debug` — a curadoria completa de conteúdo pertence a `catalog-authoring-quality`/`go-*-packs`.

### Follow-ups
- [covered: v1-integrated-acceptance] Medir progressão de granularidade no piloto.
- [open] Calibrar `EvidenceThreshold` com dados reais de sessões (hoje fixo em 3, sem piloto).
