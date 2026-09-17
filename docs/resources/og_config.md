# OgConfig

Manage an Open Graph configuration. All settings update in place. `configType` accepts the literal API values `html_css` and `templated`; these are not camel-cased. HTML/CSS configurations use `defaultOptions`. Templated configurations use `templateId`, optional `templateVersion`, and ordered `templateValuesMappings`. Destroy disables/deletes the configuration.

**Resource token:** `html-css-to-image:index:OgConfig`

See [provider setup](../../README.md#development) and the [complete YAML example](../../examples/og-config/Pulumi.yaml). The full lifecycle needs `og_configs:read`, `og_configs:create_update`, and `og_configs:delete`, plus permission to use referenced resources.

Version identifiers are **decimal strings** in Pulumi, preserving the full API int64 range in every language. Use `"9007199254740993"`, not a JavaScript number. References such as `${design.version}` already have the correct type.

Header maps are always secret outputs, even when their entries were supplied as ordinary strings. The bridge preserves sensitivity through create and refresh. Pulumi encrypts secret values in stack state using the stack’s secrets provider.

## Inputs

| Input | Type | Required | Description |
| --- | --- | --- | --- |
| `additionalHeaderOrigins` | Array of String | No | Templated only. Up to 20 additional exact HTTP(S) origins allowed to receive headers, each up to 512 UTF-8 bytes. Omission clears them. |
| `baseUrl` | String | Yes | HTTPS origin of the source website, up to 2048 characters, without a path, credentials, query, or fragment. |
| `configType` | String | Yes | Rendering source: html_css or templated. Changing type updates this configuration in place. |
| `defaultOptions` | [OgConfigDefaultOptions](#nested-objects) | No | HTML/CSS only. Default page rendering options. Unspecified fields are sent as null; the API owns render defaults. |
| `description` | String | No | Optional description, up to 1023 characters. Omission clears it. |
| `disabled` | Boolean | No | Disable serving images. Defaults to false. |
| `extractValues` | Boolean | No | HTML/CSS only. Extract image options from page metadata, overriding default_options. Omission uses false. |
| `headers` | Map of String | No | Templated only. HTTP extraction headers, up to 20. Names: 512 ASCII characters; values: 8192 UTF-8 bytes. Omission clears them. |
| `name` | String | Yes | Display name, 1–255 characters after trimming. |
| `optimizationMode` | String | No | no_optimization keeps dimensions; post_process adapts after rendering; set_viewport renders at social-platform viewport sizes. Omission lets the API choose its default (currently post_process). |
| `refreshIntervalS` | Integer | No | Seconds before cached images may refresh: 1800–31536000, subject to the plan minimum. Omission lets the API choose its default (currently 86400). |
| `templateId` | String | No | Required for templated configurations. Template ID including its t- prefix. |
| `templateValuesMappings` | Array of [OgConfigTemplateValuesMapping](#nested-objects) | No | Templated only. Up to 32 ordered mappings from page metadata to distinct template fields. |
| `templateVersion` | String | No | Templated only. Positive int64 version. Omission follows the latest version at serving time. Represented as a decimal string in Pulumi to preserve int64 precision. |

## Additional computed attributes

All inputs above are also readable resource outputs. Omitted nullable rendering inputs stay unset; they do not report effective defaults. The following are additional computed outputs.

| Attribute | Type | Description |
| --- | --- | --- |
| `id` | String | Resource identity used for import. |
| `createdAt` | String | Creation timestamp. |
| `domainId` | String | Serving identifier used in OG image URLs. Distinct from the management ID. |
| `effectiveOptimizationMode` | String | Optimization mode returned by the API. |
| `effectiveRefreshIntervalS` | Integer | Refresh interval returned by the API, in seconds. |
| `updatedAt` | String | Last update timestamp. |

## Nested objects

### OgConfigDefaultOptions

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `additionalHeaderOrigins` | Array of String | No | Additional exact HTTP or HTTPS origins allowed to receive custom headers. Supports up to 20 unique origins of up to 512 UTF-8 bytes each. Origins must use the format scheme://host[:port] without a path; duplicates are ignored.  |
| `blockConsentBanners` | Boolean | No | Attempt to block cookie/consent banners from displaying.  |
| `colorScheme` | String | No | Sets Chrome's preferred color scheme. Valid values: light or dark.  |
| `css` | String | No | CSS injected into the loaded page to override its styles.  |
| `deviceScale` | Number | No | Adjusts the pixel ratio used for the screenshot. Minimum: 0.1. Maximum: 3. HTML and template renders default to 2; URL renders default to 1.  |
| `disableTwemoji` | Boolean | No | Disables the Twemoji fallback and renders emoji using native fonts instead.  |
| `headers` | Map of String | No | HTTP headers to include on top-level page navigations to the requested URL's origin and any additional_header_origins. Supports up to 20 headers with names up to 512 ASCII characters and values up to 8192 UTF-8 bytes.  |
| `identifyAsHcti` | Boolean | No | Identify the top-level page navigation as an HCTI screenshot request using the X-HCTI-SCREENSHOT header.  |
| `includeHeadersOnSubrequests` | Boolean | No | Include custom headers on subrequests to the requested URL's origin and any additional_header_origins. Defaults to false. Requires at least one header.  |
| `maxWaitMs` | Integer | No | Sets a limit on how long to wait before taking the screenshot when the page continues loading irrelevant content. Minimum: 500. Maximum: 10000 and subject to the account plan limit.  |
| `mediaType` | String | No | Sets the CSS media type used while rendering the page. Valid values: print or screen.  |
| `msDelay` | Integer | No | Adds extra time in milliseconds before taking the screenshot so JavaScript can execute. Minimum: 0. Maximum: 10000.  |
| `proxyId` | String | No | Configured proxy ID. Must refer to an enabled proxy.  |
| `renderWhenReady` | Boolean | No | Waits until the page signals that the screenshot is ready. The image fails if the readiness signal is never sent.  |
| `selector` | String | No | A CSS selector for an element in the HTML. We’ll crop the image to this specific element.  |
| `storageDestinationId` | String | No | Configured storage destination ID. Must be enabled and allow HCTI storage.  |
| `timezone` | String | No | Sets the IANA timezone used by Chrome while rendering. Must be a recognized IANA timezone identifier.  |
| `transparentBackground` | Boolean | No | Specifies whether the image is rendered with a transparent background.  |
| `viewportHeight` | Integer | No | Sets the height of Chrome's viewport and disables automatic cropping. Minimum: 1. Maximum: 6000. Both viewport dimensions must be supplied together.  |
| `viewportLandscape` | Boolean | No | Specifies whether the emulated viewport is in landscape orientation.  |
| `viewportMobile` | Boolean | No | Specifies whether the page uses mobile viewport behavior, including its viewport meta tag.  |
| `viewportTouch` | Boolean | No | Specifies whether the emulated viewport supports touch events.  |
| `viewportWidth` | Integer | No | Sets the width of Chrome's viewport and disables automatic cropping. Minimum: 1. Maximum: 6000. Both viewport dimensions must be supplied together.  |

### OgConfigTemplateValuesMapping

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `fallback` | String | No | Metadata fallback: titles or descriptions. This selects page metadata, not a literal value.  |
| `metaKey` | String | No | Source meta tag name, up to 136 characters. Supply exactly one of metaKey or fallback.  |
| `templateKey` | String | Yes | Destination template field path, up to 128 characters. Paths must not overlap.  |

## Import

```sh
pulumi import html-css-to-image:index:OgConfig example MANAGEMENT_ID
```
