package application

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/oseiaspereira88/codinho/internal/assessment"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/evidence"
	"github.com/oseiaspereira88/codinho/internal/learning"
	"github.com/oseiaspereira88/codinho/internal/session"
)

func TestEvaluationEvidenceRejectsCrossSessionAndDrift(t *testing.T) {
	for _, scenario := range []string{"cross-session", "drift"} {
		t.Run(scenario, func(t *testing.T) {
			sessions, ws, checks := newChecksTestServices(t)
			first, err := sessions.Start(StartInput{ChallengeID: "fixture.challenge-one"})
			if err != nil {
				t.Fatal(err)
			}
			root := newFixtureGoModule(t)
			obs, err := ws.Observe(ObserveInput{SessionID: first.SessionID, StepID: first.ActiveStep, Root: root, ExpectedRevision: first.Revision})
			if err != nil {
				t.Fatal(err)
			}
			run, err := checks.Run(CheckRunInput{SessionID: first.SessionID, CheckID: "focused-tests", ExpectedRevision: obs.Revision})
			if err != nil {
				t.Fatal(err)
			}
			target, revision := first.SessionID, run.Revision
			if scenario == "cross-session" {
				second, err := sessions.Start(StartInput{ChallengeID: "fixture.challenge-one"})
				if err != nil {
					t.Fatal(err)
				}
				target, revision = second.SessionID, second.Revision
			} else if err := os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\nfunc main() { panic(1) }\n"), 0600); err != nil {
				t.Fatal(err)
			}
			before := len(ws.store.ReplayAll())
			_, err = sessions.StepEvaluate(EvaluateInput{SessionID: target, ExpectedRevision: revision, SubmissionIntent: true,
				Criteria: []CriterionInput{{Name: "tests-pass", CheckID: "focused-tests", Kind: "structural", Severity: learning.SeverityBlocking, EvidenceID: learning.EvidenceID(run.EvidenceID)}}})
			if err == nil {
				t.Fatal("unauthorized evidence accepted")
			}
			if len(ws.store.ReplayAll()) != before {
				t.Fatal("rejection appended evaluation or attempt")
			}
		})
	}
}

func checkedEvidence(t *testing.T) (*SessionService, *WorkspaceService, StartResult, CheckRunResult, string) {
	t.Helper()
	s, ws, checks := newChecksTestServices(t)
	start, err := s.Start(StartInput{ChallengeID: "fixture.challenge-one"})
	if err != nil {
		t.Fatal(err)
	}
	root := newFixtureGoModule(t)
	obs, err := ws.Observe(ObserveInput{SessionID: start.SessionID, StepID: start.ActiveStep, Root: root, ExpectedRevision: start.Revision})
	if err != nil {
		t.Fatal(err)
	}
	run, err := checks.Run(CheckRunInput{SessionID: start.SessionID, CheckID: "focused-tests", ExpectedRevision: obs.Revision})
	if err != nil || run.Outcome != "pass" {
		t.Fatalf("check: %+v %v", run, err)
	}
	return s, ws, start, run, root
}

func evaluationInput(start StartResult, run CheckRunResult) EvaluateInput {
	return EvaluateInput{SessionID: start.SessionID, ExpectedRevision: run.Revision, RequestID: "evaluate", SubmissionIntent: true, Criteria: []CriterionInput{{Name: "tests-pass", Kind: "structural", Severity: learning.SeverityBlocking, CheckID: "focused-tests", EvidenceID: learning.EvidenceID(run.EvidenceID)}}}
}

