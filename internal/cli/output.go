package cli

import (
	"encoding/json"
	"flag"
	"io"
	"strings"
)

// writeJSON encodes v as indented JSON, one value per call, so every
// command's --json output is stable and diffable (non-functional
// requirement: "saída JSON deve ser versionada e testada por golden
// files").
func writeJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// newFlagSet builds a flag.FlagSet that writes its own usage/error output
// to stderr instead of the process's real stderr, and never calls
// os.Exit on a parse error (flag.ContinueOnError), so a bad flag becomes
// an ordinary exitUsage return like every other CLI error path.
func newFlagSet(name string, stderr io.Writer) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(stderr)
	return fs
}

// extractPositional pulls non-flag tokens out of args so a positional
// argument may appear anywhere relative to flags (e.g. both
// "prepare <id> --dest x" and "prepare --dest x <id>" parse the same way),
// which the standard library's flag package does not support on its own
// (it stops parsing at the first non-flag token). valueFlags names every
// flag that consumes the following token as its value; any other
// "-"-prefixed token is treated as a standalone (typically boolean) flag.
func extractPositional(args []string, valueFlags map[string]bool) (positional, flagArgs []string) {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if !strings.HasPrefix(a, "-") {
			positional = append(positional, a)
			continue
		}
		flagArgs = append(flagArgs, a)
		name := strings.TrimLeft(a, "-")
		if strings.Contains(name, "=") {
			continue
		}
		if valueFlags[name] && i+1 < len(args) {
			i++
			flagArgs = append(flagArgs, args[i])
		}
	}
	return positional, flagArgs
}
