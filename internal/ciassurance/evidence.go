// Package ciassurance validates the evidence consumed by native delivery CI.
package ciassurance

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

type Check struct {
	ID       string   `json:"id"`
	Module   string   `json:"module"`
	Stack    string   `json:"stack"`
	Name     string   `json:"name"`
	Program  string   `json:"program"`
	Args     []string `json:"args"`
	Severity string   `json:"severity"`
	Outcome  string   `json:"outcome"`
}

type Report struct {
	SchemaVersion int       `json:"schema_version"`
	GeneratedAt   time.Time `json:"generated_at"`
	Mode          string    `json:"mode"`
	GitHead       string    `json:"git_head"`
	MatrixSHA256  string    `json:"matrix_sha256"`
	Outcome       string    `json:"outcome"`
	Checks        []Check   `json:"checks"`
}

type matrixGroup struct {
	Checks []Check `json:"checks"`
}
type matrix struct {
	Stacks    map[string]matrixGroup `json:"stacks"`
	Overrides map[string]matrixGroup `json:"moduleOverrides"`
}

// Validate requires the exact root Go check set declared by the current matrix.
// It does not infer success from counters or from the presence of a report.
func Validate(raw, matrixRaw []byte, head string, now time.Time) (Report, error) {
	var report Report
	if err := json.Unmarshal(raw, &report); err != nil {
		return report, fmt.Errorf("read validation report: %w", err)
	}
	if report.SchemaVersion != 1 || report.Mode != "strict" || report.Outcome != "pass" {
		return report, fmt.Errorf("strict passing schema-1 report required")
	}
	decodedHead, headErr := hex.DecodeString(head)
	if headErr != nil || len(decodedHead) != 20 || report.GitHead != head {
		return report, fmt.Errorf("validation commit does not match candidate")
	}
	if report.GeneratedAt.IsZero() || report.GeneratedAt.After(now.Add(time.Minute)) || now.Sub(report.GeneratedAt) > 24*time.Hour {
		return report, fmt.Errorf("validation report is stale or future-dated")
	}
	digest := sha256.Sum256(matrixRaw)
	if report.MatrixSHA256 != "sha256:"+hex.EncodeToString(digest[:]) {
		return report, fmt.Errorf("validation matrix digest does not match")
	}
	var m matrix
	if err := json.Unmarshal(matrixRaw, &m); err != nil {
		return report, fmt.Errorf("read validation matrix: %w", err)
	}
	expected := map[string]Check{}
	for _, c := range append(m.Stacks["go"].Checks, m.Overrides["."].Checks...) {
		if c.Name == "" || c.Program == "" || c.Severity != "required" {
			return report, fmt.Errorf("root CI checks must be named and required")
		}
		id := "./go/" + c.Name
		if _, exists := expected[id]; exists {
			return report, fmt.Errorf("duplicate matrix check %q", id)
		}
		expected[id] = c
	}
	if len(expected) == 0 {
		return report, fmt.Errorf("no required checks declared")
	}
	seen := map[string]bool{}
	for _, c := range report.Checks {
		e, ok := expected[c.ID]
		if !ok || seen[c.ID] {
			return report, fmt.Errorf("unknown or duplicate check %q", c.ID)
		}
		seen[c.ID] = true
		if c.Name != e.Name || c.Module != "." || c.Stack != "go" || c.Program != e.Program || !reflect.DeepEqual(c.Args, e.Args) || c.Severity != "required" || c.Outcome != "pass" {
			return report, fmt.Errorf("check %q did not execute the required passing command", c.ID)
		}
	}
	if len(seen) != len(expected) {
		return report, fmt.Errorf("required validation checks are missing")
	}
	return report, nil
}
