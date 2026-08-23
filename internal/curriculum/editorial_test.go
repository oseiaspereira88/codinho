package curriculum

import "testing"

func baseChallenge(id string) ChallengeAuthoring {
	return ChallengeAuthoring{
		ID: id, Kind: "atomic", Difficulty: "foundational",
		Competencies: CompetencyRefs{Primary: []string{"comp-a"}},
		Acceptance:   []string{"saída correta"},
	}
}

func hasRule(findings []EditorialFinding, rule EditorialRuleID) bool {
	for _, f := range findings {
		if f.Rule == rule {
			return true
		}
	}
	return false
}

func TestCheckHintOrderFlagsNonIncreasingLevels(t *testing.T) {
	ch := baseChallenge("c1")
	ch.Layers = []LayerAuthoring{{MacroSteps: []StepAuthoring{{
		ID: "s1", Kind: "micro",
		Hints: []HintAuthoring{{Level: 2, Kind: "guiding_question"}, {Level: 1, Kind: "guiding_question"}},
	}}}}
	findings := checkHintOrder("f", ch.ID, ch.Layers)
	if !hasRule(findings, RuleHintsIncreasingOrder) {
		t.Fatalf("findings = %+v, want RuleHintsIncreasingOrder", findings)
	}
}

func TestCheckHintOrderAllowsIncreasingLevels(t *testing.T) {
	ch := baseChallenge("c1")
	ch.Layers = []LayerAuthoring{{MacroSteps: []StepAuthoring{{
		ID: "s1", Kind: "micro",
		Hints: []HintAuthoring{{Level: 1, Kind: "guiding_question"}, {Level: 2, Kind: "concept_or_api"}},
	}}}}
	if findings := checkHintOrder("f", ch.ID, ch.Layers); len(findings) != 0 {
		t.Fatalf("findings = %+v, want none", findings)
	}
}

func TestCheckHintOrderFlagsEarlySolution(t *testing.T) {
	ch := baseChallenge("c1")
	ch.Layers = []LayerAuthoring{{MacroSteps: []StepAuthoring{{
		ID: "s1", Kind: "micro",
		Hints: []HintAuthoring{{Level: 3, Kind: "solution"}},
	}}}}
	findings := checkHintOrder("f", ch.ID, ch.Layers)
	if !hasRule(findings, RuleHintRevealsSolutionEarly) {
		t.Fatalf("findings = %+v, want RuleHintRevealsSolutionEarly", findings)
	}
}

func TestCheckNoEmbeddedSolutionTextFlagsCodeFence(t *testing.T) {
	ch := baseChallenge("c1")
	ch.Brief = "Implemente assim:\n```go\nfunc f() {}\n```"
	findings := checkNoEmbeddedSolutionText("f", ch)
	if !hasRule(findings, RuleBriefContainsCodeFence) {
		t.Fatalf("findings = %+v, want RuleBriefContainsCodeFence", findings)
	}
}

func TestCheckCompoundIntentFlagsConjunction(t *testing.T) {
	ch := baseChallenge("c1")
	ch.Layers = []LayerAuthoring{{MacroSteps: []StepAuthoring{{
		ID: "s1", Kind: "micro",
		Instruction: InstructionAuthoring{Objective: "Declare o tipo e também implemente o método"},
	}}}}
	findings := checkCompoundIntent("f", ch.ID, ch.Layers)
	if !hasRule(findings, RuleCompoundMicroInstruction) {
		t.Fatalf("findings = %+v, want RuleCompoundMicroInstruction", findings)
	}
	for _, fnd := range findings {
		if fnd.Severity != SeverityWarning {
			t.Fatalf("compound intent must be advisory, got severity %s", fnd.Severity)
		}
	}
}

func TestCheckCompoundIntentAllowsSingleVerb(t *testing.T) {
	ch := baseChallenge("c1")
	ch.Layers = []LayerAuthoring{{MacroSteps: []StepAuthoring{{
		ID: "s1", Kind: "micro",
		Instruction: InstructionAuthoring{Objective: "Declare o tipo User"},
	}}}}
	if findings := checkCompoundIntent("f", ch.ID, ch.Layers); len(findings) != 0 {
		t.Fatalf("findings = %+v, want none", findings)
	}
}

