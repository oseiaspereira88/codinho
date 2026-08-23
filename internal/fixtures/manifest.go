package fixtures

// ManifestFile records one materialized fixture file's identity, so a
// learner or grader can verify what was written matches what the
// challenge authored (requirement R6).
type ManifestFile struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
	Size   int64  `json:"size"`
}

// Manifest is written as manifest.json at the destination root after a
// successful Materialize call.
type Manifest struct {
	ChallengeID string         `json:"challenge_id"`
	CreatedAt   string         `json:"created_at"`
	Files       []ManifestFile `json:"files"`
}
