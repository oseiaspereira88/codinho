package checks

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/oseiaspereira88/codinho/internal/workspace"
)

func TestEditorialRealGoTestProof(t *testing.T) {
	for _, tc := range []struct{ name, source, pattern, timeout, want string }{
		{"pass", "func TestProof(t *testing.T) {}", "TestProof", "", "pass"},
		{"test-failure", "func TestProof(t *testing.T) {t.Fatal(\"expected\")}", "TestProof", "", "test_failure"},
		{"compile-failure", "func TestProof(t *testing.T) {notDefined()}", "TestProof", "", "compile_failure"},
		{"skip", "func TestProof(t *testing.T) {t.Skip(\"unavailable\")}", "TestProof", "", "skipped"},
		{"mixed-skip-failure", "func TestProof(t *testing.T) {t.Fatal(\"expected\")}; func TestSkip(t *testing.T) {t.Skip()}", "Test", "", "incomplete_failure"},
		{"no-match", "func TestProof(t *testing.T) {}", "TestMissing", "", "no_tests"},
		{"timeout", "func TestProof(t *testing.T) {}", "TestProof", "1ns", "timeout"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			for name, data := range map[string]string{"go.mod": "module proof.test/fixture\n\ngo 1.25.0\n", "proof_test.go": "package proof\nimport \"testing\"\n" + tc.source + "\n"} {
				if err := os.WriteFile(filepath.Join(dir, name), []byte(data), 0600); err != nil {
					t.Fatal(err)
				}
			}
			root, err := workspace.AuthorizeRoot(dir)
			if err != nil {
				t.Fatal(err)
			}
			r, err := Resolve(CheckSpec{ID: "proof", Runner: "go_test", Package: "./...", TestPattern: tc.pattern, Timeout: tc.timeout})
			if err != nil {
				t.Fatal(err)
			}
			got, err := NewExecutor().ExecuteEditorial(context.Background(), r, root)
			if err != nil || got != tc.want {
				t.Fatalf("proof=%s err=%v want=%s", got, err, tc.want)
			}
		})
	}
}

func TestEditorialRejectsIncompleteTestOutput(t *testing.T) {
	for _, output := range []string{
		"ok package (cached)\n",
		"{\"Action\":\"start\",\"Package\":\"p\"}\n",
		"{\"Action\":\"pass\",\"Test\":\"TestProof\",\"Package\":\"p\"}\n",
	} {
		if got := classifyTestProof(Result{Outcome: OutcomePass, Stdout: []byte(output)}); got == "pass" {
			t.Fatalf("accepted incomplete output: %s", output)
		}
	}
}

func TestEditorialBoundsCapturedOutput(t *testing.T) {
	var out boundedOutput
	input := bytes.NewReader(bytes.Repeat([]byte("x"), maxOutputBytes*4))
	n, err := io.Copy(&out, input)
	if err != nil || n != maxOutputBytes*4 {
		t.Fatalf("did not drain output: %d %v", n, err)
	}
	if out.buffer.Len() != maxOutputBytes || len(out.Bytes()) != maxOutputBytes+len(truncationSuffix) {
		t.Fatal("output capture is unbounded or did not mark truncation")
	}
}

func TestEditorialOtherRunners(t *testing.T) {
	for _, tc := range []struct{ name, runner, file, source, pattern, want string }{
		{"benchmark", "go_benchmark", "proof_test.go", "package proof\nimport \"testing\"\nfunc BenchmarkProof(b *testing.B) { for i:=0;i<b.N;i++ {} }\n", "BenchmarkProof", "pass"},
		{"benchmark-missing", "go_benchmark", "proof_test.go", "package proof\nimport \"testing\"\nfunc BenchmarkProof(b *testing.B) {}\n", "BenchmarkMissing", "no_tests"},
		{"benchmark-skip", "go_benchmark", "proof_test.go", "package proof\nimport \"testing\"\nfunc BenchmarkProof(b *testing.B) { b.Skip() }\n", "BenchmarkProof", "skipped"},
		{"gofmt-invalid", "gofmt_check", "proof.go", "not valid syntax", "", "syntax_failure"},
		{"gofmt-unformatted", "gofmt_check", "proof.go", "package proof\nfunc  F( ){ }\n", "", "format_failure"},
		{"ast-missing", "internal_ast", "proof.go", "package proof\n", "", "error"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			for name, data := range map[string]string{"go.mod": "module proof.test/fixture\n\ngo 1.25.0\n", tc.file: tc.source} {
				if err := os.WriteFile(filepath.Join(dir, name), []byte(data), 0600); err != nil {
					t.Fatal(err)
				}
			}
			root, err := workspace.AuthorizeRoot(dir)
			if err != nil {
				t.Fatal(err)
			}
			pkg := tc.file
			if tc.runner == "go_benchmark" {
				pkg = "./..."
			}
			if tc.name == "ast-missing" {
				pkg = "missing.go"
			}
			r, err := Resolve(CheckSpec{ID: "proof", Runner: tc.runner, Package: pkg, TestPattern: tc.pattern})
			if err != nil {
				t.Fatal(err)
			}
			got, err := NewExecutor().ExecuteEditorial(context.Background(), r, root)
			if err != nil || got != tc.want {
				t.Fatalf("proof=%s err=%v want=%s", got, err, tc.want)
			}
		})
	}
}
