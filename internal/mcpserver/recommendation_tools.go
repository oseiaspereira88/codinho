package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/oseiaspereira88/codinho/internal/application"
)

type learningPathRecommendArgs struct {
	ThemeIDs          []string `json:"theme_ids,omitempty" jsonschema:"explicit set of 1 to 16 theme IDs; exclusive with the legacy theme field"`
	CompetencyIDs     []string `json:"competency_ids,omitempty" jsonschema:"explicit set of 1 to 16 competency IDs; exclusive with the legacy competency field"`
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
		themes, err := subjectIDs(args.ThemeID, args.ThemeIDs)
		if err != nil {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "invalid or conflicting theme selectors", false, nil), nil
		}
		comps, err := subjectIDs(args.CompetencyID, args.CompetencyIDs)
		if err != nil {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "invalid or conflicting competency selectors", false, nil), nil
		}

		result, err := recs.RecommendSelection(application.RecommendInput{
			CompetencyIDs: comps, ThemeIDs: themes,
			TimeBudgetMinutes: args.TimeBudgetMinutes, Completed: args.Completed,
		})
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectNone, result)
		return nil, env, nil
	})
}
