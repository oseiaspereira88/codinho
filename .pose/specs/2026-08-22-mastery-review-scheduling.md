---
slug: mastery-review-scheduling
status: done
created_at: 2026-08-22
completed_at: 2026-08-23
supersedes:
depends_on: local-event-store, feedback-evaluation-progression
priority: 110
components: mastery, progress
delivers:
---

# Spec: mastery-review-scheduling

## 1. Intent

### Goal
Projetar domínio multidimensional por competência e agendar revisões espaçadas por regras transparentes.

### Business value
Distinguir conclusão momentânea de autonomia, retenção e transferência real.

### Constraints
- Não usar score opaco ou gamificação por volume.
- Solução revelada não prova autonomia.
- Retenção exige outra data; transferência exige contexto diferente.

### Non-goals
- Machine learning adaptativo.
- Ranking entre pessoas.

## 2. Requirements

### Functional
- R1: Projetar compreensão, sintaxe, implementação guiada, autônoma, depuração, explicação, retenção e transferência.
- R2: Representar não observada, introduzida, demonstra com ajuda, sem ajuda, retida e transferida.
- R3: Consumir tentativas, pistas, avaliação, reflexão, data e variante sem reescrever evidência.
- R4: Impedir promoção por uma única conclusão guiada ou por override.
- R5: Agendar revisões iniciais em 1, 3, 7, 14 e 30 dias, adaptadas por resultado.
- R6: Selecionar revisão vencida com prioridade explicável e sem bloquear sessão livre.
- R7: Expor progress_get, review_due e mastery_evidence_record.
- R8: Recalcular projeções deterministicamente a partir do log.

### Non-functional
- Cada mudança de domínio deve listar evidências causais.
- Regras e intervalos devem ser configuráveis e versionados.

### Security
- Progresso é dado local potencialmente sensível e não deve sair do dispositivo por padrão.

### Compatibility
- Mudanças na regra devem preservar projeção histórica por versão.

## 3. Technical Plan

### Affected areas
- internal/mastery/, internal/application/, internal/mcpserver/

### Artifacts
- created: internal/mastery/model.go
- created: internal/mastery/rules.go
- created: internal/mastery/scheduler.go
- created: internal/mastery/projector.go
- created: internal/mastery/rules_test.go
- created: internal/mastery/scheduler_test.go
- created: internal/mastery/projector_test.go
- created: internal/application/progress.go
- created: internal/application/progress_test.go
- created: internal/mcpserver/progress_tools.go
- created: internal/mcpserver/progress_contract_test.go
- modified: internal/mcpserver/server.go
- modified: internal/mcpserver/errors.go
- modified: internal/mcpserver/envelope.go
- modified: internal/mcpserver/contract_test.go
- modified: cmd/codinho/main.go

### Delivery targets
Nenhum novo; amplia o contrato MCP V1 planejado.

### API/contract changes
- Adicionar três tools de progresso e schemas de explicabilidade.

### Data/storage changes
- Persistir mastery_projected e review_scheduled como projeções reconstruíveis.

### Technical risks
- Regras rígidas podem recomendar prática demais.
- Variantes mal classificadas podem simular transferência.

## 4. Tasks

### Planning
- [x] Definir tabela evidência versus dimensão e estado.
- [x] Definir cálculo de revisão e desempate.

### Implementation
- [x] Implementar projector versionado.
- [x] Implementar scheduler e motivos.
- [x] Implementar queries de progresso e revisão.
- [x] Expor tools e eventos.
- [x] Criar fixtures longitudinais.

### Validation
- [x] Testar replay, mudança de regra e datas com clock fake.
- [x] Testar que ajuda, solução e override não promovem indevidamente.
- [x] Testar retenção e transferência em variantes.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Um número único esconde a natureza da lacuna.
- Options considered: XP; score 0–100; vetor de estados explicável.
- Decision: vetor de dimensões e estados derivados de evidência.
- Rationale: permite recomendação precisa e auditável.
- Consequences: UI textual é mais extensa, porém mais honesta.

