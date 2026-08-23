// Package security centralizes the local-state privacy operations
// PROJECT.md §22 requires — export and removal — plus a secret-shaped-
// content check other packages can use to prove nothing leaves this
// process's boundary unredacted (security-privacy-hardening requirement
// R4, R5, R6). It never touches a learner's actual project files or Git
// state: everything here operates only on `.codinho/state`, the
// workspace-local event log and evidence store.
package security

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/oseiaspereira88/codinho/internal/workspace"
)

// ErrPurgeNotConfirmed is returned by Purge when confirm is false, so a
// caller can never delete local state by accident (fail-closed default).
var ErrPurgeNotConfirmed = errors.New("security: purge requires explicit confirmation")

// StateDirName and its two contents are the only paths Export and Purge
// ever touch, matching cmd/codinho/main.go's own state layout.
const (
	StateDirName    = ".codinho"
	EventsFileName  = "events.jsonl"
	EvidenceDirName = "evidence"
)

// StatePath resolves the local state directory under workspaceRoot.
func StatePath(workspaceRoot string) string {
	return filepath.Join(workspaceRoot, StateDirName, "state")
}

// ContainsSecretShapedContent reports whether content matches any known
// secret pattern (requirement R4), reusing internal/workspace's own
// redaction patterns so there is exactly one place those patterns are
// authored.
func ContainsSecretShapedContent(content []byte) bool {
	return !bytes.Equal(workspace.Redact(content), content)
}

// Export copies the local state's event log and evidence store verbatim
// into destDir, for a learner who wants their own backup or to inspect
// what was recorded (requirement R6: "export... explícitos"). It never
// modifies the source, and never touches anything outside
// StatePath(workspaceRoot) — the learner's project files and Git history
// are never read or written here (requirement R9).
func Export(workspaceRoot, destDir string) error {
	stateDir := StatePath(workspaceRoot)
	if err := os.MkdirAll(destDir, 0o700); err != nil {
		return err
	}

	eventsSrc := filepath.Join(stateDir, EventsFileName)
	if _, err := os.Stat(eventsSrc); err == nil {
		if err := copyFile(eventsSrc, filepath.Join(destDir, EventsFileName), 0o600); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	evidenceSrc := filepath.Join(stateDir, EvidenceDirName)
	entries, err := os.ReadDir(evidenceSrc)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	evidenceDest := filepath.Join(destDir, EvidenceDirName)
	if err := os.MkdirAll(evidenceDest, 0o700); err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if err := copyFile(filepath.Join(evidenceSrc, e.Name()), filepath.Join(evidenceDest, e.Name()), 0o400); err != nil {
			return err
		}
	}
	return nil
}

// Purge removes the local state directory (events and evidence) under
// workspaceRoot entirely. It requires confirm to be explicitly true
// (fail-closed: no accidental destructive default) and only ever removes
// paths under StatePath — it never touches the learner's actual project
// files or `.git` (requirement R6, R9).
func Purge(workspaceRoot string, confirm bool) error {
	if !confirm {
		return ErrPurgeNotConfirmed
	}
	return os.RemoveAll(StatePath(workspaceRoot))
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}
