---
type: decision-log
slug: adr-session-cursor-review
owner: @pose-maintainers
sensitivity: public-internal
created_at: 2026-09-07
last_reviewed_at: 2026-09-07
expires_at: 2026-12-06
source_refs:
  spec: "session-tree-progression"
  workflow: ".pose/workflows/feature.md"
  commands: ["go test -race ./...", "go test ./cmd/codinho -run TestTreeProgressionOverRealStdio"]
  external_sources: []  # [{url: "", accessed_at: "YYYY-MM-DD"}]
---

# decision-log: adr-session-cursor-review

## Context

last_reviewed_at: 2026-09-07. Session-tree-progression adopts relative cursors,
ordered windows and explicit child choices. Consume knowledge:planning-audit-2026-09.

## Current state

Implemented the [ADR](../adr/2026-09-07-depth-aware-session-cursor-and-explicit-branch-selection.md)
under [session-tree-progression](../specs/2026-09-07-session-tree-progression.md).
Unit/replay, full race suite and real stdio scenarios passed; final structured
validation and independent attestation still gate closure. No multi-challenge
track or DAG join is introduced.

## Next checks

Run unit/replay and real stdio tree progression tests, structured validation and
independent review. Review before learning-track-composition changes session
scope. TTL of 90 days covers that planned integration review window.

## Risks

Coarse completion covers a subtree without inventing child evidence. Alternative
selection is exclusive; existing unmarked children remain an ordered sequence.

## Next owner

@pose-maintainers — review before multi-challenge track integration.
