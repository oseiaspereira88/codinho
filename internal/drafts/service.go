// Package drafts persists quarantined authoring data, never a published catalog.
package drafts

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"time"
)

const stream = "drafts"
const MaxSubmissions = 100
const MaxStoredBytes = 16 << 20
const Lifetime = 30 * 24 * time.Hour

var ErrUnavailable = errors.New("drafts: unavailable, removed or expired")
var ErrConflict = errors.New("drafts: request or pack identity conflict")
var ErrQuota = errors.New("drafts: cumulative quota reached; export and purge local state")

type Record struct {
	ID        string    `json:"draft_id"`
	PackID    string    `json:"pack_id"`
	Version   string    `json:"version"`
	Digest    string    `json:"digest"`
	YAML      string    `json:"pack_yaml"`
	ExpiresAt time.Time `json:"expires_at"`
}
type Service struct {
	store *eventstore.Store
	now   func() time.Time
}

func New(store *eventstore.Store) *Service {
	return &Service{store: store, now: func() time.Time { return time.Now().UTC() }}
}
func digest(s string) string { return fmt.Sprintf("%x", sha256.Sum256([]byte(s))) }
func (s *Service) Submit(raw, requestID string) (Record, curriculum.DraftValidation, error) {
	var v curriculum.DraftValidation
	if requestID == "" || len(requestID) > 200 {
		return Record{}, v, ErrConflict
	}
	hash := digest(raw)
	if ev, ok := s.store.Request(stream, requestID); ok {
		var r Record
		if ev.Type != eventstore.EventDraftSubmitted || json.Unmarshal(ev.Payload, &r) != nil || r.Digest != hash {
			return Record{}, v, ErrConflict
		}
		return r, v, nil
	}
	revision := s.store.Revision(stream)
	count, total := 0, 0
	p, v, err := curriculum.ParseDraft([]byte(raw))
	if err != nil {
		return Record{}, v, err
	}
	for _, ev := range s.store.Replay(stream) {
		if ev.Type != eventstore.EventDraftSubmitted {
			continue
		}
		var r Record
		if json.Unmarshal(ev.Payload, &r) != nil {
			return Record{}, v, ErrUnavailable
		}
		count++
		total += len(r.YAML)
		if r.PackID == p.ID {
			return Record{}, v, ErrConflict
		}
	}
	if count >= MaxSubmissions || total+len(raw) > MaxStoredBytes {
		return Record{}, v, ErrQuota
	}
	r := Record{ID: "draft_" + hash, PackID: p.ID, Version: p.Version, Digest: hash, YAML: raw, ExpiresAt: s.now().Add(Lifetime)}
	_, err = s.store.Append(stream, revision, requestID, eventstore.EventDraftSubmitted, r)
	return r, v, err
}
func (s *Service) Get(id string) (Record, error) {
	var found Record
	for _, ev := range s.store.Replay(stream) {
		switch ev.Type {
		case eventstore.EventDraftSubmitted:
			var r Record
			if json.Unmarshal(ev.Payload, &r) != nil {
				return Record{}, ErrUnavailable
			}
			if r.ID == id {
				found = r
			}
		case eventstore.EventDraftRemoved:
			var r struct {
				ID string `json:"draft_id"`
			}
			if json.Unmarshal(ev.Payload, &r) != nil {
				return Record{}, ErrUnavailable
			}
			if r.ID == id {
				return Record{}, ErrUnavailable
			}
		}
	}
	if found.ID == "" || !s.now().Before(found.ExpiresAt) || digest(found.YAML) != found.Digest {
		return Record{}, ErrUnavailable
	}
	return found, nil
}
func (s *Service) Catalog(id string) (*curriculum.Catalog, error) {
	r, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	p, _, err := curriculum.ParseDraft([]byte(r.YAML))
	if err != nil {
		return nil, err
	}
	return curriculum.DraftCatalog(p), nil
}
func (s *Service) Remove(id, requestID string, confirm bool) error {
	if !confirm || requestID == "" || len(requestID) > 200 {
		return ErrConflict
	}
	if ev, ok := s.store.Request(stream, requestID); ok {
		var p struct {
			ID string `json:"draft_id"`
		}
		if ev.Type == eventstore.EventDraftRemoved && json.Unmarshal(ev.Payload, &p) == nil && p.ID == id {
			return nil
		}
		return ErrConflict
	}
	revision := s.store.Revision(stream)
	if _, err := s.Get(id); err != nil {
		return err
	}
	_, err := s.store.Append(stream, revision, requestID, eventstore.EventDraftRemoved, map[string]string{"draft_id": id})
	return err
}
