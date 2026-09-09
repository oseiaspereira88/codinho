package curriculum

import (
	"errors"
	"slices"
	"sort"
	"strings"
)

// maxSearchTextLen bounds a caller-supplied text query so it can never cost
// more than a bounded substring scan (Security: "query de texto é dado e
// deve ter limites de tamanho e custo").
const maxSearchTextLen = 200

// ErrSearchTextTooLong is returned when Query.Text exceeds maxSearchTextLen.
var ErrSearchTextTooLong = errors.New("curriculum: search text exceeds the maximum length")

// maxSearchResults bounds how many items Search ever returns, independent
// of how many match, so a broad query can never be used to exhaust a
// caller's memory (Security).
const maxSearchResults = 100

// Query narrows a catalog search by any combination of fields (requirement
// R3). A zero-valued field is not applied; it never means "match nothing".
// Fields beyond Kind/Theme/Text only apply to KindChallenge, the catalog's
// primary browsing surface.
type Query struct {
	Kind          ItemKind
	ThemeIDs      []string
	Text          string
	CompetencyIDs []string
	Difficulty    string
	ChallengeKind string // the challenge's authored `kind` (atomic, ...), distinct from ItemKind
	MaxMinutes    int    // 0 = no limit
	Prerequisite  string // matches a challenge that lists this ID as a prerequisite
}

// Result is what Search returns: the matching items, capped at
// maxSearchResults, and whether that cap actually dropped matches
// (Security: never silently claim completeness).
type Result struct {
	Items     []Item
	Truncated bool
	Selection Selection
}

// Search returns every item matching q, in deterministic order (ID
// tie-break, non-functional requirement: stable ranking).
func (c *Catalog) Search(q Query) (Result, error) {
	if len(q.Text) > maxSearchTextLen {
		return Result{}, ErrSearchTextTooLong
	}
	var err error
	q.ThemeIDs, err = NormalizeSubjects(q.ThemeIDs)
	if err != nil {
		return Result{}, err
	}
	q.CompetencyIDs, err = NormalizeSubjects(q.CompetencyIDs)
	if err != nil {
		return Result{}, err
	}
	kind := q.Kind
	if kind == "" {
		kind = KindChallenge
	}
	text := strings.ToLower(strings.TrimSpace(q.Text))

	var out []Item
	if kind == KindChallenge {
		for _, ch := range c.challenges {
			m := c.SubjectMatch(ch.ID, q.ThemeIDs, q.CompetencyIDs)
			if len(q.ThemeIDs)+len(q.CompetencyIDs) > 0 && len(m.ThemeIDs)+len(m.CompetencyIDs) == 0 {
				continue
			}
			if !matchesChallenge(ch, q, text) {
				continue
			}
			out = append(out, Item{ID: ch.ID, Kind: KindChallenge, Title: ch.Title})
		}
	} else {
		for _, item := range c.List(kind, "") {
			if len(q.ThemeIDs)+len(q.CompetencyIDs) > 0 {
				m := c.SubjectMatch(item.ID, q.ThemeIDs, q.CompetencyIDs)
				if len(m.ThemeIDs)+len(m.CompetencyIDs) == 0 {
					continue
				}
			}
			if text != "" && !strings.Contains(strings.ToLower(item.Title), text) {
				continue
			}
			out = append(out, item)
		}
	}

	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	truncated := len(out) > maxSearchResults
	if truncated {
		out = out[:maxSearchResults]
	}
	ids := make([]string, len(out))
	for i, item := range out {
		ids[i] = item.ID
	}
	return Result{Items: out, Truncated: truncated, Selection: c.Selection(ids, q.ThemeIDs, q.CompetencyIDs)}, nil
}

func matchesChallenge(ch ChallengeAuthoring, q Query, text string) bool {
	if q.Difficulty != "" && ch.Difficulty != q.Difficulty {
		return false
	}
	if q.ChallengeKind != "" && ch.Kind != q.ChallengeKind {
		return false
	}
	if q.MaxMinutes > 0 && ch.EstimatedMinutes > q.MaxMinutes {
		return false
	}
	if q.Prerequisite != "" && !slices.Contains(ch.Prerequisites, q.Prerequisite) {
		return false
	}
	if text != "" && !strings.Contains(strings.ToLower(ch.Title), text) && !strings.Contains(strings.ToLower(ch.Brief), text) {
		return false
	}
	return true
}
