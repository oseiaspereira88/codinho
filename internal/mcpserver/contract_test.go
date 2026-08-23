package mcpserver

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/oseiaspereira88/codinho/internal/application"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/evidence"
)

const fixtureChallengeID = "fixture.challenge-one"

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func newTestCatalog(t *testing.T) *curriculum.Catalog {
	t.Helper()
	dir := t.TempDir()
	writeFile(t, dir, "manifest.yaml", "schema_version: 1\npacks:\n  - pack.yaml\n")
	writeFile(t, dir, "pack.yaml", `schema_version: 1
id: fixture-pack
version: 1.0.0
themes:
  - id: slices
    title: Slices
concepts:
  - id: slice-declaration
    title: Slice declaration
competencies:
  - id: slice-filter
    title: Filter a slice
challenges:
  - schema_version: 1
    id: fixture.challenge-one
    version: 1.0.0
    title: Fixture challenge
    kind: atomic
    difficulty: foundational
    themes: [slices]
    competencies:
      primary: [slice-filter]
    prerequisites: [slice-declaration]
    layers:
      - id: understanding
        macro_steps:
          - id: fixture.step-one
            kind: micro
            title: Fixture step
            instruction:
              objective: Declare the fixture type.
              scope: Only the declaration.
            concepts: [slice-declaration]
            evidence:
              strategies: [compile]
relations:
  - from: fixture.challenge-one
    to: slice-declaration
    kind: requires
`)
	catalog, diags, err := curriculum.Load(dir, curriculum.DefaultLimits)
	if err != nil {
		t.Fatalf("unexpected error loading fixture catalog: %v", err)
	}
	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostics loading fixture catalog: %+v", diags)
	}
	return catalog
}

// newContractClient starts a real server behind an in-memory net.Pipe
// transport and connects a real MCP client to it, so contract tests
// exercise the actual protocol round trip, not just Go function calls.
func newContractClient(t *testing.T) *mcp.ClientSession {
	t.Helper()
	catalog := newTestCatalog(t)
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
	assistanceService := application.NewAssistanceService(catalogService)
	workspaceService := application.NewWorkspaceService(store, evidenceStore)
	checksService := application.NewChecksService(store, sessionService, workspaceService)
	progressService := application.NewProgressService(store)
	recommendationService := application.NewRecommendationService(catalog, progressService)
	server := New(Deps{
		Catalog: catalogService, Session: sessionService, Assistance: assistanceService, Workspace: workspaceService,
		Checks: checksService, Progress: progressService, Recommendation: recommendationService,
	}, io.Discard)

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

	client := mcp.NewClient(&mcp.Implementation{Name: "contract-test-client", Version: "0.0.0"}, nil)
	cs, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("unexpected error connecting client: %v", err)
	}
	t.Cleanup(func() { cs.Close() })
	return cs
}

func decodeEnvelope(t *testing.T, res *mcp.CallToolResult) Envelope {
	t.Helper()
	raw, err := json.Marshal(res.StructuredContent)
	if err != nil {
		t.Fatalf("unexpected error marshaling structured content: %v", err)
	}
	var env Envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatalf("unexpected error decoding envelope: %v (raw: %s)", err, raw)
	}
	return env
}

func TestContractServerPublishesInvariantsInInstructions(t *testing.T) {
	cs := newContractClient(t)
	if got := cs.InitializeResult().Instructions; got == "" {
		t.Fatal("server did not publish instructions")
	} else if len(got) < len(coreInstructions) || got[:len(coreInstructions)] != coreInstructions {
		t.Fatalf("published instructions do not start with the core invariants: %q", got)
	}
}

