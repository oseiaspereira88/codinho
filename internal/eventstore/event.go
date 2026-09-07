// Package eventstore persists an append-only, versioned event log with
// atomic snapshots, exclusive locking and crash-safe recovery
// (PROJECT.md §18). It never interprets event payloads: encoding and
// projecting them into domain state is the consumer's responsibility
// (see local-event-store Decision 3).
package eventstore

import "encoding/json"

// CurrentEventSchemaVersion is the envelope schema version this store
// writes. Reading an event with a newer version and no registered
// Upcaster is treated as incompatible (requirement R6, compatibility
// requirement).
const CurrentEventSchemaVersion = 1

// EventType names one of the minimum events of PROJECT.md §18.3. The store
// treats Type as an opaque string; these constants exist so producers and
// consumers share the same vocabulary.
type EventType string

const (
	EventSessionStarted       EventType = "session_started"
	EventSessionPolicyChanged EventType = "session_policy_changed"
	EventInstructionIssued    EventType = "instruction_issued"
	EventGranularityChanged   EventType = "granularity_changed"
	EventObservationRecorded  EventType = "observation_recorded"
	EventHintRequested        EventType = "hint_requested"
	EventDetourStarted        EventType = "detour_started"
	EventDetourFinished       EventType = "detour_finished"
	EventFeedbackRecorded     EventType = "feedback_recorded"
	EventAttemptSubmitted     EventType = "attempt_submitted"
	EventCheckExecuted        EventType = "check_executed"
	EventEvaluationRecorded   EventType = "evaluation_recorded"
	EventEvidenceRecorded     EventType = "evidence_recorded"
	EventReflectionRecorded   EventType = "reflection_recorded"
	EventStepCompleted        EventType = "step_completed"
	EventStepAdvanced         EventType = "step_advanced"
	EventSolutionRevealed     EventType = "solution_revealed"
	EventSessionPaused        EventType = "session_paused"
	// EventSessionResumed extends PROJECT.md §18.3's minimum event list
	// (session-orchestration-disclosure, requirement R5): the minimum list
	// omits an explicit resume marker, but session_paused's counterpart is
	// needed to audit lifecycle without inferring it from other events.
	EventSessionResumed   EventType = "session_resumed"
	EventSessionFinished  EventType = "session_finished"
	EventMasteryProjected EventType = "mastery_projected"
	EventReviewScheduled  EventType = "review_scheduled"
	// EventLearnerNextStepProposed records a pure autonomy signal
	// (learning-practice-debug-modes requirement R9; PROJECT.md §21.5,
	// "capacidade de propor o próximo passo"): it never advances
	// anything on its own.
	EventLearnerNextStepProposed EventType = "learner_next_step_proposed"
)

// Event is one immutable, append-only log entry (requirement R1). Revision
// is 1-based and monotonic per StreamID; RequestID, when set, makes Append
// idempotent for retries of the same logical write (requirement R4).
type Event struct {
	SchemaVersion int             `json:"schema_version"`
	ID            string          `json:"id"`
	StreamID      string          `json:"stream_id"`
	Revision      uint64          `json:"revision"`
	Type          EventType       `json:"type"`
	RequestID     string          `json:"request_id,omitempty"`
	RecordedAt    string          `json:"recorded_at"` // RFC 3339, UTC
	Payload       json.RawMessage `json:"payload"`
}
