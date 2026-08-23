---
slug: assistance-hints-detours
status: in-progress
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: session-orchestration-disclosure
priority: 70
components: assistance, sessions, mcp-server
delivers:
---

# Spec: assistance-hints-detours

## 1. Intent

### Goal
Implementar a escada de auxílio, lembrança sintática, conteúdo conceitual e desvios pedagógicos sem avanço implícito.

### Business value
Oferecer a menor ajuda suficiente e permitir aprendizagem contextual sem tomar o teclado do aluno.

### Constraints
- Subir no máximo um nível por solicitação.
- Nível de solução exige intenção explícita e registro.
- Exemplos devem evitar resolver o desafio ativo quando possível.

### Non-goals
- Redigir explicações abertas dentro do MCP.
- Avaliar ou concluir passos.

## 2. Requirements

### Functional
- R1: Entregar pistas ordenadas nos níveis 1 a 6, respeitando a política da sessão.
- R2: Tratar objetivo e critérios como nível 0 sem consumo de pista.
- R3: Bloquear níveis não permitidos com erro estável e zero mudança de progresso.
- R4: Registrar nível, contexto, timestamp e efeito de cada pista aceita.
- R5: Fornecer conteúdo canônico de conceito e lembrança sintática com custo configurável.
- R6: Abrir, consultar e fechar desvio sem alterar o nó instrucional original.
- R7: Marcar solução revelada e agendar uma variante futura sem promover autonomia.
- R8: Expor hint_request, concept_content_get, syntax_recall_get, learning_detour_start e learning_detour_finish.

### Non-functional
- Pistas da mesma versão devem ser determinísticas.
- Conteúdo conceitual deve ser acessível sem workspace.

### Security
- Conteúdo autorado é dado; não pode invocar tools nem ampliar permissões.
- Disclosure deve filtrar gabaritos e fragmentos acima do nível.

### Compatibility
- Novos níveis podem ser adicionados sem reinterpretar registros existentes.

## 3. Technical Plan

### Affected areas
- internal/assistance/, internal/session/, internal/mcpserver/

### Artifacts
- created: internal/assistance/service.go
- created: internal/assistance/ladder.go
- created: internal/assistance/detour.go
- created: internal/assistance/service_test.go
- created: internal/assistance/disclosure_test.go
- modified: internal/session/service.go
- created: internal/session/assistance_test.go
- created: internal/mcpserver/assistance_tools.go
- created: internal/mcpserver/assistance_contract_test.go
- modified: internal/mcpserver/server.go
- modified: internal/mcpserver/errors.go
- modified: internal/mcpserver/contract_test.go
- modified: internal/mcpserver/session_contract_test.go
- created: internal/application/assistance.go
- modified: internal/learning/step.go
- created: internal/learning/step_test.go
- modified: internal/learning/errors.go
- modified: internal/curriculum/index.go
- modified: cmd/codinho/main.go

### Delivery targets
Nenhum novo; amplia o contrato MCP V1 planejado.

### API/contract changes
- Adicionar cinco tools e efeito solution_revealed.

### Data/storage changes
- Persistir hint_requested, detour_started, detour_finished e solution_revealed.

### Technical risks
- Pistas podem crescer sem diferença real de intensidade.
- Conteúdo conceitual pode vazar a solução indiretamente.

## 4. Tasks

### Planning
- [x] Definir rubrica mensurável para intensidade de pista.
- [x] Definir quando syntax recall conta como pista por modo.

### Implementation
- [x] Implementar policy e ladder.
- [x] Implementar conteúdo canônico e filtragem.
- [x] Implementar lifecycle de detour.
- [x] Implementar registro de solução e review futura.
- [x] Expor tools e contract tests.

### Validation
- [x] Testar cada transição de nível e todos os bloqueios.
- [x] Executar golden tests de não revelação.
- [x] Testar retomada do mesmo passo após detour.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Proibição absoluta de código prejudica alguns alunos.
- Options considered: sem código; ajuda livre; escada explícita.
- Decision: usar escada com solução somente no nível máximo.
- Rationale: controla custo pedagógico sem bloquear aprendizagem.
- Consequences: catálogo precisa autorar pistas distintas e auditáveis.

