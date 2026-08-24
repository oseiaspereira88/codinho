---
name: agent-batch-review
description: Runs an independent-agent pre-review (or authoring) round for a batch of work using scripts/agent-review.sh — non-interactive, sandboxed, session-resumed across rounds. Use when authoring catalog content batches or any deliverable that benefits from an adversarial second-agent pass before human review.
---

# Skill: agent-batch-review

Runs one round of the author/reviewer loop documented in
[`docs/agent-review-workflow.md`](../../../docs/agent-review-workflow.md),
using [`scripts/agent-review.sh`](../../../scripts/agent-review.sh) as the
invocation primitive. Read that doc in full before using this skill — it
has the prompt templates and the reasoning this skill assumes.

## When to use

- After authoring a batch of catalog content (challenges, concepts,
  competencies) and before handing it to the human reviewer.
- Any other deliverable where an independent agent's adversarial read
  catches objective problems (leaked solutions, unverifiable acceptance
  criteria, compound instructions, unbacked checks) cheaper than a human
  review round would.

## Steps

1. Identify the round: first round of this batch (`new`) or a follow-up
   after fixes (`resume`). Never `new` twice for the same batch — that
   throws away context and re-reads files unnecessarily (see the doc's
   "why resume" section).
2. Pick the backend agent (`codex` today; extend `scripts/agent-review.sh`
   for another CLI before picking it here).
3. Write the prompt from the doc's reviewer (or author) template, filling
   in: repo context, exactly which files/IDs changed, what to read first,
   which deterministic command proves structural validity, and — for a
   `resume` round — exactly what was fixed since the last verdict (never
   re-paste file contents the agent already read).
4. Run:
   ```
   scripts/agent-review.sh new codex /path/to/round-N.md "<prompt>"
   # or, for a follow-up:
   scripts/agent-review.sh resume codex /path/to/round-N.md "<prompt>"
   ```
5. Read the verdict. If rejected or "aprovado com ressalvas", apply the
   fixes yourself (the reviewer never edits files), re-validate
   deterministically, then loop back to step 3 with `resume`.
6. Stop when the verdict is "aprovado sem ressalvas" for every item in the
   batch, or when the user decides to proceed with known, accepted
   ressalvas. Record the round-by-round verdicts in the governing spec's
   Execution log (see `go-foundations-packs` Decision 3 / Execution log for
   the reference shape).

## What this never does

Never treat the reviewer's verdict as satisfying a project's human-review
gate (e.g. `catalog-authoring-quality`'s `reviewed_by`/`playtested`). This
skill only reduces what a human has left to check — it does not replace
them. State that explicitly in every reviewer prompt.
