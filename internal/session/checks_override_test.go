package session

import (
	"path/filepath"
	"testing"

	"github.com/oseiaspereira88/codinho/internal/assessment"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/evidence"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

// newServiceWithEvidence is like newTestService, but wires a real
// evidence.Store so StepEvaluate's check-outcome override (safe-check-
// executor Decision 3, requirement R10) has something to look up.
func newServiceWithEvidence(t *testing.T) (*Service, *evidence.Store) {
	t.Helper()
	dir := t.TempDir()
	writeTestFile(t, dir, "manifest.yaml", "schema_version: 1\npacks:\n  - pack.yaml\n")
	writeTestFile(t, dir, "pack.yaml", `schema_version: 1
id: fixture-pack
version: 1.0.0
challenges:
  - schema_version: 1
    id: fixture.challenge-one
    version: 1.0.0
    title: Fixture challenge
    kind: atomic
    difficulty: foundational
    layers:
      - id: understanding
        macro_steps:
          - id: fixture.step-one
            kind: micro
            title: Fixture step
            instruction:
              objective: Declare the fixture type.
              scope: Only the declaration.
`)
	catalog, diags, err := curriculum.Load(dir, curriculum.DefaultLimits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostics: %+v", diags)
	}
	store, err := eventstore.Open(filepath.Join(t.TempDir(), "events.jsonl"), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	evStore, err := evidence.Open(filepath.Join(t.TempDir(), "evidence"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return New(catalog, store, evStore), evStore
}

func putCheckEvidence(t *testing.T, ev *evidence.Store, outcome string) learning.EvidenceID {
	t.Helper()
	id, err := ev.Put([]byte(`{"kind":"check","check_id":"go_test","outcome":"` + outcome + `","fingerprint":"sha256-x"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return learning.EvidenceID(id)
}

func TestStepEvaluateStructuralCriterionUsesRealCheckOutcome(t *testing.T) {
	cases := []struct {
		outcome string
		want    learning.EvaluationVerdict
	}{
		{"pass", learning.VerdictMet},
		{"fail", learning.VerdictNotMet},
		{"skipped", learning.VerdictNotApplicable},
		{"error", learning.VerdictUnverifiable},
	}
	for _, c := range cases {
		t.Run(c.outcome, func(t *testing.T) {
			svc, evStore := newServiceWithEvidence(t)
			start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			evidenceID := putCheckEvidence(t, evStore, c.outcome)

			eval, err := svc.StepEvaluate(EvaluateInput{
				SessionID: start.SessionID,
				Criteria: []assessment.CriterionInput{
					{Name: "tests-pass", Kind: learning.StructuralCriterionKind, Severity: learning.SeverityBlocking, EvidenceID: evidenceID},
				},
				ExpectedRevision: start.Revision,
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := eval.Criteria[0].Verdict; got != c.want {
				t.Fatalf("verdict = %s, want %s", got, c.want)
			}
		})
	}
}

func TestStepEvaluateFallsBackToPresenceForNonCheckEvidence(t *testing.T) {
	svc, evStore := newServiceWithEvidence(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Evidence recognizably NOT from safe-check-executor (e.g. workspace-
	// observation-baselines' own "diff" kind): the override must leave
	// feedback-evaluation-progression's original behavior untouched.
	workspaceEvidenceID, err := evStore.Put([]byte(`{"kind":"diff","fingerprint":"sha256-y"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	eval, err := svc.StepEvaluate(EvaluateInput{
		SessionID: start.SessionID,
		Criteria: []assessment.CriterionInput{
			{Name: "tests-pass", Kind: learning.StructuralCriterionKind, Severity: learning.SeverityBlocking, EvidenceID: learning.EvidenceID(workspaceEvidenceID)},
		},
		ExpectedRevision: start.Revision,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := eval.Criteria[0].Verdict; got != learning.VerdictMet {
		t.Fatalf("verdict = %s, want met (presence-based fallback)", got)
	}
}

func TestStepEvaluateWithoutEvidenceStoreKeepsOriginalBehavior(t *testing.T) {
	// newPolicyTestService wires no evidence.Store (nil) — the override
	// must no-op entirely, matching feedback-evaluation-progression's
	// sealed behavior exactly.
	svc := newPolicyTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	eval, err := svc.StepEvaluate(EvaluateInput{
		SessionID: start.SessionID,
		Criteria: []assessment.CriterionInput{
			{Name: "tests-pass", Kind: learning.StructuralCriterionKind, Severity: learning.SeverityBlocking, EvidenceID: "sha256-doesnotexist"},
		},
		ExpectedRevision: start.Revision,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := eval.Criteria[0].Verdict; got != learning.VerdictMet {
		t.Fatalf("verdict = %s, want met (presence-based, no evidence store wired)", got)
	}
}
