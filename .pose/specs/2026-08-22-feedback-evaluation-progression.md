---
slug: feedback-evaluation-progression
status: done
created_at: 2026-08-22
completed_at: 2026-08-23
supersedes:
depends_on: session-orchestration-disclosure, local-event-store, mcp-stdio-foundation
priority: 80
components: assessment, sessions, mcp-server
delivers:
---

# Spec: feedback-evaluation-progression

## 1. Intent

### Goal
Implementar feedback consultivo, avaliação híbrida, reflexão, conclusão e avanço como operações independentes e auditáveis.

### Business value
Permitir experimentação segura e avaliação justa sem transformar opinião idiomática em reprovação.

### Constraints
- MCP prepara contexto; o agente redige feedback e julgamento semântico.
- Tentativa só existe com intenção de submissão.
- Avaliação nunca conclui; conclusão nunca avança.

### Non-goals
- Executar checks ou observar workspace.
- Calcular retenção longitudinal.

## 2. Requirements

### Functional
- R1: Preparar pacote de feedback com instrução, pergunta, rubrica, observações e evidências.
- R2: Registrar feedback com tipos e progress_effect sem criar avaliação ou tentativa.
- R3: Avaliar critérios determinísticos e julgamentos qualitativos separadamente.
- R4: Exigir referência de evidência e rubrica para todo julgamento qualitativo.
- R5: Produzir met, partially_met, not_met, unverifiable ou not_applicable por critério.
- R6: Classificar achados em blocking, important_non_blocking ou advisory.
- R7: Concluir apenas quando a política permitir e avançar somente em chamada separada.
- R8: Registrar override sem falsificar avaliação positiva ou domínio.
- R9: Registrar reflexão como evidência distinta da implementação.
- R10: Expor feedback_prepare, feedback_record, step_evaluate, reflection_record, step_complete e step_advance.

### Non-functional
- Toda mutação deve ser idempotente e revisionada.
- O histórico deve preservar avaliações múltiplas do mesmo passo.

### Security
- Evidências e feedback devem aplicar redaction e escopo da sessão.
- Texto qualitativo não pode conter instruções executáveis.

### Compatibility
- Novos tipos de feedback e critérios devem ser aditivos.

## 3. Technical Plan

### Affected areas
- internal/assessment/, internal/session/, internal/mcpserver/

### Artifacts
- created: internal/assessment/rubric.go
- created: internal/assessment/criteria.go
- created: internal/assessment/rubric_test.go
- created: internal/assessment/criteria_test.go
- modified: internal/session/service.go
- created: internal/session/progression.go
- created: internal/session/progression_test.go
- created: internal/session/assessment_test.go
- created: internal/mcpserver/assessment_tools.go
- created: internal/mcpserver/assessment_contract_test.go
- modified: internal/mcpserver/server.go
- modified: internal/mcpserver/envelope.go
- modified: internal/mcpserver/contract_test.go
- created: internal/application/assessment.go
- modified: internal/learning/evaluation.go
- modified: internal/learning/evidence.go
- modified: internal/learning/errors.go
- modified: internal/learning/policy_test.go

### Delivery targets
Nenhum novo; amplia o contrato MCP V1 planejado.

### API/contract changes
- Adicionar seis tools e envelopes de avaliação e achado.

### Data/storage changes
- Persistir feedback_recorded, attempt_submitted, evaluation_recorded, reflection_recorded, step_completed e step_advanced.

### Technical risks
- O agente pode produzir avaliação inconsistente com a rubrica.
- Completion override pode ser abusado e distorcer métricas.

## 4. Tasks

### Planning
- [x] Definir rubricas de Go idiomático e comunicação técnica.
- [x] Definir matriz critério versus evidência e severidade.

### Implementation
- [x] Implementar packet de feedback e registro.
- [x] Implementar avaliação híbrida e validação de evidências.
- [x] Implementar reflexão, conclusão, override e avanço.
- [x] Expor tools e allowed_actions.
- [x] Cobrir sequências inválidas e avaliações repetidas.

### Validation
- [x] Executar tabelas de transição e contract tests.
- [x] Testar feedback com progress_effect sem avanço.
- [x] Testar spoof de evidência e revisão obsoleta.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: O servidor não possui LLM, mas precisa persistir avaliação semântica.
- Options considered: fingir avaliação no servidor; não persistir; aceitar julgamento estruturado do tutor.
- Decision: receber julgamentos tipados, cada um ligado a rubrica e evidência.
- Rationale: mantém responsabilidades honestas e auditáveis.
- Consequences: a skill precisa orientar o agente e testes devem validar o envelope, não a linguagem aberta.

