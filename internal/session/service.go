// Package session orchestrates a learning session's lifecycle, single
// active instruction and progressive disclosure over the pure domain
// (internal/learning), authored content (internal/curriculum) and the
// durable event log (internal/eventstore). It never generates feedback,
// runs checks or computes mastery (non-goals of this spec).
package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/oseiaspereira88/codinho/internal/assessment"
	"github.com/oseiaspereira88/codinho/internal/assistance"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/evidence"
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
	pinned      curriculum.ChallengeAuthoring
	depth       learning.Depth
	// cleanEvaluation tracks whether the active step's most recent
	// evaluation had no blocking failure, gating step_complete's
	// requires_positive_evaluation policy (requirement R7). It resets
	// whenever the active step changes.
	cleanEvaluation bool
}

// Service orchestrates sessions: lifecycle, single active instruction,
// disclosure and granularity. It holds a session registry in memory keyed
// by SessionID; durability comes from appending every state-changing
// operation to the event store before applying it in memory.
type Service struct {
	catalog *curriculum.Catalog
	store   *eventstore.Store

	mu                sync.Mutex
	evidenceValidator EvaluationEvidenceValidator
	requestDigest     string // current command, guarded by mu; persisted with its event
	sessions          map[learning.SessionID]*record
	startResults      map[string]StartResult
	startInputs       map[string]string
	startIDs          map[string]learning.SessionID
	unrecoverable     map[learning.SessionID]bool
	hintResults       map[string]HintResult
	detourResults     map[string]DetourResult
	nextID            atomic.Uint64
}

// idempotencyKey scopes a request_id to its session, matching the
// eventstore's own (streamID, requestID) scoping, so two sessions reusing
// the same client-chosen request_id never collide.
func idempotencyKey(id learning.SessionID, requestID string) string {
	return string(id) + "\x00" + requestID
}

// New wires a Service to its catalog and event store. The evidenceStore argument
// is retained for source compatibility; only the injected validator authorizes
// evidence consumption. Missing authorization always fails closed.
func New(catalog *curriculum.Catalog, store *eventstore.Store, evidenceStore *evidence.Store, validators ...EvaluationEvidenceValidator) *Service {
	s := &Service{
		catalog:       catalog,
		store:         store,
		sessions:      map[learning.SessionID]*record{},
		startResults:  map[string]StartResult{},
		startInputs:   map[string]string{},
		startIDs:      map[string]learning.SessionID{},
		unrecoverable: map[learning.SessionID]bool{},
		hintResults:   map[string]HintResult{},
		detourResults: map[string]DetourResult{},
	}
	if len(validators) > 0 {
		s.evidenceValidator = validators[0]
	}
	s.recover()
	return s
}

