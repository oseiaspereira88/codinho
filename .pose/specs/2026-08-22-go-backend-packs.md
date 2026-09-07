---
slug: go-backend-packs
status: draft
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
- [ ] Criar matriz de conceitos, competências e distribuição 23.
- [ ] Definir fixtures herméticas e timeouts conservadores.

### Implementation
- [ ] Autorar relações, desafios, passos e pistas.
- [ ] Criar fixtures HTTP, concorrência e SQL.
- [ ] Criar três trilhas.
- [ ] Executar checks, race detector e playtests.
- [ ] Revisar abstrações e linguagem idiomática.

### Validation
- [ ] Executar catálogo e todas as fixtures repetidamente.
- [ ] Executar race detector e teste de leak.
- [ ] Auditar distribuição e cobertura cumulativa.

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
- Pendente.

### Results summary
- Nenhum pack backend publicado.

### Requirement trace
- Mapear R1–R10 a manifests, coverage, playtests e check results.

### Known gaps
- Drivers e versões devem ser fixados no início da implementação.

## 7. Final Report

### Delivered scope
Nenhum; spec draft.

### Files and modules changed
- Planejados nos três packs, testdata e docs.

### Validation executed
- Command: pose lint-spec go-backend-packs --ready-check
- Result: registrar após gate.

### Residual risks
- Flakiness será bloqueante para publicação.

### Follow-ups
- [covered: v1-integrated-acceptance] Executar missão backend em sessão de ponta a ponta.
