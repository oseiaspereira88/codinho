package eventstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ErrRevisionConflict is returned by Append when expectedRevision does not
// match the stream's current revision (optimistic concurrency, requirement
// R4).
var ErrRevisionConflict = errors.New("eventstore: revision conflict")

// seenKey scopes the idempotency cache to one stream: request_id is only
// unique per stream, never globally, so two streams reusing the same
// request_id must never collide (eventstore-idempotency-scope R1/R3).
func seenKey(streamID, requestID string) string {
	return streamID + "\x00" + requestID
}

// Store is a durable, append-only event log for one workspace. It is safe
// for concurrent use by multiple goroutines within one process; exclusive
// access across processes is the caller's responsibility via AcquireLock
// (requirement R7).
type Store struct {
	mu        sync.Mutex
	file      *os.File
	revisions map[string]uint64
	seen      map[string]Event
	events    []Event
}

// Open opens (creating if absent) the event log at path and replays it to
// rebuild the in-memory revision index and idempotency cache (requirement
// R2). A truncated final line from a prior crash is tolerated and the log
// is reopened for append from the valid prefix; any other recovery problem
// is returned so the caller can decide how to proceed rather than silently
// losing data.
func Open(path string, registry *UpcasterRegistry) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, err
	}
	result := ReadEvents(path, registry)
	if result.Err != nil {
		return nil, result.Err
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, err
	}

	s := &Store{
		file:      f,
		revisions: map[string]uint64{},
		seen:      map[string]Event{},
		events:    result.Events,
	}
	for _, ev := range result.Events {
		s.revisions[ev.StreamID] = ev.Revision
		if ev.RequestID != "" {
			s.seen[seenKey(ev.StreamID, ev.RequestID)] = ev
		}
	}
	return s, nil
}

// Close releases the underlying file handle. It does not release
// AcquireLock; callers hold and release the workspace lock independently.
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.file.Close()
}

// Revision returns the current revision of streamID, or 0 if it has no
// events yet.
func (s *Store) Revision(streamID string) uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.revisions[streamID]
}

// Append durably appends one event to streamID and returns it.
//
// expectedRevision must equal the stream's current revision (0 for a new
// stream); otherwise Append returns ErrRevisionConflict without writing
// (requirement R4, optimistic concurrency).
//
// When requestID is non-empty and was already used for this stream, Append
// returns the original event without appending again (requirement R4,
// idempotent retries).
//
// Append fsyncs before returning, so a caller-visible success always means
// the event survives a crash (requirement R1, R3: snapshots may only be
// written after this call returns successfully).
func (s *Store) Append(streamID string, expectedRevision uint64, requestID string, eventType EventType, payload any) (Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if requestID != "" {
		if existing, ok := s.seen[seenKey(streamID, requestID)]; ok {
			return existing, nil
		}
	}

	current := s.revisions[streamID]
	if expectedRevision != current {
		return Event{}, fmt.Errorf("%w: stream %s at revision %d, expected %d", ErrRevisionConflict, streamID, current, expectedRevision)
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return Event{}, fmt.Errorf("marshaling payload: %w", err)
	}

	revision := current + 1
	ev := Event{
		SchemaVersion: CurrentEventSchemaVersion,
		ID:            fmt.Sprintf("%s@%d", streamID, revision),
		StreamID:      streamID,
		Revision:      revision,
		Type:          eventType,
		RequestID:     requestID,
		RecordedAt:    time.Now().UTC().Format(time.RFC3339Nano),
		Payload:       raw,
	}

	line, err := json.Marshal(ev)
	if err != nil {
		return Event{}, fmt.Errorf("marshaling event: %w", err)
	}
	line = append(line, '\n')
	if _, err := s.file.Write(line); err != nil {
		return Event{}, fmt.Errorf("writing event: %w", err)
	}
	if err := s.file.Sync(); err != nil {
		return Event{}, fmt.Errorf("syncing event: %w", err)
	}

	s.revisions[streamID] = revision
	if requestID != "" {
		s.seen[seenKey(streamID, requestID)] = ev
	}
	s.events = append(s.events, ev)
	return ev, nil
}

// Replay returns every event for streamID in revision order.
func (s *Store) Replay(streamID string) []Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []Event
	for _, ev := range s.events {
		if ev.StreamID == streamID {
			out = append(out, ev)
		}
	}
	return out
}

// ReplayAll returns every event in the log in append order, across every
// stream (requirement R2).
func (s *Store) ReplayAll() []Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Event, len(s.events))
	copy(out, s.events)
	return out
}

// Export returns streamID's events after applying filter to each one, in
// order. filter returns the (possibly redacted) event to keep and whether
// to keep it at all, so a caller can honor an exclusion policy without this
// store knowing what that policy means (requirement R8).
func (s *Store) Export(streamID string, filter func(Event) (Event, bool)) []Event {
	var out []Event
	for _, ev := range s.Replay(streamID) {
		if redacted, keep := filter(ev); keep {
			out = append(out, redacted)
		}
	}
	return out
}
