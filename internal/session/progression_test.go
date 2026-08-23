package session

import (
	"testing"

	"github.com/oseiaspereira88/codinho/internal/curriculum"
)

func linearTree() []curriculum.StepAuthoring {
	return []curriculum.StepAuthoring{
		{ID: "macro-1", Children: []curriculum.StepAuthoring{
			{ID: "meso-1", Children: []curriculum.StepAuthoring{
				{ID: "micro-1"},
			}},
		}},
		{ID: "macro-2"},
	}
}

func TestAdvanceFromDescendsIntoSingleChild(t *testing.T) {
	next, branches, found, ok := advanceFrom(linearTree(), "macro-1", nil)
	if !found || !ok || branches != nil || next.ID != "meso-1" {
		t.Fatalf("next=%+v branches=%v found=%v ok=%v", next, branches, found, ok)
	}
}

func TestAdvanceFromBubblesToNextSiblingAtLeaf(t *testing.T) {
	next, branches, found, ok := advanceFrom(linearTree(), "micro-1", nil)
	if !found || !ok || branches != nil || next.ID != "macro-2" {
		t.Fatalf("next=%+v branches=%v found=%v ok=%v", next, branches, found, ok)
	}
}

func TestAdvanceFromReportsExhaustionAtLastStep(t *testing.T) {
	_, branches, found, ok := advanceFrom(linearTree(), "macro-2", nil)
	if !found || ok || branches != nil {
		t.Fatalf("expected found=true ok=false branches=nil at the last step, got found=%v ok=%v branches=%v", found, ok, branches)
	}
}

func TestAdvanceFromReportsBranchesWhenMultipleChildren(t *testing.T) {
	tree := []curriculum.StepAuthoring{
		{ID: "macro-1", Children: []curriculum.StepAuthoring{
			{ID: "meso-a"},
			{ID: "meso-b"},
		}},
	}
	next, branches, found, ok := advanceFrom(tree, "macro-1", nil)
	if !found || ok || next.ID != "" {
		t.Fatalf("next=%+v found=%v ok=%v", next, found, ok)
	}
	if len(branches) != 2 || branches[0].ID != "meso-a" || branches[1].ID != "meso-b" {
		t.Fatalf("unexpected branches: %+v", branches)
	}
}

func TestAdvanceFromUnknownIDNotFound(t *testing.T) {
	_, _, found, _ := advanceFrom(linearTree(), "does-not-exist", nil)
	if found {
		t.Fatal("expected found=false for an unknown step ID")
	}
}
