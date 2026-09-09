package main

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const integrationDraft = `schema_version: 1
id: synthetic-draft
version: 1.0.0
competencies:
  - {id: draft-skill, title: Practice}
tracks:
  - {id: draft-track, title: Practice path, challenge_ids: [draft-challenge]}
challenges:
  - schema_version: 1
    id: draft-challenge
    version: 1.0.0
    title: Draft challenge
    competencies: {primary: [draft-skill]}
    acceptance: [Explain the outcome]
    layers:
      - id: understanding
        macro_steps:
          - id: explain
            kind: micro
            instruction: {objective: Explain the outcome}
`

func TestDraftOverRealStdio(t *testing.T) {
	bin := buildCodinhoBinary(t)
	root := t.TempDir()
	writeIntegrationFixturePack(t, root)
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	connect := func() *mcp.ClientSession {
		cmd := exec.CommandContext(ctx, bin, "serve")
		cmd.Dir = root
		cs, err := mcp.NewClient(&mcp.Implementation{Name: "draft-test", Version: "1"}, nil).Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
		if err != nil {
			t.Fatal(err)
		}
		return cs
	}
	cs := connect()
	defer func() { cs.Close() }()
	submit := callTool(ctx, t, cs, "content_draft_submit", map[string]any{"pack_yaml": integrationDraft, "request_id": "submit"})
	data := envData(t, submit)
	id := data["draft_id"]
	retry := callTool(ctx, t, cs, "content_draft_submit", map[string]any{"pack_yaml": integrationDraft, "request_id": "submit"})
	if envData(t, retry)["draft_id"] != id {
		t.Fatal("submission retry changed identity")
	}
	listed := callTool(ctx, t, cs, "catalog_search", map[string]any{"text": "Draft challenge"})
	if strings.Contains(string(mustJSON(t, listed)), "draft-challenge") {
		t.Fatal("quarantine leaked into default catalog")
	}
	denied, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "session_start", Arguments: map[string]any{"draft_id": id, "track_id": "draft-track"}})
	if err != nil || !denied.IsError {
		t.Fatal("draft start without consent", err)
	}
	start := callTool(ctx, t, cs, "session_start", map[string]any{"draft_id": id, "accept_draft": true, "track_id": "draft-track", "request_id": "start"})
	sid := start["session_id"]
	if envData(t, start)["content_provenance"] != "draft" || envData(t, start)["content_warning"] == "" {
		t.Fatal(start)
	}
	review := callTool(ctx, t, cs, "content_draft_get", map[string]any{"draft_id": id})
	if envData(t, review)["pack_yaml"] != integrationDraft {
		t.Fatal("review export altered source")
	}
	callTool(ctx, t, cs, "content_draft_remove", map[string]any{"draft_id": id, "request_id": "remove", "confirm": true})
	cs.Close()
	cs = connect()
	status := callTool(ctx, t, cs, "session_get", map[string]any{"session_id": sid})
	if envData(t, status)["content_provenance"] != "draft" || envData(t, status)["content_warning"] == "" {
		t.Fatal("replay lost warning", status)
	}
	instruction := callTool(ctx, t, cs, "instruction_get", map[string]any{"session_id": sid})
	if !strings.Contains(string(mustJSON(t, instruction)), "Explain the outcome") {
		t.Fatal("pinned content unavailable")
	}
	denied, err = cs.CallTool(ctx, &mcp.CallToolParams{Name: "session_start", Arguments: map[string]any{"draft_id": id, "accept_draft": true, "challenge_id": "draft-challenge", "request_id": "new-start"}})
	if err != nil || !denied.IsError {
		t.Fatal("removed draft started again", err)
	}
	cs.Close()
	dest := t.TempDir()
	export := exec.CommandContext(ctx, bin, "privacy", "export", "--dest", dest)
	export.Dir = root
	if out, err := export.CombinedOutput(); err != nil {
		t.Fatalf("export: %v %s", err, out)
	}
	raw, err := os.ReadFile(filepath.Join(dest, "events.jsonl"))
	if err != nil || !strings.Contains(string(raw), "draft_submitted") || !strings.Contains(string(raw), "draft_removed") {
		t.Fatal("export omitted quarantine", err)
	}
	purge := exec.CommandContext(ctx, bin, "privacy", "purge", "--confirm")
	purge.Dir = root
	if out, err := purge.CombinedOutput(); err != nil {
		t.Fatalf("purge: %v %s", err, out)
	}
	if _, err := os.Stat(filepath.Join(root, ".codinho", "state")); !os.IsNotExist(err) {
		t.Fatal("purge left quarantine", err)
	}
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
