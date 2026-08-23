package workspace

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, dir, rel, content string) {
	t.Helper()
	full := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o700); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := os.WriteFile(full, []byte(content), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAuthorizeRootRejectsNonDirectory(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "file.txt", "x")
	if _, err := AuthorizeRoot(filepath.Join(dir, "file.txt")); err != ErrNotADirectory {
		t.Fatalf("expected ErrNotADirectory, got %v", err)
	}
}

func TestRootResolveRejectsTraversal(t *testing.T) {
	dir := t.TempDir()
	root, err := AuthorizeRoot(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := root.Resolve("../etc/passwd"); err != ErrPathEscapesRoot {
		t.Fatalf("expected ErrPathEscapesRoot, got %v", err)
	}
}

func TestRootResolveAcceptsContainedPath(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "main.go", "package main\n")
	root, err := AuthorizeRoot(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := root.Resolve("main.go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want, _ := filepath.EvalSymlinks(filepath.Join(dir, "main.go"))
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestMatchGlobSupportsDoubleStarAndSingleStar(t *testing.T) {
	cases := []struct {
		pattern, rel string
		want         bool
	}{
		{"*.go", "main.go", true},
		{"*.go", "internal/main.go", false},
		{"**/*.go", "internal/pkg/main.go", true},
		{"**/*.go", "main.go", true},
		{"internal/**", "internal/pkg/file.txt", true},
		{"internal/**", "other/file.txt", false},
	}
	for _, c := range cases {
		if got := matchGlob(c.pattern, c.rel); got != c.want {
			t.Errorf("matchGlob(%q, %q) = %v, want %v", c.pattern, c.rel, got, c.want)
		}
	}
}

func TestCaptureAndObserveDetectChanges(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "main.go", "package main\n")
	writeFile(t, dir, "keep.go", "package main\n// unchanged\n")
	root, err := AuthorizeRoot(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	baseline, err := Capture(root, []string{"**/*.go"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(baseline.Hashes) != 2 {
		t.Fatalf("expected 2 files in baseline, got %d", len(baseline.Hashes))
	}

	writeFile(t, dir, "main.go", "package main\n\nfunc main() {}\n")
	writeFile(t, dir, "new.go", "package main\n")
	if err := os.Remove(filepath.Join(dir, "keep.go")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := Observe(root, []string{"**/*.go"}, baseline)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	changes := map[string]string{}
	for _, c := range result.Changes {
		changes[c.Path] = c.Change
	}
	if changes["main.go"] != ChangeModified {
		t.Errorf("expected main.go modified, got %q", changes["main.go"])
	}
	if changes["new.go"] != ChangeAdded {
		t.Errorf("expected new.go added, got %q", changes["new.go"])
	}
	if changes["keep.go"] != ChangeRemoved {
		t.Errorf("expected keep.go removed, got %q", changes["keep.go"])
	}
}

func TestObserveNeverAttributesUnchangedFilesToTheStep(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "untouched.go", "package main\n")
	root, err := AuthorizeRoot(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	baseline, err := Capture(root, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := Observe(root, nil, baseline)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Changes) != 0 {
		t.Fatalf("expected no changes, got %+v", result.Changes)
	}
}

func TestFingerprintDetectsDriftAfterCollection(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "main.go", "package main\n")
	root, err := AuthorizeRoot(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	before, err := Fingerprint(root, []string{"**/*.go"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	after, err := Fingerprint(root, []string{"**/*.go"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if before != after {
		t.Fatalf("expected stable fingerprint across no-op recapture, got %q vs %q", before, after)
	}

	writeFile(t, dir, "main.go", "package main\n\nfunc main() {}\n")
	drifted, err := Fingerprint(root, []string{"**/*.go"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drifted == before {
		t.Fatal("expected fingerprint to change after workspace drift")
	}
}

func requireGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
}

func runGitCmd(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@test", "GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@test")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}

func TestCaptureUsesGitIdentityWhenAvailable(t *testing.T) {
	requireGit(t)
	dir := t.TempDir()
	runGitCmd(t, dir, "init", "-q")
	writeFile(t, dir, "main.go", "package main\n")
	runGitCmd(t, dir, "add", ".")
	runGitCmd(t, dir, "commit", "-q", "-m", "initial")

	root, err := AuthorizeRoot(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	baseline, err := Capture(root, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !baseline.FromGit {
		t.Fatal("expected baseline to be sourced from git")
	}
	if baseline.Commit == "" {
		t.Fatal("expected a non-empty commit")
	}
	if baseline.Dirty {
		t.Fatal("expected a clean working tree right after commit")
	}
}

func TestCaptureNeverMutatesTheGitIndexOrFileTimestamps(t *testing.T) {
	requireGit(t)
	dir := t.TempDir()
	runGitCmd(t, dir, "init", "-q")
	writeFile(t, dir, "main.go", "package main\n")
	runGitCmd(t, dir, "add", ".")
	runGitCmd(t, dir, "commit", "-q", "-m", "initial")

	indexPath := filepath.Join(dir, ".git", "index")
	before, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	mainInfo, err := os.Stat(filepath.Join(dir, "main.go"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	beforeModTime := mainInfo.ModTime()

	root, err := AuthorizeRoot(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := Capture(root, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, _, err := gitInfo(root); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	after, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("expected the Git index to be byte-for-byte unchanged after Capture")
	}
	mainInfoAfter, err := os.Stat(filepath.Join(dir, "main.go"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mainInfoAfter.ModTime().Equal(beforeModTime) {
		t.Fatal("expected main.go's mtime to be unchanged after Capture")
	}
}

func TestDiffTextReturnsRedactedTruncatedDiff(t *testing.T) {
	requireGit(t)
	dir := t.TempDir()
	runGitCmd(t, dir, "init", "-q")
	writeFile(t, dir, "main.go", "package main\n")
	runGitCmd(t, dir, "add", ".")
	runGitCmd(t, dir, "commit", "-q", "-m", "initial")

	writeFile(t, dir, "main.go", "package main\n\n// api_key: AKIAABCDEFGHIJKLMNOP\nfunc main() {}\n")

	root, err := AuthorizeRoot(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	diff, ok, err := DiffText(root, []string{"main.go"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected a diff to be available")
	}
	if contains(diff, "AKIAABCDEFGHIJKLMNOP") {
		t.Fatalf("expected the secret to be redacted from the diff, got: %s", diff)
	}
}

func TestCaptureFallsBackToHashesWithoutGit(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "main.go", "package main\n")
	root, err := AuthorizeRoot(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	baseline, err := Capture(root, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if baseline.FromGit {
		t.Fatal("expected fallback: no git repository was initialized")
	}
	if len(baseline.Hashes) != 1 {
		t.Fatalf("expected 1 hashed file, got %d", len(baseline.Hashes))
	}
}

func TestParseGoFileListsFuncsAndStructFields(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "user.go", `package main

type User struct {
	Name string
	Age  int
}

func NewUser() *User { return &User{} }

func (u *User) Greet() string { return "hi" }
`)
	decl, err := ParseGoFile(filepath.Join(dir, "user.go"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !decl.HasFunc("NewUser") {
		t.Error("expected NewUser to be a declared top-level func")
	}
	if decl.HasFunc("Greet") {
		t.Error("Greet is a method, not a top-level func")
	}
	userType, ok := decl.Type("User")
	if !ok || userType.Kind != "struct" {
		t.Fatalf("expected User struct, got %+v (ok=%v)", userType, ok)
	}
	if !decl.HasField("User", "Name") || !decl.HasField("User", "Age") {
		t.Fatalf("expected Name and Age fields, got %+v", userType.Fields)
	}
	if decl.HasField("User", "Missing") {
		t.Error("did not expect a Missing field")
	}
}
