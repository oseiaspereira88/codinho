# Step nodes de receptor nil-safe — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

O desafio `sum-list-with-nil-safe-receiver` ganhou dois nodes: verificação de
cadeias vazia, unitária e longa, e explicação da guarda antes de acessar os
campos do receptor. Fixture e desafio canônico foram preservados.

`go run ./cmd/codinho catalog validate --checks --json` passou com
`diagnostics: null`.
