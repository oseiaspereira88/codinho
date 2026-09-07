package application

import (
	"github.com/oseiaspereira88/codinho/internal/assessment"
	"github.com/oseiaspereira88/codinho/internal/learning"
	"github.com/oseiaspereira88/codinho/internal/session"
)

// Feedback, evaluation, reflection, completion and advance orchestration
// live in internal/session (state) and internal/assessment (pure
// decisions); these aliases let mcpserver depend only on application's
// stable surface (feedback-evaluation-progression).
type (
	FeedbackPacket = assessment.FeedbackPacket
	CriterionInput = assessment.CriterionInput
	EvaluateInput  = session.EvaluateInput
	EvaluateResult = session.EvaluateResult
	AdvanceOption  = session.AdvanceOption
	AdvanceResult  = session.AdvanceResult
)

// FeedbackPrepare assembles the read-only context for feedback_record's
// caller to author feedback from (requirement R1).
func (s *SessionService) FeedbackPrepare(id learning.SessionID, question string) (FeedbackPacket, error) {
	return s.svc.FeedbackPrepare(id, question)
}

// FeedbackRecord persists feedback already authored elsewhere (requirement
// R2).
func (s *SessionService) FeedbackRecord(id learning.SessionID, feedbackType learning.FeedbackType, text string, blockingOverride *bool, expectedRevision uint64, requestID string) (LifecycleResult, error) {
	return s.svc.FeedbackRecord(id, feedbackType, text, blockingOverride, expectedRevision, requestID)
}

// StepEvaluate resolves each criterion and records the evaluation,
// creating an attempt only when SubmissionIntent is true (requirement R3,
// R5, R6, R7).
func (s *SessionService) StepEvaluate(in EvaluateInput) (EvaluateResult, error) {
	return s.svc.StepEvaluate(in)
}

// ReflectionRecord persists a reflection distinct from implementation
// evidence (requirement R9).
func (s *SessionService) ReflectionRecord(id learning.SessionID, competencyID learning.CompetencyID, prompt, answer, assessmentText string, expectedRevision uint64, requestID string) (LifecycleResult, error) {
	return s.svc.ReflectionRecord(id, competencyID, prompt, answer, assessmentText, expectedRevision, requestID)
}

// StepComplete marks the active step completed when its policy is
// satisfied or override is set (requirement R7, R8).
func (s *SessionService) StepComplete(id learning.SessionID, confirm, override bool, expectedRevision uint64, requestID string) (LifecycleResult, error) {
	return s.svc.StepComplete(id, confirm, override, expectedRevision, requestID)
}

// StepAdvance activates the next permitted node, or reports branch options
// or exhaustion (requirement R7).
func (s *SessionService) StepAdvance(id learning.SessionID, override bool, expectedRevision uint64, requestID string) (AdvanceResult, error) {
	return s.svc.StepAdvance(id, override, expectedRevision, requestID)
}

// RecordQualitativeEvidence registers a cited observation without evaluating it.
func (s *SessionService) RecordQualitativeEvidence(in QualitativeEvidenceInput) (EvidenceRecordResult, error) {
	return s.svc.RecordQualitativeEvidence(in)
}

type QualitativeEvidenceInput = session.QualitativeEvidenceInput
type EvidenceRecordResult = session.EvidenceRecordResult

var ErrEvaluationEvidenceInvalid = session.ErrEvaluationEvidenceInvalid
var ErrEvaluationEvidenceStale = session.ErrEvaluationEvidenceStale
