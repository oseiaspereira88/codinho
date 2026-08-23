# Synthetic transcript: debug mode, six-stage protocol

Illustrates the six-stage debugging protocol (`internal/session.DebugProtocolStages`,
PROJECT.md §8.1 "Depuração") and granularity auto-adaptation with a
mandatory explanation (PROJECT.md §8.5). Representative envelope shapes,
not literal — the real tool contract is exercised by
`internal/mcpserver`'s contract tests and by
`cmd/codinho/cli_integration_test.go`'s real, authored `go-debug.slice-off-by-one`
challenge.

## 1. Prepare the fixture (never the tutor's own code)

**Learner:** "I have a bug in my slice code, help me find it."

Tutor tells the learner to run, in their own terminal:
`codinho workspace prepare go-debug.slice-off-by-one --dest ./practice`
The CLI materializes `main.go` (with the bug) and `manifest.json` — the
tutor never writes this file itself (Constraint: "workspace prepare é
acionado diretamente pelo usuário, não por tool MCP").

## 2. Start in debug mode, announce the protocol (rule 1)

**Tutor calls** `session_start {"challenge_id": "go-debug.slice-off-by-one", "mode": "debug"}`
→ `status: ok`, defaults applied: `depth: macro`, `help: progressive`,
`disclosure_max: 3` (logical outline — never past it, so the cause is
never handed over before a hypothesis).

**Tutor:** "Debug mode: we'll go stage by stage — reproduce, locate the
first divergence, hypothesize, observe, confirm or reject, then the
smallest fix. Let's reproduce it first: run the program and tell me what
you see."

## 3. reproduce → locate_first_divergence → hypothesize

**Learner:** "It panics: index out of range."

Tutor asks the learner to point at the exact line (`locate_first_divergence`),
then asks: "why do you think that index goes out of bounds?"
(`hypothesize` — the learner states the hypothesis, never the tutor;
PROJECT.md §8.2: "sem revelar a causa antes da hipótese").

## 4. observe → confirm_or_reject (rule 6: observe before evaluate)

**Tutor calls** `workspace_observe` with the learner's real workspace
root, then, once the learner submits deliberately,
`step_evaluate {"submission_intent": true, ...}` → criterion `met` or
`not_met`, confirming or rejecting the stated hypothesis with real
evidence.

## 5. apply_smallest_fix, done by the learner

**Learner** edits `main.go` themselves (rule 3: the tutor never edits
learner files). Once fixed, `step_complete {}` →
`progress_effect: step_completed`.

## 6. Granularity auto-adaptation, always explained (PROJECT.md §8.5)

After three debug challenges completed without help,
`progress_get {"competency_id": "debug-slice-bounds"}` shows
`state: demonstrates_without_help`, `evidence_count: 3`. Per the ladder in
`mcp-tool-routing.md` (`challenge → layer → macro → meso → micro`), that
many repeated-ease signals at the current `macro` depth means widening one
step, to `layer`.

**Tutor calls** `granularity_adjust {"depth": "layer", "reason": "3
depurações seguidas sem ajuda: ampliando para camada inteira", ...}`.

**Tutor:** "You've solved the last three bugs without a hint, so I'm
widening the window to the whole layer — you'll get less interruption
between checks. Manual adjustment always wins if you want it back."
(the explanation is mandatory, not optional — PROJECT.md §8.5: "deve
ser informado ao aluno").

## 7. Learner proposes the next step (rule 16, PROJECT.md §21.5)

**Learner:** "I think the next thing to try is the off-by-one in the
loop bound, not just the index."

**Tutor calls** `learner_next_step_propose {"step_id": "go-debug.apply-smallest-fix", ...}`
→ records the autonomy signal; the tutor still evaluates the proposal
before deciding whether to call `step_advance` — proposing never
advances anything on its own.
