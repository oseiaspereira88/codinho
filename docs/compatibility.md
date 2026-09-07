# Matriz de compatibilidade

Ver [`docs/install.md`](install.md) para pré-requisitos de instalação e
[`.github/workflows/ci.yml`](../.github/workflows/ci.yml) para os gates
que rodam a cada pull request contra esta matriz (requirement R4/R5).

## Go

- **Mínimo declarado**: `go 1.25.0` (`go.mod`).
- **Testado nesta sessão**: `go1.26.6 linux/amd64` — `go build ./...`,
  `go vet ./...` e `go test ./... -race` passam limpos.
- `GOOS=windows go vet ./...` roda sem diagnósticos a cada spec fechada
  desde o início do projeto (checagem estática cross-compile; não é
  execução real em Windows).

## Sistema operacional

| OS | Status | Evidência |
|---|---|---|
| Linux | **Suportado, testado** | Toda a suíte de testes (`go test ./... -race`) roda nesta plataforma a cada spec. |
| macOS | **Suportado, não testado nesta sessão** | Código não usa nenhuma primitiva específica de Linux (paths, syscalls e process-group handling em `internal/checks/process_unix.go` usam a API POSIX padrão do pacote `syscall`, compatível com Darwin). Sem execução real registrada. |
| Windows | **Não declarado suportado até haver evidência real** | `GOOS=windows go vet` passa (checagem estática), e `internal/checks/process_windows.go`/`internal/diagnostics/process_windows.go` implementam os pontos de divergência de plataforma (grupo de processo, verificação de PID). Sem execução real registrada — não reivindicar suporte além do que a Compatibility desta spec permite ("documentar Windows como suportado somente após evidência"). |

## SDK MCP

- `github.com/modelcontextprotocol/go-sdk` — versão fixada em `go.mod`
  (`v1.7.0`). Atualizações passam pela suíte de contrato completa
  (`internal/mcpserver/*_contract_test.go`) antes de merge.

## Schema de catálogo e de evento

- Catálogo: `schema_version: 1` (`internal/curriculum.SchemaVersion`).
  Mudança de versão exige uma spec própria (não coberta por packs
  automaticamente).
- Evento: `schema_version: 1` (`internal/eventstore.CurrentEventSchemaVersion`).
  Compatibilidade futura é via `UpcasterRegistry` — um evento mais antigo
  é migrado evento a evento durante o replay; um evento mais novo sem
  upcaster registrado falha explicitamente
  (`ErrIncompatibleEventVersion`), nunca é silenciosamente ignorado.

## Atualização de pack sem quebrar sessão fixada

Reinicie `codinho serve` para aplicar edições do catálogo a sessões novas;
não há hot-reload. Sessões iniciadas com `recovery_version: 1` preservam
política completa e cópia do desafio com SHA-256 no evento de início.
O replay restaura estado, nó, pistas, detours, avaliação e revisão sem
reconsultar o desafio atual. Remover o desafio de um catálogo válido não
remove o conteúdo fixado de sessões existentes.

Repita uma mutação confirmada com o mesmo `request_id` e entrada original
para recuperar seu resultado, mesmo após reinício/progresso posterior.
Reutilização incompatível retorna `STATE_CONFLICT`. IDs antigos permanecem
reservados. Sessões legadas sem política/conteúdo recuperável retornam
`SESSION_RECOVERY_UNAVAILABLE`; preserve o histórico e inicie outra sessão.

Baselines, roots canônicas, globs e vínculo sessão–evidência sobrevivem ao
reinício. `evidence_get` verifica atualidade no escopo persistido; root
indisponível/substituída nunca demonstra freshness. Não troque root/globs
de uma baseline existente. Isto **não** valida o consumo de evidência em
`step_evaluate`: esse gate está planejado em
[evaluation-evidence-lineage](../.pose/specs/2026-09-07-evaluation-evidence-lineage.md)
e bloqueia o aceite V1.

Use `codinho doctor` após interrupção abrupta. `SIGKILL` pode deixar lock
órfão: confirme que nenhum servidor usa o workspace antes de removê-lo,
conforme o diagnóstico. O writer com lock preserva bytes de cauda incompleta
em `events.jsonl.recovery-*` (permissão 0600) antes de truncar somente essa
cauda. Corrupção em linha terminada ou revisão inválida bloqueia startup;
inspeções da CLI são somente leitura e nunca reparam o log. Inclua backups
de recuperação na política de retenção/exclusão do estado.

Mantenha manifest e catálogo válidos para iniciar o servidor. O mecanismo
não migra payloads futuros nem corrige conteúdo corrompido por inferência.
Rollback para binário anterior perde projeção de sessão e proteção contra
colisões; não o execute como writer sobre estado novo sem estratégia de
migração. Consulte o [ADR](../.pose/adr/2026-09-06-durable-session-replay-with-pinned-content.md).

Evidência de 2026-09-07: `go test ./cmd/codinho -run TestSessionRecoveryOverRealStdio`
(update, remove, resposta perdida/cauda interrompida e legado/corrupção),
além das regressões em `internal/session`, `internal/application` e `internal/eventstore`.

## Limites de recursos aplicados

| Recurso | Limite | Onde |
|---|---|---|
| Tamanho de arquivo YAML de pack | 1 MiB | `curriculum.DefaultLimits.MaxFileBytes` |
| Profundidade de nó YAML | 32 | `curriculum.DefaultLimits.MaxNodeDepth` |
| Aliases YAML | 64 | `curriculum.DefaultLimits.MaxAliases` |
| Saída de check (stdout/stderr) | 256 KiB, com sufixo de truncamento | `internal/checks/limits.go` |
| Checks concorrentes | 4 | `internal/checks/limits.go` (`maxParallel`) |
| Timeout de check | por check, validado por `checks.Resolve` | `internal/checks/registry.go` |
| Resultado de busca no catálogo | 100 itens, com `Truncated: true` explícito | `internal/curriculum/search.go` |

## Desempenho (informativo, não um SLO de serviço cloud)

- Busca no catálogo em escala V1 (84 desafios sintéticos): p95 bem
  abaixo de 100ms — `TestSearchP95BudgetAtV1Scale`.
- Startup do `codinho serve` até a primeira resposta MCP real: bem
  abaixo de 1s — `TestServeStartupRespondsWithinOneSecond`.

Esses números são medidos no hardware desta sessão de desenvolvimento;
são thresholds informativos, não uma garantia de SLO (Non-functional:
"Benchmarks são informativos com thresholds documentados").
