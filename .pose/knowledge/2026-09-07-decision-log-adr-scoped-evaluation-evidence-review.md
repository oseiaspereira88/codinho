---
type: decision-log
slug: adr-scoped-evaluation-evidence-review
owner: @oseiaspereira
sensitivity: public-internal
created_at: 2026-09-07
last_reviewed_at: 2026-09-07
expires_at: 2026-12-06
source_refs:
  spec: "evaluation-evidence-lineage"
  workflow: ".pose/workflows/feature.md"
  commands: ["go test -race ./...", "go test ./cmd/codinho -run TestEvaluationEvidenceOverRealStdio"]
  external_sources: []  # [{url: "", accessed_at: "YYYY-MM-DD"}]
---

# decision-log: adr-scoped-evaluation-evidence-review

## Context

last_reviewed_at: 2026-09-07. A avaliação aceitava blobs sem autorização de
sessão/nó e fingerprints obsoletos; a recuperação anterior apenas protegia leitura.

## Current state

ADR aceito: injetar validador, exigir produtores confirmados, registrar origem
qualitativa e preservar retries históricos. Implementação em andamento; as
regressões iniciais cross-session e drift reproduziram aceitação indevida e
agora rejeitam sem append. E2E com restart e revisão independente inicial passaram; a revisão selada e os gates finais permanecem pré-condições do closeout.

## Next checks

Validar E2E MCP com restart, check_id divergente, blob ausente, drift entre
amostragens e nota qualitativa. Verificar segregação entre autoria e execução.
Rever o ADR se snapshots isolados ou origem externa autenticada se tornarem
necessários; TTL de 90 dias cobre essa janela de observação do contrato.

## Risks

Fingerprint é amostragem no escopo dos globs; mutex de sessão não congela
escritores externos nem detecta ABA. Registro qualitativo é declaração do tutor,
não atestado independente de veracidade. Cobertura completa de critérios
não é definida por esta remediação.

## Next owner

@oseiaspereira — conferir gatilhos de revisão até 2026-12-06.

## References

- [Spec](../specs/2026-09-07-evaluation-evidence-lineage.md).
- [ADR](../adr/2026-09-07-scoped-evaluation-evidence-and-qualitative-registration.md).
