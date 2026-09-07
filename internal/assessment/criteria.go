package assessment

import "github.com/oseiaspereira88/codinho/internal/learning"

// CriterionInput is one criterion judgment step_evaluate receives from the
// caller (the tutor agent).
type CriterionInput struct {
	Name       string
	Kind       string
	Severity   learning.FindingSeverity
	Verdict    learning.EvaluationVerdict // ignored for structural criteria; the server derives it instead
	EvidenceID learning.EvidenceID
	RubricRef  string
	CheckID    string `json:"CheckID,omitempty"`
}

// Resolve turns one CriterionInput into a domain CriterionResult. This
// server does not run or interpret checks (non-goal: "executar checks ou
// observar workspace" belongs to safe-check-executor), so a structural
// criterion's verdict is derived deterministically from evidence presence
// alone — met when evidence is cited, unverifiable otherwise — rather than
// trusting a caller-supplied verdict it cannot itself confirm. Any other
// kind is a qualitative judgment: the caller's verdict is kept, and
// learning.NewCriterionResult enforces that it cites both evidence and a
// rubric (requirement R4).
func Resolve(in CriterionInput) (learning.CriterionResult, error) {
	verdict := in.Verdict
	if in.Kind == learning.StructuralCriterionKind {
		if in.EvidenceID != "" {
			verdict = learning.VerdictMet
		} else {
			verdict = learning.VerdictUnverifiable
		}
	}
	return learning.NewCriterionResult(in.Name, in.Kind, in.Severity, verdict, in.EvidenceID, in.RubricRef)
}
