package curriculum

import (
	"fmt"
	"github.com/oseiaspereira88/codinho/internal/checks"
	"strings"
)

// PublishedPacks is the only catalog projection used by normal MCP composition.
// Input must have passed ValidatePublication; relations cannot reveal hidden IDs.
func PublishedPacks(packs []Pack) []Pack {
	out := make([]Pack, 0, len(packs))
	visible := map[string]bool{}
	for _, p := range packs {
		if p.Publication.Status != StatusPublished {
			continue
		}
		copy := p
		copy.Challenges = nil
		copy.Relations = nil
		for _, ch := range p.Challenges {
			if ch.Publication.Status == StatusPublished {
				ch.Validation = nil
				copy.Challenges = append(copy.Challenges, ch)
				visible[ch.ID] = true
			}
		}
		for _, t := range p.Themes {
			visible[t.ID] = true
		}
		for _, c := range p.Concepts {
			visible[c.ID] = true
		}
		for _, c := range p.Competencies {
			visible[c.ID] = true
		}
		for _, t := range p.Tracks {
			visible[t.ID] = true
		}
		out = append(out, copy)
	}
	for i := range out {
		for _, p := range packs {
			if p.ID != out[i].ID {
				continue
			}
			for _, r := range p.Relations {
				if visible[r.From] && visible[r.To] {
					out[i].Relations = append(out[i].Relations, r)
				}
			}
		}
	}
	return out
}

// PublishedAuthoringPacks retains private proof input for the administrative gate.
func PublishedAuthoringPacks(packs []Pack) []Pack {
	var out []Pack
	for _, p := range packs {
		if p.Publication.Status != StatusPublished {
			continue
		}
		copy := p
		copy.Challenges = nil
		for _, ch := range p.Challenges {
			if ch.Publication.Status == StatusPublished {
				copy.Challenges = append(copy.Challenges, ch)
			}
		}
		out = append(out, copy)
	}
	return out
}

func EligiblePacks(packs []Pack) []Pack {
	out := PublishedPacks(packs)
	for i := range out {
		var eligible []ChallengeAuthoring
		for _, ch := range out[i].Challenges {
			if ch.Canonical && ch.VariantOf == "" {
				eligible = append(eligible, ch)
			}
		}
		out[i].Challenges = eligible
	}
	return out
}

// CheckValidationProblems validates proof declarations without executing code.
// A missing declaration is returned even for drafts when --checks is requested.
func CheckValidationProblems(ch ChallengeAuthoring) []string {
	if len(ch.Checks) == 0 {
		return nil
	}
	v := ch.Validation
	if v == nil {
		return []string{"checks require private validation metadata and reference_fixture"}
	}
	var out []string
	if len(ch.Fixture) == 0 && len(v.BaselineFixture) == 0 {
		out = append(out, "check has no baseline fixture or reproducible alternative")
	}
	if len(v.BaselineFixture) > 0 && strings.TrimSpace(v.Justification) == "" {
		out = append(out, "alternative baseline_fixture requires justification")
	}
	if len(v.ReferenceFixture) == 0 {
		out = append(out, "reference_fixture is required")
	}
	expected := map[string]bool{}
	checkIDs := map[string]bool{}
	for _, c := range ch.Checks {
		if _, err := checks.Resolve(checks.CheckSpec{ID: c.ID, Runner: c.Runner, Package: c.Package, TestPattern: c.TestPattern, Timeout: c.Timeout}); err != nil {
			out = append(out, "unresolvable check "+c.ID)
		}
		if c.ID == "" || checkIDs[c.ID] {
			out = append(out, "check IDs must be nonempty and unique")
		}
		checkIDs[c.ID] = true
		if c.Network {
			out = append(out, "editorial proof requires network module fetching denied")
		}
	}
	for _, e := range v.Expectations {
		if !checkIDs[e.CheckID] || expected[e.CheckID] {
			out = append(out, "expectation check_id must match exactly one declared check")
		}
		expected[e.CheckID] = true
		switch e.Baseline {
		case "pass", "test_failure", "compile_failure", "format_failure", "analysis_failure", "syntax_failure":
		default:
			out = append(out, "unknown baseline expectation")
		}
		if e.Reference != "pass" {
			out = append(out, "reference expectation must be pass")
		}
	}
	for _, c := range ch.Checks {
		if !expected[c.ID] {
			out = append(out, "missing expectation for check "+c.ID)
		}
	}
	return out
}

