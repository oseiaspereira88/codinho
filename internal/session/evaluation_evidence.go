package session

import (
	"errors"
	"strings"

	"github.com/oseiaspereira88/codinho/internal/assessment"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

var (
	ErrEvaluationEvidenceInvalid = errors.New("session: evaluation evidence is unavailable or out of scope")
	ErrEvaluationEvidenceStale   = errors.New("session: evaluation evidence no longer matches the workspace")
)

type EvaluationEvidenceContext struct {
	SessionID   learning.SessionID
	StepID      learning.StepID
	ChallengeID string
	Checks      []curriculum.CheckAuthoring
}

// EvidenceLineage records the persisted producer actually authorized for a
// criterion. Outcome is used only for deterministic check verdicts.
type EvidenceLineage struct {
	Criterion   string `json:"criterion"`
	EvidenceID  string `json:"evidence_id"`
	EventID     string `json:"event_id"`
	StepID      string `json:"step_id"`
	ChallengeID string `json:"challenge_id"`
	Kind        string `json:"kind"`
	CheckID     string `json:"check_id,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	Source      string `json:"source,omitempty"`
	RubricRef   string `json:"rubric_ref,omitempty"`
	Outcome     string `json:"outcome,omitempty"`
}

// EvaluationEvidenceValidator is supplied by the application composition.
// The session reducer and historical retries do not read the workspace.
type EvaluationEvidenceValidator interface {
	Validate(EvaluationEvidenceContext, assessment.CriterionInput) (EvidenceLineage, error)
	Register(EvaluationEvidenceContext, QualitativeEvidenceInput) (string, error)
}

type QualitativeEvidenceInput struct {
	SessionID        learning.SessionID
	Source           string
	Text             string
	RubricRef        string
	ExpectedRevision uint64
	RequestID        string
}

type EvidenceRecordResult struct {
	EvidenceID string
	Revision   uint64
}

func evidenceContext(id learning.SessionID, rec *record) EvaluationEvidenceContext {
	return EvaluationEvidenceContext{SessionID: id, StepID: rec.session.ActiveStep().StepID, ChallengeID: rec.challengeID, Checks: rec.pinned.Checks}
}

func (s *Service) RecordQualitativeEvidence(in QualitativeEvidenceInput) (EvidenceRecordResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var prior EvidenceRecordResult
	if found, err := s.retryRequest(in.SessionID, in.RequestID, "RecordQualitativeEvidence", in, &prior); found || err != nil {
		return prior, err
	}
	rec, _, _, err := s.activeWindow(in.SessionID)
	if err != nil {
		return prior, err
	}
	if s.evidenceValidator == nil || strings.TrimSpace(in.Text) == "" || strings.TrimSpace(in.RubricRef) == "" || len(in.Text) > 64<<10 || len(in.RubricRef) > 1024 {
		return prior, ErrEvaluationEvidenceInvalid
	}
	switch in.Source {
	case "learner_explanation", "tutor_observation", "external_artifact":
	default:
		return prior, ErrEvaluationEvidenceInvalid
	}
	ctx := evidenceContext(in.SessionID, rec)
	id, err := s.evidenceValidator.Register(ctx, in)
	if err != nil {
		return prior, err
	}
	ev, _, err := s.mutate(in.SessionID, in.ExpectedRevision, in.RequestID, eventstore.EventEvidenceRecorded, map[string]any{
		"step_id": ctx.StepID, "challenge_id": ctx.ChallengeID, "evidence_id": id, "kind": "qualitative", "source": in.Source, "rubric_ref": in.RubricRef,
	})
	if err != nil {
		return prior, err
	}
	return EvidenceRecordResult{EvidenceID: id, Revision: ev.Revision}, nil
}
