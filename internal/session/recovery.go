package session

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/oseiaspereira88/codinho/internal/assistance"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

// ErrSessionUnrecoverable distinguishes historical sessions lacking a safe
// projection from unknown IDs. Their durable history remains untouched.
var ErrSessionUnrecoverable = errors.New("session: recovery unavailable; historical state lacks a compatible projection")

type startedPayload struct {
	Track             *trackSnapshot                `json:"track,omitempty"`
	NavigationVersion int                           `json:"navigation_version,omitempty"`
	RecoveryVersion   int                           `json:"recovery_version"`
	ChallengeID       string                        `json:"challenge_id"`
	Mode              string                        `json:"mode"`
	Input             StartInput                    `json:"input"`
	Policy            learning.SessionPolicy        `json:"policy"`
	Challenge         curriculum.ChallengeAuthoring `json:"challenge"`
	Digest            string                        `json:"content_sha256"`
}

func inputIdentity(in StartInput) string { raw, _ := json.Marshal(in); return string(raw) }
func contentDigest(ch curriculum.ChallengeAuthoring) string {
	raw, _ := json.Marshal(ch)
	return fmt.Sprintf("%x", sha256.Sum256(raw))
}
func jsonEqual(a, b []byte) bool {
	var x, y any
	return json.Unmarshal(a, &x) == nil && json.Unmarshal(b, &y) == nil && reflect.DeepEqual(x, y)
}
func pinnedChallenge(rec *record) (curriculum.ChallengeAuthoring, error) { return rec.pinned, nil }
func (s *Service) lookupError(id learning.SessionID) error {
	if s.unrecoverable[id] {
		return ErrSessionUnrecoverable
	}
	return ErrSessionNotFound
}

// retryRequest validates the original command identity before any state-based
// checks. Responses are reconstructed at the committed event, not today's node.
func (s *Service) retryRequest(id learning.SessionID, requestID, operation string, input, out any) (bool, error) {
	raw, err := json.Marshal([]any{operation, input})
	if err != nil {
		return false, err
	}
	s.requestDigest = fmt.Sprintf("%x", sha256.Sum256(raw))
	if s.unrecoverable[id] {
		return false, ErrSessionUnrecoverable
	}
	ev, ok := s.store.Request(string(id), requestID)
	if !ok {
		return false, nil
	}
	var meta struct {
		Digest string `json:"request_digest"`
	}
	if json.Unmarshal(ev.Payload, &meta) != nil || meta.Digest != s.requestDigest {
		return false, eventstore.ErrRevisionConflict
	}
	rec, err := s.recordAt(id, ev.Revision)
	if err != nil {
		return false, err
	}
	result, err := s.eventResult(rec, ev)
	if err != nil {
		return false, err
	}
	raw, err = json.Marshal(result)
	if err != nil {
		return false, err
	}
	return true, json.Unmarshal(raw, out)
}

func (s *Service) recordAt(id learning.SessionID, revision uint64) (*record, error) {
	var rec *record
	for _, ev := range s.store.Replay(string(id)) {
		if ev.Revision > revision {
			break
		}
		if ev.Type == eventstore.EventSessionStarted {
			var err error
			rec, _, err = restoreStart(ev)
			if err != nil {
				return nil, err
			}
		} else if rec != nil {
			if err := applyEvent(rec, ev); err != nil {
				return nil, err
			}
		}
	}
	if rec == nil {
		return nil, ErrSessionUnrecoverable
	}
	return rec, nil
}

