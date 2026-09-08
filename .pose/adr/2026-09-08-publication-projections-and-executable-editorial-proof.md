# ADR: Publication projections and executable editorial proof

## Status
Accepted — 2026-09-08.

## Context
Inventory currently counts drafts toward V1. Normal MCP composition exposes them.
Fixture validation accepts failed/skipped checks without an authored expectation.
JSON schemas omit fields that Go already supports and unknown YAML fields vanish.

## Decision

- Keep Load/LoadPacks as administrative inventory; compose all normal MCP services
  from PublishedPacks. serve --authoring explicitly selects inventory for playtest.
  Pack publication gates auxiliary entities; challenge publication requires its
  own human metadata and a published containing pack. Reject invalid publication
  before runtime construction, including references to invisible prerequisites.
  Missing metadata remains a draft; never fabricate review or playtest.
- Separate inventory, drafts, published and eligible coverage. canonical: true is
  an explicit curation claim; variant_of excludes derivatives from eligible counts.
  The variants array remains private alternate text and never becomes additional
  challenges. variant_of instead identifies a derivative challenge with its own
  ID and human publication metadata: it may be visible when published, but never
  counts as an additional canonical challenge.
  A versioned distribution policy lists canonical pack groups and exact kind counts;
  per-pack constraints are enforced where planned, global constraints always at V1.
  Type groups express the product's combined debugging/refactoring/review target.
  Unexpected kinds and eligible packs outside the policy fail. No invented split
  of backend/production pack groups whose individual counts are not yet planned.
- Private validation metadata declares a complete reference_fixture and one
  baseline/reference expectation per check. Baseline normally uses starter fixture;
  a private baseline_fixture is the reproducible alternative when no starter exists
  and requires a nonempty justification. Imported unverifiable receipts cannot
  replace execution. Reference must pass. Expected failures name their category.
  Run each scenario/check in its own temporary workspace with allowlisted runners,
  network module fetch denied, no shell, timeout and uncached test execution.
  This is not an OS network sandbox; authored executable fixtures are trusted local
  code, as in existing --checks. Code execution is explicit in CLI or trusted CI.
- --checks validates all authored checks; --published-checks validates every
  published check. V1 implies published checks and eligible distribution/coverage.
  Missing proof never counts as success; zero published checks reports zero and
  does not establish V1 readiness. V1 rejects zero declared published checks;
  any skipped or incomplete scenario remains unverified. Emit only IDs, versions, digests and classified
  results, never reference code or process output in editorial reports.
- Clear private validation from every catalog challenge projection and JSON
  serialization, preserving historical pinned-content digests through optional
  JSON fields. Recovery continues from existing pinned public content.
- Reject unknown YAML fields; maintain schemas and Go structural rules with a
  shared valid/invalid corpus, including relations and children_mode. Go still
  owns graph/identity/publication semantics that JSON Schema cannot express.
  Preserve legacy minimal drafts; tighter publication rules apply at publication.

## Consequences

Normal serve starts with an empty curriculum until real content is published.
Authors use --authoring; administrative inventory remains accessible. New proof
metadata is opt-in for drafts but mandatory for published checks. Existing --checks
becomes stricter intentionally and reports legacy proof gaps. No authored files or
learner sessions are deleted. A proof attests runner behavior, not pedagogical
quality or human identity. Revisit before signed external CI evidence, generated
content quarantine, richer concept metadata or release acceptance changes.

Rejected alternatives: treating all authored items as eligible; accepting any
non-infrastructure check result; exposing reference code through pinned sessions;
using AI pre-review as publication.reviewed_by.
