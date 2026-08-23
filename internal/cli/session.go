package cli

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/oseiaspereira88/codinho/internal/config"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
)

const sessionUsage = `Usage: codinho session <command>

Commands:
  inspect  Show one session's recorded event history
`

func runSession(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, sessionUsage)
		return exitUsage
	}
	switch args[0] {
	case "inspect":
		return runSessionInspect(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "codinho: unknown session subcommand %q\n\n", args[0])
		fmt.Fprint(stderr, sessionUsage)
		return exitUsage
	}
}

// sessionSummary is a read-only projection built directly from the
// session's raw event stream, not from session.Service's in-memory state:
// a session lives only in the memory of the codinho serve process that
// started it (session-orchestration-disclosure Decision 3), which is a
// different OS process from this CLI invocation. Reading the durable log
// is the only way an administrative command can inspect a session at all
// (requirement R3, "sem mutação pedagógica").
type sessionSummary struct {
	SessionID   string `json:"session_id"`
	Found       bool   `json:"found"`
	ChallengeID string `json:"challenge_id,omitempty"`
	EventCount  int    `json:"event_count"`
	Revision    uint64 `json:"revision"`
	FirstEvent  string `json:"first_event,omitempty"`
	FirstAt     string `json:"first_at,omitempty"`
	LastEvent   string `json:"last_event,omitempty"`
	LastAt      string `json:"last_at,omitempty"`
	State       string `json:"state,omitempty"`
}

func inspectSession(store *eventstore.Store, id string) sessionSummary {
	events := store.Replay(id)
	summary := sessionSummary{SessionID: id, Found: len(events) > 0, EventCount: len(events)}
	if len(events) == 0 {
		return summary
	}

	last := events[len(events)-1]
	summary.Revision = last.Revision
	summary.FirstEvent = string(events[0].Type)
	summary.FirstAt = events[0].RecordedAt
	summary.LastEvent = string(last.Type)
	summary.LastAt = last.RecordedAt

	for _, ev := range events {
		if ev.Type != eventstore.EventSessionStarted {
			continue
		}
		var payload struct {
			ChallengeID string `json:"challenge_id"`
		}
		if err := json.Unmarshal(ev.Payload, &payload); err == nil {
			summary.ChallengeID = payload.ChallengeID
		}
	}

	switch last.Type {
	case eventstore.EventSessionFinished:
		summary.State = "finished"
	case eventstore.EventSessionPaused:
		summary.State = "paused"
	default:
		summary.State = "active"
	}
	return summary
}

func runSessionInspect(args []string, stdout, stderr io.Writer) int {
	fs := newFlagSet("session inspect", stderr)
	jsonOut := fs.Bool("json", false, "print result as JSON")

	positional, flagArgs := extractPositional(args, nil)
	if err := fs.Parse(flagArgs); err != nil {
		return exitUsage
	}
	if len(positional) != 1 {
		fmt.Fprint(stderr, "Usage: codinho session inspect <session-id> [--json]\n")
		return exitUsage
	}
	id := positional[0]

	cfg := config.Load()
	store, err := openEventStore(cfg)
	if err != nil {
		fmt.Fprintf(stderr, "codinho: session inspect: %v\n", err)
		return exitError
	}
	defer store.Close()

	summary := inspectSession(store, id)
	if !summary.Found {
		if *jsonOut {
			_ = writeJSON(stdout, summary)
		} else {
			fmt.Fprintf(stdout, "session %s: not found\n", id)
		}
		return exitError
	}

	if *jsonOut {
		if err := writeJSON(stdout, summary); err != nil {
			fmt.Fprintf(stderr, "codinho: session inspect: %v\n", err)
			return exitError
		}
		return exitOK
	}
	fmt.Fprintf(stdout, "session_id: %s\n", summary.SessionID)
	fmt.Fprintf(stdout, "challenge_id: %s\n", summary.ChallengeID)
	fmt.Fprintf(stdout, "state: %s\n", summary.State)
	fmt.Fprintf(stdout, "revision: %d\n", summary.Revision)
	fmt.Fprintf(stdout, "events: %d (%s .. %s)\n", summary.EventCount, summary.FirstEvent, summary.LastEvent)
	return exitOK
}
