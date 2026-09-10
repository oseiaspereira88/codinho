# Nodes de leitura fragmentada — 2026-09-10 UTC

Spec: `go-foundations-packs`, in-progress.

A variante `go-io.read-fragmented-note` recebeu dois nodes de verificação: um cobre o acúmulo ordenado de bytes entregues em leituras parciais e outro explicita o tratamento de EOF após conteúdo. Os nodes reutilizam os conceitos `io-partial-reads` e `io-eof-boundary`, mantendo a fixture e o contrato da variante.

Validação: `go run ./cmd/codinho catalog validate --checks --json` passou com `diagnostics: null`; `git diff --check` passou.
