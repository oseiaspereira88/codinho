# Review: session-recovery-version-pinning

## Subject and decision

- Date: 2026-09-07 UTC.
- Candidate: `8c2ee6e794bde10070666ac89dd58afe0af537a6`.
- Reviewer: `agent:codex-root-resumed-review-20260907`.
- Independence: separate execution reviewing the implementation inherited
  from the preceding session; no second agent or human review is claimed.
- Decision: approved for the recovery scope defined by R1–R7 and Decision 2.
- Knowledge consumed: knowledge:planning-audit-2026-09 and
  knowledge:adr-durable-session-replay-review.

## Rules applied during review

- Change type: feature with compatibility documentation and regression tests.
- Workflow consulted: `.pose/workflows/review.md`.
- `.pose/rules/security.md`: inspect persisted content, error messages,
  evidence scope, root confinement and private tail backups.
- `.pose/rules/documentation-style.md`: reconcile spec, ADR, changelog and
  compatibility instructions with the implemented behavior.
- `.pose/rules/delivery-evidence.md` and `delivery-surface.md`: require
  immutable attribution and fresh production-entrypoint integration evidence.
- `.pose/rules/knowledge-governance.md`: validate the decision log and its TTL.
- No backend extension is installed; `pose suggest review` resolves security
  and documentation rules. Go test, race, vet, build and govulncheck were run.
- Frontend and infrastructure rules do not apply to this change.

## Criteria and evidence

| Criterion | Review and evidence |
|---|---|
| scope | 24 artifact claims reconcile with the implementation commit, zero findings; the authored changes in `packs/go-errors.yaml` are outside this commit. |
| requirements | R1–R7 have named regressions in the spec; inspected recovery, workspace and stdio tests against their assertions. |
| correctness | Inspected reducer/pre-append cloning, historical retries, request conflicts, attempt repair and revision/envelope checks; `go test -race ./...` passed. |
| cross-component-integration | `session-recovery` runs real processes through `cmd/codinho/main.go`; catalog update/removal, lost response, truncated tail and legacy/corrupt logs passed. |
| security | Cross-session retrieval and replaced-root regressions passed; inspected safe MCP diagnostic and 0600 backups. Secret-pattern scan of the 25 candidate files found no matching credentials; govulncheck found no vulnerabilities. |
| compatibility | JSONL envelope remains v1; incomplete legacy sessions fail explicitly; CLI readers use OpenReadOnly; rollout/rollback constraints are documented. |
| operability | No silent repair of terminated corruption; tail bytes are preserved before truncation; orphan-lock diagnosis remains an explicit operational step. |
| documentation | ADR, compatibility instructions, requirement trace and changelog match recovery. Evaluation consumption remains explicitly deferred to its own draft spec. |

Structured validation: [delivery-validation.json](../results/delivery-validation.json),
12 passed, zero failed/errored/skipped. `pose check --strict`, strict spec lint,
knowledge-check, skills-check, artifact-check and surface-check passed.
Technical-debt assessment returned zero markers; recurrence-check returned
zero flagged keys. `assess integrate` recognizes no MCP Go contracts, so its
empty inventory is not integration proof; the real MCP tests supply that proof.

## Tool dispositions

Run artifact-check, assess-integrate, assess-tech-debt, knowledge-check,
suggest-review for cmd/internal, surface-check and the full root validation.
The root Go module contains both cmd and internal; its full validation covers
both review components. Global assess-discover covered both; separate
component invocations were unnecessary. Defer review-check and closeout-check
until the attestation exists, then execute both before lifecycle transition.
The sealed bundle records the effective plan digest and the attestation records
each tool disposition.

## Findings and limitations

No new blocking implementation finding in this review. The pre-existing
evaluation-consumption gap remains tracked in evaluation-evidence-lineage,
priority 15, blocking V1 acceptance; this approval does not approve that gap
or close V1. Legacy replay, orphan locks, valid startup catalog and startup
growth constraints remain as described in the spec and decision log.

The installed POSE executable has the known root-module evidence defect
reported in contributions/20260823-003416-delivery-target-module-can-never-satisfy.md.
Use `/tmp/codinho-pose-current`, compiled without source changes from the clean
local POSE checkout at `bc1b8b9`, for review sealing and closeout. Its focused
`TestReviewBundleMatchesRootModuleValidationEvidenceForSubdirectoryTargets`
passed. This uses the committed upstream correction, with unchanged project
policies, delivery target, validation results and requirements.

## Sealed review

- Bundle: `rvb-a431e6d8677f4837`.
- Bundle digest: `sha256:a431e6d8677f48372ef8eb81cb0cbc8cb912b1bc592581e6fa6664d51a27b53b`.
- Plan digest: `sha256:3729c2f34abfc0ba3dfa41e51870d3bbcff7d1ff6faea6900d746c6f6832964f`.
- Attestation: `rva-339a9f4375ed88a7`.
- `review verify`: ready-to-close, fresh=true, approved=true.
- `review-check`: approved=true, fresh=true.
- Before lifecycle transition, `closeout-check` reports only status not done;
  its next action is the guarded transition. Recheck after `pose close`.
- Guarded close applied on 2026-09-07. Subsequent closeout-check returned
  lifecycle_done=true, review_approved=true, terminal=true, next_action=none.
- Refreshed assessments and project state with assess discover --update-state.
  Added the advertised public-claims command to the manual so the current
  engine's strict structure check recognizes its command reference.
