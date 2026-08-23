package session

import "github.com/oseiaspereira88/codinho/internal/curriculum"

// rootSteps returns the macro-step tree StepAdvance walks: the first
// layer with any macro steps, mirroring firstStep/deriveWindow's existing
// traversal convention (granularity.go) rather than combining every layer
// into one path.
func rootSteps(ch curriculum.ChallengeAuthoring) ([]curriculum.StepAuthoring, bool) {
	for _, layer := range ch.Layers {
		if len(layer.MacroSteps) > 0 {
			return layer.MacroSteps, true
		}
	}
	return nil, false
}

// advanceFrom walks steps in pre-order (document) traversal looking for
// currentID, then reports what comes immediately after it:
//   - a single next step when current has zero or exactly one child (in
//     which case next is that child), or no children and exactly one
//     candidate after it in the tree;
//   - a list of branch options when current has more than one child
//     (PROJECT.md §15.6 "informa opções quando houver ramificação"),
//     leaving the caller to choose rather than guessing;
//   - found=true, ok=false, no branches when current was the last step in
//     the tree (the challenge is exhausted for this traversal).
//
// siblingsAfter carries what follows current's own subtree at each
// recursion level, so reaching a childless leaf can still resolve to its
// nearest unvisited ancestor sibling.
func advanceFrom(steps []curriculum.StepAuthoring, currentID string, siblingsAfter []curriculum.StepAuthoring) (next curriculum.StepAuthoring, branches []curriculum.StepAuthoring, found, ok bool) {
	for i, s := range steps {
		if s.ID == currentID {
			switch {
			case len(s.Children) == 1:
				return s.Children[0], nil, true, true
			case len(s.Children) > 1:
				return curriculum.StepAuthoring{}, s.Children, true, false
			default:
				rest := append(append([]curriculum.StepAuthoring{}, steps[i+1:]...), siblingsAfter...)
				if len(rest) == 0 {
					return curriculum.StepAuthoring{}, nil, true, false
				}
				return rest[0], nil, true, true
			}
		}
		childRest := append(append([]curriculum.StepAuthoring{}, steps[i+1:]...), siblingsAfter...)
		if n, b, f, o := advanceFrom(s.Children, currentID, childRest); f {
			return n, b, f, o
		}
	}
	return curriculum.StepAuthoring{}, nil, false, false
}
