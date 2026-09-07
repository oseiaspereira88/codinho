package session

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/oseiaspereira88/codinho/internal/assistance"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

func TestRecoveryRestoresSessionAndStartRetry(t *testing.T) {
	svc := newTestService(t)
	in := StartInput{ChallengeID: fixtureChallengeID, RequestID: "start", Mode: learning.ModeTeaching}
	start, err := svc.Start(in)
	if err != nil {
		t.Fatal(err)
	}
	paused, err := svc.Pause(start.SessionID, start.Revision, "pause")
	if err != nil {
		t.Fatal(err)
	}
	recovered := New(svc.catalog, svc.store, nil)
	got, err := recovered.Get(start.SessionID)
	if err != nil {
		t.Fatal(err)
	}
	if got.State != learning.SessionStatePaused || got.Revision != paused.Revision {
		t.Fatalf("got %+v", got)
	}
	retry, err := recovered.Start(in)
	if err != nil || !reflect.DeepEqual(start, retry) {
		t.Fatalf("retry %+v %v", retry, err)
	}
	next, err := recovered.Start(StartInput{ChallengeID: fixtureChallengeID, RequestID: "next"})
	if err != nil || next.SessionID == start.SessionID {
		t.Fatalf("new start %+v %v", next, err)
	}
}

func TestRejectedTransitionDoesNotAppend(t *testing.T) {
	svc := newTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatal(err)
	}
	before := svc.store.Revision(string(start.SessionID))
	_, err = svc.Resume(start.SessionID, before, "invalid-resume")
	if err == nil {
		t.Fatal("expected invalid transition")
	}
	if got := svc.store.Revision(string(start.SessionID)); got != before {
		t.Fatalf("failed call appended revision %d", got)
	}
}

func TestLegacySessionDoesNotBlockNewStarts(t *testing.T) {
	svc := newTestService(t)
	_, err := svc.store.Append("ses_1", 0, "legacy", eventstore.EventSessionStarted, map[string]string{"challenge_id": fixtureChallengeID, "mode": "practice"})
	if err != nil {
		t.Fatal(err)
	}
	recovered := New(svc.catalog, svc.store, nil)
	_, err = recovered.Get("ses_1")
	if err == nil {
		t.Fatal("legacy session must not invent state")
	}
	if errors.Is(err, ErrSessionNotFound) {
		t.Fatal("legacy data needs explicit recovery diagnostic")
	}
	if _, err := recovered.Start(StartInput{ChallengeID: fixtureChallengeID}); err != nil {
		t.Fatal(err)
	}
	before := len(svc.store.ReplayAll())
	if _, err := recovered.Start(StartInput{ChallengeID: fixtureChallengeID, RequestID: "legacy"}); !errors.Is(err, ErrSessionUnrecoverable) {
		t.Fatalf("legacy retry: %v", err)
	}
	if len(svc.store.ReplayAll()) != before {
		t.Fatal("legacy retry created a session")
	}
}

func TestRecoveryRestoresPoliciesHintsDetoursAndHistoricalRetries(t *testing.T) {
	svc := newHintedTestService(t)
	limit := 25 * time.Minute
	start := startHinted(t, svc, StartInput{Mode: learning.ModeTeaching, TimeLimit: &limit})
	hint, err := svc.HintRequest(start.SessionID, false, start.Revision, "hint")
	if err != nil {
		t.Fatal(err)
	}
	detour, err := svc.DetourStart(start.SessionID, "concept", hint.Revision, "detour")
	if err != nil {
		t.Fatal(err)
	}
	paused, err := svc.Pause(start.SessionID, detour.Revision, "pause")
	if err != nil {
		t.Fatal(err)
	}
	recovered := New(nil, svc.store, nil) // No current catalog: pinned content is authoritative.
	before, after := svc.sessions[start.SessionID], recovered.sessions[start.SessionID]
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("projection differs:\n%+v\n%+v", before, after)
	}
	resumed, err := recovered.Resume(start.SessionID, paused.Revision, "resume")
	if err != nil {
		t.Fatal(err)
	}
	done, err := recovered.DetourFinish(start.SessionID, assistance.DetourResolved, resumed.Revision, "end-detour")
	if err != nil {
		t.Fatal(err)
	}
	retry, err := recovered.Pause(start.SessionID, detour.Revision, "pause")
	if err != nil || !reflect.DeepEqual(retry, paused) {
		t.Fatalf("historical pause: %+v %v", retry, err)
	}
	repeatHint, err := recovered.HintRequest(start.SessionID, false, start.Revision, "hint")
	if err != nil || !reflect.DeepEqual(hint, repeatHint) {
		t.Fatalf("hint: %+v %v", repeatHint, err)
	}
	if _, err := recovered.SyntaxRecallGet(start.SessionID, done.Revision, "hint"); !errors.Is(err, eventstore.ErrRevisionConflict) {
		t.Fatalf("cross-tool request reuse: %v", err)
	}
	if _, err := recovered.HintRequest(start.SessionID, true, start.Revision, "hint"); !errors.Is(err, eventstore.ErrRevisionConflict) {
		t.Fatalf("different input: %v", err)
	}
	if svc.store.Revision(string(start.SessionID)) != done.Revision {
		t.Fatal("retry appended")
	}
}

