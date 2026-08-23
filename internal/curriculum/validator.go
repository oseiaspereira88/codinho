package curriculum

import (
	"sort"
	"strings"
)

// DiagnosticCode is a stable, comparable identifier for a validation
// finding, independent of its human-readable Detail (requirement R4).
type DiagnosticCode string

const (
	DiagIncompatibleSchemaVersion DiagnosticCode = "incompatible_schema_version"
	DiagPathEscapesPack           DiagnosticCode = "path_escapes_pack"
	DiagUnreadablePack            DiagnosticCode = "unreadable_pack"
	DiagMalformedPack             DiagnosticCode = "malformed_pack"
	DiagDuplicateID               DiagnosticCode = "duplicate_id"
	DiagMissingReference          DiagnosticCode = "missing_reference"
	DiagPrerequisiteCycle         DiagnosticCode = "prerequisite_cycle"
	DiagCriteriaWithoutEvidence   DiagnosticCode = "criteria_without_evidence"
	DiagUnknownRelationKind       DiagnosticCode = "unknown_relation_kind"
	DiagRelationCycle             DiagnosticCode = "relation_cycle"
)

// Diagnostic reports one finding with the file, item and field it came
// from, so authors can locate and fix it (requirement R4). Blocking
// diagnostics prevent the catalog from being materialized at all.
type Diagnostic struct {
	File     string
	Item     string
	Field    string
	Code     DiagnosticCode
	Detail   string
	Blocking bool
}

// Validate checks cross-pack references, prerequisite cycles and
// criteria-without-evidence across every already schema-version-checked
// pack (requirement R3). It assumes each pack decoded successfully; the
// loader is responsible for schema-version and decode diagnostics.
func Validate(packs []Pack) []Diagnostic {
	var diags []Diagnostic

	ids := map[string]string{} // id -> file, for duplicate detection
	concepts := map[string]bool{}
	competencies := map[string]bool{}
	themes := map[string]bool{}
	challenges := map[string]bool{}

	claim := func(id, file string) *Diagnostic {
		if id == "" {
			return nil
		}
		if owner, seen := ids[id]; seen {
			return &Diagnostic{File: file, Item: id, Code: DiagDuplicateID, Detail: "already declared in " + owner, Blocking: true}
		}
		ids[id] = file
		return nil
	}

	for _, p := range packs {
		for _, t := range p.Themes {
			if d := claim(t.ID, p.File); d != nil {
				diags = append(diags, *d)
			}
			themes[t.ID] = true
		}
		for _, c := range p.Concepts {
			if d := claim(c.ID, p.File); d != nil {
				diags = append(diags, *d)
			}
			concepts[c.ID] = true
		}
		for _, c := range p.Competencies {
			if d := claim(c.ID, p.File); d != nil {
				diags = append(diags, *d)
			}
			competencies[c.ID] = true
		}
		for _, t := range p.Tracks {
			if d := claim(t.ID, p.File); d != nil {
				diags = append(diags, *d)
			}
		}
		for _, ch := range p.Challenges {
			if d := claim(ch.ID, p.File); d != nil {
				diags = append(diags, *d)
			}
			challenges[ch.ID] = true
		}
	}

	prereqs := map[string][]string{}

	for _, p := range packs {
		for _, tr := range p.Tracks {
			for _, themeID := range tr.Themes {
				if !themes[themeID] {
					diags = append(diags, Diagnostic{File: p.File, Item: tr.ID, Field: "themes", Code: DiagMissingReference, Detail: themeID})
				}
			}
		}
		for _, ch := range p.Challenges {
			prereqs[ch.ID] = ch.Prerequisites
			for _, themeID := range ch.Themes {
				if !themes[themeID] {
					diags = append(diags, Diagnostic{File: p.File, Item: ch.ID, Field: "themes", Code: DiagMissingReference, Detail: themeID})
				}
			}
			for _, id := range ch.Competencies.Primary {
				if !competencies[id] {
					diags = append(diags, Diagnostic{File: p.File, Item: ch.ID, Field: "competencies.primary", Code: DiagMissingReference, Detail: id})
				}
			}
			for _, id := range ch.Competencies.Secondary {
				if !competencies[id] {
					diags = append(diags, Diagnostic{File: p.File, Item: ch.ID, Field: "competencies.secondary", Code: DiagMissingReference, Detail: id})
				}
			}
			for _, id := range ch.Prerequisites {
				if !concepts[id] && !challenges[id] {
					diags = append(diags, Diagnostic{File: p.File, Item: ch.ID, Field: "prerequisites", Code: DiagMissingReference, Detail: id})
				}
			}
			for _, layer := range ch.Layers {
				validateSteps(p.File, ch.ID, layer.MacroSteps, concepts, &diags)
			}
		}
	}

	diags = append(diags, detectPrerequisiteCycles(prereqs)...)
	diags = append(diags, validateRelations(packs, themes, concepts, competencies, challenges)...)

	return diags
}

