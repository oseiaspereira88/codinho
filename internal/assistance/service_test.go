package assistance

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/oseiaspereira88/codinho/internal/curriculum"
)

func newTestCatalog(t *testing.T) *curriculum.Catalog {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "manifest.yaml"), []byte("schema_version: 1\npacks:\n  - pack.yaml\n"), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pack := `schema_version: 1
id: fixture-pack
version: 1.0.0
concepts:
  - id: named-types
    title: Named types
  - id: block-scope
    title: Block scope
  - id: variable-shadowing
    title: Variable shadowing
    content:
      version: 1
      explanation: "Inner declaration hides outer variable in enclosed scope."
      example:
        context: "meteorologia"
        code: "x := 1\nif true {\n\tx := 2\n\t_ = x\n}"
        explanation: "Inner x shadows outer x."
      analogy: "Like an umbrella blocking the sky view."
      relation_refs:
        - kind: relates_to
          concept_id: named-types
        - kind: applies_in
          concept_id: block-scope
relations:
  - from: variable-shadowing
    to: named-types
    kind: relates_to
  - from: variable-shadowing
    to: block-scope
    kind: applies_in
`
	if err := os.WriteFile(filepath.Join(dir, "pack.yaml"), []byte(pack), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	catalog, diags, err := curriculum.Load(dir, curriculum.DefaultLimits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostics: %+v", diags)
	}
	return catalog
}

func TestConceptContentReturnsCanonicalRecord(t *testing.T) {
	svc := New(newTestCatalog(t))

	// Legacy concept without authored content
	content, err := svc.ConceptContent("named-types")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content.ID != "named-types" || content.Title != "Named types" {
		t.Fatalf("unexpected content: %+v", content)
	}
	if content.ContentStatus != "missing" {
		t.Fatalf("expected content_status=missing, got %q", content.ContentStatus)
	}
	if content.Content != nil {
		t.Fatalf("expected nil content, got %+v", content.Content)
	}

	// Concept with authored content
	shadowing, err := svc.ConceptContent("variable-shadowing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if shadowing.ContentStatus != "available" {
		t.Fatalf("expected content_status=available, got %q", shadowing.ContentStatus)
	}
	if shadowing.Content == nil {
		t.Fatalf("expected non-nil content")
	}
	if shadowing.Content.Version != 1 {
		t.Fatalf("expected version 1, got %d", shadowing.Content.Version)
	}
	if shadowing.Content.Explanation != "Inner declaration hides outer variable in enclosed scope." {
		t.Fatalf("unexpected explanation: %q", shadowing.Content.Explanation)
	}
	if shadowing.Content.Example.Context != "meteorologia" {
		t.Fatalf("unexpected example context: %q", shadowing.Content.Example.Context)
	}
	if shadowing.Content.Analogy != "Like an umbrella blocking the sky view." {
		t.Fatalf("unexpected analogy: %q", shadowing.Content.Analogy)
	}

	// Relations ordered deterministically by kind asc, then concept_id asc
	if len(shadowing.Content.Relations) != 2 {
		t.Fatalf("expected 2 relations, got %d", len(shadowing.Content.Relations))
	}
	if shadowing.Content.Relations[0].Kind != "applies_in" || shadowing.Content.Relations[0].ConceptID != "block-scope" {
		t.Fatalf("expected first relation to be applies_in / block-scope, got %+v", shadowing.Content.Relations[0])
	}
	if shadowing.Content.Relations[0].Title != "Block scope" {
		t.Fatalf("expected resolved title 'Block scope', got %q", shadowing.Content.Relations[0].Title)
	}
	if shadowing.Content.Relations[1].Kind != "relates_to" || shadowing.Content.Relations[1].ConceptID != "named-types" {
		t.Fatalf("expected second relation to be relates_to / named-types, got %+v", shadowing.Content.Relations[1])
	}
	if shadowing.Content.Relations[1].Title != "Named types" {
		t.Fatalf("expected resolved title 'Named types', got %q", shadowing.Content.Relations[1].Title)
	}

	// Immutability: modifying returned content does not affect next call
	shadowing.Content.Explanation = "modified"
	shadowing2, _ := svc.ConceptContent("variable-shadowing")
	if shadowing2.Content.Explanation == "modified" {
		t.Fatalf("expected service output to be immutable")
	}
}

func TestConceptContentUnknownIDReturnsErrConceptNotFound(t *testing.T) {
	svc := New(newTestCatalog(t))
	_, err := svc.ConceptContent("does-not-exist")
	if !errors.Is(err, ErrConceptNotFound) {
		t.Fatalf("expected ErrConceptNotFound, got %v", err)
	}
}