func (s *Service) eventResult(rec *record, ev eventstore.Event) (any, error) {
	var p deltaPayload
	if err := json.Unmarshal(ev.Payload, &p); err != nil {
		return nil, err
	}
	switch ev.Type {
	case eventstore.EventEvidenceRecorded:
		var p struct {
			EvidenceID string `json:"evidence_id"`
		}
		if err := json.Unmarshal(ev.Payload, &p); err != nil {
			return nil, err
		}
		return EvidenceRecordResult{EvidenceID: p.EvidenceID, Revision: ev.Revision}, nil
	case eventstore.EventSessionPolicyChanged, eventstore.EventSessionPaused, eventstore.EventSessionResumed, eventstore.EventSessionFinished, eventstore.EventFeedbackRecorded, eventstore.EventReflectionRecorded, eventstore.EventStepCompleted:
		return LifecycleResult{State: rec.session.State, Revision: ev.Revision}, nil
	case eventstore.EventGranularityChanged:
		step, ok := findStep(rec.pinned, p.StepID)
		if !ok {
			return nil, ErrStepNotFound
		}
		return GranularityResult{StepID: p.StepID, Kind: step.Kind, Revision: ev.Revision}, nil
	case eventstore.EventStepAdvanced:
		if p.Done {
			return AdvanceResult{Done: true, Revision: ev.Revision}, nil
		}
		step, ok := findStep(rec.pinned, p.To)
		if !ok {
			return nil, ErrStepNotFound
		}
		return AdvanceResult{StepID: p.To, Kind: step.Kind, Revision: ev.Revision}, nil
	case eventstore.EventLearnerNextStepProposed:
		return ProposeNextStepResult{Revision: ev.Revision}, nil
	case eventstore.EventHintRequested, eventstore.EventSolutionRevealed:
		step, ok := findStep(rec.pinned, p.StepID)
		if !ok {
			return nil, ErrStepNotFound
		}
		return HintResult{StepID: learning.StepID(p.StepID), Level: p.Level, Kind: p.Kind, Objective: step.Instruction.Objective, Scope: step.Instruction.Scope, Concepts: step.Concepts, Revision: ev.Revision}, nil
	case eventstore.EventDetourStarted, eventstore.EventDetourFinished:
		if rec.session.Detour() == nil {
			return nil, ErrNoOpenDetour
		}
		return DetourResult{State: rec.session.Detour().State, Revision: ev.Revision}, nil
	case eventstore.EventEvaluationRecorded:
		var evaluation evaluationPayload
		if err := json.Unmarshal(ev.Payload, &evaluation); err != nil {
			return nil, err
		}
		revision, err := s.ensureAttempt(ev)
		if err != nil {
			return nil, err
		}
		return EvaluateResult{StepID: learning.StepID(p.StepID), Criteria: evaluation.Criteria, HasBlockingFailure: p.HasBlockingFailure, AttemptRecorded: evaluation.SubmissionIntent, Revision: revision}, nil
	}
	return nil, ErrSessionUnrecoverable
}

type evaluationPayload struct {
	StepID           string                     `json:"step_id"`
	Criteria         []learning.CriterionResult `json:"criteria"`
	SubmissionIntent bool                       `json:"submission_intent"`
	SolutionRevealed bool                       `json:"solution_revealed"`
}

// ensureAttempt completes the second part of a committed submission intent.
// Correlation uses the evaluation event ID, not a caller-controlled request ID.
func (s *Service) ensureAttempt(ev eventstore.Event) (uint64, error) {
	var p evaluationPayload
	if err := json.Unmarshal(ev.Payload, &p); err != nil {
		return 0, err
	}
	if !p.SubmissionIntent {
		return ev.Revision, nil
	}
	for _, candidate := range s.store.Replay(ev.StreamID) {
		if candidate.Type != eventstore.EventAttemptSubmitted {
			continue
		}
		var payload struct {
			EvaluationID string `json:"evaluation_event_id"`
		}
		if json.Unmarshal(candidate.Payload, &payload) == nil && payload.EvaluationID == ev.ID {
			return candidate.Revision, nil
		}
	}
	var ids []string
	for _, criterion := range p.Criteria {
		if criterion.EvidenceID != "" {
			ids = append(ids, string(criterion.EvidenceID))
		}
	}
	attempt, err := s.store.Append(ev.StreamID, s.store.Revision(ev.StreamID), "", eventstore.EventAttemptSubmitted, map[string]any{"step_id": p.StepID, "evidence_ids": ids, "solution_revealed": p.SolutionRevealed, "evaluation_event_id": ev.ID})
	return attempt.Revision, err
}

func activeProgress(id string) (*learning.StepProgress, error) {
	p := learning.NewStepProgress(learning.StepID(id))
	if err := p.Transition(learning.StepStateAvailable, false); err != nil {
		return nil, err
	}
	if err := p.Transition(learning.StepStateActive, false); err != nil {
		return nil, err
	}
	return p, nil
}

