#!/usr/bin/env bash
# Convenience wrapper for local builds. CI runs package-sdk.sh in parallel jobs.
set -euo pipefail
cd "$(dirname "$0")/.."
for language in nodejs python dotnet go java; do
  bash scripts/package-sdk.sh "$language"
done
