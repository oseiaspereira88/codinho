// Package cli implements the codinho command-line entry points.
package cli

import (
	"context"
	"fmt"
	"io"
)

// Stable process exit codes.
const (
	exitOK    = 0
	exitError = 1
	exitUsage = 2
)

const usage = `codinho - deliberate practice tutor runtime

Usage:
  codinho <command>

Commands:
  serve      Start the MCP server over stdio
  init       Scaffold an empty workspace (packs/manifest.yaml)
  doctor     Diagnose the local installation
  version    Print version and build information
  catalog    Validate, list and show curriculum items
  session    Inspect a session's recorded event history
  progress   Show and export recomputed mastery projections
  workspace  Materialize a challenge's starter fixture (prepare)
  privacy    Export or purge local state (event log and evidence)
  help       Show this help message
`

// Run dispatches args to the matching command and returns the process exit
// code. stdout carries command output; stderr carries usage and errors.
func Run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	select {
	case <-ctx.Done():
		fmt.Fprintln(stderr, "codinho: cancelled")
		return exitError
	default:
	}

	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return exitUsage
	}

	switch args[0] {
	case "serve":
		// The real entry point (cmd/codinho) intercepts "serve" before
		// calling Run, because the MCP stdio transport must own the
		// process's real stdout/stderr directly rather than the writers
		// passed here (main.go: "reserving stdout exclusively for the
		// protocol"). This case only exists so Run never reports "serve"
		// as an unknown command when a caller other than main.go passes
		// it here.
		fmt.Fprintln(stderr, "codinho: serve must be launched by the codinho binary directly (codinho serve), not through this dispatcher")
		return exitError
	case "init":
		return runInit(args[1:], stdout, stderr)
	case "version":
		return runVersion(stdout)
	case "doctor":
		return runDoctor(stdout)
	case "catalog":
		return runCatalog(args[1:], stdout, stderr)
	case "session":
		return runSession(args[1:], stdout, stderr)
	case "progress":
		return runProgress(args[1:], stdout, stderr)
	case "workspace":
		return runWorkspace(args[1:], stdout, stderr)
	case "privacy":
		return runPrivacy(args[1:], stdout, stderr)
	case "help", "-h", "--help":
		fmt.Fprint(stdout, usage)
		return exitOK
	default:
		fmt.Fprintf(stderr, "codinho: unknown command %q\n\n", args[0])
		fmt.Fprint(stderr, usage)
		return exitUsage
	}
}
