---
slug: go-backend-packs
status: in-progress
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: curriculum-graph-path-recommendation, catalog-authoring-quality, go-foundations-packs
priority: 190
components: curriculum-content
delivers:
---

# Spec: go-backend-packs

## 1. Intent

### Goal
Publicar packs de concorrência, contexto, HTTP e banco de dados com prática de lifecycle e integração.

### Business value
Levar a fluência fundamental até comportamentos backend reais sem saltar diretamente para projetos monolíticos.

### Constraints
- Publicar exatamente 23 desafios: 8 atômicos, 8 combinados, 4 fatias funcionais, 2 depurações e 1 missão.
- Adicionar pelo menos 40 conceitos, 25 competências e 120 step nodes, atingindo cumulativamente 140 conceitos, 85 competências e 420 nodes.
- Não introduzir abstração sem força concreta.

### Non-goals
- Observabilidade profunda, performance, segurança sistêmica ou arquitetura de produção.
- Framework HTTP externo como padrão.

## 2. Requirements

### Functional
- R1: Publicar go-concurrency-context, go-http e go-database.
- R2: Cobrir goroutines, ownership de channels, select, locks, races, leaks, backpressure, context e shutdown.
- R3: Cobrir servidor e cliente HTTP, JSON, limites, timeouts, middleware, lifecycle, transport e httptest.
- R4: Cobrir database/sql, pool, contexto, queries, nulls, transações, isolamento, migrations e testes.
- R5: Entregar 8 atômicos e 8 combinados com checks focalizados.
- R6: Entregar 4 fatias: agregador concorrente, handler testado, cliente resiliente e persistência transacional.
- R7: Entregar 2 depurações sobre goroutine leak e timeout HTTP.
- R8: Entregar 1 missão de serviço pequeno com shutdown e persistência.
- R9: Publicar três trilhas: Go backend, concorrência e lifecycle, HTTP e dados.
- R10: Validar race detector, cancelamento e cleanup em todo desafio pertinente.

### Non-functional
- Fixtures de integração devem ser herméticas e rápidas.
- Concorrência deve ensinar correctness antes de performance.

### Security
- HTTP e SQL incluem limites, validação e injeção como casos negativos.
- Nenhuma fixture depende de serviço remoto.

### Compatibility
- Usar net/http e database/sql compatíveis com a versão Go declarada.

## 3. Technical Plan

### Affected areas
- packs/go-concurrency-context.yaml, packs/go-http.yaml, packs/go-database.yaml

### Artifacts
- modified: packs/manifest.yaml
- created: packs/go-concurrency-context.yaml
- created: packs/go-http.yaml
- created: packs/go-database.yaml
- created: testdata/packs/backend/
- created: docs/catalog/go-backend.md
- created: docs/editorial/go-backend-inventory.md
- created: docs/editorial/go-backend-queue.json
- modified: .pose/assessments/packs.md
- created: .pose/assessments/docs-catalog.md
- modified: .pose/assessments/README.md
- modified: .pose/assessments/consolidated.md
- modified: .pose/state/components/packs.json
- created: .pose/state/components/docs-catalog.json
- modified: .pose/state/project-state.md

### Delivery targets
Nenhum tipado; conteúdo do catálogo.

### API/contract changes
- Consumir schema v1 e check types existentes.

### Data/storage changes
- Fixtures usam temporários e bancos locais descartáveis.

### Technical risks
- Testes concorrentes podem ser flakey.
- Banco local pode criar diferenças de plataforma.

## 4. Tasks

### Planning
- [x] Criar matriz de conceitos, competências e distribuição 23 em `docs/editorial/go-backend-inventory.md`.
- [x] Definir fila serializável e dependências de autoria em `docs/editorial/go-backend-queue.json`.
- [x] Definir fixtures herméticas, checks focados e timeouts conservadores por pack.
- [x] Recalibrar a matriz após os três packs autorados, sem usar variantes para mascarar déficit.

