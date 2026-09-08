---
title: "Clarify the validation umbrella evidence class in native review plans"
type: enhancement
created_at: "2026-09-08T03:44:39Z"
status: staged
privacy: sanitized-synthetic
---

# Observation

POSE 1.7.12 review verification can emit:
`review tool validate drops evidence class validation: no registered check may emit it`.
The same scope can be fresh, approved and ready to close with a passed structured
validation report. The diagnostic does not explain whether validation is an
umbrella report class or must be emitted by an individual matrix check.

# Synthetic reproduction

Use one root Go module with build, unit, integration and reachability checks,
an attributed capability spec targeting cmd/example, and the native spec-closeout
profile. Run strict validation to a structured report, refresh indexes, prepare
a review bundle and inspect its warnings. Record an independent attestation and
verify it. Observe the dropped-class warning alongside successful verification.
No project code or private configuration is needed to reproduce the ambiguity.

# Proposed improvement

Document the distinction between a structured validation report and individual
check evidence classes. If validation is an umbrella class, preserve its native
tool capability without requiring a synthetic matrix check. Otherwise explain
which concrete evidence classes satisfy the native profile and how to resolve
the warning. Add a regression covering the documented interpretation.

This is a diagnostic clarity proposal, not a claim that passed checks failed.
No gate was bypassed and no upstream submission was made.
