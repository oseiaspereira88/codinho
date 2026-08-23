package checks

import (
	"context"
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