### Decision 2
- Date: 2026-08-23
- Context: R3 diz "consumir tentativas, pistas, avaliação, reflexão,
  data e variante sem reescrever evidência". Duas leituras possíveis:
  (a) o projetor varre automaticamente TODOS os eventos de sessão
  (attempt_submitted, hint_requested, evaluation_recorded, ...) e infere
  sozinho a qual dimensão cada um pertence; (b) `mastery_evidence_record`
  é uma chamada explícita do tutor — como `feedback_record`/
  `reflection_record` já são — que CITA um `evidence_id` já existente
  (de uma tentativa, avaliação ou reflexão), sem duplicá-lo, e declara
  explicitamente `competency_id`, `dimension`, `variant`, `help_used` e
  `solution_revealed`.
- Options considered: (a) e (b) acima.
- Decision: (b). PROJECT.md §15.10 já descreve `mastery_evidence_record`
  como "escrita local" que "registra evidência válida" — o mesmo padrão
  de citar evidência que `step_evaluate` já usa para critérios
  qualitativos, não um inferidor automático de conteúdo de evento.
- Rationale: (a) exigiria um mapeamento evento→dimensão totalmente
  inventado por esta spec (PROJECT.md não especifica esse algoritmo),
  arriscando heurísticas silenciosas e não auditáveis; (b) mantém o
  mesmo princípio já estabelecido em feedback-evaluation-progression
  ("o servidor não finge ter produzido raciocínio que não tem") — quem
  decide que uma tentativa demonstra "sintaxe" vs. "implementação
  autônoma" é o tutor (que viu a interação completa), citando a
  evidência real; o servidor só aplica as regras determinísticas de
  promoção (R4, invariantes 8 e 10) sobre essa declaração.
- Consequences: "sem reescrever evidência" significa literalmente isso:
  `mastery_evidence_record` nunca cria uma nova `Evidence`/`Observation`
  de conteúdo, apenas referencia `evidence_id` de uma já registrada
  (invariante 9 preservada).

### Decision 3
- Date: 2026-08-23
- Context: R8 exige recalcular projeções deterministicamente a partir
  do log — não apenas manter estado em memória com o log como auditoria
  secundária (a solução mais simples, já usada em workspace-observation-
  baselines Decision 3 para o mesmo tipo de trade-off).
- Options considered: (a) seguir o precedente de workspace-observation-
  baselines (estado em memória, log apenas para auditoria, gap
  documentado); (b) reconstruir a projeção sempre a partir de replay
  real do log, com um motor de regras versionado.
- Decision: (b), porque R8 é explícito e testável aqui (diferente de
  workspace-observation-baselines, onde a durabilidade completa não era
  requisito textual). Todo progresso vive em UM stream global
  `"mastery"` do eventstore (progresso não é por sessão); cada evento
  `mastery_projected` carrega o `Signal` completo que o originou MAIS
  `rule_version`, então `ProgressService` reconstrói projeções por
  replay + `mastery.AdvanceState(rule_version, ...)` a cada leitura, em
  vez de confiar em um cache que poderia divergir da regra atual.
- Rationale: entrega o requisito textual real em vez de documentar mais
  um gap; o custo (reprocessar um log local pequeno a cada leitura) é
  aceitável para volume de um único aluno local.
- Consequences: `mastery_evidence_record` opera sobre a revisão do
  stream global `"mastery"` (obtida via `progress_get`), não a revisão
  de uma sessão — chamadas concorrentes de diferentes sessões
  serializam nesse único stream, aceitável para um instalador local
  single-user na V1.

### Decision 4
- Date: 2026-08-23
- Context: §21.3/§21.4 e as invariantes 8 e 10 descrevem os 6 estados
  como uma progressão, mas não detalham a máquina de transição exata.
- Options considered: contagem de repetições; matriz configurável por
  dimensão; ladder sequencial único aplicado a todas as dimensões.
- Decision: ladder sequencial único: `not_observed → introduced` em
  qualquer primeiro contato (sucesso ou não); daí, só sucesso sem
  solução revelada promove; `help_used=true` nunca ultrapassa
  `demonstrates_with_help` (satisfaz R4/invariante 10 literalmente);
  `demonstrates_without_help` só vira `retained` com nova evidência
  autônoma em outra data civil; `retained` só vira `transferred` com
  evidência autônoma em variante inédita para essa competência+dimensão.
  Falha ou solução revelada nunca regridem nem promovem.
