package main

import (
	"os/exec"
	"path/filepath"
	"testing"
)

// TestSmokeInstallScript runs the real scripts/smoke-install.sh against a
// freshly built binary, proving installation-documentation-ci requirement
// R8 ("validar todos os links, comandos e exemplos em ambiente limpo")
// with a real script execution, not just a description of what it does.
func TestSmokeInstallScript(t *testing.T) {
	bin := buildCodinhoBinary(t)
	script, err := filepath.Abs("../../scripts/smoke-install.sh")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skip("bash not available")
	}

	cmd := exec.Command("bash", script, bin)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("smoke-install.sh failed: %v\n%s", err, out)
	}
}
