package session

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/oseiaspereira88/codinho/internal/assessment"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/evidence"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

// newServiceWithEvidence supplies a real blob store without an authorization
// adapter, so tests prove storage alone never grants consumption authority.
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

func TestStepEvaluateWithoutValidatorRejectsCitations(t *testing.T) {
	for _, kind := range []string{"check", "diff"} {
		t.Run(kind, func(t *testing.T) {
			svc, evStore := newServiceWithEvidence(t)
			start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
			if err != nil {
				t.Fatal(err)
			}
			id, err := evStore.Put([]byte(`{"kind":"` + kind + `","outcome":"pass"}`))
			if err != nil {
				t.Fatal(err)
			}
			_, err = svc.StepEvaluate(EvaluateInput{SessionID: start.SessionID, ExpectedRevision: start.Revision, Criteria: []assessment.CriterionInput{{Name: "tests", Kind: "structural", Severity: learning.SeverityBlocking, EvidenceID: learning.EvidenceID(id)}}})
			if !errors.Is(err, ErrEvaluationEvidenceInvalid) {
				t.Fatalf("missing validator: %v", err)
			}
			if svc.store.Revision(string(start.SessionID)) != start.Revision {
				t.Fatal("rejected citation appended")
			}
		})
	}
}

func TestEvaluationEvidenceLegacyRetryPreservesInputIdentity(t *testing.T) {
	svc := newTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatal(err)
	}
	// Exact old JSON shape: the CheckID extension must remain omitted for old retries.
	type oldCriterion struct {
		Name, Kind string
		Severity   learning.FindingSeverity
		Verdict    learning.EvaluationVerdict
		EvidenceID learning.EvidenceID
		RubricRef  string
	}
	type oldInput struct {
		SessionID        learning.SessionID
		Criteria         []oldCriterion
		SubmissionIntent bool
		ExpectedRevision uint64
		RequestID        string
	}
	legacy := oldInput{SessionID: start.SessionID, Criteria: []oldCriterion{{Name: "legacy", Kind: "structural", Severity: learning.SeverityBlocking, EvidenceID: "historical-id"}}, ExpectedRevision: start.Revision, RequestID: "legacy-evaluation"}
	raw, _ := json.Marshal([]any{"StepEvaluate", legacy})
	digest := fmt.Sprintf("%x", sha256.Sum256(raw))
	result := learning.CriterionResult{Name: "legacy", Kind: "structural", Severity: learning.SeverityBlocking, Verdict: learning.VerdictMet, EvidenceID: "historical-id"}
	ev, err := svc.store.Append(string(start.SessionID), start.Revision, legacy.RequestID, eventstore.EventEvaluationRecorded, map[string]any{"step_id": start.ActiveStep, "criteria": []learning.CriterionResult{result}, "request_digest": digest})
	if err != nil {
		t.Fatal(err)
	}
	recovered := New(nil, svc.store, nil)
	input := EvaluateInput{SessionID: start.SessionID, Criteria: []assessment.CriterionInput{{Name: "legacy", Kind: "structural", Severity: learning.SeverityBlocking, EvidenceID: "historical-id"}}, ExpectedRevision: start.Revision, RequestID: legacy.RequestID}
	got, err := recovered.StepEvaluate(input)
	if err != nil || got.Revision != ev.Revision || len(got.Criteria) != 1 || got.Criteria[0] != result {
		t.Fatalf("legacy retry: %+v %v", got, err)
	}
	input.RequestID = "new-evaluation"
	input.ExpectedRevision = ev.Revision
	if _, err := recovered.StepEvaluate(input); !errors.Is(err, ErrEvaluationEvidenceInvalid) {
		t.Fatalf("new legacy citation accepted: %v", err)
	}
}
