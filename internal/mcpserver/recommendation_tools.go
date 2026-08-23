package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/oseiaspereira88/codinho/internal/application"
)

type learningPathRecommendArgs struct {
	CompetencyID      string   `json:"competency_id,omitempty" jsonschema:"target competency ID"`
	ThemeID           string   `json:"theme_id,omitempty" jsonschema:"target theme ID, used alongside or instead of competency_id"`
	TimeBudgetMinutes int      `json:"time_budget_minutes,omitempty" jsonschema:"available practice time in minutes; 0 means unconstrained"`
	Completed         []string `json:"completed,omitempty" jsonschema:"challenge IDs the learner has already finished, so their dependents become eligible"`
}

func registerRecommendationTools(server *mcp.Server, recs *application.RecommendationService) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "learning_path_recommend",
		Description: "Recommend a consultative, explainable path of challenges toward an objective (competency and/or theme, time budget, progress and overdue reviews). Never starts a track on its own.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args learningPathRecommendArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		result, err := recs.Recommend(application.RecommendInput{
			CompetencyID: args.CompetencyID, ThemeID: args.ThemeID,
			TimeBudgetMinutes: args.TimeBudgetMinutes, Completed: args.Completed,
		})
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectNone, map[string]any{"recommendations": result})
		return nil, env, nil
	})
}