func TestEvaluationEvidenceRejectsInvalidScopeAndBlobs(t *testing.T) {
	for _, scenario := range []string{"unknown", "unregistered", "missing-blob", "corrupt-blob", "wrong-check", "missing-check", "missing-root", "replaced-root", "wrong-step", "unrecoverable-baseline"} {
		t.Run(scenario, func(t *testing.T) {
			s, ws, start, run, root := checkedEvidence(t)
			in := evaluationInput(start, run)
			switch scenario {
			case "unknown":
				in.Criteria[0].EvidenceID = "unknown"
			case "unregistered":
				id, err := ws.evidence.Put([]byte(`{"kind":"check","check_id":"focused-tests","outcome":"pass"}`))
				if err != nil {
					t.Fatal(err)
				}
				in.Criteria[0].EvidenceID = learning.EvidenceID(id)
			case "missing-blob", "corrupt-blob":
				raw, err := ws.evidence.Get(run.EvidenceID)
				if err != nil {
					t.Fatal(err)
				}
				dir := t.TempDir()
				replacement, err := evidence.Open(dir)
				if err != nil {
					t.Fatal(err)
				}
				if _, err := replacement.Put(raw); err != nil {
					t.Fatal(err)
				}
				ws.evidence = replacement
				s.svc = session.New(nil, ws.store, replacement, &evaluationEvidenceValidator{ws.store, replacement})
				path := filepath.Join(dir, run.EvidenceID)
				if scenario == "missing-blob" {
					if err := os.Remove(path); err != nil {
						t.Fatal(err)
					}
				} else {
					if err := os.Chmod(path, 0600); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(path, []byte("tampered"), 0600); err != nil {
						t.Fatal(err)
					}
				}
			case "wrong-check":
				in.Criteria[0].CheckID = "different-check"
			case "missing-check":
				in.Criteria[0].CheckID = ""
			case "missing-root":
				if err := os.Rename(root, root+"-gone"); err != nil {
					t.Fatal(err)
				}
			case "replaced-root":
				if err := os.Rename(root, root+"-original"); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(root+"-original", root); err != nil {
					t.Fatal(err)
				}
			case "wrong-step", "unrecoverable-baseline":
				// Rebuild a synthetic log with a producer attached to another node,
				// or a legacy observation without a recoverable authorization scope.
				store, err := eventstore.Open(filepath.Join(t.TempDir(), "events.jsonl"), nil)
				if err != nil {
					t.Fatal(err)
				}
				defer store.Close()
				for _, ev := range ws.store.ReplayAll() {
					var p map[string]any
					json.Unmarshal(ev.Payload, &p)
					if scenario == "wrong-step" && ev.Type == eventstore.EventCheckExecuted {
						p["step_id"] = "another-step"
					}
					if scenario == "unrecoverable-baseline" && ev.Type == eventstore.EventObservationRecorded {
						delete(p, "recovery_version")
					}
					if _, err := store.Append(ev.StreamID, ev.Revision-1, ev.RequestID, ev.Type, p); err != nil {
						t.Fatal(err)
					}
				}
				ws.store = store
				s.svc = session.New(nil, store, ws.evidence, &evaluationEvidenceValidator{store, ws.evidence})
			}
			before := len(ws.store.ReplayAll())
			_, err := s.StepEvaluate(in)
			if !errors.Is(err, ErrEvaluationEvidenceInvalid) && !errors.Is(err, ErrEvaluationEvidenceStale) {
				t.Fatalf("rejection: %v", err)
			}
			if len(ws.store.ReplayAll()) != before {
				t.Fatal("rejection appended")
			}
		})
	}
}

func TestEvaluationEvidenceQualitativeRegistrationAndRecovery(t *testing.T) {
	s, ws, start, run, root := checkedEvidence(t)
	note := QualitativeEvidenceInput{SessionID: start.SessionID, Source: "learner_explanation", Text: "I use a pointer to preserve identity.", RubricRef: "rubric://reasoning", ExpectedRevision: run.Revision, RequestID: "register"}
	registered, err := s.RecordQualitativeEvidence(note)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ws.EvidenceGet(EvidenceGetInput{SessionID: start.SessionID, EvidenceID: registered.EvidenceID})
	if err != nil || !strings.Contains(string(got.Content), note.Text) {
		t.Fatalf("registered read: %s %v", got.Content, err)
	}
	in := EvaluateInput{SessionID: start.SessionID, ExpectedRevision: registered.Revision, RequestID: "qualitative", Criteria: []CriterionInput{{Name: "reasoning", Kind: "clarity", Severity: learning.SeverityBlocking, Verdict: learning.VerdictMet, EvidenceID: learning.EvidenceID(registered.EvidenceID), RubricRef: note.RubricRef}}}
	wrong := in
	wrong.Criteria = append([]CriterionInput(nil), in.Criteria...)
	wrong.Criteria[0].RubricRef = "rubric://other"
	if _, err := s.StepEvaluate(wrong); !errors.Is(err, ErrEvaluationEvidenceInvalid) {
		t.Fatalf("rubric mismatch: %v", err)
	}
	approved, err := s.StepEvaluate(in)
	if err != nil || approved.Criteria[0].Verdict != learning.VerdictMet {
		t.Fatalf("qualitative: %+v %v", approved, err)
	}
	if err := os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\nfunc main(){panic(1)}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	s.svc = session.New(nil, ws.store, ws.evidence, &evaluationEvidenceValidator{ws.store, ws.evidence})
	recoveredWS := NewWorkspaceService(ws.store, ws.evidence)
	read, err := recoveredWS.EvidenceGet(EvidenceGetInput{SessionID: start.SessionID, EvidenceID: registered.EvidenceID})
	if err != nil || read.Stale {
		t.Fatalf("note changed freshness at restart: %+v %v", read, err)
	}
	retry, err := s.StepEvaluate(in)
	if err != nil || !reflect.DeepEqual(retry, approved) {
		t.Fatalf("evaluation retry: %+v %v", retry, err)
	}
	repeated, err := s.RecordQualitativeEvidence(note)
	if err != nil || repeated != registered {
		t.Fatalf("registration retry: %+v %v", repeated, err)
	}
	structural := in
	structural.RequestID = "note-as-check"
	structural.ExpectedRevision = approved.Revision
	structural.Criteria = append([]CriterionInput(nil), in.Criteria...)
	structural.Criteria[0].Kind = "structural"
	unverifiable, err := s.StepEvaluate(structural)
	if err != nil || unverifiable.Criteria[0].Verdict != learning.VerdictUnverifiable {
		t.Fatalf("note must not prove a check: %+v %v", unverifiable, err)
	}
}