### Implementation
- [x] Materializar o primeiro pack `go-concurrency-context` com oito desafios, relações, passos e pistas; manter os demais packs pendentes.
- [x] Autorar relações, desafios, passos e pistas do pack HTTP; manter o pack database pendente.
- [x] Autorar relações, desafios, passos e pistas do pack database; manter `tracks: []` para a montagem global.
- [x] Declarar fixtures HTTP, concorrência e SQL herméticas nos packs, com checks focados e referência local.
- [x] Criar as três trilhas backend: concorrência/lifecycle, HTTP/serviços e `go-backend` integrada.
- [x] Executar checks, repetir o executor e rodar race detector no cenário de leak; playtests permanecem pendentes.
- [ ] Executar playtests ponta a ponta.
- [ ] Revisar abstrações e linguagem idiomática.

### Validation
- [x] Executar catálogo e todas as fixtures repetidamente.
- [x] Executar race detector e teste de leak.
- [x] Auditar distribuição e cobertura cumulativa.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Backend integrado pode esconder a competência que falhou.
- Options considered: um projeto grande; desafios isolados apenas; progressão atômico até missão.
- Decision: combinar todos os níveis com contagem fixa.
- Rationale: permite diagnóstico e transferência.
- Consequences: fixtures precisam ser reutilizáveis sem acoplar soluções.

## 6. Validation

### Strategy
Executar checks herméticos, repetição para flake, race detector e revisão humana.

### Deterministic checks
- Test: checks de todas as fixtures e go test -race onde aplicável.
- Lint: go run ./cmd/codinho catalog validate --json (executar na raiz; catálogo resolvido pelo manifest).
- Typecheck: go vet ./...
- Build: go build ./cmd/codinho
- Security / Contract: HTTP negative suite, SQL injection cases e path confinement.

### Execution log
- 2026-09-16: spec promovida de `draft` para `in-progress` após `pose lint-spec go-backend-packs --ready-check` passar.
- 2026-09-16: `pose assess discover --component packs` passou (0 LOC de código); `--component docs/catalog` passou (0 LOC); o path ainda inexistente `testdata/packs` foi rejeitado pelo discover e permanece uma lacuna planejada.
- 2026-09-16: `./codinho catalog validate --json` passou no baseline com 105 conceitos, 64 competências, 54 desafios, 366 nodes e cinco trilhas, todos draft.
- 2026-09-16: `go run ./cmd/codinho catalog validate --json` não chegou à validação por `disk quota exceeded` durante compilação temporária; a limitação foi preservada, e o binário local foi usado para a leitura do baseline.
- 2026-09-16: fila backend criada com cinco tarefas; iniciar pelo piloto `author-go-concurrency-context` em worktree isolada, sem aprovação ou integração automática do autor.
- 2026-09-16: a tentativa do piloto via `scripts/editorial.py run` terminou com `codex exec` em `exit 1` por limite de uso da CLI; a worktree não produziu patch. A causa foi registrada como limitação operacional, não como aprovação.
- 2026-09-16: fallback manual do coordenador criou `packs/go-concurrency-context.yaml` e adicionou o pack ao manifest; a revisão corrigiu uma fixture cujo baseline passava indevidamente por fechar `done` em uma goroutine auxiliar.
- 2026-09-16: `./codinho catalog validate --checks --json` passou após a correção: 55/55 checks verificados, sem diagnostics ou achados editoriais bloqueantes; cobertura global observada 119 conceitos, 73 competências, 62 desafios, 414 nodes e seis trilhas, tudo draft.
- 2026-09-16: fallback manual do coordenador criou `packs/go-http.yaml` e adicionou o pack ao manifest; a auditoria estática confirmou nove desafios canônicos, 14 conceitos, oito competências e 45 nodes, sem campos humanos preenchidos.
- 2026-09-16: `./codinho catalog validate --json` passou após a inclusão de HTTP, sem diagnostics ou achados editoriais; `./codinho catalog validate --checks --json` verificou 58/58 checks, incluindo seis cenários HTTP, com cobertura global de 133 conceitos, 81 competências, 71 desafios, 459 nodes e sete trilhas, tudo draft.
- 2026-09-16: fallback manual do coordenador criou `packs/go-database.yaml` e adicionou o pack ao manifest para validação; a auditoria estática confirmou seis desafios canônicos, 12 conceitos, oito competências e 42 nodes, sem campos humanos preenchidos e com `tracks: []`.
- 2026-09-16: `./codinho catalog validate --json` passou após a inclusão de database, sem diagnostics ou achados editoriais; `./codinho catalog validate --checks --json` verificou 60/60 checks, incluindo quatro cenários SQL locais, com cobertura global de 145 conceitos, 89 competências, 77 desafios, 501 nodes e sete trilhas, tudo draft.
- 2026-09-16: montagem integrada adicionou a trilha `go-backend` ao pack database com os 23 desafios backend em ordem de pré-requisitos e criou `docs/catalog/go-backend.md`.
- 2026-09-16: `testdata/packs/backend/README.md` registrou o ponto de extração para workspace; a fonte executável permanece embutida nos packs para evitar duplicação não governada.
- 2026-09-16: segunda execução consecutiva de `./codinho catalog validate --checks --json` verificou novamente 60/60 checks, sem diagnostics ou achados editoriais; cobertura integrada permaneceu em 145 conceitos, 89 competências, 77 desafios, 501 nodes e oito trilhas.

