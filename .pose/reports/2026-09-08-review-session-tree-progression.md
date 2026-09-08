# Review: session-tree-progression

## Subject and decision

- Date: 2026-09-08 UTC.
- Candidate: `8e151596c30124e060ca1ba6a198ff46a5469aaf`, diff against `974a59d`.
- Implementation commits: `307ccb0`, `8e15159`.
- Independent reviewer: `agent:tree_review` (`/root/tree_review`),
  gpt-5.6-luna with high reasoning, native delegated execution.
- Method: `.claude/skills/agent-batch-review/SKILL.md` and
  `docs/agent-review-workflow.md`; reviewer received read-only scope.
- Decision: approved without reservations for the final candidate.

The initial candidate reopened an exhausted earlier layer when advancing at
a coarser depth. TestTreeCoarserAdvanceNeverReopensExhaustedContainers reproduced
that failure; subtree coverage now derives from completed descendants without
creating parent completion events. The fix also preserves the selected branch
when revisiting its completed decision checkpoint. Unit and real stdio
regressions cover the correction. Final reviewer verdict: "aprovado sem
ressalvas" for `8e15159` against `974a59d`.

## Rules and knowledge

Applied `.pose/workflows/review.md`, security and documentation-style rules
returned by suggest-review for cmd/internal. Reviewed authority boundaries,
safe errors, branch identity and immutable retries. Delivery-evidence requires
actual composition plus Git attribution; knowledge-governance requires a
valid decision log and review trigger. No backend extension is installed;
the registered Go module supplies unit/race, vet, build and vulnerability checks.
UI and infrastructure rules do not apply.

Consumed knowledge:planning-audit-2026-09,
knowledge:adr-durable-session-replay-review and knowledge:adr-session-cursor-review.

## Criteria

| Criterion | Evidence and conclusion |
|---|---|
| scope | Artifact-check reconciles 25 claims and both implementation commits with zero findings. |
| requirements | R1–R6 map to named regressions in the spec, including five starting depths, ordered layers, both exclusive alternatives and restart. |
| correctness | Per-node evaluation/hint/completion state survives granularity changes; failed appends leave the cursor intact; exhausted containers do not reopen. Full suite and race checks passed. |
| cross-component-integration | TestTreeProgressionOverRealStdio runs the real binary through multiple processes and proves catalog/session/MCP composition, safe rejection, selection and retry. |
| security | Only offered direct child IDs can be selected. Revisions, lifecycle and detour gates run before append. IDs are unambiguous within the challenge. Error packets omit paths; secret-pattern scan of 25 artifacts found no credentials. |
| compatibility | Empty ChildrenMode is omitted from JSON to preserve pinned-content digests. Legacy start/advance responses replay unchanged; new navigation events have an explicit version. The intentional sequence/coarse-completion behavior change is documented. |
| operability | Last-window overrides have a durable exhaustion acknowledgment. Session finish remains explicit. Invalid choices do not mutate logs; unknown navigation versions fail recovery explicitly. |
| documentation | ADR, content-authoring guide, compatibility notes and changelog describe cursor, sequence/choice, aggregate coverage and old-writer restrictions. |

## Validation

[Structured result](../results/delivery-validation.json): generated
2026-09-08T03:39:34Z at the candidate, 14 executed and 14 passed, zero failed,
errored or skipped. Scoped provenance:
`sha256:90ab8fc9cd8ed37a7e82230e84b83ee8573ebed451ad29fbb8ca93091ae756ad`.

The author ran the full race suite and reran session/cmd with race and count=1
after the correction. The independent reviewer confirmed `go test -count=1
./...`, race for affected modules, vet, build, stdio integration and diff check.
Govulncheck found no vulnerabilities. Local `codinho catalog validate --json`
returned no structural diagnostics for 41 authored challenges; this is not
editorial publication or human playtest approval.

Strict structure/spec, knowledge and skills checks passed. Recurrence-check
found zero flagged keys and technical-debt assessment found zero markers.
Assess-integrate reports zero recognized contracts due to the known MCP Go
inventory limitation; real stdio tests provide the integration evidence.
Surface-check: one composed target, zero findings, current candidate evidence.

## Tool dispositions

Run artifact-check, assess-integrate, assess-tech-debt, knowledge-check,
suggest-review for cmd/internal, surface-check and full root validation.
The single Go module contains both review components; its full checks cover
cmd and internal. Global discover covers both, so separate component discovery
is unnecessary. Review-check and closeout-check run after attestation.

The installed POSE 1.7.12 sealed the candidate after refreshing indexes from
the structured result. No project gate, policy or delivery target was weakened.

## Residual limits

Traversal supports ordered trees and exclusive child alternatives, not DAG
joins or multiple challenges. Coarse completion covers descendants for navigation
without asserting fine-grained evidence. Do not run older writers on logs with
new navigation events. Automated review does not replace editorial/human V1
acceptance. The existing follow-up remains covered by v1-integrated-acceptance.


## Sealed review

- Bundle: `rvb-e373f801791a11a5`.
- Plan digest: `sha256:737b172eb812189e3534ae172e932c4d0491f9284c923406548580931dfea432`.
- Attestation: `rva-317b7d1ca69ad04d`.
- Review verify and review-check: fresh=true, approved=true.
- Guarded closure applied on 2026-09-08. Closeout-check returned
  lifecycle_done=true, review_approved=true, terminal=true, next_action=none.
- Refreshed assessments and project state: one Go module, 11,335 production LOC,
  10,887 test LOC, zero TODO/FIXME markers.
- Staged a sanitized local contribution about the native validation-class warning
  in POSE 1.7.12. It is diagnostic ambiguity; review and delivery gates passed.
