#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/../.."
source scripts/ci/tools.env
: "${CI_TOOL_BIN:?Set CI_TOOL_BIN to an absolute installation directory}"
case "$CI_TOOL_BIN" in /*) ;; *) echo 'CI_TOOL_BIN must be absolute' >&2; exit 1 ;; esac
case "$(uname -s)/$(uname -m)" in
  Linux/x86_64) platform=linux_amd64; digest=$POSE_LINUX_AMD64_SHA256 ;;
  Linux/aarch64|Linux/arm64) platform=linux_arm64; digest=$POSE_LINUX_ARM64_SHA256 ;;
  Darwin/x86_64) platform=darwin_amd64; digest=$POSE_DARWIN_AMD64_SHA256 ;;
  Darwin/arm64) platform=darwin_arm64; digest=$POSE_DARWIN_ARM64_SHA256 ;;
  *) echo 'Unsupported native tool platform' >&2; exit 1 ;;
esac
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
archive="pose_${POSE_VERSION}_${platform}.tar.gz"
curl --fail --silent --show-error --location --proto '=https' --tlsv1.2 \
  "https://github.com/oseiaspereira88/pose/releases/download/v${POSE_VERSION}/${archive}" -o "$work/$archive"
printf '%s  %s\n' "$digest" "$work/$archive" | shasum -a 256 -c -
tar -xzf "$work/$archive" -C "$work" pose
mkdir -p "$CI_TOOL_BIN"
install -m 755 "$work/pose" "$CI_TOOL_BIN/pose"
# Verify the pinned module against both sum.golang.org and the reviewed sum.
metadata=$(GOSUMDB=sum.golang.org go mod download -json "golang.org/x/vuln@$GOVULNCHECK_VERSION")
printf '%s\n' "$metadata" | python3 -c 'import json,sys; m=json.load(sys.stdin); sys.exit(0 if m.get("Sum")==sys.argv[1] and not m.get("Error") else 1)' "$GOVULNCHECK_SUM"
GOSUMDB=sum.golang.org GOBIN="$CI_TOOL_BIN" go install "golang.org/x/vuln/cmd/govulncheck@$GOVULNCHECK_VERSION"
"$CI_TOOL_BIN/pose" version
"$CI_TOOL_BIN/govulncheck" -version
