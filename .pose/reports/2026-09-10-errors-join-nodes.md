# Step nodes de errors.Join — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

O desafio `validate-user-joining-all-errors` ganhou dois nodes de verificação:
identidade individual das causas com `errors.Is` e retorno nil no caminho sem
violações. A árvore explicita preservação de informação e sucesso limpo sem
alterar fixture ou desafio canônico.

`go run ./cmd/codinho catalog validate --checks --json` passou com
`diagnostics: null`.
