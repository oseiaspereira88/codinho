// Package fixtures materializes a challenge's authored starter files (e.g.
// buggy code for a debug challenge) onto disk. It never runs from an MCP
// tool: the administrative CLI's workspace-prepare command is the only
// caller, so writing to a learner's project always requires an explicit,
// out-of-band, human-invoked command (administrative-cli-fixtures
// Decision 1).
package fixtures

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/oseiaspereira88/codinho/internal/curriculum"
)

// ErrEmptyFixture is returned when the challenge declares no fixture
// files: there is nothing to prepare.
var ErrEmptyFixture = errors.New("fixtures: challenge declares no fixture files")

// ErrInvalidPath is returned when an authored fixture path is empty,
// absolute, or would escape the destination via "..". The loader's
// validator already rejects this at authoring time; this is the runtime
// re-check that closes the same class of escape (requirement R5).
var ErrInvalidPath = errors.New("fixtures: fixture file path is empty or escapes the destination")

// ErrWouldOverwrite is returned by Plan (and therefore Materialize) when
// one or more target files already exist and Options.Overwrite is false
// (requirement R5: no command overwrites by default).
var ErrWouldOverwrite = errors.New("fixtures: destination already has files this operation would overwrite")

// PlannedFile is one fixture file resolved against a destination, before
// anything is written.
type PlannedFile struct {
	Path   string // relative, forward-slash, cleaned
	SHA256 string
	Size   int64
	Exists bool
}

// FixturePlan is the full effect of a Materialize call, computable without
// writing anything (requirement R7: dry-run for every multi-file write).
type FixturePlan struct {
	ChallengeID string
	Dest        string
	Files       []PlannedFile
	Conflicts   []string
}

// Options controls Plan and Materialize.
type Options struct {
	Overwrite bool
}

// Plan validates every fixture file's path against dest and reports
// conflicts, without writing anything. It returns ErrWouldOverwrite (along
// with the computed Plan, so a caller can still show what was found) when
// conflicts exist and Options.Overwrite is false.
func Plan(challenge curriculum.ChallengeAuthoring, dest string, opts Options) (FixturePlan, error) {
	if len(challenge.Fixture) == 0 {
		return FixturePlan{}, ErrEmptyFixture
	}

	realDest, err := authorizeDest(dest)
	if err != nil {
		return FixturePlan{}, err
	}

	plan := FixturePlan{ChallengeID: challenge.ID, Dest: realDest}
	seen := map[string]bool{}
	for _, f := range challenge.Fixture {
		rel, err := safeRelPath(f.Path)
		if err != nil {
			return FixturePlan{}, fmt.Errorf("fixture file %q: %w", f.Path, err)
		}
		if seen[rel] {
			return FixturePlan{}, fmt.Errorf("fixtures: duplicate fixture path %q", rel)
		}
		seen[rel] = true

		sum := sha256.Sum256([]byte(f.Content))
		target := filepath.Join(realDest, filepath.FromSlash(rel))
		_, statErr := os.Lstat(target)
		exists := statErr == nil

		plan.Files = append(plan.Files, PlannedFile{
			Path: rel, SHA256: hex.EncodeToString(sum[:]), Size: int64(len(f.Content)), Exists: exists,
		})
		if exists {
			plan.Conflicts = append(plan.Conflicts, rel)
		}
	}
	sort.Slice(plan.Files, func(i, j int) bool { return plan.Files[i].Path < plan.Files[j].Path })
	sort.Strings(plan.Conflicts)

	if len(plan.Conflicts) > 0 && !opts.Overwrite {
		return plan, ErrWouldOverwrite
	}
	return plan, nil
}

