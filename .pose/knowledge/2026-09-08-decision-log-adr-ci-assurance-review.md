---
type: decision-log
slug: adr-ci-assurance-review
owner: @pose-maintainers
sensitivity: public-internal
created_at: 2026-09-08
last_reviewed_at: 2026-09-08
expires_at: 2026-12-07
source_refs:
  spec: v1-delivery-ci-assurance
  workflow: feature
  commands: ["pose assess discover", "pose lint-spec v1-delivery-ci-assurance --ready-check"]
  external_sources: []  # [{url: "", accessed_at: "YYYY-MM-DD"}]
---

# decision-log: adr-ci-assurance-review

## Context
Required CI tools and native platform evidence are governed by the
[ADR](../adr/2026-09-08-native-ci-evidence-and-prospective-delivery-governance.md).

## Current state
Implementation in progress. Official POSE archive hashes and action commit pins
were checked on 2026-09-08. The latest remote baseline failed govulncheck against
an unpatched Go 1.25 toolchain; this is a real failure, not an approved exception.
last_reviewed_at: 2026-09-08.

## Next checks
make check now passes 19/19 locally. Five draft packs gained private validation
references: 38 checks / 76 scenarios verified; distributed fixtures, tests,
acceptance and publication fields were preserved. No pedagogical/publication
attestation was added. The first implementation commit dbd848e passed strict
artifact-check with only historical orphan warnings. MCP inventory and historical
rename are covered by tests; artifact/delivery policies are enabled. The native
report now marks source_modified and rejects mutated GitHub candidate sources.
Run native Linux/macOS CI and independent security review on the final commit.
Use /tmp/codinho-ci-tools for the pinned CI tools; keep local POSE 1.8.1 intact.
Review pins at the next Go/POSE upgrade or security advisory. The 90-day TTL
covers the initial native CI adoption and requires a later toolchain review.

## Risks
Do not turn cross-compilation, authored YAML or a pending roadmap into runtime
or V1 acceptance evidence. Preserve historical attestations during reconciliation.

## Next owner
@pose-maintainers.

## References
[Spec](../specs/2026-09-07-v1-delivery-ci-assurance.md).
