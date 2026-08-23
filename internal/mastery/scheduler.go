package mastery

import "time"

// initialLadderDays are the initial spaced-review intervals (requirement
// R5). A successful review moves to the next rung; past the ladder's end,
// intervals keep growing by lastGrowthFactor.
var initialLadderDays = []int{1, 3, 7, 14, 30}

// lastGrowthFactor grows the interval once a competency's review history
// has exhausted the initial ladder, so long-mastered competencies are not
// reviewed forever at a fixed 30-day cadence.
const lastGrowthFactor = 1.5

// minIntervalDays bounds how short a review interval can shrink to after a
// failed review: failure never "zeroes" learning (PROJECT.md §21.4), it
// only brings the next review closer.
const minIntervalDays = 1

// shrinkFactor reduces the interval after a failed review, non-punitively:
// it halves rather than resetting to the ladder's start.
const shrinkFactor = 0.5

// NextSchedule computes the next review for a competency. current is nil
// the first time a competency earns a schedule at all (its very first
// piece of qualifying evidence); success reports whether the review (or,
// for a first schedule, the originating evidence) succeeded.
func NextSchedule(competencyID string, current *Schedule, success bool, now time.Time) Schedule {
	if current == nil {
		days := initialLadderDays[0]
		return Schedule{
			CompetencyID: competencyID,
			DueAt:        now.AddDate(0, 0, days),
			IntervalDays: days,
			Reason:       "primeira evidência qualificada; primeira revisão no rung inicial da ladder",
		}
	}

	var days int
	var reason string
	if success {
		days = nextRung(current.IntervalDays)
		reason = "revisão bem-sucedida; avança para o próximo intervalo"
	} else {
		days = shrink(current.IntervalDays)
		reason = "revisão malsucedida; intervalo reduzido, não zerado"
	}
	return Schedule{CompetencyID: competencyID, DueAt: now.AddDate(0, 0, days), IntervalDays: days, Reason: reason}
}

// nextRung returns the next interval after cur along the ladder, growing
// past its end once cur reaches or exceeds the last rung.
func nextRung(cur int) int {
	for _, rung := range initialLadderDays {
		if rung > cur {
			return rung
		}
	}
	grown := int(float64(cur) * lastGrowthFactor)
	if grown <= cur {
		grown = cur + 1
	}
	return grown
}

// shrink halves cur, never below minIntervalDays.
func shrink(cur int) int {
	reduced := int(float64(cur) * shrinkFactor)
	if reduced < minIntervalDays {
		return minIntervalDays
	}
	return reduced
}

// Overdue reports whether s is due for review at now.
func (s Schedule) Overdue(now time.Time) bool { return !s.DueAt.After(now) }

// OverdueBy returns how long s has been overdue at now (zero or negative
// when not yet due).
func (s Schedule) OverdueBy(now time.Time) time.Duration { return now.Sub(s.DueAt) }