// validateRelations checks every authored RelationAuthoring (curriculum-
// graph-path-recommendation, requirement R1): an unknown Kind fails
// explicitly (Compatibility) rather than being silently indexed by
// newGraph, both endpoints must reference a real catalog item, and
// precedence-implying kinds (Decision 3) must never cycle (requirement
// R2).
func validateRelations(packs []Pack, themes, concepts, competencies, challenges map[string]bool) []Diagnostic {
	exists := func(id string) bool {
		return themes[id] || concepts[id] || competencies[id] || challenges[id]
	}

	var diags []Diagnostic
	deps := map[string][]string{}
	for _, p := range packs {
		for _, r := range p.Relations {
			kind := RelationKind(r.Kind)
			if !knownRelationKinds[kind] {
				diags = append(diags, Diagnostic{
					File: p.File, Item: r.From, Field: "relations", Code: DiagUnknownRelationKind,
					Detail: r.Kind, Blocking: true,
				})
				continue
			}
			if !exists(r.From) {
				diags = append(diags, Diagnostic{File: p.File, Item: r.From, Field: "relations.from", Code: DiagMissingReference, Detail: r.From})
			}
			if !exists(r.To) {
				diags = append(diags, Diagnostic{File: p.File, Item: r.From, Field: "relations.to", Code: DiagMissingReference, Detail: r.To})
			}
			if kind.IsPrecedence() {
				deps[r.From] = append(deps[r.From], r.To)
			}
		}
	}
	diags = append(diags, detectRelationCycles(deps)...)
	return diags
}

func validateSteps(file, challengeID string, steps []StepAuthoring, concepts map[string]bool, diags *[]Diagnostic) {
	for _, s := range steps {
		for _, id := range s.Concepts {
			if !concepts[id] {
				*diags = append(*diags, Diagnostic{File: file, Item: challengeID + "/" + s.ID, Field: "concepts", Code: DiagMissingReference, Detail: id})
			}
		}
		if len(s.Criteria) > 0 && len(s.Evidence.Strategies) == 0 {
			*diags = append(*diags, Diagnostic{File: file, Item: challengeID + "/" + s.ID, Field: "evidence.strategies", Code: DiagCriteriaWithoutEvidence, Blocking: true})
		}
		validateSteps(file, challengeID, s.Children, concepts, diags)
	}
}

// detectPrerequisiteCycles runs DFS with a recursion stack over the
// challenge prerequisite graph and reports one diagnostic per distinct
// cycle found (requirement R3).
func detectPrerequisiteCycles(prereqs map[string][]string) []Diagnostic {
	const (
		unvisited = 0
		visiting  = 1
		done      = 2
	)
	state := map[string]int{}
	var diags []Diagnostic

	var visit func(id string, path []string)
	visit = func(id string, path []string) {
		switch state[id] {
		case done:
			return
		case visiting:
			diags = append(diags, Diagnostic{Item: id, Code: DiagPrerequisiteCycle, Detail: cyclePath(path, id), Blocking: true})
			return
		}
		state[id] = visiting
		for _, dep := range prereqs[id] {
			if _, isChallenge := prereqs[dep]; isChallenge {
				visit(dep, append(path, id))
			}
		}
		state[id] = done
	}

	ids := make([]string, 0, len(prereqs))
	for id := range prereqs {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		visit(id, nil)
	}
	return diags
}

func cyclePath(path []string, closing string) string {
	var out strings.Builder
	for _, p := range path {
		out.WriteString(p)
		out.WriteString(" -> ")
	}
	out.WriteString(closing)
	return out.String()
}
