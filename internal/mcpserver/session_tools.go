package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/oseiaspereira88/codinho/internal/application"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

type sessionStartArgs struct {
	ChallengeID   string `json:"challenge_id" jsonschema:"ID of the challenge to fix for this session"`
	Mode          string `json:"mode,omitempty" jsonschema:"pedagogical mode: teaching, practice, review, debug, exploration or interview (default practice)"`
	Depth         string `json:"depth,omitempty" jsonschema:"initial depth: challenge, layer, macro, meso or micro (default micro)"`
	Help          string `json:"help,omitempty" jsonschema:"help policy: free, progressive, limited, no_code, no_hints or only_after_attempt (default progressive)"`
	Evaluation    string `json:"evaluation,omitempty" jsonschema:"evaluation policy: on_demand, on_step_complete or only_at_end (default on_demand)"`
	DisclosureMax int    `json:"disclosure_max,omitempty" jsonschema:"assistance ladder ceiling, 0-6 (default 1)"`
	RequestID     string `json:"request_id,omitempty" jsonschema:"idempotency key; retrying with the same value returns the original session"`
}

type sessionGetArgs struct {
	SessionID string `json:"session_id" jsonschema:"session ID returned by session_start"`
}

type instructionGetArgs struct {
	SessionID string `json:"session_id" jsonschema:"session ID returned by session_start"`
}

type sessionConfigureArgs struct {
	SessionID        string `json:"session_id" jsonschema:"session ID returned by session_start"`
	Help             string `json:"help,omitempty" jsonschema:"new help policy: free, progressive, limited, no_code, no_hints or only_after_attempt"`
	Evaluation       string `json:"evaluation,omitempty" jsonschema:"new evaluation policy: on_demand, on_step_complete or only_at_end"`
	ExpectedRevision uint64 `json:"expected_revision" jsonschema:"session revision this call expects, from session_get"`
	RequestID        string `json:"request_id,omitempty" jsonschema:"idempotency key for retries"`
}

type sessionLifecycleArgs struct {
	SessionID        string `json:"session_id" jsonschema:"session ID returned by session_start"`
	ExpectedRevision uint64 `json:"expected_revision" jsonschema:"session revision this call expects, from session_get"`
	RequestID        string `json:"request_id,omitempty" jsonschema:"idempotency key for retries"`
}

type granularityAdjustArgs struct {
	SessionID        string `json:"session_id" jsonschema:"session ID returned by session_start"`
	Depth            string `json:"depth" jsonschema:"target depth: challenge, layer, macro, meso or micro"`
	Reason           string `json:"reason,omitempty" jsonschema:"why this change is happening, so the student can be told (PROJECT.md §8.5: every granularity change must be explained)"`
	ExpectedRevision uint64 `json:"expected_revision" jsonschema:"session revision this call expects, from session_get"`
	RequestID        string `json:"request_id,omitempty" jsonschema:"idempotency key for retries"`
}

type learnerNextStepProposeArgs struct {
	SessionID        string `json:"session_id" jsonschema:"session ID returned by session_start"`
	StepID           string `json:"step_id" jsonschema:"the step ID the learner proposes as what comes next"`
	ExpectedRevision uint64 `json:"expected_revision" jsonschema:"session revision this call expects, from session_get"`
	RequestID        string `json:"request_id,omitempty" jsonschema:"idempotency key for retries"`
}

func toDisclosure(d application.Disclosure) *Disclosure {
	return &Disclosure{Level: d.Level, SolutionRevealed: d.SolutionRevealed}
}

