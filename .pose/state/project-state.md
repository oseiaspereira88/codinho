---
schema_version: 1
generated_at: 2026-08-24T01:05:18Z
baseline_commit: d46d2a36482c161bf4a4ec9f300112aa45e8ea2e
staleness_policy: max_age_days=7,max_commits=20
refresh_pending: 
---

# Project State

## Resumo executivo
<!-- state:curated -->

O codinho é um sistema local de prática deliberada assistida por agente,
inicialmente especializado no desenvolvimento de fluência em Go. A visão e o
escopo da V1 estão documentados em `PROJECT.md`; ainda não existe um módulo de
aplicação nem comportamento funcional implementado.

## Direção atual
<!-- state:curated -->

Execute o roadmap `codinho-v1` pelas dependências registradas nas 26 specs.
Comece por `architecture-decision-baseline`; em seguida materialize a fundação
executável e o primeiro fluxo MCP local de ponta a ponta. Mantenha catálogo,
hardening e aceite final condicionados aos gates definidos pelas respectivas
specs.

## Specs & Roadmaps
<!-- state:derived hash:e3e40e6505c4 -->

- specs: total=29 draft=7 in-progress=0 blocked=0 done=22 superseded=0 abandoned=0
- roadmaps: total=1 active=1 done=0
- últimos closeouts:
  - spec:interview-mode (2026-08-23)
  - spec:reliability-observability-compatibility (2026-08-23)
  - spec:assistance-hints-detours (2026-08-23)
  - spec:catalog-authoring-quality (2026-08-23)
  - spec:workspace-observation-baselines (2026-08-23)
  - ... e mais 17 (ver `pose_list_specs status:done`)

## Follow-ups
<!-- state:derived hash:0afbfb87d94d -->

- abertos: 6
- por criticidade: high=0 medium=0 low=0 sem-classificação=6
- vencidos (review < hoje): 0

## Capabilities
<!-- state:derived hash:80e8415e0da5 -->

- assessment: ausente (rode `pose assess init`)

## Decisões & Conhecimento
<!-- state:derived hash:60637b078948 -->

- ADRs: total=0
- knowledge: total=0 ativo=0 expirado=0

## Validação & Evidência
<!-- state:derived hash:f374377a0212 -->

- último registro: task=fechar-workspace-observation-baselines outcome=pass (2026-08-23T04:32:07Z)
- últimos 30 dias: total=16 outcome_ok=15 outcome_outro=1
- reports revisados (.md): total=9
  - report:2026-08-23-standard-fechar-workspace-observation-baselines.md
  - report:2026-08-23-standard-validate-native.md
  - report:2026-08-23-standard-fechar-feedback-evaluation-progression.md
  - report:2026-08-22-doc-audit-initialize-pose-governance-for-ailearn.md
  - report:2026-08-22-standard-fundacao-go-executavel-do-ailearn.md

## Arquitetura
<!-- state:derived hash:6dd7053e596d status:active -->

- componentes: total=1 verificados=1 completude=100.0%
- linhas_de_codigo: producao=9925 testes=8903 total=18828
- linguagens: go, shell
- saude_de_codigo: TODOs=0 FIXMEs=0 panics=0 stubs=0
- integracoes: contratos=0 ativos=0 gaps=0
- ultimos_assessments: ver artefatos em .pose/assessments/ e .pose/state/

## Docs
<!-- state:derived hash:d5892e1cac69 -->

- manifest: ausente (rode `pose docs-init`)
