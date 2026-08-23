// Package application wires the pure domain (internal/learning), curriculum
// (internal/curriculum) and persistence (internal/eventstore,
// internal/evidence) packages into the use cases the MCP server exposes. It
// is the only layer that is allowed to know about all four.
package application

import (
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/session"
)

// ErrNotFound is returned when a catalog lookup finds nothing with the
// given ID. It is the same sentinel internal/session uses for a missing
// challenge, so mapError's errors.Is checks work regardless of which layer
// raised it.
var ErrNotFound = session.ErrNotFound

// CatalogService adapts the immutable curriculum catalog for MCP tools.
type CatalogService struct {
	catalog *curriculum.Catalog
}

// NewCatalogService wraps an already-loaded, validated catalog.
func NewCatalogService(catalog *curriculum.Catalog) *CatalogService {
	return &CatalogService{catalog: catalog}
}

// SearchQuery narrows a catalog search by any combination of fields
// (requirement R3; R5, R6: "filtros básicos" from mcp-stdio-foundation
// remain the Kind/Theme subset). An empty field is not applied.
type SearchQuery struct {
	Kind          curriculum.ItemKind
	Theme         string
	Text          string
	Competency    string
	Difficulty    string
	ChallengeKind string
	MaxMinutes    int
	Prerequisite  string
}

// SearchResult is what Search returns: matching items plus whether the
// result-count cap actually dropped matches (Security: never silently
// claim completeness).
type SearchResult struct {
	Items     []curriculum.Item
	Truncated bool
}

// Search returns every sanitized item matching q. An empty Kind defaults to
// challenges, the catalog's primary browsing surface (requirement R3).
func (s *CatalogService) Search(q SearchQuery) (SearchResult, error) {
	result, err := s.catalog.Search(curriculum.Query{
		Kind: q.Kind, Theme: q.Theme, Text: q.Text, Competency: q.Competency,
		Difficulty: q.Difficulty, ChallengeKind: q.ChallengeKind, MaxMinutes: q.MaxMinutes,
		Prerequisite: q.Prerequisite,
	})
	if err != nil {
		return SearchResult{}, err
	}
	return SearchResult{Items: result.Items, Truncated: result.Truncated}, nil
}

// Get returns the sanitized item with id, or ErrNotFound.
func (s *CatalogService) Get(id string) (curriculum.Item, error) {
	item, ok := s.catalog.Get(id)
	if !ok {
		return curriculum.Item{}, ErrNotFound
	}
	return item, nil
}

// Relations returns id's outgoing and incoming relations, or ErrNotFound
// when id does not exist in the catalog at all (requirement R7).
func (s *CatalogService) Relations(id string) (out, in []curriculum.Relation, err error) {
	if _, ok := s.catalog.Get(id); !ok {
		return nil, nil, ErrNotFound
	}
	out, in = s.catalog.Relations(id)
	return out, in, nil
}
