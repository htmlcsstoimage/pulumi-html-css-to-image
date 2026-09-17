package provider

import (
	"encoding/json"
	"testing"

	pf "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfbridge"
	shim "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfshim"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
)

func urn(name string) resource.URN {
	return resource.URN("urn:pulumi:test::test::html-css-to-image:index:" + name + "::test")
}
func renderingServer(t *testing.T) (plugin.Provider, *renderFixture) {
	t.Helper()
	api := newRenderFixture(t)
	p, err := pf.NewProvider(ctx, Info("0.0.1"), pf.ProviderMetadata{PackageSchema: Schema})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { p.Close() })
	_, err = p.Configure(ctx, plugin.ConfigureRequest{Inputs: props(map[string]any{"apiId": "test-id", "apiKey": "test-key", "baseUrl": api.URL})})
	if err != nil {
		t.Fatal(err)
	}
	return p, api
}
func secretMap(t *testing.T, v resource.PropertyValue, k, want string) {
	t.Helper()
	if !v.IsSecret() {
		t.Fatal("map lost sensitivity")
	}
	if v.SecretValue().Element.ObjectValue()[resource.PropertyKey(k)].StringValue() != want {
		t.Fatal("secret map value changed")
	}
}
func TestNewResourceSchema(t *testing.T) {
	var s struct {
		Resources map[string]struct {
			Properties map[string]struct {
				Type   string
				Secret bool
			}
			RequiredInputs []string `json:"requiredInputs"`
		}
		Functions map[string]json.RawMessage
	}
	if err := json.Unmarshal(Schema, &s); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"OgConfig", "StorageDestination", "Template", "ImageHtmlCss", "ImageUrl", "ImageTemplated"} {
		if _, ok := s.Resources["html-css-to-image:index:"+name]; !ok {
			t.Fatal("missing resource", name)
		}
	}
	for _, v := range []struct{ resource, field string }{{"Template", "version"}, {"ImageTemplated", "templateVersion"}, {"ImageTemplated", "resolvedTemplateVersion"}, {"OgConfig", "templateVersion"}, {"ImageHtmlCss", "ogConfigContentVersion"}} {
		if s.Resources["html-css-to-image:index:"+v.resource].Properties[v.field].Type != "string" {
			t.Fatal("version not lossless", v)
		}
	}
	if !s.Resources["html-css-to-image:index:ImageUrl"].Properties["headers"].Secret || !s.Resources["html-css-to-image:index:ImageTemplated"].Properties["templateValues"].Secret {
		t.Fatal("missing sensitivity")
	}
	for _, k := range s.Resources["html-css-to-image:index:ImageTemplated"].RequiredInputs {
		if k == "templateVersion" {
			t.Fatal("version must remain optional")
		}
	}
	if s.Functions["html-css-to-image:index:getAwsStorageExternalId"] == nil {
		t.Fatal("missing external ID lookup")
	}
}
func TestImageBridgeLifecycles(t *testing.T) {
	for _, name := range []string{"ImageHtmlCss", "ImageUrl", "ImageTemplated"} {
		t.Run(name, func(t *testing.T) {
			p, api := renderingServer(t)
			u := urn(name)
			data := map[string]any{"html": "<h1>Hello 😀</h1>"}
			if name == "ImageUrl" {
				data = map[string]any{"url": "https://example.com", "headers": map[string]any{"Authorization": "secret 😀"}}
			}
			if name == "ImageTemplated" {
				data = map[string]any{"templateId": "t-test", "templateValues": `{"name":"😀","n":9007199254740993}`}
				api.templates = []map[string]any{{"id": "t-test", "version": int64(9007199254740993), "html": "{{name}}", "template_type": "html_css", "created_at": "now"}}
			}
			in := checked(t, p, u, nil, props(data))
			preview, err := p.Create(ctx, plugin.CreateRequest{URN: u, Properties: in, Preview: true})
			if err != nil || preview.ID != "" {
				t.Fatal("preview", err)
			}
			if api.creates != 0 {
				t.Fatal("preview wrote")
			}
			created, err := p.Create(ctx, plugin.CreateRequest{URN: u, Properties: in})
			if err != nil {
				t.Fatal(err)
			}
			if name == "ImageUrl" {
				secretMap(t, created.Properties["headers"], "Authorization", "secret 😀")
			}
			if name == "ImageTemplated" {
				secret(t, created.Properties["templateValues"], data["templateValues"].(string))
				if created.Properties["resolvedTemplateVersion"].StringValue() != "9007199254740993" || !created.Properties["templateVersion"].IsNull() {
					t.Fatal("version changed or pin invented")
				}
			}
			unchanged(t, p, u, created.ID, in, created.Properties)
			read, err := p.Read(ctx, plugin.ReadRequest{URN: u, ID: created.ID, Inputs: in, State: created.Properties})
			if err != nil {
				t.Fatal(err)
			}
			unchanged(t, p, u, created.ID, read.Inputs, read.Outputs)
			if name == "ImageUrl" {
				secretMap(t, read.Outputs["headers"], "Authorization", "secret 😀")
			}
			imported, err := p.Read(ctx, plugin.ReadRequest{URN: u, ID: created.ID})
			if err != nil {
				t.Fatal(err)
			}
			unchanged(t, p, u, created.ID, imported.Inputs, imported.Outputs)
			next := in.Copy()
			next["format"] = resource.NewStringProperty("webp")
			next = checked(t, p, u, in, next)
			diff, err := p.Diff(ctx, plugin.DiffRequest{URN: u, ID: created.ID, OldInputs: in, OldOutputs: read.Outputs, NewInputs: next})
			if err != nil || len(diff.ReplaceKeys) == 0 {
				t.Fatalf("image change must replace: %v %+v", err, diff)
			}
			_, err = p.Delete(ctx, plugin.DeleteRequest{URN: u, ID: created.ID, Inputs: in, Outputs: read.Outputs})
			if err != nil {
				t.Fatal(err)
			}
			missing, err := p.Read(ctx, plugin.ReadRequest{URN: u, ID: created.ID, Inputs: in, State: read.Outputs})
			if err != nil || missing.ID != "" {
				t.Fatal("delete/404", err)
			}
		})
	}
}
func TestTemplateBridgeVersions(t *testing.T) {
	p, api := renderingServer(t)
	u := urn("Template")
	in := checked(t, p, u, nil, props(map[string]any{"html": "<h1>{{name}}</h1>", "googleFonts": []any{"Roboto", "Open Sans"}}))
	created, err := p.Create(ctx, plugin.CreateRequest{URN: u, Properties: in})
	if err != nil {
		t.Fatal(err)
	}
	unchanged(t, p, u, created.ID, in, created.Properties)
	next := in.Copy()
	next["msDelay"] = resource.NewNumberProperty(0)
	next = checked(t, p, u, in, next)
	updated, err := p.Update(ctx, plugin.UpdateRequest{URN: u, ID: created.ID, OldInputs: in, OldOutputs: created.Properties, NewInputs: next})
	if err != nil {
		t.Fatal(err)
	}
	removed := next.Copy()
	delete(removed, "msDelay")
	removed = checked(t, p, u, next, removed)
	updated, err = p.Update(ctx, plugin.UpdateRequest{URN: u, ID: created.ID, OldInputs: next, OldOutputs: updated.Properties, NewInputs: removed})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Properties["version"].StringValue() != "9007199254740993" {
		t.Fatal("version lost precision")
	}
	if v, ok := api.templates[0]["ms_delay"]; !ok || v != nil {
		t.Fatal("removal did not send null")
	}
	unchanged(t, p, u, created.ID, removed, updated.Properties)
	imported, err := p.Read(ctx, plugin.ReadRequest{URN: u, ID: created.ID})
	if err != nil {
		t.Fatal(err)
	}
	unchanged(t, p, u, created.ID, imported.Inputs, imported.Outputs)
	imageInputs := checked(t, p, urn("ImageTemplated"), nil, props(map[string]any{"templateId": string(created.ID), "templateVersion": updated.Properties["version"].StringValue(), "templateValues": `{"name":"Hello"}`}))
	image, err := p.Create(ctx, plugin.CreateRequest{URN: urn("ImageTemplated"), Properties: imageInputs})
	if err != nil {
		t.Fatal(err)
	}
	if image.Properties["resolvedTemplateVersion"].StringValue() != "9007199254740993" {
		t.Fatal("pin rounded")
	}
	_, err = p.Delete(ctx, plugin.DeleteRequest{URN: u, ID: created.ID, Inputs: removed, Outputs: updated.Properties})
	if err != nil {
		t.Fatal(err)
	}
}
func TestOGBridgeLifecycle(t *testing.T) {
	p, api := server(t)
	u := urn("OgConfig")
	in := checked(t, p, u, nil, props(map[string]any{"name": "Website cards", "baseUrl": "https://example.com", "configType": "html_css", "defaultOptions": map[string]any{"headers": map[string]any{"Authorization": "secret"}}}))
	created, err := p.Create(ctx, plugin.CreateRequest{URN: u, Properties: in})
	if err != nil {
		t.Fatal(err)
	}
	secretMap(t, created.Properties["defaultOptions"].ObjectValue()["headers"], "Authorization", "secret")
	unchanged(t, p, u, created.ID, in, created.Properties)
	next := checked(t, p, u, in, props(map[string]any{"name": "Website cards", "baseUrl": "https://example.com", "configType": "templated", "templateId": "t-test", "templateVersion": "9007199254740993", "headers": map[string]any{"Authorization": "secret"}}))
	updated, err := p.Update(ctx, plugin.UpdateRequest{URN: u, ID: created.ID, OldInputs: in, OldOutputs: created.Properties, NewInputs: next})
	if err != nil {
		t.Fatal(err)
	}
	secretMap(t, updated.Properties["headers"], "Authorization", "secret")
	remote, _, _ := api.OGSnapshot()
	if remote.TemplateVersion == nil || *remote.TemplateVersion != 9007199254740993 {
		t.Fatal("pin rounded")
	}
	unchanged(t, p, u, created.ID, next, updated.Properties)
	imported, err := p.Read(ctx, plugin.ReadRequest{URN: u, ID: created.ID})
	if err != nil {
		t.Fatal(err)
	}
	secretMap(t, imported.Outputs["headers"], "Authorization", "secret")
	unchanged(t, p, u, created.ID, imported.Inputs, imported.Outputs)
	_, err = p.Delete(ctx, plugin.DeleteRequest{URN: u, ID: created.ID, Inputs: next, Outputs: updated.Properties})
	if err != nil {
		t.Fatal(err)
	}
}
func TestStorageBridgeLifecycleAndExternalID(t *testing.T) {
	p, api := server(t)
	u := urn("StorageDestination")
	lookup, err := p.Invoke(ctx, plugin.InvokeRequest{Tok: "html-css-to-image:index:getAwsStorageExternalId", Args: props(map[string]any{})})
	if err != nil || len(lookup.Failures) > 0 || lookup.Properties["externalId"].StringValue() != "org-external-test" || lookup.Properties["writerRoleArn"].StringValue() != "arn:aws:iam::123456789012:role/hcti-writer" {
		t.Fatal("lookup failed", err)
	}
	in := checked(t, p, u, nil, props(map[string]any{"name": "Image storage", "connectionInfo": map[string]any{"cloudflareR2": map[string]any{"bucket": "existing-bucket", "cloudflareAccountId": "0123456789abcdef0123456789abcdef", "accessKeyId": "key-id", "secretAccessKey": "secret 😀"}}}))
	created, err := p.Create(ctx, plugin.CreateRequest{URN: u, Properties: in})
	if err != nil {
		t.Fatal(err)
	}
	secret(t, created.Properties["connectionInfo"].ObjectValue()["cloudflareR2"].ObjectValue()["secretAccessKey"], "secret 😀")
	unchanged(t, p, u, created.ID, in, created.Properties)
	imported, err := p.Read(ctx, plugin.ReadRequest{URN: u, ID: created.ID})
	if err != nil {
		t.Fatal(err)
	}
	v := imported.Outputs["connectionInfo"].ObjectValue()["cloudflareR2"].ObjectValue()["secretAccessKey"]
	if !v.IsNull() && !(v.IsSecret() && v.SecretValue().Element.IsNull()) {
		t.Fatal("invented imported secret")
	}
	next := imported.Inputs.Copy()
	next["name"] = resource.NewStringProperty("Renamed storage")
	next = checked(t, p, u, imported.Inputs, next)
	updated, err := p.Update(ctx, plugin.UpdateRequest{URN: u, ID: created.ID, OldInputs: imported.Inputs, OldOutputs: imported.Outputs, NewInputs: next})
	if err != nil {
		t.Fatal(err)
	}
	_, _, _, retains := api.StorageSnapshot()
	if retains != 1 {
		t.Fatal("secret not retained")
	}
	unchanged(t, p, u, created.ID, next, updated.Properties)
	_, err = p.Delete(ctx, plugin.DeleteRequest{URN: u, ID: created.ID, Inputs: next, Outputs: updated.Properties})
	if err != nil {
		t.Fatal(err)
	}
}
func TestNewResourceUnknownPreviews(t *testing.T) {
	p, api := server(t)
	for _, tc := range []struct {
		name, field string
		data        map[string]any
	}{{"OgConfig", "defaultOptions", map[string]any{"name": "Website cards", "baseUrl": "https://example.com", "configType": "html_css"}}, {"StorageDestination", "connectionInfo", map[string]any{"name": "Image storage"}}, {"Template", "html", map[string]any{}}, {"ImageHtmlCss", "html", map[string]any{}}, {"ImageUrl", "headers", map[string]any{"url": "https://example.com"}}, {"ImageTemplated", "templateVersion", map[string]any{"templateId": "t-test", "templateValues": `{"name":"hello"}`}}} {
		t.Run(tc.name, func(t *testing.T) {
			in := props(tc.data)
			in[resource.PropertyKey(tc.field)] = resource.MakeComputed(resource.NewStringProperty(""))
			in = checked(t, p, urn(tc.name), nil, in)
			v, err := p.Create(ctx, plugin.CreateRequest{URN: urn(tc.name), Properties: in, Preview: true})
			if err != nil {
				t.Fatal(err)
			}
			if !v.Properties[resource.PropertyKey(tc.field)].ContainsUnknowns() {
				t.Fatal("preview lost unknown")
			}
		})
	}
	if _, n, _ := api.OGSnapshot(); n != 0 {
		t.Fatal("preview wrote")
	}
	if _, n, _, _ := api.StorageSnapshot(); n != 0 {
		t.Fatal("preview wrote")
	}
}

