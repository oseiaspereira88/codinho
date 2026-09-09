# Historical review closeout after engine upgrade

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
