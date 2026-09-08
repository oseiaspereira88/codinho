// Command codinho is the administrative CLI and MCP server entry point for
// the codinho runtime. Administrative commands are implemented in
// internal/cli; "serve" starts the MCP server over stdio (PROJECT.md §15),
// reserving stdout exclusively for the protocol.
package main

import (
	"context"
	"flag"
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
	"github.com/oseiaspereira88/codinho/internal/evidence"
	"github.com/oseiaspereira88/codinho/internal/mcpserver"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	args := os.Args[1:]
	if len(args) > 0 && args[0] == "serve" {
		if err := runServe(ctx, os.Stderr, args[1:]...); err != nil {
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
func runServe(ctx context.Context, stderr *os.File, args ...string) error {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.SetOutput(stderr)
	authoring := fs.Bool("authoring", false, "include local drafts for authoring/playtest")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("unexpected serve arguments")
	}
	cfg := config.Load()

	packs, diags, err := curriculum.LoadPacks(filepath.Join(cfg.WorkspaceRoot, "packs"), curriculum.DefaultLimits)
	if err != nil {
		return fmt.Errorf("loading catalog: %w", err)
	}
	diags = append(diags, curriculum.ValidatePublication(packs)...)
	for _, d := range diags {
		fmt.Fprintf(stderr, "catalog diagnostic: %+v\n", d)
	}
	for _, d := range diags {
		if d.Blocking {
			return fmt.Errorf("catalog failed to load: see diagnostics above")
		}
	}
	if !*authoring {
		packs = curriculum.PublishedPacks(packs)
	}
	catalog := curriculum.NewCatalogFromPacks(packs)

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

	evidenceStore, err := evidence.Open(filepath.Join(stateDir, "evidence"))
	if err != nil {
		return fmt.Errorf("opening evidence store: %w", err)
	}

	catalogService := application.NewCatalogService(catalog)
	sessionService := application.NewSessionService(catalogService, store, evidenceStore)
	assistanceService := application.NewAssistanceService(catalogService)
	workspaceService := application.NewWorkspaceService(store, evidenceStore)
	checksService := application.NewChecksService(store, sessionService, workspaceService)
	progressService := application.NewProgressService(store)
	recommendationService := application.NewRecommendationService(catalog, progressService)

	server := mcpserver.New(mcpserver.Deps{
		Catalog:        catalogService,
		Session:        sessionService,
		Assistance:     assistanceService,
		Workspace:      workspaceService,
		Checks:         checksService,
		Progress:       progressService,
		Recommendation: recommendationService,
	}, stderr)
	return mcpserver.Run(ctx, server, &mcp.StdioTransport{})
}
