package curriculum

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"
)

func selectionCatalog() *Catalog {
	return newCatalog([]Pack{{Challenges: []ChallengeAuthoring{
		{ID: "a", Version: "1.0.0", Themes: []string{"alpha"}},
		{ID: "b", Version: "1.0.0", Themes: []string{"beta"}, Prerequisites: []string{"z"}},
		{ID: "z", Version: "1.0.0"},
	}, Relations: []RelationAuthoring{{From: "b", To: "a", Kind: "requires"}}}})
}
func TestSelectionCoverageAndComposition(t *testing.T) {
	c := selectionCatalog()
	for _, tc := range []struct {
		themes   []string
		coverage string
		complete bool
		missing  []string
	}{
		{[]string{"alpha"}, "total", true, []string{}},
		{[]string{"alpha", "beta"}, "partial", true, []string{}},
		{[]string{"alpha", "unknown"}, "partial", false, []string{"unknown"}},
		{[]string{"unknown"}, "none", false, []string{"unknown"}},
	} {
		r, err := c.Search(Query{ThemeIDs: tc.themes})
		if err != nil {
			t.Fatal(err)
		}
		if r.Selection.Coverage != tc.coverage || r.Selection.Complete != tc.complete || !reflect.DeepEqual(r.Selection.MissingThemeIDs, tc.missing) {
			t.Fatalf("%+v", r.Selection)
		}
	}
	p, s, err := c.Compose([]string{"beta", "alpha", "alpha"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !s.Complete || s.Coverage != "partial" || !reflect.DeepEqual(p.ChallengeIDs, []string{"a", "z", "b"}) {
		t.Fatalf("%+v %+v", p, s)
	}
	p2, _, err := c.Compose([]string{"alpha", "beta"}, nil)
	if err != nil || !reflect.DeepEqual(p, p2) {
		t.Fatalf("not deterministic: %+v %v", p2, err)
	}
	chs, err := c.ResolvePath(p.ChallengeIDs)
	if err != nil || PathIdentity(chs) != p.CompositionID {
		t.Fatal("identity mismatch", err)
	}
	for _, ids := range [][]string{nil, {}, {"a", "a"}, {"missing"}, {"b", "a", "z"}} {
		if _, err := c.ResolvePath(ids); err == nil {
			t.Fatalf("accepted invalid path %v", ids)
		}
	}
}
func TestSelectionBounds(t *testing.T) {
	for _, ids := range [][]string{{}, {""}, make([]string, 17)} {
		if _, err := selectionCatalog().Search(Query{ThemeIDs: ids}); err == nil {
			t.Fatalf("accepted %v", ids)
		}
	}
}
func TestTrackMembershipValidation(t *testing.T) {
	for _, ids := range [][]string{{"unknown"}, {"a", "a"}, {}} {
		ds := validateTracks([]Pack{{Tracks: []TrackAuthoring{{ID: "t", ChallengeIDs: ids}}, Challenges: []ChallengeAuthoring{{ID: "a"}}}})
		if len(ds) == 0 || !ds[0].Blocking {
			t.Fatalf("invalid track accepted: %v", ids)
		}
	}
	if ds := validateTracks([]Pack{{Tracks: []TrackAuthoring{{ID: "legacy"}}}}); len(ds) > 0 {
		t.Fatal(ds)
	}
}

func TestCompositionBoundsAndMixedDependencyCycle(t *testing.T) {
	var chs []ChallengeAuthoring
	for i := 0; i < 101; i++ {
		chs = append(chs, ChallengeAuthoring{ID: fmt.Sprintf("c%03d", i), Themes: []string{"alpha"}})
	}
	c := newCatalog([]Pack{{Challenges: chs}})
	if _, _, err := c.Compose([]string{"alpha"}, nil); err == nil {
		t.Fatal("excessive composition accepted")
	}
	c = newCatalog([]Pack{{Challenges: []ChallengeAuthoring{{ID: "a", Themes: []string{"alpha"}, Prerequisites: []string{"b"}}, {ID: "b"}}, Relations: []RelationAuthoring{{From: "b", To: "a", Kind: "recommended_before"}}}})
	if _, _, err := c.Compose([]string{"alpha"}, nil); err == nil {
		t.Fatal("mixed dependency cycle accepted")
	}
}

func TestSelectionLatencyAtV1Scale(t *testing.T) {
	c := syntheticCatalogAtV1Scale(84)
	for _, tc := range []struct {
		name  string
		query func() error
	}{
		{"search", func() error { _, err := c.Search(Query{ThemeIDs: []string{"bench-theme"}}); return err }},
		{"compose", func() error {
			_, _, err := c.Compose([]string{"bench-theme"}, []string{"bench-competency"})
			return err
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			samples := make([]time.Duration, 100)
			for i := range samples {
				start := time.Now()
				if err := tc.query(); err != nil {
					t.Fatal(err)
				}
				samples[i] = time.Since(start)
			}
			sort.Slice(samples, func(i, j int) bool { return samples[i] < samples[j] })
			p95 := samples[94]
			t.Logf("p95=%s (100 queries, 84 indexed challenges)", p95)
			if p95 >= 100*time.Millisecond {
				t.Fatalf("p95 budget exceeded: %s", p95)
			}
		})
	}
}
