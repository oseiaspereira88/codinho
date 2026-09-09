package mcpserver

import (
	"context"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/oseiaspereira88/codinho/internal/application"
)

type draftSubmitArgs struct {
	PackYAML  string `json:"pack_yaml" jsonschema:"self-contained unreviewed authoring pack YAML, at most 256 KiB; never executed on submission"`
	RequestID string `json:"request_id" jsonschema:"required unique retry identity; reuse only with identical YAML"`
}
type draftGetArgs struct {
	DraftID string `json:"draft_id" jsonschema:"draft identity returned by content_draft_submit"`
}
type draftRemoveArgs struct {
	DraftID   string `json:"draft_id"`
	RequestID string `json:"request_id"`
	Confirm   bool   `json:"confirm" jsonschema:"explicit consent to prevent new uses; audit data remains until privacy purge"`
}

func registerContentDraftTools(server *mcp.Server, sessions *application.SessionService) {
	mcp.AddTool(server, &mcp.Tool{Name: "content_draft_submit", Description: "Validate host-authored YAML and persist it in quarantine. Does not generate, execute, materialize or publish content. Draft practice requires explicit learner consent.", Annotations: &mcp.ToolAnnotations{IdempotentHint: true}}, func(_ context.Context, req *mcp.CallToolRequest, args draftSubmitArgs) (*mcp.CallToolResult, Envelope, error) {
		id := requestIDFor(req)
		result, validation, err := sessions.Drafts().Submit(args.PackYAML, args.RequestID)
		if err != nil {
			code, msg, retry := mapError(err)
			env := errorEnvelope(id, code, msg, retry, nil)
			env.Data = validation
			return errorResult(), env, nil
		}
		return nil, okEnvelope(id, ProgressEffectNone, map[string]any{"draft_id": result.ID, "pack_id": result.PackID, "version": result.Version, "digest": result.Digest, "expires_at": result.ExpiresAt, "content_provenance": "draft", "content_warning": contentWarning("draft"), "validation": validation}), nil
	})
	mcp.AddTool(server, &mcp.Tool{Name: "content_draft_get", Description: "Explicit authoring review/export: return the quarantined source YAML, including private fixture and solution material. Do not call during learner practice.", Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true}}, func(_ context.Context, req *mcp.CallToolRequest, args draftGetArgs) (*mcp.CallToolResult, Envelope, error) {
		id := requestIDFor(req)
		result, err := sessions.Drafts().Get(args.DraftID)
		if err != nil {
			code, msg, retry := mapError(err)
			return errorResult(), errorEnvelope(id, code, msg, retry, nil), nil
		}
		return nil, okEnvelope(id, ProgressEffectNone, result), nil
	})
	mcp.AddTool(server, &mcp.Tool{Name: "content_draft_remove", Description: "With explicit confirmation, prevent new starts of this draft. Pinned sessions and audit history remain; privacy purge removes local state.", Annotations: &mcp.ToolAnnotations{IdempotentHint: true, DestructiveHint: boolPointer(true)}}, func(_ context.Context, req *mcp.CallToolRequest, args draftRemoveArgs) (*mcp.CallToolResult, Envelope, error) {
		id := requestIDFor(req)
		err := sessions.Drafts().Remove(args.DraftID, args.RequestID, args.Confirm)
		if err != nil {
			code, msg, retry := mapError(err)
			return errorResult(), errorEnvelope(id, code, msg, retry, nil), nil
		}
		return nil, okEnvelope(id, ProgressEffectNone, map[string]any{"draft_id": args.DraftID, "removed": true}), nil
	})
}

func boolPointer(value bool) *bool { return &value }
