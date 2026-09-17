#!/usr/bin/env bash
set -euo pipefail
version="$(cat VERSION)"
package="sdk/dotnet/artifacts/HtmlCssToImage.Pulumi.${version}.nupkg"
test -f "$package"
status="$(curl --silent --show-error --output /dev/null --write-out '%{http_code}' \
  "https://api.nuget.org/v3-flatcontainer/htmlcsstoimage.pulumi/${version}/htmlcsstoimage.pulumi.nuspec")"
case "$status" in
  200) echo "HtmlCssToImage.Pulumi ${version} is already available on NuGet; skipping." ;;
  404)
    # Do not use --skip-duplicate: reserved package IDs also return HTTP 409.
    dotnet nuget push "$package" --api-key "${NUGET_API_KEY:?}" --source https://api.nuget.org/v3/index.json
    ;;
  *) echo "NuGet version lookup failed: HTTP $status" >&2; exit 1 ;;
esac
