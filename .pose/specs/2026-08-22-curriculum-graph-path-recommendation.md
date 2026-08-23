---
slug: curriculum-graph-path-recommendation
status: in-progress
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: catalog-schema-loader, mastery-review-scheduling
priority: 120
components: curriculum, recommendations
delivers:
---

# Spec: curriculum-graph-path-recommendation

## 1. Intent

### Goal
Implementar o grafo curricular, pesquisa e recomendação explicável de trilhas e remediações.

### Business value
Escolher o que praticar com base em objetivo, dependências e evidências, em vez de uma lista linear fixa.

### Constraints
- Sugestões são consultivas e nunca iniciam trilha sem aceitação.
- Relações de precedência devem ser acíclicas.
- Recomendações não dependem de rede ou LLM.

### Non-goals
- Gerar desafios.
- Inferir conhecimento sem evidência.

## 2. Requirements

### Functional
- R1: Suportar requires, recommended_before, relates_to, contrasts_with, commonly_fails_with, applies_in, deepens_into e evidences.
- R2: Validar aciclicidade somente nas relações que implicam precedência.
- R3: Pesquisar itens por texto, tema, competência, dificuldade, tipo, duração e pré-requisito.
- R4: Recomendar caminho a partir de objetivo, tempo, perfil, progresso e revisões vencidas.
- R5: Explicar cada recomendação por evidência, gap, dependência e custo estimado.
- R6: Propor remediação menor e retorno ao desafio original.
- R7: Expor catalog_search, concept_relations_get e learning_path_recommend de forma completa.
- R8: Produzir resultado estável para mesma versão e estado.

### Non-functional
- Consultas p95 inferiores a 100 ms após indexação.
- Ranking deve ter tie-break determinístico.

### Security
- Query de texto é dado e deve ter limites de tamanho e custo.

### Compatibility
- Relações novas devem ser aditivas e desconhecidas falham de modo explícito.

## 3. Technical Plan

### Affected areas
- internal/curriculum/, internal/recommendation/, internal/mcpserver/

### Artifacts
- created: internal/curriculum/graph.go
- created: internal/curriculum/search.go
- created: internal/curriculum/graph_test.go
- created: internal/curriculum/search_test.go
- created: internal/curriculum/search_bench_test.go
- created: internal/curriculum/testdata/relation_cycle/manifest.yaml
- created: internal/curriculum/testdata/relation_cycle/pack.yaml
- created: internal/curriculum/testdata/unknown_relation_kind/manifest.yaml
- created: internal/curriculum/testdata/unknown_relation_kind/pack.yaml
- modified: internal/curriculum/model.go
- modified: internal/curriculum/validator.go
- modified: internal/curriculum/index.go
- modified: internal/curriculum/testdata/valid/pack.yaml
- created: internal/recommendation/model.go
- created: internal/recommendation/service.go
- created: internal/recommendation/explanation.go
- created: internal/recommendation/service_test.go
- created: internal/application/recommendation.go
- created: internal/application/recommendation_test.go
- modified: internal/application/catalog.go
- modified: internal/mastery/model.go
- modified: internal/mastery/rules_test.go
- modified: internal/mcpserver/catalog_tools.go
- created: internal/mcpserver/recommendation_tools.go
- created: internal/mcpserver/recommendation_contract_test.go
- modified: internal/mcpserver/server.go
- modified: internal/mcpserver/errors.go
- modified: internal/mcpserver/contract_test.go
- modified: cmd/codinho/main.go

### Delivery targets
Nenhum novo; amplia o contrato MCP V1 planejado.

### API/contract changes
- Completar filtros de busca e schemas de relações e recomendação.

### Data/storage changes
- Índices são projeções reconstruíveis do catálogo e do progresso.

### Technical risks
- Heurística pode favorecer pré-requisitos demais e atrasar prática.
- Relações editoriais inconsistentes degradam explicações.

## 4. Tasks

### Planning
- [x] Definir semântica e direção de cada relação.
- [x] Definir função de ranking, tie-break e limites.

### Implementation
- [x] Implementar grafo validado e vizinhança limitada.
- [x] Implementar índice de busca.
- [x] Implementar recomendador e explicações.
- [x] Implementar remediação e retorno.
- [x] Expor tools e fixtures.

### Validation
- [x] Testar ciclos, relações desconhecidas e ranking estável.
- [x] Testar perfis sem evidência e objetivos conflitantes.
- [x] Benchmarkar catálogo no tamanho-alvo da V1.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Recomendação deve ser previsível e auditável.
- Options considered: LLM; score oculto; regras determinísticas.
- Decision: ranking determinístico com decomposição de fatores.
- Rationale: permite entender e corrigir recomendações.
- Consequences: personalização sofisticada fica para versões futuras.

