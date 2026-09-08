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
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		host = "github-actions"
	}
	return json.NewEncoder(os.Stdout).Encode(struct {
		SchemaVersion int       `json:"schema_version"`
		GitHead       string    `json:"git_head"`
		OS            string    `json:"os"`
		Arch          string    `json:"arch"`
		GoVersion     string    `json:"go_version"`
		Host          string    `json:"host"`
		RunID         string    `json:"run_id,omitempty"`
		RunAttempt    string    `json:"run_attempt,omitempty"`
		GeneratedAt   time.Time `json:"generated_at"`
		Checks        int       `json:"passed_checks"`
		Outcome       string    `json:"outcome"`
	}{1, head, runtime.GOOS, runtime.GOARCH, runtime.Version(), host, os.Getenv("GITHUB_RUN_ID"), os.Getenv("GITHUB_RUN_ATTEMPT"), report.GeneratedAt, len(report.Checks), "pass"})
}
