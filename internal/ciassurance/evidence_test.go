package ciassurance

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestValidateEvidence(t *testing.T) {
	now := time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC)
	head := strings.Repeat("a", 40)
	matrixRaw := []byte(`{"stacks":{"go":{"checks":[{"name":"test","program":"go","args":["test","./..."],"severity":"required"}]}},"moduleOverrides":{}}`)
	base := Report{SchemaVersion: 1, GeneratedAt: now, Mode: "strict", GitHead: head, MatrixSHA256: fmt.Sprintf("sha256:%x", sha256.Sum256(matrixRaw)), Outcome: "pass", Checks: []Check{{ID: "./go/test", Module: ".", Stack: "go", Name: "test", Program: "go", Args: []string{"test", "./..."}, Severity: "required", Outcome: "pass"}}}
	cases := []struct {
		name   string
		change func(*Report)
		valid  bool
	}{
		{"passing", func(r *Report) {}, true},
		{"missing checks", func(r *Report) { r.Checks = nil }, false},
		{"skipped", func(r *Report) { r.Checks[0].Outcome = "skipped" }, false},
		{"failed", func(r *Report) { r.Checks[0].Outcome = "fail" }, false},
		{"unknown", func(r *Report) { r.Checks[0].ID = "./go/unknown" }, false},
		{"duplicate", func(r *Report) { r.Checks = append(r.Checks, r.Checks[0]) }, false},
		{"stale", func(r *Report) { r.GeneratedAt = now.Add(-25 * time.Hour) }, false},
		{"future", func(r *Report) { r.GeneratedAt = now.Add(2 * time.Minute) }, false},
		{"undated", func(r *Report) { r.GeneratedAt = time.Time{} }, false},
		{"wrong commit", func(r *Report) { r.GitHead = strings.Repeat("b", 40) }, false},
		{"wrong matrix", func(r *Report) { r.MatrixSHA256 = "sha256:bad" }, false},
		{"tolerant", func(r *Report) { r.Mode = "tolerant" }, false},
		{"schema", func(r *Report) { r.SchemaVersion = 2 }, false},
		{"failed report", func(r *Report) { r.Outcome = "fail" }, false},
		{"command", func(r *Report) { r.Checks[0].Program = "true" }, false},
		{"arguments", func(r *Report) { r.Checks[0].Args = []string{"version"} }, false},
		{"module", func(r *Report) { r.Checks[0].Module = "other" }, false},
		{"stack", func(r *Report) { r.Checks[0].Stack = "node" }, false},
		{"optional", func(r *Report) { r.Checks[0].Severity = "optional" }, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			raw, _ := json.Marshal(base)
			var r Report
			if err := json.Unmarshal(raw, &r); err != nil {
				t.Fatal(err)
			}
			tc.change(&r)
			raw, _ = json.Marshal(r)
			_, err := Validate(raw, matrixRaw, head, now)
			if (err == nil) != tc.valid {
				t.Fatalf("valid=%v error=%v", tc.valid, err)
			}
		})
	}
	raw, _ := json.Marshal(base)
	for _, missing := range [][]byte{nil, []byte("{"), []byte("null")} {
		if _, err := Validate(missing, matrixRaw, head, now); err == nil {
			t.Fatal("accepted missing/invalid report")
		}
	}
	if _, err := Validate(raw, matrixRaw, strings.Repeat("z", 40), now); err == nil {
		t.Fatal("accepted invalid commit")
	}
	for _, bad := range []string{`{`, `{}`, `{"stacks":{"go":{"checks":[{"name":"test","program":"go","severity":"optional"}]}}}`, `{"stacks":{"go":{"checks":[{"name":"test","program":"go","severity":"required"},{"name":"test","program":"go","severity":"required"}]}}}`} {
		r := base
		r.MatrixSHA256 = fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(bad)))
		raw, _ = json.Marshal(r)
		if _, err := Validate(raw, []byte(bad), head, now); err == nil {
			t.Fatalf("accepted invalid matrix %s", bad)
		}
	}
}
