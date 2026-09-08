package cli

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/oseiaspereira88/codinho/internal/checks"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/fixtures"
	"github.com/oseiaspereira88/codinho/internal/workspace"
)

type editorialProof struct {
	PackID           string `json:"pack_id"`
	PackVersion      string `json:"pack_version"`
	ChallengeID      string `json:"challenge_id"`
	ChallengeVersion string `json:"challenge_version"`
	CheckID          string `json:"check_id"`
	Scenario         string `json:"scenario"`
	Digest           string `json:"digest"`
	Expected         string `json:"expected"`
	Actual           string `json:"actual"`
	Verified         bool   `json:"verified"`
}

type editorialProofReport struct {
	CodeRevision   string           `json:"code_revision"`
	CodeModified   bool             `json:"code_modified"`
	GoVersion      string           `json:"go_version"`
	DeclaredChecks int              `json:"declared_checks"`
	VerifiedChecks int              `json:"verified_checks"`
	Results        []editorialProof `json:"results"`
}

func runEditorialProofs(packs []curriculum.Pack) ([]curriculum.EditorialFinding, editorialProofReport, error) {
	var findings []curriculum.EditorialFinding
	report := editorialProofReport{Results: []editorialProof{}, CodeRevision: "unknown"}
	if info, ok := debug.ReadBuildInfo(); ok {
		report.GoVersion = info.GoVersion
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				report.CodeRevision = setting.Value
			}
			if setting.Key == "vcs.modified" {
				report.CodeModified = setting.Value == "true"
			}
		}
	}
	executor := checks.NewExecutor()
	for _, p := range packs {
		for _, ch := range p.Challenges {
			report.DeclaredChecks += len(ch.Checks)
			add := func(rule curriculum.EditorialRuleID, detail string) {
				findings = append(findings, curriculum.EditorialFinding{File: p.File, Item: ch.ID, Rule: rule, Severity: curriculum.SeverityBlocking, Detail: detail})
			}
			problems := curriculum.CheckValidationProblems(ch)
			for _, problem := range problems {
				add(curriculum.RuleFixtureNotReproducible, problem)
			}
			// Resolve every check even if fixture evidence is missing: diagnose both.
			for _, check := range ch.Checks {
				resolved, err := checks.Resolve(checks.CheckSpec{ID: check.ID, Runner: check.Runner, Package: check.Package, TestPattern: check.TestPattern, Timeout: check.Timeout})
				if err != nil {
					add(curriculum.RuleCheckNotResolvable, "cannot resolve check "+check.ID)
					continue
				}
				if len(problems) > 0 {
					continue
				}
				v := ch.Validation
				var expectation curriculum.CheckExpectationAuthoring
				for _, e := range v.Expectations {
					if e.CheckID == check.ID {
						expectation = e
					}
				}
				baseline := ch.Fixture
				if len(v.BaselineFixture) > 0 {
					baseline = v.BaselineFixture
				}
				verified := true
				for _, scenario := range []struct {
					name, expected string
					files          []curriculum.FixtureFileAuthoring
				}{
					{"baseline", expectation.Baseline, baseline}, {"reference", expectation.Reference, v.ReferenceFixture},
				} {
					proof := editorialProof{PackID: p.ID, PackVersion: p.Version, ChallengeID: ch.ID, ChallengeVersion: ch.Version, CheckID: check.ID, Scenario: scenario.name, Expected: scenario.expected}
					encoded, err := json.Marshal(struct {
						Check curriculum.CheckAuthoring
						Files []curriculum.FixtureFileAuthoring
					}{check, scenario.files})
					if err != nil {
						return nil, report, err
					}
					digest := sha256.Sum256(encoded)
					proof.Digest = hex.EncodeToString(digest[:])
					actual, err := executeEditorialScenario(executor, resolved, scenario.files)
					if err != nil {
						return nil, report, err
					}
					proof.Actual = actual
					proof.Verified = actual == scenario.expected
					if !proof.Verified {
						verified = false
						add(curriculum.RuleCheckExecutionError, fmt.Sprintf("check=%s scenario=%s expected=%s actual=%s", check.ID, scenario.name, scenario.expected, actual))
					}
					report.Results = append(report.Results, proof)
				}
				if verified {
					report.VerifiedChecks++
				}
			}
		}
	}
	return findings, report, nil
}

func executeEditorialScenario(executor *checks.Executor, resolved checks.Resolved, files []curriculum.FixtureFileAuthoring) (string, error) {
	dest, err := os.MkdirTemp("", "codinho-editorial-*")
	if err != nil {
		return "error", err
	}
	defer os.RemoveAll(dest)
	if _, err = fixtures.Materialize(curriculum.ChallengeAuthoring{Fixture: files}, dest, fixtures.Options{}); err != nil {
		return "invalid_fixture", nil
	}
	root, err := workspace.AuthorizeRoot(dest)
	if err != nil {
		return "error", err
	}
	return executor.ExecuteEditorial(context.Background(), resolved, root)
}
