package provider

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// Validate checked-in YAML against generated property names and required fields.
// Examples must not silently lag behind nested schema or token changes.
func TestYAMLExamplesMatchSchema(t *testing.T) {
	var pkg struct {
		Provider struct {
			InputProperties map[string]json.RawMessage `json:"inputProperties"`
			RequiredInputs  []string                   `json:"requiredInputs"`
		}
		Resources map[string]struct {
			InputProperties map[string]json.RawMessage `json:"inputProperties"`
			RequiredInputs  []string                   `json:"requiredInputs"`
		}
		Types map[string]struct {
			Properties map[string]json.RawMessage
			Required   []string
		}
		Functions map[string]json.RawMessage
	}
	if err := json.Unmarshal(Schema, &pkg); err != nil {
		t.Fatal(err)
	}
	var validate func(map[string]any, map[string]json.RawMessage, []string, string)
	validate = func(values map[string]any, fields map[string]json.RawMessage, required []string, where string) {
		for _, k := range required {
			if _, ok := values[k]; !ok {
				t.Errorf("%s: missing %s", where, k)
			}
		}
		for k, v := range values {
			raw, ok := fields[k]
			if !ok {
				t.Errorf("%s: unknown property %s", where, k)
				continue
			}
			var spec struct {
				Ref string `json:"$ref"`
			}
			if err := json.Unmarshal(raw, &spec); err != nil {
				t.Fatal(err)
			}
			if nested, ok := v.(map[string]any); ok && strings.HasPrefix(spec.Ref, "#/types/") {
				typ := pkg.Types[strings.TrimPrefix(spec.Ref, "#/types/")]
				validate(nested, typ.Properties, typ.Required, where+"."+k)
			}
		}
	}
	files, err := filepath.Glob("../../examples/*/Pulumi.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) < 7 {
		t.Fatal("missing examples")
	}
	for _, file := range files {
		t.Run(filepath.Base(filepath.Dir(file)), func(t *testing.T) {
			data, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			var example struct {
				Resources map[string]struct {
					Type       string
					Properties map[string]any
				}
				Variables map[string]map[string]struct{ Function string }
			}
			if err := yaml.Unmarshal(data, &example); err != nil {
				t.Fatal(err)
			}
			for name, r := range example.Resources {
				if r.Type == "pulumi:providers:html-css-to-image" {
					validate(r.Properties, pkg.Provider.InputProperties, pkg.Provider.RequiredInputs, name)
					continue
				}
				spec, ok := pkg.Resources[r.Type]
				if !ok {
					t.Errorf("unknown token %s", r.Type)
					continue
				}
				validate(r.Properties, spec.InputProperties, spec.RequiredInputs, name)
			}
			for _, v := range example.Variables {
				if call, ok := v["fn::invoke"]; ok && pkg.Functions[call.Function] == nil {
					t.Errorf("unknown function %s", call.Function)
				}
			}
		})
	}
}
