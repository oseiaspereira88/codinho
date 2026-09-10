# Step nodes de clonagem e busca de Circle — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

O desafio combinado `first-circle-clone` ganhou dois nodes: verificação do
primeiro Circle em qualquer posição e verificação do isolamento do Registry
clonado. A mudança cobre ordem de busca e aliasing sem criar desafio canônico.

`go run ./cmd/codinho catalog validate --checks --json` passou com
`diagnostics: null`.
