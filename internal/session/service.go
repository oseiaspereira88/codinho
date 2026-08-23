// Package session orchestrates a learning session's lifecycle, single
// active instruction and progressive disclosure over the pure domain
// (internal/learning), authored content (internal/curriculum) and the
// durable event log (internal/eventstore). It never generates feedback,
// runs checks or computes mastery (non-goals of this spec).
package session

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/oseiaspereira88/codinho/internal/assistance"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

// ErrNotFound is returned when a challenge lookup finds nothing with the
// given ID.
var ErrNotFound = errors.New("session: not found")

// ErrSessionNotFound is returned when a session lookup finds nothing with
// the given ID.
var ErrSessionNotFound = errors.New("session: session not found")

// ErrChallengeHasNoSteps is returned when a challenge exists but authored
// no macro step to start from.
var ErrChallengeHasNoSteps = errors.New("session: challenge has no steps")

// ErrNoWindowAtDepth is returned when granularity_adjust cannot resolve any
// node at the requested depth (e.g. the challenge authored no layers).
var ErrNoWindowAtDepth = errors.New("session: no window available at that depth")

// ErrNoActiveStep is returned when a session exists but currently has no
// active instructional step to grant a hint against.
var ErrNoActiveStep = errors.New("session: no active instructional step")

// ErrNoOpenDetour is returned when DetourFinish is called and the session
// has never opened a detour.
var ErrNoOpenDetour = errors.New("session: no open detour for this session")

type record struct {
	session     *learning.LearningSession
	challengeID string
	depth       learning.Depth
}

// Service orchestrates sessions: lifecycle, single active instruction,
// disclosure and granularity. It holds a session registry in memory keyed
// by SessionID; durability comes from appending every state-changing
// operation to the event store before applying it in memory.
type Service struct {
	catalog *curriculum.Catalog
	store   *eventstore.Store

	mu            sync.Mutex
	sessions      map[learning.SessionID]*record
	startResults  map[string]StartResult
	hintResults   map[string]HintResult
	detourResults map[string]DetourResult
	nextID        atomic.Uint64
}

// idempotencyKey scopes a request_id to its session, matching the
// eventstore's own (streamID, requestID) scoping, so two sessions reusing
// the same client-chosen request_id never collide.
func idempotencyKey(id learning.SessionID, requestID string) string {
	return string(id) + "\x00" + requestID
}

// New wires a Service to its catalog and event store.
func New(catalog *curriculum.Catalog, store *eventstore.Store) *Service {
	return &Service{
		catalog:       catalog,
		store:         store,
		sessions:      map[learning.SessionID]*record{},
		startResults:  map[string]StartResult{},
		hintResults:   map[string]HintResult{},
		detourResults: map[string]DetourResult{},
	}
}

func (s *Service) challenge(id string) (curriculum.ChallengeAuthoring, error) {
	ch, ok := s.catalog.Challenge(id)
	if !ok {
		return curriculum.ChallengeAuthoring{}, ErrNotFound
	}
	return ch, nil
}

// StartInput is what a caller supplies to Start (requirement R1).
type StartInput struct {
	ChallengeID   string
	Mode          learning.PedagogicalMode
	Depth         learning.Depth
	Help          learning.HelpPolicyKind
	DisclosureMax learning.DisclosureLevel
	Evaluation    learning.EvaluationPolicyKind
	RequestID     string // idempotency key for retries (requirement R7, R8)
}

// StartResult is what Start returns on success.
type StartResult struct {
	SessionID  learning.SessionID
	ActiveStep learning.StepID
	Objective  string
	Revision   uint64
	Disclosure Disclosure
}

