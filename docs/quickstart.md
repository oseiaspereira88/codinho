# Quickstart

Do zero ao primeiro desafio completado.

## 1. Build e inicialização

```sh
go build -o codinho ./cmd/codinho
./codinho init          # cria packs/manifest.yaml vazio no diretório atual
```

Se você já está neste repositório (que já tem `packs/` autorado), pule o
`init` — ele existe para começar um catálogo próprio do zero.

## 2. Verifique a instalação

```sh
./codinho doctor
```

## 3. Configure seu host

Ver [`docs/configuration.md`](configuration.md) para Codex CLI e
extensão de IDE. Depois de configurado, seu agente já enxerga as tools
MCP do `codinho`.

## 4. Escolha um desafio

Fora de uma sessão, peça ao agente para buscar algo — ele vai chamar
`catalog_search`. Via CLI, para conferir manualmente:

```sh
./codinho catalog list
./codinho catalog show <challenge-id>
```

## 5. Se o desafio tiver código inicial (`kind: debug`)

```sh
./codinho workspace prepare <challenge-id> --dest ./meu-exercicio
```

Isso materializa a fixture (código com o bug, por exemplo) no destino
que você escolher — nunca dentro deste repositório, nunca escrito pelo
próprio agente.

## 6. Pratique

A partir daqui, é conversa com o agente: ele chama `session_start`,
anuncia o modo e a profundidade, e guia você passo a passo (ver a skill
em `.agents/skills/codinho/SKILL.md` para o contrato completo de
comportamento). Você escreve o código; o agente observa
(`workspace_observe`), avalia sob pedido (`step_evaluate`) e nunca edita
nada por você.

## 7. Acompanhe seu progresso

```sh
./codinho progress show
./codinho progress export --out progresso.json
```

## Encerrando

```sh
./codinho privacy export --dest ~/backups/codinho   # opcional, antes de apagar
./codinho privacy purge --confirm                    # remove o estado local
```
