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

// newNestedContractClient is like newContractClient but its fixture
// challenge has a macro → meso → micro chain, so granularity_adjust has
// something real to walk (requirement R6).
func newNestedContractClient(t *testing.T) *mcp.ClientSession {
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
          - id: fixture.macro-one
            kind: macro
            title: Macro step
            instruction:
              objective: Understand the fixture.
              scope: Read only.
            children:
              - id: fixture.meso-one
                kind: meso
                title: Meso step
                instruction:
                  objective: Meso objective.
                  scope: Meso scope.
                children:
                  - id: fixture.micro-one
                    kind: micro
                    title: Micro step
                    instruction:
                      objective: Declare the fixture type.
                      scope: Only the declaration.
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

	client := mcp.NewClient(&mcp.Implementation{Name: "session-contract-test-client", Version: "0.0.0"}, nil)
	cs, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("unexpected error connecting client: %v", err)
	}
	t.Cleanup(func() { cs.Close() })
	return cs
}

func startFixtureSession(t *testing.T, cs *mcp.ClientSession) Envelope {
	t.Helper()
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "session_start", Arguments: map[string]any{"challenge_id": fixtureChallengeID}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env := decodeEnvelope(t, res)
	if env.Status != "ok" {
		t.Fatalf("session_start failed: %+v", env)
	}
	return env
}

func revisionOf(t *testing.T, env Envelope) uint64 {
	t.Helper()
	data, ok := env.Data.(map[string]any)
	if !ok {
		t.Fatalf("envelope data is not an object: %+v", env.Data)
	}
	rev, ok := data["revision"].(float64) // JSON numbers decode as float64
	if !ok {
		t.Fatalf("envelope data has no numeric revision: %+v", data)
	}
	return uint64(rev)
}

func TestContractSessionConfigureChangesHelpPolicy(t *testing.T) {
	cs := newNestedContractClient(t)
	ctx := context.Background()
	start := startFixtureSession(t, cs)
	rev := revisionOf(t, start)

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "session_configure", Arguments: map[string]any{
		"session_id": start.SessionID, "help": "no_hints", "expected_revision": rev,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env := decodeEnvelope(t, res)
	if env.Status != "ok" {
		t.Fatalf("session_configure failed: %+v", env)
	}
}

func TestContractSessionConfigureRejectsStaleRevision(t *testing.T) {
	cs := newNestedContractClient(t)
	ctx := context.Background()
	start := startFixtureSession(t, cs)
	rev := revisionOf(t, start)

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "session_configure", Arguments: map[string]any{
		"session_id": start.SessionID, "help": "free", "expected_revision": rev + 99,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected IsError for a stale expected_revision")
	}
	env := decodeEnvelope(t, res)
	if env.Error == nil || env.Error.Code != string(ErrCodeStateConflict) {
		t.Fatalf("unexpected envelope: %+v", env)
	}
	if !env.Error.Retryable {
		t.Fatal("a revision conflict should be reported as retryable (fetch and retry)")
	}
}

func TestContractSessionPauseResumeFinishLifecycle(t *testing.T) {
	cs := newNestedContractClient(t)
	ctx := context.Background()
	start := startFixtureSession(t, cs)
	rev := revisionOf(t, start)

	pauseRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "session_pause", Arguments: map[string]any{"session_id": start.SessionID, "expected_revision": rev}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	paused := decodeEnvelope(t, pauseRes)
	if paused.Status != "ok" {
		t.Fatalf("session_pause failed: %+v", paused)
	}
	rev = revisionOf(t, paused)

	resumeRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "session_resume", Arguments: map[string]any{"session_id": start.SessionID, "expected_revision": rev}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resumed := decodeEnvelope(t, resumeRes)
	if resumed.Status != "ok" {
		t.Fatalf("session_resume failed: %+v", resumed)
	}
	rev = revisionOf(t, resumed)

	finishRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "session_finish", Arguments: map[string]any{"session_id": start.SessionID, "expected_revision": rev}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	finished := decodeEnvelope(t, finishRes)
	if finished.Status != "ok" {
		t.Fatalf("session_finish failed: %+v", finished)
	}
}

func TestContractGranularityAdjustWalksToMicro(t *testing.T) {
	cs := newNestedContractClient(t)
	ctx := context.Background()
	start := startFixtureSession(t, cs)
	rev := revisionOf(t, start)

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "granularity_adjust", Arguments: map[string]any{
		"session_id": start.SessionID, "depth": "micro", "expected_revision": rev,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env := decodeEnvelope(t, res)
	if env.Status != "ok" || env.ActiveNode == nil || env.ActiveNode.ID != "fixture.micro-one" {
		t.Fatalf("unexpected envelope: %+v", env)
	}

	instrRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "instruction_get", Arguments: map[string]any{"session_id": start.SessionID}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	instr := decodeEnvelope(t, instrRes)
	if instr.ActiveNode == nil || instr.ActiveNode.ID != "fixture.micro-one" {
		t.Fatalf("instruction_get did not follow the granularity adjustment: %+v", instr)
	}
}

func TestContractGranularityAdjustAcceptsReason(t *testing.T) {
	cs := newNestedContractClient(t)
	ctx := context.Background()
	start := startFixtureSession(t, cs)
	rev := revisionOf(t, start)

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "granularity_adjust", Arguments: map[string]any{
		"session_id": start.SessionID, "depth": "micro", "reason": "3 evidências sem ajuda", "expected_revision": rev,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env := decodeEnvelope(t, res)
	if env.Status != "ok" {
		t.Fatalf("unexpected envelope: %+v", env)
	}
}

func TestContractLearnerNextStepProposeRecordsWithoutAdvancing(t *testing.T) {
	cs := newNestedContractClient(t)
	ctx := context.Background()
	start := startFixtureSession(t, cs)
	rev := revisionOf(t, start)

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "learner_next_step_propose", Arguments: map[string]any{
		"session_id": start.SessionID, "step_id": "fixture.micro-one", "expected_revision": rev,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env := decodeEnvelope(t, res)
	if env.Status != "ok" {
		t.Fatalf("unexpected envelope: %+v", env)
	}

	getRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "session_get", Arguments: map[string]any{"session_id": start.SessionID}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	get := decodeEnvelope(t, getRes)
	if get.ActiveNode == nil || get.ActiveNode.ID == "fixture.micro-one" {
		t.Fatalf("proposing a step must never advance the active node: %+v", get)
	}
}

func TestContractLearnerNextStepProposeRejectsUnknownStep(t *testing.T) {
	cs := newNestedContractClient(t)
	ctx := context.Background()
	start := startFixtureSession(t, cs)
	rev := revisionOf(t, start)

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "learner_next_step_propose", Arguments: map[string]any{
		"session_id": start.SessionID, "step_id": "does-not-exist", "expected_revision": rev,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected IsError for an unknown step_id")
	}
}

func TestContractGranularityAdjustRequiresDepth(t *testing.T) {
	cs := newNestedContractClient(t)
	ctx := context.Background()
	start := startFixtureSession(t, cs)
	rev := revisionOf(t, start)

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "granularity_adjust", Arguments: map[string]any{
		"session_id": start.SessionID, "expected_revision": rev,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected IsError when depth is missing")
	}
}