type driftEvidenceValidator struct {
	inner *evaluationEvidenceValidator
	file  string
	calls int
}

func (v *driftEvidenceValidator) Register(c session.EvaluationEvidenceContext, in session.QualitativeEvidenceInput) (string, error) {
	return v.inner.Register(c, in)
}
func (v *driftEvidenceValidator) Validate(c session.EvaluationEvidenceContext, in assessment.CriterionInput) (session.EvidenceLineage, error) {
	proof, err := v.inner.Validate(c, in)
	v.calls++
	if err == nil && v.calls == 1 {
		err = os.WriteFile(v.file, []byte("package main\nfunc main(){panic(1)}\n"), 0600)
	}
	return proof, err
}

func TestEvaluationEvidenceRechecksImmediatelyBeforeAppend(t *testing.T) {
	s, ws, start, run, root := checkedEvidence(t)
	v := &driftEvidenceValidator{inner: &evaluationEvidenceValidator{ws.store, ws.evidence}, file: filepath.Join(root, "main.go")}
	s.svc = session.New(nil, ws.store, ws.evidence, v)
	_, err := s.StepEvaluate(evaluationInput(start, run))
	if !errors.Is(err, ErrEvaluationEvidenceStale) {
		t.Fatalf("drift between validation and append: %v", err)
	}
	if ws.store.Revision(string(start.SessionID)) != run.Revision {
		t.Fatal("drift appended")
	}
}

func TestEvaluationEvidenceRegistrationValidationAndRedaction(t *testing.T) {
	s, ws, start, run, _ := checkedEvidence(t)
	valid := QualitativeEvidenceInput{SessionID: start.SessionID, Source: "tutor_observation", Text: "Observed evidence", RubricRef: "rubric://clarity", ExpectedRevision: run.Revision, RequestID: "note"}
	for _, bad := range []string{"source", "empty-text", "empty-rubric", "large-text", "secret-rubric"} {
		in := valid
		switch bad {
		case "source":
			in.Source = "untrusted-source"
		case "empty-text":
			in.Text = " "
		case "empty-rubric":
			in.RubricRef = ""
		case "large-text":
			in.Text = strings.Repeat("a", (64<<10)+1)
		case "secret-rubric":
			in.RubricRef = "password=hunter2hunter2hunter2"
		}
		if _, err := s.RecordQualitativeEvidence(in); !errors.Is(err, ErrEvaluationEvidenceInvalid) {
			t.Fatalf("%s accepted: %v", bad, err)
		}
		if ws.store.Revision(string(start.SessionID)) != run.Revision {
			t.Fatal("invalid registration appended")
		}
	}
	valid.Text = "Explanation; password=hunter2hunter2hunter2"
	record, err := s.RecordQualitativeEvidence(valid)
	if err != nil {
		t.Fatal(err)
	}
	data, err := ws.EvidenceGet(EvidenceGetInput{SessionID: start.SessionID, EvidenceID: record.EvidenceID})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data.Content), "hunter2hunter2hunter2") {
		t.Fatal("registration leaked synthetic secret")
	}
	valid.Text = "Different observation"
	if _, err := s.RecordQualitativeEvidence(valid); !errors.Is(err, eventstore.ErrRevisionConflict) {
		t.Fatalf("incompatible retry: %v", err)
	}
}

