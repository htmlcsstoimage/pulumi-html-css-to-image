## Example

In a TypeScript Pulumi project, install `@html-css-to-image/pulumi` and configure your credentials:

```sh
pulumi config set html-css-to-image:apiId YOUR_API_ID
pulumi config set --secret html-css-to-image:apiKey
```

Add this to `index.ts`:

```typescript
import * as hcti from "@html-css-to-image/pulumi";

const card = new hcti.ImageHtmlCss("card", {
    html: "<h1>Hello from Pulumi</h1>",
    css: "h1 { font-family: Inter; padding: 48px; }",
    googleFonts: ["Inter"],
});

export const imageUrl = card.imageUrl;
```

Run `pulumi preview` and `pulumi up`. Open the returned URL to render the image. Creation and refresh read/save metadata without rendering bytes. Input changes replace image definitions; template edits create new versions under a stable ID.
