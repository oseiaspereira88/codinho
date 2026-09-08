---
slug: session-tree-progression
status: in-progress
created_at: 2026-09-07
completed_at:
supersedes:
depends_on: feedback-evaluation-progression, learning-practice-debug-modes, session-recovery-version-pinning, evaluation-evidence-lineage
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
Usar cursor relativo, progresso por nó e children_mode: sequence (default) ou choice. step_advance aceita next_step_id para escolher somente uma opção oferecida. Agregação conclui a janela sem fabricar eventos dos filhos; esta spec não inicia trilhas de múltiplos desafios.

## 3. Technical Plan

### Affected areas
- sessions
- learning-domain
- mcp-server

### Artifacts
- modified: internal/session/service.go
- modified: internal/session/progression.go
- modified: internal/session/granularity.go
- modified: internal/session/recovery.go
- modified: internal/session/progression_test.go
- modified: internal/session/assessment_test.go
- modified: internal/session/service_test.go
- modified: internal/session/modes_test.go
- modified: internal/curriculum/model.go
- modified: internal/curriculum/validator.go
- modified: internal/application/session.go
- modified: internal/mcpserver/errors.go
- modified: internal/application/assessment.go
- modified: internal/mcpserver/assessment_tools.go
- modified: internal/mcpserver/session_tools.go
- modified: internal/mcpserver/session_contract_test.go
- modified: .pose/indexes/validation-matrix.json
- modified: docs/compatibility.md
- modified: docs/content-authoring.md
- created: internal/session/tree_progression_test.go
- created: internal/curriculum/navigation_test.go
- created: cmd/codinho/tree_progression_integration_test.go
- created: .pose/adr/2026-09-07-depth-aware-session-cursor-and-explicit-branch-selection.md
- created: .pose/knowledge/2026-09-07-decision-log-adr-session-cursor-review.md
- created: .pose/changelogs/unreleased/session-tree-progression.md

### Delivery targets
- capability:session-tree-progression module:cmd/codinho profile:composed-capability entrypoint:cmd/codinho/main.go

### API/contract changes
Usar cursor relativo, progresso por nó e children_mode: sequence (default) ou choice. step_advance aceita next_step_id para escolher somente uma opção oferecida. Agregação conclui a janela sem fabricar eventos dos filhos; esta spec não inicia trilhas de múltiplos desafios.

### Data/storage changes
Adicionar navigation_version: 1 aos novos eventos de navegação e início. Reconstituir mapa de progresso, cursor, cobertura agrupada e escolhas por replay. Eventos anteriores mantêm projeção histórica; campos novos de catálogo usam omitempty para preservar digests legados. Nenhuma reescrita de JSONL.

### Technical risks
Agrupar nós pode apagar evidência ou revelar filhos; replay e trilhas consomem este estado e exigem teste composto no aceite.

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
- [ ] Executar pose assess integrate e validação estruturada no candidato.
- [ ] Reconciliar artifacts, surface e revisão independente antes do closeout.

## 5. Decisions

### Decision 2
- Date: 2026-09-07
- Decision: aplicar o [ADR de cursor](../adr/2026-09-07-depth-aware-session-cursor-and-explicit-branch-selection.md); consumir knowledge:adr-session-cursor-review e knowledge:adr-durable-session-replay-review.
- Rationale: profundidade muda a janela instrucional, não a identidade do progresso.
- Consequences: irmãos sem choice são obrigatórios em ordem; escolhas são explícitas e persistidas. Revisão independente segue agent-batch-review por delegação nativa.

### Decision 1
- Date: 2026-09-07
- Context: internal/session/progression.go rootSteps retorna somente a primeira layer; granularity.go deriveWindow recomeça pelo primeiro ramo e pode selecionar IDs de challenge/layer que findStep não resolve. Start escolhe firstStep mesmo quando a profundidade pedida é micro.
- Options considered: deixar implementação implícita no aceite; criar remediação focada.
- Decision: planejar remediação própria e exigir sua entrega antes do aceite.
- Rationale: v1-integrated-acceptance tem non-goal de adicionar features ou correções.
- Consequences: decisões estruturais novas exigem ADR no início da implementação.

## 6. Validation

### Strategy
Risco alto: navegação e replay cruzam catálogo, sessão e MCP. Todos os checks abaixo são obrigatórios; usar somente dados sintéticos.

