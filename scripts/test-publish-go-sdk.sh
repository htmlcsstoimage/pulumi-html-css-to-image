#!/usr/bin/env bash
set -euo pipefail
publisher="$(cd "$(dirname "$0")" && pwd)/publish-go-sdk.sh"
fixture="$(mktemp -d)"
trap 'rm -rf "$fixture"' EXIT
export GIT_AUTHOR_NAME='SDK test' GIT_AUTHOR_EMAIL='sdk@example.invalid'
export GIT_COMMITTER_NAME="$GIT_AUTHOR_NAME" GIT_COMMITTER_EMAIL="$GIT_AUTHOR_EMAIL"
git init --bare -q "$fixture/remote.git"
git init -q "$fixture/work"
cd "$fixture/work"
git remote add origin "$fixture/remote.git"
printf '0.1.0\n' > VERSION
printf 'Test license\n' > LICENSE
printf '/sdk/\n' > .gitignore
git add VERSION LICENSE .gitignore
git commit -qm 'Provider source'
head="$(git rev-parse HEAD)"
mkdir -p sdk/go
printf 'module example.com/provider/sdk/go\n\ngo 1.21\n' > sdk/go/go.mod
printf 'package sdk\n' > sdk/go/old.go
# Include staged changes to verify that the publisher does not touch the real index.
printf 'keep staged\n' > staged.txt
git add staged.txt
index="$(git write-tree)"
bash "$publisher"
git fetch -q origin sdk refs/tags/sdk/go/v0.1.0
first="$(git rev-parse FETCH_HEAD)"
test "$(git rev-parse HEAD)" = "$head"
test "$(git write-tree)" = "$index"
test "$(git show "$first:PROVIDER_SOURCE")" = "$head"
test "$(git show "$first:sdk/go/LICENSE")" = 'Test license'
bash "$publisher"
printf '// changed\n' >> sdk/go/old.go
if bash "$publisher" > "$fixture/conflict.log" 2>&1; then
  echo 'Conflicting existing tag was accepted' >&2; exit 1
fi
printf '0.1.1\n' > VERSION
git add VERSION
git commit -qm 'Next provider version'
rm sdk/go/old.go
printf 'package sdk\n' > sdk/go/new.go
bash "$publisher"
git fetch -q origin sdk
second="$(git rev-parse FETCH_HEAD)"
test "$(git rev-parse "$second^")" = "$first"
if git cat-file -e "$second:sdk/go/old.go" 2>/dev/null; then
  echo 'Deleted generated source survived' >&2; exit 1
fi
git cat-file -e "$second:sdk/go/new.go"
# A failed tag push must not advance the branch either.
cat > "$fixture/remote.git/hooks/pre-receive" <<'HOOK'
#!/usr/bin/env bash
while read -r old new ref; do
  case "$ref" in refs/tags/*) exit 1 ;; esac
done
HOOK
chmod +x "$fixture/remote.git/hooks/pre-receive"
printf '0.1.2\n' > VERSION
git add VERSION
git commit -qm 'Rejected release'
if bash "$publisher" > "$fixture/rejected.log" 2>&1; then
  echo 'Rejected tag push succeeded' >&2; exit 1
fi
test "$(git --git-dir="$fixture/remote.git" rev-parse refs/heads/sdk)" = "$second"
echo 'Go SDK publishing: first release, retry, conflict, next release, index isolation, and atomic rejection passed.'
