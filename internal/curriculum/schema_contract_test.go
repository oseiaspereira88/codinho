package curriculum

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"gopkg.in/yaml.v3"
)

// The same documents exercise the external authoring schema and production Go
// loader. Graph identity and cross-field human checks remain semantic Go rules.
func TestPublicationSchemaCorpus(t *testing.T) {
	readSchema := func(name string) (*jsonschema.Schema, error) {
		data, err := os.ReadFile(filepath.Join("..", "..", "schemas", filepath.Base(name)))
		if err != nil {
			return nil, err
		}
		var schema jsonschema.Schema
		err = json.Unmarshal(data, &schema)
		return &schema, err
	}
	schema, err := readSchema("pack.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := schema.Resolve(&jsonschema.ResolveOptions{Loader: func(u *url.URL) (*jsonschema.Schema, error) { return readSchema(u.Path) }})
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile("../../testdata/catalog-quality/schema-corpus.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Name     string
		Valid    bool
		Document map[string]any
	}
	if err = json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			schemaErr := resolved.Validate(tc.Document)
			if (schemaErr == nil) != tc.Valid {
				t.Fatalf("schema valid=%v want=%v: %v", schemaErr == nil, tc.Valid, schemaErr)
			}
			dir := t.TempDir()
			if err = os.WriteFile(filepath.Join(dir, "manifest.yaml"), []byte("schema_version: 1\npacks: [pack.yaml]\n"), 0600); err != nil {
				t.Fatal(err)
			}
			yamlData, err := yaml.Marshal(tc.Document)
			if err != nil {
				t.Fatal(err)
			}
			if err = os.WriteFile(filepath.Join(dir, "pack.yaml"), yamlData, 0600); err != nil {
				t.Fatal(err)
			}
			catalog, diags, err := Load(dir, DefaultLimits)
			valid := err == nil && catalog != nil
			if valid != tc.Valid {
				t.Fatalf("loader valid=%v want=%v: err=%v diagnostics=%+v", valid, tc.Valid, err, diags)
			}
		})
	}
}
