package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/oseiaspereira88/codinho/internal/application"
	"github.com/oseiaspereira88/codinho/internal/config"
	"github.com/oseiaspereira88/codinho/internal/mastery"
)

const testPack = `schema_version: 1
id: test-pack
version: 1.0.0
themes:
  - id: slices
    title: Slices
concepts:
  - id: slice-declaration
    title: Declaração de slice
competencies:
  - id: slice-filter
    title: Filtrar um slice
tracks: []
challenges:
  - schema_version: 1
    id: test.filter-slice
    version: 1.0.0
    title: Filtrar um slice
    kind: atomic
    difficulty: foundational
    estimated_minutes: 10
    themes: [slices]
    competencies:
      primary: [slice-filter]
      secondary: []
    prerequisites: [slice-declaration]
    brief: Produza apenas os valores aceitos.
    constraints: []
    acceptance:
      - A saída contém apenas valores aceitos.
    layers:
      - id: understanding
        macro_steps:
          - id: model.declare
            kind: micro
            action: declare
            target: named_type
            title: Declarar
            instruction:
              objective: Declarar o tipo.
              scope: Somente a declaração.
              constraints: []
            concepts: [slice-declaration]
            evidence:
              strategies: [source_inspection]
            criteria:
              - id: type-exists
                kind: structural
            completion:
              requires_positive_evaluation: true
              requires_user_confirmation: false
    checks: []
    fixture:
      - path: main.go
        content: "package main\n"
      - path: internal/helper.go
        content: "package internal\n"
`

// setupWorkspace creates a temp workspace with a valid pack and returns its
// path. It does not chdir; callers that need config.Load's cwd-based
// resolution to see it must call t.Chdir themselves.
func setupWorkspace(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	packsDir := filepath.Join(root, "packs")
	if err := os.MkdirAll(packsDir, 0o700); err != nil {
		t.Fatalf("mkdir packs: %v", err)
	}
	manifest := "schema_version: 1\npacks:\n  - test-pack.yaml\n"
	if err := os.WriteFile(filepath.Join(packsDir, "manifest.yaml"), []byte(manifest), 0o600); err != nil {
		t.Fatalf("writing manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(packsDir, "test-pack.yaml"), []byte(testPack), 0o600); err != nil {
		t.Fatalf("writing pack: %v", err)
	}
	return root
}

func run(t *testing.T, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	var out, errOut bytes.Buffer
	code = Run(context.Background(), args, &out, &errOut)
	return out.String(), errOut.String(), code
}

func TestCatalogValidateOK(t *testing.T) {
	t.Chdir(setupWorkspace(t))

	stdout, stderr, code := run(t, "catalog", "validate")
	if code != exitOK {
		t.Fatalf("code = %d, stderr = %q", code, stderr)
	}
	if !strings.Contains(stdout, "ok") {
		t.Fatalf("stdout = %q", stdout)
	}
}

func TestCatalogValidateJSONIsStableAcrossRuns(t *testing.T) {
	t.Chdir(setupWorkspace(t))

	out1, _, _ := run(t, "catalog", "validate", "--json")
	out2, _, _ := run(t, "catalog", "validate", "--json")
	if out1 != out2 {
		t.Fatalf("non-deterministic JSON output:\n%s\nvs\n%s", out1, out2)
	}
}

func TestCatalogListDefaultsToChallenges(t *testing.T) {
	t.Chdir(setupWorkspace(t))

	stdout, stderr, code := run(t, "catalog", "list")
	if code != exitOK {
		t.Fatalf("code = %d, stderr = %q", code, stderr)
	}
	if !strings.Contains(stdout, "test.filter-slice") {
		t.Fatalf("stdout missing challenge: %q", stdout)
	}
}

func TestCatalogListJSON(t *testing.T) {
	t.Chdir(setupWorkspace(t))

	stdout, _, code := run(t, "catalog", "list", "--json")
	if code != exitOK {
		t.Fatalf("code = %d", code)
	}
	var result struct {
		Items []struct{ ID string }
	}
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("decoding JSON: %v\n%s", err, stdout)
	}
	if len(result.Items) != 1 || result.Items[0].ID != "test.filter-slice" {
		t.Fatalf("items = %+v", result.Items)
	}
}

func TestCatalogShowChallenge(t *testing.T) {
	t.Chdir(setupWorkspace(t))

	stdout, stderr, code := run(t, "catalog", "show", "test.filter-slice")
	if code != exitOK {
		t.Fatalf("code = %d, stderr = %q", code, stderr)
	}
	if !strings.Contains(stdout, "test.filter-slice") || !strings.Contains(stdout, "fixture:") {
		t.Fatalf("stdout = %q", stdout)
	}
}

