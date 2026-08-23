package cli

import (
	"fmt"
	"io"

	"github.com/oseiaspereira88/codinho/internal/config"
	"github.com/oseiaspereira88/codinho/internal/diagnostics"
)

// runDoctor reports installation diagnostics only. It never executes files
// from the project workspace and never dumps the process environment,
// source code or secret-shaped content (security requirement).
func runDoctor(args []string, stdout, stderr io.Writer) int {
	fs := newFlagSet("doctor", stderr)
	jsonOut := fs.Bool("json", false, "print the diagnostic report as JSON")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}

	cfg := config.Load()
	report := diagnostics.Run(cfg.WorkspaceRoot, version, resolveCommit())

	if *jsonOut {
		if err := writeJSON(stdout, report); err != nil {
			fmt.Fprintf(stderr, "codinho: doctor: %v\n", err)
			return exitError
		}
	} else {
		fmt.Fprintf(stdout, "version: %s (%s)\n", report.Version, report.Commit)
		fmt.Fprintf(stdout, "go: %s\n", report.GoVersion)
		fmt.Fprintf(stdout, "os/arch: %s/%s\n", report.OS, report.Arch)
		fmt.Fprintf(stdout, "workspace: %s\n", report.WorkspaceRoot)
		for _, c := range report.Checks {
			fmt.Fprintf(stdout, "%s: %s", c.Name, c.Status)
			if c.Detail != "" {
				fmt.Fprintf(stdout, " (%s)", c.Detail)
			}
			fmt.Fprintln(stdout)
		}
		if report.OK() {
			fmt.Fprintln(stdout, "status: ok")
		} else {
			fmt.Fprintln(stdout, "status: error")
		}
	}

	if !report.OK() {
		return exitError
	}
	return exitOK
}
