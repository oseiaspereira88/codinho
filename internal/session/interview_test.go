package session

import (
	"errors"
	"testing"
	"time"

	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

func TestInterviewElapsedWithoutLimit(t *testing.T) {
	start := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	now := start.Add(30 * time.Minute)
	elapsed, timedOut := InterviewElapsed(start, now, nil)
	if elapsed != 30*time.Minute {
		t.Fatalf("elapsed = %s", elapsed)
	}
	if timedOut {
		t.Fatal("no limit set: must never report timed out")
	}
}

func TestInterviewElapsedRespectsLimit(t *testing.T) {
	start := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	limit := 45 * time.Minute

	before, timedOut := InterviewElapsed(start, start.Add(30*time.Minute), &limit)
	if timedOut {
		t.Fatalf("elapsed %s should not exceed limit %s", before, limit)
	}
	after, timedOut := InterviewElapsed(start, start.Add(46*time.Minute), &limit)
	if !timedOut {
		t.Fatalf("elapsed %s should exceed limit %s", after, limit)
	}
}

func startInterviewSession(t *testing.T, svc *Service, timeLimit *time.Duration) StartResult {
	t.Helper()
	result, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID, Mode: learning.ModeInterview, TimeLimit: timeLimit})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return result
}

func TestStartAcceptsTimeLimit(t *testing.T) {
	svc := newTestService(t)
	limit := 90 * time.Minute
	start := startInterviewSession(t, svc, &limit)

	rec := svc.sessions[start.SessionID]
	if rec.session.Policy.TimeLimit == nil || *rec.session.Policy.TimeLimit != limit {
		t.Fatalf("TimeLimit = %v, want %s", rec.session.Policy.TimeLimit, limit)
	}
}

func TestInterviewStatusUsesRealRecordedStartTime(t *testing.T) {
	svc := newTestService(t)
	limit := 45 * time.Minute
	start := startInterviewSession(t, svc, &limit)

	events := svc.store.Replay(string(start.SessionID))
	startedAt, err := time.Parse(time.RFC3339Nano, events[0].RecordedAt)
	if err != nil {
		t.Fatalf("parsing recorded start time: %v", err)
	}

	status, err := svc.InterviewStatus(start.SessionID, startedAt.Add(50*time.Minute))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.TimedOut {
		t.Fatal("50 minutes elapsed against a 45 minute limit must report timed out")
	}
	if status.Elapsed != 50*time.Minute {
		t.Fatalf("Elapsed = %s, want 50m", status.Elapsed)
	}
}

func TestInterviewStatusWithoutLimitNeverTimesOut(t *testing.T) {
	svc := newTestService(t)
	start := startInterviewSession(t, svc, nil)

	status, err := svc.InterviewStatus(start.SessionID, time.Now().Add(1000*time.Hour))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.TimedOut {
		t.Fatal("a session with no time limit must never time out")
	}
}

func TestRecordBlockedHintAttemptAdvancesRevisionLikeAnyOtherEvent(t *testing.T) {
	svc := newTestService(t)
	start := startInterviewSession(t, svc, nil)

	if err := svc.RecordBlockedHintAttempt(start.SessionID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// A blocked attempt is a real, durable event like any other, so it
	// advances the stream's revision; a caller still holding the stale
	// start.Revision correctly hits the same retryable conflict any
	// concurrent write produces (errors.go maps it to a retryable
	// STATE_CONFLICT: "fetch it again"), never a silent, invisible
	// mutation.
	if _, err := svc.GranularityAdjust(start.SessionID, learning.DepthMicro, "", start.Revision, ""); !errors.Is(err, eventstore.ErrRevisionConflict) {
		t.Fatalf("err = %v, want ErrRevisionConflict", err)
	}

	fresh := svc.store.Revision(string(start.SessionID))
	if _, err := svc.GranularityAdjust(start.SessionID, learning.DepthMicro, "", fresh, ""); err != nil {
		t.Fatalf("re-fetching the revision should let the call succeed: %v", err)
	}
}

func TestFinishWithReasonPersistsReason(t *testing.T) {
	svc := newTestService(t)
	start := startInterviewSession(t, svc, nil)

	if _, err := svc.FinishWithReason(start.SessionID, "timeout", start.Revision, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	report, err := svc.InterviewReport(start.SessionID, time.Now(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.FinishReason != "timeout" {
		t.Fatalf("FinishReason = %q, want timeout", report.FinishReason)
	}
}

func TestFinishStillDefaultsToEmptyReason(t *testing.T) {
	svc := newTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := svc.Finish(start.SessionID, start.Revision, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	get, err := svc.Get(start.SessionID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if get.State != learning.SessionStateCompleted {
		t.Fatalf("state = %s, want completed", get.State)
	}
}

func TestInterviewReportAssemblesFromRecordedEvidence(t *testing.T) {
	svc := newTestService(t)
	start := startInterviewSession(t, svc, nil)

	if err := svc.RecordBlockedHintAttempt(start.SessionID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rev := svc.store.Revision(string(start.SessionID))
	if _, err := svc.ReflectionRecord(start.SessionID, "", "why?", "because", "", rev, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	report, err := svc.InterviewReport(start.SessionID, time.Now(), []string{"error-handling"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.ChallengeID != fixtureChallengeID {
		t.Fatalf("ChallengeID = %s", report.ChallengeID)
	}
	if report.Hints.Blocked != 1 {
		t.Fatalf("Hints.Blocked = %d, want 1", report.Hints.Blocked)
	}
	if len(report.Reflections) != 1 || report.Reflections[0].Answer != "because" {
		t.Fatalf("Reflections = %+v", report.Reflections)
	}
	if len(report.Recommended) != 1 || report.Recommended[0] != "error-handling" {
		t.Fatalf("Recommended = %v", report.Recommended)
	}
	if report.IntegrityNote == "" {
		t.Fatal("expected a non-empty integrity note")
	}
}
