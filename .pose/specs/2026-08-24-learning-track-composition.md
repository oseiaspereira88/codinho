---
slug: learning-track-composition
status: done
created_at: 2026-08-24
completed_at: 2026-09-09
supersedes:
depends_on: curriculum-graph-path-recommendation, session-orchestration-disclosure, mastery-review-scheduling, session-recovery-version-pinning, session-tree-progression, catalog-publication-integrity
priority: 65
components: curriculum, recommendations, sessions, mcp-server
delivers: capability:learning-track-composition
---

# Spec: learning-track-composition

## 1. Intent

### Goal
Substituir a seleção por valor único (um tema/competência por chamada) por
três modos de seleção determinísticos e explícitos ao aluno: trilha completa
pré-existente, busca/recomendação por conjunto de temas com sinal de
cobertura, e composição automática de trilha de N desafios quando nenhum
item único cobre tudo o que foi pedido.

### Business value
Hoje um pedido por múltiplos assuntos ao mesmo tempo não tem como declarar
que a cobertura é parcial — o aluno recebe poucos resultados ou nenhum, sem
saber por quê. E a promessa de "trilha completa" de uma linguagem (Go na V1)
nunca chegou a ser implementada em `session_start`, mesmo estando na spec
que a antecedeu. Resolver os dois fecha o gap entre o que o domínio já
modela (trilhas, grafo de precedência) e o que o aluno de fato consegue
escolher.

### Constraints
- Ver `.pose/adr/2026-08-23-multi-subject-selection-and-path-composition.md`.
- Nenhum modo é escolhido automaticamente no lugar de outro; a escolha é
  sempre do aluno (ADR Decision 5).
- Composição de trilha usa só o grafo de precedência já existente
  (`requires`, `recommended_before`) — nenhuma heurística nova de ordenação.
- Determinismo: mesma versão de catálogo + mesmo pedido produz sempre a
  mesma cobertura e a mesma composição.

### Non-goals
- Gerar conteúdo novo (isso é `agent-authored-catalog-drafts`).
- Substituir o grafo/relações de `curriculum-graph-path-recommendation`;
  estender membership de Track somente conforme R7 e revisão do ADR.

## 2. Requirements

### Functional
- R1: `curriculum.Query` e `recommendation.Objective` aceitam um conjunto de
  temas/competências (`ThemeIDs`/`CompetencyIDs`), substituindo os campos
  singulares atuais.
- R2: Toda busca/recomendação multi-assunto retorna `coverage: total |
  partial | none`, explicando quais assuntos cada resultado cobre.
- R3: `session_start` aceita `track_id` de uma trilha pré-existente
  (ex.: "Go completo") e conduz o aluno por todos os seus desafios em ordem,
  reaproveitando o lifecycle de sessão já existente.
- R4: Quando a cobertura é `partial`, o serviço de recomendação compõe uma
  sequência de N desafios existentes cuja união cobre os assuntos pedidos,
  ordenada pelo grafo de precedência, e a expõe como uma trilha ad hoc
  (mesma forma de retorno de R3, mas sem `track_id` autorado).
- R5: `catalog_search` e `learning_path_recommend` documentam os três modos
  (trilha completa, N-assuntos, desafio único) na resposta, sem esconder
  nenhum atrás de heurística implícita.
- R6: Persistir a sequência aceita, versões e cursor da trilha; retomar após restart e update de packs sem recompor silenciosamente o caminho.
- R7: Definir membership ordenada de Track, hoje limitada a ID/Title/Themes, e um identificador estável para aceitar a composição ad hoc; rejeitar seletores conflitantes, IDs desconhecidos e conjuntos vazios/excessivos.
- R8: Retornar assuntos cobertos e ausentes mesmo quando a união de todos os candidatos não cobre o pedido; nunca rotular uma trilha incompleta como cobertura total.
- R9: Preservar os campos MCP singulares existentes com tradução para conjuntos; rejeitar combinações ambíguas. Atualizar também a CLI, que usa Query.Theme, e a skill.

### Non-functional
- Consultas multi-assunto mantêm o mesmo orçamento de latência de
  `curriculum-graph-path-recommendation` (p95 < 100ms após indexação).
- Ranking e composição continuam com tie-break determinístico (ID).

### Security
- Conjunto de temas/competências no `Query` mantém o mesmo limite de
  tamanho/custo já aplicado a `Query.Text`.

### Compatibility
- Remover os campos singulares é uma mudança interna (nenhum consumidor
  externo além do MCP server); o contrato MCP externo muda de forma aditiva
  (novos campos, nenhum campo antigo reaproveitado com significado
  diferente).

