# Installed review profile requests unproducible evidence classes

Current disposition: not reproducible in the updated instance; do not submit
as an active defect on the basis of this report.

Observed on 2026-09-08 with POSE 1.8.1. Local feedback only.

## Synthetic reproduction

1. Configure a Go module with passing required checks whose evidenceClass is
   unit or integration and a spec with an attributed delivery target.
2. Retain an installed spec-closeout profile requesting test/validation for
   correctness and requirement-trace/test for requirements.
3. Generate structured validation and refresh the delivery index.
4. Prepare a review bundle. The validate tool warns that validation cannot be
   emitted; the criterion classes remain impossible to satisfy honestly.

## Proposed solution

Offer an explicit migration to the corrected distributed profile using
unit/integration/e2e for requirements and build/unit/integration/e2e for
correctness and validate. Preserve criterion identity and requiredness, leave
sealed historical bundles immutable, and flag incompatible criterion classes
at profile validation time. Never fabricate references or waive required checks.

The instance adopted the existing corrected distribution profile's class mapping.
No upstream issue was submitted.

## Follow-up on 2026-09-09

POSE 2.0.2 also rejects the legacy backend/frontend overlays before review
preparation: contract, test, observability and validation are not registered
check evidence classes. Reproduced with the installed overlays and compared
against a fresh 2.0.2 installation in a temporary directory. Adopted its exact
overlay mappings to unit/integration/e2e/build, preserving criterion IDs,
required tools and sealed historical bundles. Recommend migrating overlays
together with the base profile. No upstream submission was made.

## Recheck after instance update — 2026-09-09

Both the installed binary and instance now report POSE 2.0.2. Strict structural
validation passes, and review bundle preparation returns prepared with no
blockers or unsupported evidence-class errors. The updated review policy has
evidence_vocabulary_reconciled_at, and the milestone profile now requests
integration instead of integration-test.

The original failure occurred with legacy profiles. Backend/frontend mappings
had already been corrected manually before this update, so this recheck does
not independently prove automatic migration of those original profiles.
There is no current reproduction warranting upstream submission. Preserve the
earlier observations as migration history rather than an unresolved defect.
