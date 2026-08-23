package learning

import "time"

// EvidenceID identifies a piece of evidence with a stable ID.
type EvidenceID CatalogID

// EvidenceKind is the shape a piece of evidence takes (PROJECT.md §7.3).
type EvidenceKind string

const (
	EvidenceKindDiff        EvidenceKind = "diff"
	EvidenceKindTest        EvidenceKind = "test"
	EvidenceKindBuild       EvidenceKind = "build"
	EvidenceKindAnswer      EvidenceKind = "answer"
	EvidenceKindExplanation EvidenceKind = "explanation"
	EvidenceKindBenchmark   EvidenceKind = "benchmark"
	EvidenceKindAnalysis    EvidenceKind = "analysis"
)

// Evidence is immutable once recorded: a correction produces a new Evidence
// rather than mutating this one (PROJECT.md §12.3 invariant 9). The zero
// value is never valid; construct via NewEvidence.
type Evidence struct {
	id   EvidenceID
	kind EvidenceKind
	ref  string // opaque pointer to the underlying artifact (diff hash, log id, ...)
}

// NewEvidence validates and constructs an immutable Evidence.
func NewEvidence(id EvidenceID, kind EvidenceKind, ref string) (Evidence, error) {
	if id == "" {
		return Evidence{}, newDomainError(ErrCodeInvalidValue, "evidence id is required")
	}
	if ref == "" {
		return Evidence{}, newDomainError(ErrCodeInvalidValue, "evidence ref is required")
	}
	return Evidence{id: id, kind: kind, ref: ref}, nil
}

func (e Evidence) ID() EvidenceID     { return e.id }
func (e Evidence) Kind() EvidenceKind { return e.kind }
func (e Evidence) Ref() string        { return e.ref }

// Observation is evidence observed without implying a deliberate attempt
// (PROJECT.md §7.3).
type Observation struct {
	StepID   StepID
	Evidence Evidence
	Observed time.Time
}

// Attempt is a deliberate submission for evaluation. SolutionRevealed marks
// that the step's solution was revealed at or before this attempt, which
// invalidates it as evidence of autonomy without discarding it as evidence
// of correctness (PROJECT.md §12.3 invariant 8).
type Attempt struct {
	StepID           StepID
	Evidence         []Evidence
	SolutionRevealed bool
}

// CountsAsAutonomyEvidence reports whether this attempt may be used to
// project autonomous mastery.
func (a Attempt) CountsAsAutonomyEvidence() bool {
	return !a.SolutionRevealed
}

// HintUsage records a hint granted at a given disclosure level. Construction
// fails if the level exceeds the session's disclosure policy (PROJECT.md
// §12.3 invariant 7).
type HintUsage struct {
	StepID StepID
	Level  DisclosureLevel
}

// NewHintUsage validates the requested level against the policy cap.
func NewHintUsage(stepID StepID, level DisclosureLevel, policy DisclosurePolicy) (HintUsage, error) {
	if !level.valid() {
		return HintUsage{}, newDomainError(ErrCodeInvalidValue, "disclosure level out of range")
	}
	if level > policy.MaxLevel {
		return HintUsage{}, newDomainError(ErrCodeDisclosureExceeded, "level exceeds session policy")
	}
	return HintUsage{StepID: stepID, Level: level}, nil
}

// Reflection is a short, non-interrogative response demonstrating reasoning
// (PROJECT.md §8.9). Its quality is evidence separate from implementation.
// Assessment is the tutor's descriptive (non-blocking, non-scoring) read of
// that quality; it never completes or advances a step (feedback-evaluation-
// progression requirement R9).
type Reflection struct {
	StepID       StepID
	CompetencyID CompetencyID
	Prompt       string
	Answer       string
	Assessment   string
}

// MasteryEvidence is one data point feeding a longitudinal mastery
// projection. This package records the point; it does not compute the
// projection itself (non-goal: "calcular domínio longitudinal").
type MasteryEvidence struct {
	CompetencyID CompetencyID
	Evidence     Evidence
	Autonomous   bool
}

// ReviewSchedule marks that a competency is due for spaced review. The
// scheduling policy itself belongs to mastery-review-scheduling; this type
// only carries the fact.
type ReviewSchedule struct {
	CompetencyID CompetencyID
	DueAt        time.Time
}
