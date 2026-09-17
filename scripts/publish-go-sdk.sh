#!/usr/bin/env bash
# Publish generated Go sources without changing the provider checkout or its index.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
version="$(cat VERSION)"
[[ "$version" =~ ^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-(alpha|beta|rc)\.(0|[1-9][0-9]*))?$ ]] || { echo 'Invalid VERSION' >&2; exit 1; }
test -f sdk/go/go.mod
source_commit="$(git rev-parse HEAD)"
tag="sdk/go/v${version}"
export GIT_AUTHOR_NAME='github-actions[bot]'
export GIT_AUTHOR_EMAIL='41898282+github-actions[bot]@users.noreply.github.com'
export GIT_COMMITTER_NAME="$GIT_AUTHOR_NAME"
export GIT_COMMITTER_EMAIL="$GIT_AUTHOR_EMAIL"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
export GIT_INDEX_FILE="$tmp/index"
git read-tree --empty
git add --force -- sdk/go LICENSE
# A nested Go module must carry its own license.
license_blob="$(git hash-object -w LICENSE)"
git update-index --add --cacheinfo "100644,$license_blob,sdk/go/LICENSE"
source_blob="$(printf '%s\n' "$source_commit" | git hash-object -w --stdin)"
git update-index --add --cacheinfo "100644,$source_blob,PROVIDER_SOURCE"
tree="$(git write-tree)"
remote_exists() {
  local status=0
  git ls-remote --exit-code origin "$1" > /dev/null || status=$?
  case "$status" in
    0) return 0 ;;
    2) return 1 ;;
    *) echo 'Could not inspect remote refs' >&2; exit "$status" ;;
  esac
}
if remote_exists "refs/tags/$tag"; then
  git fetch --no-tags origin "refs/tags/$tag"
  if [[ "$(git rev-parse 'FETCH_HEAD^{tree}')" != "$tree" ]]; then
    echo "$tag already exists with different sources or provider commit; refusing to replace it" >&2
    exit 1
  fi
  echo "$tag already published with identical sources"
  exit 0
fi
set --
if remote_exists refs/heads/sdk; then
  git fetch --no-tags origin refs/heads/sdk
  set -- -p "$(git rev-parse FETCH_HEAD)"
fi
commit="$(printf 'Go SDK v%s\n\nProvider source: %s\n' "$version" "$source_commit" | git commit-tree "$tree" "$@")"
# Both refs advance together; concurrent branch changes are rejected without force.
git push --atomic origin "$commit:refs/heads/sdk" "$commit:refs/tags/$tag"
