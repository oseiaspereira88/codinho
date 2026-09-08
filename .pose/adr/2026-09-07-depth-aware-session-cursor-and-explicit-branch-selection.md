# ADR: Depth-aware session cursor and explicit branch selection

## Status
Accepted — 2026-09-07.

## Context

Start ignores depth, traversal stops at the first layer and changing depth
recreates progress. Ordered children have been mistaken for alternatives.
Session recovery and evidence lineage require stable node identity.

## Decision

- Treat challenge/layer as synthetic instructional windows: challenge title/brief,
  layer ID as objective and challenge title as scope. Never concatenate descendant
  instructions, constraints or solutions. Macro/meso/micro use authored content.
- Traverse layers and ordinary children in document order at the requested depth.
  A completed coarse window covers its subtree for navigation, without synthetic
  child completions, evaluations or mastery events. Completing all finer windows
  exhausts their container without creating a parent completion event. Do not
  reopen an exhausted ancestor after a depth change. A completed choice
  checkpoint redirects to its selected branch when refining/coarsening; it
  does not substitute for completing that branch.
- Add optional children_mode: sequence (default) or choice to authored steps.
  Choice requires at least two children. Node IDs must be nonempty and unique
  within each challenge, including layer IDs; steps use macro/meso/micro kinds. When a finer window reaches an unselected
  choice, expose its parent as a decision checkpoint. Complete that checkpoint
  separately; step_advance then offers direct child IDs. next_step_id selects one
  offered child and activates its first window at the requested depth. Persist the
  choice; unselected alternatives are excluded, not marked completed.
- At a depth that displays the choice parent itself, ordinary completion covers
  that whole window. No lower-level choice is necessary for a coarse exercise.
- Preserve a progress/clean-evaluation record per visited node plus a cursor for
  the deepest current position. Coarsening shows the containing ancestor; refining
  returns to that cursor. Completed windows cannot be refined to generate more
  credit; use explicit advance. Depth changes never advance to a sibling.
- Keep one active instruction. Pause/finish and an open detour prevent navigation.
  Advance, including options and exhaustion, requires completion or an explicit
  override and the current revision. Invalid choices never append. An override
  records skipped coverage without creating completion or autonomy evidence.
- Record the first exhaustion acknowledgment, including last-window overrides,
  as step_advanced with done: true; do not infer session_finish. Subsequent
  exhaustion reads do not append. Prevent changes of window after exhaustion.
- Repeat completion of an already completed node is a no-op at the current
  revision, preventing duplicate completion/mastery events. Historical request
  retries retain their original response before present-state validation.
- Add navigation_version: 1 to new start/navigation/completion events. Replay
  older starts and transitions with their recorded historical semantics; new
  operations use the new cursor. Omit the new optional catalog JSON field when
  empty so old pinned-content digests remain valid. Unknown navigation versions
  fail explicitly; do not rewrite logs or silently interpret future versions.

## Consequences

Session/application/MCP contracts change additively with next_step_id and actual
node kinds. Default sibling behavior intentionally becomes ordered sequence;
consumers must author choice explicitly and must not assume macro completion
will require repeating all micro instructions. Keep explicit feedback,
evaluation, completion, advance and session_finish as distinct operations.

Rejected alternatives: preorder at all depths repeats the same exercise at each
scale; inferring alternatives from sibling count skips required work; resetting
progress destroys audit history; a mutable snapshot store duplicates JSONL state.

Review this decision before multi-challenge tracks or if authored choices require
rejoining non-tree DAGs. Those are separate contracts; this traversal supports
nested ordered trees with exclusive child choices only.
