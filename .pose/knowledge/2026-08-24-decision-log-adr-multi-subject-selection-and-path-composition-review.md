---
type: decision-log
slug: adr-multi-subject-selection-and-path-composition-review
owner: @oseiaspereira
sensitivity: public-internal
created_at: 2026-08-23
last_reviewed_at: 2026-08-23
expires_at: 2026-11-21
source_refs:
  spec: "learning-track-composition"
  workflow: ".pose/workflows/feature.md"
  commands: []
  external_sources: []
---

# decision-log: adr-multi-subject-selection-and-path-composition-review

## Context

Rastreia o gatilho de revisão do ADR
[2026-08-23-multi-subject-selection-and-path-composition](../adr/2026-08-23-multi-subject-selection-and-path-composition.md),
que define seleção por conjunto de temas/competências com sinal de
cobertura, trilha completa pré-existente e composição de trilha por N
desafios quando nenhum item único cobre tudo.

## Current state

Decisão aceita, ainda sem implementação. A spec `learning-track-composition`
está `draft`, dependente de `curriculum-graph-path-recommendation`,
`session-orchestration-disclosure` e `mastery-review-scheduling` (todas
`done`).

## Next checks

- Ao implementar `learning-track-composition`, confirmar que
  `Query.Theme`/`Objective.ThemeID` foram substituídos (não duplicados) por
  suas versões em conjunto, e que toda resposta multi-assunto declara
  cobertura `total`/`parcial`/`nenhuma` sem exceção.
- Playtest real da composição de trilha por N assuntos para confirmar que a
  ordem respeita pré-requisitos percebidos pelo aluno.

## Risks

- Composição de trilha pode gerar sequências longas/pedagogicamente ruins
  se o número de assuntos pedidos for grande — sem mitigação automática
  planejada além de reportar o custo estimado (já parte do contrato de
  `curriculum-graph-path-recommendation` R5).

## Next owner

Mesmo owner.

## References

- ADR: `.pose/adr/2026-08-23-multi-subject-selection-and-path-composition.md`
- Spec: `.pose/specs/2026-08-24-learning-track-composition.md`
- Roadmap: `.pose/roadmaps/codinho-v1.md`