func TestContractListsExactlyTheMinimalToolSlice(t *testing.T) {
	cs := newContractClient(t)
	res, err := cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]bool{
		"catalog_search": false, "catalog_get": false,
		"session_start": false, "session_get": false, "instruction_get": false,
		"session_configure": false, "session_pause": false, "session_resume": false,
		"session_finish": false, "granularity_adjust": false,
		"hint_request": false, "syntax_recall_get": false, "concept_content_get": false,
		"learning_detour_start": false, "learning_detour_finish": false,
		"feedback_prepare": false, "feedback_record": false, "step_evaluate": false,
		"reflection_record": false, "step_complete": false, "step_advance": false,
		"workspace_observe": false, "evidence_get": false, "check_run": false,
		"progress_get": false, "review_due": false, "mastery_evidence_record": false,
		"concept_relations_get": false, "learning_path_recommend": false,
	}
	for _, tool := range res.Tools {
		if _, known := want[tool.Name]; !known {
			t.Fatalf("unexpected tool exposed: %s", tool.Name)
		}
		want[tool.Name] = true
	}
	for name, seen := range want {
		if !seen {
			t.Fatalf("required tool missing: %s", name)
		}
	}
}

func TestContractCatalogSearchAndGet(t *testing.T) {
	cs := newContractClient(t)
	ctx := context.Background()

	searchRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "catalog_search", Arguments: map[string]any{"theme": "slices"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env := decodeEnvelope(t, searchRes)
	if env.Status != "ok" {
		t.Fatalf("catalog_search status = %s, want ok", env.Status)
	}

	getRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "catalog_get", Arguments: map[string]any{"id": fixtureChallengeID}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env = decodeEnvelope(t, getRes)
	if env.Status != "ok" {
		t.Fatalf("catalog_get status = %s, want ok: %+v", env.Status, env.Error)
	}
}

func TestContractCatalogGetUnknownIDReturnsItemNotFound(t *testing.T) {
	cs := newContractClient(t)
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "catalog_get", Arguments: map[string]any{"id": "does-not-exist"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected IsError for an unknown catalog ID")
	}
	env := decodeEnvelope(t, res)
	if env.Status != "error" || env.Error == nil || env.Error.Code != string(ErrCodeItemNotFound) {
		t.Fatalf("unexpected envelope: %+v", env)
	}
}

func TestContractSessionLifecycleEndToEnd(t *testing.T) {
	cs := newContractClient(t)
	ctx := context.Background()

	startRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "session_start", Arguments: map[string]any{"challenge_id": fixtureChallengeID}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	start := decodeEnvelope(t, startRes)
	if start.Status != "ok" || start.SessionID == "" || start.ActiveNode == nil {
		t.Fatalf("unexpected session_start envelope: %+v", start)
	}

	getRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "session_get", Arguments: map[string]any{"session_id": start.SessionID}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	get := decodeEnvelope(t, getRes)
	if get.Status != "ok" || get.SessionID != start.SessionID {
		t.Fatalf("unexpected session_get envelope: %+v", get)
	}

	instrRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "instruction_get", Arguments: map[string]any{"session_id": start.SessionID}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	instr := decodeEnvelope(t, instrRes)
	if instr.Status != "ok" || instr.Disclosure == nil || instr.Disclosure.SolutionRevealed {
		t.Fatalf("unexpected instruction_get envelope: %+v", instr)
	}
}

func TestContractSessionStartIsIdempotentByRequestID(t *testing.T) {
	cs := newContractClient(t)
	ctx := context.Background()
	args := map[string]any{"challenge_id": fixtureChallengeID, "request_id": "contract-retry-1"}

	first, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "session_start", Arguments: args})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "session_start", Arguments: args})
	if err != nil {
		t.Fatalf("unexpected error on retry: %v", err)
	}

	envFirst := decodeEnvelope(t, first)
	envSecond := decodeEnvelope(t, second)
	if envFirst.SessionID != envSecond.SessionID {
		t.Fatalf("retrying session_start with the same request_id produced different sessions: %s vs %s", envFirst.SessionID, envSecond.SessionID)
	}
}

func TestContractSessionGetUnknownSessionReturnsSessionNotActive(t *testing.T) {
	cs := newContractClient(t)
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "session_get", Arguments: map[string]any{"session_id": "ses_does_not_exist"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected IsError for an unknown session")
	}
	env := decodeEnvelope(t, res)
	if env.Error == nil || env.Error.Code != string(ErrCodeSessionNotActive) {
		t.Fatalf("unexpected envelope: %+v", env)
	}
}

