package application

import (
	"github.com/oseiaspereira88/codinho/internal/assistance"
	"github.com/oseiaspereira88/codinho/internal/learning"
	"github.com/oseiaspereira88/codinho/internal/session"
)

// The assistance ladder, syntax recall and detour orchestration live in
// internal/session (state) and internal/assistance (pure decisions); these
// aliases let mcpserver depend only on application's stable surface
// (assistance-hints-detours).
type (
	HintResult     = session.HintResult
	DetourResult   = session.DetourResult
	DetourOutcome  = assistance.DetourOutcome
	ConceptContent = assistance.ConceptContent
)

const (
	DetourResolved  = assistance.DetourResolved
	DetourAbandoned = assistance.DetourAbandoned
)

var (
	ErrNoActiveStep                 = session.ErrNoActiveStep
	ErrNoOpenDetour                 = session.ErrNoOpenDetour
	ErrNoHintAuthored               = assistance.ErrNoHintAuthored
	ErrHelpDisabled                 = assistance.ErrHelpDisabled
	ErrSolutionRequiresConfirmation = assistance.ErrSolutionRequiresConfirmation
	ErrConceptNotFound              = assistance.ErrConceptNotFound
)

// HintRequest grants the next disclosure rung for id's active step
// (requirement R1, R2, R3, R7).
func (s *SessionService) HintRequest(id learning.SessionID, confirmSolution bool, expectedRevision uint64, requestID string) (HintResult, error) {
	return s.svc.HintRequest(id, confirmSolution, expectedRevision, requestID)
}

// SyntaxRecallGet returns the step's authored syntax-recall rung, if any
// (requirement R5; RF-022).
func (s *SessionService) SyntaxRecallGet(id learning.SessionID, expectedRevision uint64, requestID string) (HintResult, error) {
	return s.svc.SyntaxRecallGet(id, expectedRevision, requestID)
}

// DetourStart opens a conceptual detour without changing the active step
// (requirement R6).
func (s *SessionService) DetourStart(id learning.SessionID, reason string, expectedRevision uint64, requestID string) (DetourResult, error) {
	return s.svc.DetourStart(id, reason, expectedRevision, requestID)
}

// DetourFinish closes the open detour and returns to the same active step
// (requirement R6).
func (s *SessionService) DetourFinish(id learning.SessionID, outcome DetourOutcome, expectedRevision uint64, requestID string) (DetourResult, error) {
	return s.svc.DetourFinish(id, outcome, expectedRevision, requestID)
}

// AssistanceService adapts internal/assistance for the MCP server's
// read-only concept lookup (requirement R5).
type AssistanceService struct {
	svc *assistance.Service
}

// NewAssistanceService wraps catalog's already-loaded catalog.
func NewAssistanceService(catalog *CatalogService) *AssistanceService {
	return &AssistanceService{svc: assistance.New(catalog.catalog)}
}

// ConceptContent returns the canonical record for concept id.
func (a *AssistanceService) ConceptContent(id string) (ConceptContent, error) {
	return a.svc.ConceptContent(id)
}
