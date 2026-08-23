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
)

// newHintedContractClient is like newContractClient but its fixture step
// authors a full hint ladder and a concept, so hint_request,
// syntax_recall_get and concept_content_get have real content to exercise
// (assistance-hints-detours).
func newHintedContractClient(t *testing.T) *mcp.ClientSession {
	t.Helper()
	dir := t.TempDir()
	writeFile(t, dir, "manifest.yaml", "schema_version: 1\npacks:\n  - pack.yaml\n")
	writeFile(t, dir, "pack.yaml", `schema_version: 1
id: fixture-pack
version: 1.0.0
concepts:
  - id: named-types
    title: Named types
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
          - id: fixture.macro-one
            kind: micro
            title: Macro step
            instruction:
              objective: Declare the User type.
              scope: Only the declaration.
            concepts: [named-types]
            hints:
              - level: 1
                kind: guiding_question
              - level: 2
                kind: syntax_recall
              - level: 3
                kind: logical_prose
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

	catalogService := application.NewCatalogService(catalog)
	sessionService := application.NewSessionService(catalogService, store, nil)
	assistanceService := application.NewAssistanceService(catalogService)
	server := New(Deps{Catalog: catalogService, Session: sessionService, Assistance: assistanceService}, io.Discard)

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

	client := mcp.NewClient(&mcp.Implementation{Name: "assistance-contract-test-client", Version: "0.0.0"}, nil)
	cs, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("unexpected error connecting client: %v", err)
	}
	t.Cleanup(func() { cs.Close() })
	return cs
}

func startHintedFixtureSession(t *testing.T, cs *mcp.ClientSession, extra map[string]any) Envelope {
	t.Helper()
	args := map[string]any{"challenge_id": fixtureChallengeID}
	for k, v := range extra {
		args[k] = v
	}
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "session_start", Arguments: args})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env := decodeEnvelope(t, res)
	if env.Status != "ok" {
		t.Fatalf("session_start failed: %+v", env)
	}
	return env
}

func TestContractHintRequestClimbsLadderToSolution(t *testing.T) {
	cs := newHintedContractClient(t)
	ctx := context.Background()
	start := startHintedFixtureSession(t, cs, map[string]any{"disclosure_max": 6})
	rev := revisionOf(t, start)

	for want := 1; want <= 5; want++ {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "hint_request", Arguments: map[string]any{
			"session_id": start.SessionID, "expected_revision": rev,
		}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		env := decodeEnvelope(t, res)
		if env.Status != "ok" || env.ProgressEffect != ProgressEffectNone {
			t.Fatalf("unexpected envelope at level %d: %+v", want, env)
		}
		data, ok := env.Data.(map[string]any)
		if !ok || int(data["level"].(float64)) != want {
			t.Fatalf("level = %v, want %d: %+v", data["level"], want, env)
		}
		rev = revisionOf(t, env)
	}

	// Level 6 without confirmation is rejected with zero progress change.
	blocked, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "hint_request", Arguments: map[string]any{
		"session_id": start.SessionID, "expected_revision": rev,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !blocked.IsError {
		t.Fatal("expected IsError requesting the solution without confirm_solution")
	}
	blockedEnv := decodeEnvelope(t, blocked)
	if blockedEnv.Error == nil || blockedEnv.Error.Code != string(ErrCodeInvalidInput) {
		t.Fatalf("unexpected envelope: %+v", blockedEnv)
	}

	confirmed, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "hint_request", Arguments: map[string]any{
		"session_id": start.SessionID, "expected_revision": rev, "confirm_solution": true,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	confirmedEnv := decodeEnvelope(t, confirmed)
	if confirmedEnv.Status != "ok" || confirmedEnv.ProgressEffect != ProgressEffectSolutionReveal {
		t.Fatalf("unexpected envelope revealing the solution: %+v", confirmedEnv)
	}
}

func TestContractSyntaxRecallGetFreeInTeachingMode(t *testing.T) {
	cs := newHintedContractClient(t)
	ctx := context.Background()
	start := startHintedFixtureSession(t, cs, map[string]any{"mode": "teaching"})
	rev := revisionOf(t, start)

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "syntax_recall_get", Arguments: map[string]any{
		"session_id": start.SessionID, "expected_revision": rev,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env := decodeEnvelope(t, res)
	if env.Status != "ok" {
		t.Fatalf("syntax_recall_get failed: %+v", env)
	}
	data, ok := env.Data.(map[string]any)
	if !ok || data["kind"] != "syntax_recall" {
		t.Fatalf("unexpected envelope: %+v", env)
	}
}

func TestContractConceptContentGet(t *testing.T) {
	cs := newHintedContractClient(t)
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "concept_content_get", Arguments: map[string]any{"concept_id": "named-types"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env := decodeEnvelope(t, res)
	if env.Status != "ok" {
		t.Fatalf("concept_content_get failed: %+v", env)
	}
	data, ok := env.Data.(map[string]any)
	if !ok || data["id"] != "named-types" || data["title"] != "Named types" {
		t.Fatalf("unexpected envelope: %+v", env)
	}
}

func TestContractConceptContentGetUnknownIDReturnsItemNotFound(t *testing.T) {
	cs := newHintedContractClient(t)
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "concept_content_get", Arguments: map[string]any{"concept_id": "does-not-exist"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected IsError for an unknown concept ID")
	}
	env := decodeEnvelope(t, res)
	if env.Error == nil || env.Error.Code != string(ErrCodeItemNotFound) {
		t.Fatalf("unexpected envelope: %+v", env)
	}
}

func TestContractLearningDetourStartAndFinishPreserveActiveStep(t *testing.T) {
	cs := newHintedContractClient(t)
	ctx := context.Background()
	start := startHintedFixtureSession(t, cs, nil)
	rev := revisionOf(t, start)

	startRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "learning_detour_start", Arguments: map[string]any{
		"session_id": start.SessionID, "reason": "why exported fields?", "expected_revision": rev,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	detourStart := decodeEnvelope(t, startRes)
	if detourStart.Status != "ok" {
		t.Fatalf("learning_detour_start failed: %+v", detourStart)
	}
	rev = revisionOf(t, detourStart)

	sessRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "session_get", Arguments: map[string]any{"session_id": start.SessionID}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sessDuring := decodeEnvelope(t, sessRes)
	if sessDuring.ActiveNode == nil || sessDuring.ActiveNode.ID != start.ActiveNode.ID {
		t.Fatalf("active node changed during detour: %+v", sessDuring)
	}

	finishRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "learning_detour_finish", Arguments: map[string]any{
		"session_id": start.SessionID, "outcome": "resolved", "expected_revision": rev,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	detourFinish := decodeEnvelope(t, finishRes)
	if detourFinish.Status != "ok" {
		t.Fatalf("learning_detour_finish failed: %+v", detourFinish)
	}
}

func TestContractLearningDetourFinishWithoutOpenDetourFails(t *testing.T) {
	cs := newHintedContractClient(t)
	ctx := context.Background()
	start := startHintedFixtureSession(t, cs, nil)
	rev := revisionOf(t, start)

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "learning_detour_finish", Arguments: map[string]any{
		"session_id": start.SessionID, "expected_revision": rev,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected IsError finishing a detour that was never started")
	}
}
