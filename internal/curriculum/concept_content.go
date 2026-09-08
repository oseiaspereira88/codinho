package curriculum

import (
	"strings"
	"unicode/utf8"
)

func cloneConcept(c ConceptAuthoring) ConceptAuthoring {
	if c.Content != nil {
		copy := *c.Content
		copy.RelationRefs = append([]ConceptRelationReference(nil), copy.RelationRefs...)
		c.Content = &copy
	}
	return c
}

func validateConceptContent(packs []Pack, concepts map[string]bool) []Diagnostic {
	edges := map[RelationAuthoring]bool{}
	for _, p := range packs {
		for _, r := range p.Relations {
			edges[r] = true
		}
	}
	var diags []Diagnostic
	for _, p := range packs {
		for _, c := range p.Concepts {
			if c.Content == nil {
				continue
			}
			content := c.Content
			invalid := func(field, detail string) {
				diags = append(diags, Diagnostic{File: p.File, Item: c.ID, Field: "content." + field, Code: DiagInvalidConceptContent, Detail: detail, Blocking: true})
			}
			if content.Version != 1 {
				invalid("version", "supported content version is 1")
			}
			for _, f := range []struct {
				name, value string
				limit       int
				required    bool
			}{
				{"explanation", content.Explanation, 8000, true},
				{"example.context", content.Example.Context, 200, true},
				{"example.code", content.Example.Code, 4000, true},
				{"example.explanation", content.Example.Explanation, 2000, true},
				{"analogy", content.Analogy, 2000, false},
			} {
				if (f.required && strings.TrimSpace(f.value) == "") || utf8.RuneCountInString(f.value) > f.limit {
					invalid(f.name, "text is missing, blank or exceeds the character limit")
				}
			}
			if len(content.RelationRefs) > 8 {
				invalid("relation_refs", "at most 8 relation references are allowed")
			}
			seen := map[ConceptRelationReference]bool{}
			for _, r := range content.RelationRefs {
				if seen[r] || !concepts[r.ConceptID] || !edges[RelationAuthoring{From: c.ID, To: r.ConceptID, Kind: r.Kind}] {
					invalid("relation_refs", "reference must be a unique existing outgoing relation to a concept")
				}
				seen[r] = true
			}
		}
	}
	return diags
}
