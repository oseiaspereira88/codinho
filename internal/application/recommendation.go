package application

import (
	"time"

	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/mastery"
	"github.com/oseiaspereira88/codinho/internal/recommendation"
)

// RecommendationService combines the curriculum graph with mastery-review-
// scheduling progress into an explainable, deterministic recommendation
// (curriculum-graph-path-recommendation, requirement R4).
type RecommendationService struct {
	catalog  *curriculum.Catalog
	progress *ProgressService
}

// NewRecommendationService wires a RecommendationService to the catalog and
// progress service it reads from.
func NewRecommendationService(catalog *curriculum.Catalog, progress *ProgressService) *RecommendationService {
	return &RecommendationService{catalog: catalog, progress: progress}
}

// RecommendInput is one learning_path_recommend call. Completed lists
// challenge IDs the learner has already finished. There is no separate
// global registry of challenge completions in this system — session
// completion is session-scoped and ephemeral (session-orchestration-
// disclosure), and mastery-review-scheduling tracks competency-level
// domain, not challenge-level completion — so the caller (the tutor, which
// observed the session) supplies Completed directly, the same boundary
// already established for mastery_evidence_record's cited evidence
// (Decision 5).
type RecommendInput struct {
	CompetencyID      string
	ThemeID           string
	TimeBudgetMinutes int
	Completed         []string
}

// Recommend ranks candidate challenges toward in's objective, folding in
// the learner's current mastery projection and overdue reviews when a
// ProgressService is wired (requirement R4).
func (r *RecommendationService) Recommend(in RecommendInput) ([]recommendation.Recommendation, error) {
	challenges := r.catalog.Challenges()
	challengeIDs := make(map[string]bool, len(challenges))
	for _, ch := range challenges {
		challengeIDs[ch.ID] = true
	}

	candidates := make([]recommendation.Candidate, 0, len(challenges))
	for _, ch := range challenges {
		// Prerequisites may name a concept or a challenge (catalog-schema-
		// loader). Only a challenge prerequisite gates recommendation here:
		// a concept prerequisite means "understand this first", which this
		// spec has no separate per-concept completion signal for — mastery
		// tracks competencies, not concepts (mastery-review-scheduling).
		var gating []string
		for _, p := range ch.Prerequisites {
			if challengeIDs[p] {
				gating = append(gating, p)
			}
		}
		candidates = append(candidates, recommendation.Candidate{
			ID: ch.ID, Title: ch.Title, ThemeIDs: ch.Themes,
			PrimaryCompetency: ch.Competencies.Primary, Difficulty: ch.Difficulty,
			EstimatedMinutes: ch.EstimatedMinutes, Prerequisites: gating,
		})
	}

	completed := make(map[string]bool, len(in.Completed))
	for _, id := range in.Completed {
		completed[id] = true
	}

	masteryByCompetency, err := r.masteryInputs()
	if err != nil {
		return nil, err
	}

	objective := recommendation.Objective{CompetencyID: in.CompetencyID, ThemeID: in.ThemeID, TimeBudgetMinutes: in.TimeBudgetMinutes}
	return recommendation.Recommend(objective, candidates, masteryByCompetency, completed), nil
}

// masteryInputs reduces ProgressService's full per-dimension projection and
// due-review list into the one best state and overdue-by-days figure per
// competency that internal/recommendation needs. Returns an empty map
// (never an error) when no ProgressService is wired.
func (r *RecommendationService) masteryInputs() (map[string]recommendation.Mastery, error) {
	out := map[string]recommendation.Mastery{}
	if r.progress == nil {
		return out, nil
	}

	progress, err := r.progress.Progress("")
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	due, err := r.progress.ReviewDue(now)
	if err != nil {
		return nil, err
	}
	overdueDays := make(map[string]int, len(due.Due))
	for _, s := range due.Due {
		overdueDays[s.CompetencyID] = max(int(s.OverdueBy(now).Hours()/24), 1)
	}

	for competencyID, dims := range progress.Competencies {
		best := mastery.StateNotObserved
		for _, proj := range dims {
			best = mastery.Higher(best, proj.State)
		}
		out[competencyID] = recommendation.Mastery{
			CompetencyID: competencyID, BestState: string(best), ReviewOverdueBy: overdueDays[competencyID],
		}
	}
	return out, nil
}
