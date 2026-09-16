---
type: handoff
slug: editorial-batch-review-reservations
owner: @oseiaspereira
sensitivity: public-internal
created_at: 2026-09-16
last_reviewed_at: 2026-09-16
expires_at: 2026-12-15
source_refs:
  spec: editorial-batch-engine
  workflow: review
  commands:
    - "pose review bundle spec:editorial-batch-engine --seal"
    - "pose review auto-attest rvb-19986755983b3ea2 --apply"
    - "pose review verify rvb-19986755983b3ea2"
    - "pose review-check spec:editorial-batch-engine"
    - "pose validate --strict"
  external_sources: []  # [{url: "", accessed_at: "YYYY-MM-DD"}]
---

# handoff: editorial-batch-review-reservations

## Context

A revisão da spec `editorial-batch-engine` exigiu 15 ferramentas e oito
critérios. O bundle foi selado como `rvb-19986755983b3ea2`; a atestação
automática foi registrada como `rva-7e2ea0b9f0e0ec9b`.

## Current state

`pose review verify` e `pose review-check` retornam `ready-to-close`, mas a
atestação reutiliza evidência de validação produzida no commit `f061bc...`.
O bundle registra essa diferença como reserva de frescor, além de avisos de
paths de governança sem componente mapeado. Isso não foi promovido a falha do
executor porque os checks focados do mecanismo passaram e a mudança atual é
predominantemente de workflow/documentação.

A spec permanece `in-progress`. Nenhuma revisão humana, playtest ou publicação
foi declarada.

## Next checks

Executar `pose validate --strict` no workflow de evidência nativa, com `GOCACHE`
e `GOMODCACHE` graváveis, e preservar JSON, log e exit code. Reexecutar a
validação quando rede/DNS estiver disponível para `govulncheck`. Após a janela
de 45 dias da intervenção, executar `pose recurrence-effect` e decidir se o
workflow deve ser mantido, ajustado ou descartado.

## Risks

Não usar a atestação automática como prova de uma validação atual quando o
bundle indicar commit de evidência diferente. Não fechar a spec até que os
gates humanos e as tarefas editoriais restantes tenham disposição explícita.

## Next owner

@oseiaspereira.

## References

- [spec editorial-batch-engine](../specs/2026-09-11-editorial-batch-engine.md)
- [bundle sealed](../review-bundles/rvb-19986755983b3ea2.json)
- [attestation](../review-attestations/rva-7e2ea0b9f0e0ec9b.json)
- [workflow validate-native-evidence](../workflows/validate-native-evidence.md)
- [decision log da recorrência](2026-09-16-decision-log-escalation-validate-native.md)
