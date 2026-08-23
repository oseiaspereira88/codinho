package curriculum

import (
	"fmt"
	"testing"
	"time"
)

// syntheticCatalogAtV1Scale builds a Catalog with challenge counts around
// the V1 target size (PROJECT.md: 84 challenges) without touching YAML or
// the filesystem, purely to validate the non-functional p95 budget
// ("Consultas p95 inferiores a 100 ms após indexação").
func syntheticCatalogAtV1Scale(n int) *Catalog {
	packs := []Pack{{ID: "bench", SchemaVersion: SchemaVersion}}
	for i := range n {
		packs[0].Challenges = append(packs[0].Challenges, ChallengeAuthoring{
			ID: fmt.Sprintf("bench.challenge-%d", i), Title: fmt.Sprintf("Challenge %d", i),
			Difficulty: "foundational", EstimatedMinutes: 15,
			Themes:       []string{"bench-theme"},
			Competencies: CompetencyRefs{Primary: []string{"bench-competency"}},
		})
	}
	return newCatalog(packs)
}

func BenchmarkSearchAtV1Scale(b *testing.B) {
	catalog := syntheticCatalogAtV1Scale(84)
	for b.Loop() {
		if _, err := catalog.Search(Query{Difficulty: "foundational", Competency: "bench-competency", Text: "Challenge"}); err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

// TestSearchP95BudgetAtV1Scale is a deterministic sanity check (not a
// flaky wall-clock benchmark gate) that a single Search call at V1 scale
// completes well within the 100ms p95 budget on ordinary hardware.
func TestSearchP95BudgetAtV1Scale(t *testing.T) {
	catalog := syntheticCatalogAtV1Scale(84)
	start := time.Now()
	if _, err := catalog.Search(Query{Text: "Challenge"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if elapsed := time.Since(start); elapsed > 100*time.Millisecond {
		t.Fatalf("Search took %s, want well under the 100ms p95 budget", elapsed)
	}
}
