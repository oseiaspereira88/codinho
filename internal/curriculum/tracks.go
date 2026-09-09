package curriculum

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

const MaxTrackChallenges = 100

type Path struct {
	TrackID       string   `json:"track_id,omitempty"`
	CompositionID string   `json:"composition_id,omitempty"`
	ChallengeIDs  []string `json:"challenge_ids"`
	Versions      []string `json:"versions"`
}

func cloneTrack(t TrackAuthoring) TrackAuthoring {
	t.Themes = append([]string(nil), t.Themes...)
	t.ChallengeIDs = append([]string(nil), t.ChallengeIDs...)
	return t
}
func (c *Catalog) Track(id string) (TrackAuthoring, bool) {
	t, ok := c.tracks[id]
	return cloneTrack(t), ok
}
func PathIdentity(challenges []ChallengeAuthoring) string {
	raw, _ := json.Marshal(challenges)
	return fmt.Sprintf("composition_v1_%x", sha256.Sum256(raw))
}

// dependencies uses the established From-depends-on-To relation direction.
func (c *Catalog) dependencies(id string) []string {
	deps := []string{}
	for _, p := range c.challenges[id].Prerequisites {
		if _, ok := c.challenges[p]; ok {
			deps = append(deps, p)
		}
	}
	for _, r := range c.graph.Out(id) {
		if r.Kind == RelationRequires || r.Kind == RelationRecommendedBefore {
			if _, ok := c.challenges[r.To]; ok {
				deps = append(deps, r.To)
			}
		}
	}
	sort.Strings(deps)
	return deps
}
func (c *Catalog) ResolvePath(ids []string) ([]ChallengeAuthoring, error) {
	if len(ids) == 0 || len(ids) > MaxTrackChallenges {
		return nil, ErrInvalidSelection
	}
	seen := map[string]bool{}
	out := make([]ChallengeAuthoring, 0, len(ids))
	for _, id := range ids {
		if strings.TrimSpace(id) == "" || len(id) > 200 {
			return nil, ErrInvalidSelection
		}
		ch, ok := c.Challenge(id)
		if !ok || seen[id] {
			return nil, ErrInvalidSelection
		}
		for _, p := range c.dependencies(id) {
			if !seen[p] {
				return nil, ErrInvalidSelection
			}
		}
		seen[id] = true
		out = append(out, ch)
	}
	return out, nil
}
func describePath(chs []ChallengeAuthoring) Path {
	p := Path{CompositionID: PathIdentity(chs), ChallengeIDs: []string{}, Versions: []string{}}
	for _, ch := range chs {
		p.ChallengeIDs = append(p.ChallengeIDs, ch.ID)
		p.Versions = append(p.Versions, ch.Version)
	}
	return p
}

// Compose includes every matching challenge and its prerequisite closure. It is
// a deterministic graph path, not a shortest-path or pedagogical optimizer.
func (c *Catalog) Compose(themes, comps []string) (Path, Selection, error) {
	var err error
	themes, err = NormalizeSubjects(themes)
	if err != nil {
		return Path{}, Selection{}, err
	}
	comps, err = NormalizeSubjects(comps)
	if err != nil {
		return Path{}, Selection{}, err
	}
	roots := []string{}
	for id := range c.challenges {
		m := c.SubjectMatch(id, themes, comps)
		if len(m.ThemeIDs)+len(m.CompetencyIDs) > 0 {
			roots = append(roots, id)
		}
	}
	sort.Strings(roots)
	selected := map[string]bool{}
	visiting := map[string]bool{}
	visitedCount := 0
	var include func(string) error
	include = func(id string) error {
		if visiting[id] {
			return ErrInvalidSelection
		}
		if selected[id] {
			return nil
		}
		visitedCount++
		if visitedCount > MaxTrackChallenges {
			return ErrInvalidSelection
		}
		visiting[id] = true
		for _, p := range c.dependencies(id) {
			if err := include(p); err != nil {
				return err
			}
		}
		visiting[id] = false
		selected[id] = true
		if len(selected) > MaxTrackChallenges {
			return ErrInvalidSelection
		}
		return nil
	}
	for _, id := range roots {
		if err := include(id); err != nil {
			return Path{}, Selection{}, err
		}
	}
	ids := []string{}
	done := map[string]bool{}
	for len(ids) < len(selected) {
		ready := []string{}
		for id := range selected {
			if done[id] {
				continue
			}
			ok := true
			for _, p := range c.dependencies(id) {
				ok = ok && done[p]
			}
			if ok {
				ready = append(ready, id)
			}
		}
		sort.Strings(ready)
		if len(ready) == 0 {
			return Path{}, Selection{}, ErrInvalidSelection
		}
		id := ready[0]
		ids = append(ids, id)
		done[id] = true
	}
	if len(ids) == 0 {
		return Path{ChallengeIDs: []string{}, Versions: []string{}}, c.Selection(ids, themes, comps), nil
	}
	chs, err := c.ResolvePath(ids)
	if err != nil {
		return Path{}, Selection{}, err
	}
	return describePath(chs), c.Selection(ids, themes, comps), nil
}
func validateTracks(packs []Pack) []Diagnostic {
	c := newCatalog(packs)
	var ds []Diagnostic
	for _, p := range packs {
		for _, t := range p.Tracks {
			if t.ChallengeIDs == nil {
				continue
			}
			if _, err := c.ResolvePath(t.ChallengeIDs); err != nil {
				ds = append(ds, Diagnostic{File: p.File, Item: t.ID, Field: "challenge_ids", Code: DiagMissingReference, Detail: "track requires 1–100 unique existing challenges in dependency order", Blocking: true})
			}
		}
	}
	return ds
}
