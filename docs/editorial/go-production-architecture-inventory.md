---
title: Inventário editorial de produção e arquitetura em Go
doc_type: reference
---

# Inventário editorial de produção e arquitetura em Go

Use este inventário para executar a [spec](../../.pose/specs/2026-08-22-go-production-architecture-packs.md)
e a [fila de autoria](go-production-architecture-queue.json). Ele define a
distribuição e os limites técnicos; não é evidência de qualidade pedagógica nem
de publicação.

## Baseline e alvo

Baseline em 2026-09-16, HEAD `1fde57a`, medido com
`./codinho catalog validate --json`: 145 conceitos, 89 competências, 77
desafios, 501 nodes e oito trilhas. Todos os itens permanecem em draft.

O lote acrescenta exatamente 13 desafios: dois atômicos, dois combinados, duas
fatias funcionais, quatro investigações e três missões de sistemas. A política
`packs/distribution.json` reúne `debug`, `refactoring` e `code_review` no grupo
`investigation`. O plano acrescenta 27 conceitos, 18 competências e prevê 112
nodes contextualizados; a soma real será medida depois da autoria. Não acrescente
passos sem intenção própria para alcançar contagens.

## Matriz dos desafios

| Pack | ID | Tipo | Foco verificável | Evidência principal | Nodes planejados |
|---|---|---|---|---|---:|
| `go-production` | `go-production.profile-hot-path` | atomic | Perfil local e comparação de uma transformação determinística | `go_benchmark` | 7 |
| `go-production` | `go-production.correlate-request-signals` | atomic | Propagar correlação e manter cardinalidade limitada | `go_test` + `go_benchmark` | 7 |
| `go-production` | `go-production.observable-worker-shutdown` | combined | Retirar readiness, drenar e respeitar prazo de shutdown | `go_test` | 8 |
| `go-production` | `go-production.idempotent-retry-client` | combined | Prazo RPC conceitual, retry limitado e chave idempotente | `go_test` | 8 |
| `go-production` | `go-production.operable-message-worker` | functional_slice | Worker local com fila limitada, ack explícito e sinais operacionais | `go_test_race` | 9 |
| `go-production` | `go-production.instrumented-service-boundary` | functional_slice | Limites de entrada, identidade local, correlação e logs sem segredo | `go_test` | 9 |
| `go-production` | `go-production.backpressure-redelivery-service` | systems_mission | Broker sintético, pressão, redelivery e requisito incompleto de disponibilidade | `go_test_race` | 12 |
| `go-production` | `go-production.compatible-schema-rollout` | systems_mission | Leitores mistos, rollout expand/migrate/contract e rollback | `go_test` | 12 |
| `go-production` | `go-production.partial-outage-response` | systems_mission | Falha parcial, orçamento de retries, proveniência e mitigação reversível | `go_test` | 12 |
| `go-architecture-debugging` | `go-architecture-debugging.diagnose-shared-state-race` | debug | Reproduzir e localizar corrida sem sincronizar com `sleep` | `go_test_race` | 7 |
| `go-architecture-debugging` | `go-architecture-debugging.trace-heap-retention` | debug | Seguir ownership/referências retidas sem depender de RSS ou GC timing | `go_test` | 7 |
| `go-architecture-debugging` | `go-architecture-debugging.simplify-overengineered-boundary` | refactoring | Remover abstração sem pressão real e preservar o contrato | `go_test` | 7 |
| `go-architecture-debugging` | `go-architecture-debugging.review-partial-outage-evidence` | code_review | Separar correlação de causa e identificar a mitigação sustentada pelos dados | `go_test` | 7 |
| **Total** |  | 2 atomic, 2 combined, 2 functional, 4 investigation, 3 mission |  |  | **112** |

As missões apresentam uma exigência incompleta, pelo menos dois modos de falha
e alternativas com custos explícitos. Cada uma pede justificativa de custo,
falha, operabilidade e alternativa; a rubrica aparece nos critérios/reflexões do
desafio, não em um campo de schema inventado.

## Conceitos e competências reservados

`go-production` define estes 19 conceitos:

- `local-workload-baseline`, `cpu-and-heap-profiles`, `benchmark-noise`;
- `trace-context`, `structured-event-correlation`, `cardinality-budget`;
- `readiness-vs-liveness`, `shutdown-budget`, `idempotency-key-scope`,
  `retry-budget`, `grpc-deadline-budget`, `queue-saturation`,
  `redelivery-semantics`, `expand-migrate-contract`, `failure-budget`;
- `service-identity`, `secrets-boundary`, `release-artifact-provenance`,
  `input-size-boundary`.

`go-architecture-debugging` define estes oito conceitos:

- `race-happens-before`, `reachable-heap-retention`, `demand-shaped-boundary`,
  `symptom-vs-cause`, `failure-domain`, `incident-timeline`, `causal-evidence`,
  `reversible-mitigation`.

As 11 competências de `go-production` são `profile-hot-path-with-local-baseline`,
`interpret-cpu-and-allocation-evidence`, `correlate-request-signals`,
`bound-metric-cardinality`, `separate-readiness-from-liveness`,
`drain-workers-with-deadline`, `make-retries-idempotent-and-bounded`,
`bound-queue-and-redelivery`, `roll-schema-compatibly`,
`protect-service-identity-and-secrets` e `verify-release-provenance`.

As sete competências de `go-architecture-debugging` são `diagnose-data-race`,
`trace-retained-object-lifecycle`, `remove-unsupported-abstraction`,
`reconstruct-incident-evidence`, `select-reversible-mitigation`,
`distinguish-correlation-from-causation` e `justify-trade-off-with-evidence`.

IDs de conceito e competência são namespaced pelo ID do pack e associados às
atividades por relações `applies_in`. Não se criam variantes para alcançar metas.

Temas novos, todos namespaced e sem colisões no baseline:
`go-production.performance`, `go-production.observability`,
`go-production.resilience`, `go-production.security`,
`go-architecture-debugging.investigation` e
`go-architecture-debugging.design-evolution`.

Mensageria e gRPC são abordados conceitualmente com doubles locais: ack,
redelivery, deadlines e limites de identidade; nenhum broker ou SDK gRPC é
instalado.

## Trilhas

- `go-debugging-performance`: `go-production.profile-hot-path`,
  `go-production.correlate-request-signals`,
  `go-architecture-debugging.diagnose-shared-state-race`,
  `go-architecture-debugging.trace-heap-retention`,
  `go-architecture-debugging.simplify-overengineered-boundary` e
  `go-architecture-debugging.review-partial-outage-evidence`.
- `go-architecture-production`: `go-production.profile-hot-path`,
  `go-production.correlate-request-signals`,
  `go-production.observable-worker-shutdown`,
  `go-production.idempotent-retry-client`,
  `go-production.operable-message-worker`,
  `go-production.instrumented-service-boundary`,
  `go-production.backpressure-redelivery-service`,
  `go-production.compatible-schema-rollout` e
  `go-production.partial-outage-response`.

Cada desafio pertence a uma trilha. A ordem completa e os IDs namespaced estão
nas definições dos packs; não se declaram pré-requisitos de outro pack apenas
para tornar uma trilha executável.

## Fixtures, segurança e benchmarks

- Use somente `net/http/httptest`, clocks/doubles injetados, filas/brokers
  sintéticos e dados de incidente inventados. Não acesse serviços, credenciais,
  DNS ou rede externa; mantenha `checks[].network` ausente ou `false`.
- Use TLS apenas contra servidor `httptest` com certificado local. Teste
  identidade, autenticação, limites de entrada e ausência de segredo em logs
  com valores sintéticos.
- Inclua baseline e referência por check. Cada um dos 13 desafios precisa de
  pelo menos um check declarativo executável; checks de race usam
  `go_test_race`. O check de benchmark comprova execução e emissão de medida,
  não estabelece um limite absoluto por hardware.
- Para comparar desempenho manualmente no mesmo host/toolchain, materialize as
  duas fixtures com `workspace prepare` e execute `go test -run '^$' -bench
  '<benchmark>' -benchmem -count=5 -benchtime=1s`. Compare medianas após manter
  carga, versão Go e `GOMAXPROCS` iguais. Mudanças abaixo de 15% em `ns/op` são
  inconclusivas; não transforme essa tolerância em gate global ou comparação
  entre máquinas.
- Não use RSS, pausas de GC ou sleeps como oráculos determinísticos para
  retenção/shutdown. Verifique ownership, contagem de trabalho, referências e
  limites explícitos.

Exemplos de comandos locais copy-pasteable após a autoria:

```sh
./codinho workspace prepare go-production.profile-hot-path --dest /tmp/codinho-profile-hot-path
(cd /tmp/codinho-profile-hot-path && go test -run '^$' -bench '^BenchmarkHotPath$' -benchmem -count=5 -benchtime=1s)
./codinho workspace prepare go-production.correlate-request-signals --dest /tmp/codinho-request-correlation
(cd /tmp/codinho-request-correlation && go test -run '^$' -bench '^BenchmarkRequestCorrelation$' -benchmem -count=5 -benchtime=1s)
```

## Gates humanos

Todo pack/desafio permanece draft, com autor editorial honesto. A revisão
pedagógica independente, os playtests ponta a ponta, qualquer preenchimento de
`reviewed_by`/`playtested` e a publicação permanecem pendentes. Métricas e
checks determinísticos não os substituem.