| Cenário | Comando | Evidência esperada |
|---|---|---|
| Cinco profundidades, layers e sequência/choice | go test ./internal/session -run TestTree | Cursor correto; nenhum ramo obrigatório omitido |
| Troca de granularidade, avaliação, hints e conclusão preservados | go test ./internal/session -run TestTree | Um nó ativo; sem nova conclusão/tentativa por reexibição |
| Escolha inválida, predecessor, stale revision, sessão fechada e detour | go test ./internal/session -run TestTree | Erro antes do append |
| Schema de alternativas e identidade legada | go test ./internal/curriculum ./internal/session | Enum inválido rejeitado; retry histórico estável |
| Duas layers, dois ramos, profundidades e restart reais | go test ./cmd/codinho -run TestTreeProgressionOverRealStdio | Contrato MCP, separação de efeitos, entrevista e eventos únicos |
| Regressão completa | go test -race ./...; go vet ./...; go build ./cmd/codinho | Sem races ou regressões |

Navegação não faz rede nem executa processos; indisponibilidade de append e conflitos de revisão não podem alterar o cursor. Fallback de profundidade retorna o nó mais profundo autorado; choice interrompe a descida até seleção explícita.

### Deterministic checks
- Test: go test -race ./internal/session/... ./internal/learning/... ./internal/mcpserver/...
- Lint: pose check --strict; pose lint-spec session-tree-progression --ready-check
- Typecheck: go vet ./...
- Build: go build ./cmd/codinho
- Security / Contract: pose assess integrate; pose validate --strict --json .pose/results/delivery-validation.json; pose surface-check --spec session-tree-progression --strict

### Execution log
- 2026-09-08 UTC: TestTreeStartResolvesDepth reproduziu quatro falhas no código anterior (challenge/layer/meso/micro iniciavam em macro). Após implementar o cursor, testes de sessão, aplicação e MCP passaram, incluindo preservação de avaliação e hints no replay.
- 2026-09-08 UTC: TestTreeProgressionOverRealStdio passou com cinco profundidades, duas layers, seleção de ambas as alternativas, restarts, rejeições sem append e ausência de duplicação de completion/attempt/mastery. Revisão independente agent:tree_review iniciada por delegação nativa, gpt-5.6-luna/high; não substitui gates humanos.
- 2026-09-07 UTC: criada em revisão de planejamento; implementação e gates de entrega não executados.

### Results summary
Implementação e suíte completa, incluindo race e E2E real, passaram. Atribuição Git, validação estruturada e revisão final condicionam o fechamento.

### Requirement trace
- R1 [satisfied] test:TestTreeStartResolvesDepth test:TestTreeProgressionOverRealStdio
- R2 [satisfied] test:TestTreeOrderedFrontierCrossesLayers test:TestTreeSequencesChoicesAndNegativeNavigation
- R3 [satisfied] test:TestTreePreservesProgressAcrossWindowsAndRestart test:TestTreeCoarseCompletionAndLifecycleGuards
- R4 [satisfied] test:TestTreeSequencesChoicesAndNegativeNavigation test:TestNavigationCatalogRejectsAmbiguousChoices
- R5 [satisfied] test:TestTreeCoarseCompletionAndLifecycleGuards test:TestTreeLegacyStartAndNavigationRetry test:TestTreeUnknownNavigationVersionFailsClosed
- R6 [satisfied] capability:session-tree-progression evidence:integration check:session-tree-progression test:TestTreeProgressionOverRealStdio

### Known gaps
Agrupar nós pode apagar evidência ou revelar filhos; replay e trilhas consomem este estado e exigem teste composto no aceite.

## 7. Final Report

### Delivered scope
Cursor relativo, sequência completa de layers, alternativas explícitas, progresso por nó e replay compatível implementados e testados. Gates finais de entrega ainda pendentes.

### Files and modules changed
- Sessão/replay, catálogo, aplicação/MCP, testes e documentação declarados em Artifacts.

### Validation executed
- go test ./..., go test -race ./..., E2E stdio real, strict check/spec lint/knowledge/skills passaram. assess integrate não detectou contratos MCP Go; integração é comprovada pelo E2E. assess tech-debt e recurrence-check retornaram zero achados.

### Residual risks
Validar o comportamento implementado em execução independente; não reutilizar resultado histórico como aprovação do código futuro.

### Follow-ups
- [covered: v1-integrated-acceptance] Reexecutar os requisitos desta spec no candidato composto da V1.