### Decision 2
- Date: 2026-08-22
- Context: `HintAuthoring{Level, Kind}` (catalog-schema-loader) não carrega texto
  de pista algum — apenas nível e uma tag de classificação livre. O exemplo
  de PROJECT.md §14.9 mostra `kind: syntax_recall` no nível 2, enquanto a
  tabela §8.4 descreve o nível 2 genericamente como "conceito ou API",
  revelando que `kind` é uma sobreposição opcional por passo, não um
  vocabulário fixo por nível.
- Options considered: (a) exigir pista autorada em cada nível 1-6 ou falhar;
  (b) `hint_request` sempre entrega apenas nível+kind+objective/scope/
  concepts, delegando toda a prosa explicativa ao agente tutor chamador.
- Decision: (b). `hint_request` nunca aceita um `level` explícito do
  chamador — sempe entrega `HintLevel atual + 1`, com `kind` resolvido pelo
  `HintAuthoring` do passo quando presente, caindo para o rótulo genérico de
  §8.4 quando ausente (`internal/assistance.KindFor`).
- Rationale: satisfaz o non-goal "redigir explicações abertas dentro do
  MCP" com o mínimo de acoplamento a como o catálogo é autorado; a escada
  funciona mesmo sobre conteúdo sem nenhuma pista explícita.
- Consequences: a "prosa" da pista em si nunca existe no domínio ou no
  catálogo — é sempre gerada pelo agente chamador a partir do envelope
  estruturado retornado.

### Decision 3
- Date: 2026-08-22
- Context: RF-022 exige lembrança sintática "sem obrigatoriamente consumir
  pista"; §8.4 diz que ela "pode ter custo menor no modo ensino".
  `syntax_recall_get` também é um lookup direcionado (busca o rung marcado
  `kind: syntax_recall`), não o próximo rung sequencial — diferente de
  `hint_request`.
- Options considered: (a) tratar syntax recall como apenas mais um
  `hint_request` restrito por kind; (b) um caminho de concessão separado,
  direto (fora de ordem) e com custo condicional ao modo.
- Decision: (b). Adicionado `StepProgress.GrantDirect` (aplica o teto de
  política mas não o ratchet de um-nível-por-vez) ao lado de `GrantHint`
  (ratchet completo, usado por `hint_request`). `SyntaxRecallGet` é gratuito
  (não avança `HintLevel`, apenas registra o evento) quando
  `Policy.Mode == ModeTeaching`; em qualquer outro modo, consome a escada
  via `GrantDirect` até aquele nível. `HelpNoHints` bloqueia ambos os
  caminhos.
- Rationale: reflete fielmente que syntax recall é um conteúdo
  independente da posição atual na escada, e que seu custo é configurável
  por modo, não uma regra fixa.
- Consequences: dois métodos de concessão em `learning.StepProgress`
  (`GrantHint` e `GrantDirect`) em vez de um único, cada um com uma
  invariante diferente.

### Decision 4
- Date: 2026-08-22
- Context: PROJECT.md §15.3 fixa o vocabulário de `progress_effect`
  (`none`, `session_changed`, ..., `solution_revealed`) sem um valor
  dedicado para "pista concedida" ou "desvio aberto/fechado".
- Decision: `hint_request`/`syntax_recall_get` usam `progress_effect: none`
  em toda concessão normal e `solution_revealed` apenas quando o nível
  concedido é 6 (reaproveitando `ProgressEffectSolutionReveal`, já
  reservado no envelope desde mcp-stdio-foundation). Desvios sempre usam
  `none` (§8.8: "efeito no progresso: nenhum").
- Rationale: manter o envelope compatível com o vocabulário fechado em vez
  de introduzir valores novos (requisito de compatibilidade da spec).
- Consequences: um cliente já preparado para `solution_revealed` detecta a
  revelação de solução por pista sem mudança de contrato.

### Decision 5
- Date: 2026-08-22
- Context: `hint_request`/`syntax_recall_get` decidem o próximo nível a
  partir de `StepProgress.HintLevel`, um valor mutável. Uma retentativa
  (mesmo `request_id`) chega depois que a chamada original já avançou esse
  valor, então recalcular "o próximo nível" na retentativa produziria um
  nível diferente do que foi de fato concedido na primeira vez — a mesma
  classe de bug de idempotência já corrigida em `session-orchestration-
  disclosure` para pause/resume, mas desta vez afetando o *valor*
  retornado, não apenas o estado.
