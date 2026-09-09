package mcpserver

import (
	"errors"
	"os"

	"github.com/oseiaspereira88/codinho/internal/application"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/evidence"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

// ErrorCode is one of the stable codes PROJECT.md §15.4 defines. Only the
// subset relevant to this minimal tool slice is mapped here; the remaining
// codes belong to the specs that introduce the operations that can raise
// them.
type ErrorCode string

const (
	ErrCodeInvalidInput               ErrorCode = "INVALID_INPUT"
	ErrCodeItemNotFound               ErrorCode = "ITEM_NOT_FOUND"
	ErrCodeSessionNotActive           ErrorCode = "SESSION_NOT_ACTIVE"
	ErrCodeSessionRecoveryUnavailable ErrorCode = "SESSION_RECOVERY_UNAVAILABLE"
	ErrCodeEvaluationEvidenceInvalid  ErrorCode = "EVALUATION_EVIDENCE_INVALID"
	ErrCodeEvaluationEvidenceStale    ErrorCode = "EVALUATION_EVIDENCE_STALE"
	ErrCodeStateConflict              ErrorCode = "STATE_CONFLICT"
	ErrCodeInternalError              ErrorCode = "INTERNAL_ERROR"
)

// mapError translates an internal error into a stable error code, a safe
// message and whether the caller may retry (requirement R4). The message
// never includes the original error's text, a path or an environment
// value: those could leak internals (security requirement).
func mapError(err error) (ErrorCode, string, bool) {
	switch {
	case errors.Is(err, curriculum.ErrInvalidSelection):
		return ErrCodeInvalidInput, "invalid, conflicting, stale or excessive selection", false
	case errors.Is(err, application.ErrEvaluationEvidenceInvalid):
		return ErrCodeEvaluationEvidenceInvalid, "evidence is unavailable or does not match this session, step, check or rubric", false
	case errors.Is(err, application.ErrEvaluationEvidenceStale):
		return ErrCodeEvaluationEvidenceStale, "workspace evidence is no longer current; observe and run the check again", false
	case errors.Is(err, application.ErrNotFound):
		return ErrCodeItemNotFound, "the requested item does not exist in the catalog", false
	case errors.Is(err, application.ErrSessionNotFound):
		return ErrCodeSessionNotActive, "no active session with that ID", false
	case errors.Is(err, application.ErrSessionUnrecoverable):
		return ErrCodeSessionRecoveryUnavailable, "historical session cannot be safely restored; preserve its history and start a new session", false
	case errors.Is(err, eventstore.ErrRevisionConflict):
		return ErrCodeStateConflict, "the session changed since your last read; fetch it again", true
	case errors.Is(err, application.ErrChallengeHasNoSteps):
		return ErrCodeInvalidInput, "the challenge has no authored steps to start from", false
	case errors.Is(err, application.ErrNoWindowAtDepth):
		return ErrCodeInvalidInput, "no instructional window exists at that depth for this challenge", false
	case errors.Is(err, application.ErrInvalidNextStep):
		return ErrCodeInvalidInput, "next_step_id must select an offered branch", false
	case errors.Is(err, application.ErrNavigationUnavailable):
		return ErrCodeSessionNotActive, "navigation requires an active session without an open detour", false
	case errors.Is(err, application.ErrStepNotFound):
		return ErrCodeInvalidInput, "step_id does not exist in this session's challenge", false
	case errors.Is(err, application.ErrNoActiveStep):
		return ErrCodeSessionNotActive, "the session has no active instructional step", false
	case errors.Is(err, application.ErrNoOpenDetour):
		return ErrCodeInvalidInput, "there is no open detour to finish for this session", false
	case errors.Is(err, application.ErrNoHintAuthored):
		return ErrCodeInvalidInput, "this step has no authored syntax-recall hint", false
	case errors.Is(err, application.ErrHelpDisabled):
		return ErrCodeInvalidInput, "the session's help policy disallows hints", false
	case errors.Is(err, application.ErrSolutionRequiresConfirmation):
		return ErrCodeInvalidInput, "revealing the solution requires explicit confirmation", false
	case errors.Is(err, application.ErrConceptNotFound):
		return ErrCodeItemNotFound, "the requested concept does not exist in the catalog", false
	case errors.Is(err, application.ErrWorkspaceRootInvalid):
		return ErrCodeInvalidInput, "the workspace root is not a valid, accessible directory", false
	case errors.Is(err, application.ErrEvidenceOutOfScope):
		return ErrCodeItemNotFound, "no evidence with that ID exists for this session", false
	case errors.Is(err, application.ErrCheckNotFound):
		return ErrCodeItemNotFound, "no check with that ID is declared by this session's challenge", false
	case errors.Is(err, application.ErrNoWorkspaceBaseline):
		return ErrCodeInvalidInput, "call workspace_observe for the active step before running a check", false
	case errors.Is(err, application.ErrInvalidDimension):
		return ErrCodeInvalidInput, "dimension must be one of the eight declared mastery dimensions", false
	case errors.Is(err, application.ErrInvalidCompetency):
		return ErrCodeInvalidInput, "competency_id and evidence_id are required", false
	case errors.Is(err, curriculum.ErrSearchTextTooLong):
		return ErrCodeInvalidInput, "search text is too long", false
	case errors.Is(err, evidence.ErrInvalidID), errors.Is(err, evidence.ErrIntegrityMismatch), errors.Is(err, os.ErrNotExist):
		return ErrCodeItemNotFound, "no evidence with that ID exists for this session", false
	default:
		var domainErr learning.DomainError
		if errors.As(err, &domainErr) {
			return ErrCodeInvalidInput, "the request violates a domain rule: " + string(domainErr.Code), false
		}
		return ErrCodeInternalError, "an internal error occurred", true
	}
}
