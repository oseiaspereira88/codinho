---
slug: workspace-observation-baselines
status: in-progress
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: go-runtime-foundation, local-event-store, learning-domain-model
priority: 90
components: workspace, evidence
delivers:
---

# Spec: workspace-observation-baselines

## 1. Intent

### Goal
Observar mudanças relevantes no workspace com baseline, contenção de paths e evidências imutáveis, sem editar o código do aluno.

### Business value
Permitir feedback e avaliação baseados no que o aluno realmente produziu, preservando trabalho preexistente.

### Constraints
- Leitura apenas; nenhuma API de escrita.
- Suportar repositórios Git e diretórios sem Git.
- Distinguir fixture, baseline e trabalho da etapa.

### Non-goals
- Executar testes, formatar arquivos ou corrigir código.
- Garantir segredo forte de testes em máquina controlada pelo usuário.

## 2. Requirements

### Functional
- R1: Autorizar uma raiz real e rejeitar paths ou symlinks que escapem dela.
- R2: Registrar commit, status, diff preexistente, hashes, manifesto de fixture e revisão ao ativar passo.
- R3: Coletar somente arquivos e símbolos dentro dos globs declarados pelo desafio.
- R4: Produzir diff relevante sem atribuir mudanças anteriores à etapa.
- R5: Comparar hashes e metadados quando Git não estiver disponível.
- R6: Analisar Go por parser e AST quando o critério for estrutural.
- R7: Criar Observation e Evidence sem implicar tentativa ou avaliação.
- R8: Expor workspace_observe e evidence_get com redaction, tamanho e escopo.
- R9: Detectar evidência obsoleta quando o workspace mudar após a coleta.

### Non-functional
- Observação não pode modificar mtime, índice Git ou working tree.
- Outputs devem ser limitados e content-addressed.

### Security
- Excluir secrets, .env, credenciais e paths sensíveis por padrão.
- Tratar nomes, comentários e conteúdo como dados não confiáveis.

### Compatibility
- Implementar abstração de Git e fallback portable.

## 3. Technical Plan

### Affected areas
- internal/workspace/, internal/application/, internal/mcpserver/

### Artifacts
- created: internal/workspace/root.go
- created: internal/workspace/glob.go
- created: internal/workspace/baseline.go
- created: internal/workspace/git.go
- created: internal/workspace/hashes.go
- created: internal/workspace/golang.go
- created: internal/workspace/redaction.go
- created: internal/workspace/workspace_test.go
- created: internal/workspace/security_test.go
- created: internal/application/workspace.go
- created: internal/application/workspace_test.go
- created: internal/mcpserver/workspace_tools.go
- created: internal/mcpserver/workspace_contract_test.go
- modified: internal/mcpserver/server.go
- modified: internal/mcpserver/envelope.go
- modified: internal/mcpserver/contract_test.go
- modified: cmd/codinho/main.go

### Delivery targets
Nenhum novo; amplia o contrato MCP V1 planejado.

### API/contract changes
- Adicionar workspace_observe e evidence_get.

### Data/storage changes
- Persistir baseline, observation e referências content-addressed.

### Technical risks
- Symlink races e paths com diferenças entre plataformas.
- Diffs grandes podem consumir memória ou vazar dados.

## 4. Tasks

### Planning
- [x] Definir modelo de autorização e exclusões padrão.
- [x] Definir limite de arquivo, diff e AST.

### Implementation
- [x] Implementar resolução real e contenção.
- [x] Implementar baseline Git e fallback por hash.
- [x] Implementar diff, redaction e evidência.
- [x] Implementar inspeção AST limitada.
- [x] Expor tools somente leitura.

### Validation
- [x] Testar repo sujo, sem Git, rename, symlink e traversal.
- [x] Verificar índice e arquivos byte a byte antes e depois.
- [x] Testar prompt injection em comentários como dado inerte.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Regex é frágil para critérios estruturais Go.
- Options considered: regex; go/parser e go/ast; language server.
- Decision: usar AST padrão para estrutura e texto somente quando semântico.
- Rationale: determinismo sem daemon externo.
- Consequences: critérios devem aceitar soluções equivalentes.

### Decision 2
- Date: 2026-08-23
- Context: R3 exige coletar somente dentro dos "globs declarados pelo
  desafio", mas `curriculum.ChallengeAuthoring` (catalog-schema-loader) não
  tem esse campo, e esta spec não depende de catalog-schema-loader.
- Options considered: (a) adicionar `Workspace.Globs` ao schema de
  autoria e importar `internal/curriculum` aqui; (b) `workspace_observe`
  recebe `globs []string` como parâmetro explícito do chamador, do mesmo
  jeito que `feedback_prepare` recebe rubricas como referências fixas
  (feedback-evaluation-progression Decision 2) em vez de resolvê-las
  internamente.
