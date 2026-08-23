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
	content, err := svc.ConceptContent("named-types")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content.ID != "named-types" || content.Title != "Named types" {
		t.Fatalf("unexpected content: %+v", content)
	}
}

func TestConceptContentUnknownIDReturnsErrConceptNotFound(t *testing.T) {
	svc := New(newTestCatalog(t))
	_, err := svc.ConceptContent("does-not-exist")
	if !errors.Is(err, ErrConceptNotFound) {
		t.Fatalf("expected ErrConceptNotFound, got %v", err)
	}
}
