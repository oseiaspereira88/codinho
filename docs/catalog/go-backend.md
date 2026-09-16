---
title: Catálogo backend de Go
doc_type: reference
---

# Catálogo backend de Go

Este documento registra a composição editorial dos packs backend e a evidência
determinística disponível em 2026-09-16. Ele não declara publicação: a spec
[`go-backend-packs`](../../.pose/specs/2026-08-22-go-backend-packs.md) continua
`in-progress`, todos os itens estão em draft e revisão humana/playtest continuam
pendentes.

## Composição observada

| Pack | Atomic | Combined | Functional slice | Investigação (`debug`) | Systems mission | Conceitos | Competências | Nodes |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `go-concurrency-context` | 3 | 3 | 1 | 1 | 0 | 14 | 9 | 48 |
| `go-http` | 3 | 3 | 2 | 1 | 0 | 14 | 8 | 45 |
| `go-database` | 2 | 2 | 1 | 0 | 1 | 12 | 8 | 42 |
| **Backend** | **8** | **8** | **4** | **2** | **1** | **40** | **25** | **135** |

`debug` é contado como investigação. Os 23 desafios backend são canônicos;
nenhuma variante foi usada para atingir a meta de profundidade.

## Trilhas

- `go-concurrency-lifecycle`: concorrência, ownership de channels, cancelamento,
  backpressure e shutdown.
- `go-http-services`: JSON, limites, middleware, lifecycle, cliente e timeout.
- `go-backend`: composição dos 23 desafios em ordem de pré-requisitos,
  começando por concorrência, passando por HTTP e terminando em persistência.

O pack `go-database` hospeda a trilha integrada para que a montagem preserve
uma única identidade global de percurso; os dois packs anteriores mantêm suas
trilhas específicas.

## Evidência de validação

O catálogo foi validado na raiz com o binário local `./codinho` (`codinho dev
(none)`). A compilação via `go run ./cmd/codinho` ficou limitada por `disk
quota exceeded` durante a compilação temporária; isso não foi tratado como
aprovação.

| Comando | Resultado |
|---|---|
| `./codinho catalog validate --json` | passou sem diagnostics ou achados editoriais; 23 temas, 145 conceitos, 89 competências, 8 trilhas, 77 desafios e 501 nodes. |
| `./codinho catalog validate --checks --json` | 60/60 checks verificados; baselines negativos e referências positivas conforme declarado. |

Os checks usam somente fixtures locais embutidas no YAML. O pack de
concorrência inclui `go_test_race` para o cenário de leak; HTTP usa fixtures
`net/http` sem serviço remoto; database usa drivers sintéticos locais para
parametrização/rollback e cancelamento/cleanup.

## Semântica de publicação

| Estado | Quantidade observada |
|---|---:|
| Itens inventariados | 501 nodes em 77 desafios |
| Draft | todos os 23 itens backend e seus ancestrais |
| Published | 0 |
| Eligible | 0 |

Nenhum pack ou desafio recebeu `reviewed_by`, `playtested: true` ou
`publication.status: published`. Esses campos dependem de pessoa revisora,
playtest ponta a ponta e decisão explícita de publicação.

## Pendências

- Revisão pedagógica independente e playtest das três trilhas.
- Fixtures externas para uso manual via `workspace prepare`; o registro em
  [`testdata/packs/backend/`](../../testdata/packs/backend/) documenta que a
  fonte executável atual permanece no YAML.
- Auditoria final da distribuição no candidato integrado e decisão de
  fechamento da spec.

O inventário detalhado, checkpoints e limitações ficam em
[`go-backend-inventory.md`](../editorial/go-backend-inventory.md), e a fila
serializável está em [`go-backend-queue.json`](../editorial/go-backend-queue.json).
