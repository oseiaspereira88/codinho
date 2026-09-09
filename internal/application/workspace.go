package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/evidence"
	"github.com/oseiaspereira88/codinho/internal/learning"
	"github.com/oseiaspereira88/codinho/internal/workspace"
)

// ErrWorkspaceRootInvalid is returned when the caller-supplied root does
// not exist, is not a directory, or cannot be resolved to a real path
// (requirement R1).
var ErrWorkspaceRootInvalid = errors.New("application: workspace root is invalid")

// ErrEvidenceOutOfScope is returned by EvidenceGet when evidence_id was
// never recorded for the requesting session — regardless of whether it
// exists in the store for a different session (requirement R8; security).
var ErrEvidenceOutOfScope = errors.New("application: evidence does not belong to this session")

// maxDeclaredFiles bounds how many changed .go files get an AST summary
// per observation, so a learner rewriting hundreds of files in one step
// cannot make a single call unbounded (non-functional requirement).
const maxDeclaredFiles = 50

type baselineKey struct {
	session learning.SessionID
	step    learning.StepID
}

// baselineEntry remembers not just the baseline hashes but also the root
// and globs the caller used to establish/refresh them, so a later
// check_run (safe-check-executor) can reuse the same workspace without
// repeating those parameters (safe-check-executor Decision 2).
type baselineEntry struct {
	baseline workspace.Baseline
	root     string
	globs    []string
}

type observationEvent struct {
	ContentProvenance string             `json:"content_provenance"`
	RecoveryVersion   int                `json:"recovery_version"`
	StepID            string             `json:"step_id"`
	EvidenceID        string             `json:"evidence_id"`
	Kind              string             `json:"kind"`
	Fingerprint       string             `json:"fingerprint"`
	Baseline          workspace.Baseline `json:"baseline"`
	Root              string             `json:"root"`
	Globs             []string           `json:"globs"`
}

// observationPayload is what this package stores as evidence content: a
// structured summary of what changed, never the raw before/after file
// bytes (keeps evidence small and avoids re-exposing anything the
// exclude/redact boundary in internal/workspace already stripped).
type observationPayload struct {
	Kind         string                              `json:"kind"` // "baseline" or "diff"
	Commit       string                              `json:"commit,omitempty"`
	FromGit      bool                                `json:"from_git"`
	Dirty        bool                                `json:"dirty"`
	Fixtures     []string                            `json:"fixtures,omitempty"`
	Changes      []workspace.FileChange              `json:"changes,omitempty"`
	Declarations map[string]workspace.GoDeclarations `json:"declarations,omitempty"`
	Diff         string                              `json:"diff,omitempty"`
	Fingerprint  string                              `json:"fingerprint"`
}

// WorkspaceService orchestrates read-only workspace observation: it never
// edits learner files (constraint), and creates Observation/Evidence
// without implying an attempt or evaluation (requirement R7).
type WorkspaceService struct {
	store    *eventstore.Store
	evidence *evidence.Store

	mu            sync.Mutex
	baselines     map[baselineKey]baselineEntry
	scope         map[learning.SessionID]map[string]bool
	evidenceRoots map[learning.SessionID]map[string]baselineEntry
}

// NewWorkspaceService wires a WorkspaceService to the shared event and
// evidence stores.
func NewWorkspaceService(store *eventstore.Store, evidenceStore *evidence.Store) *WorkspaceService {
	w := &WorkspaceService{
		store:         store,
		evidence:      evidenceStore,
		baselines:     map[baselineKey]baselineEntry{},
		scope:         map[learning.SessionID]map[string]bool{},
		evidenceRoots: map[learning.SessionID]map[string]baselineEntry{},
	}
	for _, ev := range store.ReplayAll() {
		if ev.Type != eventstore.EventObservationRecorded && ev.Type != eventstore.EventCheckExecuted && ev.Type != eventstore.EventEvidenceRecorded {
			continue
		}
		var p observationEvent
		if json.Unmarshal(ev.Payload, &p) != nil || p.EvidenceID == "" {
			continue
		}
		id := learning.SessionID(ev.StreamID)
		key := baselineKey{id, learning.StepID(p.StepID)}
		if ev.Type == eventstore.EventObservationRecorded && p.RecoveryVersion == 1 && p.Root != "" {
			w.baselines[key] = baselineEntry{baseline: p.Baseline, root: p.Root, globs: p.Globs}
		}
		w.RecordEvidence(id, p.EvidenceID)
		if entry, ok := w.baselines[key]; ok && ev.Type != eventstore.EventEvidenceRecorded {
			if w.evidenceRoots[id] == nil {
				w.evidenceRoots[id] = map[string]baselineEntry{}
			}
			w.evidenceRoots[id][p.EvidenceID] = entry
		}
	}
	return w
}

