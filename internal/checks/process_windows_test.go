//go:build windows

package checks

import "testing"

// waitForProcessDeath has no reliable, dependency-free equivalent on
// windows in this test suite; TestExecuteTimeoutKillsTheWholeProcessGroup
// skips itself via ProcessGroupSupported() before this could be called.
func waitForProcessDeath(t *testing.T, pidText string) {
	t.Helper()
	t.Fatal("waitForProcessDeath should be unreachable on windows: ProcessGroupSupported() is false")
}
