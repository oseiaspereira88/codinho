package curriculum_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/oseiaspereira88/codinho/internal/curriculum"
)

func TestConceptContent(t *testing.T) {
	t.Run("ValidAndLegacyConcepts", func(t *testing.T) {
		packData, err := os.ReadFile("../../testdata/catalog-quality/concept-content.yaml")
		if err != nil {
			t.Fatalf("reading concept-content.yaml: %v", err)
		}

		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "manifest.yaml"), []byte("schema_version: 1\npacks: [pack.yaml]\n"), 0600); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "pack.yaml"), packData, 0600); err != nil {
			t.Fatal(err)
		}

		catalog, diags, err := curriculum.Load(dir, curriculum.DefaultLimits)
		if err != nil {
			t.Fatalf("unexpected error loading valid concept content: %v", err)
		}
		for _, d := range diags {
			if d.Blocking {
				t.Fatalf("unexpected blocking diagnostic: %+v", d)
			}
		}

		// Valid concept with full content
		c1, ok := catalog.Concept("weather-forecast-concept")
		if !ok {
			t.Fatalf("concept not found")
		}
		if c1.Content == nil {
			t.Fatalf("expected content to be present")
		}
		if c1.Content.Version != 1 {
			t.Fatalf("expected version 1, got %d", c1.Content.Version)
		}
		if len(c1.Content.RelationRefs) != 1 {
			t.Fatalf("expected 1 relation ref, got %d", len(c1.Content.RelationRefs))
		}

		// Legacy concept without content
		cLegacy, ok := catalog.Concept("legacy-weather-concept")
		if !ok {
			t.Fatalf("legacy concept not found")
		}
		if cLegacy.Content != nil {
			t.Fatalf("expected nil content for legacy concept")
		}

		// Immutability: modifying returned concept content should not modify catalog content
		c1.Content.Explanation = "modified explanation"
		c1Again, _ := catalog.Concept("weather-forecast-concept")
		if c1Again.Content.Explanation == "modified explanation" {
			t.Fatalf("expected concept content to be immutable via defensive cloning")
		}
	})

	t.Run("ValidationRejections", func(t *testing.T) {
		makePack := func(modifyConcept func(*curriculum.ConceptAuthoring), relations []curriculum.RelationAuthoring) curriculum.Pack {
			c := curriculum.ConceptAuthoring{
				ID:    "c1",
				Title: "Concept 1",
				Content: &curriculum.ConceptContentAuthoring{
					Version:     1,
					Explanation: "Valid explanation",
					Example: curriculum.ConceptExampleAuthoring{
						Context:     "meteorologia",
						Code:        "var x = 1",
						Explanation: "declares x",
					},
					Analogy: "like a scale",
				},
			}
			if modifyConcept != nil {
				modifyConcept(&c)
			}
			return curriculum.Pack{
				SchemaVersion: 1,
				ID:            "test-pack",
				Version:       "1.0.0",
				Concepts: []curriculum.ConceptAuthoring{
					c,
					{ID: "c2", Title: "Concept 2"},
				},
				Competencies: []curriculum.CompetencyAuthoring{
					{ID: "comp1", Title: "Competency 1"},
				},
				Relations: relations,
			}
		}

		assertHasInvalidContentDiag := func(t *testing.T, pack curriculum.Pack, expectedSubstr string) {
			t.Helper()
			diags := curriculum.Validate([]curriculum.Pack{pack})
			var found bool
			for _, d := range diags {
				if d.Code == curriculum.DiagInvalidConceptContent && strings.Contains(d.Detail, expectedSubstr) {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("expected DiagInvalidConceptContent containing %q, got: %+v", expectedSubstr, diags)
			}
		}

		t.Run("UnsupportedVersion", func(t *testing.T) {
			p := makePack(func(c *curriculum.ConceptAuthoring) {
				c.Content.Version = 2
			}, nil)
			assertHasInvalidContentDiag(t, p, "version is 1")
		})

		t.Run("BlankExplanation", func(t *testing.T) {
			p := makePack(func(c *curriculum.ConceptAuthoring) {
				c.Content.Explanation = "   \n\t "
			}, nil)
			assertHasInvalidContentDiag(t, p, "text is missing, blank or exceeds")
		})

		t.Run("BlankExampleContext", func(t *testing.T) {
			p := makePack(func(c *curriculum.ConceptAuthoring) {
				c.Content.Example.Context = ""
			}, nil)
			assertHasInvalidContentDiag(t, p, "text is missing, blank or exceeds")
		})

		t.Run("BlankExampleCode", func(t *testing.T) {
			p := makePack(func(c *curriculum.ConceptAuthoring) {
				c.Content.Example.Code = " "
			}, nil)
			assertHasInvalidContentDiag(t, p, "text is missing, blank or exceeds")
		})

		t.Run("BlankExampleExplanation", func(t *testing.T) {
			p := makePack(func(c *curriculum.ConceptAuthoring) {
				c.Content.Example.Explanation = ""
			}, nil)
			assertHasInvalidContentDiag(t, p, "text is missing, blank or exceeds")
		})

		t.Run("ExplanationExceedsLimit", func(t *testing.T) {
			p := makePack(func(c *curriculum.ConceptAuthoring) {
				c.Content.Explanation = strings.Repeat("a", 8001)
			}, nil)
			assertHasInvalidContentDiag(t, p, "text is missing, blank or exceeds")
		})

		t.Run("ExampleContextExceedsLimit", func(t *testing.T) {
			p := makePack(func(c *curriculum.ConceptAuthoring) {
				c.Content.Example.Context = strings.Repeat("b", 201)
			}, nil)
			assertHasInvalidContentDiag(t, p, "text is missing, blank or exceeds")
		})

		t.Run("ExampleCodeExceedsLimit", func(t *testing.T) {
			p := makePack(func(c *curriculum.ConceptAuthoring) {
				c.Content.Example.Code = strings.Repeat("c", 4001)
			}, nil)
			assertHasInvalidContentDiag(t, p, "text is missing, blank or exceeds")
		})

		t.Run("ExampleExplanationExceedsLimit", func(t *testing.T) {
			p := makePack(func(c *curriculum.ConceptAuthoring) {
				c.Content.Example.Explanation = strings.Repeat("d", 2001)
			}, nil)
			assertHasInvalidContentDiag(t, p, "text is missing, blank or exceeds")
		})

		t.Run("AnalogyExceedsLimit", func(t *testing.T) {
			p := makePack(func(c *curriculum.ConceptAuthoring) {
				c.Content.Analogy = strings.Repeat("e", 2001)
			}, nil)
			assertHasInvalidContentDiag(t, p, "text is missing, blank or exceeds")
		})

		t.Run("TooManyRelationRefs", func(t *testing.T) {
			p := makePack(func(c *curriculum.ConceptAuthoring) {
				for i := 0; i < 9; i++ {
					c.Content.RelationRefs = append(c.Content.RelationRefs, curriculum.ConceptRelationReference{
						Kind:      "relates_to",
						ConceptID: "c2",
					})
				}
			}, nil)
			assertHasInvalidContentDiag(t, p, "at most 8 relation references")
		})

		t.Run("DuplicateRelationRef", func(t *testing.T) {
			p := makePack(func(c *curriculum.ConceptAuthoring) {
				c.Content.RelationRefs = []curriculum.ConceptRelationReference{
					{Kind: "relates_to", ConceptID: "c2"},
					{Kind: "relates_to", ConceptID: "c2"},
				}
			}, []curriculum.RelationAuthoring{
				{From: "c1", To: "c2", Kind: "relates_to"},
			})
			assertHasInvalidContentDiag(t, p, "reference must be a unique existing outgoing relation to a concept")
		})

		t.Run("RelationRefPointsToNonConcept", func(t *testing.T) {
			p := makePack(func(c *curriculum.ConceptAuthoring) {
				c.Content.RelationRefs = []curriculum.ConceptRelationReference{
					{Kind: "relates_to", ConceptID: "comp1"},
				}
			}, []curriculum.RelationAuthoring{
				{From: "c1", To: "comp1", Kind: "relates_to"},
			})
			assertHasInvalidContentDiag(t, p, "reference must be a unique existing outgoing relation to a concept")
		})

		t.Run("RelationRefNonExistingOutgoingEdge", func(t *testing.T) {
			p := makePack(func(c *curriculum.ConceptAuthoring) {
				c.Content.RelationRefs = []curriculum.ConceptRelationReference{
					{Kind: "relates_to", ConceptID: "c2"},
				}
			}, nil) // no relation in pack
			assertHasInvalidContentDiag(t, p, "reference must be a unique existing outgoing relation to a concept")
		})
	})
}
