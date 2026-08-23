package learning

import "fmt"

// ErrorCode identifies a domain failure independently of its message, so
// callers can compare and branch on it without depending on UI-facing text
// (PROJECT.md §2 requirement R7).
type ErrorCode string

const (
	ErrCodeInvalidValue                        ErrorCode = "invalid_value"
	ErrCodeInstructionAlreadyActive            ErrorCode = "instruction_already_active"
	ErrCodeNoActiveInstruction                 ErrorCode = "no_active_instruction"
	ErrCodeInvalidSessionTransition            ErrorCode = "invalid_session_transition"
	ErrCodeInvalidStepTransition               ErrorCode = "invalid_step_transition"
	ErrCodeInvalidDetourTransition             ErrorCode = "invalid_detour_transition"
	ErrCodeStepNotCompletable                  ErrorCode = "step_not_completable"
	ErrCodeAdvanceNotAllowed                   ErrorCode = "advance_not_allowed"
	ErrCodeDisclosureExceeded                  ErrorCode = "disclosure_exceeded"
	ErrCodeHintLevelSkipped                    ErrorCode = "hint_level_skipped"
	ErrCodeQualitativeJudgmentRequiresEvidence ErrorCode = "qualitative_judgment_requires_evidence"
	ErrCodeStepNotEvaluable                    ErrorCode = "step_not_evaluable"
	ErrCodeCompletionPolicyNotMet              ErrorCode = "completion_policy_not_met"
)

// DomainError is a typed, comparable error. Two DomainError values are equal
// (and errors.Is succeeds) whenever their Code matches, regardless of the
// human-readable Detail.
type DomainError struct {
	Code   ErrorCode
	Detail string
}

func (e DomainError) Error() string {
	if e.Detail == "" {
		return string(e.Code)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Detail)
}

// Is implements errors.Is comparison by Code alone, so a caller can test
// against a zero-Detail sentinel: errors.Is(err, DomainError{Code: ...}).
func (e DomainError) Is(target error) bool {
	other, ok := target.(DomainError)
	return ok && other.Code == e.Code
}

func newDomainError(code ErrorCode, detail string) DomainError {
	return DomainError{Code: code, Detail: detail}
}
