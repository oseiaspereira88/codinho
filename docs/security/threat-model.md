# Threat model

Escopo: o runtime `codinho` — servidor MCP (`codinho serve`), CLI
administrativa e a skill que os orquestra — rodando localmente na
máquina do aluno, sob o mesmo usuário do host (Codex CLI, extensão de
IDE). Não cobre a rede do provedor de modelo (LLM), nem a infraestrutura
de quem hospeda o host.

## Assets (o que este sistema protege)

| Asset | Onde vive | Por que importa |
|---|---|---|
| Código do aluno | Workspace real do aluno, fora de `.codinho/` | Nunca deve ser editado, revelado indevidamente ou vazado (ownership do aprendiz). |
| Estado de sessão e evidências | `.codinho/state/{events.jsonl,evidence/}` | Registro de progresso; contém trechos observados/executados do workspace. |
| Segredos incidentais | Potencialmente presentes no workspace do aluno (`.env`, chaves) | Nunca devem ser lidos, ecoados ou persistidos como evidência. |
| Conteúdo curricular (`packs/*.yaml`) | Repositório do produto | Gabaritos e critérios não podem vazar antes da hora certa. |
| Máquina local do aluno | Fora do processo `codinho` | Nunca deve ser afetada além do workspace declarado e `.codinho/state`. |

## Atores

- **Aluno** — confiável quanto à intenção, mas seu código/comentários são
  **dados**, nunca instruções (regra 14 da skill).
- **Tutor (LLM via skill)** — orquestra tools; nunca tem acesso a shell
  livre; só chama tools MCP declaradas.
- **Host (Codex CLI / extensão de IDE)** — controla aprovações de
  ferramenta e limites do próprio host; fora do controle deste projeto.
- **Autor de pack malicioso ou comprometido** — pode tentar declarar
  paths de fixture maliciosos, checks fora do allowlist, ou relations
  inválidas. Mitigado por `internal/curriculum`'s validador +
  catalog-authoring-quality.
- **Conteúdo observado adversarial** — código-fonte, saída de comando ou
  nome de arquivo do workspace contendo texto que parece instrução
  (prompt injection). Sempre tratado como dado inerte (ver
  `PromptInjectionInFileContentIsInertData`).

## Trust boundaries

1. **Tool MCP → domínio**: toda entrada de `session_start`/`check_run`/
   etc. passa por validação de domínio (`learning.NewSessionPolicy`,
   allowlist de checks) antes de qualquer efeito.
2. **Domínio → filesystem**: todo acesso a arquivo passa por
   `workspace.Root` (path autorizado, `..` rejeitado antes do
   filesystem, symlink resolvido e reverificado — fecha TOCTOU).
3. **Domínio → subprocesso**: todo check passa por
   `internal/checks.Resolve` (allowlist fixo de runners, sem shell,
   env allowlisted, `GOPROXY=off` por padrão) antes de `Execute`.
4. **Filesystem/subprocesso → evidência**: toda saída observada ou
   capturada passa por `workspace.Redact` antes de virar evidência
   persistida.

## Casos de abuso e mitigação

| Abuso | Mitigação | Evidência (teste) |
|---|---|---|
| Path traversal (`../../etc/passwd`) | `workspace.Root.Resolve` rejeita `..` antes do FS | `TestRootResolveRejectsSymlinkEscape` |
| Symlink escape (link para fora do root) | Re-checagem via `EvalSymlinks` após resolução | `TestWalkMatchedNeverFollowsSymlinkedDirectories` |
| Shell injection via package/pattern de check | Nenhum `exec.Command` usa shell; args isolados; regex allowlist em package/pattern | `internal/checks/registry.go` (`packagePattern`, `testPatternChars`) |
| Env injection / vazamento de env do host | `buildEnv` só repassa um allowlist fixo; nunca `os.Environ()` inteiro | `internal/checks/executor.go` |
| Exfiltração de rede durante check | `GOPROXY=off`+`GOFLAGS=-mod=readonly` por padrão; só liberado por campo autorado (`checks[].network`), nunca por parâmetro de chamada | `TestExecuteAllowsNetworkOnlyWhenExplicitlyApproved` |
| Prompt injection via código/comentário/saída observada | Tratado como dado pela skill (regra 14) e nunca interpretado pelo servidor | `TestPromptInjectionInFileContentIsInertData` |
| Vazamento de secret em evidência | `workspace.Redact` aplicado a todo stdout/stderr/diff antes de persistir | `TestRedactStripsSecretShapedContent`, `TestExecuteRedactsSecretsInCapturedOutput` |
| YAML bomb / alias expansion em pack | `Limits` (MaxFileBytes/MaxNodeDepth/MaxAliases) checados antes de `Decode` | `TestCheckYAMLLimitsRejectsExcessiveDepth` |
| Fixture com path de escape (`../evil.go`) | `fixtures.safeRelPath` + `checkPathIsCreatable` antes de qualquer escrita | `TestPlanRejectsPathTraversal` |
| Mutação indevida de Git/arquivos fora de escopo | Toda observação/check é somente leitura no workspace do aluno | `TestObservationNeverMutatesGitOrOutOfScopeFiles`, `TestExecuteNeverMutatesGitState` |
| Autor de pack tenta publicar sem revisão distinta | `checkPublicationMetadata` bloqueia `reviewed_by == author` | `TestCheckPublicationMetadataRequiresDistinctReviewer` |

## Risco residual (aceito, não escondido)

- **Sem sandbox de kernel**: um check `go_test`/`go_build` roda com os
  mesmos privilégios do processo `codinho` e do usuário local — código
  Go deliberadamente hostil (ex.: `os.RemoveAll("/")` num `TestMain`)
  tem o mesmo poder que teria rodando `go test` manualmente. Fora do
  modelo de uso suportado (Non-goal explícito desta spec e de
  safe-check-executor).
- **Confiança no toolchain `go` local**: `go vet`/`go build`/`gofmt` são
  binários do host, não sandboxed por este projeto.
- **Modo entrevista não garante integridade forte**: ver
  `docs/security/executor-limitations.md` e `internal/assessment`'s
  `integrityNote`.
- **Lock de workspace não sobrevive a kill -9**: `eventstore.AcquireLock`
  é um arquivo simples; um processo morto deixa o lock (documentado em
  `internal/eventstore/lock.go` como limite conhecido, fora do escopo
  desta spec).
