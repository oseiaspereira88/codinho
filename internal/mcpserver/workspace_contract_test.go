package mcpserver

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestContractWorkspaceObserveEstablishesBaselineThenReportsDiff(t *testing.T) {
	cs := newContractClient(t)
	ctx := context.Background()

	startRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "session_start", Arguments: map[string]any{"challenge_id": fixtureChallengeID}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	start := decodeEnvelope(t, startRes)
	if start.Status != "ok" {
		t.Fatalf("unexpected session_start envelope: %+v", start)
	}

	getRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "session_get", Arguments: map[string]any{"session_id": start.SessionID}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	get := decodeEnvelope(t, getRes)
	getData := get.Data.(map[string]any)
	startRevision := getData["revision"].(float64)

	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\n"), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	baselineRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "workspace_observe", Arguments: map[string]any{
		"session_id": start.SessionID, "step_id": start.ActiveNode.ID, "root": root, "globs": []string{"**/*.go"},
		"expected_revision": startRevision,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	baseline := decodeEnvelope(t, baselineRes)
	if baseline.Status != "ok" {
		t.Fatalf("unexpected workspace_observe (baseline) envelope: %+v (error=%+v)", baseline, baseline.Error)
	}
	data, ok := baseline.Data.(map[string]any)
	if !ok {
		t.Fatalf("unexpected data shape: %T", baseline.Data)
	}
	if data["baseline"] != true {
		t.Fatalf("expected baseline=true on first call, got %+v", data)
	}
	revision, ok := data["revision"].(float64)
	if !ok || revision != startRevision+1 {
		t.Fatalf("expected revision %v, got %+v", startRevision+1, data["revision"])
	}
	evidenceID, _ := data["evidence_id"].(string)
	if evidenceID == "" {
		t.Fatal("expected a non-empty evidence_id")
	}

	if err := os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	diffRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "workspace_observe", Arguments: map[string]any{
		"session_id": start.SessionID, "step_id": start.ActiveNode.ID, "root": root, "globs": []string{"**/*.go"},
		"expected_revision": revision,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	diff := decodeEnvelope(t, diffRes)
	if diff.Status != "ok" {
		t.Fatalf("unexpected workspace_observe (diff) envelope: %+v", diff)
	}
	diffData := diff.Data.(map[string]any)
	if diffData["baseline"] != false {
		t.Fatalf("expected baseline=false on the second call, got %+v", diffData)
	}
	changes, ok := diffData["changes"].([]any)
	if !ok || len(changes) != 1 {
		t.Fatalf("expected exactly one change, got %+v", diffData["changes"])
	}

	evidenceRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "evidence_get", Arguments: map[string]any{
		"session_id": start.SessionID, "evidence_id": evidenceID,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	evidenceEnv := decodeEnvelope(t, evidenceRes)
	if evidenceEnv.Status != "ok" {
		t.Fatalf("unexpected evidence_get envelope: %+v", evidenceEnv)
	}
}

func TestContractEvidenceGetOutOfScopeReturnsItemNotFound(t *testing.T) {
	cs := newContractClient(t)
	ctx := context.Background()

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "evidence_get", Arguments: map[string]any{
		"session_id": "some-session", "evidence_id": "sha256-0000000000000000000000000000000000000000000000000000000000000000",
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected IsError for evidence never recorded for this session")
	}
	env := decodeEnvelope(t, res)
	if env.Status != "error" || env.Error == nil || env.Error.Code != string(ErrCodeItemNotFound) {
		t.Fatalf("unexpected envelope: %+v", env)
	}
}

func TestContractWorkspaceObserveRejectsMissingRequiredArguments(t *testing.T) {
	cs := newContractClient(t)
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "workspace_observe", Arguments: map[string]any{
		"session_id": "sess", "step_id": "", "root": "", "expected_revision": 0,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected IsError when step_id and root are empty")
	}
	env := decodeEnvelope(t, res)
	if env.Status != "error" || env.Error == nil || env.Error.Code != string(ErrCodeInvalidInput) {
		t.Fatalf("unexpected envelope: %+v", env)
	}
}
