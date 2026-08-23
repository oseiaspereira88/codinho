---
slug: safe-check-executor
status: in-progress
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
- R7: Vincular resultado à revisão do workspace e marcar evidência obsoleta quando necessário, usando `workspace.Fingerprint`/`workspace.FingerprintOf` (workspace-observation-baselines) em vez de reimplementar detecção de obsolescência — o mesmo mecanismo que `evidence_get` já usa para reportar `stale`.
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
- internal/checks/, internal/application/, internal/session/, internal/curriculum/, internal/workspace/, internal/mcpserver/

### Artifacts
- created: internal/checks/model.go
- created: internal/checks/registry.go
- created: internal/checks/executor.go
- created: internal/checks/limits.go
- created: internal/checks/process_unix.go
- created: internal/checks/process_unix_test.go
- created: internal/checks/process_windows.go
- created: internal/checks/process_windows_test.go
- created: internal/checks/executor_test.go
- created: internal/checks/security_test.go
- created: internal/application/checks.go
- created: internal/application/checks_test.go
- modified: internal/application/workspace.go
- modified: internal/application/session.go
- modified: internal/application/session_test.go
- modified: internal/workspace/redaction.go
- modified: internal/curriculum/model.go
- modified: internal/session/service.go
- modified: internal/session/service_test.go
- modified: internal/session/assessment_test.go
- modified: internal/session/assistance_test.go
- created: internal/session/checks_override_test.go
- created: internal/mcpserver/check_tools.go
- created: internal/mcpserver/check_contract_test.go
- modified: internal/mcpserver/server.go
- modified: internal/mcpserver/errors.go
- modified: internal/mcpserver/contract_test.go
- modified: internal/mcpserver/assessment_contract_test.go
- modified: internal/mcpserver/assistance_contract_test.go
- modified: internal/mcpserver/session_contract_test.go
- modified: cmd/codinho/main.go

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
- [x] Definir schema fechado de cada tipo de check.
- [x] Definir matriz de limites e capacidades por plataforma.

### Implementation
- [x] Implementar registry e validadores.
- [x] Implementar executor, timeout, output cap e redaction.
- [x] Implementar process-group cancellation.
- [x] Implementar política de rede e ambiente.
- [x] Expor check_run e eventos.
- [x] Conectar resultado do check ao critério `structural` de `step_evaluate` (R10).

### Validation
- [x] Executar testes adversariais em todos os inputs.
- [x] Testar timeout, flood, filhos, sinal e revisão obsoleta.
- [x] Executar revisão de segurança independente.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: O aluno executará código próprio, mas o modelo não deve obter shell.
- Options considered: run_command; container obrigatório; checks tipados.
- Decision: checks tipados e allowlisted com limites locais.
- Rationale: superfície mínima e viável para V1 local.
- Consequences: não equivale a isolamento de código hostil e isso deve ser documentado.

### Decision 2
- Date: 2026-08-23
- Context: A API/contract changes original desta spec diz "Adicionar
  check_run com input contendo somente sessão, check_id e revisão
  esperada" — sem `root`/`globs`. Mas o check precisa rodar em algum
  diretório real e restrito, e R7 exige vinculá-lo à "revisão do
  workspace".
- Options considered: (a) `check_run` também recebe `root`/`globs` como
  `workspace_observe` recebe; (b) `check_run` reusa o `root`/`globs` já
  estabelecidos pela baseline de `workspace_observe` para o
  `(session_id, step_id ativo)` — sem repetir esses parâmetros.
- Decision: (b). `WorkspaceService` passa a lembrar `root`/`globs` junto
  do baseline por `(SessionID, StepID)`; `ChecksService` os lê de lá.
  Chamar `check_run` antes de qualquer `workspace_observe` para o passo
  ativo retorna `INVALID_INPUT` (nenhuma baseline para reusar).
- Rationale: honra literalmente "input contendo somente sessão, check_id
  e revisão esperada" (mais `request_id`, mesma convenção de
  idempotência de todo o resto do contrato) sem duplicar um parâmetro
  que o fluxo natural (observar antes de checar) já torna redundante.
