#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/../.."
source scripts/ci/tools.env
command -v pose >/dev/null
command -v govulncheck >/dev/null
pose_version=$(pose version)
pose_version=${pose_version%%$'\n'*}
case "$pose_version" in "pose $POSE_VERSION"|"$POSE_VERSION") ;; *) echo "Required POSE $POSE_VERSION, found $pose_version" >&2; exit 1 ;; esac
# Refuse a different executable even when the scanner name is present.
go version -m "$(command -v govulncheck)" | awk -v version="$GOVULNCHECK_VERSION" -v sum="$GOVULNCHECK_SUM" \
  '$1 == "mod" && $2 == "golang.org/x/vuln" && $3 == version && $4 == sum { found=1 } END { exit !found }'
pose skills-check --strict
pose docs-check
mkdir -p .pose/results
pose check --strict 2>&1 | tee .pose/results/pose-check.log
pose validate --strict --json .pose/results/delivery-validation.json 2>&1 | tee .pose/results/pose-validate.latest.log
go run ./cmd/ci-assurance > .pose/results/native-ci.json
