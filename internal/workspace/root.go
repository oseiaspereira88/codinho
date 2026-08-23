// Package workspace observes a learner's workspace read-only: it resolves
// and contains an authorized root, captures baselines and diffs relevant
// to a declared set of globs, and never edits the learner's files (PROJECT.md
// §11.2; ADR learner-code-ownership-and-safe-checks).
package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// ErrPathEscapesRoot is returned when a resolved path would leave the
// authorized root, whether through ".." segments or a symlink encountered
// along the way (requirement R1; security).
var ErrPathEscapesRoot = errors.New("workspace: path escapes authorized root")

// ErrNotADirectory is returned by AuthorizeRoot when path does not resolve
// to a directory.
var ErrNotADirectory = errors.New("workspace: root is not a directory")

// Root is an authorized, symlink-resolved real directory. Every path this
// package touches is checked against it before any filesystem access.
type Root struct {
	real string
}

// AuthorizeRoot resolves path to its real (symlink-evaluated) absolute
// form and returns a Root scoped to it. It fails if path does not exist or
// is not a directory (requirement R1).
func AuthorizeRoot(path string) (Root, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return Root{}, err
	}
	real, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return Root{}, err
	}
	info, err := os.Stat(real)
	if err != nil {
		return Root{}, err
	}
	if !info.IsDir() {
		return Root{}, ErrNotADirectory
	}
	return Root{real: real}, nil
}

// Path returns the root's real absolute path.
func (r Root) Path() string { return r.real }

// Resolve joins rel onto the root and returns its real absolute path,
// rejecting any result that escapes the root — via ".." segments (caught
// before touching the filesystem) or a symlink encountered while resolving
// an existing target (caught by re-checking after evaluation, closing the
// TOCTOU window a naive prefix check alone would leave open) (requirement
// R1; security).
func (r Root) Resolve(rel string) (string, error) {
	if rel == "" || rel == "." {
		return r.real, nil
	}
	joined := filepath.Join(r.real, rel)
	if !isWithin(r.real, joined) {
		return "", ErrPathEscapesRoot
	}
	real, err := filepath.EvalSymlinks(joined)
	if err != nil {
		return "", err
	}
	if !isWithin(r.real, real) {
		return "", ErrPathEscapesRoot
	}
	return real, nil
}

// Rel returns path relative to the root using forward slashes, regardless
// of host OS, so globs authored once behave the same on every platform.
func (r Root) Rel(path string) (string, error) {
	rel, err := filepath.Rel(r.real, path)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(rel), nil
}

func isWithin(root, candidate string) bool {
	if candidate == root {
		return true
	}
	return strings.HasPrefix(candidate, root+string(filepath.Separator))
}
