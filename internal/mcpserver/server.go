package mcpserver

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync/atomic"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/oseiaspereira88/codinho/internal/application"
)

// Deps are the already-opened dependencies a Server needs. The caller
// (cmd/codinho) owns their lifecycle: opening the catalog and event store,
// acquiring the workspace lock, and closing/releasing them after Run
// returns.
type Deps struct {
	Catalog        *application.CatalogService
	Session        *application.SessionService
	Assistance     *application.AssistanceService
	Workspace      *application.WorkspaceService
	Checks         *application.ChecksService
	Progress       *application.ProgressService
	Recommendation *application.RecommendationService
}

// New builds an MCP server exposing the minimal vertical slice of tools
// (requirement R5) with stderr-only logging, so stdout stays reserved for
// the protocol (requirement R1; security).
func New(deps Deps, stderr io.Writer) *mcp.Server {
	logger := slog.New(slog.NewTextHandler(stderr, nil))
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "codinho",
		Version: "0.1.0",
	}, &mcp.ServerOptions{
		Instructions: Instructions,
		Logger:       logger,
	})
	server.AddReceivingMiddleware(loggingMiddleware(logger))

	registerCatalogTools(server, deps.Catalog)
	registerSessionTools(server, deps.Session)
	registerAssistanceTools(server, deps.Session, deps.Assistance)
	registerAssessmentTools(server, deps.Session)
	registerWorkspaceTools(server, deps.Workspace)
	registerCheckTools(server, deps.Checks)
	registerProgressTools(server, deps.Progress)
	registerRecommendationTools(server, deps.Recommendation)

	return server
}

// Run starts server on transport and blocks until stdin closes or ctx is
// canceled, either of which is a clean shutdown, not a failure
// (requirement R1).
func Run(ctx context.Context, server *mcp.Server, transport mcp.Transport) error {
	err := server.Run(ctx, transport)
	if err == nil || ctx.Err() != nil {
		return nil
	}
	return err
}

var requestCounter atomic.Uint64

// requestIDFor assigns a per-call, server-local request ID for envelope
// reporting (PROJECT.md §15.3).
func requestIDFor(_ *mcp.CallToolRequest) string {
	return fmt.Sprintf("req_%d", requestCounter.Add(1))
}

// errorResult marks a CallToolResult as a domain-level error while leaving
// Content/StructuredContent for the SDK to populate from the handler's
// typed Envelope return value (PROJECT.md §15.4: errors are reported
// inside tool content, not as an MCP protocol-level error).
func errorResult() *mcp.CallToolResult {
	return &mcp.CallToolResult{IsError: true}
}