## 3. Technical Plan

### Affected areas
- internal/curriculum/ (search.go, graph.go)
- internal/recommendation/ (model.go, service.go)
- internal/session/ (lifecycle sobre trilha)
- internal/mcpserver/ (session_tools.go, catalog e recommendation tools)

### Artifacts
- modified: internal/curriculum/search.go
- modified: internal/curriculum/model.go
- modified: internal/recommendation/model.go
- modified: internal/recommendation/service.go
- modified: internal/session/service.go
- modified: internal/mcpserver/session_tools.go
- modified: internal/application/recommendation.go
- modified: internal/application/catalog.go
- modified: internal/cli/catalog.go
- modified: internal/mcpserver/catalog_tools.go
- modified: internal/mcpserver/recommendation_tools.go
- modified: schemas/pack.schema.json
- modified: .agents/skills/codinho/SKILL.md

- created: internal/curriculum/selection.go
- created: internal/curriculum/tracks.go
- created: internal/curriculum/selection_test.go
- created: internal/session/tracks.go
- created: internal/session/tracks_test.go
- created: internal/mcpserver/selection.go
- created: internal/mcpserver/track_contract_test.go
- created: cmd/codinho/track_integration_test.go
- modified: internal/curriculum/index.go
- modified: internal/curriculum/validator.go
- modified: internal/session/recovery.go
- modified: internal/session/progression.go
- modified: internal/mcpserver/errors.go
- modified: internal/curriculum/search_test.go
- modified: internal/recommendation/service_test.go
- modified: internal/application/recommendation_test.go
- modified: docs/content-authoring.md
- modified: .pose/indexes/validation-matrix.json
- modified: .pose/adr/2026-08-23-multi-subject-selection-and-path-composition.md

- modified: internal/mcpserver/envelope.go

- modified: internal/curriculum/search_bench_test.go
- modified: internal/mcpserver/contract_test.go

- modified: internal/application/evaluation_evidence.go
- modified: internal/application/evaluation_evidence_test.go

### Delivery targets
- capability:learning-track-composition module:cmd/codinho profile:composed-capability entrypoint:cmd/codinho/main.go

### API/contract changes
- `session_start` ganha `track_id` como alternativa a `challenge_id`.
- `catalog_search`/`learning_path_recommend` ganham `theme_ids`/
  `competency_ids` (conjunto) e campo `coverage` na resposta.

### Data/storage changes
A recomendação permanece consultiva e pode ser calculada em memória. Após
aceite explícito, a sessão persiste a sequência e suas versões. Definir
membership de Track e atualizar o ADR existente antes de implementar; o
modelo atual contém apenas ID/Title/Themes e não codifica a sequência.

### Technical risks
- Composição de trilha pode produzir sequências longas/pedagogicamente
  ruins para muitos assuntos simultâneos — ver gatilho de revisão do ADR.

## 4. Tasks

### Planning
- [x] Mapear e migrar todos os consumidores singulares, incluindo
      internal/cli/catalog.go, preservando o contrato MCP externo.

### Implementation
- [x] Trocar `Query.Theme`/`Query.Competency` por `ThemeIDs`/`CompetencyIDs`.
- [x] Trocar `Objective.ThemeID`/`CompetencyID` por versões em conjunto.
- [x] Implementar cálculo de cobertura (`total`/`partial`/`none`).
- [x] Implementar composição de trilha ad hoc a partir do grafo de
      precedência quando a cobertura for `partial`.
- [x] Implementar lifecycle de sessão sobre `track_id` pré-existente.
- [x] Atualizar contrato MCP (`session_start`, `catalog_search`,
      `learning_path_recommend`) e seus golden tests.

### Validation
- [x] Teste de cobertura `total` (um item cobre todos os assuntos pedidos).
- [x] Teste de cobertura `partial` com composição de trilha correta pelo
      grafo de precedência.
- [x] Teste de cobertura `none`.
- [x] Teste de sessão completa sobre `track_id` percorrendo todos os
      desafios em ordem.
- [x] Suíte completa de internal/curriculum, internal/recommendation,
      internal/session, internal/mcpserver sem regressão.

## 5. Decisions

### Decision 2
- Date: 2026-09-09
- Decision: extensão do ADR existente com membership, identidade de composição e replay fixado.
- Rationale: consumir knowledge:adr-multi-subject-selection-and-path-composition-review e knowledge:planning-audit-2026-09; impedir recomposição após aceite e preservar o MCP legado.
- Consequences: conjuntos limitados, seletores exclusivos e cobertura ausente explícita.


