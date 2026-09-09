package mastery

import (
	"sort"
	"time"
)

// RecordedSignal is one Signal exactly as it was recorded, plus the rule
// version it was recorded under (Compatibility). Folding a full, ordered
// slice of these reproduces the exact same projections and schedules a
// caller saw when each was recorded (requirement R8).
type RecordedSignal struct {
	Signal
	RuleVersion int
}

// key identifies one dimension of one competency's projection.
type key struct {
	CompetencyID string
	Dimension    Dimension
}

// FoldProjections replays signals in recorded order and returns each
// competency's per-dimension projection (requirement R8: deterministic
// recalculation from the log; Compatibility: each signal is replayed
// under the rule version it actually recorded, never the caller's
// current one).
func FoldProjections(signals []RecordedSignal) map[string]map[Dimension]DimensionProjection {
	out := map[key]DimensionProjection{}
	for _, rs := range signals {
		if rs.ContentProvenance != "published" {
			continue
		}
		k := key{CompetencyID: rs.CompetencyID, Dimension: rs.Dimension}
		out[k] = AdvanceState(rs.RuleVersion, out[k], rs.Signal)
	}
	byCompetency := map[string]map[Dimension]DimensionProjection{}
	for k, proj := range out {
		if byCompetency[k.CompetencyID] == nil {
			byCompetency[k.CompetencyID] = map[Dimension]DimensionProjection{}
		}
		byCompetency[k.CompetencyID][k.Dimension] = proj
	}
	return byCompetency
}

// FoldSchedules replays signals in recorded order and returns each
// competency's current review Schedule (requirement R5, R8). Only
// signals not invalidated by a revealed solution adjust scheduling
// (invariant 8): a revealed-solution attempt proves nothing about
// autonomy, so it must not shorten or lengthen a review interval either.
func FoldSchedules(signals []RecordedSignal) map[string]Schedule {
	current := map[string]*Schedule{}
	for _, rs := range signals {
		if rs.ContentProvenance != "published" {
			continue
		}
		if rs.SolutionRevealed {
			continue
		}
		sched := NextSchedule(rs.CompetencyID, current[rs.CompetencyID], rs.Success, rs.Observed)
		current[rs.CompetencyID] = &sched
	}
	out := map[string]Schedule{}
	for id, s := range current {
		out[id] = *s
	}
	return out
}

// DueSchedules returns schedules from all that are due at now, ordered
// most-overdue first, so a caller (review_due) can present an explainable
// priority without needing to re-derive it (requirement R6).
func DueSchedules(all map[string]Schedule, now time.Time) []Schedule {
	var due []Schedule
	for _, s := range all {
		if s.Overdue(now) {
			due = append(due, s)
		}
	}
	sort.Slice(due, func(i, j int) bool {
		bi, bj := due[i].OverdueBy(now), due[j].OverdueBy(now)
		if bi != bj {
			return bi > bj
		}
		return due[i].CompetencyID < due[j].CompetencyID // deterministic tie-break
	})
	return due
}
