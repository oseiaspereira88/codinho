package assistance

import (
	"errors"

	"github.com/oseiaspereira88/codinho/internal/curriculum"
)

// ErrConceptNotFound is returned when a concept lookup finds nothing with
// the given ID.
var ErrConceptNotFound = errors.New("assistance: concept not found")

// ConceptContent is the canonical, sanitized concept record
// concept_content_get exposes. It carries only what the catalog actually
// authors (PROJECT.md §7.2: id and title); adapting language, relations,
// analogies and examples to the learner's profile is the calling agent's
// job (PROJECT.md §15.7, non-goal: "redigir explicações abertas dentro do
// MCP").
type ConceptContent struct {
	ID    string
	Title string
}

// Service resolves canonical, catalog-backed assistance content. It never
// mutates session state: internal/session owns the ladder ratchet, hint
// events and detour lifecycle, calling ladder.go and detour.go directly
// under its own lock.
type Service struct {
	catalog *curriculum.Catalog
}

// New wraps an already-loaded catalog.
func New(catalog *curriculum.Catalog) *Service {
	return &Service{catalog: catalog}
}

// ConceptContent returns the canonical record for concept id (requirement
// R5).
func (s *Service) ConceptContent(id string) (ConceptContent, error) {
	c, ok := s.catalog.Concept(id)
	if !ok {
		return ConceptContent{}, ErrConceptNotFound
	}
	return ConceptContent{ID: c.ID, Title: c.Title}, nil
}
