package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/mastery"
)

// masteryStreamID is the one global eventstore stream every mastery
// signal and schedule lives on: progress is learner-wide, not scoped to
// one session (Decision 3, mastery-review-scheduling).
const masteryStreamID = "mastery"

// ErrInvalidDimension is returned when Dimension is not one of the eight
// declared dimensions (PROJECT.md §21.3).
var ErrInvalidDimension = errors.New("application: unknown mastery dimension")

// ErrInvalidCompetency is returned when CompetencyID or EvidenceID is
// empty.
var ErrInvalidCompetency = errors.New("application: competency_id and evidence_id are required")

// ProgressService orchestrates mastery-review-scheduling: it never
// invents evidence — every signal cites an evidence_id already recorded
// elsewhere (Decision 2) — and every read recomputes projections and
// schedules from the durable log rather than trusting a cache
// (requirement R8, Decision 3).
type ProgressService struct {
	store *eventstore.Store
}

// NewProgressService wires a ProgressService to the shared event store.
func NewProgressService(store *eventstore.Store) *ProgressService {
	return &ProgressService{store: store}
}

// signalPayload is the durable shape of one mastery_projected event: the
// full Signal that produced it, plus the rule version and resulting state
// (Compatibility: replay always re-derives State via the recorded
// RuleVersion, never trusts this cached copy — see recordedSignals).
type signalPayload struct {
	ContentProvenance string `json:"content_provenance"`
	CompetencyID      string `json:"competency_id"`
	Dimension         string `json:"dimension"`
	EvidenceID        string `json:"evidence_id"`
	Variant           string `json:"variant,omitempty"`
	HelpUsed          bool   `json:"help_used"`
	SolutionRevealed  bool   `json:"solution_revealed"`
	Success           bool   `json:"success"`
	ObservedAt        string `json:"observed_at"`
	RuleVersion       int    `json:"rule_version"`
	State             string `json:"state"`
}

// schedulePayload is the durable shape of one review_scheduled event
// (PROJECT.md §18.3 minimum event). It is written for every qualifying
// signal as the required audit record, but reads reconstruct current
// schedules from the mastery_projected signal log via mastery.
// FoldSchedules, not from this event, so there is exactly one source of
// truth for replay (Decision 3).
type schedulePayload struct {
	CompetencyID string `json:"competency_id"`
	DueAt        string `json:"due_at"`
	IntervalDays int    `json:"interval_days"`
	Reason       string `json:"reason"`
}

// EvidenceInput is one mastery_evidence_record call.
type EvidenceInput struct {
	CompetencyID     string
	Dimension        string
	EvidenceID       string
	Variant          string
	HelpUsed         bool
	SolutionRevealed bool
	Success          bool
	ExpectedRevision uint64
	RequestID        string
}

// EvidenceResult is what mastery_evidence_record reports.
type EvidenceResult struct {
	ContentProvenance string
	CompetencyID      string
	Dimension         string
	State             string
	Revision          uint64
}

// RecordEvidence appends sig to the mastery log and returns the resulting
// projection, computed by replaying the log including this new signal
// (requirement R8: the report is never anything but a real recalculation).
func (p *ProgressService) RecordEvidence(in EvidenceInput) (EvidenceResult, error) {
	dim := mastery.Dimension(in.Dimension)
	if !dim.Valid() {
		return EvidenceResult{}, fmt.Errorf("%w: %q", ErrInvalidDimension, in.Dimension)
	}
	if in.CompetencyID == "" || in.EvidenceID == "" {
		return EvidenceResult{}, ErrInvalidCompetency
	}

	if ev, ok := p.store.Request(masteryStreamID, in.RequestID); ok {
		var previous signalPayload
		if ev.Type != eventstore.EventMasteryProjected || json.Unmarshal(ev.Payload, &previous) != nil || previous.CompetencyID != in.CompetencyID || previous.Dimension != in.Dimension || previous.EvidenceID != in.EvidenceID || previous.Variant != in.Variant || previous.HelpUsed != in.HelpUsed || previous.SolutionRevealed != in.SolutionRevealed || previous.Success != in.Success {
			return EvidenceResult{}, eventstore.ErrRevisionConflict
		}
		provenance := previous.ContentProvenance
		if provenance == "" {
			provenance = "legacy_unreviewed"
			signals, err := p.recordedSignals()
			if err != nil {
				return EvidenceResult{}, err
			}
			state := mastery.FoldProjections(signals)[in.CompetencyID][dim].State
			if state == "" {
				state = mastery.StateNotObserved
			}
			previous.State = string(state)
		}
		return EvidenceResult{ContentProvenance: provenance, CompetencyID: previous.CompetencyID, Dimension: previous.Dimension, State: previous.State, Revision: p.Revision()}, nil
	}

	observed := time.Now().UTC()
	provenance := p.evidenceProvenance(in.EvidenceID)
	payload := signalPayload{ContentProvenance: provenance,
		CompetencyID: in.CompetencyID, Dimension: string(dim), EvidenceID: in.EvidenceID, Variant: in.Variant,
		HelpUsed: in.HelpUsed, SolutionRevealed: in.SolutionRevealed, Success: in.Success,
		ObservedAt: observed.Format(time.RFC3339Nano), RuleVersion: mastery.RuleVersionV1,
	}

	priorSignals, err := p.recordedSignals()
	if err != nil {
		return EvidenceResult{}, err
	}
	var schedulePtr *mastery.Schedule
	if priorSchedule, ok := mastery.FoldSchedules(priorSignals)[in.CompetencyID]; ok {
		schedulePtr = &priorSchedule
	}

	priorProjection := mastery.FoldProjections(priorSignals)[in.CompetencyID][dim]
	sig := mastery.Signal{ContentProvenance: provenance,
		CompetencyID: in.CompetencyID, Dimension: dim, EvidenceID: in.EvidenceID, Variant: in.Variant,
		HelpUsed: in.HelpUsed, SolutionRevealed: in.SolutionRevealed, Success: in.Success, Observed: observed,
	}
	resultProjection := mastery.AdvanceState(mastery.RuleVersionV1, priorProjection, sig)
	payload.State = string(resultProjection.State)

	ev, err := p.store.Append(masteryStreamID, in.ExpectedRevision, in.RequestID, eventstore.EventMasteryProjected, payload)
	if err != nil {
		return EvidenceResult{}, err
	}

	if !in.SolutionRevealed && provenance == "published" {
		newSchedule := mastery.NextSchedule(in.CompetencyID, schedulePtr, in.Success, observed)
		schedulePayloadValue := schedulePayload{
			CompetencyID: newSchedule.CompetencyID, DueAt: newSchedule.DueAt.Format(time.RFC3339Nano),
			IntervalDays: newSchedule.IntervalDays, Reason: newSchedule.Reason,
		}
		scheduleRequestID := ""
		if in.RequestID != "" {
			scheduleRequestID = in.RequestID + ":schedule"
		}
		if _, err := p.store.Append(masteryStreamID, ev.Revision, scheduleRequestID, eventstore.EventReviewScheduled, schedulePayloadValue); err != nil {
			return EvidenceResult{}, err
		}
	}

	return EvidenceResult{ContentProvenance: provenance,
		CompetencyID: in.CompetencyID, Dimension: string(dim), State: payload.State,
		Revision: p.store.Revision(masteryStreamID),
	}, nil
}

