# POSE Report - 2026-09-16

## Report Type
- standard

## Task
- validate-native
- Task slug: validate-native
- Spec: editorial-batch-engine
- Workflow: validate-native-evidence

## Outcome
- Outcome: partial (source: manual)

## Rules Applied
- documentation-style
- security
- knowledge-governance

## Files Changed
- pose/results/editorial-native-validation-current.json

## Validation Commands
- format bash scripts/ci/check-format.sh
- test go test -race -count=1 ./...
- vet go vet ./...
- editorial-batch-engine python3 -m unittest discover -s scripts -p test_editorial.py
- build go build ./cmd/codinho
- catalog-publication-integrity go test ./cmd/codinho -run TestPublicationIntegrity -count=1
- published-catalog-proof go run ./cmd/codinho catalog validate --published-checks --json
- session-recovery go test ./cmd/codinho -run TestSessionRecoveryOverRealStdio -count=1
- session-tree-progression go test ./cmd/codinho -run TestTreeProgressionOverRealStdio -count=1
- evaluation-evidence go test ./cmd/codinho -run TestEvaluationEvidenceOverRealStdio -count=1
- mcp-contract go test ./internal/mcpserver/... -run TestContract -count=1
- tutor-skill-routing go test ./cmd/codinho/... -run TestTutorSkillFullSessionRoutingOverRealStdio -count=1
- cli-reachability go test ./cmd/codinho/... -run TestCLIEntrypointReachesEveryCommand -count=1
- cli-e2e go test ./cmd/codinho/... -run TestCLIEndToEndWorkflow -count=1
- govulncheck govulncheck ./...
- stdout-purity go test ./cmd/codinho/... -run TestServeNeverWritesToStdoutOnStartupFailure|TestServeStdoutIsOnlyValidJSONFrames -count=1
- startup-budget go test ./cmd/codinho/... -run TestServeStartupRespondsWithinOneSecond -count=1
- smoke-install go test ./cmd/codinho/... -run TestSmokeInstallScript -count=1
- ci-assurance go test ./internal/ciassurance/... -count=1
- catalog go run ./cmd/codinho catalog validate --checks --json
- concept-content-capability go test ./cmd/codinho/... -run TestConceptContentOverRealStdio -count=1
- learning-track-composition go test ./cmd/codinho/... -run TestTrackCompositionOverRealStdio -count=1
- agent-authored-catalog-drafts go test ./cmd/codinho/... -run TestDraftOverRealStdio -count=1
- foundation-track-routing go test ./cmd/codinho/... -run TestFoundationTracksOverRealStdio -count=1

## Results
- Result: FAILURE — required check failed

## Execution Metadata
- Generated at (UTC): 2026-09-16T02:32:11Z
- Context: editorial-batch-engine-review
- Validation profile: strict
- Sequence for task/spec: 23
- Stable comparison hash: ddf6112356970212d66bac6b2f176f9b0bbe1bdc81e1ce96ea0fc83f4ee6b084

## Historical Comparison
- Previous execution: 2026-09-15T20:01:51Z
- Status: changed
- Stable field diffs:
- context: "auto-validate" -> "editorial-batch-engine-review"
- rules: "" -> "documentation-style,security,knowledge-governance"
- spec: "" -> "editorial-batch-engine"
- validation_profile: "tolerant" -> "strict"
- workflow: "" -> "validate-native-evidence"

## Risks
- Dependências Go e govulncheck exigem rede/DNS ou cache pré-aquecido; resultado não é aprovação da matriz.

## Follow-ups
- _Add next steps if needed._

## Human Review Needed
- [ ] Review functional impact
- [ ] Review validation coverage
- [ ] Approve merge
