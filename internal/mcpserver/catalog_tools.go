package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/oseiaspereira88/codinho/internal/application"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
)

// catalogSearchArgs is catalog_search's strict input schema (requirement
// R3, R5, R6).
type catalogSearchArgs struct {
	ThemeIDs      []string `json:"theme_ids,omitempty" jsonschema:"explicit set of 1 to 16 theme IDs; exclusive with the legacy theme field"`
	CompetencyIDs []string `json:"competency_ids,omitempty" jsonschema:"explicit set of 1 to 16 competency IDs; exclusive with the legacy competency field"`
	Kind          string   `json:"kind,omitempty" jsonschema:"item kind to list: theme, concept, competency, track or challenge (default challenge)"`
	Theme         string   `json:"theme,omitempty" jsonschema:"optional theme ID to narrow results"`
	Text          string   `json:"text,omitempty" jsonschema:"free text matched against title (and brief, for challenges); bounded length"`
	Competency    string   `json:"competency,omitempty" jsonschema:"challenge-only: matches a primary or secondary competency ID"`
	Difficulty    string   `json:"difficulty,omitempty" jsonschema:"challenge-only: foundational, intermediate, advanced, ..."`
	ChallengeKind string   `json:"challenge_kind,omitempty" jsonschema:"challenge-only: the authored kind (atomic, ...)"`
	MaxMinutes    int      `json:"max_minutes,omitempty" jsonschema:"challenge-only: exclude challenges estimated longer than this"`
	Prerequisite  string   `json:"prerequisite,omitempty" jsonschema:"challenge-only: matches a challenge that lists this ID as a prerequisite"`
}

// catalogGetArgs is catalog_get's strict input schema.
type catalogGetArgs struct {
	ID string `json:"id" jsonschema:"catalog item ID"`
}

// conceptRelationsGetArgs is concept_relations_get's strict input schema
// (curriculum-graph-path-recommendation, requirement R1, R7). Despite the
// name (PROJECT.md §15.10), it resolves relations for any catalog item ID,
// not concepts alone.
type conceptRelationsGetArgs struct {
	ID string `json:"id" jsonschema:"catalog item ID (concept, competency, challenge, theme or track)"`
}

func registerCatalogTools(server *mcp.Server, catalog *application.CatalogService) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "catalog_search",
		Description: "Search themes, concepts, competencies, tracks and challenges by kind, theme, text, competency, difficulty, type, duration and prerequisite.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args catalogSearchArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		themes, err := subjectIDs(args.Theme, args.ThemeIDs)
		if err != nil {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "invalid or conflicting theme selectors", false, nil), nil
		}
		comps, err := subjectIDs(args.Competency, args.CompetencyIDs)
		if err != nil {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "invalid or conflicting competency selectors", false, nil), nil
		}

		result, err := catalog.Search(application.SearchQuery{
			Kind: curriculum.ItemKind(args.Kind), ThemeIDs: themes, Text: args.Text, CompetencyIDs: comps,
			Difficulty: args.Difficulty, ChallengeKind: args.ChallengeKind, MaxMinutes: args.MaxMinutes,
			Prerequisite: args.Prerequisite,
		})
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectNone, result.Items)
		env.Selection = &result.Selection
		if result.Truncated {
			env.Warnings = []string{"result truncated: refine the query to see every match"}
		}
		return nil, env, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "catalog_get",
		Description: "Get one catalog item by ID and version, respecting the session's disclosure level.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args catalogGetArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.ID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "id is required", false, nil), nil
		}
		item, err := catalog.Get(args.ID)
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		return nil, okEnvelope(requestID, ProgressEffectNone, item), nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "concept_relations_get",
		Description: "Get every authored relation to and from one catalog item (requires, recommended_before, relates_to, contrasts_with, commonly_fails_with, applies_in, deepens_into, evidences).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args conceptRelationsGetArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.ID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "id is required", false, nil), nil
		}
		out, in, err := catalog.Relations(args.ID)
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectNone, map[string]any{"outgoing": out, "incoming": in})
		return nil, env, nil
	})
}
