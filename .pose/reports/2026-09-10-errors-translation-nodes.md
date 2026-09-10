# Step nodes de tradução de erros — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

O desafio `translate-error-at-boundary` ganhou dois nodes de verificação:
identidade pública sem a causa interna e mensagem sem o texto do timeout.
Essa separação cobre os dois canais de vazamento previstos no contrato.

`go run ./cmd/codinho catalog validate --checks --json` passou com
`diagnostics: null`.