### Decision 2
- Date: 2026-08-22
- Context: `step_evaluate` deveria "executar critérios determinísticos elegíveis"
  (PROJECT.md §21.1), mas o non-goal explícito desta spec exclui "executar
  checks ou observar workspace" (isso pertence a safe-check-executor). O
  servidor também não interpreta o conteúdo de uma evidência — `Evidence.ref`
  é opaco por design (learning-domain-model).
- Options considered: (a) confiar cegamente no verdict que o chamador
  envia mesmo para critérios `kind: structural`; (b) o servidor deriva o
  verdict de um critério estrutural de forma puramente determinística a
  partir da PRESENÇA de uma evidência citada (não do seu conteúdo) —
  `met` se citada, `unverifiable` caso contrário — ignorando qualquer
  verdict que o chamador tente enviar para esse kind.
- Decision: (b), implementado em `internal/assessment.Resolve`.
- Rationale: honra o non-goal (nenhum check é executado, nenhuma evidência é
  interpretada) sem fingir raciocínio que o servidor não tem, e ainda assim
  produz um verdict genuinamente determinístico (não confia em auto-relato
  do chamador para critérios que deveriam ser objetivos).
- Consequences: a execução REAL de checks (compilação, testes, etc. de
  §21.1) permanece para safe-check-executor; até lá, todo critério
  `structural` sem evidência citada resulta em `unverifiable`, nunca `met`
  por presunção.

### Decision 3
- Date: 2026-08-22
- Context: R6 exige achados classificados em `blocking`/
  `important_non_blocking`/`advisory`, mas `learning.CriterionResult` (de
  learning-domain-model) só tinha um campo `Blocking bool` binário, sem
  nenhum consumidor ainda.
- Decision: evoluir `CriterionResult.Blocking bool` para
  `Severity FindingSeverity` (três valores), tratando "achado" como o
  próprio resultado por critério, em vez de introduzir um tipo `Finding`
  paralelo. `HasBlockingFailure` passou a considerar `not_applicable` como
  nunca bloqueante.
- Rationale: `CriterionResult` não tinha nenhum consumidor real ainda
  (zero usos fora de um teste), então evoluir o campo é seguro e evita
  duplicar o conceito de "resultado por critério" com um tipo `Finding`
  redundante.
- Consequences: `internal/learning/policy_test.go` precisou de ajuste
  (campo renomeado); `NewCriterionResult` agora também aplica R4
  (evidência+rubrica obrigatórias para julgamento não-estrutural).

### Decision 4
- Date: 2026-08-22
- Context: `step_evaluate` pode produzir DOIS fatos distintos numa mesma
  chamada — a avaliação em si, e (se `submission_intent` for verdadeiro) uma
  tentativa — mas `eventstore.Store.Append` só grava um evento por chamada,
  e seu cache de idempotência (`s.seen`) é indexado GLOBALMENTE por
  `request_id`, não por `(stream, request_id)` (achado ao investigar
  `internal/eventstore/store.go:94-102` durante esta spec).
- Decision: dois `Append`s sequenciais por chamada quando
  `submission_intent` é verdadeiro, usando `request_id` para o primeiro
  (evaluation_recorded) e `request_id + ":attempt"` para o segundo
  (attempt_submitted) — nunca o mesmo `request_id` duas vezes na mesma
  chamada, o que colidiria com o cache global do primeiro Append e
  devolveria o evento errado.
- Rationale: preserva idempotência independente de cada fato mesmo com o
  cache sendo global-por-request_id; documentado inline em
  `internal/session/service.go`.
- Consequences: **achado colateral não corrigido nesta spec** — o cache
  `s.seen` de `eventstore.Store.Append` não verifica se o evento
  cacheado pertence ao MESMO `streamID` da chamada atual. Dois clientes
  diferentes reutilizando o mesmo `request_id` para sessões (streams)
  diferentes podem receber o evento um do outro. Isso é uma característica
  pré-existente de `local-event-store` (spec já fechada e selada), não
  introduzida aqui; ver Known Gaps e Follow-ups para a decisão de reportar
  ou corrigir.

### Decision 5
- Date: 2026-08-22
- Context: `step_advance` precisa decidir "o próximo passo" numa árvore
  autoral que permite múltiplos filhos (ramificação) — PROJECT.md §15.6 diz
  apenas "ativa o próximo nó permitido ou informa opções quando houver
  ramificação", sem detalhar o algoritmo de travessia.
- Options considered: (a) combinar todos os `layers` do desafio numa única
  trilha; (b) andar em pré-ordem apenas dentro da primeira `layer` com
  `macro_steps` não vazios — a mesma convenção que `firstStep`/
  `deriveWindow` (granularity.go, session-orchestration-disclosure) já
  usam.
- Decision: (b), implementado em `internal/session/progression.go`
  (`advanceFrom`/`rootSteps`). Um nó com exatamente um filho desce
  automaticamente; mais de um filho retorna as opções sem mutar nada; zero
  filhos sobe para o próximo irmão do ancestral mais próximo; fim da árvore
  retorna `Done: true`.