// ActiveChecks returns the checks declared by the session's fixed
// challenge and its currently active step, so a caller (safe-check-
// executor) resolves check_id exclusively within this exact session's
// version (requirement R1) without check_run needing a challenge_id
// parameter of its own.
func (s *Service) ActiveChecks(id learning.SessionID) ([]curriculum.CheckAuthoring, learning.StepID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec, ok := s.sessions[id]
	if !ok {
		return nil, "", s.lookupError(id)
	}
	challenge, err := pinnedChallenge(rec)
	if err != nil {
		return nil, "", err
	}
	var stepID learning.StepID
	if active := rec.session.ActiveStep(); active != nil {
		stepID = active.StepID
	}
	return challenge.Checks, stepID, nil
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
	// TimeLimit fixes an optional duration at session start (interview-mode
	// requirement R1: "duração"); nil means off, matching
	// learning.SessionPolicy.TimeLimit's own zero value.
	TimeLimit *time.Duration
	RequestID string // idempotency key for retries (requirement R7, R8)
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
		if id, ok := s.startIDs[in.RequestID]; ok && s.unrecoverable[id] {
			return StartResult{}, ErrSessionUnrecoverable
		}
		if cached, ok := s.startResults[in.RequestID]; ok {
			if s.startInputs[in.RequestID] != inputIdentity(in) {
				return StartResult{}, eventstore.ErrRevisionConflict
			}
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
	// Each mode supplies its own independent defaults (PROJECT.md §8.2;
	// learning-practice-debug-modes requirement R1); an explicit
	// StartInput field always overrides its mode's default (R2).
	defaults := DefaultsForMode(mode)
	depth := in.Depth
	if depth == "" {
		depth = defaults.Depth
	}
	help := in.Help
	if help == "" {
		help = defaults.Help
	}
	disclosureMax := in.DisclosureMax
	if disclosureMax == 0 {
		disclosureMax = defaults.Disclosure
	}
	evaluation := in.Evaluation
	if evaluation == "" {
		evaluation = defaults.Evaluation
	}
	policy, err := learning.NewSessionPolicy(mode, depth, help, disclosureMax, evaluation, learning.AdvanceExplicit, in.TimeLimit)
	if err != nil {
		return StartResult{}, err
	}

	id := s.nextID.Add(1)
	sessionID := learning.SessionID(fmt.Sprintf("ses_%d", id))

	// Append first: the durable record of intent to start must exist
	// before the in-memory session is exposed to callers (local-event-store
	// R1/R3 ordering contract).
	payload := startedPayload{RecoveryVersion: 1, ChallengeID: in.ChallengeID, Mode: string(mode), Input: in, Policy: policy, Challenge: challenge, Digest: contentDigest(challenge)}
	ev, err := s.store.Append(string(sessionID), 0, in.RequestID, eventstore.EventSessionStarted, payload)
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

	s.sessions[sessionID] = &record{session: domainSession, challengeID: in.ChallengeID, pinned: challenge, depth: depth}

	result := StartResult{
		SessionID:  sessionID,
		ActiveStep: progress.StepID,
		Objective:  step.Instruction.Objective,
		Revision:   ev.Revision,
		Disclosure: disclosureFor(policy.Disclosure, progress),
	}
	if in.RequestID != "" {
		s.startResults[in.RequestID] = result
		s.startInputs[in.RequestID] = inputIdentity(in)
		s.startIDs[in.RequestID] = sessionID
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
		return GetResult{}, s.lookupError(id)
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
		return Instruction{}, s.lookupError(id)
	}
	active := rec.session.ActiveStep()
	if active == nil {
		return Instruction{}, s.lookupError(id)
	}
	challenge, err := pinnedChallenge(rec)
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
	encoded, err := json.Marshal(payload)
	if err != nil {
		return ev, false, err
	}
	fields := map[string]any{}
	if string(encoded) != "null" {
		if err := json.Unmarshal(encoded, &fields); err != nil {
			return ev, false, err
		}
	}
	fields["request_digest"] = s.requestDigest
	payload = fields
	if existing, ok := s.store.Request(string(id), requestID); ok {
		raw, marshalErr := json.Marshal(payload)
		if marshalErr != nil {
			return ev, false, marshalErr
		}
		if existing.Type != eventType || !jsonEqual(existing.Payload, raw) {
			return ev, false, eventstore.ErrRevisionConflict
		}
		return existing, false, nil
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return ev, false, err
	}
	rec := *s.sessions[id]
	rec.session = rec.session.Clone()
	if err := applyEvent(&rec, eventstore.Event{Type: eventType, Payload: raw}); err != nil {
		return ev, false, err
	}
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

	var prior LifecycleResult
	if found, err := s.retryRequest(in.SessionID, in.RequestID, "Configure", in, &prior); found || err != nil {
		return prior, err
	}

	rec, ok := s.sessions[in.SessionID]
	if !ok {
		return LifecycleResult{}, s.lookupError(in.SessionID)
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
// completion from step state (requirement R5). It always records an
// empty reason; call FinishWithReason directly to record why (e.g.
// interview-mode's explicit-vs-timeout distinction).
func (s *Service) Finish(id learning.SessionID, expectedRevision uint64, requestID string) (LifecycleResult, error) {
	return s.FinishWithReason(id, "", expectedRevision, requestID)
}

func (s *Service) transitionLifecycle(id learning.SessionID, expectedRevision uint64, requestID string, eventType eventstore.EventType, to learning.SessionState) (LifecycleResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var prior LifecycleResult
	if found, err := s.retryRequest(id, requestID, "transitionLifecycle", []any{eventType}, &prior); found || err != nil {
		return prior, err
	}

	rec, ok := s.sessions[id]
	if !ok {
		return LifecycleResult{}, s.lookupError(id)
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
// reason is durably recorded on the granularity_changed event whether it
// comes from a human/skill rationale or from SuggestGranularity, so every
// change is explainable after the fact (learning-practice-debug-modes
// requirement R8, PROJECT.md §8.5: "deve ser informado ao aluno"); it may
// be empty for a purely manual adjustment with no stated reason.
func (s *Service) GranularityAdjust(id learning.SessionID, depth learning.Depth, reason string, expectedRevision uint64, requestID string) (GranularityResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var prior GranularityResult
	if found, err := s.retryRequest(id, requestID, "GranularityAdjust", []any{depth, reason}, &prior); found || err != nil {
		return prior, err
	}

	rec, ok := s.sessions[id]
	if !ok {
		return GranularityResult{}, s.lookupError(id)
	}
	challenge, err := pinnedChallenge(rec)
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
		"reason":  reason,
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
		rec.cleanEvaluation = false
	}
	return GranularityResult{StepID: window.StepID, Kind: window.Kind, Revision: ev.Revision}, nil
}

// ErrStepNotFound is returned when a proposed step ID does not exist in
// the session's fixed challenge.
var ErrStepNotFound = errors.New("session: step not found in challenge")

// ProposeNextStepResult is what ProposeNextStep returns.
type ProposeNextStepResult struct {
	Revision uint64
}

// ProposeNextStep durably records that the learner proposed stepID as
// what comes next, without advancing anything: it is a pure autonomy
// signal (PROJECT.md §21.5, "capacidade de propor o próximo passo";
// learning-practice-debug-modes requirement R9), never a substitute for
// an explicit step_advance call. stepID must name a real step in the
// session's fixed challenge, so the recorded signal is never noise.
func (s *Service) ProposeNextStep(id learning.SessionID, stepID learning.StepID, expectedRevision uint64, requestID string) (ProposeNextStepResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var prior ProposeNextStepResult
	if found, err := s.retryRequest(id, requestID, "ProposeNextStep", []any{stepID}, &prior); found || err != nil {
		return prior, err
	}

	rec, ok := s.sessions[id]
	if !ok {
		return ProposeNextStepResult{}, s.lookupError(id)
	}
	challenge, err := pinnedChallenge(rec)
	if err != nil {
		return ProposeNextStepResult{}, err
	}
	if _, ok := findStep(challenge, string(stepID)); !ok {
		return ProposeNextStepResult{}, ErrStepNotFound
	}

	ev, _, err := s.mutate(id, expectedRevision, requestID, eventstore.EventLearnerNextStepProposed, map[string]string{
		"step_id": string(stepID),
	})
	if err != nil {
		return ProposeNextStepResult{}, err
	}
	return ProposeNextStepResult{Revision: ev.Revision}, nil
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
		return nil, nil, curriculum.StepAuthoring{}, s.lookupError(id)
	}
	active := rec.session.ActiveStep()
	if active == nil {
		return nil, nil, curriculum.StepAuthoring{}, ErrNoActiveStep
	}
	challenge, err := pinnedChallenge(rec)
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

	var prior HintResult
	if found, err := s.retryRequest(id, requestID, "HintRequest", []any{confirmSolution}, &prior); found || err != nil {
		return prior, err
	}

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

	var prior HintResult
	if found, err := s.retryRequest(id, requestID, "SyntaxRecallGet", nil, &prior); found || err != nil {
		return prior, err
	}

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
	payload := map[string]any{"step_id": string(active.StepID), "level": int(level), "kind": kind, "free": free, "direct": direct}
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

	var prior DetourResult
	if found, err := s.retryRequest(id, requestID, "DetourStart", []any{reason}, &prior); found || err != nil {
		return prior, err
	}

	rec, ok := s.sessions[id]
	if !ok {
		return DetourResult{}, s.lookupError(id)
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

	var prior DetourResult
	if found, err := s.retryRequest(id, requestID, "DetourFinish", []any{outcome}, &prior); found || err != nil {
		return prior, err
	}

	rec, ok := s.sessions[id]
	if !ok {
		return DetourResult{}, s.lookupError(id)
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

// FeedbackPrepare assembles the read-only context feedback_record's caller
// authors the actual feedback from (requirement R1). It never mutates
// session state or writes an event.
func (s *Service) FeedbackPrepare(id learning.SessionID, question string) (assessment.FeedbackPacket, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, _, step, err := s.activeWindow(id)
	if err != nil {
		return assessment.FeedbackPacket{}, err
	}
	return assessment.PrepareFeedback(step, question), nil
}

// FeedbackRecord persists feedback already authored elsewhere (the tutor
// agent), without approving, consuming an attempt, completing or advancing
// the step (requirement R2; PROJECT.md §8.6 invariant 3).
func (s *Service) FeedbackRecord(id learning.SessionID, feedbackType learning.FeedbackType, text string, blockingOverride *bool, expectedRevision uint64, requestID string) (LifecycleResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var prior LifecycleResult
	if found, err := s.retryRequest(id, requestID, "FeedbackRecord", []any{feedbackType, text, blockingOverride}, &prior); found || err != nil {
		return prior, err
	}

	rec, active, _, err := s.activeWindow(id)
	if err != nil {
		return LifecycleResult{}, err
	}
	fb, err := learning.NewFeedbackRecord(active.StepID, feedbackType, text, blockingOverride)
	if err != nil {
		return LifecycleResult{}, err
	}
	ev, _, err := s.mutate(id, expectedRevision, requestID, eventstore.EventFeedbackRecorded, map[string]any{
		"step_id": string(fb.StepID), "type": string(fb.Type), "blocking": fb.Blocking,
	})
	if err != nil {
		return LifecycleResult{}, err
	}
	return LifecycleResult{State: rec.session.State, Revision: ev.Revision}, nil
}

// EvaluateInput is what StepEvaluate receives (requirement R3).
type EvaluateInput struct {
	SessionID        learning.SessionID
	Criteria         []assessment.CriterionInput
	SubmissionIntent bool
	ExpectedRevision uint64
	RequestID        string
}

// EvaluateResult is what StepEvaluate returns.
type EvaluateResult struct {
	StepID             learning.StepID
	Criteria           []learning.CriterionResult
	HasBlockingFailure bool
	AttemptRecorded    bool
	Revision           uint64
}

// verdictForCheckOutcome maps a safe-check-executor Outcome string to the
// verdict step_evaluate reports (requirement R10).
func verdictForCheckOutcome(outcome string) learning.EvaluationVerdict {
	switch outcome {
	case "pass":
		return learning.VerdictMet
	case "fail":
		return learning.VerdictNotMet
	case "skipped":
		return learning.VerdictNotApplicable
	default: // "error" or anything unrecognized: infra failure proves nothing either way
		return learning.VerdictUnverifiable
	}
}

func (s *Service) StepEvaluate(in EvaluateInput) (EvaluateResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var prior EvaluateResult
	if found, err := s.retryRequest(in.SessionID, in.RequestID, "StepEvaluate", in, &prior); found || err != nil {
		return prior, err
	}

	rec, active, _, err := s.activeWindow(in.SessionID)
	if err != nil {
		return EvaluateResult{}, err
	}
	if active.State != learning.StepStateActive && active.State != learning.StepStateEvaluated {
		return EvaluateResult{}, learning.DomainError{Code: learning.ErrCodeStepNotEvaluable, Detail: string(active.State)}
	}

	ctx := evidenceContext(in.SessionID, rec)
	lineage := make([]EvidenceLineage, 0, len(in.Criteria))
	results := make([]learning.CriterionResult, 0, len(in.Criteria))
	for _, c := range in.Criteria {
		r, err := assessment.Resolve(c)
		if err != nil {
			return EvaluateResult{}, err
		}
		if c.EvidenceID != "" {
			if s.evidenceValidator == nil {
				return EvaluateResult{}, ErrEvaluationEvidenceInvalid
			}
			proof, err := s.evidenceValidator.Validate(ctx, c)
			if err != nil {
				return EvaluateResult{}, err
			}
			lineage = append(lineage, proof)
			if r.Kind == learning.StructuralCriterionKind {
				r.Verdict = learning.VerdictUnverifiable
				if proof.Kind == "check" {
					r.Verdict = verdictForCheckOutcome(proof.Outcome)
				}
			}
		}
		results = append(results, r)
	}
	blocking := (learning.Evaluation{StepID: active.StepID, Criteria: results}).HasBlockingFailure()
	// Sample again immediately before append. The revision guard rejects
	// concurrent event producers; external filesystem writers remain outside it.
	for _, c := range in.Criteria {
		if c.EvidenceID != "" {
			if _, err := s.evidenceValidator.Validate(ctx, c); err != nil {
				return EvaluateResult{}, err
			}
		}
	}

	ev, fresh, err := s.mutate(in.SessionID, in.ExpectedRevision, in.RequestID, eventstore.EventEvaluationRecorded, map[string]any{
		"step_id":              string(active.StepID),
		"criteria":             results,
		"submission_intent":    in.SubmissionIntent,
		"has_blocking_failure": blocking,
		"evidence_lineage":     lineage,
		"solution_revealed":    active.SolutionRevealed,
	})
	if err != nil {
		return EvaluateResult{}, err
	}
	if fresh {
		if active.State == learning.StepStateActive {
			if err := active.Transition(learning.StepStateEvaluated, false); err != nil {
				return EvaluateResult{}, err
			}
		}
		rec.cleanEvaluation = !blocking
	}

	revision, err := s.ensureAttempt(ev)
	if err != nil {
		return EvaluateResult{}, err
	}

	return EvaluateResult{
		StepID:             active.StepID,
		Criteria:           results,
		HasBlockingFailure: blocking,
		AttemptRecorded:    in.SubmissionIntent,
		Revision:           revision,
	}, nil
}

// ReflectionRecord persists a short reasoning answer and its descriptive
// assessment, distinct from implementation evidence (requirement R9;
// PROJECT.md §8.9). It never changes step or session state.
func (s *Service) ReflectionRecord(id learning.SessionID, competencyID learning.CompetencyID, prompt, answer, assessmentText string, expectedRevision uint64, requestID string) (LifecycleResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var prior LifecycleResult
	if found, err := s.retryRequest(id, requestID, "ReflectionRecord", []any{competencyID, prompt, answer, assessmentText}, &prior); found || err != nil {
		return prior, err
	}

	rec, active, _, err := s.activeWindow(id)
	if err != nil {
		return LifecycleResult{}, err
	}
	if answer == "" {
		return LifecycleResult{}, learning.DomainError{Code: learning.ErrCodeInvalidValue, Detail: "reflection answer is required"}
	}
	ev, _, err := s.mutate(id, expectedRevision, requestID, eventstore.EventReflectionRecorded, map[string]any{
		"step_id": string(active.StepID), "competency_id": string(competencyID),
		"prompt": prompt, "answer": answer, "assessment": assessmentText,
	})
	if err != nil {
		return LifecycleResult{}, err
	}
	return LifecycleResult{State: rec.session.State, Revision: ev.Revision}, nil
}

// StepComplete marks the active step completed when the step's authored
// completion policy is satisfied (or override is set), never activating
// the next node (requirement R7, R8; PROJECT.md §8.6 "Conclusão").
func (s *Service) StepComplete(id learning.SessionID, confirm, override bool, expectedRevision uint64, requestID string) (LifecycleResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var prior LifecycleResult
	if found, err := s.retryRequest(id, requestID, "StepComplete", []any{confirm, override}, &prior); found || err != nil {
		return prior, err
	}

	rec, active, step, err := s.activeWindow(id)
	if err != nil {
		return LifecycleResult{}, err
	}
	if !override {
		if step.Completion.RequiresPositiveEvaluation && !rec.cleanEvaluation {
			return LifecycleResult{}, learning.DomainError{Code: learning.ErrCodeCompletionPolicyNotMet, Detail: "requires_positive_evaluation"}
		}
		if step.Completion.RequiresUserConfirmation && !confirm {
			return LifecycleResult{}, learning.DomainError{Code: learning.ErrCodeCompletionPolicyNotMet, Detail: "requires_user_confirmation"}
		}
	}

	ev, fresh, err := s.mutate(id, expectedRevision, requestID, eventstore.EventStepCompleted, map[string]any{
		"step_id": string(active.StepID), "override": override, "clean_evaluation": rec.cleanEvaluation, "confirmed": confirm,
	})
	if err != nil {
		return LifecycleResult{}, err
	}
	if fresh {
		if err := active.Complete(override); err != nil {
			return LifecycleResult{}, err
		}
	}
	return LifecycleResult{State: rec.session.State, Revision: ev.Revision}, nil
}

// AdvanceOption is one candidate next step when the active step branches
// into more than one child (PROJECT.md §15.6 "informa opções quando houver
// ramificação").
type AdvanceOption struct {
	StepID string `json:"step_id"`
	Kind   string `json:"kind"`
}

// AdvanceResult is what StepAdvance returns: either a single activated
// next step, a list of branch options to choose from, or Done when the
// challenge tree is exhausted. Branches and Done never mutate state or
// consume expectedRevision/requestID — only activating a single next step
// does.
type AdvanceResult struct {
	StepID   string
	Kind     string
	Branches []AdvanceOption
	Done     bool
	Revision uint64
}

// StepAdvance activates the next permitted node in the authored tree's
// document order, requiring the current active step to be completed
// unless override is set (requirement R7; PROJECT.md §8.6 "Avanço").
func (s *Service) StepAdvance(id learning.SessionID, override bool, expectedRevision uint64, requestID string) (AdvanceResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var prior AdvanceResult
	if found, err := s.retryRequest(id, requestID, "StepAdvance", []any{override}, &prior); found || err != nil {
		return prior, err
	}

	rec, ok := s.sessions[id]
	if !ok {
		return AdvanceResult{}, s.lookupError(id)
	}
	active := rec.session.ActiveStep()
	if active == nil {
		return AdvanceResult{}, ErrNoActiveStep
	}
	challenge, err := pinnedChallenge(rec)
	if err != nil {
		return AdvanceResult{}, err
	}
	roots, ok := rootSteps(challenge)
	if !ok {
		return AdvanceResult{}, ErrNoWindowAtDepth
	}
	next, branches, found, ok := advanceFrom(roots, string(active.StepID), nil)
	if !found {
		return AdvanceResult{}, ErrNotFound
	}
	if len(branches) > 0 {
		opts := make([]AdvanceOption, len(branches))
		for i, b := range branches {
			opts[i] = AdvanceOption{StepID: b.ID, Kind: b.Kind}
		}
		return AdvanceResult{Branches: opts, Revision: s.store.Revision(string(id))}, nil
	}
	if !ok {
		return AdvanceResult{Done: true, Revision: s.store.Revision(string(id))}, nil
	}

	ev, fresh, err := s.mutate(id, expectedRevision, requestID, eventstore.EventStepAdvanced, map[string]any{
		"from": string(active.StepID), "to": next.ID, "override": override,
	})
	if err != nil {
		return AdvanceResult{}, err
	}
	if fresh {
		nextProgress := learning.NewStepProgress(learning.StepID(next.ID))
		if err := nextProgress.Transition(learning.StepStateAvailable, false); err != nil {
			return AdvanceResult{}, err
		}
		if err := nextProgress.Transition(learning.StepStateActive, false); err != nil {
			return AdvanceResult{}, err
		}
		if err := rec.session.Advance(nextProgress, override); err != nil {
			return AdvanceResult{}, err
		}
		rec.cleanEvaluation = false
	}
	return AdvanceResult{StepID: next.ID, Kind: next.Kind, Revision: ev.Revision}, nil
}
