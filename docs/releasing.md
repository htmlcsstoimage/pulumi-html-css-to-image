# Releasing

`VERSION` controls the plugin and schema version. Push a new version to `main`; after **CI** succeeds, **Release** checks out that tested commit and creates its `v`-prefixed tag. GoReleaser uploads six platform archives, checksums, and `schema.json` to a draft, then publishes it after all uploads succeed. Published versions are skipped and never replaced. An unpublished tag can only be retried from the original commit.

The plugin download location is `github://api.github.com/htmlcsstoimage/pulumi-html-css-to-image`. GitHub Actions uses its built-in token; this flow requires no Terraform signing key or registry token. Set repository visibility to public before publishing.

## Registry listing

After the first GitHub release succeeds, submit the package to Pulumi Registry using this entry in `pulumi/registry`'s `community-packages/package-list.json`:

```json
{
  "repoSlug": "htmlcsstoimage/pulumi-html-css-to-image",
  "schemaFile": "internal/provider/schema.json"
}
```

Follow [Pulumi's registration instructions](https://www.pulumi.com/docs/iac/guides/building-extending/packages/publishing-packages/). The publisher is `htmlcsstoimage`; the package name is `html-css-to-image`. Include a company contact and register the publisher mapping if requested.

## Language package publishing

The same `release.yml` workflow generates SDKs, checks the committed schema for drift, and builds all packages before creating a release. It uses Pulumi's [package publisher](https://github.com/pulumi/pulumi-package-publisher) as a packaging reference, with registry-native trusted publishing for authentication.

Configure these identities for owner **htmlcsstoimage**, repository **pulumi-html-css-to-image**, workflow **release.yml**, with no environment name:

| Registry | Package | Setup |
| --- | --- | --- |
| npm | `@html-css-to-image/pulumi` | Add a GitHub Actions trusted publisher in package settings. A new package needs its initial publication before configuring that policy. |
| PyPI | `pulumi-html-css-to-image` | Add a pending trusted publisher for the new project, or a publisher on the existing project. |
| NuGet | `Pulumi.HtmlCssToImage` | Create a trusted publishing policy for this package and make `NUGET_USER` available as a repository or organization Actions secret. This is the account name, not a long-lived API key. |
| Maven Central | `com.htmlcsstoimage:pulumi` | Verify `com.htmlcsstoimage`; set `MAVEN_CENTRAL_USERNAME`, `MAVEN_CENTRAL_PASSWORD`, `GPG_PRIVATE_KEY`, and `GPG_PASSPHRASE`. |
| Go | `github.com/htmlcsstoimage/pulumi-html-css-to-image/sdk/go` | CI commits generated `sdk/go` sources to the `sdk` branch and creates the matching `sdk/go/vVERSION` tag there. |

Use the existing HCTI registry accounts. Their existing client policies do not automatically authorize a new repository or package. npm uses Node 24 with OIDC and provenance; PyPI uses `pypa/gh-action-pypi-publish`; NuGet uses `NuGet/login` to obtain a short-lived publishing key. No npm, PyPI, or NuGet publishing token is stored in GitHub.

Maven Central uses Portal token credentials and GPG signing. The Java SDK includes the main JAR, sources, Javadocs, and POM. Publish the signing **public** key to a [Central-supported keyserver](https://central.sonatype.org/publish/requirements/gpg/); registering it with Terraform alone does not do this. The workflow waits for Central publication and only skips a version when its POM is already available on Maven Central. If a previous deployment is still processing, check its status in the Central Portal before retrying.

Go publishing uses `scripts/publish-go-sdk.sh` to create a commit containing the generated module, license, and `PROVIDER_SOURCE` commit ID. It preserves the provider checkout and index, advances the `sdk` branch and module tag atomically, and refuses to overwrite a tag with different contents. Retrying an identical publication is a no-op. The workflow needs permission to create and update the `sdk` branch and create `sdk/go/*` tags; repository rules must permit these bot pushes.

The GitHub release stays a draft until all language packages and the Go tag are published. A partial registry publication cannot be rolled back; retries skip existing versions and fail on authentication or other publishing errors. npm prereleases use the `next` distribution tag; stable versions use `latest`.

See the registry setup instructions for [npm](https://docs.npmjs.com/trusted-publishers/), [PyPI](https://docs.pypi.org/trusted-publishers/), and [NuGet](https://learn.microsoft.com/en-us/nuget/nuget-org/trusted-publishing).

## Local packaging checks

```sh
GOWORK=off make build
GOWORK=off make test
GOWORK=off make vet
goreleaser check
GOWORK=off goreleaser release --snapshot --clean --skip=publish
```

Never overwrite an existing published tag or release. If a draft upload fails, retry its original workflow; remove only the unpublished draft if partial assets prevent recovery.
