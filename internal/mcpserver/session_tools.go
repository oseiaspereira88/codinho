package mcpserver

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/oseiaspereira88/codinho/internal/application"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

type sessionStartArgs struct {
	TrackID          string   `json:"track_id,omitempty" jsonschema:"authored track ID, exclusive with challenge_id and composition_id"`
	CompositionID    string   `json:"composition_id,omitempty" jsonschema:"composition ID offered by learning_path_recommend, requires challenge_ids"`
	ChallengeIDs     []string `json:"challenge_ids,omitempty" jsonschema:"ordered IDs accepted with composition_id"`
	ChallengeID      string   `json:"challenge_id,omitempty" jsonschema:"ID of the challenge to fix for this session"`
	Mode             string   `json:"mode,omitempty" jsonschema:"pedagogical mode: teaching, practice, review, debug, exploration or interview (default practice)"`
	Depth            string   `json:"depth,omitempty" jsonschema:"initial depth: challenge, layer, macro, meso or micro (default micro)"`
	Help             string   `json:"help,omitempty" jsonschema:"help policy: free, progressive, limited, no_code, no_hints or only_after_attempt (default progressive)"`
	Evaluation       string   `json:"evaluation,omitempty" jsonschema:"evaluation policy: on_demand, on_step_complete or only_at_end (default on_demand)"`
	DisclosureMax    int      `json:"disclosure_max,omitempty" jsonschema:"assistance ladder ceiling, 0-6 (default 1)"`
	TimeLimitSeconds int      `json:"time_limit_seconds,omitempty" jsonschema:"optional session duration in seconds, off by default; interview-mode uses this for its optional timer"`
	RequestID        string   `json:"request_id,omitempty" jsonschema:"idempotency key; retrying with the same value returns the original session"`
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

type sessionFinishArgs struct {
	SessionID        string `json:"session_id" jsonschema:"session ID returned by session_start"`
	Reason           string `json:"reason,omitempty" jsonschema:"why the session is finishing, e.g. an explicit learner choice or timeout"`
	ExpectedRevision uint64 `json:"expected_revision" jsonschema:"session revision this call expects, from session_get"`
	RequestID        string `json:"request_id,omitempty" jsonschema:"idempotency key for retries"`
}

type interviewStatusArgs struct {
	SessionID string `json:"session_id" jsonschema:"session ID returned by session_start"`
}

type interviewReportArgs struct {
	SessionID   string   `json:"session_id" jsonschema:"session ID returned by session_start"`
	Recommended []string `json:"recommended,omitempty" jsonschema:"step or competency IDs worth revisiting, e.g. from learning_path_recommend"`
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
		Description: "Start an explicitly chosen challenge, authored track or accepted composition; pin its content and activate only the first instructional step.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, IdempotentHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args sessionStartArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.ChallengeID == "" && args.TrackID == "" && args.CompositionID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "challenge_id, track_id or composition_id is required", false, nil), nil
		}
		var timeLimit *time.Duration
		if args.TimeLimitSeconds > 0 {
			d := time.Duration(args.TimeLimitSeconds) * time.Second
			timeLimit = &d
		}
		result, err := sessions.Start(application.StartInput{
			TrackID: args.TrackID, CompositionID: args.CompositionID, ChallengeIDs: args.ChallengeIDs,
			ChallengeID:   args.ChallengeID,
			Mode:          learning.PedagogicalMode(args.Mode),
			Depth:         learning.Depth(args.Depth),
			Help:          learning.HelpPolicyKind(args.Help),
			Evaluation:    learning.EvaluationPolicyKind(args.Evaluation),
			DisclosureMax: learning.DisclosureLevel(args.DisclosureMax),
			TimeLimit:     timeLimit,
			RequestID:     args.RequestID,
		})
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectSessionChanged, map[string]any{
			"track":     result.Track,
			"objective": result.Objective,
			"revision":  result.Revision,
		})
		env.SessionID = string(result.SessionID)
		env.ActiveNode = &ActiveNode{ID: string(result.ActiveStep), Kind: result.Kind}
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
			"track":    result.Track,
			"state":    string(result.State),
			"revision": result.Revision,
		})
		env.SessionID = string(result.SessionID)
		env.Disclosure = toDisclosure(result.Disclosure)
		if result.ActiveStep != "" {
			env.ActiveNode = &ActiveNode{ID: string(result.ActiveStep), Kind: result.Kind}
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
		env.ActiveNode = &ActiveNode{ID: string(instruction.StepID), Kind: instruction.Kind}
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

	mcp.AddTool(server, &mcp.Tool{
		Name:        "session_finish",
		Description: "Finish a session without inferring completion from step state. reason records why (e.g. explicit learner choice, or timeout), for an interview report to read back.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, IdempotentHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args sessionFinishArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.SessionID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "session_id is required", false, nil), nil
		}
		result, err := sessions.FinishWithReason(learning.SessionID(args.SessionID), args.Reason, args.ExpectedRevision, args.RequestID)
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectSessionChanged, map[string]any{"state": string(result.State), "revision": result.Revision})
		env.SessionID = args.SessionID
		return nil, env, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "granularity_adjust",
		Description: "Show an ancestor or return to the current finer cursor, preserving each node's progress. Never advances to a sibling or rewrites the canonical tree.",
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

	mcp.AddTool(server, &mcp.Tool{
		Name:        "interview_status",
		Description: "Report elapsed time and whether the session's time limit (if any) has been reached, from its real recorded start time.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args interviewStatusArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.SessionID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "session_id is required", false, nil), nil
		}
		status, err := sessions.InterviewStatus(learning.SessionID(args.SessionID), time.Now().UTC())
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		data := map[string]any{"elapsed_seconds": status.Elapsed.Seconds(), "timed_out": status.TimedOut}
		if status.TimeLimit != nil {
			data["time_limit_seconds"] = status.TimeLimit.Seconds()
		}
		env := okEnvelope(requestID, ProgressEffectNone, data)
		env.SessionID = args.SessionID
		return nil, env, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "interview_report",
		Description: "Assemble the descriptive, non-scalar end-of-interview report from what was actually recorded: evaluations, granted-vs-blocked hints and reflections. Never reduces performance to a single score.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args interviewReportArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.SessionID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "session_id is required", false, nil), nil
		}
		report, err := sessions.InterviewReport(learning.SessionID(args.SessionID), time.Now().UTC(), args.Recommended)
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectNone, map[string]any{
			"challenge_id":    report.ChallengeID,
			"elapsed_seconds": report.Elapsed.Seconds(),
			"timed_out":       report.TimedOut,
			"finish_reason":   report.FinishReason,
			"evaluations":     report.Evaluations,
			"hints":           report.Hints,
			"reflections":     report.Reflections,
			"gaps":            report.Gaps,
			"recommended":     report.Recommended,
			"integrity_note":  report.IntegrityNote,
		})
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
