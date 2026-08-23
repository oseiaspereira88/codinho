package curriculum

import "fmt"

// Coverage is a projection of catalog composition, independent of any
// editorial finding (requirement R6).
type Coverage struct {
	Themes          int
	Concepts        int
	Competencies    int
	Tracks          int
	Challenges      int
	StepNodes       int
	ByDifficulty    map[string]int
	ByChallengeKind map[string]int
}

// ProjectCoverage computes Coverage from catalog's public API alone, so
// it never needs access to Catalog's private indexes.
func ProjectCoverage(catalog *Catalog) Coverage {
	cov := Coverage{
		Themes:          len(catalog.List(KindTheme, "")),
		Concepts:        len(catalog.List(KindConcept, "")),
		Competencies:    len(catalog.List(KindCompetency, "")),
		Tracks:          len(catalog.List(KindTrack, "")),
		ByDifficulty:    map[string]int{},
		ByChallengeKind: map[string]int{},
	}
	challenges := catalog.Challenges()
	cov.Challenges = len(challenges)
	for _, ch := range challenges {
		cov.ByDifficulty[ch.Difficulty]++
		cov.ByChallengeKind[ch.Kind]++
		cov.StepNodes += countStepNodes(ch.Layers)
	}
	return cov
}

func countStepNodes(layers []LayerAuthoring) int {
	var count func(steps []StepAuthoring) int
	count = func(steps []StepAuthoring) int {
		n := len(steps)
		for _, s := range steps {
			n += count(s.Children)
		}
		return n
	}
	n := 0
	for _, l := range layers {
		n += count(l.MacroSteps)
	}
	return n
}

// RuleV1GateBelowThreshold marks a dimension of Coverage below the V1
// roadmap's completion criteria.
const RuleV1GateBelowThreshold EditorialRuleID = "v1_gate_below_threshold"

// V1Thresholds are the roadmap's stated completion criteria (requirement
// R9; codinho-v1 roadmap "Critérios de conclusão"): 160 concepts, 100
// competencies, 84 challenges, 12 tracks, 500 step nodes.
var V1Thresholds = Coverage{Concepts: 160, Competencies: 100, Challenges: 84, Tracks: 12, StepNodes: 500}

// CheckV1Gate reports every Coverage dimension below V1Thresholds
// (requirement R9). It is never run as part of ordinary catalog
// validate: it is meant for the V1 acceptance gate, called explicitly
// once a full curriculum is authored, so an in-progress catalog (like
// this repository's own, still far below every threshold) never fails a
// routine build.
func CheckV1Gate(cov Coverage) []EditorialFinding {
	var out []EditorialFinding
	check := func(name string, got, want int) {
		if got < want {
			out = append(out, EditorialFinding{
				Item: "v1-gate", Rule: RuleV1GateBelowThreshold, Severity: SeverityBlocking,
				Detail:     fmt.Sprintf("%s = %d, want >= %d", name, got, want),
				Suggestion: "ampliar autoria de conteúdo nessa dimensão antes do aceite V1",
			})
		}
	}
	check("concepts", cov.Concepts, V1Thresholds.Concepts)
	check("competencies", cov.Competencies, V1Thresholds.Competencies)
	check("challenges", cov.Challenges, V1Thresholds.Challenges)
	check("tracks", cov.Tracks, V1Thresholds.Tracks)
	check("step_nodes", cov.StepNodes, V1Thresholds.StepNodes)
	return out
}

// RuleTypeDistributionMismatch marks a challenge kind or difficulty
// whose count does not exactly match a pack spec's declared
// distribution.
const RuleTypeDistributionMismatch EditorialRuleID = "type_distribution_mismatch"

// TypeDistribution declares the exact expected count per challenge kind
// or difficulty a pack spec defines (requirement R10). A nil or empty
// map means that dimension is not yet constrained — no pack spec has
// declared numbers for it yet.
type TypeDistribution struct {
	ByChallengeKind map[string]int
	ByDifficulty    map[string]int
}

// CheckTypeDistribution compares cov against want, reporting every exact
// mismatch (requirement R10).
func CheckTypeDistribution(cov Coverage, want TypeDistribution) []EditorialFinding {
	var out []EditorialFinding
	for kind, wantN := range want.ByChallengeKind {
		if got := cov.ByChallengeKind[kind]; got != wantN {
			out = append(out, EditorialFinding{
				Item: "distribution:kind:" + kind, Rule: RuleTypeDistributionMismatch, Severity: SeverityBlocking,
				Detail: fmt.Sprintf("challenge kind %q = %d, want exactly %d", kind, got, wantN),
			})
		}
	}
	for difficulty, wantN := range want.ByDifficulty {
		if got := cov.ByDifficulty[difficulty]; got != wantN {
			out = append(out, EditorialFinding{
				Item: "distribution:difficulty:" + difficulty, Rule: RuleTypeDistributionMismatch, Severity: SeverityBlocking,
				Detail: fmt.Sprintf("difficulty %q = %d, want exactly %d", difficulty, got, wantN),
			})
		}
	}
	return out
}
