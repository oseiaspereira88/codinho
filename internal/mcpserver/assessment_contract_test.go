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

// newPolicyContractClient is like newContractClient but its fixture step
// requires a positive evaluation to complete, so step_evaluate/
// step_complete/step_advance have a real policy to exercise
// (feedback-evaluation-progression).
func newPolicyContractClient(t *testing.T) *mcp.ClientSession {
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
            completion:
              requires_positive_evaluation: true
          - id: fixture.step-two
            kind: micro
            title: Second fixture step
            instruction:
              objective: Second objective.
              scope: Second scope.
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

	client := mcp.NewClient(&mcp.Implementation{Name: "assessment-contract-test-client", Version: "0.0.0"}, nil)
	cs, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("unexpected error connecting client: %v", err)
	}
	t.Cleanup(func() { cs.Close() })
	return cs
}

func TestContractFeedbackPrepareAndRecord(t *testing.T) {
	cs := newPolicyContractClient(t)
	ctx := context.Background()
	start := startFixtureSession(t, cs)
	rev := revisionOf(t, start)

	prepRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "feedback_prepare", Arguments: map[string]any{"session_id": start.SessionID, "question": "why exported fields?"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	prep := decodeEnvelope(t, prepRes)
	if prep.Status != "ok" {
		t.Fatalf("feedback_prepare failed: %+v", prep)
	}

	recRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "feedback_record", Arguments: map[string]any{
		"session_id": start.SessionID, "type": "explanation", "text": "exported fields are visible across packages", "expected_revision": rev,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec := decodeEnvelope(t, recRes)
	if rec.Status != "ok" || rec.ProgressEffect != ProgressEffectFeedbackRecorded {
		t.Fatalf("unexpected feedback_record envelope: %+v", rec)
	}
}

func TestContractStepEvaluateCompleteAdvanceFlow(t *testing.T) {
	cs := newPolicyContractClient(t)
	ctx := context.Background()
	start := startFixtureSession(t, cs)
	rev := revisionOf(t, start)

	// Evaluating with a blocking, uncited structural criterion resolves to
	// unverifiable, which blocks completion without override.
	evalRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "step_evaluate", Arguments: map[string]any{
		"session_id": start.SessionID, "expected_revision": rev,
		"criteria": []map[string]any{{"name": "compiles", "kind": "structural", "severity": "blocking"}},
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	eval := decodeEnvelope(t, evalRes)
	if eval.Status != "ok" || eval.ProgressEffect != ProgressEffectEvaluationRecorded {
		t.Fatalf("unexpected step_evaluate envelope: %+v", eval)
	}
	rev = revisionOf(t, eval)

	blockedRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "step_complete", Arguments: map[string]any{"session_id": start.SessionID, "expected_revision": rev}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !blockedRes.IsError {
		t.Fatal("expected IsError completing a step with a blocking evaluation failure")
	}
	// The error envelope carries no data; the rejected call still appended
	// its event (this codebase's established append-first-then-validate
	// ordering), so the caller must re-fetch the revision via session_get.
	getRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "session_get", Arguments: map[string]any{"session_id": start.SessionID}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rev = revisionOf(t, decodeEnvelope(t, getRes))

	completeRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "step_complete", Arguments: map[string]any{"session_id": start.SessionID, "override": true, "expected_revision": rev}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	complete := decodeEnvelope(t, completeRes)
	if complete.Status != "ok" || complete.ProgressEffect != ProgressEffectStepCompleted {
		t.Fatalf("unexpected step_complete envelope: %+v", complete)
	}
	rev = revisionOf(t, complete)

	advanceRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "step_advance", Arguments: map[string]any{"session_id": start.SessionID, "expected_revision": rev}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	advance := decodeEnvelope(t, advanceRes)
	if advance.Status != "ok" || advance.ProgressEffect != ProgressEffectStepAdvanced || advance.ActiveNode == nil || advance.ActiveNode.ID != "fixture.step-two" {
		t.Fatalf("unexpected step_advance envelope: %+v", advance)
	}
}

func TestContractReflectionRecord(t *testing.T) {
	cs := newPolicyContractClient(t)
	ctx := context.Background()
	start := startFixtureSession(t, cs)
	rev := revisionOf(t, start)

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "reflection_record", Arguments: map[string]any{
		"session_id": start.SessionID, "prompt": "why a pointer?", "answer": "to avoid copying", "expected_revision": rev,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env := decodeEnvelope(t, res)
	if env.Status != "ok" {
		t.Fatalf("reflection_record failed: %+v", env)
	}
}

func TestContractStepAdvanceRequiresSessionID(t *testing.T) {
	cs := newPolicyContractClient(t)
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "step_advance", Arguments: map[string]any{"expected_revision": 0}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected IsError when session_id is missing")
	}
}
