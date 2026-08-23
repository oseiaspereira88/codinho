//go:build unix

package checks

import (
	"os/exec"
	"syscall"
)

// ProcessGroupSupported reports whether this platform can guarantee that
// killing a check also kills every process it spawned (Compatibility:
// platform-specific guardrails must report real support, never pretend).
func ProcessGroupSupported() bool { return true }

// setProcessGroup starts cmd in its own process group, so killProcessGroup
// can terminate the whole tree instead of only the direct child
// (requirement R4, non-functional: a canceled check must not leak a
// process).
func setProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

// killProcessGroup signals the entire process group cmd started, once its
// PID is known. Called after cmd.Start succeeds.
func killProcessGroup(cmd *exec.Cmd) {
	if cmd.Process == nil {
		return
	}
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
}
