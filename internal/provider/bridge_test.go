package provider

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/htmlcsstoimage/pulumi-html-css-to-image/internal/testapi"
	pf "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfbridge"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
)

var ctx = context.Background()

const proxyURN resource.URN = "urn:pulumi:test::test::html-css-to-image:index:Proxy::test"
const keyURN resource.URN = "urn:pulumi:test::test::html-css-to-image:index:ApiKey::test"

func server(t *testing.T) (plugin.Provider, *testapi.Server) {
	t.Helper()
	api := testapi.New(t)
	p, err := pf.NewProvider(ctx, Info("0.0.1"), pf.ProviderMetadata{PackageSchema: Schema})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = p.Close() })
	_, err = p.Configure(ctx, plugin.ConfigureRequest{Inputs: props(map[string]any{"apiId": "test-id", "apiKey": "test-key", "baseUrl": api.URL})})
	if err != nil {
		t.Fatal(err)
	}
	return p, api
}
func props(v map[string]any) resource.PropertyMap { return resource.NewPropertyMapFromMap(v) }
func checked(t *testing.T, p plugin.Provider, u resource.URN, old, news resource.PropertyMap) resource.PropertyMap {
	t.Helper()
	r, err := p.Check(ctx, plugin.CheckRequest{URN: u, Olds: old, News: news, AllowUnknowns: true})
	if err != nil || len(r.Failures) > 0 {
		t.Fatalf("check: %v %v", err, r.Failures)
	}
	return r.Properties
}
func unchanged(t *testing.T, p plugin.Provider, u resource.URN, id resource.ID, inputs, state resource.PropertyMap) {
	t.Helper()
	d, err := p.Diff(ctx, plugin.DiffRequest{URN: u, ID: id, OldInputs: inputs, OldOutputs: state, NewInputs: inputs})
	if err != nil || d.Changes != plugin.DiffNone {
		t.Fatalf("unexpected diff: %v %+v", err, d)
	}
}
func secret(t *testing.T, v resource.PropertyValue, want string) {
	t.Helper()
	if !v.IsSecret() || v.SecretValue().Element.StringValue() != want {
		t.Fatal("missing or incorrect secret")
	}
}
func TestProxyLifecycle(t *testing.T) {
	p, api := server(t)
	inputs := checked(t, p, proxyURN, nil, props(map[string]any{"name": "  initial proxy  ", "url": "https://proxy.example", "authentication": map[string]any{"username": "user", "password": "secret"}, "bypassHosts": []any{"EXAMPLE.COM", "https://example.com/path"}}))
	preview, err := p.Create(ctx, plugin.CreateRequest{URN: proxyURN, Properties: inputs, Preview: true})
	if err != nil {
		t.Fatal(err)
	}
	if preview.ID != "" {
		t.Fatal("preview returned ID")
	}
	if n, _ := api.Counts(); n != 0 {
		t.Fatal("preview wrote")
	}
	created, err := p.Create(ctx, plugin.CreateRequest{URN: proxyURN, Properties: inputs})
	if err != nil {
		t.Fatal(err)
	}
	if created.ID != "proxy-test" {
		t.Fatal("incorrect ID")
	}
	secret(t, created.Properties["authentication"].ObjectValue()["password"], "secret")
	if !created.Properties["port"].IsNull() {
		t.Fatal("injected port default")
	}
	unchanged(t, p, proxyURN, created.ID, inputs, created.Properties)
	read, err := p.Read(ctx, plugin.ReadRequest{URN: proxyURN, ID: created.ID, Inputs: inputs, State: created.Properties})
	if err != nil {
		t.Fatal(err)
	}
	secret(t, read.Outputs["authentication"].ObjectValue()["password"], "secret")
	unchanged(t, p, proxyURN, created.ID, read.Inputs, read.Outputs)
	imported, err := p.Read(ctx, plugin.ReadRequest{URN: proxyURN, ID: created.ID})
	if err != nil {
		t.Fatal(err)
	}
	auth := imported.Outputs["authentication"].ObjectValue()
	if v := auth["password"]; !v.IsNull() && !(v.IsSecret() && v.SecretValue().Element.IsNull()) {
		t.Fatal("import invented password")
	}
	next := imported.Inputs.Copy()
	next["name"] = resource.NewStringProperty("renamed proxy")
	next = checked(t, p, proxyURN, imported.Inputs, next)
	updated, err := p.Update(ctx, plugin.UpdateRequest{URN: proxyURN, ID: created.ID, OldInputs: imported.Inputs, OldOutputs: imported.Outputs, NewInputs: next})
	if err != nil {
		t.Fatal(err)
	}
	if _, retains := api.Counts(); retains != 1 {
		t.Fatal("imported password not retained")
	}
	removed := next.Copy()
	delete(removed, "authentication")
	removed = checked(t, p, proxyURN, next, removed)
	updated, err = p.Update(ctx, plugin.UpdateRequest{URN: proxyURN, ID: created.ID, OldInputs: next, OldOutputs: updated.Properties, NewInputs: removed})
	if err != nil {
		t.Fatal(err)
	}
	unchanged(t, p, proxyURN, created.ID, removed, updated.Properties)
	_, err = p.Delete(ctx, plugin.DeleteRequest{URN: proxyURN, ID: created.ID, Inputs: removed, Outputs: updated.Properties})
	if err != nil {
		t.Fatal(err)
	}
	missing, err := p.Read(ctx, plugin.ReadRequest{URN: proxyURN, ID: created.ID, Inputs: removed, State: updated.Properties})
	if err != nil || missing.ID != "" {
		t.Fatalf("404 read: %v", err)
	}
}
func TestAPIKeyLifecycle(t *testing.T) {
	p, api := server(t)
	inputs := checked(t, p, keyURN, nil, props(map[string]any{"permissions": []any{"images:create"}}))
	_, err := p.Create(ctx, plugin.CreateRequest{URN: keyURN, Properties: inputs, Preview: true})
	if err != nil {
		t.Fatal(err)
	}
	if _, n := api.KeySnapshot(); n != 0 {
		t.Fatal("preview wrote")
	}
	created, err := p.Create(ctx, plugin.CreateRequest{URN: keyURN, Properties: inputs})
	if err != nil {
		t.Fatal(err)
	}
	secret(t, created.Properties["apiKey"], "create-only-secret")
	if !created.Properties["name"].IsNull() {
		t.Fatal("injected generated name")
	}
	unchanged(t, p, keyURN, created.ID, inputs, created.Properties)
	read, err := p.Read(ctx, plugin.ReadRequest{URN: keyURN, ID: created.ID, Inputs: inputs, State: created.Properties})
	if err != nil {
		t.Fatal(err)
	}
	secret(t, read.Outputs["apiKey"], "create-only-secret")
	imported, err := p.Read(ctx, plugin.ReadRequest{URN: keyURN, ID: created.ID})
	if err != nil {
		t.Fatal(err)
	}
	v := imported.Outputs["apiKey"]
	if !v.IsNull() && !(v.IsSecret() && v.SecretValue().Element.IsNull()) {
		t.Fatal("import invented secret")
	}
	next := checked(t, p, keyURN, inputs, props(map[string]any{"allFuturePermissions": true}))
	updated, err := p.Update(ctx, plugin.UpdateRequest{URN: keyURN, ID: created.ID, OldInputs: inputs, OldOutputs: read.Outputs, NewInputs: next})
	if err != nil {
		t.Fatal(err)
	}
	secret(t, updated.Properties["apiKey"], "create-only-secret")
	api.ExpandKeyGrants()
	read, err = p.Read(ctx, plugin.ReadRequest{URN: keyURN, ID: created.ID, Inputs: next, State: updated.Properties})
	if err != nil {
		t.Fatal(err)
	}
	unchanged(t, p, keyURN, created.ID, next, read.Outputs)
	_, err = p.Delete(ctx, plugin.DeleteRequest{URN: keyURN, ID: created.ID, Inputs: next, Outputs: read.Outputs})
	if err != nil {
		t.Fatal(err)
	}
	key, _ := api.KeySnapshot()
	if key.Enabled {
		t.Fatal("destroy did not disable key")
	}
}
func TestUnknownsStayUnknown(t *testing.T) {
	for _, tc := range []struct {
		u      resource.URN
		field  string
		inputs map[string]any
	}{
		{keyURN, "permissions", map[string]any{}},
		{proxyURN, "bypassHosts", map[string]any{"name": "test proxy", "url": "https://proxy.example"}},
		{proxyURN, "authentication", map[string]any{"name": "test proxy", "url": "https://proxy.example"}},
	} {
		t.Run(tc.field, func(t *testing.T) {
			p, api := server(t)
			in := props(tc.inputs)
			in[resource.PropertyKey(tc.field)] = resource.MakeComputed(resource.NewStringProperty(""))
			in = checked(t, p, tc.u, nil, in)
			if !in[resource.PropertyKey(tc.field)].ContainsUnknowns() {
				t.Fatal("check lost unknown")
			}
			preview, err := p.Create(ctx, plugin.CreateRequest{URN: tc.u, Properties: in, Preview: true})
			if err != nil {
				t.Fatal(err)
			}
			if !preview.Properties[resource.PropertyKey(tc.field)].ContainsUnknowns() {
				t.Fatal("preview lost unknown")
			}
			if n, _ := api.Counts(); n != 0 {
				t.Fatal("preview wrote proxy")
			}
			if _, n := api.KeySnapshot(); n != 0 {
				t.Fatal("preview wrote key")
			}
		})
	}
}
func TestSchema(t *testing.T) {
	var s struct {
		Resources map[string]struct {
			InputProperties map[string]json.RawMessage `json:"inputProperties"`
			Properties      map[string]struct {
				Secret      bool
				Description string
			}
		}
		Types map[string]struct {
			Properties map[string]struct{ Secret bool }
		}
	}
	if err := json.Unmarshal(Schema, &s); err != nil {
		t.Fatal(err)
	}
	if len(s.Resources) != 8 {
		t.Fatal("expected exactly eight resources")
	}
	key := s.Resources["html-css-to-image:index:ApiKey"]
	if !key.Properties["apiKey"].Secret || key.InputProperties["apiKey"] != nil {
		t.Fatal("key must be a secret output only")
	}
	if key.Properties["permissions"].Description == "" {
		t.Fatal("missing upstream field description")
	}
	passwordSecret := false
	for _, typ := range s.Types {
		if typ.Properties["password"].Secret {
			passwordSecret = true
		}
	}
	if !passwordSecret {
		t.Fatal("missing password sensitivity")
	}
}

