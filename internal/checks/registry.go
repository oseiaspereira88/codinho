package checks

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// ErrUnknownKind is returned when a CheckSpec's Runner is not one of the
// allowlisted Kind values (ADR learner-code-ownership-and-safe-checks).
var ErrUnknownKind = errors.New("checks: unknown runner kind")

// ErrInvalidPackage is returned when Package escapes the workspace, is
// absolute, or could be mistaken for a flag by the invoked program.
var ErrInvalidPackage = errors.New("checks: invalid package path")

// ErrInvalidTestPattern is returned when TestPattern contains characters
// that have no legitimate use in a Go test-name regexp.
var ErrInvalidTestPattern = errors.New("checks: invalid test pattern")

// ErrInvalidTimeout is returned when Timeout does not parse as a positive
// duration within the safe ceiling.
var ErrInvalidTimeout = errors.New("checks: invalid timeout")

// packagePattern allows the shapes Go module-relative package paths and
// plain file paths actually take: letters, digits, '.', '_', '-', '/'. It
// rejects anything else outright, closing off shell/flag metacharacters
// before any exec call is built (requirement R3; security).
var packagePattern = regexp.MustCompile(`^[A-Za-z0-9._/-]+$`)

// testPatternChars allows a Go regexp's legitimate alphabet plus common
// anchors, rejecting shell metacharacters and control bytes.
var testPatternChars = regexp.MustCompile(`^[A-Za-z0-9_^$.*+?()|\[\]{}\\-]+$`)

var allowedKinds = map[Kind]bool{
	KindGoTest:      true,
	KindGoTestRace:  true,
	KindGoVet:       true,
	KindGofmtCheck:  true,
	KindGoBuild:     true,
	KindGoBenchmark: true,
	KindInternalAST: true,
}

// Resolve validates spec against the fixed allowlist and shape rules,
// returning a Resolved ready for Execute. It never trusts spec.Package or
// spec.TestPattern beyond what this validation confirms (requirement R2,
// R3, R9): nothing here derives a program or argument from anything the
// caller of check_run supplied directly — spec itself comes only from the
// challenge's own fixed catalog (requirement R1).
func Resolve(spec CheckSpec) (Resolved, error) {
	kind := Kind(spec.Runner)
	if !allowedKinds[kind] {
		return Resolved{}, fmt.Errorf("%w: %q", ErrUnknownKind, spec.Runner)
	}

	pkg := spec.Package
	if pkg == "" {
		pkg = "./..."
	}
	if err := validatePackage(pkg); err != nil {
		return Resolved{}, err
	}

	timeout, err := resolveTimeout(spec.Timeout)
	if err != nil {
		return Resolved{}, err
	}

	if kind == KindInternalAST {
		return Resolved{ID: spec.ID, Kind: kind, path: pkg, Timeout: timeout}, nil
	}

	pattern := spec.TestPattern
	if pattern != "" {
		if err := validateTestPattern(pattern); err != nil {
			return Resolved{}, err
		}
	}

	program, args := buildArgs(kind, pkg, pattern)
	return Resolved{
		ID:              spec.ID,
		Kind:            kind,
		Program:         program,
		Args:            args,
		Timeout:         timeout,
		NetworkApproved: spec.NetworkApproved,
	}, nil
}

func validatePackage(pkg string) error {
	if !packagePattern.MatchString(pkg) {
		return fmt.Errorf("%w: %q", ErrInvalidPackage, pkg)
	}
	if strings.HasPrefix(pkg, "-") {
		return fmt.Errorf("%w: must not look like a flag: %q", ErrInvalidPackage, pkg)
	}
	if strings.HasPrefix(pkg, "/") {
		return fmt.Errorf("%w: must be relative to the workspace root: %q", ErrInvalidPackage, pkg)
	}
	for seg := range strings.SplitSeq(pkg, "/") {
		if seg == ".." {
			return fmt.Errorf("%w: must not traverse above the workspace root: %q", ErrInvalidPackage, pkg)
		}
	}
	return nil
}

func validateTestPattern(pattern string) error {
	if len(pattern) > 200 {
		return fmt.Errorf("%w: too long", ErrInvalidTestPattern)
	}
	if strings.HasPrefix(pattern, "-") {
		return fmt.Errorf("%w: must not look like a flag: %q", ErrInvalidTestPattern, pattern)
	}
	if !testPatternChars.MatchString(pattern) {
		return fmt.Errorf("%w: %q", ErrInvalidTestPattern, pattern)
	}
	return nil
}

func resolveTimeout(declared string) (time.Duration, error) {
	if declared == "" {
		return defaultTimeout, nil
	}
	d, err := time.ParseDuration(declared)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidTimeout, err)
	}
	if d <= 0 || d > maxTimeout {
		return 0, fmt.Errorf("%w: %s exceeds the safe ceiling of %s", ErrInvalidTimeout, d, maxTimeout)
	}
	return d, nil
}

// buildArgs returns the fixed program and argv for kind, never
// interpolating pkg/pattern into a shell string — each becomes its own
// exec.Cmd argument.
func buildArgs(kind Kind, pkg, pattern string) (program string, args []string) {
	switch kind {
	case KindGoTest:
		args = []string{"test"}
		if pattern != "" {
			args = append(args, "-run", pattern)
		}
		return "go", append(args, pkg)
	case KindGoTestRace:
		args = []string{"test", "-race"}
		if pattern != "" {
			args = append(args, "-run", pattern)
		}
		return "go", append(args, pkg)
	case KindGoVet:
		return "go", []string{"vet", pkg}
	case KindGoBuild:
		return "go", []string{"build", pkg}
	case KindGoBenchmark:
		bench := pattern
		if bench == "" {
			bench = "."
		}
		return "go", []string{"test", "-bench", bench, "-run", "^$", pkg}
	case KindGofmtCheck:
		return "gofmt", []string{"-l", pkg}
	default:
		return "", nil
	}
}
