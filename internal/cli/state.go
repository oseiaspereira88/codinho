package cli

import (
	"path/filepath"

	"github.com/oseiaspereira88/codinho/internal/config"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
)

// openEventStore opens the workspace's event log for administrative
// commands. It never calls Append: every CLI command that touches the
// store is read-only (requirement R3, "sem mutação pedagógica"). This is
// intentionally separate from codinho serve's exclusive AcquireLock, so a
// read-only inspection can run without contending with a running server
// for the workspace lock.
func openEventStore(cfg config.Config) (*eventstore.Store, error) {
	stateDir := filepath.Join(cfg.WorkspaceRoot, ".codinho", "state")
	return eventstore.OpenReadOnly(filepath.Join(stateDir, "events.jsonl"), nil)
}
