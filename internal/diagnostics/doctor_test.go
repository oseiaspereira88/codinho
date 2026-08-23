package diagnostics

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func setupValidWorkspace(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	packsDir := filepath.Join(root, "packs")
	if err := os.MkdirAll(packsDir, 0o700); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(packsDir, "manifest.yaml"), []byte("schema_version: 1\npacks: []\n"), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return root
}

func findCheck(report Report, name string) (Check, bool) {
	for _, c := range report.Checks {
		if c.Name == name {
			return c, true
		}
	}
	return Check{}, false
}

func TestRunReportsOKForAHealthyWorkspace(t *testing.T) {
	root := setupValidWorkspace(t)
	report := Run(root, "v1.2.3", "abc123")

	if !report.OK() {
		t.Fatalf("expected a healthy report, got %+v", report.Checks)
	}
	if report.Version != "v1.2.3" || report.Commit != "abc123" {
		t.Fatalf("report = %+v", report)
	}
	if c, ok := findCheck(report, "catalog"); !ok || c.Status != StatusOK {
		t.Fatalf("catalog check = %+v", c)
	}
	if c, ok := findCheck(report, "lock"); !ok || c.Status != StatusOK {
		t.Fatalf("lock check = %+v", c)
	}
}

func TestRunFlagsMissingWorkspaceRoot(t *testing.T) {
	report := Run(filepath.Join(t.TempDir(), "does-not-exist"), "v1", "c1")
	if report.OK() {
		t.Fatal("expected a non-OK report for a missing workspace root")
	}
	c, ok := findCheck(report, "workspace_root")
	if !ok || c.Status != StatusError {
		t.Fatalf("workspace_root check = %+v", c)
	}
}

func TestRunFlagsBlockingCatalogDiagnostics(t *testing.T) {
	root := t.TempDir()
	packsDir := filepath.Join(root, "packs")
	if err := os.MkdirAll(packsDir, 0o700); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// A challenge with an invalid version is a blocking diagnostic.
	pack := `schema_version: 1
id: bad
version: "not-a-version"
challenges: []
`
	if err := os.WriteFile(filepath.Join(packsDir, "manifest.yaml"), []byte("schema_version: 1\npacks:\n  - pack.yaml\n"), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(packsDir, "pack.yaml"), []byte(pack), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	report := Run(root, "v1", "c1")
	c, ok := findCheck(report, "catalog")
	if !ok || c.Status != StatusError {
		t.Fatalf("catalog check = %+v, want StatusError for an invalid pack version", c)
	}
}

func TestCheckLockDetectsOrphanedLock(t *testing.T) {
	root := setupValidWorkspace(t)
	stateDir := filepath.Join(root, ".codinho", "state")
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// A PID astronomically unlikely to be a real running process.
	if err := os.WriteFile(filepath.Join(stateDir, "lock"), []byte("999999999\n"), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	report := Run(root, "v1", "c1")
	c, ok := findCheck(report, "lock")
	if !ok || c.Status != StatusError {
		t.Fatalf("lock check = %+v, want StatusError for an orphaned lock", c)
	}
	if report.OK() {
		t.Fatal("an orphaned lock must make the overall report non-OK")
	}
}

func TestCheckLockAcceptsALiveProcess(t *testing.T) {
	root := setupValidWorkspace(t)
	stateDir := filepath.Join(root, ".codinho", "state")
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// This test process's own PID is definitely alive.
	if err := os.WriteFile(filepath.Join(stateDir, "lock"), []byte(strconv.Itoa(os.Getpid())+"\n"), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	report := Run(root, "v1", "c1")
	c, ok := findCheck(report, "lock")
	if !ok || c.Status != StatusOK {
		t.Fatalf("lock check = %+v, want StatusOK for a live process", c)
	}
}

func TestCheckStateWritableDetectsReadOnlyState(t *testing.T) {
	root := setupValidWorkspace(t)
	stateDir := filepath.Join(root, ".codinho", "state")
	if err := os.MkdirAll(filepath.Dir(stateDir), 0o700); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := os.Mkdir(stateDir, 0o500); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Cleanup(func() { os.Chmod(stateDir, 0o700) })

	if os.Getuid() == 0 {
		t.Skip("running as root bypasses permission bits")
	}

	report := Run(root, "v1", "c1")
	c, ok := findCheck(report, "state_writable")
	if !ok || c.Status != StatusError {
		t.Fatalf("state_writable check = %+v, want StatusError for a read-only state dir", c)
	}
}

func TestReportNeverExposesFullEnvironment(t *testing.T) {
	t.Setenv("CODINHO_TEST_SECRET_MARKER", "should-never-appear")
	root := setupValidWorkspace(t)
	report := Run(root, "v1", "c1")

	for _, c := range report.Checks {
		if c.Detail == "" {
			continue
		}
		if c.Detail == "should-never-appear" {
			t.Fatalf("check %s leaked an environment value", c.Name)
		}
	}
}
