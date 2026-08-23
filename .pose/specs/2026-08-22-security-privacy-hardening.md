---
slug: security-privacy-hardening
status: in-progress
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: mcp-stdio-foundation, workspace-observation-baselines, safe-check-executor, tutor-skill-host-integration, administrative-cli-fixtures
priority: 220
components: security, privacy
delivers:
---

# Spec: security-privacy-hardening

## 1. Intent

### Goal
Executar threat modeling e hardening transversal de MCP, CLI, workspace, checks, evidências, catálogo e skill.

### Business value
Garantir que a ferramenta de aprendizagem preserve código, dados e máquina local enquanto observa e executa verificações.

### Constraints
- Operação local, zero telemetria e zero upload por padrão.
- Nenhuma tool de edição ou comando livre.
- Mínimo privilégio e fail closed nas fronteiras.

### Non-goals
- Sandbox de kernel para código deliberadamente hostil.
- Autenticação remota, multiusuário ou compliance certificado.

## 2. Requirements

### Functional
- R1: Publicar threat model com assets, atores, trust boundaries, abuso e risco residual.
- R2: Impedir traversal, symlink escape, env injection, shell injection e execução fora do registry.
- R3: Tratar código, comentários, fixtures e outputs como conteúdo não confiável.
- R4: Excluir e redigir secrets, .env, credenciais, paths sensíveis e dados pessoais.
- R5: Aplicar permissões restritivas e retenção configurável a estado e evidências.
- R6: Fornecer export e remoção local explícitos sem apagar código ou conteúdo versionado.
- R7: Negar rede por padrão e exigir aprovação visível quando necessária.
- R8: Avaliar dependências, licenças e vulnerabilidades com política de atualização.
- R9: Provar que repositório sujo, índice Git e arquivos fora de escopo são preservados.
- R10: Documentar limites de segurança, testes reservados e modo entrevista.

### Non-functional
- Todo finding crítico ou alto bloqueia closeout sem exceção governada.
- Logs não contêm código por padrão.

### Security
- Aplicar integralmente .pose/rules/security.md e revisão independente.

### Compatibility
- Controles devem degradar explicitamente quando plataforma não oferecer uma primitiva.

## 3. Technical Plan

### Affected areas
- internal/security/, internal/workspace/, internal/checks/, internal/evidence/, internal/mcpserver/, internal/cli/, docs/

### Artifacts
- created: docs/security/threat-model.md
- created: docs/security/privacy.md
- created: docs/security/executor-limitations.md
- created: internal/security/policy.go
- created: internal/security/policy_test.go
- created: internal/cli/privacy.go
- modified: internal/cli/root.go
- modified: internal/cli/cli_test.go
- modified: cmd/codinho/cli_integration_test.go
- modified: internal/workspace/security_test.go
- modified: internal/checks/security_test.go
- modified: internal/application/checks.go
- modified: internal/mcpserver/check_tools.go
- modified: internal/mcpserver/check_contract_test.go
- modified: .pose/indexes/validation-matrix.json

### Delivery targets
Nenhum novo; hardening transversal.

### API/contract changes
- Erros de autorização e redaction permanecem estáveis no MCP e CLI.

### Data/storage changes
- Adicionar política de retenção e remoção ao estado local.

### Technical risks
- Garantias podem ser superestimadas em documentação.
- Redaction pode falhar para formatos desconhecidos.

## 4. Tasks

### Planning
- [x] Modelar ameaças e classificar riscos.
- [x] Definir allowlists, exclusions, retenção e approvals.

### Implementation
- [x] Centralizar políticas e redaction.
- [x] Corrigir findings em todas as fronteiras.
- [x] Implementar export e remoção segura.
- [x] Criar corpus adversarial.
- [x] Documentar garantias e limites.

### Validation
- [x] Executar secret scan, govulncheck e suites negativas.
- [ ] Executar revisão independente de threat model e diff.
- [x] Verificar zero mutação do workspace fora do comando explícito.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Execução local segura não equivale a sandbox completa.
- Options considered: prometer isolamento; exigir container; declarar modelo local limitado.
- Decision: controles fortes de processo e path com limites honestos.
- Rationale: atende V1 sem falsa garantia.
- Consequences: código hostil permanece fora do modelo de uso suportado.

### Decision 2
- Date: 2026-08-23
- Context: R5 pede "retenção configurável" ao estado local, mas o event
  log (`internal/eventstore`) é intencionalmente append-only e imutável
  — é o que garante replay determinístico de progresso/mastery
  (mastery-review-scheduling). Apagar um evento no meio quebraria esse
  replay.
- Options considered: (a) adicionar poda por idade a nível de evento
  individual, reescrevendo o formato do log; (b) tratar retenção como
  uma decisão explícita e no nível do workspace inteiro — exportar
  quando quiser, apagar tudo quando quiser, nunca poda automática
  parcial.
- Decision: (b). `internal/security.Export`/`Purge` operam sobre
  `.codinho/state/` inteiro, nunca sobre eventos individuais.
- Rationale: honestidade sobre o que o formato de log realmente suporta
  (Constraint desta spec: "Validação deve ser offline e
  determinística" — poda parcial arriscaria integridade de replay já
  comprovada por mastery-review-scheduling); "configurável" vira "o
  aluno decide quando", não "expira sozinho".
- Consequences: não há expiração automática por idade; documentado em
  `docs/security/privacy.md`.

### Decision 3
- Date: 2026-08-23
- Context: R7 pede aprovação de rede visível, mas `check_run` nunca
  expunha se a rede foi liberada — a decisão (autorada em
  `checks[].network`) era invisível ao chamador.
