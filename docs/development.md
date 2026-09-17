# Development

Requires Go 1.25.11+, Pulumi CLI 3.245.0+, and the language toolchains for SDK compilation. The Terraform provider is pinned to published `v0.1.0`; no sibling checkout is required.

```sh
GOWORK=off make build
GOWORK=off make test
GOWORK=off make vet
GOWORK=off make sdk
```

Bridge tests use fixture APIs and dummy credentials. They cover lifecycle behavior, secret propagation, import, replacement, template lookups and version precision. The framework adapter presents template versions as strings and translates them to exact Terraform int64 values without float64 rounding.

All generated SDKs are ignored on `main`. Release CI publishes Go source under `sdk/go` on the dedicated `sdk` branch; module tags point to those generated commits. The other SDKs are uploaded to their language registries. CI runs one SDK per job through `.github/workflows/sdk.yaml`. To build one locally, run `make sdk-dotnet` and `bash scripts/package-sdk.sh dotnet` (or use `nodejs`, `python`, `go`, or `java`). Run `bash scripts/package-sdks.sh` to build and validate all five SDKs before release; this requires Node 24, Python 3.12, .NET 10, JDK 21, and Maven 3.9+. Java targets Java 11. `packaging/java/pom.xml` supplies the Maven build and Central publishing configuration; `scripts/prepare-java-sdk.mjs` adds the generated SDK version and plugin metadata. To use a development plugin, install the built binary into an isolated `PULUMI_HOME` with `pulumi plugin install resource html-css-to-image 0.1.2 --file bin/pulumi-resource-html-css-to-image`. Use a private local backend and test credentials for live tests. Do not commit state or credentials.

Update mappings in `internal/provider/resources.go`, then regenerate schema and SDKs after an upstream release. `VERSION` must match the generated schema and release binary. See [releasing](releasing.md) for the GitHub release and Registry setup.
