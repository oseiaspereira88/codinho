# Configuração do servidor MCP

O `codinho` expõe seu servidor MCP via `codinho serve`, sobre stdio —
o mesmo transporte que Codex CLI e extensões de IDE já sabem usar.
`stdout` é reservado exclusivamente para o protocolo (nunca escreva
nada em `stdout` se você estender este código — ver
[`docs/security/executor-limitations.md`](security/executor-limitations.md)
e `docs/troubleshooting.md`).

## Codex CLI

Ver [`.agents/skills/codinho/references/codex-configuration.md`](../.agents/skills/codinho/references/codex-configuration.md)
para o bloco de configuração completo (`[mcp_servers.codinho]`, sem
paths privados).

## Extensão de IDE

Qualquer extensão com suporte a servidores MCP locais via stdio aceita
a mesma forma geral:

```json
{
  "mcpServers": {
    "codinho": {
      "command": "codinho",
      "args": ["serve"]
    }
  }
}
```

Ajuste `command` para o caminho completo do binário se `codinho` não
estiver no `PATH` da extensão. `codinho serve` deve ser iniciado a
partir do diretório do projeto que contém `packs/` (é onde ele procura
o catálogo e onde cria `.codinho/state/`).

## Skill

Aponte seu host para `.agents/skills/codinho/` (formato Agent Skills
padrão). Ver `.agents/skills/codinho/SKILL.md` para o contrato completo
e `references/mcp-tool-routing.md` para a lista de tools.

## Variáveis de ambiente

Nenhuma é necessária. O `codinho` não lê configuração de ambiente além
do `PATH` (para localizar `go`/`gofmt`, usados por `check_run`) e não
faz nenhuma chamada de rede por padrão (ver
[`docs/security/privacy.md`](security/privacy.md)).
