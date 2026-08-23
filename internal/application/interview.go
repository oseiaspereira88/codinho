package application

import (
	"time"

	"github.com/oseiaspereira88/codinho/internal/assessment"
	"github.com/oseiaspereira88/codinho/internal/learning"
	"github.com/oseiaspereira88/codinho/internal/session"
)

// InterviewStatus and InterviewReport live in internal/session (state
// replay) and internal/assessment (the pure descriptive report),
// respectively; these aliases let mcpserver depend only on application's
// stable surface (interview-mode).
type (
	InterviewStatus = session.InterviewStatus
	InterviewReport = assessment.InterviewReport
)

// InterviewStatus reports elapsed time and whether the session's time
// limit (if any) has been reached (requirement R3, R5).
func (s *SessionService) InterviewStatus(id learning.SessionID, now time.Time) (InterviewStatus, error) {
	return s.svc.InterviewStatus(id, now)
}

// RecordBlockedHintAttempt durably logs that a hint request was blocked
// by the session's help policy (requirement R4).
func (s *SessionService) RecordBlockedHintAttempt(id learning.SessionID) error {
	return s.svc.RecordBlockedHintAttempt(id)
}

// FinishWithReason transitions a session to completed and records why —
// an explicit learner action, or "timeout" (requirement R5).
func (s *SessionService) FinishWithReason(id learning.SessionID, reason string, expectedRevision uint64, requestID string) (LifecycleResult, error) {
	return s.svc.FinishWithReason(id, reason, expectedRevision, requestID)
}

// InterviewReport assembles the descriptive, non-scalar end-of-interview
// report from what was actually recorded (requirement R6, R7).
// recommended is caller-supplied (e.g. from learning_path_recommend).
func (s *SessionService) InterviewReport(id learning.SessionID, now time.Time, recommended []string) (InterviewReport, error) {
	return s.svc.InterviewReport(id, now, recommended)
}
