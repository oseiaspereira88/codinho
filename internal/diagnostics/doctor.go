// Package diagnostics implements codinho doctor's checks: binary
// version, catalog validity, local state health (including orphaned-lock
// detection), workspace reachability and required toolchain binaries
// (reliability-observability-compatibility requirement R2, R4). It never
// reads or reports the full process environment, source code, or any
// secret-shaped content (security requirement).
package diagnostics

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/oseiaspereira88/codinho/internal/curriculum"
)

// Status is one check's outcome.
type Status string

const (
	StatusOK    Status = "ok"
	StatusWarn  Status = "warn"
	StatusError Status = "error"
)

// Check is one diagnostic finding: a name, its status, and a detail
// message safe to print (never a path outside workspaceRoot, never an
// environment value, never source code).
type Check struct {
	Name   string `json:"name"`
	Status Status `json:"status"`
	Detail string `json:"detail,omitempty"`
}

// Report is codinho doctor's full output (requirement R10: "relatório
// diagnóstico sanitizado").
type Report struct {
	Version       string  `json:"version"`
	Commit        string  `json:"commit"`
	GoVersion     string  `json:"go_version"`
	OS            string  `json:"os"`
	Arch          string  `json:"arch"`
	WorkspaceRoot string  `json:"workspace_root"`
	Checks        []Check `json:"checks"`
}

// OK reports whether every check in the report passed (no error-status
// checks) — warnings alone do not fail the report.
func (r Report) OK() bool {
	for _, c := range r.Checks {
		if c.Status == StatusError {
			return false
		}
	}
	return true
}

// Run executes every check against workspaceRoot and returns the full
// report.
func Run(workspaceRoot, version, commit string) Report {
	return Report{
		Version:       version,
		Commit:        commit,
		GoVersion:     runtime.Version(),
		OS:            runtime.GOOS,
		Arch:          runtime.GOARCH,
		WorkspaceRoot: workspaceRoot,
		Checks: []Check{
			checkWorkspaceRoot(workspaceRoot),
			checkCatalog(workspaceRoot),
			checkStateWritable(workspaceRoot),
			checkLock(workspaceRoot),
			checkEvidence(workspaceRoot),
			checkToolchain("go"),
			checkToolchain("gofmt"),
		},
	}
}

func checkWorkspaceRoot(root string) Check {
	info, err := os.Stat(root)
	if err != nil {
		return Check{Name: "workspace_root", Status: StatusError, Detail: "workspace root is not accessible"}
	}
	if !info.IsDir() {
		return Check{Name: "workspace_root", Status: StatusError, Detail: "workspace root is not a directory"}
	}
	return Check{Name: "workspace_root", Status: StatusOK}
}

func checkCatalog(root string) Check {
	// LoadPacks returns a nil []Pack both when zero packs are legitimately
	// declared and when a blocking diagnostic prevented loading — never
	// used as a success/failure sentinel here; the diagnostics list is
	// the only authoritative signal.
	_, diags, err := curriculum.LoadPacks(filepath.Join(root, "packs"), curriculum.DefaultLimits)
	if err != nil {
		return Check{Name: "catalog", Status: StatusError, Detail: "packs/manifest.yaml could not be read"}
	}
	blocking := 0
	for _, d := range diags {
		if d.Blocking {
			blocking++
		}
	}
	if blocking > 0 {
		return Check{Name: "catalog", Status: StatusError, Detail: "catalog has blocking diagnostics; run codinho catalog validate"}
	}
	if len(diags) > 0 {
		return Check{Name: "catalog", Status: StatusWarn, Detail: "catalog loads with non-blocking diagnostics"}
	}
	return Check{Name: "catalog", Status: StatusOK}
}

func checkStateWritable(root string) Check {
	stateDir := filepath.Join(root, ".codinho", "state")
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		return Check{Name: "state_writable", Status: StatusError, Detail: "cannot create .codinho/state"}
	}
	probe := filepath.Join(stateDir, ".doctor-probe")
	if err := os.WriteFile(probe, []byte("ok"), 0o600); err != nil {
		return Check{Name: "state_writable", Status: StatusError, Detail: "state directory is not writable"}
	}
	os.Remove(probe)
	return Check{Name: "state_writable", Status: StatusOK}
}

// checkLock detects a lock file whose recorded PID no longer belongs to
// a live process — an orphaned lock left behind by a killed `codinho
// serve` (requirement R2). It never removes the lock itself: recovery is
// the operator's decision.
func checkLock(root string) Check {
	lockPath := filepath.Join(root, ".codinho", "state", "lock")
	data, err := os.ReadFile(lockPath)
	if err != nil {
		if os.IsNotExist(err) {
			return Check{Name: "lock", Status: StatusOK, Detail: "no lock held"}
		}
		return Check{Name: "lock", Status: StatusWarn, Detail: "lock file could not be read"}
	}
	pid, parseErr := strconv.Atoi(strings.TrimSpace(string(data)))
	if parseErr != nil {
		return Check{Name: "lock", Status: StatusWarn, Detail: "lock file content is not a valid PID"}
	}
	if processAlive(pid) {
		return Check{Name: "lock", Status: StatusOK, Detail: "held by a running codinho serve process"}
	}
	return Check{
		Name: "lock", Status: StatusError,
		Detail: "orphaned: the process that held this lock is no longer running; remove .codinho/state/lock manually once you've confirmed no codinho serve is running",
	}
}

func checkEvidence(root string) Check {
	evidenceDir := filepath.Join(root, ".codinho", "state", "evidence")
	entries, err := os.ReadDir(evidenceDir)
	if err != nil {
		if os.IsNotExist(err) {
			return Check{Name: "evidence", Status: StatusOK, Detail: "no evidence recorded yet"}
		}
		return Check{Name: "evidence", Status: StatusWarn, Detail: "evidence directory could not be read"}
	}
	return Check{Name: "evidence", Status: StatusOK, Detail: strconv.Itoa(len(entries)) + " evidence blob(s)"}
}

func checkToolchain(binary string) Check {
	if _, err := exec.LookPath(binary); err != nil {
		return Check{Name: "toolchain:" + binary, Status: StatusError, Detail: binary + " not found on PATH"}
	}
	return Check{Name: "toolchain:" + binary, Status: StatusOK}
}
