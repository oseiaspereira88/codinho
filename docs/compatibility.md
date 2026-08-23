# Matriz de compatibilidade

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

`codinho serve` carrega o catálogo **uma única vez**, na inicialização
(`cmd/codinho/main.go:runServe`) — não há hot-reload. Uma sessão em
andamento referencia esse mesmo ponteiro de catálogo durante toda a
sua vida útil; editar `packs/*.yaml` em disco enquanto `codinho serve`
está rodando não afeta nenhuma sessão ativa (nem cria inconsistência),
porque não existe nenhum caminho de código que releia o catálogo depois
do startup. Para que uma edição de pack tenha efeito, reinicie
`codinho serve` — e como sessões já são apenas em memória (não
sobrevivem a reinício de processo, ver session-orchestration-disclosure
Decision 3), isso nunca mistura conteúdo antigo e novo na mesma sessão.

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