// Start fixes challengeID, creates a new session at its first authored
// macro step, and durably records the session start. Calling Start again
// with the same RequestID returns the same result without creating a
// second session (requirement R4, R7, R8: idempotent retry).
func (s *Service) Start(in StartInput) (StartResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if in.RequestID != "" {
		if cached, ok := s.startResults[in.RequestID]; ok {
			return cached, nil
		}
	}

	challenge, err := s.challenge(in.ChallengeID)
	if err != nil {
		return StartResult{}, err
	}
	step, ok := firstStep(challenge)
	if !ok {
		return StartResult{}, ErrChallengeHasNoSteps
	}

	mode := in.Mode
	if mode == "" {
		mode = learning.ModePractice
	}
	depth := in.Depth
	if depth == "" {
		depth = learning.DepthMicro
	}
	help := in.Help
	if help == "" {
		help = learning.HelpProgressive
	}
	disclosureMax := in.DisclosureMax
	if disclosureMax == 0 {
		disclosureMax = learning.DisclosureGuidingQuestion
	}
	evaluation := in.Evaluation
	if evaluation == "" {
		evaluation = learning.EvaluationOnDemand
	}
	policy, err := learning.NewSessionPolicy(mode, depth, help, disclosureMax, evaluation, learning.AdvanceExplicit, nil)
	if err != nil {
		return StartResult{}, err
	}

	id := s.nextID.Add(1)
	sessionID := learning.SessionID(fmt.Sprintf("ses_%d", id))

	// Append first: the durable record of intent to start must exist
	// before the in-memory session is exposed to callers (local-event-store
	// R1/R3 ordering contract).
	ev, err := s.store.Append(string(sessionID), 0, in.RequestID, eventstore.EventSessionStarted, map[string]string{
		"challenge_id": in.ChallengeID,
		"mode":         string(mode),
	})
	if err != nil {
		return StartResult{}, err
	}

	domainSession := learning.NewLearningSession(sessionID, policy, learning.CatalogRef{})
	if err := domainSession.Transition(learning.SessionStateActive); err != nil {
		return StartResult{}, err
	}
	progress := learning.NewStepProgress(learning.StepID(step.ID))
	if err := progress.Transition(learning.StepStateAvailable, false); err != nil {
		return StartResult{}, err
	}
	if err := progress.Transition(learning.StepStateActive, false); err != nil {
		return StartResult{}, err
	}
	if err := domainSession.SetActiveInstruction(progress); err != nil {
		return StartResult{}, err
	}

	s.sessions[sessionID] = &record{session: domainSession, challengeID: in.ChallengeID, depth: depth}

	result := StartResult{
		SessionID:  sessionID,
		ActiveStep: progress.StepID,
		Objective:  step.Instruction.Objective,
		Revision:   ev.Revision,
		Disclosure: disclosureFor(policy.Disclosure, progress),
	}
	if in.RequestID != "" {
		s.startResults[in.RequestID] = result
	}
	return result, nil
}

// GetResult is what Get returns (requirement R2).
type GetResult struct {
	SessionID  learning.SessionID
	State      learning.SessionState
	ActiveStep learning.StepID
	Revision   uint64
	Disclosure Disclosure
}

// Get returns the current state of an in-memory session, or
// ErrSessionNotFound.
func (s *Service) Get(id learning.SessionID) (GetResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.getLocked(id)
}

func (s *Service) getLocked(id learning.SessionID) (GetResult, error) {
	rec, ok := s.sessions[id]
	if !ok {
		return GetResult{}, ErrSessionNotFound
	}
	active := rec.session.ActiveStep()
	var stepID learning.StepID
	if active != nil {
		stepID = active.StepID
	}
	return GetResult{
		SessionID:  id,
		State:      rec.session.State,
		ActiveStep: stepID,
		Revision:   s.store.Revision(string(id)),
		Disclosure: disclosureFor(rec.session.Policy.Disclosure, active),
	}, nil
}

// Instruction is the single disclosed instruction for a session's active
// step: objective and scope only, never constraints, children or a
// solution (requirement R3; PROJECT.md §15.6 "instruction_get").
type Instruction struct {
	StepID     learning.StepID
	Objective  string
	Scope      string
	Disclosure Disclosure
}

// Instruction returns the active step's instruction for session id.
func (s *Service) Instruction(id learning.SessionID) (Instruction, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rec, ok := s.sessions[id]
	if !ok {
		return Instruction{}, ErrSessionNotFound
	}
	active := rec.session.ActiveStep()
	if active == nil {
		return Instruction{}, ErrSessionNotFound
	}
	challenge, err := s.challenge(rec.challengeID)
	if err != nil {
		return Instruction{}, err
	}
	step, ok := findStep(challenge, string(active.StepID))
	if !ok {
		return Instruction{}, ErrNotFound
	}
	return Instruction{
		StepID:     active.StepID,
		Objective:  step.Instruction.Objective,
		Scope:      step.Instruction.Scope,
		Disclosure: disclosureFor(rec.session.Policy.Disclosure, active),
	}, nil
}

