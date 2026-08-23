package application

import (
	"encoding/json"
	"errors"
	"fmt"
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

	mu        sync.Mutex
	baselines map[baselineKey]workspace.Baseline
	scope     map[learning.SessionID]map[string]bool
}

// NewWorkspaceService wires a WorkspaceService to the shared event and
// evidence stores.
func NewWorkspaceService(store *eventstore.Store, evidenceStore *evidence.Store) *WorkspaceService {
	return &WorkspaceService{
		store:     store,
		evidence:  evidenceStore,
		baselines: map[baselineKey]workspace.Baseline{},
		scope:     map[learning.SessionID]map[string]bool{},
	}
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
	root, err := workspace.AuthorizeRoot(in.Root)
	if err != nil {
		return ObserveResult{}, fmt.Errorf("%w: %w", ErrWorkspaceRootInvalid, err)
	}

	key := baselineKey{session: in.SessionID, step: in.StepID}
	w.mu.Lock()
	baseline, hasBaseline := w.baselines[key]
	w.mu.Unlock()

	var (
		result      workspace.ObserveResult
		newBaseline workspace.Baseline
		isBaseline  bool
	)
	if hasBaseline {
		result, err = workspace.Observe(root, in.Globs, baseline)
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

	ev, err := w.store.Append(string(in.SessionID), in.ExpectedRevision, in.RequestID, eventstore.EventObservationRecorded, map[string]any{
		"step_id":     string(in.StepID),
		"evidence_id": evidenceID,
		"kind":        payload.Kind,
		"fingerprint": payload.Fingerprint,
	})
	if err != nil {
		return ObserveResult{}, err
	}

	w.mu.Lock()
	if isBaseline {
		w.baselines[key] = newBaseline
	}
	if w.scope[in.SessionID] == nil {
		w.scope[in.SessionID] = map[string]bool{}
	}
	w.scope[in.SessionID][evidenceID] = true
	w.mu.Unlock()

	learningEvidence, err := learning.NewEvidence(learning.EvidenceID(evidenceID), learning.EvidenceKindDiff, evidenceID)
	if err != nil {
		return ObserveResult{}, err
	}
	observation := learning.Observation{StepID: in.StepID, Evidence: learningEvidence, Observed: time.Now().UTC()}

	return ObserveResult{
		Observation: observation,
		EvidenceID:  evidenceID,
		Baseline:    isBaseline,
		Changes:     result.Changes,
		FromGit:     result.Current.FromGit,
		Dirty:       result.Current.Dirty,
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
	w.mu.Unlock()
	if !known {
		return EvidenceGetResult{}, ErrEvidenceOutOfScope
	}

	data, err := w.evidence.Get(in.EvidenceID)
	if err != nil {
		return EvidenceGetResult{}, err
	}

	var stale bool
	if in.Root != "" {
		if root, rerr := workspace.AuthorizeRoot(in.Root); rerr == nil {
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
