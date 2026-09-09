package mcpserver

import (
	"context"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"testing"
)

func TestContractDraftSubmissionRejectsInvalidPayloadWithoutPublishing(t *testing.T) {
	cs := newContractClient(t)
	for _, raw := range []string{"[", "schema_version: 1\nid: synthetic\nversion: 1.0.0\npublication: {status: published}\n"} {
		res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "content_draft_submit", Arguments: map[string]any{"pack_yaml": raw, "request_id": "invalid"}})
		if err != nil {
			t.Fatal(err)
		}
		env := decodeEnvelope(t, res)
		if !res.IsError || env.Error == nil || env.Error.Code != string(ErrCodeInvalidInput) {
			t.Fatalf("invalid draft accepted: %+v", env)
		}
		if env.Data == nil {
			t.Fatal("missing actionable validation diagnostics")
		}
	}
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "content_draft_get", Arguments: map[string]any{"draft_id": "missing"}})
	if err != nil {
		t.Fatal(err)
	}
	if env := decodeEnvelope(t, res); !res.IsError || env.Error.Code != string(ErrCodeItemNotFound) {
		t.Fatal(env)
	}
}
