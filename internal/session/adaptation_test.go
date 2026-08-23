package session

import (
	"testing"

	"github.com/oseiaspereira88/codinho/internal/learning"
	"github.com/oseiaspereira88/codinho/internal/mastery"
)

func TestSuggestGranularityRequiresEvidenceThreshold(t *testing.T) {
	proj := mastery.DimensionProjection{State: mastery.StateDemonstratesWithoutHelp, EvidenceCount: EvidenceThreshold - 1}
	_, reason, ok := SuggestGranularity(learning.DepthMicro, proj)
	if ok {
		t.Fatalf("expected no suggestion below threshold, got reason %q", reason)
	}
}

func TestSuggestGranularityWidensOnRepeatedEase(t *testing.T) {
	proj := mastery.DimensionProjection{State: mastery.StateDemonstratesWithoutHelp, EvidenceCount: EvidenceThreshold}
	depth, reason, ok := SuggestGranularity(learning.DepthMicro, proj)
	if !ok {
		t.Fatal("expected a suggestion at threshold")
	}
	if depth != learning.DepthMeso {
		t.Fatalf("depth = %s, want %s (micro -> meso -> macro per PROJECT.md §8.5)", depth, learning.DepthMeso)
	}
	if reason == "" {
		t.Fatal("PROJECT.md §8.5 requires every automatic change to be explained")
	}
}

func TestSuggestGranularityNarrowsOnRepeatedDifficulty(t *testing.T) {
	proj := mastery.DimensionProjection{State: mastery.StateIntroduced, EvidenceCount: EvidenceThreshold}
	depth, reason, ok := SuggestGranularity(learning.DepthMacro, proj)
	if !ok {
		t.Fatal("expected a suggestion at threshold")
	}
	if depth != learning.DepthMeso {
		t.Fatalf("depth = %s, want %s (macro -> meso -> micro per PROJECT.md §8.5)", depth, learning.DepthMeso)
	}
	if reason == "" {
		t.Fatal("expected a non-empty explanation")
	}
}

func TestSuggestGranularityHasNoNeighborBeyondTheLadderEdges(t *testing.T) {
	widen := mastery.DimensionProjection{State: mastery.StateRetained, EvidenceCount: EvidenceThreshold}
	if _, _, ok := SuggestGranularity(learning.DepthChallenge, widen); ok {
		t.Fatal("challenge is already the coarsest depth; there is nothing to widen to")
	}

	narrow := mastery.DimensionProjection{State: mastery.StateNotObserved, EvidenceCount: EvidenceThreshold}
	if _, _, ok := SuggestGranularity(learning.DepthMicro, narrow); ok {
		t.Fatal("micro is already the finest depth; there is nothing to narrow to")
	}
}

func TestSuggestGranularityNeverMutatesAnything(t *testing.T) {
	// Purely a documentation-level check that the function has no side
	// effects observable through the session service: calling it never
	// touches a *Service at all, so there is nothing to assert beyond
	// its return values being pure functions of the inputs.
	proj := mastery.DimensionProjection{State: mastery.StateDemonstratesWithoutHelp, EvidenceCount: EvidenceThreshold}
	d1, r1, _ := SuggestGranularity(learning.DepthMicro, proj)
	d2, r2, _ := SuggestGranularity(learning.DepthMicro, proj)
	if d1 != d2 || r1 != r2 {
		t.Fatalf("SuggestGranularity is not deterministic: (%s,%q) vs (%s,%q)", d1, r1, d2, r2)
	}
}
