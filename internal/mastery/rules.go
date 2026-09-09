package mastery

import "time"

// AdvanceState applies sig to current under ruleVersion, returning the
// resulting projection. It never regresses and never promotes on a failed
// or solution-revealed signal (invariant 8); a help_used signal alone
// never reaches beyond demonstrates_with_help, however many times it
// repeats (requirement R4, invariant 10: no single guided completion, or
// any number of them, promotes past that ceiling); demonstrates_without_
// help only becomes retained given a later calendar date, and retained
// only becomes transferred given an unseen variant (Decision 4).
//
// ruleVersion is accepted (not yet branched on) so a future rule table can
// be added without changing every caller or silently reinterpreting
// history recorded under RuleVersionV1 (Compatibility).
func AdvanceState(ruleVersion int, current DimensionProjection, sig Signal) DimensionProjection {
	next := current
	next.RuleVersion = ruleVersion
	if current.SeenVariants == nil {
		next.SeenVariants = map[string]bool{}
	} else {
		next.SeenVariants = current.SeenVariants
	}
	if current.State == "" {
		next.State = StateNotObserved
	}

	if sig.ContentProvenance != "published" {
		return next
	}

	next.EvidenceCount++
	next.LastEvidenceID = sig.EvidenceID

	if !sig.Success || sig.SolutionRevealed {
		// First contact still counts as "introduced" even when it fails
		// or the solution was revealed: the competency is no longer
		// unobserved. Nothing else changes (no promotion, no regression).
		if rank[next.State] < rank[StateIntroduced] {
			next.State = StateIntroduced
		}
		return next
	}

	target := StateDemonstratesWithHelp
	if !sig.HelpUsed {
		target = StateDemonstratesWithoutHelp
		if rank[current.State] >= rank[StateDemonstratesWithoutHelp] && !sameCalendarDay(current.LastObserved, sig.Observed) {
			target = StateRetained
		}
		if rank[current.State] >= rank[StateRetained] && sig.Variant != "" && !current.SeenVariants[sig.Variant] {
			target = StateTransferred
		}
	}
	if rank[target] > rank[next.State] {
		next.State = target
	}
	next.LastObserved = sig.Observed
	if sig.Variant != "" {
		next.SeenVariants[sig.Variant] = true
	}
	return next
}

func sameCalendarDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}
