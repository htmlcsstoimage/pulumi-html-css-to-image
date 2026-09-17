#!/usr/bin/env bash
# Build the generated SDKs into the layouts used by Pulumi's package publisher.
set -euo pipefail
cd "$(dirname "$0")/.."
version="$(cat VERSION)"
export HCTI_RELEASE_VERSION="$version"
node <<'JS'
const fs = require('fs');
for (const file of ['internal/provider/schema.json', 'sdk/nodejs/package.json']) {
  const actual = JSON.parse(fs.readFileSync(file)).version;
  if (actual !== process.env.HCTI_RELEASE_VERSION) throw new Error(`${file}: unexpected version ${actual}`);
}
const version = process.env.HCTI_RELEASE_VERSION;
const pep440 = version.replace(/-(alpha|beta|rc)\.(\d+)$/, (_, kind, n) => ({alpha: 'a', beta: 'b', rc: 'rc'}[kind]) + n);
const pythonVersion = fs.readFileSync('sdk/python/setup.py', 'utf8').match(/^VERSION = "([^"]+)"/m)?.[1];
if (pythonVersion !== pep440) throw new Error(`Python version ${pythonVersion} does not match ${pep440}`);
const dotnetVersion = fs.readFileSync('sdk/dotnet/Pulumi.HtmlCssToImage.csproj', 'utf8').match(/<Version>([^<]+)<\/Version>/)?.[1];
if (dotnetVersion !== version) throw new Error(`.NET version ${dotnetVersion} does not match ${version}`);
JS
(
  cd sdk/nodejs
  npm install --ignore-scripts
  npm run build
  cp package.json README.md ../../LICENSE bin/
  npm pack ./bin --dry-run
)
# The Python generator converts prerelease versions to PEP 440 (e.g. 0.1.0b1).
python -m pip install build twine
python -m build sdk/python
python -m twine check sdk/python/dist/*
test "$(cat sdk/dotnet/version.txt)" = "$version"
dotnet pack sdk/dotnet/Pulumi.HtmlCssToImage.csproj --configuration Release --output sdk/dotnet/artifacts
GOWORK=off go -C sdk/go test ./...
node scripts/prepare-java-sdk.mjs
mvn --batch-mode --no-transfer-progress -f sdk/java/pom.xml clean verify
