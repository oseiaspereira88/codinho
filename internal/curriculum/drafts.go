package curriculum

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/oseiaspereira88/codinho/internal/checks"
	"github.com/oseiaspereira88/codinho/internal/security"
	"gopkg.in/yaml.v3"
	"io"
	"strings"
)

const MaxDraftBytes = 256 << 10

var ErrInvalidDraft = errors.New("curriculum: invalid draft")

type DraftValidation struct {
	Diagnostics []Diagnostic       `json:"diagnostics"`
	Editorial   []EditorialFinding `json:"editorial"`
}

// ParseDraft validates data only. It never reads fixture paths or executes checks.
func ParseDraft(data []byte) (Pack, DraftValidation, error) {
	var p Pack
	var result DraftValidation
	fail := func(detail string) (Pack, DraftValidation, error) {
		result.Diagnostics = append(result.Diagnostics, Diagnostic{File: "draft.yaml", Code: DiagMalformedPack, Detail: detail, Blocking: true})
		return Pack{}, result, ErrInvalidDraft
	}
	if len(data) == 0 || len(data) > MaxDraftBytes {
		return fail("draft must contain 1–262144 bytes")
	}
	if security.ContainsSecretShapedContent(data) {
		return fail("draft contains secret-shaped content")
	}
	var doc yaml.Node
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&doc); err != nil {
		return fail("malformed YAML")
	}
	var extra yaml.Node
	if err := decoder.Decode(&extra); err != io.EOF {
		return fail("submission requires exactly one YAML document")
	}
	if err := checkYAMLLimits(&doc, DefaultLimits); err != nil {
		return fail("YAML depth or alias limit exceeded")
	}
	if err := decodeStrict(&doc, &p); err != nil {
		return fail("unknown field or invalid authoring type")
	}
	if p.SchemaVersion != SchemaVersion {
		return fail("unsupported pack schema_version")
	}
	draftOnly := func(meta PublicationAuthoring) bool {
		return (meta.Status == "" || meta.Status == "draft") && meta.ReviewedBy == "" && !meta.Playtested
	}
	if !draftOnly(p.Publication) {
		return fail("submission cannot claim publication, review or playtest")
	}
	if strings.TrimSpace(p.ID) == "" || len(p.ID) > 200 || len(p.Challenges) == 0 || len(p.Challenges) > MaxTrackChallenges {
		return fail("pack needs an ID and 1–100 challenges")
	}
	p.Publication.Status = "draft"
	if p.Publication.Author == "" {
		p.Publication.Author = "host-agent"
	}
	p.File = "draft.yaml"
	for i := range p.Challenges {
		ch := &p.Challenges[i]
		if ch.SchemaVersion != SchemaVersion || strings.TrimSpace(ch.ID) == "" || len(ch.ID) > 200 {
			return fail("challenge requires a supported schema_version and an ID")
		}
		if !draftOnly(ch.Publication) {
			return fail("challenge cannot claim publication, review or playtest")
		}
		ch.Publication.Status = "draft"
		if ch.Publication.Author == "" {
			ch.Publication.Author = p.Publication.Author
		}
		for _, check := range ch.Checks {
			if _, err := checks.Resolve(checks.CheckSpec{ID: check.ID, Runner: check.Runner, Package: check.Package, TestPattern: check.TestPattern, Timeout: check.Timeout, NetworkApproved: check.Network}); err != nil {
				return fail(fmt.Sprintf("check %q does not resolve through the allowed runners", check.ID))
			}
		}
	}
	result.Diagnostics = Validate([]Pack{p})
	result.Diagnostics = append(result.Diagnostics, ValidatePublication([]Pack{p})...)
	result.Editorial = RunEditorialChecks([]Pack{p})
	for _, d := range result.Diagnostics {
		if d.Blocking || d.Code == DiagMissingReference {
			return Pack{}, result, ErrInvalidDraft
		}
	}
	for _, d := range result.Editorial {
		if d.Severity == SeverityBlocking {
			return Pack{}, result, ErrInvalidDraft
		}
	}
	return p, result, nil
}
func DraftCatalog(p Pack) *Catalog { return newCatalog([]Pack{p}) }
