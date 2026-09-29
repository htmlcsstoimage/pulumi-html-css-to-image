package provider

import (
	"encoding/json"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
)

func TestRequestOverrideSchema(t *testing.T) {
	const enumRef = "#/types/html-css-to-image:index:RequestOverrideResourceType"
	var spec struct {
		Types map[string]struct {
			Type       string                   `json:"type"`
			Enum       []struct{ Value string } `json:"enum"`
			Properties map[string]struct {
				Items struct {
					Ref string `json:"$ref"`
				} `json:"items"`
			} `json:"properties"`
		} `json:"types"`
	}
	if err := json.Unmarshal(Schema, &spec); err != nil {
		t.Fatal(err)
	}
	resourceType := spec.Types["html-css-to-image:index:RequestOverrideResourceType"]
	if resourceType.Type != "string" {
		t.Fatal("resource type enum must serialize as a string")
	}
	want := []string{"beacon", "document", "stylesheet", "image", "image_set", "media", "font", "script", "text_track", "xhr", "fetch", "event_source", "manifest", "ping", "img", "other"}
	if len(resourceType.Enum) != len(want) {
		t.Fatalf("resource type enum has %d values, want %d", len(resourceType.Enum), len(want))
	}
	for i, value := range want {
		if resourceType.Enum[i].Value != value {
			t.Fatalf("resource type enum value %d = %q, want %q", i, resourceType.Enum[i].Value, value)
		}
	}
	for _, name := range []string{
		"ImageHtmlCssRequestOverride", "ImageUrlRequestOverride", "TemplateRequestOverride",
		"OgConfigDefaultOptionsRequestOverride", "getTemplateRequestOverride",
	} {
		key := "html-css-to-image:index/" + name + ":" + name
		if got := spec.Types[key].Properties["resourceTypes"].Items.Ref; got != enumRef {
			t.Errorf("%s resourceTypes item ref = %q, want %q", key, got, enumRef)
		}
	}
}

func TestRequestOverridesBridge(t *testing.T) {
	rules := []any{map[string]any{
		"action":        "block",
		"url":           "*://analytics.example.com/*",
		"resourceTypes": []any{"script", "image_set"},
	}}
	for _, tc := range []struct {
		name string
		data map[string]any
	}{
		{"ImageHtmlCss", map[string]any{"html": "<h1>Hello</h1>", "requestOverrides": rules}},
		{"ImageUrl", map[string]any{"url": "https://example.com", "requestOverrides": rules}},
		{"Template", map[string]any{"html": "<h1>{{title}}</h1>", "requestOverrides": rules}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p, api := renderingServer(t)
			u := urn(tc.name)
			in := checked(t, p, u, nil, props(tc.data))
			created, err := p.Create(ctx, plugin.CreateRequest{URN: u, Properties: in})
			if err != nil {
				t.Fatal(err)
			}
			api.mu.Lock()
			var saved map[string]any
			if tc.name == "Template" {
				saved = api.templates[0]
			} else {
				saved = api.images[string(created.ID)]
			}
			api.mu.Unlock()
			assertRequestOverride(t, saved["request_overrides"])
			assertPulumiRequestOverride(t, created.Properties["requestOverrides"])
			unchanged(t, p, u, created.ID, in, created.Properties)
		})
	}

	p, api := server(t)
	u := urn("OgConfig")
	in := checked(t, p, u, nil, props(map[string]any{
		"name": "Website cards", "baseUrl": "https://example.com", "configType": "html_css",
		"defaultOptions": map[string]any{"requestOverrides": rules},
	}))
	created, err := p.Create(ctx, plugin.CreateRequest{URN: u, Properties: in})
	if err != nil {
		t.Fatal(err)
	}
	remote, _, _ := api.OGSnapshot()
	if remote.DefaultOptions == nil || len(remote.DefaultOptions.RequestOverrides) != 1 {
		t.Fatal("OG default options lost request overrides")
	}
	assertPulumiRequestOverride(t, created.Properties["defaultOptions"].ObjectValue()["requestOverrides"])
	unchanged(t, p, u, created.ID, in, created.Properties)
}

func assertRequestOverride(t *testing.T, value any) {
	t.Helper()
	rules, ok := value.([]any)
	if !ok || len(rules) != 1 {
		t.Fatalf("missing API request overrides: %#v", value)
	}
	rule, ok := rules[0].(map[string]any)
	if !ok || rule["action"] != "block" || rule["url"] != "*://analytics.example.com/*" {
		t.Fatalf("incorrect API request override: %#v", rules[0])
	}
	resources, ok := rule["resource_types"].([]any)
	if !ok || len(resources) != 2 || resources[0] != "script" || resources[1] != "image_set" {
		t.Fatalf("incorrect API resource types: %#v", rule["resource_types"])
	}
}

func assertPulumiRequestOverride(t *testing.T, value resource.PropertyValue) {
	t.Helper()
	rules := value.ArrayValue()
	if len(rules) != 1 {
		t.Fatal("missing Pulumi request overrides")
	}
	rule := rules[0].ObjectValue()
	if rule["action"].StringValue() != "block" || rule["url"].StringValue() != "*://analytics.example.com/*" {
		t.Fatal("Pulumi rule changed")
	}
	resources := rule["resourceTypes"].ArrayValue()
	if len(resources) != 2 || resources[0].StringValue() != "script" || resources[1].StringValue() != "image_set" {
		t.Fatal("Pulumi resource type strings changed")
	}
}
