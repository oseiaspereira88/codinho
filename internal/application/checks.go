package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/oseiaspereira88/codinho/internal/checks"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/learning"
	"github.com/oseiaspereira88/codinho/internal/workspace"
)

// ErrCheckNotFound is returned when check_id does not name a check
// declared by the session's fixed challenge (requirement R1).
var ErrCheckNotFound = errors.New("application: check not found for this session's challenge")

// checkPayload is what this package stores as evidence content for a
// check_run. Its "kind" and "fingerprint" fields are recognized generically
// by WorkspaceService.EvidenceGet's staleness check (requirement R9) and by
// session.Service's structural-criterion override (requirement R10) — both
// only ever look for these two fields by name, never at this type.
type checkPayload struct {
	Kind        string `json:"kind"` // always "check"
	CheckID     string `json:"check_id"`
	Outcome     string `json:"outcome"`
	ExitCode    int    `json:"exit_code"`
	Stdout      string `json:"stdout,omitempty"`
	Stderr      string `json:"stderr,omitempty"`
	Fingerprint string `json:"fingerprint"`
	// NetworkApproved records whether this run had network access,
	// decided only by the challenge's own authored Network field — never
	// a caller-supplied override — so the durable record makes an
	// otherwise invisible decision auditable (security-privacy-hardening
	// requirement R7: "exigir aprovação visível quando necessária").
	NetworkApproved bool `json:"network_approved"`
}

// ChecksService orchestrates safe-check-executor: resolving check_id
// against the session's fixed challenge (requirement R1), executing it
// inside the workspace root already established by workspace_observe
// (Decision 2), and recording the outcome as scoped, redacted evidence.
type ChecksService struct {
	mu        sync.Mutex
	store     *eventstore.Store
	sessions  *SessionService
	workspace *WorkspaceService
	executor  *checks.Executor
}

// NewChecksService wires a ChecksService to the shared event store and to
// the session/workspace services it reads context from.
func NewChecksService(store *eventstore.Store, sessions *SessionService, ws *WorkspaceService) *ChecksService {
	return &ChecksService{store: store, sessions: sessions, workspace: ws, executor: checks.NewExecutor()}
}

// CheckRunInput is one check_run call. It carries no free command or
// workspace parameter (requirement R9): check_id is resolved from the
// session's own fixed catalog (requirement R1), and the workspace root is
// whatever workspace_observe already established for the active step
// (Decision 2).
type CheckRunInput struct {
	SessionID        learning.SessionID
	CheckID          string
	ExpectedRevision uint64
	RequestID        string
}

// CheckRunResult is what check_run reports.
type CheckRunResult struct {
	Outcome         string
	EvidenceID      string
	Fingerprint     string
	Revision        uint64
	NetworkApproved bool
}

// Run resolves CheckID against the session's active challenge, executes
// it inside the workspace root/globs already observed for the active
// step, and records a redacted, size-capped, fingerprinted evidence blob
// scoped to this session (requirements R1, R6, R7, R8, R9).
func (c *ChecksService) Run(in CheckRunInput) (CheckRunResult, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ev, ok := c.store.Request(string(in.SessionID), in.RequestID); ok {
		var previous struct {
			CheckID    string `json:"check_id"`
			EvidenceID string `json:"evidence_id"`
		}
		if ev.Type != eventstore.EventCheckExecuted || json.Unmarshal(ev.Payload, &previous) != nil || previous.CheckID != in.CheckID {
			return CheckRunResult{}, eventstore.ErrRevisionConflict
		}
		evidence, err := c.workspace.EvidenceGet(EvidenceGetInput{SessionID: in.SessionID, EvidenceID: previous.EvidenceID})
		if err != nil {
			return CheckRunResult{}, err
		}
		var payload checkPayload
		if err := json.Unmarshal(evidence.Content, &payload); err != nil {
			return CheckRunResult{}, err
		}
		return CheckRunResult{Outcome: payload.Outcome, EvidenceID: previous.EvidenceID, Fingerprint: payload.Fingerprint, Revision: ev.Revision, NetworkApproved: payload.NetworkApproved}, nil
	}
	authored, stepID, err := c.sessions.ActiveChecks(in.SessionID)
	if err != nil {
		return CheckRunResult{}, err
	}
	var found *curriculum.CheckAuthoring
	for i := range authored {
		if authored[i].ID == in.CheckID {
			found = &authored[i]
			break
		}
	}
	if found == nil {
		return CheckRunResult{}, fmt.Errorf("%w: %q", ErrCheckNotFound, in.CheckID)
	}

	rootPath, globs, err := c.workspace.RootFor(in.SessionID, stepID)
	if err != nil {
		return CheckRunResult{}, err
	}
	root, err := workspace.AuthorizeRoot(rootPath)
	if err != nil {
		return CheckRunResult{}, fmt.Errorf("%w: %w", ErrWorkspaceRootInvalid, err)
	}

	resolved, err := checks.Resolve(checks.CheckSpec{
		ID: found.ID, Runner: found.Runner, Package: found.Package,
		TestPattern: found.TestPattern, Timeout: found.Timeout, NetworkApproved: found.Network,
	})
	if err != nil {
		return CheckRunResult{}, err
	}

	result, err := c.executor.Execute(context.Background(), resolved, root)
	if err != nil {
		return CheckRunResult{}, err
	}

	fingerprint, err := workspace.Fingerprint(root, globs)
	if err != nil {
		return CheckRunResult{}, err
	}

	payload := checkPayload{
		Kind: "check", CheckID: found.ID, Outcome: string(result.Outcome), ExitCode: result.ExitCode,
		Stdout: string(result.Stdout), Stderr: string(result.Stderr), Fingerprint: fingerprint,
		NetworkApproved: found.Network,
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return CheckRunResult{}, err
	}
	evidenceID, err := c.workspace.PutEvidence(raw)
	if err != nil {
		return CheckRunResult{}, err
	}

	ev, err := c.store.Append(string(in.SessionID), in.ExpectedRevision, in.RequestID, eventstore.EventCheckExecuted, map[string]any{
		"step_id":     string(stepID),
		"check_id":    found.ID,
		"outcome":     string(result.Outcome),
		"evidence_id": evidenceID,
		"fingerprint": fingerprint,
	})
	if err != nil {
		return CheckRunResult{}, err
	}
	c.workspace.RecordEvidence(in.SessionID, evidenceID)

	return CheckRunResult{
		Outcome: string(result.Outcome), EvidenceID: evidenceID, Fingerprint: fingerprint, Revision: ev.Revision,
		NetworkApproved: found.Network,
	}, nil
}
