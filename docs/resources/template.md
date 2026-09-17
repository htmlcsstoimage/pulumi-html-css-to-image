# Template

Manage an HTML/CSS template. Edits create a new version under the same ID, including changes to its name and description. Refresh reads the latest version and detects external edits; unchanged content does not create another version. Destroy deletes all versions. Block templates cannot be managed or imported into this resource.

**Resource token:** `html-css-to-image:index:Template`

See [provider setup](../../README.md#development) and the [complete YAML example](../../examples/template/Pulumi.yaml). The full lifecycle needs `templates:read`, `templates:create_update`, and `templates:delete`, plus permission to use referenced resources.

Unspecified nullable rendering settings are sent as explicit JSON null. The API owns rendering defaults; the provider does not fill in omitted inputs. Removing a setting resets it to API-default behavior.

Version identifiers are **decimal strings** in Pulumi, preserving the full API int64 range in every language. Use `"9007199254740993"`, not a JavaScript number. References such as `${design.version}` already have the correct type.

## Inputs

| Input | Type | Required | Description |
| --- | --- | --- | --- |
| `colorScheme` | String | No | Sets the preferred color scheme: light or dark. |
| `css` | String | No | CSS used to style images rendered from the template. Handlebars expressions are not supported in CSS; put dynamic CSS in the HTML instead. |
| `description` | String | No | An optional description of the template for your reference. Maximum: 1024 characters. |
| `deviceScale` | Number | No | Adjusts the pixel ratio used for the screenshot. Minimum: 0.1. Maximum: 3. HTML and template renders default to 2; URL renders default to 1. |
| `disableTwemoji` | Boolean | No | Disables the Twemoji fallback and renders emoji using native fonts instead. |
| `googleFonts` | Array of String | No | Google fonts to load. Separate multiple fonts with a pipe, such as 'Roboto&#124;OpenSans', and set font-family in the CSS to use them. |
| `html` | String | Yes | HTML to render for the template. Use Handlebars placeholders for values that will be substituted when an image is rendered. Must be non-empty, contain at least one Handlebars placeholder, and compile as valid Handlebars. |
| `jumboMaxHeight` | Integer | No | Maximum height of the rendered image in jumbo mode. Jumbo rendering consumes additional renders and requires jumbo_max_width. Supply both jumbo dimensions. Each must be greater than 0 and no more than 80000; at least one must exceed 8000; total area cannot exceed 400000000 pixels. |
| `jumboMaxWidth` | Integer | No | Maximum width of the rendered image in jumbo mode. Jumbo rendering consumes additional renders and requires jumbo_max_height. Supply both jumbo dimensions. Each must be greater than 0 and no more than 80000; at least one must exceed 8000; total area cannot exceed 400000000 pixels. |
| `maxWaitMs` | Integer | No | Sets a limit on how long to wait before taking the screenshot when the page continues loading irrelevant content. Minimum: 500. Maximum: 10000 and subject to the account plan limit. |
| `mediaType` | String | No | CSS media type to emulate: screen or print. |
| `msDelay` | Integer | No | Adds extra time in milliseconds before taking the screenshot so JavaScript can execute. Minimum: 0. Maximum: 10000. |
| `name` | String | No | The name of the template, used to identify it in your account. Maximum: 64 characters. |
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
| `templateType` | String | Template type; this resource manages html_css templates. |
| `updatedAt` | String | Last update timestamp. |
| `version` | String | Latest saved template version. Represented as a decimal string in Pulumi to preserve int64 precision. |

## Import

```sh
pulumi import html-css-to-image:index:Template example t-TEMPLATE_ID
```