// ProgressResult is what progress_get reports: one competency's state
// across every dimension it has evidence for.
type ProgressResult struct {
	Competencies map[string]map[string]mastery.DimensionProjection
	Revision     uint64
}

// Progress recomputes and returns the full projection from the log
// (requirement R8), optionally restricted to one competency.
func (p *ProgressService) Progress(competencyID string) (ProgressResult, error) {
	signals, err := p.recordedSignals()
	if err != nil {
		return ProgressResult{}, err
	}
	byCompetency := mastery.FoldProjections(signals)
	out := map[string]map[string]mastery.DimensionProjection{}
	for id, dims := range byCompetency {
		if competencyID != "" && id != competencyID {
			continue
		}
		byDim := map[string]mastery.DimensionProjection{}
		for dim, proj := range dims {
			byDim[string(dim)] = proj
		}
		out[id] = byDim
	}
	return ProgressResult{Competencies: out, Revision: p.store.Revision(masteryStreamID)}, nil
}

// ReviewDueResult is what review_due reports.
type ReviewDueResult struct {
	Due      []mastery.Schedule
	Revision uint64
}

// ReviewDue recomputes every competency's schedule from the log and
// returns those due at now, most-overdue first (requirements R6, R8). It
// never touches session state: a due review is a recommendation, never a
// gate on starting a free session.
func (p *ProgressService) ReviewDue(now time.Time) (ReviewDueResult, error) {
	signals, err := p.recordedSignals()
	if err != nil {
		return ReviewDueResult{}, err
	}
	schedules := mastery.FoldSchedules(signals)
	return ReviewDueResult{Due: mastery.DueSchedules(schedules, now), Revision: p.store.Revision(masteryStreamID)}, nil
}

// Revision returns the mastery stream's current revision, so a caller can
// obtain one before its first RecordEvidence call without recomputing the
// full projection.
func (p *ProgressService) Revision() uint64 { return p.store.Revision(masteryStreamID) }

// recordedSignals replays every mastery_projected event on the mastery
// stream, in order, decoded back into mastery.RecordedSignal.
func (p *ProgressService) recordedSignals() ([]mastery.RecordedSignal, error) {
	events := p.store.Replay(masteryStreamID)
	out := make([]mastery.RecordedSignal, 0, len(events))
	for _, ev := range events {
		if ev.Type != eventstore.EventMasteryProjected {
			continue
		}
		var payload signalPayload
		if err := json.Unmarshal(ev.Payload, &payload); err != nil {
			return nil, fmt.Errorf("decoding mastery_projected event %s: %w", ev.ID, err)
		}
		observed, err := time.Parse(time.RFC3339Nano, payload.ObservedAt)
		if err != nil {
			return nil, fmt.Errorf("parsing observed_at for event %s: %w", ev.ID, err)
		}
		out = append(out, mastery.RecordedSignal{
			Signal: mastery.Signal{ContentProvenance: payload.ContentProvenance,
				CompetencyID: payload.CompetencyID, Dimension: mastery.Dimension(payload.Dimension),
				EvidenceID: payload.EvidenceID, Variant: payload.Variant, HelpUsed: payload.HelpUsed,
				SolutionRevealed: payload.SolutionRevealed, Success: payload.Success, Observed: observed,
			},
			RuleVersion: payload.RuleVersion,
		})
	}
	return out, nil
}
