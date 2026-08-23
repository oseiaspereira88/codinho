---
slug: administrative-cli-fixtures
status: done
created_at: 2026-08-22
completed_at: 2026-08-23
supersedes:
depends_on: go-runtime-foundation, catalog-schema-loader, local-event-store, workspace-observation-baselines, safe-check-executor
priority: 140
components: cli, fixtures
delivers: surface:codinho-cli
---

# Spec: administrative-cli-fixtures

## 1. Intent

### Goal
Completar a CLI administrativa e preparar fixtures de desafios com consentimento, preview e proteção contra sobrescrita.

### Business value
Tornar instalação, autoria, diagnóstico e suporte utilizáveis sem duplicar a conversa pedagógica.

### Constraints
- CLI é determinística e não chama LLM.
- workspace prepare é acionado diretamente pelo usuário, não por tool MCP.
- Nenhum comando sobrescreve por padrão.

### Non-goals
- TUI interativa ou chat.
- Gerenciar sessões como interface pedagógica primária.

## 2. Requirements

### Functional
- R1: Expor serve, init, doctor, version e ajuda.
- R2: Expor catalog validate, list e show com saída humana e JSON estável.
- R3: Expor session inspect, progress show e progress export sem mutação pedagógica.
- R4: Implementar workspace prepare com destination explícito e plano de escrita.
- R5: Recusar destino inseguro, path traversal, symlink escape e overwrite por padrão.
- R6: Materializar fixture reproduzível e registrar manifesto com hashes.
- R7: Fornecer dry-run para toda operação que escrever múltiplos arquivos.
- R8: Usar códigos de saída e mensagens consistentes, com logs em stderr.

### Non-functional
- Comandos de leitura devem funcionar offline.
- Saída JSON deve ser versionada e testada por golden files.

### Security
- Nunca incluir secrets em export ou doctor.
- Validar todas as entradas de path e permissões.

### Compatibility
- Comandos mantêm aliases somente quando documentados; mudanças quebráveis exigem versão maior.

## 3. Technical Plan

### Affected areas
- internal/cli/, internal/fixtures/, cmd/codinho/

### Artifacts
- modified: internal/cli/root.go
- created: internal/cli/catalog.go
- created: internal/cli/session.go
- created: internal/cli/progress.go
- created: internal/cli/workspace.go
- created: internal/cli/init.go
- created: internal/cli/output.go
- created: internal/cli/state.go
- created: internal/fixtures/manifest.go
- created: internal/fixtures/materialize.go
- created: internal/fixtures/materialize_test.go
- created: internal/cli/cli_test.go
- created: cmd/codinho/cli_integration_test.go
- modified: internal/curriculum/model.go
- modified: internal/curriculum/validator.go
- modified: .pose/indexes/validation-matrix.json

### Delivery targets
- surface:codinho-cli module:cmd/codinho profile:cli-surface entrypoint:cmd/codinho/main.go

### API/contract changes
- Estabilizar comandos, flags, JSON e códigos de saída da CLI V1.

### Data/storage changes
- Manifestos de fixture e exports explícitos.

### Technical risks
- Fixture pode sobrescrever trabalho por erro de resolução.
- CLI pode duplicar lógica dos casos de uso.

## 4. Tasks

### Planning
- [x] Definir árvore de comandos e schemas JSON.
- [x] Definir protocolo de preview, confirmação e overwrite.

### Implementation
- [x] Implementar comandos de catálogo, sessão e progresso.
- [x] Implementar materialização transacional de fixture.
- [x] Implementar outputs e códigos de saída.
- [x] Reutilizar casos de uso do núcleo.
- [x] Criar testes CLI end-to-end.

### Validation
- [x] Testar destinos vazios, não vazios, symlinks e falha intermediária.
- [x] Testar goldens humana/JSON e compatibilidade.
- [x] Executar surface-check com reachability.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Debug challenges precisam de código inicial, mas o tutor não pode editar o aluno.
- Options considered: tool MCP de escrita; download manual; comando CLI explícito.
- Decision: workspace prepare somente pela CLI com preview e confirmação.
- Rationale: separa consentimento material de tutoria.
- Consequences: o fluxo de onboarding deve ensinar o comando.

