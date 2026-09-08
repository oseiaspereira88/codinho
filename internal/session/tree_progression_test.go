package session

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oseiaspereira88/codinho/internal/assistance"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/learning"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"reflect"
	"testing"
)

func TestTreeStartResolvesDepth(t *testing.T) {
	for _, tc := range []struct {
		depth learning.Depth
		id    string
	}{{learning.DepthChallenge, fixtureChallengeID}, {learning.DepthLayer, "understanding"}, {learning.DepthMacro, "fixture.macro-one"}, {learning.DepthMeso, "fixture.meso-one"}, {learning.DepthMicro, "fixture.micro-one"}} {
		t.Run(string(tc.depth), func(t *testing.T) {
			svc := newTestService(t)
			start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID, Depth: tc.depth})
			if err != nil {
				t.Fatal(err)
			}
			if string(start.ActiveStep) != tc.id {
				t.Fatalf("got %s, want %s", start.ActiveStep, tc.id)
			}
			if _, err := svc.Instruction(start.SessionID); err != nil {
				t.Fatal(err)
			}
		})
	}
}

// TestTreePreservesProgressAcrossWindowsAndRestart exercises the actual event
// reducer, including failed appends and the clean-evaluation/hint flags.
func TestTreePreservesProgressAcrossWindowsAndRestart(t *testing.T) {
	svc := newTreeService(t, false)
	start, err := svc.Start(StartInput{ChallengeID: "tree", Depth: learning.DepthMicro, Help: learning.HelpProgressive, DisclosureMax: learning.DisclosureSolution})
	mustTree(t, err)
	sid := start.SessionID
	rev := start.Revision
	if start.ActiveStep != "a1" {
		t.Fatal(start)
	}
	h, err := svc.HintRequest(sid, false, rev, "hint")
	mustTree(t, err)
	rev = h.Revision
	eval, err := svc.StepEvaluate(EvaluateInput{SessionID: sid, ExpectedRevision: rev, RequestID: "eval", SubmissionIntent: true})
	mustTree(t, err)
	rev = eval.Revision
	// Failed revision must not change cursor, progress or clean evaluation.
	before := svc.sessions[sid].clone()
	_, err = svc.GranularityAdjust(sid, learning.DepthLayer, "", rev+1, "stale")
	if err == nil || !reflect.DeepEqual(before, svc.sessions[sid]) {
		t.Fatalf("failed adjustment mutated state: %v", err)
	}
	for _, d := range []learning.Depth{learning.DepthMacro, learning.DepthLayer, learning.DepthChallenge, learning.DepthMicro} {
		r, err := svc.GranularityAdjust(sid, d, "keep position", rev, string(d))
		mustTree(t, err)
		rev = r.Revision
	}
	rec := svc.sessions[sid]
	p := rec.session.ActiveStep()
	if p.StepID != "a1" || p.State != learning.StepStateEvaluated || p.HintLevel != 1 || !rec.cleanEvaluation {
		t.Fatalf("progress reset: %+v clean=%v", p, rec.cleanEvaluation)
	}
	svc = New(svc.catalog, svc.store, nil)
	rec = svc.sessions[sid]
	if rec == nil {
		t.Fatal("recovery failed")
	}
	p = rec.session.ActiveStep()
	if p.State != learning.StepStateEvaluated || p.HintLevel != 1 || !rec.cleanEvaluation {
		t.Fatalf("recovery reset: %+v", p)
	}
	done, err := svc.StepComplete(sid, true, false, rev, "complete-a1")
	mustTree(t, err)
	rev = done.Revision
	duplicate, err := svc.StepComplete(sid, true, true, rev, "different-complete")
	mustTree(t, err)
	if duplicate.Revision != rev {
		t.Fatal("duplicate completion appended")
	}
	next, err := svc.StepAdvance(sid, false, rev, "next")
	mustTree(t, err)
	rev = next.Revision
	if next.StepID != "a2" {
		t.Fatal(next)
	}
	for _, d := range []learning.Depth{learning.DepthLayer, learning.DepthChallenge, learning.DepthMeso, learning.DepthMicro} {
		r, err := svc.GranularityAdjust(sid, d, "keep a2", rev, "")
		mustTree(t, err)
		rev = r.Revision
	}
	if svc.sessions[sid].session.ActiveStep().StepID != "a2" {
		t.Fatal("restarted first branch")
	}
	if svc.sessions[sid].progress["a1"].Step.State != learning.StepStateCompleted {
		t.Fatal("lost completed child")
	}
	if retry, err := svc.StepEvaluate(EvaluateInput{SessionID: sid, ExpectedRevision: h.Revision, RequestID: "eval", SubmissionIntent: true}); err != nil || !reflect.DeepEqual(retry, eval) {
		t.Fatalf("historical retry: %+v %v", retry, err)
	}
	counts := map[eventstore.EventType]int{}
	for _, e := range svc.store.Replay(string(sid)) {
		counts[e.Type]++
	}
	if counts[eventstore.EventStepCompleted] != 1 || counts[eventstore.EventAttemptSubmitted] != 1 {
		t.Fatal(counts)
	}
}

