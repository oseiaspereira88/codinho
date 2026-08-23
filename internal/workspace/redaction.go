package workspace

import (
	"path/filepath"
	"regexp"
	"strings"
)

// defaultExcludeNames lists base names never collected, even when a
// challenge's glob would otherwise match them (security requirement:
// exclude secrets, .env, credentials and paths sensitive by default).
var defaultExcludeNames = map[string]bool{
	".env":             true,
	".env.local":       true,
	".git":             true,
	"id_rsa":           true,
	"id_ed25519":       true,
	"credentials.json": true,
	"credentials.yaml": true,
	"credentials.yml":  true,
	"secrets.yaml":     true,
	"secrets.yml":      true,
	".netrc":           true,
	".npmrc":           true,
	".aws":             true,
}

// defaultExcludeSuffixes lists extensions never collected regardless of
// name (private key material and archives are never meaningfully
// diffable evidence, and keys must never be read at all).
var defaultExcludeSuffixes = []string{".pem", ".key", ".pfx", ".p12"}

// isExcluded reports whether rel (a root-relative, forward-slash path)
// must never be read, independent of any glob that would otherwise match
// it.
func isExcluded(rel string) bool {
	for seg := range strings.SplitSeq(rel, "/") {
		if defaultExcludeNames[seg] {
			return true
		}
	}
	base := filepath.Base(rel)
	if defaultExcludeNames[base] {
		return true
	}
	for _, suf := range defaultExcludeSuffixes {
		if strings.HasSuffix(base, suf) {
			return true
		}
	}
	return false
}

// secretPatterns match common secret shapes so they never leave this
// package in evidence, even when the excluded-file boundary above misses
// a secret embedded inside an otherwise-legitimate source file (security
// requirement: treat content as untrusted).
var secretPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(api[_-]?key|secret|password|token)\s*[:=]\s*["']?[A-Za-z0-9/+_\-]{12,}["']?`),
	regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
	regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----[\s\S]*?-----END [A-Z ]*PRIVATE KEY-----`),
	regexp.MustCompile(`ghp_[A-Za-z0-9]{36}`),
}

const redactedPlaceholder = "[REDACTED]"

// Redact replaces every secret-shaped substring of content with a fixed
// placeholder before it becomes evidence. Exported so other packages that
// capture process output (safe-check-executor) reuse the same secret
// patterns instead of duplicating them.
func Redact(content []byte) []byte { return redact(content) }

// redact is the unexported implementation Redact and this package's own
// callers share.
func redact(content []byte) []byte {
	out := content
	for _, p := range secretPatterns {
		out = p.ReplaceAll(out, []byte(redactedPlaceholder))
	}
	return out
}
