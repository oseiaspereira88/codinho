---
type: decision-log
slug: adr-durable-session-replay-review
owner: @oseiaspereira
sensitivity: public-internal
created_at: 2026-09-07
last_reviewed_at: 2026-09-07
expires_at: 2026-12-06
source_refs:
  spec: session-recovery-version-pinning
  workflow: .pose/workflows/bugfix.md
  commands: ["go test -race ./...", "go test ./cmd/codinho -run TestSessionRecoveryOverRealStdio"]
  external_sources: []  # [{url: "", accessed_at: "YYYY-MM-DD"}]
---

# decision-log: adr-durable-session-replay-review

## Context

O log persistia eventos sem dados suficientes para reconstruir sessões.
Revisão de 2026-09-07 escolheu replay de deltas com início versionado,
conteúdo fixado, digest de request e baselines persistidas.

## Current state

Implementação e testes Go completos passaram. Revisão independente corrigiu
retries entre ferramentas, corrupção/revisões, CLI read-only e reparação de
tentativa interrompida. Gates POSE finais ainda condicionam o closeout.
O consumo de evidência em avaliação continua sem garantia de origem/freshness;
evaluation-evidence-lineage foi planejada e bloqueia aceite V1.

## Next checks

Executar pose validate --strict com resultado estruturado, reconciliar artefatos
no commit e selar revisão independente. Medir startup com log representativo
antes de ampliar retenção. Rever o ADR se houver migração de payload ou
degradação mensurada de startup. TTL de 90 dias justificado por essa janela
de observação de volume/compatibilidade, sem transformar memória em regra.

## Risks

Sessões legadas não são reconstruídas por inferência; preserve seus eventos.
SIGKILL deixa lock órfão de resolução manual. Catálogo válido é pré-requisito
do startup. Backups de cauda integram retenção/exclusão de estado. Não confundir
isolamento de EvidenceGet com autorização para StepEvaluate.

## Next owner

@oseiaspereira — confirmar disposition da lacuna de consumo e priorizar sua
spec (15) após recuperação (10), antes da progressão em árvore (20).

## References

- [Recuperação](../specs/2026-09-07-session-recovery-version-pinning.md).
- [ADR](../adr/2026-09-06-durable-session-replay-with-pinned-content.md).
- [Consumo de evidências](../specs/2026-09-07-evaluation-evidence-lineage.md).
