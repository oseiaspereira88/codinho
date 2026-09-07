---
slug: catalog-publication-integrity
status: draft
created_at: 2026-09-07
completed_at:
supersedes:
depends_on: catalog-authoring-quality, catalog-schema-loader, safe-check-executor
priority: 30
components: curriculum, cli, mcp-server
delivers: capability:catalog-publication-integrity
---

# Spec: catalog-publication-integrity

## 1. Intent

### Goal
Separar inventário autorado de conteúdo publicável e ligar os gates editoriais à CLI e à CI.

### Business value
Fechar uma lacuna verificável da V1 antes de ampliar a superfície que depende dela.

### Constraints
- Preservar autoria do aluno, operação local, JSONL e separação pedagógica.
- Usar o [relatório de auditoria](../reports/2026-09-07-doc-audit-auditoria-do-planejamento-v1.md) como baseline, não como evidência de entrega.
- Consumir knowledge:planning-audit-2026-09 antes de iniciar a implementação.

### Non-goals
- Reabrir ou reescrever atestações históricas de specs done.
- Publicar release, alterar código do aluno ou ampliar a V1 para nuvem/multiusuário.

## 2. Requirements

### Functional
- R1: Projetar separadamente inventário, rascunhos, publicações e conteúdo elegível aos gates; contar nos 84 apenas desafios curados, sem duplicar variantes ou protótipos.
- R2: Aplicar distribuição declarativa por pack e global na CLI, incluindo tipos inesperados; ligar CheckTypeDistribution ao caminho de produção e testar que uma contagem divergente falha.
- R3: Rejeitar publicação inválida antes do uso normal pelo MCP; busca/recomendação padrão só mostram conteúdo publicado, com modo explícito de autoria/playtest para rascunhos locais.
- R4: Executar checks com expectativa declarada de baseline e solução de referência; falha pedagógica prevista é distinta de timeout, skip, teste inexistente e erro de compilação inesperado. Gabaritos ficam fora da projeção pública.
- R5: Detectar check sem fixture ou evidência de execução alternativa e registrar justificativa por desafio; comprovar 100% dos checks publicáveis em CI ou por fixture.
- R6: Validar um corpus compartilhado de documentos válidos/inválidos contra schemas JSON e loader Go, incluindo publication, relations, fixture, variantes e campos desconhecidos.
- R7: Definir comportamento para metadados legados ausentes sem promovê-los a published; preservar IDs e sessões existentes, sem excluir trabalho autorado.
- R8: Provar pela CLI real que catálogo numericamente suficiente mas sem revisão/playtest falha o gate V1; manter validação estrutural incremental utilizável.

### Non-functional
- Manter determinismo e falhas explícitas; provar o caminho composto, não apenas helpers isolados.
- Identificar resultados por versão, commit e cenário; skips obrigatórios não são sucesso.

### Security
- Não copiar estado real do aluno para fixtures, logs ou relatórios.
- Aplicar confinamento de paths, redaction e consentimento nos novos caminhos.

### Compatibility
Registrar ADR para visibilidade, distribuição e expectativas de checks antes do código. Reusar publication e o gate humano existentes; geração dinâmica/quarentena fica em agent-authored-catalog-drafts.

## 3. Technical Plan

### Affected areas
- curriculum
- cli
- mcp-server

### Artifacts
- modified: internal/curriculum/coverage.go
- modified: internal/curriculum/coverage_test.go
- modified: internal/curriculum/index.go
- modified: internal/curriculum/loader.go
- modified: internal/curriculum/editorial.go
- modified: internal/cli/catalog.go
- modified: internal/cli/editorial.go
- modified: internal/cli/editorial_test.go
- modified: schemas/catalog.schema.json
- modified: schemas/challenge.schema.json
- created: testdata/catalog-quality/distribution.json
- created: internal/curriculum/schema_contract_test.go
- modified: docs/content-authoring.md
- modified: .pose/indexes/validation-matrix.json

