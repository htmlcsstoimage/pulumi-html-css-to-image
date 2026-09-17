package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	upstream "github.com/htmlcsstoimage/terraform-provider-html-css-to-image/provider"
)

func versionTestResource(t *testing.T) *versionResource {
	t.Helper()
	for _, f := range (versionProvider{upstream.New("test", "HCTIPulumi/test")}).Resources(ctx) {
		r := f()
		var m resource.MetadataResponse
		r.Metadata(ctx, resource.MetadataRequest{}, &m)
		if m.TypeName == "htmlcsstoimage_template" {
			return r.(*versionResource)
		}
	}
	t.Fatal("template not found")
	return nil
}
func versionTestRaw(r *versionResource, value any) tftypes.Value {
	typ := r.bridged.Type().TerraformType(ctx).(tftypes.Object)
	a := map[string]tftypes.Value{}
	for k, v := range typ.AttributeTypes {
		a[k] = tftypes.NewValue(v, nil)
	}
	a["version"] = tftypes.NewValue(tftypes.String, value)
	return tftypes.NewValue(typ, a)
}
func TestVersionAdapterExactAndImmutable(t *testing.T) {
	r := versionTestResource(t)
	for _, v := range []any{nil, tftypes.UnknownValue, "0", "9007199254740993", "9223372036854775807", "-9223372036854775808"} {
		raw := versionTestRaw(r, v)
		original := raw.Copy()
		var d diag.Diagnostics
		encoded := r.convert(ctx, raw, false, &d)
		if d.HasError() {
			t.Fatal(d)
		}
		decoded := r.convert(ctx, encoded, true, &d)
		if d.HasError() || !decoded.Equal(original) || !raw.Equal(original) {
			t.Fatalf("lossy or mutating conversion for %v: %v", v, d)
		}
	}
}
func TestVersionAdapterRejectsInvalidInput(t *testing.T) {
	r := versionTestResource(t)
	for _, s := range []string{"9007199254740993.0", "9223372036854775808", "1e3", " 1", "01", "secret-value"} {
		raw := versionTestRaw(r, s)
		var d diag.Diagnostics
		r.convert(ctx, raw, false, &d)
		if !d.HasError() {
			t.Fatal("accepted invalid version")
		}
		for _, v := range d {
			if v.Detail() == s {
				t.Fatal("input leaked")
			}
		}
	}
}
func TestVersionAdapterPreservesStateOnExistingError(t *testing.T) {
	r := versionTestResource(t)
	var d diag.Diagnostics
	raw := versionTestRaw(r, "9007199254740993")
	numeric := r.convert(ctx, raw, false, &d)
	d.AddError("API error", "An upstream read failed.")
	result := r.convert(ctx, numeric, true, &d)
	if !result.Equal(raw) {
		t.Fatal("upstream error erased saved state")
	}
}
func TestAPILiteralsNotTranslated(t *testing.T) {
	input := `<span pulumi-lang-nodejs="htmlCss" pulumi-lang-hcl="html_css">htmlCss</span>`
	if got := literalDocs(input); got != "html_css" {
		t.Fatal(got)
	}
}
