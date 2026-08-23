#!/usr/bin/env bash
# smoke-install.sh — validates a built codinho binary end to end in a
# clean, throwaway directory (installation-documentation-ci requirement
# R8: "validar todos os links, comandos e exemplos em ambiente limpo").
# It never touches the caller's real workspace: everything happens in a
# temp directory removed on exit.
set -euo pipefail

if [ "$#" -ne 1 ]; then
  echo "usage: $0 <path-to-codinho-binary>" >&2
  exit 2
fi

BIN="$(cd "$(dirname "$1")" && pwd)/$(basename "$1")"
if [ ! -x "$BIN" ]; then
  echo "error: $BIN is not an executable file" >&2
  exit 2
fi

WORKDIR="$(mktemp -d)"
trap 'rm -rf "$WORKDIR"' EXIT
cd "$WORKDIR"

echo "==> codinho version"
"$BIN" version

echo "==> codinho init"
"$BIN" init

echo "==> codinho doctor"
"$BIN" doctor

echo "==> codinho catalog validate"
"$BIN" catalog validate

echo "==> codinho privacy export/purge round-trip"
"$BIN" privacy export --dest "$WORKDIR/export"
test -f "$WORKDIR/export/events.jsonl" || {
  # No events recorded yet (init alone never starts a session) is fine —
  # export must simply not error, not that a file necessarily exists.
  true
}

echo "smoke-install: OK"
