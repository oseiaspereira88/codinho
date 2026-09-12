package curriculum

import (
	"slices"
	"sort"
	"testing"
)

var foundationPackIDs = map[string]bool{
	"go-first-steps": true,
	"go-core":        true,
	"go-data-text":   true,
	"go-errors":      true,
	"go-io":          true,
	"go-testing":     true,
	"go-type-design": true,
}

const (
	historicalFoundationPrototype = "go-data.slice-filter-preserve-input"
	outOfScopePrototype           = "go-debug.slice-off-by-one"
	minimumFoundationNodes        = 300
)

var essentialFoundationCompetencies = []string{
	"inspect-go-toolchain-context",
	"control-closure-capture",
	"define-package-api-boundary",
	"control-test-cleanup-lifecycle",
	"clone-struct-with-independent-slice",
	"return-true-nil-interface",
	"distinguish-zero-value-from-absence",
	"wrap-sentinel-with-context",
	"testing-design-deterministic-time",
	"testing-detect-aliasing",
}

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

// TestFoundationCatalogAudit is the executable ledger for the seven
// foundational packs. It intentionally keeps inventory, variants and
// prototypes separate: a variant or a preserved historical prototype cannot
// pay the contextualized-node target for canonical foundation content.
func TestFoundationCatalogAudit(t *testing.T) {
	packs, diags, err := LoadPacks("../../packs", DefaultLimits)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range diags {
		if d.Blocking {
			t.Fatalf("blocking catalog diagnostic: %+v", d)
		}
	}

	byPack := map[string]Pack{}
	for _, p := range packs {
		byPack[p.ID] = p
	}
	if len(byPack) != 8 {
		t.Fatalf("loaded packs = %d, want seven foundation packs plus debugging", len(byPack))
	}
	for id := range foundationPackIDs {
		if _, ok := byPack[id]; !ok {
			t.Fatalf("missing foundation pack %s", id)
		}
	}

	inventory := ProjectCoverage(NewCatalogFromPacks(packs))
	publication := ProjectPublicationCoverage(packs)
	if inventory.Concepts != 105 || inventory.Competencies != 64 || inventory.Challenges != 54 || inventory.StepNodes != 343 {
		t.Fatalf("global inventory = %+v, want concepts=105 competencies=64 challenges=54 nodes=343", inventory)
	}
	if publication.Published.Challenges != 0 || publication.Eligible.Challenges != 0 {
		t.Fatalf("published projection = %+v, want no published or eligible challenges", publication)
	}

	foundationConcepts := map[string]bool{}
	foundationCompetencies := map[string]bool{}
	baseChallenges := 0
	baseNodes := 0
	creditedNodes := 0
	variants := 0
	kindCounts := map[string]int{}
	competencyContexts := map[string]map[string]bool{}
	for _, p := range packs {
		isFoundation := foundationPackIDs[p.ID]
		for _, concept := range p.Concepts {
			if !isFoundation {
				continue
			}
			foundationConcepts[concept.ID] = true
		}
		for _, competency := range p.Competencies {
			if !isFoundation {
				continue
			}
			foundationCompetencies[competency.ID] = true
		}
		for _, ch := range p.Challenges {
			if !isFoundation {
				continue
			}
			for _, competency := range append(append([]string{}, ch.Competencies.Primary...), ch.Competencies.Secondary...) {
				if competencyContexts[competency] == nil {
					competencyContexts[competency] = map[string]bool{}
				}
				competencyContexts[competency][ch.ID] = true
			}
			if ch.VariantOf != "" {
				variants++
				continue
			}
			baseChallenges++
			kindCounts[ch.Kind]++
			nodes := countStepNodes(ch.Layers)
			baseNodes += nodes
			if ch.ID != historicalFoundationPrototype {
				creditedNodes += nodes
			}
		}
	}

	if len(foundationConcepts) != 104 || len(foundationCompetencies) != 60 {
		t.Fatalf("foundation coverage = concepts=%d competencies=%d, want 104/60", len(foundationConcepts), len(foundationCompetencies))
	}
	if baseChallenges != 45 || variants != 8 || baseNodes != 281 || creditedNodes != 280 {
		t.Fatalf("foundation ledger = bases=%d variants=%d nodes=%d credited=%d, want 45/8/281/280", baseChallenges, variants, baseNodes, creditedNodes)
	}
	if kindCounts["atomic"] != 33 || kindCounts["combined"] != 10 || kindCounts["functional_slice"] != 2 {
		t.Fatalf("foundation base kinds including preserved prototype = %+v, want atomic=33 combined=10 functional_slice=2", kindCounts)
	}

	prototypePacks := map[string]string{
		historicalFoundationPrototype: "go-first-steps",
		outOfScopePrototype:           "go-debugging",
	}
	for id, expectedPackID := range prototypePacks {
		found := false
		for _, ch := range byPack[expectedPackID].Challenges {
			if ch.ID != id {
				continue
			}
			found = true
			if ch.VariantOf != "" {
				t.Errorf("prototype %s must remain a base challenge, got variant of %s", id, ch.VariantOf)
			}
		}
		if !found {
			t.Fatalf("prototype %s is missing from expected pack %s", id, expectedPackID)
		}
	}

	competencies := append([]string{}, essentialFoundationCompetencies...)
	sort.Strings(competencies)
	for _, competency := range competencies {
		contexts := make([]string, 0, len(competencyContexts[competency]))
		for challenge := range competencyContexts[competency] {
			contexts = append(contexts, challenge)
		}
		sort.Strings(contexts)
		if len(contexts) < 2 {
			t.Errorf("essential competency %s has %d contexts, want at least 2: %v", competency, len(contexts), contexts)
		}
		t.Logf("essential competency context: %s => %v", competency, contexts)
	}
	t.Logf("foundation audit: concepts=%d competencies=%d bases=%d variants=%d nodes=%d credited_nodes=%d", len(foundationConcepts), len(foundationCompetencies), baseChallenges, variants, baseNodes, creditedNodes)
	if creditedNodes < minimumFoundationNodes {
		t.Logf("foundation contextualized-node gate remains below minimum: got=%d want>=%d deficit=%d", creditedNodes, minimumFoundationNodes, minimumFoundationNodes-creditedNodes)
	}
}
