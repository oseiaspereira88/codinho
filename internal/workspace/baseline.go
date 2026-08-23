package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
)

// Baseline captures a root's relevant state at one point in time: the Git
// identity when available, content hashes otherwise, restricted to the
// declared globs (requirements R2, R5).
type Baseline struct {
	Commit   string // "" when the root has no Git identity
	Dirty    bool
	FromGit  bool
	Hashes   map[string]string // rel path -> "sha256-<hex>"
	Fixtures []string          // rel paths present at capture time
}

// Capture walks root, restricted to globs, and records a Baseline. It
// never writes to root: no mtime, Git index or working-tree change
// (non-functional requirement).
func Capture(root Root, globs []string) (Baseline, error) {
	files, err := walkMatched(root, globs)
	if err != nil {
		return Baseline{}, err
	}
	b := Baseline{Hashes: map[string]string{}}
	fixtures := make([]string, 0, len(files))
	for _, f := range files {
		b.Hashes[f.Path] = f.Hash
		fixtures = append(fixtures, f.Path)
	}
	sort.Strings(fixtures)
	b.Fixtures = fixtures

	if info, ok, gitErr := gitInfo(root); gitErr == nil && ok {
		b.Commit = info.Commit
		b.Dirty = info.Dirty
		b.FromGit = true
	}
	return b, nil
}

// FileChange is one file's status between a Baseline and a later capture.
type FileChange struct {
	Path   string
	Change string // "added", "modified" or "removed"
	Hash   string // new content hash; "" when Change is "removed"
}

const (
	ChangeAdded    = "added"
	ChangeModified = "modified"
	ChangeRemoved  = "removed"
)

// ObserveResult is the relevant diff since baseline, restricted to globs,
// never attributing a change that predates the baseline to the step
// (requirement R4).
type ObserveResult struct {
	Baseline    Baseline
	Current     Baseline
	Changes     []FileChange
	Fingerprint string
}

// Observe captures root's current state and diffs it against baseline.
func Observe(root Root, globs []string, baseline Baseline) (ObserveResult, error) {
	current, err := Capture(root, globs)
	if err != nil {
		return ObserveResult{}, err
	}
	var changes []FileChange
	for path, hash := range current.Hashes {
		if prev, existed := baseline.Hashes[path]; !existed {
			changes = append(changes, FileChange{Path: path, Change: ChangeAdded, Hash: hash})
		} else if prev != hash {
			changes = append(changes, FileChange{Path: path, Change: ChangeModified, Hash: hash})
		}
	}
	for path := range baseline.Hashes {
		if _, still := current.Hashes[path]; !still {
			changes = append(changes, FileChange{Path: path, Change: ChangeRemoved})
		}
	}
	sort.Slice(changes, func(i, j int) bool { return changes[i].Path < changes[j].Path })
	return ObserveResult{
		Baseline:    baseline,
		Current:     current,
		Changes:     changes,
		Fingerprint: FingerprintOf(current),
	}, nil
}

// FingerprintOf returns b's stable digest directly, for a caller that
// already holds a Baseline (e.g. one just established by Capture) and
// does not need to walk the filesystem again via Fingerprint.
func FingerprintOf(b Baseline) string { return fingerprintOf(b) }

// Fingerprint recaptures root's current state restricted to globs and
// returns its fingerprint, without needing a full Baseline in hand. A
// caller compares this against a previously recorded fingerprint to
// detect that the workspace moved after evidence was collected
// (requirement R9).
func Fingerprint(root Root, globs []string) (string, error) {
	current, err := Capture(root, globs)
	if err != nil {
		return "", err
	}
	return fingerprintOf(current), nil
}

// fingerprintOf hashes the sorted (path, hash) pairs of a Baseline into
// one stable digest, so two Baselines with the same content compare equal
// regardless of walk order.
func fingerprintOf(b Baseline) string {
	paths := make([]string, 0, len(b.Hashes))
	for p := range b.Hashes {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	h := sha256.New()
	for _, p := range paths {
		h.Write([]byte(p))
		h.Write([]byte{0})
		h.Write([]byte(b.Hashes[p]))
		h.Write([]byte{0})
	}
	return "sha256-" + hex.EncodeToString(h.Sum(nil))
}
