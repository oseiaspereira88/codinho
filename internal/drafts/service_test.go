package drafts

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/oseiaspereira88/codinho/internal/eventstore"
)

const testDraft = `schema_version: 1
id: synthetic-draft
version: 1.0.0
competencies:
  - id: synthetic-skill
    title: Practice
challenges:
  - schema_version: 1
    id: synthetic-challenge
    version: 1.0.0
    competencies:
      primary: [synthetic-skill]
    acceptance: [Explain the outcome]
    layers:
      - id: understanding
        macro_steps:
          - id: explain
            kind: micro
            instruction: {objective: Explain the outcome}
`

func TestDraftReplayIdempotencyAndRemoval(t *testing.T) {
	path := filepath.Join(t.TempDir(), "events.jsonl")
	store, err := eventstore.Open(path, nil)
	if err != nil {
		t.Fatal(err)
	}
	s := New(store)
	r, v, err := s.Submit(testDraft, "submit")
	if err != nil {
		t.Fatalf("%v %+v", err, v)
	}
	retry, _, err := s.Submit(testDraft, "submit")
	if err != nil || retry != r || store.Revision(stream) != 1 {
		t.Fatal("retry not idempotent", err)
	}
	if _, _, err = s.Submit(testDraft+"\n", "submit"); !errors.Is(err, ErrConflict) {
		t.Fatal("changed request accepted", err)
	}
	if _, _, err = s.Submit(testDraft, "another"); !errors.Is(err, ErrConflict) {
		t.Fatal("pack overwrite accepted", err)
	}
	if err = store.Close(); err != nil {
		t.Fatal(err)
	}
	store, err = eventstore.Open(path, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	s = New(store)
	restored, err := s.Get(r.ID)
	if err != nil || restored != r {
		t.Fatal("replay mismatch", err)
	}
	if _, err = s.Catalog(r.ID); err != nil {
		t.Fatal(err)
	}
	if err = s.Remove(r.ID, "remove", false); !errors.Is(err, ErrConflict) {
		t.Fatal("missing consent accepted", err)
	}
	if err = s.Remove(r.ID, "remove", true); err != nil {
		t.Fatal(err)
	}
	if err = s.Remove(r.ID, "remove", true); err != nil {
		t.Fatal("remove retry", err)
	}
	if _, err = s.Catalog(r.ID); !errors.Is(err, ErrUnavailable) {
		t.Fatal("removed draft accessible", err)
	}
	if _, _, err = s.Submit(testDraft, "new-request"); !errors.Is(err, ErrConflict) {
		t.Fatal("removed pack identity reused", err)
	}
}
func TestDraftExpiryAndCumulativeQuota(t *testing.T) {
	store, err := eventstore.Open(filepath.Join(t.TempDir(), "events.jsonl"), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	s := New(store)
	now := time.Date(2026, 9, 9, 0, 0, 0, 0, time.UTC)
	s.now = func() time.Time { return now }
	r, _, err := s.Submit(testDraft, "one")
	if err != nil {
		t.Fatal(err)
	}
	now = now.Add(Lifetime)
	if _, err = s.Get(r.ID); !errors.Is(err, ErrUnavailable) {
		t.Fatal("expired draft accessible", err)
	}
	// Synthetic persisted history exercises the cumulative cap without 100 costly authoring submissions.
	for i := 1; i < MaxSubmissions; i++ {
		_, err = store.Append(stream, store.Revision(stream), "", eventstore.EventDraftSubmitted, Record{PackID: "historical"})
		if err != nil {
			t.Fatal(err)
		}
	}
	if _, _, err = s.Submit(strings.Replace(testDraft, "synthetic-draft", "another-pack", 1), "over-quota"); !errors.Is(err, ErrQuota) {
		t.Fatal("quota ignored", err)
	}
}
