package learning

import (
	"errors"
	"testing"
)

func TestGrantHintClimbsOneRungAtATime(t *testing.T) {
	p := NewStepProgress("step_1")
	policy := DisclosurePolicy{MaxLevel: DisclosureSolution}

	if err := p.GrantHint(DisclosureGuidingQuestion, policy); err != nil {
		t.Fatalf("unexpected error granting level 1: %v", err)
	}
	if p.HintLevel != DisclosureGuidingQuestion {
		t.Fatalf("HintLevel = %d, want %d", p.HintLevel, DisclosureGuidingQuestion)
	}
	if err := p.GrantHint(DisclosureConceptOrAPI, policy); err != nil {
		t.Fatalf("unexpected error granting level 2: %v", err)
	}
	if p.HintLevel != DisclosureConceptOrAPI {
		t.Fatalf("HintLevel = %d, want %d", p.HintLevel, DisclosureConceptOrAPI)
	}
}

func TestGrantHintRejectsSkippingARung(t *testing.T) {
	p := NewStepProgress("step_1")
	policy := DisclosurePolicy{MaxLevel: DisclosureSolution}

	if err := p.GrantHint(DisclosureLogicalOutline, policy); !errors.Is(err, DomainError{Code: ErrCodeHintLevelSkipped}) {
		t.Fatalf("expected ErrCodeHintLevelSkipped, got %v", err)
	}
	if p.HintLevel != DisclosureNone {
		t.Fatalf("HintLevel changed on a rejected grant: %d", p.HintLevel)
	}
}

func TestGrantHintRejectsAboveSessionCap(t *testing.T) {
	p := NewStepProgress("step_1")
	policy := DisclosurePolicy{MaxLevel: DisclosureGuidingQuestion}

	if err := p.GrantHint(DisclosureGuidingQuestion, policy); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := p.GrantHint(DisclosureConceptOrAPI, policy); !errors.Is(err, DomainError{Code: ErrCodeDisclosureExceeded}) {
		t.Fatalf("expected ErrCodeDisclosureExceeded, got %v", err)
	}
}

func TestGrantDirectBypassesRatchetButNotCap(t *testing.T) {
	p := NewStepProgress("step_1")
	policy := DisclosurePolicy{MaxLevel: DisclosureLogicalOutline}

	// GrantDirect targets a specific rung out of ladder order (syntax
	// recall's use case), which GrantHint would reject as a skip.
	if err := p.GrantDirect(DisclosureConceptOrAPI, policy); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.HintLevel != DisclosureConceptOrAPI {
		t.Fatalf("HintLevel = %d, want %d", p.HintLevel, DisclosureConceptOrAPI)
	}
	// A lower level than already reached never regresses HintLevel.
	if err := p.GrantDirect(DisclosureGuidingQuestion, policy); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.HintLevel != DisclosureConceptOrAPI {
		t.Fatalf("HintLevel regressed to %d", p.HintLevel)
	}
	if err := p.GrantDirect(DisclosureSolution, policy); !errors.Is(err, DomainError{Code: ErrCodeDisclosureExceeded}) {
		t.Fatalf("expected ErrCodeDisclosureExceeded above the session cap, got %v", err)
	}
}
