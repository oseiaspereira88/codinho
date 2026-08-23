package assessment

import (
	"errors"
	"testing"

	"github.com/oseiaspereira88/codinho/internal/learning"
)

func TestResolveStructuralCriterionDerivesVerdictFromEvidencePresence(t *testing.T) {
	met, err := Resolve(CriterionInput{Name: "compiles", Kind: learning.StructuralCriterionKind, Severity: learning.SeverityBlocking, EvidenceID: "ev_1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if met.Verdict != learning.VerdictMet {
		t.Fatalf("verdict = %s, want met when evidence is cited", met.Verdict)
	}

	unverifiable, err := Resolve(CriterionInput{Name: "compiles", Kind: learning.StructuralCriterionKind, Severity: learning.SeverityBlocking})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if unverifiable.Verdict != learning.VerdictUnverifiable {
		t.Fatalf("verdict = %s, want unverifiable with no evidence", unverifiable.Verdict)
	}
}

func TestResolveStructuralCriterionIgnoresCallerVerdict(t *testing.T) {
	result, err := Resolve(CriterionInput{Name: "compiles", Kind: learning.StructuralCriterionKind, Severity: learning.SeverityBlocking, Verdict: learning.VerdictNotMet, EvidenceID: "ev_1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Verdict != learning.VerdictMet {
		t.Fatalf("verdict = %s, a structural criterion must not trust the caller's verdict", result.Verdict)
	}
}

func TestResolveQualitativeCriterionRequiresEvidenceAndRubric(t *testing.T) {
	_, err := Resolve(CriterionInput{Name: "idiomatic-go", Kind: "qualitative", Severity: learning.SeverityAdvisory, Verdict: learning.VerdictPartiallyMet})
	if !errors.Is(err, learning.DomainError{Code: learning.ErrCodeQualitativeJudgmentRequiresEvidence}) {
		t.Fatalf("expected ErrCodeQualitativeJudgmentRequiresEvidence, got %v", err)
	}

	result, err := Resolve(CriterionInput{
		Name: "idiomatic-go", Kind: "qualitative", Severity: learning.SeverityAdvisory,
		Verdict: learning.VerdictPartiallyMet, EvidenceID: "ev_1", RubricRef: "rubric://idiomatic-go",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Verdict != learning.VerdictPartiallyMet {
		t.Fatalf("verdict = %s, a qualitative criterion must keep the caller's verdict", result.Verdict)
	}
}
