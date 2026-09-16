---
schema_version: 1
generated_at: 2026-09-16T03:41:06Z
baseline_commit: 4067f89f02aedba8c9b02cf75f6ed46872396fcd
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
<!-- state:derived hash:f51da1136d98 -->

- specs: total=36 draft=4 in-progress=2 blocked=0 done=30 superseded=0 abandoned=0
- roadmaps: total=1 active=1 done=0
- últimos closeouts:
  - spec:agent-authored-catalog-drafts (2026-09-09)
  - spec:learning-track-composition (2026-09-09)
  - spec:concept-content-authoring (2026-09-09)
  - spec:v1-delivery-ci-assurance (2026-09-08)
  - spec:session-tree-progression (2026-09-08)
  - ... e mais 25 (ver `pose_list_specs status:done`)

## Follow-ups
<!-- state:derived hash:1ab12e4b3c46 -->

- abertos: 10
- por criticidade: high=5 medium=5 low=0 sem-classificação=0
- vencidos (review < hoje): 1
  - spec:session-recovery-version-pinning (owner:@oseiaspereira review:2026-09-14)

## Capabilities
<!-- state:derived hash:80e8415e0da5 -->

- assessment: ausente (rode `pose assess init`)

## Decisões & Conhecimento
<!-- state:derived hash:b230590fe157 -->

- ADRs: total=13
  - adr:2026-09-11-editorial-worktrees-and-explicit-sessions.md
  - adr:2026-09-08-versioned-canonical-concept-content.md
  - adr:2026-09-08-publication-projections-and-executable-editorial-proof.md
  - adr:2026-09-08-native-ci-evidence-and-prospective-delivery-governance.md
  - adr:2026-09-07-scoped-evaluation-evidence-and-qualitative-registration.md
- knowledge: total=18 ativo=18 expirado=0

## Validação & Evidência
<!-- state:derived hash:a5ae1c75b718 -->

- último registro: task=validate-native outcome=partial (2026-09-16T02:32:11Z)
- últimos 30 dias: total=38 outcome_ok=24 outcome_outro=14
- reports revisados (.md): total=92
  - report:2026-09-16-standard-validate-native.md
  - report:2026-09-15-standard-validate-native.md
  - report:2026-09-10-data-text-csv-quote-node.md
  - report:2026-09-10-data-text-map-nil-node.md
  - report:2026-09-10-data-text-truncate-node.md

## Arquitetura
<!-- state:derived hash:059dc41e1f5b status:active -->

- componentes: total=10 verificados=10 completude=100.0%
- linhas_de_codigo: producao=16669 testes=18400 total=35069
- linguagens: go, python, shell
- saude_de_codigo: TODOs=0 FIXMEs=0 panics=0 stubs=0
- integracoes: contratos=0 ativos=0 gaps=0
- divida_tecnica: total=0 coberta=0 descoberta=0
- ultimos_assessments: ver artefatos em .pose/assessments/ e .pose/state/

## Docs
<!-- state:derived hash:c646f4fadde7 -->

- manifest: presente, profile=cli, roots=docs
- docs: declaradas=16 não-declaradas=0 vencidas=0 erros=0 avisos=0