### Delivery targets
- capability:catalog-publication-integrity module:cmd/codinho profile:composed-capability entrypoint:cmd/codinho/main.go

### API/contract changes
Registrar ADR para visibilidade, distribuição e expectativas de checks antes do código. Reusar publication e o gate humano existentes; geração dinâmica/quarentena fica em agent-authored-catalog-drafts.

### Data/storage changes
Documentar os campos e migração exigidos pelos requisitos antes do primeiro incremento. Nenhuma migração é executada nesta rodada de planejamento.

### Technical risks
O catálogo atual não tem published: aplicar filtro exige um fluxo explícito de playtest para continuar autoria. Schema JSON não deve virar fonte concorrente de regras.

## 4. Tasks

### Planning
- [ ] Reproduzir o achado e revisar contratos/ADRs aplicáveis.
- [ ] Completar decisões de formato e plano de testes negativos antes de modificar código.
- [ ] Reconciliar esta lista de artefatos com os arquivos efetivos; declarar arquivos adicionais antes de alterá-los.

### Implementation
- [ ] Implementar primeiro o menor fluxo que fecha a lacuna.
- [ ] Integrar entradas reais, persistência/compatibilidade e diagnósticos.
- [ ] Atualizar documentação e checks declarativos junto com o contrato.

### Validation
- [ ] Executar os cenários de cada R-ID, incluindo negativos.
- [ ] Executar pose assess integrate e validação estruturada no candidato.
- [ ] Reconciliar artifacts, surface e revisão independente antes do closeout.

## 5. Decisions

### Decision 1
- Date: 2026-09-07
- Context: ProjectCoverage conta todos os desafios, incluindo 41 itens sem publicação. CheckTypeDistribution é chamado apenas por testes. runChecksAgainstFixtures ignora desafios sem fixture e aceita fail/skipped como sucesso de infraestrutura; schemas JSON e loader Go não têm teste de equivalência.
- Options considered: deixar implementação implícita no aceite; criar remediação focada.
- Decision: planejar remediação própria e exigir sua entrega antes do aceite.
- Rationale: v1-integrated-acceptance tem non-goal de adicionar features ou correções.
- Consequences: esta spec permanece draft; decisões estruturais novas exigem ADR no início da implementação.

## 6. Validation

### Strategy
Validar cada contrato com cenário positivo e negativo pela entrada de produção; usar somente dados sintéticos.

### Deterministic checks
- Test: go test -race ./internal/curriculum/... ./internal/cli/... ./internal/mcpserver/...
- Lint: pose check --strict; pose lint-spec catalog-publication-integrity --ready-check
- Typecheck: go vet ./...
- Build: go build ./cmd/codinho
- Security / Contract: pose assess integrate; pose validate --strict --json .pose/results/delivery-validation.json; pose surface-check --spec catalog-publication-integrity --strict

### Execution log
- 2026-09-07 UTC: criada em revisão de planejamento; implementação e gates de entrega não executados.

### Results summary
Escopo proposto com requisitos verificáveis. A validação atual do produto está no relatório; não prova os novos requisitos.

### Requirement trace
Preencher R1–R8 com evidência por cenário durante a implementação e no closeout.

### Known gaps
O catálogo atual não tem published: aplicar filtro exige um fluxo explícito de playtest para continuar autoria. Schema JSON não deve virar fonte concorrente de regras.

## 7. Final Report

### Delivered scope
Somente planejamento; nenhuma funcionalidade desta spec foi entregue.

### Files and modules changed
- Esta spec; dependências no roadmap e no aceite integrado.

### Validation executed
- Planejamento sujeito a pose lint-spec --ready-check e pose check --strict nesta auditoria.

### Residual risks
Validar o comportamento implementado em execução independente; não reutilizar resultado histórico como aprovação do código futuro.

### Follow-ups
- [covered: v1-integrated-acceptance] Reexecutar os requisitos desta spec no candidato composto da V1.
