---
title: Inventário editorial dos packs backend de Go
doc_type: reference
---

# Inventário editorial dos packs backend de Go

Este inventário governa a autoria da [spec de backend](../../.pose/specs/2026-08-22-go-backend-packs.md)
e a [fila estruturada](go-backend-queue.json). Ele separa metas de autoria,
contagem de catálogo e gates humanos; não declara publicação.

## Baseline e alvo

Baseline medido em 2026-09-16 com `./codinho catalog validate --json`, no
commit `4067f89`: 105 conceitos, 64 competências, 54 entradas de desafio,
366 nodes e cinco trilhas. Tudo permanece em draft.

Checkpoint 1 medido após a materialização de `go-concurrency-context`:
119 conceitos, 73 competências, 62 entradas de desafio, 414 nodes e seis
trilhas. O acréscimo observado é de 14 conceitos, nove competências, oito
desafios e 48 nodes; a contagem ainda inclui somente drafts.

Checkpoint 2 medido após a materialização de `go-http`:
133 conceitos, 81 competências, 71 entradas de desafio, 459 nodes e sete
trilhas. O pack HTTP acrescentou 14 conceitos, oito competências, nove
desafios e 45 nodes; seus 58 checks declarados foram verificados, sem
diagnostics ou achados editoriais. A contagem ainda inclui somente drafts.

O backend adicionará 23 desafios canônicos, pelo menos 40 conceitos, 25
competências e 120 nodes creditáveis. Como os fundamentos já excedem a meta
intermediária de 140 conceitos, 85 competências e 420 nodes quando o backend
começar, a contagem será acompanhada sem inflar a árvore para atingir um
número.

| Pack | Atômicos | Combinados | Fatias | Investigação | Missão | Conceitos | Competências | Nodes mínimos | Trilha própria |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `go-concurrency-context` | 3 | 3 | 1 | 1 | 0 | 14 | 9 | 42 | concorrência e lifecycle |
| `go-http` | 3 | 3 | 2 | 1 | 0 | 14 | 8 | 40 | HTTP e serviços |
| `go-database` | 2 | 2 | 1 | 0 | 1 | 12 | 8 | 40 | Go backend |
| **Total** | **8** | **8** | **4** | **2** | **1** | **40** | **25** | **122** | **3** |

`debug` é o tipo de catálogo usado para investigação e conta no grupo
`investigation` da política de distribuição. Variantes não serão usadas para
substituir desafios canônicos ou profundidade das bases.

## Matriz de desafios

### `go-concurrency-context`

- Atômicos: `worker-lifecycle-with-waitgroup`, `channel-ownership-and-close`,
  `cancel-select-loop`.
- Combinados: `bounded-worker-pool`, `pipeline-with-backpressure`,
  `shutdown-on-first-error`.
- Fatia: `concurrent-log-aggregator`.
- Investigação: `reproduce-and-fix-goroutine-leak`.

### `go-http`

- Atômicos: `parse-json-request`, `limit-http-body`,
  `propagate-request-timeout`.
- Combinados: `compose-json-handler-middleware`,
  `configure-http-server-lifecycle`, `build-resilient-http-client`.
- Fatias: `tested-json-handler`, `resilient-service-client`.
- Investigação: `reproduce-http-timeout-leak`.

### `go-database`

- Atômicos: `parameterize-query-and-null`, `bound-query-to-context`.
- Combinados: `transaction-with-isolation`, `pool-and-migration-lifecycle`.
- Fatia: `transactional-repository`.
- Missão: `graceful-service-with-persistence`.

## Critérios de evidência

Cada desafio deve permanecer em draft e conter briefing sem solução,
competência primária, aceitação observável, hints progressivos, reflexão e
fixtures/referência locais quando declarar checks. Concorrência deve incluir
`go_test_race` nos desafios pertinentes; HTTP e SQL devem provar limites,
cancelamento, parametrização e ausência de dependência remota.

O binário local `./codinho` é a rota de validação editorial desta execução.
`go run ./cmd/codinho catalog validate` ficou limitado por `disk quota
exceeded` durante recompilação temporária; essa limitação não conta como
aprovação. A validação final deve registrar comando, versão do binário e
resultado estruturado.

## Gates mantidos fora da automação

- Revisão pedagógica por pessoa diferente do autor.
- Playtest ponta a ponta dos desafios e das três trilhas.
- Preenchimento de `reviewed_by`, `playtested` ou `publication.status`.
- Decisão de publicação e fechamento da spec.
