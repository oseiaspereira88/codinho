---
schema_version: 1
generated_at: 2026-08-23T04:01:25Z
baseline_commit: c1e10cedb4324ecafbefd5083feb14c2a0c9537a
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
<!-- state:derived hash:63e702e33686 -->

- specs: total=27 draft=18 in-progress=0 blocked=0 done=9 superseded=0 abandoned=0
- roadmaps: total=1 active=1 done=0
- últimos closeouts:
  - spec:assistance-hints-detours (2026-08-23)
  - spec:feedback-evaluation-progression (2026-08-23)
  - spec:local-event-store (2026-08-23)
  - spec:mcp-stdio-foundation (2026-08-23)
  - spec:session-orchestration-disclosure (2026-08-23)
  - ... e mais 4 (ver `pose_list_specs status:done`)

## Follow-ups
<!-- state:derived hash:8acb74144df3 -->

- abertos: 0
- por criticidade: high=0 medium=0 low=0 sem-classificação=0
- vencidos (review < hoje): 0

## Capabilities
<!-- state:derived hash:80e8415e0da5 -->

- assessment: ausente (rode `pose assess init`)

## Decisões & Conhecimento
<!-- state:derived hash:60637b078948 -->

- ADRs: total=0
- knowledge: total=0 ativo=0 expirado=0

## Validação & Evidência
<!-- state:derived hash:d453798f36e8 -->

- último registro: task=validate-native outcome=pass (2026-08-23T03:58:17Z)
- últimos 30 dias: total=14 outcome_ok=13 outcome_outro=1
- reports revisados (.md): total=8
  - report:2026-08-23-standard-validate-native.md
  - report:2026-08-23-standard-fechar-feedback-evaluation-progression.md
  - report:2026-08-22-doc-audit-initialize-pose-governance-for-ailearn.md
  - report:2026-08-22-standard-fundacao-go-executavel-do-ailearn.md
  - report:2026-08-22-standard-create-complete-ailearn-v1-spec-portfolio.md

## Arquitetura
<!-- state:derived hash:7c4507618bce status:active -->

- componentes: total=1 verificados=1 completude=100.0%
- linhas_de_codigo: producao=7447 testes=6160 total=13607
- linguagens: go
- saude_de_codigo: TODOs=0 FIXMEs=0 panics=0 stubs=0
- ultimos_assessments: ver artefatos em .pose/assessments/ e .pose/state/

## Docs
<!-- state:derived hash:d5892e1cac69 -->

- manifest: ausente (rode `pose docs-init`)