### Decision 2
- Date: 2026-08-23
- Context: `session inspect` (R3) precisa funcionar a partir de um
  processo de CLI separado de `codinho serve`, mas `session.Service`
  mantém sessões apenas em memória, no processo que as criou
  (session-orchestration-disclosure Decision 3) — nunca reconstruídas por
  replay.
- Options considered: (a) adicionar reconstrução completa de sessão via
  replay ao `internal/session`, tratando essa spec como gatilho para
  reabrir aquele domínio; (b) ler diretamente o stream de eventos bruto
  da sessão para uma projeção de diagnóstico somente leitura, sem tocar
  `internal/session`.
- Decision: (b). `session inspect` replaya o stream do event store pelo
  session_id e projeta um resumo (challenge_id, estado inferido pelo
  último evento de ciclo de vida, contagem de eventos, revisão) sem
  nunca chamar `session.Service`.
- Rationale: reconstrução completa de sessão em memória é escopo de uma
  spec de sessão dedicada, não desta; a projeção somente leitura já
  satisfaz R3 ("sem mutação pedagógica") com evidência real do log
  durável, sem duplicar nem alterar o domínio de sessão existente.
- Consequences: `session inspect` mostra o histórico de eventos e um
  estado inferido, não o objeto de sessão vivo; uma spec futura que
  precise de replay completo de sessão deve revisitar este ponto.

### Decision 3
- Date: 2026-08-23
- Context: `workspace prepare` (R4, R6) precisa de conteúdo de arquivo
  inicial por desafio, mas `curriculum.ChallengeAuthoring` não tinha
  nenhum campo para isso, e nenhum pack autoral ainda declara um desafio
  `debug` (essa autoria pertence a `learning-practice-debug-modes`).
- Options considered: (a) bloquear `workspace prepare` até
  `learning-practice-debug-modes` definir o formato de fixture; (b)
  adicionar um campo `fixture` mínimo e genérico a
  `ChallengeAuthoring` agora, validado no loader, deixando a autoria de
  conteúdo real para specs futuras.
- Decision: (b). Adicionado `ChallengeAuthoring.Fixture
  []FixtureFileAuthoring{Path, Content}`, validado por
  `validateFixture` (rejeita path vazio/absoluto/travessia e
  duplicatas, ambos bloqueantes).
- Rationale: a materialização de fixture é a entrega desta spec; a
  autoria de desafios `debug` com fixture populado é entrega de
  `learning-practice-debug-modes` — o campo vazio em packs existentes
  não quebra nada e libera o mecanismo genérico para ser reusado por
  qualquer spec de conteúdo futura.
- Consequences: `workspace prepare <id>` retorna `ErrEmptyFixture`
  para qualquer desafio sem `fixture` declarado (todo o catálogo atual),
  o que é esperado até packs de debug existirem.

## 6. Validation

### Strategy
CLI real (binário compilado, subprocessos reais) em diretórios
temporários, mais testes internos de `internal/cli` e `internal/fixtures`
cobrindo destinos vazios, conflitantes, travessia de path e falha
intermediária.

### Deterministic checks
- Test: go test ./internal/cli/... ./internal/fixtures/... ./internal/curriculum/...
- Test (reachability/e2e real): go test ./cmd/codinho/... -run "TestCLIEntrypointReachesEveryCommand|TestCLIEndToEndWorkflow"
- Lint: gofmt -l internal/cli internal/fixtures internal/curriculum cmd/codinho
- Typecheck: go vet ./... ; GOOS=windows go vet ./...
- Build: go build ./...
- Security: govulncheck ./...
- Contract: pose surface-check --spec administrative-cli-fixtures --strict.

### Execution log
- `go build ./...` → ok (2026-08-23).
- `gofmt -l internal/cli internal/fixtures internal/curriculum cmd/codinho` → sem saída (2026-08-23).
- `go vet ./...` e `GOOS=windows go vet ./...` → sem diagnósticos (2026-08-23).
- `go test ./... -race` → `ok` em todos os pacotes com testes, incluindo `internal/cli` (novo) e `internal/fixtures` (novo), sem data races (2026-08-23).
- `go test ./cmd/codinho/... -run "TestCLIEntrypointReachesEveryCommand|TestCLIEndToEndWorkflow" -v` → PASS; binário real compilado e invocado via subprocesso para todo comando top-level e um fluxo `init → catalog validate/list/show → workspace prepare → session inspect → progress show` (2026-08-23).
- `govulncheck ./...` → "No vulnerabilities found." (2026-08-23).
- `pose validate --strict --json .pose/results/delivery-validation.json` → checks de módulo `.` incluindo `cli-reachability` (evidenceClass reachability) e `cli-e2e` (evidenceClass e2e), registrados em `.pose/indexes/validation-matrix.json` (2026-08-23).

