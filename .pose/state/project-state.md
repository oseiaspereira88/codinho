---
schema_version: 1
generated_at: 2026-09-08T03:44:13Z
baseline_commit: 8e151596c30124e060ca1ba6a198ff46a5469aaf
staleness_policy: max_age_days=7,max_commits=20
refresh_pending: 
---

# Project State

## Resumo executivo
<!-- state:curated -->

O codinho é um sistema local de prática deliberada assistida por agente,
inicialmente especializado em Go. A visão e o escopo V1 estão em PROJECT.md.
Recuperação de sessão, autorização de evidências e navegação completa da árvore
foram implementadas e encerradas com revisão e evidências. O candidato 8e15159
passou os 14 checks estruturados; consulte
`.pose/reports/2026-09-08-review-session-tree-progression.md`.
O catálogo local contém 41 desafios autorados e passou validação estrutural;
publicação editorial e aceite humano V1 continuam pendentes.

## Direção atual
<!-- state:curated -->

Priorize catalog-publication-integrity (30), v1-delivery-ci-assurance (35)
e concept-content-authoring (40), respeitando a DAG vigente. Depois execute
composição de trilhas e quarentena de rascunhos. A autoria curricular pode
continuar com os gates de revisão/publicação aplicáveis. Use os cut criteria
do roadmap codinho-v1; fechamento de uma remediação não encerra a V1.

## Specs & Roadmaps
<!-- state:derived hash:61ae00177ccd -->

- specs: total=35 draft=9 in-progress=1 blocked=0 done=25 superseded=0 abandoned=0
- roadmaps: total=1 active=1 done=0
- últimos closeouts:
  - spec:session-tree-progression (2026-09-08)
  - spec:evaluation-evidence-lineage (2026-09-07)
  - spec:session-recovery-version-pinning (2026-09-07)
  - spec:eventstore-idempotency-scope (2026-08-24)
  - spec:safe-check-executor (2026-08-23)
  - ... e mais 20 (ver `pose_list_specs status:done`)

## Follow-ups
<!-- state:derived hash:2ebc3185c034 -->

- abertos: 7
- por criticidade: high=5 medium=2 low=0 sem-classificação=0
- vencidos (review < hoje): 0

## Capabilities
<!-- state:derived hash:80e8415e0da5 -->

- assessment: ausente (rode `pose assess init`)

## Decisões & Conhecimento
<!-- state:derived hash:7ea9462f2c34 -->

- ADRs: total=9
  - adr:2026-09-07-scoped-evaluation-evidence-and-qualitative-registration.md
  - adr:2026-09-07-depth-aware-session-cursor-and-explicit-branch-selection.md
  - adr:2026-09-06-durable-session-replay-with-pinned-content.md
  - adr:2026-08-23-multi-subject-selection-and-path-composition.md
  - adr:2026-08-23-agent-generated-catalog-content-as-a-first-class-mode.md
- knowledge: total=11 ativo=11 expirado=0

## Validação & Evidência
<!-- state:derived hash:5a7df8680967 -->

- último registro: task=auditoria-do-planejamento-v1 outcome=partial (2026-09-07T00:43:12Z)
- últimos 30 dias: total=17 outcome_ok=15 outcome_outro=2
- reports revisados (.md): total=13
  - report:2026-09-08-review-session-tree-progression.md
  - report:2026-09-07-review-evaluation-evidence.md
  - report:2026-09-07-review-session-recovery.md
  - report:2026-09-07-doc-audit-auditoria-do-planejamento-v1.md
  - report:2026-08-23-standard-fechar-workspace-observation-baselines.md

## Arquitetura
<!-- state:derived hash:26bc9a5cafff status:unavailable -->

GraphForge export local ainda não é publicado por nenhum produtor neste repositório; seção indisponível nesta versão (spec pose-project-state-artifact, Não-objetivos e Compatibilidade).

## Docs
<!-- state:derived hash:d5892e1cac69 -->

- manifest: ausente (rode `pose docs-init`)
