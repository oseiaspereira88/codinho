package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"gopkg.in/yaml.v3"
)

func writePublicationPack(t *testing.T, root string, p curriculum.Pack) {
	t.Helper()
	dir := filepath.Join(root, "packs")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	data, err := yaml.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	for name, data := range map[string][]byte{"manifest.yaml": []byte("schema_version: 1\npacks: [pack.yaml]\n"), "pack.yaml": data} {
		if err := os.WriteFile(filepath.Join(dir, name), data, 0600); err != nil {
			t.Fatal(err)
		}
	}
}

func syntheticPublicationPack() curriculum.Pack {
	meta := curriculum.PublicationAuthoring{Status: "published", Author: "synthetic-author", ReviewedBy: "synthetic-reviewer", Playtested: true}
	ch := curriculum.ChallengeAuthoring{SchemaVersion: 1, ID: "public-challenge", Version: "1.0.0", Title: "Public practice", Kind: "atomic", Difficulty: "foundational", Canonical: true, Publication: meta, Acceptance: []string{"Observable result"}, Competencies: curriculum.CompetencyRefs{Primary: []string{"competency"}}, Layers: []curriculum.LayerAuthoring{{ID: "layer", MacroSteps: []curriculum.StepAuthoring{{ID: "step", Kind: "micro", Instruction: curriculum.InstructionAuthoring{Objective: "Practice one action"}}}}}, Fixture: []curriculum.FixtureFileAuthoring{{Path: "main.go", Content: "package main\n"}}, Checks: []curriculum.CheckAuthoring{{ID: "parse", Runner: "internal_ast", Package: "main.go"}}, Validation: &curriculum.CheckValidationAuthoring{ReferenceFixture: []curriculum.FixtureFileAuthoring{{Path: "main.go", Content: "package main\n// PRIVATE-REFERENCE-SENTINEL\n"}}, Expectations: []curriculum.CheckExpectationAuthoring{{CheckID: "parse", Baseline: "pass", Reference: "pass"}}}}
	draft := ch
	draft.ID = "hidden-draft"
	draft.Title = "Draft practice"
	draft.Publication = curriculum.PublicationAuthoring{}
	return curriculum.Pack{SchemaVersion: 1, ID: "fixture", Version: "1.0.0", Publication: meta, Competencies: []curriculum.CompetencyAuthoring{{ID: "competency", Title: "Competency"}}, Challenges: []curriculum.ChallengeAuthoring{ch, draft}, Relations: []curriculum.RelationAuthoring{{From: ch.ID, To: draft.ID, Kind: "relates_to"}}}
}

func TestPublicationIntegrityOverRealStdio(t *testing.T) {
	bin := buildCodinhoBinary(t)
	root := t.TempDir()
	pack := syntheticPublicationPack()
	writePublicationPack(t, root, pack)
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	connect := func(authoring bool) *mcp.ClientSession {
		args := []string{"serve"}
		if authoring {
			args = append(args, "--authoring")
		}
		cmd := exec.CommandContext(ctx, bin, args...)
		cmd.Dir = root
		cs, err := mcp.NewClient(&mcp.Implementation{Name: "publication-test", Version: "1"}, nil).Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
		if err != nil {
			t.Fatal(err)
		}
		return cs
	}
	cs := connect(false)
	for _, req := range []struct {
		tool string
		args map[string]any
	}{
		{"catalog_search", map[string]any{}},
		{"learning_path_recommend", map[string]any{"competency_id": "competency"}},
		{"concept_relations_get", map[string]any{"id": "public-challenge"}},
	} {
		env := callTool(ctx, t, cs, req.tool, req.args)
		data, _ := json.Marshal(env)
		if strings.Contains(string(data), "hidden-draft") || strings.Contains(string(data), "PRIVATE-REFERENCE-SENTINEL") {
			t.Fatalf("%s leaked private content: %s", req.tool, data)
		}
		if req.tool == "catalog_search" && !strings.Contains(string(data), "public-challenge") {
			t.Fatalf("published content missing: %s", data)
		}
	}
	for _, req := range []struct {
		tool string
		args map[string]any
	}{
		{"catalog_get", map[string]any{"id": "hidden-draft"}},
		{"session_start", map[string]any{"challenge_id": "hidden-draft"}},
	} {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: req.tool, Arguments: req.args})
		if err != nil {
			t.Fatal(err)
		}
		env, ok := res.StructuredContent.(map[string]any)
		if !ok || env["status"] == "ok" {
			t.Fatalf("draft accepted through %s: %+v", req.tool, res)
		}
	}
	cs.Close()
	cs = connect(true)
	env := callTool(ctx, t, cs, "catalog_search", map[string]any{})
	data, _ := json.Marshal(env)
	if !strings.Contains(string(data), "hidden-draft") {
		t.Fatal("authoring mode omitted drafts")
	}
	started := callTool(ctx, t, cs, "session_start", map[string]any{"challenge_id": "hidden-draft", "depth": "micro"})
	sessionID := started["session_id"]
	cs.Close()
	cs = connect(false)
	callTool(ctx, t, cs, "session_get", map[string]any{"session_id": sessionID})
	cs.Close()
	log, err := os.ReadFile(filepath.Join(root, ".codinho", "state", "events.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(log), "PRIVATE-REFERENCE-SENTINEL") {
		t.Fatal("reference solution leaked into pinned session event")
	}
	// A claimed publication without review must fail before even starting MCP.
	pack.Challenges[0].Publication.ReviewedBy = ""
	writePublicationPack(t, root, pack)
	cmd := exec.CommandContext(ctx, bin, "serve")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err == nil || !strings.Contains(string(out), "invalid_publication") {
		t.Fatalf("invalid publication started: %s %v", out, err)
	}
}

