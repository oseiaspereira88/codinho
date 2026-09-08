package cli

import (
	"fmt"
	"io"
	"path/filepath"
	"sort"

	"github.com/oseiaspereira88/codinho/internal/config"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
)

const catalogUsage = `Usage: codinho catalog <command>

Commands:
  validate  Load every pack and print diagnostics
  list      List catalog items, optionally filtered
  show      Show one catalog item in full
`

func runCatalog(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, catalogUsage)
		return exitUsage
	}
	switch args[0] {
	case "validate":
		return runCatalogValidate(args[1:], stdout, stderr)
	case "list":
		return runCatalogList(args[1:], stdout, stderr)
	case "show":
		return runCatalogShow(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "codinho: unknown catalog subcommand %q\n\n", args[0])
		fmt.Fprint(stderr, catalogUsage)
		return exitUsage
	}
}

// loadCatalog loads administrative inventory, including drafts. Normal serve
// additionally selects the published projection for every runtime service.
func loadCatalog(cfg config.Config) (*curriculum.Catalog, []curriculum.Diagnostic, error) {
	return curriculum.Load(filepath.Join(cfg.WorkspaceRoot, "packs"), curriculum.DefaultLimits)
}

type diagnosticOut struct {
	File     string `json:"file"`
	Item     string `json:"item,omitempty"`
	Field    string `json:"field,omitempty"`
	Code     string `json:"code"`
	Detail   string `json:"detail,omitempty"`
	Blocking bool   `json:"blocking"`
}

type editorialFindingOut struct {
	File       string `json:"file,omitempty"`
	Item       string `json:"item,omitempty"`
	Rule       string `json:"rule"`
	Severity   string `json:"severity"`
	Detail     string `json:"detail,omitempty"`
	Suggestion string `json:"suggestion,omitempty"`
}

type catalogValidateOut struct {
	Diagnostics []diagnosticOut                `json:"diagnostics"`
	Editorial   []editorialFindingOut          `json:"editorial"`
	Coverage    curriculum.Coverage            `json:"coverage"`
	V1Gate      []editorialFindingOut          `json:"v1_gate,omitempty"`
	Publication curriculum.PublicationCoverage `json:"publication"`
	Proof       *editorialProofReport          `json:"proof,omitempty"`
}

