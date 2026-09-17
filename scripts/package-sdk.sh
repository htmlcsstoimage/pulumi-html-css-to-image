#!/usr/bin/env bash
# Build one generated SDK; each language can run on its own CI runner.
set -euo pipefail
cd "$(dirname "$0")/.."
language="${1:?Specify nodejs, python, dotnet, go, or java}"
export HCTI_RELEASE_VERSION="$(cat VERSION)"
node scripts/check-sdk-version.mjs "$language"
case "$language" in
  nodejs)
    cd sdk/nodejs
    npm install --ignore-scripts
    npm run build
    cp package.json README.md ../../LICENSE bin/
    npm pack ./bin --dry-run
    ;;
  python)
    python -m pip install build twine
    python -m build sdk/python
    python -m twine check sdk/python/dist/*
    ;;
  dotnet)
    # GeneratePackageOnBuild suppresses pack's implicit build on fresh runners.
    dotnet pack sdk/dotnet/Pulumi.HtmlCssToImage.csproj --configuration Release -p:GeneratePackageOnBuild=false --output sdk/dotnet/artifacts
    ;;
  go) GOWORK=off go -C sdk/go test ./... ;;
  java) mvn --batch-mode --no-transfer-progress -f sdk/java/pom.xml clean verify ;;
  *) echo "Unknown SDK: $language" >&2; exit 1 ;;
esac
