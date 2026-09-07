---
slug: session-recovery-version-pinning
status: draft
created_at: 2026-09-07
completed_at:
supersedes:
depends_on: eventstore-idempotency-scope, session-orchestration-disclosure, workspace-observation-baselines, reliability-observability-compatibility
priority: 10
components: sessions, eventstore, workspace, mcp-server
delivers: capability:session-recovery
---

# Spec: session-recovery-version-pinning

## 1. Intent

### Goal
Restaurar sessões, conteúdo fixado e escopo de evidências após reinício real do servidor.

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
- R1: Após fechar e reabrir codinho serve sobre o mesmo estado, restaurar sessão, política completa, profundidade, nó, avaliações, pistas, detours e revisão confirmados.
- R2: Preservar IDs e idempotência após restart: retry retorna o resultado original; request novo cria sessão distinta; rejeitar reutilização incompatível sem anexar evento.
- R3: Persistir identidade e conteúdo recuperável do catálogo fixado; atualização ou remoção de pack não troca silenciosamente instruções, checks ou critérios de sessões existentes.
- R4: Restaurar baselines, roots autorizadas e vínculo sessão–evidência; evidência externa ou obsoleta continua bloqueada após restart.
- R5: Validar transições antes do append; falha de domínio não pode deixar evento confirmado que o replay não consegue aplicar. Testar crash entre append e aplicação.
- R6: Ler logs legados sem inventar políticas/conteúdo ausentes; retornar diagnóstico explícito de sessão irrecuperável e permitir novas sessões sem apagar o histórico.
- R7: Provar os cenários por dois processos MCP distintos, incluindo atualização de pack entre processos, corrupção/truncamento e retry após perda da resposta.

### Non-functional
- Manter determinismo e falhas explícitas; provar o caminho composto, não apenas helpers isolados.
- Identificar resultados por versão, commit e cenário; skips obrigatórios não são sucesso.

### Security
- Não copiar estado real do aluno para fixtures, logs ou relatórios.
- Aplicar confinamento de paths, redaction e consentimento nos novos caminhos.

### Compatibility
Documentar em ADR antes da implementação a extensão dos payloads, snapshot do catálogo e política para eventos legados incompletos. Preservar JSONL e interfaces MCP existentes; nenhum banco novo.

## 3. Technical Plan

### Affected areas
- sessions
- eventstore
- workspace
- mcp-server

### Artifacts
- modified: internal/session/service.go
- created: internal/session/recovery.go
- created: internal/session/recovery_test.go
- modified: internal/application/workspace.go
- modified: internal/application/session.go
- modified: internal/eventstore/event.go
- modified: internal/curriculum/index.go
- modified: cmd/codinho/main.go
- created: cmd/codinho/recovery_integration_test.go
- modified: schemas/event.schema.json
- modified: docs/compatibility.md
- modified: .pose/indexes/validation-matrix.json

### Delivery targets
- capability:session-recovery module:cmd/codinho profile:composed-capability entrypoint:cmd/codinho/main.go

### API/contract changes
Documentar em ADR antes da implementação a extensão dos payloads, snapshot do catálogo e política para eventos legados incompletos. Preservar JSONL e interfaces MCP existentes; nenhum banco novo.

### Data/storage changes
Documentar os campos e migração exigidos pelos requisitos antes do primeiro incremento. Nenhuma migração é executada nesta rodada de planejamento.

### Technical risks
Replay parcial pode perder consentimento, atribuir evidência incorretamente ou reutilizar IDs. A ausência de dados em eventos antigos não admite recuperação perfeita por inferência.

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
- Context: O probe MCP de 2026-09-07 UTC iniciou ses_1; após reinício limpo, session_get retornou SESSION_NOT_ACTIVE e um novo session_start retornou STATE_CONFLICT. internal/session.New cria mapas e contador vazios; Start persiste apenas challenge_id/mode e usa CatalogRef vazio. A recuperação de JSONL existente não recompõe os serviços.
- Options considered: deixar implementação implícita no aceite; criar remediação focada.
- Decision: planejar remediação própria e exigir sua entrega antes do aceite.
- Rationale: v1-integrated-acceptance tem non-goal de adicionar features ou correções.
- Consequences: esta spec permanece draft; decisões estruturais novas exigem ADR no início da implementação.

## 6. Validation

### Strategy
Validar cada contrato com cenário positivo e negativo pela entrada de produção; usar somente dados sintéticos.

### Deterministic checks
- Test: go test -race ./internal/session/... ./internal/application/... ./internal/eventstore/... ./cmd/codinho/...
- Lint: pose check --strict; pose lint-spec session-recovery-version-pinning --ready-check
- Typecheck: go vet ./...
- Build: go build ./cmd/codinho
- Security / Contract: pose assess integrate; pose validate --strict --json .pose/results/delivery-validation.json; pose surface-check --spec session-recovery-version-pinning --strict

### Execution log
- 2026-09-07 UTC: criada em revisão de planejamento; implementação e gates de entrega não executados.

### Results summary
Escopo proposto com requisitos verificáveis. A validação atual do produto está no relatório; não prova os novos requisitos.

### Requirement trace
Preencher R1–R7 com evidência por cenário durante a implementação e no closeout.

### Known gaps
Replay parcial pode perder consentimento, atribuir evidência incorretamente ou reutilizar IDs. A ausência de dados em eventos antigos não admite recuperação perfeita por inferência.

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
