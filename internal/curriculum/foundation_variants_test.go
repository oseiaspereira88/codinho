package curriculum

import (
	"slices"
	"testing"
)

// Authored variants must retain their origin and cannot inflate base coverage.
func TestFoundationVariantOrigins(t *testing.T) {
	packs, _, err := LoadPacks("../../packs", DefaultLimits)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range packs {
		if p.ID != "go-io" && p.ID != "go-testing" {
			continue
		}
		byID := map[string]ChallengeAuthoring{}
		for _, ch := range p.Challenges {
			byID[ch.ID] = ch
		}
		variants := 0
		for _, ch := range p.Challenges {
			if ch.VariantOf == "" {
				continue
			}
			variants++
			parent, ok := byID[ch.VariantOf]
			if !ok || parent.VariantOf != "" {
				t.Fatalf("%s: missing base origin %s", ch.ID, ch.VariantOf)
			}
			if ch.Canonical || ch.Kind != "atomic" || len(ch.Competencies.Primary) != 1 {
				t.Fatalf("%s: invalid atomic variant classification", ch.ID)
			}
			skill := ch.Competencies.Primary[0]
			if !slices.Contains(parent.Competencies.Primary, skill) && !slices.Contains(parent.Competencies.Secondary, skill) {
				t.Fatalf("%s: skill %s absent from origin", ch.ID, skill)
			}
		}
		if variants != 4 {
			t.Fatalf("%s: got %d contextual variants, want 4", p.ID, variants)
		}
	}
}
