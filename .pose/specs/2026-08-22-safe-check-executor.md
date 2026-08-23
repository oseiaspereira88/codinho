---
slug: safe-check-executor
status: draft
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: catalog-schema-loader, workspace-observation-baselines, feedback-evaluation-progression
priority: 100
components: checks, workspace, mcp-server
delivers:
---

# Spec: safe-check-executor

## 1. Intent

### Goal
Executar verificações pré-declaradas por check_id com isolamento de argumentos, recursos, ambiente e filesystem.

### Business value
Gerar evidência objetiva sem oferecer ao agente uma ferramenta de shell arbitrária.

### Constraints
- Executar diretamente por os/exec, nunca por shell.
- Resolver programa e argumentos somente do catálogo fixado.
- Rede é negada por padrão.

### Non-goals
- Executar código totalmente não confiável com garantia de sandbox de kernel.
- Permitir comandos ad hoc solicitados pelo modelo.

## 2. Requirements

### Functional
- R1: Resolver check_id exclusivamente na versão de desafio da sessão.
- R2: Suportar go_test, go_test_race, go_vet, gofmt_check, go_build, go_benchmark e verificadores internos.
- R3: Validar programa, argumentos, cwd real, ambiente permitido e ausência de traversal.
- R4: Aplicar timeout, limite de output, paralelismo e cancelamento da árvore de processos.
- R5: Negar rede salvo check explicitamente marcado e aprovado.
- R6: Capturar stdout e stderr truncados, aplicar redaction e criar evidência imutável.
- R7: Vincular resultado à revisão do workspace e marcar evidência obsoleta quando necessário.
- R8: Retornar pass, fail, error ou skipped sem confundir falha de infra com falha do aluno.
- R9: Expor check_run sem parâmetro de comando livre.
- R10: Alimentar o critério `structural` de `step_evaluate` (feedback-evaluation-progression) com o resultado real do check (pass/fail/error/skipped), substituindo a derivação por mera presença de evidência (Decision 2 de feedback-evaluation-progression) por um verdict determinístico baseado na execução.

### Non-functional
- Um check cancelado não pode deixar processo filho.
- Limites devem ser configuráveis dentro de tetos seguros.

### Security
- Allowlist fechada de executáveis e variáveis.
- Testar traversal, env injection, output flood e process leak.

### Compatibility
- Guardrails específicos de plataforma ficam atrás de adapter e reportam suporte real.

## 3. Technical Plan

### Affected areas
- internal/checks/, internal/workspace/, internal/mcpserver/

### Artifacts
- created: internal/checks/model.go
- created: internal/checks/registry.go
- created: internal/checks/executor.go
- created: internal/checks/limits.go
- created: internal/checks/process_unix.go
- created: internal/checks/process_windows.go
- created: internal/checks/executor_test.go
- created: internal/checks/security_test.go
- modified: internal/mcpserver/workspace_tools.go
- created: internal/mcpserver/check_contract_test.go

### Delivery targets
Nenhum novo; amplia o contrato MCP V1 planejado.

### API/contract changes
- Adicionar check_run com input contendo somente sessão, check_id e revisão esperada.

### Data/storage changes
- Persistir check_executed e sua evidência redigida.

### Technical risks
- os/exec não é sandbox completo.
- Encerramento de grupos varia entre plataformas.
- Checks mal autorados podem ser perigosos.

## 4. Tasks

### Planning
- [ ] Definir schema fechado de cada tipo de check.
- [ ] Definir matriz de limites e capacidades por plataforma.

### Implementation
- [ ] Implementar registry e validadores.
- [ ] Implementar executor, timeout, output cap e redaction.
- [ ] Implementar process-group cancellation.
- [ ] Implementar política de rede e ambiente.
- [ ] Expor check_run e eventos.
- [ ] Conectar resultado do check ao critério `structural` de `step_evaluate` (R10).

### Validation
- [ ] Executar testes adversariais em todos os inputs.
- [ ] Testar timeout, flood, filhos, sinal e revisão obsoleta.
- [ ] Executar revisão de segurança independente.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: O aluno executará código próprio, mas o modelo não deve obter shell.
- Options considered: run_command; container obrigatório; checks tipados.
- Decision: checks tipados e allowlisted com limites locais.
- Rationale: superfície mínima e viável para V1 local.
- Consequences: não equivale a isolamento de código hostil e isso deve ser documentado.

## 6. Validation

### Strategy
Usar binaries de fixture controlados para simular ataques e estados de processo.

### Deterministic checks
- Test: go test ./internal/checks/... ./internal/mcpserver/...
- Lint: gofmt -l internal/checks
- Typecheck: go vet ./internal/checks/...
- Build: go test ./internal/checks/... -run ^$
- Security / Contract: go test -race; negative suite; govulncheck; revisão independente.

### Execution log
- Pendente.

### Results summary
- Nenhum executor entregue.

### Requirement trace
- Mapear R1–R10 a checks de segurança, lifecycle e contrato.

### Known gaps
- Isolamento forte para código adversarial está fora da V1.

## 7. Final Report

### Delivered scope
Nenhum; planejamento.

### Files and modules changed
- Planejados em checks, workspace e mcpserver.

### Validation executed
- Command: pose lint-spec safe-check-executor --ready-check
- Result: registrar após gate.

### Residual risks
- O modelo de ameaça local deve permanecer explícito na documentação.

### Follow-ups
- [covered: security-privacy-hardening] Executar threat model e hardening final.