- Rationale: é a leitura mais literal dos textos de R4/R5/invariantes 8
  e 10, sem inventar contagens ou pesos não especificados.
- Consequences: uma dimensão nunca "pula" `retained` para chegar em
  `transferred` sem antes ter sido retida — leitura sequencial estrita
  dos estados descritivos listados em ordem no PROJECT.md.

## 6. Validation

### Strategy
Usar timelines sintéticas, clocks determinísticos e property tests de monotonicidade condicional.

### Deterministic checks
- Test: go test ./internal/mastery/... ./internal/application/... ./internal/mcpserver/...
- Lint: gofmt -l internal/mastery internal/application internal/mcpserver cmd/codinho
- Typecheck: go vet ./...
- Build: go build ./...
- Security / Contract: go test -race; govulncheck; progresso é local (nenhuma chamada de rede em todo o pacote).

### Execution log
- `gofmt -l ...` (todos os pacotes tocados) → saída vazia (2026-08-23).
- `go vet ./...` (linux e `GOOS=windows`) → sem findings (2026-08-23).
- `go build ./...` → ok (2026-08-23).
- `go test ./... -race` → `ok` em todos os pacotes com testes, sem data races (2026-08-23).
- `govulncheck ./...` → "No vulnerabilities found." (2026-08-23).
- Smoke test real via `mcp.CommandTransport` contra `codinho serve` e
  `packs/go-first-steps.yaml`: `session_start` → `step_evaluate`
  (citando um evidence_id) → `progress_get` (revisão inicial 0,
  `competencies: {}`) → `mastery_evidence_record` citando essa mesma
  evidência (`competency_id: slice-filter`, `dimension:
  autonomous_implementation`, `success: true`) → `state:
  demonstrates_without_help` de primeira (sucesso autônomo sem ajuda
  não passa por `introduced`/`demonstrates_with_help`, Decision 4) →
  `progress_get` reflete exatamente esse estado com
  `evidence_count: 1` → `review_due` não retorna nada (revisão agendada
  1 dia à frente) → segunda `mastery_evidence_record` na mesma variante
  e mesmo instante (mesmo dia civil) permanece em
  `demonstrates_without_help`, provando que retenção exige data
  diferente mesmo pelo protocolo real, não só em teste unitário
  (2026-08-23).

### Results summary
- `internal/mastery` (puro, sem estado): `AdvanceState` implementa o
  ladder sequencial de 6 estados sobre 8 dimensões independentes;
  `help_used=true` nunca ultrapassa `demonstrates_with_help` mesmo após
  10 sucessos guiados seguidos (requirement R4, invariante 10, testado
  explicitamente); `solution_revealed` e falha nunca promovem nem
  regridem (invariante 8); `retained` exige uma data civil diferente da
  última evidência e `transferred` exige uma variante nunca vista para
  aquela competência+dimensão. `NextSchedule` implementa a ladder 1-3-
  7-14-30 dias (requirement R5), cresce geometricamente após o fim da
  ladder e reduz (nunca zera) após falha. `FoldProjections`/
  `FoldSchedules` replay uma lista ordenada de `RecordedSignal` e
  reproduzem exatamente o mesmo resultado que produção sequencial via
  `AdvanceState`/`NextSchedule` — testado diretamente e via reabertura
  real do event store (ver `TestProgressServiceRecomputesAfterReopeningTheStore`).
- `internal/application.ProgressService` mantém TODO o progresso num
  único stream global `"mastery"` do eventstore (progresso não é por
  sessão, Decision 3); `RecordEvidence` nunca cria uma nova `Evidence`
  de conteúdo — só referencia um `evidence_id` já existente (Decision
  2, invariante 9); cada chamada bem-sucedida sem solução revelada
  grava dois eventos (`mastery_projected` e `review_scheduled`,
  PROJECT.md §18.3) na mesma revisão do stream. `Progress`/`ReviewDue`
  nunca leem um cache: recomputam via replay completo a cada chamada
  (requirement R8), verificado inclusive fechando e reabrindo o
  eventstore entre a escrita e a leitura.
