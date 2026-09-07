---
slug: concept-content-authoring
status: draft
created_at: 2026-09-07
completed_at:
supersedes:
depends_on: assistance-hints-detours, catalog-authoring-quality
priority: 40
components: curriculum, assistance, mcp-server
delivers: capability:canonical-concept-content
---

# Spec: concept-content-authoring

## 1. Intent

### Goal
Entregar explicações canônicas autoradas para conceitos, com exemplos fora da solução e relações úteis.

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
- R1: Definir campos versionados para explicação, exemplo, analogia opcional e referências a relações; manter documentos antigos carregáveis com indicação explícita de conteúdo ausente.
- R2: Retornar conteúdo autorado e relações por concept_content_get com limites, ordem determinística e sem gerar texto no servidor.
- R3: Impedir acesso a solução, fixtures reservadas ou pistas superiores por esse caminho; exemplos devem ter contexto distinto do desafio ativo.
- R4: Autorar e revisar uma amostra de linguagem, erros e testes; delegar a cobertura integral às specs de packs, evitando nova contagem duplicada.
- R5: Cobrir paridade schema/loader, campos ausentes, referência inválida, tamanho excessivo e consulta MCP sem mutação de sessão.

### Non-functional
- Manter determinismo e falhas explícitas; provar o caminho composto, não apenas helpers isolados.
- Identificar resultados por versão, commit e cenário; skips obrigatórios não são sucesso.

### Security
- Não copiar estado real do aluno para fixtures, logs ou relatórios.
- Aplicar confinamento de paths, redaction e consentimento nos novos caminhos.

### Compatibility
Registrar ADR da extensão aditiva de ConceptAuthoring antes de codificar. Não introduzir LLM no servidor; consumir o contrato de relações existente.

## 3. Technical Plan

### Affected areas
- curriculum
- assistance
- mcp-server

### Artifacts
- modified: internal/curriculum/model.go
- modified: internal/curriculum/validator.go
- modified: internal/assistance/service.go
- modified: internal/assistance/service_test.go
- modified: internal/mcpserver/assistance_tools.go
- modified: internal/mcpserver/assistance_contract_test.go
- modified: schemas/catalog.schema.json
- modified: docs/content-authoring.md
- created: testdata/catalog-quality/concept-content.yaml
- modified: .pose/indexes/validation-matrix.json

### Delivery targets
- capability:canonical-concept-content module:cmd/codinho profile:composed-capability entrypoint:cmd/codinho/main.go

### API/contract changes
Registrar ADR da extensão aditiva de ConceptAuthoring antes de codificar. Não introduzir LLM no servidor; consumir o contrato de relações existente.

### Data/storage changes
Documentar os campos e migração exigidos pelos requisitos antes do primeiro incremento. Nenhuma migração é executada nesta rodada de planejamento.

### Technical risks
Número de conceitos não prova qualidade explicativa; exemplos mal escolhidos podem entregar a solução. A amostra não substitui a curadoria do catálogo inteiro.

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
- Context: ConceptAuthoring e assistance.ConceptContent possuem apenas ID/Title; o follow-up aberto de assistance-hints-detours continua sem cobertura, apesar de catalog-authoring-quality estar done. PROJECT.md RF-021 e §15.7 pedem conteúdo canônico.
- Options considered: deixar implementação implícita no aceite; criar remediação focada.
- Decision: planejar remediação própria e exigir sua entrega antes do aceite.
- Rationale: v1-integrated-acceptance tem non-goal de adicionar features ou correções.
- Consequences: esta spec permanece draft; decisões estruturais novas exigem ADR no início da implementação.

## 6. Validation

### Strategy
Validar cada contrato com cenário positivo e negativo pela entrada de produção; usar somente dados sintéticos.

### Deterministic checks
- Test: go test ./internal/curriculum/... ./internal/assistance/... ./internal/mcpserver/...
- Lint: pose check --strict; pose lint-spec concept-content-authoring --ready-check
- Typecheck: go vet ./...
- Build: go build ./cmd/codinho
- Security / Contract: pose assess integrate; pose validate --strict --json .pose/results/delivery-validation.json; pose surface-check --spec concept-content-authoring --strict

### Execution log
- 2026-09-07 UTC: criada em revisão de planejamento; implementação e gates de entrega não executados.

### Results summary
Escopo proposto com requisitos verificáveis. A validação atual do produto está no relatório; não prova os novos requisitos.

### Requirement trace
Preencher R1–R5 com evidência por cenário durante a implementação e no closeout.

### Known gaps
Número de conceitos não prova qualidade explicativa; exemplos mal escolhidos podem entregar a solução. A amostra não substitui a curadoria do catálogo inteiro.

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