func TestReadErrorsDoNotRemoveState(t *testing.T) {
	for _, u := range []resource.URN{proxyURN, keyURN} {
		for _, status := range []int{403, 429, 500} {
			t.Run(string(u.Type())+"/"+strconv.Itoa(status), func(t *testing.T) {
				p, api := server(t)
				data := map[string]any{"name": "test proxy", "url": "https://proxy.example"}
				if u == keyURN {
					data = map[string]any{"permissions": []any{}}
				}
				in := checked(t, p, u, nil, props(data))
				created, err := p.Create(ctx, plugin.CreateRequest{URN: u, Properties: in})
				if err != nil {
					t.Fatal(err)
				}
				api.Status(status)
				_, err = p.Read(ctx, plugin.ReadRequest{URN: u, ID: created.ID, Inputs: in, State: created.Properties})
				if err == nil {
					t.Fatal("read error reported success")
				}
			})
		}
	}
}

func TestUnknownUpdatePreview(t *testing.T) {
	p, api := server(t)
	in := checked(t, p, keyURN, nil, props(map[string]any{"permissions": []any{}}))
	created, err := p.Create(ctx, plugin.CreateRequest{URN: keyURN, Properties: in})
	if err != nil {
		t.Fatal(err)
	}
	next := in.Copy()
	next["permissions"] = resource.MakeComputed(resource.NewStringProperty(""))
	next = checked(t, p, keyURN, in, next)
	preview, err := p.Update(ctx, plugin.UpdateRequest{URN: keyURN, ID: created.ID, OldInputs: in, OldOutputs: created.Properties, NewInputs: next, Preview: true})
	if err != nil {
		t.Fatal(err)
	}
	if !preview.Properties["permissions"].ContainsUnknowns() {
		t.Fatal("update preview lost unknown permissions")
	}
	if _, n := api.KeySnapshot(); n != 1 {
		t.Fatal("update preview wrote")
	}
	secret(t, preview.Properties["apiKey"], "create-only-secret")
}

func TestValidationAndMalformedRead(t *testing.T) {
	p, api := server(t)
	in := props(map[string]any{"allFuturePermissions": false})
	c, err := p.Check(ctx, plugin.CheckRequest{URN: keyURN, News: in})
	if err == nil && len(c.Failures) == 0 {
		// Terraform's cross-field validation runs during planning, not Check.
		_, err = p.Create(ctx, plugin.CreateRequest{URN: keyURN, Properties: c.Properties, Preview: true})
		if err == nil {
			t.Fatal("missing permissions accepted")
		}
	}
	if _, n := api.KeySnapshot(); n != 0 {
		t.Fatal("invalid input wrote")
	}
	in = checked(t, p, keyURN, nil, props(map[string]any{"permissions": []any{}}))
	created, err := p.Create(ctx, plugin.CreateRequest{URN: keyURN, Properties: in})
	if err != nil {
		t.Fatal(err)
	}
	api.KeyBody(`{"id":"key-management"}`)
	_, err = p.Read(ctx, plugin.ReadRequest{URN: keyURN, ID: created.ID, Inputs: in, State: created.Properties})
	if err == nil {
		t.Fatal("malformed read treated as valid state")
	}
}
