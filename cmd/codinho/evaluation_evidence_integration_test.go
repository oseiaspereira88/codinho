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

func TestEvaluationEvidenceOverRealStdio(t *testing.T) {
	bin := buildCodinhoBinary(t)
	root, learner := t.TempDir(), t.TempDir()
	writeIntegrationFixturePack(t, root)
	for name, content := range map[string]string{"go.mod": "module fixture\n\ngo 1.25\n", "main.go": "package main\nfunc main(){}\n", "main_test.go": "package main\nimport \"testing\"\nfunc TestOK(t *testing.T){}\n"} {
		if err := os.WriteFile(filepath.Join(learner, name), []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	connect := func() *mcp.ClientSession {
		cmd := exec.CommandContext(ctx, bin, "serve")
		cmd.Dir = root
		cs, err := mcp.NewClient(&mcp.Implementation{Name: "evaluation-lineage-test", Version: "1"}, nil).Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
		if err != nil {
			t.Fatal(err)
		}
		return cs
	}
	first := connect()
	defer first.Close()
	start := callTool(ctx, t, first, "session_start", map[string]any{"challenge_id": "fixture.challenge-one"})
	sid := start["session_id"].(string)
	step := start["active_node"].(map[string]any)["id"].(string)
	other := callTool(ctx, t, first, "session_start", map[string]any{"challenge_id": "fixture.challenge-one"})
	otherID := other["session_id"].(string)
	obs := callTool(ctx, t, first, "workspace_observe", map[string]any{"session_id": sid, "step_id": step, "root": learner, "expected_revision": revisionOf(t, start)})
	run := callTool(ctx, t, first, "check_run", map[string]any{"session_id": sid, "check_id": "focused-tests", "expected_revision": revisionOf(t, obs)})
	eid := envData(t, run)["evidence_id"].(string)
	evalArgs := func(session string, revision float64, id, check, request string) map[string]any {
		return map[string]any{"session_id": session, "expected_revision": revision, "request_id": request, "submission_intent": true, "criteria": []map[string]any{{"name": "tests-pass", "kind": "structural", "severity": "blocking", "check_id": check, "evidence_id": id}}}
	}
	logPath := filepath.Join(root, ".codinho", "state", "events.jsonl")
	reject := func(cs *mcp.ClientSession, args map[string]any, code string) {
		t.Helper()
		before := eventstore.ReadEvents(logPath, nil)
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "step_evaluate", Arguments: args})
		if err != nil {
			t.Fatal(err)
		}
		raw, _ := json.Marshal(res.StructuredContent)
		if !res.IsError || !strings.Contains(string(raw), code) {
			t.Fatalf("expected %s: %s", code, raw)
		}
		if strings.Contains(string(raw), learner) || strings.Contains(string(raw), root) {
			t.Fatalf("error leaked path: %s", raw)
		}
		after := eventstore.ReadEvents(logPath, nil)
		if before.Err != nil || after.Err != nil || len(before.Events) != len(after.Events) {
			t.Fatal("rejected evaluation appended an event")
		}
	}
	reject(first, evalArgs(otherID, revisionOf(t, other), eid, "focused-tests", "cross-session"), "EVALUATION_EVIDENCE_INVALID")
	reject(first, evalArgs(sid, revisionOf(t, run), "unknown", "focused-tests", "unknown"), "EVALUATION_EVIDENCE_INVALID")
	reject(first, evalArgs(sid, revisionOf(t, run), eid, "different-check", "wrong-check"), "EVALUATION_EVIDENCE_INVALID")
	originalArgs := evalArgs(sid, revisionOf(t, run), eid, "focused-tests", "original-evaluation")
	approved := callTool(ctx, t, first, "step_evaluate", originalArgs)
	if envData(t, approved)["has_blocking_failure"] != false {
		t.Fatalf("fresh check not accepted: %+v", approved)
	}
	for _, ev := range eventstore.ReadEvents(logPath, nil).Events {
		if ev.Type == eventstore.EventEvaluationRecorded {
			var p struct {
				Lineage []struct {
					EventID string `json:"event_id"`
					CheckID string `json:"check_id"`
				} `json:"evidence_lineage"`
			}
			if err := json.Unmarshal(ev.Payload, &p); err != nil {
				t.Fatal(err)
			}
			if len(p.Lineage) != 1 || p.Lineage[0].EventID == "" || p.Lineage[0].CheckID != "focused-tests" {
				t.Fatal("evaluation lost producer lineage")
			}
		}
	}
	registerArgs := map[string]any{"session_id": sid, "source": "external_artifact", "text": "Observed explanation: errors retain their cause.", "rubric_ref": "rubric://reasoning", "expected_revision": revisionOf(t, approved), "request_id": "register-note"}
	note := callTool(ctx, t, first, "evidence_record", registerArgs)
	noteID := envData(t, note)["evidence_id"].(string)
	qualitative := map[string]any{"session_id": sid, "expected_revision": revisionOf(t, note), "request_id": "qualitative-evaluation", "criteria": []map[string]any{{"name": "reasoning", "kind": "clarity", "severity": "blocking", "verdict": "met", "evidence_id": noteID, "rubric_ref": "rubric://reasoning"}}}
	judgment := callTool(ctx, t, first, "step_evaluate", qualitative)
	if err := os.WriteFile(filepath.Join(learner, "main.go"), []byte("package main\nfunc main(){panic(1)}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	reject(first, evalArgs(sid, revisionOf(t, judgment), eid, "focused-tests", "stale-before-restart"), "EVALUATION_EVIDENCE_STALE")
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	second := connect()
	defer second.Close()
	reject(second, evalArgs(otherID, revisionOf(t, other), eid, "focused-tests", "cross-after-restart"), "EVALUATION_EVIDENCE_INVALID")
	reject(second, evalArgs(sid, revisionOf(t, judgment), eid, "focused-tests", "stale-after-restart"), "EVALUATION_EVIDENCE_STALE")
	retry := callTool(ctx, t, second, "step_evaluate", originalArgs)
	if !reflect.DeepEqual(envData(t, retry), envData(t, approved)) {
		t.Fatal("historical check retry changed")
	}
	retryNote := callTool(ctx, t, second, "evidence_record", registerArgs)
	if !reflect.DeepEqual(envData(t, retryNote), envData(t, note)) {
		t.Fatal("registration retry changed")
	}
	retryJudgment := callTool(ctx, t, second, "step_evaluate", qualitative)
	if !reflect.DeepEqual(envData(t, retryJudgment), envData(t, judgment)) {
		t.Fatal("qualitative retry changed")
	}
	callTool(ctx, t, second, "evidence_get", map[string]any{"session_id": sid, "evidence_id": noteID})
	qualitative["session_id"], qualitative["expected_revision"], qualitative["request_id"] = otherID, revisionOf(t, other), "foreign-note"
	reject(second, qualitative, "EVALUATION_EVIDENCE_INVALID")
	newObs := callTool(ctx, t, second, "workspace_observe", map[string]any{"session_id": sid, "step_id": step, "root": learner, "expected_revision": revisionOf(t, judgment)})
	newRun := callTool(ctx, t, second, "check_run", map[string]any{"session_id": sid, "check_id": "focused-tests", "expected_revision": revisionOf(t, newObs)})
	callTool(ctx, t, second, "step_evaluate", evalArgs(sid, revisionOf(t, newRun), envData(t, newRun)["evidence_id"].(string), "focused-tests", "fresh-after-restart"))
}
