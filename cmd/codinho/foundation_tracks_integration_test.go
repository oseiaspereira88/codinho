package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestFoundationTracksOverRealStdio checks the actual authored catalog, not
// another synthetic graph. Override traversal is a routing check, not playtest.
func TestFoundationTracksOverRealStdio(t *testing.T) {
	bin := buildCodinhoBinary(t)
	root := t.TempDir()
	if err := os.CopyFS(filepath.Join(root, "packs"), os.DirFS("../../packs")); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, "serve", "--authoring")
	cmd.Dir = root
	cs, err := mcp.NewClient(&mcp.Implementation{Name: "foundation-tracks-test", Version: "1"}, nil).Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer cs.Close()
	listed := callTool(ctx, t, cs, "catalog_search", map[string]any{"kind": "track"})
	items, ok := listed["data"].([]any)
	if !ok || len(items) != 5 {
		t.Fatalf("expected five authored tracks: %+v", listed)
	}
	for _, id := range []string{"go-from-zero", "go-oop-transition", "go-practical-fluency", "go-data-and-types", "go-testing-and-design"} {
		t.Run(id, func(t *testing.T) {
			current := callTool(ctx, t, cs, "session_start", map[string]any{"track_id": id, "depth": "challenge", "request_id": id})
			sid := current["session_id"]
			data := envData(t, current)
			track := data["track"].(map[string]any)
			total := int(track["total"].(float64))
			if total < 2 || total > 100 || data["content_provenance"] != "draft" {
				t.Fatalf("unexpected authored path: %+v", data)
			}
			seen := map[string]bool{}
			for cursor := 0; cursor < total; cursor++ {
				status := callTool(ctx, t, cs, "session_get", map[string]any{"session_id": sid})
				data = envData(t, status)
				track = data["track"].(map[string]any)
				challenge := track["challenge_id"].(string)
				if int(track["cursor"].(float64)) != cursor || challenge == "" || seen[challenge] {
					t.Fatalf("invalid traversal: %+v", track)
				}
				seen[challenge] = true
				instruction := callTool(ctx, t, cs, "instruction_get", map[string]any{"session_id": sid})
				if objective, _ := envData(t, instruction)["objective"].(string); objective == "" {
					t.Fatal("unreachable authored instruction", instruction)
				}
				if cursor+1 < total {
					callTool(ctx, t, cs, "step_advance", map[string]any{"session_id": sid, "expected_revision": revisionOf(t, status), "override": true})
				}
			}
		})
	}
}
