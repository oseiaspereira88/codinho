# Node de casos mistos de validade — 2026-09-10 UTC

Spec: `go-foundations-packs`, in-progress.

A variante `go-testing.check-ticket-boundary` recebeu uma verificação de conjunto misto, cobrindo simultaneamente expiração anterior, igualdade e expiração futura. O node explicita que somente o caso estritamente futuro permanece válido.

Validação: `go run ./cmd/codinho catalog validate --checks --json` passou com `diagnostics: null`; `git diff --check` passou.
