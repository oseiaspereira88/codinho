package mcpserver

import (
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/oseiaspereira88/codinho/internal/application"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

type workspaceObserveArgs struct {
	SessionID        string   `json:"session_id" jsonschema:"session ID returned by session_start"`
	StepID           string   `json:"step_id" jsonschema:"the step this observation belongs to"`
	Root             string   `json:"root" jsonschema:"absolute path to the learner's authorized workspace root"`
	Globs            []string `json:"globs,omitempty" jsonschema:"glob patterns (supporting **) restricting which files are collected; empty matches every non-excluded file"`
	ExpectedRevision uint64   `json:"expected_revision" jsonschema:"session revision this call expects, from session_get"`
	RequestID        string   `json:"request_id,omitempty" jsonschema:"idempotency key for retries"`
}

type evidenceGetArgs struct {
	SessionID  string   `json:"session_id" jsonschema:"session ID returned by session_start"`
	EvidenceID string   `json:"evidence_id" jsonschema:"evidence ID returned by workspace_observe"`
	Root       string   `json:"root,omitempty" jsonschema:"workspace root, to check whether it changed since this evidence was collected"`
	Globs      []string `json:"globs,omitempty" jsonschema:"the same globs used when the evidence was collected, for the staleness check"`
}

func registerWorkspaceTools(server *mcp.Server, workspaces *application.WorkspaceService) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "workspace_observe",
		Description: "Observe the learner's workspace read-only: the first call for a step establishes its baseline, later calls report the relevant diff restricted to globs. Never edits files, never implies an attempt or evaluation.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args workspaceObserveArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.SessionID == "" || args.StepID == "" || args.Root == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "session_id, step_id and root are required", false, nil), nil
		}
		result, err := workspaces.Observe(application.ObserveInput{
			SessionID:        learning.SessionID(args.SessionID),
			StepID:           learning.StepID(args.StepID),
			Root:             args.Root,
			Globs:            args.Globs,
			ExpectedRevision: args.ExpectedRevision,
			RequestID:        args.RequestID,
		})
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectNone, map[string]any{
			"baseline":    result.Baseline,
			"changes":     result.Changes,
			"from_git":    result.FromGit,
			"dirty":       result.Dirty,
			"evidence_id": result.EvidenceID,
			"fingerprint": result.Fingerprint,
			"revision":    result.Revision,
		})
		env.SessionID = args.SessionID
		env.ActiveNode = &ActiveNode{ID: args.StepID}
		return nil, env, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "evidence_get",
		Description: "Fetch previously recorded workspace evidence by ID, scoped to the requesting session. Reports whether the workspace changed since the evidence was collected.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args evidenceGetArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.SessionID == "" || args.EvidenceID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "session_id and evidence_id are required", false, nil), nil
		}
		result, err := workspaces.EvidenceGet(application.EvidenceGetInput{
			SessionID:  learning.SessionID(args.SessionID),
			EvidenceID: args.EvidenceID,
			Root:       args.Root,
			Globs:      args.Globs,
		})
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectNone, map[string]any{
			"content": json.RawMessage(result.Content),
			"size":    result.Size,
			"stale":   result.Stale,
		})
		env.SessionID = args.SessionID
		return nil, env, nil
	})
}