func TestEveryTerraformResourceHasMapping(t *testing.T) {
	info := Info("test")
	info.P.ResourcesMap().Range(func(name string, _ shim.Resource) bool {
		if info.Resources[name] == nil {
			t.Errorf("unmapped Terraform resource: %s", name)
		}
		return true
	})
	info.P.DataSourcesMap().Range(func(name string, _ shim.Resource) bool {
		if info.DataSources[name] == nil {
			t.Errorf("unmapped Terraform data source: %s", name)
		}
		return true
	})
}

func TestTemplateLookupVersionsAreExact(t *testing.T) {
	p, api := renderingServer(t)
	const large = "9007199254740993"
	api.mu.Lock()
	api.templates = []map[string]any{
		{"id": "t-test", "version": json.Number(large), "template_type": "html_css", "html": "<h1>{{title}}</h1>", "created_at": "2026-01-01T00:00:00Z", "updated_at": "2026-01-02T00:00:00Z"},
		{"id": "t-test", "version": json.Number("42"), "template_type": "html_css", "html": "<h1>{{title}}</h1>", "created_at": "2026-01-01T00:00:00Z", "updated_at": "2026-01-01T00:00:00Z"},
	}
	api.mu.Unlock()
	for _, tc := range []struct {
		args    map[string]any
		version string
	}{
		{map[string]any{"id": "t-test"}, large},
		{map[string]any{"id": "t-test", "version": large}, large},
		{map[string]any{"id": "t-test", "version": "42"}, "42"},
	} {
		result, err := p.Invoke(ctx, plugin.InvokeRequest{Tok: "html-css-to-image:index:getTemplate", Args: props(tc.args)})
		if err != nil || len(result.Failures) != 0 {
			t.Fatal("lookup failed", err, result.Failures)
		}
		if result.Properties["version"].StringValue() != tc.version {
			t.Fatal("version precision lost", result.Properties)
		}
	}
	result, err := p.Invoke(ctx, plugin.InvokeRequest{Tok: "html-css-to-image:index:getTemplateVersions", Args: props(map[string]any{"id": "t-test", "limit": 1})})
	if err != nil || len(result.Failures) != 0 {
		t.Fatal("versions failed", err, result.Failures)
	}
	versions := result.Properties["versions"].ArrayValue()
	if len(versions) != 1 || versions[0].ObjectValue()["version"].StringValue() != large {
		t.Fatal("limit or precision lost", result.Properties)
	}
	result, err = p.Invoke(ctx, plugin.InvokeRequest{Tok: "html-css-to-image:index:getTemplateVersions", Args: props(map[string]any{"id": "t-test"})})
	if err != nil || len(result.Failures) != 0 || result.Properties["limit"].NumberValue() != 1000 {
		t.Fatal("default limit lost", err, result.Failures)
	}
	_, err = p.Invoke(ctx, plugin.InvokeRequest{Tok: "html-css-to-image:index:getTemplate", Args: props(map[string]any{"id": "t-test", "version": "9007199254740993.0"})})
	if err == nil {
		t.Fatal("invalid version accepted")
	}
}
