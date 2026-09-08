package main

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
)

const conceptContentPack = `schema_version: 1
id: test-concept-content-pack
version: 1.0.0
themes:
  - id: meteorology
    title: Meteorologia
concepts:
  - id: barometric-pressure
    title: Pressão barométrica
    content:
      version: 1
      explanation: "A pressão barométrica mede a força por unidade de área exercida pela atmosfera."
      example:
        context: "meteorologia"
        code: "type Reading struct {\n\tHectopascals float64\n}"
        explanation: "Medição de pressão atmosférica em hPa."
      analogy: "Como o peso de uma coluna de ar sobre um ponto."
      relation_refs:
        - kind: relates_to
          concept_id: weather-station
  - id: weather-station
    title: Estação meteorológica
  - id: legacy-thermometer
    title: Termômetro legado
competencies:
  - id: read-barometer
    title: Ler barômetro digital
challenges:
  - schema_version: 1
    id: challenge.weather-reading
    version: 1.0.0
    title: Leitura meteorológica
    kind: atomic
    difficulty: foundational
    competencies:
      primary: [read-barometer]
    acceptance: ["Retorna leitura sem erro"]
    layers:
      - id: understanding
        macro_steps:
          - id: step.read
            kind: micro
            title: Executar leitura
            instruction:
              objective: Declarar a leitura barométrica.
              scope: Apenas o tipo de leitura.
            concepts: [barometric-pressure]
            hints:
              - level: 1
                kind: guiding_question
              - level: 2
                kind: syntax_recall
relations:
  - from: barometric-pressure
    to: weather-station
    kind: relates_to
`

func TestConceptContentOverRealStdio(t *testing.T) {
	bin := buildCodinhoBinary(t)
	root := t.TempDir()
	packsDir := filepath.Join(root, "packs")
	if err := os.MkdirAll(packsDir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(packsDir, "manifest.yaml"), []byte("schema_version: 1\npacks: [pack.yaml]\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(packsDir, "pack.yaml"), []byte(conceptContentPack), 0600); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin, "serve", "--authoring")
	cmd.Dir = root
	cs, err := mcp.NewClient(&mcp.Implementation{Name: "concept-content-test", Version: "1"}, nil).Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
	if err != nil {
		t.Fatalf("connecting over stdio: %v", err)
	}
	defer cs.Close()

	// 1. Start a session
	startEnv := callTool(ctx, t, cs, "session_start", map[string]any{
		"challenge_id":   "challenge.weather-reading",
		"disclosure_max": 6,
	})
	sessionID, ok := startEnv["session_id"].(string)
	if !ok || sessionID == "" {
		t.Fatalf("expected session_id, got: %+v", startEnv)
	}
	rev := revisionOf(t, startEnv)

	// 2. Request hint level 1
	hint1Env := callTool(ctx, t, cs, "hint_request", map[string]any{
		"session_id":        sessionID,
		"expected_revision": rev,
	})
	if hint1Env["status"] != "ok" || hint1Env["progress_effect"] != "none" {
		t.Fatalf("unexpected hint1 envelope: %+v", hint1Env)
	}
	hint1Data := hint1Env["data"].(map[string]any)
	if hint1Data["level"].(float64) != 1 {
		t.Fatalf("expected hint level 1, got %v", hint1Data["level"])
	}

	// Snapshot state before concept_content_get
	sessBefore := callTool(ctx, t, cs, "session_get", map[string]any{"session_id": sessionID})
	eventsLog := filepath.Join(root, ".codinho", "state", "events.jsonl")
	readEvents := func() []eventstore.Event {
		t.Helper()
		r := eventstore.ReadEvents(eventsLog, nil)
		if r.Err != nil {
			t.Fatal(r.Err)
		}
		return r.Events
	}
	eventsBefore := readEvents()

	// 3. Call concept_content_get for authored concept
	conceptEnv := callTool(ctx, t, cs, "concept_content_get", map[string]any{"concept_id": "barometric-pressure"})
	if conceptEnv["status"] != "ok" {
		t.Fatalf("concept_content_get failed: %+v", conceptEnv)
	}
	if conceptEnv["progress_effect"] != "none" {
		t.Fatalf("expected progress_effect=none, got %v", conceptEnv["progress_effect"])
	}

	data, ok := conceptEnv["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data map in concept response: %+v", conceptEnv)
	}
	if data["id"] != "barometric-pressure" || data["content_status"] != "available" {
		t.Fatalf("unexpected concept data: %+v", data)
	}
	content, ok := data["content"].(map[string]any)
	if !ok {
		t.Fatalf("expected content map: %+v", data)
	}
	if content["explanation"] != "A pressão barométrica mede a força por unidade de área exercida pela atmosfera." {
		t.Fatalf("unexpected explanation: %+v", content)
	}
	relations, ok := content["relations"].([]any)
	if !ok || len(relations) != 1 {
		t.Fatalf("unexpected relations: %+v", content)
	}
	rel0 := relations[0].(map[string]any)
	if rel0["kind"] != "relates_to" || rel0["concept_id"] != "weather-station" || rel0["title"] != "Estação meteorológica" {
		t.Fatalf("unexpected relation 0: %+v", rel0)
	}

	// 4. Call concept_content_get for legacy concept
	legacyEnv := callTool(ctx, t, cs, "concept_content_get", map[string]any{"concept_id": "legacy-thermometer"})
	if legacyEnv["status"] != "ok" || legacyEnv["progress_effect"] != "none" {
		t.Fatalf("unexpected legacy envelope: %+v", legacyEnv)
	}
	legacyData := legacyEnv["data"].(map[string]any)
	if legacyData["id"] != "legacy-thermometer" || legacyData["content_status"] != "missing" || legacyData["content"] != nil {
		t.Fatalf("unexpected legacy concept data: %+v", legacyData)
	}

	// 5. Verify that session state, revision, and events are 100% identical after concept queries
	sessAfter := callTool(ctx, t, cs, "session_get", map[string]any{"session_id": sessionID})
	if !reflect.DeepEqual(sessBefore["data"], sessAfter["data"]) ||
		!reflect.DeepEqual(sessBefore["active_node"], sessAfter["active_node"]) ||
		!reflect.DeepEqual(sessBefore["disclosure"], sessAfter["disclosure"]) ||
		sessBefore["session_id"] != sessAfter["session_id"] {
		t.Fatalf("session changed after concept queries: before=%+v, after=%+v", sessBefore, sessAfter)
	}

	eventsAfter := readEvents()
	if len(eventsBefore) != len(eventsAfter) {
		t.Fatalf("events recorded during read-only concept query: before=%d, after=%d", len(eventsBefore), len(eventsAfter))
	}

	// 6. Verify hint ladder progression continues seamlessly (level 2 comes next with expected_revision unchanged by concept queries)
	revAfterHint := revisionOf(t, hint1Env)
	hint2Env := callTool(ctx, t, cs, "hint_request", map[string]any{
		"session_id":        sessionID,
		"expected_revision": revAfterHint,
	})
	hintData := hint2Env["data"].(map[string]any)
	if hintData["level"].(float64) != 2 {
		t.Fatalf("expected hint level 2, got %v", hintData["level"])
	}

	// 7. Non-disclosure check: verify no leakage of solutions/fixtures
	rawConcept, _ := json.Marshal(conceptEnv)
	for _, forbidden := range []string{"challenge.weather-reading", "solution", "Retorna leitura sem erro"} {
		if strings.Contains(string(rawConcept), forbidden) {
			t.Fatalf("concept_content_get leaked %q: %s", forbidden, rawConcept)
		}
	}
}