- Consequences: a skill do tutor deve orientar `workspace_observe` antes
  de `check_run` no mesmo passo; isso já é a ordem natural do fluxo
  descrito em PROJECT.md §21.1.

### Decision 3
- Date: 2026-08-23
- Context: R10 exige que `step_evaluate` derive o verdict de um
  critério `structural` do resultado real do check, mas
  `internal/assessment.Resolve` (artifact selado de
  feedback-evaluation-progression) hoje deriva esse verdict só da
  presença de evidência, por design (Decision 2 daquela spec).
- Options considered: (a) mudar a assinatura de `assessment.Resolve`
  para aceitar um lookup de outcome, tocando um artifact selado de uma
  spec já fechada; (b) aplicar a substituição de verdict em
  `internal/session.Service.StepEvaluate`, DEPOIS de chamar
  `assessment.Resolve` sem alterá-lo, quando a evidência citada for
  reconhecidamente uma evidência de check (`kind: "check"` no JSON).
- Decision: (b). `assessment.Resolve`/`criteria.go` permanecem
  intocados; `StepEvaluate` ganha uma dependência opcional em
  `*evidence.Store` e, só para critérios `structural` cujo evidence_id
  resolve para `kind: "check"`, reconstrói o `CriterionResult` com o
  verdict mapeado do `outcome` real (`pass→met`, `fail→not_met`,
  `skipped→not_applicable`, `error`/desconhecido→`unverifiable`).
- Rationale: evita reabrir/amendar uma spec já selada para uma mudança
  que pertence à evolução desta; mantém `assessment.Resolve` como a
  única fonte de verdade documentada para o caso SEM check real (evidência
  de workspace ou nenhuma), e `StepEvaluate` como o único lugar que
  conhece as duas fontes de evidência (workspace e check).
- Consequences: `session.New` ganha um parâmetro `*evidence.Store`
  (nil-safe: sem ele, o comportamento antigo de feedback-evaluation-
  progression é preservado sem alteração observável).

### Decision 4
- Date: 2026-08-23
- Context: R5 exige negar rede por padrão, salvo aprovação explícita;
  R9 exige que `check_run` não aceite parâmetro de comando/config livre,
  então a aprovação não pode vir da chamada.
- Decision: `curriculum.CheckAuthoring` ganha um campo aditivo
  `network bool` (default false); só um check cujo autor marcou
  `network: true` no pack recebe ambiente de rede liberado.
- Rationale: mantém a promessa de "nenhum parâmetro livre" em
  `check_run` — toda a configuração do check já vem fixada no catálogo
  antes da sessão existir.
- Consequences: catalog-schema-loader (spec já fechada) ganha um campo
  YAML novo e compatível; nenhum pack existente precisa mudar (default
  false = comportamento seguro).

## 6. Validation

### Strategy
Usar binaries de fixture controlados para simular ataques e estados de processo.

### Deterministic checks
- Test: go test ./internal/checks/... ./internal/application/... ./internal/session/... ./internal/mcpserver/...
- Lint: gofmt -l internal/checks internal/application internal/session internal/curriculum internal/workspace internal/mcpserver cmd/codinho
- Typecheck: go vet ./...
- Build: go build ./...
- Security / Contract: go test -race (traversal, env injection, output flood, process leak, redaction, network denial); govulncheck.

### Execution log
- `gofmt -l ...` (todos os pacotes tocados) → saída vazia (2026-08-23).
- `go vet ./...` (linux e `GOOS=windows`) → sem findings (2026-08-23).
- `go build ./...` → ok (2026-08-23).
- `go test ./... -race` → `ok` em todos os pacotes com testes, sem data races (2026-08-23).
- `govulncheck ./...` → "No vulnerabilities found." (2026-08-23).
- Smoke test real via `mcp.CommandTransport` contra `codinho serve` e
  `packs/go-first-steps.yaml` (desafio
  `go-data.slice-filter-preserve-input`, check `focused-tests`), com um
  módulo Go real (fora do repositório) implementando `Filter`:
  `session_start` → `session_get` → `workspace_observe` (baseline) →
  `check_run` (`focused-tests`, `outcome: pass` contra implementação
  correta) → `step_evaluate` (critério `structural` citando a evidência
  do check, `verdict: met`, `has_blocking_failure: false`) → implementação
  quebrada propositalmente → `check_run` de novo (`outcome: fail`, prova
  de que o check re-executa em vez de cachear) → `step_evaluate`
  (`verdict: not_met`, `has_blocking_failure: true`). Confirma R1–R10 de
  ponta a ponta com processos reais, não apenas chamadas de função
  (2026-08-23).

