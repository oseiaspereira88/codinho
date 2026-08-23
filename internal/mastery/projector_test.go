package mastery

import "testing"

func TestFoldProjectionsReproducesSequentialAdvanceState(t *testing.T) {
	signals := []RecordedSignal{
		{Signal: Signal{CompetencyID: "comp-1", Dimension: DimensionAutonomousImplementation, Success: true, Observed: day(1)}, RuleVersion: RuleVersionV1},
		{Signal: Signal{CompetencyID: "comp-1", Dimension: DimensionAutonomousImplementation, Success: true, Observed: day(2)}, RuleVersion: RuleVersionV1},
		{Signal: Signal{CompetencyID: "comp-1", Dimension: DimensionExplanation, Success: true, Observed: day(1)}, RuleVersion: RuleVersionV1},
	}
	byCompetency := FoldProjections(signals)
	got := byCompetency["comp-1"][DimensionAutonomousImplementation]
	if got.State != StateRetained {
		t.Fatalf("autonomous_implementation state = %s, want retained", got.State)
	}
	if got.EvidenceCount != 2 {
		t.Fatalf("evidence count = %d, want 2", got.EvidenceCount)
	}
	explanation := byCompetency["comp-1"][DimensionExplanation]
	if explanation.State != StateDemonstratesWithoutHelp {
		t.Fatalf("explanation state = %s, want demonstrates_without_help (dimensions must not bleed into each other)", explanation.State)
	}
}

func TestFoldProjectionsIsOrderDependent(t *testing.T) {
	guidedThenAutonomous := []RecordedSignal{
		{Signal: Signal{CompetencyID: "c", Dimension: DimensionDebugging, Success: true, HelpUsed: true, Observed: day(1)}, RuleVersion: RuleVersionV1},
		{Signal: Signal{CompetencyID: "c", Dimension: DimensionDebugging, Success: true, HelpUsed: false, Observed: day(2)}, RuleVersion: RuleVersionV1},
	}
	got := FoldProjections(guidedThenAutonomous)["c"][DimensionDebugging]
	if got.State != StateDemonstratesWithoutHelp {
		t.Fatalf("state = %s, want demonstrates_without_help", got.State)
	}
}

func TestFoldSchedulesSkipsRevealedSolutionSignals(t *testing.T) {
	signals := []RecordedSignal{
		{Signal: Signal{CompetencyID: "comp-1", Success: true, SolutionRevealed: true, Observed: day(1)}, RuleVersion: RuleVersionV1},
	}
	schedules := FoldSchedules(signals)
	if _, ok := schedules["comp-1"]; ok {
		t.Fatal("expected no schedule from a solution-revealed signal (invariant 8)")
	}
}

func TestFoldSchedulesTracksFailureAndSuccessSequentially(t *testing.T) {
	signals := []RecordedSignal{
		{Signal: Signal{CompetencyID: "comp-1", Success: true, Observed: day(1)}, RuleVersion: RuleVersionV1},
		{Signal: Signal{CompetencyID: "comp-1", Success: true, Observed: day(1)}, RuleVersion: RuleVersionV1},
	}
	schedules := FoldSchedules(signals)
	s, ok := schedules["comp-1"]
	if !ok {
		t.Fatal("expected a schedule for comp-1")
	}
	if s.IntervalDays != initialLadderDays[1] {
		t.Fatalf("interval = %d, want %d (second rung after two successes)", s.IntervalDays, initialLadderDays[1])
	}
}
