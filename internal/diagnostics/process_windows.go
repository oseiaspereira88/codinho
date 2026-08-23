//go:build windows

package diagnostics

import "os"

// processAlive probes pid by attempting to open a handle to it: on
// Windows, os.FindProcess itself fails for a PID that does not exist,
// unlike on Unix where it always succeeds (requirement R2, orphaned-lock
// detection; Compatibility: "degradar explicitamente quando plataforma
// não oferecer uma primitiva" — signal 0 has no Windows equivalent, so
// this is the platform-appropriate substitute, not a degraded stub).
func processAlive(pid int) bool {
	_, err := os.FindProcess(pid)
	return err == nil
}
