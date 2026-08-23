package security

import (
	"os"
	"path/filepath"
	"testing"
)

func TestContainsSecretShapedContentDetectsKnownPatterns(t *testing.T) {
	if !ContainsSecretShapedContent([]byte("api_key: sk_live_abcdefghijklmnop")) {
		t.Fatal("expected a secret-shaped api_key assignment to be detected")
	}
	if ContainsSecretShapedContent([]byte("package main\n\nfunc main() {}\n")) {
		t.Fatal("ordinary Go source must never be flagged")
	}
}

func seedState(t *testing.T, root string) {
	t.Helper()
	stateDir := StatePath(root)
	if err := os.MkdirAll(filepath.Join(stateDir, EvidenceDirName), 0o700); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(stateDir, EventsFileName), []byte(`{"id":"ev1"}`+"\n"), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(stateDir, EvidenceDirName, "abc123"), []byte("evidence content"), 0o400); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExportCopiesEventsAndEvidenceWithoutTouchingSource(t *testing.T) {
	root := t.TempDir()
	seedState(t, root)
	dest := t.TempDir()

	if err := Export(root, dest); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gotEvents, err := os.ReadFile(filepath.Join(dest, EventsFileName))
	if err != nil || string(gotEvents) != `{"id":"ev1"}`+"\n" {
		t.Fatalf("events = %q, err = %v", gotEvents, err)
	}
	gotEvidence, err := os.ReadFile(filepath.Join(dest, EvidenceDirName, "abc123"))
	if err != nil || string(gotEvidence) != "evidence content" {
		t.Fatalf("evidence = %q, err = %v", gotEvidence, err)
	}

	// The source must be untouched.
	if _, err := os.Stat(filepath.Join(StatePath(root), EventsFileName)); err != nil {
		t.Fatalf("source events file was removed or modified: %v", err)
	}
}

func TestExportWithNoStateIsANoOp(t *testing.T) {
	root := t.TempDir()
	dest := t.TempDir()
	if err := Export(root, dest); err != nil {
		t.Fatalf("unexpected error exporting from an empty workspace: %v", err)
	}
}

func TestPurgeRequiresExplicitConfirmation(t *testing.T) {
	root := t.TempDir()
	seedState(t, root)

	if err := Purge(root, false); err != ErrPurgeNotConfirmed {
		t.Fatalf("err = %v, want ErrPurgeNotConfirmed", err)
	}
	if _, err := os.Stat(filepath.Join(StatePath(root), EventsFileName)); err != nil {
		t.Fatalf("state was removed despite confirm=false: %v", err)
	}
}

func TestPurgeRemovesOnlyLocalState(t *testing.T) {
	root := t.TempDir()
	seedState(t, root)
	// A file that looks like the learner's own project content, sitting
	// right next to .codinho — Purge must never touch it.
	projectFile := filepath.Join(root, "main.go")
	if err := os.WriteFile(projectFile, []byte("package main\n"), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := Purge(root, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(StatePath(root)); !os.IsNotExist(err) {
		t.Fatalf("state directory still exists after purge: %v", err)
	}
	if _, err := os.Stat(projectFile); err != nil {
		t.Fatalf("purge touched the learner's project file: %v", err)
	}
}
