# Step nodes de tempo determinístico — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

A variante `measure-reservation-remaining` ganhou dois nodes de verificação:
resultado não negativo nas fronteiras temporais e uso do relógio exatamente uma
vez. O lote reforça testabilidade determinística sem alterar desafios canônicos
ou fixtures.

`go run ./cmd/codinho catalog validate --checks --json` passou com
`diagnostics: null`.
