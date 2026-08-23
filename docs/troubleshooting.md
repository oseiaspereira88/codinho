# Troubleshooting

Primeiro passo sempre: `codinho doctor` (ou `codinho doctor --json` para
um relatório estruturado). Ele nunca imprime código, ambiente completo
ou segredos — só o necessário para diagnosticar.

Instalação do zero: [`docs/install.md`](install.md). Configuração de
host: [`docs/configuration.md`](configuration.md).

## `lock: orphaned`

**Sintoma**: `codinho doctor` reporta `lock: error (orphaned: the
process that held this lock is no longer running...)`.

**Causa**: um `codinho serve` anterior foi encerrado abruptamente
(`kill -9`, queda de energia, etc.) sem liberar
`.codinho/state/lock` — arquivos de lock não sobrevivem a esse tipo de
encerramento (limite conhecido e documentado de
`internal/eventstore.Lock`).

**Resolução**: confirme que nenhum `codinho serve` real está rodando
(`ps aux | grep codinho`), então remova o arquivo manualmente:

```sh
rm .codinho/state/lock
```

`codinho doctor` nunca remove esse arquivo sozinho — é uma decisão do
operador, nunca automática.

## `catalog: error`

**Sintoma**: `codinho doctor` ou `codinho catalog validate` reporta
erro de catálogo.

**Resolução**: rode `codinho catalog validate` (sem `--json` primeiro,
para ler as mensagens) — cada diagnóstico aponta arquivo, item, campo e
regra. Ver `docs/content-authoring.md` para o significado de cada regra.

## `state_writable: error`

**Sintoma**: `.codinho/state` não pode ser criado ou não é gravável.

**Resolução**: verifique permissões do diretório do workspace; o
processo `codinho` precisa poder criar `.codinho/state` com permissão
`0700`. Rodar como um usuário sem permissão de escrita no diretório do
projeto é a causa mais comum.

## `toolchain:go` ou `toolchain:gofmt` ausente

**Sintoma**: `codinho doctor` reporta que `go` ou `gofmt` não está no
PATH.

**Resolução**: `check_run` para checks `go_test`/`go_vet`/`gofmt_check`
etc. precisa desses binários no PATH do processo `codinho serve`. Se
você usa uma versão de Go instalada fora do PATH padrão do sistema,
inicie `codinho serve` a partir de um shell que já tenha o Go
configurado.

## `codinho serve` não responde / trava

- Confirme que nada escreve para stdout além do próprio processo MCP —
  `stdout` é reservado exclusivamente para o protocolo
  (`TestServeStdoutIsOnlyValidJSONFrames` prova isso no CI). Se você
  modificou o código e adicionou um `fmt.Println` por engano em
  qualquer caminho alcançável por `codinho serve`, isso quebra o
  protocolo silenciosamente do lado do host.
- Verifique se `.codinho/state/lock` já existe (outro `codinho serve` já
  está rodando neste workspace) — `codinho doctor` detecta isso.

## Uma sessão "desapareceu"

Sessões existem **somente na memória** do processo `codinho serve` que
as criou (session-orchestration-disclosure Decision 3) — não
sobrevivem a reinício do processo. Isso é esperado, não um bug: é o
mesmo modelo de qualquer servidor MCP stdio de vida curta, atado à
conversa do host. Progresso de competência (`progress_get`,
`review_due`) é durável e sobrevive normalmente, porque vive no event
log compartilhado, não na sessão em si.

## Quero apagar/exportar meu progresso local

Ver `docs/security/privacy.md` — `codinho privacy export --dest <path>`
e `codinho privacy purge --confirm`.