func runCatalogValidate(args []string, stdout, stderr io.Writer) int {
	fs := newFlagSet("catalog validate", stderr)
	jsonOut := fs.Bool("json", false, "print diagnostics as JSON")
	v1Gate := fs.Bool("v1-gate", false, "also check coverage against the V1 roadmap thresholds (requirement R9); off by default since an in-progress catalog is expected to be below them")
	runChecks := fs.Bool("checks", false, "execute declared baseline and reference expectations for every authored check")
	publishedChecks := fs.Bool("published-checks", false, "execute every published check against baseline and reference")
	distribution := fs.String("distribution", "", "validate eligible counts against a JSON policy; V1 defaults to packs/distribution.json")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}

	cfg := config.Load()
	packs, diags, err := curriculum.LoadPacks(filepath.Join(cfg.WorkspaceRoot, "packs"), curriculum.DefaultLimits)
	if err != nil {
		fmt.Fprintf(stderr, "codinho: catalog validate: %v\n", err)
		return exitError
	}

	blocking := false
	for _, d := range diags {
		if d.Blocking {
			blocking = true
		}
	}

	var editorial []curriculum.EditorialFinding
	var coverage curriculum.Coverage
	var publication curriculum.PublicationCoverage
	var proof *editorialProofReport
	var v1Findings []curriculum.EditorialFinding
	if !blocking {
		editorial = curriculum.RunEditorialChecks(packs)
		publication = curriculum.ProjectPublicationCoverage(packs)
		coverage = publication.Inventory
		publicationDiags := curriculum.ValidatePublication(packs)
		diags = append(diags, publicationDiags...)
		if len(publicationDiags) > 0 {
			blocking = true
		}
		for _, f := range editorial {
			if f.Severity == curriculum.SeverityBlocking {
				blocking = true
			}
		}
		if *v1Gate {
			v1Findings = curriculum.CheckV1Gate(publication.Eligible)
			if len(v1Findings) > 0 {
				blocking = true
			}
		}
		if *v1Gate && *distribution == "" {
			*distribution = filepath.Join(cfg.WorkspaceRoot, "packs", "distribution.json")
		}
		if *distribution != "" {
			policy, err := curriculum.LoadDistributionPolicy(*distribution)
			if err != nil {
				fmt.Fprintf(stderr, "codinho: distribution: %v\n", err)
				return exitError
			}
			distributionFindings := policy.Check(packs)
			v1Findings = append(v1Findings, distributionFindings...)
			if len(distributionFindings) > 0 {
				blocking = true
			}
		}
		if *runChecks || *publishedChecks || *v1Gate {
			selected := packs
			if !*runChecks {
				selected = curriculum.PublishedAuthoringPacks(packs)
			}
			checkFindings, report, err := runEditorialProofs(selected)
			if err != nil {
				fmt.Fprintf(stderr, "codinho: catalog validate: %v\n", err)
				return exitError
			}
			proof = &report
			editorial = append(editorial, checkFindings...)
			if len(checkFindings) > 0 {
				blocking = true
			}
		}
	}

	if *jsonOut {
		out := catalogValidateOut{Coverage: coverage, Publication: publication, Proof: proof}
		for _, d := range diags {
			out.Diagnostics = append(out.Diagnostics, diagnosticOut{File: d.File, Item: d.Item, Field: d.Field, Code: string(d.Code), Detail: d.Detail, Blocking: d.Blocking})
		}
		for _, f := range editorial {
			out.Editorial = append(out.Editorial, editorialFindingOut{File: f.File, Item: f.Item, Rule: string(f.Rule), Severity: string(f.Severity), Detail: f.Detail, Suggestion: f.Suggestion})
		}
		for _, f := range v1Findings {
			out.V1Gate = append(out.V1Gate, editorialFindingOut{Item: f.Item, Rule: string(f.Rule), Severity: string(f.Severity), Detail: f.Detail, Suggestion: f.Suggestion})
		}
		if err := writeJSON(stdout, out); err != nil {
			fmt.Fprintf(stderr, "codinho: catalog validate: %v\n", err)
			return exitError
		}
	} else {
		fmt.Fprintf(stdout, "challenges: inventory=%d drafts=%d published=%d eligible=%d\n", publication.Inventory.Challenges, publication.Drafts.Challenges, publication.Published.Challenges, publication.Eligible.Challenges)
		if proof != nil {
			fmt.Fprintf(stdout, "checks: declared=%d verified=%d\n", proof.DeclaredChecks, proof.VerifiedChecks)
		}
		if len(diags) == 0 && len(editorial) == 0 && len(v1Findings) == 0 {
			fmt.Fprintln(stdout, "catalog: ok, no diagnostics")
		}
		for _, d := range diags {
			level := "warn"
			if d.Blocking {
				level = "error"
			}
			fmt.Fprintf(stdout, "%s: file=%s item=%s field=%s code=%s detail=%s\n", level, d.File, d.Item, d.Field, d.Code, d.Detail)
		}
		for _, f := range editorial {
			level := "warn"
			if f.Severity == curriculum.SeverityBlocking {
				level = "error"
			}
			fmt.Fprintf(stdout, "%s: file=%s item=%s rule=%s detail=%s suggestion=%s\n", level, f.File, f.Item, f.Rule, f.Detail, f.Suggestion)
		}
		for _, f := range v1Findings {
			fmt.Fprintf(stdout, "error: item=%s rule=%s detail=%s\n", f.Item, f.Rule, f.Detail)
		}
	}

	if blocking {
		return exitError
	}
	return exitOK
}

