package mcpserver

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/oseiaspereira88/codinho/internal/application"
)

type progressGetArgs struct {
	CompetencyID string `json:"competency_id,omitempty" jsonschema:"restrict the result to one competency; omit for every competency with evidence"`
}

type reviewDueArgs struct{}

type masteryEvidenceRecordArgs struct {
	CompetencyID     string `json:"competency_id" jsonschema:"the competency this evidence relates to"`
	Dimension        string `json:"dimension" jsonschema:"understanding, syntax_recall, guided_implementation, autonomous_implementation, debugging, explanation, retention or transfer"`
	EvidenceID       string `json:"evidence_id" jsonschema:"an evidence ID already recorded elsewhere (an attempt, evaluation or reflection) — never new content"`
	Variant          string `json:"variant,omitempty" jsonschema:"identifies the challenge/variant this evidence came from, so transfer can be detected against an unseen one"`
	HelpUsed         bool   `json:"help_used,omitempty" jsonschema:"true when a hint or detour was used before this evidence"`
	SolutionRevealed bool   `json:"solution_revealed,omitempty" jsonschema:"true when the solution was revealed at or before this evidence; never promotes domain or scheduling"`
	Success          bool   `json:"success,omitempty" jsonschema:"true when the cited evidence demonstrates the dimension successfully"`
	ExpectedRevision uint64 `json:"expected_revision" jsonschema:"mastery stream revision this call expects, from progress_get or review_due"`
	RequestID        string `json:"request_id,omitempty" jsonschema:"idempotency key for retries"`
}

func registerProgressTools(server *mcp.Server, progress *application.ProgressService) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "progress_get",
		Description: "Present learner-wide progress by competency and dimension, recomputed from the durable log every call. Never a session-scoped view.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args progressGetArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		result, err := progress.Progress(args.CompetencyID)
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectNone, map[string]any{
			"competencies": result.Competencies,
			"revision":     result.Revision,
		})
		return nil, env, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "review_due",
		Description: "List overdue spaced reviews with an explainable priority (most overdue first). Advisory only: never blocks starting a free session.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, _ reviewDueArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		result, err := progress.ReviewDue(time.Now().UTC())
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectNone, map[string]any{
			"due":      result.Due,
			"revision": result.Revision,
		})
		return nil, env, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "mastery_evidence_record",
		Description: "Record one piece of mastery evidence citing an evidence_id already recorded elsewhere, and recalculate the competency's projection by transparent, deterministic rules. Never lets a caller set domain arbitrarily, and never promotes on a single guided completion or a revealed solution.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, IdempotentHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args masteryEvidenceRecordArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.CompetencyID == "" || args.Dimension == "" || args.EvidenceID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "competency_id, dimension and evidence_id are required", false, nil), nil
		}
		result, err := progress.RecordEvidence(application.EvidenceInput{
			CompetencyID:     args.CompetencyID,
			Dimension:        args.Dimension,
			EvidenceID:       args.EvidenceID,
			Variant:          args.Variant,
			HelpUsed:         args.HelpUsed,
			SolutionRevealed: args.SolutionRevealed,
			Success:          args.Success,
			ExpectedRevision: args.ExpectedRevision,
			RequestID:        args.RequestID,
		})
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectMasteryProjected, map[string]any{
			"competency_id": result.CompetencyID,
			"dimension":     result.Dimension,
			"state":         result.State,
			"revision":      result.Revision,
		})
		return nil, env, nil
	})
}