func registerSessionTools(server *mcp.Server, sessions *application.SessionService) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "session_start",
		Description: "Fix a challenge, mode and policies, and activate the session's first instructional step.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, IdempotentHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args sessionStartArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.ChallengeID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "challenge_id is required", false, nil), nil
		}
		result, err := sessions.Start(application.StartInput{
			ChallengeID:   args.ChallengeID,
			Mode:          learning.PedagogicalMode(args.Mode),
			Depth:         learning.Depth(args.Depth),
			Help:          learning.HelpPolicyKind(args.Help),
			Evaluation:    learning.EvaluationPolicyKind(args.Evaluation),
			DisclosureMax: learning.DisclosureLevel(args.DisclosureMax),
			RequestID:     args.RequestID,
		})
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectSessionChanged, map[string]any{
			"objective": result.Objective,
			"revision":  result.Revision,
		})
		env.SessionID = string(result.SessionID)
		env.ActiveNode = &ActiveNode{ID: string(result.ActiveStep), Kind: "micro"}
		env.Disclosure = toDisclosure(result.Disclosure)
		env.AllowedActions = []string{"instruction_get", "session_get", "session_configure", "session_pause", "session_finish", "granularity_adjust", "learner_next_step_propose"}
		return nil, env, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "session_get",
		Description: "Return a session's state, active node, revision and disclosure.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args sessionGetArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.SessionID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "session_id is required", false, nil), nil
		}
		result, err := sessions.Get(learning.SessionID(args.SessionID))
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectNone, map[string]any{
			"state":    string(result.State),
			"revision": result.Revision,
		})
		env.SessionID = string(result.SessionID)
		env.Disclosure = toDisclosure(result.Disclosure)
		if result.ActiveStep != "" {
			env.ActiveNode = &ActiveNode{ID: string(result.ActiveStep), Kind: "micro"}
		}
		return nil, env, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "instruction_get",
		Description: "Deliver a single instruction at the authorized depth. Never returns future children or the solution.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args instructionGetArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.SessionID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "session_id is required", false, nil), nil
		}
		instruction, err := sessions.Instruction(learning.SessionID(args.SessionID))
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectNone, map[string]string{
			"objective": instruction.Objective,
			"scope":     instruction.Scope,
		})
		env.SessionID = args.SessionID
		env.ActiveNode = &ActiveNode{ID: string(instruction.StepID), Kind: "micro"}
		env.Disclosure = toDisclosure(instruction.Disclosure)
		return nil, env, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "session_configure",
		Description: "Change only mutable session properties: help policy, evaluation policy. Rejects a stale expected_revision.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, IdempotentHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args sessionConfigureArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.SessionID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "session_id is required", false, nil), nil
		}
		in := application.ConfigureInput{
			SessionID:        learning.SessionID(args.SessionID),
			ExpectedRevision: args.ExpectedRevision,
			RequestID:        args.RequestID,
		}
		if args.Help != "" {
			h := learning.HelpPolicyKind(args.Help)
			in.Help = &h
		}
		if args.Evaluation != "" {
			e := learning.EvaluationPolicyKind(args.Evaluation)
			in.Evaluation = &e
		}
		result, err := sessions.Configure(in)
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectSessionChanged, map[string]any{"state": string(result.State), "revision": result.Revision})
		env.SessionID = args.SessionID
		return nil, env, nil
	})

	registerLifecycleTool(server, "session_pause", "Pause an active session without inferring anything about step completion.", sessions.Pause)
	registerLifecycleTool(server, "session_resume", "Resume a paused session.", sessions.Resume)
	registerLifecycleTool(server, "session_finish", "Finish a session without inferring completion from step state.", sessions.Finish)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "granularity_adjust",
		Description: "Move the session's instructional window to a new depth by walking the canonical step tree, without rewriting it.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, IdempotentHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args granularityAdjustArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.SessionID == "" || args.Depth == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "session_id and depth are required", false, nil), nil
		}
		result, err := sessions.GranularityAdjust(learning.SessionID(args.SessionID), learning.Depth(args.Depth), args.Reason, args.ExpectedRevision, args.RequestID)
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectNone, map[string]any{"revision": result.Revision})
		env.SessionID = args.SessionID
		env.ActiveNode = &ActiveNode{ID: result.StepID, Kind: result.Kind}
		return nil, env, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "learner_next_step_propose",
		Description: "Record that the learner proposed a step as what comes next: a pure autonomy signal that never advances anything on its own.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, IdempotentHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args learnerNextStepProposeArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.SessionID == "" || args.StepID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "session_id and step_id are required", false, nil), nil
		}
		result, err := sessions.ProposeNextStep(learning.SessionID(args.SessionID), learning.StepID(args.StepID), args.ExpectedRevision, args.RequestID)
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectNone, map[string]any{"revision": result.Revision})
		env.SessionID = args.SessionID
		return nil, env, nil
	})
}

// lifecycleFunc is the shape shared by SessionService.Pause, Resume and
// Finish, letting registerLifecycleTool avoid repeating three near-identical
// tool registrations.
type lifecycleFunc func(id learning.SessionID, expectedRevision uint64, requestID string) (application.LifecycleResult, error)

func registerLifecycleTool(server *mcp.Server, name, description string, fn lifecycleFunc) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        name,
		Description: description,
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, IdempotentHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args sessionLifecycleArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.SessionID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "session_id is required", false, nil), nil
		}
		result, err := fn(learning.SessionID(args.SessionID), args.ExpectedRevision, args.RequestID)
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectSessionChanged, map[string]any{"state": string(result.State), "revision": result.Revision})
		env.SessionID = args.SessionID
		return nil, env, nil
	})
}
