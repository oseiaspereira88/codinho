# Step nodes de snapshot de leituras — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

A variante `snapshot-readings` ganhou dois nodes de verificação: isolamento
bidirecional entre entrada e saída e retorno vazio não nil para entradas sem
elementos. A fixture e a variante foram preservadas.

`go run ./cmd/codinho catalog validate --checks --json` passou com
`diagnostics: null`.
