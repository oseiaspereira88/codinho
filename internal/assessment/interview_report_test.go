package assessment_test

import (
	"testing"
	"time"

	"github.com/oseiaspereira88/codinho/internal/assessment"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

func TestBuildInterviewReportNeverCollapsesToASingleScore(t *testing.T) {
	evaluations := []assessment.InterviewEvaluation{
		{StepID: "step-1", Criteria: []learning.CriterionResult{
			{Name: "compiles", Kind: learning.StructuralCriterionKind, Verdict: learning.VerdictMet},
			{Name: "handles-error", Kind: learning.StructuralCriterionKind, Verdict: learning.VerdictNotMet},
		}},
	}
	report := assessment.BuildInterviewReport("go-interviews.example", 45*time.Minute, false, "explicit",
		evaluations, assessment.InterviewHintUsage{Granted: 1, Blocked: 2}, nil, []string{"error-handling"})

	if report.ChallengeID != "go-interviews.example" {
		t.Fatalf("ChallengeID = %s", report.ChallengeID)
	}
	if len(report.Gaps) != 1 || report.Gaps[0] != "step-1/handles-error" {
		t.Fatalf("Gaps = %v, want exactly the not-met criterion", report.Gaps)
	}
	if report.Hints.Granted != 1 || report.Hints.Blocked != 2 {
		t.Fatalf("Hints = %+v", report.Hints)
	}
	if report.IntegrityNote == "" {
		t.Fatal("PROJECT.md §22 requires an integrity disclosure on every report")
	}
	if len(report.Recommended) != 1 || report.Recommended[0] != "error-handling" {
		t.Fatalf("Recommended = %v, want the caller-supplied list echoed back", report.Recommended)
	}
}

func TestBuildInterviewReportGapsIgnoresMetCriteria(t *testing.T) {
	evaluations := []assessment.InterviewEvaluation{
		{StepID: "step-1", Criteria: []learning.CriterionResult{
			{Name: "compiles", Kind: learning.StructuralCriterionKind, Verdict: learning.VerdictMet},
		}},
	}
	report := assessment.BuildInterviewReport("x", time.Minute, false, "explicit", evaluations, assessment.InterviewHintUsage{}, nil, nil)
	if len(report.Gaps) != 0 {
		t.Fatalf("Gaps = %v, want none when everything is met", report.Gaps)
	}
}

func TestBuildInterviewReportRecordsTimeoutAndReason(t *testing.T) {
	report := assessment.BuildInterviewReport("x", time.Hour, true, "timeout", nil, assessment.InterviewHintUsage{}, nil, nil)
	if !report.TimedOut {
		t.Fatal("TimedOut should be true")
	}
	if report.FinishReason != "timeout" {
		t.Fatalf("FinishReason = %s, want timeout", report.FinishReason)
	}
}
