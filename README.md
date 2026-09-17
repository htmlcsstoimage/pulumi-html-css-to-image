# HTML/CSS to Image Pulumi provider

Manage HTML/CSS to Image resources with TypeScript, Python, Go, C#, Java, or Pulumi YAML.

[Website](https://htmlcsstoimage.com) · [Documentation](https://docs.htmlcsstoimage.com/) · [Management API](https://docs.htmlcsstoimage.com/management-api/) · [API reference](https://htmlcsstoimage.com/api-docs/)

Create images and reusable templates alongside application API keys, proxies, storage destinations, and Open Graph configurations. The provider shares its implementation with the [Terraform provider](https://github.com/htmlcsstoimage/terraform-provider-html-css-to-image).

## Installation

Install the SDK for your Pulumi project's language:

| Language | Install |
| --- | --- |
| TypeScript / JavaScript | `npm install @html-css-to-image/pulumi` |
| Python | `pip install pulumi-html-css-to-image` |
| Go | `go get github.com/htmlcsstoimage/pulumi-html-css-to-image/sdk/go` |
| C# / .NET | `dotnet add package HtmlCssToImage.Pulumi` |

For Java, add this Maven dependency:

```xml
<dependency>
  <groupId>com.htmlcsstoimage</groupId>
  <artifactId>pulumi</artifactId>
  <version>0.1.2</version>
</dependency>
```

The SDK automatically downloads the matching provider plugin from GitHub Releases. YAML programs select the plugin directly, as shown below.

## Example

Set `HCTI_API_ID` and `HCTI_API_KEY` in your environment. This `Pulumi.yaml` creates an HTML/CSS image definition:

```yaml
name: hcti-image
runtime: yaml
resources:
  hcti:
    type: pulumi:providers:html-css-to-image
    defaultProvider: true
    options:
      version: 0.1.2
      pluginDownloadURL: github://api.github.com/htmlcsstoimage/pulumi-html-css-to-image
  card:
    type: html-css-to-image:index:ImageHtmlCss
    properties:
      html: '<h1>Hello from Pulumi</h1>'
      css: 'h1 { font-family: Inter; padding: 48px; }'
      googleFonts: [Inter]
outputs:
  imageUrl: ${card.imageUrl}
```

Run `pulumi stack init dev`, `pulumi preview`, and `pulumi up`. Open the returned URL to render the image. Creation and refresh read/save metadata without rendering bytes. Input changes replace image definitions; template edits create new versions under a stable ID.

## Configuration

| Setting | Environment fallback | Description |
| --- | --- | --- |
| `html-css-to-image:apiId` | `HCTI_API_ID` | API ID for your organization. |
| `html-css-to-image:apiKey` | `HCTI_API_KEY` | API key; use `pulumi config set --secret` when storing it in stack configuration. |
| `html-css-to-image:baseUrl` | None | Optional API origin; defaults to `https://hcti.io`. |

The key needs the read/write/delete permissions for the resources you manage. Credentials and headers retain Pulumi secret markings. Version identifiers are decimal strings to preserve all int64 digits. Requests identify themselves with `HCTIGo/<sdk-version> HCTIPulumi/<provider-version>`.

## Resources and lookups

- [Proxy](docs/resources/proxy.md), [ApiKey](docs/resources/api_key.md), [OgConfig](docs/resources/og_config.md), and [StorageDestination](docs/resources/storage_destination.md).
- [Template](docs/resources/template.md), [ImageHtmlCss](docs/resources/image_html_css.md), [ImageUrl](docs/resources/image_url.md), and [ImageTemplated](docs/resources/image_templated.md).
- [getTemplate](docs/functions/get_template.md): latest or pinned template settings.
- [getTemplateVersions](docs/functions/get_template_versions.md): newest-first history with automatic pagination and a default limit of 1000.
- [getAwsStorageExternalId](docs/functions/get_aws_storage_external_id.md): external ID and writer role ARN for an AWS trust policy.

See [examples](examples) for complete YAML programs. In .NET, the API key secret is `ApiKey.Value` and the URL image operation URL is `ImageUrl.RenderingUrl` to avoid member/type name collisions.

## Contributing

See [development](docs/development.md) and [releasing](docs/releasing.md). API lifecycle behavior belongs in the Terraform provider; this repository maintains the bridge, exact version conversions, SDK metadata, and Pulumi examples.

---

[HTML/CSS to Image](https://htmlcsstoimage.com) · [Documentation](https://docs.htmlcsstoimage.com/) · [Management API](https://docs.htmlcsstoimage.com/management-api/) · [Pulumi provider repository](https://github.com/htmlcsstoimage/pulumi-html-css-to-image) · [Report an issue](https://github.com/htmlcsstoimage/pulumi-html-css-to-image/issues)
