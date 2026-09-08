---
slug: concept-content-authoring
status: in-progress
created_at: 2026-09-07
completed_at:
supersedes:
depends_on: assistance-hints-detours, catalog-authoring-quality
priority: 40
components: internal/curriculum, internal/assistance, internal/mcpserver, cmd/codinho, schemas, packs, docs
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
- modified: schemas/pack.schema.json
- created: schemas/concept-content.schema.json
- modified: internal/curriculum/index.go
- created: internal/curriculum/concept_content.go
- created: internal/curriculum/concept_content_test.go
- created: cmd/codinho/concept_content_integration_test.go
- modified: packs/go-first-steps.yaml
- modified: packs/go-errors.yaml
- modified: packs/go-core.yaml
- created: .pose/adr/2026-09-08-versioned-canonical-concept-content.md
- created: .pose/knowledge/2026-09-08-decision-log-adr-concept-content-review.md
- modified: docs/content-authoring.md
- created: testdata/catalog-quality/concept-content.yaml
- modified: .pose/indexes/validation-matrix.json

### Delivery targets
- capability:canonical-concept-content module:cmd/codinho profile:composed-capability entrypoint:cmd/codinho/main.go

### API/contract changes
Registrar ADR da extensão aditiva de ConceptAuthoring antes de codificar. Não introduzir LLM no servidor; consumir o contrato de relações existente.

### Data/storage changes
Adicionar ConceptAuthoring.content opcional (null equivale a ausente), version=1, explanation, example {context, code, explanation}, analogy opcional e relation_refs [{kind, concept_id}]. Limites Unicode: 8000/200/4000/2000/2000 caracteres respectivamente; até 8 relações. Exigir texto não branco. Referências devem corresponder a arestas de saída existentes para conceitos, sem duplicatas. Retornar content_status available/missing e projeção pública com relações ordenadas por kind/concept_id. Não persistir nem gerar conteúdo. Copiar o conteúdo no índice e nas consultas para preservar imutabilidade. Incrementar versões dos três packs alterados sem mudar desafios/publicação.

### Technical risks
Número de conceitos não prova qualidade explicativa; exemplos mal escolhidos podem entregar a solução. A amostra não substitui a curadoria do catálogo inteiro.

## 4. Tasks

### Planning
- [x] Reproduzir o achado e revisar contratos/ADRs aplicáveis.
- [x] Completar decisões de formato e plano de testes negativos antes de modificar código.
- [x] Reconciliar esta lista de artefatos com os arquivos efetivos; declarar arquivos adicionais antes de alterá-los.

### Implementation
- [x] Implementar primeiro o menor fluxo que fecha a lacuna.
- [x] Integrar entradas reais, persistência/compatibilidade e diagnósticos.
- [x] Atualizar documentação e checks declarativos junto com o contrato.

### Validation
- [x] Executar os cenários de cada R-ID, incluindo negativos.
- [x] Executar pose assess integrate e validação estruturada no candidato.
- [x] Reconciliar artifacts, surface e revisão independente antes do closeout.

## 5. Decisions

### Decision 2
- Date: 2026-09-08
- Decision: adotar o ADR versioned-canonical-concept-content e consumir knowledge:planning-audit-2026-09 e knowledge:adr-concept-content-review.
- Rationale: a tool hoje retorna somente id/title; conteúdo público deve ser autorado separadamente das respostas reservadas. Estender pack.schema.json, pois catalog.schema.json contém somente o manifest.
- Consequences: manter compatibilidade de conceitos antigos, validar limites antes de servir e não alterar sessão/publicação. A amostra não substitui os gates editoriais dos packs.


### Decision 1
- Date: 2026-09-07
- Context: ConceptAuthoring e assistance.ConceptContent possuem apenas ID/Title; o follow-up aberto de assistance-hints-detours continua sem cobertura, apesar de catalog-authoring-quality estar done. PROJECT.md RF-021 e §15.7 pedem conteúdo canônico.
- Options considered: deixar implementação implícita no aceite; criar remediação focada.
- Decision: planejar remediação própria e exigir sua entrega antes do aceite.
- Rationale: v1-integrated-acceptance tem non-goal de adicionar features ou correções.
- Consequences: esta spec permanece draft; decisões estruturais novas exigem ADR no início da implementação.

## 6. Validation

### Strategy
Risco médio: extensão aditiva de schema e resposta MCP, sem alteração de sessão.

| Cenário | Comando obrigatório | Evidência esperada |
|---|---|---|
| Schema e loader | go test ./internal/curriculum -run TestConceptContent -count=1 | Conteúdo válido/legado, null, versão desconhecida, obrigatórios, campos reservados, limites Unicode e referências inválidas. |
| Projeção e imutabilidade | go test ./internal/assistance -count=1 | Conteúdo autorado exato, relações ordenadas, cópias independentes, ausente explícito, ID desconhecido. |
| MCP e não divulgação | go test ./internal/mcpserver -run TestContractConcept -count=1 | Envelope none, contrato aditivo e nenhuma fixture/solução/pista superior. |
| Composição stdio | go test ./cmd/codinho -run TestConceptContentOverRealStdio -count=1 | Loader → aplicação → MCP real; sessão, revisão, hint ladder e events.jsonl idênticos antes/depois. |
| Amostra e gates | make check | Conceitos de linguagem, erros e testes disponíveis no catálogo autoral; suíte race, vet, build, scanner e provas editoriais preservadas. |

