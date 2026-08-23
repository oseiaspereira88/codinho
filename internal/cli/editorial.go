package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/oseiaspereira88/codinho/internal/checks"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/fixtures"
	"github.com/oseiaspereira88/codinho/internal/workspace"
)

// runChecksAgainstFixtures materializes every challenge that declares
// both a fixture and checks into a real temporary directory and runs
// each check for real against it (catalog-authoring-quality requirement
// R5: "executar... checks contra fixtures reproduzíveis"). It never
// assumes the fixture passes its checks — an intentionally buggy debug
// fixture is expected to fail — so only an infrastructure-level error
// (the fixture cannot even be materialized, a check does not resolve, or
// running it errors outright) is reported; Outcome pass/fail/skipped are
// all healthy, reproducible outcomes.
func runChecksAgainstFixtures(packs []curriculum.Pack) ([]curriculum.EditorialFinding, error) {
	var findings []curriculum.EditorialFinding
	executor := checks.NewExecutor()

	for _, p := range packs {
		for _, ch := range p.Challenges {
			if len(ch.Fixture) == 0 || len(ch.Checks) == 0 {
				continue
			}

			dest, err := os.MkdirTemp("", "codinho-catalog-quality-*")
			if err != nil {
				return nil, err
			}
			materializeErr := func() error {
				defer os.RemoveAll(dest)
				if _, err := fixtures.Materialize(ch, dest, fixtures.Options{}); err != nil {
					findings = append(findings, curriculum.EditorialFinding{
						File: p.File, Item: ch.ID, Rule: curriculum.RuleFixtureNotReproducible, Severity: curriculum.SeverityBlocking,
						Detail: err.Error(), Suggestion: "corrija os paths/conteúdo declarados em fixture",
					})
					return nil
				}
				root, err := workspace.AuthorizeRoot(dest)
				if err != nil {
					return err
				}
				for _, c := range ch.Checks {
					resolved, err := checks.Resolve(checks.CheckSpec{
						ID: c.ID, Runner: c.Runner, Package: c.Package, TestPattern: c.TestPattern, Timeout: c.Timeout, NetworkApproved: c.Network,
					})
					if err != nil {
						findings = append(findings, curriculum.EditorialFinding{
							File: p.File, Item: ch.ID + "/" + c.ID, Rule: curriculum.RuleCheckNotResolvable, Severity: curriculum.SeverityBlocking,
							Detail: err.Error(), Suggestion: "corrija runner/package/test_pattern/timeout do check",
						})
						continue
					}
					result, err := executor.Execute(context.Background(), resolved, root)
					if err != nil {
						findings = append(findings, curriculum.EditorialFinding{
							File: p.File, Item: ch.ID + "/" + c.ID, Rule: curriculum.RuleCheckExecutionError, Severity: curriculum.SeverityBlocking,
							Detail: err.Error(),
						})
						continue
					}
					if result.Outcome == checks.OutcomeError {
						findings = append(findings, curriculum.EditorialFinding{
							File: p.File, Item: ch.ID + "/" + c.ID, Rule: curriculum.RuleCheckExecutionError, Severity: curriculum.SeverityBlocking,
							Detail:     fmt.Sprintf("outcome=error exit_code=%d", result.ExitCode),
							Suggestion: "verifique se a fixture forma um workspace válido para este runner (ex.: go.mod para go_test)",
						})
					}
				}
				return nil
			}()
			if materializeErr != nil {
				return nil, materializeErr
			}
		}
	}
	return findings, nil
}
