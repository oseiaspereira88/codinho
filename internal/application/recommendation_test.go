package application

import (
	"path/filepath"
	"testing"

	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
)

func newRecommendationTestCatalog(t *testing.T) *curriculum.Catalog {
	t.Helper()
	dir := t.TempDir()
	writeTestFile(t, dir, "manifest.yaml", "schema_version: 1\npacks:\n  - pack.yaml\n")
	writeTestFile(t, dir, "pack.yaml", `schema_version: 1
id: fixture-pack
version: 1.0.0
themes:
  - id: slices
    title: Slices
competencies:
  - id: comp-a
    title: Competency A
challenges:
  - schema_version: 1
    id: challenge-big
    version: 1.0.0
    title: Big challenge
    kind: atomic
    difficulty: foundational
    estimated_minutes: 30
    themes: [slices]
    competencies:
      primary: [comp-a]
  - schema_version: 1
    id: challenge-small
    version: 1.0.0
    title: Small challenge
    kind: atomic
    difficulty: foundational
    estimated_minutes: 10
    themes: [slices]
    competencies:
      primary: [comp-a]
  - schema_version: 1
    id: challenge-blocked
    version: 1.0.0
    title: Blocked challenge
    kind: atomic
    difficulty: foundational
    themes: [slices]
    prerequisites: [challenge-big]
`)
	catalog, diags, err := curriculum.Load(dir, curriculum.DefaultLimits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostics: %+v", diags)
	}
	return catalog
}

func TestRecommendationServiceWithoutProgressStillExcludesBlockedChallenge(t *testing.T) {
	catalog := newRecommendationTestCatalog(t)
	svc := NewRecommendationService(catalog, nil)

	got, err := svc.Recommend(RecommendInput{CompetencyID: "comp-a"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// No progress wired means no evidence at all for comp-a, which is
	// itself below the remediation threshold (not_observed < demonstrates_
	// without_help), so the smaller challenge is recommended as
	// remediation ahead of the bigger one — 3 entries, blocked excluded.
	if len(got) != 3 {
		t.Fatalf("expected 3 entries (remediation + 2 eligible, blocked excluded), got %+v", got)
	}
	for _, r := range got {
		if r.ChallengeID == "challenge-blocked" {
			t.Fatal("challenge-blocked has an unmet prerequisite and must be excluded")
		}
	}
}

func TestRecommendationServiceExcludesUnmetPrerequisiteUntilCompleted(t *testing.T) {
	catalog := newRecommendationTestCatalog(t)
	svc := NewRecommendationService(catalog, nil)

	withoutCompletion, err := svc.Recommend(RecommendInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range withoutCompletion {
		if r.ChallengeID == "challenge-blocked" {
			t.Fatal("challenge-blocked must be excluded before its prerequisite is completed")
		}
	}

	withCompletion, err := svc.Recommend(RecommendInput{Completed: []string{"challenge-big"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, r := range withCompletion {
		if r.ChallengeID == "challenge-blocked" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected challenge-blocked to become eligible once its prerequisite is completed")
	}
}

func TestRecommendationServiceFoldsInMasteryAndReviewDue(t *testing.T) {
	catalog := newRecommendationTestCatalog(t)
	path := filepath.Join(t.TempDir(), "events.jsonl")
	store, err := eventstore.Open(path, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer store.Close()
	progress := NewProgressService(store)
	if _, err := progress.RecordEvidence(EvidenceInput{
		CompetencyID: "comp-a", Dimension: "guided_implementation", EvidenceID: "ev1", HelpUsed: true, Success: true,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	svc := NewRecommendationService(catalog, progress)
	got, err := svc.Recommend(RecommendInput{CompetencyID: "comp-a"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// comp-a is only demonstrates_with_help: below the remediation
	// threshold, so the smaller challenge should be recommended first, as
	// a remediation step for the bigger one.
	if len(got) != 3 {
		t.Fatalf("expected 3 entries (remediation + 2 eligible), got %+v", got)
	}
	if got[0].ChallengeID != "challenge-small" || got[0].Explanation.RemediationFor != "challenge-big" {
		t.Fatalf("expected challenge-small first as remediation, got %+v", got[0])
	}
}