func TestRecoveryRestoresGranularityEvaluationAndAdvance(t *testing.T) {
	svc := newTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatal(err)
	}
	granularity, err := svc.GranularityAdjust(start.SessionID, learning.DepthMicro, "focus", start.Revision, "granularity")
	if err != nil {
		t.Fatal(err)
	}
	recovered := New(nil, svc.store, nil)
	if !reflect.DeepEqual(svc.sessions[start.SessionID], recovered.sessions[start.SessionID]) {
		t.Fatal("depth/node not restored")
	}
	evalIn := EvaluateInput{SessionID: start.SessionID, ExpectedRevision: granularity.Revision, RequestID: "evaluation", SubmissionIntent: true}
	eval, err := recovered.StepEvaluate(evalIn)
	if err != nil {
		t.Fatal(err)
	}
	complete, err := recovered.StepComplete(start.SessionID, true, false, eval.Revision, "complete")
	if err != nil {
		t.Fatal(err)
	}
	again := New(nil, svc.store, nil)
	retry, err := again.StepEvaluate(evalIn)
	if err != nil || !reflect.DeepEqual(eval, retry) {
		t.Fatalf("evaluation retry after completion: %+v %v", retry, err)
	}
	if again.store.Revision(string(start.SessionID)) != complete.Revision {
		t.Fatal("evaluation retry appended")
	}

	// Exercise real advancement separately from the current tree-navigation gap.
	svc = newPolicyTestService(t)
	start, err = svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatal(err)
	}
	advance, err := svc.StepAdvance(start.SessionID, true, start.Revision, "advance")
	if err != nil {
		t.Fatal(err)
	}
	again = New(nil, svc.store, nil)
	repeated, err := again.StepAdvance(start.SessionID, true, start.Revision, "advance")
	if err != nil || !reflect.DeepEqual(advance, repeated) {
		t.Fatalf("advance retry: %+v %v", repeated, err)
	}
}

func TestRecoveryCompletesInterruptedSubmission(t *testing.T) {
	svc := newTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatal(err)
	}
	in := EvaluateInput{SessionID: start.SessionID, ExpectedRevision: start.Revision, RequestID: "evaluate", SubmissionIntent: true}
	original, err := svc.StepEvaluate(in)
	if err != nil {
		t.Fatal(err)
	}
	// Copy only the durable prefix preceding attempt append, simulating a crash
	// after evaluation fsync and before its secondary write/application.
	interrupted, err := eventstore.Open(filepath.Join(t.TempDir(), "events.jsonl"), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer interrupted.Close()
	for _, ev := range svc.store.ReplayAll() {
		if ev.Type == eventstore.EventAttemptSubmitted {
			break
		}
		if _, err := interrupted.Append(ev.StreamID, ev.Revision-1, ev.RequestID, ev.Type, ev.Payload); err != nil {
			t.Fatal(err)
		}
	}
	recovered := New(nil, interrupted, nil)
	retry, err := recovered.StepEvaluate(in)
	if err != nil || !reflect.DeepEqual(original, retry) {
		t.Fatalf("lost response retry: %+v %v", retry, err)
	}
	twice := New(nil, interrupted, nil)
	if twice.store.Revision(string(start.SessionID)) != original.Revision {
		t.Fatal("repair duplicated attempt")
	}
}

func TestRecoveryRejectsDamagedPinnedContent(t *testing.T) {
	svc := newTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID, RequestID: "original"})
	if err != nil {
		t.Fatal(err)
	}
	var payload startedPayload
	if err := json.Unmarshal(svc.store.Replay(string(start.SessionID))[0].Payload, &payload); err != nil {
		t.Fatal(err)
	}
	payload.Challenge.Title = "changed without matching digest"
	store, err := eventstore.Open(filepath.Join(t.TempDir(), "events.jsonl"), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	if _, err := store.Append("ses_1", 0, "original", eventstore.EventSessionStarted, payload); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Append("ses_1", 1, "pending", eventstore.EventEvaluationRecorded, map[string]any{"step_id": "fixture.macro-one", "request_digest": "synthetic", "submission_intent": true}); err != nil {
		t.Fatal(err)
	}
	recovered := New(svc.catalog, store, nil)
	if _, err := recovered.Get("ses_1"); !errors.Is(err, ErrSessionUnrecoverable) {
		t.Fatalf("damaged content: %v", err)
	}
	if _, err := recovered.Start(StartInput{ChallengeID: fixtureChallengeID, RequestID: "original"}); !errors.Is(err, ErrSessionUnrecoverable) {
		t.Fatalf("damaged start retry: %v", err)
	}
	if store.Revision("ses_1") != 2 {
		t.Fatal("startup modified incompatible history")
	}
}
