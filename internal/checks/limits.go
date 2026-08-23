package checks

import "time"

// defaultTimeout applies when a CheckSpec declares no timeout.
const defaultTimeout = 30 * time.Second

// maxTimeout bounds any declared timeout, so a misauthored check can never
// hang the server indefinitely (non-functional requirement: limits must
// be configurable within safe ceilings).
const maxTimeout = 5 * time.Minute

// maxOutputBytes caps how much of stdout/stderr this package ever
// captures per stream, independent of redaction (non-functional
// requirement: outputs must be limited).
const maxOutputBytes = 256 * 1024

// maxParallel bounds how many checks this package's Executor runs at
// once, regardless of how many callers ask concurrently (requirement R4).
const maxParallel = 4

// truncationSuffix marks output cut off by maxOutputBytes.
var truncationSuffix = []byte("\n[TRUNCATED]")

func capOutput(b []byte) []byte {
	if len(b) <= maxOutputBytes {
		return b
	}
	out := make([]byte, 0, maxOutputBytes+len(truncationSuffix))
	out = append(out, b[:maxOutputBytes]...)
	out = append(out, truncationSuffix...)
	return out
}
