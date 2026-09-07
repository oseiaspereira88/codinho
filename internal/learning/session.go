package learning

// SessionID identifies a learning session with a stable ID.
type SessionID CatalogID

// SessionState is the lifecycle of a LearningSession (PROJECT.md §13.1).
type SessionState string

const (
	SessionStateCreated   SessionState = "created"
	SessionStateActive    SessionState = "active"
	SessionStatePaused    SessionState = "paused"
	SessionStateCompleted SessionState = "completed"
	SessionStateAbandoned SessionState = "abandoned"
)

// sessionTransitions enumerates every edge of PROJECT.md §13.1:
//
//	created → active → completed
//	                 ↘ abandoned
//	active ↔ paused
//
// A completed or abandoned session is immutable: it has no outgoing edges.
var sessionTransitions = map[SessionState]map[SessionState]bool{
	SessionStateCreated:   {SessionStateActive: true},
	SessionStateActive:    {SessionStateCompleted: true, SessionStateAbandoned: true, SessionStatePaused: true},
	SessionStatePaused:    {SessionStateActive: true},
	SessionStateCompleted: {},
	SessionStateAbandoned: {},
}

// ValidSessionTransition reports whether moving a session from one state to
// another is allowed.
func ValidSessionTransition(from, to SessionState) bool {
	return sessionTransitions[from][to]
}

// LearningSession fixes a profile, content and policies; controls the
// active node, disclosure and detours; and holds exactly one active
// progression line (PROJECT.md §12.1).
type LearningSession struct {
	ID      SessionID
	State   SessionState
	Policy  SessionPolicy
	Catalog CatalogRef
	active  *StepProgress
	detour  *LearningDetour
}

// Clone isolates mutable state for validation before a durable append.
func (s *LearningSession) Clone() *LearningSession {
	c := *s
	if s.active != nil {
		active := *s.active
		c.active = &active
	}
	if s.detour != nil {
		detour := *s.detour
		c.detour = &detour
	}
	if s.Policy.TimeLimit != nil {
		limit := *s.Policy.TimeLimit
		c.Policy.TimeLimit = &limit
	}
	return &c
}

// CatalogRef pins the exact catalog version a session fixed at start
// (PROJECT.md §12.3 invariant 1).
type CatalogRef struct {
	TrackID TrackID
	Version ContentVersion
}

// NewLearningSession starts a session in SessionStateCreated with no active
// instruction.
func NewLearningSession(id SessionID, policy SessionPolicy, catalog CatalogRef) *LearningSession {
	return &LearningSession{ID: id, State: SessionStateCreated, Policy: policy, Catalog: catalog}
}

// Transition moves the session to a new state.
func (s *LearningSession) Transition(to SessionState) error {
	if !ValidSessionTransition(s.State, to) {
		return newDomainError(ErrCodeInvalidSessionTransition, string(s.State)+" -> "+string(to))
	}
	s.State = to
	return nil
}

// ActiveStep returns the currently active step progress, or nil when no
// instruction is active.
func (s *LearningSession) ActiveStep() *StepProgress {
	return s.active
}

// SetActiveInstruction sets the session's single active instruction. It
// fails when one is already active, enforcing that a session has at most
// one active instruction at a time (PROJECT.md §12.3 invariant 2,
// requirement R4).
func (s *LearningSession) SetActiveInstruction(step *StepProgress) error {
	if s.active != nil {
		return newDomainError(ErrCodeInstructionAlreadyActive, string(s.active.StepID))
	}
	s.active = step
	return nil
}

// ClearActiveInstruction releases the active instruction so another one may
// be set, typically after the active step completes or is skipped.
func (s *LearningSession) ClearActiveInstruction() {
	s.active = nil
}

// Advance activates next as the session's new active instruction. It
// requires the current active step to be completed (or an explicit
// override), and never runs implicitly from evaluation or completion
// (PROJECT.md §8.6 "Avanço", §12.3 invariant 6, requirement R4).
func (s *LearningSession) Advance(next *StepProgress, override bool) error {
	if s.active != nil && s.active.State != StepStateCompleted && !override {
		return newDomainError(ErrCodeAdvanceNotAllowed, string(s.active.StepID))
	}
	s.active = nil
	return s.SetActiveInstruction(next)
}

// Detour returns the session's open learning detour, or nil.
func (s *LearningSession) Detour() *LearningDetour {
	return s.detour
}

// OpenDetour opens a conceptual detour without changing the active
// instruction: the original node stays active throughout (PROJECT.md §8.8,
// §13.3).
func (s *LearningSession) OpenDetour(reason string) (*LearningDetour, error) {
	if s.active == nil {
		return nil, newDomainError(ErrCodeNoActiveInstruction, "detour requires an active step to return to")
	}
	if s.detour != nil && s.detour.State != DetourStateReturned {
		return nil, newDomainError(ErrCodeInvalidDetourTransition, "a detour is already open")
	}
	d := &LearningDetour{Reason: reason, State: DetourStateOpened}
	s.detour = d
	return d, nil
}

// DetourState is the lifecycle of a LearningDetour (PROJECT.md §13.3).
type DetourState string

const (
	DetourStateOpened    DetourState = "opened"
	DetourStateActive    DetourState = "active"
	DetourStateResolved  DetourState = "resolved"
	DetourStateAbandoned DetourState = "abandoned"
	DetourStateReturned  DetourState = "returned"
)

var detourTransitions = map[DetourState]map[DetourState]bool{
	DetourStateOpened:    {DetourStateActive: true},
	DetourStateActive:    {DetourStateResolved: true, DetourStateAbandoned: true},
	DetourStateResolved:  {DetourStateReturned: true},
	DetourStateAbandoned: {DetourStateReturned: true},
	DetourStateReturned:  {},
}

// LearningDetour is a temporary conceptual detour that returns to the same
// active step without affecting its progress.
type LearningDetour struct {
	Reason string
	State  DetourState
}

// Transition moves the detour to a new state.
func (d *LearningDetour) Transition(to DetourState) error {
	if !detourTransitions[d.State][to] {
		return newDomainError(ErrCodeInvalidDetourTransition, string(d.State)+" -> "+string(to))
	}
	d.State = to
	return nil
}
