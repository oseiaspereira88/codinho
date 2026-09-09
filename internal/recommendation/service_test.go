package recommendation

import "testing"

func TestRecommendPrioritizesCompetencyMatchOverThemeMatch(t *testing.T) {
	candidates := []Candidate{
		{ID: "theme-only", ThemeIDs: []string{"slices"}},
		{ID: "competency-match", PrimaryCompetency: []string{"slice-filter"}},
	}
	objective := Objective{CompetencyIDs: []string{"slice-filter"}, ThemeIDs: []string{"slices"}}
	got := Recommend(objective, candidates, nil, nil)
	if len(got) != 2 || got[0].ChallengeID != "competency-match" {
		t.Fatalf("expected competency-match first, got %+v", got)
	}
}

func TestRecommendExcludesCandidatesWithUnmetPrerequisites(t *testing.T) {
	candidates := []Candidate{
		{ID: "blocked", Prerequisites: []string{"not-done"}},
		{ID: "free"},
	}
	got := Recommend(Objective{}, candidates, nil, map[string]bool{})
	if len(got) != 1 || got[0].ChallengeID != "free" {
		t.Fatalf("expected only 'free', got %+v", got)
	}
}

func TestRecommendIncludesCandidateOnceItsPrerequisiteIsCompleted(t *testing.T) {
	candidates := []Candidate{{ID: "next", Prerequisites: []string{"done"}}}
	got := Recommend(Objective{}, candidates, nil, map[string]bool{"done": true})
	if len(got) != 1 || got[0].ChallengeID != "next" {
		t.Fatalf("expected 'next' to be eligible, got %+v", got)
	}
}

func TestRecommendPrioritizesOverdueReview(t *testing.T) {
	candidates := []Candidate{
		{ID: "fresh", PrimaryCompetency: []string{"comp-a"}},
		{ID: "due", PrimaryCompetency: []string{"comp-b"}},
	}
	mastery := map[string]Mastery{"comp-b": {CompetencyID: "comp-b", ReviewOverdueBy: 5}}
	got := Recommend(Objective{}, candidates, mastery, nil)
	if len(got) != 2 || got[0].ChallengeID != "due" {
		t.Fatalf("expected the overdue-review candidate first, got %+v", got)
	}
}

func TestRecommendDeprioritizesAlreadyMasteredCompetencies(t *testing.T) {
	candidates := []Candidate{
		{ID: "mastered", PrimaryCompetency: []string{"comp-a"}},
		{ID: "unknown", PrimaryCompetency: []string{"comp-b"}},
	}
	mastery := map[string]Mastery{"comp-a": {CompetencyID: "comp-a", BestState: "retained"}}
	got := Recommend(Objective{}, candidates, mastery, nil)
	if len(got) != 2 || got[0].ChallengeID != "unknown" {
		t.Fatalf("expected 'unknown' (no evidence yet) ahead of an already-retained competency, got %+v", got)
	}
}

func TestRecommendTieBreaksDeterministicallyByID(t *testing.T) {
	candidates := []Candidate{{ID: "b"}, {ID: "a"}, {ID: "c"}}
	got := Recommend(Objective{}, candidates, nil, nil)
	if len(got) != 3 || got[0].ChallengeID != "a" || got[1].ChallengeID != "b" || got[2].ChallengeID != "c" {
		t.Fatalf("expected alphabetical tie-break, got %+v", got)
	}
}

func TestRecommendIsStableAcrossRepeatedCalls(t *testing.T) {
	candidates := []Candidate{
		{ID: "x", PrimaryCompetency: []string{"comp-a"}, EstimatedMinutes: 20},
		{ID: "y", PrimaryCompetency: []string{"comp-a"}, EstimatedMinutes: 10},
	}
	objective := Objective{CompetencyIDs: []string{"comp-a"}, TimeBudgetMinutes: 15}
	first := Recommend(objective, candidates, nil, nil)
	for range 5 {
		next := Recommend(objective, candidates, nil, nil)
		if len(next) != len(first) {
			t.Fatalf("unstable result length")
		}
		for i := range first {
			if first[i].ChallengeID != next[i].ChallengeID || first[i].Score != next[i].Score {
				t.Fatalf("requirement R8 violated: unstable ranking across identical calls: %+v vs %+v", first, next)
			}
		}
	}
}

func TestRecommendPrependsSmallerRemediationBelowThreshold(t *testing.T) {
	candidates := []Candidate{
		{ID: "big", PrimaryCompetency: []string{"comp-a"}, EstimatedMinutes: 30},
		{ID: "small", PrimaryCompetency: []string{"comp-a"}, EstimatedMinutes: 10},
	}
	objective := Objective{CompetencyIDs: []string{"comp-a"}}
	mastery := map[string]Mastery{"comp-a": {CompetencyID: "comp-a", BestState: "introduced"}}
	got := Recommend(objective, candidates, mastery, nil)
	if len(got) != 3 {
		t.Fatalf("expected big, small, and small-as-remediation (3 entries), got %+v", got)
	}
	if got[0].ChallengeID != "small" || got[0].Explanation.RemediationFor != "big" {
		t.Fatalf("expected 'small' first as remediation for 'big', got %+v", got[0])
	}
	if got[1].ChallengeID != "big" {
		t.Fatalf("expected the original top recommendation right after remediation, got %+v", got[1])
	}
}

func TestRecommendSkipsRemediationOnceAutonomousWithoutHelp(t *testing.T) {
	candidates := []Candidate{
		{ID: "big", PrimaryCompetency: []string{"comp-a"}, EstimatedMinutes: 30},
		{ID: "small", PrimaryCompetency: []string{"comp-a"}, EstimatedMinutes: 10},
	}
	objective := Objective{CompetencyIDs: []string{"comp-a"}}
	mastery := map[string]Mastery{"comp-a": {CompetencyID: "comp-a", BestState: "demonstrates_without_help"}}
	got := Recommend(objective, candidates, mastery, nil)
	if len(got) != 2 {
		t.Fatalf("expected no remediation once past the threshold, got %+v", got)
	}
}

func TestExplanationReportsEvidenceGapDependencyAndCost(t *testing.T) {
	c := Candidate{ID: "x", PrimaryCompetency: []string{"comp-a"}, EstimatedMinutes: 15, Prerequisites: []string{"p1"}}
	m := Mastery{CompetencyID: "comp-a", BestState: "introduced"}
	e := explain(c, m, map[string]bool{}, "")
	if e.CostMinutes != 15 {
		t.Fatalf("cost = %d, want 15", e.CostMinutes)
	}
	if e.Gap == "nenhum gap conhecido" {
		t.Fatal("expected a real gap below the remediation threshold")
	}
	if e.Dependency == "sem pré-requisitos pendentes" {
		t.Fatal("expected the unmet prerequisite to be reported")
	}
}
