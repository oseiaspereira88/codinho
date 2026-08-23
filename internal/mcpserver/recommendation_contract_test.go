package mcpserver

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestContractLearningPathRecommendExcludesUnmetPrerequisite(t *testing.T) {
	cs := newContractClient(t)
	ctx := context.Background()

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "learning_path_recommend", Arguments: map[string]any{
		"competency_id": "slice-filter",
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env := decodeEnvelope(t, res)
	if env.Status != "ok" {
		t.Fatalf("unexpected envelope: %+v (error=%+v)", env, env.Error)
	}
	recs, _ := env.Data.(map[string]any)["recommendations"].([]any)
	if len(recs) != 1 {
		t.Fatalf("expected exactly the fixture challenge, got %+v", recs)
	}
	first := recs[0].(map[string]any)
	if first["challenge_id"] != fixtureChallengeID {
		t.Fatalf("expected %s, got %+v", fixtureChallengeID, first)
	}
	explanation := first["explanation"].(map[string]any)
	for _, field := range []string{"evidence", "gap", "dependency", "cost_minutes"} {
		if _, ok := explanation[field]; !ok {
			t.Fatalf("expected explanation.%s to be present, got %+v", field, explanation)
		}
	}
}

func TestContractConceptRelationsGetReturnsBothDirections(t *testing.T) {
	cs := newContractClient(t)
	ctx := context.Background()

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "concept_relations_get", Arguments: map[string]any{
		"id": "slice-declaration",
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env := decodeEnvelope(t, res)
	if env.Status != "ok" {
		t.Fatalf("unexpected envelope: %+v (error=%+v)", env, env.Error)
	}
	data := env.Data.(map[string]any)
	incoming, _ := data["incoming"].([]any)
	if len(incoming) != 1 {
		t.Fatalf("expected 1 incoming relation (requires, from the fixture challenge), got %+v", data)
	}
}

func TestContractConceptRelationsGetUnknownIDReturnsItemNotFound(t *testing.T) {
	cs := newContractClient(t)
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "concept_relations_get", Arguments: map[string]any{
		"id": "does-not-exist",
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected IsError for an unknown ID")
	}
	env := decodeEnvelope(t, res)
	if env.Error == nil || env.Error.Code != string(ErrCodeItemNotFound) {
		t.Fatalf("unexpected envelope: %+v", env)
	}
}

func TestContractCatalogSearchByDifficultyAndCompetency(t *testing.T) {
	cs := newContractClient(t)
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "catalog_search", Arguments: map[string]any{
		"difficulty": "foundational", "competency": "slice-filter",
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env := decodeEnvelope(t, res)
	if env.Status != "ok" {
		t.Fatalf("unexpected envelope: %+v", env)
	}
	items, _ := env.Data.([]any)
	if len(items) != 1 {
		t.Fatalf("expected 1 match, got %+v", items)
	}
}

func TestContractCatalogSearchRejectsOverlongText(t *testing.T) {
	cs := newContractClient(t)
	longText := ""
	for range 300 {
		longText += "a"
	}
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "catalog_search", Arguments: map[string]any{
		"text": longText,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected IsError for an overlong search text")
	}
	env := decodeEnvelope(t, res)
	if env.Error == nil || env.Error.Code != string(ErrCodeInvalidInput) {
		t.Fatalf("unexpected envelope: %+v", env)
	}
}
