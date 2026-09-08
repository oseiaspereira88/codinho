package mcpserver

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// The governance inventory must match the actual protocol surface, not a
// hand-maintained assertion that the integration detector found no gaps.
func TestContractGovernanceInventoryMatchesServer(t *testing.T) {
	raw, err := os.ReadFile("../../.pose/contracts/mcp-stdio.json")
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		SchemaVersion   int      `json:"schema_version"`
		Protocol        string   `json:"protocol"`
		Transport       string   `json:"transport"`
		Provider        string   `json:"provider"`
		Entrypoint      string   `json:"entrypoint"`
		Consumer        string   `json:"consumer"`
		ValidationCheck string   `json:"validation_check"`
		Tools           []string `json:"tools"`
	}
	if err := json.Unmarshal(raw, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.SchemaVersion != 1 || manifest.Protocol != "mcp" || manifest.Transport != "stdio" || manifest.ValidationCheck != "mcp-contract" {
		t.Fatal("invalid MCP inventory identity")
	}
	for _, path := range []string{manifest.Provider, manifest.Entrypoint, manifest.Consumer} {
		if !filepath.IsLocal(path) {
			t.Fatalf("non-local contract path %q", path)
		}
		if info, err := os.Stat(filepath.Join("../..", path)); err != nil || info.IsDir() {
			t.Fatalf("contract file missing %q", path)
		}
	}
	want := map[string]bool{}
	for _, name := range manifest.Tools {
		if name == "" || want[name] {
			t.Fatalf("invalid/duplicate tool %q", name)
		}
		want[name] = true
	}
	cs := newContractClient(t)
	result, err := cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if result.NextCursor != "" {
		t.Fatal("inventory comparison requires all pages")
	}
	for _, tool := range result.Tools {
		if !want[tool.Name] {
			t.Fatalf("undeclared/duplicate runtime tool %q", tool.Name)
		}
		delete(want, tool.Name)
	}
	if len(want) != 0 {
		t.Fatalf("declared tools absent from server: %v", want)
	}
}
