package mcpserver

import (
	"context"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"testing"
)

func TestContractTrackSelection(t *testing.T) {
	cs := newContractClient(t)
	for _, tc := range []struct {
		tool string
		args map[string]any
	}{
		{"catalog_search", map[string]any{"theme": "slices", "theme_ids": []string{"slices"}}},
		{"catalog_search", map[string]any{"theme_ids": []string{}}},
		{"learning_path_recommend", map[string]any{"competency_id": "slice-filter", "competency_ids": []string{"slice-filter"}}},
		{"session_start", map[string]any{}},
		{"session_start", map[string]any{"track_id": "unknown"}},
		{"session_start", map[string]any{"track_id": "track", "challenge_id": fixtureChallengeID}},
	} {
		res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: tc.tool, Arguments: tc.args})
		if err != nil {
			t.Fatal(err)
		}
		if !res.IsError {
			t.Fatalf("accepted %s %+v", tc.tool, tc.args)
		}
	}
	for _, args := range []map[string]any{{"competency": "slice-filter"}, {"competency_ids": []string{"slice-filter"}}} {
		res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "catalog_search", Arguments: args})
		if err != nil {
			t.Fatal(err)
		}
		env := decodeEnvelope(t, res)
		if env.Status != "ok" || env.Selection == nil || env.Selection.Coverage != "total" {
			t.Fatalf("%+v", env)
		}
		if _, ok := env.Data.([]any); !ok {
			t.Fatalf("legacy data array changed: %+v", env)
		}
	}
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "learning_path_recommend", Arguments: map[string]any{"competency_ids": []string{"slice-filter", "missing"}}})
	if err != nil {
		t.Fatal(err)
	}
	env := decodeEnvelope(t, res)
	if env.Status != "ok" {
		t.Fatalf("%+v", env)
	}
	data := env.Data.(map[string]any)
	selection := data["selection"].(map[string]any)
	if selection["complete"] != false || selection["coverage"] != "partial" {
		t.Fatal(selection)
	}
	if len(selection["modes"].([]any)) != 3 {
		t.Fatal("selection modes missing")
	}
}
