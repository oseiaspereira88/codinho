package curriculum

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// DistributionPolicy is versioned separately from authored learning content.
// Pack groups allow an aggregate planned target without inventing per-pack splits.
type DistributionPolicy struct {
	SchemaVersion int                 `json:"schema_version"`
	KindGroups    map[string][]string `json:"kind_groups,omitempty"`
	Global        TypeDistribution    `json:"global"`
	Packs         []PackDistribution  `json:"packs"`
}

type PackDistribution struct {
	IDs      []string         `json:"ids"`
	Expected TypeDistribution `json:"expected"`
}

func LoadDistributionPolicy(path string) (DistributionPolicy, error) {
	var p DistributionPolicy
	f, err := os.Open(path)
	if err != nil {
		return p, err
	}
	defer f.Close()
	dec := json.NewDecoder(io.LimitReader(f, 1<<20))
	dec.DisallowUnknownFields()
	if err = dec.Decode(&p); err != nil {
		return p, err
	}
	var trailing any
	if err = dec.Decode(&trailing); err != io.EOF {
		return p, fmt.Errorf("distribution must contain exactly one JSON document")
	}
	if err = p.Validate(); err != nil {
		return p, err
	}
	return p, nil
}

func (p DistributionPolicy) Validate() error {
	if p.SchemaVersion != 1 || len(p.Global.ByChallengeKind) == 0 || len(p.Packs) == 0 {
		return fmt.Errorf("distribution requires schema_version 1, global kinds and pack targets")
	}
	seenKinds := map[string]bool{}
	for name, kinds := range p.KindGroups {
		if strings.TrimSpace(name) == "" || len(kinds) == 0 {
			return fmt.Errorf("empty kind group")
		}
		for _, kind := range kinds {
			if kind == name || seenKinds[kind] || strings.TrimSpace(kind) == "" {
				return fmt.Errorf("overlapping or invalid kind groups")
			}
			if _, nested := p.KindGroups[kind]; nested {
				return fmt.Errorf("nested kind groups are not supported")
			}
			seenKinds[kind] = true
		}
	}
	seenIDs := map[string]bool{}
	validateCounts := func(t TypeDistribution) error {
		if len(t.ByChallengeKind) == 0 {
			return fmt.Errorf("pack kind distribution cannot be empty")
		}
		for _, m := range []map[string]int{t.ByChallengeKind, t.ByDifficulty} {
			for k, n := range m {
				if k == "" || n < 0 {
					return fmt.Errorf("distribution keys must be nonempty and counts nonnegative")
				}
			}
		}
		return nil
	}
	if err := validateCounts(p.Global); err != nil {
		return err
	}
	for _, pack := range p.Packs {
		if len(pack.IDs) == 0 {
			return fmt.Errorf("pack target requires IDs")
		}
		for _, id := range pack.IDs {
			if strings.TrimSpace(id) == "" || seenIDs[id] {
				return fmt.Errorf("pack IDs must be nonempty and occur in one target")
			}
			seenIDs[id] = true
		}
		if err := validateCounts(pack.Expected); err != nil {
			return err
		}
	}
	totals := map[string]int{}
	for _, pack := range p.Packs {
		for kind, n := range pack.Expected.ByChallengeKind {
			totals[kind] += n
		}
	}
	if len(CheckTypeDistribution(Coverage{ByChallengeKind: totals}, TypeDistribution{ByChallengeKind: p.Global.ByChallengeKind})) > 0 {
		return fmt.Errorf("pack kind targets must sum to global kind target")
	}
	return nil
}

func (p DistributionPolicy) Check(packs []Pack) []EditorialFinding {
	eligible := EligiblePacks(packs)
	project := func(packs []Pack) Coverage {
		cov := ProjectCoverage(newCatalog(packs))
		for group, kinds := range p.KindGroups {
			for _, kind := range kinds {
				cov.ByChallengeKind[group] += cov.ByChallengeKind[kind]
				delete(cov.ByChallengeKind, kind)
			}
		}
		// Absent grouped kinds must not look like unexpected zero-count categories.
		for kind, n := range cov.ByChallengeKind {
			if n == 0 {
				delete(cov.ByChallengeKind, kind)
			}
		}
		return cov
	}
	out := CheckTypeDistribution(project(eligible), p.Global)
	allowed := map[string]bool{}
	for _, target := range p.Packs {
		selected := []Pack{}
		for _, id := range target.IDs {
			exists := false
			for _, authored := range packs {
				if authored.ID == id {
					exists = true
				}
			}
			if !exists {
				out = append(out, EditorialFinding{Item: "pack:" + id, Rule: RuleTypeDistributionMismatch, Severity: SeverityBlocking, Detail: "policy target pack does not exist in inventory"})
			}
			allowed[id] = true
			for _, pack := range eligible {
				if pack.ID == id {
					selected = append(selected, pack)
				}
			}
		}
		for _, f := range CheckTypeDistribution(project(selected), target.Expected) {
			f.Item = "pack:" + strings.Join(target.IDs, ",") + ":" + f.Item
			out = append(out, f)
		}
	}
	for _, pack := range eligible {
		if !allowed[pack.ID] && (len(pack.Challenges) > 0 || len(pack.Concepts) > 0 || len(pack.Competencies) > 0 || len(pack.Tracks) > 0) {
			out = append(out, EditorialFinding{File: pack.File, Item: pack.ID, Rule: RuleTypeDistributionMismatch, Severity: SeverityBlocking, Detail: "published pack is outside the canonical distribution policy"})
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Item < out[j].Item })
	return out
}
