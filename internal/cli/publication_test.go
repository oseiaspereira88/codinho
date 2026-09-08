package cli

import (
	"encoding/json"
	"testing"
)

func TestPublishedChecksEmptyCatalogReportsZero(t *testing.T) {
	t.Chdir(setupWorkspace(t))
	stdout, stderr, code := run(t, "catalog", "validate", "--published-checks", "--json")
	if code != exitOK {
		t.Fatalf("incremental published gate: %d %s %s", code, stdout, stderr)
	}
	var report catalogValidateOut
	if err := json.Unmarshal([]byte(stdout), &report); err != nil {
		t.Fatal(err)
	}
	if report.Proof == nil || report.Proof.DeclaredChecks != 0 || report.Proof.VerifiedChecks != 0 || report.Publication.Published.Challenges != 0 {
		t.Fatalf("zero publication misrepresented: %s", stdout)
	}
}
