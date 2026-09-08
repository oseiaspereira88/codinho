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
Executable CI and native evidence reviewed after commit b9b8868. Consult the
spec lifecycle and sealed attestation for the current formal closeout state.
Official POSE archive hashes and action commit pins
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
Native run 34273134508 passed on commit 6497259: Linux/amd64 and macOS/arm64,
Go 1.27.1, 19 checks each, source_modified=false. Artifacts were downloaded and
verified. Recheck the final commit and obtain independent security review.
Use /tmp/codinho-ci-tools for the pinned CI tools; keep local POSE 1.8.1 intact.
Run 34283684515 passed on b9b8868 in Linux/amd64 and macOS/arm64 with Go 1.27.1,
19 checks each and source_modified=false. Downloaded reports, gate logs and
native metadata were verified; readiness documentation records their hashes.
The separate review pass uses the same actor in a later execution, as allowed by
reviewer_independence=same-actor-separate-execution; it is not a human audit.
Refresh `pose index` before bundle preparation so it consumes current validation.
Component filters match no standalone Go modules: use the root matrix, which
actually executes all 19 checks. Never count an empty filtered run as coverage.
The installed review profile used unproducible test/validation classes. Adopted
the corrected distribution mapping without removing criteria or required tools;
the spec records the decision and the local contribution records reproduction.
Review pins at the next Go/POSE upgrade or security advisory. The 90-day TTL
covers the initial native CI adoption and requires a later toolchain review.

## Risks
Do not turn cross-compilation, authored YAML or a pending roadmap into runtime
or V1 acceptance evidence. Preserve historical attestations during reconciliation.

## Next owner
@pose-maintainers.

## References
[Spec](../specs/2026-09-07-v1-delivery-ci-assurance.md).
