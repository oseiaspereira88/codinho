---
type: decision-log
slug: adr-concept-content-review
owner: @pose-maintainers
sensitivity: public-internal
created_at: 2026-09-08
last_reviewed_at: 2026-09-08
expires_at: 2026-12-07
source_refs:
  spec: "concept-content-authoring"
  workflow: ".pose/workflows/feature.md"
  commands: ["go test ./...", "go test ./internal/curriculum -run TestConceptContent", "go test ./cmd/codinho -run TestConceptContentOverRealStdio"]
  external_sources: []
---

# decision-log: adr-concept-content-review

## Context
last_reviewed_at: 2026-09-08. Consume knowledge:planning-audit-2026-09.
ConceptAuthoring previously held only id/title; canonical public explanations, examples and relation references are required for deliberate practice without disclosing challenge solutions or generating prose in the server.

## Current state
[ADR](../adr/2026-09-08-versioned-canonical-concept-content.md)
accepted. Implementation adds optional versioned content to ConceptAuthoring, schema validation, defensive copying, and projection over concept_content_get tool.

## Next checks
Verify schema/loader parity, legacy missing status, Unicode bounds, non-disclosure over real stdio, and lack of session mutation. Review before content version 2, external references, or personalization; 90-day TTL covers that window.

## Risks
Authors must carefully choose example contexts distinct from challenges; structural checks cannot automatically verify pedagogical differentiation.

## Next owner
@pose-maintainers — review before v1-integrated-acceptance.
