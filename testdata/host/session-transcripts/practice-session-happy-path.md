# Synthetic transcript: practice session, happy path

Illustrates the tool routing order in `mcp-tool-routing.md`, with
representative (not literal) envelope shapes matching the real MCP
contract. This is documentation, not a Go test — the actual protocol
shapes are exercised by `internal/mcpserver`'s own contract tests and by
the tutor-skill-host-integration smoke test run against the real
`codinho serve` binary over `mcp.CommandTransport`.

## 1. Discovery

**Learner:** "I want to practice filtering slices in Go."

**Tutor calls** `catalog_search {"text": "filtrar", "competency": "slice-filter"}`
→ `status: ok`, one match: `go-data.slice-filter-preserve-input`.

**Tutor:** "Found it — 'Filtrar valores sem alterar a entrada', about 15
minutes, foundational difficulty. Want to start?"

## 2. Start and announce (rule 1)

**Tutor calls** `session_start {"challenge_id": "go-data.slice-filter-preserve-input"}`
→ `status: ok, session_id: "ses_1", revision: 1, active_node: {id: "model.declare-user-struct"}`.

**Tutor:** "Starting in practice mode, micro depth: small steps with
frequent feedback. First step: declare the `User` struct as specified."

## 3. Observe before evaluate (rule 6)

Learner writes code. Tutor calls `workspace_observe` with the learner's
real workspace root — first call establishes the baseline
(`progress_effect: none`, `baseline: true`).

## 4. Evaluate only with explicit intent (rule 7)

**Learner:** "Okay, I think that's right, check it."

**Tutor calls** `step_evaluate {"submission_intent": true, "criteria": [...]}`
→ `progress_effect: attempt_recorded`, criterion verdict `met`.

**Tutor:** "That records an attempt — the struct declaration checks out."
(rule 9: explaining progress_effect)

## 5. Consultative feedback (rule 8)

**Tutor calls** `feedback_prepare` then `feedback_record {"type": "idiom", ...}`
with an observation about field naming — `progress_effect:
feedback_recorded`, never `step_completed`.

## 6. Complete and advance are separate (invariant 5)

**Tutor calls** `step_complete {}` → `progress_effect: step_completed`.

**Learner:** "Actually, let me try a cleaner version first."

**Tutor:** does **not** call `step_advance` (rule 15 — respects the
request to stay on the step).

Later, learner: "Okay, ready to move on."

**Tutor calls** `step_advance {}` → `progress_effect: step_advanced`,
new `active_node`.

## 7. Domain evidence (only from real interaction)

**Tutor calls** `mastery_evidence_record {"competency_id":
"slice-filter", "dimension": "autonomous_implementation", "evidence_id":
"<the step_evaluate criterion's evidence_id>", "success": true}`
→ `state: demonstrates_without_help`.

## 8. Session end

**Learner:** "That's enough for today."

**Tutor calls** `session_finish {}`, then fills
`assets/session-summary-template.md` from data returned by the actual
tool calls above — never from memory or invention.
