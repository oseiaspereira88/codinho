---
title: "codinho"
doc_type: reference
---

# codinho

## Microaprendizado de código assistido por IA

O **codinho** é um tutor de prática deliberada para desenvolver e
recuperar fluência em Go, sem deixar a IA substituir o raciocínio, a
decomposição e a escrita de código do próprio aluno. Ele roda
inteiramente na sua máquina: um servidor MCP local em Go, uma skill que
ensina um agente (Codex CLI, extensão de IDE) a orquestrá-lo, e uma CLI
administrativa.

## Início rápido

```sh
go build -o codinho ./cmd/codinho
./codinho init
./codinho doctor
```

Depois disso, siga [`docs/quickstart.md`](docs/quickstart.md) para
conectar o servidor MCP ao seu host (Codex CLI ou extensão de IDE) e
começar seu primeiro desafio.

## Documentação

| Documento | Conteúdo |
|---|---|
| [`docs/install.md`](docs/install.md) | Pré-requisitos, build, instalação, atualização e remoção. |
| [`docs/quickstart.md`](docs/quickstart.md) | Do zero ao primeiro desafio completado. |
| [`docs/configuration.md`](docs/configuration.md) | Configuração do servidor MCP no Codex CLI e na extensão de IDE. |
| [`docs/content-authoring.md`](docs/content-authoring.md) | Como autorar packs de conteúdo curricular. |
| [`docs/content-review-checklist.md`](docs/content-review-checklist.md) | Checklist de revisão antes de publicar um desafio. |
| [`docs/compatibility.md`](docs/compatibility.md) | Matriz de Go, SO, SDK MCP e schemas suportados. |
| [`docs/troubleshooting.md`](docs/troubleshooting.md) | Diagnóstico de problemas comuns (`codinho doctor` primeiro). |
| [`docs/security/threat-model.md`](docs/security/threat-model.md) | Modelo de ameaças, trust boundaries e limites honestos. |
| [`docs/security/privacy.md`](docs/security/privacy.md) | O que é coletado, retenção, export e remoção. |

## Componentes

- **Servidor MCP** (`codinho serve`) — currículo, sessões, evidências e
  progresso, tudo local (`internal/mcpserver`, `internal/session`,
  `internal/eventstore`).
- **Skill do tutor** (`.agents/skills/codinho/`) — ensina o agente do
  host a rotear as tools MCP de forma pedagogicamente consistente.
- **CLI administrativa** (`codinho`) — `init`, `catalog`, `session`,
  `progress`, `workspace`, `privacy`, `doctor` (ver
  [`docs/install.md`](docs/install.md)).
- **Catálogo curricular** (`packs/*.yaml`) — temas, conceitos,
  competências, trilhas e desafios versionados.

## Desenvolvimento

O repositório usa POSE para governar o trabalho de agentes e a
evolução das especificações. Leia [`AGENTS.md`](AGENTS.md) antes de
alterar o projeto e use [`POSE.md`](POSE.md) como manual operacional.

- [`PROJECT.md`](PROJECT.md) é a referência principal da visão da V1.
- [`.pose/specs/`](.pose/specs/) contém as specs incrementais de implementação.
- [`.pose/roadmaps/codinho-v1.md`](.pose/roadmaps/codinho-v1.md) é a sequência governada de entrega.
- [`.pose/changelogs/unreleased/`](.pose/changelogs/unreleased/) acumula o changelog até o corte de release.

```sh
make build   # go build ./cmd/codinho
make test    # go test ./... -race
make lint    # gofmt -l . && go vet ./...
make check   # tudo acima + govulncheck + catalog validate
```

Ver [`Makefile`](Makefile) para a lista completa de alvos.
