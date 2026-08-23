---
slug: catalog-authoring-quality
status: in-progress
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: catalog-schema-loader, safe-check-executor, administrative-cli-fixtures, curriculum-graph-path-recommendation
priority: 170
components: curriculum, cli, editorial
delivers:
---

# Spec: catalog-authoring-quality

## 1. Intent

### Goal
Entregar o workflow, linter e gates editoriais que tornam conceitos, competências, desafios, passos, pistas, checks e trilhas publicáveis.

### Business value
Impedir que a meta de cobertura produza conteúdo superficial, não avaliável ou que revele soluções.

### Constraints
- Conteúdo publicado exige revisão humana e execução por pessoa diferente do autor.
- Contagens não substituem cobertura ou qualidade.
- Validação deve ser offline e determinística.

### Non-goals
- Escrever os packs da V1.
- Avaliar estilo pedagógico apenas por heurística.

## 2. Requirements

### Functional
- R1: Validar schema, IDs, referências, ciclos, versões e paths de fixtures.
- R2: Aplicar teste de intenção única a micropassos e sinalizar múltiplos verbos independentes.
- R3: Detectar ausência de competência, critério, evidência, borda, pista ou reflexão exigida.
- R4: Verificar ordem crescente das pistas e possível gabarito no briefing ou níveis inferiores.
- R5: Executar ou validar todos os checks contra fixtures reproduzíveis.
- R6: Projetar cobertura por tema, conceito, competência, tipo, nível, variante e step node, e verificar cobertura e coerência das relações do grafo curricular (curriculum-graph-path-recommendation): todo item publicável deve ter ao menos uma relação declarada além de `requires`/`prerequisites` puro, e relações `contrasts_with`/`commonly_fails_with` devem ser recíprocas ou justificadas.
- R7: Exigir metadados de autoria, revisão humana e playtest antes de status published.
- R8: Expor codinho catalog validate com saída humana, JSON e códigos estáveis.
- R9: Falhar o gate V1 abaixo de 160 conceitos, 100 competências, 84 desafios, 12 trilhas e 500 step nodes.
- R10: Verificar a distribuição exata de tipos definida nas specs de packs.

### Non-functional
- Mensagens devem apontar pack, arquivo, item, regra e correção possível.
- Validação completa deve caber no CI sem rede.

### Security
- Fixtures e checks passam pelas mesmas regras de path e execução segura.
- Scanner detecta secrets e instruções inseguras.

### Compatibility
- Regras editoriais possuem IDs e severidade configurável sem alterar resultados antigos silenciosamente.

## 3. Technical Plan

### Affected areas
- internal/curriculum/, internal/cli/, docs/, schemas/

### Artifacts
- created: internal/curriculum/editorial.go
- created: internal/curriculum/coverage.go
- created: internal/curriculum/playtest.go
- created: internal/curriculum/editorial_test.go
- created: internal/curriculum/coverage_test.go
- modified: internal/curriculum/model.go
- modified: internal/curriculum/validator.go
- modified: internal/curriculum/loader.go
- modified: internal/curriculum/index.go
- modified: internal/cli/catalog.go
- created: internal/cli/editorial.go
- created: internal/cli/editorial_test.go
- modified: internal/cli/cli_test.go
- modified: cmd/codinho/cli_integration_test.go
- modified: schemas/challenge.schema.json
- created: schemas/editorial.schema.json
- created: docs/content-authoring.md
- created: docs/content-review-checklist.md
- created: testdata/catalog-quality/editorial/manifest.yaml
- created: testdata/catalog-quality/editorial/pack.yaml
- created: testdata/catalog-quality/structural/manifest.yaml
- created: testdata/catalog-quality/structural/pack.yaml

### Delivery targets
Nenhum; gate interno de conteúdo.

### API/contract changes
- Estabilizar findings do validador e formato JSON de cobertura.

### Data/storage changes
- Adicionar metadados editoriais versionados aos packs.

### Technical risks
- Heurísticas podem gerar falsos positivos.
- Playtest humano pode virar campo marcado sem evidência.

## 4. Tasks

### Planning
- [x] Definir IDs, severidades e política de supressão.
- [x] Definir prova de playtest e reviewer distinto.

### Implementation
- [x] Implementar regras estruturais e editoriais.
- [x] Implementar execução controlada de fixtures.
- [x] Implementar projeção de cobertura e gates V1.
- [x] Integrar CLI e CI.
- [x] Documentar autoria e revisão.

