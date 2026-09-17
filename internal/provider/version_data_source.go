package provider

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func (p versionProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	factories := p.Provider.DataSources(ctx)
	for i, factory := range factories {
		f := factory
		factories[i] = func() datasource.DataSource {
			d := f()
			var m datasource.MetadataResponse
			d.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "htmlcsstoimage"}, &m)
			if m.TypeName != "htmlcsstoimage_template" && m.TypeName != "htmlcsstoimage_template_versions" {
				return d
			}
			wrapped := &versionDataSource{DataSource: d}
			var s datasource.SchemaResponse
			wrapped.Schema(ctx, datasource.SchemaRequest{}, &s)
			return wrapped
		}
	}
	return factories
}

type versionDataSource struct {
	datasource.DataSource
	upstream, bridged schema.Schema
}

func versionDataAttributes(attrs map[string]schema.Attribute) map[string]schema.Attribute {
	out := map[string]schema.Attribute{}
	for k, v := range attrs {
		out[k] = v
		if k == "version" {
			a := v.(schema.Int64Attribute)
			out[k] = schema.StringAttribute{Required: a.Required, Optional: a.Optional, Computed: a.Computed, MarkdownDescription: a.MarkdownDescription + " Represented as a decimal string to preserve int64 precision."}
		} else if a, ok := v.(schema.ListNestedAttribute); ok {
			a.NestedObject.Attributes = versionDataAttributes(a.NestedObject.Attributes)
			out[k] = a
		}
	}
	return out
}
func (d *versionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	if d.upstream.Attributes == nil {
		d.DataSource.Schema(ctx, req, resp)
		d.upstream = resp.Schema
		d.bridged = resp.Schema
		d.bridged.Attributes = versionDataAttributes(d.upstream.Attributes)
	}
	resp.Schema = d.bridged
}
func (d *versionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if c, ok := d.DataSource.(datasource.DataSourceWithConfigure); ok {
		c.Configure(ctx, req, resp)
	}
}

// Convert by target type, including version fields inside lists, without float64.
func convertVersionValue(raw tftypes.Value, target tftypes.Type) (tftypes.Value, error) {
	if raw.Type() == nil || raw.IsNull() {
		return tftypes.NewValue(target, nil), nil
	}
	if !raw.IsKnown() {
		return tftypes.NewValue(target, tftypes.UnknownValue), nil
	}
	if raw.Type().Equal(target) {
		return raw, nil
	}
	switch t := target.(type) {
	case tftypes.Object:
		var values map[string]tftypes.Value
		if err := raw.As(&values); err != nil {
			return raw, err
		}
		out := map[string]tftypes.Value{}
		for k, v := range values {
			n, err := convertVersionValue(v, t.AttributeTypes[k])
			if err != nil {
				return raw, err
			}
			out[k] = n
		}
		return tftypes.NewValue(target, out), nil
	case tftypes.List:
		var values []tftypes.Value
		if err := raw.As(&values); err != nil {
			return raw, err
		}
		for i, v := range values {
			n, err := convertVersionValue(v, t.ElementType)
			if err != nil {
				return raw, err
			}
			values[i] = n
		}
		return tftypes.NewValue(target, values), nil
	}
	if target.Equal(tftypes.String) && raw.Type().Equal(tftypes.Number) {
		var n big.Float
		if err := raw.As(&n); err != nil {
			return raw, err
		}
		i, accuracy := n.Int64()
		if accuracy != big.Exact {
			return raw, fmt.Errorf("version is outside the int64 range")
		}
		return tftypes.NewValue(target, strconv.FormatInt(i, 10)), nil
	}
	if target.Equal(tftypes.Number) && raw.Type().Equal(tftypes.String) {
		var s string
		if err := raw.As(&s); err != nil {
			return raw, err
		}
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil || strconv.FormatInt(i, 10) != s {
			return raw, fmt.Errorf("version must be a canonical decimal int64 string")
		}
		return tftypes.NewValue(target, i), nil
	}
	return raw, fmt.Errorf("unsupported version conversion")
}
func (d *versionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	raw, err := convertVersionValue(req.Config.Raw, d.upstream.Type().TerraformType(ctx))
	if err != nil {
		resp.Diagnostics.AddError("Invalid template version", err.Error())
		return
	}
	req.Config = tfsdk.Config{Schema: d.upstream, Raw: raw}
	upstream := datasource.ReadResponse{State: tfsdk.State{Schema: d.upstream, Raw: tftypes.NewValue(d.upstream.Type().TerraformType(ctx), nil)}}
	d.DataSource.Read(ctx, req, &upstream)
	resp.Diagnostics.Append(upstream.Diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
	raw, err = convertVersionValue(upstream.State.Raw, d.bridged.Type().TerraformType(ctx))
	if err != nil {
		resp.Diagnostics.AddError("Invalid template version", err.Error())
		return
	}
	resp.State = tfsdk.State{Schema: d.bridged, Raw: raw}
}
