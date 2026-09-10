# Step nodes de errors.As — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

O desafio `extract-field-with-errors-as` ganhou dois nodes de verificação:
wrapping em múltiplas camadas e ausência de `ValidationError`, incluindo nil.
O fluxo explicita a busca profunda e o retorno seguro sem alterar a fixture.

`go run ./cmd/codinho catalog validate --checks --json` passou com
`diagnostics: null`.