func runCatalogList(args []string, stdout, stderr io.Writer) int {
	fs := newFlagSet("catalog list", stderr)
	kind := fs.String("kind", "", "item kind: theme, concept, competency, track, challenge (default challenge)")
	theme := fs.String("theme", "", "restrict to items under this theme")
	jsonOut := fs.Bool("json", false, "print result as JSON")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}

	cfg := config.Load()
	catalog, _, err := loadCatalog(cfg)
	if err != nil {
		fmt.Fprintf(stderr, "codinho: catalog list: %v\n", err)
		return exitError
	}
	if catalog == nil {
		fmt.Fprintln(stderr, "codinho: catalog list: catalog failed to load; run catalog validate for details")
		return exitError
	}

	result, err := catalog.Search(curriculum.Query{Kind: curriculum.ItemKind(*kind), Theme: *theme})
	if err != nil {
		fmt.Fprintf(stderr, "codinho: catalog list: %v\n", err)
		return exitError
	}

	if *jsonOut {
		if err := writeJSON(stdout, result); err != nil {
			fmt.Fprintf(stderr, "codinho: catalog list: %v\n", err)
			return exitError
		}
	} else {
		for _, item := range result.Items {
			fmt.Fprintf(stdout, "%s\t%s\t%s\n", item.ID, item.Kind, item.Title)
		}
		if result.Truncated {
			fmt.Fprintln(stderr, "codinho: catalog list: results truncated")
		}
	}
	return exitOK
}

func runCatalogShow(args []string, stdout, stderr io.Writer) int {
	fs := newFlagSet("catalog show", stderr)
	jsonOut := fs.Bool("json", false, "print result as JSON")

	positional, flagArgs := extractPositional(args, nil)
	if err := fs.Parse(flagArgs); err != nil {
		return exitUsage
	}
	if len(positional) != 1 {
		fmt.Fprint(stderr, "Usage: codinho catalog show <id> [--json]\n")
		return exitUsage
	}
	id := positional[0]

	cfg := config.Load()
	catalog, _, err := loadCatalog(cfg)
	if err != nil {
		fmt.Fprintf(stderr, "codinho: catalog show: %v\n", err)
		return exitError
	}
	if catalog == nil {
		fmt.Fprintln(stderr, "codinho: catalog show: catalog failed to load; run catalog validate for details")
		return exitError
	}

	if ch, ok := catalog.Challenge(id); ok {
		if *jsonOut {
			if err := writeJSON(stdout, ch); err != nil {
				fmt.Fprintf(stderr, "codinho: catalog show: %v\n", err)
				return exitError
			}
			return exitOK
		}
		printChallenge(stdout, ch)
		return exitOK
	}

	item, ok := catalog.Get(id)
	if !ok {
		fmt.Fprintf(stderr, "codinho: catalog show: %q not found\n", id)
		return exitError
	}
	if *jsonOut {
		if err := writeJSON(stdout, item); err != nil {
			fmt.Fprintf(stderr, "codinho: catalog show: %v\n", err)
			return exitError
		}
		return exitOK
	}
	fmt.Fprintf(stdout, "id: %s\nkind: %s\ntitle: %s\n", item.ID, item.Kind, item.Title)
	return exitOK
}

func printChallenge(stdout io.Writer, ch curriculum.ChallengeAuthoring) {
	fmt.Fprintf(stdout, "id: %s\ntitle: %s\nkind: %s\ndifficulty: %s\nestimated_minutes: %d\n",
		ch.ID, ch.Title, ch.Kind, ch.Difficulty, ch.EstimatedMinutes)
	fmt.Fprintf(stdout, "themes: %v\n", ch.Themes)
	fmt.Fprintf(stdout, "competencies.primary: %v\n", ch.Competencies.Primary)
	fmt.Fprintf(stdout, "competencies.secondary: %v\n", ch.Competencies.Secondary)
	fmt.Fprintf(stdout, "prerequisites: %v\n", ch.Prerequisites)
	fmt.Fprintf(stdout, "brief: %s\n", ch.Brief)
	checks := make([]string, 0, len(ch.Checks))
	for _, c := range ch.Checks {
		checks = append(checks, c.ID)
	}
	sort.Strings(checks)
	fmt.Fprintf(stdout, "checks: %v\n", checks)
	if len(ch.Fixture) > 0 {
		paths := make([]string, 0, len(ch.Fixture))
		for _, f := range ch.Fixture {
			paths = append(paths, f.Path)
		}
		sort.Strings(paths)
		fmt.Fprintf(stdout, "fixture: %v\n", paths)
	}
}