// ValidatePublication rejects invalid publication before any normal MCP service
// is constructed. It never runs authored code and never promotes missing metadata.
func ValidatePublication(packs []Pack) []Diagnostic {
	var out []Diagnostic
	add := func(file, item, field, detail string) {
		out = append(out, Diagnostic{File: file, Item: item, Field: field, Code: "invalid_publication", Detail: detail, Blocking: true})
	}
	public := newCatalog(PublishedPacks(packs))
	all := newCatalog(packs)
	packIDs := map[string]bool{}
	metadata := func(file, id string, p PublicationAuthoring) {
		if p.Status != "" && p.Status != "draft" && p.Status != StatusPublished {
			add(file, id, "publication.status", "unknown publication status")
		}
		for _, f := range checkPublicationMetadata(file, ChallengeAuthoring{ID: id, Publication: p}) {
			add(file, id, "publication", string(f.Rule))
		}
	}
	for _, p := range packs {
		if p.ID == "" || packIDs[p.ID] {
			add(p.File, p.ID, "id", "pack IDs must be nonempty and unique")
		}
		packIDs[p.ID] = true
		metadata(p.File, p.ID, p.Publication)
		for _, ch := range p.Challenges {
			metadata(p.File, ch.ID, ch.Publication)
			if ch.VariantOf != "" {
				if ch.VariantOf == ch.ID {
					add(p.File, ch.ID, "variant_of", "variant cannot reference itself")
				}
				if _, ok := all.Challenge(ch.VariantOf); !ok {
					add(p.File, ch.ID, "variant_of", "variant parent does not exist")
				}
			}
			if ch.Validation != nil {
				for _, files := range [][]FixtureFileAuthoring{ch.Validation.BaselineFixture, ch.Validation.ReferenceFixture} {
					ds := validateFixture(p.File, ChallengeAuthoring{ID: ch.ID, Fixture: files})
					for i := range ds {
						ds[i].Field = "validation.fixture"
					}
					out = append(out, ds...)
				}
			}
			if ch.Publication.Status != StatusPublished {
				continue
			}
			if ch.SchemaVersion != SchemaVersion {
				add(p.File, ch.ID, "schema_version", "published challenge requires supported schema version")
			}
			if strings.TrimSpace(ch.ID) == "" || strings.TrimSpace(ch.Title) == "" || strings.TrimSpace(ch.Difficulty) == "" {
				add(p.File, ch.ID, "identity", "published challenge requires id, title and difficulty")
			}
			switch ch.Kind {
			case "atomic", "combined", "functional_slice", "debug", "debugging", "refactoring", "code_review", "systems_mission", "simulation":
			default:
				add(p.File, ch.ID, "kind", "unknown published challenge kind")
			}
			if p.Publication.Status != StatusPublished {
				add(p.File, ch.ID, "publication", "published challenge requires published containing pack")
			}
			for _, problem := range CheckValidationProblems(ch) {
				add(p.File, ch.ID, "validation", problem)
			}
			for _, f := range RunEditorialChecks([]Pack{{File: p.File, Challenges: []ChallengeAuthoring{ch}}}) {
				if f.Severity == SeverityBlocking {
					add(p.File, ch.ID, "editorial", string(f.Rule))
				}
			}
			refs := append([]string{}, ch.Themes...)
			refs = append(refs, ch.Competencies.Primary...)
			refs = append(refs, ch.Competencies.Secondary...)
			refs = append(refs, ch.Prerequisites...)
			if ch.VariantOf != "" {
				refs = append(refs, ch.VariantOf)
			}
			walkLayers(ch.Layers, func(s StepAuthoring) { refs = append(refs, s.Concepts...) })
			for _, id := range refs {
				if _, ok := public.Get(id); !ok {
					add(p.File, ch.ID, "references", fmt.Sprintf("reference %q is not published", id))
				}
			}
		}
		if p.Publication.Status == StatusPublished {
			for _, t := range p.Tracks {
				for _, id := range t.Themes {
					if _, ok := public.Get(id); !ok {
						add(p.File, t.ID, "themes", "track references unpublished theme")
					}
				}
			}
			for _, r := range p.Relations {
				if _, ok := public.Get(r.From); !ok {
					continue
				}
				if RelationKind(r.Kind).IsPrecedence() {
					if _, ok := public.Get(r.To); !ok {
						add(p.File, r.From, "relations", "published item depends on unpublished relation target")
					}
				}
			}
		}
	}
	return out
}
