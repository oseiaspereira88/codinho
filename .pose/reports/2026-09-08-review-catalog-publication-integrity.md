# Review: catalog-publication-integrity

## Subject and decision

- Date: 2026-09-08 UTC.
- Candidate: `ef456abf43a22c2502004aaf824c6e5992c09a72`, diff against
  `2d390450d6773ddfc0e08a30a0b91e7ea07d361f`.
- Implementation commits: `f9ca9c8`, `ef456ab`.
- Independent reviewer: `agent:publication_review` (`/root/publication_review`),
  gpt-5.6-luna with high reasoning, native delegated execution.
- Method: `.claude/skills/agent-batch-review/SKILL.md` and
  `docs/agent-review-workflow.md`; reviewer receives read-only source scope.
- Decision: approved without technical reservations for ef456ab.

Round 1 rejected the initial candidate for a V1 gate accepting zero declared
checks, generic timeout reporting, mixed fail/skip reporting and a missing pack
with a zero-count policy target. Each has a regression and correction in ef456ab.
The author also reproduced three published-schema discrepancies: missing title,
unknown challenge kind and future challenge schema version. Published content now
rejects each; legacy minimal drafts remain loadable.

The reviewer withdrew an initial variant_of visibility concern after checking R1:
a derivative with its own ID and publication is visible but excluded from the
canonical count. The older variants text array remains private. ADR/docs and real
MCP tests now distinguish them, including checks of pinned JSONL for private text.

Round 2 confirmed every source correction and passed affected race/count=1,
vet/build and diff checks. Its only reservation was confirmation of final structured
validation. The reviewer independently read the 16/16 result for ef456ab and sealed
bundle, withdrew that reservation, and returned: "aprovado sem ressalvas técnicas".
Human publication/playtest remains a separate requirement, not an open finding.

## Rules and knowledge

Applied feature/ADR/test-plan, review and delivery-evidence workflows, security,
documentation-style and knowledge-governance rules. No Go backend extension is
installed; the registered module supplies tests, vet, build and vulnerability
checks. CI adds an actual published-proof command and artifact without changing
release workflows or weakening gates.

Consumed knowledge:planning-audit-2026-09,
knowledge:adr-local-versioned-catalog-and-event-state-review and
knowledge:adr-publication-review. ADR documents visibility, canonical eligibility,
private reproducible proof, human metadata and compatibility before implementation.

## Criteria

| Criterion | Evidence and conclusion |
|---|---|
| scope | Artifact-check reconciles 39 claims and 40 observed paths including the spec, across two commits, with zero findings. |
| requirements | R1–R8 map to named tests in the spec. CLI and stdio tests exercise actual composition, not only helper projections. |
| correctness | Global/per-pack distribution detects unexpected kinds, missing packs and inconsistent policy totals. V1 uses eligible content and requires published check proof. Each check has isolated baseline/reference execution; no tests, skip, mixed failure, timeout, truncation and unexpected outcomes fail. |
| cross-component-integration | TestPublicationIntegrityOverRealStdio covers search, recommendation, relations, get/start rejection, explicit authoring and recovery. TestPublicationIntegrityCLI covers successful published proof, wrong reference, sufficient drafts and sufficient publication with zero checks. |
| security | Reference fixtures are cleared from catalog records and omitted from session JSON. Reports contain IDs, versions, digests and outcome categories only. Fixture paths remain confined; manifest symlink escapes and unknown fields fail. Allowlisted execution has timeout and bounded capture. Network module fetch denial is not represented as an OS network sandbox. |
| compatibility | Missing publication remains draft. Existing process fixtures explicitly request authoring. Old replay/evidence/tree tests pass without rewriting IDs or JSONL; optional empty fields preserve historical content digests. |
| operability | Reports distinguish inventory/drafts/published/eligible and declared/verified checks, identify binary revision or unknown explicitly, and include Go version. Zero publication is reported as zero, not V1 readiness. |
| documentation | ADR, authoring guide, compatibility notes, versioned policy, shared schema corpus and changelog describe the stricter commands and human publication requirements. |

## Validation

[Structured result](../results/delivery-validation.json) generated
2026-09-08T12:55:02Z on ef456ab: 16 executed, 16 passed, zero failures,
errors or skips. Govulncheck found no vulnerabilities. Scoped provenance:
`sha256:aad9e2679bd78a2dd31f58b56a447e1e937f26e4f15a3ddc635d4ae72baa0b50`.
Surface-check verifies one composed target with zero findings. The initial reviewer
invocation without the govulncheck PATH failed environmentally; it was replaced
by this complete run and is not passing evidence.

Full `go test -race ./...` passed on the initial implementation. After corrections,
`go test -race ./internal/curriculum/... ./internal/checks/... ./internal/cli/...
./cmd/codinho/... -count=1` passed. The initial independent reviewer also ran full
race/count=1, vet/build and composed tests successfully.

Strict structure/spec and knowledge/skills checks passed. Recurrence-check found
zero flagged keys; tech-debt assessment found zero markers. Assess-integrate
reports zero recognized contracts under the known MCP Go inventory limitation;
real stdio tests supply integration evidence. Final surface check consumed the
fresh structured result and passed.

## Tool dispositions

Run artifact-check, assess-integrate, assess-tech-debt, knowledge-check,
suggest-review for cmd/internal/packs/schemas/testdata, surface-check and full root validation.
The root Go module covers cmd, internal, packs, schemas and testdata through
compiled tests, corpus validation and actual catalog loading. Global discover covers their
assessment needs. Review-check and closeout-check run after attestation.

Refresh pose index after structured validation so review consumes current evidence.
Known validation evidence-class/unmapped-component warnings do not replace the
actual gate result; the sanitized prior contribution already records the warning.

## Residual limits

The real catalog has 41 authored drafts and zero published challenges. No human
review/playtest or real publication is fabricated. Runners verify mechanical
behavior, not pedagogical quality. Schema corpus covers the shared structural
contract; Go additionally owns graph, identity and publication semantics. Local
fixtures are trusted executable code, not an OS-sandbox guarantee. Backend and
production distribution uses planned groups until individual splits are authored.
V1 acceptance remains covered by the existing v1-integrated-acceptance follow-up.


## Sealed review and closeout

- Bundle: `rvb-838a0d1e9391698a`.
- Plan digest: `sha256:d815afe523f1b9cce3bd06bacae75823679c8d95d32b51b8e365da1bb1392c03`.
- Attestation: `rva-f35c1de4d1992136`, recorded by native auto-attest after
  the independent final verdict and review report were present.
- Review verify/review-check: fresh=true, approved=true.
- Guarded `pose close` applied on 2026-09-08. Closeout-check confirms
  lifecycle_done=true, review_approved=true, terminal=true, next_action=none.
- Final discovery: one Go module, 12,127 production LOC and 11,524 test LOC;
  zero TODO/FIXME/debt markers. Project state refreshed after closure.
- Existing follow-up remains covered by v1-integrated-acceptance. No source,
  policy, review metadata or human publication requirement was weakened to close.
