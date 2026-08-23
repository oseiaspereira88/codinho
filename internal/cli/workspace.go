package cli

import (
	"fmt"
	"io"

	"github.com/oseiaspereira88/codinho/internal/config"
	"github.com/oseiaspereira88/codinho/internal/fixtures"
)

const workspaceUsage = `Usage: codinho workspace <command>

Commands:
  prepare  Materialize a challenge's starter fixture into a destination
`

func runWorkspace(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, workspaceUsage)
		return exitUsage
	}
	switch args[0] {
	case "prepare":
		return runWorkspacePrepare(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "codinho: unknown workspace subcommand %q\n\n", args[0])
		fmt.Fprint(stderr, workspaceUsage)
		return exitUsage
	}
}

// runWorkspacePrepare is invoked directly by the operator, never by an MCP
// tool: writing starter code into a learner's project always requires an
// explicit, out-of-band command with an explicit destination
// (administrative-cli-fixtures Decision 1; Constraint: "workspace prepare
// é acionado diretamente pelo usuário, não por tool MCP").
func runWorkspacePrepare(args []string, stdout, stderr io.Writer) int {
	fs := newFlagSet("workspace prepare", stderr)
	dest := fs.String("dest", "", "destination directory (required)")
	dryRun := fs.Bool("dry-run", false, "show the plan without writing anything")
	overwrite := fs.Bool("overwrite", false, "allow overwriting existing files at dest")
	jsonOut := fs.Bool("json", false, "print the plan/result as JSON")

	positional, flagArgs := extractPositional(args, map[string]bool{"dest": true})
	if err := fs.Parse(flagArgs); err != nil {
		return exitUsage
	}
	if len(positional) != 1 {
		fmt.Fprint(stderr, "Usage: codinho workspace prepare <challenge-id> --dest <path> [--dry-run] [--overwrite] [--json]\n")
		return exitUsage
	}
	if *dest == "" {
		fmt.Fprintln(stderr, "codinho: workspace prepare: --dest is required")
		return exitUsage
	}
	challengeID := positional[0]

	cfg := config.Load()
	catalog, _, err := loadCatalog(cfg)
	if err != nil {
		fmt.Fprintf(stderr, "codinho: workspace prepare: %v\n", err)
		return exitError
	}
	if catalog == nil {
		fmt.Fprintln(stderr, "codinho: workspace prepare: catalog failed to load; run catalog validate for details")
		return exitError
	}
	challenge, ok := catalog.Challenge(challengeID)
	if !ok {
		fmt.Fprintf(stderr, "codinho: workspace prepare: challenge %q not found\n", challengeID)
		return exitError
	}

	opts := fixtures.Options{Overwrite: *overwrite}

	if *dryRun {
		plan, err := fixtures.Plan(challenge, *dest, opts)
		if err != nil && err != fixtures.ErrWouldOverwrite {
			fmt.Fprintf(stderr, "codinho: workspace prepare: %v\n", err)
			return exitError
		}
		if *jsonOut {
			if encErr := writeJSON(stdout, plan); encErr != nil {
				fmt.Fprintf(stderr, "codinho: workspace prepare: %v\n", encErr)
				return exitError
			}
		} else {
			printPlan(stdout, plan)
		}
		if err == fixtures.ErrWouldOverwrite {
			fmt.Fprintln(stderr, "codinho: workspace prepare: destination has conflicts; pass --overwrite to allow")
			return exitError
		}
		return exitOK
	}

	manifest, err := fixtures.Materialize(challenge, *dest, opts)
	if err != nil {
		if err == fixtures.ErrWouldOverwrite {
			fmt.Fprintln(stderr, "codinho: workspace prepare: destination has conflicts; pass --overwrite to allow")
		} else {
			fmt.Fprintf(stderr, "codinho: workspace prepare: %v\n", err)
		}
		return exitError
	}

	if *jsonOut {
		if err := writeJSON(stdout, manifest); err != nil {
			fmt.Fprintf(stderr, "codinho: workspace prepare: %v\n", err)
			return exitError
		}
		return exitOK
	}
	fmt.Fprintf(stdout, "prepared %d file(s) for %s at %s\n", len(manifest.Files), manifest.ChallengeID, *dest)
	for _, f := range manifest.Files {
		fmt.Fprintf(stdout, "  %s (%d bytes, sha256:%s)\n", f.Path, f.Size, f.SHA256)
	}
	return exitOK
}

func printPlan(stdout io.Writer, plan fixtures.FixturePlan) {
	fmt.Fprintf(stdout, "plan for %s -> %s\n", plan.ChallengeID, plan.Dest)
	for _, f := range plan.Files {
		conflict := ""
		if f.Exists {
			conflict = " (conflict)"
		}
		fmt.Fprintf(stdout, "  %s (%d bytes, sha256:%s)%s\n", f.Path, f.Size, f.SHA256, conflict)
	}
}