### Decision 2
- Date: 2026-08-23
- Context: `ChallengeAuthoring.Prerequisites` (catalog-schema-loader) já
  existe e já é validado contra ciclo (`validator.go`,
  `DiagPrerequisiteCycle`), mas só entre desafios e só implica uma
  relação (precedência). R1 exige 8 tipos de relação entre QUALQUER
  item do catálogo (conceito, competência, desafio, tema).
- Options considered: (a) substituir `Prerequisites` pelo novo grafo
  genérico; (b) manter `Prerequisites` como está (mecanismo já selado
  de catalog-schema-loader) e adicionar uma seção `relations:` nova e
  aditiva no pack, cobrindo os 8 tipos, independente de
  `Prerequisites`.
- Decision: (b). `Prerequisites` continua validando só precedência
  entre desafios, como sempre validou; a seção `relations:` é aditiva e
  cobre todo o resto (incluindo `requires` entre quaisquer itens, que
  pode se sobrepor semanticamente a `Prerequisites` para desafios, mas
  sem conflito mecânico — são validações independentes).
- Rationale: evitar migrar/depreciar um mecanismo já selado e testado
  de uma spec fechada só para generalizar o tipo; a Compatibility desta
  spec já exige que relações novas sejam aditivas.
- Consequences: um desafio pode ter precedência dupla-declarada
  (`prerequisites:` E uma relação `requires:`) sem que o servidor
  reconcilie as duas — cabe à autoria manter consistência (mesma
  fronteira de responsabilidade de catalog-authoring-quality).

### Decision 3
- Date: 2026-08-23
- Context: R2 exige aciclicidade "somente nas relações que implicam
  precedência", mas PROJECT.md não diz quais das 8 são essas.
- Decision: `requires`, `recommended_before` e `deepens_into` implicam
  ordem (o alvo depende logicamente da origem) e são validadas contra
  ciclo; `relates_to`, `contrasts_with`, `commonly_fails_with`,
  `applies_in` e `evidences` são simétricas/associativas por natureza e
  não são validadas contra ciclo.
- Rationale: leitura literal do nome de cada relação — "X requires Y" e
  "X deepens into Y" descrevem uma direção de dependência real;
  "X relates to Y" ou "X contrasts with Y" não.
- Consequences: uma relação desconhecida (fora das 8) falha
  explicitamente no load do pack (Compatibility), nunca é ignorada
  silenciosamente.

### Decision 4
- Date: 2026-08-23
- Context: R6 pede "remediação menor e retorno ao desafio original",
  mas as 8 relações de R1 não incluem um tipo "remediation" — a
  mais próxima (`commonly_fails_with`) descreve correlação de falha,
  não necessariamente um desafio mais simples da mesma competência.
- Decision: remediação é uma ESTRATÉGIA do recomendador, não um dado de
  currículo novo: quando o gap de domínio (mastery-review-scheduling)
  para a competência-alvo está abaixo de `demonstrates_without_help`, o
  recomendador busca, entre desafios que compartilham a mesma
  competência primária, o de menor `estimated_minutes`/dificuldade mais
  baixa como remediação, e mantém o desafio original na trilha logo
  depois.
- Rationale: não inventa uma nona relação de currículo não pedida por
  R1; reusa dados já existentes (competências, dificuldade, duração) e
  a projeção de domínio já entregue por mastery-review-scheduling.
- Consequences: a qualidade da remediação depende de existirem desafios
  menores cobrindo a mesma competência no catálogo autorado — mesmo
  risco já registrado em Technical risks e no Known Gap desta spec.

### Decision 5
- Date: 2026-08-23
- Context: R4 pede recomendação a partir de "objetivo, tempo, perfil,
  progresso e revisões vencidas", mas não existe (nem esta spec cria)
  um registro global de "quais desafios o aluno já concluiu" — sessão
  (session-orchestration-disclosure) é efêmera e por sessão; domínio
  (mastery-review-scheduling) é por competência, não por desafio.
- Decision: `learning_path_recommend` recebe `completed` (lista de
  challenge_id) como parâmetro explícito do chamador (o tutor, que
  observou as sessões), em vez de esta spec introduzir um novo registro
  durável de conclusões.
- Rationale: mesmo princípio já usado em `mastery_evidence_record`
  (Decision 2) e em `check_run` (safe-check-executor Decision 2): o
  servidor não infere um dado que não persiste, o chamador o fornece.
  Criar um registro de conclusões novo pertence a uma spec própria, não
  a um efeito colateral desta.
- Consequences: sem `completed`, desafios com pré-requisito de outro
  desafio nunca aparecem na recomendação — comportamento correto, mas
  que exige o chamador manter essa lista entre chamadas.

