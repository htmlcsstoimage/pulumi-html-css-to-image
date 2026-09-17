# ImageTemplated

Create a saved image from a template and values. `templateVersion` is optional: omission selects latest at creation and does not replace the image when the template later changes. Reference a managed template’s `version` to replace images with that template. `resolvedTemplateVersion` reports the actual saved version without inserting a pin. Inputs replace the image, deduplication is disabled, and refresh never renders it.

**Resource token:** `html-css-to-image:index:ImageTemplated`

See [provider setup](../../README.md#development) and the [complete YAML example](../../examples/images/Pulumi.yaml). The full lifecycle needs `images:read`, `images:create`, and `images:delete`, plus permission to use referenced resources.

Version identifiers are **decimal strings** in Pulumi, preserving the full API int64 range in every language. Use `"9007199254740993"`, not a JavaScript number. References such as `${design.version}` already have the correct type.

`templateValues` is a sensitive JSON object **string**. Use `JSON.stringify` in TypeScript or your language’s JSON serializer. Equality ignores key ordering and whitespace while preserving numeric precision. Supply very large JSON numbers as exact JSON text where the language would round them.

## Inputs

| Input | Type | Required | Description |
| --- | --- | --- | --- |
| `format` | String | No | Rendering URL format: png, jpg, jpeg, webp, or pdf. Does not restrict later renders to this format. |
| `templateId` | String | Yes | Template ID including the t- prefix. |
| `templateValues` | String | Yes | Values substituted into the template for this render. Must be a non-empty JSON object. Include the values needed by the template. |
| `templateVersion` | String | No | Optional version to pin. Omission selects the latest version when creating the image. Represented as a decimal string in Pulumi to preserve int64 precision. |

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
| `resolvedTemplateVersion` | String | Actual template version used at creation; does not pin an omitted template_version. Represented as a decimal string in Pulumi to preserve int64 precision. |
| `storageDestinationHctiStorageDisabled` | Boolean | Whether HCTI storage was disabled when the image was created. |
| `savedToStorageDestinationAt` | String | Last base-image save to custom storage, when available. |

## Import

```sh
pulumi import html-css-to-image:index:ImageTemplated example IMAGE_ID
```

`imageUrl` is an operation URL. Normal output uses public GET; custom-storage-only output requires authenticated PUT to `/v1/store/{id}` and has no `publicUrl`. Creation/refresh do not invoke either operation. Import cannot recover the requested `format`, which is not saved in image metadata; setting one afterward replaces the image.

Import leaves `templateVersion` unset and fills `resolvedTemplateVersion`. The rendering URL’s mode comes directly from the image’s saved storage metadata. No template or storage-destination reads are needed.
