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

func TestSessionRecoveryOverRealStdio(t *testing.T) {
	bin := buildCodinhoBinary(t)
	t.Run("legacy-diagnostic-and-corrupt-log", func(t *testing.T) {
		root := t.TempDir()
		writeIntegrationFixturePack(t, root)
		path := filepath.Join(root, ".codinho", "state", "events.jsonl")
		store, err := eventstore.Open(path, nil)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := store.Append("ses_1", 0, "legacy", eventstore.EventSessionStarted, map[string]string{"challenge_id": "fixture.challenge-one", "mode": "practice"}); err != nil {
			t.Fatal(err)
		}
		store.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, bin, "serve", "--authoring")
		cmd.Dir = root
		cs, err := mcp.NewClient(&mcp.Implementation{Name: "legacy-test", Version: "1"}, nil).Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
		if err != nil {
			t.Fatal(err)
		}
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "session_get", Arguments: map[string]any{"session_id": "ses_1"}})
		if err != nil {
			t.Fatal(err)
		}
		raw, _ := json.Marshal(res.StructuredContent)
		if !strings.Contains(string(raw), "SESSION_RECOVERY_UNAVAILABLE") {
			t.Fatalf("legacy diagnostic: %s", raw)
		}
		fresh := callTool(ctx, t, cs, "session_start", map[string]any{"challenge_id": "fixture.challenge-one", "request_id": "fresh"})
		if fresh["session_id"] == "ses_1" {
			t.Fatal("legacy ID reused")
		}
		cs.Close()
		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o600)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := f.WriteString("corrupted complete line\n"); err != nil {
			t.Fatal(err)
		}
		f.Close()
		before, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		broken := exec.CommandContext(ctx, bin, "serve", "--authoring")
		broken.Dir = root
		if output, err := broken.CombinedOutput(); err == nil || !strings.Contains(string(output), "corrupted log") {
			t.Fatalf("corrupt startup: %s %v", output, err)
		}
		after, err := os.ReadFile(path)
		if err != nil || string(before) != string(after) {
			t.Fatal("corrupt history modified")
		}
	})
	for _, change := range []string{"update", "remove"} {
		t.Run(change, func(t *testing.T) {
			root := t.TempDir()
			writeIntegrationFixturePack(t, root)
			learner := t.TempDir()
			if err := os.WriteFile(filepath.Join(learner, "main.go"), []byte("package main\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			connect := func() *mcp.ClientSession {
				cmd := exec.CommandContext(ctx, bin, "serve", "--authoring")
				cmd.Dir = root
				client := mcp.NewClient(&mcp.Implementation{Name: "recovery-test", Version: "1"}, nil)
				cs, err := client.Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
				if err != nil {
					t.Fatal(err)
				}
				return cs
			}
			first := connect()
			args := map[string]any{"challenge_id": "fixture.challenge-one", "mode": "teaching", "time_limit_seconds": 900, "request_id": "original-start"}
			start := callTool(ctx, t, first, "session_start", args)
			sid := start["session_id"].(string)
			step := start["active_node"].(map[string]any)["id"].(string)
			instruction := envData(t, callTool(ctx, t, first, "instruction_get", map[string]any{"session_id": sid}))
			observeArgs := map[string]any{"session_id": sid, "step_id": step, "root": learner, "globs": []string{"**/*.go"}, "expected_revision": revisionOf(t, start), "request_id": "baseline"}
			baseline := callTool(ctx, t, first, "workspace_observe", observeArgs)
			paused := callTool(ctx, t, first, "session_pause", map[string]any{"session_id": sid, "expected_revision": revisionOf(t, baseline), "request_id": "pause"})
			if err := first.Close(); err != nil {
				t.Fatal(err)
			}

			packPath := filepath.Join(root, "packs", "pack.yaml")
			pack, err := os.ReadFile(packPath)
			if err != nil {
				t.Fatal(err)
			}
			if change == "update" {
				pack = []byte(strings.ReplaceAll(string(pack), "Declare the fixture type.", "Implement a DIFFERENT objective."))
			} else {
				pack = []byte(strings.ReplaceAll(string(pack), "fixture.challenge-one", "fixture.replacement"))
			}
			if err := os.WriteFile(packPath, pack, 0o600); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(learner, "main.go"), []byte("package main\nfunc main() {}\n"), 0o600); err != nil {
				t.Fatal(err)
			}

			second := connect()
			defer second.Close()
			current := callTool(ctx, t, second, "session_get", map[string]any{"session_id": sid})
			if envData(t, current)["state"] != "paused" || revisionOf(t, current) != revisionOf(t, paused) {
				t.Fatalf("state not restored: %+v", current)
			}
			oldInstruction := envData(t, callTool(ctx, t, second, "instruction_get", map[string]any{"session_id": sid}))
			if !reflect.DeepEqual(instruction, oldInstruction) {
				t.Fatalf("pinned instruction changed: %+v", oldInstruction)
			}
			retry := callTool(ctx, t, second, "session_start", args)
			if retry["session_id"] != sid || !reflect.DeepEqual(envData(t, retry), envData(t, start)) {
				t.Fatalf("start retry: %+v", retry)
			}
			repeatedBaseline := callTool(ctx, t, second, "workspace_observe", observeArgs)
			if !reflect.DeepEqual(envData(t, baseline), envData(t, repeatedBaseline)) {
				t.Fatalf("baseline retry changed: %+v", repeatedBaseline)
			}
			oldEvidence := envData(t, callTool(ctx, t, second, "evidence_get", map[string]any{"session_id": sid, "evidence_id": envData(t, baseline)["evidence_id"]}))
			if oldEvidence["stale"] != true {
				t.Fatalf("staleness lost: %+v", oldEvidence)
			}
			observeArgs["expected_revision"], observeArgs["request_id"] = revisionOf(t, current), "diff"
			diff := envData(t, callTool(ctx, t, second, "workspace_observe", observeArgs))
			if diff["baseline"] != false {
				t.Fatalf("baseline reinitialized: %+v", diff)
			}
			challenge := "fixture.challenge-one"
			if change == "remove" {
				challenge = "fixture.replacement"
			}
			fresh := callTool(ctx, t, second, "session_start", map[string]any{"challenge_id": challenge, "request_id": "new-start"})
			if fresh["session_id"] == sid {
				t.Fatal("session ID reused")
			}
			if change == "update" && envData(t, fresh)["objective"] != "Implement a DIFFERENT objective." {
				t.Fatalf("new catalog not used: %+v", fresh)
			}
		})
	}
	t.Run("lost-response-and-truncated-tail", func(t *testing.T) {
		root := t.TempDir()
		writeIntegrationFixturePack(t, root)
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, bin, "serve", "--authoring")
		cmd.Dir = root
		stdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			t.Fatal(err)
		}
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		defer func() {
			if cmd.ProcessState == nil {
				cmd.Process.Kill()
				cmd.Wait()
			}
		}()
		encoder, decoder := json.NewEncoder(stdin), json.NewDecoder(stdout)
		if err := encoder.Encode(map[string]any{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": map[string]any{"protocolVersion": "2025-11-25", "capabilities": map[string]any{}, "clientInfo": map[string]any{"name": "crash-test", "version": "1"}}}); err != nil {
			t.Fatal(err)
		}
		var initialized map[string]any
		if err := decoder.Decode(&initialized); err != nil {
			t.Fatal(err)
		}
		if err := encoder.Encode(map[string]any{"jsonrpc": "2.0", "method": "notifications/initialized"}); err != nil {
			t.Fatal(err)
		}
		args := map[string]any{"challenge_id": "fixture.challenge-one", "request_id": "lost-response"}
		if err := encoder.Encode(map[string]any{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": map[string]any{"name": "session_start", "arguments": args}}); err != nil {
			t.Fatal(err)
		}
		logPath := filepath.Join(root, ".codinho", "state", "events.jsonl")
		deadline := time.Now().Add(5 * time.Second)
		for len(eventstore.ReadEvents(logPath, nil).Events) == 0 {
			if time.Now().After(deadline) {
				t.Fatal("start never committed")
			}
			time.Sleep(5 * time.Millisecond)
		}
		// Never consume the start response: lose the transport after the append.
		if err := cmd.Process.Kill(); err != nil {
			t.Fatal(err)
		}
		_ = cmd.Wait()
		// Preserve the existing operational policy: SIGKILL leaves an orphaned
		// lock. Only remove this synthetic lock after Wait proves its owner exited.
		if err := os.Remove(filepath.Join(root, ".codinho", "state", "lock")); err != nil {
			t.Fatal(err)
		}
		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0o600)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := f.WriteString(`{"schema_version":1,"id":"partial`); err != nil {
			t.Fatal(err)
		}
		f.Close()
		retryCmd := exec.CommandContext(ctx, bin, "serve", "--authoring")
		retryCmd.Dir = root
		cs, err := mcp.NewClient(&mcp.Implementation{Name: "retry-test", Version: "1"}, nil).Connect(ctx, &mcp.CommandTransport{Command: retryCmd}, nil)
		if err != nil {
			t.Fatal(err)
		}
		defer cs.Close()
		start := callTool(ctx, t, cs, "session_start", args)
		if start["session_id"] != "ses_1" || revisionOf(t, start) != 1 {
			t.Fatalf("lost response duplicated: %+v", start)
		}
		callTool(ctx, t, cs, "session_pause", map[string]any{"session_id": "ses_1", "expected_revision": 1})
		if result := eventstore.ReadEvents(logPath, nil); result.Err != nil || result.Truncated || len(result.Events) != 2 {
			t.Fatalf("tail repair failed: %+v", result)
		}
	})
}
