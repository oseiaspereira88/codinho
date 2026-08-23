package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/oseiaspereira88/codinho/internal/application"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

type feedbackPrepareArgs struct {
	SessionID string `json:"session_id" jsonschema:"session ID returned by session_start"`
	Question  string `json:"question,omitempty" jsonschema:"the learner's own question, if any"`
}

type feedbackRecordArgs struct {
	SessionID        string `json:"session_id" jsonschema:"session ID returned by session_start"`
	Type             string `json:"type" jsonschema:"confirmation, question, explanation, relation, suggestion, idiom, risk, violation or error"`
	Text             string `json:"text" jsonschema:"the feedback text already authored by the caller"`
	Blocking         *bool  `json:"blocking,omitempty" jsonschema:"override the type's default blocks-completion classification"`
	ExpectedRevision uint64 `json:"expected_revision" jsonschema:"session revision this call expects, from session_get"`
	RequestID        string `json:"request_id,omitempty" jsonschema:"idempotency key for retries"`
}

type criterionArg struct {
	Name       string `json:"name" jsonschema:"criterion name"`
	Kind       string `json:"kind,omitempty" jsonschema:"structural (server-derived from evidence presence) or a qualitative kind judged by the caller"`
	Severity   string `json:"severity" jsonschema:"blocking, important_non_blocking or advisory"`
	Verdict    string `json:"verdict,omitempty" jsonschema:"met, partially_met, not_met, unverifiable or not_applicable; ignored for structural criteria"`
	EvidenceID string `json:"evidence_id,omitempty" jsonschema:"required for qualitative criteria"`
	RubricRef  string `json:"rubric_ref,omitempty" jsonschema:"required for qualitative criteria, e.g. rubric://idiomatic-go"`
}

type stepEvaluateArgs struct {
	SessionID        string         `json:"session_id" jsonschema:"session ID returned by session_start"`
	Criteria         []criterionArg `json:"criteria,omitempty" jsonschema:"criterion judgments for this evaluation"`
	SubmissionIntent bool           `json:"submission_intent,omitempty" jsonschema:"true when the learner deliberately submitted this for evaluation, creating an attempt"`
	ExpectedRevision uint64         `json:"expected_revision" jsonschema:"session revision this call expects, from session_get"`
	RequestID        string         `json:"request_id,omitempty" jsonschema:"idempotency key for retries"`
}

type reflectionRecordArgs struct {
	SessionID        string `json:"session_id" jsonschema:"session ID returned by session_start"`
	CompetencyID     string `json:"competency_id,omitempty" jsonschema:"the competency this reflection relates to"`
	Prompt           string `json:"prompt" jsonschema:"the reflection prompt asked"`
	Answer           string `json:"answer" jsonschema:"the learner's answer"`
	Assessment       string `json:"assessment,omitempty" jsonschema:"the tutor's descriptive, non-scoring read of the answer's quality"`
	ExpectedRevision uint64 `json:"expected_revision" jsonschema:"session revision this call expects, from session_get"`
	RequestID        string `json:"request_id,omitempty" jsonschema:"idempotency key for retries"`
}

type stepCompleteArgs struct {
	SessionID        string `json:"session_id" jsonschema:"session ID returned by session_start"`
	Confirm          bool   `json:"confirm,omitempty" jsonschema:"required true when the step's authored policy requires user confirmation"`
	Override         bool   `json:"override,omitempty" jsonschema:"bypass the step's completion policy explicitly"`
	ExpectedRevision uint64 `json:"expected_revision" jsonschema:"session revision this call expects, from session_get"`
	RequestID        string `json:"request_id,omitempty" jsonschema:"idempotency key for retries"`
}

type stepAdvanceArgs struct {
	SessionID        string `json:"session_id" jsonschema:"session ID returned by session_start"`
	Override         bool   `json:"override,omitempty" jsonschema:"activate the next node even if the current one is not completed"`
	ExpectedRevision uint64 `json:"expected_revision" jsonschema:"session revision this call expects, from session_get"`
	RequestID        string `json:"request_id,omitempty" jsonschema:"idempotency key for retries"`
}

