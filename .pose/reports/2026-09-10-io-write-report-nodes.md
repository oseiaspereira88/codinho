# Nodes de escrita de relatório — 2026-09-10 UTC

Spec: `go-foundations-packs`, in-progress.

A variante `go-io.write-result-report` recebeu verificações para o newline final do JSON e para a preservação da identidade do erro retornado pelo `Writer`, inclusive quando parte dos bytes já foi aceita.

Validação: `go run ./cmd/codinho catalog validate --checks --json` passou com `diagnostics: null`; `git diff --check` passou.