### Decision 1
- Date: 2026-08-23
- Context: como sinalizar cobertura parcial sem inventar conteúdo.
- Options considered: (a) devolver lista vazia/curta sem explicação (status
  quo); (b) compor trilha ad hoc de itens reais; (c) sintetizar um item
  virtual concatenando trechos de vários desafios.
- Decision: (b).
- Rationale: preserva granularidade de sessão e integridade de evidência de
  maestria — ver ADR `multi-subject-selection-and-path-composition`.
- Consequences: toda resposta multi-assunto precisa carregar `coverage`
  explicitamente, nunca implícito pelo tamanho da lista.

## 6. Validation

### Strategy
Corpus de catálogo de teste com temas sobrepostos e não sobrepostos,
cobrindo os três valores de `coverage` e uma trilha autorada completa.

| Cenário obrigatório | Comando | Evidência |
|---|---|---|
| Conjuntos e composição | go test ./internal/curriculum ./internal/recommendation ./internal/application | total/partial/none, ausentes, ordem, limites, determinismo |
| Sessão e replay | go test ./internal/session | trilha inteira, restart, packs alterados, retry e cursor |
| Contrato e composição real | go test ./internal/mcpserver ./cmd/codinho | seletores conflitantes, campos antigos, stdio real |
| Regressão e segurança | pose validate --strict --json .pose/results/delivery-validation.json | race, vet, build, scanner, contratos e integração passam |

### Deterministic checks
- Test: go test ./internal/curriculum/... ./internal/recommendation/... ./internal/session/... ./internal/mcpserver/...
- Lint: gofmt -l internal/curriculum internal/recommendation internal/session internal/mcpserver
- Typecheck: go vet ./internal/curriculum/... ./internal/recommendation/... ./internal/session/... ./internal/mcpserver/...
- Build: go build ./...

### Execution log
- 2026-09-09: implementados conjuntos, cobertura, composição por dependências, membership e lifecycle fixado de trilhas. Contratos MCP legados mantidos; CLI e skill atualizadas.
- 2026-09-09: matriz com 21 checks passou; revisão acrescentou proteção de custo e ciclo misto e testes de conclusão explícita. Candidato ae0f7ae validado: 21/21 checks, sem skips. Revisão e closeout aprovados; consultar report:.pose/reports/2026-09-09-review-learning-track-composition.md.

### Results summary
Busca, recomendação e sessões sobre trilhas implementadas; evidência final registrada no relatório de revisão.

### Requirement trace
- R1 [satisfied] test:TestSelectionCoverageAndComposition report:internal/recommendation/service_test.go
- R2 [satisfied] test:TestSelectionCoverageAndComposition test:TestContractTrackSelection
- R3 [satisfied] test:TestTrackExplicitCompletionTraversesAllChallenges test:TestTrackCompositionOverRealStdio
- R4 [satisfied] test:TestSelectionCoverageAndComposition test:TestCompositionBoundsAndMixedDependencyCycle
- R5 [satisfied] test:TestContractTrackSelection report:.agents/skills/codinho/SKILL.md
- R6 [satisfied] test:TestTrackSessionReplayAndBoundaries test:TestTrackCompositionOverRealStdio
- R7 [satisfied] test:TestTrackSelectorsAndComposition test:TestTrackMembershipValidation test:TestSelectionBounds
- R8 [satisfied] test:TestSelectionCoverageAndComposition test:TestContractTrackSelection
- R9 [satisfied] test:TestContractTrackSelection report:internal/cli/catalog.go report:.agents/skills/codinho/SKILL.md

### Known gaps
Playtest humano nos dois hosts continua no aceite V1. A composição inclui todos os candidatos e pré-requisitos, sem prometer caminho mínimo ou qualidade pedagógica autorreferendada.

## 7. Final Report

### Delivered scope
Conjuntos de assuntos com cobertura explícita; composição determinística; trilhas autoradas e compostas com versões/cursor fixados; compatibilidade MCP, CLI e skill; isolamento de evidências entre desafios.

### Files and modules changed
- Artefatos reconciliados na seção Technical Plan; inclui persistência/replay, testes de contrato e integração stdio.

### Validation executed
- Command: pose lint-spec learning-track-composition --ready-check
- Result: passed; usar relatório de revisão e validação estruturada para o candidato final.

### Residual risks
- Nenhum adicional além do já descrito em Technical risks.

### Follow-ups
- [covered: v1-integrated-acceptance] Provar seleção, execução e retomada de trilha autorada e composta nos dois hosts.
