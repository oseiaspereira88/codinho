---
slug: learning-track-composition
status: draft
created_at: 2026-08-24
completed_at:
supersedes:
depends_on: curriculum-graph-path-recommendation, session-orchestration-disclosure, mastery-review-scheduling
priority: 65
components: curriculum, recommendations, sessions, mcp-server
delivers:
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
- Mudar o formato de autoria de trilha (`Track`/relations) já validado por
  `curriculum-graph-path-recommendation`.

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

### Delivery targets
Nenhum novo; amplia o contrato MCP V1 já entregue por `mcp-stdio-foundation`.

### API/contract changes
- `session_start` ganha `track_id` como alternativa a `challenge_id`.
- `catalog_search`/`learning_path_recommend` ganham `theme_ids`/
  `competency_ids` (conjunto) e campo `coverage` na resposta.

### Data/storage changes
Nenhuma no formato persistido de packs; a trilha ad hoc de R4 não é
autorada, é calculada em memória a cada chamada.

### Technical risks
- Composição de trilha pode produzir sequências longas/pedagogicamente
  ruins para muitos assuntos simultâneos — ver gatilho de revisão do ADR.

## 4. Tasks

### Planning
- [ ] Confirmar que nenhum consumidor atual depende do formato singular de
      `Query.Theme`/`Objective.ThemeID` fora de `internal/mcpserver`.

### Implementation
- [ ] Trocar `Query.Theme`/`Query.Competency` por `ThemeIDs`/`CompetencyIDs`.
- [ ] Trocar `Objective.ThemeID`/`CompetencyID` por versões em conjunto.
- [ ] Implementar cálculo de cobertura (`total`/`partial`/`none`).
- [ ] Implementar composição de trilha ad hoc a partir do grafo de
      precedência quando a cobertura for `partial`.
- [ ] Implementar lifecycle de sessão sobre `track_id` pré-existente.
- [ ] Atualizar contrato MCP (`session_start`, `catalog_search`,
      `learning_path_recommend`) e seus golden tests.

### Validation
- [ ] Teste de cobertura `total` (um item cobre todos os assuntos pedidos).
- [ ] Teste de cobertura `partial` com composição de trilha correta pelo
      grafo de precedência.
- [ ] Teste de cobertura `none`.
- [ ] Teste de sessão completa sobre `track_id` percorrendo todos os
      desafios em ordem.
- [ ] Suíte completa de internal/curriculum, internal/recommendation,
      internal/session, internal/mcpserver sem regressão.

## 5. Decisions

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

### Deterministic checks
- Test: go test ./internal/curriculum/... ./internal/recommendation/... ./internal/session/... ./internal/mcpserver/...
- Lint: gofmt -l internal/curriculum internal/recommendation internal/session internal/mcpserver
- Typecheck: go vet ./internal/curriculum/... ./internal/recommendation/... ./internal/session/... ./internal/mcpserver/...
- Build: go build ./...

### Execution log
- Pendente.

### Results summary
Nenhuma implementação ainda; spec criada para sequenciar o trabalho.

### Requirement trace
- Mapear R1–R5 a testes de cobertura e composição.

### Known gaps
- Nenhum até a implementação começar.

## 7. Final Report

### Delivered scope
Nenhum; spec draft aguardando implementação.

### Files and modules changed
- Planejados nas áreas afetadas acima.

### Validation executed
- Command: pose lint-spec learning-track-composition --ready-check
- Result: registrar após validação.

### Residual risks
- Nenhum adicional além do já descrito em Technical risks.

### Follow-ups
- [open]
