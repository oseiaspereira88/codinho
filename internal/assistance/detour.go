package assistance

import "github.com/oseiaspereira88/codinho/internal/learning"

// DetourOutcome is how an open detour concludes.
type DetourOutcome string

const (
	DetourResolved  DetourOutcome = "resolved"
	DetourAbandoned DetourOutcome = "abandoned"
)

func (o DetourOutcome) state() learning.DetourState {
	if o == DetourAbandoned {
		return learning.DetourStateAbandoned
	}
	return learning.DetourStateResolved
}

// StartDetour opens sess's detour and immediately activates it: a single
// learning_detour_start call takes it straight to "active" because
// starting a detour means engaging with it right away (PROJECT.md §8.8).
// The session's active step is never touched (requirement R6).
func StartDetour(sess *learning.LearningSession, reason string) (*learning.LearningDetour, error) {
	d, err := sess.OpenDetour(reason)
	if err != nil {
		return nil, err
	}
	if err := d.Transition(learning.DetourStateActive); err != nil {
		return nil, err
	}
	return d, nil
}

// FinishDetour closes an active detour with outcome and returns it to the
// session's original active step without altering that step (requirement
// R6).
func FinishDetour(d *learning.LearningDetour, outcome DetourOutcome) error {
	if err := d.Transition(outcome.state()); err != nil {
		return err
	}
	return d.Transition(learning.DetourStateReturned)
}
