# Review: learning-track-composition

Reviewed on 2026-09-09 against a7b2009 and ae0f7ae. The review covers
multi-subject selection, authored/composed tracks and durable advancement.
Human host playtests and V1 catalog publication remain in integrated acceptance.

## Rules applied during review

- Change type: Go feature, additive MCP/schema contracts and tutor documentation.
- Workflow: `.pose/workflows/review.md`.
- Security: bounded selectors and paths; no generated content, learner-file
  writes or execution during recommendation/start; evidence cannot cross a
  challenge boundary even when step/check IDs coincide.
- Documentation style: spec, ADR refinement, CLI and tutor guidance agree on
  explicit selection, union completeness and immutable accepted sequences.
- Delivery evidence/surface: require current registered stdio integration.
- Rule suggestions for internal, MCP, CLI, schemas and docs returned security
  and documentation-style. No frontend/infra changes or new dependencies.
- Consumed planning-audit-2026-09 and the existing multi-subject ADR decision log.

## Criteria and findings

- Scope: internal singular Query/Objective fields replaced by sets. Legacy MCP
  fields translate to sets; ambiguous combinations fail. The catalog data array
  is preserved and selection metadata is additive.
- Requirements: R1–R9 map to tests in the spec. Coverage distinguishes a complete
  single result from a partial union and reports missing subjects. Composition
  follows existing dependency direction with deterministic ID ties.
- Correctness: authored membership validates uniqueness, existence, bounds and
  dependency order. Combined prerequisite/relation cycles and excessive closure
  fail explicitly. Composition IDs bind accepted ordered challenge snapshots.
- Compatibility: legacy tracks remain browsable; empty membership cannot start.
  Historical single-challenge starts preserve their replay and request identity.
  Optional path state is versioned and hashed, with a validated event cursor.
- Security: full snapshots remain private; public path/status expose IDs,
  versions and cursor, not future instructions. Challenge transitions reset node
  progress/disclosure. Old evidence producers are excluded at path boundaries.
- Operability: missing/conflicting/stale selections return bounded errors.
  Restart uses pinned snapshots instead of silently consulting new pack content.
- Cross-component integration: real stdio accepts authored and composed paths,
  advances, restarts with the catalog replaced, retries the original start and
  finishes traversal. Normal evaluation/completion/advance is separately tested.
- Documentation: three choices and explicit acceptance are described in the
  tutor skill; CLI adds --themes/--competencies and keeps --theme.
- No unresolved critical/high findings identified in this pass. Review additions
  include mixed-cycle/cost limits, evidence-boundary isolation and measured p95.

## Validation and tool dispositions

Required: root `pose validate --strict` executes the single Go module's matrix,
including the registered learning-track-composition stdio check. Component
filters are not counted as independent tests of nonexistent nested modules.
Artifact reconciliation, skills-check, knowledge-check, docs-check and strict
spec lint passed. Integrate ran but recognizes zero MCP contracts; the actual
MCP contract suite is the relevant contract evidence.

Recommended: assess discover, assess tech-debt and per-component rule suggestions
ran. No TODO/FIXME/panic/stub markers were found. Recurrence scan found zero
flagged keys among 18 records over 14 days. Review/closeout completion tools are
executed after sealing and attestation.

Performance on the local Linux/amd64 runner, 84 indexed challenges:
search 28022 ns/op; composition 422569 ns/op. The explicit 100-query measurement
reported search p95 45.747µs and composition p95 601.877µs. A permanent test gates
both p95 values below 100ms. These are synthetic local measurements, not host
playtest results or production telemetry.

## Residual scope

Composition includes all matching candidates and prerequisite closure, up to
100 challenges; it does not claim a minimal or pedagogically optimized path.
The existing follow-up points to v1-integrated-acceptance for human host tests.
No release or catalog publication is claimed by this review.

## Final gate

Decision: approved. All 21 required checks passed without skips on
 ae0f7ae56e5a597ffaedfb03bbbe00098279e87a. Structured evidence is retained at
`.pose/results/delivery-validation.json`; surface-check reported zero findings.
Bundle rvb-a7544352b351efd7 binds plan digest
sha256:d0c57cb9f306ccfac7228ec88ef1f394cb7b9ae803e0c18b45437a81522185f9.
Attestation rva-566fee7941ed6c5b passed review verify and review-check;
guarded `pose close spec:learning-track-composition` succeeded on 2026-09-09.
The decision is an agent review pass, not human pedagogical acceptance.