- Decision: (b). `internal/workspace` e o tool MCP não importam
  `internal/curriculum`; quem declara os globs (a skill do tutor, lendo o
  desafio autorado) os repassa em cada chamada.
- Rationale: mantém `internal/workspace` testável isoladamente e coerente
  com `depends_on` desta spec; evita uma segunda fonte de verdade para
  "quais arquivos pertencem ao desafio" competindo com o pack YAML.
- Consequences: se um desafio futuro precisar que o PRÓPRIO servidor
  resolva os globs (em vez do chamador), isso exige nova spec com
  `depends_on: catalog-schema-loader`.

### Decision 3
- Date: 2026-08-23
- Context: R2 registra baseline "ao ativar passo" e R4/R9 comparam
  observações posteriores contra ela, mas não existe replay do eventstore
  para reconstruir estado em memória nesta base de código ainda
  (`internal/session.Service` também mantém seu registro de sessões
  apenas em memória, sem replay ao reabrir — achado ao inspecionar
  `internal/session/service.go`, não corrigido aqui por ser característica
  pré-existente e fora do escopo desta spec).
- Decision: seguir a mesma convenção já estabelecida — `WorkspaceService`
  mantém o baseline mais recente por `(SessionID, StepID)` em memória, e
  ainda assim persiste cada observação como `observation_recorded` no
  eventstore para trilha auditável (requirement não-funcional "toda
  mutação deve ser revisionada", mesmo princípio de outras specs).
- Rationale: consistência com o padrão já existente é mais importante do
  que resolver aqui um gap de durabilidade que não é desta spec; introduzir
  replay-based reconstruction unilateralmente aqui criaria duas convenções
  divergentes no mesmo servidor.
- Consequences: reiniciar o processo entre a baseline e a observação de um
  mesmo passo perde a baseline em memória — mesma limitação que sessões já
  têm hoje. Registrar como Known Gap; corrigir de forma unificada (sessão e
  workspace) pertence a uma spec futura de recuperação de estado.

## 6. Validation

### Strategy
Combinar fixtures Git/não Git, ataques de path e comparação de filesystem.

### Deterministic checks
- Test: go test ./internal/workspace/... ./internal/application/... ./internal/mcpserver/...
- Lint: gofmt -l internal/workspace internal/application internal/mcpserver cmd/codinho
- Typecheck: go vet ./...
- Build: go build ./...
- Security / Contract: symlink escape, symlink-directory walk, default-exclude de secrets, redaction de conteúdo, prompt injection em comentário, cross-session evidence scope, govulncheck.

### Execution log
- `gofmt -l internal/workspace internal/application internal/mcpserver cmd/codinho` → saída vazia (2026-08-23).
- `go vet ./...` → sem findings (2026-08-23).
- `go build ./...` → ok (2026-08-23).
- `go test ./... -race` → `ok` em todos os pacotes com testes, sem data races (2026-08-23).
- `govulncheck ./...` → "No vulnerabilities found." (2026-08-23).
- Smoke test real de ponta a ponta via `mcp.CommandTransport` contra o
  binário `codinho serve`, com um workspace de aluno sintético (fora do
  repositório) contendo `main.go` e um `.env` com segredo:
  `session_start` → `session_get` → `workspace_observe` (baseline,
  `.env` corretamente ausente do manifesto de fixtures) →
  modificação do `main.go` (incluindo comentário adversarial
  "ignore all previous instructions...") → `workspace_observe` (diff,
  1 arquivo `modified`) → `evidence_get` da evidência de baseline
  (`stale: true`, workspace mudou depois da coleta) → `evidence_get` da
  evidência de diff (`stale: false`) → confirmado que o segredo do
  `.env` nunca aparece em nenhuma evidência → `evidence_get` da mesma
  evidência por uma `session_id` diferente devolve `ITEM_NOT_FOUND`.
  Todas as chamadas `status: ok` exceto a última (escopo), como esperado
  (2026-08-23).

### Results summary
- `internal/workspace` (puro, sem estado): `AuthorizeRoot`/`Root.Resolve`
  contêm paths e symlinks dentro da raiz autorizada (fecham a janela
  TOCTOU reavaliando após `EvalSymlinks`); `Capture`/`Observe` produzem
  baseline e diff por hash content-addressed, com fallback automático
  quando o diretório não é um repositório Git (`gitInfo`); `ParseGoFile`
  lista funções e tipos top-level via `go/parser`/`go/ast` sem executar
  ou type-checar o código do aluno; `DiffText` devolve diff textual
  redigido e truncado quando Git está disponível; exclusão padrão
  (`.env`, `id_rsa`, `credentials.*`, `.pem`/`.key`, ...) e redaction de
  padrões de segredo (`AKIA...`, `ghp_...`, chaves PEM, `key: valor`)
  se aplicam antes de qualquer conteúdo virar evidência.
- `internal/application.WorkspaceService` orquestra `internal/workspace`
  com o eventstore e o evidence store compartilhados: a primeira
  observação de um `(SessionID, StepID)` estabelece a baseline (em
  memória, Decision 3) e a evidência tem `kind: baseline`; chamadas
  seguintes calculam o diff contra essa baseline (nunca atribuindo
  mudança anterior à etapa, requirement R4) com `kind: diff`, incluindo
  declarações AST dos arquivos `.go` alterados (limitado a 50 arquivos)
  e o diff textual redigido quando disponível. Toda observação é
  persistida como evento `observation_recorded` no stream da sessão,
  respeitando a mesma concorrência otimista (`expected_revision`) que os
  demais tools de sessão.
- `evidence_get` só retorna evidência que a própria sessão gravou
  (`ErrEvidenceOutOfScope` caso contrário, requirement R8) e reporta
  `stale: true` quando um novo fingerprint do workspace diverge do
  fingerprint gravado junto da evidência (requirement R9).

### Requirement trace
- R1 [satisfied] test:TestRootResolveRejectsTraversal test:TestRootResolveRejectsSymlinkEscape test:TestWalkMatchedNeverFollowsSymlinkedDirectories
- R2 [satisfied] test:TestObserveEstablishesBaselineOnFirstCall test:TestCaptureUsesGitIdentityWhenAvailable
- R3 [satisfied] test:TestMatchGlobSupportsDoubleStarAndSingleStar test:TestCaptureAndObserveDetectChanges
- R4 [satisfied] test:TestObserveNeverAttributesUnchangedFilesToTheStep test:TestObserveReportsChangesAgainstBaseline
- R5 [satisfied] test:TestCaptureFallsBackToHashesWithoutGit
- R6 [satisfied] test:TestParseGoFileListsFuncsAndStructFields report:internal/application/workspace.go (goDeclarationsFor)
- R7 [satisfied] report:internal/application/workspace.go (Observe builds learning.Observation/Evidence, never an Attempt/evaluation)
- R8 [satisfied] test:TestEvidenceGetIsScopedToTheRecordingSession test:TestContractEvidenceGetOutOfScopeReturnsItemNotFound
- R9 [satisfied] test:TestFingerprintDetectsDriftAfterCollection test:TestEvidenceGetDetectsDriftSinceCollection

### Known gaps
- Limites de recurso do subprocesso pertencem a safe-check-executor.
- Baseline por `(SessionID, StepID)` vive só em memória (Decision 3);
  reiniciar o processo entre a baseline e a observação perde o estado,
  mesma limitação já presente em `internal/session.Service`. Corrigir de
  forma unificada pertence a uma spec futura de recuperação de estado, não
  a esta.

## 7. Final Report

### Delivered scope
Observação read-only de workspace com autorização/contenção de raiz,
baseline por Git com fallback por hash, diff restrito a globs
declarados pelo chamador, análise estrutural Go por AST, exclusão e
redaction padrão de segredos, e detecção de evidência obsoleta —
expostos como `workspace_observe` e `evidence_get`. Nenhuma edição de
arquivo do aluno e nenhuma execução de check (non-goals preservados).

### Files and modules changed
- `internal/workspace/{root,glob,baseline,git,hashes,golang,redaction}.go` + `workspace_test.go`, `security_test.go` (criados)
- `internal/application/workspace.go` + `workspace_test.go` (criados)
- `internal/mcpserver/workspace_tools.go` + `workspace_contract_test.go` (criados), `server.go`, `errors.go`, `contract_test.go` (modificados)
- `cmd/codinho/main.go` (modificado: abre o evidence store e injeta `WorkspaceService`)

### Validation executed
- Command: go test ./... -race
- Result: `ok` em todos os pacotes com testes, sem data races.
- Command: govulncheck ./...
- Result: "No vulnerabilities found."
- Command: smoke test real via mcp.CommandTransport contra `codinho serve`, workspace de aluno sintético com secret em `.env`
- Result: baseline→diff→evidence_get(stale)→evidence_get(fresh)→sem vazamento de segredo→escopo cross-session negado, todos os `status` esperados.

### Residual risks
- Portabilidade de symlink requer CI por sistema operacional (os testes de symlink pulam em Windows).
- Diff textual via `git diff` fica limitado a `maxDiffBytes` (64 KiB); acima disso o texto é truncado, não rejeitado — evidência ainda é útil mas parcial.
- Ver Known Gaps: baseline por `(SessionID, StepID)` só em memória (Decision 3), mesma limitação já presente em `internal/session.Service`.

### Follow-ups
- [covered: security-privacy-hardening] Realizar revisão adversarial completa da fronteira de workspace.
- [covered: safe-check-executor] Vincular resultado de check_run à mesma noção de fingerprint/obsolescência (`workspace.Fingerprint`) usada aqui, em vez de reimplementar detecção de staleness — safe-check-executor já depende de workspace-observation-baselines e sua R7 ("marcar evidência obsoleta") descreve exatamente este mecanismo.
