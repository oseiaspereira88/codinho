---
type: decision-log
slug: adr-agent-generated-catalog-content-as-a-first-class-mode-review
owner: @oseiaspereira
sensitivity: public-internal
created_at: 2026-08-23
last_reviewed_at: 2026-08-23
expires_at: 2026-11-21
source_refs:
  spec: "agent-authored-catalog-drafts"
  workflow: ".pose/workflows/feature.md"
  commands: []
  external_sources: []
---

# decision-log: adr-agent-generated-catalog-content-as-a-first-class-mode-review

## Context

Rastreia o gatilho de revisão do ADR
[2026-08-23-agent-generated-catalog-content-as-a-first-class-mode](../adr/2026-08-23-agent-generated-catalog-content-as-a-first-class-mode.md),
que torna geração de conteúdo pelo agente-tutor conectado ao MCP um modo de
seleção de primeira classe (trilha, N-desafios ou desafio único), sempre
entrando como `publication.status: draft` em quarentena, com evidência de
sessão marcada por proveniência (`content_provenance: draft`) e nunca
promovida à ladder de maestria revisada até o rascunho virar `published`
pelo mesmo funil autor/revisor de `catalog-authoring-quality`.

## Current state

Decisão aceita, ainda sem implementação. A spec
`agent-authored-catalog-drafts` está `draft`, dependente de
`catalog-authoring-quality`, `learning-track-composition` e
`tutor-skill-host-integration` (as duas primeiras `done`/`draft`
respectivamente).

## Next checks

- Ao implementar, confirmar que a nova tool de submissão de rascunho reusa
  `internal/curriculum.Validate`/`RunEditorialChecks` sem duplicar regras.
- Confirmar que nenhuma sessão sobre rascunho grava evento de mastery sem o
  campo de proveniência `draft`, e que a projeção de `mastery-review-
  scheduling` ignora evidência de proveniência `draft` ao calcular estado.
- Confirmar que o servidor MCP continua sem chamar LLM diretamente — a tool
  nova só recebe e valida o que o agente do host já gerou.

## Risks

- Volume de rascunhos gerados em sessão real pode superar a capacidade de
  revisão de um revisor único (fase 1) antes da curadoria em duas camadas
  (fase 2) estar pronta — mitigação: acompanhar contagem de rascunhos
  pendentes e antecipar a fase 2 se necessário.
- Aluno pode confundir sessão sobre rascunho com conteúdo oficial revisado
  se a skill não deixar o aviso explícito e persistente na interação.

## Next owner

Mesmo owner.

## References

- ADR: `.pose/adr/2026-08-23-agent-generated-catalog-content-as-a-first-class-mode.md`
- ADR irmã: `.pose/adr/2026-08-23-multi-subject-selection-and-path-composition.md`
- Spec: `.pose/specs/2026-08-24-agent-authored-catalog-drafts.md`
- Spec: `.pose/specs/2026-08-22-catalog-authoring-quality.md`
- Roadmap: `.pose/roadmaps/codinho-v1.md`