### Validation
- [x] Criar uma fixture negativa por regra.
- [ ] Executar mutation tests do validador onde viável.
- [x] Revisar falsos positivos em amostra humana.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Parte da qualidade é objetiva e parte exige julgamento.
- Options considered: gate só humano; heurística bloqueante total; camadas determinística e editorial.
- Decision: bloquear regras objetivas e exigir evidência humana para qualidade semântica.
- Rationale: automação não deve fingir compreender pedagogia.
- Consequences: publicação possui workflow de revisão explícito.

### Decision 2
- Date: 2026-08-23
- Context: checar checks reais contra fixtures (R5) exige
  `internal/fixtures` e `internal/checks`, mas ambos já importam
  `internal/curriculum` (fixtures para `ChallengeAuthoring`, checks
  indiretamente via workspace) — colocar essa lógica dentro de
  `internal/curriculum` criaria um ciclo de import.
- Options considered: (a) inverter a dependência, injetando uma
  interface de "materializar+executar" em `internal/curriculum`; (b)
  implementar a execução real na camada `internal/cli`, que já pode
  importar curriculum, fixtures e checks juntos sem ciclo, e manter em
  `internal/curriculum` só os IDs de regra e severidades.
- Decision: (b). `internal/cli/editorial.go` implementa
  `runChecksAgainstFixtures`, ligado a `catalog validate --checks`;
  `internal/curriculum/editorial.go` só declara
  `RuleFixtureNotReproducible`/`RuleCheckNotResolvable`/
  `RuleCheckExecutionError` como vocabulário compartilhado.
- Rationale: evita inversão de dependência artificial só para manter um
  arquivo num pacote "errado"; a camada CLI já é o lugar natural de
  qualquer coisa que spawna subprocessos.
- Consequences: `RunEditorialChecks` (puro, em curriculum) nunca inclui
  os achados de R5 sozinho — um caller que só chama essa função (sem
  passar por `catalog validate --checks`) não executa nada, o que é
  intencional (R5 exige subprocessos reais, opt-in via flag).

### Decision 3
- Date: 2026-08-23
- Context: R9 (gate V1: 160 conceitos, 100 competências, 84 desafios, 12
  trilhas, 500 step nodes) e R10 (distribuição exata por spec de pack)
  não podem ser satisfeitos por este catálogo em progresso (hoje: 1
  tema, ~11 conceitos, 3 competências, 2 desafios).
- Options considered: (a) rodar o gate V1 por padrão em `catalog
  validate`, bloqueando todo build atual; (b) manter o gate V1 e a
  distribuição exata como checks explícitos, opt-in (`--v1-gate`), nunca
  parte do fluxo padrão.
- Decision: (b).
- Rationale: um gate de conclusão da V1 não deve quebrar o
  desenvolvimento incremental do próprio conteúdo que ele mede; ele
  pertence ao momento do aceite V1 (v1-integrated-acceptance), não a
  cada `catalog validate`.
- Consequences: `codinho catalog validate` (sem flags) nunca falha por
  causa dos limiares V1; `--v1-gate` só é usado deliberadamente perto do
  aceite.

## 6. Validation

### Strategy
Usar corpus positivo/negativo, checks reais de fixture e auditoria amostral.

### Deterministic checks
- Test: go test ./internal/curriculum/... ./internal/cli/...
- Lint: codinho catalog validate packs
- Typecheck: go vet ./internal/curriculum/... ./internal/cli/...
- Build: go build ./cmd/codinho
- Security / Contract: fixture confinement, secret scan e check registry.

### Execution log
- `go build ./...` → ok (2026-08-23).
- `gofmt -l .` → sem saída (2026-08-23).
- `go vet ./...` e `GOOS=windows go vet ./...` → sem diagnósticos (2026-08-23).
- `go test ./... -race` → `ok` em todos os pacotes, incluindo `internal/curriculum` (editorial.go/coverage.go/playtest.go novos, corpus negativo em `testdata/catalog-quality/`) e `internal/cli` (editorial.go novo, `--checks`/`--v1-gate`), sem data races (2026-08-23).
- `go run ./cmd/codinho catalog validate` no catálogo real (`packs/`) → exit 0, só avisos `relation_isolated` (esperado: nenhum item ainda declara relations além de requires/prerequisites) (2026-08-23).
- `govulncheck ./...` → "No vulnerabilities found." (2026-08-23).

