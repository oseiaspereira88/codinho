package curriculum

import "testing"

func TestProjectCoverageCountsEverything(t *testing.T) {
	packs := []Pack{{
		File:         "f",
		Themes:       []ThemeAuthoring{{ID: "t1"}},
		Concepts:     []ConceptAuthoring{{ID: "c1"}, {ID: "c2"}},
		Competencies: []CompetencyAuthoring{{ID: "comp1"}},
		Tracks:       []TrackAuthoring{{ID: "trk1"}},
		Challenges: []ChallengeAuthoring{
			{
				ID: "ch1", Kind: "atomic", Difficulty: "foundational",
				Layers: []LayerAuthoring{{MacroSteps: []StepAuthoring{
					{ID: "s1", Kind: "macro", Children: []StepAuthoring{{ID: "s1.1", Kind: "micro"}}},
				}}},
			},
			{ID: "ch2", Kind: "debug", Difficulty: "intermediate"},
		},
	}}
	catalog := newCatalog(packs)
	cov := ProjectCoverage(catalog)

	if cov.Themes != 1 || cov.Concepts != 2 || cov.Competencies != 1 || cov.Tracks != 1 || cov.Challenges != 2 {
		t.Fatalf("cov = %+v", cov)
	}
	if cov.StepNodes != 2 {
		t.Fatalf("StepNodes = %d, want 2 (s1 + s1.1)", cov.StepNodes)
	}
	if cov.ByChallengeKind["atomic"] != 1 || cov.ByChallengeKind["debug"] != 1 {
		t.Fatalf("ByChallengeKind = %+v", cov.ByChallengeKind)
	}
	if cov.ByDifficulty["foundational"] != 1 || cov.ByDifficulty["intermediate"] != 1 {
		t.Fatalf("ByDifficulty = %+v", cov.ByDifficulty)
	}
}

func TestCheckV1GateFlagsEveryDimensionBelowThreshold(t *testing.T) {
	findings := CheckV1Gate(Coverage{})
	// Every one of the five thresholded dimensions starts at zero, so
	// every one must be flagged (requirement R9).
	if len(findings) != 5 {
		t.Fatalf("findings = %+v, want exactly 5 (one per thresholded dimension)", findings)
	}
	for _, f := range findings {
		if f.Severity != SeverityBlocking {
			t.Fatalf("V1 gate findings must block, got %s", f.Severity)
		}
	}
}

func TestCheckV1GateSatisfiedHasNoFindings(t *testing.T) {
	if findings := CheckV1Gate(V1Thresholds); len(findings) != 0 {
		t.Fatalf("findings = %+v, want none when coverage exactly meets thresholds", findings)
	}
}

func TestCheckTypeDistributionFlagsMismatch(t *testing.T) {
	cov := Coverage{ByChallengeKind: map[string]int{"atomic": 3}, ByDifficulty: map[string]int{"foundational": 2}}
	want := TypeDistribution{ByChallengeKind: map[string]int{"atomic": 5}, ByDifficulty: map[string]int{"foundational": 2}}
	findings := CheckTypeDistribution(cov, want)
	if len(findings) != 1 || findings[0].Rule != RuleTypeDistributionMismatch {
		t.Fatalf("findings = %+v", findings)
	}
}

func TestCheckTypeDistributionEmptyWantNeverFlags(t *testing.T) {
	cov := Coverage{ByChallengeKind: map[string]int{"atomic": 1}}
	if findings := CheckTypeDistribution(cov, TypeDistribution{}); len(findings) != 0 {
		t.Fatalf("findings = %+v, want none when no distribution is declared yet", findings)
	}
}

func TestTypeDistributionIncludesUnexpectedKindsDeterministically(t *testing.T) {
	cov := Coverage{ByChallengeKind: map[string]int{"atomic": 1, "unexpected": 1}, ByDifficulty: map[string]int{"advanced": 1, "foundational": 1}}
	want := TypeDistribution{ByChallengeKind: map[string]int{"atomic": 1}, ByDifficulty: map[string]int{"foundational": 2}}
	for i := 0; i < 5; i++ {
		findings := CheckTypeDistribution(cov, want)
		if len(findings) != 3 || findings[0].Item != "distribution:kind:unexpected" || findings[1].Item != "distribution:difficulty:advanced" || findings[2].Item != "distribution:difficulty:foundational" {
			t.Fatalf("incomplete or nondeterministic distribution: %+v", findings)
		}
	}
	if len(CheckTypeDistribution(cov, TypeDistribution{})) != 0 {
		t.Fatal("unspecified distribution must remain unconstrained")
	}
}
