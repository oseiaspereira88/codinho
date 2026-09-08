package main

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
)

const treePack = `schema_version: 1
id: tree-pack
version: 1.0.0
challenges:
  - schema_version: 1
    id: tree
    version: 1.0.0
    title: Synthetic tree
    brief: Practice one window at a time.
    layers:
      - id: layer-one
        macro_steps:
          - id: macro-one
            kind: macro
            instruction: {objective: First group, scope: First layer}
            children:
              - id: meso-one
                kind: meso
                instruction: {objective: First task, scope: First group}
                children:
                  - id: a1
                    kind: micro
                    instruction: {objective: First action, scope: Only a1}
                  - id: a2
                    kind: micro
                    instruction: {objective: Second action, scope: Only a2}
      - id: layer-two
        macro_steps:
          - id: macro-two
            kind: macro
            instruction: {objective: Second group, scope: Second layer}
            children:
              - id: choose
                kind: meso
                children_mode: choice
                instruction: {objective: Select an approach, scope: One alternative}
                children:
                  - id: branch-a
                    kind: meso
                    children:
                      - id: b1
                        kind: micro
                        instruction: {objective: Approach A, scope: Only b1}
                  - id: branch-b
                    kind: meso
                    children:
                      - id: c1
                        kind: micro
                        instruction: {objective: Approach B, scope: Only c1}
              - id: last
                kind: micro
                instruction: {objective: Final action, scope: Only last}
`

