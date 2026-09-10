# Nodes de CSV citado — 2026-09-10 UTC

Spec: `go-foundations-packs`, in-progress.

A variante `go-io.read-quoted-attendees` recebeu verificações para preservar vírgulas dentro de campos entre aspas e para garantir que registros malformados não produzam lista parcial. O lote reforça o conceito `io-csv-quoted-fields` e mantém o contrato existente.

Validação: `go run ./cmd/codinho catalog validate --checks --json` passou com `diagnostics: null`; `git diff --check` passou.
