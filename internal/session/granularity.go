package session

import (
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

// Window identifies one instruction; the canonical tree is never rewritten.
type Window struct{ StepID, Kind string }

func challengeTree(ch curriculum.ChallengeAuthoring) curriculum.StepAuthoring {
	root := curriculum.StepAuthoring{ID: ch.ID, Kind: "challenge", Instruction: curriculum.InstructionAuthoring{Objective: ch.Title, Scope: ch.Brief}}
	for _, l := range ch.Layers {
		root.Children = append(root.Children, curriculum.StepAuthoring{ID: l.ID, Kind: "layer", Instruction: curriculum.InstructionAuthoring{Objective: l.ID, Scope: ch.Title}, Children: l.MacroSteps})
	}
	return root
}

func depthRank(kind string) int {
	switch kind {
	case "challenge":
		return 0
	case "layer":
		return 1
	case "macro":
		return 2
	case "meso":
		return 3
	case "micro":
		return 4
	}
	return -1
}

func deriveWindow(ch curriculum.ChallengeAuthoring, depth learning.Depth) (Window, bool) {
	if depthRank(string(depth)) < 0 {
		return Window{}, false
	}
	n := challengeTree(ch)
	for depthRank(n.Kind) < depthRank(string(depth)) && len(n.Children) > 0 && n.ChildrenMode != "choice" {
		n = n.Children[0]
	}
	return Window{n.ID, n.Kind}, n.ID != ""
}

func firstStep(ch curriculum.ChallengeAuthoring) (curriculum.StepAuthoring, bool) {
	for _, l := range ch.Layers {
		if len(l.MacroSteps) > 0 {
			return l.MacroSteps[0], true
		}
	}
	return curriculum.StepAuthoring{}, false
}

func findStep(ch curriculum.ChallengeAuthoring, id string) (curriculum.StepAuthoring, bool) {
	return searchSteps([]curriculum.StepAuthoring{challengeTree(ch)}, id)
}
func searchSteps(steps []curriculum.StepAuthoring, id string) (curriculum.StepAuthoring, bool) {
	for _, s := range steps {
		if s.ID == id {
			return s, true
		}
		if n, ok := searchSteps(s.Children, id); ok {
			return n, true
		}
	}
	return curriculum.StepAuthoring{}, false
}

func nodePath(node curriculum.StepAuthoring, id string) []curriculum.StepAuthoring {
	if node.ID == id {
		return []curriculum.StepAuthoring{node}
	}
	for _, c := range node.Children {
		if p := nodePath(c, id); len(p) > 0 {
			return append([]curriculum.StepAuthoring{node}, p...)
		}
	}
	return nil
}

func (r *record) windowAt(depth learning.Depth) (Window, error) {
	if depthRank(string(depth)) < 0 {
		return Window{}, ErrNoWindowAtDepth
	}
	active := r.session.ActiveStep()
	path := nodePath(challengeTree(r.pinned), r.cursor)
	if len(path) == 0 && active != nil {
		path = nodePath(challengeTree(r.pinned), string(active.StepID))
	}
	if len(path) == 0 {
		return Window{}, ErrNoActiveStep
	}
	for _, n := range path {
		if r.covered[n.ID] && depthRank(string(depth)) > depthRank(n.Kind) {
			return Window{}, ErrNoWindowAtDepth
		}
		if depthRank(n.Kind) >= depthRank(string(depth)) {
			return Window{n.ID, n.Kind}, nil
		}
	}
	n := path[len(path)-1]
	// Re-exposing an already completed cursor is not an implicit advance.
	if !r.completed(n.ID) {
		next := r.frontier(n, depth)
		if next.node.ID != "" {
			n = next.node
		} else if next.parent.ID != "" {
			n = next.parent
		}
	}
	return Window{n.ID, n.Kind}, nil
}
