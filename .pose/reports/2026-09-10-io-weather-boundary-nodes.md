# Nodes de fronteira JSON em previsão — 2026-09-10 UTC

Spec: `go-foundations-packs`, in-progress.

A variante `go-io.decode-weather-fields` recebeu verificações para rejeitar `null` como previsão e para bloquear um segundo valor JSON após o documento principal. Os nodes reforçam os conceitos de campos estritos e fronteira EOF sem alterar a fixture.

Validação: `go run ./cmd/codinho catalog validate --checks --json` passou com `diagnostics: null`; `git diff --check` passou.