### Results summary
Gate de qualidade editorial entregue: validação estrutural existente
(schema/IDs/referências/ciclos) ganhou checagem de versão
(`invalid_version`); novo pacote de regras editoriais (pistas
crescentes, sem gabarito prematuro, competência/aceite obrigatórios,
metadados de publicação com revisor distinto, coerência de relações);
projeção de cobertura e gate V1 explícito (opt-in); execução real de
checks contra fixtures materializadas (opt-in, `--checks`). `codinho
catalog validate` permanece a única superfície (R8), agora reportando
diagnósticos estruturais + achados editoriais + cobertura.

### Requirement trace
- R1 [satisfied] internal/curriculum/validator.go (DiagInvalidVersion) + já existente (schema/IDs/referências/ciclos/fixture paths).
- R2 [satisfied] checkCompoundIntent (aviso) + test:TestCheckCompoundIntentFlagsConjunction.
- R3 [satisfied] checkCoverageGaps (competência/aceite bloqueantes; hints/reflexão avisos) + testes correspondentes.
- R4 [satisfied] checkHintOrder + checkNoEmbeddedSolutionText + testes correspondentes.
- R5 [satisfied] internal/cli/editorial.go (runChecksAgainstFixtures, real via internal_ast) + test:TestRunChecksAgainstFixturesExecutesRealInternalASTCheck, test:TestCatalogValidateChecksFlagExecutesFixtureChecks.
- R6 [satisfied] coverage.go (ProjectCoverage) + checkRelationCoherence + testes correspondentes.
- R7 [satisfied] playtest.go (checkPublicationMetadata) + testes correspondentes.
- R8 [satisfied] internal/cli/catalog.go (catalog validate: humano/JSON, exit codes estáveis) + test:TestCatalogValidateFlagsMissingCompetencyAsBlocking.
- R9 [satisfied] coverage.go (CheckV1Gate, opt-in `--v1-gate`) + test:TestCheckV1GateFlagsEveryDimensionBelowThreshold.
- R10 [satisfied] coverage.go (CheckTypeDistribution, genérico — nenhuma spec de pack ainda declara números exatos) + test:TestCheckTypeDistributionFlagsMismatch.

### Known gaps
- Mutation testing do validador não foi executado (ferramenta de mutation testing para Go não está no toolchain atual); cobertura por corpus negativo (`testdata/catalog-quality/`) é o substituto usado.
- R10 (distribuição exata) tem o mecanismo pronto mas nenhuma spec de pack (go-foundations-packs etc.) ainda declarou uma distribuição concreta para checar contra.
- `--checks` só prova reprodutibilidade de infraestrutura (materializa + resolve + executa sem erro); não afirma que o outcome (`pass`/`fail`) é o esperado — isso é revisão humana (Decision 1).

## 7. Final Report

### Delivered scope
Gate de qualidade editorial: validação estrutural ampliada (versões),
regras editoriais (pistas, briefing, cobertura, publicação, relações),
projeção de cobertura, gate V1 opt-in, execução real de checks contra
fixtures, documentação de autoria e revisão.

### Files and modules changed
- internal/curriculum/{editorial,coverage,playtest}.go (novos) + testes.
- internal/curriculum/{model,validator,loader,index}.go: campo Publication, DiagInvalidVersion, LoadPacks/NewCatalogFromPacks exportados.
- internal/cli/{catalog,editorial}.go: catalog validate com achados editoriais/cobertura, flags --checks e --v1-gate.
- schemas/challenge.schema.json (fixture, publication, network) + schemas/editorial.schema.json (novo).
- docs/{content-authoring,content-review-checklist}.md.
- testdata/catalog-quality/{editorial,structural}/ (corpus negativo, um exemplo por regra).

### Validation executed
- Command: go test ./... -race
- Result: ok em todos os pacotes.
- Command: go run ./cmd/codinho catalog validate (catálogo real)
- Result: exit 0, só avisos esperados.
- Command: govulncheck ./...
- Result: No vulnerabilities found.

### Residual risks
- Revisão humana continua sendo o gargalo correto para publicação.
- Heurística de intenção composta (`compound_micro_instruction`) é conservadora — pode deixar passar casos reais (por isso é aviso, nunca bloqueante).

### Follow-ups
- [covered: v1-integrated-acceptance] Reconciliar contagem, qualidade e playtests de todos os packs.