### Results summary
- `internal/checks` (puro, sem estado de sessão): `Resolve` valida
  `Runner` contra a allowlist fixa e `Package`/`TestPattern` contra
  traversal, path absoluto e formato de flag antes de qualquer
  `os/exec`; `Executor.Execute` roda em processo próprio (`Setpgid` em
  unix), aplica timeout via `context`, mata o grupo de processos inteiro
  no timeout (provado matando um `sleep 30` filho de um `go test`),
  aplica `maxOutputBytes` e `workspace.Redact` em stdout/stderr, nega
  rede por padrão (`GOPROXY=off`, `GOFLAGS=-mod=readonly`) salvo
  `NetworkApproved`, e nunca herda variáveis de ambiente fora da
  allowlist. `gofmt_check` é tratado como caso especial (exit code
  sempre 0; presença de stdout é que indica falha). `go_test`/
  `go_test_race` reconhecem "nada rodou" como `skipped`, não `pass`.
  `internal_ast` roda em processo (parser Go) para "verificadores
  internos" (R2).
- `internal/application.ChecksService.Run` resolve `check_id`
  exclusivamente contra `SessionService.ActiveChecks` (challenge fixado
  da sessão, requirement R1), reusa `root`/`globs` que
  `WorkspaceService` já guardava da última `workspace_observe` para o
  passo ativo (Decision 2 — `check_run` não recebe esses parâmetros),
  calcula fingerprint pós-execução e grava evidência via
  `WorkspaceService.PutEvidence`/`RecordEvidence`, reusando o mesmo
  `evidence_get` e a mesma detecção de obsolescência já entregues por
  workspace-observation-baselines (requirements R7, R8, R9).
- `internal/session.Service.checkOutcomeOverride` (Decision 3) troca o
  verdict de um critério `structural` para `met`/`not_met`/
  `not_applicable`/`unverifiable` conforme o `outcome` real quando a
  evidência citada é reconhecidamente de um check (`kind: "check"`),
  sem alterar `internal/assessment.Resolve` (artifact selado de
  feedback-evaluation-progression) — comprovado por 4 subtestes
  (pass/fail/skipped/error) mais 2 testes de não-regressão (evidência de
  workspace continua "met" por presença; sem `evidence.Store` wired,
  comportamento idêntico ao pré-existente).

### Requirement trace
- R1 [satisfied] report:internal/session/service.go (ActiveChecks) test:TestChecksRunRejectsUnknownCheckID
- R2 [satisfied] test:TestExecuteGoVetPassAndFail test:TestExecuteGoBuildFailsOnBrokenCode test:TestExecuteGofmtCheckDetectsUnformattedCode test:TestExecuteInternalASTChecksStructuralValidity
- R3 [satisfied] test:TestResolveRejectsPackageTraversal test:TestResolveRejectsInvalidTestPattern test:TestExecuteInternalASTRejectsPathEscapingRoot
- R4 [satisfied] test:TestExecuteTimeoutKillsTheWholeProcessGroup test:TestExecuteCapsOutputSize
- R5 [satisfied] test:TestExecuteDeniesNetworkByDefault test:TestExecuteAllowsNetworkOnlyWhenExplicitlyApproved
- R6 [satisfied] test:TestExecuteRedactsSecretsInCapturedOutput report:internal/checks/limits.go (capOutput)
- R7 [satisfied] report:internal/application/checks.go (fingerprint via workspace.Fingerprint) test:TestChecksRunEndToEndAfterWorkspaceObserve
- R8 [satisfied] test:TestChecksRunRequiresAPriorWorkspaceObservation report:internal/application/workspace.go (RootFor/RecordEvidence scope reuse)
- R9 [satisfied] report:internal/mcpserver/check_tools.go (checkRunArgs: session_id, check_id, expected_revision, request_id apenas)
- R10 [satisfied] test:TestStepEvaluateStructuralCriterionUsesRealCheckOutcome test:TestChecksRunFeedsStepEvaluateStructuralVerdict

