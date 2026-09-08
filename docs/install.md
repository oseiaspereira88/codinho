---
title: "Instalação"
doc_type: howto
---

# Instalação

## Pré-requisitos

- Go 1.25 ou mais recente (ver [`docs/compatibility.md`](compatibility.md)).
- `git` (usado por `workspace_observe`/`check_run` quando o desafio está
  num repositório Git; opcional caso contrário).
- Um host MCP: [Codex CLI](https://github.com/openai/codex) ou uma
  extensão de IDE com suporte a servidores MCP locais.

Nenhuma outra dependência externa é necessária — o `codinho` não tem
dependências de rede em tempo de execução.

## Build

```sh
git clone <este repositório>
cd codinho
go build -o codinho ./cmd/codinho
```

Isso produz um único binário `codinho` (comando administrativo e
servidor MCP, ver `cmd/codinho/main.go`). Coloque-o no seu `PATH`, ou
referencie o caminho completo na configuração do host.

## Instalação

Não há instalador nem gerenciador de pacotes na V1 (Constraint:
"Distribuição V1 é binário, skill e configuração local; não é
plugin/marketplace/auto-update"). Os três passos são:

1. **Binário**: `go build -o codinho ./cmd/codinho` (acima), ou baixe um
   binário já compilado da release, se disponível para sua plataforma
   (ver [`docs/compatibility.md`](compatibility.md) para o que já foi
   testado).
2. **Skill**: aponte seu host para `.agents/skills/codinho/` (o formato
   é Agent Skills padrão — a maioria dos hosts a descobre automaticamente
   se o repositório estiver no workspace, ou você pode copiá-la para o
   diretório de skills do seu host).
3. **Configuração MCP**: ver [`docs/configuration.md`](configuration.md).

## Verificação

```sh
./codinho doctor
```

Deve reportar `status: ok` (ou apontar exatamente o que está faltando —
catálogo, toolchain, permissões de `.codinho/state`).

## Atualização

Substitua o binário (`go build -o codinho ./cmd/codinho` novamente após
`git pull`) e reinicie qualquer `codinho serve` em execução — o
catálogo é carregado uma única vez na inicialização (ver
[`docs/compatibility.md`](compatibility.md#atualização-de-pack-sem-quebrar-sessão-fixada)),
então uma sessão MCP ativa não é afetada até o próximo restart.

## Remoção

```sh
./codinho privacy purge --confirm  # estado local (.codinho/state)
rm codinho                        # remova o binário por último
```

Isso nunca toca seu código nem o repositório Git do desafio — ver
[`docs/security/privacy.md`](security/privacy.md).
