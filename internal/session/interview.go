package session

import (
	"encoding/json"
	"time"

	"github.com/oseiaspereira88/codinho/internal/assessment"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

// InterviewElapsed reports how much wall-clock time has passed since
// startedAt and whether that exceeds limit. now is supplied by the
// caller rather than read internally (the MCP tool layer passes
// time.Now(); a test passes a fixed value), the same pattern
// ProgressService.ReviewDue(now time.Time) already establishes, so this
// logic is deterministic and testable without touching the system clock
// (non-functional requirement: "Clock deve ser injetável e
// determinístico em testes").
func InterviewElapsed(startedAt, now time.Time, limit *time.Duration) (elapsed time.Duration, timedOut bool) {
	elapsed = now.Sub(startedAt)
	if limit == nil {
		return elapsed, false
	}
	return elapsed, elapsed >= *limit
}

// InterviewStatus is what Service.InterviewStatus returns.
type InterviewStatus struct {
	StartedAt time.Time
	Elapsed   time.Duration
	TimeLimit *time.Duration
	TimedOut  bool
}

// InterviewStatus reports elapsed time and whether the session's
// TimeLimit (if any) has been reached, using the session's own real,
// tamper-evident session_started event as the start time — this server
// never invents or manipulates that timestamp (requirement R3: "sem
// manipular relógio do sistema"). Applicable to any session with a
// TimeLimit set, not only interview mode.
func (s *Service) InterviewStatus(id learning.SessionID, now time.Time) (InterviewStatus, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rec, ok := s.sessions[id]
	if !ok {
		return InterviewStatus{}, ErrSessionNotFound
	}
	events := s.store.Replay(string(id))
	if len(events) == 0 {
		return InterviewStatus{}, ErrSessionNotFound
	}
	startedAt, err := time.Parse(time.RFC3339Nano, events[0].RecordedAt)
	if err != nil {
		return InterviewStatus{}, err
	}
	elapsed, timedOut := InterviewElapsed(startedAt, now, rec.session.Policy.TimeLimit)
	return InterviewStatus{
		StartedAt: startedAt, Elapsed: elapsed, TimeLimit: rec.session.Policy.TimeLimit, TimedOut: timedOut,
	}, nil
}

// RecordBlockedHintAttempt durably logs that a hint request was blocked
// by the session's help policy, without revealing anything about the
// step (PROJECT.md §9.5: "registra pedidos de ajuda bloqueados ou
// autorizados"; requirement R4). It takes no expected_revision from the
// caller — HintRequest already returned the "help disabled" error before
// this runs, so the caller has no revision to supply — but it still
// advances the stream's real revision like any other durable event; a
// caller holding an older revision hits the same retryable
// eventstore.ErrRevisionConflict any concurrent write would produce, and
// resolves it the same way (re-fetch via session_get).
func (s *Service) RecordBlockedHintAttempt(id learning.SessionID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.sessions[id]; !ok {
		return ErrSessionNotFound
	}
	current := s.store.Revision(string(id))
	_, err := s.store.Append(string(id), current, "", eventstore.EventHintRequested, map[string]any{"blocked": true})
	return err
}

// FinishWithReason transitions to completed and records why (an explicit
// learner action, or "timeout"), so an interview report can distinguish
// them (requirement R5). Finish calls this with an empty reason, so every
// existing caller keeps its exact prior behavior.
func (s *Service) FinishWithReason(id learning.SessionID, reason string, expectedRevision uint64, requestID string) (LifecycleResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rec, ok := s.sessions[id]
	if !ok {
		return LifecycleResult{}, ErrSessionNotFound
	}
	ev, fresh, err := s.mutate(id, expectedRevision, requestID, eventstore.EventSessionFinished, map[string]string{"reason": reason})
	if err != nil {
		return LifecycleResult{}, err
	}
	if fresh {
		if err := rec.session.Transition(learning.SessionStateCompleted); err != nil {
			return LifecycleResult{}, err
		}
	}
	return LifecycleResult{State: rec.session.State, Revision: ev.Revision}, nil
}

// InterviewReport replays id's full event stream and assembles a
// descriptive, non-scalar report (requirement R6, R7) from what was
// actually recorded: evaluations, granted-vs-blocked hints and
// reflections. now is the caller-supplied wall clock (requirement R3);
// recommended is a caller-supplied list of follow-up step/competency IDs
// (e.g. from learning_path_recommend) — this method never invents one.
func (s *Service) InterviewReport(id learning.SessionID, now time.Time, recommended []string) (assessment.InterviewReport, error) {
	s.mu.Lock()
	rec, ok := s.sessions[id]
	if !ok {
		s.mu.Unlock()
		return assessment.InterviewReport{}, ErrSessionNotFound
	}
	challengeID := rec.challengeID
	timeLimit := rec.session.Policy.TimeLimit
	s.mu.Unlock()

	events := s.store.Replay(string(id))
	if len(events) == 0 {
		return assessment.InterviewReport{}, ErrSessionNotFound
	}
	startedAt, err := time.Parse(time.RFC3339Nano, events[0].RecordedAt)
	if err != nil {
		return assessment.InterviewReport{}, err
	}
	elapsed, timedOut := InterviewElapsed(startedAt, now, timeLimit)

	var (
		evaluations  []assessment.InterviewEvaluation
		reflections  []assessment.InterviewReflection
		hints        assessment.InterviewHintUsage
		finishReason string
	)
	for _, ev := range events {
		switch ev.Type {
		case eventstore.EventEvaluationRecorded:
			var payload struct {
				StepID   string                     `json:"step_id"`
				Criteria []learning.CriterionResult `json:"criteria"`
			}
			if err := json.Unmarshal(ev.Payload, &payload); err != nil {
				return assessment.InterviewReport{}, err
			}
			evaluations = append(evaluations, assessment.InterviewEvaluation{StepID: payload.StepID, Criteria: payload.Criteria})
		case eventstore.EventHintRequested:
			var payload struct {
				Blocked bool `json:"blocked"`
			}
			if err := json.Unmarshal(ev.Payload, &payload); err != nil {
				return assessment.InterviewReport{}, err
			}
			if payload.Blocked {
				hints.Blocked++
			} else {
				hints.Granted++
			}
		case eventstore.EventReflectionRecorded:
			var payload struct {
				StepID string `json:"step_id"`
				Prompt string `json:"prompt"`
				Answer string `json:"answer"`
			}
			if err := json.Unmarshal(ev.Payload, &payload); err != nil {
				return assessment.InterviewReport{}, err
			}
			reflections = append(reflections, assessment.InterviewReflection(payload))
		case eventstore.EventSessionFinished:
			var payload struct {
				Reason string `json:"reason"`
			}
			if err := json.Unmarshal(ev.Payload, &payload); err == nil {
				finishReason = payload.Reason
			}
		}
	}

	return assessment.BuildInterviewReport(challengeID, elapsed, timedOut, finishReason, evaluations, hints, reflections, recommended), nil
}
