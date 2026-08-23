// Package mastery projects multidimensional domain per competency and
// schedules spaced review by transparent, versioned rules (PROJECT.md
// §21.3, §21.4). It never computes an opaque score: every state change
// traces back to a cited Signal, and every rule application is versioned
// so a later rule change never rewrites past projections (Compatibility).
package mastery

import "time"

// Dimension is one of the eight independent axes a competency is observed
// on (PROJECT.md §21.3). They are never inferred from one another.
type Dimension string

const (
	DimensionUnderstanding            Dimension = "understanding"
	DimensionSyntaxRecall             Dimension = "syntax_recall"
	DimensionGuidedImplementation     Dimension = "guided_implementation"
	DimensionAutonomousImplementation Dimension = "autonomous_implementation"
	DimensionDebugging                Dimension = "debugging"
	DimensionExplanation              Dimension = "explanation"
	DimensionRetention                Dimension = "retention"
	DimensionTransfer                 Dimension = "transfer"
)

// Valid reports whether d is one of the eight declared dimensions.
func (d Dimension) Valid() bool {
	switch d {
	case DimensionUnderstanding, DimensionSyntaxRecall, DimensionGuidedImplementation,
		DimensionAutonomousImplementation, DimensionDebugging, DimensionExplanation,
		DimensionRetention, DimensionTransfer:
		return true
	}
	return false
}

// State is one of the six descriptive states a Dimension can be in
// (PROJECT.md §21.3). States form a strict ladder: reaching one implies
// every earlier one was already reached (Decision 4).
type State string

const (
	StateNotObserved             State = "not_observed"
	StateIntroduced              State = "introduced"
	StateDemonstratesWithHelp    State = "demonstrates_with_help"
	StateDemonstratesWithoutHelp State = "demonstrates_without_help"
	StateRetained                State = "retained"
	StateTransferred             State = "transferred"
)

// rank orders State for comparison; higher is further along the ladder.
var rank = map[State]int{
	StateNotObserved:             0,
	StateIntroduced:              1,
	StateDemonstratesWithHelp:    2,
	StateDemonstratesWithoutHelp: 3,
	StateRetained:                4,
	StateTransferred:             5,
}

// Signal is one piece of mastery evidence a caller (mastery_evidence_
// record) asserts, citing an already-recorded Evidence rather than
// creating a new one (Decision 2; invariant 9: evidence is immutable).
type Signal struct {
	CompetencyID     string
	Dimension        Dimension
	EvidenceID       string
	Variant          string
	HelpUsed         bool
	SolutionRevealed bool
	Success          bool
	Observed         time.Time
}

// DimensionProjection is one competency's current state in one dimension,
// plus the minimum audit trail AdvanceState needs to decide retention and
// transfer promotions.
type DimensionProjection struct {
	State          State           `json:"state"`
	RuleVersion    int             `json:"rule_version"`
	LastEvidenceID string          `json:"last_evidence_id,omitempty"`
	LastObserved   time.Time       `json:"last_observed,omitzero"`
	SeenVariants   map[string]bool `json:"seen_variants,omitempty"`
	EvidenceCount  int             `json:"evidence_count"`
}

// RuleVersion identifies the promotion table a projection was computed
// with (Compatibility: rule changes must preserve historical projections
// by version). V1 is the only version this package implements.
const RuleVersionV1 = 1

// Schedule is one competency's next spaced-review date (requirement R5).
// Scheduling is per competency, not per dimension: a competency is due
// for review as a whole once its strongest evidence ages out.
type Schedule struct {
	CompetencyID string    `json:"competency_id"`
	DueAt        time.Time `json:"due_at"`
	IntervalDays int       `json:"interval_days"`
	Reason       string    `json:"reason"`
}