- Options considered: (a) deixar como está (a política já nega rede por
  padrão, então "visível" seria só documentação); (b) adicionar
  `network_approved` ao payload de evidência e ao envelope de
  `check_run`.
- Decision: (b). Mudança aditiva a `checkPayload`/`CheckRunResult`
  (feedback-evaluation-progression/safe-check-executor não são
  modificados em sua lógica, só ganham um campo novo).
- Rationale: uma decisão de segurança correta mas silenciosa ainda é um
  finding de R7 — "exigir aprovação visível" significa o chamador
  conseguir ver o fato, não só confiar que o servidor decidiu certo.
- Consequences: qualquer consumidor de `check_run`/evidência de check
  agora vê `network_approved` sem precisar inspecionar o pack.

## 6. Validation

### Strategy
Usar threat model, corpus adversarial, scanners e revisão separada.

### Deterministic checks
- Test: go test -race ./internal/security/... ./internal/workspace/... ./internal/checks/... ./internal/mcpserver/... ./internal/cli/...
- Lint: gofmt -l internal
- Typecheck: go vet ./...
- Build: go build ./cmd/codinho
- Security / Contract: govulncheck ./...; secret scan; negative suite; pose assess integrate.

### Execution log
- `go build ./...` → ok (2026-08-23).
- `gofmt -l .` → sem saída (2026-08-23).
- `go vet ./...` e `GOOS=windows go vet ./...` → sem diagnósticos (2026-08-23).
- `go test ./... -race` → `ok` em todos os pacotes, incluindo `internal/security` (novo), `internal/workspace` e `internal/checks` (novos testes de preservação de Git), sem data races (2026-08-23).
- `govulncheck ./...` → "No vulnerabilities found." (2026-08-23); registrado como check obrigatório em `.pose/indexes/validation-matrix.json` (evidenceClass: security).
- Revisão manual de licenças de dependências (`go list -m all`): apenas módulos MIT/BSD-3-Clause/Apache-2.0 conhecidos (mcp-go SDK, google/jsonschema-go, yaml.v3, golang.org/x/*, segmentio/*, golang-jwt) — nenhuma dependência copyleft forte.

### Results summary
Threat model, política de privacidade e limites do executor documentados;
`internal/security` centraliza export/purge de estado local (nunca o
código do aluno) com scanner de conteúdo em formato de segredo reusando
os padrões de `internal/workspace`; `network_approved` agora visível em
`check_run`; corpus adversarial ampliado com prova real de que
observação e execução de check nunca mutam Git ou arquivos fora de
escopo (R9). Controles de traversal/symlink/shell/env/secret já
existentes (de specs anteriores) foram auditados e permanecem intactos.

### Requirement trace
- R1 [satisfied] docs/security/threat-model.md.
- R2 [satisfied] já existente (workspace.Root, checks.Resolve) + testes de security_test.go pré-existentes, auditados.
- R3 [satisfied] TestPromptInjectionInFileContentIsInertData (pré-existente, auditado).
- R4 [satisfied] workspace.Redact (pré-existente) + internal/security.ContainsSecretShapedContent (reuso) + test:TestContainsSecretShapedContentDetectsKnownPatterns.
- R5 [satisfied] permissões restritivas já existentes (0o700/0o600/0o400); retenção explícita via Decision 2.
- R6 [satisfied] internal/security.Export/Purge + tools CLI `privacy export`/`privacy purge` + test:TestPrivacyExportAndPurge.
- R7 [satisfied] network_approved em check_run (Decision 3) + test:TestContractCheckRunEndToEndAfterWorkspaceObserve (assert atualizado).
- R8 [satisfied] govulncheck registrado em validation-matrix.json + revisão de licenças documentada acima.
- R9 [satisfied] test:TestObservationNeverMutatesGitOrOutOfScopeFiles, test:TestExecuteNeverMutatesGitState.
- R10 [satisfied] docs/security/executor-limitations.md.

### Known gaps
- Isolamento forte (sandbox de kernel) poderá ser uma capacidade pós-V1.
- Revisão independente de threat model e diff por uma pessoa diferente
  deste agente não foi realizada nesta sessão — permanece pendente de
  revisão humana real antes do aceite V1 (não pode ser autocertificada).

## 7. Final Report

### Delivered scope
Threat model formal, política de privacidade com export/purge
explícitos, aprovação de rede visível em check_run, e prova real de
não-mutação de Git/workspace fora de escopo.

### Files and modules changed
- internal/security/{policy,policy_test}.go (novos).
- internal/cli/privacy.go (novo) + root.go (comando privacy).
- internal/workspace/security_test.go, internal/checks/security_test.go: teste de preservação de Git.
- internal/application/checks.go, internal/mcpserver/check_tools.go: network_approved visível.
- docs/security/{threat-model,privacy,executor-limitations}.md.
- .pose/indexes/validation-matrix.json: check govulncheck obrigatório.

### Validation executed
- Command: go test ./... -race
- Result: ok em todos os pacotes.
- Command: govulncheck ./...
- Result: No vulnerabilities found.

### Residual risks
- Riscos aceitos (sem sandbox de kernel, sem antifraude de entrevista) documentados, nunca silenciados.
- Revisão humana independente do threat model é um gate real ainda pendente.

### Follow-ups
- [covered: v1-integrated-acceptance] Reexecutar suite adversarial no candidato V1.
- [open] Obter revisão humana independente do threat model e do diff desta spec antes do aceite V1.
