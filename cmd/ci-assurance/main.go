// ci-assurance is a repository administration tool, not a learner command.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/oseiaspereira88/codinho/internal/ciassurance"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	raw, err := os.ReadFile(".pose/results/delivery-validation.json")
	if err != nil {
		return err
	}
	matrix, err := os.ReadFile(".pose/indexes/validation-matrix.json")
	if err != nil {
		return err
	}
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return err
	}
	head := strings.TrimSpace(string(out))
	if candidate := os.Getenv("GITHUB_SHA"); candidate != "" && candidate != head {
		return fmt.Errorf("checkout does not match GitHub candidate")
	}
	report, err := ciassurance.Validate(raw, matrix, head, time.Now().UTC())
	if err != nil {
		return err
	}
	host := "local"
	modified, err := candidateModified(".")
	if err != nil {
		return err
	}
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		host = "github-actions"
		if modified {
			return fmt.Errorf("CI candidate source changed during validation")
		}
	}
	return json.NewEncoder(os.Stdout).Encode(struct {
		SchemaVersion  int       `json:"schema_version"`
		GitHead        string    `json:"git_head"`
		OS             string    `json:"os"`
		Arch           string    `json:"arch"`
		GoVersion      string    `json:"go_version"`
		Host           string    `json:"host"`
		RunID          string    `json:"run_id,omitempty"`
		RunAttempt     string    `json:"run_attempt,omitempty"`
		GeneratedAt    time.Time `json:"generated_at"`
		Checks         int       `json:"passed_checks"`
		Outcome        string    `json:"outcome"`
		SourceModified bool      `json:"source_modified"`
	}{1, head, runtime.GOOS, runtime.GOARCH, runtime.Version(), host, os.Getenv("GITHUB_RUN_ID"), os.Getenv("GITHUB_RUN_ATTEMPT"), report.GeneratedAt, len(report.Checks), "pass", modified})
}

// Generated validation/assessment reports do not change the candidate source.
// Include untracked authored inputs so a local run cannot silently attest HEAD.
func candidateModified(root string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain=v1", "--untracked-files=all", "--",
		"cmd", "internal", "packs", "scripts", ".github", "Makefile", "go.mod", "go.sum",
		"README.md", "PROJECT.md", "docs", ".agents/skills/codinho", ".pose/adr", ".pose/specs", ".pose/contracts", ".pose/policy", ".pose/docs.json",
		".pose/indexes/module-metadata.json", ".pose/indexes/validation-matrix.json")
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("inspect candidate source: %w", err)
	}
	return len(out) > 0, nil
}
