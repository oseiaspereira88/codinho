package session

import (
	"encoding/json"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/drafts"
)

// trackSnapshot is private durable state. Public status never exposes future
// instructions, solutions, checks or fixture contents.
type trackSnapshot struct {
	Version       int                             `json:"version"`
	TrackID       string                          `json:"track_id,omitempty"`
	CompositionID string                          `json:"composition_id,omitempty"`
	Challenges    []curriculum.ChallengeAuthoring `json:"challenges"`
	Digest        string                          `json:"digest"`
}
type TrackStatus struct {
	TrackID          string `json:"track_id,omitempty"`
	CompositionID    string `json:"composition_id,omitempty"`
	ChallengeID      string `json:"challenge_id"`
	ChallengeVersion string `json:"challenge_version"`
	Cursor           int    `json:"cursor"`
	Total            int    `json:"total"`
}

func (r *record) trackStatus() *TrackStatus {
	if r.track == nil {
		return nil
	}
	return &TrackStatus{TrackID: r.track.TrackID, CompositionID: r.track.CompositionID, ChallengeID: r.challengeID, ChallengeVersion: r.pinned.Version, Cursor: r.trackCursor, Total: len(r.track.Challenges)}
}
func (s *Service) resolveStart(in StartInput) (curriculum.ChallengeAuthoring, *trackSnapshot, error) {
	catalog := s.catalog
	if in.DraftID != "" {
		if !in.AcceptDraft {
			return curriculum.ChallengeAuthoring{}, nil, curriculum.ErrInvalidSelection
		}
		var err error
		catalog, err = drafts.New(s.store).Catalog(in.DraftID)
		if err != nil {
			return curriculum.ChallengeAuthoring{}, nil, err
		}
	} else if in.AcceptDraft {
		return curriculum.ChallengeAuthoring{}, nil, curriculum.ErrInvalidSelection
	}
	count := 0
	for _, id := range []string{in.ChallengeID, in.TrackID, in.CompositionID} {
		if id != "" {
			count++
		}
	}
	if count != 1 || (in.CompositionID == "" && in.ChallengeIDs != nil) {
		return curriculum.ChallengeAuthoring{}, nil, curriculum.ErrInvalidSelection
	}
	if in.ChallengeID != "" {
		ch, ok := catalog.Challenge(in.ChallengeID)
		if !ok {
			return ch, nil, ErrNotFound
		}
		return ch, nil, nil
	}
	ids := in.ChallengeIDs
	if in.TrackID != "" {
		t, ok := catalog.Track(in.TrackID)
		if !ok {
			return curriculum.ChallengeAuthoring{}, nil, ErrNotFound
		}
		ids = t.ChallengeIDs
	}
	chs, err := catalog.ResolvePath(ids)
	if err != nil {
		return curriculum.ChallengeAuthoring{}, nil, err
	}
	digest := curriculum.PathIdentity(chs)
	if in.CompositionID != "" && in.CompositionID != digest {
		return curriculum.ChallengeAuthoring{}, nil, curriculum.ErrInvalidSelection
	}
	// Own immutable snapshots independently from caller and catalog slices.
	raw, err := json.Marshal(chs)
	if err != nil {
		return curriculum.ChallengeAuthoring{}, nil, err
	}
	if err = json.Unmarshal(raw, &chs); err != nil {
		return curriculum.ChallengeAuthoring{}, nil, err
	}
	track := &trackSnapshot{Version: 1, TrackID: in.TrackID, CompositionID: in.CompositionID, Challenges: chs, Digest: digest}
	return chs[0], track, nil
}
func (r *record) nextChallenge() (navigation, error) {
	if r.track == nil || r.trackCursor+1 >= len(r.track.Challenges) {
		return navigation{}, nil
	}
	r.trackCursor++
	r.pinned = r.track.Challenges[r.trackCursor]
	r.challengeID = r.pinned.ID
	r.session.ClearActiveInstruction()
	r.progress = map[string]savedProgress{}
	r.covered = map[string]bool{}
	r.choices = map[string]string{}
	r.cleanEvaluation = false
	r.exhausted = false
	w, ok := deriveWindow(r.pinned, r.depth)
	if !ok {
		return navigation{}, ErrNoWindowAtDepth
	}
	n, ok := findStep(r.pinned, w.StepID)
	if !ok {
		return navigation{}, ErrStepNotFound
	}
	r.cursor = n.ID
	return navigation{node: n}, nil
}

func cloneStartResult(r StartResult) StartResult {
	if r.Track != nil {
		status := *r.Track
		r.Track = &status
	}
	return r
}

// snapshotProvenance is fixed for the complete path, never upgraded on replay.
func snapshotProvenance(ch curriculum.ChallengeAuthoring, track *trackSnapshot) string {
	if ch.Publication.Status != "published" {
		return "draft"
	}
	if track != nil {
		for _, member := range track.Challenges {
			if member.Publication.Status != "published" {
				return "draft"
			}
		}
	}
	return "published"
}
