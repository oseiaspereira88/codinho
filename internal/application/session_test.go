package application

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

func newTestServices(t *testing.T) (*CatalogService, *SessionService) {
	t.Helper()
	dir := t.TempDir()
	manifest := "schema_version: 1\npacks:\n  - pack.yaml\n"
	pack := `schema_version: 1
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
`
	writeTestFile(t, dir, "manifest.yaml", manifest)
	writeTestFile(t, dir, "pack.yaml", pack)

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

	catalogService := NewCatalogService(catalog)
	return catalogService, NewSessionService(catalogService, store, nil)
}

func TestSessionServiceStartAssignsActiveStep(t *testing.T) {
	_, sessions := newTestServices(t)
	result, err := sessions.Start(StartInput{ChallengeID: "fixture.challenge-one"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.SessionID == "" || result.ActiveStep != "fixture.step-one" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestSessionServiceStartIsIdempotentByRequestID(t *testing.T) {
	_, sessions := newTestServices(t)
	in := StartInput{ChallengeID: "fixture.challenge-one", RequestID: "retry-1"}

	first, err := sessions.Start(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := sessions.Start(in)
	if err != nil {
		t.Fatalf("unexpected error on retry: %v", err)
	}
	if first.SessionID != second.SessionID {
		t.Fatalf("retry produced a different session: %s vs %s", first.SessionID, second.SessionID)
	}

	// Regression guard: a distinct RequestID must still create a distinct
	// session, or idempotency would degenerate into always-reuse.
	third, err := sessions.Start(StartInput{ChallengeID: "fixture.challenge-one", RequestID: "retry-2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if third.SessionID == first.SessionID {
		t.Fatal("a different RequestID must not reuse the same session")
	}
}

func TestSessionServiceStartRejectsUnknownChallenge(t *testing.T) {
	_, sessions := newTestServices(t)
	_, err := sessions.Start(StartInput{ChallengeID: "does-not-exist"})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestSessionServiceGetAndInstructionRoundTrip(t *testing.T) {
	_, sessions := newTestServices(t)
	start, err := sessions.Start(StartInput{ChallengeID: "fixture.challenge-one"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	get, err := sessions.Get(start.SessionID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if get.State != learning.SessionStateActive || get.ActiveStep != start.ActiveStep {
		t.Fatalf("unexpected get result: %+v", get)
	}

	instr, err := sessions.Instruction(start.SessionID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if instr.Objective != "Declare the fixture type." {
		t.Fatalf("unexpected instruction: %+v", instr)
	}
}

func TestSessionServiceGetUnknownSessionFails(t *testing.T) {
	_, sessions := newTestServices(t)
	_, err := sessions.Get("ses_does_not_exist")
	if !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("expected ErrSessionNotFound, got %v", err)
	}
}

func writeTestFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
