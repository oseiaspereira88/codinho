package checks

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/oseiaspereira88/codinho/internal/workspace"
)

// allowedEnvNames is the minimum environment a Go toolchain invocation
// needs to function locally, and nothing else (ADR learner-code-ownership-
// and-safe-checks: "ambiente mínimo com allowlist de variáveis").
var allowedEnvNames = map[string]bool{
	"PATH":        true,
	"HOME":        true,
	"GOCACHE":     true,
	"GOPATH":      true,
	"GOMODCACHE":  true,
	"TMPDIR":      true,
	"TEMP":        true,
	"TMP":         true,
	"USERPROFILE": true, // windows equivalent of HOME, harmless elsewhere
}

// Executor runs Resolved checks with a bounded degree of parallelism
// (requirement R4).
type Executor struct {
	sem chan struct{}
}

// NewExecutor returns an Executor bounded to maxParallel concurrent runs.
func NewExecutor() *Executor {
	return &Executor{sem: make(chan struct{}, maxParallel)}
}

// Execute runs r inside root, never through a shell, contained to root's
// real directory, with a minimal allowlisted environment, network denied
// unless r.NetworkApproved, a hard timeout, and the whole process group
// killed on cancellation (requirements R3, R4, R5, R6, R8).
func (e *Executor) Execute(ctx context.Context, r Resolved, root workspace.Root) (Result, error) {
	e.sem <- struct{}{}
	defer func() { <-e.sem }()

	if r.Kind == KindInternalAST {
		return executeInternalAST(root, r.path)
	}

	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	cmd := exec.Command(r.Program, r.Args...)
	cmd.Dir = root.Path()
	cmd.Env = buildEnv(r.NetworkApproved)
	setProcessGroup(cmd)

	var stdout, stderr boundedOutput
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return Result{Outcome: OutcomeError, Stderr: []byte(safeStartError(err))}, nil
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	var runErr error
	select {
	case runErr = <-done:
	case <-ctx.Done():
		killProcessGroup(cmd)
		<-done // reap: avoid leaving a zombie behind
		runErr = ctx.Err()
	}

	out := workspace.Redact(capOutput(stdout.Bytes()))
	errOut := workspace.Redact(capOutput(stderr.Bytes()))

	if errors.Is(runErr, context.DeadlineExceeded) {
		return Result{Outcome: OutcomeError, Stdout: out, Stderr: append(errOut, []byte("\n[TIMEOUT]")...)}, nil
	}

	exitCode := 0
	if runErr != nil {
		var exitErr *exec.ExitError
		if errors.As(runErr, &exitErr) {
			exitCode = exitErr.ExitCode()
		} else {
			// The process never produced an exit code we can trust
			// (killed by a signal outside our own cancellation, I/O
			// failure, ...): infra failure, not a learner failure.
			return Result{Outcome: OutcomeError, Stdout: out, Stderr: errOut}, nil
		}
	}

	if r.Kind == KindGofmtCheck {
		// `gofmt -l` always exits 0; it reports unformatted files by
		// listing their paths on stdout, never through the exit code.
		if len(bytes.TrimSpace(out)) > 0 {
			return Result{Outcome: OutcomeFail, ExitCode: exitCode, Stdout: out, Stderr: errOut}, nil
		}
		return Result{Outcome: OutcomePass, ExitCode: exitCode, Stdout: out, Stderr: errOut}, nil
	}

	if exitCode == 0 {
		if looksSkipped(out) {
			return Result{Outcome: OutcomeSkipped, ExitCode: exitCode, Stdout: out, Stderr: errOut}, nil
		}
		return Result{Outcome: OutcomePass, ExitCode: exitCode, Stdout: out, Stderr: errOut}, nil
	}
	return Result{Outcome: OutcomeFail, ExitCode: exitCode, Stdout: out, Stderr: errOut}, nil
}

// looksSkipped recognizes go's own stable "nothing ran" messages, so a
// pattern matching zero tests is reported as skipped rather than a false
// pass (requirement R8).
func looksSkipped(stdout []byte) bool {
	s := string(stdout)
	return strings.Contains(s, "no test files") ||
		strings.Contains(s, "testing: warning: no tests to run") ||
		strings.Contains(s, "[no tests to run]")
}

// buildEnv returns the minimal allowlisted environment for a check
// invocation. GOPROXY/GOFLAGS deny network module resolution unless the
// challenge explicitly approved it (requirement R5).
func buildEnv(networkApproved bool) []string {
	env := make([]string, 0, len(allowedEnvNames)+2)
	for name := range allowedEnvNames {
		if v, ok := os.LookupEnv(name); ok {
			env = append(env, name+"="+v)
		}
	}
	if !networkApproved {
		// GOPROXY=off denies module-fetch network access; GOFLAGS=-mod=mod
		// would let `go` silently rewrite the learner's go.mod/go.sum to
		// satisfy a missing dependency, which is itself an edit this
		// package must never make (constraint: leitura apenas).
		// -mod=readonly fails fast instead.
		env = append(env, "GOPROXY=off", "GOFLAGS=-mod=readonly")
	}
	return env
}

// safeStartError reports that a check failed to even start, without
// leaking the resolved absolute path or environment details a raw error
// string might contain (security requirement).
func safeStartError(err error) string {
	if errors.Is(err, exec.ErrNotFound) {
		return "check runner executable not found"
	}
	return "check failed to start"
}

// executeInternalAST runs the in-process structural verifier: it parses
// path as Go source and reports pass/fail without spawning anything
// (requirement R2: "verificadores internos").
func executeInternalAST(root workspace.Root, path string) (Result, error) {
	full, err := root.Resolve(path)
	if err != nil {
		return Result{Outcome: OutcomeError, Stderr: []byte("path escapes the authorized workspace root")}, nil
	}
	if _, err := workspace.ParseGoFile(full); err != nil {
		return Result{Outcome: OutcomeFail, Stderr: workspace.Redact([]byte(err.Error()))}, nil
	}
	return Result{Outcome: OutcomePass}, nil
}

// boundedOutput drains the subprocess stream while retaining at most the cap.
// Both streams have separate writers; each is written by one exec copy goroutine.
type boundedOutput struct {
	buffer    bytes.Buffer
	truncated bool
}

func (b *boundedOutput) Write(data []byte) (int, error) {
	n := len(data)
	remaining := maxOutputBytes - b.buffer.Len()
	if len(data) > remaining {
		data = data[:remaining]
		b.truncated = true
	}
	_, err := b.buffer.Write(data)
	return n, err
}

func (b *boundedOutput) Bytes() []byte {
	if b.truncated {
		return append(b.buffer.Bytes(), truncationSuffix...)
	}
	return b.buffer.Bytes()
}
