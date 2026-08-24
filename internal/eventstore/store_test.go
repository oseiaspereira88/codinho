package eventstore

import (
	"errors"
	"path/filepath"
	"sync"
	"testing"
)

func openTestStore(t *testing.T) (*Store, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "events.jsonl")
	s, err := Open(path, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s, path
}

func TestAppendAssignsMonotonicRevisions(t *testing.T) {
	s, _ := openTestStore(t)

	ev1, err := s.Append("ses_1", 0, "", EventSessionStarted, map[string]string{"mode": "practice"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ev1.Revision != 1 {
		t.Fatalf("revision = %d, want 1", ev1.Revision)
	}

	ev2, err := s.Append("ses_1", 1, "", EventSessionPaused, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ev2.Revision != 2 {
		t.Fatalf("revision = %d, want 2", ev2.Revision)
	}
	if s.Revision("ses_1") != 2 {
		t.Fatalf("Revision() = %d, want 2", s.Revision("ses_1"))
	}
}

func TestAppendRejectsRevisionConflict(t *testing.T) {
	s, _ := openTestStore(t)

	if _, err := s.Append("ses_1", 0, "", EventSessionStarted, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err := s.Append("ses_1", 0, "", EventSessionPaused, nil)
	if !errors.Is(err, ErrRevisionConflict) {
		t.Fatalf("expected ErrRevisionConflict, got %v", err)
	}
	if s.Revision("ses_1") != 1 {
		t.Fatalf("a rejected append must not advance the revision, got %d", s.Revision("ses_1"))
	}
}

func TestAppendIsIdempotentByRequestID(t *testing.T) {
	s, _ := openTestStore(t)

	first, err := s.Append("ses_1", 0, "req-1", EventSessionStarted, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := s.Append("ses_1", 0, "req-1", EventSessionStarted, nil)
	if err != nil {
		t.Fatalf("unexpected error on retry: %v", err)
	}
	if first.ID != second.ID || first.Revision != second.Revision {
		t.Fatalf("retry with the same request_id produced a different event: %+v vs %+v", first, second)
	}
	if s.Revision("ses_1") != 1 {
		t.Fatalf("idempotent retry must not advance the revision, got %d", s.Revision("ses_1"))
	}
}

func TestAppendIdempotencyIsScopedPerStream(t *testing.T) {
	s, _ := openTestStore(t)

	first, err := s.Append("ses_1", 0, "shared-req", EventSessionStarted, map[string]string{"owner": "ses_1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := s.Append("ses_2", 0, "shared-req", EventSessionStarted, map[string]string{"owner": "ses_2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first.StreamID == second.StreamID {
		t.Fatalf("test setup error: expected distinct streams, got %s twice", first.StreamID)
	}
	if first.ID == second.ID {
		t.Fatalf("two different streams reusing the same request_id must never collapse into one event: %+v vs %+v", first, second)
	}
	if s.Revision("ses_2") != 1 {
		t.Fatalf("ses_2's own append must not be swallowed by ses_1's cached request_id, got revision %d", s.Revision("ses_2"))
	}

	// A real retry within ses_2 (same stream, same request_id) must still be idempotent.
	retry, err := s.Append("ses_2", 0, "shared-req", EventSessionStarted, map[string]string{"owner": "ses_2"})
	if err != nil {
		t.Fatalf("unexpected error on retry: %v", err)
	}
	if retry.ID != second.ID {
		t.Fatalf("retry with the same (stream_id, request_id) produced a different event: %+v vs %+v", retry, second)
	}
}

func TestReplayAndReplayAll(t *testing.T) {
	s, _ := openTestStore(t)
	mustAppend(t, s, "ses_1", 0, EventSessionStarted)
	mustAppend(t, s, "ses_2", 0, EventSessionStarted)
	mustAppend(t, s, "ses_1", 1, EventSessionPaused)

	ses1 := s.Replay("ses_1")
	if len(ses1) != 2 || ses1[0].Type != EventSessionStarted || ses1[1].Type != EventSessionPaused {
		t.Fatalf("Replay(ses_1) = %+v", ses1)
	}

	all := s.ReplayAll()
	if len(all) != 3 {
		t.Fatalf("ReplayAll() length = %d, want 3", len(all))
	}
}

func TestExportAppliesFilter(t *testing.T) {
	s, _ := openTestStore(t)
	mustAppend(t, s, "ses_1", 0, EventSessionStarted)
	mustAppend(t, s, "ses_1", 1, EventSolutionRevealed)
	mustAppend(t, s, "ses_1", 2, EventSessionFinished)

	exported := s.Export("ses_1", func(ev Event) (Event, bool) {
		if ev.Type == EventSolutionRevealed {
			return Event{}, false // excluded by policy
		}
		return ev, true
	})
	if len(exported) != 2 {
		t.Fatalf("Export length = %d, want 2", len(exported))
	}
	for _, ev := range exported {
		if ev.Type == EventSolutionRevealed {
			t.Fatal("excluded event leaked through Export")
		}
	}
}

func TestOpenRecoversRevisionsAndIdempotencyFromExistingLog(t *testing.T) {
	path := filepath.Join(t.TempDir(), "events.jsonl")
	s1, err := Open(path, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := s1.Append("ses_1", 0, "req-1", EventSessionStarted, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := s1.Append("ses_1", 1, "", EventSessionPaused, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := s1.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s2, err := Open(path, nil)
	if err != nil {
		t.Fatalf("unexpected error reopening: %v", err)
	}
	defer s2.Close()

	if s2.Revision("ses_1") != 2 {
		t.Fatalf("recovered revision = %d, want 2", s2.Revision("ses_1"))
	}
	replayed := s2.Replay("ses_1")
	if len(replayed) != 2 {
		t.Fatalf("recovered replay length = %d, want 2", len(replayed))
	}

	// The request_id index must also survive recovery.
	again, err := s2.Append("ses_1", 2, "req-1", EventSessionStarted, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if again.Revision != 1 {
		t.Fatalf("retrying req-1 after reopen should return the original revision 1 event, got revision %d", again.Revision)
	}
	if s2.Revision("ses_1") != 2 {
		t.Fatalf("idempotent retry after reopen must not advance the revision, got %d", s2.Revision("ses_1"))
	}
}

func TestAppendConcurrentSameStreamNeverDuplicatesRevisions(t *testing.T) {
	s, _ := openTestStore(t)

	const workers = 8
	var wg sync.WaitGroup
	successes := make(chan uint64, workers)
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			// Every worker races for revision 1; exactly one should win.
			if ev, err := s.Append("ses_1", 0, "", EventSessionStarted, nil); err == nil {
				successes <- ev.Revision
			}
		}()
	}
	wg.Wait()
	close(successes)

	count := 0
	for rev := range successes {
		if rev != 1 {
			t.Fatalf("winning append got revision %d, want 1", rev)
		}
		count++
	}
	if count != 1 {
		t.Fatalf("expected exactly one winner among %d concurrent appends, got %d", workers, count)
	}
}

// AppendAndSnapshot documents and exercises the ordering contract of
// requirement R3: a snapshot may only be written after its event has been
// synced. This is not enforced by the store (snapshot content is
// consumer-owned, Decision 3); it is a pattern the consumer must follow.
func TestSnapshotOnlyWrittenAfterEventSynced(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "events.jsonl")
	snapPath := filepath.Join(dir, "snapshot.json")

	s, err := Open(logPath, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer s.Close()

	if _, err := ReadSnapshot(snapPath); !errors.Is(err, ErrSnapshotMissing) {
		t.Fatalf("expected ErrSnapshotMissing before any snapshot exists, got %v", err)
	}

	ev, err := s.Append("ses_1", 0, "", EventSessionStarted, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Append already fsynced by the time it returns, so it is now safe to
	// derive and persist a projection snapshot.
	projection := []byte(`{"stream_id":"` + ev.StreamID + `","revision":1}`)
	if err := WriteSnapshotAtomic(snapPath, projection); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := ReadSnapshot(snapPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != string(projection) {
		t.Fatalf("snapshot content = %s, want %s", data, projection)
	}
}

func TestReadSnapshotDetectsInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snapshot.json")
	if err := WriteSnapshotAtomic(path, []byte("not json")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err := ReadSnapshot(path)
	if !errors.Is(err, ErrSnapshotInvalid) {
		t.Fatalf("expected ErrSnapshotInvalid, got %v", err)
	}
}

func TestLockIsExclusiveAndReleasable(t *testing.T) {
	path := filepath.Join(t.TempDir(), "lock")

	lock, err := AcquireLock(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := AcquireLock(path); !errors.Is(err, ErrLockHeld) {
		t.Fatalf("expected ErrLockHeld while lock is outstanding, got %v", err)
	}

	if err := lock.Release(); err != nil {
		t.Fatalf("unexpected error releasing: %v", err)
	}

	lock2, err := AcquireLock(path)
	if err != nil {
		t.Fatalf("expected to reacquire after release, got %v", err)
	}
	_ = lock2.Release()
}

func mustAppend(t *testing.T, s *Store, streamID string, expectedRevision uint64, eventType EventType) {
	t.Helper()
	if _, err := s.Append(streamID, expectedRevision, "", eventType, nil); err != nil {
		t.Fatalf("unexpected error appending %s: %v", eventType, err)
	}
}
