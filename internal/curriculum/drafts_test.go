package curriculum

import (
	"errors"
	"strings"
	"testing"
)

const validDraft = `schema_version: 1
id: synthetic-draft
version: 1.0.0
competencies:
  - id: synthetic-skill
    title: Practice
challenges:
  - schema_version: 1
    id: synthetic-challenge
    version: 1.0.0
    competencies:
      primary: [synthetic-skill]
    acceptance: [Explain the outcome]
    layers:
      - id: understanding
        macro_steps:
          - id: explain
            kind: micro
            instruction: {objective: Explain the outcome}
`

func TestParseDraftQuarantinesUnreviewedContent(t *testing.T) {
	p, v, err := ParseDraft([]byte(validDraft))
	if err != nil {
		t.Fatalf("%v: %+v", err, v)
	}
	if p.Publication.Status != "draft" || p.Challenges[0].Publication.Status != "draft" {
		t.Fatal("not quarantined")
	}
	if len(PublishedPacks([]Pack{p})) != 0 {
		t.Fatal("draft leaked into published catalog")
	}
	if _, ok := DraftCatalog(p).Challenge("synthetic-challenge"); !ok {
		t.Fatal("explicit catalog missing draft")
	}
}
func TestParseDraftRejectsUnsafeOrInvalidAuthoring(t *testing.T) {
	for name, raw := range map[string]string{
		"unknown field":     validDraft + "unknown: true\n",
		"schema":            strings.Replace(validDraft, "schema_version: 1", "schema_version: 99", 1),
		"challenge schema":  strings.ReplaceAll(validDraft, "schema_version: 1", "schema_version: 99"),
		"review claim":      validDraft + "publication: {status: draft, reviewed_by: someone}\n",
		"published claim":   validDraft + "publication: {status: published}\n",
		"playtest claim":    validDraft + "publication: {playtested: true}\n",
		"missing reference": strings.Replace(validDraft, "primary: [synthetic-skill]", "primary: [missing]", 1),
		"editorial gap":     strings.Replace(validDraft, "    acceptance: [Explain the outcome]\n", "", 1),
		"oversize":          strings.Repeat("x", MaxDraftBytes+1),
		"malformed":         "[",
	} {
		t.Run(name, func(t *testing.T) {
			_, v, err := ParseDraft([]byte(raw))
			if !errors.Is(err, ErrInvalidDraft) {
				t.Fatalf("accepted invalid draft: %+v", v)
			}
		})
	}
}
