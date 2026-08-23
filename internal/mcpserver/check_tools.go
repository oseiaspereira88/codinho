package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/oseiaspereira88/codinho/internal/application"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

type checkRunArgs struct {
	SessionID        string `json:"session_id" jsonschema:"session ID returned by session_start"`
	CheckID          string `json:"check_id" jsonschema:"a check_id declared by the session's fixed challenge"`
	ExpectedRevision uint64 `json:"expected_revision" jsonschema:"session revision this call expects, from session_get"`
	RequestID        string `json:"request_id,omitempty" jsonschema:"idempotency key for retries"`
}

func registerCheckTools(server *mcp.Server, checksSvc *application.ChecksService) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "check_run",
		Description: "Execute a check_id already declared by the session's fixed challenge, inside the workspace root the last workspace_observe call established for the active step. No free command parameter: only pre-declared, allowlisted checks ever run.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, IdempotentHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args checkRunArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.SessionID == "" || args.CheckID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "session_id and check_id are required", false, nil), nil
		}
		result, err := checksSvc.Run(application.CheckRunInput{
			SessionID:        learning.SessionID(args.SessionID),
			CheckID:          args.CheckID,
			ExpectedRevision: args.ExpectedRevision,
			RequestID:        args.RequestID,
		})
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectNone, map[string]any{
			"outcome":     result.Outcome,
			"evidence_id": result.EvidenceID,
			"fingerprint": result.Fingerprint,
			"revision":    result.Revision,
		})
		env.SessionID = args.SessionID
		return nil, env, nil
	})
}
