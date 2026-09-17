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

In a Pulumi project for your language, install the SDK and configure credentials:

```sh
pulumi config set html-css-to-image:apiId YOUR_API_ID
pulumi config set --secret html-css-to-image:apiKey
```

```go
package main

import (
    hcti "github.com/htmlcsstoimage/pulumi-html-css-to-image/sdk/go"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        card, err := hcti.NewImageHtmlCss(ctx, "card", &hcti.ImageHtmlCssArgs{
            Html: pulumi.String("<h1>Hello from Pulumi</h1>"),
        })
        if err != nil {
            return err
        }
        ctx.Export("imageUrl", card.ImageUrl)
        return nil
    })
}
```

Run `pulumi preview` and `pulumi up`. Open the returned URL to render the image. Creation and refresh save/read metadata without rendering bytes. Input changes replace image definitions; template edits create new versions under a stable ID.

## Configuration

| Setting | Environment fallback | Description |
| --- | --- | --- |
| `html-css-to-image:apiId` | `HCTI_API_ID` | API ID for your organization. |
| `html-css-to-image:apiKey` | `HCTI_API_KEY` | API key; use `pulumi config set --secret` when storing it in stack configuration. |
| `html-css-to-image:baseUrl` | None | Optional API origin; defaults to `https://hcti.io`. |

The key needs the read/write/delete permissions for the resources you manage. Credentials and headers retain Pulumi secret markings. Version identifiers are decimal strings to preserve all int64 digits. Requests identify themselves with `HCTIGo/<sdk-version> HCTIPulumi/<provider-version>`.

## Resources and lookups

- [Proxy](https://github.com/htmlcsstoimage/pulumi-html-css-to-image/tree/main/docs/resources/proxy.md), [ApiKey](https://github.com/htmlcsstoimage/pulumi-html-css-to-image/tree/main/docs/resources/api_key.md), [OgConfig](https://github.com/htmlcsstoimage/pulumi-html-css-to-image/tree/main/docs/resources/og_config.md), and [StorageDestination](https://github.com/htmlcsstoimage/pulumi-html-css-to-image/tree/main/docs/resources/storage_destination.md).
- [Template](https://github.com/htmlcsstoimage/pulumi-html-css-to-image/tree/main/docs/resources/template.md), [ImageHtmlCss](https://github.com/htmlcsstoimage/pulumi-html-css-to-image/tree/main/docs/resources/image_html_css.md), [ImageUrl](https://github.com/htmlcsstoimage/pulumi-html-css-to-image/tree/main/docs/resources/image_url.md), and [ImageTemplated](https://github.com/htmlcsstoimage/pulumi-html-css-to-image/tree/main/docs/resources/image_templated.md).
- [getTemplate](https://github.com/htmlcsstoimage/pulumi-html-css-to-image/tree/main/docs/functions/get_template.md): latest or pinned template settings.
- [getTemplateVersions](https://github.com/htmlcsstoimage/pulumi-html-css-to-image/tree/main/docs/functions/get_template_versions.md): newest-first history with automatic pagination and a default limit of 1000.
- [getAwsStorageExternalId](https://github.com/htmlcsstoimage/pulumi-html-css-to-image/tree/main/docs/functions/get_aws_storage_external_id.md): external ID and writer role ARN for an AWS trust policy.

See [examples](https://github.com/htmlcsstoimage/pulumi-html-css-to-image/tree/main/examples) for complete YAML programs. In .NET, the API key secret is `ApiKey.Value` and the URL image operation URL is `ImageUrl.RenderingUrl` to avoid member/type name collisions.

## Contributing

See [development](https://github.com/htmlcsstoimage/pulumi-html-css-to-image/tree/main/docs/development.md) and [releasing](https://github.com/htmlcsstoimage/pulumi-html-css-to-image/tree/main/docs/releasing.md). API lifecycle behavior belongs in the Terraform provider; this repository maintains the bridge, exact version conversions, SDK metadata, and Pulumi examples.

---

[HTML/CSS to Image](https://htmlcsstoimage.com) · [Documentation](https://docs.htmlcsstoimage.com/) · [Management API](https://docs.htmlcsstoimage.com/management-api/) · [Pulumi provider repository](https://github.com/htmlcsstoimage/pulumi-html-css-to-image) · [Report an issue](https://github.com/htmlcsstoimage/pulumi-html-css-to-image/issues)
