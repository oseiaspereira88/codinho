# Node de entrada vazia em defer — 2026-09-10 UTC

Spec: `go-foundations-packs`, in-progress.

O desafio `go-core.close-resources-in-defer-order` recebeu uma verificação explícita para entrada nil ou vazia, garantindo que nenhum defer seja registrado e que o retorno preserve nil conforme o contrato.

Validação: `go run ./cmd/codinho catalog validate --checks --json` passou com `diagnostics: null`; `git diff --check` passou.