Exemplos usam contextos autorados diferentes dos desafios (meteorologia,
reserva de planetário e calendário), verificados editorialmente. O servidor
não pode provar diferença semântica: o modelo fechado impede navegar para
soluções/fixtures/pistas, e a revisão autoral impede copiá-las em prosa.
O tutor recebe dados públicos, nunca uma alegação de isolamento contra autor
malicioso. Nenhum acesso de rede, timeout externo ou migração é introduzido.

### Deterministic checks
- Test: go test ./internal/curriculum/... ./internal/assistance/... ./internal/mcpserver/...
- Lint: pose check --strict; pose lint-spec concept-content-authoring --ready-check
- Typecheck: go vet ./...
- Build: go build ./cmd/codinho
- Security / Contract: pose assess integrate; pose validate --strict --json .pose/results/delivery-validation.json; pose surface-check --spec concept-content-authoring --strict

### Execution log
- 2026-09-07 UTC: criada em revisão de planejamento; baseline registrado.
- 2026-09-08 UTC: implementada extensão de schema (schemas/concept-content.schema.json, pack.schema.json) e modelo Go (internal/curriculum), validação de limites Unicode e integridade de relações, projeção canônica com ordenação determinística e imutabilidade (internal/assistance), testes MCP de contrato e não-divulgação (internal/mcpserver), teste de composição sobre stdio real (cmd/codinho/concept_content_integration_test.go), documentação autoral (docs/content-authoring.md) e amostra revisada nos packs (go-first-steps, go-errors, go-core).

### Results summary
Todos os cenários obrigatórios e verificações determinísticas executadas com sucesso.

### Requirement trace
- R1 [satisfied] report:internal/curriculum/concept_content_test.go report:schemas/concept-content.schema.json report:schemas/pack.schema.json
- R2 [satisfied] report:internal/assistance/service_test.go report:internal/mcpserver/assistance_contract_test.go
- R3 [satisfied] report:internal/mcpserver/assistance_contract_test.go report:cmd/codinho/concept_content_integration_test.go
- R4 [satisfied] report:packs/go-first-steps.yaml report:packs/go-errors.yaml report:packs/go-core.yaml
- R5 [satisfied] report:internal/curriculum/concept_content_test.go report:cmd/codinho/concept_content_integration_test.go

### Known gaps
Número de conceitos não prova qualidade explicativa; exemplos mal escolhidos podem entregar a solução. A amostra não substitui a curadoria do catálogo inteiro.

## 7. Final Report

### Delivered scope
- Suporte a conteúdo canônico versionado para conceitos no catálogo, com version=1, explanation, example (context, code, explanation), analogy opcional e relation_refs com validação de limites e integridade de grafo.
- Projeção e ferramenta concept_content_get com content_status (available/missing) sem mutação de sessão, sem vazamento de fixtures ou soluções e sem geração de texto no servidor.
- Amostra autoral nos packs go-first-steps, go-errors e go-core usando contextos distintos dos desafios (meteorologia, reserva de planetário, calendário).
- Testes unitários, de contrato MCP e de ponta a ponta sobre stdio real.
- Documentação de autoria e atualização do validation matrix.

### Files and modules changed
- internal/curriculum/model.go
- internal/curriculum/validator.go
- internal/curriculum/index.go
- internal/curriculum/concept_content.go
- internal/curriculum/concept_content_test.go
- internal/assistance/service.go
- internal/assistance/service_test.go
- internal/mcpserver/assistance_tools.go
- internal/mcpserver/assistance_contract_test.go
- schemas/pack.schema.json
- schemas/concept-content.schema.json
- testdata/catalog-quality/concept-content.yaml
- cmd/codinho/concept_content_integration_test.go
- packs/go-first-steps.yaml
- packs/go-errors.yaml
- packs/go-core.yaml
- .pose/adr/2026-09-08-versioned-canonical-concept-content.md
- .pose/knowledge/2026-09-08-decision-log-adr-concept-content-review.md
- docs/content-authoring.md
- .pose/indexes/validation-matrix.json

### Validation executed
- go test ./... -race -count=1
- bash scripts/ci/check-format.sh
- go vet ./...
- go build ./cmd/codinho
- go test ./cmd/codinho -run TestConceptContentOverRealStdio -count=1
- pose assess integrate
- pose check --strict
- pose docs-check
- pose lint-spec concept-content-authoring --ready-check

### Residual risks
Validação integral de packs futuros depende de curadoria humana contínua dos exemplos contra vazamento semântico.

### Follow-ups
- [covered: v1-integrated-acceptance] Reexecutar os requisitos desta spec no candidato composto da V1.
