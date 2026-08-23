package cli

import (
	"testing"

	"github.com/oseiaspereira88/codinho/internal/curriculum"
)

func TestRunChecksAgainstFixturesSkipsChallengesWithoutBoth(t *testing.T) {
	packs := []curriculum.Pack{{File: "f", Challenges: []curriculum.ChallengeAuthoring{
		{ID: "no-fixture", Checks: []curriculum.CheckAuthoring{{ID: "c", Runner: "internal_ast", Package: "main.go"}}},
		{ID: "no-checks", Fixture: []curriculum.FixtureFileAuthoring{{Path: "main.go", Content: "package main\n"}}},
	}}}
	findings, err := runChecksAgainstFixtures(packs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("findings = %+v, want none: a challenge needs both fixture and checks to be exercised", findings)
	}
}

func TestRunChecksAgainstFixturesExecutesRealInternalASTCheck(t *testing.T) {
	packs := []curriculum.Pack{{File: "f", Challenges: []curriculum.ChallengeAuthoring{{
		ID:      "ch1",
		Fixture: []curriculum.FixtureFileAuthoring{{Path: "main.go", Content: "package main\n\nfunc main() {}\n"}},
		Checks:  []curriculum.CheckAuthoring{{ID: "parses", Runner: "internal_ast", Package: "main.go"}},
	}}}}
	findings, err := runChecksAgainstFixtures(packs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("findings = %+v, want none: valid Go syntax must parse cleanly", findings)
	}
}

func TestRunChecksAgainstFixturesToleratesAFailingOutcome(t *testing.T) {
	// A syntactically invalid fixture is a legitimate reproducible
	// outcome (Outcome fail), never a blocking finding: this gate only
	// blocks on infrastructure-level errors, never on the fixture's own
	// content failing its check (a debug challenge's fixture is *meant*
	// to fail until fixed).
	packs := []curriculum.Pack{{File: "f", Challenges: []curriculum.ChallengeAuthoring{{
		ID:      "ch1",
		Fixture: []curriculum.FixtureFileAuthoring{{Path: "main.go", Content: "not valid go syntax {{{"}},
		Checks:  []curriculum.CheckAuthoring{{ID: "parses", Runner: "internal_ast", Package: "main.go"}},
	}}}}
	findings, err := runChecksAgainstFixtures(packs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("findings = %+v, want none: a failing outcome is reproducible, not an infrastructure error", findings)
	}
}

func TestRunChecksAgainstFixturesFlagsUnresolvableCheck(t *testing.T) {
	packs := []curriculum.Pack{{File: "f", Challenges: []curriculum.ChallengeAuthoring{{
		ID:      "ch1",
		Fixture: []curriculum.FixtureFileAuthoring{{Path: "main.go", Content: "package main\n"}},
		Checks:  []curriculum.CheckAuthoring{{ID: "bogus", Runner: "not_a_real_runner"}},
	}}}}
	findings, err := runChecksAgainstFixtures(packs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 1 || findings[0].Rule != curriculum.RuleCheckNotResolvable {
		t.Fatalf("findings = %+v, want exactly one RuleCheckNotResolvable", findings)
	}
}