// TestContractMissingRequiredArgumentIsRejectedBySchema documents that the
// strict input schema (requirement R6) rejects a call missing a required
// argument before the handler ever runs: this is an SDK-level protocol
// error, not our Envelope, and never leaks internal detail (requirement
// R4, security).
func TestContractMissingRequiredArgumentIsRejectedBySchema(t *testing.T) {
	cs := newContractClient(t)
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "session_start", Arguments: map[string]any{}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected IsError when challenge_id is missing")
	}
	if res.StructuredContent != nil {
		t.Fatalf("schema-level rejection should not carry our Envelope, got %+v", res.StructuredContent)
	}
}

// TestContractEmptyRequiredArgumentReturnsInvalidInput covers the case the
// schema alone cannot catch (an empty string still satisfies type:string),
// which is exactly why the handler also validates (requirement R4).
func TestContractEmptyRequiredArgumentReturnsInvalidInput(t *testing.T) {
	cs := newContractClient(t)
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "session_start", Arguments: map[string]any{"challenge_id": ""}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected IsError when challenge_id is empty")
	}
	env := decodeEnvelope(t, res)
	if env.Error == nil || env.Error.Code != string(ErrCodeInvalidInput) {
		t.Fatalf("unexpected envelope: %+v", env)
	}
}

func newTestServer(t *testing.T) *mcp.Server {
	t.Helper()
	catalog := newTestCatalog(t)
	store, err := eventstore.Open(filepath.Join(t.TempDir(), "events.jsonl"), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	catalogService := application.NewCatalogService(catalog)
	sessionService := application.NewSessionService(catalogService, store, nil)
	return New(Deps{Catalog: catalogService, Session: sessionService}, io.Discard)
}

// TestRunReturnsCleanlyOnContextCancel documents and verifies half of
// requirement R1: cancellation ends the server without treating it as a
// failure.
func TestRunReturnsCleanlyOnContextCancel(t *testing.T) {
	server := newTestServer(t)
	clientTransport, serverTransport := mcp.NewInMemoryTransports()

	ctx, cancel := context.WithCancel(context.Background())
	runErr := make(chan error, 1)
	go func() { runErr <- Run(ctx, server, serverTransport) }()

	client := mcp.NewClient(&mcp.Implementation{Name: "cancel-test-client", Version: "0.0.0"}, nil)
	cs, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer cs.Close()

	cancel()

	select {
	case err := <-runErr:
		if err != nil {
			t.Fatalf("Run must return nil on context cancellation, got %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not return within 5s of context cancellation")
	}
}

// pipeReadWriteCloser adapts an io.Reader/io.Writer pair (from io.Pipe) to
// the io.ReadCloser/io.WriteCloser pair mcp.IOTransport expects, so this
// test can close the "stdin" side exactly like a real client disconnecting.
type pipeReadCloser struct{ io.Reader }

func (pipeReadCloser) Close() error { return nil }

// TestRunReturnsCleanlyWhenStdinCloses verifies the other half of
// requirement R1: a client-side EOF (stdin closing) is also a clean
// shutdown, not a failure that Run reports as an error.
func TestRunReturnsCleanlyWhenStdinCloses(t *testing.T) {
	server := newTestServer(t)
	stdinReader, stdinWriter := io.Pipe()
	transport := &mcp.IOTransport{Reader: pipeReadCloser{stdinReader}, Writer: nopWriteCloser{io.Discard}}

	runErr := make(chan error, 1)
	go func() { runErr <- Run(context.Background(), server, transport) }()

	// Closing the write side of the pipe delivers EOF to the server's
	// reader, exactly like a client closing stdin.
	if err := stdinWriter.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	select {
	case err := <-runErr:
		if err != nil {
			t.Fatalf("Run must return nil when stdin closes cleanly, got %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not return within 5s of stdin closing")
	}
}

type nopWriteCloser struct{ io.Writer }

func (nopWriteCloser) Close() error { return nil }
