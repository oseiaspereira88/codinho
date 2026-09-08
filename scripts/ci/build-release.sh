#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/../.."
version=${1:-}
# Closed grammar; no spaces, substitutions, shell metacharacters or ldflag tokens.
if [[ ! "$version" =~ ^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-[A-Za-z0-9]+([.-][A-Za-z0-9]+)*)?$ ]] || [ "${#version}" -gt 80 ]; then
  echo 'Invalid build version; expected vMAJOR.MINOR.PATCH[-PRERELEASE]' >&2
  exit 2
fi
os=${GOOS:-$(go env GOOS)}
arch=${GOARCH:-$(go env GOARCH)}
case "$os/$arch" in linux/amd64|linux/arm64|darwin/amd64|darwin/arm64) ;; *) echo 'Unsupported build target' >&2; exit 2 ;; esac
mkdir -p dist
out="codinho-${version}-${os}-${arch}"
CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" go build -trimpath \
  -ldflags "-X github.com/oseiaspereira88/codinho/internal/cli.version=$version" \
  -o "dist/$out" ./cmd/codinho
(cd dist && shasum -a 256 "$out" > "$out.sha256")
printf 'Built %s/%s; this is build evidence, not native execution evidence.\n' "$os" "$arch"
