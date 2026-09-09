package application

import (
	"encoding/json"
	"github.com/oseiaspereira88/codinho/internal/eventstore"
)

// evidenceProvenance follows immutable producer events and their owning starts.
// Unknown or ambiguous origins cannot prove reviewed publication.
func (p *ProgressService) evidenceProvenance(id string) string {
	starts := map[string]string{}
	provenance := ""
	for _, ev := range p.store.ReplayAll() {
		var data struct {
			EvidenceID string `json:"evidence_id"`
			Provenance string `json:"content_provenance"`
		}
		if json.Unmarshal(ev.Payload, &data) != nil {
			continue
		}
		if ev.Type == eventstore.EventSessionStarted {
			if data.Provenance != "published" && data.Provenance != "draft" {
				data.Provenance = "legacy_unreviewed"
			}
			starts[ev.StreamID] = data.Provenance
		}
		if ev.Type != eventstore.EventEvidenceRecorded && ev.Type != eventstore.EventObservationRecorded && ev.Type != eventstore.EventCheckExecuted {
			continue
		}
		if data.EvidenceID != id {
			continue
		}
		origin := starts[ev.StreamID]
		if origin == "" {
			origin = "legacy_unreviewed"
		}
		if provenance != "" && provenance != origin {
			return "legacy_unreviewed"
		}
		provenance = origin
	}
	if provenance == "" {
		return "legacy_unreviewed"
	}
	return provenance
}

func sessionProvenance(store *eventstore.Store, id string) string {
	for _, ev := range store.Replay(id) {
		if ev.Type != eventstore.EventSessionStarted {
			continue
		}
		var data struct {
			Provenance string `json:"content_provenance"`
		}
		if json.Unmarshal(ev.Payload, &data) == nil && (data.Provenance == "published" || data.Provenance == "draft") {
			return data.Provenance
		}
		break
	}
	return "legacy_unreviewed"
}
