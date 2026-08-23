package mastery

import "testing"

func TestNextScheduleFirstReviewUsesFirstLadderRung(t *testing.T) {
	s := NextSchedule("comp-1", nil, true, day(1))
	if s.IntervalDays != initialLadderDays[0] {
		t.Fatalf("interval = %d, want %d", s.IntervalDays, initialLadderDays[0])
	}
	if !s.DueAt.Equal(day(1).AddDate(0, 0, initialLadderDays[0])) {
		t.Fatalf("due at = %v, want %v", s.DueAt, day(1).AddDate(0, 0, initialLadderDays[0]))
	}
}

func TestNextScheduleSuccessClimbsTheLadder(t *testing.T) {
	s := NextSchedule("comp-1", nil, true, day(1))
	for _, want := range initialLadderDays[1:] {
		s = NextSchedule("comp-1", &s, true, day(1))
		if s.IntervalDays != want {
			t.Fatalf("interval = %d, want %d", s.IntervalDays, want)
		}
	}
}

func TestNextScheduleGrowsPastTheLadderEnd(t *testing.T) {
	s := Schedule{CompetencyID: "comp-1", IntervalDays: 30}
	next := NextSchedule("comp-1", &s, true, day(1))
	if next.IntervalDays <= 30 {
		t.Fatalf("expected growth past the ladder's last rung, got %d", next.IntervalDays)
	}
}

func TestNextScheduleFailureShrinksButNeverBelowMinimum(t *testing.T) {
	s := Schedule{CompetencyID: "comp-1", IntervalDays: 1}
	next := NextSchedule("comp-1", &s, false, day(1))
	if next.IntervalDays < minIntervalDays {
		t.Fatalf("interval = %d, must never go below %d (failure never zeroes learning)", next.IntervalDays, minIntervalDays)
	}
}

func TestNextScheduleFailureReducesFromALargerInterval(t *testing.T) {
	s := Schedule{CompetencyID: "comp-1", IntervalDays: 30}
	next := NextSchedule("comp-1", &s, false, day(1))
	if next.IntervalDays >= 30 {
		t.Fatalf("expected the interval to shrink after failure, got %d", next.IntervalDays)
	}
}

func TestDueSchedulesOrdersMostOverdueFirstWithDeterministicTieBreak(t *testing.T) {
	now := day(10)
	all := map[string]Schedule{
		"comp-b": {CompetencyID: "comp-b", DueAt: day(5)},  // 5 days overdue
		"comp-a": {CompetencyID: "comp-a", DueAt: day(5)},  // 5 days overdue, tie with comp-b
		"comp-c": {CompetencyID: "comp-c", DueAt: day(8)},  // 2 days overdue
		"comp-d": {CompetencyID: "comp-d", DueAt: day(12)}, // not due yet
	}
	due := DueSchedules(all, now)
	if len(due) != 3 {
		t.Fatalf("expected 3 due schedules, got %d: %+v", len(due), due)
	}
	if due[0].CompetencyID != "comp-a" || due[1].CompetencyID != "comp-b" {
		t.Fatalf("expected comp-a and comp-b (tied, alphabetical) first, got %+v", due)
	}
	if due[2].CompetencyID != "comp-c" {
		t.Fatalf("expected comp-c last among due, got %+v", due)
	}
}