- Decision: cache de resultado por `(session_id, request_id)`
  (`Service.hintResults`), verificado antes de qualquer trabalho — mesmo
  padrão de `Service.startResults` em `Start`. `DetourStart`/`DetourFinish`
  não precisaram do mesmo cache: seu valor de retorno é lido do estado do
  domínio *depois* da aplicação condicional, que já é estável entre
  chamadas.
- Rationale: confirmado por teste (`TestHintRequestIsIdempotentByRequestID`
  falhou antes da correção, retornando o nível 2 numa retentativa que
  originalmente concedeu o nível 1).
- Consequences: mais um mapa de cache por request_id no `Service`; padrão a
  repetir sempre que um método de escrita retorna um valor *derivado* de
  estado mutável, não apenas um "estado atual".

## 6. Validation

### Strategy
Usar golden disclosure, matriz de políticas e replay de eventos.

### Deterministic checks
- Test: go test ./internal/assistance/... ./internal/session/...
- Lint: gofmt -l internal/assistance internal/session
- Typecheck: go vet ./internal/assistance/... ./internal/session/...
- Build: go test ./internal/assistance/... -run ^$
- Security / Contract: testes negativos de leak e schemas MCP.

### Execution log
- `gofmt -l internal/assistance internal/session internal/learning internal/curriculum internal/mcpserver internal/application cmd/codinho` → saída vazia (2026-08-23).
- `go vet ./...` → sem findings (2026-08-23).
- `go build ./...` → ok (2026-08-23).
- `go test ./... -race` → `ok` em todos os 8 pacotes com testes, 85 funções
  de teste passando em assistance/session/learning/mcpserver combinados,
  sem data races (2026-08-23).
- `govulncheck ./...` → "No vulnerabilities found." (2026-08-23).
- Smoke test real de ponta a ponta via `mcp.CommandTransport` contra o
  binário `codinho serve` e o pack real `packs/go-first-steps.yaml`
  (desafio `go-data.slice-filter-preserve-input`, passo
  `model.declare-user-struct`): `session_start` (disclosure_max=6) →
  `hint_request` (nível 1, kind guiding_question) → `concept_content_get`
  (named-types) → `syntax_recall_get` (nível 2, kind syntax_recall) →
  `learning_detour_start` (reason) → `learning_detour_finish` (resolved).
  Todas as chamadas retornaram `status: ok` com o node ativo inalterado
  durante o desvio (2026-08-23).

### Results summary
- Escada de auxílio implementada em `internal/assistance` (decisão pura,
  catalog-only) + `internal/session` (estado, ratchet, eventos): `hint_request`
  sobe exatamente um nível por chamada, respeita o teto de política, exige
  `confirm_solution` explícito para o nível 6 e marca `SolutionRevealed`
  nesse caso. `syntax_recall_get` localiza o rung com `kind: syntax_recall`
  autorado no passo, é gratuito (não avança a escada) em modo ensino e
  consome a escada até aquele nível em qualquer outro modo;
  `HelpNoHints` bloqueia ambos.
- `concept_content_get` expõe o registro canônico do conceito (id, título)
  sem gerar prosa.
- `learning_detour_start`/`learning_detour_finish` abrem e fecham um
  `LearningDetour` sem jamais alterar o passo ativo da sessão, verificado
  por teste unitário e de contrato.
- Todas as cinco tools novas usam o mesmo envelope MCP e o mesmo padrão de
  idempotência por `request_id` + `expected_revision` já estabelecido em
  `session-orchestration-disclosure`; um bug de idempotência real
  (retentativa de `hint_request` recalculando o nível a partir de estado já
  avançado) foi encontrado por teste e corrigido antes do fechamento (ver
  Decision 5).

