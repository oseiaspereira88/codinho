---
title: Catálogo de produção e arquitetura em Go
doc_type: reference
---

# Catálogo de produção e arquitetura em Go

Este documento registra a composição observada dos packs de produção e
arquitetura após sua inclusão no manifest. Ele é um registro editorial, não uma
declaração de publicação: os dois packs e seus 13 desafios permanecem em
`draft`, e revisão humana, playtest e publicação ainda estão pendentes.
Consulte a [spec](../../.pose/specs/2026-08-22-go-production-architecture-packs.md),
o [inventário editorial](../editorial/go-production-architecture-inventory.md)
e a [fila de autoria](../editorial/go-production-architecture-queue.json) para
o planejamento que originou esta integração.

## Contagens observadas

A medição foi feita em 2026-09-16 com `go run ./cmd/codinho catalog validate
--json`, na raiz do repositório, depois de acrescentar os dois packs ao
manifest. O baseline anterior à integração é o snapshot registrado no
inventário editorial (`1fde57a`). Nodes são contados recursivamente nas árvores
`macro_steps`; a camada `layer` não é um node.

| Escopo | Temas | Conceitos | Competências | Trilhas | Desafios | Nodes |
|---|---:|---:|---:|---:|---:|---:|
| Baseline antes da integração | 23 | 145 | 89 | 8 | 77 | 501 |
| `go-production` | 4 | 19 | 11 | 1 | 9 | 84 |
| `go-architecture-debugging` | 2 | 8 | 7 | 1 | 4 | 28 |
| **Catálogo integrado** | **29** | **172** | **107** | **10** | **90** | **613** |

Os dois packs acrescentam exatamente seis temas, 27 conceitos, 18
competências, duas trilhas, 13 desafios e 112 nodes. A distribuição interna
dos desafios novos é `go-production`: 2 `atomic`, 2 `combined`, 2
`functional_slice` e 3 `systems_mission`; `go-architecture-debugging`: 2
`debug`, 1 `refactoring` e 1 `code_review`. Para a política de distribuição,
os quatro desafios de investigação, distribuídos entre `debug`, `refactoring`
e `code_review`, são tratados no grupo `investigation` conforme
`packs/distribution.json`.

## IDs e tipos dos desafios

| Pack | ID | Tipo | Check(s) declarados |
|---|---|---|---|
| `go-production` | `go-production.profile-hot-path` | `atomic` | `hot-path-benchmark` (`go_benchmark`) |
| `go-production` | `go-production.correlate-request-signals` | `atomic` | `correlation-tests` (`go_test`); `correlation-benchmark` (`go_benchmark`) |
| `go-production` | `go-production.observable-worker-shutdown` | `combined` | `shutdown-tests` (`go_test`) |
| `go-production` | `go-production.idempotent-retry-client` | `combined` | `retry-tests` (`go_test`) |
| `go-production` | `go-production.operable-message-worker` | `functional_slice` | `worker-race-tests` (`go_test_race`) |
| `go-production` | `go-production.instrumented-service-boundary` | `functional_slice` | `boundary-tests` (`go_test`) |
| `go-production` | `go-production.backpressure-redelivery-service` | `systems_mission` | `backpressure-race-tests` (`go_test_race`) |
| `go-production` | `go-production.compatible-schema-rollout` | `systems_mission` | `schema-rollout-tests` (`go_test`) |
| `go-production` | `go-production.partial-outage-response` | `systems_mission` | `outage-tests` (`go_test`) |
| `go-architecture-debugging` | `go-architecture-debugging.diagnose-shared-state-race` | `debug` | `shared-state-race` (`go_test_race`) |
| `go-architecture-debugging` | `go-architecture-debugging.trace-heap-retention` | `debug` | `heap-retention-tests` (`go_test`) |
| `go-architecture-debugging` | `go-architecture-debugging.simplify-overengineered-boundary` | `refactoring` | `boundary-refactoring-tests` (`go_test`) |
| `go-architecture-debugging` | `go-architecture-debugging.review-partial-outage-evidence` | `code_review` | `outage-evidence-tests` (`go_test`) |

Cada desafio declara uma fixture inicial e uma `reference_fixture` privada no
próprio pack. Nos checks dos dois packs, o baseline dos desafios de correção
tem resultado esperado `test_failure` e a referência `pass`; os dois
benchmarks esperam `pass` em ambos os lados. `test_failure` no baseline é a
falha intencional que o check deve reproduzir, não uma falha da validação.

## Relações e trilhas

As relações explícitas dos packs são `applies_in`: há 81 arestas em
`go-production` e 20 em `go-architecture-debugging`, totalizando 101 arestas
que ligam temas, conceitos e competências aos desafios. As dependências
sequenciais permanecem nos campos `prerequisites`/`requires` dos desafios; por
exemplo, o diagnóstico de corrida reutiliza
`go-production.correlate-request-signals` como pré-requisito. Não foi criado um
pré-requisito artificial apenas para conectar as trilhas.