- Rationale: `layers` representam perspectivas paralelas do mesmo desafio,
  não uma sequência — combiná-las não tem semântica clara ainda, e reusar a
  convenção já estabelecida evita duas noções incompatíveis de "próximo
  passo" convivendo no mesmo código.
- Consequences: `step_advance` nunca move a "janela" para uma `layer`
  diferente da que a sessão já está usando; sessões que precisem trocar de
  layer devem usar `granularity_adjust` ou um `session_start` novo.

## 6. Validation

### Strategy
Testar operações isoladas e sequências completas, incluindo feedback sem efeito e overrides.

### Deterministic checks
- Test: go test ./internal/assessment/... ./internal/session/...
- Lint: gofmt -l internal/assessment internal/session
- Typecheck: go vet ./internal/assessment/... ./internal/session/...
- Build: go test ./internal/assessment/... -run ^$
- Security / Contract: schema tests, redaction e rejeição de evidence IDs fora da sessão.

### Execution log
- `gofmt -l internal/assessment internal/session internal/learning internal/curriculum internal/mcpserver internal/application cmd/codinho` → saída vazia (2026-08-23).
- `go vet ./...` → sem findings (2026-08-23).
- `go build ./...` → ok (2026-08-23).
- `go test ./... -race` → `ok` em todos os pacotes com testes, 100 funções
  de teste passando em assessment/session/learning/mcpserver combinados,
  sem data races (2026-08-23).
- `govulncheck ./...` → "No vulnerabilities found." (2026-08-23).
- Smoke test real de ponta a ponta via `mcp.CommandTransport` contra o
  binário `codinho serve` e `packs/go-first-steps.yaml` (desafio
  `go-data.slice-filter-preserve-input`): `session_start` →
  `feedback_prepare` → `feedback_record` → `step_evaluate`
  (critério estrutural sem evidência → `unverifiable`, `submission_intent`
  gerando attempt) → `reflection_record` → `step_complete` (override) →
  `step_advance` (fim da árvore, `done: true`). Todas `status: ok`
  (2026-08-23). Este smoke test também revelou e motivou a correção real
  descrita abaixo.

### Results summary
- `internal/assessment` (puro, sem estado): `PrepareFeedback` monta o
  pacote de contexto (objetivo, escopo, pergunta, referências fixas de
  rubrica) sem gerar prosa; `Resolve` deriva o verdict de critérios
  `structural` pela presença de evidência (nunca confiando no verdict do
  chamador para esse kind) e preserva o julgamento do chamador para
  critérios qualitativos, exigindo evidência+rubrica via
  `learning.NewCriterionResult` (R4).
- `internal/session` ganhou `FeedbackPrepare`, `FeedbackRecord`,
  `StepEvaluate`, `ReflectionRecord`, `StepComplete` e `StepAdvance`,
  todos seguindo o padrão append-then-validate e idempotência por
  `request_id` já estabelecido. `StepEvaluate` registra `evaluation_recorded`
  e, quando `submission_intent`, um segundo evento `attempt_submitted`
  independentemente idempotente (Decision 4). `StepComplete` valida a
  política autoral do passo (`requires_positive_evaluation` via
  `record.cleanEvaluation`, `requires_user_confirmation` via um parâmetro
  `confirm` explícito) antes de tentar a transição de domínio, e sempre
  registra `override`/`clean_evaluation`/`confirmed` no evento (R8).
  `StepAdvance` caminha a árvore autoral em pré-ordem
  (`internal/session/progression.go`), descendo automaticamente em nós com
  um único filho, retornando opções sem mutar nada quando há ramificação, e
  `Done: true` ao esgotar a árvore.
- Corrigido durante a implementação (achado pelos próprios testes, não
  pelo smoke test): o payload de `criteria` no evento e no envelope MCP
  vazava nomes de campo Go em PascalCase (`EvidenceID`, `RubricRef`, ...)
  em vez de snake_case, por `learning.CriterionResult` não ter tags JSON.
  Corrigido adicionando tags `json:"..."` diretamente no tipo (mesmo
  precedente de `eventstore.Event`), e o mesmo em `session.AdvanceOption`.