### Known gaps
- Isolamento forte para código adversarial está fora da V1 (non-goal
  explícito).
- `Package`/`TestPattern` são validados como string (allowlist de
  caracteres, sem traversal), mas um pacote Go dentro da raiz que seja
  ele próprio um symlink para fora dela não é revalidado por
  `internal/checks` antes de `go build`/`go test` percorrê-lo — mesma
  classe de risco residual já registrada em workspace-observation-
  baselines para caminhos de arquivo individuais, aqui não fechada para
  pacotes Go inteiros. Mitigação real requer resolver o pacote para um
  diretório concreto e chamar `workspace.Root.Resolve` nele antes de
  invocar `go`, o que fica para security-privacy-hardening (covered
  follow-up abaixo).
- Suporte a process-group no Windows é best-effort (`Process.Kill` só no
  processo direto); `ProcessGroupSupported() == false` reporta isso
  honestamente em vez de fingir suporte total (Compatibility).

## 7. Final Report

### Delivered scope
Execução de checks pré-declarados (`go_test`, `go_test_race`, `go_vet`,
`gofmt_check`, `go_build`, `go_benchmark`, `internal_ast`) por `check_id`,
com allowlist fechada, contenção de path, ambiente mínimo, timeout,
cancelamento de grupo de processos, cap de output, redaction e negação
de rede por padrão — exposto como `check_run`. Conecta o resultado real
do check ao critério `structural` de `step_evaluate`
(feedback-evaluation-progression), fechando o gap que aquela spec havia
documentado. Nenhum parâmetro de comando livre; nenhuma edição de
arquivo do aluno (non-goals preservados).

### Files and modules changed
- `internal/checks/{model,registry,executor,limits,process_unix,process_windows}.go` + testes (criados)
- `internal/application/checks.go` + `checks_test.go` (criados), `workspace.go` (RootFor/PutEvidence/RecordEvidence, baseline agora guarda root/globs), `session.go` (ActiveChecks, NewSessionService ganha evidenceStore)
- `internal/session/service.go` (ActiveChecks, checkOutcomeOverride, `New` ganha `*evidence.Store`) + testes ajustados, `checks_override_test.go` (criado)
- `internal/curriculum/model.go` (`CheckAuthoring.Network`)
- `internal/workspace/redaction.go` (`Redact` exportado)
- `internal/mcpserver/check_tools.go` + `check_contract_test.go` (criados), `server.go`, `errors.go`, `contract_test.go` (ajustados)
- `cmd/codinho/main.go` (injeta `ChecksService`)

### Validation executed
- Command: go test ./... -race
- Result: `ok` em todos os pacotes com testes, sem data races.
- Command: govulncheck ./...
- Result: "No vulnerabilities found."
- Command: smoke test real via mcp.CommandTransport contra `codinho serve`, `packs/go-first-steps.yaml`, módulo Go real fora do repositório
- Result: check_run real (pass e fail, comprovando re-execução) alimentando corretamente o verdict de step_evaluate (R10).

### Residual risks
- O modelo de ameaça local (sem sandbox de kernel) permanece explícito no Non-goals e nesta seção.
- Ver Known Gaps: pacote Go que seja ele próprio um symlink para fora da raiz não é revalidado antes de `go build`/`go test` percorrê-lo.
- Suporte a process-group no Windows é best-effort, reportado honestamente por `ProcessGroupSupported() == false`.

### Follow-ups
- [covered: security-privacy-hardening] Executar threat model e hardening final.
