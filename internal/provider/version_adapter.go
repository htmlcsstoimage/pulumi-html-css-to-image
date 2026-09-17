package provider

import (
	"context"
	"math/big"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	framework "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// The bridge's string-over-number encoder in v3.139.0 parses through float64,
// rounding int64 inputs above 2^53. Present actual string attributes to the bridge
// and convert directly at the framework boundary. All lifecycle work stays upstream.
type versionProvider struct{ framework.Provider }

func (p versionProvider) Resources(ctx context.Context) []func() resource.Resource {
	factories := p.Provider.Resources(ctx)
	for i, factory := range factories {
		f := factory
		factories[i] = func() resource.Resource {
			r := f()
			var m resource.MetadataResponse
			r.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "htmlcsstoimage"}, &m)
			fields := versionNames(m.TypeName)
			if len(fields) == 0 {
				return r
			}
			wrapped := &versionResource{Resource: r, fields: fields}
			var sch resource.SchemaResponse
			wrapped.Schema(ctx, resource.SchemaRequest{}, &sch)
			return wrapped
		}
	}
	return factories
}
func versionNames(name string) []string {
	switch name {
	case "htmlcsstoimage_template":
		return []string{"version"}
	case "htmlcsstoimage_og_config":
		return []string{"template_version"}
	case "htmlcsstoimage_image_templated":
		return []string{"template_version", "resolved_template_version", "og_config_content_version"}
	case "htmlcsstoimage_image_html_css", "htmlcsstoimage_image_url":
		return []string{"og_config_content_version"}
	}
	return nil
}

type versionResource struct {
	resource.Resource
	fields            []string
	upstream, bridged schema.Schema
}

