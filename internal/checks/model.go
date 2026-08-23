// Package checks executes pre-declared, allowlisted verifications by
// check_id, isolating arguments, resources, environment and filesystem
// (ADR learner-code-ownership-and-safe-checks). It never accepts a free
// command string and never runs anything through a shell.
package checks

import "time"

// Kind is one of the allowlisted check runners (ADR learner-code-
// ownership-and-safe-checks; PROJECT.md §20.2). Any value outside this set
// is rejected by Resolve, never executed.
type Kind string

const (
	KindGoTest      Kind = "go_test"
	KindGoTestRace  Kind = "go_test_race"
	KindGoVet       Kind = "go_vet"
	KindGofmtCheck  Kind = "gofmt_check"
	KindGoBuild     Kind = "go_build"
	KindGoBenchmark Kind = "go_benchmark"
	// KindInternalAST is an in-process structural verifier: it parses a Go
	// file and reports pass/fail without spawning any subprocess.
	KindInternalAST Kind = "internal_ast"
)

// CheckSpec is what a caller (internal/application) supplies, translated
// verbatim from a challenge's authored internal/curriculum.CheckAuthoring.
// It carries no session or workspace path: Resolve validates it in
// isolation, and Execute binds it to a specific, already-authorized
// workspace.Root.
type CheckSpec struct {
	ID              string
	Runner          string
	Package         string
	TestPattern     string
	Timeout         string
	NetworkApproved bool
}

// Outcome is one of the stable results a check can report (requirement
// R8): it never confuses infrastructure failure (Error) with the
// learner's code failing the check (Fail).
type Outcome string

const (
	OutcomePass    Outcome = "pass"
	OutcomeFail    Outcome = "fail"
	OutcomeError   Outcome = "error"
	OutcomeSkipped Outcome = "skipped"
)

// Resolved is a CheckSpec that passed allowlist and shape validation,
// ready for Execute. Its zero value must never be executed; only Resolve
// constructs one.
type Resolved struct {
	ID              string
	Kind            Kind
	Program         string
	Args            []string
	Timeout         time.Duration
	NetworkApproved bool
	// path is set only for KindInternalAST: the file to parse, relative to
	// the workspace root.
	path string
}

// Result is what Execute reports for one run.
type Result struct {
	Outcome  Outcome
	ExitCode int
	Stdout   []byte
	Stderr   []byte
}
