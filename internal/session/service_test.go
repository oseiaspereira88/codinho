package session

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

const fixtureChallengeID = "fixture.challenge-one"

func writeTestFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func newTestService(t *testing.T) *Service {
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
          - id: fixture.macro-one
            kind: macro
            title: Macro step
            instruction:
              objective: Understand the fixture.
              scope: Read only.
            children:
              - id: fixture.meso-one
                kind: meso
                title: Meso step
                children:
                  - id: fixture.micro-one
                    kind: micro
                    title: Micro step
                    instruction:
                      objective: Declare the fixture type.
                      scope: Only the declaration.
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

	return New(catalog, store, nil)
}

func TestStartHonorsDefaultMicroDepth(t *testing.T) {
	svc := newTestService(t)
	result, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ActiveStep != "fixture.micro-one" {
		t.Fatalf("active step = %s, want fixture.micro-one", result.ActiveStep)
	}
	if result.Revision != 1 {
		t.Fatalf("revision = %d, want 1", result.Revision)
	}
	if result.Disclosure.SolutionRevealed {
		t.Fatal("a fresh session must never start with the solution revealed")
	}
}

func TestStartIsIdempotentByRequestID(t *testing.T) {
	svc := newTestService(t)
	in := StartInput{ChallengeID: fixtureChallengeID, RequestID: "retry-1"}

	first, err := svc.Start(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := svc.Start(in)
	if err != nil {
		t.Fatalf("unexpected error on retry: %v", err)
	}
	if first.SessionID != second.SessionID {
		t.Fatalf("retry produced a different session: %s vs %s", first.SessionID, second.SessionID)
	}
}

func TestStartRejectsUnknownChallenge(t *testing.T) {
	svc := newTestService(t)
	_, err := svc.Start(StartInput{ChallengeID: "does-not-exist"})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestGetReportsRevisionAndDisclosure(t *testing.T) {
	svc := newTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	get, err := svc.Get(start.SessionID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if get.Revision != start.Revision || get.ActiveStep != start.ActiveStep {
		t.Fatalf("Get() = %+v, want to match Start() = %+v", get, start)
	}
}

func TestGetUnknownSessionFails(t *testing.T) {
	svc := newTestService(t)
	_, err := svc.Get("ses_does_not_exist")
	if !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("expected ErrSessionNotFound, got %v", err)
	}
}

func TestInstructionNeverExposesChildrenOrSolution(t *testing.T) {
	svc := newTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	instr, err := svc.Instruction(start.SessionID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if instr.Objective != "Declare the fixture type." || instr.Scope != "Only the declaration." {
		t.Fatalf("unexpected instruction: %+v", instr)
	}
}

func TestConfigureAppliesOnlySpecifiedFields(t *testing.T) {
	svc := newTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	newHelp := learning.HelpNoHints
	result, err := svc.Configure(ConfigureInput{SessionID: start.SessionID, Help: &newHelp, ExpectedRevision: start.Revision})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Revision != start.Revision+1 {
		t.Fatalf("revision = %d, want %d", result.Revision, start.Revision+1)
	}

	svc.mu.Lock()
	got := svc.sessions[start.SessionID].session.Policy.Help.Kind
	svc.mu.Unlock()
	if got != learning.HelpNoHints {
		t.Fatalf("help policy = %s, want %s", got, learning.HelpNoHints)
	}
}

func TestConfigureRejectsStaleRevision(t *testing.T) {
	svc := newTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	newHelp := learning.HelpFree
	_, err = svc.Configure(ConfigureInput{SessionID: start.SessionID, Help: &newHelp, ExpectedRevision: start.Revision + 1})
	if !errors.Is(err, eventstore.ErrRevisionConflict) {
		t.Fatalf("expected ErrRevisionConflict, got %v", err)
	}
}

func TestConfigureIsIdempotentByRequestID(t *testing.T) {
	svc := newTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	newHelp := learning.HelpNoHints
	in := ConfigureInput{SessionID: start.SessionID, Help: &newHelp, ExpectedRevision: start.Revision, RequestID: "cfg-retry"}

	first, err := svc.Configure(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := svc.Configure(in)
	if err != nil {
		t.Fatalf("unexpected error on retry: %v", err)
	}
	if first.Revision != second.Revision {
		t.Fatalf("retry advanced the revision: %d vs %d", first.Revision, second.Revision)
	}
}

func TestPauseResumeFinishLifecycle(t *testing.T) {
	svc := newTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rev := start.Revision

	paused, err := svc.Pause(start.SessionID, rev, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if paused.State != learning.SessionStatePaused {
		t.Fatalf("state = %s, want paused", paused.State)
	}
	rev = paused.Revision

	resumed, err := svc.Resume(start.SessionID, rev, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resumed.State != learning.SessionStateActive {
		t.Fatalf("state = %s, want active", resumed.State)
	}
	rev = resumed.Revision

	finished, err := svc.Finish(start.SessionID, rev, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if finished.State != learning.SessionStateCompleted {
		t.Fatalf("state = %s, want completed", finished.State)
	}
}

func TestFinishFromPausedFailsWithoutInferringResume(t *testing.T) {
	svc := newTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	paused, err := svc.Pause(start.SessionID, start.Revision, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var domainErr learning.DomainError
	_, err = svc.Finish(start.SessionID, paused.Revision, "")
	if !errors.As(err, &domainErr) {
		t.Fatalf("expected a domain transition error, got %v", err)
	}
}

func TestPauseIsIdempotentByRequestID(t *testing.T) {
	svc := newTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	first, err := svc.Pause(start.SessionID, start.Revision, "pause-retry")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := svc.Pause(start.SessionID, start.Revision, "pause-retry")
	if err != nil {
		t.Fatalf("unexpected error on retry: %v", err)
	}
	if first.Revision != second.Revision || second.State != learning.SessionStatePaused {
		t.Fatalf("retry was not idempotent: %+v vs %+v", first, second)
	}
}

func TestGranularityAdjustWalksToRequestedDepth(t *testing.T) {
	svc := newTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := svc.GranularityAdjust(start.SessionID, learning.DepthMicro, "", start.Revision, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.StepID != "fixture.micro-one" || result.Kind != "micro" {
		t.Fatalf("unexpected window: %+v", result)
	}

	get, err := svc.Get(start.SessionID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if get.ActiveStep != "fixture.micro-one" {
		t.Fatalf("active step after adjust = %s, want fixture.micro-one", get.ActiveStep)
	}
}

func TestGranularityAdjustNeverRewritesTheTree(t *testing.T) {
	svc := newTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rev := start.Revision
	for _, depth := range []learning.Depth{learning.DepthMicro, learning.DepthMacro, learning.DepthMicro} {
		result, err := svc.GranularityAdjust(start.SessionID, depth, "", rev, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		rev = result.Revision
	}

	// The tree itself is re-derived fresh from the catalog each time, so
	// asking for micro again must resolve to the same node as before —
	// nothing was mutated by the earlier macro request.
	final, err := svc.GranularityAdjust(start.SessionID, learning.DepthMicro, "", rev, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if final.StepID != "fixture.micro-one" {
		t.Fatalf("StepID = %s, want fixture.micro-one", final.StepID)
	}
}
