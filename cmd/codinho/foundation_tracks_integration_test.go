package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
)

var foundationTrackIDs = []string{
	"go-from-zero",
	"go-oop-transition",
	"go-practical-fluency",
	"go-data-and-types",
	"go-testing-and-design",
}

var foundationTrackPaths = map[string][]string{
	"go-from-zero": {
		"go-first-steps.declare-a-minimal-module",
		"go-first-steps.enumerate-weekdays-with-iota",
		"go-first-steps.convert-celsius-to-fahrenheit",
		"go-first-steps.clamp-int-to-byte",
		"go-first-steps.format-rate-limit-message",
		"go-first-steps.avoid-shadowing-named-returns",
		"go-data-text.filter-without-mutating-input",
		"go-core.sum-variadic-numbers",
		"go-data-text.lookup-map-value-with-comma-ok",
		"go-errors.wrap-sentinel-error-with-context",
		"go-errors.maintain-memory-catalog",
	},
	"go-oop-transition": {
		"go-core.increment-counter-with-pointer-receiver",
		"go-core.guard-invariant-with-unexported-field",
		"go-core.hide-counter-type-behind-interface",
		"go-type-design.select-validator-true-nil",
		"go-type-design.sum-list-with-nil-safe-receiver",
		"go-type-design.clone-inventory-independently",
		"go-type-design.first-circle-clone",
		"go-core.configure-server-with-functional-options",
		"go-testing.select-active-tokens-with-injected-clock",
		"go-errors.maintain-memory-catalog",
	},
	"go-practical-fluency": {
		"go-core.find-first-value-at-least",
		"go-data-text.dedupe-preserving-first-occurrence",
		"go-data-text.sum-valid-integers-safely",
		"go-errors.wrap-sentinel-error-with-context",
		"go-errors.extract-field-with-errors-as",
		"go-errors.classify-wrapped-errors",
		"go-io.decode-strict-config",
		"go-io.import-csv-inventory",
		"go-errors.maintain-memory-catalog",
	},
	"go-data-and-types": {
		"go-data-text.filter-without-mutating-input",
		"go-data-text.preallocate-slice-with-zero-length",
		"go-data-text.truncate-bytes-at-rune-boundary",
		"go-data-text.title-first-letters-preserving-spacing",
		"go-data-text.quote-csv-field-when-needed",
		"go-data-text.format-names-as-csv-fields",
		"go-type-design.clone-inventory-independently",
		"go-type-design.contains-generic-comparable",
		"go-type-design.map-generic-transform",
		"go-type-design.generic-node-values-nil-safe",
	},
	"go-testing-and-design": {
		"go-type-design.clone-inventory-independently",
		"go-core.guard-invariant-with-unexported-field",
		"go-testing.select-active-tokens-with-injected-clock",
		"go-errors.wrap-sentinel-error-with-context",
		"go-errors.validate-user-joining-all-errors",
		"go-errors.classify-wrapped-errors",
		"go-errors.maintain-memory-catalog",
	},
}

func assertFoundationTrackAuthoring(t *testing.T, catalog *curriculum.Catalog, id string, want []string) {
	t.Helper()
	track, ok := catalog.Track(id)
	if !ok {
		t.Fatalf("missing authored track %q", id)
	}
	if !reflect.DeepEqual(track.ChallengeIDs, want) {
		t.Fatalf("track %q path = %v, want %v", id, track.ChallengeIDs, want)
	}
	if _, err := catalog.ResolvePath(track.ChallengeIDs); err != nil {
		t.Fatalf("track %q does not satisfy authored dependency order: %v", id, err)
	}

	positions := make(map[string]int, len(track.ChallengeIDs))
	for position, challengeID := range track.ChallengeIDs {
		positions[challengeID] = position
	}
	for position, challengeID := range track.ChallengeIDs {
		challenge, ok := catalog.Challenge(challengeID)
		if !ok {
			t.Fatalf("track %q references missing challenge %q", id, challengeID)
		}
		for _, prerequisite := range challenge.Prerequisites {
			if _, ok := catalog.Concept(prerequisite); ok {
				continue
			}
			if _, ok := catalog.Challenge(prerequisite); !ok {
				t.Fatalf("challenge %q references missing prerequisite %q", challengeID, prerequisite)
			}
			prerequisitePosition, ok := positions[prerequisite]
			if !ok || prerequisitePosition >= position {
				t.Fatalf("track %q places prerequisite %q after %q", id, prerequisite, challengeID)
			}
		}
	}
}

