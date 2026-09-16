# POSE Report - 2026-09-16

## Report Type
- standard

## Task
- validate-native
- Task slug: validate-native

## Outcome
- Outcome: fail (source: derived)

## Rules Applied
- _Not provided_

## Files Changed
- pose/indexes/delivery-integrity.json
- .pose/reports/2026-09-16-standard-validate-native.md
- .pose/reports/history/standard-validate-native.jsonl
- .pose/specs/2026-08-22-go-backend-packs.md
- cmd/codinho/foundation_tracks_integration_test.go
- internal/curriculum/foundation_variants_test.go

## Validation Commands
- bash scripts/ci/check-format.sh
- go test -race -count=1 ./...
- go vet ./...
- python3 -m unittest discover -s scripts -p test_editorial.py
- go build ./cmd/codinho
- go test ./cmd/codinho -run TestPublicationIntegrity -count=1
- go run ./cmd/codinho catalog validate --published-checks --json
- go test ./cmd/codinho -run TestSessionRecoveryOverRealStdio -count=1
- go test ./cmd/codinho -run TestTreeProgressionOverRealStdio -count=1
- go test ./cmd/codinho -run TestEvaluationEvidenceOverRealStdio -count=1
- go test ./internal/mcpserver/... -run TestContract -count=1
- go test ./cmd/codinho/... -run TestTutorSkillFullSessionRoutingOverRealStdio -count=1
- go test ./cmd/codinho/... -run TestCLIEntrypointReachesEveryCommand -count=1
- go test ./cmd/codinho/... -run TestCLIEndToEndWorkflow -count=1
- govulncheck ./...
- go test ./cmd/codinho/... -run TestServeNeverWritesToStdoutOnStartupFailure|TestServeStdoutIsOnlyValidJSONFrames -count=1
- go test ./cmd/codinho/... -run TestServeStartupRespondsWithinOneSecond -count=1
- go test ./cmd/codinho/... -run TestSmokeInstallScript -count=1
- go test ./internal/ciassurance/... -count=1
- go run ./cmd/codinho catalog validate --checks --json
- go test ./cmd/codinho/... -run TestConceptContentOverRealStdio -count=1
- go test ./cmd/codinho/... -run TestTrackCompositionOverRealStdio -count=1
- go test ./cmd/codinho/... -run TestDraftOverRealStdio -count=1
- go test ./cmd/codinho/... -run TestFoundationTracksOverRealStdio -count=1

## Results
- - [pass] ./go/format (0.3s)
- - [pass] ./go/test (48.1s)
- - [pass] ./go/vet (0.1s)
- - [pass] ./go/editorial-batch-engine (2.0s)
- - [pass] ./go/build (0.1s)
- - [pass] ./go/catalog-publication-integrity (1.2s)
- - [pass] ./go/published-catalog-proof (0.2s)
- - [pass] ./go/session-recovery (0.8s)
- - [pass] ./go/session-tree-progression (0.8s)
- - [pass] ./go/evaluation-evidence (0.8s)
- - [pass] ./go/mcp-contract (1.2s)
- - [pass] ./go/tutor-skill-routing (0.7s)
- - [pass] ./go/cli-reachability (0.6s)
- - [pass] ./go/cli-e2e (0.6s)
- - [fail] ./go/govulncheck (1.0s) exit=1
- - [pass] ./go/stdout-purity (0.9s)
- - [pass] ./go/startup-budget (0.6s)
- - [pass] ./go/smoke-install (0.6s)
- - [pass] ./go/ci-assurance (0.1s)
- - [pass] ./go/catalog (16.1s)
- - [pass] ./go/concept-content-capability (0.6s)
- - [pass] ./go/learning-track-composition (0.6s)
- - [pass] ./go/agent-authored-catalog-drafts (0.6s)
- - [pass] ./go/foundation-track-routing (0.8s)
- Result: FAILURE (required check failed)

## Execution Metadata
- Generated at (UTC): 2026-09-16T05:30:19Z
- Context: auto-validate
- Validation profile: strict
- Sequence for task/spec: 27
- Stable comparison hash: 5b47855e60f64e73728abd99582eb01357a94f0c289ad7fa9125d680a322e54f

## Historical Comparison
- Previous execution: 2026-09-16T05:24:18Z
- Status: changed
- Stable field diffs:
- validation_profile: "tolerant" -> "strict"

## Risks
- _No risks provided_

## Follow-ups
- _Add next steps if needed._

## Human Review Needed
- [ ] Review functional impact
- [ ] Review validation coverage
- [ ] Approve merge
