package curriculum

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNavigationCatalogRejectsAmbiguousChoices(t *testing.T) {
	for _, tc := range []struct {
		name, mode, second string
		bad                bool
	}{
		{"default", "", "b", false}, {"sequence", "sequence", "b", false}, {"choice", "choice", "b", false}, {"unknown", "random", "b", true}, {"single", "choice", "", true}, {"duplicate", "choice", "a", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pack := "schema_version: 1\nid: pack\nversion: 1.0.0\nchallenges:\n  - schema_version: 1\n    id: ch\n    version: 1.0.0\n    layers:\n      - id: layer\n        macro_steps:\n          - id: macro\n            kind: macro\n            children_mode: " + tc.mode + "\n            children:\n              - id: a\n                kind: micro\n"
			if tc.second != "" {
				pack += "              - id: " + tc.second + "\n                kind: micro\n"
			}
			dir := t.TempDir()
			for name, data := range map[string]string{"manifest.yaml": "schema_version: 1\npacks: [pack.yaml]\n", "pack.yaml": pack} {
				if err := os.WriteFile(filepath.Join(dir, name), []byte(data), 0600); err != nil {
					t.Fatal(err)
				}
			}
			catalog, diags, err := Load(dir, DefaultLimits)
			if tc.bad {
				found := false
				for _, d := range diags {
					found = found || (d.Code == DiagInvalidNavigation && d.Blocking)
				}
				if !found || catalog != nil {
					t.Fatalf("accepted malformed navigation: %v %+v %v", catalog, diags, err)
				}
			} else if err != nil || len(diags) > 0 {
				t.Fatalf("valid tree: %+v %v", diags, err)
			}
		})
	}
}

func TestNavigationEmptyFieldPreservesLegacyJSON(t *testing.T) {
	raw, err := json.Marshal(StepAuthoring{ID: "old", Kind: "macro"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "ChildrenMode") {
		t.Fatal("empty field changes legacy pinned-content digest")
	}
	var roundtrip StepAuthoring
	if err := json.Unmarshal(raw, &roundtrip); err != nil {
		t.Fatal(err)
	}
	again, _ := json.Marshal(roundtrip)
	if string(raw) != string(again) {
		t.Fatal("legacy JSON changed")
	}
}
