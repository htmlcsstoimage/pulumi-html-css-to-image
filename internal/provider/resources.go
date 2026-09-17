package provider

import (
	upstream "github.com/htmlcsstoimage/terraform-provider-html-css-to-image/provider"
	pf "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfbridge"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"
)

const Name = "html-css-to-image"

// Info maps public Pulumi names. Resource behavior and field descriptions come
// from Terraform; no API lifecycle logic belongs in this package.
func Info(version string) tfbridge.ProviderInfo {
	upstreamLicense := tfbridge.MITLicenseType
	return tfbridge.ProviderInfo{
		P:    pf.ShimProvider(versionProvider{upstream.New(version, "HCTIPulumi/"+version)}),
		Name: Name, ResourcePrefix: "htmlcsstoimage", Version: version,
		GitHubOrg: "htmlcsstoimage", DisplayName: "HTML/CSS to Image",
		Description:       "Create images and manage templates, API keys, storage destinations, proxies, and Open Graph configurations with the HTML/CSS to Image provider for Pulumi.",
		Keywords:          []string{"category/utility", "html", "css", "image", "screenshot", "pdf", "templates", "open-graph"},
		LogoURL:           "https://raw.githubusercontent.com/htmlcsstoimage/pulumi-html-css-to-image/main/assets/logo.png",
		TFProviderLicense: &upstreamLicense,
		License:           "MIT", Publisher: "htmlcsstoimage", Homepage: "https://htmlcsstoimage.com",
		Repository:          "https://github.com/htmlcsstoimage/pulumi-html-css-to-image",
		PluginDownloadURL:   "github://api.github.com/htmlcsstoimage/pulumi-html-css-to-image",
		MetadataInfo:        tfbridge.NewProviderMetadata(nil),
		SchemaPostProcessor: preserveAPILiterals,
		Config: map[string]*tfbridge.SchemaInfo{
			"api_id":   {Name: "apiId", Default: &tfbridge.DefaultInfo{EnvVars: []string{"HCTI_API_ID"}}},
			"api_key":  {Name: "apiKey", Default: &tfbridge.DefaultInfo{EnvVars: []string{"HCTI_API_KEY"}}},
			"base_url": {Name: "baseUrl"},
		},
		Resources: map[string]*tfbridge.ResourceInfo{
			"htmlcsstoimage_proxy":               {Tok: "html-css-to-image:index:Proxy"},
			"htmlcsstoimage_api_key":             {Tok: "html-css-to-image:index:ApiKey", Fields: map[string]*tfbridge.SchemaInfo{"api_key": {CSharpName: "Value"}}},
			"htmlcsstoimage_og_config":           {Tok: "html-css-to-image:index:OgConfig"},
			"htmlcsstoimage_storage_destination": {Tok: "html-css-to-image:index:StorageDestination"},
			"htmlcsstoimage_template":            {Tok: "html-css-to-image:index:Template"},
			"htmlcsstoimage_image_html_css":      {Tok: "html-css-to-image:index:ImageHtmlCss"},
			"htmlcsstoimage_image_url":           {Tok: "html-css-to-image:index:ImageUrl", Fields: map[string]*tfbridge.SchemaInfo{"image_url": {CSharpName: "RenderingUrl"}}},
			"htmlcsstoimage_image_templated":     {Tok: "html-css-to-image:index:ImageTemplated"},
		},
		DataSources: map[string]*tfbridge.DataSourceInfo{
			"htmlcsstoimage_aws_storage_external_id": {Tok: "html-css-to-image:index:getAwsStorageExternalId"},
			"htmlcsstoimage_template":                {Tok: "html-css-to-image:index:getTemplate"},
			"htmlcsstoimage_template_versions":       {Tok: "html-css-to-image:index:getTemplateVersions"},
		},
		JavaScript: &tfbridge.JavaScriptInfo{PackageName: "@html-css-to-image/pulumi", RespectSchemaVersion: true, Dependencies: map[string]string{"@pulumi/pulumi": "^3.0.0"}},
		Python:     &tfbridge.PythonInfo{PackageName: "pulumi_html_css_to_image", RespectSchemaVersion: true, Requires: map[string]string{"pulumi": ">=3.0.0,<4.0.0"}},
		Golang:     &tfbridge.GolangInfo{ImportBasePath: "github.com/htmlcsstoimage/pulumi-html-css-to-image/sdk/go", RespectSchemaVersion: true, RootPackageName: "htmlcsstoimage"},
		CSharp:     &tfbridge.CSharpInfo{RootNamespace: "Pulumi", RespectSchemaVersion: true, PackageReferences: map[string]string{"Pulumi": "3.*"}},
		Java:       &tfbridge.JavaInfo{BasePackage: "com.htmlcsstoimage", Packages: map[string]string{"html-css-to-image": "pulumi"}},
	}
}
