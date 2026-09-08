package checks

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/oseiaspereira88/codinho/internal/workspace"
)

// ExecuteEditorial returns a proof category, never process output. Test proof
// requires uncached JSON test events, so exit zero without tests is not a pass.
// This does not change learner-facing Execute's historical outcome contract.
func (e *Executor) ExecuteEditorial(ctx context.Context, r Resolved, root workspace.Root) (string, error) {
	if r.Kind == KindGoTest || r.Kind == KindGoTestRace || r.Kind == KindGoBenchmark {
		r.Args = append([]string{r.Args[0], "-json", "-count=1"}, r.Args[1:]...)
	}
	if r.Kind == KindGoBenchmark {
		r.Args = append([]string{r.Args[0], "-benchtime=1x"}, r.Args[1:]...)
	}
	if r.Kind == KindInternalAST {
		path, err := root.Resolve(r.path)
		if err != nil {
			return "error", nil
		}
		if info, err := os.Stat(path); err != nil || !info.Mode().IsRegular() {
			return "error", nil
		}
	}
	r.NetworkApproved = false
	result, err := e.Execute(ctx, r, root)
	if err != nil {
		return "error", err
	}
	if bytes.Contains(result.Stdout, truncationSuffix) || bytes.Contains(result.Stderr, truncationSuffix) {
		return "unverifiable_output", nil
	}
	if result.Outcome == OutcomeError {
		if bytes.Contains(result.Stderr, []byte("[TIMEOUT]")) {
			return "timeout", nil
		}
		return "error", nil
	}
	if r.Kind == KindGoTest || r.Kind == KindGoTestRace || r.Kind == KindGoBenchmark {
		return classifyTestProof(result, r.Kind == KindGoBenchmark), nil
	}
	if r.Kind == KindGofmtCheck && result.ExitCode != 0 {
		return "syntax_failure", nil
	}
	if result.Outcome == OutcomeSkipped {
		return "skipped", nil
	}
	if result.Outcome == OutcomePass {
		return "pass", nil
	}
	switch r.Kind {
	case KindGoBuild:
		return "compile_failure", nil
	case KindGoVet:
		return "analysis_failure", nil
	case KindGofmtCheck:
		return "format_failure", nil
	case KindInternalAST:
		return "syntax_failure", nil
	default:
		return "error", nil
	}
}

func classifyTestProof(r Result, benchmark ...bool) string {
	dec := json.NewDecoder(bytes.NewReader(r.Stdout))
	ran, failed, packages := 0, 0, 0
	skipped, buildFailed := false, false
	started := map[string]bool{}
	ended := map[string]bool{}
	packageFailed := map[string]bool{}
	testFailed := map[string]bool{}
	for {
		var ev struct{ Action, Package, Test, Output, FailedBuild string }
		err := dec.Decode(&ev)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "unverifiable_output"
		}
		if ev.FailedBuild != "" || ev.Action == "build-fail" {
			buildFailed = true
		}
		if len(benchmark) > 0 && benchmark[0] && ev.Action == "output" && strings.HasPrefix(ev.Test, "Benchmark") {
			fields := strings.Fields(ev.Output)
			if len(fields) >= 4 && strings.HasPrefix(fields[0], ev.Test) && fields[3] == "ns/op" {
				if n, err := strconv.Atoi(fields[1]); err == nil && n > 0 {
					ran++
				}
			}
		}
		if ev.Test != "" {
			switch ev.Action {
			case "pass":
				ran++
			case "fail":
				ran++
				failed++
				testFailed[ev.Package] = true
			case "skip":
				skipped = true
			}
		} else {
			switch ev.Action {
			case "start":
				started[ev.Package] = true
			case "pass", "fail", "skip":
				ended[ev.Package] = true
				packages++
				if ev.Action == "fail" {
					packageFailed[ev.Package] = true
				}
				if ev.Action == "skip" {
					skipped = true
				}
			}
		}
		if strings.Contains(ev.Output, "[build failed]") {
			buildFailed = true
		}
	}
	if skipped {
		if failed > 0 || buildFailed || len(packageFailed) > 0 {
			return "incomplete_failure"
		}
		return "skipped"
	}
	if buildFailed {
		return "compile_failure"
	}
	if packages == 0 {
		return "unverifiable_output"
	}
	for pkg := range started {
		if !ended[pkg] {
			return "unverifiable_output"
		}
	}
	for pkg := range packageFailed {
		if !testFailed[pkg] {
			return "unexpected_package_failure"
		}
	}
	if ran == 0 {
		return "no_tests"
	}
	if r.ExitCode != 0 {
		if failed > 0 {
			return "test_failure"
		}
		return "unexpected_package_failure"
	}
	if failed > 0 {
		return "unverifiable_output"
	}
	return "pass"
}
