package curriculum

import "testing"

func TestLoadValidPackIndexesRelationsBothDirections(t *testing.T) {
	catalog, diags, err := Load("testdata/valid", DefaultLimits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostics: %+v", diags)
	}

	out, in := catalog.Relations("fixture.challenge-one")
	if len(out) != 2 {
		t.Fatalf("expected 2 outgoing relations, got %+v", out)
	}
	if out[0].Kind != RelationRelatesTo || out[0].To != "slice-filter" {
		t.Fatalf("expected relates_to first (sorted by kind), got %+v", out[0])
	}
	if out[1].Kind != RelationRequires || out[1].To != "slice-declaration" {
		t.Fatalf("expected requires second, got %+v", out[1])
	}
	if len(in) != 0 {
		t.Fatalf("expected no incoming relations for fixture.challenge-one, got %+v", in)
	}

	_, inForDecl := catalog.Relations("slice-declaration")
	if len(inForDecl) != 1 || inForDecl[0].From != "fixture.challenge-one" {
		t.Fatalf("expected slice-declaration to see the incoming requires edge, got %+v", inForDecl)
	}
}

func TestLoadRejectsUnknownRelationKind(t *testing.T) {
	assertBlockingDiagnostic(t, "testdata/unknown_relation_kind", DiagUnknownRelationKind)
}

func TestLoadDetectsRelationCycleOnlyForPrecedenceKinds(t *testing.T) {
	assertBlockingDiagnostic(t, "testdata/relation_cycle", DiagRelationCycle)
}

func TestValidateAllowsCyclesAmongNonPrecedenceRelations(t *testing.T) {
	packs := []Pack{{
		ID: "p", SchemaVersion: SchemaVersion,
		Concepts: []ConceptAuthoring{{ID: "a", Title: "A"}, {ID: "b", Title: "B"}},
		Relations: []RelationAuthoring{
			{From: "a", To: "b", Kind: string(RelationRelatesTo)},
			{From: "b", To: "a", Kind: string(RelationRelatesTo)},
		},
	}}
	diags := Validate(packs)
	if hasDiagnostic(diags, DiagRelationCycle) {
		t.Fatalf("relates_to is symmetric by nature and must never be cycle-checked, got %+v", diags)
	}
}

func TestRelationKindIsPrecedenceMatchesDecision3(t *testing.T) {
	precedence := []RelationKind{RelationRequires, RelationRecommendedBefore, RelationDeepensInto}
	for _, k := range precedence {
		if !k.IsPrecedence() {
			t.Errorf("%s should be precedence", k)
		}
	}
	symmetric := []RelationKind{RelationRelatesTo, RelationContrastsWith, RelationCommonlyFailsWith, RelationAppliesIn, RelationEvidences}
	for _, k := range symmetric {
		if k.IsPrecedence() {
			t.Errorf("%s should not be precedence", k)
		}
	}
}
