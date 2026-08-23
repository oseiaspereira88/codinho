package mcpserver

import (
	"context"
	"io"
	"path/filepath"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/oseiaspereira88/codinho/internal/application"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/evidence"
)

// newChecksContractClient is like newContractClient, but its fixture
// challenge declares a real check (safe-check-executor), so check_run has
// something legitimate to resolve and execute.
func newChecksContractClient(t *testing.T) *mcp.ClientSession {
	t.Helper()
	dir := t.TempDir()
	writeFile(t, dir, "manifest.yaml", "schema_version: 1\npacks:\n  - pack.yaml\n")
	writeFile(t, dir, "pack.yaml", `schema_version: 1
id: fixture-pack
version: 1.0.0
challenges:
  - schema_version: 1
    id: fixture.challenge-one
    version: 1.0.0
    title: Fixture challenge
    kind: atomic
    difficulty: foundational
    layers:
      - id: understanding
        macro_steps:
          - id: fixture.step-one
            kind: micro
            title: Fixture step
            instruction:
              objective: Declare the fixture type.
              scope: Only the declaration.
    checks:
      - id: focused-tests
        runner: go_test
        package: ./...
`)
	catalog, diags, err := curriculum.Load(dir, curriculum.DefaultLimits)
	if err != nil {
		t.Fatalf("unexpected error loading fixture catalog: %v", err)
	}
	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostics loading fixture catalog: %+v", diags)
	}

	store, err := eventstore.Open(filepath.Join(t.TempDir(), "events.jsonl"), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	evidenceStore, err := evidence.Open(filepath.Join(t.TempDir(), "evidence"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	catalogService := application.NewCatalogService(catalog)
	sessionService := application.NewSessionService(catalogService, store, evidenceStore)
	workspaceService := application.NewWorkspaceService(store, evidenceStore)
	checksService := application.NewChecksService(store, sessionService, workspaceService)
	server := New(Deps{Catalog: catalogService, Session: sessionService, Workspace: workspaceService, Checks: checksService}, io.Discard)

	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = server.Run(ctx, serverTransport)
	}()
	t.Cleanup(func() {
		cancel()
		<-done
	})

	client := mcp.NewClient(&mcp.Implementation{Name: "checks-contract-test-client", Version: "0.0.0"}, nil)
	cs, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("unexpected error connecting client: %v", err)
	}
	t.Cleanup(func() { cs.Close() })
	return cs
}

func newFixtureGoModuleDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	writeFile(t, dir, "go.mod", "module fixture\n\ngo 1.25\n")
	writeFile(t, dir, "main.go", "package main\n\nfunc main() {}\n")
	writeFile(t, dir, "main_test.go", "package main\n\nimport \"testing\"\n\nfunc TestOK(t *testing.T) {}\n")
	return dir
}

func TestContractCheckRunEndToEndAfterWorkspaceObserve(t *testing.T) {
	cs := newChecksContractClient(t)
	ctx := context.Background()

	start := startFixtureSession(t, cs)
	rev := revisionOf(t, start)
	activeNode := start.ActiveNode.ID

	root := newFixtureGoModuleDir(t)
	obsRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "workspace_observe", Arguments: map[string]any{
		"session_id": start.SessionID, "step_id": activeNode, "root": root, "globs": []string{"**/*.go"}, "expected_revision": rev,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	obs := decodeEnvelope(t, obsRes)
	if obs.Status != "ok" {
		t.Fatalf("unexpected workspace_observe envelope: %+v (error=%+v)", obs, obs.Error)
	}
	rev = revisionOf(t, obs)

	runRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "check_run", Arguments: map[string]any{
		"session_id": start.SessionID, "check_id": "focused-tests", "expected_revision": rev,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	run := decodeEnvelope(t, runRes)
	if run.Status != "ok" {
		t.Fatalf("unexpected check_run envelope: %+v (error=%+v)", run, run.Error)
	}
	runData := run.Data.(map[string]any)
	if runData["outcome"] != "pass" {
		t.Fatalf("expected outcome=pass, got %+v", runData)
	}
	evidenceID, _ := runData["evidence_id"].(string)
	if evidenceID == "" {
		t.Fatal("expected a non-empty evidence_id")
	}

	evRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "evidence_get", Arguments: map[string]any{
		"session_id": start.SessionID, "evidence_id": evidenceID,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ev := decodeEnvelope(t, evRes)
	if ev.Status != "ok" {
		t.Fatalf("unexpected evidence_get envelope: %+v", ev)
	}
}

func TestContractCheckRunWithoutWorkspaceObserveReturnsInvalidInput(t *testing.T) {
	cs := newChecksContractClient(t)
	ctx := context.Background()
	start := startFixtureSession(t, cs)
	rev := revisionOf(t, start)

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "check_run", Arguments: map[string]any{
		"session_id": start.SessionID, "check_id": "focused-tests", "expected_revision": rev,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected IsError when no workspace_observe happened yet")
	}
	env := decodeEnvelope(t, res)
	if env.Error == nil || env.Error.Code != string(ErrCodeInvalidInput) {
		t.Fatalf("unexpected envelope: %+v", env)
	}
}

func TestContractCheckRunUnknownCheckIDReturnsItemNotFound(t *testing.T) {
	cs := newChecksContractClient(t)
	ctx := context.Background()
	start := startFixtureSession(t, cs)
	rev := revisionOf(t, start)

	root := newFixtureGoModuleDir(t)
	obsRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "workspace_observe", Arguments: map[string]any{
		"session_id": start.SessionID, "step_id": start.ActiveNode.ID, "root": root, "expected_revision": rev,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	obs := decodeEnvelope(t, obsRes)
	rev = revisionOf(t, obs)

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "check_run", Arguments: map[string]any{
		"session_id": start.SessionID, "check_id": "does-not-exist", "expected_revision": rev,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected IsError for an unknown check_id")
	}
	env := decodeEnvelope(t, res)
	if env.Error == nil || env.Error.Code != string(ErrCodeItemNotFound) {
		t.Fatalf("unexpected envelope: %+v", env)
	}
}
