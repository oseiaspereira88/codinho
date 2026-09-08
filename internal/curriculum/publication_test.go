package curriculum

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func publicationPack() Pack {
	meta := PublicationAuthoring{Status: StatusPublished, Author: "synthetic-author", ReviewedBy: "synthetic-reviewer", Playtested: true}
	return Pack{ID: "fixture", Version: "1.0.0", Publication: meta, Competencies: []CompetencyAuthoring{{ID: "competency"}}, Challenges: []ChallengeAuthoring{
		{SchemaVersion: 1, ID: "canonical", Version: "1.0.0", Title: "Challenge", Difficulty: "foundational", Kind: "atomic", Canonical: true, Publication: meta, Competencies: CompetencyRefs{Primary: []string{"competency"}}, Acceptance: []string{"observable"}},
		{ID: "draft", Version: "1.0.0", Kind: "atomic"},
	}}
}

func TestPublicationCoverageVisibilityAndPrivateProof(t *testing.T) {
	p := publicationPack()
	prototype := p.Challenges[0]
	prototype.ID = "prototype"
	prototype.Canonical = false
	variant := p.Challenges[0]
	variant.ID = "variant"
	variant.VariantOf = "canonical"
	p.Challenges = append(p.Challenges, prototype, variant)
	p.Challenges[0].Variants = []string{"private-alternative"}
	p.Challenges[0].Validation = &CheckValidationAuthoring{ReferenceFixture: []FixtureFileAuthoring{{Path: "secret.go", Content: "private-reference"}}}
	p.Relations = []RelationAuthoring{{From: "canonical", To: "draft", Kind: "relates_to"}, {From: "canonical", To: "competency", Kind: "evidences"}}
	cov := ProjectPublicationCoverage([]Pack{p})
	if cov.Inventory.Challenges != 4 || cov.Drafts.Challenges != 1 || cov.Published.Challenges != 3 || cov.Eligible.Challenges != 1 {
		t.Fatalf("coverage: %+v", cov)
	}
	public := newCatalog(PublishedPacks([]Pack{p}))
	if _, ok := public.Get("draft"); ok {
		t.Fatal("draft leaked")
	}
	out, _ := public.Relations("canonical")
	if len(out) != 1 || out[0].To != "competency" {
		t.Fatalf("relations leaked draft: %+v", out)
	}
	for _, catalog := range []*Catalog{public, newCatalog([]Pack{p})} {
		ch, _ := catalog.Challenge("canonical")
		if ch.Validation != nil || ch.Variants != nil {
			t.Fatal("private record leaked")
		}
		data, _ := json.Marshal(catalog.Challenges())
		if strings.Contains(string(data), "private-") {
			t.Fatalf("private payload leaked: %s", data)
		}
	}
	data, _ := json.Marshal(p.Challenges[0])
	if strings.Contains(string(data), "private-reference") {
		t.Fatal("private solution serialized into pinned content")
	}
	if p.Challenges[0].Validation == nil {
		t.Fatal("projection mutated administrative input")
	}
}

func TestPublicationRejectsInvalidMetadataAndHiddenDependencies(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func(*Pack)
	}{
		{"review", func(p *Pack) { p.Challenges[0].Publication.ReviewedBy = " SYNTHETIC-AUTHOR " }},
		{"blank-author", func(p *Pack) { p.Publication.Author = "  " }},
		{"playtest", func(p *Pack) { p.Challenges[0].Publication.Playtested = false }},
		{"draft-pack", func(p *Pack) { p.Publication = PublicationAuthoring{} }},
		{"unknown-status", func(p *Pack) { p.Challenges[1].Publication.Status = "ready" }},
		{"prerequisite", func(p *Pack) { p.Challenges[0].Prerequisites = []string{"draft"} }},
		{"graph", func(p *Pack) { p.Relations = []RelationAuthoring{{From: "canonical", To: "draft", Kind: "requires"}} }},
		{"missing-proof", func(p *Pack) { p.Challenges[0].Checks = []CheckAuthoring{{ID: "check", Runner: "internal_ast"}} }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := publicationPack()
			tc.mutate(&p)
			if len(ValidatePublication([]Pack{p})) == 0 {
				t.Fatal("invalid publication accepted")
			}
		})
	}
	p := publicationPack()
	if ds := ValidatePublication([]Pack{p}); len(ds) != 0 {
		t.Fatalf("valid publication rejected: %+v", ds)
	}
}

func TestDistributionProductionPolicyAndUnexpectedTypes(t *testing.T) {
	policy, err := LoadDistributionPolicy("../../packs/distribution.json")
	if err != nil {
		t.Fatal(err)
	}
	if policy.Global.ByChallengeKind["atomic"] != 42 {
		t.Fatal("wrong product target")
	}
	fixture, err := LoadDistributionPolicy("../../testdata/catalog-quality/distribution.json")
	if err != nil {
		t.Fatal(err)
	}
	p := publicationPack()
	if f := fixture.Check([]Pack{p}); len(f) != 0 {
		t.Fatalf("valid distribution: %+v", f)
	}
	unexpected := p.Challenges[0]
	unexpected.ID = "extra"
	unexpected.Kind = "unplanned"
	p.Challenges = append(p.Challenges, unexpected)
	if f := fixture.Check([]Pack{p}); len(f) != 2 {
		t.Fatalf("global AND pack must flag unexpected type: %+v", f)
	}
	p = publicationPack()
	p.ID = "prototype-pack"
	if len(fixture.Check([]Pack{p})) == 0 {
		t.Fatal("unplanned pack accepted")
	}
	fixture.Packs[0].Expected.ByChallengeKind["atomic"] = 2
	if fixture.Validate() == nil {
		t.Fatal("inconsistent global/pack policy accepted")
	}
}

func TestLoaderRejectsUnknownManifestAndSymlinkEscape(t *testing.T) {
	dir := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.yaml")
	if err := os.WriteFile(outside, []byte("schema_version: 1\nid: outside\nversion: 1.0.0\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(dir, "pack.yaml")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	write := func(s string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, "manifest.yaml"), []byte(s), 0600); err != nil {
			t.Fatal(err)
		}
	}
	write("schema_version: 1\npacks: [pack.yaml]\nunknown: true\n")
	if c, _, err := Load(dir, DefaultLimits); err == nil && c != nil {
		t.Fatal("unknown manifest field accepted")
	}
	write("schema_version: 1\npacks: [pack.yaml]\n")
	if c, diags, err := Load(dir, DefaultLimits); err != nil || c != nil || !hasDiagnostic(diags, DiagPathEscapesPack) {
		t.Fatalf("symlink escape: %v %+v %v", c, diags, err)
	}
	write("schema_version: 2\npacks: []\n")
	if c, diags, err := Load(dir, DefaultLimits); err != nil || c != nil || !hasDiagnostic(diags, DiagIncompatibleSchemaVersion) {
		t.Fatalf("future manifest was not blocking: %v %+v %v", c, diags, err)
	}
}

func TestDistributionRejectsAbsentZeroTarget(t *testing.T) {
	policy, err := LoadDistributionPolicy("../../testdata/catalog-quality/distribution.json")
	if err != nil {
		t.Fatal(err)
	}
	policy.Packs = append(policy.Packs, PackDistribution{IDs: []string{"absent"}, Expected: TypeDistribution{ByChallengeKind: map[string]int{"atomic": 0}}})
	if err = policy.Validate(); err != nil {
		t.Fatal(err)
	}
	if len(policy.Check([]Pack{publicationPack()})) == 0 {
		t.Fatal("absent pack silently accepted for zero target")
	}
}
