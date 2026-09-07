package eventstore

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

// ErrIncompatibleEventVersion is returned when a persisted event's
// SchemaVersion is newer than CurrentEventSchemaVersion and no Upcaster
// bridges the gap (compatibility requirement).
var ErrIncompatibleEventVersion = errors.New("eventstore: incompatible event schema version")

// ErrCorruptedLog is returned when a line in the middle of the log fails to
// parse. It is distinct from a truncated final line (requirement R6).
var ErrCorruptedLog = errors.New("eventstore: corrupted log entry")

// Upcaster transforms an event's raw JSON from one schema version to the
// next. Registering one for (Type, fromVersion) lets recovery bridge old
// events forward instead of rejecting them (compatibility requirement:
// "eventos possuem versão e upcasters explícitos").
type Upcaster func(raw json.RawMessage) (json.RawMessage, error)

type upcastKey struct {
	eventType   EventType
	fromVersion int
}

// UpcasterRegistry maps a versioned event shape to the function that
// upgrades it by exactly one schema version.
type UpcasterRegistry struct {
	upcasters map[upcastKey]Upcaster
}

// NewUpcasterRegistry returns an empty registry.
func NewUpcasterRegistry() *UpcasterRegistry {
	return &UpcasterRegistry{upcasters: map[upcastKey]Upcaster{}}
}

// Register adds an upcaster from fromVersion to fromVersion+1 for eventType.
func (r *UpcasterRegistry) Register(eventType EventType, fromVersion int, up Upcaster) {
	r.upcasters[upcastKey{eventType, fromVersion}] = up
}

func (r *UpcasterRegistry) upcast(ev *Event) error {
	for ev.SchemaVersion < CurrentEventSchemaVersion {
		up, ok := r.upcasters[upcastKey{ev.Type, ev.SchemaVersion}]
		if !ok {
			return fmt.Errorf("%w: %s at version %d", ErrIncompatibleEventVersion, ev.Type, ev.SchemaVersion)
		}
		raw, err := up(ev.Payload)
		if err != nil {
			return fmt.Errorf("upcasting %s from version %d: %w", ev.Type, ev.SchemaVersion, err)
		}
		ev.Payload = raw
		ev.SchemaVersion++
	}
	if ev.SchemaVersion > CurrentEventSchemaVersion {
		return fmt.Errorf("%w: %s at version %d, store supports %d", ErrIncompatibleEventVersion, ev.Type, ev.SchemaVersion, CurrentEventSchemaVersion)
	}
	return nil
}

// RecoveryResult reports what ReadEvents was able to recover. Events always
// holds every entry successfully parsed before any problem was hit
// (requirement R6: "sem perder o prefixo válido").
type RecoveryResult struct {
	Events []Event
	// Truncated is true when the last line of the log was incomplete
	// (e.g. the process was killed mid-write). This is not an error: the
	// valid prefix is intact and appending may safely resume.
	Truncated    bool
	validBytes   int64
	tail         []byte
	needsNewline bool
	// Err is non-nil only for a problem that is not a truncated final
	// line: corruption in the middle of the log, or an event whose schema
	// version has no path to CurrentEventSchemaVersion.
	Err error
}

// ReadEvents parses path as a JSONL event log, upcasting each event via
// registry as needed. It never returns fewer valid events than actually
// exist before the first problem encountered.
func ReadEvents(path string, registry *UpcasterRegistry) RecoveryResult {
	if registry == nil {
		registry = NewUpcasterRegistry()
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return RecoveryResult{}
		}
		return RecoveryResult{Err: err}
	}
	defer f.Close()

	reader := bufio.NewReaderSize(f, 1<<20)
	result := RecoveryResult{}
	revisions := map[string]uint64{}
	ids := map[string]bool{}
	requests := map[string]bool{}
	for i := 1; ; i++ {
		line, readErr := reader.ReadSlice('\n')
		if readErr != nil && readErr != io.EOF {
			result.Err = readErr
			return result
		}
		if len(line) == 0 {
			return result
		}
		if len(bytes.TrimSpace(line)) == 0 {
			result.validBytes += int64(len(line))
			if readErr == io.EOF {
				result.needsNewline = true
				return result
			}
			continue
		}
		var ev Event
		if err := json.Unmarshal(line, &ev); err != nil {
			if readErr == io.EOF {
				result.Truncated = true
				result.tail = bytes.Clone(line)
				return result
			}
			result.Err = fmt.Errorf("%w: line %d: %v", ErrCorruptedLog, i, err)
			return result
		}
		if err := registry.upcast(&ev); err != nil {
			result.Err = err
			return result
		}
		key := seenKey(ev.StreamID, ev.RequestID)
		if ev.StreamID == "" || ev.ID == "" || ev.Type == "" || ev.Revision != revisions[ev.StreamID]+1 || ids[ev.ID] || (ev.RequestID != "" && requests[key]) {
			result.Err = fmt.Errorf("%w: invalid envelope or revision at line %d", ErrCorruptedLog, i)
			return result
		}
		revisions[ev.StreamID] = ev.Revision
		ids[ev.ID] = true
		if ev.RequestID != "" {
			requests[key] = true
		}
		result.Events = append(result.Events, ev)
		result.validBytes += int64(len(line))
		if readErr == io.EOF {
			result.needsNewline = true
			return result
		}
	}
}