// Materialize plans, then writes, every fixture file for challenge into
// dest, plus a manifest.json recording each file's hash (requirement R6).
// Files are staged in a temp directory under dest and moved into place
// only after every one of them is written successfully, so a failure
// partway through never leaves a half-written file at its final path
// (Technical risk: "Fixture pode sobrescrever trabalho por erro de
// resolução").
func Materialize(challenge curriculum.ChallengeAuthoring, dest string, opts Options) (Manifest, error) {
	plan, err := Plan(challenge, dest, opts)
	if err != nil {
		return Manifest{}, err
	}

	content := make(map[string]string, len(challenge.Fixture))
	for _, f := range challenge.Fixture {
		rel, err := safeRelPath(f.Path)
		if err != nil {
			return Manifest{}, err
		}
		content[rel] = f.Content
	}

	tmp, err := os.MkdirTemp(plan.Dest, ".fixture-tmp-*")
	if err != nil {
		return Manifest{}, err
	}
	defer os.RemoveAll(tmp)

	for _, pf := range plan.Files {
		staged := filepath.Join(tmp, filepath.FromSlash(pf.Path))
		if err := os.MkdirAll(filepath.Dir(staged), 0o700); err != nil {
			return Manifest{}, err
		}
		if err := os.WriteFile(staged, []byte(content[pf.Path]), 0o600); err != nil {
			return Manifest{}, err
		}
	}

	// Validate every final target directory can actually be created before
	// renaming anything, so a conflict discovered on file N never leaves
	// files 1..N-1 already committed at their final path.
	for _, pf := range plan.Files {
		final := filepath.Join(plan.Dest, filepath.FromSlash(pf.Path))
		if err := checkPathIsCreatable(plan.Dest, filepath.Dir(final)); err != nil {
			return Manifest{}, err
		}
	}

	for _, pf := range plan.Files {
		staged := filepath.Join(tmp, filepath.FromSlash(pf.Path))
		final := filepath.Join(plan.Dest, filepath.FromSlash(pf.Path))
		if err := os.MkdirAll(filepath.Dir(final), 0o700); err != nil {
			return Manifest{}, err
		}
		// Re-check containment after directory creation, closing the same
		// TOCTOU window internal/workspace's Root.Resolve closes.
		realParent, err := filepath.EvalSymlinks(filepath.Dir(final))
		if err != nil {
			return Manifest{}, err
		}
		if !isWithin(plan.Dest, realParent) {
			return Manifest{}, ErrInvalidPath
		}
		if err := os.Rename(staged, final); err != nil {
			return Manifest{}, err
		}
	}

	manifest := Manifest{ChallengeID: challenge.ID, CreatedAt: time.Now().UTC().Format(time.RFC3339)}
	for _, pf := range plan.Files {
		manifest.Files = append(manifest.Files, ManifestFile{Path: pf.Path, SHA256: pf.SHA256, Size: pf.Size})
	}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return Manifest{}, err
	}
	if err := os.WriteFile(filepath.Join(plan.Dest, "manifest.json"), data, 0o600); err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}

// authorizeDest creates dest if absent, then resolves it to a real
// (symlink-evaluated) absolute directory, mirroring
// internal/workspace.AuthorizeRoot.
func authorizeDest(dest string) (string, error) {
	if err := os.MkdirAll(dest, 0o700); err != nil {
		return "", err
	}
	abs, err := filepath.Abs(dest)
	if err != nil {
		return "", err
	}
	real, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(real)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("fixtures: destination %q is not a directory", dest)
	}
	return real, nil
}

// safeRelPath cleans raw and rejects anything empty, absolute, or that
// would escape the destination via "..".
func safeRelPath(raw string) (string, error) {
	if raw == "" {
		return "", ErrInvalidPath
	}
	clean := path.Clean(filepath.ToSlash(raw))
	if path.IsAbs(clean) || clean == "." || clean == ".." || strings.HasPrefix(clean, "../") {
		return "", ErrInvalidPath
	}
	return clean, nil
}

// checkPathIsCreatable walks up from dir toward root and fails if any
// ancestor already exists as a non-directory, without creating anything.
func checkPathIsCreatable(root, dir string) error {
	for len(dir) >= len(root) {
		info, err := os.Lstat(dir)
		if err == nil {
			if !info.IsDir() {
				return fmt.Errorf("fixtures: %s exists and is not a directory", dir)
			}
			return nil
		}
		if !os.IsNotExist(err) {
			return err
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return nil
		}
		dir = parent
	}
	return nil
}

func isWithin(root, candidate string) bool {
	if candidate == root {
		return true
	}
	return strings.HasPrefix(candidate, root+string(filepath.Separator))
}