func registerAssessmentTools(server *mcp.Server, sessions *application.SessionService) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "feedback_prepare",
		Description: "Assemble instruction, question, rubric references and scope for the caller to author feedback from. This server never drafts feedback itself.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args feedbackPrepareArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.SessionID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "session_id is required", false, nil), nil
		}
		packet, err := sessions.FeedbackPrepare(learning.SessionID(args.SessionID), args.Question)
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectNone, map[string]any{
			"objective": packet.Objective, "scope": packet.Scope, "question": packet.Question, "rubric_refs": packet.RubricRefs,
		})
		env.SessionID = args.SessionID
		env.ActiveNode = &ActiveNode{ID: packet.StepID}
		return nil, env, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "feedback_record",
		Description: "Record feedback already authored elsewhere, with its type. Never approves, evaluates, completes or advances.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, IdempotentHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args feedbackRecordArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.SessionID == "" || args.Text == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "session_id and text are required", false, nil), nil
		}
		result, err := sessions.FeedbackRecord(learning.SessionID(args.SessionID), learning.FeedbackType(args.Type), args.Text, args.Blocking, args.ExpectedRevision, args.RequestID)
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectFeedbackRecorded, map[string]any{"revision": result.Revision})
		env.SessionID = args.SessionID
		return nil, env, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "step_evaluate",
		Description: "Resolve criteria (deterministically for structural ones, from the caller's cited judgment for qualitative ones) and record the evaluation. Creates an attempt only when submission_intent is true. Never completes or advances.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, IdempotentHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args stepEvaluateArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.SessionID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "session_id is required", false, nil), nil
		}
		criteria := make([]application.CriterionInput, len(args.Criteria))
		for i, c := range args.Criteria {
			criteria[i] = application.CriterionInput{
				Name: c.Name, Kind: c.Kind, Severity: learning.FindingSeverity(c.Severity),
				Verdict: learning.EvaluationVerdict(c.Verdict), EvidenceID: learning.EvidenceID(c.EvidenceID), RubricRef: c.RubricRef,
			}
		}
		result, err := sessions.StepEvaluate(application.EvaluateInput{
			SessionID: learning.SessionID(args.SessionID), Criteria: criteria, SubmissionIntent: args.SubmissionIntent,
			ExpectedRevision: args.ExpectedRevision, RequestID: args.RequestID,
		})
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		effect := ProgressEffectEvaluationRecorded
		if result.AttemptRecorded {
			effect = ProgressEffectAttemptRecorded
		}
		env := okEnvelope(requestID, effect, map[string]any{
			"criteria": result.Criteria, "has_blocking_failure": result.HasBlockingFailure,
			"attempt_recorded": result.AttemptRecorded, "revision": result.Revision,
		})
		env.SessionID = args.SessionID
		env.ActiveNode = &ActiveNode{ID: string(result.StepID)}
		return nil, env, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "reflection_record",
		Description: "Record a reflection answer, its competency and a descriptive assessment, without altering code or step state.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, IdempotentHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args reflectionRecordArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.SessionID == "" || args.Answer == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "session_id and answer are required", false, nil), nil
		}
		result, err := sessions.ReflectionRecord(learning.SessionID(args.SessionID), learning.CompetencyID(args.CompetencyID), args.Prompt, args.Answer, args.Assessment, args.ExpectedRevision, args.RequestID)
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectNone, map[string]any{"revision": result.Revision})
		env.SessionID = args.SessionID
		return nil, env, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "step_complete",
		Description: "Mark the active step completed when its authored policy is satisfied, or override explicitly. Never activates the next node.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, IdempotentHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args stepCompleteArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.SessionID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "session_id is required", false, nil), nil
		}
		result, err := sessions.StepComplete(learning.SessionID(args.SessionID), args.Confirm, args.Override, args.ExpectedRevision, args.RequestID)
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectStepCompleted, map[string]any{"state": string(result.State), "revision": result.Revision})
		env.SessionID = args.SessionID
		return nil, env, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "step_advance",
		Description: "Activate the next permitted node in document order, or report branch options when the step splits, or that the challenge is exhausted.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, IdempotentHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args stepAdvanceArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.SessionID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "session_id is required", false, nil), nil
		}
		result, err := sessions.StepAdvance(learning.SessionID(args.SessionID), args.Override, args.ExpectedRevision, args.RequestID)
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		if len(result.Branches) > 0 {
			env := okEnvelope(requestID, ProgressEffectNone, map[string]any{"branches": result.Branches, "revision": result.Revision})
			env.SessionID = args.SessionID
			return nil, env, nil
		}
		if result.Done {
			env := okEnvelope(requestID, ProgressEffectNone, map[string]any{"done": true, "revision": result.Revision})
			env.SessionID = args.SessionID
			return nil, env, nil
		}
		env := okEnvelope(requestID, ProgressEffectStepAdvanced, map[string]any{"revision": result.Revision})
		env.SessionID = args.SessionID
		env.ActiveNode = &ActiveNode{ID: result.StepID, Kind: result.Kind}
		return nil, env, nil
	})
}