// TestFoundationTracksOverRealStdio checks the actual authored catalog, not
// another synthetic graph. Override traversal is a routing check, not playtest:
// it bypasses learner completion only to inspect every authored position.
func TestFoundationTracksOverRealStdio(t *testing.T) {
	bin := buildCodinhoBinary(t)
	root := t.TempDir()
	if err := os.CopyFS(filepath.Join(root, "packs"), os.DirFS("../../packs")); err != nil {
		t.Fatal(err)
	}
	catalog, diagnostics, err := curriculum.Load(filepath.Join(root, "packs"), curriculum.DefaultLimits)
	if err != nil {
		t.Fatal(err)
	}
	for _, diagnostic := range diagnostics {
		if diagnostic.Blocking {
			t.Fatalf("authored catalog has blocking diagnostic: %+v", diagnostic)
		}
	}
	if catalog == nil {
		t.Fatal("authored catalog was not materialized")
	}
	for _, id := range foundationTrackIDs {
		assertFoundationTrackAuthoring(t, catalog, id, foundationTrackPaths[id])
	}
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, "serve", "--authoring")
	cmd.Dir = root
	cs, err := mcp.NewClient(&mcp.Implementation{Name: "foundation-tracks-test", Version: "1"}, nil).Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer cs.Close()
	listed := callTool(ctx, t, cs, "catalog_search", map[string]any{"kind": "track"})
	items, ok := listed["data"].([]any)
	if !ok || len(items) != 5 {
		t.Fatalf("expected five authored tracks: %+v", listed)
	}
	listedIDs := make(map[string]bool, len(items))
	for _, item := range items {
		fields, ok := item.(map[string]any)
		if !ok {
			t.Fatalf("unexpected track search item: %#v", item)
		}
		id, _ := fields["ID"].(string)
		if id == "" {
			id, _ = fields["id"].(string)
		}
		if id == "" {
			t.Fatalf("track search item has no ID: %#v", item)
		}
		listedIDs[id] = true
	}
	for _, id := range foundationTrackIDs {
		if !listedIDs[id] {
			t.Fatalf("track search omitted %q: %+v", id, listed)
		}
	}
	for _, id := range foundationTrackIDs {
		t.Run(id, func(t *testing.T) {
			current := callTool(ctx, t, cs, "session_start", map[string]any{"track_id": id, "depth": "challenge", "request_id": id})
			sid := current["session_id"]
			data := envData(t, current)
			track := data["track"].(map[string]any)
			total := int(track["total"].(float64))
			want := foundationTrackPaths[id]
			if track["track_id"] != id || total != len(want) || data["content_provenance"] != "draft" {
				t.Fatalf("unexpected authored path: %+v", data)
			}
			seen := map[string]bool{}
			got := make([]string, 0, total)
			for cursor := 0; cursor < total; cursor++ {
				status := callTool(ctx, t, cs, "session_get", map[string]any{"session_id": sid})
				data = envData(t, status)
				track = data["track"].(map[string]any)
				challenge := track["challenge_id"].(string)
				if int(track["cursor"].(float64)) != cursor || challenge == "" || seen[challenge] {
					t.Fatalf("invalid traversal: %+v", track)
				}
				seen[challenge] = true
				got = append(got, challenge)
				instruction := callTool(ctx, t, cs, "instruction_get", map[string]any{"session_id": sid})
				if objective, _ := envData(t, instruction)["objective"].(string); objective == "" {
					t.Fatal("unreachable authored instruction", instruction)
				}
				if cursor+1 < total {
					callTool(ctx, t, cs, "step_advance", map[string]any{"session_id": sid, "expected_revision": revisionOf(t, status), "override": true})
				}
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("MCP traversal for %q = %v, want %v", id, got, want)
			}
		})
	}
}
