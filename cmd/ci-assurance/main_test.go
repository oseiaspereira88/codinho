package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCandidateModifiedSeparatesSourceAndGeneratedEvidence(t *testing.T) {
	root := t.TempDir()
	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %s (%v)", args, out, err)
		}
	}
	write := func(path, body string) {
		t.Helper()
		p := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0600); err != nil {
			t.Fatal(err)
		}
	}
	assert := func(want bool) {
		t.Helper()
		got, err := candidateModified(root)
		if err != nil || got != want {
			t.Fatalf("modified=%v want=%v err=%v", got, want, err)
		}
	}
	git("init", "-q")
	write("cmd/app/main.go", "package main\n")
	write(".pose/results/delivery-validation.json", "{}\n")
	git("add", ".")
	git("-c", "user.name=fixture", "-c", "user.email=fixture@example.invalid", "-c", "commit.gpgsign=false", "commit", "-qm", "fixture")
	assert(false)
	write(".pose/results/delivery-validation.json", "{\"outcome\":\"pass\"}\n")
	write(".pose/assessments/new.md", "generated\n")
	assert(false)
	write("cmd/app/main.go", "package main\n// changed\n")
	assert(true)
	write("cmd/app/main.go", "package main\n")
	assert(false)
	write("internal/new/new.go", "package new\n")
	assert(true)
	if _, err := candidateModified(t.TempDir()); err == nil {
		t.Fatal("non-repository must fail inspection")
	}
}
