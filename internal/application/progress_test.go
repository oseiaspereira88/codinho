package application

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/mastery"
)

func newProgressTestService(t *testing.T) (*ProgressService, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "events.jsonl")
	store, err := eventstore.Open(path, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return NewProgressService(store), path
}

func TestRecordEvidenceRejectsInvalidDimension(t *testing.T) {
	svc, _ := newProgressTestService(t)
	_, err := svc.RecordEvidence(EvidenceInput{CompetencyID: "c", Dimension: "not-a-dimension", EvidenceID: "e1"})
	if !errors.Is(err, ErrInvalidDimension) {
		t.Fatalf("expected ErrInvalidDimension, got %v", err)
	}
}

func TestRecordEvidenceRejectsMissingCompetencyOrEvidence(t *testing.T) {
	svc, _ := newProgressTestService(t)
	_, err := svc.RecordEvidence(EvidenceInput{Dimension: string(mastery.DimensionExplanation)})
	if !errors.Is(err, ErrInvalidCompetency) {
		t.Fatalf("expected ErrInvalidCompetency, got %v", err)
	}
}

func TestRecordEvidenceProducesExpectedStateProgression(t *testing.T) {
	svc, _ := newProgressTestService(t)

	first, err := svc.RecordEvidence(EvidenceInput{
		CompetencyID: "slice-filter", Dimension: string(mastery.DimensionAutonomousImplementation),
		EvidenceID: "ev1", Success: true, ExpectedRevision: 0,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if first.State != string(mastery.StateDemonstratesWithoutHelp) {
		t.Fatalf("state = %s, want demonstrates_without_help", first.State)
	}
	if first.Revision != 2 { // mastery_projected + review_scheduled, one call each
		t.Fatalf("revision = %d, want 2", first.Revision)
	}

	second, err := svc.RecordEvidence(EvidenceInput{
		CompetencyID: "slice-filter", Dimension: string(mastery.DimensionAutonomousImplementation),
		EvidenceID: "ev2", Success: true, ExpectedRevision: first.Revision,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Same-day repeat: rules_test.go covers the same-day-vs-later-day
	// distinction directly; here we only check the plumbing advances.
	if second.State != string(mastery.StateDemonstratesWithoutHelp) && second.State != string(mastery.StateRetained) {
		t.Fatalf("unexpected state after second signal: %s", second.State)
	}
}

func TestRecordEvidenceRejectsRevisionConflict(t *testing.T) {
	svc, _ := newProgressTestService(t)
	if _, err := svc.RecordEvidence(EvidenceInput{CompetencyID: "c", Dimension: string(mastery.DimensionExplanation), EvidenceID: "e1", ExpectedRevision: 0}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err := svc.RecordEvidence(EvidenceInput{CompetencyID: "c", Dimension: string(mastery.DimensionExplanation), EvidenceID: "e2", ExpectedRevision: 0})
	if !errors.Is(err, eventstore.ErrRevisionConflict) {
		t.Fatalf("expected ErrRevisionConflict, got %v", err)
	}
}

func TestProgressFiltersByCompetency(t *testing.T) {
	svc, _ := newProgressTestService(t)
	r1, err := svc.RecordEvidence(EvidenceInput{CompetencyID: "comp-a", Dimension: string(mastery.DimensionExplanation), EvidenceID: "e1", Success: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := svc.RecordEvidence(EvidenceInput{CompetencyID: "comp-b", Dimension: string(mastery.DimensionExplanation), EvidenceID: "e2", Success: true, ExpectedRevision: r1.Revision}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	all, err := svc.Progress("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(all.Competencies) != 2 {
		t.Fatalf("expected 2 competencies, got %d", len(all.Competencies))
	}

	filtered, err := svc.Progress("comp-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filtered.Competencies) != 1 {
		t.Fatalf("expected 1 competency, got %d: %+v", len(filtered.Competencies), filtered.Competencies)
	}
	if _, ok := filtered.Competencies["comp-a"]; !ok {
		t.Fatalf("expected comp-a in filtered result, got %+v", filtered.Competencies)
	}
}

func TestReviewDueExcludesSolutionRevealedAndNotYetDue(t *testing.T) {
	svc, _ := newProgressTestService(t)
	if _, err := svc.RecordEvidence(EvidenceInput{
		CompetencyID: "revealed", Dimension: string(mastery.DimensionAutonomousImplementation),
		EvidenceID: "e1", Success: true, SolutionRevealed: true,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r2, err := svc.RecordEvidence(EvidenceInput{
		CompetencyID: "not-revealed", Dimension: string(mastery.DimensionAutonomousImplementation),
		EvidenceID: "e2", Success: true, ExpectedRevision: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = r2

	due, err := svc.ReviewDue(day(1)) // long before the 1-day-out schedule is due
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, s := range due.Due {
		if s.CompetencyID == "revealed" {
			t.Fatal("expected no schedule for a solution-revealed signal (invariant 8)")
		}
	}
	if len(due.Due) != 0 {
		t.Fatalf("expected nothing due yet, got %+v", due.Due)
	}
}

func day(n int) time.Time { return time.Date(2026, 1, n, 12, 0, 0, 0, time.UTC) }

func TestReviewDueIncludesOverdueSchedules(t *testing.T) {
	svc, _ := newProgressTestService(t)
	if _, err := svc.RecordEvidence(EvidenceInput{
		CompetencyID: "comp-a", Dimension: string(mastery.DimensionAutonomousImplementation), EvidenceID: "e1", Success: true,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	farFuture := time.Now().UTC().AddDate(0, 1, 0) // a month out: the 1-day schedule is long overdue
	due, err := svc.ReviewDue(farFuture)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(due.Due) != 1 || due.Due[0].CompetencyID != "comp-a" {
		t.Fatalf("expected comp-a due, got %+v", due.Due)
	}
}

func TestProgressServiceRecomputesAfterReopeningTheStore(t *testing.T) {
	svc, path := newProgressTestService(t)
	if _, err := svc.RecordEvidence(EvidenceInput{
		CompetencyID: "comp-a", Dimension: string(mastery.DimensionAutonomousImplementation), EvidenceID: "e1", Success: true,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reopened, err := eventstore.Open(path, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer reopened.Close()
	restarted := NewProgressService(reopened)

	got, err := restarted.Progress("comp-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	proj, ok := got.Competencies["comp-a"][string(mastery.DimensionAutonomousImplementation)]
	if !ok || proj.State != mastery.StateDemonstratesWithoutHelp {
		t.Fatalf("expected the projection to survive a store reopen (requirement R8), got %+v", got.Competencies)
	}
}
