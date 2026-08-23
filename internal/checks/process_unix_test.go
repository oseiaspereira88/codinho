//go:build unix

package checks

import (
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

// waitForProcessDeath polls until pid (given as text, possibly with
// trailing whitespace) no longer exists, proving killProcessGroup reached
// a grandchild the check itself spawned (requirement R4, non-functional:
// a canceled check must not leak a process).
func waitForProcessDeath(t *testing.T, pidText string) {
	t.Helper()
	pid, err := strconv.Atoi(strings.TrimSpace(pidText))
	if err != nil {
		t.Fatalf("unexpected error parsing pid: %v", err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if err := syscall.Kill(pid, 0); err != nil {
			return // ESRCH: the process is gone
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("expected grandchild pid %d to be dead after process-group kill, but it is still alive", pid)
}
