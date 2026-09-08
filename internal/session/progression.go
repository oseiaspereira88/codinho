package session

import (
	"errors"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

var ErrInvalidNextStep = errors.New("session: next step is not an offered option")
var ErrNavigationUnavailable = errors.New("session: navigation requires an active session without an open detour")

type savedProgress struct {
	Step  learning.StepProgress
	Clean bool
}

func (r *record) clone() *record {
	c := *r
	c.session = r.session.Clone()
	c.progress = map[string]savedProgress{}
	for k, v := range r.progress {
		c.progress[k] = v
	}
	c.covered = map[string]bool{}
	for k, v := range r.covered {
		c.covered[k] = v
	}
	c.choices = map[string]string{}
	for k, v := range r.choices {
		c.choices[k] = v
	}
	return &c
}
func (r *record) saveWindow() {
	if r.progress == nil {
		r.progress = map[string]savedProgress{}
	}
	if p := r.session.ActiveStep(); p != nil {
		r.progress[string(p.StepID)] = savedProgress{*p, r.cleanEvaluation}
	}
}
func (r *record) completed(id string) bool {
	if p := r.session.ActiveStep(); p != nil && string(p.StepID) == id {
		return p.State == learning.StepStateCompleted || p.State == learning.StepStateSkipped
	}
	p := r.progress[id].Step
	return p.State == learning.StepStateCompleted || p.State == learning.StepStateSkipped
}
func (r *record) activate(id string) error {
	r.saveWindow()
	saved, ok := r.progress[id]
	if !ok {
		p, err := activeProgress(id)
		if err != nil {
			return err
		}
		saved.Step = *p
	}
	r.session.ClearActiveInstruction()
	r.cleanEvaluation = saved.Clean
	return r.session.SetActiveInstruction(&saved.Step)
}
func (r *record) navigationAllowed() error {
	if r.session.State != learning.SessionStateActive {
		return ErrNavigationUnavailable
	}
	if d := r.session.Detour(); d != nil && d.State != learning.DetourStateReturned {
		return ErrNavigationUnavailable
	}
	return nil
}

type navigation struct {
	node, parent curriculum.StepAuthoring
	options      []curriculum.StepAuthoring
}

// frontier resolves the first unfinished window in document order. Containers
// do not acquire synthetic completion events when their finer windows finish.
func (r *record) frontier(n curriculum.StepAuthoring, depth learning.Depth) navigation {
	if r.covered[n.ID] {
		return navigation{}
	}
	if depthRank(n.Kind) >= depthRank(string(depth)) || len(n.Children) == 0 {
		if r.completed(n.ID) {
			return navigation{}
		}
		return navigation{node: n}
	}
	if n.ChildrenMode == "choice" {
		chosen := r.choices[n.ID]
		if chosen == "" {
			if !r.completed(n.ID) {
				return navigation{node: n}
			}
			return navigation{parent: n, options: n.Children}
		}
		for _, c := range n.Children {
			if c.ID == chosen {
				return r.frontier(c, depth)
			}
		}
		return navigation{}
	}
	for _, c := range n.Children {
		if next := r.frontier(c, depth); next.node.ID != "" || len(next.options) > 0 {
			return next
		}
	}
	return navigation{}
}
func (r *record) coversWindow(n curriculum.StepAuthoring) bool {
	return depthRank(n.Kind) >= depthRank(string(r.depth)) || len(n.Children) == 0
}

// nextNavigation validates completion before returning even options/exhaustion.
// The returned record is isolated; a rejected append cannot mutate live state.
func (r *record) nextNavigation(override bool, selected string) (navigation, *record, error) {
	if err := r.navigationAllowed(); err != nil {
		return navigation{}, nil, err
	}
	active := r.session.ActiveStep()
	if active == nil {
		return navigation{}, nil, ErrNoActiveStep
	}
	if active.State != learning.StepStateCompleted && !override {
		return navigation{}, nil, learning.DomainError{Code: learning.ErrCodeAdvanceNotAllowed}
	}
	c := r.clone()
	c.saveWindow()
	if override && active.State != learning.StepStateCompleted {
		p := c.session.ActiveStep()
		p.State = learning.StepStateSkipped
		c.saveWindow()
		n, _ := findStep(c.pinned, string(p.StepID))
		if c.coversWindow(n) {
			c.covered[n.ID] = true
		}
	}
	next := c.frontier(challengeTree(c.pinned), c.depth)
	if selected != "" {
		valid := false
		for _, o := range next.options {
			if o.ID == selected {
				valid = true
				break
			}
		}
		if !valid {
			return navigation{}, nil, ErrInvalidNextStep
		}
		c.choices[next.parent.ID] = selected
		next = c.frontier(challengeTree(c.pinned), c.depth)
	}
	return next, c, nil
}

func (s *Service) checkNavigationRevision(id learning.SessionID, expected uint64) error {
	if s.store.Revision(string(id)) != expected {
		return eventstore.ErrRevisionConflict
	}
	return nil
}
