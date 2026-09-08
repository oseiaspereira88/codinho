package assistance

import (
	"errors"
	"sort"

	"github.com/oseiaspereira88/codinho/internal/curriculum"
)

// ErrConceptNotFound is returned when a concept lookup finds nothing with
// the given ID.
var ErrConceptNotFound = errors.New("assistance: concept not found")

// ConceptContent exposes only public authored concept content. The tutor may
// adapt it; the server never generates prose or consults challenge solutions.
type ConceptContent struct {
	ID            string                   `json:"id"`
	Title         string                   `json:"title"`
	ContentStatus string                   `json:"content_status"`
	Content       *CanonicalConceptContent `json:"content"`
}

type CanonicalConceptContent struct {
	Version     int                                `json:"version"`
	Explanation string                             `json:"explanation"`
	Example     curriculum.ConceptExampleAuthoring `json:"example"`
	Analogy     string                             `json:"analogy,omitempty"`
	Relations   []ConceptRelation                  `json:"relations"`
}

type ConceptRelation struct {
	Kind      string `json:"kind"`
	ConceptID string `json:"concept_id"`
	Title     string `json:"title"`
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
	out := ConceptContent{ID: c.ID, Title: c.Title, ContentStatus: "missing"}
	if c.Content == nil {
		return out, nil
	}
	content := c.Content
	public := &CanonicalConceptContent{Version: content.Version, Explanation: content.Explanation, Example: content.Example, Analogy: content.Analogy, Relations: []ConceptRelation{}}
	for _, ref := range content.RelationRefs {
		if target, ok := s.catalog.Concept(ref.ConceptID); ok {
			public.Relations = append(public.Relations, ConceptRelation{Kind: ref.Kind, ConceptID: ref.ConceptID, Title: target.Title})
		}
	}
	sort.Slice(public.Relations, func(i, j int) bool {
		a, b := public.Relations[i], public.Relations[j]
		if a.Kind != b.Kind {
			return a.Kind < b.Kind
		}
		return a.ConceptID < b.ConceptID
	})
	out.ContentStatus, out.Content = "available", public
	return out, nil
}
