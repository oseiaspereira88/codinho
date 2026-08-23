package recommendation

import (
	"slices"
	"sort"
)

// remediationThreshold is the mastery state below which the recommender
// looks for a smaller step before the objective's own best-matching
// candidate (Decision 4, requirement R6). demonstrates_without_help is the
// first state that proves the learner can do it unaided; anything short
// of that is treated as a real gap worth remediating.
const remediationThreshold = "demonstrates_without_help"

// Recommend ranks candidates toward objective and returns them ordered
// most-recommended first, tie-broken by ChallengeID for a stable result
// across identical inputs (requirement R8). completed marks challenge IDs
// the learner has already finished, so unmet Prerequisites can exclude a
// candidate (requirement R2's acyclic graph is a precondition this
// function assumes, not one it re-validates).
func Recommend(objective Objective, candidates []Candidate, mastery map[string]Mastery, completed map[string]bool) []Recommendation {
	var eligible []Candidate
	for _, c := range candidates {
		if prerequisitesMet(c, completed) {
			eligible = append(eligible, c)
		}
	}

	scored := make([]Recommendation, 0, len(eligible))
	for _, c := range eligible {
		scored = append(scored, score(c, objective, mastery, completed))
	}
	sortRecommendations(scored)

	return withRemediation(scored, eligible, objective, mastery, completed)
}

func prerequisitesMet(c Candidate, completed map[string]bool) bool {
	for _, p := range c.Prerequisites {
		if !completed[p] {
			return false
		}
	}
	return true
}

func matchesObjective(c Candidate, objective Objective) (competencyMatch, themeMatch bool) {
	competencyMatch = objective.CompetencyID != "" && slices.Contains(c.PrimaryCompetency, objective.CompetencyID)
	themeMatch = objective.ThemeID != "" && slices.Contains(c.ThemeIDs, objective.ThemeID)
	return
}

func score(c Candidate, objective Objective, mastery map[string]Mastery, completed map[string]bool) Recommendation {
	s := 0
	competencyMatch, themeMatch := matchesObjective(c, objective)
	switch {
	case competencyMatch:
		s += 100
	case themeMatch:
		s += 50
	}

	var best Mastery
	for _, comp := range c.PrimaryCompetency {
		if m, ok := mastery[comp]; ok && stateRank[m.BestState] > stateRank[best.BestState] {
			best = m
		}
		if m, ok := mastery[comp]; ok && m.ReviewOverdueBy > best.ReviewOverdueBy {
			best.ReviewOverdueBy = m.ReviewOverdueBy
		}
	}
	if best.ReviewOverdueBy > 0 {
		s += 30 // a due review takes priority over fresh material
	}
	if stateRank[best.BestState] >= stateRank["retained"] {
		s -= 50 // already well-mastered: deprioritize in favor of growth
	}

	if objective.TimeBudgetMinutes > 0 {
		if c.EstimatedMinutes > 0 && c.EstimatedMinutes <= objective.TimeBudgetMinutes {
			s += 10
		} else if c.EstimatedMinutes > objective.TimeBudgetMinutes {
			s -= 20
		}
	}

	return Recommendation{ChallengeID: c.ID, Score: s, Explanation: explain(c, best, completed, "")}
}

func sortRecommendations(recs []Recommendation) {
	sort.Slice(recs, func(i, j int) bool {
		if recs[i].Score != recs[j].Score {
			return recs[i].Score > recs[j].Score
		}
		return recs[i].ChallengeID < recs[j].ChallengeID
	})
}

// withRemediation prepends a smaller same-competency candidate ahead of
// the top recommendation when the objective competency's mastery is below
// remediationThreshold (Decision 4, requirement R6), then returns to the
// original top recommendation right after it.
func withRemediation(ranked []Recommendation, eligible []Candidate, objective Objective, mastery map[string]Mastery, completed map[string]bool) []Recommendation {
	if len(ranked) == 0 || objective.CompetencyID == "" {
		return ranked
	}
	m := mastery[objective.CompetencyID]
	if stateRank[m.BestState] >= stateRank[remediationThreshold] {
		return ranked
	}

	top := ranked[0]
	var topCandidate Candidate
	for _, c := range eligible {
		if c.ID == top.ChallengeID {
			topCandidate = c
			break
		}
	}

	var smaller *Candidate
	for _, c := range eligible {
		if c.ID == top.ChallengeID || !slices.Contains(c.PrimaryCompetency, objective.CompetencyID) {
			continue
		}
		if c.EstimatedMinutes == 0 || c.EstimatedMinutes >= topCandidate.EstimatedMinutes {
			continue
		}
		if smaller == nil || c.EstimatedMinutes < smaller.EstimatedMinutes || (c.EstimatedMinutes == smaller.EstimatedMinutes && c.ID < smaller.ID) {
			cc := c
			smaller = &cc
		}
	}
	if smaller == nil {
		return ranked
	}

	remediation := Recommendation{
		ChallengeID: smaller.ID,
		Score:       top.Score + 1, // ranks just ahead of the original
		Explanation: explain(*smaller, m, completed, top.ChallengeID),
	}
	out := make([]Recommendation, 0, len(ranked)+1)
	out = append(out, remediation)
	out = append(out, ranked...)
	return out
}