// ObserveInput is one workspace_observe call.
type ObserveInput struct {
	SessionID        learning.SessionID
	StepID           learning.StepID
	Root             string
	Globs            []string
	ExpectedRevision uint64
	RequestID        string
}

// ObserveResult is what workspace_observe reports.
type ObserveResult struct {
	Observation learning.Observation
	EvidenceID  string
	Baseline    bool // true when this call established the baseline, rather than diffing against one
	Changes     []workspace.FileChange
	FromGit     bool
	Dirty       bool
	Fingerprint string
	Revision    uint64
}

// Observe records the workspace's current state for (SessionID, StepID).
// The first call for a given step establishes its baseline (requirement
// R2); every later call diffs the current state against that baseline,
// restricted to Globs, never attributing a change that predates the
// baseline to this step (requirement R4).
func (w *WorkspaceService) Observe(in ObserveInput) (ObserveResult, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	root, err := workspace.AuthorizeRoot(in.Root)
	if err != nil {
		return ObserveResult{}, fmt.Errorf("%w: %w", ErrWorkspaceRootInvalid, err)
	}

	key := baselineKey{session: in.SessionID, step: in.StepID}
	entry, hasBaseline := w.baselines[key]
	if hasBaseline && (entry.root != root.Path() || !slices.Equal(entry.globs, in.Globs)) {
		return ObserveResult{}, ErrWorkspaceRootInvalid
	}
	if ev, ok := w.store.Request(string(in.SessionID), in.RequestID); ok {
		var previous observationEvent
		if ev.Type != eventstore.EventObservationRecorded || json.Unmarshal(ev.Payload, &previous) != nil || previous.RecoveryVersion != 1 || previous.Root != root.Path() || previous.StepID != string(in.StepID) || !slices.Equal(previous.Globs, in.Globs) {
			return ObserveResult{}, eventstore.ErrRevisionConflict
		}
		return w.observationResult(ev, previous.EvidenceID, in.StepID)
	}

	var (
		result      workspace.ObserveResult
		newBaseline workspace.Baseline
		isBaseline  bool
	)
	if hasBaseline {
		result, err = workspace.Observe(root, in.Globs, entry.baseline)
	} else {
		newBaseline, err = workspace.Capture(root, in.Globs)
		if err == nil {
			result = workspace.ObserveResult{Baseline: newBaseline, Current: newBaseline}
			isBaseline = true
		}
	}
	if err != nil {
		return ObserveResult{}, err
	}

	payload := observationPayload{
		Kind:        observationKind(isBaseline),
		Commit:      result.Current.Commit,
		FromGit:     result.Current.FromGit,
		Dirty:       result.Current.Dirty,
		Changes:     result.Changes,
		Fingerprint: fingerprint(result, isBaseline),
	}
	if isBaseline {
		payload.Fixtures = result.Current.Fixtures
	} else {
		payload.Declarations = goDeclarationsFor(root, result.Changes)
		if diff, ok, derr := diffTextFor(root, result.Changes); derr == nil && ok {
			payload.Diff = diff
		}
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return ObserveResult{}, err
	}
	evidenceID, err := w.evidence.Put(raw)
	if err != nil {
		return ObserveResult{}, err
	}

	if isBaseline {
		entry = baselineEntry{baseline: newBaseline, root: root.Path(), globs: slices.Clone(in.Globs)}
	}
	ev, err := w.store.Append(string(in.SessionID), in.ExpectedRevision, in.RequestID, eventstore.EventObservationRecorded, observationEvent{ContentProvenance: sessionProvenance(w.store, string(in.SessionID)),
		RecoveryVersion: 1, StepID: string(in.StepID), EvidenceID: evidenceID, Kind: payload.Kind, Fingerprint: payload.Fingerprint,
		Baseline: entry.baseline, Root: entry.root, Globs: entry.globs,
	})
	if err != nil {
		return ObserveResult{}, err
	}

	w.baselines[key] = entry
	if w.scope[in.SessionID] == nil {
		w.scope[in.SessionID] = map[string]bool{}
	}
	w.scope[in.SessionID][evidenceID] = true
	if w.evidenceRoots[in.SessionID] == nil {
		w.evidenceRoots[in.SessionID] = map[string]baselineEntry{}
	}
	w.evidenceRoots[in.SessionID][evidenceID] = entry
	return w.observationResult(ev, evidenceID, in.StepID)
}