// mutate appends eventType/payload to id's stream, guarding it with
// expectedRevision and requestID (requirement R7, R8), and reports whether
// the append was genuinely new (fresh == true) or an idempotent replay of
// an earlier call (fresh == false). Callers must skip re-applying their
// domain-state change when fresh is false, or a retry would be rejected by
// the very state machine it is supposed to be idempotent against.
//
// A retry necessarily repeats the same expectedRevision the client already
// had (it never saw the first response), so comparing the returned event's
// Revision to expectedRevision+1 cannot tell a fresh write from a replay:
// both satisfy that equation. Only the stream's revision *before* this call
// distinguishes them — it still equals expectedRevision on a truly fresh
// call, and has already moved past it once the first call succeeded.
func (s *Service) mutate(id learning.SessionID, expectedRevision uint64, requestID string, eventType eventstore.EventType, payload any) (ev eventstore.Event, fresh bool, err error) {
	before := s.store.Revision(string(id))
	ev, err = s.store.Append(string(id), expectedRevision, requestID, eventType, payload)
	if err != nil {
		return ev, false, err
	}
	return ev, before == expectedRevision, nil
}

// ConfigureInput changes only the mutable session properties PROJECT.md
// §15.6 allows: help policy, evaluation policy and time limit. Granularity
// has its own tool, GranularityAdjust (requirement R4, R6).
type ConfigureInput struct {
	SessionID        learning.SessionID
	Help             *learning.HelpPolicyKind
	Evaluation       *learning.EvaluationPolicyKind
	ExpectedRevision uint64
	RequestID        string
}

// LifecycleResult is what Configure, Pause, Resume and Finish return.
type LifecycleResult struct {
	State    learning.SessionState
	Revision uint64
}

// Configure applies only the fields set in in, validating the resulting
// policy as a whole before committing it (requirement R4).
func (s *Service) Configure(in ConfigureInput) (LifecycleResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rec, ok := s.sessions[in.SessionID]
	if !ok {
		return LifecycleResult{}, ErrSessionNotFound
	}

	next := rec.session.Policy
	changed := map[string]string{}
	if in.Help != nil {
		next.Help = learning.HelpPolicy{Kind: *in.Help}
		changed["help"] = string(*in.Help)
	}
	if in.Evaluation != nil {
		next.Evaluation = *in.Evaluation
		changed["evaluation"] = string(*in.Evaluation)
	}
	validated, err := learning.NewSessionPolicy(next.Mode, next.InitialDepth, next.Help.Kind, next.Disclosure.MaxLevel, next.Evaluation, next.Advance, next.TimeLimit)
	if err != nil {
		return LifecycleResult{}, err
	}

	ev, fresh, err := s.mutate(in.SessionID, in.ExpectedRevision, in.RequestID, eventstore.EventSessionPolicyChanged, changed)
	if err != nil {
		return LifecycleResult{}, err
	}
	if fresh {
		rec.session.Policy = validated
	}
	return LifecycleResult{State: rec.session.State, Revision: ev.Revision}, nil
}

// Pause transitions a session to paused without inferring anything about
// step completion (requirement R5).
func (s *Service) Pause(id learning.SessionID, expectedRevision uint64, requestID string) (LifecycleResult, error) {
	return s.transitionLifecycle(id, expectedRevision, requestID, eventstore.EventSessionPaused, learning.SessionStatePaused)
}

// Resume transitions a paused session back to active.
func (s *Service) Resume(id learning.SessionID, expectedRevision uint64, requestID string) (LifecycleResult, error) {
	return s.transitionLifecycle(id, expectedRevision, requestID, eventstore.EventSessionResumed, learning.SessionStateActive)
}

// Finish transitions an active session to completed without inferring
// completion from step state (requirement R5).
func (s *Service) Finish(id learning.SessionID, expectedRevision uint64, requestID string) (LifecycleResult, error) {
	return s.transitionLifecycle(id, expectedRevision, requestID, eventstore.EventSessionFinished, learning.SessionStateCompleted)
}

func (s *Service) transitionLifecycle(id learning.SessionID, expectedRevision uint64, requestID string, eventType eventstore.EventType, to learning.SessionState) (LifecycleResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rec, ok := s.sessions[id]
	if !ok {
		return LifecycleResult{}, ErrSessionNotFound
	}

	ev, fresh, err := s.mutate(id, expectedRevision, requestID, eventType, nil)
	if err != nil {
		return LifecycleResult{}, err
	}
	if fresh {
		if err := rec.session.Transition(to); err != nil {
			return LifecycleResult{}, err
		}
	}
	return LifecycleResult{State: rec.session.State, Revision: ev.Revision}, nil
}

// GranularityResult is what GranularityAdjust returns.
type GranularityResult struct {
	StepID   string
	Kind     string
	Revision uint64
}