func TestCheckCoverageGapsFlagsMissingCompetencyAndAcceptance(t *testing.T) {
	ch := ChallengeAuthoring{ID: "c1", Kind: "atomic"}
	findings := checkCoverageGaps("f", ch)
	if !hasRule(findings, RuleMissingCompetency) || !hasRule(findings, RuleMissingAcceptance) || !hasRule(findings, RuleMissingReflection) {
		t.Fatalf("findings = %+v", findings)
	}
}

func TestCheckCoverageGapsFlagsGatedStepWithoutHints(t *testing.T) {
	ch := baseChallenge("c1")
	ch.Layers = []LayerAuthoring{{MacroSteps: []StepAuthoring{{
		ID: "s1", Kind: "micro",
		Completion: CompletionAuthoring{RequiresPositiveEvaluation: true},
	}}}}
	findings := checkCoverageGaps("f", ch)
	if !hasRule(findings, RuleMissingHintsForGatedStep) {
		t.Fatalf("findings = %+v, want RuleMissingHintsForGatedStep", findings)
	}
}

func TestCheckCoverageGapsSatisfiedHasNoFindings(t *testing.T) {
	ch := baseChallenge("c1")
	ch.Layers = []LayerAuthoring{{MacroSteps: []StepAuthoring{{
		ID: "s1", Kind: "micro",
		Completion: CompletionAuthoring{RequiresPositiveEvaluation: true},
		Hints:      []HintAuthoring{{Level: 1, Kind: "guiding_question"}},
		Reflection: ReflectionAuthoring{Optional: []string{"Por que?"}},
	}}}}
	if findings := checkCoverageGaps("f", ch); len(findings) != 0 {
		t.Fatalf("findings = %+v, want none", findings)
	}
}

func TestCheckPublicationMetadataIgnoresDrafts(t *testing.T) {
	ch := baseChallenge("c1")
	if findings := checkPublicationMetadata("f", ch); len(findings) != 0 {
		t.Fatalf("an unpublished draft must never be flagged, got %+v", findings)
	}
}

func TestCheckPublicationMetadataRequiresDistinctReviewer(t *testing.T) {
	ch := baseChallenge("c1")
	ch.Publication = PublicationAuthoring{Status: StatusPublished, Author: "alice", ReviewedBy: "alice", Playtested: true}
	findings := checkPublicationMetadata("f", ch)
	if !hasRule(findings, RulePublishedSameReviewer) {
		t.Fatalf("findings = %+v, want RulePublishedSameReviewer", findings)
	}
}

func TestCheckPublicationMetadataRequiresPlaytest(t *testing.T) {
	ch := baseChallenge("c1")
	ch.Publication = PublicationAuthoring{Status: StatusPublished, Author: "alice", ReviewedBy: "bob"}
	findings := checkPublicationMetadata("f", ch)
	if !hasRule(findings, RulePublishedWithoutPlaytest) {
		t.Fatalf("findings = %+v, want RulePublishedWithoutPlaytest", findings)
	}
}

func TestCheckPublicationMetadataSatisfied(t *testing.T) {
	ch := baseChallenge("c1")
	ch.Publication = PublicationAuthoring{Status: StatusPublished, Author: "alice", ReviewedBy: "bob", Playtested: true}
	if findings := checkPublicationMetadata("f", ch); len(findings) != 0 {
		t.Fatalf("findings = %+v, want none", findings)
	}
}

func TestCheckRelationCoherenceFlagsIsolatedItem(t *testing.T) {
	packs := []Pack{{File: "f", Concepts: []ConceptAuthoring{{ID: "concept-a"}}}}
	findings := checkRelationCoherence(packs)
	if !hasRule(findings, RuleRelationIsolated) {
		t.Fatalf("findings = %+v, want RuleRelationIsolated", findings)
	}
}

