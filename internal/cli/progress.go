package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/oseiaspereira88/codinho/internal/application"
	"github.com/oseiaspereira88/codinho/internal/config"
)

const progressUsage = `Usage: codinho progress <command>

Commands:
  show    Print recomputed mastery projections
  export  Write recomputed mastery projections to a file
`

func runProgress(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, progressUsage)
		return exitUsage
	}
	switch args[0] {
	case "show":
		return runProgressShow(args[1:], stdout, stderr)
	case "export":
		return runProgressExport(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "codinho: unknown progress subcommand %q\n\n", args[0])
		fmt.Fprint(stderr, progressUsage)
		return exitUsage
	}
}

func runProgressShow(args []string, stdout, stderr io.Writer) int {
	fs := newFlagSet("progress show", stderr)
	competency := fs.String("competency", "", "restrict to one competency id")
	jsonOut := fs.Bool("json", false, "print result as JSON")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}

	result, err := recomputeProgress(*competency)
	if err != nil {
		fmt.Fprintf(stderr, "codinho: progress show: %v\n", err)
		return exitError
	}

	if *jsonOut {
		if err := writeJSON(stdout, result); err != nil {
			fmt.Fprintf(stderr, "codinho: progress show: %v\n", err)
			return exitError
		}
		return exitOK
	}

	competencyIDs := make([]string, 0, len(result.Competencies))
	for id := range result.Competencies {
		competencyIDs = append(competencyIDs, id)
	}
	sort.Strings(competencyIDs)
	for _, id := range competencyIDs {
		dims := result.Competencies[id]
		dimIDs := make([]string, 0, len(dims))
		for d := range dims {
			dimIDs = append(dimIDs, d)
		}
		sort.Strings(dimIDs)
		for _, d := range dimIDs {
			proj := dims[d]
			fmt.Fprintf(stdout, "%s\t%s\t%s\n", id, d, proj.State)
		}
	}
	fmt.Fprintf(stdout, "revision: %d\n", result.Revision)
	return exitOK
}

func runProgressExport(args []string, stdout, stderr io.Writer) int {
	fs := newFlagSet("progress export", stderr)
	out := fs.String("out", "", "file to write the JSON export to (required)")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if *out == "" {
		fmt.Fprintln(stderr, "codinho: progress export: --out is required")
		return exitUsage
	}

	result, err := recomputeProgress("")
	if err != nil {
		fmt.Fprintf(stderr, "codinho: progress export: %v\n", err)
		return exitError
	}

	payload := struct {
		ExportedAt string                     `json:"exported_at"`
		Revision   uint64                     `json:"revision"`
		Result     application.ProgressResult `json:"result"`
	}{ExportedAt: time.Now().UTC().Format(time.RFC3339), Revision: result.Revision, Result: result}

	encoded, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		fmt.Fprintf(stderr, "codinho: progress export: %v\n", err)
		return exitError
	}
	if err := os.WriteFile(*out, encoded, 0o600); err != nil {
		fmt.Fprintf(stderr, "codinho: progress export: %v\n", err)
		return exitError
	}
	fmt.Fprintf(stdout, "exported progress to %s\n", *out)
	return exitOK
}

func recomputeProgress(competencyID string) (application.ProgressResult, error) {
	cfg := config.Load()
	store, err := openEventStore(cfg)
	if err != nil {
		return application.ProgressResult{}, err
	}
	defer store.Close()

	progress := application.NewProgressService(store)
	return progress.Progress(competencyID)
}
