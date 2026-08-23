package curriculum

import "slices"

// Item is the sanitized, public projection of any catalog entity, returned
// by Get and List. It never exposes reserved fields (requirement R6).
type Item struct {
	ID    string
	Kind  ItemKind
	Title string
}

// Catalog is an immutable, queryable snapshot built once by Load. Every
// session that starts fixes one Catalog and never observes a different one
// (PROJECT.md §12.1, §12.3 invariant 1).
type Catalog struct {
	themes       map[string]ThemeAuthoring
	concepts     map[string]ConceptAuthoring
	competencies map[string]CompetencyAuthoring
	tracks       map[string]TrackAuthoring
	challenges   map[string]ChallengeAuthoring
}

func newCatalog(packs []Pack) *Catalog {
	c := &Catalog{
		themes:       map[string]ThemeAuthoring{},
		concepts:     map[string]ConceptAuthoring{},
		competencies: map[string]CompetencyAuthoring{},
		tracks:       map[string]TrackAuthoring{},
		challenges:   map[string]ChallengeAuthoring{},
	}
	for _, p := range packs {
		for _, t := range p.Themes {
			c.themes[t.ID] = t
		}
		for _, cn := range p.Concepts {
			c.concepts[cn.ID] = cn
		}
		for _, cp := range p.Competencies {
			c.competencies[cp.ID] = cp
		}
		for _, tr := range p.Tracks {
			c.tracks[tr.ID] = tr
		}
		for _, ch := range p.Challenges {
			c.challenges[ch.ID] = ch
		}
	}
	return c
}

// Get returns the sanitized item with the given ID, regardless of kind
// (requirement R6).
func (c *Catalog) Get(id string) (Item, bool) {
	if t, ok := c.themes[id]; ok {
		return Item{ID: t.ID, Kind: KindTheme, Title: t.Title}, true
	}
	if cn, ok := c.concepts[id]; ok {
		return Item{ID: cn.ID, Kind: KindConcept, Title: cn.Title}, true
	}
	if cp, ok := c.competencies[id]; ok {
		return Item{ID: cp.ID, Kind: KindCompetency, Title: cp.Title}, true
	}
	if tr, ok := c.tracks[id]; ok {
		return Item{ID: tr.ID, Kind: KindTrack, Title: tr.Title}, true
	}
	if ch, ok := c.challenges[id]; ok {
		return Item{ID: ch.ID, Kind: KindChallenge, Title: ch.Title}, true
	}
	return Item{}, false
}

// List returns every sanitized item of the given kind, optionally narrowed
// to a theme (requirement R6: "filtros básicos").
func (c *Catalog) List(kind ItemKind, theme string) []Item {
	var out []Item
	switch kind {
	case KindTheme:
		for _, t := range c.themes {
			out = append(out, Item{ID: t.ID, Kind: KindTheme, Title: t.Title})
		}
	case KindConcept:
		for _, cn := range c.concepts {
			out = append(out, Item{ID: cn.ID, Kind: KindConcept, Title: cn.Title})
		}
	case KindCompetency:
		for _, cp := range c.competencies {
			out = append(out, Item{ID: cp.ID, Kind: KindCompetency, Title: cp.Title})
		}
	case KindTrack:
		for _, tr := range c.tracks {
			if theme != "" && !slices.Contains(tr.Themes, theme) {
				continue
			}
			out = append(out, Item{ID: tr.ID, Kind: KindTrack, Title: tr.Title})
		}
	case KindChallenge:
		for _, ch := range c.challenges {
			if theme != "" && !slices.Contains(ch.Themes, theme) {
				continue
			}
			out = append(out, Item{ID: ch.ID, Kind: KindChallenge, Title: ch.Title})
		}
	}
	return out
}

// Concept returns the authoring record for id, or false when id is not a
// concept (assistance-hints-detours requirement R5).
func (c *Catalog) Concept(id string) (ConceptAuthoring, bool) {
	cn, ok := c.concepts[id]
	return cn, ok
}

// Challenge returns the full authoring record for id with reserved fields
// cleared, or false when id is not a challenge (requirement R6).
func (c *Catalog) Challenge(id string) (ChallengeAuthoring, bool) {
	ch, ok := c.challenges[id]
	if !ok {
		return ChallengeAuthoring{}, false
	}
	ch.Variants = nil // reserved: excluded from the public view
	return ch, true
}