func restoreStart(ev eventstore.Event) (*record, startedPayload, error) {
	var p startedPayload
	if err := json.Unmarshal(ev.Payload, &p); err != nil {
		return nil, p, err
	}
	if p.NavigationVersion < 0 || p.NavigationVersion > 1 || p.RecoveryVersion != 1 || p.ChallengeID == "" || p.Challenge.ID != p.ChallengeID || p.Digest != contentDigest(p.Challenge) {
		return nil, p, ErrSessionUnrecoverable
	}
	if p.Track != nil {
		if p.Track.Version != 1 || len(p.Track.Challenges) == 0 || len(p.Track.Challenges) > curriculum.MaxTrackChallenges || p.Track.Digest != curriculum.PathIdentity(p.Track.Challenges) || contentDigest(p.Track.Challenges[0]) != p.Digest {
			return nil, p, ErrSessionUnrecoverable
		}
	}
	policy, err := learning.NewSessionPolicy(p.Policy.Mode, p.Policy.InitialDepth, p.Policy.Help.Kind, p.Policy.Disclosure.MaxLevel, p.Policy.Evaluation, p.Policy.Advance, p.Policy.TimeLimit)
	if err != nil {
		return nil, p, err
	}
	step, ok := firstStep(p.Challenge)
	if p.NavigationVersion == 1 {
		w, found := deriveWindow(p.Challenge, policy.InitialDepth)
		ok = found
		step, _ = findStep(p.Challenge, w.StepID)
	}
	if !ok {
		return nil, p, ErrChallengeHasNoSteps
	}
	domain := learning.NewLearningSession(learning.SessionID(ev.StreamID), policy, learning.CatalogRef{})
	if err := domain.Transition(learning.SessionStateActive); err != nil {
		return nil, p, err
	}
	active, err := activeProgress(step.ID)
	if err != nil {
		return nil, p, err
	}
	if err := domain.SetActiveInstruction(active); err != nil {
		return nil, p, err
	}
	return &record{track: p.Track, session: domain, challengeID: p.ChallengeID, pinned: p.Challenge, depth: policy.InitialDepth, cursor: step.ID, progress: map[string]savedProgress{}, covered: map[string]bool{}, choices: map[string]string{}}, p, nil
}

// recover projects only persisted facts, never consulting today's catalog.
func (s *Service) recover() {
	for _, ev := range s.store.ReplayAll() {
		id := learning.SessionID(ev.StreamID)
		if n, err := strconv.ParseUint(strings.TrimPrefix(ev.StreamID, "ses_"), 10, 64); err == nil && n > s.nextID.Load() {
			s.nextID.Store(n)
		}
		if ev.Type == eventstore.EventSessionStarted {
			if ev.RequestID != "" {
				s.startIDs[ev.RequestID] = id
			}
			if s.unrecoverable[id] {
				continue
			}
			rec, p, err := restoreStart(ev)
			if err != nil {
				s.unrecoverable[id] = true
				continue
			}
			s.sessions[id] = rec
			if ev.RequestID != "" {
				step, _ := findStep(p.Challenge, string(rec.session.ActiveStep().StepID))
				s.startResults[ev.RequestID] = StartResult{Track: rec.trackStatus(), SessionID: id, ActiveStep: rec.session.ActiveStep().StepID, Kind: step.Kind, Objective: step.Instruction.Objective, Revision: ev.Revision, Disclosure: disclosureFor(p.Policy.Disclosure, rec.session.ActiveStep())}
				s.startInputs[ev.RequestID] = inputIdentity(p.Input)
			}
			continue
		}
		rec, ok := s.sessions[id]
		if !ok || s.unrecoverable[id] {
			continue
		}
		if err := applyEvent(rec, ev); err != nil {
			s.unrecoverable[id] = true
			delete(s.sessions, id)
			continue
		}
		if ev.RequestID != "" && (ev.Type == eventstore.EventHintRequested || ev.Type == eventstore.EventSolutionRevealed) {
			var p deltaPayload
			if json.Unmarshal(ev.Payload, &p) == nil && !p.Blocked {
				step, ok := findStep(rec.pinned, p.StepID)
				if ok {
					s.hintResults[idempotencyKey(id, ev.RequestID)] = HintResult{StepID: learning.StepID(p.StepID), Level: p.Level, Kind: p.Kind, Objective: step.Instruction.Objective, Scope: step.Instruction.Scope, Concepts: step.Concepts, Revision: ev.Revision}
				}
			}
		}
	}
	// Complete pending submission intents only after the entire stream has
	// passed projection validation. Never append to an incompatible history.
	for _, ev := range s.store.ReplayAll() {
		id := learning.SessionID(ev.StreamID)
		if ev.Type != eventstore.EventEvaluationRecorded || s.sessions[id] == nil || s.unrecoverable[id] {
			continue
		}
		if _, err := s.ensureAttempt(ev); err != nil {
			s.unrecoverable[id] = true
			delete(s.sessions, id)
		}
	}
}

