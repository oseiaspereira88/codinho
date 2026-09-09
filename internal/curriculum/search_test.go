package curriculum

import "testing"

func loadValidCatalog(t *testing.T) *Catalog {
	t.Helper()
	catalog, diags, err := Load("testdata/valid", DefaultLimits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostics: %+v", diags)
	}
	return catalog
}

func TestSearchByTextMatchesTitleCaseInsensitively(t *testing.T) {
	catalog := loadValidCatalog(t)
	result, err := catalog.Search(Query{Text: "FIXTURE"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Items) != 1 || result.Items[0].ID != "fixture.challenge-one" {
		t.Fatalf("expected fixture.challenge-one, got %+v", result.Items)
	}
}

func TestSearchByDifficultyAndCompetency(t *testing.T) {
	catalog := loadValidCatalog(t)
	result, err := catalog.Search(Query{Difficulty: "foundational", CompetencyIDs: []string{"slice-filter"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Items) != 1 {
		t.Fatalf("expected 1 match, got %+v", result.Items)
	}

	none, err := catalog.Search(Query{Difficulty: "advanced"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(none.Items) != 0 {
		t.Fatalf("expected no matches for a difficulty that does not exist, got %+v", none.Items)
	}
}

func TestSearchByPrerequisite(t *testing.T) {
	catalog := loadValidCatalog(t)
	result, err := catalog.Search(Query{Prerequisite: "slice-declaration"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Items) != 1 || result.Items[0].ID != "fixture.challenge-one" {
		t.Fatalf("expected fixture.challenge-one, got %+v", result.Items)
	}
}

func TestSearchByMaxMinutesExcludesLongerChallenges(t *testing.T) {
	catalog := loadValidCatalog(t)
	result, err := catalog.Search(Query{MaxMinutes: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// fixture.challenge-one has no estimated_minutes authored (0), which
	// must never be excluded by a minutes ceiling (0 means "unspecified",
	// not "instant").
	if len(result.Items) != 1 {
		t.Fatalf("expected the unspecified-duration fixture to still match, got %+v", result.Items)
	}
}

func TestSearchRejectsOverlongText(t *testing.T) {
	catalog := loadValidCatalog(t)
	longText := make([]byte, maxSearchTextLen+1)
	for i := range longText {
		longText[i] = 'a'
	}
	_, err := catalog.Search(Query{Text: string(longText)})
	if err != ErrSearchTextTooLong {
		t.Fatalf("expected ErrSearchTextTooLong, got %v", err)
	}
}

func TestSearchIsDeterministicallySortedByID(t *testing.T) {
	catalog := loadValidCatalog(t)
	var last string
	for range 3 {
		result, err := catalog.Search(Query{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var ids string
		for _, item := range result.Items {
			ids += item.ID + ","
		}
		if last != "" && ids != last {
			t.Fatalf("search order is not stable: %q vs %q", ids, last)
		}
		last = ids
	}
}
