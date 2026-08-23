package mcpserver

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/oseiaspereira88/codinho/internal/application"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/evidence"
)

// TestLoggingMiddlewareEmitsStructuredLineWithoutPayload proves
// reliability-observability-compatibility requirement R3: one line per
// call, to stderr only, carrying tool/status/duration/size but never the
// request or response payload.
func TestLoggingMiddlewareEmitsStructuredLineWithoutPayload(t *testing.T) {
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
	workspaceService := application.NewWorkspaceService(store, evidenceStore)
	checksService := application.NewChecksService(store, sessionService, workspaceService)
	progressService := application.NewProgressService(store)

	var stderr bytes.Buffer
	server := New(Deps{
		Catalog: catalogService, Session: sessionService,
		Assistance: application.NewAssistanceService(catalogService),
		Workspace:  workspaceService, Checks: checksService, Progress: progressService,
		Recommendation: application.NewRecommendationService(catalog, progressService),
	}, &stderr)

	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.Run(ctx, serverTransport)

	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "0.0.0"}, nil)
	cs, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("unexpected error connecting client: %v", err)
	}
	defer cs.Close()

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "catalog_search", Arguments: map[string]any{"text": "fixture"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.IsError {
		t.Fatalf("unexpected tool error: %+v", res)
	}

	log := stderr.String()
	if !strings.Contains(log, "tool=catalog_search") {
		t.Fatalf("expected a log line naming the tool, got: %q", log)
	}
	if !strings.Contains(log, "status=ok") {
		t.Fatalf("expected status=ok, got: %q", log)
	}
	if !strings.Contains(log, "correlation_id=") || !strings.Contains(log, "duration_ms=") || !strings.Contains(log, "size_bytes=") {
		t.Fatalf("expected correlation_id/duration_ms/size_bytes fields, got: %q", log)
	}
	if strings.Contains(log, "fixture.challenge-one") {
		t.Fatalf("log must never carry response payload content, got: %q", log)
	}
}
