package session

import (
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/learning"
	"testing"
)

func TestTreeOrderedFrontierCrossesLayers(t *testing.T) {
	svc := newTestService(t)
	ch, _ := svc.catalog.Challenge(fixtureChallengeID)
	ch.Layers = append(ch.Layers, curriculum.LayerAuthoring{ID: "second", MacroSteps: []curriculum.StepAuthoring{{ID: "macro-2", Kind: "macro", Children: []curriculum.StepAuthoring{{ID: "micro-2", Kind: "micro"}}}}})
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID, Depth: learning.DepthMicro})
	if err != nil {
		t.Fatal(err)
	}
	// Pin the synthetic second layer into this isolated domain fixture.
	rec := svc.sessions[start.SessionID]
	rec.pinned = ch
	done, err := svc.StepComplete(start.SessionID, false, true, start.Revision, "")
	if err != nil {
		t.Fatal(err)
	}
	next, err := svc.StepAdvance(start.SessionID, false, done.Revision, "")
	if err != nil {
		t.Fatal(err)
	}
	if next.Done || next.StepID != "micro-2" {
		t.Fatalf("lost second layer: %+v", next)
	}
}
