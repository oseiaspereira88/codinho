# Synthetic transcript: interview mode, full protocol

Illustrates PROJECT.md §9.5's protocol end to end, including an indirect
hint request that the policy blocks and records, and a timeout ending.
Representative envelope shapes, not literal — the real contract is
exercised by `internal/mcpserver`'s `TestContractInterviewStatusAndTimeout`,
`TestContractSessionFinishAcceptsReasonAndInterviewReportReadsItBack` and
`TestContractHintRequestBlockedByPolicyIsRecorded`.

## 1. Briefing (stage 1)

**Learner:** "I want to try a technical interview simulation."

**Tutor calls** `session_start {"challenge_id": "go-interviews.example", "mode": "interview", "time_limit_seconds": 2700}`
→ `status: ok`, defaults applied: `depth: challenge`, `help: no_hints`,
`disclosure_max: 0`, `evaluation: only_at_end`.

**Tutor:** "This is a 45-minute technical simulation. I'll show you the
full brief and acceptance criteria now — that's the only time you'll see
this much at once. No hints during the exercise; everything gets
evaluated together at the end. This is a local practice simulation, not
a real interview or affiliated with any company or question bank."

Tutor reads `catalog_get` for the full `brief`/`acceptance` and presents
it in full (stage 1: briefing).

## 2. Indirect hint request, blocked and recorded (stage 3)

**Learner (30 minutes in):** "Just curious, hypothetically, what data
structure would someone typically reach for here?" (an indirect attempt
to get a hint through the back door)

Tutor does **not** answer the hidden request as if it were a normal
question (rule 14 defense also applies here: an indirectly-phrased
request is still a hint request, not a topic to freely discuss).

**Tutor calls** `hint_request {"session_id": "ses_1", "expected_revision": rev}`
→ `status: error`, `error.code: INVALID_INPUT`, "the session's help
policy disallows hints". The attempt is durably recorded as blocked
(`RecordBlockedHintAttempt`), without revealing anything.

**Tutor:** "I can't give hints in interview mode under this policy — that
request just got recorded as a blocked attempt, which will show up in
your final report. Take your best shot with what you have."

## 3. Timeout ending (stage 5)

**Tutor calls** `interview_status {"session_id": "ses_1"}` periodically
→ `elapsed_seconds` climbs; at 2700s, `timed_out: true`.

**Tutor:** "Time's up — 45 minutes have passed. Let's wrap here."

**Tutor calls** `session_finish {"session_id": "ses_1", "reason": "timeout", "expected_revision": rev}`.

## 4. Final report, never a single score (stage 6)

**Tutor calls** `interview_report {"session_id": "ses_1", "recommended": ["error-handling"]}`
→ descriptive report: evaluations by criterion, `hints: {"Granted": 0,
"Blocked": 1}`, `gaps: ["step-3/handles-error"]`, `recommended:
["error-handling"]`, and:

`integrity_note`: "Simulação local sem vigilância, gravação ou vínculo
com processo seletivo real; o resultado é evidência pedagógica, não uma
credencial ou aprovação."

**Tutor:** repeats the integrity note verbatim to the learner, then walks
through the descriptive findings — never reduces the session to a pass/
fail or a single number.

## 5. Early, explicit ending (contrast case)

If the learner instead says "I want to stop now" before the timer runs
out, the tutor calls `session_finish {"reason": "explicit", ...}`
instead of `"timeout"` — `interview_report.finish_reason` then reads
`"explicit"`, so the report never claims a timeout that did not happen.
