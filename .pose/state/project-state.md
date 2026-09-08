---
schema_version: 1
generated_at: 2026-09-08T13:00:24Z
baseline_commit: ef456abf43a22c2502004aaf824c6e5992c09a72
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
<!-- state:derived hash:e02c27b1022b -->

- specs: total=35 draft=8 in-progress=1 blocked=0 done=26 superseded=0 abandoned=0
- roadmaps: total=1 active=1 done=0
- últimos closeouts:
  - spec:session-tree-progression (2026-09-08)
  - spec:catalog-publication-integrity (2026-09-08)
  - spec:evaluation-evidence-lineage (2026-09-07)
  - spec:session-recovery-version-pinning (2026-09-07)
  - spec:eventstore-idempotency-scope (2026-08-24)
  - ... e mais 21 (ver `pose_list_specs status:done`)

## Follow-ups
<!-- state:derived hash:2ebc3185c034 -->

- abertos: 7
- por criticidade: high=5 medium=2 low=0 sem-classificação=0
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
<!-- state:derived hash:d7d24da3fffa -->

- último registro: task=auditoria-do-planejamento-v1 outcome=partial (2026-09-07T00:43:12Z)
- últimos 30 dias: total=17 outcome_ok=15 outcome_outro=2
- reports revisados (.md): total=14
  - report:2026-09-08-review-catalog-publication-integrity.md
  - report:2026-09-08-review-session-tree-progression.md
  - report:2026-09-07-review-evaluation-evidence.md
  - report:2026-09-07-review-session-recovery.md
  - report:2026-09-07-doc-audit-auditoria-do-planejamento-v1.md

## Arquitetura
<!-- state:derived hash:26bc9a5cafff status:unavailable -->

GraphForge export local ainda não é publicado por nenhum produtor neste repositório; seção indisponível nesta versão (spec pose-project-state-artifact, Não-objetivos e Compatibilidade).

## Docs
<!-- state:derived hash:d5892e1cac69 -->

- manifest: ausente (rode `pose docs-init`)
