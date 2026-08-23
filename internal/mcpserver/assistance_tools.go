package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/oseiaspereira88/codinho/internal/application"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

type hintRequestArgs struct {
	SessionID        string `json:"session_id" jsonschema:"session ID returned by session_start"`
	ConfirmSolution  bool   `json:"confirm_solution,omitempty" jsonschema:"must be true when the next rung is level 6 (commented solution); required intent per PROJECT.md §8.4"`
	ExpectedRevision uint64 `json:"expected_revision" jsonschema:"session revision this call expects, from session_get"`
	RequestID        string `json:"request_id,omitempty" jsonschema:"idempotency key for retries"`
}

type syntaxRecallGetArgs struct {
	SessionID        string `json:"session_id" jsonschema:"session ID returned by session_start"`
	ExpectedRevision uint64 `json:"expected_revision" jsonschema:"session revision this call expects, from session_get"`
	RequestID        string `json:"request_id,omitempty" jsonschema:"idempotency key for retries"`
}

type conceptContentGetArgs struct {
	ConceptID string `json:"concept_id" jsonschema:"catalog concept ID"`
}

type learningDetourStartArgs struct {
	SessionID        string `json:"session_id" jsonschema:"session ID returned by session_start"`
	Reason           string `json:"reason" jsonschema:"the conceptual question or reason prompting the detour"`
	ExpectedRevision uint64 `json:"expected_revision" jsonschema:"session revision this call expects, from session_get"`
	RequestID        string `json:"request_id,omitempty" jsonschema:"idempotency key for retries"`
}

type learningDetourFinishArgs struct {
	SessionID        string `json:"session_id" jsonschema:"session ID returned by session_start"`
	Outcome          string `json:"outcome,omitempty" jsonschema:"how the detour concluded: resolved or abandoned (default resolved)"`
	ExpectedRevision uint64 `json:"expected_revision" jsonschema:"session revision this call expects, from session_get"`
	RequestID        string `json:"request_id,omitempty" jsonschema:"idempotency key for retries"`
}

// hintEnvelope shapes a HintResult into the shared Envelope, using
// ProgressEffectSolutionReveal only when the granted rung is the solution
// (PROJECT.md §15.3's fixed progress_effect vocabulary).
func hintEnvelope(requestID string, result application.HintResult) Envelope {
	effect := ProgressEffectNone
	if result.Level == learning.DisclosureSolution {
		effect = ProgressEffectSolutionReveal
	}
	env := okEnvelope(requestID, effect, map[string]any{
		"level":     int(result.Level),
		"kind":      result.Kind,
		"objective": result.Objective,
		"scope":     result.Scope,
		"concepts":  result.Concepts,
		"revision":  result.Revision,
	})
	env.ActiveNode = &ActiveNode{ID: string(result.StepID)}
	return env
}

func registerAssistanceTools(server *mcp.Server, sessions *application.SessionService, assist *application.AssistanceService) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "hint_request",
		Description: "Deliver the next allowed hint, climbing the assistance ladder by at most one level. Reaching the solution rung requires confirm_solution.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, IdempotentHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args hintRequestArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.SessionID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "session_id is required", false, nil), nil
		}
		result, err := sessions.HintRequest(learning.SessionID(args.SessionID), args.ConfirmSolution, args.ExpectedRevision, args.RequestID)
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := hintEnvelope(requestID, result)
		env.SessionID = args.SessionID
		return nil, env, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "syntax_recall_get",
		Description: "Return the active step's authored syntax-recall content. Free in teaching mode; otherwise consumes the ladder up to that rung.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, IdempotentHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args syntaxRecallGetArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.SessionID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "session_id is required", false, nil), nil
		}
		result, err := sessions.SyntaxRecallGet(learning.SessionID(args.SessionID), args.ExpectedRevision, args.RequestID)
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := hintEnvelope(requestID, result)
		env.SessionID = args.SessionID
		return nil, env, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "concept_content_get",
		Description: "Return the canonical catalog record for a concept. The calling agent adapts language to the learner's profile.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args conceptContentGetArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.ConceptID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "concept_id is required", false, nil), nil
		}
		content, err := assist.ConceptContent(args.ConceptID)
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		return nil, okEnvelope(requestID, ProgressEffectNone, map[string]string{"id": content.ID, "title": content.Title}), nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "learning_detour_start",
		Description: "Open a conceptual detour without changing the session's active step.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, IdempotentHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args learningDetourStartArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.SessionID == "" || args.Reason == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "session_id and reason are required", false, nil), nil
		}
		result, err := sessions.DetourStart(learning.SessionID(args.SessionID), args.Reason, args.ExpectedRevision, args.RequestID)
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectNone, map[string]any{"state": string(result.State), "revision": result.Revision})
		env.SessionID = args.SessionID
		return nil, env, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "learning_detour_finish",
		Description: "Close the session's open detour and return to the same active step without altering it.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, IdempotentHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args learningDetourFinishArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.SessionID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "session_id is required", false, nil), nil
		}
		outcome := application.DetourOutcome(args.Outcome)
		if outcome == "" {
			outcome = application.DetourResolved
		}
		if outcome != application.DetourResolved && outcome != application.DetourAbandoned {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "outcome must be resolved or abandoned", false, nil), nil
		}
		result, err := sessions.DetourFinish(learning.SessionID(args.SessionID), outcome, args.ExpectedRevision, args.RequestID)
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectNone, map[string]any{"state": string(result.State), "revision": result.Revision})
		env.SessionID = args.SessionID
		return nil, env, nil
	})
}
