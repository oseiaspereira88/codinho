// Command codinho is the administrative CLI and MCP server entry point for
// the codinho runtime. Administrative commands are implemented in
// internal/cli; "serve" starts the MCP server over stdio (PROJECT.md §15),
// reserving stdout exclusively for the protocol.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/oseiaspereira88/codinho/internal/application"
	"github.com/oseiaspereira88/codinho/internal/cli"
	"github.com/oseiaspereira88/codinho/internal/config"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/mcpserver"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	args := os.Args[1:]
	if len(args) > 0 && args[0] == "serve" {
		if err := runServe(ctx, os.Stderr); err != nil {
			fmt.Fprintf(os.Stderr, "codinho serve: %v\n", err)
			os.Exit(1)
		}
		return
	}

	os.Exit(cli.Run(ctx, args, os.Stdout, os.Stderr))
}

// runServe wires the catalog, event store and workspace lock, then blocks
// running the MCP server until stdin closes or ctx is canceled
// (requirement R1). It never writes to stdout: mcp.StdioTransport owns
// stdout for the protocol, and every diagnostic here goes to stderr.
func runServe(ctx context.Context, stderr *os.File) error {
	cfg := config.Load()

	catalog, diags, err := curriculum.Load(filepath.Join(cfg.WorkspaceRoot, "packs"), curriculum.DefaultLimits)
	if err != nil {
		return fmt.Errorf("loading catalog: %w", err)
	}
	for _, d := range diags {
		fmt.Fprintf(stderr, "catalog diagnostic: %+v\n", d)
	}
	if catalog == nil {
		return fmt.Errorf("catalog failed to load: see diagnostics above")
	}

	stateDir := filepath.Join(cfg.WorkspaceRoot, ".codinho", "state")
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		return fmt.Errorf("preparing state directory: %w", err)
	}

	lock, err := eventstore.AcquireLock(filepath.Join(stateDir, "lock"))
	if err != nil {
		return fmt.Errorf("acquiring workspace lock: %w", err)
	}
	defer lock.Release()

	store, err := eventstore.Open(filepath.Join(stateDir, "events.jsonl"), nil)
	if err != nil {
		return fmt.Errorf("opening event store: %w", err)
	}
	defer store.Close()

	catalogService := application.NewCatalogService(catalog)
	sessionService := application.NewSessionService(catalogService, store)
	assistanceService := application.NewAssistanceService(catalogService)

	server := mcpserver.New(mcpserver.Deps{Catalog: catalogService, Session: sessionService, Assistance: assistanceService}, stderr)
	return mcpserver.Run(ctx, server, &mcp.StdioTransport{})
}
