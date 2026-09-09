# Historical review closeout after engine upgrade

Current disposition: original strict-check failure no longer reproducible
after instance update; do not submit as an active defect.

Observed on 2026-09-09 with POSE 2.0.2 in an instance last updated by 1.8.1.
Local staged feedback only; no upstream submission.

## Synthetic reproduction

Open an existing project containing completed specs and sealed review bundles
created with an earlier engine. Retain historical attestations. Run the new
engine's strict structural check: completed scopes may require fresh reviews,
while a newly sealed, validated scope can pass its own closeout checks.

## Observed limitation and proposal

The diagnostic reports only that a fresh review is needed, without identifying
which historical input or compatibility rule invalidated each scope. Expose
per-scope mismatch reasons and an explicit migration plan. Preserve immutable
historical attestations; do not auto-approve, downgrade required checks or
restate delivery. Root cause has not been isolated by this report.

## Recheck after instance update — 2026-09-09

The binary and instance now both report POSE 2.0.2. `pose check --strict`
passes with no historical closeout errors. The updated policy records
evidence_vocabulary_reconciled_at. A completed scope's closeout check returns
terminal=true, lifecycle_done=true and no blockers, without resealing or
replacing its historical attestation during this recheck.

Direct `pose review verify` can still report superseded bundles, with deltas
for changed criteria or source/evidence inputs. This is distinct from the
original global structural failure: the closeout policy accepts the completed
scope while direct comparison identifies differences from its sealed subject.
Do not describe those different command results as proof that the original
26-error failure persists.

The evidence supports a legacy-instance/migration explanation, without proving
that version mismatch alone caused every original finding. No upstream issue
was submitted; retain this report as historical context.