func TestTreeSequencesChoicesAndNegativeNavigation(t *testing.T) {
	for _, choose := range []string{"branch-a", "branch-b"} {
		t.Run(choose, func(t *testing.T) {
			svc := newTreeService(t, true)
			start, err := svc.Start(StartInput{ChallengeID: "tree", Depth: learning.DepthMicro})
			mustTree(t, err)
			sid := start.SessionID
			rev := start.Revision
			reject := func(f func() error) {
				t.Helper()
				before := svc.sessions[sid].clone()
				count := len(svc.store.Replay(string(sid)))
				if err := f(); err == nil {
					t.Fatal("expected rejection")
				}
				if count != len(svc.store.Replay(string(sid))) || !reflect.DeepEqual(before, svc.sessions[sid]) {
					t.Fatal("rejection mutated state")
				}
			}
			reject(func() error { _, e := svc.StepAdvance(sid, false, rev, ""); return e })
			reject(func() error { _, e := svc.StepAdvance(sid, true, rev, "", "outside"); return e })
			complete := func() {
				t.Helper()
				e, err := svc.StepEvaluate(EvaluateInput{SessionID: sid, ExpectedRevision: rev})
				mustTree(t, err)
				d, err := svc.StepComplete(sid, true, false, e.Revision, "")
				mustTree(t, err)
				rev = d.Revision
			}
			advance := func(want string) {
				t.Helper()
				n, err := svc.StepAdvance(sid, false, rev, "")
				mustTree(t, err)
				rev = n.Revision
				if n.StepID != want || n.Done || len(n.Branches) > 0 {
					t.Fatalf("want %s: %+v", want, n)
				}
			}
			complete()
			advance("a2")
			complete()
			advance("choose")
			reject(func() error { _, e := svc.StepAdvance(sid, false, rev, ""); return e })
			complete()
			opts, err := svc.StepAdvance(sid, false, rev, "options")
			mustTree(t, err)
			if len(opts.Branches) != 2 || opts.Revision != rev {
				t.Fatal(opts)
			}
			reject(func() error { _, e := svc.StepAdvance(sid, false, rev-1, ""); return e })
			reject(func() error { _, e := svc.StepAdvance(sid, false, rev, "", "b1"); return e })
			reject(func() error { _, e := svc.StepAdvance(sid, true, rev, "", "last"); return e })
			n, err := svc.StepAdvance(sid, false, rev, "choose-request", choose)
			mustTree(t, err)
			expected := "b1"
			if choose == "branch-b" {
				expected = "c1"
			}
			if n.StepID != expected {
				t.Fatal(n)
			}
			chooseRev := rev
			rev = n.Revision
			svc = New(svc.catalog, svc.store, nil)
			retry, err := svc.StepAdvance(sid, false, chooseRev, "choose-request", choose)
			mustTree(t, err)
			if !reflect.DeepEqual(n, retry) {
				t.Fatal("choice retry changed")
			}
			reject(func() error { _, e := svc.StepAdvance(sid, false, chooseRev, "choose-request", "other"); return e })
			for _, d := range []learning.Depth{learning.DepthChallenge, learning.DepthLayer, learning.DepthMicro} {
				x, err := svc.GranularityAdjust(sid, d, "", rev, "")
				mustTree(t, err)
				rev = x.Revision
			}
			if svc.sessions[sid].session.ActiveStep().StepID != learning.StepID(expected) {
				t.Fatal("lost selected branch")
			}
			complete()
			advance("last")
			complete()
			end, err := svc.StepAdvance(sid, false, rev, "exhaust")
			mustTree(t, err)
			if !end.Done {
				t.Fatal(end)
			}
			rev = end.Revision
			state, err := svc.Get(sid)
			mustTree(t, err)
			if state.State != learning.SessionStateActive {
				t.Fatal("advance finished session")
			}
			svc = New(svc.catalog, svc.store, nil)
			again, err := svc.StepAdvance(sid, false, rev, "other-exhaust")
			mustTree(t, err)
			if !again.Done || again.Revision != rev {
				t.Fatal(again)
			}
		})
	}
}

