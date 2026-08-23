package application

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/evidence"
	"github.com/oseiaspereira88/codinho/internal/learning"
	"github.com/oseiaspereira88/codinho/internal/workspace"
)

func newTestWorkspaceService(t *testing.T) *WorkspaceService {
	t.Helper()
	dir := t.TempDir()
	store, err := eventstore.Open(filepath.Join(dir, "events.jsonl"), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	evStore, err := evidence.Open(filepath.Join(dir, "evidence"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return NewWorkspaceService(store, evStore)
}

func writeWorkspaceFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestObserveEstablishesBaselineOnFirstCall(t *testing.T) {
	svc := newTestWorkspaceService(t)
	root := t.TempDir()
	writeWorkspaceFile(t, root, "main.go", "package main\n")

	result, err := svc.Observe(ObserveInput{
		SessionID: "sess-1", StepID: "step-1", Root: root, Globs: []string{"**/*.go"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Baseline {
		t.Fatal("expected the first observation to establish a baseline")
	}
	if len(result.Changes) != 0 {
		t.Fatalf("expected no changes on baseline, got %+v", result.Changes)
	}
	if result.Observation.Evidence.ID() == "" {
		t.Fatal("expected a non-empty evidence ID")
	}
	if result.Revision != 1 {
		t.Fatalf("expected revision 1, got %d", result.Revision)
	}
}

func TestObserveReportsChangesAgainstBaseline(t *testing.T) {
	svc := newTestWorkspaceService(t)
	root := t.TempDir()
	writeWorkspaceFile(t, root, "main.go", "package main\n")

	first, err := svc.Observe(ObserveInput{SessionID: "sess-1", StepID: "step-1", Root: root, Globs: []string{"**/*.go"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	writeWorkspaceFile(t, root, "main.go", "package main\n\nfunc main() {}\n")
	second, err := svc.Observe(ObserveInput{
		SessionID: "sess-1", StepID: "step-1", Root: root, Globs: []string{"**/*.go"}, ExpectedRevision: first.Revision,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if second.Baseline {
		t.Fatal("expected the second observation to diff against the baseline, not re-establish it")
	}
	if len(second.Changes) != 1 || second.Changes[0].Path != "main.go" || second.Changes[0].Change != workspace.ChangeModified {
		t.Fatalf("expected exactly one modified change for main.go, got %+v", second.Changes)
	}
	if second.Revision != 2 {
		t.Fatalf("expected revision 2, got %d", second.Revision)
	}
}

func TestObserveRejectsRevisionConflict(t *testing.T) {
	svc := newTestWorkspaceService(t)
	root := t.TempDir()
	writeWorkspaceFile(t, root, "main.go", "package main\n")

	if _, err := svc.Observe(ObserveInput{SessionID: "sess-1", StepID: "step-1", Root: root}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err := svc.Observe(ObserveInput{SessionID: "sess-1", StepID: "step-1", Root: root, ExpectedRevision: 0})
	if !errors.Is(err, eventstore.ErrRevisionConflict) {
		t.Fatalf("expected ErrRevisionConflict, got %v", err)
	}
}

func TestObserveRejectsInvalidRoot(t *testing.T) {
	svc := newTestWorkspaceService(t)
	_, err := svc.Observe(ObserveInput{SessionID: "sess-1", StepID: "step-1", Root: filepath.Join(t.TempDir(), "does-not-exist")})
	if !errors.Is(err, ErrWorkspaceRootInvalid) {
		t.Fatalf("expected ErrWorkspaceRootInvalid, got %v", err)
	}
}

func TestEvidenceGetIsScopedToTheRecordingSession(t *testing.T) {
	svc := newTestWorkspaceService(t)
	root := t.TempDir()
	writeWorkspaceFile(t, root, "main.go", "package main\n")

	result, err := svc.Observe(ObserveInput{SessionID: "sess-1", StepID: "step-1", Root: root})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := svc.EvidenceGet(EvidenceGetInput{SessionID: "sess-1", EvidenceID: result.EvidenceID}); err != nil {
		t.Fatalf("unexpected error for the owning session: %v", err)
	}
	_, err = svc.EvidenceGet(EvidenceGetInput{SessionID: learning.SessionID("sess-2"), EvidenceID: result.EvidenceID})
	if !errors.Is(err, ErrEvidenceOutOfScope) {
		t.Fatalf("expected ErrEvidenceOutOfScope for a different session, got %v", err)
	}
}

func TestEvidenceGetDetectsDriftSinceCollection(t *testing.T) {
	svc := newTestWorkspaceService(t)
	root := t.TempDir()
	writeWorkspaceFile(t, root, "main.go", "package main\n")

	result, err := svc.Observe(ObserveInput{SessionID: "sess-1", StepID: "step-1", Root: root, Globs: []string{"**/*.go"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fresh, err := svc.EvidenceGet(EvidenceGetInput{SessionID: "sess-1", EvidenceID: result.EvidenceID, Root: root, Globs: []string{"**/*.go"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fresh.Stale {
		t.Fatal("expected fresh evidence to not be stale")
	}

	writeWorkspaceFile(t, root, "main.go", "package main\n\nfunc main() {}\n")
	drifted, err := svc.EvidenceGet(EvidenceGetInput{SessionID: "sess-1", EvidenceID: result.EvidenceID, Root: root, Globs: []string{"**/*.go"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drifted.Stale {
		t.Fatal("expected evidence to be reported stale after the workspace changed")
	}
}
