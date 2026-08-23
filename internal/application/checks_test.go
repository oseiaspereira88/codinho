package application

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/oseiaspereira88/codinho/internal/checks"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/evidence"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

func newChecksTestServices(t *testing.T) (*SessionService, *WorkspaceService, *ChecksService) {
	t.Helper()
	dir := t.TempDir()
	writeTestFile(t, dir, "manifest.yaml", "schema_version: 1\npacks:\n  - pack.yaml\n")
	writeTestFile(t, dir, "pack.yaml", `schema_version: 1
id: fixture-pack
version: 1.0.0
challenges:
  - schema_version: 1
    id: fixture.challenge-one
    version: 1.0.0
    title: Fixture challenge
    kind: atomic
    difficulty: foundational
    layers:
      - id: understanding
        macro_steps:
          - id: fixture.step-one
            kind: micro
            title: Fixture step
            instruction:
              objective: Declare the fixture type.
              scope: Only the declaration.
    checks:
      - id: focused-tests
        runner: go_test
        package: ./...
`)
	catalog, diags, err := curriculum.Load(dir, curriculum.DefaultLimits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostics: %+v", diags)
	}

	store, err := eventstore.Open(filepath.Join(t.TempDir(), "events.jsonl"), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	evStore, err := evidence.Open(filepath.Join(t.TempDir(), "evidence"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	catalogService := NewCatalogService(catalog)
	sessions := NewSessionService(catalogService, store, evStore)
	ws := NewWorkspaceService(store, evStore)
	checksSvc := NewChecksService(store, sessions, ws)
	return sessions, ws, checksSvc
}

func newFixtureGoModule(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module fixture\n\ngo 1.25\n"), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "main_test.go"), []byte("package main\n\nimport \"testing\"\n\nfunc TestOK(t *testing.T) {}\n"), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return dir
}

func TestChecksRunEndToEndAfterWorkspaceObserve(t *testing.T) {
	sessions, ws, checksSvc := newChecksTestServices(t)
	start, err := sessions.Start(StartInput{ChallengeID: "fixture.challenge-one"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	root := newFixtureGoModule(t)
	obs, err := ws.Observe(ObserveInput{SessionID: start.SessionID, StepID: start.ActiveStep, Root: root, Globs: []string{"**/*.go"}, ExpectedRevision: start.Revision})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := checksSvc.Run(CheckRunInput{SessionID: start.SessionID, CheckID: "focused-tests", ExpectedRevision: obs.Revision})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Outcome != string(checks.OutcomePass) {
		t.Fatalf("expected pass, got %s", result.Outcome)
	}
	if result.EvidenceID == "" {
		t.Fatal("expected a non-empty evidence ID")
	}

	// The same evidence_get already exposed by WorkspaceService must scope
	// and serve check-produced evidence too (requirement R8).
	got, err := ws.EvidenceGet(EvidenceGetInput{SessionID: start.SessionID, EvidenceID: result.EvidenceID, Root: root, Globs: []string{"**/*.go"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Stale {
		t.Fatal("expected fresh check evidence to not be stale immediately after Run")
	}
}

func TestChecksRunRejectsUnknownCheckID(t *testing.T) {
	sessions, ws, checksSvc := newChecksTestServices(t)
	start, err := sessions.Start(StartInput{ChallengeID: "fixture.challenge-one"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	root := newFixtureGoModule(t)
	if _, err := ws.Observe(ObserveInput{SessionID: start.SessionID, StepID: start.ActiveStep, Root: root, ExpectedRevision: start.Revision}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = checksSvc.Run(CheckRunInput{SessionID: start.SessionID, CheckID: "does-not-exist", ExpectedRevision: 1})
	if !errors.Is(err, ErrCheckNotFound) {
		t.Fatalf("expected ErrCheckNotFound, got %v", err)
	}
}

func TestChecksRunRequiresAPriorWorkspaceObservation(t *testing.T) {
	sessions, _, checksSvc := newChecksTestServices(t)
	start, err := sessions.Start(StartInput{ChallengeID: "fixture.challenge-one"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = checksSvc.Run(CheckRunInput{SessionID: start.SessionID, CheckID: "focused-tests", ExpectedRevision: start.Revision})
	if !errors.Is(err, ErrNoWorkspaceBaseline) {
		t.Fatalf("expected ErrNoWorkspaceBaseline, got %v", err)
	}
}

func TestChecksRunReportsFailOnBrokenLearnerCode(t *testing.T) {
	sessions, ws, checksSvc := newChecksTestServices(t)
	start, err := sessions.Start(StartInput{ChallengeID: "fixture.challenge-one"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	root := newFixtureGoModule(t)
	if err := os.WriteFile(filepath.Join(root, "main_test.go"), []byte("package main\n\nimport \"testing\"\n\nfunc TestFails(t *testing.T) { t.Fatal(\"boom\") }\n"), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	obs, err := ws.Observe(ObserveInput{SessionID: start.SessionID, StepID: start.ActiveStep, Root: root, ExpectedRevision: start.Revision})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := checksSvc.Run(CheckRunInput{SessionID: start.SessionID, CheckID: "focused-tests", ExpectedRevision: obs.Revision})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Outcome != string(checks.OutcomeFail) {
		t.Fatalf("expected fail, got %s", result.Outcome)
	}
}

func TestChecksRunFeedsStepEvaluateStructuralVerdict(t *testing.T) {
	sessions, ws, checksSvc := newChecksTestServices(t)
	start, err := sessions.Start(StartInput{ChallengeID: "fixture.challenge-one"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	root := newFixtureGoModule(t)
	obs, err := ws.Observe(ObserveInput{SessionID: start.SessionID, StepID: start.ActiveStep, Root: root, ExpectedRevision: start.Revision})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	runResult, err := checksSvc.Run(CheckRunInput{SessionID: start.SessionID, CheckID: "focused-tests", ExpectedRevision: obs.Revision})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	eval, err := sessions.StepEvaluate(EvaluateInput{
		SessionID: start.SessionID,
		Criteria: []CriterionInput{
			{Name: "tests-pass", Kind: learning.StructuralCriterionKind, Severity: learning.SeverityBlocking, EvidenceID: learning.EvidenceID(runResult.EvidenceID)},
		},
		ExpectedRevision: runResult.Revision,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if eval.Criteria[0].Verdict != learning.VerdictMet {
		t.Fatalf("verdict = %s, want met (real check outcome: pass)", eval.Criteria[0].Verdict)
	}
	if eval.HasBlockingFailure {
		t.Fatal("expected no blocking failure when the real check passed")
	}
}
