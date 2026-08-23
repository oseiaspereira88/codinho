package session

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/oseiaspereira88/codinho/internal/assistance"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

// newHintedTestService is like newTestService but its fixture step authors
// a full hint ladder (guiding_question, syntax_recall override, logical
// prose), matching PROJECT.md §14.9's example, so hint_request and
// syntax_recall_get have real content to grant.
func newHintedTestService(t *testing.T) *Service {
	t.Helper()
	dir := t.TempDir()
	writeTestFile(t, dir, "manifest.yaml", "schema_version: 1\npacks:\n  - pack.yaml\n")
	writeTestFile(t, dir, "pack.yaml", `schema_version: 1
id: fixture-pack
version: 1.0.0
concepts:
  - id: named-types
    title: Named types
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
          - id: fixture.macro-one
            kind: micro
            title: Macro step
            instruction:
              objective: Declare the User type.
              scope: Only the declaration.
            concepts: [named-types]
            hints:
              - level: 1
                kind: guiding_question
              - level: 2
                kind: syntax_recall
              - level: 3
                kind: logical_prose
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

func startHinted(t *testing.T, svc *Service, in StartInput) StartResult {
	t.Helper()
	if in.ChallengeID == "" {
		in.ChallengeID = fixtureChallengeID
	}
	start, err := svc.Start(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return start
}

func TestHintRequestClimbsLadderAndStopsAtSessionCap(t *testing.T) {
	svc := newHintedTestService(t)
	start := startHinted(t, svc, StartInput{DisclosureMax: learning.DisclosureLogicalOutline})
	rev := start.Revision

	first, err := svc.HintRequest(start.SessionID, false, rev, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if first.Level != learning.DisclosureGuidingQuestion || first.Kind != "guiding_question" {
		t.Fatalf("unexpected first hint: %+v", first)
	}
	rev = first.Revision

	second, err := svc.HintRequest(start.SessionID, false, rev, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if second.Level != learning.DisclosureConceptOrAPI || second.Kind != assistance.KindSyntaxRecall {
		t.Fatalf("unexpected second hint: %+v", second)
	}
	rev = second.Revision

	third, err := svc.HintRequest(start.SessionID, false, rev, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if third.Level != learning.DisclosureLogicalOutline {
		t.Fatalf("unexpected third hint: %+v", third)
	}
	rev = third.Revision

	_, err = svc.HintRequest(start.SessionID, false, rev, "")
	if !errors.Is(err, learning.DomainError{Code: learning.ErrCodeDisclosureExceeded}) {
		t.Fatalf("expected ErrCodeDisclosureExceeded at the session cap, got %v", err)
	}
}

func TestHintRequestSolutionRequiresConfirmation(t *testing.T) {
	svc := newHintedTestService(t)
	start := startHinted(t, svc, StartInput{DisclosureMax: learning.DisclosureSolution})
	rev := start.Revision

	// Climb to level 5 first (levels 1-3 authored, 4-5 fall back to the
	// generic §8.4 default).
	for i := 0; i < 5; i++ {
		result, err := svc.HintRequest(start.SessionID, false, rev, "")
		if err != nil {
			t.Fatalf("unexpected error climbing to level %d: %v", i+1, err)
		}
		rev = result.Revision
	}

	if _, err := svc.HintRequest(start.SessionID, false, rev, ""); !errors.Is(err, assistance.ErrSolutionRequiresConfirmation) {
		t.Fatalf("expected ErrSolutionRequiresConfirmation, got %v", err)
	}

	result, err := svc.HintRequest(start.SessionID, true, rev, "")
	if err != nil {
		t.Fatalf("unexpected error with confirmation: %v", err)
	}
	if result.Level != learning.DisclosureSolution {
		t.Fatalf("level = %d, want solution", result.Level)
	}

	get, err := svc.Get(start.SessionID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !get.Disclosure.SolutionRevealed {
		t.Fatal("revealing the solution rung must mark the step's disclosure as solution-revealed")
	}
}

func TestHintRequestBlockedByNoHintsPolicy(t *testing.T) {
	svc := newHintedTestService(t)
	help := learning.HelpNoHints
	start := startHinted(t, svc, StartInput{Help: help, DisclosureMax: learning.DisclosureSolution})

	if _, err := svc.HintRequest(start.SessionID, false, start.Revision, ""); !errors.Is(err, assistance.ErrHelpDisabled) {
		t.Fatalf("expected ErrHelpDisabled, got %v", err)
	}
}

func TestHintRequestIsIdempotentByRequestID(t *testing.T) {
	svc := newHintedTestService(t)
	start := startHinted(t, svc, StartInput{DisclosureMax: learning.DisclosureLogicalOutline})

	first, err := svc.HintRequest(start.SessionID, false, start.Revision, "hint-retry")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := svc.HintRequest(start.SessionID, false, start.Revision, "hint-retry")
	if err != nil {
		t.Fatalf("unexpected error on retry: %v", err)
	}
	if first.Level != second.Level || first.Revision != second.Revision {
		t.Fatalf("retry was not idempotent: %+v vs %+v", first, second)
	}
}

func TestSyntaxRecallGetIsFreeInTeachingMode(t *testing.T) {
	svc := newHintedTestService(t)
	start := startHinted(t, svc, StartInput{Mode: learning.ModeTeaching, DisclosureMax: learning.DisclosureLogicalOutline})

	result, err := svc.SyntaxRecallGet(start.SessionID, start.Revision, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Level != learning.DisclosureConceptOrAPI || result.Kind != assistance.KindSyntaxRecall {
		t.Fatalf("unexpected result: %+v", result)
	}

	get, err := svc.Get(start.SessionID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	svc.mu.Lock()
	hintLevel := svc.sessions[start.SessionID].session.ActiveStep().HintLevel
	svc.mu.Unlock()
	if hintLevel != learning.DisclosureNone {
		t.Fatalf("teaching-mode syntax recall must not climb the ladder, HintLevel = %d", hintLevel)
	}
	_ = get
}

func TestSyntaxRecallGetConsumesLadderOutsideTeachingMode(t *testing.T) {
	svc := newHintedTestService(t)
	start := startHinted(t, svc, StartInput{Mode: learning.ModePractice, DisclosureMax: learning.DisclosureLogicalOutline})

	if _, err := svc.SyntaxRecallGet(start.SessionID, start.Revision, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	svc.mu.Lock()
	hintLevel := svc.sessions[start.SessionID].session.ActiveStep().HintLevel
	svc.mu.Unlock()
	if hintLevel != learning.DisclosureConceptOrAPI {
		t.Fatalf("HintLevel = %d, want %d", hintLevel, learning.DisclosureConceptOrAPI)
	}
}

func TestSyntaxRecallGetNoAuthoredHintFails(t *testing.T) {
	svc := newTestService(t) // fixture has no authored hints
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := svc.SyntaxRecallGet(start.SessionID, start.Revision, ""); !errors.Is(err, assistance.ErrNoHintAuthored) {
		t.Fatalf("expected ErrNoHintAuthored, got %v", err)
	}
}

func TestDetourStartAndFinishDoNotChangeActiveStep(t *testing.T) {
	svc := newHintedTestService(t)
	start := startHinted(t, svc, StartInput{})

	startResult, err := svc.DetourStart(start.SessionID, "why exported fields?", start.Revision, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if startResult.State != learning.DetourStateActive {
		t.Fatalf("state = %s, want active", startResult.State)
	}

	get, err := svc.Get(start.SessionID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if get.ActiveStep != start.ActiveStep {
		t.Fatalf("active step changed during detour: %s vs %s", get.ActiveStep, start.ActiveStep)
	}

	finishResult, err := svc.DetourFinish(start.SessionID, assistance.DetourResolved, startResult.Revision, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if finishResult.State != learning.DetourStateReturned {
		t.Fatalf("state = %s, want returned", finishResult.State)
	}

	get, err = svc.Get(start.SessionID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if get.ActiveStep != start.ActiveStep {
		t.Fatalf("active step changed after detour finish: %s vs %s", get.ActiveStep, start.ActiveStep)
	}
}

func TestDetourFinishWithoutOpenDetourFails(t *testing.T) {
	svc := newHintedTestService(t)
	start := startHinted(t, svc, StartInput{})

	if _, err := svc.DetourFinish(start.SessionID, assistance.DetourResolved, start.Revision, ""); !errors.Is(err, ErrNoOpenDetour) {
		t.Fatalf("expected ErrNoOpenDetour, got %v", err)
	}
}

func TestDetourStartTwiceWithoutFinishingFails(t *testing.T) {
	svc := newHintedTestService(t)
	start := startHinted(t, svc, StartInput{})

	first, err := svc.DetourStart(start.SessionID, "first question", start.Revision, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var domainErr learning.DomainError
	_, err = svc.DetourStart(start.SessionID, "second question", first.Revision, "")
	if !errors.As(err, &domainErr) || domainErr.Code != learning.ErrCodeInvalidDetourTransition {
		t.Fatalf("expected ErrCodeInvalidDetourTransition, got %v", err)
	}
}

func TestDetourFinishIsIdempotentByRequestID(t *testing.T) {
	svc := newHintedTestService(t)
	start := startHinted(t, svc, StartInput{})

	startResult, err := svc.DetourStart(start.SessionID, "why?", start.Revision, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	first, err := svc.DetourFinish(start.SessionID, assistance.DetourResolved, startResult.Revision, "finish-retry")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := svc.DetourFinish(start.SessionID, assistance.DetourResolved, startResult.Revision, "finish-retry")
	if err != nil {
		t.Fatalf("unexpected error on retry: %v", err)
	}
	if first.Revision != second.Revision || second.State != learning.DetourStateReturned {
		t.Fatalf("retry was not idempotent: %+v vs %+v", first, second)
	}
}