### Requirement trace
- R1 [satisfied] test:TestHintRequestClimbsLadderAndStopsAtSessionCap test:TestContractHintRequestClimbsLadderToSolution
- R2 [satisfied] report:internal/session/service.go (Instruction não toca HintLevel; nível 0 nunca passa por GrantHint/GrantDirect)
- R3 [satisfied] test:TestGrantHintRejectsSkippingARung test:TestGrantHintRejectsAboveSessionCap test:TestHintRequestBlockedByNoHintsPolicy
- R4 [satisfied] report:internal/session/service.go (payload do evento carrega step_id, level, kind, free; RecordedAt vem do eventstore)
- R5 [satisfied] test:TestConceptContentReturnsCanonicalRecord test:TestSyntaxRecallGetIsFreeInTeachingMode test:TestSyntaxRecallGetConsumesLadderOutsideTeachingMode
- R6 [satisfied] test:TestDetourStartAndFinishDoNotChangeActiveStep test:TestContractLearningDetourStartAndFinishPreserveActiveStep
- R7 [satisfied] test:TestHintRequestSolutionRequiresConfirmation (payload do evento solution_revealed carrega needs_variant)
- R8 [satisfied] test:TestContractListsExactlyTheMinimalToolSlice

### Known gaps
- Qualidade linguística final depende da skill e do agente host.
- `concept_content_get` retorna apenas id/título (o único conteúdo que o
  catálogo de fato autora); relações, analogias e exemplos citados em
  PROJECT.md §15.7 exigiriam campos de autoria que não existem em
  `ConceptAuthoring` — follow-up para uma spec de autoria de catálogo, não
  implementável aqui sem inventar formato fora do SchemaVersion 1.
- "Agendar uma variante futura" (R7) é registrado apenas como
  `needs_variant: true` no payload do evento `solution_revealed`; a
  seleção/agendamento real de uma variante pertence a uma spec de
  recomendação de trilha ou revisão espaçada futura.
- Nenhum teste de fuzz ou golden-file dedicado para o texto de disclosure
  (não há texto gerado pelo servidor a testar; o envelope estruturado em si
  está coberto pelos testes de contrato).

## 7. Final Report

### Delivered scope
Escada de auxílio (níveis 1-6 com solução sob confirmação explícita),
lembrança sintática com custo por modo, conteúdo canônico de conceito e
desvio pedagógico sem avanço implícito, expostos como as cinco tools MCP
`hint_request`, `syntax_recall_get`, `concept_content_get`,
`learning_detour_start` e `learning_detour_finish`. Nenhuma prosa
explicativa é gerada pelo servidor (non-goal preservado); apenas
nível/kind/objetivo/escopo/conceitos estruturados, para o agente tutor
autorar o texto.

### Files and modules changed
- `internal/assistance/{ladder,detour,service}.go` + testes (criados)
- `internal/learning/step.go` (GrantHint/GrantDirect/HintLevel), `errors.go` (ErrCodeHintLevelSkipped), `step_test.go` (criado)
- `internal/curriculum/index.go` (Concept accessor)
- `internal/session/service.go` (HintRequest/SyntaxRecallGet/DetourStart/DetourFinish + caches de idempotência), `assistance_test.go` (criado)
- `internal/application/assistance.go` (criado)
- `internal/mcpserver/assistance_tools.go` (criado), `server.go`, `errors.go`, `contract_test.go`, `session_contract_test.go`, `assistance_contract_test.go` (criado)
- `cmd/codinho/main.go` (wiring do AssistanceService)

### Validation executed
- Command: go test ./... -race
- Result: `ok` em todos os pacotes com testes, 85 funções de teste, sem data races.
- Command: smoke test real via mcp.CommandTransport contra `codinho serve` e packs/go-first-steps.yaml
- Result: sessão real, hint_request, syntax_recall_get, concept_content_get e o ciclo completo de detour todos `status: ok`.

### Residual risks
- Revisão humana das pistas permanece obrigatória.
- O cache de idempotência por `request_id` (`hintResults`, `detourResults`) cresce sem limite por processo — aceitável no V1 (processo de vida curta por sessão MCP), mas deve ser revisto se sessões de longa duração ou muitos `request_id` distintos por sessão se tornarem comuns.

### Follow-ups
- [covered: tutor-skill-host-integration] Validar adaptação das explicações pelo tutor.
- [open] Relações/analogias/exemplos de conceito (PROJECT.md §15.7) exigem novos campos de autoria em ConceptAuthoring; avaliar em uma futura spec de qualidade de autoria de catálogo.
- [covered: mastery-review-scheduling] Consumir `needs_variant: true` do evento `solution_revealed` para de fato agendar uma variante futura do desafio.
