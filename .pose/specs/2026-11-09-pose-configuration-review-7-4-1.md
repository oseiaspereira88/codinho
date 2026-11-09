---
slug: pose-configuration-review-7-4-1
status: draft
created_at: 2026-11-09
completed_at:
supersedes:
depends_on:
remediates:
priority: 1
components:
task_type: feature
surface: minimal
delivers:
changelog: none
---

# Spec: Review the capabilities POSE 7.4.1 brings

## 1. Intent

### Goal

Decide, for each capability this project has not decided under POSE 7.4.1, whether to adopt, decline or defer it. Nothing changes until the maintainer answers.

### Business value

An update brings capabilities; deciding each one explicitly keeps none in limbo and records why it is on or off.

### Constraints

Each decision is an action request to the `maintainer` role (`pose action list`, `pose setup`). An answer is applied with `pose adopt --request <act-id> --apply` or by `pose setup`; adopting dates the capability with the day it is applied, so earlier work is not re-judged.

### Non-goals

Changing anything the maintainer did not answer.

## 2. Requirements

### Functional

- R1: The maintainer shall adopt, decline or defer `agency-readiness` — an action request that is not satisfied refuses the start, close or release it restricts (since 7.0.0; recommended: adopt, as new instances do).
- R2: The maintainer shall adopt, decline or defer `atomic-start` — a spec moves from draft to in-progress with `pose start --apply`, which records the baseline it is reconciled against (since 6.0.0; recommended: adopt, as new instances do).
- R3: The maintainer shall adopt, decline or defer `causality-closeout` — a review that answers for observed structure maps each material fact to a requirement or decision, with its own reason (since 6.0.0; recommended: adopt, as new instances do).
- R4: The maintainer shall adopt, decline or defer `contract-nodes` — a change to a started spec's requirements, assumptions or decisions is recorded with `pose amend` (since 6.0.0; recommended: adopt, as new instances do).
- R5: The maintainer shall adopt, decline or defer `definition-of-ready` — a spec created on or after the date cannot move to in-progress until Intent, Requirements with stable ids and Technical Plan are filled (since 4.0.0; recommended: adopt, as new instances do).
- R6: The maintainer shall adopt, decline or defer `overlay:engineering-judgment` — a scope selected by its declared risk is also reviewed for engineering judgment (since 6.0.0; no recommendation).
- R7: The maintainer shall adopt, decline or defer `overlay:high-criticality-review` — a scope touching a critical component carries the stricter review its criticality asks for (since 6.0.0; no recommendation).
- R8: The maintainer shall adopt, decline or defer `overlay:structural-materiality` — a spec that changes dependencies, components, public or governance contracts owes a structural review (since 6.0.0; no recommendation).
- R9: The maintainer shall adopt, decline or defer `qualified-artifact-refs` — a spec, roadmap or milestone of another project is referenced as `xref:<project>/…` and resolved through the bound project (since 6.0.0; no recommendation).
- R10: The maintainer shall adopt, decline or defer `spec-authority-transfer` — a spec's authority moves to another project with `pose spec-transfer`, leaving a redirect every reference follows (since 6.0.0; no recommendation).

### Non-functional

- None.

### Security

- None.

### Compatibility

- None.

## 3. Technical Plan

### Affected areas

POSE policy only.

### Artifacts

- created: .pose/specs/2026-11-09-pose-configuration-review-7-4-1.md
- modified: .pose/policy/adoption-decisions.json

### Technical risks

- None.

## 6. Validation

### Strategy

`pose setup` lists no capability decision pending.

### Deterministic checks

#### Health
- Command: `pose doctor`
- Expected: no `setup.capabilities` step

### Requirement trace

### Known gaps

## 7. Final Report

### Delivered scope

Not started.

### Residual risks

### Follow-ups
