#!/usr/bin/env bash
# Check the bridge's directory-name convention cannot change generated docs.
set -euo pipefail
cd "$(dirname "$0")/.."
fixture="$(mktemp -d)"
trap 'rm -rf "$fixture"' EXIT
export GOWORK=off
go build -o "$fixture/schema" ./cmd/schema
version="$(cat VERSION)"
for name in provider custom-checkout; do
  mkdir -p "$fixture/$name/internal/provider"
  cp go.mod go.sum "$fixture/$name/"
  (
    cd "$fixture/$name"
    "$fixture/schema" -version "$version"
  )
  cmp internal/provider/schema.json "$fixture/$name/internal/provider/schema.json"
done
echo 'Schema matches in provider and custom-checkout directories.'
