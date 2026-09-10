---
schema_version: 1
generated_at: 2026-09-09T13:28:01Z
baseline_commit: ae0f7ae56e5a597ffaedfb03bbbe00098279e87a
staleness_policy: max_age_days=7,max_commits=20
refresh_pending: 
---

# Project State

## Resumo executivo
<!-- state:curated -->

O codinho é um sistema local de prática deliberada assistida por agente,
inicialmente especializado em Go. A visão e o escopo V1 estão em PROJECT.md.
Recuperação de sessão, autorização de evidências, navegação da árvore e
integridade de publicação foram implementadas e encerradas com revisão.
O candidato ef456ab passou 16 checks estruturados, suites de race e revisão
independente; consulte `.pose/reports/2026-09-08-review-catalog-publication-integrity.md`.
O catálogo real contém 41 desafios em rascunho, zero publicados. Use
`serve --authoring` para playtest local; publicação e aceite humano V1 seguem pendentes.

## Direção atual
<!-- state:curated -->

Priorize v1-delivery-ci-assurance (35) e concept-content-authoring (40),
respeitando a DAG vigente. Depois execute composição de trilhas e quarentena
de rascunhos. A autoria curricular pode continuar com revisão e playtest reais,
provas de baseline/referência e distribuição canônica em packs/distribution.json.
Use os cut criteria do roadmap codinho-v1; esta remediação não encerra a V1.

## Specs & Roadmaps
<!-- state:derived hash:50db194798ca -->

- specs: total=35 draft=5 in-progress=1 blocked=0 done=29 superseded=0 abandoned=0
- roadmaps: total=1 active=1 done=0
- últimos closeouts:
  - spec:learning-track-composition (2026-09-09)
  - spec:concept-content-authoring (2026-09-09)
  - spec:v1-delivery-ci-assurance (2026-09-08)
  - spec:session-tree-progression (2026-09-08)
  - spec:catalog-publication-integrity (2026-09-08)
  - ... e mais 24 (ver `pose_list_specs status:done`)

## Follow-ups
<!-- state:derived hash:a7066737ad47 -->

- abertos: 10
- por criticidade: high=5 medium=5 low=0 sem-classificação=0
- vencidos (review < hoje): 0

## Capabilities
<!-- state:derived hash:80e8415e0da5 -->

- assessment: ausente (rode `pose assess init`)

## Decisões & Conhecimento
<!-- state:derived hash:9b887c6ba15e -->

- ADRs: total=10
  - adr:2026-09-08-publication-projections-and-executable-editorial-proof.md
  - adr:2026-09-07-scoped-evaluation-evidence-and-qualitative-registration.md
  - adr:2026-09-07-depth-aware-session-cursor-and-explicit-branch-selection.md
  - adr:2026-09-06-durable-session-replay-with-pinned-content.md
  - adr:2026-08-23-multi-subject-selection-and-path-composition.md
- knowledge: total=12 ativo=12 expirado=0

## Validação & Evidência
<!-- state:derived hash:eeb7c47d4ac1 -->

- último registro: task=v1-delivery-ci-review outcome=pass (2026-09-08T22:05:19Z)
- últimos 30 dias: total=18 outcome_ok=16 outcome_outro=2
- reports revisados (.md): total=17
  - report:2026-09-09-review-learning-track-composition.md
  - report:2026-09-09-review-concept-content-authoring.md
  - report:2026-09-08-standard-v1-delivery-ci-review.md
  - report:2026-09-08-review-catalog-publication-integrity.md
  - report:2026-09-08-review-session-tree-progression.md

## Arquitetura
<!-- state:derived hash:10dff9420419 status:active -->

- componentes: total=10 verificados=10 completude=100.0%
- linhas_de_codigo: producao=15417 testes=17912 total=33329
- linguagens: go, shell
- saude_de_codigo: TODOs=0 FIXMEs=0 panics=0 stubs=0
- integracoes: contratos=0 ativos=0 gaps=0
- divida_tecnica: total=0 coberta=0 descoberta=0
- ultimos_assessments: ver artefatos em .pose/assessments/ e .pose/state/

## Docs
<!-- state:derived hash:d5892e1cac69 -->

- manifest: ausente (rode `pose docs-init`)
