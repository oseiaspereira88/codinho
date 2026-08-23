package mcpserver

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestContractMasteryEvidenceRecordThenProgressGetAndReviewDue(t *testing.T) {
	cs := newContractClient(t)
	ctx := context.Background()

	revRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "progress_get", Arguments: map[string]any{}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	revEnv := decodeEnvelope(t, revRes)
	if revEnv.Status != "ok" {
		t.Fatalf("unexpected progress_get envelope: %+v", revEnv)
	}
	rev := revEnv.Data.(map[string]any)["revision"].(float64)

	recordRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "mastery_evidence_record", Arguments: map[string]any{
		"competency_id": "slice-filter", "dimension": "autonomous_implementation", "evidence_id": "ev-1",
		"success": true, "expected_revision": rev,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	recordEnv := decodeEnvelope(t, recordRes)
	if recordEnv.Status != "ok" {
		t.Fatalf("unexpected mastery_evidence_record envelope: %+v (error=%+v)", recordEnv, recordEnv.Error)
	}
	if recordEnv.ProgressEffect != ProgressEffectMasteryProjected {
		t.Fatalf("progress_effect = %s, want mastery_projected", recordEnv.ProgressEffect)
	}
	recordData := recordEnv.Data.(map[string]any)
	if recordData["state"] != "demonstrates_without_help" {
		t.Fatalf("state = %v, want demonstrates_without_help", recordData["state"])
	}

	progressRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "progress_get", Arguments: map[string]any{"competency_id": "slice-filter"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	progressEnv := decodeEnvelope(t, progressRes)
	if progressEnv.Status != "ok" {
		t.Fatalf("unexpected progress_get envelope: %+v", progressEnv)
	}
	competencies := progressEnv.Data.(map[string]any)["competencies"].(map[string]any)
	if _, ok := competencies["slice-filter"]; !ok {
		t.Fatalf("expected slice-filter in progress, got %+v", competencies)
	}

	dueRes, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "review_due", Arguments: map[string]any{}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	dueEnv := decodeEnvelope(t, dueRes)
	if dueEnv.Status != "ok" {
		t.Fatalf("unexpected review_due envelope: %+v", dueEnv)
	}
	// Freshly scheduled 1 day out: nothing is due yet.
	due, _ := dueEnv.Data.(map[string]any)["due"].([]any)
	if len(due) != 0 {
		t.Fatalf("expected nothing due yet, got %+v", due)
	}
}

func TestContractMasteryEvidenceRecordRejectsUnknownDimension(t *testing.T) {
	cs := newContractClient(t)
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "mastery_evidence_record", Arguments: map[string]any{
		"competency_id": "c", "dimension": "not-a-real-dimension", "evidence_id": "e1", "expected_revision": 0,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected IsError for an unknown dimension")
	}
	env := decodeEnvelope(t, res)
	if env.Error == nil || env.Error.Code != string(ErrCodeInvalidInput) {
		t.Fatalf("unexpected envelope: %+v", env)
	}
}

func TestContractMasteryEvidenceRecordGuidedNeverExceedsDemonstratesWithHelp(t *testing.T) {
	cs := newContractClient(t)
	ctx := context.Background()
	rev := 0.0

	for i := 0; i < 5; i++ {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "mastery_evidence_record", Arguments: map[string]any{
			"competency_id": "guided-comp", "dimension": "guided_implementation", "evidence_id": "ev",
			"success": true, "help_used": true, "expected_revision": rev,
		}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		env := decodeEnvelope(t, res)
		if env.Status != "ok" {
			t.Fatalf("unexpected envelope: %+v (error=%+v)", env, env.Error)
		}
		data := env.Data.(map[string]any)
		if data["state"] != "demonstrates_with_help" {
			t.Fatalf("iteration %d: state = %v, want demonstrates_with_help (requirement R4)", i, data["state"])
		}
		rev = data["revision"].(float64)
	}
}
