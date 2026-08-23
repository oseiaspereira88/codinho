package fixtures_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/fixtures"
)

func challengeWithFixture(files ...curriculum.FixtureFileAuthoring) curriculum.ChallengeAuthoring {
	return curriculum.ChallengeAuthoring{ID: "debug.off-by-one", Fixture: files}
}

func TestPlanRejectsEmptyFixture(t *testing.T) {
	dest := t.TempDir()
	_, err := fixtures.Plan(curriculum.ChallengeAuthoring{ID: "no-fixture"}, dest, fixtures.Options{})
	if err != fixtures.ErrEmptyFixture {
		t.Fatalf("err = %v, want ErrEmptyFixture", err)
	}
}

func TestPlanDoesNotWrite(t *testing.T) {
	dest := t.TempDir()
	ch := challengeWithFixture(curriculum.FixtureFileAuthoring{Path: "main.go", Content: "package main\n"})

	plan, err := fixtures.Plan(ch, dest, fixtures.Options{})
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	if len(plan.Files) != 1 || plan.Files[0].Path != "main.go" {
		t.Fatalf("plan.Files = %+v", plan.Files)
	}
	entries, err := os.ReadDir(dest)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("Plan wrote to dest, entries = %v", entries)
	}
}

func TestMaterializeWritesFilesAndManifest(t *testing.T) {
	dest := t.TempDir()
	ch := challengeWithFixture(
		curriculum.FixtureFileAuthoring{Path: "main.go", Content: "package main\n"},
		curriculum.FixtureFileAuthoring{Path: "internal/util.go", Content: "package internal\n"},
	)

	manifest, err := fixtures.Materialize(ch, dest, fixtures.Options{})
	if err != nil {
		t.Fatalf("Materialize: %v", err)
	}
	if manifest.ChallengeID != ch.ID || len(manifest.Files) != 2 {
		t.Fatalf("manifest = %+v", manifest)
	}

	got, err := os.ReadFile(filepath.Join(dest, "main.go"))
	if err != nil || string(got) != "package main\n" {
		t.Fatalf("main.go = %q, err = %v", got, err)
	}
	got, err = os.ReadFile(filepath.Join(dest, "internal", "util.go"))
	if err != nil || string(got) != "package internal\n" {
		t.Fatalf("internal/util.go = %q, err = %v", got, err)
	}

	raw, err := os.ReadFile(filepath.Join(dest, "manifest.json"))
	if err != nil {
		t.Fatalf("reading manifest.json: %v", err)
	}
	var onDisk fixtures.Manifest
	if err := json.Unmarshal(raw, &onDisk); err != nil {
		t.Fatalf("decoding manifest.json: %v", err)
	}
	if onDisk.ChallengeID != ch.ID || len(onDisk.Files) != 2 {
		t.Fatalf("manifest.json = %+v", onDisk)
	}

	// No leftover staging directory.
	entries, err := os.ReadDir(dest)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	for _, e := range entries {
		if filepath.Ext(e.Name()) == "" && e.IsDir() && e.Name() != "internal" {
			t.Fatalf("leftover staging entry: %s", e.Name())
		}
	}
}

func TestMaterializeRefusesOverwriteByDefault(t *testing.T) {
	dest := t.TempDir()
	if err := os.WriteFile(filepath.Join(dest, "main.go"), []byte("existing"), 0o600); err != nil {
		t.Fatalf("seeding existing file: %v", err)
	}
	ch := challengeWithFixture(curriculum.FixtureFileAuthoring{Path: "main.go", Content: "package main\n"})

	_, err := fixtures.Materialize(ch, dest, fixtures.Options{})
	if err != fixtures.ErrWouldOverwrite {
		t.Fatalf("err = %v, want ErrWouldOverwrite", err)
	}
	got, err := os.ReadFile(filepath.Join(dest, "main.go"))
	if err != nil || string(got) != "existing" {
		t.Fatalf("existing file was modified: %q, err = %v", got, err)
	}
}

func TestMaterializeOverwriteTrueReplacesContent(t *testing.T) {
	dest := t.TempDir()
	if err := os.WriteFile(filepath.Join(dest, "main.go"), []byte("existing"), 0o600); err != nil {
		t.Fatalf("seeding existing file: %v", err)
	}
	ch := challengeWithFixture(curriculum.FixtureFileAuthoring{Path: "main.go", Content: "package main\n"})

	if _, err := fixtures.Materialize(ch, dest, fixtures.Options{Overwrite: true}); err != nil {
		t.Fatalf("Materialize: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(dest, "main.go"))
	if err != nil || string(got) != "package main\n" {
		t.Fatalf("main.go = %q, err = %v", got, err)
	}
}

func TestPlanRejectsPathTraversal(t *testing.T) {
	dest := t.TempDir()
	ch := challengeWithFixture(curriculum.FixtureFileAuthoring{Path: "../evil.go", Content: "x"})

	_, err := fixtures.Plan(ch, dest, fixtures.Options{})
	if err == nil {
		t.Fatal("expected error for path traversal")
	}

	parent := filepath.Dir(dest)
	if _, statErr := os.Stat(filepath.Join(parent, "evil.go")); statErr == nil {
		t.Fatal("traversal path was written outside dest")
	}
}

func TestPlanRejectsAbsolutePath(t *testing.T) {
	dest := t.TempDir()
	ch := challengeWithFixture(curriculum.FixtureFileAuthoring{Path: "/etc/passwd", Content: "x"})

	_, err := fixtures.Plan(ch, dest, fixtures.Options{})
	if err == nil {
		t.Fatal("expected error for absolute path")
	}
}

func TestMaterializeFailureLeavesNoPartialConflictState(t *testing.T) {
	dest := t.TempDir()
	// Second file's target directory collides with a pre-existing file,
	// so MkdirAll for it must fail after the first file already staged.
	if err := os.WriteFile(filepath.Join(dest, "pkg"), []byte("not a dir"), 0o600); err != nil {
		t.Fatalf("seeding conflicting file: %v", err)
	}
	ch := challengeWithFixture(
		curriculum.FixtureFileAuthoring{Path: "main.go", Content: "package main\n"},
		curriculum.FixtureFileAuthoring{Path: "pkg/util.go", Content: "package pkg\n"},
	)

	_, err := fixtures.Materialize(ch, dest, fixtures.Options{})
	if err == nil {
		t.Fatal("expected error when a target directory cannot be created")
	}
	if _, statErr := os.Stat(filepath.Join(dest, "main.go")); statErr == nil {
		t.Fatal("main.go was written despite overall failure")
	}
}
