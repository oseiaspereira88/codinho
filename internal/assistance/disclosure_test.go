package assistance

import (
	"errors"
	"testing"

	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

func fixtureStep() curriculum.StepAuthoring {
	return curriculum.StepAuthoring{
		ID: "model.declare-user-struct",
		Hints: []curriculum.HintAuthoring{
			{Level: 1, Kind: "guiding_question"},
			{Level: 2, Kind: KindSyntaxRecall},
			{Level: 3, Kind: "logical_prose"},
		},
	}
}

func TestNextHintClimbsOneLevelAndResolvesKind(t *testing.T) {
	step := fixtureStep()

	level, kind, err := NextHint(step, learning.DisclosureNone, learning.HelpProgressive, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if level != learning.DisclosureGuidingQuestion || kind != "guiding_question" {
		t.Fatalf("got level=%d kind=%q, want level=1 kind=guiding_question", level, kind)
	}

	level, kind, err = NextHint(step, learning.DisclosureGuidingQuestion, learning.HelpProgressive, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if level != learning.DisclosureConceptOrAPI || kind != KindSyntaxRecall {
		t.Fatalf("got level=%d kind=%q, want level=2 kind=%s (authored override)", level, kind, KindSyntaxRecall)
	}
}

func TestNextHintFallsBackToGenericKindWhenUnauthored(t *testing.T) {
	step := curriculum.StepAuthoring{} // no authored hints at all
	level, kind, err := NextHint(step, learning.DisclosureNone, learning.HelpProgressive, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if level != learning.DisclosureGuidingQuestion || kind != "guiding_question" {
		t.Fatalf("got level=%d kind=%q, want the §8.4 generic default", level, kind)
	}
}

func TestNextHintRequiresConfirmationForSolution(t *testing.T) {
	step := fixtureStep()
	_, _, err := NextHint(step, learning.DisclosureSkeleton, learning.HelpProgressive, false)
	if !errors.Is(err, ErrSolutionRequiresConfirmation) {
		t.Fatalf("expected ErrSolutionRequiresConfirmation, got %v", err)
	}
	level, kind, err := NextHint(step, learning.DisclosureSkeleton, learning.HelpProgressive, true)
	if err != nil {
		t.Fatalf("unexpected error with confirmation: %v", err)
	}
	if level != learning.DisclosureSolution || kind != "solution" {
		t.Fatalf("got level=%d kind=%q, want level=6 kind=solution", level, kind)
	}
}

func TestNextHintBlockedByNoHintsPolicy(t *testing.T) {
	step := fixtureStep()
	_, _, err := NextHint(step, learning.DisclosureNone, learning.HelpNoHints, false)
	if !errors.Is(err, ErrHelpDisabled) {
		t.Fatalf("expected ErrHelpDisabled, got %v", err)
	}
}

func TestSyntaxRecallLevelFindsAuthoredRung(t *testing.T) {
	level, ok := SyntaxRecallLevel(fixtureStep())
	if !ok || level != learning.DisclosureConceptOrAPI {
		t.Fatalf("got level=%d ok=%v, want level=2 ok=true", level, ok)
	}
}

func TestSyntaxRecallLevelNotFoundWhenUnauthored(t *testing.T) {
	_, ok := SyntaxRecallLevel(curriculum.StepAuthoring{})
	if ok {
		t.Fatal("expected no syntax-recall rung on a step with no hints")
	}
}

func TestStartDetourOpensAndActivates(t *testing.T) {
	policy, err := learning.NewSessionPolicy(learning.ModePractice, learning.DepthMicro, learning.HelpProgressive, learning.DisclosureConceptOrAPI, learning.EvaluationOnDemand, learning.AdvanceExplicit, nil)
	if err != nil {
		t.Fatalf("unexpected error building policy: %v", err)
	}
	sess := learning.NewLearningSession("ses_1", policy, learning.CatalogRef{})
	if err := sess.SetActiveInstruction(learning.NewStepProgress("step_1")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	d, err := StartDetour(sess, "why exported fields?")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.State != learning.DetourStateActive {
		t.Fatalf("detour state = %s, want active", d.State)
	}
	if sess.ActiveStep().StepID != "step_1" {
		t.Fatal("starting a detour must not change the active step")
	}
}

func TestFinishDetourResolvesAndReturns(t *testing.T) {
	sess := learning.NewLearningSession("ses_1", learning.SessionPolicy{}, learning.CatalogRef{})
	if err := sess.SetActiveInstruction(learning.NewStepProgress("step_1")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	d, err := StartDetour(sess, "why exported fields?")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := FinishDetour(d, DetourResolved); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.State != learning.DetourStateReturned {
		t.Fatalf("detour state = %s, want returned", d.State)
	}
}

func TestFinishDetourAbandoned(t *testing.T) {
	sess := learning.NewLearningSession("ses_1", learning.SessionPolicy{}, learning.CatalogRef{})
	if err := sess.SetActiveInstruction(learning.NewStepProgress("step_1")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	d, err := StartDetour(sess, "unrelated tangent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := FinishDetour(d, DetourAbandoned); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.State != learning.DetourStateReturned {
		t.Fatalf("detour state = %s, want returned", d.State)
	}
}
