# Installed review profile requests unproducible evidence classes

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
