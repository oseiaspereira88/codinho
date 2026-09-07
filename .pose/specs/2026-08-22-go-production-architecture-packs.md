---
slug: go-production-architecture-packs
status: draft
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: curriculum-graph-path-recommendation, catalog-authoring-quality, go-backend-packs
priority: 200
components: curriculum-content
delivers:
---

# Spec: go-production-architecture-packs

## 1. Intent

### Goal
Publicar packs de produção, arquitetura e depuração que completem a taxonomia e o volume curricular da V1.

### Business value
Praticar decisões seniores, investigação, operabilidade e evolução de sistemas em vez de apenas aumentar o número de arquivos.

### Constraints
- Publicar exatamente 13 desafios: 2 atômicos, 2 combinados, 2 fatias funcionais, 4 depuração/refatoração/revisão e 3 missões.
- Adicionar pelo menos 20 conceitos, 15 competências e 80 step nodes, atingindo 160 conceitos, 100 competências e 500 nodes cumulativos.
- Missões devem conter requisito incompleto, falhas e trade-offs.

### Non-goals
- Cobertura exaustiva de todo broker, banco, cloud ou framework.
- Ensinar padrões por nome sem força de design.

## 2. Requirements

### Functional
- R1: Publicar go-production e go-architecture-debugging.
- R2: Cobrir mensageria e gRPC conceituais, observabilidade, performance, runtime, segurança, arquitetura e depuração.
- R3: Entregar 2 atômicos sobre profiling e correlação observável.
- R4: Entregar 2 combinados sobre shutdown observável e cliente idempotente.
- R5: Entregar 2 fatias sobre worker operável e serviço instrumentado.
- R6: Entregar 4 desafios de depuração, refatoração ou revisão sobre race, retenção, overengineering e incidente.
- R7: Entregar 3 missões sobre backpressure, migração compatível e resiliência operacional.
- R8: Publicar duas trilhas: depuração e performance; arquitetura e produção.
- R9: Exigir justificativa de custo, falha, operabilidade e alternativa em toda missão.
- R10: Comprovar as metas cumulativas de 160 conceitos, 100 competências e 500 step nodes.

### Non-functional
- Desafios seniores podem ter pouco código e devem avaliar decisão.
- Checks de performance usam baseline local e tolerâncias documentadas.

### Security
- Incluir input limits, secrets, TLS, auth, supply chain e disponibilidade em casos pertinentes.
- Não usar serviços reais ou credenciais.

### Compatibility
- Contratos e migrações devem incluir backward compatibility e rollout progressivo.

## 3. Technical Plan

### Affected areas
- packs/go-production.yaml, packs/go-architecture-debugging.yaml

### Artifacts
- modified: packs/manifest.yaml
- created: packs/go-production.yaml
- created: packs/go-architecture-debugging.yaml
- created: testdata/packs/production/
- created: docs/catalog/go-production-architecture.md

### Delivery targets
Nenhum tipado; conteúdo do catálogo.

### API/contract changes
- Consumir schemas e checks existentes; extensão exige spec própria.

### Data/storage changes
- Fixtures sintéticas de logs, traces, perfis, incidentes e migrações.

### Technical risks
- Missões podem ficar subjetivas demais para avaliar.
- Benchmarks podem variar por hardware.

## 4. Tasks

### Planning
- [ ] Definir matriz dos 13 desafios e critérios semânticos.
- [ ] Definir fixtures portáveis e tolerâncias de benchmark.

### Implementation
- [ ] Autorar conceitos, competências, relações e desafios.
- [ ] Autorar rubricas de trade-off e investigação.
- [ ] Criar fixtures e checks.
- [ ] Criar duas trilhas.
- [ ] Realizar playtests com revisão sênior.

### Validation
- [ ] Executar catálogo, fixtures, race e benchmarks.
- [ ] Auditar distribuição e metas cumulativas.
- [ ] Revisar segurança, observabilidade e compatibilidade.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Senioridade não é proporcional ao tamanho do projeto.
- Options considered: projetos maiores; algoritmos difíceis; decisões sob restrições.
- Decision: usar investigações e missões com falhas e trade-offs.
- Rationale: aproxima competências reais de engenharia.
- Consequences: rubricas e revisão humana têm papel maior.

## 6. Validation

### Strategy
Combinar checks determinísticos, rubricas, fixtures de incidente e playtest independente.

### Deterministic checks
- Test: checks de todas as fixtures, race detector e benchmarks controlados.
- Lint: go run ./cmd/codinho catalog validate --json (executar na raiz; catálogo resolvido pelo manifest).
- Typecheck: go vet ./...
- Build: go build ./cmd/codinho
- Security / Contract: threat cases, compatibility fixtures e secret scan.

### Execution log
- Pendente.

### Results summary
- Nenhum pack de produção publicado.

### Requirement trace
- Mapear R1–R10 a manifests, coverage, rubricas, playtests e checks.

### Known gaps
- Tolerâncias de performance serão específicas por desafio, não globais.

## 7. Final Report

### Delivered scope
Nenhum; conteúdo planejado.

### Files and modules changed
- Planejados em dois packs, testdata e docs.

### Validation executed
- Command: pose lint-spec go-production-architecture-packs --ready-check
- Result: registrar após gate.

### Residual risks
- Critérios semânticos precisam de calibração interavaliador.

### Follow-ups
- [covered: v1-integrated-acceptance] Auditar metas cumulativas e uma missão completa.
