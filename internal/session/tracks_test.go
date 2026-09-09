package session

import (
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/learning"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func trackService(t *testing.T) *Service {
	t.Helper()
	base := newTestService(t)
	ch, _ := base.catalog.Challenge(fixtureChallengeID)
	a := ch
	a.ID = "a"
	b := ch
	b.ID = "b"
	b.Version = "2.0.0"
	b.Prerequisites = []string{"a"}
	pack := curriculum.Pack{SchemaVersion: 1, ID: "tracks", Version: "1.0.0", Challenges: []curriculum.ChallengeAuthoring{a, b}, Tracks: []curriculum.TrackAuthoring{{ID: "track", Title: "Track", ChallengeIDs: []string{"a", "b"}}}}
	data, err := yaml.Marshal(pack)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	writeTestFile(t, dir, "manifest.yaml", "schema_version: 1\npacks: [pack.yaml]\n")
	if err := os.WriteFile(filepath.Join(dir, "pack.yaml"), data, 0600); err != nil {
		t.Fatal(err)
	}
	cat, ds, err := curriculum.Load(dir, curriculum.DefaultLimits)
	if err != nil || cat == nil {
		t.Fatalf("%v %+v", err, ds)
	}
	store, err := eventstore.Open(filepath.Join(t.TempDir(), "events.jsonl"), nil)
	if err != nil {
		t.Fatal(err)
	}
	return New(cat, store, nil)
}
func TestTrackSessionReplayAndBoundaries(t *testing.T) {
	s := trackService(t)
	in := StartInput{TrackID: "track", Depth: learning.DepthChallenge, RequestID: "start"}
	start, err := s.Start(in)
	if err != nil {
		t.Fatal(err)
	}
	if start.Track == nil || start.Track.Cursor != 0 || start.Track.Total != 2 {
		t.Fatal(start)
	}
	if _, err = s.StepAdvance(start.SessionID, false, start.Revision, "denied"); err == nil {
		t.Fatal("advanced without completion")
	}
	if _, err = s.StepAdvance(start.SessionID, true, start.Revision+1, "stale"); err == nil {
		t.Fatal("stale revision accepted")
	}
	first, err := s.StepAdvance(start.SessionID, true, start.Revision, "next")
	if err != nil || first.Done || first.StepID != "b" {
		t.Fatalf("%+v %v", first, err)
	}
	// Replay with an unrelated catalog: every future challenge comes from the start.
	s = New(newTestService(t).catalog, s.store, nil)
	got, err := s.Get(start.SessionID)
	if err != nil || got.Track == nil || got.Track.Cursor != 1 || got.Track.ChallengeVersion != "2.0.0" {
		t.Fatalf("%+v %v", got, err)
	}
	retry, err := s.StepAdvance(start.SessionID, true, start.Revision, "next")
	if err != nil || !reflect.DeepEqual(first, retry) {
		t.Fatalf("retry %+v %v", retry, err)
	}
	again, err := s.Start(in)
	if err != nil || !reflect.DeepEqual(start, again) {
		t.Fatalf("start retry %+v %v", again, err)
	}
	end, err := s.StepAdvance(start.SessionID, true, got.Revision, "end")
	if err != nil || !end.Done {
		t.Fatalf("%+v %v", end, err)
	}
	s = New(s.catalog, s.store, nil)
	end2, err := s.StepAdvance(start.SessionID, true, end.Revision, "end")
	if err != nil || !reflect.DeepEqual(end, end2) {
		t.Fatalf("end retry %+v %v", end2, err)
	}
}
func TestTrackSelectorsAndComposition(t *testing.T) {
	s := trackService(t)
	chs, err := s.catalog.ResolvePath([]string{"a", "b"})
	if err != nil {
		t.Fatal(err)
	}
	for _, in := range []StartInput{{}, {TrackID: "missing"}, {TrackID: "track", ChallengeID: "a"}, {CompositionID: "bad", ChallengeIDs: []string{"a", "b"}}, {ChallengeIDs: []string{"a"}}, {TrackID: "track", ChallengeIDs: []string{"a"}}} {
		if _, err := s.Start(in); err == nil {
			t.Fatalf("accepted %+v", in)
		}
	}
	in := StartInput{CompositionID: curriculum.PathIdentity(chs), ChallengeIDs: []string{"a", "b"}, Depth: learning.DepthMicro}
	start, err := s.Start(in)
	if err != nil {
		t.Fatal(err)
	}
	h, err := s.HintRequest(start.SessionID, false, start.Revision, "h")
	if err != nil {
		t.Fatal(err)
	}
	next, err := s.StepAdvance(start.SessionID, true, h.Revision, "n")
	if err != nil || next.Done {
		t.Fatalf("%+v %v", next, err)
	}
	// Both challenges deliberately reuse step IDs; disclosure must reset.
	got, err := s.Get(start.SessionID)
	if err != nil || got.Track.Cursor != 1 || got.Disclosure.SolutionRevealed {
		t.Fatalf("%+v %v", got, err)
	}
	active := s.sessions[start.SessionID].session.ActiveStep()
	if active.HintLevel != 0 {
		t.Fatal("hint crossed challenge boundary")
	}
}

func TestTrackExplicitCompletionTraversesAllChallenges(t *testing.T) {
	s := trackService(t)
	start, err := s.Start(StartInput{TrackID: "track", Depth: learning.DepthChallenge})
	if err != nil {
		t.Fatal(err)
	}
	rev := start.Revision
	for i := 0; i < 2; i++ {
		eval, err := s.StepEvaluate(EvaluateInput{SessionID: start.SessionID, ExpectedRevision: rev})
		if err != nil {
			t.Fatal(err)
		}
		done, err := s.StepComplete(start.SessionID, true, false, eval.Revision, "")
		if err != nil {
			t.Fatal(err)
		}
		next, err := s.StepAdvance(start.SessionID, false, done.Revision, "")
		if err != nil || next.Done != (i == 1) {
			t.Fatalf("index %d: %+v %v", i, next, err)
		}
		rev = next.Revision
	}
	got, err := s.Get(start.SessionID)
	if err != nil || got.State != learning.SessionStateActive {
		t.Fatal("tree exhaustion implicitly finished session", err)
	}
}

func TestSnapshotProvenanceIncludesEveryTrackMember(t *testing.T) {
	reviewed := curriculum.ChallengeAuthoring{Publication: curriculum.PublicationAuthoring{Status: "published"}}
	draft := curriculum.ChallengeAuthoring{}
	if snapshotProvenance(reviewed, nil) != "published" {
		t.Fatal("reviewed snapshot misclassified")
	}
	if snapshotProvenance(reviewed, &trackSnapshot{Challenges: []curriculum.ChallengeAuthoring{reviewed, draft}}) != "draft" {
		t.Fatal("future draft member not quarantined")
	}
}
