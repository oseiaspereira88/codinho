#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/../.."
# Include new product/test Go files before staging. Malformed editorial fixtures
# live below testdata and are inputs, not formatting targets.
failed=0
while IFS= read -r -d '' file; do
  case "$file" in */testdata/*|testdata/*|packs/*) continue ;; esac
  if ! output=$(gofmt -l "$file"); then exit 1; fi
  if [ -n "$output" ]; then printf '%s\n' "$output"; failed=1; fi
done < <(git ls-files --cached --others --exclude-standard -z -- '*.go')
exit "$failed"