### Requirement trace
- R1 [satisfied] test:TestFeedbackPrepareAssemblesPacketWithoutMutating test:TestPrepareFeedbackAssemblesPacketFromStep
- R2 [satisfied] test:TestFeedbackRecordNeverCompletesOrAdvances
- R3 [satisfied] test:TestStepEvaluateStructuralCriterionDrivesBlockingFailure test:TestResolveQualitativeCriterionRequiresEvidenceAndRubric
- R4 [satisfied] test:TestResolveQualitativeCriterionRequiresEvidenceAndRubric report:internal/learning/evaluation.go (NewCriterionResult)
- R5 [satisfied] report:internal/learning/evaluation.go (VerdictNotApplicable) test:TestStepEvaluateStructuralCriterionDrivesBlockingFailure
- R6 [satisfied] report:internal/learning/evaluation.go (FindingSeverity, HasBlockingFailure)
- R7 [satisfied] test:TestStepCompleteRequiresPositiveEvaluationPolicy test:TestStepAdvanceReachesTheEndOfTheTree
- R8 [satisfied] report:internal/session/service.go (StepComplete payload: override, clean_evaluation, confirmed)
- R9 [satisfied] test:TestReflectionRecordDoesNotChangeStepState report:internal/learning/evidence.go (Reflection.CompetencyID, Assessment)
- R10 [satisfied] test:TestContractListsExactlyTheMinimalToolSlice

### Known gaps
- Checks determinísticos reais (compilação, testes, vet, race, etc. de
  §21.1) serão conectados em safe-check-executor; até lá, `step_evaluate`
  só sabe dizer "evidência foi citada ou não" para critérios `structural`
  (Decision 2).
- **`eventstore.Store.Append`'s idempotency cache (`s.seen`) is keyed
  globally by `request_id`, not by `(stream_id, request_id)`** (found
  while designing StepEvaluate's two-event append, Decision 4). Two
  different sessions reusing the same client-chosen `request_id` could
  receive each other's cached event. This is a pre-existing characteristic
  of the already-closed, sealed `local-event-store` spec, not introduced
  here — spawned as its own spec, `eventstore-idempotency-scope`, rather
  than reopening sealed evidence within this one.
- `learning.Attempt` (o objeto de domínio) não é materializado em memória
  por `StepEvaluate` — apenas o evento de auditoria é persistido, já que a
  resolução de CONTEÚDO de evidência (não apenas do ID) pertence a uma
  spec futura (workspace-observation-baselines/safe-check-executor).
- `step_advance` só caminha dentro da primeira `layer` com `macro_steps`
  (Decision 5); não há noção de avançar entre `layers` diferentes ainda.

## 7. Final Report

### Delivered scope
Feedback consultivo, avaliação híbrida (determinística por presença de
evidência + qualitativa por julgamento citado), reflexão, conclusão
guiada por política autoral e avanço por travessia de árvore, expostos
como `feedback_prepare`, `feedback_record`, `step_evaluate`,
`reflection_record`, `step_complete` e `step_advance`. Nenhuma execução de
check ou observação de workspace (non-goals preservados).

### Files and modules changed
- `internal/assessment/{rubric,criteria}.go` + testes (criados)
- `internal/learning/evaluation.go` (VerdictNotApplicable, FindingSeverity, CriterionResult+NewCriterionResult, tags JSON), `evidence.go` (Reflection.CompetencyID/Assessment), `errors.go` (3 novos códigos), `policy_test.go` (ajustado)
- `internal/session/service.go` (6 novos métodos + cleanEvaluation), `progression.go` (criado: advanceFrom/rootSteps), `progression_test.go`, `assessment_test.go` (criados)
- `internal/application/assessment.go` (criado)
- `internal/mcpserver/assessment_tools.go` (criado), `server.go`, `envelope.go` (5 novos progress_effect), `contract_test.go`, `assessment_contract_test.go` (criado)

### Validation executed
- Command: go test ./... -race
- Result: `ok` em todos os pacotes com testes, 100 funções de teste, sem data races.
- Command: smoke test real via mcp.CommandTransport contra `codinho serve` e packs/go-first-steps.yaml
- Result: fluxo completo feedback→evaluate→reflect→complete→advance, todos `status: ok`.

### Residual risks
- Rubricas (`rubric://idiomatic-go`, `rubric://technical-communication`) são apenas referências fixas; seu conteúdo real vive fora deste servidor (`.agents/skills/codinho/references/feedback-rubric.md`) e não é validado por este código.
- O cache de idempotência por `request_id` (`hintResults`/`detourResults`, e agora implicitamente via `request_id + ":attempt"`) cresce sem limite por processo — mesmo risco já registrado em assistance-hints-detours.
- Ver Known Gaps: colisão potencial de `request_id` entre sessões distintas no eventstore.

### Follow-ups
- [covered: v1-integrated-acceptance] Comprovar separação operacional em E2E.
- [spawned: eventstore-idempotency-scope] Corrigir `eventstore.Store.Append`'s `s.seen` para ser indexado por `(stream_id, request_id)` em vez de `request_id` global — pré-existente de local-event-store, superfície aqui.
- [covered: safe-check-executor] Conectar execução real de checks determinísticos a `step_evaluate`'s critérios `structural`.
