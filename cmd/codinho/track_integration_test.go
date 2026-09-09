package main

import (
	"context"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

const trackPack = `schema_version: 1
id: tracks
version: 1.0.0
themes:
  - {id: alpha, title: Alpha}
  - {id: beta, title: Beta}
tracks:
  - id: authored
    title: Authored path
    challenge_ids: [a, b]
challenges:
  - schema_version: 1
    id: a
    version: 1.0.0
    title: A
    kind: atomic
    difficulty: foundational
    themes: [alpha]
    layers:
      - id: a-layer
        macro_steps:
          - id: a-step
            kind: macro
            instruction: {objective: Alpha objective, scope: Alpha scope}
  - schema_version: 1
    id: b
    version: 2.0.0
    title: B
    kind: atomic
    difficulty: foundational
    themes: [beta]
    prerequisites: [a]
    layers:
      - id: b-layer
        macro_steps:
          - id: b-step
            kind: macro
            instruction: {objective: Beta objective, scope: Beta scope}
`

func TestTrackCompositionOverRealStdio(t *testing.T) {
	bin := buildCodinhoBinary(t)
	root := t.TempDir()
	packs := filepath.Join(root, "packs")
	if err := os.MkdirAll(packs, 0700); err != nil {
		t.Fatal(err)
	}
	for name, data := range map[string]string{"manifest.yaml": "schema_version: 1\npacks: [pack.yaml]\n", "pack.yaml": trackPack} {
		if err := os.WriteFile(filepath.Join(packs, name), []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	connect := func() *mcp.ClientSession {
		cmd := exec.CommandContext(ctx, bin, "serve", "--authoring")
		cmd.Dir = root
		cs, err := mcp.NewClient(&mcp.Implementation{Name: "track-test", Version: "1"}, nil).Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
		if err != nil {
			t.Fatal(err)
		}
		return cs
	}
	cs := connect()
	rec := callTool(ctx, t, cs, "learning_path_recommend", map[string]any{"theme_ids": []string{"beta", "alpha"}})
	if rec["status"] != "ok" {
		t.Fatal(rec)
	}
	data := rec["data"].(map[string]any)
	selection := data["selection"].(map[string]any)
	if selection["coverage"] != "partial" || selection["complete"] != true {
		t.Fatal(selection)
	}
	path := data["path"].(map[string]any)
	ids := path["challenge_ids"].([]any)
	if len(ids) != 2 || ids[0] != "a" || ids[1] != "b" {
		t.Fatal(path)
	}
	start := callTool(ctx, t, cs, "session_start", map[string]any{"composition_id": path["composition_id"], "challenge_ids": ids, "depth": "challenge", "request_id": "composed"})
	if start["status"] != "ok" {
		t.Fatal(start)
	}
	authored := callTool(ctx, t, cs, "session_start", map[string]any{"track_id": "authored", "depth": "challenge"})
	if authored["status"] != "ok" {
		t.Fatal(authored)
	}
	sid := start["session_id"].(string)
	next := callTool(ctx, t, cs, "step_advance", map[string]any{"session_id": sid, "expected_revision": revisionOf(t, start), "override": true, "request_id": "next"})
	if next["status"] != "ok" {
		t.Fatal(next)
	}
	if err := cs.Close(); err != nil {
		t.Fatal(err)
	}
	// Remove the accepted catalog sequence entirely; replay must not recompose.
	if err := os.WriteFile(filepath.Join(packs, "pack.yaml"), []byte("schema_version: 1\nid: empty\nversion: 3.0.0\n"), 0600); err != nil {
		t.Fatal(err)
	}
	cs = connect()
	defer cs.Close()
	got := callTool(ctx, t, cs, "session_get", map[string]any{"session_id": sid})
	if got["status"] != "ok" {
		t.Fatal(got)
	}
	track := got["data"].(map[string]any)["track"].(map[string]any)
	if track["cursor"] != float64(1) || track["challenge_id"] != "b" || track["challenge_version"] != "2.0.0" {
		t.Fatal(track)
	}
	retry := callTool(ctx, t, cs, "session_start", map[string]any{"composition_id": path["composition_id"], "challenge_ids": ids, "depth": "challenge", "request_id": "composed"})
	if retry["session_id"] != sid || revisionOf(t, retry) != revisionOf(t, start) {
		t.Fatal(retry)
	}
	end := callTool(ctx, t, cs, "step_advance", map[string]any{"session_id": sid, "expected_revision": revisionOf(t, got), "override": true})
	if end["status"] != "ok" || end["data"].(map[string]any)["done"] != true {
		t.Fatal(end)
	}
}