### Results summary
CLI administrativa completa (serve/init/doctor/version/help, catalog
validate/list/show, session inspect, progress show/export, workspace
prepare) e mecanismo de fixture reproduzível entregues; conteúdo real de
fixture para desafios `debug` fica para `learning-practice-debug-modes`
(Decision 3).

### Requirement trace
- R1 [satisfied] internal/cli/root.go (usage lista serve/init/doctor/version/help) + test:TestHelpListsAllCommands, test:TestCLIEntrypointReachesEveryCommand.
- R2 [satisfied] internal/cli/catalog.go (validate/list/show) + test:TestCatalogValidateOK, test:TestCatalogListJSON, test:TestCatalogShowChallenge.
- R3 [satisfied] internal/cli/session.go, internal/cli/progress.go + test:TestSessionInspectReportsRecordedEvents, test:TestProgressShowAndExportReflectRecordedEvidence.
- R4 [satisfied] internal/cli/workspace.go + test:TestWorkspacePrepareMaterializesAndProtectsAgainstOverwrite.
- R5 [satisfied] internal/fixtures/materialize.go (safeRelPath, authorizeDest, ErrWouldOverwrite) + test:TestPlanRejectsPathTraversal, test:TestPlanRejectsAbsolutePath, test:TestMaterializeRefusesOverwriteByDefault.
- R6 [satisfied] internal/fixtures/manifest.go + test:TestMaterializeWritesFilesAndManifest.
- R7 [satisfied] internal/fixtures.Plan (--dry-run) + test:TestWorkspacePrepareDryRunDoesNotWrite, test:TestPlanDoesNotWrite.
- R8 [satisfied] exit codes exitOK/exitError/exitUsage consistentes em internal/cli/root.go + todos os testes de internal/cli.

### Known gaps
- Instaladores e distribuição pertencem a installation-documentation-ci.
- Conteúdo real de fixture (desafios `debug`) pertence a
  learning-practice-debug-modes (Decision 3); o mecanismo de
  materialização já está pronto e testado.
- Sem golden files formais para saída humana; a saída JSON é testada por
  decodificação e estabilidade entre execuções (`TestCatalogValidateJSONIsStableAcrossRuns`),
  não por comparação byte-a-byte contra um arquivo golden.

## 7. Final Report

### Delivered scope
CLI administrativa completa e mecanismo de fixture reproduzível com
consentimento explícito, preview (`--dry-run`) e proteção contra
sobrescrita por padrão.

### Files and modules changed
- internal/cli/{root,catalog,session,progress,workspace,init,output,state}.go (novos comandos, reutilizando internal/application e internal/curriculum).
- internal/fixtures/{manifest,materialize}.go (materialização transacional com staging + rename).
- internal/curriculum/{model,validator}.go (campo `Fixture` e validação de path).
- cmd/codinho/cli_integration_test.go (reachability + e2e reais via subprocesso).
- .pose/indexes/validation-matrix.json (checks cli-reachability, cli-e2e).

### Validation executed
- Command: go test ./... -race
- Result: ok em todos os pacotes.
- Command: go test ./cmd/codinho/... -run "TestCLIEntrypointReachesEveryCommand|TestCLIEndToEndWorkflow"
- Result: PASS (evidência real de composição/reachability e fluxo e2e).
- Command: govulncheck ./...
- Result: No vulnerabilities found.

### Residual risks
- Windows path semantics exigem CI futuro (GOOS=windows go vet passou, mas sem execução real em Windows).
- Fixture ainda não tem conteúdo autoral real; primeira autoria real acontece em learning-practice-debug-modes.

### Follow-ups
- [covered: installation-documentation-ci] Documentar e testar instalação da CLI.
- [covered: learning-practice-debug-modes] Autorar desafios `debug` com fixture real e validar workspace prepare de ponta a ponta com conteúdo real.
