package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// writeCLIFixturePack seeds dir with a pack that declares a fixture, so
// workspace prepare has something real to materialize.
func writeCLIFixturePack(t *testing.T, dir string) {
	t.Helper()
	packsDir := filepath.Join(dir, "packs")
	if err := os.MkdirAll(packsDir, 0o700); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(packsDir, "manifest.yaml"), []byte("schema_version: 1\npacks:\n  - pack.yaml\n"), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pack := `schema_version: 1
id: cli-fixture-pack
version: 1.0.0
competencies:
  - id: comp-a
    title: Competency A
challenges:
  - schema_version: 1
    id: cli.fixture-challenge
    version: 1.0.0
    title: CLI fixture challenge
    kind: debug
    difficulty: foundational
    competencies:
      primary: [comp-a]
    acceptance:
      - O bug relatado deixa de ocorrer.
    layers:
      - id: understanding
        macro_steps:
          - id: cli.step-one
            kind: micro
            title: Fixture step
            instruction:
              objective: Fix the bug.
              scope: Only the reported bug.
    fixture:
      - path: main.go
        content: "package main\n"
`
	if err := os.WriteFile(filepath.Join(packsDir, "pack.yaml"), []byte(pack), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func runCodinho(t *testing.T, bin, dir string, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	var out, errOut strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err == nil {
		return out.String(), errOut.String(), 0
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return out.String(), errOut.String(), exitErr.ExitCode()
	}
	t.Fatalf("running codinho %v: %v", args, err)
	return "", "", -1
}

// TestCLIEntrypointReachesEveryCommand spawns the real compiled binary
// (the same composed entrypoint cmd/codinho/main.go declares as
// codinho-cli's reachability target) and confirms every documented
// top-level command actually dispatches to real command logic rather than
// falling through to "unknown command" (administrative-cli-fixtures R1;
// evidenceClass: reachability).
func TestCLIEntrypointReachesEveryCommand(t *testing.T) {
	bin := buildCodinhoBinary(t)
	dir := t.TempDir()

	for _, args := range [][]string{
		{"version"},
		{"doctor"},
		{"help"},
		{"catalog"},
		{"session"},
		{"progress"},
		{"workspace"},
		{"privacy"},
	} {
		_, stderr, _ := runCodinho(t, bin, dir, args...)
		if strings.Contains(stderr, "unknown command") {
			t.Fatalf("command %v was not reachable: %s", args, stderr)
		}
	}
}

// TestCLIEndToEndWorkflow drives the real binary through init, catalog
// inspection and workspace prepare against real files on disk, so the
// composed CLI is proven end to end rather than only through in-process
// internal/cli tests (administrative-cli-fixtures requirements R1-R7;
// evidenceClass: e2e).
func TestCLIEndToEndWorkflow(t *testing.T) {
	bin := buildCodinhoBinary(t)
	workspace := t.TempDir()

	stdout, stderr, code := runCodinho(t, bin, workspace, "init")
	if code != 0 {
		t.Fatalf("init: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}

	writeCLIFixturePack(t, workspace)

	stdout, stderr, code = runCodinho(t, bin, workspace, "catalog", "validate")
	if code != 0 {
		t.Fatalf("catalog validate: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}

	stdout, stderr, code = runCodinho(t, bin, workspace, "catalog", "list", "--json")
	if code != 0 || !strings.Contains(stdout, "cli.fixture-challenge") {
		t.Fatalf("catalog list: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}

	stdout, stderr, code = runCodinho(t, bin, workspace, "catalog", "show", "cli.fixture-challenge")
	if code != 0 || !strings.Contains(stdout, "fixture:") {
		t.Fatalf("catalog show: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}

	dest := filepath.Join(t.TempDir(), "prepared")
	stdout, stderr, code = runCodinho(t, bin, workspace, "workspace", "prepare", "cli.fixture-challenge", "--dest", dest)
	if code != 0 {
		t.Fatalf("workspace prepare: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if _, err := os.Stat(filepath.Join(dest, "main.go")); err != nil {
		t.Fatalf("workspace prepare did not materialize main.go: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "manifest.json")); err != nil {
		t.Fatalf("workspace prepare did not write manifest.json: %v", err)
	}

	// Overwrite protection holds through the real binary too.
	_, stderr, code = runCodinho(t, bin, workspace, "workspace", "prepare", "cli.fixture-challenge", "--dest", dest)
	if code == 0 {
		t.Fatalf("expected non-zero exit on conflicting prepare, stderr=%q", stderr)
	}

	_, stderr, code = runCodinho(t, bin, workspace, "session", "inspect", "ses_never_started")
	if code == 0 {
		t.Fatalf("expected non-zero exit for a session that was never started, stderr=%q", stderr)
	}

	stdout, stderr, code = runCodinho(t, bin, workspace, "progress", "show")
	if code != 0 {
		t.Fatalf("progress show: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
}

// TestWorkspacePrepareMaterializesRealAuthoredDebugChallenge runs the real
// binary against the project's own packs/ (not a synthetic fixture),
// proving workspace prepare end to end against a real, authored
// kind:debug challenge's fixture content (learning-practice-debug-modes
// requirement R10).
func TestWorkspacePrepareMaterializesRealAuthoredDebugChallenge(t *testing.T) {
	bin := buildCodinhoBinary(t)
	repoRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("resolving repo root: %v", err)
	}

	stdout, stderr, code := runCodinho(t, bin, repoRoot, "catalog", "show", "go-debug.slice-off-by-one")
	if code != 0 || !strings.Contains(stdout, "fixture:") {
		t.Fatalf("catalog show: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}

	dest := t.TempDir()
	stdout, stderr, code = runCodinho(t, bin, repoRoot, "workspace", "prepare", "go-debug.slice-off-by-one", "--dest", dest)
	if code != 0 {
		t.Fatalf("workspace prepare: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}

	got, err := os.ReadFile(filepath.Join(dest, "main.go"))
	if err != nil {
		t.Fatalf("reading materialized main.go: %v", err)
	}
	if !strings.Contains(string(got), "func lastElement") || !strings.Contains(string(got), "values[len(values)]") {
		t.Fatalf("materialized content does not match the authored fixture:\n%s", got)
	}
}
