---
title: "Make artifact-check provenance scope explicit"
type: "enhancement"
module: ""
created_at: "2026-09-15T08:15:08Z"
status: staged
upstream: "https://github.com/oseiaspereira88/pose"
privacy: sanitized-synthetic
---

# Contribution Draft: Make artifact-check provenance scope explicit

## Observation

In a repository with a long-lived spec and several commits carrying the same
`POSE-Spec:` trailer, `pose artifact-check --spec <slug> --strict` reports the
historical trailer range together with orphan findings for many tracked paths
that are unrelated to the current batch. The output does not make the selected
change-set boundary prominent enough to distinguish current attribution from
repository-wide ownership diagnostics.

When an explicit `--from <rev> --to <rev>` range is used to inspect one
documentation-only batch, the command still compares that small observed range
with every artifact action declared by the spec. Unchanged historical artifacts
then appear as `action-mismatch`, even though the range is intentionally narrow.
The tolerant mode changes the exit code but does not clarify which findings are
expected consequences of the selected range.

## Synthetic reproduction

1. Create a small repository with a spec declaring `docs/catalog.md` and
   `packs/base.yaml` in `### Artifacts`.
2. Make an initial commit for both paths with `POSE-Spec: sample-spec`.
3. Make a later commit that changes only `docs/catalog.md`, using the same
   trailer.
4. Run both `pose artifact-check --spec sample-spec --strict` and
   `pose artifact-check --spec sample-spec --from <later-parent> --to <later>`.

The cumulative invocation should explain its full historical range and any
undeclared paths. The explicit-range invocation should report the one changed
path without presenting untouched historical declarations as a new mismatch.
No project code, private host, credential, or customer data is needed.

## Proposed improvement

Expose the provenance mode in the result (`cumulative-trailers` versus
`explicit-range`) and report findings per change set. For an explicit range,
scope action comparison to observed paths while retaining a separate, clearly
labelled completeness check for the spec's full artifact inventory. Scope
`orphan` diagnostics to the selected ownership context, or provide an explicit
repository-wide ownership mode when that is the intended audit.

Add regression coverage for the three-commit synthetic repository and preserve
the distinction between a real undeclared path and a historical artifact that
was simply outside the selected range. Do not downgrade genuine undeclared
changes or silently broaden a batch's attribution.

This is local feedback only; no upstream submission was made and no POSE gate
was bypassed.
