# ImageUrl

Create a saved webpage image definition. `url` is the source webpage; `imageUrl` is the rendering operation URL. Creation does not render bytes. Input changes replace the image; deduplication is disabled. A changed webpage alone does not cause replacement. Refresh uses metadata, never the rendering URL.

**Resource token:** `html-css-to-image:index:ImageUrl`

See [provider setup](../../README.md#development) and the [complete YAML example](../../examples/images/Pulumi.yaml). The full lifecycle needs `images:read`, `images:create`, and `images:delete`, plus permission to use referenced resources.

Unspecified nullable rendering settings are sent as explicit JSON null. The API owns rendering defaults; the provider does not fill in omitted inputs. Removing a setting resets it to API-default behavior.

Header maps are always secret outputs, even when their entries were supplied as ordinary strings. The bridge preserves sensitivity through create and refresh. Pulumi encrypts secret values in stack state using the stack’s secrets provider.

In .NET, the `imageUrl` output is named `RenderingUrl` because C# forbids a member named `ImageUrl` inside the `ImageUrl` class. The input `Url` remains the source webpage.

## Inputs

| Input | Type | Required | Description |
| --- | --- | --- | --- |
| `additionalHeaderOrigins` | Array of String | No | Additional exact HTTP or HTTPS origins allowed to receive custom headers. Supports up to 20 unique origins of up to 512 UTF-8 bytes each. Origins must use the format scheme://host[:port] without a path; duplicates are ignored. For GET and form-encoded requests, repeat this parameter for each origin. |
| `blockConsentBanners` | Boolean | No | Attempt to block cookie/consent banners from displaying. |
| `colorScheme` | String | No | Sets the preferred color scheme: light or dark. |
| `css` | String | No | CSS injected into the loaded URL to override styles on the page. |
| `deviceScale` | Number | No | Adjusts the pixel ratio used for the screenshot. Minimum: 0.1. Maximum: 3. HTML and template renders default to 2; URL renders default to 1. |
| `disableTwemoji` | Boolean | No | Disables the Twemoji fallback and renders emoji using native fonts instead. |
| `format` | String | No | Rendering URL format: png, jpg, jpeg, webp, or pdf. Does not restrict later renders to this format. |
| `fullScreen` | Boolean | No | Take a screenshot of the entire screen after scrolling down and back to the top. |
| `headers` | Map of String | No | HTTP headers to include on top-level page navigations to the requested URL's origin and any additional_header_origins. Supports up to 20 headers with names up to 512 ASCII characters and values up to 8192 UTF-8 bytes. For GET and form-encoded requests, repeat this parameter using the format `headers=name:value`. |
| `identifyAsHcti` | Boolean | No | Identify the top-level page navigation as an HCTI screenshot request using the X-HCTI-SCREENSHOT header. |
| `includeHeadersOnSubrequests` | Boolean | No | Include custom headers on subrequests to the requested URL's origin and any additional_header_origins. Defaults to false. Requires at least one header. |
| `jumboMaxHeight` | Integer | No | Maximum height of the rendered image in jumbo mode. Jumbo rendering consumes additional renders and requires jumbo_max_width. Supply both jumbo dimensions. Each must be greater than 0 and no more than 80000; at least one must exceed 8000; total area cannot exceed 400000000 pixels. |
| `jumboMaxWidth` | Integer | No | Maximum width of the rendered image in jumbo mode. Jumbo rendering consumes additional renders and requires jumbo_max_height. Supply both jumbo dimensions. Each must be greater than 0 and no more than 80000; at least one must exceed 8000; total area cannot exceed 400000000 pixels. |
| `maxRenderOnce` | Boolean | No | Ensure the image is only ever rendered and saved one time. This is an advanced option not applicable to most requests. |
| `maxWaitMs` | Integer | No | Sets a limit on how long to wait before taking the screenshot when the page continues loading irrelevant content. Minimum: 500. Maximum: 10000 and subject to the account plan limit. |
| `mediaType` | String | No | CSS media type to emulate: screen or print. |
| `metadata` | Map of String | No | Custom key-value metadata stored with the image. |
| `msDelay` | Integer | No | Adds extra time in milliseconds before taking the screenshot so JavaScript can execute. Minimum: 0. Maximum: 10000. |
| `pdfOptions` | [ImageUrlPdfOptions](#nested-objects) | No | Rendering option. |
| `proxyId` | String | No | Specifies which configured organization proxy to use when rendering. |
| `renderWhenReady` | Boolean | No | Waits until the page signals that the screenshot is ready. The image fails if the readiness signal is never sent. |
| `selector` | String | No | A CSS selector for an element in the HTML. We’ll crop the image to this specific element. |
| `storageDestinationId` | String | No | Specifies which configured organization storage destination receives the rendered image. |
| `timezone` | String | No | Sets the IANA timezone used by Chrome while rendering. Must be a recognized IANA timezone identifier. |
| `transparentBackground` | Boolean | No | Specifies whether the image is rendered with a transparent background. |
| `url` | String | Yes | Public HTTP or HTTPS URL to capture. Required for URL image requests. |
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

### ImageUrlPdfOptions

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `margins` | [ImageUrlPdfOptionsMargins](#nested-objects) | No | PDF margins. Supply all four sides with units.  |
| `pageHeight` | String | No | Page height with px, in, cm, mm, or pt units.  |
| `pageWidth` | String | No | Page width with px, in, cm, mm, or pt units.  |
| `printBackground` | Boolean | No | Print page backgrounds. Omission lets the API choose.  |
| `scale` | Number | No | PDF scale from 0.1 to 2.  |

### ImageUrlPdfOptionsMargins

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `bottom` | String | Yes | Margin with px, in, cm, mm, or pt units.  |
| `left` | String | Yes | Margin with px, in, cm, mm, or pt units.  |
| `right` | String | Yes | Margin with px, in, cm, mm, or pt units.  |
| `top` | String | Yes | Margin with px, in, cm, mm, or pt units.  |

## Import

```sh
pulumi import html-css-to-image:index:ImageUrl example IMAGE_ID
```

`imageUrl` is an operation URL. Normal output uses public GET; custom-storage-only output requires authenticated PUT to `/v1/store/{id}` and has no `publicUrl`. Creation/refresh do not invoke either operation. Import cannot recover the requested `format`, which is not saved in image metadata; setting one afterward replaces the image.
