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
make check was executed locally: 18/19 checks passed, including race, vet,
build, scanner, MCP/CLI and smoke. catalog failed because draft checks lack
validation/reference_fixture metadata (38 declared, zero verified). Keep that
gate failing until editorial remediation supplies real fixtures. Negative
ci-assurance tests and docs-check (13 docs) pass. Native CI, independent security
review, contract mapping, prospective policies and provenance remain pending.
The report describes a dirty worktree, not an immutable candidate attestation.
Resume with the execution log in the spec; do not reinstall or overwrite local
POSE 1.8.1 (the CI-pinned tools are isolated under /tmp/codinho-ci-tools).
Review pins at the next Go/POSE upgrade or security advisory. The 90-day TTL
covers the initial native CI adoption and requires a later toolchain review.

## Risks
Do not turn cross-compilation, authored YAML or a pending roadmap into runtime
or V1 acceptance evidence. Preserve historical attestations during reconciliation.

## Next owner
@pose-maintainers.

## References
[Spec](../specs/2026-09-07-v1-delivery-ci-assurance.md).
