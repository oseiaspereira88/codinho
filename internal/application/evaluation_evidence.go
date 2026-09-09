package application

import (
	"encoding/json"

	"github.com/oseiaspereira88/codinho/internal/assessment"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/evidence"
	"github.com/oseiaspereira88/codinho/internal/session"
	"github.com/oseiaspereira88/codinho/internal/workspace"
)

type evaluationEvidenceValidator struct {
	store    *eventstore.Store
	evidence *evidence.Store
}

type evaluationEvidencePayload struct {
	Kind        string `json:"kind"`
	SessionID   string `json:"session_id"`
	StepID      string `json:"step_id"`
	ChallengeID string `json:"challenge_id"`
	EvidenceID  string `json:"evidence_id"`
	CheckID     string `json:"check_id"`
	Outcome     string `json:"outcome"`
	Fingerprint string `json:"fingerprint"`
	Source      string `json:"source"`
	Text        string `json:"text"`
	RubricRef   string `json:"rubric_ref"`
}

func (v *evaluationEvidenceValidator) Register(ctx session.EvaluationEvidenceContext, in session.QualitativeEvidenceInput) (string, error) {
	if v.evidence == nil || string(workspace.Redact([]byte(in.RubricRef))) != in.RubricRef {
		return "", session.ErrEvaluationEvidenceInvalid
	}
	p := evaluationEvidencePayload{Kind: "qualitative", SessionID: string(ctx.SessionID), StepID: string(ctx.StepID), ChallengeID: ctx.ChallengeID, Source: in.Source, Text: string(workspace.Redact([]byte(in.Text))), RubricRef: in.RubricRef}
	raw, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return v.evidence.Put(raw)
}

func (v *evaluationEvidenceValidator) Validate(ctx session.EvaluationEvidenceContext, c assessment.CriterionInput) (session.EvidenceLineage, error) {
	var empty session.EvidenceLineage
	if v.evidence == nil {
		return empty, session.ErrEvaluationEvidenceInvalid
	}
	raw, err := v.evidence.Get(string(c.EvidenceID))
	if err != nil {
		return empty, session.ErrEvaluationEvidenceInvalid
	}
	var blob evaluationEvidencePayload
	if json.Unmarshal(raw, &blob) != nil {
		return empty, session.ErrEvaluationEvidenceInvalid
	}
	var scope observationEvent
	var producer eventstore.Event
	var event evaluationEvidencePayload
	var producerScope observationEvent
	trackCursor := 0
	for _, ev := range v.store.Replay(string(ctx.SessionID)) {
		if ev.Type == eventstore.EventStepAdvanced {
			var boundary struct {
				Cursor int `json:"track_cursor"`
			}
			if json.Unmarshal(ev.Payload, &boundary) == nil && boundary.Cursor > trackCursor {
				trackCursor = boundary.Cursor
				producer = eventstore.Event{}
			}
		}

		var p evaluationEvidencePayload
		if json.Unmarshal(ev.Payload, &p) != nil || p.StepID != string(ctx.StepID) {
			continue
		}
		if ev.Type == eventstore.EventObservationRecorded {
			var observed observationEvent
			if json.Unmarshal(ev.Payload, &observed) == nil && observed.RecoveryVersion == 1 {
				scope = observed
			}
		}
		if p.EvidenceID != string(c.EvidenceID) {
			continue
		}
		if ev.Type != eventstore.EventObservationRecorded && ev.Type != eventstore.EventCheckExecuted && ev.Type != eventstore.EventEvidenceRecorded {
			continue
		}
		producer, event, producerScope = ev, p, scope
	}
	if producer.ID == "" {
		return empty, session.ErrEvaluationEvidenceInvalid
	}
	lineage := session.EvidenceLineage{Criterion: c.Name, EvidenceID: string(c.EvidenceID), EventID: producer.ID, StepID: string(ctx.StepID), ChallengeID: ctx.ChallengeID, Kind: blob.Kind, RubricRef: c.RubricRef}
	if producer.Type == eventstore.EventEvidenceRecorded {
		if blob.Kind != "qualitative" || blob.SessionID != string(ctx.SessionID) || blob.StepID != string(ctx.StepID) || blob.ChallengeID != ctx.ChallengeID || event.ChallengeID != ctx.ChallengeID || blob.Source != event.Source || blob.RubricRef != event.RubricRef || c.CheckID != "" {
			return empty, session.ErrEvaluationEvidenceInvalid
		}
		if c.Kind != "structural" && c.RubricRef != blob.RubricRef {
			return empty, session.ErrEvaluationEvidenceInvalid
		}
		lineage.Source = blob.Source
		return lineage, nil
	}
	if producer.Type == eventstore.EventCheckExecuted {
		if blob.Kind != "check" || blob.CheckID == "" || event.CheckID != blob.CheckID || event.Outcome != blob.Outcome {
			return empty, session.ErrEvaluationEvidenceInvalid
		}
		known := false
		for _, check := range ctx.Checks {
			if check.ID == blob.CheckID {
				known = true
				break
			}
		}
		if !known || (c.Kind == "structural" && c.CheckID != blob.CheckID) || (c.CheckID != "" && c.CheckID != blob.CheckID) {
			return empty, session.ErrEvaluationEvidenceInvalid
		}
		lineage.CheckID, lineage.Outcome = blob.CheckID, blob.Outcome
	} else if (blob.Kind != "baseline" && blob.Kind != "diff") || blob.Kind != event.Kind || c.CheckID != "" {
		return empty, session.ErrEvaluationEvidenceInvalid
	}
	if producerScope.Root == "" || blob.Fingerprint == "" || blob.Fingerprint != event.Fingerprint {
		return empty, session.ErrEvaluationEvidenceInvalid
	}
	root, err := workspace.AuthorizeRoot(producerScope.Root)
	if err != nil || root.Path() != producerScope.Root {
		return empty, session.ErrEvaluationEvidenceStale
	}
	fingerprint, err := workspace.Fingerprint(root, producerScope.Globs)
	if err != nil || fingerprint != blob.Fingerprint {
		return empty, session.ErrEvaluationEvidenceStale
	}
	lineage.Fingerprint = fingerprint
	return lineage, nil
}
