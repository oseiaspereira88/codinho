package learning

// FeedbackType classifies consultative feedback (PROJECT.md §8.7). Feedback
// never mutates step state regardless of type (invariant 3) — Blocking only
// describes whether it should stop a completion request from *elsewhere*.
type FeedbackType string

const (
	FeedbackConfirmation FeedbackType = "confirmation"
	FeedbackQuestion     FeedbackType = "question"
	FeedbackExplanation  FeedbackType = "explanation"
	FeedbackRelation     FeedbackType = "relation"
	FeedbackSuggestion   FeedbackType = "suggestion"
	FeedbackIdiom        FeedbackType = "idiom"
	FeedbackRisk         FeedbackType = "risk"
	FeedbackViolation    FeedbackType = "violation"
	FeedbackError        FeedbackType = "error"
)

// blocksCompletionByDefault mirrors the "Bloqueia conclusão?" column of
// PROJECT.md §8.7. Idiom and Risk are context-dependent ("normalmente não" /
// "conforme critério") and default to non-blocking; a caller with more
// context can override via FeedbackRecord.Blocking.
func (t FeedbackType) blocksCompletionByDefault() bool {
	switch t {
	case FeedbackViolation, FeedbackError:
		return true
	default:
		return false
	}
}

func (t FeedbackType) valid() bool {
	switch t {
	case FeedbackConfirmation, FeedbackQuestion, FeedbackExplanation, FeedbackRelation,
		FeedbackSuggestion, FeedbackIdiom, FeedbackRisk, FeedbackViolation, FeedbackError:
		return true
	}
	return false
}

// FeedbackRecord is consultative: it describes observations, strengths,
// risks, suggestions and questions, approves nothing, consumes no attempt,
// and never completes or advances a step (PROJECT.md §8.6 invariant 3).
type FeedbackRecord struct {
	StepID   StepID
	Type     FeedbackType
	Text     string
	Blocking bool
}

// NewFeedbackRecord validates Type and applies the table default for
// Blocking unless overridden.
func NewFeedbackRecord(stepID StepID, feedbackType FeedbackType, text string, blockingOverride *bool) (FeedbackRecord, error) {
	if !feedbackType.valid() {
		return FeedbackRecord{}, newDomainError(ErrCodeInvalidValue, "feedback type: "+string(feedbackType))
	}
	if text == "" {
		return FeedbackRecord{}, newDomainError(ErrCodeInvalidValue, "feedback text is required")
	}
	blocking := feedbackType.blocksCompletionByDefault()
	if blockingOverride != nil {
		blocking = *blockingOverride
	}
	return FeedbackRecord{StepID: stepID, Type: feedbackType, Text: text, Blocking: blocking}, nil
}

// EvaluationVerdict is the outcome of comparing evidence with criteria
// (PROJECT.md §8.6 "Avaliação", §21.2).
type EvaluationVerdict string

const (
	VerdictMet           EvaluationVerdict = "met"
	VerdictPartiallyMet  EvaluationVerdict = "partially_met"
	VerdictNotMet        EvaluationVerdict = "not_met"
	VerdictUnverifiable  EvaluationVerdict = "unverifiable"
	VerdictNotApplicable EvaluationVerdict = "not_applicable"
)

func (v EvaluationVerdict) valid() bool {
	switch v {
	case VerdictMet, VerdictPartiallyMet, VerdictNotMet, VerdictUnverifiable, VerdictNotApplicable:
		return true
	}
	return false
}

// FindingSeverity classifies an achado (finding) surfaced while evaluating
// a criterion (PROJECT.md §21.2).
type FindingSeverity string

const (
	SeverityBlocking             FindingSeverity = "blocking"
	SeverityImportantNonBlocking FindingSeverity = "important_non_blocking"
	SeverityAdvisory             FindingSeverity = "advisory"
)

func (s FindingSeverity) valid() bool {
	switch s {
	case SeverityBlocking, SeverityImportantNonBlocking, SeverityAdvisory:
		return true
	}
	return false
}

// CriterionResult judges one evaluation criterion. EvidenceID and
// RubricRef are required for every qualitative judgment (requirement R4,
// PROJECT.md §21.1: "Cada julgamento qualitativo deve citar uma evidência
// observada e uma rubrica"); Kind == "structural" is the one exemption,
// since its verdict is derived by the caller from evidence presence
// alone, not semantic judgment.
type CriterionResult struct {
	Name       string            `json:"name"`
	Kind       string            `json:"kind"`
	Severity   FindingSeverity   `json:"severity"`
	Verdict    EvaluationVerdict `json:"verdict"`
	EvidenceID EvidenceID        `json:"evidence_id,omitempty"`
	RubricRef  string            `json:"rubric_ref,omitempty"`
}

// StructuralCriterionKind marks a criterion whose verdict this system
// derives deterministically rather than from semantic judgment (PROJECT.md
// §21.1 "Determinística pelo servidor").
const StructuralCriterionKind = "structural"

// NewCriterionResult validates severity and verdict, and enforces that
// every non-structural (qualitative) judgment cites both evidence and a
// rubric (requirement R4).
func NewCriterionResult(name, kind string, severity FindingSeverity, verdict EvaluationVerdict, evidenceID EvidenceID, rubricRef string) (CriterionResult, error) {
	if name == "" {
		return CriterionResult{}, newDomainError(ErrCodeInvalidValue, "criterion name is required")
	}
	if !severity.valid() {
		return CriterionResult{}, newDomainError(ErrCodeInvalidValue, "finding severity: "+string(severity))
	}
	if !verdict.valid() {
		return CriterionResult{}, newDomainError(ErrCodeInvalidValue, "evaluation verdict: "+string(verdict))
	}
	if kind != StructuralCriterionKind && (evidenceID == "" || rubricRef == "") {
		return CriterionResult{}, newDomainError(ErrCodeQualitativeJudgmentRequiresEvidence, name)
	}
	return CriterionResult{Name: name, Kind: kind, Severity: severity, Verdict: verdict, EvidenceID: evidenceID, RubricRef: rubricRef}, nil
}

// Evaluation compares evidence with criteria. It never completes or
// advances a step by itself (PROJECT.md §12.3 invariant 4); those are
// separate, explicit operations on StepProgress/LearningSession.
type Evaluation struct {
	StepID    StepID
	Criteria  []CriterionResult
	AttemptID *EvidenceID // set only when the learner deliberately submitted
}

// HasBlockingFailure reports whether any blocking-severity criterion was
// not met (not_applicable never blocks: the criterion simply does not
// apply here).
func (e Evaluation) HasBlockingFailure() bool {
	for _, c := range e.Criteria {
		if c.Severity == SeverityBlocking && c.Verdict != VerdictMet && c.Verdict != VerdictNotApplicable {
			return true
		}
	}
	return false
}
