// Package recommendation ranks candidate challenges toward an objective by
// deterministic, decomposable factors (Decision 1) — never an LLM, never a
// hidden score. It knows nothing about curriculum.Catalog or eventstore
// directly: a caller (internal/application) normalizes catalog and mastery
// data into this package's inputs, the same boundary already established
// by internal/workspace and internal/checks.
package recommendation

// Candidate is one challenge already resolved by the caller, reduced to
// what ranking needs (requirement R4).
type Candidate struct {
	ID                string
	Title             string
	ThemeIDs          []string
	PrimaryCompetency []string
	Difficulty        string
	EstimatedMinutes  int
	Prerequisites     []string
}

// Objective is what the caller wants to work toward (requirement R4).
type Objective struct {
	CompetencyID      string
	ThemeID           string
	TimeBudgetMinutes int // 0 = no budget declared
}

// Mastery is the caller's reduced mastery-review-scheduling state for one
// competency: the strongest dimension state reached, and how many days a
// review has been overdue (0 when not due).
type Mastery struct {
	CompetencyID    string
	BestState       string // a mastery.State value, e.g. "demonstrates_without_help"
	ReviewOverdueBy int    // days; 0 or negative means not due
}

// stateRank mirrors internal/mastery's ladder without importing it, so this
// package stays decoupled from the mastery domain's own types (same
// boundary as Candidate/Objective/Mastery being plain data).
var stateRank = map[string]int{
	"":                          0,
	"not_observed":              0,
	"introduced":                1,
	"demonstrates_with_help":    2,
	"demonstrates_without_help": 3,
	"retained":                  4,
	"transferred":               5,
}

// Explanation justifies one Recommendation by evidence, gap, dependency and
// cost (requirement R5) — never a bare number.
type Explanation struct {
	Evidence       string `json:"evidence"`
	Gap            string `json:"gap"`
	Dependency     string `json:"dependency"`
	CostMinutes    int    `json:"cost_minutes"`
	RemediationFor string `json:"remediation_for,omitempty"` // non-empty when this candidate is a smaller step before RemediationFor (requirement R6)
}

// Recommendation is one ranked candidate with its explanation.
type Recommendation struct {
	ChallengeID string      `json:"challenge_id"`
	Score       int         `json:"score"`
	Explanation Explanation `json:"explanation"`
}
