package cli

import (
	"fmt"
	"io"

	"github.com/oseiaspereira88/codinho/internal/config"
	"github.com/oseiaspereira88/codinho/internal/security"
)

const privacyUsage = `Usage: codinho privacy <command>

Commands:
  export  Copy the local event log (including quarantined drafts) and evidence store to a destination
  purge   Remove all local state (event log (including quarantined drafts) and evidence store)
`

func runPrivacy(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, privacyUsage)
		return exitUsage
	}
	switch args[0] {
	case "export":
		return runPrivacyExport(args[1:], stdout, stderr)
	case "purge":
		return runPrivacyPurge(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "codinho: unknown privacy subcommand %q\n\n", args[0])
		fmt.Fprint(stderr, privacyUsage)
		return exitUsage
	}
}

// runPrivacyExport copies the local event log (including quarantined drafts) and evidence store
// verbatim to dest, never touching the learner's project files or Git
// state (security-privacy-hardening requirement R6, R9).
func runPrivacyExport(args []string, stdout, stderr io.Writer) int {
	fs := newFlagSet("privacy export", stderr)
	dest := fs.String("dest", "", "destination directory (required)")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if *dest == "" {
		fmt.Fprintln(stderr, "codinho: privacy export: --dest is required")
		return exitUsage
	}

	cfg := config.Load()
	if err := security.Export(cfg.WorkspaceRoot, *dest); err != nil {
		fmt.Fprintf(stderr, "codinho: privacy export: %v\n", err)
		return exitError
	}
	fmt.Fprintf(stdout, "exported local state to %s\n", *dest)
	return exitOK
}

// runPrivacyPurge removes only the local `.codinho/state` directory
// (event log (including quarantined drafts) and evidence store); it always requires --confirm and never
// touches the learner's project files (fail-closed default; requirement
// R6, R9).
func runPrivacyPurge(args []string, stdout, stderr io.Writer) int {
	fs := newFlagSet("privacy purge", stderr)
	confirm := fs.Bool("confirm", false, "required: confirms permanent removal of local state")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}

	cfg := config.Load()
	if err := security.Purge(cfg.WorkspaceRoot, *confirm); err != nil {
		fmt.Fprintf(stderr, "codinho: privacy purge: %v\n", err)
		return exitError
	}
	fmt.Fprintln(stdout, "removed local state (event log (including quarantined drafts) and evidence store)")
	return exitOK
}