// GranularityAdjust moves the session's instructional window to depth by
// walking the canonical, unmodified step tree (requirement R6; Decision 1).
func (s *Service) GranularityAdjust(id learning.SessionID, depth learning.Depth, expectedRevision uint64, requestID string) (GranularityResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rec, ok := s.sessions[id]
	if !ok {
		return GranularityResult{}, ErrSessionNotFound
	}
	challenge, err := s.challenge(rec.challengeID)
	if err != nil {
		return GranularityResult{}, err
	}
	window, ok := deriveWindow(challenge, depth)
	if !ok {
		return GranularityResult{}, ErrNoWindowAtDepth
	}

	ev, fresh, err := s.mutate(id, expectedRevision, requestID, eventstore.EventGranularityChanged, map[string]string{
		"depth":   string(depth),
		"step_id": window.StepID,
	})
	if err != nil {
		return GranularityResult{}, err
	}
	if fresh {
		next := learning.NewStepProgress(learning.StepID(window.StepID))
		if err := next.Transition(learning.StepStateAvailable, false); err != nil {
			return GranularityResult{}, err
		}
		if err := next.Transition(learning.StepStateActive, false); err != nil {
			return GranularityResult{}, err
		}
		rec.session.ClearActiveInstruction()
		if err := rec.session.SetActiveInstruction(next); err != nil {
			return GranularityResult{}, err
		}
		rec.depth = depth
	}
	return GranularityResult{StepID: window.StepID, Kind: window.Kind, Revision: ev.Revision}, nil
}

// HintResult is what HintRequest and SyntaxRecallGet return: enough
// structured context (level, kind, objective, scope, concepts) for the
// calling tutor agent to author the actual pista text within the
// authorized envelope, without this server writing prose itself
// (non-goal: "redigir explicações abertas dentro do MCP").
type HintResult struct {
	StepID    learning.StepID
	Level     learning.DisclosureLevel
	Kind      string
	Objective string
	Scope     string
	Concepts  []string
	Revision  uint64
}

// activeWindow resolves id's active step and its authored content. Callers
// must already hold s.mu.
func (s *Service) activeWindow(id learning.SessionID) (*record, *learning.StepProgress, curriculum.StepAuthoring, error) {
	rec, ok := s.sessions[id]
	if !ok {
		return nil, nil, curriculum.StepAuthoring{}, ErrSessionNotFound
	}
	active := rec.session.ActiveStep()
	if active == nil {
		return nil, nil, curriculum.StepAuthoring{}, ErrNoActiveStep
	}
	challenge, err := s.challenge(rec.challengeID)
	if err != nil {
		return nil, nil, curriculum.StepAuthoring{}, err
	}
	step, ok := findStep(challenge, string(active.StepID))
	if !ok {
		return nil, nil, curriculum.StepAuthoring{}, ErrNotFound
	}
	return rec, active, step, nil
}

// HintRequest grants the next disclosure rung for id's active step,
// climbing the ladder by at most one level per call and enforcing the
// session's disclosure cap and the explicit confirmation reaching the
// solution requires (requirement R1, R2, R3, R7; PROJECT.md §8.4).
func (s *Service) HintRequest(id learning.SessionID, confirmSolution bool, expectedRevision uint64, requestID string) (HintResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if requestID != "" {
		if cached, ok := s.hintResults[idempotencyKey(id, requestID)]; ok {
			return cached, nil
		}
	}

	rec, active, step, err := s.activeWindow(id)
	if err != nil {
		return HintResult{}, err
	}
	next, kind, err := assistance.NextHint(step, active.HintLevel, rec.session.Policy.Help.Kind, confirmSolution)
	if err != nil {
		return HintResult{}, err
	}
	result, err := s.grantDisclosure(id, rec, active, step, next, kind, false, false, expectedRevision, requestID)
	if err != nil {
		return HintResult{}, err
	}
	if requestID != "" {
		s.hintResults[idempotencyKey(id, requestID)] = result
	}
	return result, nil
}

