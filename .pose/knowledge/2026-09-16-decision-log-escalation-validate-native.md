---
type: decision-log
slug: escalation-validate-native
owner: @oseiaspereira
sensitivity: public-internal
created_at: 2026-09-16
last_reviewed_at: 2026-09-16
expires_at: 2026-12-15
source_refs:
  spec: editorial-batch-engine
  workflow: validate-native-evidence
  commands:
    - "pose recurrence-check --tolerant --window-days 30 --threshold 3"
    - "pose stats workflows --since-days 30"
    - "pose stats tasks --since-days 30 --json"
    - "python3 -m unittest discover -s scripts -p test_editorial.py"
    - "pose check --strict"
    - "pose validate --strict"
  external_sources: []  # [{url: "", accessed_at: "YYYY-MM-DD"}]
---

# decision-log: escalation-validate-native

## Context

`pose recurrence-check --tolerant --window-days 30 --threshold 3` sinalizou
`validate-native` com seis falhas em trinta dias. A agregação registrou 16
passagens e seis falhas em 22 execuções. O histórico recente mistura
execuções estritas e tolerantes, resultados derivados e relatórios manuais sem
comandos ou saída capturada.

As falhas reproduzidas nesta execução incluem escrita de cache de módulos em
um caminho somente leitura durante o build e falha de DNS ao consultar
`govulncheck`. Portanto, o sinal comprova uma fragilidade de evidência e de
ambiente, mas não comprova regressão no mecanismo editorial.

## Current state

Classificação: domínio de processo/validação. O workflow
`.pose/workflows/validate-native-evidence.md` passou a exigir caches graváveis,
`GOFLAGS=-mod=readonly -buildvcs=false`, logs/JSON preservados e classificação
explícita de limitações ambientais. O workflow de revisão e o guia editorial
apontam para essa fonte.

Não promoveremos checks nem alteraremos a matriz antes de um ciclo de pelo
menos quatro semanas com 95% de sucesso e evidência completa. A spec
`editorial-batch-engine` continua `in-progress`.

## Next checks

Durante a janela de 45 dias, executar o novo workflow para cada revisão que
use a matriz nativa e registrar a intervenção com
`pose recurrence-effect --register`. Ao fim da janela, executar
`pose recurrence-effect` e decidir `keep`, `adjust` ou `discard` com base na
taxa de recorrência e na qualidade dos artefatos.

## Risks

`govulncheck` continua dependente de DNS/rede e o cache de módulos pode não
estar pré-aquecido. Mitigação: usar CI ou ambiente com rede, ou um cache
gravável previamente populado; preservar a saída quando isso não for possível.
O `artifact-check --strict` também mantém avisos históricos de órfãos já
registrados na contribuição de proveniência; isso é uma reserva separada e
não deve ser confundido com `validate-native`.

## Next owner

@oseiaspereira; revisar o efeito após a janela de observação.

## References

- [spec editorial-batch-engine](../specs/2026-09-11-editorial-batch-engine.md)
- [workflow validate-native-evidence](../workflows/validate-native-evidence.md)
- [workflow recurrence-escalation](../workflows/recurrence-escalation.md)
- [contribuição sobre provenance do artifact-check](../contributions/20260915-081508-make-artifact-check-provenance-scope-exp.md)
