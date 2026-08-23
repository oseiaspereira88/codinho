package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
)

// FileHash is one file's content hash relative to a Root, as of one
// observation.
type FileHash struct {
	Path string // root-relative, forward-slash
	Hash string // "sha256-<hex>", matching internal/evidence's ID shape
	Size int64
}

// hashFor returns the content-addressed hash of data, in the same shape
// internal/evidence.Store uses, so a workspace hash and an evidence ID are
// directly comparable.
func hashFor(data []byte) string {
	sum := sha256.Sum256(data)
	return "sha256-" + hex.EncodeToString(sum[:])
}

// walkMatched walks root read-only and returns the hash of every regular
// file whose root-relative path matches globs and is not covered by
// defaultExcludeNames/defaultExcludeSuffixes (requirement R3; security).
// Symlinks are neither followed nor collected: a symlink inside the root
// could point outside it, and this package never needs to read through
// one (requirement R1).
func walkMatched(root Root, globs []string) ([]FileHash, error) {
	var out []FileHash
	err := filepath.WalkDir(root.Path(), func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type()&fs.ModeSymlink != 0 {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		rel, relErr := root.Rel(p)
		if relErr != nil {
			return relErr
		}
		if d.IsDir() {
			if rel != "." && isExcluded(rel) {
				return filepath.SkipDir
			}
			return nil
		}
		if isExcluded(rel) || !matchAny(globs, rel) {
			return nil
		}
		data, readErr := os.ReadFile(p)
		if readErr != nil {
			return readErr
		}
		out = append(out, FileHash{Path: rel, Hash: hashFor(data), Size: int64(len(data))})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}
