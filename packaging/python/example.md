## Example

In a Python Pulumi project, install `pulumi-html-css-to-image` and configure your credentials:

```sh
pulumi config set html-css-to-image:apiId YOUR_API_ID
pulumi config set --secret html-css-to-image:apiKey
```

Add this to `__main__.py`:

```python
import pulumi
import pulumi_html_css_to_image as hcti

card = hcti.ImageHtmlCss(
    "card",
    html="<h1>Hello from Pulumi</h1>",
    css="h1 { font-family: Inter; padding: 48px; }",
    google_fonts=["Inter"],
)

pulumi.export("imageUrl", card.image_url)
```

Run `pulumi preview` and `pulumi up`. Open the returned URL to render the image. Creation and refresh read/save metadata without rendering bytes. Input changes replace image definitions; template edits create new versions under a stable ID.
