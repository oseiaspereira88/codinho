package mastery

import (
	"testing"
	"time"
)

func day(n int) time.Time {
	return time.Date(2026, 1, n, 12, 0, 0, 0, time.UTC)
}

func TestAdvanceStateFirstContactAlwaysIntroduces(t *testing.T) {
	cases := []Signal{
		{Success: false, Observed: day(1)},
		{Success: true, SolutionRevealed: true, Observed: day(1)},
	}
	for _, sig := range cases {
		got := AdvanceState(RuleVersionV1, DimensionProjection{}, sig)
		if got.State != StateIntroduced {
			t.Errorf("sig=%+v: state = %s, want introduced", sig, got.State)
		}
	}
}

func TestAdvanceStateGuidedSuccessNeverExceedsDemonstratesWithHelp(t *testing.T) {
	proj := DimensionProjection{}
	for i := range 10 {
		proj = AdvanceState(RuleVersionV1, proj, Signal{Success: true, HelpUsed: true, Observed: day(1 + i)})
	}
	if proj.State != StateDemonstratesWithHelp {
		t.Fatalf("after 10 guided successes, state = %s, want demonstrates_with_help (requirement R4, invariant 10)", proj.State)
	}
}

func TestAdvanceStateAutonomousSuccessReachesDemonstratesWithoutHelpDirectly(t *testing.T) {
	proj := AdvanceState(RuleVersionV1, DimensionProjection{}, Signal{Success: true, HelpUsed: false, Observed: day(1)})
	if proj.State != StateDemonstratesWithoutHelp {
		t.Fatalf("state = %s, want demonstrates_without_help", proj.State)
	}
}

func TestAdvanceStateSolutionRevealedNeverPromotes(t *testing.T) {
	proj := DimensionProjection{State: StateDemonstratesWithHelp}
	got := AdvanceState(RuleVersionV1, proj, Signal{Success: true, HelpUsed: false, SolutionRevealed: true, Observed: day(1)})
	if got.State != StateDemonstratesWithHelp {
		t.Fatalf("state = %s, want unchanged demonstrates_with_help (invariant 8)", got.State)
	}
}

func TestAdvanceStateFailureNeverRegresses(t *testing.T) {
	proj := DimensionProjection{State: StateDemonstratesWithoutHelp}
	got := AdvanceState(RuleVersionV1, proj, Signal{Success: false, Observed: day(1)})
	if got.State != StateDemonstratesWithoutHelp {
		t.Fatalf("state = %s, want unchanged demonstrates_without_help", got.State)
	}
}

func TestAdvanceStateRetentionRequiresADifferentCalendarDay(t *testing.T) {
	proj := AdvanceState(RuleVersionV1, DimensionProjection{}, Signal{Success: true, Observed: day(1)})
	if proj.State != StateDemonstratesWithoutHelp {
		t.Fatalf("precondition failed: state = %s", proj.State)
	}

	sameDay := AdvanceState(RuleVersionV1, proj, Signal{Success: true, Observed: day(1).Add(2 * time.Hour)})
	if sameDay.State != StateDemonstratesWithoutHelp {
		t.Fatalf("same-day repeat: state = %s, want still demonstrates_without_help", sameDay.State)
	}

	laterDay := AdvanceState(RuleVersionV1, proj, Signal{Success: true, Observed: day(2)})
	if laterDay.State != StateRetained {
		t.Fatalf("later-day repeat: state = %s, want retained", laterDay.State)
	}
}

func TestAdvanceStateTransferRequiresAnUnseenVariant(t *testing.T) {
	proj := AdvanceState(RuleVersionV1, DimensionProjection{}, Signal{Success: true, Variant: "v1", Observed: day(1)})
	proj = AdvanceState(RuleVersionV1, proj, Signal{Success: true, Variant: "v1", Observed: day(2)})
	if proj.State != StateRetained {
		t.Fatalf("precondition failed: state = %s", proj.State)
	}

	sameVariant := AdvanceState(RuleVersionV1, proj, Signal{Success: true, Variant: "v1", Observed: day(3)})
	if sameVariant.State != StateRetained {
		t.Fatalf("same variant: state = %s, want still retained", sameVariant.State)
	}

	newVariant := AdvanceState(RuleVersionV1, proj, Signal{Success: true, Variant: "v2", Observed: day(3)})
	if newVariant.State != StateTransferred {
		t.Fatalf("new variant: state = %s, want transferred", newVariant.State)
	}
}

func TestHigherReturnsTheFurtherAlongState(t *testing.T) {
	if got := Higher(StateIntroduced, StateRetained); got != StateRetained {
		t.Fatalf("Higher(introduced, retained) = %s, want retained", got)
	}
	if got := Higher(StateTransferred, StateNotObserved); got != StateTransferred {
		t.Fatalf("Higher(transferred, not_observed) = %s, want transferred", got)
	}
}

func TestAdvanceStateEvidenceCountAlwaysIncrements(t *testing.T) {
	proj := DimensionProjection{}
	for i := 1; i <= 5; i++ {
		proj = AdvanceState(RuleVersionV1, proj, Signal{Success: false, Observed: day(i)})
		if proj.EvidenceCount != i {
			t.Fatalf("evidence count = %d, want %d", proj.EvidenceCount, i)
		}
	}
}
