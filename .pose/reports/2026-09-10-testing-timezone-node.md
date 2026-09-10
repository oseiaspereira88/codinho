# Node de equivalência de fusos — 2026-09-10 UTC

Spec: `go-foundations-packs`, in-progress.

O desafio `go-testing.select-active-tokens-with-injected-clock` recebeu um node que verifica a validade de instantes equivalentes representados em fusos distintos. A regra temporal continua baseada no instante absoluto e no relógio injetado.

Validação: `go run ./cmd/codinho catalog validate --checks --json` passou com `diagnostics: null`; `git diff --check` passou.
