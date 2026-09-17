# getTemplateVersions

List versions of an existing template, newest first. Requires `templates:read`.

**Function token:** `html-css-to-image:index:getTemplateVersions`

```yaml
variables:
  history:
    fn::invoke:
      function: html-css-to-image:index:getTemplateVersions
      arguments:
        id: t-YOUR_TEMPLATE_ID
        limit: 10
outputs:
  templateVersions: ${history.versions}
```

`limit` is optional, defaults to **1000**, and must be positive. Pagination is automatic and stops at the limit or the end of the history. The result includes the effective `limit` and a `versions` array of version strings, names, descriptions, template types, and timestamps. It excludes template contents. Each version is an exact decimal string, including identifiers above JavaScript's safe integer range.
