package session

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/oseiaspereira88/codinho/internal/assessment"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

// newPolicyTestService's fixture has three sibling micro steps at the
// root: a plain one with no completion policy, one requiring a positive
// evaluation, and one requiring user confirmation — plus a fourth that
// branches into two children, so step_advance's every outcome (descend,
// bubble to sibling, branch, exhaust) has real content to walk.
func newPolicyTestService(t *testing.T) *Service {
	t.Helper()
	dir := t.TempDir()
	writeTestFile(t, dir, "manifest.yaml", "schema_version: 1\npacks:\n  - pack.yaml\n")
	writeTestFile(t, dir, "pack.yaml", `schema_version: 1
id: fixture-pack
version: 1.0.0
challenges:
  - schema_version: 1
    id: fixture.challenge-one
    version: 1.0.0
    title: Fixture challenge
    kind: atomic
    difficulty: foundational
    layers:
      - id: understanding
        macro_steps:
          - id: fixture.step-plain
            kind: micro
            title: Plain step
            instruction:
              objective: Plain objective.
              scope: Plain scope.
          - id: fixture.step-needs-eval
            kind: micro
            title: Needs positive evaluation
            instruction:
              objective: Eval objective.
              scope: Eval scope.
            completion:
              requires_positive_evaluation: true
          - id: fixture.step-needs-confirm
            kind: micro
            title: Needs user confirmation
            instruction:
              objective: Confirm objective.
              scope: Confirm scope.
            completion:
              requires_user_confirmation: true
          - id: fixture.step-branch
            kind: micro
            title: Branches
            instruction:
              objective: Branch objective.
              scope: Branch scope.
            children:
              - id: fixture.branch-a
                kind: micro
                title: Branch A
              - id: fixture.branch-b
                kind: micro
                title: Branch B
`)
	catalog, diags, err := curriculum.Load(dir, curriculum.DefaultLimits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostics: %+v", diags)
	}
	store, err := eventstore.Open(filepath.Join(t.TempDir(), "events.jsonl"), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return New(catalog, store, nil)
}

func TestFeedbackPrepareAssemblesPacketWithoutMutating(t *testing.T) {
	svc := newPolicyTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	packet, err := svc.FeedbackPrepare(start.SessionID, "why exported fields?")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if packet.Objective != "Plain objective." || packet.Question != "why exported fields?" {
		t.Fatalf("unexpected packet: %+v", packet)
	}
	if got := svc.store.Revision(string(start.SessionID)); got != start.Revision {
		t.Fatalf("feedback_prepare must not mutate the stream, revision = %d, want %d", got, start.Revision)
	}
}

func TestFeedbackRecordNeverCompletesOrAdvances(t *testing.T) {
	svc := newPolicyTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := svc.FeedbackRecord(start.SessionID, learning.FeedbackViolation, "this violates the constraint", nil, start.Revision, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Revision != start.Revision+1 {
		t.Fatalf("revision = %d, want %d", result.Revision, start.Revision+1)
	}
	get, err := svc.Get(start.SessionID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if get.ActiveStep != start.ActiveStep {
		t.Fatal("feedback must never change the active step")
	}
}

func TestStepEvaluateThenCompletePlainStep(t *testing.T) {
	svc := newPolicyTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Completing before evaluating must fail: the step is still Active. The
	// event is still appended before the domain check runs (this
	// codebase's established append-first-then-validate ordering), so the
	// stream's revision advances even on this rejected call.
	if _, err := svc.StepComplete(start.SessionID, false, false, start.Revision, ""); !errors.As(err, new(learning.DomainError)) {
		t.Fatalf("expected a domain error completing an unevaluated step, got %v", err)
	}
	rev := svc.store.Revision(string(start.SessionID))

	eval, err := svc.StepEvaluate(EvaluateInput{SessionID: start.SessionID, ExpectedRevision: rev})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if eval.HasBlockingFailure {
		t.Fatal("no criteria were submitted; there should be no blocking failure")
	}

	complete, err := svc.StepComplete(start.SessionID, false, false, eval.Revision, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if complete.Revision != eval.Revision+1 {
		t.Fatalf("revision = %d, want %d", complete.Revision, eval.Revision+1)
	}
}

func TestStepEvaluateStructuralCriterionDrivesBlockingFailure(t *testing.T) {
	svc := newPolicyTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	eval, err := svc.StepEvaluate(EvaluateInput{
		SessionID: start.SessionID,
		Criteria: []assessment.CriterionInput{
			{Name: "compiles", Kind: learning.StructuralCriterionKind, Severity: learning.SeverityBlocking},
		},
		ExpectedRevision: start.Revision,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !eval.HasBlockingFailure {
		t.Fatal("a structural criterion with no cited evidence must resolve to unverifiable, which blocks completion")
	}
	if eval.Criteria[0].Verdict != learning.VerdictUnverifiable {
		t.Fatalf("verdict = %s, want unverifiable", eval.Criteria[0].Verdict)
	}
}

func TestStepEvaluateWithSubmissionIntentRecordsAttempt(t *testing.T) {
	svc := newPolicyTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	eval, err := svc.StepEvaluate(EvaluateInput{SessionID: start.SessionID, SubmissionIntent: true, ExpectedRevision: start.Revision})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !eval.AttemptRecorded {
		t.Fatal("expected AttemptRecorded when submission_intent is true")
	}
	if eval.Revision != start.Revision+2 {
		t.Fatalf("revision = %d, want %d (evaluation + attempt)", eval.Revision, start.Revision+2)
	}
}

func TestStepEvaluateIsIdempotentByRequestIDIncludingAttempt(t *testing.T) {
	svc := newPolicyTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	in := EvaluateInput{SessionID: start.SessionID, SubmissionIntent: true, ExpectedRevision: start.Revision, RequestID: "eval-retry"}

	first, err := svc.StepEvaluate(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := svc.StepEvaluate(in)
	if err != nil {
		t.Fatalf("unexpected error on retry: %v", err)
	}
	if first.Revision != second.Revision || !second.AttemptRecorded {
		t.Fatalf("retry was not idempotent: %+v vs %+v", first, second)
	}
}

func TestStepCompleteRequiresPositiveEvaluationPolicy(t *testing.T) {
	svc := newPolicyTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	eval, err := svc.StepEvaluate(EvaluateInput{SessionID: start.SessionID, ExpectedRevision: start.Revision})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	complete, err := svc.StepComplete(start.SessionID, false, false, eval.Revision, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	advance, err := svc.StepAdvance(start.SessionID, false, complete.Revision, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if advance.StepID != "fixture.step-needs-eval" {
		t.Fatalf("StepID = %s, want fixture.step-needs-eval", advance.StepID)
	}

	failedEval, err := svc.StepEvaluate(EvaluateInput{
		SessionID:        start.SessionID,
		Criteria:         []assessment.CriterionInput{{Name: "compiles", Kind: learning.StructuralCriterionKind, Severity: learning.SeverityBlocking}},
		ExpectedRevision: advance.Revision,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !failedEval.HasBlockingFailure {
		t.Fatal("expected a blocking failure with no evidence cited")
	}

	if _, err := svc.StepComplete(start.SessionID, false, false, failedEval.Revision, ""); !errors.Is(err, learning.DomainError{Code: learning.ErrCodeCompletionPolicyNotMet}) {
		t.Fatalf("expected ErrCodeCompletionPolicyNotMet, got %v", err)
	}

	overridden, err := svc.StepComplete(start.SessionID, false, true, failedEval.Revision, "")
	if err != nil {
		t.Fatalf("unexpected error with override: %v", err)
	}
	if overridden.Revision != failedEval.Revision+1 {
		t.Fatalf("revision = %d, want %d", overridden.Revision, failedEval.Revision+1)
	}
}

func TestStepCompleteRequiresUserConfirmationPolicy(t *testing.T) {
	svc := newPolicyTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	eval, err := svc.StepEvaluate(EvaluateInput{SessionID: start.SessionID, ExpectedRevision: start.Revision})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	complete, err := svc.StepComplete(start.SessionID, false, false, eval.Revision, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	advance, err := svc.StepAdvance(start.SessionID, false, complete.Revision, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	advance2, err := svc.StepAdvance(start.SessionID, true, advance.Revision, "") // step-needs-eval isn't completed; override to reach step-needs-confirm
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if advance2.StepID != "fixture.step-needs-confirm" {
		t.Fatalf("StepID = %s, want fixture.step-needs-confirm", advance2.StepID)
	}

	eval2, err := svc.StepEvaluate(EvaluateInput{SessionID: start.SessionID, ExpectedRevision: advance2.Revision})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := svc.StepComplete(start.SessionID, false, false, eval2.Revision, ""); !errors.Is(err, learning.DomainError{Code: learning.ErrCodeCompletionPolicyNotMet}) {
		t.Fatalf("expected ErrCodeCompletionPolicyNotMet without confirmation, got %v", err)
	}
	if _, err := svc.StepComplete(start.SessionID, true, false, eval2.Revision, ""); err != nil {
		t.Fatalf("unexpected error with confirm=true: %v", err)
	}
}

func TestStepAdvanceReportsBranchesWithoutMutating(t *testing.T) {
	svc := newPolicyTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Walk (with overrides, since these steps aren't actually completed)
	// straight to fixture.step-branch.
	rev := start.Revision
	for range 3 {
		res, err := svc.StepAdvance(start.SessionID, true, rev, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		rev = res.Revision
	}
	get, err := svc.Get(start.SessionID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if get.ActiveStep != "fixture.step-branch" {
		t.Fatalf("active step = %s, want fixture.step-branch", get.ActiveStep)
	}

	result, err := svc.StepAdvance(start.SessionID, true, rev, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Branches) != 2 || result.Branches[0].StepID != "fixture.branch-a" || result.Branches[1].StepID != "fixture.branch-b" {
		t.Fatalf("unexpected branches: %+v", result)
	}
	if got := svc.store.Revision(string(start.SessionID)); got != rev {
		t.Fatalf("reporting branches must not append an event, revision = %d, want %d", got, rev)
	}
}

func TestStepAdvanceReachesTheEndOfTheTree(t *testing.T) {
	svc := newTestService(t) // nested macro -> meso -> micro, single children throughout
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rev := start.Revision
	for range 2 {
		res, err := svc.StepAdvance(start.SessionID, true, rev, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		rev = res.Revision
	}
	get, err := svc.Get(start.SessionID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if get.ActiveStep != "fixture.micro-one" {
		t.Fatalf("active step = %s, want fixture.micro-one", get.ActiveStep)
	}

	result, err := svc.StepAdvance(start.SessionID, true, rev, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Done {
		t.Fatalf("expected Done at the end of the tree, got %+v", result)
	}
}

func TestReflectionRecordDoesNotChangeStepState(t *testing.T) {
	svc := newPolicyTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := svc.ReflectionRecord(start.SessionID, "slice-filter", "why a pointer?", "because it avoids copying", "clear reasoning", start.Revision, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Revision != start.Revision+1 {
		t.Fatalf("revision = %d, want %d", result.Revision, start.Revision+1)
	}
	get, err := svc.Get(start.SessionID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if get.ActiveStep != start.ActiveStep {
		t.Fatal("reflection must never change the active step")
	}
}

func TestReflectionRecordRequiresAnswer(t *testing.T) {
	svc := newPolicyTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := svc.ReflectionRecord(start.SessionID, "", "why?", "", "", start.Revision, ""); !errors.Is(err, learning.DomainError{Code: learning.ErrCodeInvalidValue}) {
		t.Fatalf("expected ErrCodeInvalidValue, got %v", err)
	}
}
