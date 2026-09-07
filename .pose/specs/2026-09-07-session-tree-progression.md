---
slug: session-tree-progression
status: draft
created_at: 2026-09-07
completed_at:
supersedes:
depends_on: feedback-evaluation-progression, learning-practice-debug-modes
priority: 20
components: sessions, learning-domain, mcp-server
delivers: capability:session-tree-progression
---

# Spec: session-tree-progression

## 1. Intent

### Goal
Percorrer a árvore completa do desafio e preservar a posição ao mudar granularidade.

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
- R1: Iniciar uma sessão na profundidade solicitada com nó e kind coerentes; instruction_get resolve challenge, layer, macro, meso e micro sem ITEM_NOT_FOUND artificial.
- R2: Concluir e avançar explicitamente por todas as layers, respeitando predecessores e escolhas; não declarar desafio esgotado na primeira layer.
- R3: Ajustar granularidade a partir da posição atual, preservando avaliações, filhos já concluídos e uma única instrução ativa; não reiniciar no primeiro ramo.
- R4: Distinguir sequência obrigatória de alternativas; uma ramificação devolvida por step_advance deve ter uma operação pública capaz de escolher o próximo nó permitido.
- R5: Preservar a separação feedback/avaliação/conclusão/avanço, overrides auditáveis e bloqueio de disclosure em entrevista durante toda navegação.
- R6: Validar duas layers, dois ramos e ao menos três profundidades por contrato MCP; rejeitar escolhas fora da árvore e provar ausência de conclusão/maestria duplicada.

### Non-functional
- Manter determinismo e falhas explícitas; provar o caminho composto, não apenas helpers isolados.
- Identificar resultados por versão, commit e cenário; skips obrigatórios não são sucesso.

### Security
- Não copiar estado real do aluno para fixtures, logs ou relatórios.
- Aplicar confinamento de paths, redaction e consentimento nos novos caminhos.

### Compatibility
Definir em ADR a semântica de cursor/progresso e escolha de ramo antes de alterar contratos. Reusar o lifecycle; esta spec não inicia trilhas de múltiplos desafios.

## 3. Technical Plan

### Affected areas
- sessions
- learning-domain
- mcp-server

### Artifacts
- modified: internal/session/service.go
- modified: internal/session/progression.go
- modified: internal/session/granularity.go
- modified: internal/session/progression_test.go
- modified: internal/session/disclosure_test.go
- modified: internal/learning/session.go
- modified: internal/application/assessment.go
- modified: internal/mcpserver/assessment_tools.go
- modified: internal/mcpserver/session_tools.go
- modified: internal/mcpserver/session_contract_test.go
- modified: .pose/indexes/validation-matrix.json

### Delivery targets
- capability:session-tree-progression module:cmd/codinho profile:composed-capability entrypoint:cmd/codinho/main.go

### API/contract changes
Definir em ADR a semântica de cursor/progresso e escolha de ramo antes de alterar contratos. Reusar o lifecycle; esta spec não inicia trilhas de múltiplos desafios.

### Data/storage changes
Documentar os campos e migração exigidos pelos requisitos antes do primeiro incremento. Nenhuma migração é executada nesta rodada de planejamento.

### Technical risks
Agrupar nós pode apagar evidência ou revelar filhos; replay e trilhas consomem este estado e exigem teste composto no aceite.

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
- Context: internal/session/progression.go rootSteps retorna somente a primeira layer; granularity.go deriveWindow recomeça pelo primeiro ramo e pode selecionar IDs de challenge/layer que findStep não resolve. Start escolhe firstStep mesmo quando a profundidade pedida é micro.
- Options considered: deixar implementação implícita no aceite; criar remediação focada.
- Decision: planejar remediação própria e exigir sua entrega antes do aceite.
- Rationale: v1-integrated-acceptance tem non-goal de adicionar features ou correções.
- Consequences: esta spec permanece draft; decisões estruturais novas exigem ADR no início da implementação.

## 6. Validation

### Strategy
Validar cada contrato com cenário positivo e negativo pela entrada de produção; usar somente dados sintéticos.

### Deterministic checks
- Test: go test -race ./internal/session/... ./internal/learning/... ./internal/mcpserver/...
- Lint: pose check --strict; pose lint-spec session-tree-progression --ready-check
- Typecheck: go vet ./...
- Build: go build ./cmd/codinho
- Security / Contract: pose assess integrate; pose validate --strict --json .pose/results/delivery-validation.json; pose surface-check --spec session-tree-progression --strict

### Execution log
- 2026-09-07 UTC: criada em revisão de planejamento; implementação e gates de entrega não executados.

### Results summary
Escopo proposto com requisitos verificáveis. A validação atual do produto está no relatório; não prova os novos requisitos.

### Requirement trace
Preencher R1–R6 com evidência por cenário durante a implementação e no closeout.

### Known gaps
Agrupar nós pode apagar evidência ou revelar filhos; replay e trilhas consomem este estado e exigem teste composto no aceite.

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