func TestPublicationIntegrityCLI(t *testing.T) {
	bin := buildCodinhoBinary(t)
	t.Run("published-proof", func(t *testing.T) {
		root := t.TempDir()
		p := syntheticPublicationPack()
		writePublicationPack(t, root, p)
		policy, err := os.ReadFile("../../testdata/catalog-quality/distribution.json")
		if err != nil {
			t.Fatal(err)
		}
		policyPath := filepath.Join(root, "packs", "distribution.json")
		if err = os.WriteFile(policyPath, policy, 0600); err != nil {
			t.Fatal(err)
		}
		cmd := exec.Command(bin, "catalog", "validate", "--published-checks", "--distribution", policyPath, "--json")
		cmd.Dir = root
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("valid published gate: %s %v", out, err)
		}
		var result struct {
			Publication curriculum.PublicationCoverage
			Proof       struct {
				DeclaredChecks int `json:"declared_checks"`
				VerifiedChecks int `json:"verified_checks"`
			}
		}
		if err = json.Unmarshal(out, &result); err != nil {
			t.Fatal(err)
		}
		if result.Publication.Inventory.Challenges != 2 || result.Publication.Eligible.Challenges != 1 || result.Proof.DeclaredChecks != 1 || result.Proof.VerifiedChecks != 1 {
			t.Fatalf("incorrect report: %s", out)
		}
		if strings.Contains(string(out), "PRIVATE-REFERENCE-SENTINEL") {
			t.Fatal("CLI leaked reference")
		}
		p.Challenges[0].Validation.ReferenceFixture[0].Content = "invalid syntax"
		writePublicationPack(t, root, p)
		cmd = exec.Command(bin, "catalog", "validate", "--published-checks", "--json")
		cmd.Dir = root
		out, err = cmd.CombinedOutput()
		if err == nil || !strings.Contains(string(out), "syntax_failure") {
			t.Fatalf("bad reference accepted: %s %v", out, err)
		}
	})
	t.Run("numerically-sufficient-drafts", func(t *testing.T) {
		root := t.TempDir()
		p := curriculum.Pack{SchemaVersion: 1, ID: "go-first-steps", Version: "1.0.0"}
		for i := 0; i < 160; i++ {
			p.Concepts = append(p.Concepts, curriculum.ConceptAuthoring{ID: fmt.Sprintf("concept-%d", i), Title: "Concept"})
		}
		for i := 0; i < 100; i++ {
			p.Competencies = append(p.Competencies, curriculum.CompetencyAuthoring{ID: fmt.Sprintf("competency-%d", i), Title: "Competency"})
		}
		for i := 0; i < 12; i++ {
			p.Tracks = append(p.Tracks, curriculum.TrackAuthoring{ID: fmt.Sprintf("track-%d", i), Title: "Track"})
		}
		for i := 0; i < 84; i++ {
			ch := curriculum.ChallengeAuthoring{SchemaVersion: 1, ID: fmt.Sprintf("challenge-%d", i), Version: "1.0.0", Title: "Draft", Kind: "atomic", Difficulty: "foundational", Acceptance: []string{"Observable result"}, Competencies: curriculum.CompetencyRefs{Primary: []string{"competency-0"}}, Canonical: true}
			steps := []curriculum.StepAuthoring{}
			for j := 0; j < 6; j++ {
				steps = append(steps, curriculum.StepAuthoring{ID: fmt.Sprintf("step-%d", j), Kind: "micro"})
			}
			ch.Layers = []curriculum.LayerAuthoring{{ID: "layer", MacroSteps: steps}}
			p.Challenges = append(p.Challenges, ch)
		}
		writePublicationPack(t, root, p)
		policy, err := os.ReadFile("../../packs/distribution.json")
		if err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(filepath.Join(root, "packs", "distribution.json"), policy, 0600); err != nil {
			t.Fatal(err)
		}
		cmd := exec.Command(bin, "catalog", "validate", "--json")
		cmd.Dir = root
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("incremental validation blocked drafts: %s %v", out, err)
		}
		var result struct {
			Publication curriculum.PublicationCoverage
		}
		if err = json.Unmarshal(out, &result); err != nil {
			t.Fatal(err)
		}
		if f := curriculum.CheckV1Gate(result.Publication.Inventory); len(f) != 0 {
			t.Fatalf("fixture not numerically sufficient: %+v", f)
		}
		cmd = exec.Command(bin, "catalog", "validate", "--v1-gate", "--json")
		cmd.Dir = root
		out, err = cmd.CombinedOutput()
		if err == nil {
			t.Fatal("numerically sufficient drafts passed V1")
		}
		if err = json.Unmarshal(out, &result); err != nil {
			t.Fatalf("non-JSON gate output: %s", out)
		}
		if result.Publication.Eligible.Challenges != 0 || !strings.Contains(string(out), "v1_gate_below_threshold") {
			t.Fatalf("did not gate on publication: %s", out)
		}
	})
}
