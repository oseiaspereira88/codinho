# Nodes de chamadas do relógio — 2026-09-10 UTC

Spec: `go-foundations-packs`, in-progress.

A variante `go-testing.count-deadline-clock-calls` recebeu verificações para classificar relógio ausente sem panic e para confirmar que um lote inteiro usa uma única observação do instante. O contrato de contagem e a fixture permanecem inalterados.

Validação: `go run ./cmd/codinho catalog validate --checks --json` passou com `diagnostics: null`; `git diff --check` passou.
