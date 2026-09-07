---
schema_version: 1
generated_at: 2026-09-07T00:53:58Z
baseline_commit: d7587442bc9f2f5ba0fd751a3586b35006a41803
staleness_policy: max_age_days=7,max_commits=20
refresh_pending: 
---

# Project State

## Resumo executivo
<!-- state:curated -->

O codinho é um sistema local de prática deliberada assistida por agente,
inicialmente especializado no desenvolvimento de fluência em Go. A visão e o
escopo da V1 estão documentados em `PROJECT.md`. Há CLI, servidor MCP stdio,
domínio, catálogo e testes executáveis. A auditoria de 2026-09-07 UTC confirmou
validação estrita e race verdes, mas reproduziu falha de retomada após restart.
O catálogo do worktree contém 41 desafios autorados e zero publicados; a V1
permanece em implementação. Consulte o relatório
`.pose/reports/2026-09-07-doc-audit-auditoria-do-planejamento-v1.md`.

## Direção atual
<!-- state:curated -->

Priorize `session-recovery-version-pinning`, `session-tree-progression` e
`catalog-publication-integrity`. Prepare gates de CI e conteúdo canônico;
depois execute composição de trilhas e quarentena de rascunhos. A autoria
curricular pode continuar em paralelo, com revisão humana antes da publicação.
Use as dependências e os oito critérios declarativos de `codinho-v1` para o
aceite final; status históricos done não demonstram composição atual.

## Specs & Roadmaps
<!-- state:derived hash:8b8897981bae -->

- specs: total=34 draft=11 in-progress=1 blocked=0 done=22 superseded=0 abandoned=0
- roadmaps: total=1 active=1 done=0
- últimos closeouts:
  - spec:eventstore-idempotency-scope (2026-08-24)
  - spec:interview-mode (2026-08-23)
  - spec:learning-practice-debug-modes (2026-08-23)
  - spec:catalog-authoring-quality (2026-08-23)
  - spec:workspace-observation-baselines (2026-08-23)
  - ... e mais 17 (ver `pose_list_specs status:done`)

## Follow-ups
<!-- state:derived hash:9b01ea69d647 -->

- abertos: 6
- por criticidade: high=4 medium=2 low=0 sem-classificação=0
- vencidos (review < hoje): 0

## Capabilities
<!-- state:derived hash:80e8415e0da5 -->

- assessment: ausente (rode `pose assess init`)

## Decisões & Conhecimento
<!-- state:derived hash:0ab3afd5fafc -->

- ADRs: total=6
  - adr:2026-08-23-multi-subject-selection-and-path-composition.md
  - adr:2026-08-23-agent-generated-catalog-content-as-a-first-class-mode.md
  - adr:2026-08-22-local-versioned-catalog-and-event-state.md
  - adr:2026-08-22-learner-code-ownership-and-safe-checks.md
  - adr:2026-08-22-independent-pedagogical-transitions.md
- knowledge: total=8 ativo=8 expirado=0

## Validação & Evidência
<!-- state:derived hash:57eabb181e4c -->

- último registro: task=auditoria-do-planejamento-v1 outcome=partial (2026-09-07T00:43:12Z)
- últimos 30 dias: total=17 outcome_ok=15 outcome_outro=2
- reports revisados (.md): total=10
  - report:2026-09-07-doc-audit-auditoria-do-planejamento-v1.md
  - report:2026-08-23-standard-fechar-workspace-observation-baselines.md
  - report:2026-08-23-standard-validate-native.md
  - report:2026-08-23-standard-fechar-feedback-evaluation-progression.md
  - report:2026-08-22-doc-audit-initialize-pose-governance-for-ailearn.md

## Arquitetura
<!-- state:derived hash:26bc9a5cafff status:unavailable -->

GraphForge export local ainda não é publicado por nenhum produtor neste repositório; seção indisponível nesta versão (spec pose-project-state-artifact, Não-objetivos e Compatibilidade).

## Docs
<!-- state:derived hash:d5892e1cac69 -->

- manifest: ausente (rode `pose docs-init`)
