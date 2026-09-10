# Node de relógio ausente em reservas — 2026-09-10 UTC

Spec: `go-foundations-packs`, in-progress.

A variante `go-testing.measure-reservation-remaining` recebeu um node dedicado à fronteira de dependência ausente, verificando `ErrMissingClock` e duração zero sem chamada inválida. Os nodes temporais existentes permanecem preservados.

Validação: `go run ./cmd/codinho catalog validate --checks --json` passou com `diagnostics: null`; `git diff --check` passou.
