package session

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/oseiaspereira88/codinho/internal/learning"
)

func TestDefaultsForModeAreIndependentPerMode(t *testing.T) {
	cases := []struct {
		mode learning.PedagogicalMode
		want ModeDefaults
	}{
		{learning.ModeTeaching, ModeDefaults{learning.DepthMacro, learning.HelpProgressive, learning.DisclosureConceptOrAPI, learning.EvaluationOnDemand}},
		{learning.ModePractice, ModeDefaults{learning.DepthMicro, learning.HelpProgressive, learning.DisclosureGuidingQuestion, learning.EvaluationOnDemand}},
		{learning.ModeReview, ModeDefaults{learning.DepthMeso, learning.HelpLimited, learning.DisclosureGuidingQuestion, learning.EvaluationOnStepComplete}},
		{learning.ModeDebug, ModeDefaults{learning.DepthMacro, learning.HelpProgressive, learning.DisclosureLogicalOutline, learning.EvaluationOnDemand}},
		{learning.ModeExploration, ModeDefaults{learning.DepthLayer, learning.HelpFree, learning.DisclosureLogicalOutline, learning.EvaluationOnDemand}},
		{learning.ModeInterview, ModeDefaults{learning.DepthChallenge, learning.HelpNoHints, learning.DisclosureNone, learning.EvaluationOnlyAtEnd}},
		{learning.PedagogicalMode("bogus"), ModeDefaults{learning.DepthMicro, learning.HelpProgressive, learning.DisclosureGuidingQuestion, learning.EvaluationOnDemand}},
	}
	for _, tc := range cases {
		if got := DefaultsForMode(tc.mode); got != tc.want {
			t.Errorf("DefaultsForMode(%s) = %+v, want %+v", tc.mode, got, tc.want)
		}
	}
}

func TestModeDefaultsAreAllDistinctAcrossTheFiveNonInterviewModes(t *testing.T) {
	modes := []learning.PedagogicalMode{
		learning.ModeTeaching, learning.ModePractice, learning.ModeReview, learning.ModeDebug, learning.ModeExploration,
	}
	seen := map[ModeDefaults]learning.PedagogicalMode{}
	for _, m := range modes {
		d := DefaultsForMode(m)
		if prior, ok := seen[d]; ok {
			t.Fatalf("modes %s and %s share identical defaults %+v; each mode must be independently configurable", prior, m, d)
		}
		seen[d] = m
	}
}

func TestDebugProtocolStagesIsOrderedAndComplete(t *testing.T) {
	want := []string{"reproduce", "locate_first_divergence", "hypothesize", "observe", "confirm_or_reject", "apply_smallest_fix"}
	got := DebugProtocolStages()
	if len(got) != len(want) {
		t.Fatalf("stages = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("stage[%d] = %s, want %s", i, got[i], want[i])
		}
	}
}

func TestStartAppliesModeDefaultsWhenFieldsUnset(t *testing.T) {
	svc := newTestService(t)
	result, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID, Mode: learning.ModeTeaching})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec := svc.sessions[result.SessionID]
	policy := rec.session.Policy
	if policy.InitialDepth != learning.DepthMacro {
		t.Errorf("depth = %s, want %s (teaching default)", policy.InitialDepth, learning.DepthMacro)
	}
	if policy.Help.Kind != learning.HelpProgressive {
		t.Errorf("help = %s, want %s", policy.Help.Kind, learning.HelpProgressive)
	}
	if policy.Disclosure.MaxLevel != learning.DisclosureConceptOrAPI {
		t.Errorf("disclosure = %d, want %d", policy.Disclosure.MaxLevel, learning.DisclosureConceptOrAPI)
	}
	if policy.Evaluation != learning.EvaluationOnDemand {
		t.Errorf("evaluation = %s, want %s", policy.Evaluation, learning.EvaluationOnDemand)
	}
}

func TestStartExplicitFieldsOverrideModeDefaults(t *testing.T) {
	svc := newTestService(t)
	result, err := svc.Start(StartInput{
		ChallengeID: fixtureChallengeID,
		Mode:        learning.ModeTeaching,
		Depth:       learning.DepthMicro,
		Help:        learning.HelpNoHints,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec := svc.sessions[result.SessionID]
	policy := rec.session.Policy
	if policy.InitialDepth != learning.DepthMicro {
		t.Errorf("explicit depth was overridden by mode default: got %s", policy.InitialDepth)
	}
	if policy.Help.Kind != learning.HelpNoHints {
		t.Errorf("explicit help was overridden by mode default: got %s", policy.Help.Kind)
	}
	// Disclosure was left unset, so teaching's default still applies —
	// mode defaults are independent per dimension (requirement R1).
	if policy.Disclosure.MaxLevel != learning.DisclosureConceptOrAPI {
		t.Errorf("disclosure = %d, want teaching default %d", policy.Disclosure.MaxLevel, learning.DisclosureConceptOrAPI)
	}
}

func TestGranularityAdjustPersistsReason(t *testing.T) {
	svc := newTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	const reason = "3 evidências sem ajuda: ampliando para macro"
	if _, err := svc.GranularityAdjust(start.SessionID, learning.DepthMacro, reason, start.Revision, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	events := svc.store.Replay(string(start.SessionID))
	last := events[len(events)-1]
	var payload struct {
		Reason string `json:"reason"`
	}
	if err := json.Unmarshal(last.Payload, &payload); err != nil {
		t.Fatalf("decoding payload: %v", err)
	}
	if payload.Reason != reason {
		t.Fatalf("persisted reason = %q, want %q", payload.Reason, reason)
	}
}

func TestProposeNextStepRecordsSignalWithoutAdvancing(t *testing.T) {
	svc := newTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := svc.ProposeNextStep(start.SessionID, "fixture.micro-one", start.Revision, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Revision <= start.Revision {
		t.Fatalf("revision did not advance: %d", result.Revision)
	}

	get, err := svc.Get(start.SessionID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if get.ActiveStep != "fixture.macro-one" {
		t.Fatalf("active step changed to %s; proposing a step must never advance it", get.ActiveStep)
	}
}

func TestProposeNextStepRejectsUnknownStep(t *testing.T) {
	svc := newTestService(t)
	start, err := svc.Start(StartInput{ChallengeID: fixtureChallengeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := svc.ProposeNextStep(start.SessionID, "does-not-exist", start.Revision, ""); !errors.Is(err, ErrStepNotFound) {
		t.Fatalf("err = %v, want ErrStepNotFound", err)
	}
}
