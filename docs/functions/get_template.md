# getTemplate

Read an existing HTML/CSS template by `id`, optionally pinned with `version`. Requires `templates:read`.

**Function token:** `html-css-to-image:index:getTemplate`

```yaml
variables:
  template:
    fn::invoke:
      function: html-css-to-image:index:getTemplate
      arguments:
        id: t-YOUR_TEMPLATE_ID
outputs:
  templateVersion: ${template.version}
  templateHtml: ${template.html}
```

Omit `version` to read the latest version on each invocation. Pass a decimal string to pin a specific version. The returned `version` is also a string, preserving all int64 digits. The result includes saved rendering settings, template type, and timestamps. Unset settings remain null.

Referencing `template.version` in an image's `templateVersion` makes template changes trigger image replacement. Lookups do not manage or render the template. Missing templates, missing versions, and authorization failures are errors.
