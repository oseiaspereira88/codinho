package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/oseiaspereira88/codinho/internal/config"
)

const manifestScaffold = `schema_version: 1
packs: []
`

// runInit scaffolds an empty packs/manifest.yaml under the workspace root
// so catalog validate/list/show and codinho serve have something to load
// (requirement R1). It never overwrites an existing manifest by default
// (requirement R5).
func runInit(args []string, stdout, stderr io.Writer) int {
	fs := newFlagSet("init", stderr)
	dest := fs.String("dest", ".", "workspace directory to initialize")
	force := fs.Bool("force", false, "overwrite an existing packs/manifest.yaml")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}

	cfg := config.Load()
	root := *dest
	if !filepath.IsAbs(root) {
		root = filepath.Join(cfg.WorkspaceRoot, root)
	}

	packsDir := filepath.Join(root, "packs")
	manifestPath := filepath.Join(packsDir, "manifest.yaml")

	if _, err := os.Stat(manifestPath); err == nil && !*force {
		fmt.Fprintf(stderr, "codinho: init: %s already exists; pass --force to overwrite\n", manifestPath)
		return exitError
	}

	if err := os.MkdirAll(packsDir, 0o700); err != nil {
		fmt.Fprintf(stderr, "codinho: init: %v\n", err)
		return exitError
	}
	if err := os.WriteFile(manifestPath, []byte(manifestScaffold), 0o600); err != nil {
		fmt.Fprintf(stderr, "codinho: init: %v\n", err)
		return exitError
	}

	fmt.Fprintf(stdout, "initialized workspace at %s\n", root)
	fmt.Fprintf(stdout, "created %s\n", manifestPath)
	return exitOK
}