### Decision 6
- Date: 2026-08-23
- Context: `ChallengeAuthoring.Prerequisites` (catalog-schema-loader)
  pode referenciar um `concept` OU um `challenge` (validator.go aceita
  ambos). `completed` (Decision 5) só faz sentido para desafios — não
  há um sinal de "conceito já entendido" nesta base de código.
- Decision: `RecommendationService` só usa como filtro de elegibilidade
  os itens de `Prerequisites` que resolvem para outro `challenge` no
  catálogo; um pré-requisito que é um `concept` nunca bloqueia
  elegibilidade aqui (mas continua aparecendo na explicação de
  `Dependency` como já demonstrado pela spec original de
  catalog-schema-loader).
- Rationale: sem isso, TODO desafio com qualquer pré-requisito de
  conceito (o caso mais comum, ver `packs/go-first-steps.yaml`) seria
  permanentemente inelegível, já que "conceito concluído" não é um
  fato que `completed` consegue expressar.
- Consequences: um desafio com pré-requisito conceitual aparece na
  recomendação mesmo que o aluno nunca tenha visto aquele conceito —
  aceitável para V1 porque a instrução do próprio passo já cobre o
  conceito antes de exigi-lo (fluxo normal de sessão), e é a mesma
  fronteira de responsabilidade documentada no Known Gap desta spec.

## 6. Validation

### Strategy
Usar grafos sintéticos, golden rankings e benchmarks no volume V1.

### Deterministic checks
- Test: go test ./internal/curriculum/... ./internal/recommendation/... ./internal/application/... ./internal/mcpserver/...
- Lint: gofmt -l internal/curriculum internal/recommendation internal/application internal/mastery internal/mcpserver cmd/codinho
- Typecheck: go vet ./...
- Build: go build ./...
- Security / Contract: go test -race; texto de busca acima do limite rejeitado; govulncheck.

### Execution log
- `gofmt -l ...` (todos os pacotes tocados) → saída vazia (2026-08-23).
- `go vet ./...` (linux e `GOOS=windows`) → sem findings (2026-08-23).
- `go build ./...` → ok (2026-08-23).
- `go test ./... -race` → `ok` em todos os pacotes com testes, sem data races (2026-08-23).
- `govulncheck ./...` → "No vulnerabilities found." (2026-08-23).
- `go test ./internal/curriculum/... -bench BenchmarkSearchAtV1Scale -benchtime=1000x` → ~70µs/op com 84 desafios sintéticos (tamanho-alvo V1), muito abaixo do orçamento p95 de 100ms (2026-08-23).
- Smoke test real via mcp.CommandTransport contra `codinho serve` e
  `packs/go-first-steps.yaml`: `catalog_search` com `difficulty` +
  `competency` + `max_minutes` retorna exatamente o desafio esperado;
  `catalog_search` por `text: "Filtrar"` também; `concept_relations_get`
  em `slice-declaration` retorna `incoming`/`outgoing` vazios sem erro
  (o pack real ainda não autora `relations:`); `learning_path_recommend`
  com `competency_id: slice-filter` retorna o desafio com explicação
  completa (evidence/gap/dependency/cost_minutes) e score mais alto
  quando `time_budget_minutes` casa com `estimated_minutes`; pré-requisito
  conceitual (`slice-declaration`) corretamente NÃO bloqueia elegibilidade
  (Decision 6) (2026-08-23).

### Results summary
- `internal/curriculum` ganha `graph.go` (8 tipos de relação, aciclicidade
  só para `requires`/`recommended_before`/`deepens_into` — Decision 3 —
  kind desconhecido falha explicitamente no load, Compatibility) e
  `search.go` (filtros por texto, tema, competência, dificuldade, tipo,
  duração e pré-requisito, resultado ordenado por ID e limitado a 100
  itens com sinalização de truncamento, nunca silenciosa).
  `Prerequisites` (catalog-schema-loader) permanece intocado e validado
  como sempre (Decision 2): a seção `relations:` é inteiramente aditiva.
- `internal/recommendation` (puro, sem I/O) implementa ranking
  determinístico decomposto em fatores (match de competência/tema,
  revisão vencida, domínio já atingido, orçamento de tempo) com
  tie-break por ID (requirement R8, testado com 5 repetições idênticas),
  exclusão por pré-requisito não satisfeito, e remediação (Decision 4):
  quando o domínio da competência-alvo está abaixo de
  `demonstrates_without_help`, antepõe o desafio menor da mesma
  competência antes do original.