func TestEvaluationEvidenceRejectsDriftDuringCheck(t *testing.T) {
	s, ws, checks := newChecksTestServices(t)
	start, err := s.Start(StartInput{ChallengeID: "fixture.challenge-one"})
	if err != nil {
		t.Fatal(err)
	}
	root := newFixtureGoModule(t)
	// Synthetic trusted fixture: demonstrate that a passing test can change
	// the files whose pre-execution state it was supposed to establish.
	code := "package main\nimport (\"testing\";\"os\")\nfunc TestMutates(t *testing.T){if err:=os.WriteFile(\"main.go\",[]byte(\"package main\\nfunc main(){panic(1)}\\n\"),0600);err!=nil{t.Fatal(err)}}\n"
	if err := os.WriteFile(filepath.Join(root, "main_test.go"), []byte(code), 0600); err != nil {
		t.Fatal(err)
	}
	obs, err := ws.Observe(ObserveInput{SessionID: start.SessionID, StepID: start.ActiveStep, Root: root, ExpectedRevision: start.Revision})
	if err != nil {
		t.Fatal(err)
	}
	_, err = checks.Run(CheckRunInput{SessionID: start.SessionID, CheckID: "focused-tests", ExpectedRevision: obs.Revision})
	if !errors.Is(err, ErrEvaluationEvidenceStale) {
		t.Fatalf("passing mutation accepted: %v", err)
	}
	if ws.store.Revision(string(start.SessionID)) != obs.Revision {
		t.Fatal("drifting check appended")
	}
}

func TestEvaluationEvidenceCheckOutcomesAndObservationLimits(t *testing.T) {
	s, ws, start, run, _ := checkedEvidence(t)
	for outcome, want := range map[string]learning.EvaluationVerdict{"pass": learning.VerdictMet, "fail": learning.VerdictNotMet, "skipped": learning.VerdictNotApplicable, "error": learning.VerdictUnverifiable} {
		raw, err := ws.evidence.Get(run.EvidenceID)
		if err != nil {
			t.Fatal(err)
		}
		var payload checkPayload
		if err := json.Unmarshal(raw, &payload); err != nil {
			t.Fatal(err)
		}
		payload.Outcome = outcome
		raw, err = json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		id, err := ws.evidence.Put(raw)
		if err != nil {
			t.Fatal(err)
		}
		producer, err := ws.store.Append(string(start.SessionID), ws.store.Revision(string(start.SessionID)), "", eventstore.EventCheckExecuted, map[string]any{"step_id": start.ActiveStep, "check_id": "focused-tests", "outcome": outcome, "evidence_id": id, "fingerprint": payload.Fingerprint})
		if err != nil {
			t.Fatal(err)
		}
		in := evaluationInput(start, CheckRunResult{EvidenceID: id, Revision: producer.Revision})
		in.RequestID = "evaluate-" + outcome
		got, err := s.StepEvaluate(in)
		if err != nil || got.Criteria[0].Verdict != want {
			t.Fatalf("%s: %+v %v", outcome, got, err)
		}
	}
	var id string
	for _, ev := range ws.store.Replay(string(start.SessionID)) {
		if ev.Type == eventstore.EventObservationRecorded {
			var p observationEvent
			json.Unmarshal(ev.Payload, &p)
			id = p.EvidenceID
			break
		}
	}
	in := evaluationInput(start, CheckRunResult{EvidenceID: id, Revision: ws.store.Revision(string(start.SessionID))})
	in.RequestID = "observation"
	in.Criteria[0].CheckID = ""
	got, err := s.StepEvaluate(in)
	if err != nil || got.Criteria[0].Verdict != learning.VerdictUnverifiable {
		t.Fatalf("baseline proved arbitrary structural claim: %+v %v", got, err)
	}
}

func TestTrackBoundaryRejectsPreviousChallengeCheck(t *testing.T) {
	s, ws, start, run, _ := checkedEvidence(t)
	checks, step, err := s.ActiveChecks(start.SessionID)
	if err != nil {
		t.Fatal(err)
	}
	v := &evaluationEvidenceValidator{store: ws.store, evidence: ws.evidence}
	ctx := session.EvaluationEvidenceContext{SessionID: start.SessionID, StepID: step, ChallengeID: "fixture.challenge-one", Checks: checks}
	criterion := assessment.CriterionInput{Name: "test", Kind: "structural", CheckID: "focused-tests", EvidenceID: learning.EvidenceID(run.EvidenceID)}
	if _, err := v.Validate(ctx, criterion); err != nil {
		t.Fatal(err)
	}
	// The durable boundary separates challenges even if step/check IDs coincide.
	if _, err := ws.store.Append(string(start.SessionID), run.Revision, "boundary", eventstore.EventStepAdvanced, map[string]any{"track_cursor": 1}); err != nil {
		t.Fatal(err)
	}
	ctx.ChallengeID = "second"
	if _, err := v.Validate(ctx, criterion); !errors.Is(err, session.ErrEvaluationEvidenceInvalid) {
		t.Fatalf("old evidence accepted: %v", err)
	}
}