// SyntaxRecallGet returns the step's authored syntax-recall rung, if any
// (RF-022). In teaching mode it is free: it does not climb the ladder or
// count as a consumed hint, only recording usage for audit (PROJECT.md
// §8.4 "custo menor no modo ensino"); in every other mode it consumes the
// ladder up to that rung like a targeted hint grant.
func (s *Service) SyntaxRecallGet(id learning.SessionID, expectedRevision uint64, requestID string) (HintResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if requestID != "" {
		if cached, ok := s.hintResults[idempotencyKey(id, requestID)]; ok {
			return cached, nil
		}
	}

	rec, active, step, err := s.activeWindow(id)
	if err != nil {
		return HintResult{}, err
	}
	if rec.session.Policy.Help.Kind == learning.HelpNoHints {
		return HintResult{}, assistance.ErrHelpDisabled
	}
	level, ok := assistance.SyntaxRecallLevel(step)
	if !ok {
		return HintResult{}, assistance.ErrNoHintAuthored
	}
	free := rec.session.Policy.Mode == learning.ModeTeaching
	result, err := s.grantDisclosure(id, rec, active, step, level, assistance.KindSyntaxRecall, true, free, expectedRevision, requestID)
	if err != nil {
		return HintResult{}, err
	}
	if requestID != "" {
		s.hintResults[idempotencyKey(id, requestID)] = result
	}
	return result, nil
}

// grantDisclosure appends the hint/solution event and, on a fresh (non-
// replay) call, applies the ladder mutation in memory. direct selects
// GrantDirect (syntax_recall_get, which targets a specific rung out of
// ladder order) over GrantHint (hint_request's one-rung-at-a-time climb).
// Callers must already hold s.mu.
func (s *Service) grantDisclosure(id learning.SessionID, rec *record, active *learning.StepProgress, step curriculum.StepAuthoring, level learning.DisclosureLevel, kind string, direct, free bool, expectedRevision uint64, requestID string) (HintResult, error) {
	eventType := eventstore.EventHintRequested
	payload := map[string]any{"step_id": string(active.StepID), "level": int(level), "kind": kind, "free": free}
	if level == learning.DisclosureSolution {
		eventType = eventstore.EventSolutionRevealed
		payload["needs_variant"] = true
	}

	ev, fresh, err := s.mutate(id, expectedRevision, requestID, eventType, payload)
	if err != nil {
		return HintResult{}, err
	}
	if fresh && !free {
		var grantErr error
		if direct {
			grantErr = active.GrantDirect(level, rec.session.Policy.Disclosure)
		} else {
			grantErr = active.GrantHint(level, rec.session.Policy.Disclosure)
		}
		if grantErr != nil {
			return HintResult{}, grantErr
		}
		if level == learning.DisclosureSolution {
			active.RevealSolution()
		}
	}
	return HintResult{
		StepID:    active.StepID,
		Level:     level,
		Kind:      kind,
		Objective: step.Instruction.Objective,
		Scope:     step.Instruction.Scope,
		Concepts:  step.Concepts,
		Revision:  ev.Revision,
	}, nil
}

// DetourResult is what DetourStart and DetourFinish return.
type DetourResult struct {
	State    learning.DetourState
	Revision uint64
}

// DetourStart opens a conceptual detour without changing the session's
// active step (requirement R6; PROJECT.md §8.8).
func (s *Service) DetourStart(id learning.SessionID, reason string, expectedRevision uint64, requestID string) (DetourResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rec, ok := s.sessions[id]
	if !ok {
		return DetourResult{}, ErrSessionNotFound
	}
	ev, fresh, err := s.mutate(id, expectedRevision, requestID, eventstore.EventDetourStarted, map[string]string{"reason": reason})
	if err != nil {
		return DetourResult{}, err
	}
	if fresh {
		if _, err := assistance.StartDetour(rec.session, reason); err != nil {
			return DetourResult{}, err
		}
	}
	d := rec.session.Detour()
	if d == nil {
		return DetourResult{}, ErrNoOpenDetour
	}
	return DetourResult{State: d.State, Revision: ev.Revision}, nil
}

// DetourFinish closes the session's open detour with outcome ("resolved"
// or "abandoned") and returns to the same active step without altering it
// (requirement R6).
func (s *Service) DetourFinish(id learning.SessionID, outcome assistance.DetourOutcome, expectedRevision uint64, requestID string) (DetourResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rec, ok := s.sessions[id]
	if !ok {
		return DetourResult{}, ErrSessionNotFound
	}
	ev, fresh, err := s.mutate(id, expectedRevision, requestID, eventstore.EventDetourFinished, map[string]string{"outcome": string(outcome)})
	if err != nil {
		return DetourResult{}, err
	}
	d := rec.session.Detour()
	if d == nil {
		return DetourResult{}, ErrNoOpenDetour
	}
	if fresh {
		if err := assistance.FinishDetour(d, outcome); err != nil {
			return DetourResult{}, err
		}
	}
	return DetourResult{State: d.State, Revision: ev.Revision}, nil
}
