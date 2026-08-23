//go:build windows

package checks

import "os/exec"

// ProcessGroupSupported reports whether this platform can guarantee that
// killing a check also kills every process it spawned. Windows support
// here is best-effort only (Process.Kill on the direct child): a full job-
// object implementation is out of scope for this V1 (Compatibility:
// platform-specific guardrails must report real support, never pretend).
func ProcessGroupSupported() bool { return false }

// setProcessGroup is a no-op on this platform: see ProcessGroupSupported.
func setProcessGroup(cmd *exec.Cmd) {}

// killProcessGroup kills only the direct child process on this platform;
// grandchildren it spawned may survive (documented limitation, see
// ProcessGroupSupported).
func killProcessGroup(cmd *exec.Cmd) {
	if cmd.Process == nil {
		return
	}
	_ = cmd.Process.Kill()
}
