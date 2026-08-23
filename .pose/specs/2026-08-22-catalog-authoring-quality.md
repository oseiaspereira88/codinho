---
slug: catalog-authoring-quality
status: draft
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
- modified: internal/cli/catalog.go
- created: schemas/editorial.schema.json
- created: docs/content-authoring.md
- created: docs/content-review-checklist.md
- created: testdata/catalog-quality/

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
- [ ] Definir IDs, severidades e política de supressão.
- [ ] Definir prova de playtest e reviewer distinto.

### Implementation
- [ ] Implementar regras estruturais e editoriais.
- [ ] Implementar execução controlada de fixtures.
- [ ] Implementar projeção de cobertura e gates V1.
- [ ] Integrar CLI e CI.
- [ ] Documentar autoria e revisão.

### Validation
- [ ] Criar uma fixture negativa por regra.
- [ ] Executar mutation tests do validador onde viável.
- [ ] Revisar falsos positivos em amostra humana.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Parte da qualidade é objetiva e parte exige julgamento.
- Options considered: gate só humano; heurística bloqueante total; camadas determinística e editorial.
- Decision: bloquear regras objetivas e exigir evidência humana para qualidade semântica.
- Rationale: automação não deve fingir compreender pedagogia.
- Consequences: publicação possui workflow de revisão explícito.

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
- Pendente.

### Results summary
- Nenhum gate editorial entregue.

### Requirement trace
- Mapear R1–R10 a rule IDs, fixtures e relatório de cobertura.

### Known gaps
- Rubricas podem evoluir após playtests da V1.

## 7. Final Report

### Delivered scope
Nenhum; spec draft.

### Files and modules changed
- Planejados em curriculum, CLI, schemas, docs e testdata.

### Validation executed
- Command: pose lint-spec catalog-authoring-quality --ready-check
- Result: registrar após gate.

### Residual risks
- Revisão humana continua sendo o gargalo correto para publicação.

### Follow-ups
- [covered: v1-integrated-acceptance] Reconciliar contagem, qualidade e playtests de todos os packs.
