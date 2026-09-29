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
		case "html_css", "no_optimization", "post_process", "set_viewport",
			"block", "beacon", "document", "stylesheet", "image", "image_set", "media",
			"font", "script", "text_track", "xhr", "fetch", "event_source",
			"manifest", "ping", "img", "other":
			return original
		}
		return span
	})
}
func postProcessSchema(s *schema.PackageSpec) {
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
		if strings.Contains(k, "RequestOverride:") {
			if property, ok := v.Properties["resourceTypes"]; ok {
				property.Items = &schema.TypeSpec{Ref: "#/types/html-css-to-image:index:RequestOverrideResourceType"}
				v.Properties["resourceTypes"] = property
			}
		}
		s.Types[k] = v
	}
	resources := []struct{ name, value string }{
		{"Beacon", "beacon"}, {"Document", "document"}, {"Stylesheet", "stylesheet"},
		{"Image", "image"}, {"ImageSet", "image_set"}, {"Media", "media"},
		{"Font", "font"}, {"Script", "script"}, {"TextTrack", "text_track"},
		{"Xhr", "xhr"}, {"Fetch", "fetch"}, {"EventSource", "event_source"},
		{"Manifest", "manifest"}, {"Ping", "ping"}, {"Img", "img"}, {"Other", "other"},
	}
	values := make([]schema.EnumValueSpec, len(resources))
	for i, resource := range resources {
		values[i] = schema.EnumValueSpec{Name: resource.name, Value: resource.value}
	}
	s.Types["html-css-to-image:index:RequestOverrideResourceType"] = schema.ComplexTypeSpec{
		ObjectTypeSpec: schema.ObjectTypeSpec{Type: "string", Description: "Browser resource type matched by a request override rule."},
		Enum:           values,
	}
}
