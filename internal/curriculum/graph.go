package curriculum

import "sort"

// RelationKind is one of the eight authored relation types connecting
// catalog items (curriculum-graph-path-recommendation, requirement R1).
// This is the only vocabulary the loader accepts; any other value fails
// explicitly rather than being silently ignored (Compatibility).
type RelationKind string

const (
	RelationRequires          RelationKind = "requires"
	RelationRecommendedBefore RelationKind = "recommended_before"
	RelationRelatesTo         RelationKind = "relates_to"
	RelationContrastsWith     RelationKind = "contrasts_with"
	RelationCommonlyFailsWith RelationKind = "commonly_fails_with"
	RelationAppliesIn         RelationKind = "applies_in"
	RelationDeepensInto       RelationKind = "deepens_into"
	RelationEvidences         RelationKind = "evidences"
)

var knownRelationKinds = map[RelationKind]bool{
	RelationRequires: true, RelationRecommendedBefore: true, RelationRelatesTo: true,
	RelationContrastsWith: true, RelationCommonlyFailsWith: true, RelationAppliesIn: true,
	RelationDeepensInto: true, RelationEvidences: true,
}

// precedenceRelationKinds are the relations requirement R2 requires to be
// acyclic: their direction expresses a real dependency of From on To
// (Decision 3). The other five are symmetric or associative by nature and
// are never cycle-checked.
var precedenceRelationKinds = map[RelationKind]bool{
	RelationRequires: true, RelationRecommendedBefore: true, RelationDeepensInto: true,
}

// IsPrecedence reports whether k implies a dependency order (requirement
// R2, Decision 3).
func (k RelationKind) IsPrecedence() bool { return precedenceRelationKinds[k] }

// Relation is one authored edge between two catalog items.
type Relation struct {
	From string       `json:"from"`
	To   string       `json:"to"`
	Kind RelationKind `json:"kind"`
}

// Graph indexes every authored Relation for lookup by either endpoint,
// deterministically ordered (non-functional requirement: stable results).
type Graph struct {
	out map[string][]Relation
	in  map[string][]Relation
}

func newGraph(packs []Pack) *Graph {
	g := &Graph{out: map[string][]Relation{}, in: map[string][]Relation{}}
	for _, p := range packs {
		for _, r := range p.Relations {
			kind := RelationKind(r.Kind)
			if !knownRelationKinds[kind] {
				continue // rejected at Validate time; never indexed
			}
			rel := Relation{From: r.From, To: r.To, Kind: kind}
			g.out[rel.From] = append(g.out[rel.From], rel)
			g.in[rel.To] = append(g.in[rel.To], rel)
		}
	}
	for id := range g.out {
		sortRelations(g.out[id])
	}
	for id := range g.in {
		sortRelations(g.in[id])
	}
	return g
}

func sortRelations(rels []Relation) {
	sort.Slice(rels, func(i, j int) bool {
		if rels[i].Kind != rels[j].Kind {
			return rels[i].Kind < rels[j].Kind
		}
		if rels[i].To != rels[j].To {
			return rels[i].To < rels[j].To
		}
		return rels[i].From < rels[j].From
	})
}

// Out returns every relation authored FROM id, in deterministic order.
func (g *Graph) Out(id string) []Relation { return g.out[id] }

// In returns every relation authored pointing TO id, in deterministic
// order.
func (g *Graph) In(id string) []Relation { return g.in[id] }

// detectRelationCycles walks only precedence-implying edges (deps) and
// reports a blocking diagnostic for each cycle found, mirroring
// detectPrerequisiteCycles's algorithm (validator.go) so both use the same
// well-tested traversal shape.
func detectRelationCycles(deps map[string][]string) []Diagnostic {
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
			diags = append(diags, Diagnostic{Item: id, Code: DiagRelationCycle, Detail: cyclePath(path, id), Blocking: true})
			return
		}
		state[id] = visiting
		for _, dep := range deps[id] {
			visit(dep, append(path, id))
		}
		state[id] = done
	}

	ids := make([]string, 0, len(deps))
	for id := range deps {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		visit(id, nil)
	}
	return diags
}