func TestCatalogShowNotFound(t *testing.T) {
	t.Chdir(setupWorkspace(t))

	_, stderr, code := run(t, "catalog", "show", "does-not-exist")
	if code != exitError {
		t.Fatalf("code = %d", code)
	}
	if !strings.Contains(stderr, "not found") {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestWorkspacePrepareDryRunDoesNotWrite(t *testing.T) {
	t.Chdir(setupWorkspace(t))
	dest := t.TempDir()

	stdout, stderr, code := run(t, "workspace", "prepare", "test.filter-slice", "--dest", dest, "--dry-run")
	if code != exitOK {
		t.Fatalf("code = %d, stderr = %q", code, stderr)
	}
	if !strings.Contains(stdout, "main.go") {
		t.Fatalf("stdout = %q", stdout)
	}
	entries, _ := os.ReadDir(dest)
	if len(entries) != 0 {
		t.Fatalf("dry-run wrote files: %v", entries)
	}
}

func TestWorkspacePrepareMaterializesAndProtectsAgainstOverwrite(t *testing.T) {
	t.Chdir(setupWorkspace(t))
	dest := t.TempDir()

	_, stderr, code := run(t, "workspace", "prepare", "test.filter-slice", "--dest", dest)
	if code != exitOK {
		t.Fatalf("code = %d, stderr = %q", code, stderr)
	}
	if _, err := os.Stat(filepath.Join(dest, "main.go")); err != nil {
		t.Fatalf("main.go not written: %v", err)
	}

	_, stderr, code = run(t, "workspace", "prepare", "test.filter-slice", "--dest", dest)
	if code != exitError {
		t.Fatalf("second prepare should fail without --overwrite, code = %d", code)
	}
	if !strings.Contains(stderr, "conflict") {
		t.Fatalf("stderr = %q", stderr)
	}

	_, stderr, code = run(t, "workspace", "prepare", "test.filter-slice", "--dest", dest, "--overwrite")
	if code != exitOK {
		t.Fatalf("overwrite should succeed, code = %d, stderr = %q", code, stderr)
	}
}

func TestWorkspacePrepareUnknownChallenge(t *testing.T) {
	t.Chdir(setupWorkspace(t))
	dest := t.TempDir()

	_, stderr, code := run(t, "workspace", "prepare", "does-not-exist", "--dest", dest)
	if code != exitError {
		t.Fatalf("code = %d", code)
	}
	if !strings.Contains(stderr, "not found") {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestSessionInspectNotFound(t *testing.T) {
	t.Chdir(setupWorkspace(t))

	_, _, code := run(t, "session", "inspect", "ses_missing")
	if code != exitError {
		t.Fatalf("code = %d", code)
	}
}

func TestSessionInspectReportsRecordedEvents(t *testing.T) {
	root := setupWorkspace(t)
	t.Chdir(root)

	store, err := openEventStore(config.Config{WorkspaceRoot: root})
	if err != nil {
		t.Fatalf("openEventStore: %v", err)
	}
	if _, err := store.Append("ses_1", 0, "", "session_started", map[string]string{
		"challenge_id": "test.filter-slice", "mode": "practice",
	}); err != nil {
		t.Fatalf("seeding event: %v", err)
	}
	store.Close()

	stdout, stderr, code := run(t, "session", "inspect", "ses_1")
	if code != exitOK {
		t.Fatalf("code = %d, stderr = %q", code, stderr)
	}
	if !strings.Contains(stdout, "test.filter-slice") || !strings.Contains(stdout, "active") {
		t.Fatalf("stdout = %q", stdout)
	}
}

func TestProgressShowAndExportReflectRecordedEvidence(t *testing.T) {
	root := setupWorkspace(t)
	t.Chdir(root)

	store, err := openEventStore(config.Config{WorkspaceRoot: root})
	if err != nil {
		t.Fatalf("openEventStore: %v", err)
	}
	svc := application.NewProgressService(store)
	if _, err := svc.RecordEvidence(application.EvidenceInput{
		CompetencyID: "slice-filter",
		Dimension:    string(mastery.DimensionExplanation),
		EvidenceID:   "ev1",
		Success:      true,
	}); err != nil {
		t.Fatalf("RecordEvidence: %v", err)
	}
	store.Close()

	stdout, stderr, code := run(t, "progress", "show")
	if code != exitOK {
		t.Fatalf("code = %d, stderr = %q", code, stderr)
	}
	if !strings.Contains(stdout, "slice-filter") {
		t.Fatalf("stdout = %q", stdout)
	}

	out := filepath.Join(t.TempDir(), "export.json")
	_, stderr, code = run(t, "progress", "export", "--out", out)
	if code != exitOK {
		t.Fatalf("code = %d, stderr = %q", code, stderr)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("reading export: %v", err)
	}
	if !strings.Contains(string(data), "slice-filter") {
		t.Fatalf("export = %s", data)
	}
}

func TestInitScaffoldsManifestAndRefusesOverwrite(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)

	stdout, stderr, code := run(t, "init")
	if code != exitOK {
		t.Fatalf("code = %d, stderr = %q", code, stderr)
	}
	if !strings.Contains(stdout, "initialized") {
		t.Fatalf("stdout = %q", stdout)
	}
	manifestPath := filepath.Join(root, "packs", "manifest.yaml")
	if _, err := os.Stat(manifestPath); err != nil {
		t.Fatalf("manifest not created: %v", err)
	}

	_, stderr, code = run(t, "init")
	if code != exitError {
		t.Fatalf("second init should fail without --force, code = %d", code)
	}
	if !strings.Contains(stderr, "already exists") {
		t.Fatalf("stderr = %q", stderr)
	}

	_, stderr, code = run(t, "init", "--force")
	if code != exitOK {
		t.Fatalf("init --force should succeed, code = %d, stderr = %q", code, stderr)
	}
}

func TestServeIsRejectedByDispatcher(t *testing.T) {
	_, stderr, code := run(t, "serve")
	if code != exitError {
		t.Fatalf("code = %d", code)
	}
	if !strings.Contains(stderr, "codinho serve") {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestHelpListsAllCommands(t *testing.T) {
	stdout, _, code := run(t, "help")
	if code != exitOK {
		t.Fatalf("code = %d", code)
	}
	for _, cmd := range []string{"serve", "init", "catalog", "session", "progress", "workspace"} {
		if !strings.Contains(stdout, cmd) {
			t.Fatalf("help missing command %q: %q", cmd, stdout)
		}
	}
}
