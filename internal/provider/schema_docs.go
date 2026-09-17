package provider

import (
	"regexp"
	"strings"

	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
)

// Documentation translation must not camel-case API enum values.
var translatedSpan = regexp.MustCompile(`<span[^>]*pulumi-lang-hcl="([^"]*)"[^>]*>.*?</span>`)

func literalDocs(s string) string {
	return translatedSpan.ReplaceAllStringFunc(s, func(span string) string {
		original := translatedSpan.FindStringSubmatch(span)[1]
		switch strings.Trim(strings.TrimSpace(original), "`") {
		case "html_css", "no_optimization", "post_process", "set_viewport":
			return original
		}
		return span
	})
}
func preserveAPILiterals(s *schema.PackageSpec) {
	fix := func(props map[string]schema.PropertySpec) {
		for k, v := range props {
			v.Description = literalDocs(v.Description)
			props[k] = v
		}
	}
	for k, v := range s.Resources {
		v.Description = literalDocs(v.Description)
		fix(v.Properties)
		fix(v.InputProperties)
		if v.StateInputs != nil {
			fix(v.StateInputs.Properties)
		}
		s.Resources[k] = v
	}
	for k, v := range s.Types {
		v.Description = literalDocs(v.Description)
		fix(v.Properties)
		s.Types[k] = v
	}
}
