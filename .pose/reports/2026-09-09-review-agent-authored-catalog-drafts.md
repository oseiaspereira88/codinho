# Review: agent-authored-catalog-drafts

Reviewed on 2026-09-09 in a separate review pass over the immutable range
adb3235..44bd0dc. The historical base commit 2c8b990 lacked a spec trailer;
`pose report --change-from/--change-to` records its attribution without rewriting
Git. Its accompanying instance update is historical context, not new product scope.

## Rules applied during review

Feature, Go/MCP contract and tutor documentation review using security,
documentation-style and delivery-evidence rules. No new dependency or direct LLM
call. The plan requires same-actor-separate-execution; this review is a separate
pass by the implementing agent, not a claimed independent human review.

## Criteria and findings

- Scope and requirements: host-authored packs enter a separate drafts event
  stream. Three MCP tools manage submission, explicit authoring review/export
  and confirmed logical removal. Four selection modes are exposed by the skill.
- Correctness: strict one-document YAML decoding, compatible schema, bounded
  bytes/depth/aliases, shared structural/editorial/publication checks, secret
  detection and allowlisted runner resolution occur before append. Submission
  neither materializes fixtures nor executes checks.
- Data integrity: digest/version and request identity are durable. Pack identity
  cannot be overwritten or reused after removal. Concurrent append receipts are
  checked against the requested digest/removal ID; revision conflicts fail.
- Compatibility: optional draft selectors preserve prior start identities.
  Complete track provenance is fixed conservatively, including future members.
  Replay uses pinned snapshots even after draft removal/expiry. Legacy evidence
  remains auditable but cannot prove reviewed publication.
- Security: explicit accept_draft is required for quarantined session starts.
  Reviewed mastery derives provenance from the owning session of the recorded
  producer, never a client claim. Missing or ambiguous origins are excluded.
  Evidence events retain provenance. Retry preserves original provenance and
  legacy retries do not present historical unreviewed state as reviewed mastery.
- Operability: 256 KiB per submission, 100 lifetime submissions/16 MiB raw YAML
  and 30-day availability bound the quarantine. Removal is a tombstone, not quota
  reclamation; existing privacy export/purge covers the shared durable log.
- Integration: TestDraftOverRealStdio exercises normal production serve, valid
  submit/retry, hidden default search, consent rejection, track start, explicit
  authoring export, removal, restart, pinned instruction, denied new start and
  privacy export/purge. Existing real-session evaluation tests demonstrate draft
  evidence does not affect reviewed mastery. The shared publication gate has
  both synthetic promotion tests and the existing publication integration suite.
- Documentation: user choice, persistent warning, private authoring export,
  quotas, migration semantics and human promotion are documented. No human
  playtest, pedagogical approval or publication is asserted by these tests.

Review corrections are implemented: unsupported schema and multiple YAML
messages fail; internal private validation fixtures use publication validators;
legacy retry responses remain conservative. No unresolved critical/high findings
were identified. The remaining human host acceptance is already assigned to
v1-integrated-acceptance in the spec's existing follow-up.

## Validation and tools

`pose validate --strict --json .pose/results/delivery-validation.json` passed on
44bd0dc: 22 executed, 22 passed, no failed/error/skipped checks. This includes
race tests, formatting, vet/build, security scanning, publication and real MCP
contracts. Filtered nested modules are not counted as separate validation: the
repository has one root Go module, covered by the complete matrix.

Artifact-check, surface-check, skills-check, docs-check, knowledge-check and
strict spec lint passed. `pose index` refreshed review evidence before sealing.
Assess integrate recognizes zero contracts; the actual MCP contract suite and
registered stdio integration are the substantive integration evidence. Discover,
tech-debt and component rule suggestions were run; no uncovered debt markers.
Recurrence scan: zero flagged keys. Review-check and closeout-check are deferred
until attestation as required by the plan.

Bundle: rvb-aa9f1bede3dde740.
Decision: approved for the implemented draft capability; human catalog publication
and host acceptance remain governed separately.
