---
title: HTML/CSS to Image
meta_desc: Create images and manage templates, API keys, storage destinations, proxies, and Open Graph configurations with Pulumi.
layout: package
---

The HTML/CSS to Image provider manages images and reusable templates alongside your application's infrastructure.

[Website](https://htmlcsstoimage.com) · [Documentation](https://docs.htmlcsstoimage.com/) · [Management API](https://docs.htmlcsstoimage.com/management-api/)

## Installation

- JavaScript / TypeScript: `npm install @html-css-to-image/pulumi`
- Python: `pip install pulumi-html-css-to-image`
- Go: `go get github.com/htmlcsstoimage/pulumi-html-css-to-image/sdk/go`
- Java: Maven dependency `com.htmlcsstoimage:pulumi:0.1.0`.
- .NET: `dotnet add package Pulumi.HtmlCssToImage`

The SDK downloads the matching provider plugin automatically.

## Example usage

```typescript
import * as hcti from "@html-css-to-image/pulumi";

const card = new hcti.ImageHtmlCss("card", {
    html: "<h1>Hello from Pulumi</h1>",
    css: "h1 { font-family: Inter; padding: 48px; }",
    googleFonts: ["Inter"],
});

export const imageUrl = card.imageUrl;
```

Open the returned URL to render the image. Creation saves an image definition; changes to image inputs replace that definition.

## Configuration

Set `HCTI_API_ID` and `HCTI_API_KEY`, or use Pulumi configuration:

```sh
pulumi config set html-css-to-image:apiId YOUR_API_ID
pulumi config set --secret html-css-to-image:apiKey YOUR_API_KEY
```

Use an API key with permissions to manage the resources in your stack. Find complete programs in the [examples directory](https://github.com/htmlcsstoimage/pulumi-html-css-to-image/tree/main/examples).

## Java example

```java
import com.htmlcsstoimage.pulumi.ImageHtmlCss;
import com.htmlcsstoimage.pulumi.ImageHtmlCssArgs;
import com.pulumi.Pulumi;

public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            var card = new ImageHtmlCss("card", ImageHtmlCssArgs.builder()
                    .html("<h1>Hello from Pulumi</h1>")
                    .css("h1 { font-family: Inter; padding: 48px; }")
                    .googleFonts("Inter")
                    .build());
            ctx.export("imageUrl", card.imageUrl());
        });
    }
}
```
