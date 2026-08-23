//go:build !windows

package diagnostics

import (
	"os"
	"syscall"
)

// processAlive probes pid with signal 0, which the kernel validates
// without actually delivering anything (standard liveness-check idiom on
// Unix; requirement R2, orphaned-lock detection).
func processAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return proc.Signal(syscall.Signal(0)) == nil
}
