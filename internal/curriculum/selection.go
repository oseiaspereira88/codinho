package curriculum

import (
	"errors"
	"slices"
	"sort"
	"strings"
)

var ErrInvalidSelection = errors.New("curriculum: invalid or excessive selection")

// Selection explains coverage only for the returned candidates, never hidden
// matches beyond a result cap. Modes are choices, not automatic transitions.
type Selection struct {
	Coverage             string         `json:"coverage"`
	CoveredThemeIDs      []string       `json:"covered_theme_ids"`
	MissingThemeIDs      []string       `json:"missing_theme_ids"`
	CoveredCompetencyIDs []string       `json:"covered_competency_ids"`
	MissingCompetencyIDs []string       `json:"missing_competency_ids"`
	Complete             bool           `json:"complete"`
	Matches              []SubjectMatch `json:"matches"`
	Modes                []string       `json:"modes"`
}
type SubjectMatch struct {
	ID            string   `json:"id"`
	ThemeIDs      []string `json:"theme_ids"`
	CompetencyIDs []string `json:"competency_ids"`
}

func NormalizeSubjects(ids []string) ([]string, error) {
	if len(ids) > 16 || (ids != nil && len(ids) == 0) {
		return nil, ErrInvalidSelection
	}
	out := append([]string(nil), ids...)
	for _, id := range out {
		if strings.TrimSpace(id) == "" || len(id) > 200 {
			return nil, ErrInvalidSelection
		}
	}
	sort.Strings(out)
	return slices.Compact(out), nil
}
func intersect(a, b []string) []string {
	out := []string{}
	for _, id := range a {
		if slices.Contains(b, id) {
			out = append(out, id)
		}
	}
	return out
}
func (c *Catalog) SubjectMatch(id string, themes, comps []string) SubjectMatch {
	m := SubjectMatch{ID: id, ThemeIDs: []string{}, CompetencyIDs: []string{}}
	if ch, ok := c.challenges[id]; ok {
		m.ThemeIDs = intersect(themes, ch.Themes)
		m.CompetencyIDs = intersect(comps, append(append([]string{}, ch.Competencies.Primary...), ch.Competencies.Secondary...))
	}
	if tr, ok := c.tracks[id]; ok {
		m.ThemeIDs = intersect(themes, tr.Themes)
	}
	return m
}
func (c *Catalog) Selection(ids, themes, comps []string) Selection {
	s := Selection{Coverage: "none", CoveredThemeIDs: []string{}, MissingThemeIDs: []string{}, CoveredCompetencyIDs: []string{}, MissingCompetencyIDs: []string{}, Matches: []SubjectMatch{}, Modes: []string{"track", "subjects", "challenge", "generate"}}
	for _, id := range ids {
		m := c.SubjectMatch(id, themes, comps)
		s.Matches = append(s.Matches, m)
		if len(m.ThemeIDs)+len(m.CompetencyIDs) > 0 {
			if s.Coverage == "none" {
				s.Coverage = "partial"
			}
		}
		if len(themes)+len(comps) > 0 && len(m.ThemeIDs) == len(themes) && len(m.CompetencyIDs) == len(comps) {
			s.Coverage = "total"
		}
	}
	for _, t := range themes {
		found := false
		for _, m := range s.Matches {
			found = found || slices.Contains(m.ThemeIDs, t)
		}
		if found {
			s.CoveredThemeIDs = append(s.CoveredThemeIDs, t)
		} else {
			s.MissingThemeIDs = append(s.MissingThemeIDs, t)
		}
	}
	for _, t := range comps {
		found := false
		for _, m := range s.Matches {
			found = found || slices.Contains(m.CompetencyIDs, t)
		}
		if found {
			s.CoveredCompetencyIDs = append(s.CoveredCompetencyIDs, t)
		} else {
			s.MissingCompetencyIDs = append(s.MissingCompetencyIDs, t)
		}
	}
	s.Complete = len(themes)+len(comps) > 0 && len(s.MissingThemeIDs)+len(s.MissingCompetencyIDs) == 0
	return s
}
