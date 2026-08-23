package workspace

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"
	"time"
)

// envPath returns the host PATH so git can be found, without inheriting
// any other environment variable from this process (least privilege).
func envPath() string { return os.Getenv("PATH") }

// StatusEntry is one line of `git status --porcelain`.
type StatusEntry struct {
	Path   string // root-relative, forward-slash
	Status string // e.g. "M", "??", "A", "D"
}

// GitInfo captures a git-backed root's identity at observation time
// (requirement R2).
type GitInfo struct {
	Commit string
	Dirty  bool
	Status []StatusEntry
}

// gitTimeout bounds every git invocation this package makes: it is a
// read-only identity probe, never the learner's own test run, so it must
// never be allowed to hang the caller.
const gitTimeout = 5 * time.Second

// gitInfo probes root for a Git identity. ok is false when root is not
// inside a Git working tree, so the caller can fall back to hash-based
// baselines (requirement R5) instead of treating the absence of Git as an
// error.
func gitInfo(root Root) (info GitInfo, ok bool, err error) {
	if _, gitErr := runGit(root, "rev-parse", "--is-inside-work-tree"); gitErr != nil {
		return GitInfo{}, false, nil
	}
	commit, err := runGit(root, "rev-parse", "HEAD")
	if err != nil {
		// A Git working tree with no commits yet: still usable, just with
		// an empty identity and everything reported as untracked/dirty.
		commit = ""
	}
	statusOut, err := runGit(root, "status", "--porcelain")
	if err != nil {
		return GitInfo{}, false, err
	}
	entries := parseStatus(statusOut)
	return GitInfo{Commit: strings.TrimSpace(commit), Dirty: len(entries) > 0, Status: entries}, true, nil
}

// runGit executes a fixed git subcommand directly via os/exec — never
// through a shell, and never with any argument derived from learner or
// caller input beyond the already-authorized root itself (ADR learner-
// code-ownership-and-safe-checks; security requirement).
func runGit(root Root, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gitTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = root.Path()
	cmd.Env = []string{"PATH=" + envPath()}
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}

// maxDiffBytes bounds how much diff text this package ever returns, so a
// single huge change can never make evidence unbounded (non-functional
// requirement: outputs must be limited).
const maxDiffBytes = 64 * 1024

// DiffText returns the working-tree diff for paths against HEAD, redacted
// and truncated to maxDiffBytes, when root has a Git identity. ok is false
// when root is not a Git working tree or there is nothing to diff, so the
// caller falls back to the hash-level FileChange list alone (requirement
// R4, R5).
func DiffText(root Root, paths []string) (diff string, ok bool, err error) {
	if len(paths) == 0 {
		return "", false, nil
	}
	if _, gitErr := runGit(root, "rev-parse", "--is-inside-work-tree"); gitErr != nil {
		return "", false, nil
	}
	args := append([]string{"diff", "--unified=3", "--"}, paths...)
	out, err := runGit(root, args...)
	if err != nil {
		return "", false, err
	}
	if out == "" {
		return "", false, nil
	}
	redacted := redact([]byte(out))
	if len(redacted) > maxDiffBytes {
		redacted = append(redacted[:maxDiffBytes], []byte("\n[TRUNCATED]")...)
	}
	return string(redacted), true, nil
}

func parseStatus(porcelain string) []StatusEntry {
	var entries []StatusEntry
	for line := range strings.SplitSeq(porcelain, "\n") {
		if len(line) < 4 {
			continue
		}
		status := strings.TrimSpace(line[:2])
		path := line[3:]
		// A rename/copy line reports "old -> new"; the new path is what
		// exists in the working tree now.
		if idx := strings.Index(path, " -> "); idx >= 0 {
			path = path[idx+4:]
		}
		entries = append(entries, StatusEntry{Path: path, Status: status})
	}
	return entries
}