func (r *versionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	if r.upstream.Attributes != nil {
		resp.Schema = r.bridged
		return
	}
	r.Resource.Schema(ctx, req, resp)
	r.upstream = resp.Schema
	r.bridged = resp.Schema
	attrs := map[string]schema.Attribute{}
	for k, v := range resp.Schema.Attributes {
		attrs[k] = v
	}
	r.bridged.Attributes = attrs
	for _, k := range r.fields {
		v := attrs[k].(schema.Int64Attribute)
		attrs[k] = schema.StringAttribute{Required: v.Required, Optional: v.Optional, Computed: v.Computed, Sensitive: v.Sensitive, MarkdownDescription: v.MarkdownDescription + " Represented as a decimal string in Pulumi to preserve int64 precision."}
	}
	resp.Schema = r.bridged
}
func (r *versionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if v, ok := r.Resource.(resource.ResourceWithConfigure); ok {
		v.Configure(ctx, req, resp)
	}
}
func (r *versionResource) convert(ctx context.Context, raw tftypes.Value, toStrings bool, d *diag.Diagnostics) tftypes.Value {
	previousErrors := d.ErrorsCount()
	target := r.upstream
	if toStrings {
		target = r.bridged
	}
	typ := target.Type().TerraformType(ctx)
	if raw.Type() == nil || raw.IsNull() {
		return tftypes.NewValue(typ, nil)
	}
	if !raw.IsKnown() {
		return tftypes.NewValue(typ, tftypes.UnknownValue)
	}
	var values map[string]tftypes.Value
	if err := raw.As(&values); err != nil {
		d.AddError("Unable to convert version fields", err.Error())
		return tftypes.NewValue(typ, nil)
	}
	copied := make(map[string]tftypes.Value, len(values))
	for k, v := range values {
		copied[k] = v
	}
	values = copied
	for _, k := range r.fields {
		v := values[k]
		t := tftypes.Number
		if toStrings {
			t = tftypes.String
		}
		if v.IsNull() {
			values[k] = tftypes.NewValue(t, nil)
			continue
		}
		if !v.IsKnown() {
			values[k] = tftypes.NewValue(t, tftypes.UnknownValue)
			continue
		}
		if toStrings {
			var n big.Float
			if err := v.As(&n); err != nil {
				d.AddError("Invalid version", "Expected an int64 version.")
				continue
			}
			i, accuracy := n.Int64()
			if accuracy != big.Exact {
				d.AddError("Invalid version", "Version is outside the int64 range.")
				continue
			}
			values[k] = tftypes.NewValue(t, strconv.FormatInt(i, 10))
		} else {
			var s string
			if err := v.As(&s); err != nil {
				d.AddError("Invalid version", "Expected a decimal string.")
				continue
			}
			i, err := strconv.ParseInt(s, 10, 64)
			if err != nil || strconv.FormatInt(i, 10) != s {
				d.AddError("Invalid version", "Use a canonical decimal int64 string, without spaces or a decimal point.")
				continue
			}
			values[k] = tftypes.NewValue(t, i)
		}
	}
	if d.ErrorsCount() > previousErrors {
		return tftypes.NewValue(typ, nil)
	}
	return tftypes.NewValue(typ, values)
}
func (r *versionResource) state(ctx context.Context, s tfsdk.State, out bool, d *diag.Diagnostics) tfsdk.State {
	raw := r.convert(ctx, s.Raw, out, d)
	sch := r.upstream
	if out {
		sch = r.bridged
	}
	return tfsdk.State{Raw: raw, Schema: sch}
}
func (r *versionResource) plan(ctx context.Context, s tfsdk.Plan, out bool, d *diag.Diagnostics) tfsdk.Plan {
	raw := r.convert(ctx, s.Raw, out, d)
	sch := r.upstream
	if out {
		sch = r.bridged
	}
	return tfsdk.Plan{Raw: raw, Schema: sch}
}
func (r *versionResource) config(ctx context.Context, s tfsdk.Config, d *diag.Diagnostics) tfsdk.Config {
	return tfsdk.Config{Raw: r.convert(ctx, s.Raw, false, d), Schema: r.upstream}
}
func (r *versionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	req.Config = r.config(ctx, req.Config, &resp.Diagnostics)
	req.Plan = r.plan(ctx, req.Plan, false, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.State = r.state(ctx, resp.State, false, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	r.Resource.Create(ctx, req, resp)
	resp.State = r.state(ctx, resp.State, true, &resp.Diagnostics)
}
func (r *versionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	req.State = r.state(ctx, req.State, false, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.State = r.state(ctx, resp.State, false, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	r.Resource.Read(ctx, req, resp)
	resp.State = r.state(ctx, resp.State, true, &resp.Diagnostics)
}
func (r *versionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	req.Config = r.config(ctx, req.Config, &resp.Diagnostics)
	req.Plan = r.plan(ctx, req.Plan, false, &resp.Diagnostics)
	req.State = r.state(ctx, req.State, false, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.State = r.state(ctx, resp.State, false, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	r.Resource.Update(ctx, req, resp)
	resp.State = r.state(ctx, resp.State, true, &resp.Diagnostics)
}
func (r *versionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	req.State = r.state(ctx, req.State, false, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.State = r.state(ctx, resp.State, false, &resp.Diagnostics)
	r.Resource.Delete(ctx, req, resp)
	resp.State = r.state(ctx, resp.State, true, &resp.Diagnostics)
}
func (r *versionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.State = r.state(ctx, resp.State, false, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	r.Resource.(resource.ResourceWithImportState).ImportState(ctx, req, resp)
	resp.State = r.state(ctx, resp.State, true, &resp.Diagnostics)
}
func (r *versionResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	req.Config = r.config(ctx, req.Config, &resp.Diagnostics)
	req.Plan = r.plan(ctx, req.Plan, false, &resp.Diagnostics)
	req.State = r.state(ctx, req.State, false, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Plan = r.plan(ctx, resp.Plan, false, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if v, ok := r.Resource.(resource.ResourceWithModifyPlan); ok {
		v.ModifyPlan(ctx, req, resp)
	}
	resp.Plan = r.plan(ctx, resp.Plan, true, &resp.Diagnostics)
}
func (r *versionResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	req.Config = r.config(ctx, req.Config, &resp.Diagnostics)
	if !resp.Diagnostics.HasError() {
		if v, ok := r.Resource.(resource.ResourceWithValidateConfig); ok {
			v.ValidateConfig(ctx, req, resp)
		}
	}
}