- `internal/application.RecommendationService` traduz
  `curriculum.Catalog`/`ProgressService` para as entradas puras do
  recomendador; `completed` (desafios concluídos) é fornecido pelo
  chamador, nunca inferido (Decision 5); pré-requisito que é um
  `concept` nunca bloqueia elegibilidade, só um `challenge` bloqueia
  (Decision 6).
- `internal/mcpserver` expõe `catalog_search` completo (R3, R7),
  `concept_relations_get` (R1, R7) e `learning_path_recommend` (R4, R5,
  R6, R7) — nenhum inicia trilha sozinho (Constraint), apenas consultivo.

### Requirement trace
- R1 [satisfied] report:internal/curriculum/graph.go (8 RelationKind constants) test:TestLoadValidPackIndexesRelationsBothDirections
- R2 [satisfied] test:TestLoadDetectsRelationCycleOnlyForPrecedenceKinds test:TestValidateAllowsCyclesAmongNonPrecedenceRelations
- R3 [satisfied] test:TestSearchByTextMatchesTitleCaseInsensitively test:TestSearchByDifficultyAndCompetency test:TestSearchByPrerequisite test:TestSearchByMaxMinutesExcludesLongerChallenges
- R4 [satisfied] test:TestRecommendationServiceFoldsInMasteryAndReviewDue test:TestRecommendPrioritizesOverdueReview
- R5 [satisfied] test:TestExplanationReportsEvidenceGapDependencyAndCost
- R6 [satisfied] test:TestRecommendPrependsSmallerRemediationBelowThreshold test:TestRecommendSkipsRemediationOnceAutonomousWithoutHelp
- R7 [satisfied] report:internal/mcpserver/catalog_tools.go report:internal/mcpserver/recommendation_tools.go
- R8 [satisfied] test:TestRecommendIsStableAcrossRepeatedCalls test:TestSearchIsDeterministicallySortedByID

### Known gaps
- Qualidade depende da curadoria dos 160 conceitos e 100 competências;
  `packs/go-first-steps.yaml` ainda não autora nenhuma `relations:`
  (confirmado no smoke test — `concept_relations_get` retorna vazio para
  conteúdo real).
- Sem um registro global de desafios concluídos (Decision 5), a
  qualidade de `learning_path_recommend` depende do chamador manter
  `completed` corretamente entre chamadas.

## 7. Final Report

### Delivered scope
Grafo curricular com 8 tipos de relação (aciclicidade só onde implica
precedência), busca completa por texto/tema/competência/dificuldade/
tipo/duração/pré-requisito, e recomendação de trilha explicável e
determinística que combina catálogo, domínio (mastery-review-
scheduling) e revisões vencidas — incluindo remediação automática por
desafio menor quando o domínio da competência-alvo ainda é raso.
Expostos como `concept_relations_get`, `learning_path_recommend` e um
`catalog_search` completo. Nenhuma geração de desafio, nenhuma inferência
de domínio sem evidência (non-goals preservados).

### Files and modules changed
- `internal/curriculum/{graph,search}.go` + testes e fixtures de teste (criados), `model.go`, `validator.go`, `index.go`, `testdata/valid/pack.yaml` (modificados)
- `internal/recommendation/{model,service,explanation}.go` + `service_test.go` (criados)
- `internal/application/recommendation.go` + `recommendation_test.go` (criados), `catalog.go` (modificado)
- `internal/mastery/model.go` (`Higher`, modificado — decisão de reusar a ladder já existente em vez de duplicá-la em `internal/recommendation`)
- `internal/mcpserver/recommendation_tools.go` + `recommendation_contract_test.go` (criados), `catalog_tools.go`, `server.go`, `errors.go`, `contract_test.go` (modificados)
- `cmd/codinho/main.go` (injeta `RecommendationService`)

### Validation executed
- Command: go test ./... -race
- Result: `ok` em todos os pacotes com testes, sem data races.
- Command: govulncheck ./...
- Result: "No vulnerabilities found."
- Command: go test ./internal/curriculum/... -bench BenchmarkSearchAtV1Scale -benchtime=1000x
- Result: ~70µs/op no tamanho-alvo V1 (84 desafios), dentro do orçamento p95 de 100ms.
- Command: smoke test real via mcp.CommandTransport contra `codinho serve`, `packs/go-first-steps.yaml`
- Result: catalog_search (filtros novos), concept_relations_get e learning_path_recommend consistentes ponta a ponta contra o pack real.

### Residual risks
- Ranking inicial (pesos 100/50/30/-50/10/-20) deverá ser reavaliado após uso real — heurística conservadora assumida nesta entrega.
- Ver Known Gaps: sem registro global de conclusões, e o pack real ainda não usa `relations:`.

### Follow-ups
- [covered: catalog-authoring-quality] Validar cobertura e coerência das relações.
