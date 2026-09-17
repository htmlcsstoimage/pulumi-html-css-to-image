# ImageHtmlCss

Create a saved HTML/CSS image definition. Creation does not render bytes. Input changes replace the image; each create disables deduplication to guarantee independent ownership. Refresh reads metadata without calling a rendering URL. Destroy accepts asynchronous deletion.

**Resource token:** `html-css-to-image:index:ImageHtmlCss`

See [provider setup](../../README.md#development) and the [complete YAML example](../../examples/images/Pulumi.yaml). The full lifecycle needs `images:read`, `images:create`, and `images:delete`, plus permission to use referenced resources.

Unspecified nullable rendering settings are sent as explicit JSON null. The API owns rendering defaults; the provider does not fill in omitted inputs. Removing a setting resets it to API-default behavior.

## Inputs

| Input | Type | Required | Description |
| --- | --- | --- | --- |
| `colorScheme` | String | No | Sets the preferred color scheme: light or dark. |
| `css` | String | No | CSS used to style the rendered HTML. |
| `deviceScale` | Number | No | Adjusts the pixel ratio used for the screenshot. Minimum: 0.1. Maximum: 3. HTML and template renders default to 2; URL renders default to 1. |
| `disableTwemoji` | Boolean | No | Disables the Twemoji fallback and renders emoji using native fonts instead. |
| `format` | String | No | Rendering URL format: png, jpg, jpeg, webp, or pdf. Does not restrict later renders to this format. |
| `googleFonts` | Array of String | No | Google fonts to load. Separate multiple fonts with a pipe, such as 'Roboto&#124;OpenSans', and set font-family in the CSS to use them. |
| `html` | String | Yes | HTML to render and take a screenshot of. HTML fragments are rendered in a wrapper document unless a complete HTML document is supplied. Required for HTML image requests. |
| `jumboMaxHeight` | Integer | No | Maximum height of the rendered image in jumbo mode. Jumbo rendering consumes additional renders and requires jumbo_max_width. Supply both jumbo dimensions. Each must be greater than 0 and no more than 80000; at least one must exceed 8000; total area cannot exceed 400000000 pixels. |
| `jumboMaxWidth` | Integer | No | Maximum width of the rendered image in jumbo mode. Jumbo rendering consumes additional renders and requires jumbo_max_height. Supply both jumbo dimensions. Each must be greater than 0 and no more than 80000; at least one must exceed 8000; total area cannot exceed 400000000 pixels. |
| `maxRenderOnce` | Boolean | No | Ensure the image is only ever rendered and saved one time. This is an advanced option not applicable to most requests. |
| `maxWaitMs` | Integer | No | Sets a limit on how long to wait before taking the screenshot when the page continues loading irrelevant content. Minimum: 500. Maximum: 10000 and subject to the account plan limit. |
| `mediaType` | String | No | CSS media type to emulate: screen or print. |
| `metadata` | Map of String | No | Custom key-value metadata stored with the image. |
| `msDelay` | Integer | No | Adds extra time in milliseconds before taking the screenshot so JavaScript can execute. Minimum: 0. Maximum: 10000. |
| `pdfOptions` | [ImageHtmlCssPdfOptions](#nested-objects) | No | Rendering option. |
| `proxyId` | String | No | Specifies which configured organization proxy to use when rendering. |
| `renderWhenReady` | Boolean | No | Waits until the page signals that the screenshot is ready. The image fails if the readiness signal is never sent. |
| `selector` | String | No | A CSS selector for an element in the HTML. We’ll crop the image to this specific element. |
| `storageDestinationId` | String | No | Specifies which configured organization storage destination receives the rendered image. |
| `timezone` | String | No | Sets the IANA timezone used by Chrome while rendering. Must be a recognized IANA timezone identifier. |
| `transparentBackground` | Boolean | No | Specifies whether the image is rendered with a transparent background. |
| `viewportHeight` | Integer | No | Sets the height of Chrome's viewport and disables automatic cropping. Minimum: 1. Maximum: 6000. Both viewport dimensions must be supplied together. |
| `viewportLandscape` | Boolean | No | Specifies whether the emulated viewport is in landscape orientation. |
| `viewportMobile` | Boolean | No | Specifies whether the page uses mobile viewport behavior, including its viewport meta tag. |
| `viewportTouch` | Boolean | No | Specifies whether the emulated viewport supports touch events. |
| `viewportWidth` | Integer | No | Sets the width of Chrome's viewport and disables automatic cropping. Minimum: 1. Maximum: 6000. Both viewport dimensions must be supplied together. |

## Additional computed attributes

All inputs above are also readable resource outputs. Omitted nullable rendering inputs stay unset; they do not report effective defaults. The following are additional computed outputs.

| Attribute | Type | Description |
| --- | --- | --- |
| `id` | String | Resource identity used for import. |
| `createdAt` | String | Creation timestamp. |
| `imageUrl` | String | Rendering operation URL. Creation does not render bytes. |
| `lastRenderStoredAt` | String | Last HCTI stored-render timestamp, when available. |
| `ogConfigContentVersion` | String | Associated OG content version, if any. Represented as a decimal string in Pulumi to preserve int64 precision. |
| `ogConfigId` | String | Associated OG configuration ID, if any. |
| `publicUrl` | String | Public GET rendering URL, or null when HCTI storage is disabled. |
| `renderMethod` | String | HTTP method required to render: GET or PUT. |
| `renderRequiresAuth` | Boolean | Whether rendering requires API authentication. |
| `savedToStorageDestinationAt` | String | Last base-image save to custom storage, when available. |
| `storageDestinationHctiStorageDisabled` | Boolean | Whether rendered output is saved only to custom storage. |

## Nested objects

### ImageHtmlCssPdfOptions

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `margins` | [ImageHtmlCssPdfOptionsMargins](#nested-objects) | No | PDF margins. Supply all four sides with units.  |
| `pageHeight` | String | No | Page height with px, in, cm, mm, or pt units.  |
| `pageWidth` | String | No | Page width with px, in, cm, mm, or pt units.  |
| `printBackground` | Boolean | No | Print page backgrounds. Omission lets the API choose.  |
| `scale` | Number | No | PDF scale from 0.1 to 2.  |

### ImageHtmlCssPdfOptionsMargins

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `bottom` | String | Yes | Margin with px, in, cm, mm, or pt units.  |
| `left` | String | Yes | Margin with px, in, cm, mm, or pt units.  |
| `right` | String | Yes | Margin with px, in, cm, mm, or pt units.  |
| `top` | String | Yes | Margin with px, in, cm, mm, or pt units.  |

## Import

```sh
pulumi import html-css-to-image:index:ImageHtmlCss example IMAGE_ID
```

`imageUrl` is an operation URL. Normal output uses public GET; custom-storage-only output requires authenticated PUT to `/v1/store/{id}` and has no `publicUrl`. Creation/refresh do not invoke either operation. Import cannot recover the requested `format`, which is not saved in image metadata; setting one afterward replaces the image.
