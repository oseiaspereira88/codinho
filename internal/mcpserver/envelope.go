// Package mcpserver implements the codinho MCP server over stdio
// (PROJECT.md §15). It never calls an LLM and never edits learner files; it
// exposes a minimal vertical slice of read-mostly tools over the domain,
// curriculum and event-store packages.
package mcpserver

// ProgressEffect is one of the stable values PROJECT.md §15.3 defines for
// the envelope's progress_effect field.
type ProgressEffect string

const (
	ProgressEffectNone               ProgressEffect = "none"
	ProgressEffectSessionChanged     ProgressEffect = "session_changed"
	ProgressEffectFeedbackRecorded   ProgressEffect = "feedback_recorded"
	ProgressEffectAttemptRecorded    ProgressEffect = "attempt_recorded"
	ProgressEffectEvaluationRecorded ProgressEffect = "evaluation_recorded"
	ProgressEffectStepCompleted      ProgressEffect = "step_completed"
	ProgressEffectStepAdvanced       ProgressEffect = "step_advanced"
	ProgressEffectSolutionReveal     ProgressEffect = "solution_revealed"
	ProgressEffectMasteryProjected   ProgressEffect = "mastery_projected"
)

// ActiveNode identifies the session's current instructional node.
type ActiveNode struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
}

// Disclosure reports the session's current disclosure state.
type Disclosure struct {
	Level            string `json:"level"`
	SolutionRevealed bool   `json:"solution_revealed"`
}

// ErrorInfo is the envelope's error shape (PROJECT.md §15.4). Message is
// always safe to show a client: it never contains paths, environment
// values or internal error text (requirement R4; security).
type ErrorInfo struct {
	Code      string         `json:"code"`
	Message   string         `json:"message"`
	Retryable bool           `json:"retryable"`
	Details   map[string]any `json:"details,omitempty"`
}

// Envelope is the single response shape every tool in this server returns
// (PROJECT.md §15.3, §15.4).
type Envelope struct {
	Status         string         `json:"status"`
	RequestID      string         `json:"request_id,omitempty"`
	SessionID      string         `json:"session_id,omitempty"`
	ActiveNode     *ActiveNode    `json:"active_node,omitempty"`
	ProgressEffect ProgressEffect `json:"progress_effect"`
	Disclosure     *Disclosure    `json:"disclosure,omitempty"`
	AllowedActions []string       `json:"allowed_actions,omitempty"`
	Data           any            `json:"data,omitempty"`
	Warnings       []string       `json:"warnings,omitempty"`
	Error          *ErrorInfo     `json:"error,omitempty"`
}

// okEnvelope builds a successful response.
func okEnvelope(requestID string, effect ProgressEffect, data any) Envelope {
	return Envelope{
		Status:         "ok",
		RequestID:      requestID,
		ProgressEffect: effect,
		Data:           data,
	}
}

// errorEnvelope builds a failed response. It never sets Status to anything
// but "error", so a caller can always branch on Status alone.
func errorEnvelope(requestID string, code ErrorCode, message string, retryable bool, details map[string]any) Envelope {
	return Envelope{
		Status:         "error",
		RequestID:      requestID,
		ProgressEffect: ProgressEffectNone,
		Error: &ErrorInfo{
			Code:      string(code),
			Message:   message,
			Retryable: retryable,
			Details:   details,
		},
	}
}
