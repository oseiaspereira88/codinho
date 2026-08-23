package curriculum

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Limits bound how permissive the YAML parser is, so an authored (or
// compromised) pack cannot exhaust memory or CPU via deep nesting or alias
// expansion (PROJECT.md §2 security requirement).
type Limits struct {
	MaxFileBytes int64
	MaxNodeDepth int
	MaxAliases   int
}

// DefaultLimits are generous enough for hand-authored curriculum content and
// tight enough to reject adversarial YAML.
var DefaultLimits = Limits{
	MaxFileBytes: 1 << 20, // 1 MiB
	MaxNodeDepth: 32,
	MaxAliases:   64,
}

// Load reads packs/manifest.yaml under dir, loads every listed pack within
// Limits, and returns the validated, immutable Catalog together with every
// diagnostic collected along the way. Load never partially trusts a pack: a
// blocking diagnostic (see Diagnostic.Blocking) prevents catalog
// construction so a broken pack can never be materialized silently.
func Load(dir string, limits Limits) (*Catalog, []Diagnostic, error) {
	packs, diags, err := LoadPacks(dir, limits)
	if err != nil {
		return nil, nil, err
	}
	for _, d := range diags {
		if d.Blocking {
			return nil, diags, nil
		}
	}
	return newCatalog(packs), diags, nil
}

// LoadPacks reads and validates every pack under dir the same way Load
// does, but returns the raw []Pack instead of an assembled Catalog. It
// exists for callers that need authored data Catalog's public API does
// not expose in bulk — e.g. catalog-authoring-quality's editorial checks,
// which need Pack.Relations and per-pack File attribution across every
// item, not just one Catalog.Challenge(id) at a time. It returns a nil
// []Pack (not an error) when a blocking diagnostic exists, matching
// Load's own "never partially trust a pack" contract.
func LoadPacks(dir string, limits Limits) ([]Pack, []Diagnostic, error) {
	manifestPath := filepath.Join(dir, "manifest.yaml")
	manifestNode, err := readLimitedYAML(manifestPath, limits)
	if err != nil {
		return nil, nil, fmt.Errorf("reading manifest: %w", err)
	}
	var manifest PackManifest
	if err := manifestNode.Decode(&manifest); err != nil {
		return nil, nil, fmt.Errorf("decoding manifest: %w", err)
	}
	if manifest.SchemaVersion != SchemaVersion {
		return nil, []Diagnostic{{
			File: manifestPath, Code: DiagIncompatibleSchemaVersion,
			Detail: fmt.Sprintf("manifest schema_version %d, loader supports %d", manifest.SchemaVersion, SchemaVersion),
		}}, nil
	}

	// Sort pack filenames before loading so filesystem enumeration order
	// never affects the result (non-functional requirement).
	names := append([]string{}, manifest.Packs...)
	sort.Strings(names)

	var packs []Pack
	var diags []Diagnostic
	for _, name := range names {
		path, err := confinedPath(dir, name)
		if err != nil {
			diags = append(diags, Diagnostic{File: name, Code: DiagPathEscapesPack, Detail: err.Error(), Blocking: true})
			continue
		}
		node, err := readLimitedYAML(path, limits)
		if err != nil {
			diags = append(diags, Diagnostic{File: name, Code: DiagUnreadablePack, Detail: err.Error(), Blocking: true})
			continue
		}
		var pack Pack
		if err := node.Decode(&pack); err != nil {
			diags = append(diags, Diagnostic{File: name, Code: DiagMalformedPack, Detail: err.Error(), Blocking: true})
			continue
		}
		pack.File = name
		if pack.SchemaVersion != SchemaVersion {
			diags = append(diags, Diagnostic{
				File: name, Item: pack.ID, Field: "schema_version", Code: DiagIncompatibleSchemaVersion,
				Detail:   fmt.Sprintf("pack schema_version %d, loader supports %d", pack.SchemaVersion, SchemaVersion),
				Blocking: true,
			})
			continue
		}
		packs = append(packs, pack)
	}

	diags = append(diags, Validate(packs)...)
	for _, d := range diags {
		if d.Blocking {
			return nil, diags, nil
		}
	}

	return packs, diags, nil
}

// readLimitedYAML reads path, rejects it if it exceeds limits.MaxFileBytes,
// parses it into a yaml.Node, and rejects the node if it exceeds
// limits.MaxNodeDepth or limits.MaxAliases before any Decode call ever runs
// (Decode is where alias expansion and deep recursion would otherwise cost
// memory and CPU).
func readLimitedYAML(path string, limits Limits) (*yaml.Node, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.Size() > limits.MaxFileBytes {
		return nil, fmt.Errorf("%s exceeds max size of %d bytes", path, limits.MaxFileBytes)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	if err := checkYAMLLimits(&doc, limits); err != nil {
		return nil, err
	}
	return &doc, nil
}

func checkYAMLLimits(node *yaml.Node, limits Limits) error {
	aliases := 0
	var walk func(n *yaml.Node, depth int) error
	walk = func(n *yaml.Node, depth int) error {
		if depth > limits.MaxNodeDepth {
			return fmt.Errorf("yaml nesting exceeds max depth of %d", limits.MaxNodeDepth)
		}
		if n.Kind == yaml.AliasNode {
			aliases++
			if aliases > limits.MaxAliases {
				return fmt.Errorf("yaml alias count exceeds max of %d", limits.MaxAliases)
			}
		}
		for _, c := range n.Content {
			if err := walk(c, depth+1); err != nil {
				return err
			}
		}
		return nil
	}
	return walk(node, 0)
}

// confinedPath resolves name relative to dir and rejects any path that
// would escape dir, so a manifest entry can never read a fixture or file
// outside the pack (PROJECT.md §2 security requirement).
func confinedPath(dir, name string) (string, error) {
	joined := filepath.Join(dir, name)
	rel, err := filepath.Rel(dir, joined)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return "", fmt.Errorf("pack path %q escapes catalog directory", name)
	}
	return joined, nil
}
