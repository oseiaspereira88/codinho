---
type: decision-log
slug: adr-publication-review
owner: @pose-maintainers
sensitivity: public-internal
created_at: 2026-09-08
last_reviewed_at: 2026-09-08
expires_at: 2026-12-07
source_refs:
  spec: "catalog-publication-integrity"
  workflow: ".pose/workflows/feature.md"
  commands: ["go test ./...", "go test ./cmd/codinho -run TestPublicationIntegrity"]
  external_sources: []
---

# decision-log: adr-publication-review

## Context
last_reviewed_at: 2026-09-08. Consume knowledge:planning-audit-2026-09.
Inventory, publication and executed editorial proof need separate projections.

## Current state
[ADR](../adr/2026-09-08-publication-projections-and-executable-editorial-proof.md)
accepted. Implementation and independent validation govern delivery.

## Next checks
Verify published-only MCP, private reference fixtures, explicit expected failures,
legacy replay and full CLI gates. Review before generated content or external CI
receipts; 90-day TTL covers that integration window.

## Risks
Metadata cannot establish human identity or pedagogical quality automatically.

## Next owner
@pose-maintainers — review before release acceptance or signed proof contracts.