func (w *WorkspaceService) observationResult(ev eventstore.Event, evidenceID string, stepID learning.StepID) (ObserveResult, error) {
	data, err := w.evidence.Get(evidenceID)
	if err != nil {
		return ObserveResult{}, err
	}
	var payload observationPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return ObserveResult{}, err
	}
	learningEvidence, err := learning.NewEvidence(learning.EvidenceID(evidenceID), learning.EvidenceKindDiff, evidenceID)
	if err != nil {
		return ObserveResult{}, err
	}
	observed, err := time.Parse(time.RFC3339Nano, ev.RecordedAt)
	if err != nil {
		return ObserveResult{}, err
	}
	observation := learning.Observation{StepID: stepID, Evidence: learningEvidence, Observed: observed}

	return ObserveResult{
		Observation: observation,
		EvidenceID:  evidenceID,
		Baseline:    payload.Kind == "baseline",
		Changes:     payload.Changes,
		FromGit:     payload.FromGit,
		Dirty:       payload.Dirty,
		Fingerprint: payload.Fingerprint,
		Revision:    ev.Revision,
	}, nil
}

// EvidenceGetInput is one evidence_get call. Root and Globs are optional:
// when supplied, the result reports whether the workspace drifted since
// the evidence was collected (requirement R9).
type EvidenceGetInput struct {
	SessionID  learning.SessionID
	EvidenceID string
	Root       string
	Globs      []string
}

// EvidenceGetResult is what evidence_get reports.
type EvidenceGetResult struct {
	Content []byte
	Size    int
	Stale   bool
}

// EvidenceGet returns previously recorded evidence, scoped to the session
// that recorded it (requirement R8; security) and already redacted at
// write time.
func (w *WorkspaceService) EvidenceGet(in EvidenceGetInput) (EvidenceGetResult, error) {
	w.mu.Lock()
	known := w.scope[in.SessionID] != nil && w.scope[in.SessionID][in.EvidenceID]
	entry, hasRoot := w.evidenceRoots[in.SessionID][in.EvidenceID]
	w.mu.Unlock()
	if !known {
		for _, ev := range w.store.Replay(string(in.SessionID)) {
			if ev.Type != eventstore.EventEvidenceRecorded {
				continue
			}
			var p observationEvent
			if json.Unmarshal(ev.Payload, &p) == nil && p.EvidenceID == in.EvidenceID {
				known = true
				break
			}
		}
	}
	if !known {
		return EvidenceGetResult{}, ErrEvidenceOutOfScope
	}

	data, err := w.evidence.Get(in.EvidenceID)
	if err != nil {
		return EvidenceGetResult{}, err
	}

	if hasRoot {
		if in.Root != "" {
			root, err := workspace.AuthorizeRoot(in.Root)
			if err != nil || root.Path() != entry.root || !slices.Equal(in.Globs, entry.globs) {
				return EvidenceGetResult{}, ErrWorkspaceRootInvalid
			}
		}
		in.Root, in.Globs = entry.root, entry.globs
	}
	var stale bool
	if in.Root != "" {
		stale = true // An unavailable/replaced root cannot establish freshness.
		if root, rerr := workspace.AuthorizeRoot(in.Root); rerr == nil && (!hasRoot || root.Path() == entry.root) {
			var payload observationPayload
			if json.Unmarshal(data, &payload) == nil && payload.Fingerprint != "" {
				if fp, ferr := workspace.Fingerprint(root, in.Globs); ferr == nil {
					stale = fp != payload.Fingerprint
				}
			}
		}
	}

	return EvidenceGetResult{Content: data, Size: len(data), Stale: stale}, nil
}

