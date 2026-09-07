package application

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/oseiaspereira88/codinho/internal/eventstore"
)

func TestWorkspaceRecoveryBaselineScopeAndRetry(t *testing.T) {
	svc := newTestWorkspaceService(t)
	root := t.TempDir()
	writeWorkspaceFile(t, root, "main.go", "package main\n")
	input := ObserveInput{SessionID: "ses_1", StepID: "step_1", Root: root, Globs: []string{"**/*.go"}, RequestID: "baseline"}
	baseline, err := svc.Observe(input)
	if err != nil {
		t.Fatal(err)
	}
	recovered := NewWorkspaceService(svc.store, svc.evidence)
	writeWorkspaceFile(t, root, "main.go", "package main\nfunc main() {}\n")
	retry, err := recovered.Observe(input)
	if err != nil || !reflect.DeepEqual(baseline, retry) {
		t.Fatalf("retry: %+v %v", retry, err)
	}
	input.RequestID, input.ExpectedRevision = "diff", baseline.Revision
	diff, err := recovered.Observe(input)
	if err != nil || diff.Baseline || len(diff.Changes) != 1 {
		t.Fatalf("restored diff: %+v %v", diff, err)
	}
	if _, err := recovered.EvidenceGet(EvidenceGetInput{SessionID: "ses_2", EvidenceID: baseline.EvidenceID}); !errors.Is(err, ErrEvidenceOutOfScope) {
		t.Fatalf("cross-session: %v", err)
	}
	old, err := recovered.EvidenceGet(EvidenceGetInput{SessionID: "ses_1", EvidenceID: baseline.EvidenceID})
	if err != nil || !old.Stale {
		t.Fatalf("baseline must be stale: %+v %v", old, err)
	}
	fresh, err := recovered.EvidenceGet(EvidenceGetInput{SessionID: "ses_1", EvidenceID: diff.EvidenceID})
	if err != nil || fresh.Stale {
		t.Fatalf("diff must be current: %+v %v", fresh, err)
	}
	input.Root = t.TempDir()
	input.ExpectedRevision = diff.Revision
	if _, err := recovered.Observe(input); !errors.Is(err, ErrWorkspaceRootInvalid) {
		t.Fatalf("changed scope: %v", err)
	}
}

func TestWorkspaceRecoveryRejectsReplacedRoot(t *testing.T) {
	svc := newTestWorkspaceService(t)
	parent := t.TempDir()
	root := filepath.Join(parent, "learner")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatal(err)
	}
	writeWorkspaceFile(t, root, "main.go", "package main\n")
	baseline, err := svc.Observe(ObserveInput{SessionID: "ses_1", StepID: "step_1", Root: root})
	if err != nil {
		t.Fatal(err)
	}
	recovered := NewWorkspaceService(svc.store, svc.evidence)
	if err := os.Rename(root, filepath.Join(parent, "original")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(t.TempDir(), root); err != nil {
		t.Fatal(err)
	}
	if _, _, err := recovered.RootFor("ses_1", "step_1"); !errors.Is(err, ErrWorkspaceRootInvalid) {
		t.Fatalf("followed replaced root: %v", err)
	}
	stale, err := recovered.EvidenceGet(EvidenceGetInput{SessionID: "ses_1", EvidenceID: baseline.EvidenceID})
	if err != nil || !stale.Stale {
		t.Fatalf("replaced root freshness: %+v %v", stale, err)
	}
}

func TestWorkspaceRecoveryCheckEvidenceAndRetry(t *testing.T) {
	// Broken source after the first check proves retry does not rerun it.
	sessions, ws, svc := newChecksTestServices(t)
	start, err := sessions.Start(StartInput{ChallengeID: "fixture.challenge-one"})
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	writeWorkspaceFile(t, root, "go.mod", "module fixture\n\ngo 1.25\n")
	writeWorkspaceFile(t, root, "main.go", "package main\nfunc main() {}\n")
	baseline, err := ws.Observe(ObserveInput{SessionID: start.SessionID, StepID: start.ActiveStep, Root: root, Globs: []string{"**/*.go"}, ExpectedRevision: start.Revision, RequestID: "observe"})
	if err != nil {
		t.Fatal(err)
	}
	input := CheckRunInput{SessionID: start.SessionID, CheckID: "focused-tests", ExpectedRevision: baseline.Revision, RequestID: "check"}
	first, err := svc.Run(input)
	if err != nil {
		t.Fatal(err)
	}
	writeWorkspaceFile(t, root, "main.go", "broken Go\n")
	prior, err := ws.EvidenceGet(EvidenceGetInput{SessionID: start.SessionID, EvidenceID: first.EvidenceID})
	if err != nil || !prior.Stale {
		t.Fatalf("runtime check scope: %+v %v", prior, err)
	}
	recoveredWS := NewWorkspaceService(ws.store, ws.evidence)
	recovered := NewChecksService(ws.store, sessions, recoveredWS)
	second, err := recovered.Run(input)
	if err != nil || !reflect.DeepEqual(first, second) {
		t.Fatalf("check reran: %+v %v", second, err)
	}
	if ws.store.Revision(string(start.SessionID)) != first.Revision {
		t.Fatal("retry appended")
	}
	input.RequestID = "observe"
	if _, err := recovered.Run(input); !errors.Is(err, eventstore.ErrRevisionConflict) {
		t.Fatalf("cross-tool retry: %v", err)
	}
	input.RequestID, input.CheckID = "check", "different"
	if _, err := recovered.Run(input); !errors.Is(err, eventstore.ErrRevisionConflict) {
		t.Fatalf("changed check ID: %v", err)
	}
}