func TestTreeCoarseCompletionAndLifecycleGuards(t *testing.T) {
	svc := newTreeService(t, false)
	start, err := svc.Start(StartInput{ChallengeID: "tree", Depth: learning.DepthLayer})
	mustTree(t, err)
	sid := start.SessionID
	d, err := svc.StepComplete(sid, false, true, start.Revision, "")
	mustTree(t, err)
	if _, err := svc.GranularityAdjust(sid, learning.DepthMicro, "", d.Revision, ""); err == nil {
		t.Fatal("refined a completed aggregate")
	}
	n, err := svc.StepAdvance(sid, false, d.Revision, "")
	mustTree(t, err)
	if n.StepID != "layer-two" || n.Kind != "layer" {
		t.Fatal(n)
	}
	detour, err := svc.DetourStart(sid, "concept", n.Revision, "")
	mustTree(t, err)
	if _, err := svc.GranularityAdjust(sid, learning.DepthMicro, "", detour.Revision, ""); err == nil {
		t.Fatal("left open detour")
	}
	finished, err := svc.DetourFinish(sid, assistance.DetourResolved, detour.Revision, "")
	mustTree(t, err)
	paused, err := svc.Pause(sid, finished.Revision, "")
	mustTree(t, err)
	if _, err := svc.StepAdvance(sid, true, paused.Revision, ""); err == nil {
		t.Fatal("advanced paused session")
	}
	resumed, err := svc.Resume(sid, paused.Revision, "")
	mustTree(t, err)
	ended, err := svc.Finish(sid, resumed.Revision, "")
	mustTree(t, err)
	if _, err := svc.GranularityAdjust(sid, learning.DepthMicro, "", ended.Revision, ""); err == nil {
		t.Fatal("navigated finished session")
	}
	count := 0
	for _, e := range svc.store.Replay(string(sid)) {
		if e.Type == eventstore.EventStepCompleted {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("synthetic descendant completions: %d", count)
	}
}

func TestTreeLegacyStartAndNavigationRetry(t *testing.T) {
	svc := newTestService(t)
	ch, _ := svc.catalog.Challenge(fixtureChallengeID)
	policy, err := learning.NewSessionPolicy(learning.ModePractice, learning.DepthMicro, learning.HelpProgressive, learning.DisclosureSolution, learning.EvaluationOnDemand, learning.AdvanceExplicit, nil)
	mustTree(t, err)
	in := StartInput{ChallengeID: fixtureChallengeID, RequestID: "old-start"}
	_, err = svc.store.Append("ses_7", 0, "old-start", eventstore.EventSessionStarted, startedPayload{RecoveryVersion: 1, ChallengeID: ch.ID, Challenge: ch, Digest: contentDigest(ch), Input: in, Policy: policy})
	mustTree(t, err)
	// Original StepAdvance input encoding, with no optional next_step_id.
	raw, _ := json.Marshal([]any{"StepAdvance", []any{true}})
	digest := fmt.Sprintf("%x", sha256.Sum256(raw))
	_, err = svc.store.Append("ses_7", 1, "old-advance", eventstore.EventStepAdvanced, map[string]any{"from": "fixture.macro-one", "to": "fixture.meso-one", "override": true, "request_digest": digest})
	mustTree(t, err)
	svc = New(svc.catalog, svc.store, nil)
	start, err := svc.Start(in)
	mustTree(t, err)
	if start.ActiveStep != "fixture.macro-one" {
		t.Fatal("rewrote old initial position")
	}
	old, err := svc.StepAdvance("ses_7", true, 1, "old-advance")
	mustTree(t, err)
	if old.StepID != "fixture.meso-one" || old.Revision != 2 {
		t.Fatal(old)
	}
	x, err := svc.GranularityAdjust("ses_7", learning.DepthMicro, "continue", 2, "new")
	mustTree(t, err)
	if x.StepID != "fixture.micro-one" {
		t.Fatal(x)
	}
	svc = New(svc.catalog, svc.store, nil)
	old, err = svc.StepAdvance("ses_7", true, 1, "old-advance")
	mustTree(t, err)
	if old.StepID != "fixture.meso-one" {
		t.Fatal(old)
	}
}

func TestTreeUnknownNavigationVersionFailsClosed(t *testing.T) {
	svc := newTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	mustTree(t, err)
	_, err = svc.store.Append(string(start.SessionID), start.Revision, "future", eventstore.EventGranularityChanged, map[string]any{"navigation_version": 2, "depth": "macro", "step_id": "fixture.macro-one"})
	mustTree(t, err)
	svc = New(svc.catalog, svc.store, nil)
	if _, err := svc.Get(start.SessionID); !errors.Is(err, ErrSessionUnrecoverable) {
		t.Fatalf("accepted future version: %v", err)
	}
}

func mustTree(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func newTreeService(t *testing.T, choice bool) *Service {
	t.Helper()
	leaf := func(id string) curriculum.StepAuthoring {
		return curriculum.StepAuthoring{ID: id, Kind: "micro", Instruction: curriculum.InstructionAuthoring{Objective: "Do " + id, Scope: "Only " + id}, Completion: curriculum.CompletionAuthoring{RequiresUserConfirmation: true}}
	}
	first := curriculum.StepAuthoring{ID: "macro-one", Kind: "macro", Children: []curriculum.StepAuthoring{{ID: "meso-one", Kind: "meso", Children: []curriculum.StepAuthoring{leaf("a1"), leaf("a2")}}}}
	second := curriculum.StepAuthoring{ID: "macro-two", Kind: "macro", Children: []curriculum.StepAuthoring{leaf("last")}}
	if choice {
		second.Children = append([]curriculum.StepAuthoring{{ID: "choose", Kind: "meso", ChildrenMode: "choice", Instruction: curriculum.InstructionAuthoring{Objective: "Choose an approach"}, Children: []curriculum.StepAuthoring{{ID: "branch-a", Kind: "meso", Children: []curriculum.StepAuthoring{leaf("b1")}}, {ID: "branch-b", Kind: "meso", Children: []curriculum.StepAuthoring{leaf("c1")}}}}}, second.Children...)
	}
	ch := curriculum.ChallengeAuthoring{SchemaVersion: 1, ID: "tree", Version: "1.0.0", Title: "Tree", Brief: "Synthetic tree", Layers: []curriculum.LayerAuthoring{{ID: "layer-one", MacroSteps: []curriculum.StepAuthoring{first}}, {ID: "layer-two", MacroSteps: []curriculum.StepAuthoring{second}}}}
	pack := curriculum.Pack{SchemaVersion: 1, ID: "tree-pack", Version: "1.0.0", Challenges: []curriculum.ChallengeAuthoring{ch}}
	raw, err := yaml.Marshal(pack)
	mustTree(t, err)
	dir := t.TempDir()
	writeTestFile(t, dir, "manifest.yaml", "schema_version: 1\npacks: [pack.yaml]\n")
	writeTestFile(t, dir, "pack.yaml", string(raw))
	catalog, diags, err := curriculum.Load(dir, curriculum.DefaultLimits)
	mustTree(t, err)
	if len(diags) > 0 {
		t.Fatal(diags)
	}
	store, err := eventstore.Open(filepath.Join(t.TempDir(), "events.jsonl"), nil)
	mustTree(t, err)
	t.Cleanup(func() { store.Close() })
	return New(catalog, store, nil)
}

func TestTreeFailedAppendPreservesSuspendedProgress(t *testing.T) {
	svc := newTreeService(t, false)
	start, err := svc.Start(StartInput{ChallengeID: "tree", Depth: learning.DepthMicro})
	mustTree(t, err)
	ev, err := svc.StepEvaluate(EvaluateInput{SessionID: start.SessionID, ExpectedRevision: start.Revision})
	mustTree(t, err)
	before := svc.sessions[start.SessionID].clone()
	mustTree(t, svc.store.Close())
	if _, err := svc.GranularityAdjust(start.SessionID, learning.DepthLayer, "", ev.Revision, "cannot-write"); err == nil {
		t.Fatal("append unexpectedly succeeded")
	}
	if !reflect.DeepEqual(before, svc.sessions[start.SessionID]) {
		t.Fatal("failed append changed progress/cursor")
	}
}

func TestTreeOverrideAtLastWindowIsAudited(t *testing.T) {
	svc := newTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID, Depth: learning.DepthMicro})
	mustTree(t, err)
	end, err := svc.StepAdvance(start.SessionID, true, start.Revision, "skip-last")
	mustTree(t, err)
	if !end.Done || end.Revision != start.Revision+1 {
		t.Fatal(end)
	}
	events := svc.store.Replay(string(start.SessionID))
	last := events[len(events)-1]
	var p deltaPayload
	mustTree(t, json.Unmarshal(last.Payload, &p))
	if last.Type != eventstore.EventStepAdvanced || !p.Override || !p.Done {
		t.Fatal("last override unaudited")
	}
	svc = New(svc.catalog, svc.store, nil)
	retry, err := svc.StepAdvance(start.SessionID, true, start.Revision, "skip-last")
	mustTree(t, err)
	if !reflect.DeepEqual(end, retry) {
		t.Fatal(retry)
	}
	if _, err := svc.GranularityAdjust(start.SessionID, learning.DepthMacro, "", end.Revision, ""); err == nil {
		t.Fatal("resurrected exhausted tree")
	}
	if _, err := svc.StepComplete(start.SessionID, true, true, end.Revision, ""); err == nil {
		t.Fatal("manufactured completion after skip")
	}
	for _, e := range svc.store.Replay(string(start.SessionID)) {
		if e.Type == eventstore.EventStepCompleted {
			t.Fatal("override created completion")
		}
	}
}