// ErrNoWorkspaceBaseline is returned when a caller asks to reuse the
// workspace root/globs for (session, step) before any workspace_observe
// call established one (safe-check-executor Decision 2).
var ErrNoWorkspaceBaseline = errors.New("application: no workspace baseline established for this step yet")

// RootFor returns the root and globs most recently used to observe
// (sessionID, stepID), so a caller (safe-check-executor) can run a check
// against the same workspace without repeating those parameters.
func (w *WorkspaceService) RootFor(sessionID learning.SessionID, stepID learning.StepID) (root string, globs []string, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	entry, ok := w.baselines[baselineKey{session: sessionID, step: stepID}]
	if !ok {
		return "", nil, ErrNoWorkspaceBaseline
	}
	resolved, err := workspace.AuthorizeRoot(entry.root)
	if err != nil || resolved.Path() != entry.root {
		return "", nil, ErrWorkspaceRootInvalid
	}
	return entry.root, slices.Clone(entry.globs), nil
}

// PutEvidence stores raw content-addressed evidence and returns its ID,
// for a caller (safe-check-executor) that produces its own evidence
// payload shape but shares this service's evidence store.
func (w *WorkspaceService) PutEvidence(raw []byte) (string, error) {
	return w.evidence.Put(raw)
}

// RecordEvidence scopes evidenceID to sessionID, so a caller (safe-check-
// executor) that put its own evidence still has it enforced by
// EvidenceGet's per-session scope (requirement R8).
func (w *WorkspaceService) RecordEvidence(sessionID learning.SessionID, evidenceID string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.scope[sessionID] == nil {
		w.scope[sessionID] = map[string]bool{}
	}
	w.scope[sessionID][evidenceID] = true
	// Check evidence has the same root scope before and after replay.
	for _, ev := range w.store.Replay(string(sessionID)) {
		if ev.Type != eventstore.EventCheckExecuted {
			continue
		}
		var p observationEvent
		if json.Unmarshal(ev.Payload, &p) != nil || p.EvidenceID != evidenceID {
			continue
		}
		if entry, ok := w.baselines[baselineKey{sessionID, learning.StepID(p.StepID)}]; ok {
			if w.evidenceRoots[sessionID] == nil {
				w.evidenceRoots[sessionID] = map[string]baselineEntry{}
			}
			w.evidenceRoots[sessionID][evidenceID] = entry
		}
	}
}

func observationKind(isBaseline bool) string {
	if isBaseline {
		return "baseline"
	}
	return "diff"
}

// fingerprint prefers the ObserveResult's own fingerprint; a freshly
// established baseline computes one directly from its own hashes since
// workspace.Observe was never called to produce one.
func fingerprint(result workspace.ObserveResult, isBaseline bool) string {
	if !isBaseline {
		return result.Fingerprint
	}
	return workspace.FingerprintOf(result.Current)
}

// goDeclarationsFor parses every changed, still-present .go file up to
// maxDeclaredFiles, so structural criteria have AST facts to check against
// without executing the learner's code (requirement R6).
func goDeclarationsFor(root workspace.Root, changes []workspace.FileChange) map[string]workspace.GoDeclarations {
	decls := map[string]workspace.GoDeclarations{}
	for _, c := range changes {
		if len(decls) >= maxDeclaredFiles {
			break
		}
		if c.Change == workspace.ChangeRemoved || !strings.HasSuffix(c.Path, ".go") {
			continue
		}
		full, err := root.Resolve(c.Path)
		if err != nil {
			continue
		}
		decl, err := workspace.ParseGoFile(full)
		if err != nil {
			continue // not (yet) syntactically valid Go: no structural facts, not a failure
		}
		decls[c.Path] = decl
	}
	return decls
}

// diffTextFor requests a textual diff only for changed, still-present
// files, so DiffText's pathspec never grows to include the whole tree.
func diffTextFor(root workspace.Root, changes []workspace.FileChange) (string, bool, error) {
	paths := make([]string, 0, len(changes))
	for _, c := range changes {
		if c.Change != workspace.ChangeRemoved {
			paths = append(paths, c.Path)
		}
	}
	return workspace.DiffText(root, paths)
}