func TestTreeProgressionOverRealStdio(t *testing.T) {
	bin := buildCodinhoBinary(t)
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "packs"), 0700); err != nil {
		t.Fatal(err)
	}
	for name, data := range map[string]string{"manifest.yaml": "schema_version: 1\npacks: [pack.yaml]\n", "pack.yaml": treePack} {
		if err := os.WriteFile(filepath.Join(root, "packs", name), []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	connect := func() *mcp.ClientSession {
		cmd := exec.CommandContext(ctx, bin, "serve", "--authoring")
		cmd.Dir = root
		cs, err := mcp.NewClient(&mcp.Implementation{Name: "tree-progression-test", Version: "1"}, nil).Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
		if err != nil {
			t.Fatal(err)
		}
		return cs
	}
	cs := connect()
	defer func() { cs.Close() }()
	node := func(env map[string]any, id, kind string) {
		t.Helper()
		n, ok := env["active_node"].(map[string]any)
		if !ok || n["id"] != id || n["kind"] != kind {
			t.Fatalf("want %s/%s: %+v", id, kind, env)
		}
	}
	for depth, id := range map[string]string{"challenge": "tree", "layer": "layer-one", "macro": "macro-one", "meso": "meso-one", "micro": "a1"} {
		e := callTool(ctx, t, cs, "session_start", map[string]any{"challenge_id": "tree", "depth": depth})
		node(e, id, depth)
		args := map[string]any{"session_id": e["session_id"]}
		node(callTool(ctx, t, cs, "session_get", args), id, depth)
		instruction := callTool(ctx, t, cs, "instruction_get", args)
		node(instruction, id, depth)
		raw, _ := json.Marshal(instruction)
		if strings.Contains(string(raw), "children") || strings.Contains(string(raw), "branch-b") {
			t.Fatal("instruction leaked tree")
		}
	}
	progressBefore := callTool(ctx, t, cs, "progress_get", map[string]any{})
	log := filepath.Join(root, ".codinho", "state", "events.jsonl")
	readEvents := func() []eventstore.Event {
		t.Helper()
		r := eventstore.ReadEvents(log, nil)
		if r.Err != nil {
			t.Fatal(r.Err)
		}
		return r.Events
	}
	reject := func(name string, args map[string]any, code string) {
		t.Helper()
		before := len(readEvents())
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: args})
		if err != nil {
			t.Fatal(err)
		}
		raw, _ := json.Marshal(res.StructuredContent)
		if !res.IsError || !strings.Contains(string(raw), code) {
			t.Fatalf("expected %s: %s", code, raw)
		}
		if strings.Contains(string(raw), root) {
			t.Fatal("error leaked path")
		}
		if len(readEvents()) != before {
			t.Fatal("rejection appended")
		}
	}
	for _, choice := range []string{"branch-a", "branch-b"} {
		start := callTool(ctx, t, cs, "session_start", map[string]any{"challenge_id": "tree", "depth": "micro", "request_id": "start-" + choice})
		sid := start["session_id"].(string)
		rev := revisionOf(t, start)
		mutate := func(name string, args map[string]any) map[string]any {
			t.Helper()
			if args == nil {
				args = map[string]any{}
			}
			args["session_id"] = sid
			args["expected_revision"] = rev
			r := callTool(ctx, t, cs, name, args)
			rev = revisionOf(t, r)
			return r
		}
		reject("step_advance", map[string]any{"session_id": sid, "expected_revision": rev}, "INVALID_INPUT")
		reject("step_advance", map[string]any{"session_id": sid, "expected_revision": rev, "override": true, "next_step_id": "last"}, "INVALID_INPUT")
		feedback := mutate("feedback_record", map[string]any{"type": "question", "text": "How can this be checked?"})
		if feedback["progress_effect"] != "feedback_recorded" {
			t.Fatal(feedback)
		}
		node(callTool(ctx, t, cs, "session_get", map[string]any{"session_id": sid}), "a1", "micro")
		e := mutate("step_evaluate", map[string]any{"submission_intent": true, "request_id": "evaluation-" + choice})
		if e["progress_effect"] != "attempt_recorded" {
			t.Fatal(e)
		}
		node(callTool(ctx, t, cs, "session_get", map[string]any{"session_id": sid}), "a1", "micro")
		for _, d := range []string{"macro", "layer", "challenge", "micro"} {
			mutate("granularity_adjust", map[string]any{"depth": d})
		}
		// Completion succeeds after window changes without a second evaluation.
		completeArgs := map[string]any{"session_id": sid, "expected_revision": rev, "request_id": "complete-" + choice}
		complete := callTool(ctx, t, cs, "step_complete", completeArgs)
		rev = revisionOf(t, complete)
		duplicate := mutate("step_complete", map[string]any{"override": true, "request_id": "duplicate"})
		if rev != revisionOf(t, complete) {
			t.Fatal("duplicate completion changed revision")
		}
		_ = duplicate
		node(mutate("step_advance", nil), "a2", "micro")
		for _, d := range []string{"challenge", "layer", "micro"} {
			mutate("granularity_adjust", map[string]any{"depth": d})
		}
		node(callTool(ctx, t, cs, "instruction_get", map[string]any{"session_id": sid}), "a2", "micro")
		if choice == "branch-a" {
			if err := cs.Close(); err != nil {
				t.Fatal(err)
			}
			cs = connect()
			node(callTool(ctx, t, cs, "session_get", map[string]any{"session_id": sid}), "a2", "micro")
			retry := callTool(ctx, t, cs, "step_complete", completeArgs)
			if revisionOf(t, retry) != revisionOf(t, complete) {
				t.Fatal("historical complete retry changed")
			}
		}
		finishNode := func() {
			t.Helper()
			mutate("step_evaluate", nil)
			c := mutate("step_complete", nil)
			if c["progress_effect"] != "step_completed" {
				t.Fatal(c)
			}
		}
		finishNode()
		node(mutate("step_advance", nil), "choose", "meso")
		finishNode()
		options := mutate("step_advance", nil)
		opts := envData(t, options)["branches"].([]any)
		if len(opts) != 2 || options["progress_effect"] != "none" {
			t.Fatal(options)
		}
		reject("step_advance", map[string]any{"session_id": sid, "expected_revision": rev, "next_step_id": "b1"}, "INVALID_INPUT")
		reject("step_advance", map[string]any{"session_id": sid, "expected_revision": rev - 1, "next_step_id": choice}, "STATE_CONFLICT")
		chosenArgs := map[string]any{"session_id": sid, "expected_revision": rev, "request_id": "choose-" + choice, "next_step_id": choice}
		selected := callTool(ctx, t, cs, "step_advance", chosenArgs)
		rev = revisionOf(t, selected)
		want := "b1"
		if choice == "branch-b" {
			want = "c1"
		}
		node(selected, want, "micro")
		if choice == "branch-a" {
			if err := cs.Close(); err != nil {
				t.Fatal(err)
			}
			cs = connect()
			node(callTool(ctx, t, cs, "step_advance", chosenArgs), want, "micro")
		}
		for _, d := range []string{"layer", "meso", "micro"} {
			mutate("granularity_adjust", map[string]any{"depth": d})
		}
		node(callTool(ctx, t, cs, "instruction_get", map[string]any{"session_id": sid}), want, "micro")
		finishNode()
		node(mutate("step_advance", nil), "last", "micro")
		finishNode()
		end := mutate("step_advance", map[string]any{"request_id": "end"})
		if envData(t, end)["done"] != true {
			t.Fatal(end)
		}
		if envData(t, callTool(ctx, t, cs, "session_get", map[string]any{"session_id": sid}))["state"] != "active" {
			t.Fatal("advance implicitly finished session")
		}
		counts := map[string]int{}
		attempts := 0
		for _, ev := range readEvents() {
			if ev.StreamID != sid {
				continue
			}
			if ev.Type == eventstore.EventStepCompleted {
				var p struct {
					StepID string `json:"step_id"`
				}
				if err := json.Unmarshal(ev.Payload, &p); err != nil {
					t.Fatal(err)
				}
				counts[p.StepID]++
			}
			if ev.Type == eventstore.EventAttemptSubmitted {
				attempts++
			}
		}
		if len(counts) != 5 || attempts != 1 {
			t.Fatalf("completion/attempt count: %v %d", counts, attempts)
		}
		for id, n := range counts {
			if n != 1 {
				t.Fatalf("duplicate %s: %d", id, n)
			}
		}
	}
	// Changing scale in a later layer must not resurrect earlier containers.
	coarse := callTool(ctx, t, cs, "session_start", map[string]any{"challenge_id": "tree", "depth": "micro"})
	coarseID := coarse["session_id"].(string)
	coarseRev := revisionOf(t, coarse)
	for _, want := range []string{"a2", "choose"} {
		d := callTool(ctx, t, cs, "step_complete", map[string]any{"session_id": coarseID, "override": true, "expected_revision": coarseRev})
		n := callTool(ctx, t, cs, "step_advance", map[string]any{"session_id": coarseID, "expected_revision": revisionOf(t, d)})
		kind := "micro"
		if want == "choose" {
			kind = "meso"
		}
		node(n, want, kind)
		coarseRev = revisionOf(t, n)
	}
	w := callTool(ctx, t, cs, "granularity_adjust", map[string]any{"session_id": coarseID, "depth": "layer", "expected_revision": coarseRev})
	node(w, "layer-two", "layer")
	d := callTool(ctx, t, cs, "step_complete", map[string]any{"session_id": coarseID, "override": true, "expected_revision": revisionOf(t, w)})
	end := callTool(ctx, t, cs, "step_advance", map[string]any{"session_id": coarseID, "expected_revision": revisionOf(t, d)})
	if envData(t, end)["done"] != true {
		t.Fatal("coarser advance reopened an earlier layer")
	}
	progressAfter := callTool(ctx, t, cs, "progress_get", map[string]any{})
	if !reflect.DeepEqual(envData(t, progressBefore), envData(t, progressAfter)) {
		t.Fatal("navigation fabricated mastery")
	}
	// Interview policy survives both coarser and finer navigation and restart.
	interview := callTool(ctx, t, cs, "session_start", map[string]any{"challenge_id": "tree", "mode": "interview", "depth": "micro"})
	sid := interview["session_id"].(string)
	rev := revisionOf(t, interview)
	for _, d := range []string{"layer", "micro"} {
		x := callTool(ctx, t, cs, "granularity_adjust", map[string]any{"session_id": sid, "depth": d, "expected_revision": rev})
		rev = revisionOf(t, x)
	}
	if err := cs.Close(); err != nil {
		t.Fatal(err)
	}
	cs = connect()
	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "hint_request", Arguments: map[string]any{"session_id": sid, "expected_revision": rev, "confirm_solution": true}})
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsError {
		t.Fatal("interview disclosed a hint")
	}
}