### Results summary
- Nenhum pack backend foi publicado. Os três packs backend foram materializados, integrados à trilha e documentados como drafts; a extração de fixtures para workspace e os gates humanos continuam pendentes. A contagem autorada acrescenta 23 desafios, 40 conceitos, 25 competências e 135 nodes, mantendo os gates humanos pendentes.

### Requirement trace
- R1: `packs/go-concurrency-context.yaml`, `packs/go-http.yaml` e `packs/go-database.yaml` estão no manifest; status permanece draft.
- R2: concorrência/contexto/lifecycle estão cobertos pelos oito desafios do primeiro pack.
- R3: HTTP/JSON/limites/timeouts/middleware/client estão cobertos pelos nove desafios HTTP.
- R4: database/sql, pool, contexto, queries, nulls, transações, migrations e cleanup estão cobertos pelos seis desafios SQL.
- R5–R8: distribuição integrada é 8 atomic, 8 combined, 4 functional_slice, 2 debug e 1 systems_mission; checks locais cobrem rollback, cancelamento, HTTP e leak.
- R9: `go-concurrency-lifecycle`, `go-http-services` e `go-backend` existem no catálogo.
- R10: `catalog validate --checks` verificou 60/60 cenários; o leak usa `go_test_race`; nenhum resultado humano foi inferido.

### Known gaps
- O pack SQL usa driver sintético local e Go mínimo 1.25 nas fixtures; uma integração com driver real permanece fora desta spec.
- `testdata/packs/backend/` contém somente o registro de extração; duplicar fixtures para workspace exige fluxo manual de preparação e não é gate do validador catalogal.
- A revisão pedagógica, o playtest e a publicação continuam fora da automação.
- O autor delegado não pôde concluir por limite de uso da CLI; o scaffold manual não substitui uma revisão independente nem o playtest.

### Knowledge consumed

- `knowledge:agent-batch-review-pattern` — revisão independente, limites de paths e retomada de sessões.
- `knowledge:go-foundations-io-batch` — baseline editorial, distinção entre drafts/variantes e pendências humanas.

## 7. Final Report

### Delivered scope
Parcial: os três packs backend estão em draft, com validação estrutural e checks determinísticos; montagem global, fixtures externas e gates humanos permanecem pendentes.

### Files and modules changed
- `packs/go-concurrency-context.yaml`, `packs/go-http.yaml`, `packs/go-database.yaml` e `packs/manifest.yaml`.
- `docs/catalog/go-backend.md`, `testdata/packs/backend/README.md`, inventário e fila editorial, além dos artefatos de assessment/state registrados no planejamento.

### Validation executed
- Command: `./codinho catalog validate --json`
- Result: passou sem diagnostics/editorial, cobertura global 133 conceitos, 81 competências, 71 desafios, 459 nodes e sete trilhas.
- Command: `./codinho catalog validate --checks --json`
- Result: 58/58 checks verificados; todos os itens permanecem em draft.
- Result adicional: após database, 60/60 checks verificados; todos os itens permanecem em draft.
- Result integrado: a trilha `go-backend` contém 23 IDs únicos em ordem válida; a segunda execução repetiu 60/60 checks.

### Residual risks
- Flakiness será bloqueante para publicação.

### Follow-ups
- [covered: v1-integrated-acceptance] Executar missão backend em sessão de ponta a ponta.
