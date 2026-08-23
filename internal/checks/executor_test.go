package checks

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/oseiaspereira88/codinho/internal/workspace"
)

func writeFixtureFile(t *testing.T, dir, rel, content string) {
	t.Helper()
	full := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o700); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := os.WriteFile(full, []byte(content), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// newFixtureModule creates a minimal, valid Go module so go_test/go_vet/
// go_build/gofmt_check have something real to run against, without
// needing network access (no external dependencies).
func newFixtureModule(t *testing.T) (dir string, root workspace.Root) {
	t.Helper()
	dir = t.TempDir()
	writeFixtureFile(t, dir, "go.mod", "module fixture\n\ngo 1.25\n")
	writeFixtureFile(t, dir, "main.go", "package main\n\nfunc main() {}\n")
	root, err := workspace.AuthorizeRoot(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return dir, root
}

func TestResolveRejectsUnknownKind(t *testing.T) {
	if _, err := Resolve(CheckSpec{ID: "x", Runner: "run_command"}); err == nil {
		t.Fatal("expected an error for an unknown runner")
	}
}

func TestResolveRejectsPackageTraversal(t *testing.T) {
	cases := []string{"../../etc/passwd", "/etc/passwd", "-o=evil"}
	for _, pkg := range cases {
		if _, err := Resolve(CheckSpec{ID: "x", Runner: string(KindGoTest), Package: pkg}); err == nil {
			t.Errorf("expected Resolve to reject package %q", pkg)
		}
	}
}

func TestResolveRejectsInvalidTestPattern(t *testing.T) {
	cases := []string{"-run=all; rm -rf /", "test`whoami`", strings.Repeat("a", 300)}
	for _, pattern := range cases {
		if _, err := Resolve(CheckSpec{ID: "x", Runner: string(KindGoTest), TestPattern: pattern}); err == nil {
			t.Errorf("expected Resolve to reject test pattern %q", pattern)
		}
	}
}

func TestResolveRejectsTimeoutOutsideSafeCeiling(t *testing.T) {
	cases := []string{"not-a-duration", "-5s", "0s", "1h"}
	for _, timeout := range cases {
		if _, err := Resolve(CheckSpec{ID: "x", Runner: string(KindGoTest), Timeout: timeout}); err == nil {
			t.Errorf("expected Resolve to reject timeout %q", timeout)
		}
	}
}

func TestResolveDefaultsPackageAndTimeout(t *testing.T) {
	r, err := Resolve(CheckSpec{ID: "x", Runner: string(KindGoVet)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Timeout != defaultTimeout {
		t.Fatalf("expected default timeout, got %v", r.Timeout)
	}
	if got := r.Args; len(got) != 2 || got[1] != "./..." {
		t.Fatalf("expected default package ./..., got %+v", got)
	}
}

func TestExecuteGoVetPassAndFail(t *testing.T) {
	_, root := newFixtureModule(t)
	e := NewExecutor()

	resolved, err := Resolve(CheckSpec{ID: "vet", Runner: string(KindGoVet)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := e.Execute(context.Background(), resolved, root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Outcome != OutcomePass {
		t.Fatalf("expected pass, got %s (stderr=%s)", result.Outcome, result.Stderr)
	}
}

func TestExecuteGoBuildFailsOnBrokenCode(t *testing.T) {
	dir, root := newFixtureModule(t)
	writeFixtureFile(t, dir, "main.go", "package main\n\nfunc main() { this is not go }\n")
	e := NewExecutor()

	resolved, err := Resolve(CheckSpec{ID: "build", Runner: string(KindGoBuild)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := e.Execute(context.Background(), resolved, root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Outcome != OutcomeFail {
		t.Fatalf("expected fail, got %s", result.Outcome)
	}
}

func TestExecuteGofmtCheckDetectsUnformattedCode(t *testing.T) {
	dir, root := newFixtureModule(t)
	writeFixtureFile(t, dir, "main.go", "package main\nfunc main(){\nx:=1\n_=x\n}\n")
	e := NewExecutor()

	resolved, err := Resolve(CheckSpec{ID: "fmt", Runner: string(KindGofmtCheck), Package: "main.go"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := e.Execute(context.Background(), resolved, root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Outcome != OutcomeFail {
		t.Fatalf("expected fail (unformatted), got %s stdout=%s", result.Outcome, result.Stdout)
	}
	if !strings.Contains(string(result.Stdout), "main.go") {
		t.Fatalf("expected gofmt to list main.go, got %s", result.Stdout)
	}
}

func TestExecuteGoTestRunsAPatternAndReportsSkippedWhenNothingMatches(t *testing.T) {
	dir, root := newFixtureModule(t)
	writeFixtureFile(t, dir, "main_test.go", `package main

import "testing"

func TestReal(t *testing.T) {}
`)
	e := NewExecutor()

	resolved, err := Resolve(CheckSpec{ID: "test", Runner: string(KindGoTest), TestPattern: "^TestDoesNotExist$"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := e.Execute(context.Background(), resolved, root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = dir
	if result.Outcome != OutcomeSkipped {
		t.Fatalf("expected skipped, got %s stdout=%s", result.Outcome, result.Stdout)
	}
}

func TestExecuteTimeoutKillsTheWholeProcessGroup(t *testing.T) {
	if !ProcessGroupSupported() {
		t.Skip("process group cancellation is not supported on this platform")
	}
	dir, root := newFixtureModule(t)
	pidFile := filepath.Join(dir, "child.pid")
	writeFixtureFile(t, dir, "main_test.go", `package main

import (
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"
)

func TestHangsForever(t *testing.T) {
	cmd := exec.Command("sleep", "30")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start sleep: %v", err)
	}
	os.WriteFile("`+strings.ReplaceAll(pidFile, `\`, `\\`)+`", []byte(strconv.Itoa(cmd.Process.Pid)), 0o600)
	time.Sleep(30 * time.Second)
}
`)
	e := NewExecutor()
	resolved, err := Resolve(CheckSpec{ID: "hang", Runner: string(KindGoTest), TestPattern: "^TestHangsForever$", Timeout: "1s"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	start := time.Now()
	result, err := e.Execute(context.Background(), resolved, root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if elapsed := time.Since(start); elapsed > 10*time.Second {
		t.Fatalf("expected the timeout to bound execution near 1s, took %s", elapsed)
	}
	if result.Outcome != OutcomeError {
		t.Fatalf("expected error (timeout), got %s", result.Outcome)
	}

	pidBytes, err := os.ReadFile(pidFile)
	if err != nil {
		t.Fatalf("expected the grandchild to have written its pid file: %v", err)
	}
	waitForProcessDeath(t, string(pidBytes))
}

func TestExecuteCapsOutputSize(t *testing.T) {
	dir, root := newFixtureModule(t)
	writeFixtureFile(t, dir, "main_test.go", `package main

import (
	"strings"
	"testing"
)

func TestFloodsOutput(t *testing.T) {
	t.Error(strings.Repeat("A", 2_000_000))
}
`)
	_ = dir
	e := NewExecutor()
	resolved, err := Resolve(CheckSpec{ID: "flood", Runner: string(KindGoTest), TestPattern: "^TestFloodsOutput$"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := e.Execute(context.Background(), resolved, root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Stdout) > maxOutputBytes+len(truncationSuffix)+64 {
		t.Fatalf("expected stdout to be capped near %d bytes, got %d", maxOutputBytes, len(result.Stdout))
	}
	if !strings.Contains(string(result.Stdout), "TRUNCATED") {
		t.Fatal("expected a truncation marker in the capped output")
	}
}

func TestExecuteDeniesNetworkByDefault(t *testing.T) {
	dir, root := newFixtureModule(t)
	writeFixtureFile(t, dir, "main_test.go", `package main

import (
	"os"
	"testing"
)

func TestReportsGoproxy(t *testing.T) {
	if os.Getenv("GOPROXY") != "off" {
		t.Fatalf("expected GOPROXY=off, got %q", os.Getenv("GOPROXY"))
	}
}
`)
	_ = dir
	e := NewExecutor()
	resolved, err := Resolve(CheckSpec{ID: "net", Runner: string(KindGoTest), TestPattern: "^TestReportsGoproxy$"})
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

func TestExecuteNeverLeaksUnallowlistedEnvironmentVariables(t *testing.T) {
	t.Setenv("CODINHO_TEST_SECRET", "should-not-leak")
	dir, root := newFixtureModule(t)
	writeFixtureFile(t, dir, "main_test.go", `package main

import (
	"os"
	"testing"
)

func TestEnvIsAllowlisted(t *testing.T) {
	if v := os.Getenv("CODINHO_TEST_SECRET"); v != "" {
		t.Fatalf("expected an empty environment for non-allowlisted vars, got %q", v)
	}
}
`)
	_ = dir
	e := NewExecutor()
	resolved, err := Resolve(CheckSpec{ID: "env", Runner: string(KindGoTest), TestPattern: "^TestEnvIsAllowlisted$"})
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

func TestExecuteInternalASTChecksStructuralValidity(t *testing.T) {
	dir, root := newFixtureModule(t)
	e := NewExecutor()

	resolved, err := Resolve(CheckSpec{ID: "ast", Runner: string(KindInternalAST), Package: "main.go"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := e.Execute(context.Background(), resolved, root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Outcome != OutcomePass {
		t.Fatalf("expected pass for valid Go, got %s", result.Outcome)
	}

	writeFixtureFile(t, dir, "broken.go", "package main\nthis is not valid go\n")
	resolvedBroken, err := Resolve(CheckSpec{ID: "ast", Runner: string(KindInternalAST), Package: "broken.go"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultBroken, err := e.Execute(context.Background(), resolvedBroken, root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resultBroken.Outcome != OutcomeFail {
		t.Fatalf("expected fail for invalid Go, got %s", resultBroken.Outcome)
	}
}