func TestCheckRelationCoherenceAcceptsAnyNonPrecedenceRelation(t *testing.T) {
	packs := []Pack{{
		File:     "f",
		Concepts: []ConceptAuthoring{{ID: "a"}, {ID: "b"}},
		Relations: []RelationAuthoring{
			{From: "a", To: "b", Kind: string(RelationRelatesTo)},
		},
	}}
	if findings := checkRelationCoherence(packs); hasRule(findings, RuleRelationIsolated) {
		t.Fatalf("findings = %+v, want no RuleRelationIsolated", findings)
	}
}

func TestCheckRelationCoherenceFlagsNonReciprocalContrast(t *testing.T) {
	packs := []Pack{{
		File:     "f",
		Concepts: []ConceptAuthoring{{ID: "a"}, {ID: "b"}},
		Relations: []RelationAuthoring{
			{From: "a", To: "b", Kind: string(RelationContrastsWith)},
		},
	}}
	findings := checkRelationCoherence(packs)
	if !hasRule(findings, RuleRelationNotReciprocal) {
		t.Fatalf("findings = %+v, want RuleRelationNotReciprocal", findings)
	}
}

func TestCheckRelationCoherenceAcceptsReciprocalContrast(t *testing.T) {
	packs := []Pack{{
		File:     "f",
		Concepts: []ConceptAuthoring{{ID: "a"}, {ID: "b"}},
		Relations: []RelationAuthoring{
			{From: "a", To: "b", Kind: string(RelationContrastsWith)},
			{From: "b", To: "a", Kind: string(RelationContrastsWith)},
		},
	}}
	if findings := checkRelationCoherence(packs); hasRule(findings, RuleRelationNotReciprocal) {
		t.Fatalf("findings = %+v, want no RuleRelationNotReciprocal", findings)
	}
}

// TestNegativeFixtureCorpusTriggersEveryRule loads
// testdata/catalog-quality (never listed in the real packs/manifest.yaml)
// and confirms each authored negative example actually fires the rule
// its comment names, so the corpus itself never drifts silently out of
// sync with the checks it documents (Validation task: "Criar uma fixture
// negativa por regra").
func TestNegativeFixtureCorpusTriggersEveryRule(t *testing.T) {
	_, structuralDiags, err := LoadPacks("../../testdata/catalog-quality/structural", DefaultLimits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	blockingDiagCodes := map[DiagnosticCode]bool{}
	for _, d := range structuralDiags {
		blockingDiagCodes[d.Code] = true
	}
	if !blockingDiagCodes[DiagInvalidVersion] {
		t.Error("expected DiagInvalidVersion from neg.invalid-version")
	}
	if !blockingDiagCodes[DiagInvalidFixturePath] {
		t.Error("expected DiagInvalidFixturePath from neg.fixture-path-traversal")
	}

	packs, diags, err := LoadPacks("../../testdata/catalog-quality/editorial", DefaultLimits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if packs == nil {
		t.Fatalf("editorial fixture pack failed to load: %+v", diags)
	}

	findings := RunEditorialChecks(packs)
	wantRules := []EditorialRuleID{
		RuleMissingCompetency, RuleMissingAcceptance, RuleBriefContainsCodeFence,
		RuleHintsIncreasingOrder, RuleHintRevealsSolutionEarly, RuleCompoundMicroInstruction,
		RuleMissingHintsForGatedStep, RuleMissingReflection,
		RulePublishedSameReviewer, RulePublishedWithoutPlaytest,
		RuleRelationNotReciprocal,
	}
	for _, want := range wantRules {
		if !hasRule(findings, want) {
			t.Errorf("negative fixture corpus never triggered rule %q", want)
		}
	}
}

func TestRunEditorialChecksAggregatesAcrossPacks(t *testing.T) {
	ch := ChallengeAuthoring{ID: "c1", Kind: "atomic"}
	packs := []Pack{{File: "f", Challenges: []ChallengeAuthoring{ch}}}
	findings := RunEditorialChecks(packs)
	if !hasRule(findings, RuleMissingCompetency) {
		t.Fatalf("findings = %+v, want RuleMissingCompetency from the aggregated run", findings)
	}
}
