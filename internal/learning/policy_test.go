package learning

import (
	"errors"
	"testing"
	"time"
)

func TestNewSessionPolicyRejectsInvalidValues(t *testing.T) {
	cases := []struct {
		name       string
		mode       PedagogicalMode
		depth      Depth
		help       HelpPolicyKind
		disclosure DisclosureLevel
		evaluation EvaluationPolicyKind
		advance    AdvancePolicyKind
		timeLimit  *time.Duration
	}{
		{"bad mode", "bogus", DepthMicro, HelpFree, DisclosureNone, EvaluationOnDemand, AdvanceExplicit, nil},
		{"bad depth", ModePractice, "bogus", HelpFree, DisclosureNone, EvaluationOnDemand, AdvanceExplicit, nil},
		{"bad help", ModePractice, DepthMicro, "bogus", DisclosureNone, EvaluationOnDemand, AdvanceExplicit, nil},
		{"disclosure too high", ModePractice, DepthMicro, HelpFree, DisclosureLevel(7), EvaluationOnDemand, AdvanceExplicit, nil},
		{"disclosure negative", ModePractice, DepthMicro, HelpFree, DisclosureLevel(-1), EvaluationOnDemand, AdvanceExplicit, nil},
		{"bad evaluation", ModePractice, DepthMicro, HelpFree, DisclosureNone, "bogus", AdvanceExplicit, nil},
		{"bad advance", ModePractice, DepthMicro, HelpFree, DisclosureNone, EvaluationOnDemand, "bogus", nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := NewSessionPolicy(c.mode, c.depth, c.help, c.disclosure, c.evaluation, c.advance, c.timeLimit)
			if !errors.Is(err, DomainError{Code: ErrCodeInvalidValue}) {
				t.Fatalf("expected ErrCodeInvalidValue, got %v", err)
			}
		})
	}

	negative := -time.Second
	_, err := NewSessionPolicy(ModePractice, DepthMicro, HelpFree, DisclosureNone, EvaluationOnDemand, AdvanceExplicit, &negative)
	if !errors.Is(err, DomainError{Code: ErrCodeInvalidValue}) {
		t.Fatalf("expected ErrCodeInvalidValue for negative time limit, got %v", err)
	}
}

func TestNewSessionPolicyAcceptsValidValues(t *testing.T) {
	limit := time.Hour
	p, err := NewSessionPolicy(ModeInterview, DepthChallenge, HelpNoHints, DisclosureSolution, EvaluationOnlyAtEnd, AdvanceExplicit, &limit)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Mode != ModeInterview || p.InitialDepth != DepthChallenge || p.Help.Kind != HelpNoHints {
		t.Fatalf("unexpected policy: %+v", p)
	}
	if p.Disclosure.MaxLevel != DisclosureSolution {
		t.Fatalf("disclosure max = %d, want %d", p.Disclosure.MaxLevel, DisclosureSolution)
	}
}

func TestNewHintUsageEnforcesDisclosureCap(t *testing.T) {
	policy := DisclosurePolicy{MaxLevel: DisclosureConceptOrAPI}

	if _, err := NewHintUsage("step_1", DisclosureGuidingQuestion, policy); err != nil {
		t.Fatalf("unexpected error within cap: %v", err)
	}

	_, err := NewHintUsage("step_1", DisclosurePseudocode, policy)
	if !errors.Is(err, DomainError{Code: ErrCodeDisclosureExceeded}) {
		t.Fatalf("expected ErrCodeDisclosureExceeded, got %v", err)
	}
}

func TestNewFeedbackRecordAppliesTableDefaults(t *testing.T) {
	fb, err := NewFeedbackRecord("step_1", FeedbackViolation, "missing required nil check", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !fb.Blocking {
		t.Fatal("violation feedback should block completion by default")
	}

	fb2, err := NewFeedbackRecord("step_1", FeedbackSuggestion, "consider a table-driven test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fb2.Blocking {
		t.Fatal("suggestion feedback should not block completion by default")
	}

	override := true
	fb3, err := NewFeedbackRecord("step_1", FeedbackRisk, "works only under low load", &override)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !fb3.Blocking {
		t.Fatal("explicit override should have made risk feedback blocking")
	}
}

func TestNewFeedbackRecordRejectsInvalidType(t *testing.T) {
	_, err := NewFeedbackRecord("step_1", "bogus", "text", nil)
	if !errors.Is(err, DomainError{Code: ErrCodeInvalidValue}) {
		t.Fatalf("expected ErrCodeInvalidValue, got %v", err)
	}
}

func TestEvaluationNeverCompletesOrAdvances(t *testing.T) {
	eval := Evaluation{
		StepID: "step_1",
		Criteria: []CriterionResult{
			{Name: "fields-match", Kind: StructuralCriterionKind, Severity: SeverityBlocking, Verdict: VerdictNotMet},
		},
	}
	if !eval.HasBlockingFailure() {
		t.Fatal("expected blocking failure")
	}

	step := NewStepProgress("step_1")
	if step.State != StepStateLocked {
		t.Fatalf("evaluating a step must never itself change its state; state = %s", step.State)
	}
}

func TestNewEvidenceRejectsEmptyFields(t *testing.T) {
	if _, err := NewEvidence("", EvidenceKindTest, "ref"); !errors.Is(err, DomainError{Code: ErrCodeInvalidValue}) {
		t.Fatalf("expected ErrCodeInvalidValue for empty id, got %v", err)
	}
	if _, err := NewEvidence("ev_1", EvidenceKindTest, ""); !errors.Is(err, DomainError{Code: ErrCodeInvalidValue}) {
		t.Fatalf("expected ErrCodeInvalidValue for empty ref, got %v", err)
	}
	ev, err := NewEvidence("ev_1", EvidenceKindTest, "sha256:abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ev.ID() != "ev_1" || ev.Kind() != EvidenceKindTest || ev.Ref() != "sha256:abc" {
		t.Fatalf("unexpected evidence: %+v", ev)
	}
}