As duas trilhas integradas preservam esta ordem:

- `go-debugging-performance`: `go-production.profile-hot-path` →
  `go-production.correlate-request-signals` →
  `go-architecture-debugging.diagnose-shared-state-race` →
  `go-architecture-debugging.trace-heap-retention` →
  `go-architecture-debugging.simplify-overengineered-boundary` →
  `go-architecture-debugging.review-partial-outage-evidence`.
- `go-architecture-production`: `go-production.profile-hot-path` →
  `go-production.correlate-request-signals` →
  `go-production.observable-worker-shutdown` →
  `go-production.idempotent-retry-client` →
  `go-production.operable-message-worker` →
  `go-production.instrumented-service-boundary` →
  `go-production.backpressure-redelivery-service` →
  `go-production.compatible-schema-rollout` →
  `go-production.partial-outage-response`.

Os temas de produção cobrem performance, observabilidade, resiliência e
segurança. O pack de arquitetura cobre investigação e evolução de design.
Mensageria e gRPC aparecem como contratos e doubles conceituais; nenhum broker
ou SDK gRPC externo é instalado.

## Checks e baseline local

Na execução integrada, `go run ./cmd/codinho catalog validate --json` passou
sem diagnostics ou achados editoriais. `go run ./cmd/codinho catalog validate
--checks --json` verificou 74/74 checks declarados no catálogo completo; os
dois packs novos contribuíram 14 declarações e 28 cenários baseline/referência,
todos verificados. O executor registrou `go1.27.1-X:nodwarf5`; os `go.mod` das
fixtures declaram Go 1.25.0, portanto comparações de desempenho devem manter a
toolchain efetiva constante.

Os checks `go_test_race` cobrem a corrida compartilhada, o worker e o
backpressure. Os checks `go_benchmark` cobrem o hot path e a correlação; eles
comprovam que a medição executa e emite `ns/op`/alocações, sem impor um limite
absoluto dependente de hardware. Para uma comparação manual pareada no mesmo
host, prepare as duas fixtures e repita cinco vezes:

```sh
go run ./cmd/codinho workspace prepare go-production.profile-hot-path --dest /tmp/codinho-profile-hot-path
(cd /tmp/codinho-profile-hot-path && go test -run '^$' -bench '^BenchmarkHotPath$' -benchmem -count=5 -benchtime=1s)
go run ./cmd/codinho workspace prepare go-production.correlate-request-signals --dest /tmp/codinho-request-correlation
(cd /tmp/codinho-request-correlation && go test -run '^$' -bench '^BenchmarkRequestCorrelation$' -benchmem -count=5 -benchtime=1s)
```

Registre carga, versão Go e `GOMAXPROCS` ao interpretar as medianas. Uma
diferença inferior a 15% em `ns/op` é inconclusiva; não compare hosts
distintos, não transforme a medição em gate global e não trate benchmark como
prova de correção funcional.

## Limites de segurança e hermeticidade

- Conforme o cenário, as fixtures usam dados inventados, biblioteca padrão,
  `net/http/httptest`, clocks/doubles injetados, filas e brokers sintéticos e
  `network: false`.
  Não acessam serviços, credenciais, DNS, rede externa, broker real ou SDK
  gRPC.
- Os cenários de fronteira verificam limites de entrada, identidade local e
  ausência de segredo sintético em logs. Isso demonstra o comportamento da
  fixture; não é certificação de autenticação, TLS, gestão de segredos,
  disponibilidade ou supply chain em produção.
- Retenção de heap é avaliada por ownership, referências e contagens
  explícitas. RSS, timing de GC e `sleep` não são oráculos. Os resultados de
  race e dos doubles também não substituem uma revisão de implantação real.

## Estado editorial e pendências

O pack `go-production`, o pack `go-architecture-debugging` e cada um dos 13
desafios listados acima estão com `publication.status: draft` e autoria
editorial `codex`. A projeção observada confirma que todo o inventário está em
rascunho:

| Projeção | Temas | Conceitos | Competências | Trilhas | Desafios | Nodes |
|---|---:|---:|---:|---:|---:|---:|
| Inventário | 29 | 172 | 107 | 10 | 90 | 613 |
| Draft | 29 | 172 | 107 | 10 | 90 | 613 |
| Published | 0 | 0 | 0 | 0 | 0 | 0 |
| Eligible | 0 | 0 | 0 | 0 | 0 | 0 |

`reviewed_by` não foi preenchido e `playtested` não foi marcado como `true`.
Continuam pendentes a revisão pedagógica independente, o playtest humano
ponta a ponta e a decisão explícita de publicação. Nenhum desses gates é
substituído pelos checks automatizados ou pelas contagens observadas.
