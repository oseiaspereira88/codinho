package cli

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/oseiaspereira88/codinho/internal/curriculum"
)

func proofChallenge() curriculum.ChallengeAuthoring {
	good := []curriculum.FixtureFileAuthoring{{Path: "main.go", Content: "package main\n// PRIVATE-REFERENCE\n"}}
	return curriculum.ChallengeAuthoring{ID: "proof", Version: "1.0.0", Fixture: []curriculum.FixtureFileAuthoring{{Path: "main.go", Content: "package main\n"}}, Checks: []curriculum.CheckAuthoring{{ID: "parse", Runner: "internal_ast", Package: "main.go"}}, Validation: &curriculum.CheckValidationAuthoring{ReferenceFixture: good, Expectations: []curriculum.CheckExpectationAuthoring{{CheckID: "parse", Baseline: "pass", Reference: "pass"}}}}
}

func TestEditorialExpectedBaselineAndReference(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func(*curriculum.ChallengeAuthoring)
		valid  bool
	}{
		{"pass", func(*curriculum.ChallengeAuthoring) {}, true},
		{"expected-syntax", func(ch *curriculum.ChallengeAuthoring) {
			ch.Fixture[0].Content = "bad syntax"
			ch.Validation.Expectations[0].Baseline = "syntax_failure"
		}, true},
		{"unexpected-failure", func(ch *curriculum.ChallengeAuthoring) { ch.Fixture[0].Content = "bad syntax" }, false},
		{"bad-reference", func(ch *curriculum.ChallengeAuthoring) { ch.Validation.ReferenceFixture[0].Content = "bad syntax" }, false},
		{"missing-fixture", func(ch *curriculum.ChallengeAuthoring) { ch.Fixture = nil }, false},
		{"missing-proof", func(ch *curriculum.ChallengeAuthoring) { ch.Validation = nil }, false},
		{"unjustified-alternative", func(ch *curriculum.ChallengeAuthoring) { ch.Validation.BaselineFixture = ch.Fixture; ch.Fixture = nil }, false},
		{"reproducible-alternative", func(ch *curriculum.ChallengeAuthoring) {
			ch.Validation.BaselineFixture = ch.Fixture
			ch.Validation.Justification = "Challenge starts from an empty student workspace"
			ch.Fixture = nil
		}, true},
		{"unknown-runner", func(ch *curriculum.ChallengeAuthoring) { ch.Checks[0].Runner = "shell" }, false},
		{"missing-expectation", func(ch *curriculum.ChallengeAuthoring) { ch.Validation.Expectations = nil }, false},
		{"duplicate-expectation", func(ch *curriculum.ChallengeAuthoring) {
			ch.Validation.Expectations = append(ch.Validation.Expectations, ch.Validation.Expectations[0])
		}, false},
		{"network", func(ch *curriculum.ChallengeAuthoring) { ch.Checks[0].Network = true }, false},
		{"escape", func(ch *curriculum.ChallengeAuthoring) { ch.Validation.ReferenceFixture[0].Path = "../secret.go" }, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ch := proofChallenge()
			tc.mutate(&ch)
			findings, report, err := runEditorialProofs([]curriculum.Pack{{ID: "fixture", Version: "1.0.0", Challenges: []curriculum.ChallengeAuthoring{ch}}})
			if err != nil {
				t.Fatal(err)
			}
			if (len(findings) == 0) != tc.valid {
				t.Fatalf("valid=%v findings=%+v", tc.valid, findings)
			}
			if report.DeclaredChecks != 1 || (report.VerifiedChecks == 1) != tc.valid {
				t.Fatalf("false completion: %+v", report)
			}
			data, _ := json.Marshal(report)
			if strings.Contains(string(data), "PRIVATE-REFERENCE") {
				t.Fatal("reference leaked")
			}
			for _, proof := range report.Results {
				if len(proof.Digest) != 64 || proof.ChallengeVersion != "1.0.0" {
					t.Fatalf("unattributed proof: %+v", proof)
				}
			}
		})
	}
}

func TestEditorialNoChecksIsExplicitZero(t *testing.T) {
	findings, report, err := runEditorialProofs([]curriculum.Pack{{Challenges: []curriculum.ChallengeAuthoring{{ID: "no-checks"}}}})
	if err != nil || len(findings) != 0 || report.DeclaredChecks != 0 || report.VerifiedChecks != 0 {
		t.Fatalf("%+v %+v %v", findings, report, err)
	}
}
