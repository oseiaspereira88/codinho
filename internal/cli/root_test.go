package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunNoArgsPrintsUsageToStderr(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), nil, &stdout, &stderr)

	if code != exitUsage {
		t.Fatalf("exit code = %d, want %d", code, exitUsage)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout should stay empty on usage error, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "Usage:") {
		t.Fatalf("stderr missing usage text: %q", stderr.String())
	}
}

func TestRunUnknownCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), []string{"bogus"}, &stdout, &stderr)

	if code != exitUsage {
		t.Fatalf("exit code = %d, want %d", code, exitUsage)
	}
	if !strings.Contains(stderr.String(), `unknown command "bogus"`) {
		t.Fatalf("stderr missing error text: %q", stderr.String())
	}
}

func TestRunVersionIsDeterministic(t *testing.T) {
	var out1, out2, stderr bytes.Buffer
	Run(context.Background(), []string{"version"}, &out1, &stderr)
	Run(context.Background(), []string{"version"}, &out2, &stderr)

	if out1.String() != out2.String() {
		t.Fatalf("version output not deterministic: %q vs %q", out1.String(), out2.String())
	}
	if out1.Len() == 0 {
		t.Fatal("version produced no output")
	}
}

func TestRunDoctorReportsOK(t *testing.T) {
	t.Chdir(setupWorkspace(t))

	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), []string{"doctor"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("exit code = %d, want %d, stdout = %q, stderr = %q", code, exitOK, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "status: ok") {
		t.Fatalf("doctor output missing status: %q", stdout.String())
	}
}

func TestRunDoctorJSONReportsChecksAndNeverLeaksEnv(t *testing.T) {
	t.Setenv("CODINHO_TEST_SECRET_MARKER", "should-never-appear-in-doctor-output")
	t.Chdir(setupWorkspace(t))

	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), []string{"doctor", "--json"}, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("exit code = %d, want %d, stderr = %q", code, exitOK, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"catalog"`) {
		t.Fatalf("doctor --json missing catalog check: %q", stdout.String())
	}
	if strings.Contains(stdout.String(), "should-never-appear-in-doctor-output") {
		t.Fatal("doctor --json leaked an environment value")
	}
}

func TestRunDoctorDetectsOrphanedLock(t *testing.T) {
	root := setupWorkspace(t)
	t.Chdir(root)
	stateDir := filepath.Join(root, ".codinho", "state")
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(stateDir, "lock"), []byte("999999999\n"), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), []string{"doctor"}, &stdout, &stderr)
	if code != exitError {
		t.Fatalf("exit code = %d, want exitError for an orphaned lock, stdout = %q", code, stdout.String())
	}
	if !strings.Contains(stdout.String(), "orphaned") {
		t.Fatalf("doctor output missing orphaned-lock detail: %q", stdout.String())
	}
}

func TestRunHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), []string{"help"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("exit code = %d, want %d", code, exitOK)
	}
	if !strings.Contains(stdout.String(), "Usage:") {
		t.Fatalf("help output missing usage: %q", stdout.String())
	}
}

func TestRunRespectsCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var stdout, stderr bytes.Buffer
	code := Run(ctx, []string{"doctor"}, &stdout, &stderr)

	if code != exitError {
		t.Fatalf("exit code = %d, want %d", code, exitError)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout should stay empty when cancelled, got %q", stdout.String())
	}
}