- `internal/mcpserver` expõe `progress_get`, `review_due` e
  `mastery_evidence_record` (requirement R7) sem `session_id`: progresso
  é global ao aluno local, não por sessão. `review_due` nunca bloqueia
  `session_start` (requirement R6) — não há acoplamento algum entre os
  dois tools.

### Requirement trace
- R1 [satisfied] report:internal/mastery/model.go (8 constantes de Dimension)
- R2 [satisfied] report:internal/mastery/model.go (6 constantes de State)
- R3 [satisfied] report:internal/application/progress.go (RecordEvidence cita evidence_id, nunca recria evidência — Decision 2)
- R4 [satisfied] test:TestAdvanceStateGuidedSuccessNeverExceedsDemonstratesWithHelp test:TestContractMasteryEvidenceRecordGuidedNeverExceedsDemonstratesWithHelp
- R5 [satisfied] test:TestNextScheduleFirstReviewUsesFirstLadderRung test:TestNextScheduleSuccessClimbsTheLadder test:TestNextScheduleFailureShrinksButNeverBelowMinimum
- R6 [satisfied] report:internal/mastery/projector.go (DueSchedules, ordenação explicável) test:TestDueSchedulesOrdersMostOverdueFirstWithDeterministicTieBreak
- R7 [satisfied] report:internal/mcpserver/progress_tools.go (progress_get, review_due, mastery_evidence_record)
- R8 [satisfied] test:TestProgressServiceRecomputesAfterReopeningTheStore test:TestFoldProjectionsReproducesSequentialAdvanceState

### Known gaps
- Calibração dos intervalos (ladder 1-3-7-14-30, fator de crescimento
  1.5x, redução 0.5x) precisará de evidência de uso real — heurísticas
  conservadoras assumidas nesta entrega (residual risk já previsto).
- Mapeamento de qual `competency_id`/`dimension`/`variant` uma
  interação específica representa é decisão do chamador (tutor), não
  verificado pelo servidor além do formato — um tutor mal orientado
  poderia citar a dimensão errada; mitigação é a skill do tutor, não
  código aqui (mesma fronteira de responsabilidade que
  `feedback_prepare`/`step_evaluate` já estabelecem para julgamento
  qualitativo).

## 7. Final Report

### Delivered scope
Projeção multidimensional de domínio por competência (8 dimensões × 6
estados, regras versionadas e determinísticas) e agendamento de revisão
espaçada (ladder 1-3-7-14-30 dias, adaptado por resultado) — expostos
como `progress_get`, `review_due` e `mastery_evidence_record`. Nenhuma
promoção por conclusão guiada isolada, override ou solução revelada
(invariantes 8 e 10 preservadas); toda projeção é recomputável do zero a
partir do log de eventos.

### Files and modules changed
- `internal/mastery/{model,rules,scheduler,projector}.go` + testes (criados)
- `internal/application/progress.go` + `progress_test.go` (criados)
- `internal/mcpserver/progress_tools.go` + `progress_contract_test.go` (criados), `server.go`, `errors.go`, `envelope.go`, `contract_test.go` (modificados)
- `cmd/codinho/main.go` (injeta `ProgressService`)

### Validation executed
- Command: go test ./... -race
- Result: `ok` em todos os pacotes com testes, sem data races.
- Command: govulncheck ./...
- Result: "No vulnerabilities found."
- Command: smoke test real via mcp.CommandTransport contra `codinho serve`, `packs/go-first-steps.yaml`
- Result: step_evaluate → mastery_evidence_record → progress_get → review_due consistentes ponta a ponta, incluindo a regra de "mesma data civil" via protocolo real.

### Residual risks
- Recomendações iniciais (ladder, fatores de crescimento/redução) usam heurísticas conservadoras; calibração real depende de uso longitudinal (Known Gaps).
- O mapeamento evidência→dimensão/competência é responsabilidade do tutor que chama `mastery_evidence_record`, não verificado semanticamente pelo servidor (mesma fronteira já estabelecida para julgamento qualitativo em `step_evaluate`).

### Follow-ups
- [covered: v1-integrated-acceptance] Validar retenção e transferência em sessões reais.
