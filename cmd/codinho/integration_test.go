package main

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// buildCodinhoBinary compiles the current package once per test run and
// returns the path to the binary, so every test in this file avoids
// paying `go build`'s cost more than once.
func buildCodinhoBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "codinho")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("building codinho: %v\n%s", err, out)
	}
	return bin
}

func writeIntegrationFixturePack(t *testing.T, dir string) {
	t.Helper()
	packsDir := filepath.Join(dir, "packs")
	if err := os.MkdirAll(packsDir, 0o700); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(packsDir, "manifest.yaml"), []byte("schema_version: 1\npacks:\n  - pack.yaml\n"), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pack := `schema_version: 1
id: fixture-pack
version: 1.0.0
competencies:
  - id: comp-a
    title: Competency A
challenges:
  - schema_version: 1
    id: fixture.challenge-one
    version: 1.0.0
    title: Fixture challenge
    kind: atomic
    difficulty: foundational
    competencies:
      primary: [comp-a]
    layers:
      - id: understanding
        macro_steps:
          - id: fixture.step-one
            kind: micro
            title: Fixture step
            instruction:
              objective: Declare the fixture type.
              scope: Only the declaration.
    checks:
      - id: focused-tests
        runner: go_test
        package: ./...
`
	if err := os.WriteFile(filepath.Join(packsDir, "pack.yaml"), []byte(pack), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func callTool(ctx context.Context, t *testing.T, cs *mcp.ClientSession, name string, args map[string]any) map[string]any {
	t.Helper()
	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		t.Fatalf("%s: transport error: %v", name, err)
	}
	if res.IsError {
		raw, _ := json.Marshal(res.Content)
		t.Fatalf("%s returned a protocol-level error: %s", name, raw)
	}
	m, ok := res.StructuredContent.(map[string]any)
	if !ok {
		t.Fatalf("%s: unexpected structured content shape: %#v", name, res.StructuredContent)
	}
	if m["status"] != "ok" {
		t.Fatalf("%s: unexpected envelope status: %+v", name, m)
	}
	return m
}

func envData(t *testing.T, env map[string]any) map[string]any {
	t.Helper()
	data, ok := env["data"].(map[string]any)
	if !ok {
		t.Fatalf("envelope has no data object: %+v", env)
	}
	return data
}

func revisionOf(t *testing.T, env map[string]any) float64 {
	t.Helper()
	rev, ok := envData(t, env)["revision"].(float64)
	if !ok {
		t.Fatalf("envelope data has no numeric revision: %+v", env)
	}
	return rev
}

// TestTutorSkillFullSessionRoutingOverRealStdio exercises the exact tool
// sequence the codinho skill documents in
// .agents/skills/codinho/references/mcp-tool-routing.md — discovery,
// session start, observe, check, evaluate, reflect, complete, advance,
// and mastery/review — against the real compiled binary over the same
// stdio transport a real host (Codex CLI, IDE extension) uses
// (tutor-skill-host-integration requirement R9, Decision 2).
func TestTutorSkillFullSessionRoutingOverRealStdio(t *testing.T) {
	bin := buildCodinhoBinary(t)
	workspaceRoot := t.TempDir()
	writeIntegrationFixturePack(t, workspaceRoot)

	learnerRoot := t.TempDir()
	if err := os.WriteFile(filepath.Join(learnerRoot, "go.mod"), []byte("module fixture\n\ngo 1.25\n"), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(learnerRoot, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(learnerRoot, "main_test.go"), []byte("package main\n\nimport \"testing\"\n\nfunc TestOK(t *testing.T) {}\n"), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()
	cmd := exec.Command(bin, "serve")
	cmd.Dir = workspaceRoot
	transport := &mcp.CommandTransport{Command: cmd}
	client := mcp.NewClient(&mcp.Implementation{Name: "tutor-skill-integration-test", Version: "0.0.0"}, nil)
	cs, err := client.Connect(ctx, transport, nil)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer cs.Close()

	// 1. Discovery (catalog_search, learning_path_recommend).
	callTool(ctx, t, cs, "catalog_search", map[string]any{"competency": "comp-a"})

	// 2. Start and announce.
	start := callTool(ctx, t, cs, "session_start", map[string]any{"challenge_id": "fixture.challenge-one"})
	sid, _ := start["session_id"].(string)
	stepID, _ := start["active_node"].(map[string]any)["id"].(string)
	rev := revisionOf(t, start)

	// 3. Observe before evaluate.
	baseline := callTool(ctx, t, cs, "workspace_observe", map[string]any{
		"session_id": sid, "step_id": stepID, "root": learnerRoot, "globs": []string{"**/*.go"}, "expected_revision": rev,
	})
	rev = revisionOf(t, baseline)

	// 4. Run the declared check.
	checkRun := callTool(ctx, t, cs, "check_run", map[string]any{
		"session_id": sid, "check_id": "focused-tests", "expected_revision": rev,
	})
	rev = revisionOf(t, checkRun)
	checkData := envData(t, checkRun)
	if checkData["outcome"] != "pass" {
		t.Fatalf("expected the check to pass against valid fixture code, got %+v", checkData)
	}
	evidenceID, _ := checkData["evidence_id"].(string)

	// 5. Evaluate with explicit submission_intent, citing the check's evidence.
	eval := callTool(ctx, t, cs, "step_evaluate", map[string]any{
		"session_id": sid, "expected_revision": rev, "submission_intent": true,
		"criteria": []map[string]any{
			{"name": "tests-pass", "check_id": "focused-tests", "kind": "structural", "severity": "blocking", "evidence_id": evidenceID},
		},
	})
	rev = revisionOf(t, eval)
	evalData := envData(t, eval)
	criteria, _ := evalData["criteria"].([]any)
	if len(criteria) != 1 || criteria[0].(map[string]any)["verdict"] != "met" {
		t.Fatalf("expected verdict=met from the real check outcome (requirement R10 of safe-check-executor), got %+v", evalData)
	}

	// 6. Reflect, complete, advance — separate calls (invariant 5).
	reflection := callTool(ctx, t, cs, "reflection_record", map[string]any{
		"session_id": sid, "competency_id": "comp-a", "prompt": "What did you learn?", "answer": "Zero values.", "expected_revision": rev,
	})
	rev = revisionOf(t, reflection)

	complete := callTool(ctx, t, cs, "step_complete", map[string]any{"session_id": sid, "expected_revision": rev})
	rev = revisionOf(t, complete)

	advance := callTool(ctx, t, cs, "step_advance", map[string]any{"session_id": sid, "expected_revision": rev})
	_ = advance

	// 7. Domain evidence and progress/review — global tools, no session_id.
	progressBefore := callTool(ctx, t, cs, "progress_get", map[string]any{})
	progressRev := revisionOf(t, progressBefore)

	record := callTool(ctx, t, cs, "mastery_evidence_record", map[string]any{
		"competency_id": "comp-a", "dimension": "autonomous_implementation", "evidence_id": evidenceID,
		"success": true, "expected_revision": progressRev,
	})
	recordData := envData(t, record)
	if recordData["state"] != "demonstrates_without_help" {
		t.Fatalf("expected demonstrates_without_help for a first autonomous success, got %+v", recordData)
	}

	progressAfter := callTool(ctx, t, cs, "progress_get", map[string]any{"competency_id": "comp-a"})
	competencies, _ := envData(t, progressAfter)["competencies"].(map[string]any)
	if _, ok := competencies["comp-a"]; !ok {
		t.Fatalf("expected comp-a to appear in progress after recording evidence, got %+v", competencies)
	}

	due := callTool(ctx, t, cs, "review_due", map[string]any{})
	dueList, _ := envData(t, due)["due"].([]any)
	if len(dueList) != 0 {
		t.Fatalf("expected nothing due immediately after scheduling, got %+v", dueList)
	}
}
