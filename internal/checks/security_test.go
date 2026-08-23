package checks

import (
	"context"
	"os/exec"
	"strings"
	"testing"
)

func TestExecuteAllowsNetworkOnlyWhenExplicitlyApproved(t *testing.T) {
	_, root := newFixtureModule(t)
	writeFixtureFile(t, root.Path(), "main_test.go", `package main

import (
	"os"
	"testing"
)

func TestReportsGoproxy(t *testing.T) {
	if v := os.Getenv("GOPROXY"); v == "off" {
		t.Fatalf("expected GOPROXY to be unset when network is approved, got %q", v)
	}
}
`)
	e := NewExecutor()
	resolved, err := Resolve(CheckSpec{ID: "net", Runner: string(KindGoTest), TestPattern: "^TestReportsGoproxy$", NetworkApproved: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := e.Execute(context.Background(), resolved, root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Outcome != OutcomePass {
		t.Fatalf("expected pass, got %s stdout=%s", result.Outcome, result.Stdout)
	}
}

func TestExecuteRedactsSecretsInCapturedOutput(t *testing.T) {
	_, root := newFixtureModule(t)
	writeFixtureFile(t, root.Path(), "main_test.go", `package main

import "testing"

func TestPrintsSecret(t *testing.T) {
	t.Error("api_key: AKIAABCDEFGHIJKLMNOPQRST")
}
`)
	e := NewExecutor()
	resolved, err := Resolve(CheckSpec{ID: "secret", Runner: string(KindGoTest), TestPattern: "^TestPrintsSecret$"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := e.Execute(context.Background(), resolved, root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(string(result.Stdout), "AKIAABCDEFGHIJKLMNOPQRST") {
		t.Fatalf("expected the secret to be redacted from captured stdout, got: %s", result.Stdout)
	}
}

func TestExecuteReportsErrorRatherThanFailWhenTheRunnerBinaryIsMissing(t *testing.T) {
	_, root := newFixtureModule(t)
	e := NewExecutor()
	// A hand-built Resolved bypassing Resolve's allowlist, simulating a
	// misconfigured environment where the resolved program is unusable —
	// R8 requires this to surface as infra error, never as the learner's
	// check failing.
	resolved := Resolved{ID: "missing", Kind: KindGoTest, Program: "codinho-definitely-not-a-real-binary", Args: []string{"test"}, Timeout: defaultTimeout}
	result, err := e.Execute(context.Background(), resolved, root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Outcome != OutcomeError {
		t.Fatalf("expected error outcome for a missing runner binary, got %s", result.Outcome)
	}
}

func TestExecuteInternalASTRejectsPathEscapingRoot(t *testing.T) {
	resolved, err := Resolve(CheckSpec{ID: "ast", Runner: string(KindInternalAST), Package: "../../../etc/passwd"})
	if err == nil {
		t.Fatalf("expected Resolve to reject a traversal path outright, got %+v", resolved)
	}
}

// TestExecuteNeverMutatesGitState proves security-privacy-hardening
// requirement R9: running a real check against a dirty Git working tree
// never changes its status or HEAD — Execute only ever reads the
// workspace and writes nothing back to it.
func TestExecuteNeverMutatesGitState(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
	dir, root := newFixtureModule(t)
	runGit(t, dir, "init", "-q")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "-c", "user.email=test@test", "-c", "user.name=test", "commit", "-q", "-m", "initial")
	writeFixtureFile(t, dir, "untracked.txt", "scratch")

	statusBefore := gitOutput(t, dir, "status", "--porcelain")
	headBefore := gitOutput(t, dir, "rev-parse", "HEAD")

	resolved, err := Resolve(CheckSpec{ID: "vet", Runner: string(KindGoVet), Package: "./..."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := NewExecutor().Execute(context.Background(), resolved, root); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := gitOutput(t, dir, "status", "--porcelain"); got != statusBefore {
		t.Fatalf("git status changed:\nbefore: %q\nafter:  %q", statusBefore, got)
	}
	if got := gitOutput(t, dir, "rev-parse", "HEAD"); got != headBefore {
		t.Fatalf("HEAD moved: before=%q after=%q", headBefore, got)
	}
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}

func gitOutput(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}
