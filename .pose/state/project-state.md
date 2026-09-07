---
schema_version: 1
generated_at: 2026-09-07T14:42:50Z
baseline_commit: b978493870ff03b4133bc845b037be296ab154b6
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
<!-- state:derived hash:60862ac2c241 -->

- specs: total=35 draft=11 in-progress=1 blocked=0 done=23 superseded=0 abandoned=0
- roadmaps: total=1 active=1 done=0
- últimos closeouts:
  - spec:session-recovery-version-pinning (2026-09-07)
  - spec:eventstore-idempotency-scope (2026-08-24)
  - spec:interview-mode (2026-08-23)
  - spec:security-privacy-hardening (2026-08-23)
  - spec:workspace-observation-baselines (2026-08-23)
  - ... e mais 18 (ver `pose_list_specs status:done`)

## Follow-ups
<!-- state:derived hash:2ebc3185c034 -->

- abertos: 7
- por criticidade: high=5 medium=2 low=0 sem-classificação=0
- vencidos (review < hoje): 0

## Capabilities
<!-- state:derived hash:80e8415e0da5 -->

- assessment: ausente (rode `pose assess init`)

## Decisões & Conhecimento
<!-- state:derived hash:21f861ec45d8 -->

- ADRs: total=7
  - adr:2026-09-06-durable-session-replay-with-pinned-content.md
  - adr:2026-08-23-multi-subject-selection-and-path-composition.md
  - adr:2026-08-23-agent-generated-catalog-content-as-a-first-class-mode.md
  - adr:2026-08-22-local-versioned-catalog-and-event-state.md
  - adr:2026-08-22-learner-code-ownership-and-safe-checks.md
- knowledge: total=9 ativo=9 expirado=0

## Validação & Evidência
<!-- state:derived hash:4d8da336d134 -->

- último registro: task=auditoria-do-planejamento-v1 outcome=partial (2026-09-07T00:43:12Z)
- últimos 30 dias: total=17 outcome_ok=15 outcome_outro=2
- reports revisados (.md): total=11
  - report:2026-09-07-review-session-recovery.md
  - report:2026-09-07-doc-audit-auditoria-do-planejamento-v1.md
  - report:2026-08-23-standard-fechar-workspace-observation-baselines.md
  - report:2026-08-23-standard-validate-native.md
  - report:2026-08-23-standard-fechar-feedback-evaluation-progression.md

## Arquitetura
<!-- state:derived hash:26bc9a5cafff status:unavailable -->

GraphForge export local ainda não é publicado por nenhum produtor neste repositório; seção indisponível nesta versão (spec pose-project-state-artifact, Não-objetivos e Compatibilidade).

## Docs
<!-- state:derived hash:d5892e1cac69 -->

- manifest: ausente (rode `pose docs-init`)