type deltaPayload struct {
	TrackCursor        int                           `json:"track_cursor"`
	Done               bool                          `json:"done"`
	NavigationVersion  int                           `json:"navigation_version"`
	Cursor             string                        `json:"cursor"`
	Covered            bool                          `json:"covered"`
	ChoiceParent       string                        `json:"choice_parent"`
	ChoiceID           string                        `json:"choice_id"`
	StepID             string                        `json:"step_id"`
	From               string                        `json:"from"`
	To                 string                        `json:"to"`
	Depth              learning.Depth                `json:"depth"`
	Help               learning.HelpPolicyKind       `json:"help"`
	Evaluation         learning.EvaluationPolicyKind `json:"evaluation"`
	Reason             string                        `json:"reason"`
	Outcome            assistance.DetourOutcome      `json:"outcome"`
	Level              learning.DisclosureLevel      `json:"level"`
	Kind               string                        `json:"kind"`
	Free               bool                          `json:"free"`
	Direct             bool                          `json:"direct"`
	Blocked            bool                          `json:"blocked"`
	Override           bool                          `json:"override"`
	HasBlockingFailure bool                          `json:"has_blocking_failure"`
}

// applyEvent is shared by replay and pre-append validation. It never performs IO.
func applyEvent(rec *record, ev eventstore.Event) error {
	var p deltaPayload
	if err := json.Unmarshal(ev.Payload, &p); err != nil {
		return err
	}
	if p.NavigationVersion < 0 || p.NavigationVersion > 1 {
		return ErrSessionUnrecoverable
	}
	active := rec.session.ActiveStep()
	if p.StepID != "" && ev.Type != eventstore.EventGranularityChanged && ev.Type != eventstore.EventLearnerNextStepProposed && ev.Type != eventstore.EventObservationRecorded && ev.Type != eventstore.EventCheckExecuted && ev.Type != eventstore.EventAttemptSubmitted {
		if active == nil || string(active.StepID) != p.StepID {
			return ErrNoActiveStep
		}
	}
	switch ev.Type {
	case eventstore.EventSessionPaused:
		return rec.session.Transition(learning.SessionStatePaused)
	case eventstore.EventSessionResumed:
		return rec.session.Transition(learning.SessionStateActive)
	case eventstore.EventSessionFinished:
		return rec.session.Transition(learning.SessionStateCompleted)
	case eventstore.EventSessionPolicyChanged:
		policy := rec.session.Policy
		if p.Help != "" {
			policy.Help.Kind = p.Help
		}
		if p.Evaluation != "" {
			policy.Evaluation = p.Evaluation
		}
		validated, err := learning.NewSessionPolicy(policy.Mode, policy.InitialDepth, policy.Help.Kind, policy.Disclosure.MaxLevel, policy.Evaluation, policy.Advance, policy.TimeLimit)
		if err != nil {
			return err
		}
		rec.session.Policy = validated
	case eventstore.EventGranularityChanged, eventstore.EventStepAdvanced:
		if p.NavigationVersion == 1 {
			if err := rec.navigationAllowed(); err != nil {
				return err
			}
			if ev.Type == eventstore.EventGranularityChanged {
				if rec.exhausted {
					return ErrNavigationUnavailable
				}
				w, err := rec.windowAt(p.Depth)
				if err != nil {
					return err
				}
				cursor := rec.cursor
				if cursor == "" && active != nil {
					cursor = string(active.StepID)
				}
				if len(nodePath(challengeTree(rec.pinned), w.StepID)) > len(nodePath(challengeTree(rec.pinned), cursor)) {
					cursor = w.StepID
				}
				if w.StepID != p.StepID || cursor != p.Cursor {
					return ErrSessionUnrecoverable
				}
				if err := rec.activate(p.StepID); err != nil {
					return err
				}
				rec.depth = p.Depth
				rec.cursor = cursor
			} else {
				if active == nil || string(active.StepID) != p.From {
					return ErrNoActiveStep
				}
				next, projected, err := rec.nextNavigation(p.Override, p.ChoiceID)
				if err != nil {
					return err
				}
				if projected.trackCursor != p.TrackCursor || next.node.ID != p.To || p.Done != (p.To == "") || len(next.options) > 0 {
					return ErrInvalidNextStep
				}
				parent := ""
				for key, v := range projected.choices {
					if rec.choices[key] != v {
						parent = key
					}
				}
				if parent != p.ChoiceParent {
					return ErrInvalidNextStep
				}
				*rec = *projected
				if p.Done {
					rec.exhausted = true
					return nil
				}
				if err := rec.activate(p.To); err != nil {
					return err
				}
				rec.cursor = p.To
			}
			return nil
		}
		// Historical navigation resets its window exactly as originally recorded.
		rec.saveWindow()
		nextID := p.StepID
		if ev.Type == eventstore.EventStepAdvanced {
			nextID = p.To
			if active == nil || string(active.StepID) != p.From {
				return ErrNoActiveStep
			}
		}
		if _, ok := findStep(rec.pinned, nextID); !ok {
			return ErrStepNotFound
		}
		next, err := activeProgress(nextID)
		if err != nil {
			return err
		}
		if ev.Type == eventstore.EventStepAdvanced {
			if err := rec.session.Advance(next, p.Override); err != nil {
				return err
			}
		} else {
			rec.session.ClearActiveInstruction()
			if err := rec.session.SetActiveInstruction(next); err != nil {
				return err
			}
			rec.depth = p.Depth
		}
		rec.cleanEvaluation = false
		rec.cursor = nextID
	case eventstore.EventHintRequested, eventstore.EventSolutionRevealed:
		if p.Blocked || p.Free {
			return nil
		}
		if active == nil {
			return ErrNoActiveStep
		}
		var err error
		if p.Direct {
			err = active.GrantDirect(p.Level, rec.session.Policy.Disclosure)
		} else {
			err = active.GrantHint(p.Level, rec.session.Policy.Disclosure)
		}
		if err != nil {
			return err
		}
		if p.Level == learning.DisclosureSolution {
			active.RevealSolution()
		}
	case eventstore.EventDetourStarted:
		_, err := assistance.StartDetour(rec.session, p.Reason)
		return err
	case eventstore.EventDetourFinished:
		if rec.session.Detour() == nil {
			return ErrNoOpenDetour
		}
		return assistance.FinishDetour(rec.session.Detour(), p.Outcome)
	case eventstore.EventEvaluationRecorded:
		if active == nil {
			return ErrNoActiveStep
		}
		if active.State != learning.StepStateEvaluated {
			if err := active.Transition(learning.StepStateEvaluated, false); err != nil {
				return err
			}
		}
		rec.cleanEvaluation = !p.HasBlockingFailure
	case eventstore.EventStepCompleted:
		if active == nil {
			return ErrNoActiveStep
		}
		if p.NavigationVersion == 1 {
			if rec.exhausted {
				return ErrNavigationUnavailable
			}
			if err := rec.navigationAllowed(); err != nil {
				return err
			}
			node, _ := findStep(rec.pinned, p.StepID)
			if p.Covered != rec.coversWindow(node) {
				return ErrSessionUnrecoverable
			}
		}
		if err := active.Complete(p.Override); err != nil {
			return err
		}
		if p.NavigationVersion == 1 && p.Covered {
			rec.covered[p.StepID] = true
		}
		return nil
	case eventstore.EventEvidenceRecorded, eventstore.EventObservationRecorded, eventstore.EventCheckExecuted, eventstore.EventFeedbackRecorded, eventstore.EventAttemptSubmitted, eventstore.EventReflectionRecorded, eventstore.EventLearnerNextStepProposed, eventstore.EventMasteryProjected, eventstore.EventReviewScheduled:
		// These events carry evidence or projections, not session state mutations.
	default:
		return fmt.Errorf("%w: unsupported event %s", ErrSessionUnrecoverable, ev.Type)
	}
	return nil
}
