---
slug: v1-integrated-acceptance
status: draft
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: architecture-decision-baseline, go-runtime-foundation, learning-domain-model, catalog-schema-loader, local-event-store, mcp-stdio-foundation, session-orchestration-disclosure, assistance-hints-detours, feedback-evaluation-progression, workspace-observation-baselines, safe-check-executor, mastery-review-scheduling, curriculum-graph-path-recommendation, tutor-skill-host-integration, administrative-cli-fixtures, learning-practice-debug-modes, interview-mode, catalog-authoring-quality, go-foundations-packs, go-backend-packs, go-production-architecture-packs, go-interviews-pack, security-privacy-hardening, reliability-observability-compatibility, installation-documentation-ci
priority: 250
components: integration, acceptance
delivers:
---

# Spec: v1-integrated-acceptance

## 1. Intent

### Goal
Comprovar por composição que a V1 atende integralmente o PROJECT.md e está pronta para um release governado.

### Business value
Evitar declarar conclusão com componentes isolados, conteúdo contado ou documentação sem uso real.

### Constraints
- Todas as specs predecessoras devem estar done e com evidência atribuída.
- Aceite exige Codex CLI e extensão IDE reais.
- Resultado parcial não pode ser chamado de V1 concluída.

### Non-goals
- Adicionar features, conteúdo ou refactors durante o gate.
- Publicar release automaticamente.

## 2. Requirements

### Functional
- R1: Executar sessão microguiada completa com iniciar, instruir, observar, avaliar, concluir e avançar.
- R2: Provar feedback sem progresso, pista limitada, desvio com retorno e solução revelada registrada.
- R3: Provar falha e correção de check, retomada após restart e sessão fixada após update de pack.
- R4: Executar modos ensino, prática, revisão, depuração, exploração e entrevista conforme policies.
- R5: Comprovar progressão de micro para maior autonomia e domínio sem promoção indevida.
- R6: Validar CLI, fixture preparation, export, doctor, MCP stdio e skill nos dois hosts.
- R7: Validar 160 conceitos, 100 competências, 84 desafios, 12 trilhas, 500 nodes e distribuição exata.
- R8: Executar suite adversarial de workspace, executor, disclosure, redaction e preservação.
- R9: Executar instalação limpa, recuperação, compatibilidade e documentação.
- R10: Reconciliar todos os critérios de aceite e Definition of Done da seção 30 e 37 de PROJECT.md.
- R11: Produzir relatório de piloto com sessões reais e limitações, sem dados pessoais.
- R12: Preparar bundle de revisão, attestation independente, surface checks, roadmap check e plano de release.
- R13: Comprovar em E2E a separação operacional entre feedback/avaliação (feedback-evaluation-progression) e execução de checks/observação de workspace (safe-check-executor, workspace-observation-baselines) — nenhuma avaliação deve depender de execução implícita de código do aluno, e nenhum check deve concluir ou avançar um passo por conta própria.

### Non-functional
- Todas as evidências devem ser atuais, reproduzíveis e ligadas a commit.
- Flakes, skips obrigatórios ou achados críticos bloqueiam o aceite.

### Security
- Reexecutar threat suite, scanners e revisão de privacidade no candidato exato.

### Compatibility
- Aceitar apenas plataformas e hosts com evidência no candidato.

## 3. Technical Plan

### Affected areas
- tests/e2e/, tests/security/, docs/acceptance/, .pose/reports/, todos os delivery targets V1

### Artifacts
- created: tests/e2e/session-guided/
- created: tests/e2e/session-resume/
- created: tests/e2e/modes/
- created: tests/e2e/host-cli/
- created: tests/e2e/host-ide/
- created: tests/e2e/installation/
- created: tests/security/v1-adversarial/
- created: docs/acceptance/v1-requirement-matrix.md
- created: docs/acceptance/v1-pilot-report.md
- created: docs/acceptance/v1-release-readiness.md

### Delivery targets
O alvo planejado é a capability composta `codinho-v1`, com módulo
`.agents/skills/codinho`, perfil `composed-capability` e entrypoint
`.agents/skills/codinho/SKILL.md`. Registrá-lo como alvo tipado quando esses
caminhos forem materializados, antes do closeout desta spec.

### API/contract changes
- Nenhuma; qualquer mudança retorna à spec proprietária e invalida a evidência afetada.

### Data/storage changes
- Somente fixtures e relatórios sanitizados.

### Technical risks
- Evidências de specs anteriores podem estar obsoletas.
- Piloto pequeno pode não demonstrar retenção de 30 dias.
- Host IDE pode exigir validação manual reproduzível.

## 4. Tasks

### Planning
- [ ] Congelar candidato, matriz de requisitos e conjunto de evidências.
- [ ] Confirmar prontidão de todas as dependências.

### Implementation
- [ ] Implementar apenas harnesses e relatórios de aceite.
- [ ] Executar os dez cenários E2E de PROJECT.md e variações de host.
- [ ] Executar piloto sanitizado e registrar limitações.
- [ ] Reconciliar conteúdo, segurança, compatibilidade e docs.
- [ ] Preparar review bundle e plano de release.

### Validation
- [ ] Executar pose validate --strict --report.
- [ ] Executar pose assess integrate, tech-debt e discover --update-state.
- [ ] Executar surface-check para todos os targets e roadmap-check codinho-v1.
- [ ] Exigir revisão e attestation independentes antes de closeout.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Contagem de componentes não prova produto composto.
- Options considered: fechar por spec; checklist manual; gate integrado com evidência.
- Decision: uma spec terminal depende de todas e valida a composição.
- Rationale: torna a Definition of Done verificável.
- Consequences: qualquer regressão relevante reabre a evidência final.

## 6. Validation

### Strategy
Executar matriz bidirecional requisito versus evidência em candidato imutável, com revisão separada.

### Deterministic checks
- Test: go test -race ./... e todos os harnesses tests/e2e.
- Lint: pose check --strict; catalog validate; skills-check; docs-check.
- Typecheck: go vet ./...
- Build: build e smoke install da matriz suportada.
- Security / Contract: govulncheck; secret scan; pose assess integrate; surface-check; roadmap-check.

### Execution log
- Pendente; dependências não entregues.

### Results summary
- Nenhum aceite executado e nenhuma prontidão declarada.

### Requirement trace
- Mapear R1–R13 e cada critério de PROJECT.md a check, test, report, surface e commit.

### Known gaps
- Retenção de 30 dias pode exigir janela real; se ausente, a V1 não deve alegar esse resultado.

## 7. Final Report

### Delivered scope
Nenhum; esta é a spec terminal de aceite em estado draft.

### Files and modules changed
- Planejados em harnesses e documentação de aceite.

### Validation executed
- Command: pose lint-spec v1-integrated-acceptance --ready-check
- Result: registrar após gate.

### Residual risks
- Nenhuma concessão de gate pode ser escondida; toda waiver deve aparecer na matriz.

### Follow-ups
- [wont-do: nenhum follow-up pós-V1 pode ser definido antes do aceite] Planejar versões futuras somente após o fechamento governado.
